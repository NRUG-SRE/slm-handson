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
	// SLMハンズオン用に在庫減少が無効化されたため、常に成功することを確認
	tests := []struct {
		name         string
		initialStock int
		decreaseBy   int
	}{
		{
			name:         "在庫管理無効化確認1",
			initialStock: 10,
			decreaseBy:   3,
		},
		{
			name:         "在庫管理無効化確認2",
			initialStock: 0,
			decreaseBy:   5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			product := NewProduct("テスト商品", "説明", 1000, "image.jpg", tt.initialStock)
			originalStock := product.Stock
			originalUpdatedAt := product.UpdatedAt

			time.Sleep(10 * time.Millisecond)
			err := product.DecreaseStock(tt.decreaseBy)

			// SLMハンズオン用に在庫減少は無効化されているため、常に成功する
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
			// 在庫は変更されないことを確認
			if product.Stock != originalStock {
				t.Errorf("Stock should not change: got %v, want %v", product.Stock, originalStock)
			}
			// UpdatedAtは更新されることを確認
			if !product.UpdatedAt.After(originalUpdatedAt) {
				t.Error("更新日時が更新されていません")
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
	// SLMハンズオン用に在庫チェックが無効化されたため、常にtrueを返すことを確認
	tests := []struct {
		name         string
		stock        int
		requestedQty int
	}{
		{
			name:         "在庫チェック無効化確認1",
			stock:        10,
			requestedQty: 5,
		},
		{
			name:         "在庫チェック無効化確認2",
			stock:        0,
			requestedQty: 1,
		},
		{
			name:         "在庫チェック無効化確認3",
			stock:        3,
			requestedQty: 100,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			product := NewProduct("テスト商品", "説明", 1000, "image.jpg", tt.stock)
			// SLMハンズオン用に在庫チェックが無効化されているため、常にtrueを返す
			if result := product.IsAvailable(tt.requestedQty); result != true {
				t.Errorf("IsAvailable(%d) = %v, want %v", tt.requestedQty, result, true)
			}
		})
	}
}