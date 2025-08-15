package memory

import (
	"context"
	"sync"

	"github.com/NRUG-SRE/slm-handson/backend/internal/domain/entity"
	"github.com/NRUG-SRE/slm-handson/backend/internal/domain/repository"
)

type cartRepository struct {
	carts map[string]*entity.Cart
	mutex sync.RWMutex
}

func NewCartRepository() repository.CartRepository {
	return &cartRepository{
		carts: make(map[string]*entity.Cart),
		mutex: sync.RWMutex{},
	}
}

func (r *cartRepository) GetByID(ctx context.Context, id string) (*entity.Cart, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	cart, exists := r.carts[id]
	if !exists {
		return nil, entity.ErrItemNotFound
	}

	return cart, nil
}

func (r *cartRepository) GetOrCreate(ctx context.Context, id string) (*entity.Cart, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	cart, exists := r.carts[id]
	if !exists {
		cart = entity.NewCart()
		cart.ID = id // 指定されたIDを使用
		r.carts[id] = cart
	}

	return cart, nil
}

func (r *cartRepository) Save(ctx context.Context, cart *entity.Cart) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	r.carts[cart.ID] = cart
	return nil
}

func (r *cartRepository) Delete(ctx context.Context, id string) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, exists := r.carts[id]; !exists {
		return entity.ErrItemNotFound
	}

	delete(r.carts, id)
	return nil
}

func (r *cartRepository) Clear(ctx context.Context, id string) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	cart, exists := r.carts[id]
	if !exists {
		return entity.ErrItemNotFound
	}

	cart.Clear()
	return nil
}
