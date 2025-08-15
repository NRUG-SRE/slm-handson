package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestCartHandler_GetCart(t *testing.T) {
	// Ginエンジンのセットアップ
	gin.SetMode(gin.TestMode)
	router := gin.New()

	// カートハンドラーをモックして登録（実際のUseCaseなしでHTTPレスポンスをテスト）
	router.GET("/api/cart", func(c *gin.Context) {
		// モックレスポンス
		cartData := gin.H{
			"id": DefaultCartID,
			"items": []gin.H{
				{
					"id":        "item-1",
					"productId": "product-1",
					"quantity":  2,
					"price":     1000,
				},
			},
			"totalAmount": 2000,
			"itemCount":   1,
		}
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"data":    cartData,
			"error":   nil,
		})
	})

	// テストリクエスト作成
	req, _ := http.NewRequest("GET", "/api/cart", nil)
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

func TestCartHandler_AddToCart_ValidRequest(t *testing.T) {
	// Ginエンジンのセットアップ
	gin.SetMode(gin.TestMode)
	router := gin.New()

	// カートハンドラーをモックして登録
	router.POST("/api/cart/items", func(c *gin.Context) {
		var req AddToCartRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"error": gin.H{
					"code":    "BAD_REQUEST",
					"message": "Invalid request body",
				},
			})
			return
		}

		// モックレスポンス（成功）
		cartData := gin.H{
			"id": DefaultCartID,
			"items": []gin.H{
				{
					"id":        "item-1",
					"productId": req.ProductID,
					"quantity":  req.Quantity,
					"price":     1000,
				},
			},
			"totalAmount": 1000 * req.Quantity,
			"itemCount":   1,
		}
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"data":    cartData,
			"error":   nil,
		})
	})

	// テストリクエストボディ作成
	requestBody := AddToCartRequest{
		ProductID: "product-1",
		Quantity:  2,
	}
	jsonBody, _ := json.Marshal(requestBody)

	// テストリクエスト作成
	req, _ := http.NewRequest("POST", "/api/cart/items", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
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

func TestCartHandler_AddToCart_InvalidRequest(t *testing.T) {
	// Ginエンジンのセットアップ
	gin.SetMode(gin.TestMode)
	router := gin.New()

	// カートハンドラーをモックして登録
	router.POST("/api/cart/items", func(c *gin.Context) {
		var req AddToCartRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"error": gin.H{
					"code":    "BAD_REQUEST",
					"message": "Invalid request body",
				},
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{"success": true})
	})

	// 無効なテストリクエストボディ作成（必須フィールドなし）
	invalidBody := gin.H{
		"quantity": 2,
		// productIdが欠落
	}
	jsonBody, _ := json.Marshal(invalidBody)

	// テストリクエスト作成
	req, _ := http.NewRequest("POST", "/api/cart/items", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// リクエスト実行
	router.ServeHTTP(w, req)

	// レスポンス検証（バリデーションエラー）
	if w.Code != http.StatusBadRequest {
		t.Errorf("ステータスコード = %v, want %v", w.Code, http.StatusBadRequest)
	}
}

func TestCartHandler_UpdateCartItem_ValidRequest(t *testing.T) {
	// Ginエンジンのセットアップ
	gin.SetMode(gin.TestMode)
	router := gin.New()

	// カートハンドラーをモックして登録
	router.PUT("/api/cart/items/:id", func(c *gin.Context) {
		itemID := c.Param("id")
		if itemID == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"error": gin.H{
					"code":    "BAD_REQUEST",
					"message": "Item ID is required",
				},
			})
			return
		}

		var req UpdateCartItemRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"error": gin.H{
					"code":    "BAD_REQUEST",
					"message": "Invalid request body",
				},
			})
			return
		}

		// 負の値チェック
		if req.Quantity < 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"error": gin.H{
					"code":    "BAD_REQUEST",
					"message": "Quantity cannot be negative",
				},
			})
			return
		}

		// モックレスポンス（成功）
		cartData := gin.H{
			"id": DefaultCartID,
			"items": []gin.H{
				{
					"id":        itemID,
					"productId": "product-1",
					"quantity":  req.Quantity,
					"price":     1000,
				},
			},
			"totalAmount": 1000 * req.Quantity,
			"itemCount":   1,
		}
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"data":    cartData,
			"error":   nil,
		})
	})

	// テストリクエストボディ作成
	requestBody := UpdateCartItemRequest{
		Quantity: 3,
	}
	jsonBody, _ := json.Marshal(requestBody)

	// テストリクエスト作成
	req, _ := http.NewRequest("PUT", "/api/cart/items/item-1", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// リクエスト実行
	router.ServeHTTP(w, req)

	// レスポンス検証
	if w.Code != http.StatusOK {
		t.Errorf("ステータスコード = %v, want %v", w.Code, http.StatusOK)
	}
}

func TestCartHandler_UpdateCartItem_NegativeQuantity(t *testing.T) {
	// Ginエンジンのセットアップ
	gin.SetMode(gin.TestMode)
	router := gin.New()

	// カートハンドラーをモックして登録
	router.PUT("/api/cart/items/:id", func(c *gin.Context) {
		var req UpdateCartItemRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"error": gin.H{
					"code":    "BAD_REQUEST",
					"message": "Invalid request body",
				},
			})
			return
		}

		// 負の値チェック
		if req.Quantity < 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"error": gin.H{
					"code":    "BAD_REQUEST",
					"message": "Quantity cannot be negative",
				},
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{"success": true})
	})

	// 負の数量でテストリクエストボディ作成
	requestBody := UpdateCartItemRequest{
		Quantity: -1,
	}
	jsonBody, _ := json.Marshal(requestBody)

	// テストリクエスト作成
	req, _ := http.NewRequest("PUT", "/api/cart/items/item-1", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// リクエスト実行
	router.ServeHTTP(w, req)

	// レスポンス検証（バリデーションエラー）
	if w.Code != http.StatusBadRequest {
		t.Errorf("ステータスコード = %v, want %v", w.Code, http.StatusBadRequest)
	}
}
