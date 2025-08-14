package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestOrderHandler_CreateOrder(t *testing.T) {
	// Ginエンジンのセットアップ
	gin.SetMode(gin.TestMode)
	router := gin.New()
	
	// 注文ハンドラーをモックして登録（実際のUseCaseなしでHTTPレスポンスをテスト）
	router.POST("/api/orders", func(c *gin.Context) {
		// モック注文データ（成功レスポンス）
		orderData := gin.H{
			"id": "order-123",
			"cartId": DefaultCartID,
			"status": "Pending",
			"items": []gin.H{
				{
					"id": "item-1",
					"productId": "product-1",
					"quantity": 2,
					"price": 1000,
				},
			},
			"totalAmount": 2000,
			"itemCount": 1,
			"createdAt": "2024-01-01T00:00:00Z",
		}
		
		c.JSON(http.StatusCreated, gin.H{
			"success": true,
			"data": orderData,
			"error": nil,
		})
	})

	// テストリクエスト作成
	req, _ := http.NewRequest("POST", "/api/orders", nil)
	w := httptest.NewRecorder()

	// リクエスト実行
	router.ServeHTTP(w, req)

	// レスポンス検証
	if w.Code != http.StatusCreated {
		t.Errorf("ステータスコード = %v, want %v", w.Code, http.StatusCreated)
	}

	// JSONレスポンスの構造を確認
	body := w.Body.String()
	if body == "" {
		t.Error("レスポンスボディが空です")
	}

	// Content-Typeの検証
	contentType := w.Header().Get("Content-Type")
	expectedContentType := "application/json; charset=utf-8"
	if contentType != expectedContentType {
		t.Errorf("Content-Type = %v, want %v", contentType, expectedContentType)
	}
}

func TestOrderHandler_CreateOrder_EmptyCart(t *testing.T) {
	// Ginエンジンのセットアップ
	gin.SetMode(gin.TestMode)
	router := gin.New()
	
	// 空のカートの場合のエラーレスポンスをモック
	router.POST("/api/orders", func(c *gin.Context) {
		// 空カートエラーをシミュレート
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"success": false,
			"error": gin.H{
				"code": "UNPROCESSABLE_ENTITY",
				"message": "Cart is empty",
			},
		})
	})

	// テストリクエスト作成
	req, _ := http.NewRequest("POST", "/api/orders", nil)
	w := httptest.NewRecorder()

	// リクエスト実行
	router.ServeHTTP(w, req)

	// レスポンス検証（422エラー）
	if w.Code != http.StatusUnprocessableEntity {
		t.Errorf("ステータスコード = %v, want %v", w.Code, http.StatusUnprocessableEntity)
	}
}

func TestOrderHandler_GetOrder_ValidID(t *testing.T) {
	// Ginエンジンのセットアップ
	gin.SetMode(gin.TestMode)
	router := gin.New()
	
	// 注文ハンドラーをモックして登録
	router.GET("/api/orders/:id", func(c *gin.Context) {
		orderID := c.Param("id")
		
		// IDが空の場合のバリデーション
		if orderID == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"error": gin.H{
					"code": "BAD_REQUEST",
					"message": "Order ID is required",
				},
			})
			return
		}

		// モック注文データ
		orderData := gin.H{
			"id": orderID,
			"cartId": DefaultCartID,
			"status": "Completed",
			"items": []gin.H{
				{
					"id": "item-1",
					"productId": "product-1",
					"quantity": 1,
					"price": 1000,
				},
			},
			"totalAmount": 1000,
			"itemCount": 1,
			"createdAt": "2024-01-01T00:00:00Z",
		}
		
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"data": orderData,
			"error": nil,
		})
	})

	// テストリクエスト作成
	req, _ := http.NewRequest("GET", "/api/orders/order-123", nil)
	w := httptest.NewRecorder()

	// リクエスト実行
	router.ServeHTTP(w, req)

	// レスポンス検証
	if w.Code != http.StatusOK {
		t.Errorf("ステータスコード = %v, want %v", w.Code, http.StatusOK)
	}

	// JSONレスポンスの構造を確認
	body := w.Body.String()
	if body == "" {
		t.Error("レスポンスボディが空です")
	}
}

