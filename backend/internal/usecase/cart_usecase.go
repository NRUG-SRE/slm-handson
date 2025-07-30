package usecase

import (
	"context"
	"fmt"

	"github.com/NRUG-SRE/slm-handson/backend/internal/domain/entity"
	"github.com/NRUG-SRE/slm-handson/backend/internal/domain/repository"
)

type CartUseCase struct {
	cartRepo    repository.CartRepository
	productRepo repository.ProductRepository
}

func NewCartUseCase(cartRepo repository.CartRepository, productRepo repository.ProductRepository) *CartUseCase {
	return &CartUseCase{
		cartRepo:    cartRepo,
		productRepo: productRepo,
	}
}

func (uc *CartUseCase) GetCart(ctx context.Context, cartID string) (*entity.Cart, error) {
	if cartID == "" {
		return nil, entity.ErrInvalidInput
	}

	cart, err := uc.cartRepo.GetOrCreate(ctx, cartID)
	if err != nil {
		return nil, fmt.Errorf("failed to get cart: %w", err)
	}

	return cart, nil
}

func (uc *CartUseCase) AddToCart(ctx context.Context, cartID, productID string, quantity int) (*entity.Cart, error) {
	if cartID == "" || productID == "" || quantity <= 0 {
		return nil, entity.ErrInvalidInput
	}

	// カートを取得または作成
	cart, err := uc.cartRepo.GetOrCreate(ctx, cartID)
	if err != nil {
		return nil, fmt.Errorf("failed to get cart: %w", err)
	}

	// 商品情報を取得
	product, err := uc.productRepo.GetByID(ctx, productID)
	if err != nil {
		return nil, fmt.Errorf("failed to get product: %w", err)
	}

	// カートに商品を追加
	if err := cart.AddItem(product, quantity); err != nil {
		return nil, fmt.Errorf("failed to add item to cart: %w", err)
	}

	// カートを保存
	if err := uc.cartRepo.Save(ctx, cart); err != nil {
		return nil, fmt.Errorf("failed to save cart: %w", err)
	}

	return cart, nil
}

func (uc *CartUseCase) UpdateCartItem(ctx context.Context, cartID, itemID string, quantity int) (*entity.Cart, error) {
	if cartID == "" || itemID == "" {
		return nil, entity.ErrInvalidInput
	}

	// カートを取得
	cart, err := uc.cartRepo.GetByID(ctx, cartID)
	if err != nil {
		return nil, fmt.Errorf("failed to get cart: %w", err)
	}

	// カートアイテムの数量を更新
	if err := cart.UpdateItemQuantity(itemID, quantity); err != nil {
		return nil, fmt.Errorf("failed to update cart item: %w", err)
	}

	// カートを保存
	if err := uc.cartRepo.Save(ctx, cart); err != nil {
		return nil, fmt.Errorf("failed to save cart: %w", err)
	}

	return cart, nil
}

func (uc *CartUseCase) RemoveFromCart(ctx context.Context, cartID, itemID string) (*entity.Cart, error) {
	if cartID == "" || itemID == "" {
		return nil, entity.ErrInvalidInput
	}

	// カートを取得
	cart, err := uc.cartRepo.GetByID(ctx, cartID)
	if err != nil {
		return nil, fmt.Errorf("failed to get cart: %w", err)
	}

	// カートからアイテムを削除
	if err := cart.RemoveItem(itemID); err != nil {
		return nil, fmt.Errorf("failed to remove cart item: %w", err)
	}

	// カートを保存
	if err := uc.cartRepo.Save(ctx, cart); err != nil {
		return nil, fmt.Errorf("failed to save cart: %w", err)
	}

	return cart, nil
}

func (uc *CartUseCase) ClearCart(ctx context.Context, cartID string) error {
	if cartID == "" {
		return entity.ErrInvalidInput
	}

	if err := uc.cartRepo.Clear(ctx, cartID); err != nil {
		return fmt.Errorf("failed to clear cart: %w", err)
	}

	return nil
}