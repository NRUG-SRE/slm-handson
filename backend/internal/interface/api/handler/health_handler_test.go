package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestHealthHandler_HealthCheck(t *testing.T) {
	// Ginエンジンのセットアップ
	gin.SetMode(gin.TestMode)
	router := gin.New()
	
	// ヘルスチェックハンドラーを登録
	healthHandler := NewHealthHandler()
	router.GET("/health", healthHandler.HealthCheck)

	// テストリクエスト作成
	req, _ := http.NewRequest("GET", "/health", nil)
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