package entity

import (
	"testing"
	"time"
)

func TestNewCart(t *testing.T) {
	cart := NewCart()

	if cart == nil {
		t.Fatal("カートが作成されませんでした")
	}
	if cart.ID == "" {
		t.Error("IDが生成されていません")
	}
	if len(cart.Items) != 0 {
		t.Errorf("Items length = %v, want 0", len(cart.Items))
	}
	if cart.TotalAmount != 0 {
		t.Errorf("TotalAmount = %v, want 0", cart.TotalAmount)
	}
	if cart.CreatedAt.IsZero() {
		t.Error("CreatedAtが設定されていません")
	}
	if cart.UpdatedAt.IsZero() {
		t.Error("UpdatedAtが設定されていません")
	}
}

func TestNewCartItem(t *testing.T) {
	product := NewProduct("テスト商品", "説明", 1000, "image.jpg", 10)
	cartItem := NewCartItem(product.ID, product, 3)

	if cartItem == nil {
		t.Fatal("カートアイテムが作成されませんでした")
	}
	if cartItem.ID == "" {
		t.Error("IDが生成されていません")
	}
	if cartItem.ProductID != product.ID {
		t.Errorf("ProductID = %v, want %v", cartItem.ProductID, product.ID)
	}
	if cartItem.Product != product {
		t.Error("Productが設定されていません")
	}
	if cartItem.Quantity != 3 {
		t.Errorf("Quantity = %v, want 3", cartItem.Quantity)
	}
	if cartItem.CreatedAt.IsZero() {
		t.Error("CreatedAtが設定されていません")
	}
	if cartItem.UpdatedAt.IsZero() {
		t.Error("UpdatedAtが設定されていません")
	}
}

func TestCart_AddItem(t *testing.T) {
	tests := []struct {
		name          string
		setupCart     func() *Cart
		product       *Product
		quantity      int
		expectError   bool
		expectedItems int
		expectedTotal int
	}{
		{
			name: "新規商品追加",
			setupCart: func() *Cart {
				return NewCart()
			},
			product:       NewProduct("商品A", "説明A", 1000, "imageA.jpg", 10),
			quantity:      2,
			expectError:   false,
			expectedItems: 1,
			expectedTotal: 2000,
		},
		{
			name: "既存商品の数量追加",
			setupCart: func() *Cart {
				cart := NewCart()
				product := NewProduct("商品B", "説明B", 500, "imageB.jpg", 10)
				product.ID = "product-b" // IDを固定
				cart.AddItem(product, 1)
				return cart
			},
			product: func() *Product {
				p := NewProduct("商品B", "説明B", 500, "imageB.jpg", 10)
				p.ID = "product-b" // 同じIDに設定
				return p
			}(),
			quantity:      2,
			expectError:   false,
			expectedItems: 1,
			expectedTotal: 1500,
		},
		{
			name: "在庫不足エラー",
			setupCart: func() *Cart {
				return NewCart()
			},
			product:       NewProduct("商品C", "説明C", 1000, "imageC.jpg", 2),
			quantity:      3,
			expectError:   true,
			expectedItems: 0,
			expectedTotal: 0,
		},
		{
			name: "既存商品で合計が在庫超過",
			setupCart: func() *Cart {
				cart := NewCart()
				product := NewProduct("商品D", "説明D", 500, "imageD.jpg", 5)
				product.ID = "product-d" // IDを固定
				cart.AddItem(product, 3)
				return cart
			},
			product: func() *Product {
				p := NewProduct("商品D", "説明D", 500, "imageD.jpg", 5)
				p.ID = "product-d" // 同じIDに設定
				return p
			}(),
			quantity:      3,
			expectError:   true,
			expectedItems: 1,
			expectedTotal: 1500, // 元の状態のまま
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cart := tt.setupCart()
			originalUpdatedAt := cart.UpdatedAt

			time.Sleep(10 * time.Millisecond)
			err := cart.AddItem(tt.product, tt.quantity)

			if tt.expectError {
				if err != ErrInsufficientStock {
					t.Errorf("Expected ErrInsufficientStock, got %v", err)
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if !cart.UpdatedAt.After(originalUpdatedAt) {
					t.Error("UpdatedAtが更新されていません")
				}
			}

			if len(cart.Items) != tt.expectedItems {
				t.Errorf("Items length = %v, want %v", len(cart.Items), tt.expectedItems)
			}
			if cart.TotalAmount != tt.expectedTotal {
				t.Errorf("TotalAmount = %v, want %v", cart.TotalAmount, tt.expectedTotal)
			}
		})
	}
}

