package memory

import (
	"context"
	"sync"
	"testing"

	"github.com/NRUG-SRE/slm-handson/backend/internal/domain/entity"
)

func TestCartRepository_GetOrCreate(t *testing.T) {
	repo := NewCartRepository()
	ctx := context.Background()

	tests := []struct {
		name   string
		cartID string
	}{
		{
			name:   "新しいカートを作成",
			cartID: "cart-123",
		},
		{
			name:   "既存のカートを取得",
			cartID: "cart-123", // 同じIDを再度使用
		},
	}

	var firstCart *entity.Cart

	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cart, err := repo.GetOrCreate(ctx, tt.cartID)

			if err != nil {
				t.Errorf("予期しないエラー: %v", err)
			}
			if cart == nil {
				t.Fatal("カートがnilです")
			}
			if cart.ID != tt.cartID {
				t.Errorf("ID = %v, want %v", cart.ID, tt.cartID)
			}

			if i == 0 {
				// 最初の呼び出し - 新しいカートが作成される
				firstCart = cart
				if !cart.IsEmpty() {
					t.Error("新しいカートは空であるべきです")
				}
			} else {
				// 2回目の呼び出し - 同じカートインスタンスが返される
				if cart != firstCart {
					t.Error("同じカートインスタンスが返されるべきです")
				}
			}
		})
	}
}

func TestCartRepository_GetByID(t *testing.T) {
	repo := NewCartRepository()
	ctx := context.Background()

	// テスト用のカートを作成
	testCart, _ := repo.GetOrCreate(ctx, "cart-test")

	tests := []struct {
		name        string
		cartID      string
		expectError bool
		checkResult func(t *testing.T, cart *entity.Cart)
	}{
		{
			name:        "存在するカートを取得",
			cartID:      "cart-test",
			expectError: false,
			checkResult: func(t *testing.T, cart *entity.Cart) {
				if cart == nil {
					t.Fatal("カートがnilです")
				}
				if cart.ID != "cart-test" {
					t.Errorf("ID = %v, want cart-test", cart.ID)
				}
				if cart != testCart {
					t.Error("同じカートインスタンスが返されるべきです")
				}
			},
		},
		{
			name:        "存在しないカートを取得",
			cartID:      "nonexistent-cart",
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
			cart, err := repo.GetByID(ctx, tt.cartID)

			if tt.expectError {
				if err != entity.ErrItemNotFound {
					t.Errorf("Expected ErrItemNotFound, got %v", err)
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

func TestCartRepository_Save(t *testing.T) {
	repo := NewCartRepository()
	ctx := context.Background()

	// テスト用のカートを作成
	cart := entity.NewCart()
	cart.ID = "cart-save-test"

	// テスト用の商品を追加
	product := entity.NewProduct("テスト商品", "説明", 1000, "image.jpg", 10)
	cart.AddItem(product, 2)

	// カートを保存
	err := repo.Save(ctx, cart)
	if err != nil {
		t.Errorf("予期しないエラー: %v", err)
	}

	// 保存されたカートを取得して確認
	savedCart, err := repo.GetByID(ctx, cart.ID)
	if err != nil {
		t.Errorf("保存されたカートの取得でエラー: %v", err)
	}
	if savedCart == nil {
		t.Fatal("保存されたカートがnilです")
	}
	if savedCart.ID != cart.ID {
		t.Errorf("ID = %v, want %v", savedCart.ID, cart.ID)
	}
	if len(savedCart.Items) != 1 {
		t.Errorf("アイテム数 = %v, want 1", len(savedCart.Items))
	}
	if savedCart.TotalAmount != 2000 {
		t.Errorf("合計金額 = %v, want 2000", savedCart.TotalAmount)
	}
}

func TestCartRepository_Delete(t *testing.T) {
	repo := NewCartRepository()
	ctx := context.Background()

	// テスト用のカートを作成
	testCart, _ := repo.GetOrCreate(ctx, "cart-delete-test")

	tests := []struct {
		name        string
		cartID      string
		expectError bool
	}{
		{
			name:        "存在するカートを削除",
			cartID:      testCart.ID,
			expectError: false,
		},
		{
			name:        "存在しないカートを削除",
			cartID:      "nonexistent-cart",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.Delete(ctx, tt.cartID)

			if tt.expectError {
				if err != entity.ErrItemNotFound {
					t.Errorf("Expected ErrItemNotFound, got %v", err)
				}
			} else {
				if err != nil {
					t.Errorf("予期しないエラー: %v", err)
				}

				// 削除されたカートが取得できないことを確認
				_, err := repo.GetByID(ctx, tt.cartID)
				if err != entity.ErrItemNotFound {
					t.Errorf("削除されたカートが取得できてしまいました")
				}
			}
		})
	}
}

func TestCartRepository_Clear(t *testing.T) {
	repo := NewCartRepository()
	ctx := context.Background()

	// テスト用のカートを作成
	testCart, _ := repo.GetOrCreate(ctx, "cart-clear-test")

	// カートに商品を追加
	product1 := entity.NewProduct("商品1", "説明1", 1000, "image1.jpg", 10)
	product2 := entity.NewProduct("商品2", "説明2", 2000, "image2.jpg", 10)
	testCart.AddItem(product1, 2)
	testCart.AddItem(product2, 1)
	repo.Save(ctx, testCart)

	tests := []struct {
		name        string
		cartID      string
		expectError bool
	}{
		{
			name:        "存在するカートをクリア",
			cartID:      testCart.ID,
			expectError: false,
		},
		{
			name:        "存在しないカートをクリア",
			cartID:      "nonexistent-cart",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.Clear(ctx, tt.cartID)

			if tt.expectError {
				if err != entity.ErrItemNotFound {
					t.Errorf("Expected ErrItemNotFound, got %v", err)
				}
			} else {
				if err != nil {
					t.Errorf("予期しないエラー: %v", err)
				}

				// クリアされたカートを取得して確認
				clearedCart, err := repo.GetByID(ctx, tt.cartID)
				if err != nil {
					t.Errorf("クリアされたカートの取得でエラー: %v", err)
				}
				if !clearedCart.IsEmpty() {
					t.Error("カートが空になっていません")
				}
				if clearedCart.TotalAmount != 0 {
					t.Errorf("合計金額 = %v, want 0", clearedCart.TotalAmount)
				}
			}
		})
	}
}

// 並行処理の安全性をテスト
func TestCartRepository_ConcurrentAccess(t *testing.T) {
	repo := NewCartRepository()
	ctx := context.Background()

	const numGoroutines = 10
	const operationsPerGoroutine = 10

	var wg sync.WaitGroup
	wg.Add(numGoroutines * 3) // GetOrCreate, Save, GetByIDの3つの操作

	// 並行でGetOrCreateを実行
	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			defer wg.Done()
			for j := 0; j < operationsPerGoroutine; j++ {
				cartID := "concurrent-cart"
				cart, err := repo.GetOrCreate(ctx, cartID)
				if err != nil {
					t.Errorf("並行GetOrCreateでエラー: %v", err)
				}
				if cart == nil {
					t.Error("カートがnilです")
				}
				if cart.ID != cartID {
					t.Errorf("ID = %v, want %v", cart.ID, cartID)
				}
			}
		}(i)
	}

	// 並行でSaveを実行
	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			defer wg.Done()
			for j := 0; j < operationsPerGoroutine; j++ {
				cart, _ := repo.GetOrCreate(ctx, "concurrent-cart")
				// カートに商品を追加してから保存
				product := entity.NewProduct("商品", "説明", 100, "image.jpg", 10)
				cart.AddItem(product, 1)
				err := repo.Save(ctx, cart)
				if err != nil {
					t.Errorf("並行Saveでエラー: %v", err)
				}
			}
		}(i)
	}

	// 並行でGetByIDを実行
	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			defer wg.Done()
			for j := 0; j < operationsPerGoroutine; j++ {
				// まずカートが存在することを確認
				repo.GetOrCreate(ctx, "concurrent-cart")

				_, err := repo.GetByID(ctx, "concurrent-cart")
				if err != nil {
					t.Errorf("並行GetByIDでエラー: %v", err)
				}
			}
		}(i)
	}

	wg.Wait()

	// 最終的にカートが破損していないことを確認
	finalCart, err := repo.GetByID(ctx, "concurrent-cart")
	if err != nil {
		t.Errorf("最終確認でエラー: %v", err)
	}
	if finalCart == nil {
		t.Error("カートが破損しています")
	}
	if finalCart.ID != "concurrent-cart" {
		t.Error("カートのIDが破損しています")
	}
}

