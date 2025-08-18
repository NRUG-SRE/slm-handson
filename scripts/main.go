package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"
)

// Product represents a product from the API
type Product struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Price int    `json:"price"`
}

// APIResponse represents the standard API response format
type APIResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
}

// CartItem represents an item in the cart
type CartItem struct {
	ProductID string `json:"productId"`
	Quantity  int    `json:"quantity"`
}

// CartResponse represents the cart response
type CartResponse struct {
	Items []CartItem `json:"items"`
	Total int        `json:"total"`
}

// AccessGenerator はSLMハンズオン用のユーザーアクセス生成器
type AccessGenerator struct {
	frontendURL    string
	apiBaseURL     string
	interval       time.Duration
	duration       time.Duration
	userAgent      string
	httpClient     *http.Client
	accessCount    int
	successCount   int
	journeyCount   int
	completeCount  int
	startTime      time.Time
	rand           *rand.Rand
	sessionID      string  // セッション管理用
}

// NewAccessGenerator は新しいアクセス生成器を作成
func NewAccessGenerator() *AccessGenerator {
	// 環境変数から設定を取得
	targetURL := getEnv("TARGET_URL", "http://localhost:3000")
	// Docker環境では api-server コンテナにアクセス
	apiBaseURL := strings.Replace(targetURL, "frontend:3000", "api-server:8080/api", 1)
	if strings.Contains(targetURL, "localhost") {
		apiBaseURL = strings.Replace(targetURL, ":3000", ":8080/api", 1)
	}
	intervalSec := getEnvInt("ACCESS_INTERVAL", 10)
	durationSec := getEnvInt("DURATION", 300)

	return &AccessGenerator{
		frontendURL: targetURL,
		apiBaseURL:  apiBaseURL,
		interval:    time.Duration(intervalSec) * time.Second,
		duration:    time.Duration(durationSec) * time.Second,
		userAgent:   "SLM-Handson-Access-Generator/1.0",
		httpClient: &http.Client{
			Timeout: 15 * time.Second,
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse // リダイレクトを自動追従しない
			},
		},
		rand:      rand.New(rand.NewSource(time.Now().UnixNano())),
		sessionID: fmt.Sprintf("session-%d", time.Now().UnixNano()),
	}
}

// getEnv は環境変数を取得（デフォルト値付き）
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvInt は環境変数を整数として取得（デフォルト値付き）
func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