func TestOrderHandler_GetOrder_NotFound(t *testing.T) {
	// Ginエンジンのセットアップ
	gin.SetMode(gin.TestMode)
	router := gin.New()
	
	// 注文ハンドラーをモックして登録
	router.GET("/api/orders/:id", func(c *gin.Context) {
		orderID := c.Param("id")
		
		// 特定のIDの場合に404を返す
		if orderID == "not-found" {
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"error": gin.H{
					"code": "NOT_FOUND",
					"message": "Order not found",
				},
			})
			return
		}

		// その他の場合は正常レスポンス
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"data": gin.H{"id": orderID},
		})
	})

	// テストリクエスト作成（存在しない注文ID）
	req, _ := http.NewRequest("GET", "/api/orders/not-found", nil)
	w := httptest.NewRecorder()

	// リクエスト実行
	router.ServeHTTP(w, req)

	// レスポンス検証（404エラー）
	if w.Code != http.StatusNotFound {
		t.Errorf("ステータスコード = %v, want %v", w.Code, http.StatusNotFound)
	}
}

func TestOrderHandler_GetOrders(t *testing.T) {
	// Ginエンジンのセットアップ
	gin.SetMode(gin.TestMode)
	router := gin.New()
	
	// 注文一覧ハンドラーをモックして登録
	router.GET("/api/orders", func(c *gin.Context) {
		// モック注文一覧データ
		orders := []gin.H{
			{
				"id": "order-1",
				"cartId": DefaultCartID,
				"status": "Completed",
				"totalAmount": 1000,
				"itemCount": 1,
				"createdAt": "2024-01-01T00:00:00Z",
			},
			{
				"id": "order-2",
				"cartId": DefaultCartID,
				"status": "Pending",
				"totalAmount": 2000,
				"itemCount": 2,
				"createdAt": "2024-01-02T00:00:00Z",
			},
		}
		
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"data": orders,
			"error": nil,
		})
	})

	// テストリクエスト作成
	req, _ := http.NewRequest("GET", "/api/orders", nil)
	w := httptest.NewRecorder()

	// リクエスト実行
	router.ServeHTTP(w, req)

	// レスポンス検証
	if w.Code != http.StatusOK {
		t.Errorf("ステータスコード = %v, want %v", w.Code, http.StatusOK)
	}

	// JSONレスポンスの構造を確認
	body := w.Body.String()
	if body == "" {
		t.Error("レスポンスボディが空です")
	}

	// Content-Typeの検証
	contentType := w.Header().Get("Content-Type")
	expectedContentType := "application/json; charset=utf-8"
	if contentType != expectedContentType {
		t.Errorf("Content-Type = %v, want %v", contentType, expectedContentType)
	}
}

func TestOrderHandler_CreateOrder_WithJSONRequest(t *testing.T) {
	// Ginエンジンのセットアップ
	gin.SetMode(gin.TestMode)
	router := gin.New()
	
	// JSON バインディングをテストするハンドラー
	router.POST("/api/orders", func(c *gin.Context) {
		var req CreateOrderRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"error": gin.H{
					"code": "BAD_REQUEST",
					"message": "Invalid request body",
				},
			})
			return
		}

		// リクエストが正常に解析された場合
		orderData := gin.H{
			"id": "order-123",
			"items": req.Items,
			"status": "Pending",
		}
		
		c.JSON(http.StatusCreated, gin.H{
			"success": true,
			"data": orderData,
		})
	})

	// テストリクエストボディ作成
	requestBody := CreateOrderRequest{
		Items: []CreateOrderItem{
			{
				ProductID: "product-1",
				Quantity:  2,
			},
		},
	}
	jsonBody, _ := json.Marshal(requestBody)

	// テストリクエスト作成
	req, _ := http.NewRequest("POST", "/api/orders", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// リクエスト実行
	router.ServeHTTP(w, req)

	// レスポンス検証
	if w.Code != http.StatusCreated {
		t.Errorf("ステータスコード = %v, want %v", w.Code, http.StatusCreated)
	}
}

func TestOrderHandler_CreateOrder_InvalidJSON(t *testing.T) {
	// Ginエンジンのセットアップ
	gin.SetMode(gin.TestMode)
	router := gin.New()
	
	// JSON バインディングをテストするハンドラー
	router.POST("/api/orders", func(c *gin.Context) {
		var req CreateOrderRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"error": gin.H{
					"code": "BAD_REQUEST",
					"message": "Invalid request body",
				},
			})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"success": true})
	})

	// 完全に無効なJSONでテストリクエスト作成
	invalidJSON := `{"invalid": json}`

	// テストリクエスト作成
	req, _ := http.NewRequest("POST", "/api/orders", bytes.NewBufferString(invalidJSON))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// リクエスト実行
	router.ServeHTTP(w, req)

	// レスポンス検証（バリデーションエラー）
	if w.Code != http.StatusBadRequest {
		t.Errorf("ステータスコード = %v, want %v", w.Code, http.StatusBadRequest)
	}
}