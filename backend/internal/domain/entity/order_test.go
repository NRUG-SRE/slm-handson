package entity

import (
	"testing"
	"time"
)

func TestNewOrder(t *testing.T) {
	tests := []struct {
		name          string
		setupCart     func() *Cart
		expectError   bool
		expectedItems int
		expectedTotal int
	}{
		{
			name: "正常な注文作成",
			setupCart: func() *Cart {
				cart := NewCart()
				product1 := NewProduct("商品1", "説明1", 1000, "image1.jpg", 10)
				product2 := NewProduct("商品2", "説明2", 500, "image2.jpg", 10)
				cart.AddItem(product1, 2)
				cart.AddItem(product2, 3)
				return cart
			},
			expectError:   false,
			expectedItems: 2,
			expectedTotal: 3500,
		},
		{
			name: "単一商品の注文",
			setupCart: func() *Cart {
				cart := NewCart()
				product := NewProduct("商品", "説明", 2000, "image.jpg", 5)
				cart.AddItem(product, 1)
				return cart
			},
			expectError:   false,
			expectedItems: 1,
			expectedTotal: 2000,
		},
		{
			name: "空のカートからの注文（エラー）",
			setupCart: func() *Cart {
				return NewCart()
			},
			expectError:   true,
			expectedItems: 0,
			expectedTotal: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cart := tt.setupCart()
			order, err := NewOrder(cart)

			if tt.expectError {
				if err != ErrEmptyCart {
					t.Errorf("Expected ErrEmptyCart, got %v", err)
				}
				if order != nil {
					t.Error("注文が作成されてしまいました")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if order == nil {
					t.Fatal("注文が作成されませんでした")
				}
				if order.ID == "" {
					t.Error("IDが生成されていません")
				}
				if len(order.Items) != tt.expectedItems {
					t.Errorf("Items length = %v, want %v", len(order.Items), tt.expectedItems)
				}
				if order.TotalAmount != tt.expectedTotal {
					t.Errorf("TotalAmount = %v, want %v", order.TotalAmount, tt.expectedTotal)
				}
				if order.Status != OrderStatusPending {
					t.Errorf("Status = %v, want %v", order.Status, OrderStatusPending)
				}
				if order.CreatedAt.IsZero() {
					t.Error("CreatedAtが設定されていません")
				}
				if order.UpdatedAt.IsZero() {
					t.Error("UpdatedAtが設定されていません")
				}

				// 各OrderItemの検証
				for i, orderItem := range order.Items {
					if orderItem.ID == "" {
						t.Errorf("OrderItem[%d].ID が生成されていません", i)
					}
					if orderItem.ProductID == "" {
						t.Errorf("OrderItem[%d].ProductID が設定されていません", i)
					}
					if orderItem.Product == nil {
						t.Errorf("OrderItem[%d].Product が設定されていません", i)
					}
					if orderItem.Quantity <= 0 {
						t.Errorf("OrderItem[%d].Quantity = %v, want > 0", i, orderItem.Quantity)
					}
					if orderItem.Price != orderItem.Product.Price {
						t.Errorf("OrderItem[%d].Price = %v, want %v", i, orderItem.Price, orderItem.Product.Price)
					}
					if orderItem.CreatedAt.IsZero() {
						t.Errorf("OrderItem[%d].CreatedAt が設定されていません", i)
					}
				}
			}
		})
	}
}

