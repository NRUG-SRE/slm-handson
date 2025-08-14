package entity

import (
	"testing"
	"time"
)

func TestNewProduct(t *testing.T) {
	tests := []struct {
		name        string
		productName string
		description string
		price       int
		imageURL    string
		stock       int
	}{
		{
			name:        "正常な商品作成",
			productName: "テスト商品",
			description: "これはテスト商品です",
			price:       1000,
			imageURL:    "https://example.com/image.jpg",
			stock:       10,
		},
		{
			name:        "在庫0の商品作成",
			productName: "在庫切れ商品",
			description: "在庫がない商品",
			price:       2000,
			imageURL:    "https://example.com/image2.jpg",
			stock:       0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			product := NewProduct(tt.productName, tt.description, tt.price, tt.imageURL, tt.stock)

			if product == nil {
				t.Fatal("商品が作成されませんでした")
			}
			if product.ID == "" {
				t.Error("IDが生成されていません")
			}
			if product.Name != tt.productName {
				t.Errorf("Name = %v, want %v", product.Name, tt.productName)
			}
			if product.Description != tt.description {
				t.Errorf("Description = %v, want %v", product.Description, tt.description)
			}
			if product.Price != tt.price {
				t.Errorf("Price = %v, want %v", product.Price, tt.price)
			}
			if product.ImageURL != tt.imageURL {
				t.Errorf("ImageURL = %v, want %v", product.ImageURL, tt.imageURL)
			}
			if product.Stock != tt.stock {
				t.Errorf("Stock = %v, want %v", product.Stock, tt.stock)
			}
			if product.CreatedAt.IsZero() {
				t.Error("CreatedAtが設定されていません")
			}
			if product.UpdatedAt.IsZero() {
				t.Error("UpdatedAtが設定されていません")
			}
		})
	}
}

func TestProduct_UpdateStock(t *testing.T) {
	product := NewProduct("テスト商品", "説明", 1000, "image.jpg", 10)
	originalUpdatedAt := product.UpdatedAt

	// 少し待機してから更新
	time.Sleep(10 * time.Millisecond)
	product.UpdateStock(20)

	if product.Stock != 20 {
		t.Errorf("Stock = %v, want %v", product.Stock, 20)
	}
	if !product.UpdatedAt.After(originalUpdatedAt) {
		t.Error("UpdatedAtが更新されていません")
	}
}

func TestProduct_DecreaseStock(t *testing.T) {
	tests := []struct {
		name          string
		initialStock  int
		decreaseBy    int
		expectedStock int
		expectError   bool
	}{
		{
			name:          "正常な在庫減少",
			initialStock:  10,
			decreaseBy:    3,
			expectedStock: 7,
			expectError:   false,
		},
		{
			name:          "在庫と同数の減少",
			initialStock:  5,
			decreaseBy:    5,
			expectedStock: 0,
			expectError:   false,
		},
		{
			name:          "在庫不足エラー",
			initialStock:  3,
			decreaseBy:    5,
			expectedStock: 3, // 変更されない
			expectError:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			product := NewProduct("テスト商品", "説明", 1000, "image.jpg", tt.initialStock)
			originalUpdatedAt := product.UpdatedAt

			time.Sleep(10 * time.Millisecond)
			err := product.DecreaseStock(tt.decreaseBy)

			if tt.expectError {
				if err != ErrInsufficientStock {
					t.Errorf("Expected ErrInsufficientStock, got %v", err)
				}
				if product.Stock != tt.expectedStock {
					t.Errorf("Stock should not change on error: got %v, want %v", product.Stock, tt.expectedStock)
				}
				if product.UpdatedAt != originalUpdatedAt {
					t.Error("UpdatedAt should not change on error")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if product.Stock != tt.expectedStock {
					t.Errorf("Stock = %v, want %v", product.Stock, tt.expectedStock)
				}
				if !product.UpdatedAt.After(originalUpdatedAt) {
					t.Error("UpdatedAtが更新されていません")
				}
			}
		})
	}
}

func TestProduct_IncreaseStock(t *testing.T) {
	product := NewProduct("テスト商品", "説明", 1000, "image.jpg", 10)
	originalUpdatedAt := product.UpdatedAt

	time.Sleep(10 * time.Millisecond)
	product.IncreaseStock(5)

	if product.Stock != 15 {
		t.Errorf("Stock = %v, want %v", product.Stock, 15)
	}
	if !product.UpdatedAt.After(originalUpdatedAt) {
		t.Error("UpdatedAtが更新されていません")
	}
}

func TestProduct_IsInStock(t *testing.T) {
	tests := []struct {
		name     string
		stock    int
		expected bool
	}{
		{
			name:     "在庫あり",
			stock:    10,
			expected: true,
		},
		{
			name:     "在庫1個",
			stock:    1,
			expected: true,
		},
		{
			name:     "在庫なし",
			stock:    0,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			product := NewProduct("テスト商品", "説明", 1000, "image.jpg", tt.stock)
			if result := product.IsInStock(); result != tt.expected {
				t.Errorf("IsInStock() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestProduct_IsAvailable(t *testing.T) {
	tests := []struct {
		name           string
		stock          int
		requestedQty   int
		expectedResult bool
	}{
		{
			name:           "十分な在庫あり",
			stock:          10,
			requestedQty:   5,
			expectedResult: true,
		},
		{
			name:           "ちょうど同じ在庫",
			stock:          5,
			requestedQty:   5,
			expectedResult: true,
		},
		{
			name:           "在庫不足",
			stock:          3,
			requestedQty:   5,
			expectedResult: false,
		},
		{
			name:           "在庫0で要求",
			stock:          0,
			requestedQty:   1,
			expectedResult: false,
		},
		{
			name:           "0個の要求",
			stock:          5,
			requestedQty:   0,
			expectedResult: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			product := NewProduct("テスト商品", "説明", 1000, "image.jpg", tt.stock)
			if result := product.IsAvailable(tt.requestedQty); result != tt.expectedResult {
				t.Errorf("IsAvailable(%d) = %v, want %v", tt.requestedQty, result, tt.expectedResult)
			}
		})
	}
}