package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/NRUG-SRE/slm-handson/backend/internal/domain/entity"
	"github.com/NRUG-SRE/slm-handson/backend/internal/usecase/mocks"
)

func TestCartUseCase_GetCart(t *testing.T) {
	tests := []struct {
		name        string
		cartID      string
		setupMock   func() *mocks.MockCartRepository
		expectError bool
		checkResult func(t *testing.T, cart *entity.Cart)
	}{
		{
			name:   "正常にカートを取得",
			cartID: "cart-123",
			setupMock: func() *mocks.MockCartRepository {
				mock := &mocks.MockCartRepository{}
				mock.GetOrCreateFunc = func(ctx context.Context, id string) (*entity.Cart, error) {
					cart := entity.NewCart()
					cart.ID = id
					return cart, nil
				}
				return mock
			},
			expectError: false,
			checkResult: func(t *testing.T, cart *entity.Cart) {
				if cart == nil {
					t.Fatal("カートがnilです")
				}
				if cart.ID != "cart-123" {
					t.Errorf("ID = %v, want cart-123", cart.ID)
				}
			},
		},
		{
			name:   "空のIDでエラー",
			cartID: "",
			setupMock: func() *mocks.MockCartRepository {
				return &mocks.MockCartRepository{}
			},
			expectError: true,
			checkResult: func(t *testing.T, cart *entity.Cart) {
				if cart != nil {
					t.Error("カートがnilであるべきです")
				}
			},
		},
		{
			name:   "リポジトリエラー",
			cartID: "cart-123",
			setupMock: func() *mocks.MockCartRepository {
				mock := &mocks.MockCartRepository{}
				mock.GetOrCreateFunc = func(ctx context.Context, id string) (*entity.Cart, error) {
					return nil, errors.New("database error")
				}
				return mock
			},
			expectError: true,
			checkResult: func(t *testing.T, cart *entity.Cart) {
				if cart != nil {
					t.Error("カートがnilであるべきです")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCartRepo := tt.setupMock()
			mockProductRepo := &mocks.MockProductRepository{}
			uc := NewCartUseCase(mockCartRepo, mockProductRepo)
			ctx := context.Background()

			cart, err := uc.GetCart(ctx, tt.cartID)

			if tt.expectError {
				if err == nil {
					t.Error("エラーが期待されましたが、エラーが発生しませんでした")
				}
			} else {
				if err != nil {
					t.Errorf("予期しないエラー: %v", err)
				}
			}

			tt.checkResult(t, cart)
		})
	}
}

func TestCartUseCase_AddToCart(t *testing.T) {
	tests := []struct {
		name             string
		cartID           string
		productID        string
		quantity         int
		setupCartMock    func() *mocks.MockCartRepository
		setupProductMock func() *mocks.MockProductRepository
		expectError      bool
		checkResult      func(t *testing.T, cart *entity.Cart)
	}{
		{
			name:      "正常にカートに商品を追加",
			cartID:    "cart-123",
			productID: "product-456",
			quantity:  2,
			setupCartMock: func() *mocks.MockCartRepository {
				mock := &mocks.MockCartRepository{}
				mock.GetOrCreateFunc = func(ctx context.Context, id string) (*entity.Cart, error) {
					cart := entity.NewCart()
					cart.ID = id
					return cart, nil
				}
				mock.SaveFunc = func(ctx context.Context, cart *entity.Cart) error {
					return nil
				}
				return mock
			},
			setupProductMock: func() *mocks.MockProductRepository {
				mock := &mocks.MockProductRepository{}
				mock.GetByIDFunc = func(ctx context.Context, id string) (*entity.Product, error) {
					if id == "product-456" {
						product := entity.NewProduct("テスト商品", "説明", 1000, "image.jpg", 10)
						product.ID = id
						return product, nil
					}
					return nil, entity.ErrProductNotFound
				}
				return mock
			},
			expectError: false,
			checkResult: func(t *testing.T, cart *entity.Cart) {
				if cart == nil {
					t.Fatal("カートがnilです")
				}
				if len(cart.Items) != 1 {
					t.Errorf("アイテム数 = %v, want 1", len(cart.Items))
				}
				if cart.Items[0].Quantity != 2 {
					t.Errorf("数量 = %v, want 2", cart.Items[0].Quantity)
				}
			},
		},
		{
			name:      "空のカートIDでエラー",
			cartID:    "",
			productID: "product-456",
			quantity:  1,
			setupCartMock: func() *mocks.MockCartRepository {
				return &mocks.MockCartRepository{}
			},
			setupProductMock: func() *mocks.MockProductRepository {
				return &mocks.MockProductRepository{}
			},
			expectError: true,
			checkResult: func(t *testing.T, cart *entity.Cart) {
				if cart != nil {
					t.Error("カートがnilであるべきです")
				}
			},
		},
		{
			name:      "空の商品IDでエラー",
			cartID:    "cart-123",
			productID: "",
			quantity:  1,
			setupCartMock: func() *mocks.MockCartRepository {
				return &mocks.MockCartRepository{}
			},
			setupProductMock: func() *mocks.MockProductRepository {
				return &mocks.MockProductRepository{}
			},
			expectError: true,
			checkResult: func(t *testing.T, cart *entity.Cart) {
				if cart != nil {
					t.Error("カートがnilであるべきです")
				}
			},
		},
		{
			name:      "数量が0以下でエラー",
			cartID:    "cart-123",
			productID: "product-456",
			quantity:  0,
			setupCartMock: func() *mocks.MockCartRepository {
				return &mocks.MockCartRepository{}
			},
			setupProductMock: func() *mocks.MockProductRepository {
				return &mocks.MockProductRepository{}
			},
			expectError: true,
			checkResult: func(t *testing.T, cart *entity.Cart) {
				if cart != nil {
					t.Error("カートがnilであるべきです")
				}
			},
		},
		{
			name:      "存在しない商品でエラー",
			cartID:    "cart-123",
			productID: "nonexistent",
			quantity:  1,
			setupCartMock: func() *mocks.MockCartRepository {
				mock := &mocks.MockCartRepository{}
				mock.GetOrCreateFunc = func(ctx context.Context, id string) (*entity.Cart, error) {
					return entity.NewCart(), nil
				}
				return mock
			},
			setupProductMock: func() *mocks.MockProductRepository {
				mock := &mocks.MockProductRepository{}
				mock.GetByIDFunc = func(ctx context.Context, id string) (*entity.Product, error) {
					return nil, entity.ErrProductNotFound
				}
				return mock
			},
			expectError: true,
			checkResult: func(t *testing.T, cart *entity.Cart) {
				if cart != nil {
					t.Error("カートがnilであるべきです")
				}
			},
		},
		{
			name:      "在庫不足でエラー",
			cartID:    "cart-123",
			productID: "product-456",
			quantity:  15, // 在庫10を超える
			setupCartMock: func() *mocks.MockCartRepository {
				mock := &mocks.MockCartRepository{}
				mock.GetOrCreateFunc = func(ctx context.Context, id string) (*entity.Cart, error) {
					return entity.NewCart(), nil
				}
				return mock
			},
			setupProductMock: func() *mocks.MockProductRepository {
				mock := &mocks.MockProductRepository{}
				mock.GetByIDFunc = func(ctx context.Context, id string) (*entity.Product, error) {
					product := entity.NewProduct("テスト商品", "説明", 1000, "image.jpg", 10)
					product.ID = id
					return product, nil
				}
				return mock
			},
			expectError: true,
			checkResult: func(t *testing.T, cart *entity.Cart) {
				if cart != nil {
					t.Error("カートがnilであるべきです")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCartRepo := tt.setupCartMock()
			mockProductRepo := tt.setupProductMock()
			uc := NewCartUseCase(mockCartRepo, mockProductRepo)
			ctx := context.Background()

			cart, err := uc.AddToCart(ctx, tt.cartID, tt.productID, tt.quantity)

			if tt.expectError {
				if err == nil {
					t.Error("エラーが期待されましたが、エラーが発生しませんでした")
				}
			} else {
				if err != nil {
					t.Errorf("予期しないエラー: %v", err)
				}
			}

			tt.checkResult(t, cart)
		})
	}
}

func TestCartUseCase_UpdateCartItem(t *testing.T) {
	tests := []struct {
		name        string
		cartID      string
		itemID      string
		quantity    int
		setupMock   func() *mocks.MockCartRepository
		expectError bool
		checkResult func(t *testing.T, cart *entity.Cart)
	}{
		{
			name:     "正常にアイテム数量を更新",
			cartID:   "cart-123",
			itemID:   "item-456",
			quantity: 3,
			setupMock: func() *mocks.MockCartRepository {
				mock := &mocks.MockCartRepository{}
				mock.GetByIDFunc = func(ctx context.Context, id string) (*entity.Cart, error) {
					cart := entity.NewCart()
					cart.ID = id
					// テスト用の商品とアイテムを作成
					product := entity.NewProduct("テスト商品", "説明", 1000, "image.jpg", 10)
					item := entity.NewCartItem(product.ID, product, 1)
					item.ID = "item-456"
					cart.Items = append(cart.Items, item)
					cart.TotalAmount = 1000
					return cart, nil
				}
				mock.SaveFunc = func(ctx context.Context, cart *entity.Cart) error {
					return nil
				}
				return mock
			},
			expectError: false,
			checkResult: func(t *testing.T, cart *entity.Cart) {
				if cart == nil {
					t.Fatal("カートがnilです")
				}
				if len(cart.Items) != 1 {
					t.Errorf("アイテム数 = %v, want 1", len(cart.Items))
				}
				if cart.Items[0].Quantity != 3 {
					t.Errorf("数量 = %v, want 3", cart.Items[0].Quantity)
				}
			},
		},
		{
			name:     "空のカートIDでエラー",
			cartID:   "",
			itemID:   "item-456",
			quantity: 2,
			setupMock: func() *mocks.MockCartRepository {
				return &mocks.MockCartRepository{}
			},
			expectError: true,
			checkResult: func(t *testing.T, cart *entity.Cart) {
				if cart != nil {
					t.Error("カートがnilであるべきです")
				}
			},
		},
		{
			name:     "空のアイテムIDでエラー",
			cartID:   "cart-123",
			itemID:   "",
			quantity: 2,
			setupMock: func() *mocks.MockCartRepository {
				return &mocks.MockCartRepository{}
			},
			expectError: true,
			checkResult: func(t *testing.T, cart *entity.Cart) {
				if cart != nil {
					t.Error("カートがnilであるべきです")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCartRepo := tt.setupMock()
			mockProductRepo := &mocks.MockProductRepository{}
			uc := NewCartUseCase(mockCartRepo, mockProductRepo)
			ctx := context.Background()

			cart, err := uc.UpdateCartItem(ctx, tt.cartID, tt.itemID, tt.quantity)

			if tt.expectError {
				if err == nil {
					t.Error("エラーが期待されましたが、エラーが発生しませんでした")
				}
			} else {
				if err != nil {
					t.Errorf("予期しないエラー: %v", err)
				}
			}

			tt.checkResult(t, cart)
		})
	}
}

func TestCartUseCase_RemoveFromCart(t *testing.T) {
	tests := []struct {
		name        string
		cartID      string
		itemID      string
		setupMock   func() *mocks.MockCartRepository
		expectError bool
		checkResult func(t *testing.T, cart *entity.Cart)
	}{
		{
			name:   "正常にアイテムを削除",
			cartID: "cart-123",
			itemID: "item-456",
			setupMock: func() *mocks.MockCartRepository {
				mock := &mocks.MockCartRepository{}
				mock.GetByIDFunc = func(ctx context.Context, id string) (*entity.Cart, error) {
					cart := entity.NewCart()
					cart.ID = id
					// テスト用の商品とアイテムを作成
					product := entity.NewProduct("テスト商品", "説明", 1000, "image.jpg", 10)
					item := entity.NewCartItem(product.ID, product, 2)
					item.ID = "item-456"
					cart.Items = append(cart.Items, item)
					cart.TotalAmount = 2000
					return cart, nil
				}
				mock.SaveFunc = func(ctx context.Context, cart *entity.Cart) error {
					return nil
				}
				return mock
			},
			expectError: false,
			checkResult: func(t *testing.T, cart *entity.Cart) {
				if cart == nil {
					t.Fatal("カートがnilです")
				}
				if len(cart.Items) != 0 {
					t.Errorf("アイテム数 = %v, want 0", len(cart.Items))
				}
				if cart.TotalAmount != 0 {
					t.Errorf("合計金額 = %v, want 0", cart.TotalAmount)
				}
			},
		},
		{
			name:   "空のカートIDでエラー",
			cartID: "",
			itemID: "item-456",
			setupMock: func() *mocks.MockCartRepository {
				return &mocks.MockCartRepository{}
			},
			expectError: true,
			checkResult: func(t *testing.T, cart *entity.Cart) {
				if cart != nil {
					t.Error("カートがnilであるべきです")
				}
			},
		},
		{
			name:   "空のアイテムIDでエラー",
			cartID: "cart-123",
			itemID: "",
			setupMock: func() *mocks.MockCartRepository {
				return &mocks.MockCartRepository{}
			},
			expectError: true,
			checkResult: func(t *testing.T, cart *entity.Cart) {
				if cart != nil {
					t.Error("カートがnilであるべきです")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCartRepo := tt.setupMock()
			mockProductRepo := &mocks.MockProductRepository{}
			uc := NewCartUseCase(mockCartRepo, mockProductRepo)
			ctx := context.Background()

			cart, err := uc.RemoveFromCart(ctx, tt.cartID, tt.itemID)

			if tt.expectError {
				if err == nil {
					t.Error("エラーが期待されましたが、エラーが発生しませんでした")
				}
			} else {
				if err != nil {
					t.Errorf("予期しないエラー: %v", err)
				}
			}

			tt.checkResult(t, cart)
		})
	}
}

func TestCartUseCase_ClearCart(t *testing.T) {
	tests := []struct {
		name        string
		cartID      string
		setupMock   func() *mocks.MockCartRepository
		expectError bool
	}{
		{
			name:   "正常にカートをクリア",
			cartID: "cart-123",
			setupMock: func() *mocks.MockCartRepository {
				mock := &mocks.MockCartRepository{}
				mock.ClearFunc = func(ctx context.Context, id string) error {
					return nil
				}
				return mock
			},
			expectError: false,
		},
		{
			name:   "空のカートIDでエラー",
			cartID: "",
			setupMock: func() *mocks.MockCartRepository {
				return &mocks.MockCartRepository{}
			},
			expectError: true,
		},
		{
			name:   "リポジトリエラー",
			cartID: "cart-123",
			setupMock: func() *mocks.MockCartRepository {
				mock := &mocks.MockCartRepository{}
				mock.ClearFunc = func(ctx context.Context, id string) error {
					return errors.New("database error")
				}
				return mock
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCartRepo := tt.setupMock()
			mockProductRepo := &mocks.MockProductRepository{}
			uc := NewCartUseCase(mockCartRepo, mockProductRepo)
			ctx := context.Background()

			err := uc.ClearCart(ctx, tt.cartID)

			if tt.expectError {
				if err == nil {
					t.Error("エラーが期待されましたが、エラーが発生しませんでした")
				}
			} else {
				if err != nil {
					t.Errorf("予期しないエラー: %v", err)
				}
			}
		})
	}
}
