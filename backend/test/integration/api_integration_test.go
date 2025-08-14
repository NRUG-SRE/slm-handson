package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/NRUG-SRE/slm-handson/backend/internal/domain/entity"
	"github.com/NRUG-SRE/slm-handson/backend/internal/domain/repository"
	"github.com/NRUG-SRE/slm-handson/backend/internal/infrastructure/monitoring"
	"github.com/NRUG-SRE/slm-handson/backend/internal/infrastructure/persistence/memory"
	"github.com/NRUG-SRE/slm-handson/backend/internal/interface/api/handler"
	"github.com/NRUG-SRE/slm-handson/backend/internal/usecase"
)

// TestSetup はE2Eテスト用のセットアップを行う
func setupTestApplication() *gin.Engine {
	// Ginをテストモードに設定
	gin.SetMode(gin.TestMode)

	// リポジトリの初期化（インメモリ実装）
	productRepo := memory.NewProductRepository()
	cartRepo := memory.NewCartRepository()
	orderRepo := memory.NewOrderRepository()

	// テスト用の商品データを初期化（シードデータがすでにあるのでスキップ）
	// setupTestProducts(productRepo)

	// ユースケースの初期化
	productUseCase := usecase.NewProductUseCase(productRepo)
	cartUseCase := usecase.NewCartUseCase(cartRepo, productRepo)
	orderUseCase := usecase.NewOrderUseCase(orderRepo, cartRepo, productRepo)

	// New Relicクライアント（テスト用 - 環境変数なしで初期化）
	nrClient, _ := monitoring.NewNewRelicClient()

	// ハンドラーの初期化
	healthHandler := handler.NewHealthHandler()
	productHandler := handler.NewProductHandler(productUseCase, nrClient)
	cartHandler := handler.NewCartHandler(cartUseCase, nrClient)
	orderHandler := handler.NewOrderHandler(orderUseCase, nrClient)

	// テスト用のシンプルなルーター設定
	return setupTestRouter(healthHandler, productHandler, cartHandler, orderHandler)
}

// setupTestRouter はテスト用に軽量なルーターをセットアップする
func setupTestRouter(
	healthHandler *handler.HealthHandler,
	productHandler *handler.ProductHandler,
	cartHandler *handler.CartHandler,
	orderHandler *handler.OrderHandler,
) *gin.Engine {
	engine := gin.New()
	
	// 最小限のミドルウェア
	engine.Use(gin.Recovery())
	
	// ヘルスチェックエンドポイント
	engine.GET("/health", healthHandler.HealthCheck)

	// APIルートグループ
	apiV1 := engine.Group("/api")
	{
		// 商品関連エンドポイント
		apiV1.GET("/products", productHandler.GetProducts)
		apiV1.GET("/products/:id", productHandler.GetProduct)

		// カート関連エンドポイント
		apiV1.GET("/cart", cartHandler.GetCart)
		apiV1.POST("/cart/items", cartHandler.AddToCart)
		apiV1.PUT("/cart/items/:id", cartHandler.UpdateCartItem)

		// 注文関連エンドポイント
		apiV1.POST("/orders", orderHandler.CreateOrder)
		apiV1.GET("/orders/:id", orderHandler.GetOrder)
		apiV1.GET("/orders", orderHandler.GetOrders)

		// SLMデモ用エンドポイント
		apiV1.GET("/v1/error", productHandler.TriggerError)
	}

	return engine
}

// setupTestProducts はテスト用の商品データを設定する
func setupTestProducts(repo repository.ProductRepository) {
	ctx := context.Background()
	
	products := []*entity.Product{
		{
			ID:          "product-1",
			Name:        "Test Product 1",
			Description: "This is test product 1",
			Price:       1000,
			Stock:       100,
			ImageURL:    "/images/product1.svg",
		},
		{
			ID:          "product-2",
			Name:        "Test Product 2",
			Description: "This is test product 2",
			Price:       2000,
			Stock:       50,
			ImageURL:    "/images/product2.svg",
		},
	}

	for _, product := range products {
		repo.Create(ctx, product)
	}
}

