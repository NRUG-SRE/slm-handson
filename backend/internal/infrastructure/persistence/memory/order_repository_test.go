package memory

import (
	"context"
	"sync"
	"testing"

	"github.com/NRUG-SRE/slm-handson/backend/internal/domain/entity"
)

func createTestOrder() *entity.Order {
	cart := entity.NewCart()
	product := entity.NewProduct("テスト商品", "説明", 1000, "image.jpg", 10)
	cart.AddItem(product, 2)
	order, _ := entity.NewOrder(cart)
	return order
}

func TestOrderRepository_GetAll(t *testing.T) {
	repo := NewOrderRepository()
	ctx := context.Background()

	// 初期状態では注文が存在しない
	orders, err := repo.GetAll(ctx)
	if err != nil {
		t.Errorf("予期しないエラー: %v", err)
	}
	if len(orders) != 0 {
		t.Errorf("注文数 = %v, want 0", len(orders))
	}

	// テスト用の注文を作成
	order1 := createTestOrder()
	order2 := createTestOrder()

	repo.Create(ctx, order1)
	repo.Create(ctx, order2)

	// 全注文を取得
	orders, err = repo.GetAll(ctx)
	if err != nil {
		t.Errorf("予期しないエラー: %v", err)
	}
	if len(orders) != 2 {
		t.Errorf("注文数 = %v, want 2", len(orders))
	}

	// 注文IDが正しいことを確認
	orderIDs := make(map[string]bool)
	for _, order := range orders {
		orderIDs[order.ID] = true
	}
	if !orderIDs[order1.ID] {
		t.Error("order1が見つかりません")
	}
	if !orderIDs[order2.ID] {
		t.Error("order2が見つかりません")
	}
}

func TestOrderRepository_GetByID(t *testing.T) {
	repo := NewOrderRepository()
	ctx := context.Background()

	// テスト用の注文を作成
	testOrder := createTestOrder()
	repo.Create(ctx, testOrder)

	tests := []struct {
		name        string
		orderID     string
		expectError bool
		checkResult func(t *testing.T, order *entity.Order)
	}{
		{
			name:    "存在する注文を取得",
			orderID: testOrder.ID,
			expectError: false,
			checkResult: func(t *testing.T, order *entity.Order) {
				if order == nil {
					t.Fatal("注文がnilです")
				}
				if order.ID != testOrder.ID {
					t.Errorf("ID = %v, want %v", order.ID, testOrder.ID)
				}
				if order.TotalAmount != testOrder.TotalAmount {
					t.Errorf("TotalAmount = %v, want %v", order.TotalAmount, testOrder.TotalAmount)
				}
				if len(order.Items) != len(testOrder.Items) {
					t.Errorf("Items length = %v, want %v", len(order.Items), len(testOrder.Items))
				}
			},
		},
		{
			name:    "存在しない注文を取得",
			orderID: "nonexistent-order",
			expectError: true,
			checkResult: func(t *testing.T, order *entity.Order) {
				if order != nil {
					t.Error("注文がnilであるべきです")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			order, err := repo.GetByID(ctx, tt.orderID)

			if tt.expectError {
				if err != entity.ErrOrderNotFound {
					t.Errorf("Expected ErrOrderNotFound, got %v", err)
				}
			} else {
				if err != nil {
					t.Errorf("予期しないエラー: %v", err)
				}
			}

			tt.checkResult(t, order)
		})
	}
}

func TestOrderRepository_Create(t *testing.T) {
	repo := NewOrderRepository()
	ctx := context.Background()

	// 作成前の注文数を確認
	beforeOrders, _ := repo.GetAll(ctx)
	beforeCount := len(beforeOrders)

	// 新しい注文を作成
	newOrder := createTestOrder()

	err := repo.Create(ctx, newOrder)
	if err != nil {
		t.Errorf("予期しないエラー: %v", err)
	}

	// 作成後の注文数を確認
	afterOrders, _ := repo.GetAll(ctx)
	afterCount := len(afterOrders)

	if afterCount != beforeCount+1 {
		t.Errorf("注文数 = %v, want %v", afterCount, beforeCount+1)
	}

	// 作成した注文を取得して確認
	createdOrder, err := repo.GetByID(ctx, newOrder.ID)
	if err != nil {
		t.Errorf("作成した注文の取得でエラー: %v", err)
	}
	if createdOrder.ID != newOrder.ID {
		t.Errorf("ID = %v, want %v", createdOrder.ID, newOrder.ID)
	}
	if createdOrder.TotalAmount != newOrder.TotalAmount {
		t.Errorf("TotalAmount = %v, want %v", createdOrder.TotalAmount, newOrder.TotalAmount)
	}
}

func TestOrderRepository_Update(t *testing.T) {
	repo := NewOrderRepository()
	ctx := context.Background()

	// テスト用の注文を作成
	testOrder := createTestOrder()
	repo.Create(ctx, testOrder)

	tests := []struct {
		name        string
		order       *entity.Order
		expectError bool
	}{
		{
			name: "存在する注文を更新",
			order: func() *entity.Order {
				o := testOrder
				o.Complete() // ステータスを変更
				return o
			}(),
			expectError: false,
		},
		{
			name: "存在しない注文を更新",
			order: func() *entity.Order {
				o := createTestOrder()
				o.ID = "nonexistent-order"
				return o
			}(),
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.Update(ctx, tt.order)

			if tt.expectError {
				if err != entity.ErrOrderNotFound {
					t.Errorf("Expected ErrOrderNotFound, got %v", err)
				}
			} else {
				if err != nil {
					t.Errorf("予期しないエラー: %v", err)
				}

				// 更新された注文を取得して確認
				updatedOrder, err := repo.GetByID(ctx, tt.order.ID)
				if err != nil {
					t.Errorf("更新された注文の取得でエラー: %v", err)
				}
				if updatedOrder.Status != tt.order.Status {
					t.Errorf("Status = %v, want %v", updatedOrder.Status, tt.order.Status)
				}
			}
		})
	}
}