func TestCart_UpdateItemQuantity(t *testing.T) {
	tests := []struct {
		name          string
		setupCart     func() (*Cart, string) // カートとアイテムIDを返す
		newQuantity   int
		expectError   bool
		expectedItems int
		expectedTotal int
	}{
		{
			name: "数量を増やす",
			setupCart: func() (*Cart, string) {
				cart := NewCart()
				product := NewProduct("商品", "説明", 1000, "image.jpg", 10)
				cart.AddItem(product, 2)
				return cart, cart.Items[0].ID
			},
			newQuantity:   5,
			expectError:   false,
			expectedItems: 1,
			expectedTotal: 5000,
		},
		{
			name: "数量を減らす",
			setupCart: func() (*Cart, string) {
				cart := NewCart()
				product := NewProduct("商品", "説明", 1000, "image.jpg", 10)
				cart.AddItem(product, 5)
				return cart, cart.Items[0].ID
			},
			newQuantity:   2,
			expectError:   false,
			expectedItems: 1,
			expectedTotal: 2000,
		},
		{
			name: "数量を0にする（削除）",
			setupCart: func() (*Cart, string) {
				cart := NewCart()
				product := NewProduct("商品", "説明", 1000, "image.jpg", 10)
				cart.AddItem(product, 3)
				return cart, cart.Items[0].ID
			},
			newQuantity:   0,
			expectError:   false,
			expectedItems: 0,
			expectedTotal: 0,
		},
		{
			name: "在庫超過エラー",
			setupCart: func() (*Cart, string) {
				cart := NewCart()
				product := NewProduct("商品", "説明", 1000, "image.jpg", 5)
				cart.AddItem(product, 2)
				return cart, cart.Items[0].ID
			},
			newQuantity:   6,
			expectError:   true,
			expectedItems: 1,
			expectedTotal: 2000, // 変更されない
		},
		{
			name: "存在しないアイテムID",
			setupCart: func() (*Cart, string) {
				cart := NewCart()
				product := NewProduct("商品", "説明", 1000, "image.jpg", 10)
				cart.AddItem(product, 2)
				return cart, "invalid-id"
			},
			newQuantity:   3,
			expectError:   true,
			expectedItems: 1,
			expectedTotal: 2000,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cart, itemID := tt.setupCart()
			originalUpdatedAt := cart.UpdatedAt

			time.Sleep(10 * time.Millisecond)
			err := cart.UpdateItemQuantity(itemID, tt.newQuantity)

			if tt.expectError {
				if err == nil {
					t.Error("エラーが期待されましたが、エラーが発生しませんでした")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if !cart.UpdatedAt.After(originalUpdatedAt) {
					t.Error("UpdatedAtが更新されていません")
				}
			}

			if len(cart.Items) != tt.expectedItems {
				t.Errorf("Items length = %v, want %v", len(cart.Items), tt.expectedItems)
			}
			if cart.TotalAmount != tt.expectedTotal {
				t.Errorf("TotalAmount = %v, want %v", cart.TotalAmount, tt.expectedTotal)
			}
		})
	}
}

