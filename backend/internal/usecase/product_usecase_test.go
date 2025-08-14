package usecase

import (
	"context"
	"errors"
	"os"
	"testing"

	"github.com/NRUG-SRE/slm-handson/backend/internal/domain/entity"
	"github.com/NRUG-SRE/slm-handson/backend/internal/usecase/mocks"
)

func TestProductUseCase_GetAllProducts(t *testing.T) {
	// テスト用に環境変数を設定（エラーとレスポンス時間のシミュレーションを無効化）
	os.Setenv("ERROR_RATE", "0")
	os.Setenv("RESPONSE_TIME_MIN", "0")
	os.Setenv("RESPONSE_TIME_MAX", "0")
	defer func() {
		os.Unsetenv("ERROR_RATE")
		os.Unsetenv("RESPONSE_TIME_MIN")
		os.Unsetenv("RESPONSE_TIME_MAX")
	}()

	tests := []struct {
		name          string
		setupMock     func() *mocks.MockProductRepository
		expectError   bool
		expectedCount int
	}{
		{
			name: "正常に全商品を取得",
			setupMock: func() *mocks.MockProductRepository {
				mock := &mocks.MockProductRepository{}
				mock.GetAllFunc = func(ctx context.Context) ([]*entity.Product, error) {
					return []*entity.Product{
						entity.NewProduct("商品1", "説明1", 1000, "image1.jpg", 10),
						entity.NewProduct("商品2", "説明2", 2000, "image2.jpg", 20),
					}, nil
				}
				return mock
			},
			expectError:   false,
			expectedCount: 2,
		},
		{
			name: "リポジトリエラー",
			setupMock: func() *mocks.MockProductRepository {
				mock := &mocks.MockProductRepository{}
				mock.GetAllFunc = func(ctx context.Context) ([]*entity.Product, error) {
					return nil, errors.New("database error")
				}
				return mock
			},
			expectError:   true,
			expectedCount: 0,
		},
		{
			name: "空の商品リスト",
			setupMock: func() *mocks.MockProductRepository {
				mock := &mocks.MockProductRepository{}
				mock.GetAllFunc = func(ctx context.Context) ([]*entity.Product, error) {
					return []*entity.Product{}, nil
				}
				return mock
			},
			expectError:   false,
			expectedCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := tt.setupMock()
			uc := NewProductUseCase(mockRepo)
			ctx := context.Background()

			products, err := uc.GetAllProducts(ctx)

			if tt.expectError {
				if err == nil {
					t.Error("エラーが期待されましたが、エラーが発生しませんでした")
				}
			} else {
				if err != nil {
					t.Errorf("予期しないエラー: %v", err)
				}
				if len(products) != tt.expectedCount {
					t.Errorf("商品数 = %v, want %v", len(products), tt.expectedCount)
				}
			}

			// リポジトリメソッドが呼ばれたことを確認
			if len(mockRepo.GetAllCalls) != 1 {
				t.Errorf("GetAllの呼び出し回数 = %v, want 1", len(mockRepo.GetAllCalls))
			}
		})
	}
}