func TestOrderRepository_Delete(t *testing.T) {
	repo := NewOrderRepository()
	ctx := context.Background()

	// テスト用の注文を作成
	testOrder := createTestOrder()
	repo.Create(ctx, testOrder)

	tests := []struct {
		name        string
		orderID     string
		expectError bool
	}{
		{
			name:    "存在する注文を削除",
			orderID: testOrder.ID,
			expectError: false,
		},
		{
			name:    "存在しない注文を削除",
			orderID: "nonexistent-order",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 削除前の注文数を確認
			beforeOrders, _ := repo.GetAll(ctx)
			beforeCount := len(beforeOrders)

			err := repo.Delete(ctx, tt.orderID)

			if tt.expectError {
				if err != entity.ErrOrderNotFound {
					t.Errorf("Expected ErrOrderNotFound, got %v", err)
				}

				// エラーの場合は注文数が変わらないことを確認
				afterOrders, _ := repo.GetAll(ctx)
				if len(afterOrders) != beforeCount {
					t.Errorf("注文数が変わってしまいました: %v -> %v", beforeCount, len(afterOrders))
				}
			} else {
				if err != nil {
					t.Errorf("予期しないエラー: %v", err)
				}

				// 削除後の注文数を確認
				afterOrders, _ := repo.GetAll(ctx)
				afterCount := len(afterOrders)

				if afterCount != beforeCount-1 {
					t.Errorf("注文数 = %v, want %v", afterCount, beforeCount-1)
				}

				// 削除された注文が取得できないことを確認
				_, err := repo.GetByID(ctx, tt.orderID)
				if err != entity.ErrOrderNotFound {
					t.Errorf("削除された注文が取得できてしまいました")
				}
			}
		})
	}
}

