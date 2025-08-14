package memory

import (
	"context"
	"sync"
	"testing"

	"github.com/NRUG-SRE/slm-handson/backend/internal/domain/entity"
)

func TestProductRepository_GetAll(t *testing.T) {
	repo := NewProductRepository()
	ctx := context.Background()

	products, err := repo.GetAll(ctx)

	if err != nil {
		t.Errorf("予期しないエラー: %v", err)
	}
	if len(products) != 6 { // 初期データの数
		t.Errorf("商品数 = %v, want 6", len(products))
	}

	// 初期データの一部を確認
	found := false
	for _, product := range products {
		if product.Name == "ワイヤレスヘッドホン" {
			found = true
			if product.Price != 25000 {
				t.Errorf("価格 = %v, want 25000", product.Price)
			}
			break
		}
	}
	if !found {
		t.Error("ワイヤレスヘッドホンが見つかりません")
	}
}

func TestProductRepository_GetByID(t *testing.T) {
	repo := NewProductRepository()
	ctx := context.Background()

	// 最初に全商品を取得してIDを確認
	products, _ := repo.GetAll(ctx)
	if len(products) == 0 {
		t.Fatal("テスト用の商品が存在しません")
	}

	tests := []struct {
		name        string
		productID   string
		expectError bool
		checkResult func(t *testing.T, product *entity.Product)
	}{
		{
			name:      "存在する商品を取得",
			productID: products[0].ID,
			expectError: false,
			checkResult: func(t *testing.T, product *entity.Product) {
				if product == nil {
					t.Fatal("商品がnilです")
				}
				if product.ID != products[0].ID {
					t.Errorf("ID = %v, want %v", product.ID, products[0].ID)
				}
			},
		},
		{
			name:      "存在しない商品",
			productID: "nonexistent-id",
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
			product, err := repo.GetByID(ctx, tt.productID)

			if tt.expectError {
				if err != entity.ErrProductNotFound {
					t.Errorf("Expected ErrProductNotFound, got %v", err)
				}
			} else {
				if err != nil {
					t.Errorf("予期しないエラー: %v", err)
				}
			}

			tt.checkResult(t, product)
		})
	}
}

func TestProductRepository_Create(t *testing.T) {
	repo := NewProductRepository()
	ctx := context.Background()

	// 新しい商品を作成
	newProduct := entity.NewProduct(
		"テスト商品",
		"テスト用の商品説明",
		1000,
		"/images/test.svg",
		5,
	)

	// 作成前の商品数を確認
	beforeProducts, _ := repo.GetAll(ctx)
	beforeCount := len(beforeProducts)

	// 商品を作成
	err := repo.Create(ctx, newProduct)
	if err != nil {
		t.Errorf("予期しないエラー: %v", err)
	}

	// 作成後の商品数を確認
	afterProducts, _ := repo.GetAll(ctx)
	afterCount := len(afterProducts)

	if afterCount != beforeCount+1 {
		t.Errorf("商品数 = %v, want %v", afterCount, beforeCount+1)
	}

	// 作成した商品を取得して確認
	createdProduct, err := repo.GetByID(ctx, newProduct.ID)
	if err != nil {
		t.Errorf("作成した商品の取得でエラー: %v", err)
	}
	if createdProduct.Name != newProduct.Name {
		t.Errorf("Name = %v, want %v", createdProduct.Name, newProduct.Name)
	}
}

func TestProductRepository_Update(t *testing.T) {
	repo := NewProductRepository()
	ctx := context.Background()

	// 既存の商品を取得
	products, _ := repo.GetAll(ctx)
	if len(products) == 0 {
		t.Fatal("テスト用の商品が存在しません")
	}

	tests := []struct {
		name        string
		product     *entity.Product
		expectError bool
	}{
		{
			name: "存在する商品を更新",
			product: func() *entity.Product {
				p := products[0]
				p.Name = "更新された商品名"
				p.Price = 99999
				return p
			}(),
			expectError: false,
		},
		{
			name: "存在しない商品を更新",
			product: func() *entity.Product {
				p := entity.NewProduct("存在しない", "説明", 1000, "image.jpg", 10)
				p.ID = "nonexistent-id"
				return p
			}(),
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.Update(ctx, tt.product)

			if tt.expectError {
				if err != entity.ErrProductNotFound {
					t.Errorf("Expected ErrProductNotFound, got %v", err)
				}
			} else {
				if err != nil {
					t.Errorf("予期しないエラー: %v", err)
				}

				// 更新された商品を取得して確認
				updatedProduct, err := repo.GetByID(ctx, tt.product.ID)
				if err != nil {
					t.Errorf("更新された商品の取得でエラー: %v", err)
				}
				if updatedProduct.Name != tt.product.Name {
					t.Errorf("Name = %v, want %v", updatedProduct.Name, tt.product.Name)
				}
				if updatedProduct.Price != tt.product.Price {
					t.Errorf("Price = %v, want %v", updatedProduct.Price, tt.product.Price)
				}
			}
		})
	}
}