func TestProductUseCase_GetProductByID(t *testing.T) {
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
		name        string
		productID   string
		setupMock   func() *mocks.MockProductRepository
		expectError bool
		checkResult func(t *testing.T, product *entity.Product)
	}{
		{
			name:      "正常に商品を取得",
			productID: "product-123",
			setupMock: func() *mocks.MockProductRepository {
				mock := &mocks.MockProductRepository{}
				mock.GetByIDFunc = func(ctx context.Context, id string) (*entity.Product, error) {
					if id == "product-123" {
						product := entity.NewProduct("テスト商品", "説明", 1500, "image.jpg", 5)
						product.ID = id
						return product, nil
					}
					return nil, entity.ErrProductNotFound
				}
				return mock
			},
			expectError: false,
			checkResult: func(t *testing.T, product *entity.Product) {
				if product == nil {
					t.Fatal("商品がnilです")
				}
				if product.ID != "product-123" {
					t.Errorf("ID = %v, want product-123", product.ID)
				}
				if product.Name != "テスト商品" {
					t.Errorf("Name = %v, want テスト商品", product.Name)
				}
			},
		},
		{
			name:      "空のIDでエラー",
			productID: "",
			setupMock: func() *mocks.MockProductRepository {
				mock := &mocks.MockProductRepository{}
				return mock
			},
			expectError: true,
			checkResult: func(t *testing.T, product *entity.Product) {
				if product != nil {
					t.Error("商品がnilであるべきです")
				}
			},
		},
		{
			name:      "存在しない商品",
			productID: "nonexistent",
			setupMock: func() *mocks.MockProductRepository {
				mock := &mocks.MockProductRepository{}
				mock.GetByIDFunc = func(ctx context.Context, id string) (*entity.Product, error) {
					return nil, entity.ErrProductNotFound
				}
				return mock
			},
			expectError: true,
			checkResult: func(t *testing.T, product *entity.Product) {
				if product != nil {
					t.Error("商品がnilであるべきです")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := tt.setupMock()
			uc := NewProductUseCase(mockRepo)
			ctx := context.Background()

			product, err := uc.GetProductByID(ctx, tt.productID)

			if tt.expectError {
				if err == nil {
					t.Error("エラーが期待されましたが、エラーが発生しませんでした")
				}
			} else {
				if err != nil {
					t.Errorf("予期しないエラー: %v", err)
				}
			}

			tt.checkResult(t, product)

			// 空のIDの場合はリポジトリが呼ばれないことを確認
			if tt.productID == "" {
				if len(mockRepo.GetByIDCalls) != 0 {
					t.Errorf("GetByIDが呼ばれるべきではありません")
				}
			} else {
				if len(mockRepo.GetByIDCalls) != 1 {
					t.Errorf("GetByIDの呼び出し回数 = %v, want 1", len(mockRepo.GetByIDCalls))
				}
			}
		})
	}
}

func TestProductUseCase_CreateProduct(t *testing.T) {
	tests := []struct {
		name        string
		productName string
		description string
		price       int
		imageURL    string
		stock       int
		setupMock   func() *mocks.MockProductRepository
		expectError bool
	}{
		{
			name:        "正常に商品を作成",
			productName: "新商品",
			description: "新商品の説明",
			price:       3000,
			imageURL:    "new.jpg",
			stock:       15,
			setupMock: func() *mocks.MockProductRepository {
				mock := &mocks.MockProductRepository{}
				mock.CreateFunc = func(ctx context.Context, product *entity.Product) error {
					return nil
				}
				return mock
			},
			expectError: false,
		},
		{
			name:        "名前が空でバリデーションエラー",
			productName: "",
			description: "説明",
			price:       1000,
			imageURL:    "image.jpg",
			stock:       10,
			setupMock: func() *mocks.MockProductRepository {
				return &mocks.MockProductRepository{}
			},
			expectError: true,
		},
		{
			name:        "価格が0以下でバリデーションエラー",
			productName: "商品",
			description: "説明",
			price:       0,
			imageURL:    "image.jpg",
			stock:       10,
			setupMock: func() *mocks.MockProductRepository {
				return &mocks.MockProductRepository{}
			},
			expectError: true,
		},
		{
			name:        "在庫が負の値でバリデーションエラー",
			productName: "商品",
			description: "説明",
			price:       1000,
			imageURL:    "image.jpg",
			stock:       -1,
			setupMock: func() *mocks.MockProductRepository {
				return &mocks.MockProductRepository{}
			},
			expectError: true,
		},
		{
			name:        "リポジトリエラー",
			productName: "商品",
			description: "説明",
			price:       1000,
			imageURL:    "image.jpg",
			stock:       10,
			setupMock: func() *mocks.MockProductRepository {
				mock := &mocks.MockProductRepository{}
				mock.CreateFunc = func(ctx context.Context, product *entity.Product) error {
					return errors.New("database error")
				}
				return mock
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := tt.setupMock()
			uc := NewProductUseCase(mockRepo)
			ctx := context.Background()

			product, err := uc.CreateProduct(ctx, tt.productName, tt.description, tt.price, tt.imageURL, tt.stock)

			if tt.expectError {
				if err == nil {
					t.Error("エラーが期待されましたが、エラーが発生しませんでした")
				}
				if product != nil {
					t.Error("エラー時は商品がnilであるべきです")
				}
			} else {
				if err != nil {
					t.Errorf("予期しないエラー: %v", err)
				}
				if product == nil {
					t.Error("商品が作成されませんでした")
				} else {
					if product.Name != tt.productName {
						t.Errorf("Name = %v, want %v", product.Name, tt.productName)
					}
					if product.Price != tt.price {
						t.Errorf("Price = %v, want %v", product.Price, tt.price)
					}
					if product.Stock != tt.stock {
						t.Errorf("Stock = %v, want %v", product.Stock, tt.stock)
					}
				}
			}
		})
	}
}

