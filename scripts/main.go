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
	targetURL      string
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
		targetURL:  targetURL,
		apiBaseURL: apiBaseURL,
		interval:   time.Duration(intervalSec) * time.Second,
		duration:   time.Duration(durationSec) * time.Second,
		userAgent:  "SLM-Handson-Access-Generator/1.0",
		httpClient: &http.Client{
			Timeout: 15 * time.Second,
		},
		rand: rand.New(rand.NewSource(time.Now().UnixNano())),
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

// makeRequest はHTTPリクエストを送信
func (ag *AccessGenerator) makeRequest(method, url, description string, body []byte) (int, time.Duration, []byte) {
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
	if method == "GET" {
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

	// レスポンスボディを読み込み
	respBody := make([]byte, 0)
	if resp.Body != nil {
		buf := make([]byte, 1024)
		for {
			n, err := resp.Body.Read(buf)
			if n > 0 {
				respBody = append(respBody, buf[:n]...)
			}
			if err != nil {
				break
			}
		}
	}

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		log.Printf("✅ %s | %d | %v", description, resp.StatusCode, responseTime.Round(time.Millisecond))
	} else {
		log.Printf("⚠️  %s | %d | %v", description, resp.StatusCode, responseTime.Round(time.Millisecond))
	}

	return resp.StatusCode, responseTime, respBody
}

// fetchProducts は商品一覧を取得
func (ag *AccessGenerator) fetchProducts() ([]Product, bool) {
	statusCode, _, body := ag.makeRequest("GET", ag.apiBaseURL+"/products", "商品一覧取得 (GET /api/products)", nil)
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

// fetchProductDetail は商品詳細を取得
func (ag *AccessGenerator) fetchProductDetail(productID string) bool {
	url := fmt.Sprintf("%s/products/%s", ag.apiBaseURL, productID)
	statusCode, _, _ := ag.makeRequest("GET", url, fmt.Sprintf("商品詳細取得 (GET /api/products/%s)", productID), nil)
	return statusCode == 200
}

// addToCart は商品をカートに追加
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

	statusCode, _, _ := ag.makeRequest("POST", ag.apiBaseURL+"/cart/items", fmt.Sprintf("カート追加 (POST /api/cart/items) - 商品ID:%s", productID), body)
	return statusCode >= 200 && statusCode < 300
}

// fetchCart はカート内容を取得
func (ag *AccessGenerator) fetchCart() bool {
	statusCode, _, _ := ag.makeRequest("GET", ag.apiBaseURL+"/cart", "カート内容取得 (GET /api/cart)", nil)
	return statusCode == 200
}

// createOrder は注文を作成
func (ag *AccessGenerator) createOrder() bool {
	// 空のJSONボディで注文作成（カート内容から自動作成）
	body := []byte("{}")
	statusCode, _, _ := ag.makeRequest("POST", ag.apiBaseURL+"/orders", "注文作成 (POST /api/orders)", body)
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

	// 1. TOPページ訪問 → 商品一覧取得
	log.Printf("📱 1. TOPページ訪問")
	products, success := ag.fetchProducts()
	if !success || len(products) == 0 {
		log.Printf("❌ ジャーニー失敗: 商品一覧取得エラー")
		return false
	}

	ag.simulateUserThinking("商品閲覧")

	// 2. ランダムな商品の詳細ページを表示
	selectedProduct := products[ag.rand.Intn(len(products))]
	log.Printf("👀 2. 商品詳細ページ表示 (商品ID: %s)", selectedProduct.ID)
	if !ag.fetchProductDetail(selectedProduct.ID) {
		log.Printf("❌ ジャーニー失敗: 商品詳細取得エラー")
		return false
	}

	ag.simulateUserThinking("商品検討")

	// 3. カートに追加
	log.Printf("🛍️  3. カートに商品追加")
	if !ag.addToCart(selectedProduct.ID) {
		log.Printf("❌ ジャーニー失敗: カート追加エラー")
		return false
	}

	ag.simulateUserThinking("カート確認")

	// 4. カートページ表示
	log.Printf("🛒 4. カートページ表示")
	if !ag.fetchCart() {
		log.Printf("❌ ジャーニー失敗: カート内容取得エラー")
		return false
	}

	ag.simulateUserThinking("決済検討")

	// 5. 決済ページ → カート内容再確認
	log.Printf("💳 5. 決済ページ表示")
	if !ag.fetchCart() {
		log.Printf("❌ ジャーニー失敗: 決済時カート確認エラー")
		return false
	}

	ag.simulateUserThinking("注文確認")

	// 6. 注文確定
	log.Printf("✅ 6. 注文確定")
	if !ag.createOrder() {
		log.Printf("❌ ジャーニー失敗: 注文作成エラー")
		return false
	}

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
}

// Run はメインのアクセス生成ループを実行
func (ag *AccessGenerator) Run() {
	log.Printf("🚀 SLM ハンズオン ユーザーアクセス生成開始")
	log.Printf("   フロントエンドURL: %s", ag.targetURL)
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