// メモリリークのテスト（大量のカート作成・削除）
func TestCartRepository_MemoryManagement(t *testing.T) {
	repo := NewCartRepository()
	ctx := context.Background()

	const numCarts = 1000

	// 大量のカートを作成
	for i := 0; i < numCarts; i++ {
		cartID := "memory-test-cart-" + string(rune(i))
		cart, err := repo.GetOrCreate(ctx, cartID)
		if err != nil {
			t.Errorf("カート作成でエラー: %v", err)
		}
		if cart == nil {
			t.Error("カートがnilです")
		}

		// 商品を追加
		product := entity.NewProduct("商品", "説明", 1000, "image.jpg", 10)
		cart.AddItem(product, 1)
		repo.Save(ctx, cart)
	}

	// 作成したカートの半分を削除
	for i := 0; i < numCarts/2; i++ {
		cartID := "memory-test-cart-" + string(rune(i))
		err := repo.Delete(ctx, cartID)
		if err != nil {
			t.Errorf("カート削除でエラー: %v", err)
		}
	}

	// 残りのカートが正常にアクセスできることを確認
	for i := numCarts / 2; i < numCarts; i++ {
		cartID := "memory-test-cart-" + string(rune(i))
		cart, err := repo.GetByID(ctx, cartID)
		if err != nil {
			t.Errorf("残存カートの取得でエラー: %v", err)
		}
		if cart == nil {
			t.Error("残存カートがnilです")
		}
		if len(cart.Items) != 1 {
			t.Errorf("アイテム数 = %v, want 1", len(cart.Items))
		}
	}

	// 削除されたカートがアクセスできないことを確認
	for i := 0; i < numCarts/2; i++ {
		cartID := "memory-test-cart-" + string(rune(i))
		_, err := repo.GetByID(ctx, cartID)
		if err != entity.ErrItemNotFound {
			t.Errorf("削除されたカートがアクセスできてしまいました: %v", cartID)
		}
	}
}
