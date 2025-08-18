package usecase

import (
	"context"
	"errors"
	"os"
	"testing"

	"github.com/NRUG-SRE/slm-handson/backend/internal/domain/entity"
	"github.com/NRUG-SRE/slm-handson/backend/internal/usecase/mocks"
)

func TestOrderUseCase_CreateOrder(t *testing.T) {
	// エラーシミュレーションを無効化
	os.Setenv("ERROR_RATE", "0")
	os.Setenv("RESPONSE_TIME_MIN", "0")
	os.Setenv("RESPONSE_TIME_MAX", "0")
	defer func() {
		os.Unsetenv("ERROR_RATE")
		os.Unsetenv("RESPONSE_TIME_MIN")
		os.Unsetenv("RESPONSE_TIME_MAX")
	}()

	tests := []struct {
		name             string
		cartID           string
		setupOrderMock   func() *mocks.MockOrderRepository
		setupCartMock    func() *mocks.MockCartRepository
		setupProductMock func() *mocks.MockProductRepository
		expectError      bool
		checkResult      func(t *testing.T, order *entity.Order)
	}{
		{
			name:   "正常に注文を作成",
			cartID: "cart-123",
			setupOrderMock: func() *mocks.MockOrderRepository {
				mock := &mocks.MockOrderRepository{}
				mock.CreateFunc = func(ctx context.Context, order *entity.Order) error {
					return nil
				}
				return mock
			},
			setupCartMock: func() *mocks.MockCartRepository {
				mock := &mocks.MockCartRepository{}
				mock.GetByIDFunc = func(ctx context.Context, id string) (*entity.Cart, error) {
					cart := entity.NewCart()
					cart.ID = id
					// テスト用の商品とアイテムを作成
					product := entity.NewProduct("テスト商品", "説明", 1500, "image.jpg", 10)
					product.ID = "product-123"
					cart.AddItem(product, 2)
					return cart, nil
				}
				mock.SaveFunc = func(ctx context.Context, cart *entity.Cart) error {
					return nil
				}
				return mock
			},
			setupProductMock: func() *mocks.MockProductRepository {
				mock := &mocks.MockProductRepository{}
				mock.DecreaseStockFunc = func(ctx context.Context, id string, quantity int) error {
					return nil
				}
				return mock
			},
			expectError: false,
			checkResult: func(t *testing.T, order *entity.Order) {
				if order == nil {
					t.Fatal("注文がnilです")
				}
				if order.Status != entity.OrderStatusPending {
					t.Errorf("Status = %v, want %v", order.Status, entity.OrderStatusPending)
				}
				if len(order.Items) != 1 {
					t.Errorf("Items length = %v, want 1", len(order.Items))
				}
				if order.TotalAmount != 3000 {
					t.Errorf("TotalAmount = %v, want 3000", order.TotalAmount)
				}
			},
		},
		{
			name:   "空のカートIDでエラー",
			cartID: "",
			setupOrderMock: func() *mocks.MockOrderRepository {
				return &mocks.MockOrderRepository{}
			},
			setupCartMock: func() *mocks.MockCartRepository {
				return &mocks.MockCartRepository{}
			},
			setupProductMock: func() *mocks.MockProductRepository {
				return &mocks.MockProductRepository{}
			},
			expectError: true,
			checkResult: func(t *testing.T, order *entity.Order) {
				if order != nil {
					t.Error("注文がnilであるべきです")
				}
			},
		},
		{
			name:   "空のカートでエラー",
			cartID: "cart-123",
			setupOrderMock: func() *mocks.MockOrderRepository {
				return &mocks.MockOrderRepository{}
			},
			setupCartMock: func() *mocks.MockCartRepository {
				mock := &mocks.MockCartRepository{}
				mock.GetByIDFunc = func(ctx context.Context, id string) (*entity.Cart, error) {
					cart := entity.NewCart()
					cart.ID = id
					return cart, nil
				}
				return mock
			},
			setupProductMock: func() *mocks.MockProductRepository {
				return &mocks.MockProductRepository{}
			},
			expectError: true,
			checkResult: func(t *testing.T, order *entity.Order) {
				if order != nil {
					t.Error("注文がnilであるべきです")
				}
			},
		},
		// SLMハンズオン用に在庫チェックと在庫減少が無効化されたため、在庫関連エラーテストケースは削除
		{
			name:   "注文保存でエラー",
			cartID: "cart-123",
			setupOrderMock: func() *mocks.MockOrderRepository {
				mock := &mocks.MockOrderRepository{}
				mock.CreateFunc = func(ctx context.Context, order *entity.Order) error {
					return errors.New("database error")
				}
				return mock
			},
			setupCartMock: func() *mocks.MockCartRepository {
				mock := &mocks.MockCartRepository{}
				mock.GetByIDFunc = func(ctx context.Context, id string) (*entity.Cart, error) {
					cart := entity.NewCart()
					cart.ID = id
					product := entity.NewProduct("テスト商品", "説明", 1000, "image.jpg", 10)
					product.ID = "product-123"
					cart.AddItem(product, 2)
					return cart, nil
				}
				return mock
			},
			setupProductMock: func() *mocks.MockProductRepository {
				mock := &mocks.MockProductRepository{}
				mock.DecreaseStockFunc = func(ctx context.Context, id string, quantity int) error {
					return nil
				}
				return mock
			},
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
			mockOrderRepo := tt.setupOrderMock()
			mockCartRepo := tt.setupCartMock()
			mockProductRepo := tt.setupProductMock()
			uc := NewOrderUseCase(mockOrderRepo, mockCartRepo, mockProductRepo)
			ctx := context.Background()

			order, err := uc.CreateOrder(ctx, tt.cartID)

			if tt.expectError {
				if err == nil {
					t.Error("エラーが期待されましたが、エラーが発生しませんでした")
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

func TestOrderUseCase_GetOrder(t *testing.T) {
	tests := []struct {
		name        string
		orderID     string
		setupMock   func() *mocks.MockOrderRepository
		expectError bool
		checkResult func(t *testing.T, order *entity.Order)
	}{
		{
			name:    "正常に注文を取得",
			orderID: "order-123",
			setupMock: func() *mocks.MockOrderRepository {
				mock := &mocks.MockOrderRepository{}
				mock.GetByIDFunc = func(ctx context.Context, id string) (*entity.Order, error) {
					if id == "order-123" {
						cart := entity.NewCart()
						product := entity.NewProduct("テスト商品", "説明", 1000, "image.jpg", 10)
						cart.AddItem(product, 1)
						order, _ := entity.NewOrder(cart)
						order.ID = id
						return order, nil
					}
					return nil, entity.ErrOrderNotFound
				}
				return mock
			},
			expectError: false,
			checkResult: func(t *testing.T, order *entity.Order) {
				if order == nil {
					t.Fatal("注文がnilです")
				}
				if order.ID != "order-123" {
					t.Errorf("ID = %v, want order-123", order.ID)
				}
			},
		},
		{
			name:    "空のIDでエラー",
			orderID: "",
			setupMock: func() *mocks.MockOrderRepository {
				return &mocks.MockOrderRepository{}
			},
			expectError: true,
			checkResult: func(t *testing.T, order *entity.Order) {
				if order != nil {
					t.Error("注文がnilであるべきです")
				}
			},
		},
		{
			name:    "存在しない注文",
			orderID: "nonexistent",
			setupMock: func() *mocks.MockOrderRepository {
				mock := &mocks.MockOrderRepository{}
				mock.GetByIDFunc = func(ctx context.Context, id string) (*entity.Order, error) {
					return nil, entity.ErrOrderNotFound
				}
				return mock
			},
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
			mockOrderRepo := tt.setupMock()
			mockCartRepo := &mocks.MockCartRepository{}
			mockProductRepo := &mocks.MockProductRepository{}
			uc := NewOrderUseCase(mockOrderRepo, mockCartRepo, mockProductRepo)
			ctx := context.Background()

			order, err := uc.GetOrder(ctx, tt.orderID)

			if tt.expectError {
				if err == nil {
					t.Error("エラーが期待されましたが、エラーが発生しませんでした")
				}
			} else {
				if err != nil {
					t.Errorf("予期しないエラー: %v", err)
				}
			}

			tt.checkResult(t, order)

			// 空のIDの場合はリポジトリが呼ばれないことを確認
			if tt.orderID == "" {
				if len(mockOrderRepo.GetByIDCalls) != 0 {
					t.Errorf("GetByIDが呼ばれるべきではありません")
				}
			} else {
				if len(mockOrderRepo.GetByIDCalls) != 1 {
					t.Errorf("GetByIDの呼び出し回数 = %v, want 1", len(mockOrderRepo.GetByIDCalls))
				}
			}
		})
	}
}

func TestOrderUseCase_GetAllOrders(t *testing.T) {
	tests := []struct {
		name          string
		setupMock     func() *mocks.MockOrderRepository
		expectError   bool
		expectedCount int
	}{
		{
			name: "正常に全注文を取得",
			setupMock: func() *mocks.MockOrderRepository {
				mock := &mocks.MockOrderRepository{}
				mock.GetAllFunc = func(ctx context.Context) ([]*entity.Order, error) {
					// テスト用の注文を作成
					orders := make([]*entity.Order, 2)
					for i := 0; i < 2; i++ {
						cart := entity.NewCart()
						product := entity.NewProduct("商品", "説明", 1000, "image.jpg", 10)
						cart.AddItem(product, 1)
						order, _ := entity.NewOrder(cart)
						orders[i] = order
					}
					return orders, nil
				}
				return mock
			},
			expectError:   false,
			expectedCount: 2,
		},
		{
			name: "リポジトリエラー",
			setupMock: func() *mocks.MockOrderRepository {
				mock := &mocks.MockOrderRepository{}
				mock.GetAllFunc = func(ctx context.Context) ([]*entity.Order, error) {
					return nil, errors.New("database error")
				}
				return mock
			},
			expectError:   true,
			expectedCount: 0,
		},
		{
			name: "空の注文リスト",
			setupMock: func() *mocks.MockOrderRepository {
				mock := &mocks.MockOrderRepository{}
				mock.GetAllFunc = func(ctx context.Context) ([]*entity.Order, error) {
					return []*entity.Order{}, nil
				}
				return mock
			},
			expectError:   false,
			expectedCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockOrderRepo := tt.setupMock()
			mockCartRepo := &mocks.MockCartRepository{}
			mockProductRepo := &mocks.MockProductRepository{}
			uc := NewOrderUseCase(mockOrderRepo, mockCartRepo, mockProductRepo)
			ctx := context.Background()

			orders, err := uc.GetAllOrders(ctx)

			if tt.expectError {
				if err == nil {
					t.Error("エラーが期待されましたが、エラーが発生しませんでした")
				}
			} else {
				if err != nil {
					t.Errorf("予期しないエラー: %v", err)
				}
				if len(orders) != tt.expectedCount {
					t.Errorf("注文数 = %v, want %v", len(orders), tt.expectedCount)
				}
			}

			// リポジトリメソッドが呼ばれたことを確認
			if len(mockOrderRepo.GetAllCalls) != 1 {
				t.Errorf("GetAllの呼び出し回数 = %v, want 1", len(mockOrderRepo.GetAllCalls))
			}
		})
	}
}