func TestOrder_Complete(t *testing.T) {
	tests := []struct {
		name         string
		setupOrder   func() *Order
		expectError  bool
		expectedStatus OrderStatus
	}{
		{
			name: "Pendingから完了",
			setupOrder: func() *Order {
				cart := NewCart()
				product := NewProduct("商品", "説明", 1000, "image.jpg", 10)
				cart.AddItem(product, 1)
				order, _ := NewOrder(cart)
				return order
			},
			expectError:    false,
			expectedStatus: OrderStatusCompleted,
		},
		{
			name: "既に完了済み（エラー）",
			setupOrder: func() *Order {
				cart := NewCart()
				product := NewProduct("商品", "説明", 1000, "image.jpg", 10)
				cart.AddItem(product, 1)
				order, _ := NewOrder(cart)
				order.Complete()
				return order
			},
			expectError:    true,
			expectedStatus: OrderStatusCompleted,
		},
		{
			name: "失敗状態から完了（エラー）",
			setupOrder: func() *Order {
				cart := NewCart()
				product := NewProduct("商品", "説明", 1000, "image.jpg", 10)
				cart.AddItem(product, 1)
				order, _ := NewOrder(cart)
				order.Fail()
				return order
			},
			expectError:    true,
			expectedStatus: OrderStatusFailed,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			order := tt.setupOrder()
			originalUpdatedAt := order.UpdatedAt

			time.Sleep(10 * time.Millisecond)
			err := order.Complete()

			if tt.expectError {
				if err != ErrInvalidOrderStatus {
					t.Errorf("Expected ErrInvalidOrderStatus, got %v", err)
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if !order.UpdatedAt.After(originalUpdatedAt) {
					t.Error("UpdatedAtが更新されていません")
				}
			}

			if order.Status != tt.expectedStatus {
				t.Errorf("Status = %v, want %v", order.Status, tt.expectedStatus)
			}
		})
	}
}

func TestOrder_Fail(t *testing.T) {
	tests := []struct {
		name         string
		setupOrder   func() *Order
		expectError  bool
		expectedStatus OrderStatus
	}{
		{
			name: "Pendingから失敗",
			setupOrder: func() *Order {
				cart := NewCart()
				product := NewProduct("商品", "説明", 1000, "image.jpg", 10)
				cart.AddItem(product, 1)
				order, _ := NewOrder(cart)
				return order
			},
			expectError:    false,
			expectedStatus: OrderStatusFailed,
		},
		{
			name: "完了済みから失敗（エラー）",
			setupOrder: func() *Order {
				cart := NewCart()
				product := NewProduct("商品", "説明", 1000, "image.jpg", 10)
				cart.AddItem(product, 1)
				order, _ := NewOrder(cart)
				order.Complete()
				return order
			},
			expectError:    true,
			expectedStatus: OrderStatusCompleted,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			order := tt.setupOrder()
			originalUpdatedAt := order.UpdatedAt

			time.Sleep(10 * time.Millisecond)
			err := order.Fail()

			if tt.expectError {
				if err != ErrInvalidOrderStatus {
					t.Errorf("Expected ErrInvalidOrderStatus, got %v", err)
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if !order.UpdatedAt.After(originalUpdatedAt) {
					t.Error("UpdatedAtが更新されていません")
				}
			}

			if order.Status != tt.expectedStatus {
				t.Errorf("Status = %v, want %v", order.Status, tt.expectedStatus)
			}
		})
	}
}

func TestOrder_Cancel(t *testing.T) {
	tests := []struct {
		name         string
		setupOrder   func() *Order
		expectError  bool
		expectedStatus OrderStatus
	}{
		{
			name: "Pendingからキャンセル",
			setupOrder: func() *Order {
				cart := NewCart()
				product := NewProduct("商品", "説明", 1000, "image.jpg", 10)
				cart.AddItem(product, 1)
				order, _ := NewOrder(cart)
				return order
			},
			expectError:    false,
			expectedStatus: OrderStatusCanceled,
		},
		{
			name: "失敗状態からキャンセル",
			setupOrder: func() *Order {
				cart := NewCart()
				product := NewProduct("商品", "説明", 1000, "image.jpg", 10)
				cart.AddItem(product, 1)
				order, _ := NewOrder(cart)
				order.Fail()
				return order
			},
			expectError:    false,
			expectedStatus: OrderStatusCanceled,
		},
		{
			name: "完了済みからキャンセル（エラー）",
			setupOrder: func() *Order {
				cart := NewCart()
				product := NewProduct("商品", "説明", 1000, "image.jpg", 10)
				cart.AddItem(product, 1)
				order, _ := NewOrder(cart)
				order.Complete()
				return order
			},
			expectError:    true,
			expectedStatus: OrderStatusCompleted,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			order := tt.setupOrder()
			originalUpdatedAt := order.UpdatedAt

			time.Sleep(10 * time.Millisecond)
			err := order.Cancel()

			if tt.expectError {
				if err != ErrInvalidOrderStatus {
					t.Errorf("Expected ErrInvalidOrderStatus, got %v", err)
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if !order.UpdatedAt.After(originalUpdatedAt) {
					t.Error("UpdatedAtが更新されていません")
				}
			}

			if order.Status != tt.expectedStatus {
				t.Errorf("Status = %v, want %v", order.Status, tt.expectedStatus)
			}
		})
	}
}