// TestE2E_CompleteUserJourney は完全なユーザージャーニーをテストする
func TestE2E_CompleteUserJourney(t *testing.T) {
	// アプリケーションをセットアップ
	app := setupTestApplication()

	// 1. ヘルスチェック
	t.Run("HealthCheck", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/health", nil)
		w := httptest.NewRecorder()
		app.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("HealthCheck失敗: ステータスコード = %v, want %v", w.Code, http.StatusOK)
		}
	})

	// 2. 商品一覧の取得
	t.Run("GetProducts", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/api/products", nil)
		w := httptest.NewRecorder()
		app.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("商品一覧取得失敗: ステータスコード = %v, want %v", w.Code, http.StatusOK)
		}

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		
		if !response["success"].(bool) {
			t.Error("商品一覧取得: レスポンスのsuccessがfalse")
		}
		
		data := response["data"].([]interface{})
		if len(data) < 1 {
			t.Errorf("商品一覧取得: 商品数 = %v, want > 0", len(data))
		}
	})

	// 3. 特定商品の詳細取得（商品IDを動的に取得）
	var firstProductID, secondProductID string
	t.Run("GetProduct", func(t *testing.T) {
		// まず商品一覧から商品のIDを取得
		req, _ := http.NewRequest("GET", "/api/products", nil)
		w := httptest.NewRecorder()
		app.ServeHTTP(w, req)

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		data := response["data"].([]interface{})
		firstProduct := data[0].(map[string]interface{})
		firstProductID = firstProduct["id"].(string)
		
		if len(data) > 1 {
			secondProduct := data[1].(map[string]interface{})
			secondProductID = secondProduct["id"].(string)
		}

		// 取得したIDで商品詳細を取得
		req, _ = http.NewRequest("GET", "/api/products/"+firstProductID, nil)
		w = httptest.NewRecorder()
		app.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("商品詳細取得失敗: ステータスコード = %v, want %v", w.Code, http.StatusOK)
		}

		var detailResponse map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &detailResponse)
		
		if !detailResponse["success"].(bool) {
			t.Error("商品詳細取得: レスポンスのsuccessがfalse")
		}
		
		detailData := detailResponse["data"].(map[string]interface{})
		if detailData["id"] != firstProductID {
			t.Errorf("商品詳細取得: ID = %v, want %v", detailData["id"], firstProductID)
		}
	})

	// 4. カートに商品を追加
	t.Run("AddToCart", func(t *testing.T) {
		addRequest := map[string]interface{}{
			"productId": firstProductID,
			"quantity":  2,
		}
		jsonBody, _ := json.Marshal(addRequest)

		req, _ := http.NewRequest("POST", "/api/cart/items", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		app.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("カート追加失敗: ステータスコード = %v, want %v", w.Code, http.StatusOK)
		}

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		
		if !response["success"].(bool) {
			t.Error("カート追加: レスポンスのsuccessがfalse")
		}
	})

	// 5. カート内容を確認
	t.Run("GetCart", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/api/cart", nil)
		w := httptest.NewRecorder()
		app.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("カート取得失敗: ステータスコード = %v, want %v", w.Code, http.StatusOK)
		}

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		
		if !response["success"].(bool) {
			t.Error("カート取得: レスポンスのsuccessがfalse")
		}

		data := response["data"].(map[string]interface{})
		items := data["items"].([]interface{})
		if len(items) != 1 {
			t.Errorf("カート取得: アイテム数 = %v, want 1", len(items))
		}

		// アイテムの詳細確認
		item := items[0].(map[string]interface{})
		if item["productId"] != firstProductID {
			t.Errorf("カート取得: ProductID = %v, want %v", item["productId"], firstProductID)
		}
		if item["quantity"] != float64(2) {
			t.Errorf("カート取得: Quantity = %v, want 2", item["quantity"])
		}
	})

	// 6. 更に商品を追加
	t.Run("AddAnotherProduct", func(t *testing.T) {
		addRequest := map[string]interface{}{
			"productId": secondProductID,
			"quantity":  1,
		}
		jsonBody, _ := json.Marshal(addRequest)

		req, _ := http.NewRequest("POST", "/api/cart/items", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		app.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("2つ目の商品追加失敗: ステータスコード = %v, want %v", w.Code, http.StatusOK)
		}
	})

	// 7. 注文を作成（購入）
	var orderID string
	t.Run("CreateOrder", func(t *testing.T) {
		req, _ := http.NewRequest("POST", "/api/orders", nil)
		w := httptest.NewRecorder()
		app.ServeHTTP(w, req)

		if w.Code != http.StatusCreated {
			t.Errorf("注文作成失敗: ステータスコード = %v, want %v", w.Code, http.StatusCreated)
		}

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		
		if !response["success"].(bool) {
			t.Error("注文作成: レスポンスのsuccessがfalse")
		}

		data := response["data"].(map[string]interface{})
		orderID = data["id"].(string)
		if orderID == "" {
			t.Error("注文作成: Order IDが空")
		}

		// 注文内容の確認
		items := data["items"].([]interface{})
		if len(items) != 2 {
			t.Errorf("注文作成: アイテム数 = %v, want 2", len(items))
		}

		// 合計金額の確認（product-1: 1000円×2個 + product-2: 2000円×1個 = 4000円）
		totalAmount := data["totalAmount"].(float64)
		if totalAmount != 4000 {
			t.Errorf("注文作成: 合計金額 = %v, want 4000", totalAmount)
		}
	})

	// 8. 作成した注文の詳細を取得
	t.Run("GetOrder", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/api/orders/"+orderID, nil)
		w := httptest.NewRecorder()
		app.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("注文詳細取得失敗: ステータスコード = %v, want %v", w.Code, http.StatusOK)
		}

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		
		if !response["success"].(bool) {
			t.Error("注文詳細取得: レスポンスのsuccessがfalse")
		}

		data := response["data"].(map[string]interface{})
		if data["id"] != orderID {
			t.Errorf("注文詳細取得: ID = %v, want %v", data["id"], orderID)
		}
		
		// ステータスがpendingであることを確認
		if data["status"] != "pending" {
			t.Errorf("注文詳細取得: Status = %v, want pending", data["status"])
		}
	})

	// 9. 全注文一覧を取得
	t.Run("GetAllOrders", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/api/orders", nil)
		w := httptest.NewRecorder()
		app.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("注文一覧取得失敗: ステータスコード = %v, want %v", w.Code, http.StatusOK)
		}

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		
		if !response["success"].(bool) {
			t.Error("注文一覧取得: レスポンスのsuccessがfalse")
		}

		data := response["data"].([]interface{})
		if len(data) != 1 {
			t.Errorf("注文一覧取得: 注文数 = %v, want 1", len(data))
		}
	})

	// 10. カートが空になっていることを確認
	t.Run("VerifyEmptyCart", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/api/cart", nil)
		w := httptest.NewRecorder()
		app.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("カート確認失敗: ステータスコード = %v, want %v", w.Code, http.StatusOK)
		}

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		
		data := response["data"].(map[string]interface{})
		items := data["items"].([]interface{})
		if len(items) != 0 {
			t.Errorf("カート確認: アイテム数 = %v, want 0 (カートが空でない)", len(items))
		}
	})
}

