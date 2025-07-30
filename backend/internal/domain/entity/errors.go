package entity

import "errors"

// ドメインエラー定義
var (
	// Product関連エラー
	ErrProductNotFound    = errors.New("product not found")
	ErrInsufficientStock  = errors.New("insufficient stock")

	// Cart関連エラー
	ErrItemNotFound = errors.New("item not found in cart")
	ErrEmptyCart    = errors.New("cart is empty")

	// Order関連エラー
	ErrOrderNotFound      = errors.New("order not found")
	ErrInvalidOrderStatus = errors.New("invalid order status transition")

	// 一般的なエラー
	ErrInvalidInput = errors.New("invalid input")
)