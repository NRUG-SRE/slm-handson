package presenter

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	return gin.New()
}

func TestSuccessResponse(t *testing.T) {
	router := setupRouter()

	router.GET("/test", func(c *gin.Context) {
		data := map[string]string{
			"message": "test successful",
		}
		SuccessResponse(c, http.StatusOK, data)
	})

	req, _ := http.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// ステータスコードの検証
	if w.Code != http.StatusOK {
		t.Errorf("ステータスコード = %v, want %v", w.Code, http.StatusOK)
	}

	// レスポンスボディの検証
	var response Response
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("JSONのパースでエラー: %v", err)
	}

	if !response.Success {
		t.Error("Successがtrueになっていません")
	}

	if response.Data == nil {
		t.Error("Dataが設定されていません")
	}

	if response.Error != nil {
		t.Error("Errorが設定されてしまいました")
	}

	// Content-Typeの検証
	contentType := w.Header().Get("Content-Type")
	expectedContentType := "application/json; charset=utf-8"
	if contentType != expectedContentType {
		t.Errorf("Content-Type = %v, want %v", contentType, expectedContentType)
	}
}

func TestErrorResponse(t *testing.T) {
	router := setupRouter()

	router.GET("/error", func(c *gin.Context) {
		ErrorResponse(c, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid input provided")
	})

	req, _ := http.NewRequest("GET", "/error", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// ステータスコードの検証
	if w.Code != http.StatusBadRequest {
		t.Errorf("ステータスコード = %v, want %v", w.Code, http.StatusBadRequest)
	}

	// レスポンスボディの検証
	var response Response
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("JSONのパースでエラー: %v", err)
	}

	if response.Success {
		t.Error("Successがfalseになっていません")
	}

	if response.Data != nil {
		t.Error("Dataが設定されてしまいました")
	}

	if response.Error == nil {
		t.Error("Errorが設定されていません")
	}

	if response.Error.Code != "VALIDATION_ERROR" {
		t.Errorf("Error.Code = %v, want VALIDATION_ERROR", response.Error.Code)
	}

	if response.Error.Message != "Invalid input provided" {
		t.Errorf("Error.Message = %v, want 'Invalid input provided'", response.Error.Message)
	}
}

func TestBadRequestResponse(t *testing.T) {
	router := setupRouter()

	router.GET("/badrequest", func(c *gin.Context) {
		BadRequestResponse(c, "Missing required field")
	})

	req, _ := http.NewRequest("GET", "/badrequest", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// ステータスコードの検証
	if w.Code != http.StatusBadRequest {
		t.Errorf("ステータスコード = %v, want %v", w.Code, http.StatusBadRequest)
	}

	// レスポンスボディの検証
	var response Response
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("JSONのパースでエラー: %v", err)
	}

	if response.Success {
		t.Error("Successがfalseになっていません")
	}

	if response.Error == nil {
		t.Error("Errorが設定されていません")
	}

	if response.Error.Code != "BAD_REQUEST" {
		t.Errorf("Error.Code = %v, want BAD_REQUEST", response.Error.Code)
	}

	if response.Error.Message != "Missing required field" {
		t.Errorf("Error.Message = %v, want 'Missing required field'", response.Error.Message)
	}
}

func TestNotFoundResponse(t *testing.T) {
	router := setupRouter()

	router.GET("/notfound", func(c *gin.Context) {
		NotFoundResponse(c, "Resource not found")
	})

	req, _ := http.NewRequest("GET", "/notfound", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// ステータスコードの検証
	if w.Code != http.StatusNotFound {
		t.Errorf("ステータスコード = %v, want %v", w.Code, http.StatusNotFound)
	}

	// レスポンスボディの検証
	var response Response
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("JSONのパースでエラー: %v", err)
	}

	if response.Error == nil {
		t.Error("Errorが設定されていません")
	}

	if response.Error.Code != "NOT_FOUND" {
		t.Errorf("Error.Code = %v, want NOT_FOUND", response.Error.Code)
	}
}

func TestInternalServerErrorResponse(t *testing.T) {
	router := setupRouter()

	router.GET("/servererror", func(c *gin.Context) {
		InternalServerErrorResponse(c, "Database connection failed")
	})

	req, _ := http.NewRequest("GET", "/servererror", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// ステータスコードの検証
	if w.Code != http.StatusInternalServerError {
		t.Errorf("ステータスコード = %v, want %v", w.Code, http.StatusInternalServerError)
	}

	// レスポンスボディの検証
	var response Response
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("JSONのパースでエラー: %v", err)
	}

	if response.Error == nil {
		t.Error("Errorが設定されていません")
	}

	if response.Error.Code != "INTERNAL_SERVER_ERROR" {
		t.Errorf("Error.Code = %v, want INTERNAL_SERVER_ERROR", response.Error.Code)
	}
}

func TestUnprocessableEntityResponse(t *testing.T) {
	router := setupRouter()

	router.GET("/unprocessable", func(c *gin.Context) {
		UnprocessableEntityResponse(c, "Business rule violation")
	})

	req, _ := http.NewRequest("GET", "/unprocessable", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// ステータスコードの検証
	if w.Code != http.StatusUnprocessableEntity {
		t.Errorf("ステータスコード = %v, want %v", w.Code, http.StatusUnprocessableEntity)
	}

	// レスポンスボディの検証
	var response Response
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("JSONのパースでエラー: %v", err)
	}

	if response.Error == nil {
		t.Error("Errorが設定されていません")
	}

	if response.Error.Code != "UNPROCESSABLE_ENTITY" {
		t.Errorf("Error.Code = %v, want UNPROCESSABLE_ENTITY", response.Error.Code)
	}
}