// 並行処理の安全性をテスト
func TestOrderRepository_ConcurrentAccess(t *testing.T) {
	repo := NewOrderRepository()
	ctx := context.Background()

	const numGoroutines = 10
	const operationsPerGoroutine = 10

	var wg sync.WaitGroup
	wg.Add(numGoroutines * 3) // Create, Update, GetByIDの3つの操作

	// 並行でCreateを実行
	orderIDs := make([]string, numGoroutines*operationsPerGoroutine)

	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			defer wg.Done()
			for j := 0; j < operationsPerGoroutine; j++ {
				order := createTestOrder()
				err := repo.Create(ctx, order)
				if err != nil {
					t.Errorf("並行Createでエラー: %v", err)
				}
				orderIDs[id*operationsPerGoroutine+j] = order.ID
			}
		}(i)
	}

	// 並行でUpdateを実行
	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			defer wg.Done()
			for j := 0; j < operationsPerGoroutine; j++ {
				// まず注文を作成
				order := createTestOrder()
				repo.Create(ctx, order)
				
				// その後更新
				order.Complete()
				err := repo.Update(ctx, order)
				if err != nil {
					t.Errorf("並行Updateでエラー: %v", err)
				}
			}
		}(i)
	}

	// 並行でGetAllを実行
	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			defer wg.Done()
			for j := 0; j < operationsPerGoroutine; j++ {
				_, err := repo.GetAll(ctx)
				if err != nil {
					t.Errorf("並行GetAllでエラー: %v", err)
				}
			}
		}(i)
	}

	wg.Wait()

	// 最終的にデータが破損していないことを確認
	finalOrders, err := repo.GetAll(ctx)
	if err != nil {
		t.Errorf("最終確認でエラー: %v", err)
	}
	if finalOrders == nil {
		t.Error("注文リストが破損しています")
	}

	// 作成された注文が取得できることを確認（一部のみチェック）
	for i := 0; i < 10 && i < len(orderIDs); i++ {
		if orderIDs[i] != "" {
			order, err := repo.GetByID(ctx, orderIDs[i])
			if err != nil {
				t.Errorf("作成された注文の取得でエラー: %v", err)
			}
			if order == nil {
				t.Error("作成された注文がnilです")
			}
		}
	}
}

// 大量データ処理のテスト
func TestOrderRepository_LargeDataSet(t *testing.T) {
	repo := NewOrderRepository()
	ctx := context.Background()

	const numOrders = 1000

	// 大量の注文を作成
	for i := 0; i < numOrders; i++ {
		order := createTestOrder()
		err := repo.Create(ctx, order)
		if err != nil {
			t.Errorf("注文作成でエラー: %v", err)
		}
	}

	// 全注文を取得
	orders, err := repo.GetAll(ctx)
	if err != nil {
		t.Errorf("全注文取得でエラー: %v", err)
	}
	if len(orders) != numOrders {
		t.Errorf("注文数 = %v, want %v", len(orders), numOrders)
	}

	// 各注文が正しく保存されていることを確認（サンプリング）
	for i := 0; i < 100; i++ {
		order := orders[i]
		retrievedOrder, err := repo.GetByID(ctx, order.ID)
		if err != nil {
			t.Errorf("注文の取得でエラー: %v", err)
		}
		if retrievedOrder.ID != order.ID {
			t.Errorf("ID = %v, want %v", retrievedOrder.ID, order.ID)
		}
		if len(retrievedOrder.Items) == 0 {
			t.Error("注文アイテムが空です")
		}
	}

	// 一部の注文を削除
	for i := 0; i < 100; i++ {
		order := orders[i]
		err := repo.Delete(ctx, order.ID)
		if err != nil {
			t.Errorf("注文削除でエラー: %v", err)
		}
	}

	// 削除後の注文数を確認
	remainingOrders, err := repo.GetAll(ctx)
	if err != nil {
		t.Errorf("残存注文取得でエラー: %v", err)
	}
	if len(remainingOrders) != numOrders-100 {
		t.Errorf("残存注文数 = %v, want %v", len(remainingOrders), numOrders-100)
	}
}