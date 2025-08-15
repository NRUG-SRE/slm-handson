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

type ProductUseCase struct {
	productRepo repository.ProductRepository
}

func NewProductUseCase(productRepo repository.ProductRepository) *ProductUseCase {
	return &ProductUseCase{
		productRepo: productRepo,
	}
}

func (uc *ProductUseCase) GetAllProducts(ctx context.Context) ([]*entity.Product, error) {
	// SLMデモ用のレスポンス時間調整
	uc.simulateResponseTime()

	// SLMデモ用のランダムエラー生成
	if uc.shouldSimulateError() {
		return nil, fmt.Errorf("simulated product service error")
	}

	products, err := uc.productRepo.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get products: %w", err)
	}

	return products, nil
}

func (uc *ProductUseCase) GetProductByID(ctx context.Context, id string) (*entity.Product, error) {
	// SLMデモ用のレスポンス時間調整
	uc.simulateResponseTime()

	// SLMデモ用のランダムエラー生成
	if uc.shouldSimulateError() {
		return nil, fmt.Errorf("simulated product service error")
	}

	if id == "" {
		return nil, entity.ErrInvalidInput
	}

	product, err := uc.productRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get product: %w", err)
	}

	return product, nil
}

func (uc *ProductUseCase) CreateProduct(ctx context.Context, name, description string, price int, imageURL string, stock int) (*entity.Product, error) {
	// バリデーション
	if name == "" || description == "" || price <= 0 || stock < 0 {
		return nil, entity.ErrInvalidInput
	}

	product := entity.NewProduct(name, description, price, imageURL, stock)

	if err := uc.productRepo.Create(ctx, product); err != nil {
		return nil, fmt.Errorf("failed to create product: %w", err)
	}

	return product, nil
}

func (uc *ProductUseCase) UpdateProduct(ctx context.Context, product *entity.Product) error {
	// バリデーション
	if product == nil || product.ID == "" {
		return entity.ErrInvalidInput
	}

	if err := uc.productRepo.Update(ctx, product); err != nil {
		return fmt.Errorf("failed to update product: %w", err)
	}

	return nil
}

func (uc *ProductUseCase) DeleteProduct(ctx context.Context, id string) error {
	if id == "" {
		return entity.ErrInvalidInput
	}

	if err := uc.productRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete product: %w", err)
	}

	return nil
}

// SLMデモ用のレスポンス時間シミュレーション
func (uc *ProductUseCase) simulateResponseTime() {
	minTime := uc.getEnvInt("RESPONSE_TIME_MIN", 50)
	maxTime := uc.getEnvInt("RESPONSE_TIME_MAX", 500)

	// 一定確率で遅いレスポンスを生成
	slowEndpointRate := uc.getEnvFloat("SLOW_ENDPOINT_RATE", 0.2)
	if rand.Float64() < slowEndpointRate {
		// 遅いエンドポイントの場合は最大時間の2-3倍にする
		maxTime = maxTime * (2 + rand.Intn(2))
	}

	if maxTime > minTime {
		responseTime := minTime + rand.Intn(maxTime-minTime)
		time.Sleep(time.Duration(responseTime) * time.Millisecond)
	}
}

// SLMデモ用のランダムエラー生成
func (uc *ProductUseCase) shouldSimulateError() bool {
	errorRate := uc.getEnvFloat("ERROR_RATE", 0.1)
	return rand.Float64() < errorRate
}

func (uc *ProductUseCase) getEnvInt(key string, defaultValue int) int {
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

func (uc *ProductUseCase) getEnvFloat(key string, defaultValue float64) float64 {
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