// makeRequest はHTTPリクエストを送信（内部API用）
func (ag *AccessGenerator) makeRequest(method, url, description string, body []byte, isHTML bool) (int, time.Duration, []byte) {
	var req *http.Request
	var err error

	if body != nil {
		req, err = http.NewRequest(method, url, bytes.NewBuffer(body))
	} else {
		req, err = http.NewRequest(method, url, nil)
	}

	if err != nil {
		log.Printf("❌ %s | REQUEST_ERROR: %v", description, err)
		return 0, 0, nil
	}

	req.Header.Set("User-Agent", ag.userAgent)
	req.Header.Set("X-Session-ID", ag.sessionID)
	
	if isHTML {
		// HTMLページリクエスト
		req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
		req.Header.Set("Accept-Language", "ja,en;q=0.9")
	} else if method == "GET" {
		req.Header.Set("Accept", "application/json,text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	} else {
		req.Header.Set("Accept", "application/json")
		req.Header.Set("Content-Type", "application/json")
	}

	start := time.Now()
	resp, err := ag.httpClient.Do(req)
	responseTime := time.Since(start)

	if err != nil {
		log.Printf("❌ %s | ERROR: %v", description, err)
		return 0, responseTime, nil
	}
	defer resp.Body.Close()

	// レスポンスボディを読み込み（最初の1KBのみ、HTMLは破棄）
	respBody := make([]byte, 0)
	if !isHTML && resp.Body != nil {
		buf := make([]byte, 1024)
		for {
			n, err := resp.Body.Read(buf)
			if n > 0 {
				respBody = append(respBody, buf[:n]...)
			}
			if err != nil || len(respBody) > 10240 { // 10KB以上は読まない
				break
			}
		}
	}

	if resp.StatusCode >= 200 && resp.StatusCode < 400 {
		log.Printf("✅ %s | %d | %v", description, resp.StatusCode, responseTime.Round(time.Millisecond))
	} else {
		log.Printf("⚠️  %s | %d | %v", description, resp.StatusCode, responseTime.Round(time.Millisecond))
	}

	return resp.StatusCode, responseTime, respBody
}

// visitFrontendPage はフロントエンドページを訪問（RUM計測のため）
func (ag *AccessGenerator) visitFrontendPage(path, description string) bool {
	url := ag.frontendURL + path
	statusCode, _, _ := ag.makeRequest("GET", url, description, nil, true)
	// HTMLページは200-399を成功とする（リダイレクト含む）
	return statusCode >= 200 && statusCode < 400
}

// fetchProductsAPI は商品一覧をAPI経由で取得（商品IDリスト取得用）
func (ag *AccessGenerator) fetchProductsAPI() ([]Product, bool) {
	statusCode, _, body := ag.makeRequest("GET", ag.apiBaseURL+"/products", "商品一覧API (GET /api/products)", nil, false)
	if statusCode != 200 {
		return nil, false
	}

	var response APIResponse
	if err := json.Unmarshal(body, &response); err != nil {
		log.Printf("❌ APIレスポンスJSONパース失敗: %v", err)
		return nil, false
	}

	if !response.Success {
		log.Printf("❌ APIエラー: success=false")
		return nil, false
	}

	// data フィールドを []Product にマップ
	dataBytes, _ := json.Marshal(response.Data)
	var products []Product
	if err := json.Unmarshal(dataBytes, &products); err != nil {
		log.Printf("❌ 商品一覧JSONパース失敗: %v", err)
		return nil, false
	}

	return products, true
}

// addToCart は商品をカートに追加（API直接呼び出し）
func (ag *AccessGenerator) addToCart(productID string) bool {
	cartItem := CartItem{
		ProductID: productID,
		Quantity:  ag.rand.Intn(3) + 1, // 1-3個をランダム選択
	}

	body, err := json.Marshal(cartItem)
	if err != nil {
		log.Printf("❌ カート追加JSONエンコード失敗: %v", err)
		return false
	}

	statusCode, _, _ := ag.makeRequest("POST", ag.apiBaseURL+"/cart/items", fmt.Sprintf("カート追加API (POST /api/cart/items) - 商品ID:%s", productID), body, false)
	return statusCode >= 200 && statusCode < 300
}

// createOrder は注文を作成（API直接呼び出し）
func (ag *AccessGenerator) createOrder() bool {
	// 空のJSONボディで注文作成（カート内容から自動作成）
	body := []byte("{}")
	statusCode, _, _ := ag.makeRequest("POST", ag.apiBaseURL+"/orders", "注文作成API (POST /api/orders)", body, false)
	return statusCode >= 200 && statusCode < 300
}

// simulateUserThinking はユーザーの思考時間をシミュレート
func (ag *AccessGenerator) simulateUserThinking(action string) {
	// 1-5秒のランダムな待機時間でユーザーの行動をシミュレート
	waitTime := time.Duration(ag.rand.Intn(4)+1) * time.Second
	log.Printf("💭 %s中... (%v)", action, waitTime.Round(time.Second))
	time.Sleep(waitTime)
}

// performUserJourney は完全なユーザージャーニーを実行
func (ag *AccessGenerator) performUserJourney() bool {
	log.Printf("🛒 ユーザージャーニー開始 (#%d)", ag.journeyCount+1)

	// 1. TOPページ訪問（Frontend）- RUM計測
	log.Printf("📱 1. TOPページ訪問")
	if !ag.visitFrontendPage("/", "TOPページ (GET /)") {
		log.Printf("❌ ジャーニー失敗: TOPページアクセスエラー")
		return false
	}
	
	// API経由で商品リストを取得（ID取得のため）
	products, success := ag.fetchProductsAPI()
	if !success || len(products) == 0 {
		log.Printf("❌ ジャーニー失敗: 商品一覧取得エラー")
		return false
	}

	ag.simulateUserThinking("商品閲覧")

	// 2. ランダムな商品の詳細ページを表示（Frontend）- RUM計測
	selectedProduct := products[ag.rand.Intn(len(products))]
	log.Printf("👀 2. 商品詳細ページ表示 (商品ID: %s)", selectedProduct.ID)
	productPath := fmt.Sprintf("/products/%s", selectedProduct.ID)
	if !ag.visitFrontendPage(productPath, fmt.Sprintf("商品詳細ページ (GET %s)", productPath)) {
		log.Printf("❌ ジャーニー失敗: 商品詳細ページアクセスエラー")
		return false
	}

	ag.simulateUserThinking("商品検討")

	// 3. カートに追加（API直接）- カート操作は内部API
	log.Printf("🛍️  3. カートに商品追加")
	if !ag.addToCart(selectedProduct.ID) {
		log.Printf("❌ ジャーニー失敗: カート追加エラー")
		return false
	}

	ag.simulateUserThinking("カート確認")

	// 4. カートページ表示（Frontend）- RUM計測
	log.Printf("🛒 4. カートページ表示")
	if !ag.visitFrontendPage("/cart", "カートページ (GET /cart)") {
		log.Printf("❌ ジャーニー失敗: カートページアクセスエラー")
		return false
	}

	ag.simulateUserThinking("決済検討")

	// 5. 決済ページ表示（Frontend）- RUM計測
	log.Printf("💳 5. 決済ページ表示")
	if !ag.visitFrontendPage("/checkout", "決済ページ (GET /checkout)") {
		log.Printf("❌ ジャーニー失敗: 決済ページアクセスエラー")
		return false
	}

	ag.simulateUserThinking("注文確認")

	// 6. 注文確定（API直接）- 注文処理は内部API
	log.Printf("✅ 6. 注文確定")
	if !ag.createOrder() {
		log.Printf("❌ ジャーニー失敗: 注文作成エラー")
		return false
	}

	// 7. 注文完了ページ（仮想）- 実際にはTOPにリダイレクト
	log.Printf("🎊 7. 注文完了画面表示")
	ag.visitFrontendPage("/", "注文完了後TOPページリダイレクト (GET /)")

	log.Printf("🎉 ユーザージャーニー完了! 商品ID:%s → 注文完了", selectedProduct.ID)
	return true
}

// printStatistics は統計情報を表示
func (ag *AccessGenerator) printStatistics() {
	if ag.journeyCount == 0 {
		return
	}

	completionRate := float64(ag.completeCount) / float64(ag.journeyCount) * 100
	elapsed := time.Since(ag.startTime)

	log.Printf("📊 統計 | ジャーニー数: %d | 完了数: %d | 完了率: %.1f%% | 経過時間: %v",
		ag.journeyCount, ag.completeCount, completionRate, elapsed.Round(time.Second))
}

// printFinalStatistics は最終統計を表示
func (ag *AccessGenerator) printFinalStatistics() {
	elapsed := time.Since(ag.startTime)
	completionRate := float64(0)
	if ag.journeyCount > 0 {
		completionRate = float64(ag.completeCount) / float64(ag.journeyCount) * 100
	}

	fmt.Println(strings.Repeat("=", 70))
	log.Printf("📈 最終統計")
	fmt.Println(strings.Repeat("=", 70))
	log.Printf("実行時間: %v", elapsed.Round(time.Second))
	log.Printf("総ユーザージャーニー数: %d", ag.journeyCount)
	log.Printf("完了したジャーニー数: %d", ag.completeCount)
	log.Printf("ジャーニー完了率: %.1f%%", completionRate)
	log.Printf("🏁 SLM ハンズオン ユーザージャーニー完了")
	log.Printf("💡 New Relic UIでSLO/SLI監視データを確認してください")
	log.Printf("📊 APM（バックエンド）とRUM（フロントエンド）両方のデータが収集されました")
}

// Run はメインのアクセス生成ループを実行
func (ag *AccessGenerator) Run() {
	log.Printf("🚀 SLM ハンズオン ユーザーアクセス生成開始")
	log.Printf("   フロントエンドURL: %s", ag.frontendURL)
	log.Printf("   API URL: %s", ag.apiBaseURL)
	log.Printf("   ジャーニー間隔: %v", ag.interval)
	log.Printf("   実行時間: %v", ag.duration)

	// シグナルハンドリング設定
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	ag.startTime = time.Now()
	ticker := time.NewTicker(ag.interval)
	defer ticker.Stop()

	timeoutChan := time.After(ag.duration)

	fmt.Println(strings.Repeat("=", 70))
	log.Printf("🎯 SLO監視用ECサイトユーザージャーニー開始")
	log.Printf("🛒 フロー: TOPページ → 商品詳細 → カート追加 → カート確認 → 決済 → 注文完了")
	log.Printf("📱 Frontend（RUM）: ページ表示でブラウザメトリクスを収集")
	log.Printf("⚙️  Backend（APM）: API呼び出しでサーバーメトリクスを収集")
	fmt.Println(strings.Repeat("=", 70))

	for {
		select {
		case <-sigChan:
			log.Printf("⏹️  シグナルを受信しました。停止中...")
			ag.printFinalStatistics()
			return

		case <-timeoutChan:
			log.Printf("⏰ 指定時間(%v)が経過しました", ag.duration)
			ag.printFinalStatistics()
			return

		case <-ticker.C:
			// 完全なユーザージャーニーを実行
			ag.journeyCount++
			if ag.performUserJourney() {
				ag.completeCount++
			}

			// 統計情報を表示
			if ag.journeyCount%3 == 0 {
				ag.printStatistics()
			}
		}
	}
}

func main() {
	generator := NewAccessGenerator()
	generator.Run()
}