func TestOrder_GetItemCount(t *testing.T) {
	tests := []struct {
		name          string
		setupOrder    func() *Order
		expectedCount int
	}{
		{
			name: "単一商品",
			setupOrder: func() *Order {
				cart := NewCart()
				product := NewProduct("商品", "説明", 1000, "image.jpg", 10)
				cart.AddItem(product, 3)
				order, _ := NewOrder(cart)
				return order
			},
			expectedCount: 3,
		},
		{
			name: "複数商品",
			setupOrder: func() *Order {
				cart := NewCart()
				product1 := NewProduct("商品1", "説明1", 1000, "image1.jpg", 10)
				product2 := NewProduct("商品2", "説明2", 500, "image2.jpg", 10)
				cart.AddItem(product1, 2)
				cart.AddItem(product2, 3)
				order, _ := NewOrder(cart)
				return order
			},
			expectedCount: 5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			order := tt.setupOrder()
			count := order.GetItemCount()
			if count != tt.expectedCount {
				t.Errorf("GetItemCount() = %v, want %v", count, tt.expectedCount)
			}
		})
	}
}

func TestOrder_StatusCheckers(t *testing.T) {
	tests := []struct {
		name           string
		setupOrder     func() *Order
		isCompleted    bool
		isPending      bool
		isFailed       bool
		isCanceled     bool
	}{
		{
			name: "Pending状態",
			setupOrder: func() *Order {
				cart := NewCart()
				product := NewProduct("商品", "説明", 1000, "image.jpg", 10)
				cart.AddItem(product, 1)
				order, _ := NewOrder(cart)
				return order
			},
			isCompleted: false,
			isPending:   true,
			isFailed:    false,
			isCanceled:  false,
		},
		{
			name: "Completed状態",
			setupOrder: func() *Order {
				cart := NewCart()
				product := NewProduct("商品", "説明", 1000, "image.jpg", 10)
				cart.AddItem(product, 1)
				order, _ := NewOrder(cart)
				order.Complete()
				return order
			},
			isCompleted: true,
			isPending:   false,
			isFailed:    false,
			isCanceled:  false,
		},
		{
			name: "Failed状態",
			setupOrder: func() *Order {
				cart := NewCart()
				product := NewProduct("商品", "説明", 1000, "image.jpg", 10)
				cart.AddItem(product, 1)
				order, _ := NewOrder(cart)
				order.Fail()
				return order
			},
			isCompleted: false,
			isPending:   false,
			isFailed:    true,
			isCanceled:  false,
		},
		{
			name: "Canceled状態",
			setupOrder: func() *Order {
				cart := NewCart()
				product := NewProduct("商品", "説明", 1000, "image.jpg", 10)
				cart.AddItem(product, 1)
				order, _ := NewOrder(cart)
				order.Cancel()
				return order
			},
			isCompleted: false,
			isPending:   false,
			isFailed:    false,
			isCanceled:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			order := tt.setupOrder()

			if order.IsCompleted() != tt.isCompleted {
				t.Errorf("IsCompleted() = %v, want %v", order.IsCompleted(), tt.isCompleted)
			}
			if order.IsPending() != tt.isPending {
				t.Errorf("IsPending() = %v, want %v", order.IsPending(), tt.isPending)
			}
			if order.IsFailed() != tt.isFailed {
				t.Errorf("IsFailed() = %v, want %v", order.IsFailed(), tt.isFailed)
			}
			if order.IsCanceled() != tt.isCanceled {
				t.Errorf("IsCanceled() = %v, want %v", order.IsCanceled(), tt.isCanceled)
			}
		})
	}
}