func TestCart_RemoveItem(t *testing.T) {
	tests := []struct {
		name          string
		setupCart     func() (*Cart, string)
		expectError   bool
		expectedItems int
		expectedTotal int
	}{
		{
			name: "アイテムを削除",
			setupCart: func() (*Cart, string) {
				cart := NewCart()
				product1 := NewProduct("商品1", "説明1", 1000, "image1.jpg", 10)
				product2 := NewProduct("商品2", "説明2", 500, "image2.jpg", 10)
				cart.AddItem(product1, 2)
				cart.AddItem(product2, 3)
				return cart, cart.Items[0].ID
			},
			expectError:   false,
			expectedItems: 1,
			expectedTotal: 1500,
		},
		{
			name: "存在しないアイテムID",
			setupCart: func() (*Cart, string) {
				cart := NewCart()
				product := NewProduct("商品", "説明", 1000, "image.jpg", 10)
				cart.AddItem(product, 2)
				return cart, "invalid-id"
			},
			expectError:   true,
			expectedItems: 1,
			expectedTotal: 2000,
		},
		{
			name: "最後のアイテムを削除",
			setupCart: func() (*Cart, string) {
				cart := NewCart()
				product := NewProduct("商品", "説明", 1000, "image.jpg", 10)
				cart.AddItem(product, 2)
				return cart, cart.Items[0].ID
			},
			expectError:   false,
			expectedItems: 0,
			expectedTotal: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cart, itemID := tt.setupCart()
			originalUpdatedAt := cart.UpdatedAt

			time.Sleep(10 * time.Millisecond)
			err := cart.RemoveItem(itemID)

			if tt.expectError {
				if err != ErrItemNotFound {
					t.Errorf("Expected ErrItemNotFound, got %v", err)
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if !cart.UpdatedAt.After(originalUpdatedAt) {
					t.Error("UpdatedAtが更新されていません")
				}
			}

			if len(cart.Items) != tt.expectedItems {
				t.Errorf("Items length = %v, want %v", len(cart.Items), tt.expectedItems)
			}
			if cart.TotalAmount != tt.expectedTotal {
				t.Errorf("TotalAmount = %v, want %v", cart.TotalAmount, tt.expectedTotal)
			}
		})
	}
}

func TestCart_Clear(t *testing.T) {
	cart := NewCart()
	product1 := NewProduct("商品1", "説明1", 1000, "image1.jpg", 10)
	product2 := NewProduct("商品2", "説明2", 500, "image2.jpg", 10)
	cart.AddItem(product1, 2)
	cart.AddItem(product2, 3)

	originalUpdatedAt := cart.UpdatedAt
	time.Sleep(10 * time.Millisecond)

	cart.Clear()

	if len(cart.Items) != 0 {
		t.Errorf("Items length = %v, want 0", len(cart.Items))
	}
	if cart.TotalAmount != 0 {
		t.Errorf("TotalAmount = %v, want 0", cart.TotalAmount)
	}
	if !cart.UpdatedAt.After(originalUpdatedAt) {
		t.Error("UpdatedAtが更新されていません")
	}
}

func TestCart_GetItemCount(t *testing.T) {
	tests := []struct {
		name          string
		setupCart     func() *Cart
		expectedCount int
	}{
		{
			name: "空のカート",
			setupCart: func() *Cart {
				return NewCart()
			},
			expectedCount: 0,
		},
		{
			name: "単一商品",
			setupCart: func() *Cart {
				cart := NewCart()
				product := NewProduct("商品", "説明", 1000, "image.jpg", 10)
				cart.AddItem(product, 3)
				return cart
			},
			expectedCount: 3,
		},
		{
			name: "複数商品",
			setupCart: func() *Cart {
				cart := NewCart()
				product1 := NewProduct("商品1", "説明1", 1000, "image1.jpg", 10)
				product2 := NewProduct("商品2", "説明2", 500, "image2.jpg", 10)
				cart.AddItem(product1, 2)
				cart.AddItem(product2, 3)
				return cart
			},
			expectedCount: 5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cart := tt.setupCart()
			count := cart.GetItemCount()
			if count != tt.expectedCount {
				t.Errorf("GetItemCount() = %v, want %v", count, tt.expectedCount)
			}
		})
	}
}

func TestCart_IsEmpty(t *testing.T) {
	tests := []struct {
		name        string
		setupCart   func() *Cart
		expectEmpty bool
	}{
		{
			name: "空のカート",
			setupCart: func() *Cart {
				return NewCart()
			},
			expectEmpty: true,
		},
		{
			name: "商品があるカート",
			setupCart: func() *Cart {
				cart := NewCart()
				product := NewProduct("商品", "説明", 1000, "image.jpg", 10)
				cart.AddItem(product, 1)
				return cart
			},
			expectEmpty: false,
		},
		{
			name: "クリア後のカート",
			setupCart: func() *Cart {
				cart := NewCart()
				product := NewProduct("商品", "説明", 1000, "image.jpg", 10)
				cart.AddItem(product, 1)
				cart.Clear()
				return cart
			},
			expectEmpty: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cart := tt.setupCart()
			isEmpty := cart.IsEmpty()
			if isEmpty != tt.expectEmpty {
				t.Errorf("IsEmpty() = %v, want %v", isEmpty, tt.expectEmpty)
			}
		})
	}
}
