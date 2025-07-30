package memory

import (
	"context"
	"sync"

	"github.com/NRUG-SRE/slm-handson/backend/internal/domain/entity"
	"github.com/NRUG-SRE/slm-handson/backend/internal/domain/repository"
)

type orderRepository struct {
	orders map[string]*entity.Order
	mutex  sync.RWMutex
}

func NewOrderRepository() repository.OrderRepository {
	return &orderRepository{
		orders: make(map[string]*entity.Order),
		mutex:  sync.RWMutex{},
	}
}

func (r *orderRepository) GetAll(ctx context.Context) ([]*entity.Order, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	orders := make([]*entity.Order, 0, len(r.orders))
	for _, order := range r.orders {
		orders = append(orders, order)
	}

	return orders, nil
}

func (r *orderRepository) GetByID(ctx context.Context, id string) (*entity.Order, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	order, exists := r.orders[id]
	if !exists {
		return nil, entity.ErrOrderNotFound
	}

	return order, nil
}

func (r *orderRepository) Create(ctx context.Context, order *entity.Order) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	r.orders[order.ID] = order
	return nil
}

func (r *orderRepository) Update(ctx context.Context, order *entity.Order) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, exists := r.orders[order.ID]; !exists {
		return entity.ErrOrderNotFound
	}

	r.orders[order.ID] = order
	return nil
}

func (r *orderRepository) Delete(ctx context.Context, id string) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, exists := r.orders[id]; !exists {
		return entity.ErrOrderNotFound
	}

	delete(r.orders, id)
	return nil
}