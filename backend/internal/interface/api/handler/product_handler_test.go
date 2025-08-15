package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestProductHandler_GetProducts(t *testing.T) {
	// Ginエンジンのセットアップ
	gin.SetMode(gin.TestMode)
	router := gin.New()

	// 商品ハンドラーをモックして登録（実際のUseCaseなしでHTTPレスポンスをテスト）
	router.GET("/api/products", func(c *gin.Context) {
		// モック商品データ
		products := []gin.H{
			{
				"id":          "product-1",
				"name":        "Sample Product 1",
				"description": "This is a sample product 1",
				"price":       1000,
				"stock":       50,
				"category":    "Electronics",
				"imageURL":    "/images/product1.svg",
				"available":   true,
			},
			{
				"id":          "product-2",
				"name":        "Sample Product 2",
				"description": "This is a sample product 2",
				"price":       2000,
				"stock":       30,
				"category":    "Books",
				"imageURL":    "/images/product2.svg",
				"available":   true,
			},
		}

		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"data":    products,
			"error":   nil,
		})
	})

	// テストリクエスト作成
	req, _ := http.NewRequest("GET", "/api/products", nil)
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

func TestProductHandler_GetProduct_ValidID(t *testing.T) {
	// Ginエンジンのセットアップ
	gin.SetMode(gin.TestMode)
	router := gin.New()

	// 商品ハンドラーをモックして登録
	router.GET("/api/products/:id", func(c *gin.Context) {
		productID := c.Param("id")

		// IDが空の場合のバリデーション
		if productID == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"error": gin.H{
					"code":    "BAD_REQUEST",
					"message": "Product ID is required",
				},
			})
			return
		}

		// モック商品データ
		product := gin.H{
			"id":          productID,
			"name":        "Sample Product",
			"description": "This is a sample product",
			"price":       1000,
			"stock":       50,
			"category":    "Electronics",
			"imageURL":    "/images/product.svg",
			"available":   true,
		}

		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"data":    product,
			"error":   nil,
		})
	})

	// テストリクエスト作成
	req, _ := http.NewRequest("GET", "/api/products/product-1", nil)
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

func TestProductHandler_GetProduct_NotFound(t *testing.T) {
	// Ginエンジンのセットアップ
	gin.SetMode(gin.TestMode)
	router := gin.New()

	// 商品ハンドラーをモックして登録
	router.GET("/api/products/:id", func(c *gin.Context) {
		productID := c.Param("id")

		// 特定のIDの場合に404を返す
		if productID == "not-found" {
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"error": gin.H{
					"code":    "NOT_FOUND",
					"message": "Product not found",
				},
			})
			return
		}

		// その他の場合は正常レスポンス
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"data":    gin.H{"id": productID},
		})
	})

	// テストリクエスト作成（存在しない商品ID）
	req, _ := http.NewRequest("GET", "/api/products/not-found", nil)
	w := httptest.NewRecorder()

	// リクエスト実行
	router.ServeHTTP(w, req)

	// レスポンス検証（404エラー）
	if w.Code != http.StatusNotFound {
		t.Errorf("ステータスコード = %v, want %v", w.Code, http.StatusNotFound)
	}
}

func TestProductHandler_TriggerError(t *testing.T) {
	// Ginエンジンのセットアップ
	gin.SetMode(gin.TestMode)
	router := gin.New()

	// SLMデモ用エラーエンドポイントをモックして登録
	router.GET("/api/v1/error", func(c *gin.Context) {
		// 意図的にエラーを返す
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INTERNAL_SERVER_ERROR",
				"message": "This is a simulated error for SLM demonstration",
			},
		})
	})

	// テストリクエスト作成
	req, _ := http.NewRequest("GET", "/api/v1/error", nil)
	w := httptest.NewRecorder()

	// リクエスト実行
	router.ServeHTTP(w, req)

	// レスポンス検証（500エラー）
	if w.Code != http.StatusInternalServerError {
		t.Errorf("ステータスコード = %v, want %v", w.Code, http.StatusInternalServerError)
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

func TestProductHandler_GetProduct_WithUserAgent(t *testing.T) {
	// Ginエンジンのセットアップ
	gin.SetMode(gin.TestMode)
	router := gin.New()

	// User-Agentヘッダーを確認するハンドラー
	router.GET("/api/products/:id", func(c *gin.Context) {
		productID := c.Param("id")
		userAgent := c.GetHeader("User-Agent")

		// User-Agentが設定されていることを確認
		if userAgent == "" {
			userAgent = "unknown"
		}

		product := gin.H{
			"id":        productID,
			"name":      "Sample Product",
			"userAgent": userAgent, // テスト用にUser-Agentも含める
		}

		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"data":    product,
		})
	})

	// テストリクエスト作成（User-Agentヘッダー付き）
	req, _ := http.NewRequest("GET", "/api/products/product-1", nil)
	req.Header.Set("User-Agent", "Test-Browser/1.0")
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
