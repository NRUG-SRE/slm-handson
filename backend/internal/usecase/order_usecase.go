package usecase

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/NRUG-SRE/slm-handson/backend/internal/domain/entity"
	"github.com/NRUG-SRE/slm-handson/backend/internal/domain/repository"
)

type OrderUseCase struct {
	orderRepo   repository.OrderRepository
	cartRepo    repository.CartRepository
	productRepo repository.ProductRepository
}

func NewOrderUseCase(
	orderRepo repository.OrderRepository,
	cartRepo repository.CartRepository,
	productRepo repository.ProductRepository,
) *OrderUseCase {
	return &OrderUseCase{
		orderRepo:   orderRepo,
		cartRepo:    cartRepo,
		productRepo: productRepo,
	}
}

func (uc *OrderUseCase) CreateOrder(ctx context.Context, cartID string) (*entity.Order, error) {
	// SLMデモ用のレスポンス時間調整
	uc.simulateResponseTime()

	// SLMデモ用のランダムエラー生成
	if uc.shouldSimulateError() {
		return nil, fmt.Errorf("simulated order processing error")
	}

	if cartID == "" {
		return nil, entity.ErrInvalidInput
	}

	// カートを取得
	cart, err := uc.cartRepo.GetByID(ctx, cartID)
	if err != nil {
		return nil, fmt.Errorf("failed to get cart: %w", err)
	}

	if cart.IsEmpty() {
		return nil, entity.ErrEmptyCart
	}

	// SLMハンズオン用に在庫チェックを無効化
	// 在庫確認は行わず、すべての商品が利用可能として処理

	// 注文を作成
	order, err := entity.NewOrder(cart)
	if err != nil {
		return nil, fmt.Errorf("failed to create order: %w", err)
	}

	// SLMハンズオン用に在庫減少処理を無効化
	// 在庫は減らさず、注文のみ作成

	// 注文を保存
	if err := uc.orderRepo.Create(ctx, order); err != nil {
		return nil, fmt.Errorf("failed to save order: %w", err)
	}

	// カートをクリア
	cart.Clear()
	if err := uc.cartRepo.Save(ctx, cart); err != nil {
		// ログに記録するが、注文は成功として扱う
		fmt.Printf("warning: failed to clear cart after order creation: %v\n", err)
	}

	// 決済処理のシミュレーション（非同期処理を模擬）
	go uc.processPayment(context.Background(), order)

	return order, nil
}

func (uc *OrderUseCase) GetOrder(ctx context.Context, orderID string) (*entity.Order, error) {
	if orderID == "" {
		return nil, entity.ErrInvalidInput
	}

	order, err := uc.orderRepo.GetByID(ctx, orderID)
	if err != nil {
		return nil, fmt.Errorf("failed to get order: %w", err)
	}

	return order, nil
}

func (uc *OrderUseCase) GetAllOrders(ctx context.Context) ([]*entity.Order, error) {
	orders, err := uc.orderRepo.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get orders: %w", err)
	}

	return orders, nil
}

// 決済処理のシミュレーション
func (uc *OrderUseCase) processPayment(ctx context.Context, order *entity.Order) {
	// 決済処理時間をシミュレート
	processingTime := 2 + rand.Intn(8) // 2-10秒
	time.Sleep(time.Duration(processingTime) * time.Second)

	// ランダムに決済成功/失敗を決定
	paymentSuccessRate := 0.9 // 90%の成功率
	if rand.Float64() < paymentSuccessRate {
		order.Complete()
	} else {
		order.Fail()
		// 在庫を戻す処理（実際の実装では必要）
		uc.restoreStock(ctx, order)
	}

	// 注文状態を更新
	if err := uc.orderRepo.Update(ctx, order); err != nil {
		fmt.Printf("error: failed to update order status: %v\n", err)
	}
}

func (uc *OrderUseCase) restoreStock(ctx context.Context, order *entity.Order) {
	for _, item := range order.Items {
		if err := uc.productRepo.IncreaseStock(ctx, item.ProductID, item.Quantity); err != nil {
			fmt.Printf("error: failed to restore stock for product %s: %v\n", item.ProductID, err)
		}
	}
}

// SLMデモ用のレスポンス時間シミュレーション
func (uc *OrderUseCase) simulateResponseTime() {
	minTime := uc.getEnvInt("RESPONSE_TIME_MIN", 50)
	maxTime := uc.getEnvInt("RESPONSE_TIME_MAX", 500)

	// 注文処理は通常より時間がかかる
	maxTime = maxTime * 2

	// 一定確率で遅いレスポンスを生成
	slowEndpointRate := uc.getEnvFloat("SLOW_ENDPOINT_RATE", 0.2)
	if rand.Float64() < slowEndpointRate {
		maxTime = maxTime * (2 + rand.Intn(2))
	}

	if maxTime > minTime {
		responseTime := minTime + rand.Intn(maxTime-minTime)
		time.Sleep(time.Duration(responseTime) * time.Millisecond)
	}
}

// SLMデモ用のランダムエラー生成
func (uc *OrderUseCase) shouldSimulateError() bool {
	errorRate := uc.getEnvFloat("ERROR_RATE", 0.1)
	// 注文処理はより高いエラー率
	return rand.Float64() < (errorRate * 1.5)
}

func (uc *OrderUseCase) getEnvInt(key string, defaultValue int) int {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}

	intValue, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}

	return intValue
}

func (uc *OrderUseCase) getEnvFloat(key string, defaultValue float64) float64 {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}

	floatValue, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return defaultValue
	}

	return floatValue
}