func TestProductRepository_Delete(t *testing.T) {
	repo := NewProductRepository()
	ctx := context.Background()

	// 既存の商品を取得
	products, _ := repo.GetAll(ctx)
	if len(products) == 0 {
		t.Fatal("テスト用の商品が存在しません")
	}

	tests := []struct {
		name        string
		productID   string
		expectError bool
	}{
		{
			name:      "存在する商品を削除",
			productID: products[0].ID,
			expectError: false,
		},
		{
			name:      "存在しない商品を削除",
			productID: "nonexistent-id",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 削除前の商品数を確認
			beforeProducts, _ := repo.GetAll(ctx)
			beforeCount := len(beforeProducts)

			err := repo.Delete(ctx, tt.productID)

			if tt.expectError {
				if err != entity.ErrProductNotFound {
					t.Errorf("Expected ErrProductNotFound, got %v", err)
				}

				// エラーの場合は商品数が変わらないことを確認
				afterProducts, _ := repo.GetAll(ctx)
				if len(afterProducts) != beforeCount {
					t.Errorf("商品数が変わってしまいました: %v -> %v", beforeCount, len(afterProducts))
				}
			} else {
				if err != nil {
					t.Errorf("予期しないエラー: %v", err)
				}

				// 削除後の商品数を確認
				afterProducts, _ := repo.GetAll(ctx)
				afterCount := len(afterProducts)

				if afterCount != beforeCount-1 {
					t.Errorf("商品数 = %v, want %v", afterCount, beforeCount-1)
				}

				// 削除された商品が取得できないことを確認
				_, err := repo.GetByID(ctx, tt.productID)
				if err != entity.ErrProductNotFound {
					t.Errorf("削除された商品が取得できてしまいました")
				}
			}
		})
	}
}

func TestProductRepository_UpdateStock(t *testing.T) {
	repo := NewProductRepository()
	ctx := context.Background()

	// 既存の商品を取得
	products, _ := repo.GetAll(ctx)
	if len(products) == 0 {
		t.Fatal("テスト用の商品が存在しません")
	}

	tests := []struct {
		name        string
		productID   string
		newStock    int
		expectError bool
	}{
		{
			name:      "存在する商品の在庫を更新",
			productID: products[0].ID,
			newStock:  100,
			expectError: false,
		},
		{
			name:      "存在しない商品の在庫を更新",
			productID: "nonexistent-id",
			newStock:  50,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.UpdateStock(ctx, tt.productID, tt.newStock)

			if tt.expectError {
				if err != entity.ErrProductNotFound {
					t.Errorf("Expected ErrProductNotFound, got %v", err)
				}
			} else {
				if err != nil {
					t.Errorf("予期しないエラー: %v", err)
				}

				// 更新された在庫を確認
				product, err := repo.GetByID(ctx, tt.productID)
				if err != nil {
					t.Errorf("商品の取得でエラー: %v", err)
				}
				if product.Stock != tt.newStock {
					t.Errorf("Stock = %v, want %v", product.Stock, tt.newStock)
				}
			}
		})
	}
}