// TestE2E_ErrorScenarios はエラーシナリオをテストする
func TestE2E_ErrorScenarios(t *testing.T) {
	app := setupTestApplication()

	// 存在しない商品の取得
	t.Run("GetNonExistentProduct", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/api/products/non-existent", nil)
		w := httptest.NewRecorder()
		app.ServeHTTP(w, req)

		if w.Code != http.StatusNotFound {
			t.Errorf("存在しない商品取得: ステータスコード = %v, want %v", w.Code, http.StatusNotFound)
		}
	})

	// 存在しない商品をカートに追加
	t.Run("AddNonExistentProductToCart", func(t *testing.T) {
		addRequest := map[string]interface{}{
			"productId": "non-existent",
			"quantity":  1,
		}
		jsonBody, _ := json.Marshal(addRequest)

		req, _ := http.NewRequest("POST", "/api/cart/items", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		app.ServeHTTP(w, req)

		if w.Code != http.StatusNotFound {
			t.Errorf("存在しない商品のカート追加: ステータスコード = %v, want %v", w.Code, http.StatusNotFound)
		}
	})

	// 無効なQuantityでカートに追加
	t.Run("AddInvalidQuantityToCart", func(t *testing.T) {
		addRequest := map[string]interface{}{
			"productId": "product-1",
			"quantity":  0,
		}
		jsonBody, _ := json.Marshal(addRequest)

		req, _ := http.NewRequest("POST", "/api/cart/items", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		app.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("無効なQuantityのカート追加: ステータスコード = %v, want %v", w.Code, http.StatusBadRequest)
		}
	})

	// 空のカートで注文作成
	t.Run("CreateOrderWithEmptyCart", func(t *testing.T) {
		req, _ := http.NewRequest("POST", "/api/orders", nil)
		w := httptest.NewRecorder()
		app.ServeHTTP(w, req)

		if w.Code != http.StatusUnprocessableEntity {
			t.Errorf("空カートでの注文作成: ステータスコード = %v, want %v", w.Code, http.StatusUnprocessableEntity)
		}
	})

	// 存在しない注文の取得
	t.Run("GetNonExistentOrder", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/api/orders/non-existent", nil)
		w := httptest.NewRecorder()
		app.ServeHTTP(w, req)

		if w.Code != http.StatusNotFound {
			t.Errorf("存在しない注文取得: ステータスコード = %v, want %v", w.Code, http.StatusNotFound)
		}
	})
}

