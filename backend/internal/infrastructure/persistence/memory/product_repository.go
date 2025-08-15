package memory

import (
	"context"
	"sync"

	"github.com/NRUG-SRE/slm-handson/backend/internal/domain/entity"
	"github.com/NRUG-SRE/slm-handson/backend/internal/domain/repository"
)

type productRepository struct {
	products map[string]*entity.Product
	mutex    sync.RWMutex
}

func NewProductRepository() repository.ProductRepository {
	repo := &productRepository{
		products: make(map[string]*entity.Product),
		mutex:    sync.RWMutex{},
	}

	// 初期データを投入
	repo.seedData()

	return repo
}

func (r *productRepository) seedData() {
	products := []*entity.Product{
		entity.NewProduct(
			"ワイヤレスヘッドホン",
			"高音質なノイズキャンセリング機能付きワイヤレスヘッドホン",
			25000,
			"/images/headphones.svg",
			10,
		),
		entity.NewProduct(
			"スマートウォッチ",
			"フィットネストラッキング機能付きの最新スマートウォッチ",
			35000,
			"/images/smartwatch.svg",
			5,
		),
		entity.NewProduct(
			"ポータブルスピーカー",
			"防水機能付きの高音質Bluetoothスピーカー",
			12000,
			"/images/speaker.svg",
			15,
		),
		entity.NewProduct(
			"ワイヤレスキーボード",
			"人間工学に基づいたデザインのワイヤレスキーボード",
			8500,
			"/images/keyboard.svg",
			20,
		),
		entity.NewProduct(
			"4K Webカメラ",
			"リモートワークに最適な高画質Webカメラ",
			15000,
			"/images/webcam.svg",
			8,
		),
		entity.NewProduct(
			"USB-C ハブ",
			"7つのポートを備えた多機能USB-Cハブ",
			6500,
			"/images/usb-hub.svg",
			0,
		),
	}

	for _, product := range products {
		r.products[product.ID] = product
	}
}

func (r *productRepository) GetAll(ctx context.Context) ([]*entity.Product, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	products := make([]*entity.Product, 0, len(r.products))
	for _, product := range r.products {
		products = append(products, product)
	}

	return products, nil
}

func (r *productRepository) GetByID(ctx context.Context, id string) (*entity.Product, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	product, exists := r.products[id]
	if !exists {
		return nil, entity.ErrProductNotFound
	}

	return product, nil
}

func (r *productRepository) Create(ctx context.Context, product *entity.Product) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	r.products[product.ID] = product
	return nil
}

func (r *productRepository) Update(ctx context.Context, product *entity.Product) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, exists := r.products[product.ID]; !exists {
		return entity.ErrProductNotFound
	}

	r.products[product.ID] = product
	return nil
}

func (r *productRepository) Delete(ctx context.Context, id string) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, exists := r.products[id]; !exists {
		return entity.ErrProductNotFound
	}

	delete(r.products, id)
	return nil
}

func (r *productRepository) UpdateStock(ctx context.Context, id string, newStock int) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	product, exists := r.products[id]
	if !exists {
		return entity.ErrProductNotFound
	}

	product.UpdateStock(newStock)
	return nil
}

func (r *productRepository) DecreaseStock(ctx context.Context, id string, quantity int) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	product, exists := r.products[id]
	if !exists {
		return entity.ErrProductNotFound
	}

	return product.DecreaseStock(quantity)
}

func (r *productRepository) IncreaseStock(ctx context.Context, id string, quantity int) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	product, exists := r.products[id]
	if !exists {
		return entity.ErrProductNotFound
	}

	product.IncreaseStock(quantity)
	return nil
}