func TestProductRepository_DecreaseStock(t *testing.T) {
	repo := NewProductRepository()
	ctx := context.Background()

	// テスト用の商品を作成
	testProduct := entity.NewProduct("在庫テスト商品", "説明", 1000, "image.jpg", 10)
	repo.Create(ctx, testProduct)

	tests := []struct {
		name        string
		productID   string
		quantity    int
		expectError bool
		expectedStock int
	}{
		{
			name:        "正常な在庫減少",
			productID:   testProduct.ID,
			quantity:    3,
			expectError: false,
			expectedStock: 7,
		},
		{
			name:        "在庫不足",
			productID:   testProduct.ID,
			quantity:    10, // 現在の在庫7より多い
			expectError: true,
			expectedStock: 7, // 変更されない
		},
		{
			name:        "存在しない商品",
			productID:   "nonexistent-id",
			quantity:    1,
			expectError: true,
			expectedStock: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.DecreaseStock(ctx, tt.productID, tt.quantity)

			if tt.expectError {
				if err == nil {
					t.Error("エラーが期待されましたが、エラーが発生しませんでした")
				}
			} else {
				if err != nil {
					t.Errorf("予期しないエラー: %v", err)
				}

				// 在庫が正しく減少したことを確認
				product, err := repo.GetByID(ctx, tt.productID)
				if err != nil {
					t.Errorf("商品の取得でエラー: %v", err)
				}
				if product.Stock != tt.expectedStock {
					t.Errorf("Stock = %v, want %v", product.Stock, tt.expectedStock)
				}
			}
		})
	}
}

func TestProductRepository_IncreaseStock(t *testing.T) {
	repo := NewProductRepository()
	ctx := context.Background()

	// テスト用の商品を作成
	testProduct := entity.NewProduct("在庫増加テスト商品", "説明", 1000, "image.jpg", 5)
	repo.Create(ctx, testProduct)

	tests := []struct {
		name        string
		productID   string
		quantity    int
		expectError bool
		expectedStock int
	}{
		{
			name:        "正常な在庫増加",
			productID:   testProduct.ID,
			quantity:    3,
			expectError: false,
			expectedStock: 8,
		},
		{
			name:        "存在しない商品",
			productID:   "nonexistent-id",
			quantity:    1,
			expectError: true,
			expectedStock: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.IncreaseStock(ctx, tt.productID, tt.quantity)

			if tt.expectError {
				if err != entity.ErrProductNotFound {
					t.Errorf("Expected ErrProductNotFound, got %v", err)
				}
			} else {
				if err != nil {
					t.Errorf("予期しないエラー: %v", err)
				}

				// 在庫が正しく増加したことを確認
				product, err := repo.GetByID(ctx, tt.productID)
				if err != nil {
					t.Errorf("商品の取得でエラー: %v", err)
				}
				if product.Stock != tt.expectedStock {
					t.Errorf("Stock = %v, want %v", product.Stock, tt.expectedStock)
				}
			}
		})
	}
}

// 並行処理の安全性をテスト
func TestProductRepository_ConcurrentAccess(t *testing.T) {
	repo := NewProductRepository()
	ctx := context.Background()

	// テスト用の商品を作成
	testProduct := entity.NewProduct("並行処理テスト商品", "説明", 1000, "image.jpg", 100)
	repo.Create(ctx, testProduct)

	const numGoroutines = 10
	const operationsPerGoroutine = 10

	var wg sync.WaitGroup
	wg.Add(numGoroutines * 2) // 読み取りと書き込みの両方

	// 並行読み取りテスト
	for i := 0; i < numGoroutines; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < operationsPerGoroutine; j++ {
				_, err := repo.GetByID(ctx, testProduct.ID)
				if err != nil {
					t.Errorf("並行読み取りでエラー: %v", err)
				}
			}
		}()
	}

	// 並行書き込みテスト（在庫の増減）
	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			defer wg.Done()
			for j := 0; j < operationsPerGoroutine; j++ {
				if id%2 == 0 {
					repo.IncreaseStock(ctx, testProduct.ID, 1)
				} else {
					repo.DecreaseStock(ctx, testProduct.ID, 1)
				}
			}
		}(i)
	}

	wg.Wait()

	// 最終的にデータが破損していないことを確認
	finalProduct, err := repo.GetByID(ctx, testProduct.ID)
	if err != nil {
		t.Errorf("最終確認でエラー: %v", err)
	}
	if finalProduct == nil {
		t.Error("商品が破損しています")
	}

	// 在庫が負の値になっていないことを確認
	if finalProduct.Stock < 0 {
		t.Errorf("在庫が負の値になっています: %v", finalProduct.Stock)
	}
}