// TestE2E_ConcurrentAccess は並行アクセスのテストを行う
func TestE2E_ConcurrentAccess(t *testing.T) {
	app := setupTestApplication()

	// 複数のゴルーチンで同時にカートに商品を追加
	t.Run("ConcurrentCartAdditions", func(t *testing.T) {
		// まず商品一覧から商品IDを取得
		req, _ := http.NewRequest("GET", "/api/products", nil)
		w := httptest.NewRecorder()
		app.ServeHTTP(w, req)

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		data := response["data"].([]interface{})
		firstProduct := data[0].(map[string]interface{})
		productID := firstProduct["id"].(string)

		numGoroutines := 10
		done := make(chan bool, numGoroutines)

		for i := 0; i < numGoroutines; i++ {
			go func() {
				defer func() { done <- true }()

				addRequest := map[string]interface{}{
					"productId": productID,
					"quantity":  1,
				}
				jsonBody, _ := json.Marshal(addRequest)

				req, _ := http.NewRequest("POST", "/api/cart/items", bytes.NewBuffer(jsonBody))
				req.Header.Set("Content-Type", "application/json")
				w := httptest.NewRecorder()
				app.ServeHTTP(w, req)

				if w.Code != http.StatusOK {
					t.Errorf("並行カート追加失敗: ステータスコード = %v", w.Code)
				}
			}()
		}

		// 全てのゴルーチンの完了を待つ
		for i := 0; i < numGoroutines; i++ {
			<-done
		}

		// カートの最終状態を確認
		cartReq, _ := http.NewRequest("GET", "/api/cart", nil)
		cartW := httptest.NewRecorder()
		app.ServeHTTP(cartW, cartReq)

		var cartResponse map[string]interface{}
		json.Unmarshal(cartW.Body.Bytes(), &cartResponse)
		
		cartData := cartResponse["data"].(map[string]interface{})
		cartItems := cartData["items"].([]interface{})
		
		if len(cartItems) != 1 {
			t.Errorf("並行処理後のカート: アイテム種類数 = %v, want 1", len(cartItems))
		}

		// 数量が正しく累積されていることを確認
		cartItem := cartItems[0].(map[string]interface{})
		quantity := cartItem["quantity"].(float64)
		if quantity != float64(numGoroutines) {
			t.Errorf("並行処理後のカート: 数量 = %v, want %v", quantity, numGoroutines)
		}
	})
}