func TestProductUseCase_UpdateProduct(t *testing.T) {
	tests := []struct {
		name        string
		product     *entity.Product
		setupMock   func() *mocks.MockProductRepository
		expectError bool
	}{
		{
			name: "正常に商品を更新",
			product: func() *entity.Product {
				p := entity.NewProduct("更新商品", "更新説明", 2000, "update.jpg", 20)
				p.ID = "product-123"
				return p
			}(),
			setupMock: func() *mocks.MockProductRepository {
				mock := &mocks.MockProductRepository{}
				mock.UpdateFunc = func(ctx context.Context, product *entity.Product) error {
					return nil
				}
				return mock
			},
			expectError: false,
		},
		{
			name:    "nilの商品でエラー",
			product: nil,
			setupMock: func() *mocks.MockProductRepository {
				return &mocks.MockProductRepository{}
			},
			expectError: true,
		},
		{
			name: "空のIDでエラー",
			product: func() *entity.Product {
				p := entity.NewProduct("商品", "説明", 1000, "image.jpg", 10)
				p.ID = ""
				return p
			}(),
			setupMock: func() *mocks.MockProductRepository {
				return &mocks.MockProductRepository{}
			},
			expectError: true,
		},
		{
			name: "リポジトリエラー",
			product: func() *entity.Product {
				p := entity.NewProduct("商品", "説明", 1000, "image.jpg", 10)
				p.ID = "product-123"
				return p
			}(),
			setupMock: func() *mocks.MockProductRepository {
				mock := &mocks.MockProductRepository{}
				mock.UpdateFunc = func(ctx context.Context, product *entity.Product) error {
					return errors.New("database error")
				}
				return mock
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := tt.setupMock()
			uc := NewProductUseCase(mockRepo)
			ctx := context.Background()

			err := uc.UpdateProduct(ctx, tt.product)

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

func TestProductUseCase_DeleteProduct(t *testing.T) {
	tests := []struct {
		name        string
		productID   string
		setupMock   func() *mocks.MockProductRepository
		expectError bool
	}{
		{
			name:      "正常に商品を削除",
			productID: "product-123",
			setupMock: func() *mocks.MockProductRepository {
				mock := &mocks.MockProductRepository{}
				mock.DeleteFunc = func(ctx context.Context, id string) error {
					return nil
				}
				return mock
			},
			expectError: false,
		},
		{
			name:      "空のIDでエラー",
			productID: "",
			setupMock: func() *mocks.MockProductRepository {
				return &mocks.MockProductRepository{}
			},
			expectError: true,
		},
		{
			name:      "リポジトリエラー",
			productID: "product-123",
			setupMock: func() *mocks.MockProductRepository {
				mock := &mocks.MockProductRepository{}
				mock.DeleteFunc = func(ctx context.Context, id string) error {
					return errors.New("database error")
				}
				return mock
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := tt.setupMock()
			uc := NewProductUseCase(mockRepo)
			ctx := context.Background()

			err := uc.DeleteProduct(ctx, tt.productID)

			if tt.expectError {
				if err == nil {
					t.Error("エラーが期待されましたが、エラーが発生しませんでした")
				}
			} else {
				if err != nil {
					t.Errorf("予期しないエラー: %v", err)
				}
			}

			// 空のIDの場合はリポジトリが呼ばれないことを確認
			if tt.productID == "" {
				if len(mockRepo.DeleteCalls) != 0 {
					t.Errorf("Deleteが呼ばれるべきではありません")
				}
			} else {
				if len(mockRepo.DeleteCalls) != 1 {
					t.Errorf("Deleteの呼び出し回数 = %v, want 1", len(mockRepo.DeleteCalls))
				}
			}
		})
	}
}