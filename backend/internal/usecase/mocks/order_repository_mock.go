package mocks

import (
	"context"

	"github.com/NRUG-SRE/slm-handson/backend/internal/domain/entity"
)

// MockOrderRepository はOrderRepositoryのモック実装
type MockOrderRepository struct {
	GetAllFunc  func(ctx context.Context) ([]*entity.Order, error)
	GetByIDFunc func(ctx context.Context, id string) (*entity.Order, error)
	CreateFunc  func(ctx context.Context, order *entity.Order) error
	UpdateFunc  func(ctx context.Context, order *entity.Order) error
	DeleteFunc  func(ctx context.Context, id string) error

	// 呼び出し記録用
	GetAllCalls  []context.Context
	GetByIDCalls []struct{ Ctx context.Context; ID string }
	CreateCalls  []struct{ Ctx context.Context; Order *entity.Order }
	UpdateCalls  []struct{ Ctx context.Context; Order *entity.Order }
	DeleteCalls  []struct{ Ctx context.Context; ID string }
}

func (m *MockOrderRepository) GetAll(ctx context.Context) ([]*entity.Order, error) {
	m.GetAllCalls = append(m.GetAllCalls, ctx)
	if m.GetAllFunc != nil {
		return m.GetAllFunc(ctx)
	}
	return nil, nil
}

func (m *MockOrderRepository) GetByID(ctx context.Context, id string) (*entity.Order, error) {
	m.GetByIDCalls = append(m.GetByIDCalls, struct{ Ctx context.Context; ID string }{ctx, id})
	if m.GetByIDFunc != nil {
		return m.GetByIDFunc(ctx, id)
	}
	return nil, nil
}

func (m *MockOrderRepository) Create(ctx context.Context, order *entity.Order) error {
	m.CreateCalls = append(m.CreateCalls, struct{ Ctx context.Context; Order *entity.Order }{ctx, order})
	if m.CreateFunc != nil {
		return m.CreateFunc(ctx, order)
	}
	return nil
}

func (m *MockOrderRepository) Update(ctx context.Context, order *entity.Order) error {
	m.UpdateCalls = append(m.UpdateCalls, struct{ Ctx context.Context; Order *entity.Order }{ctx, order})
	if m.UpdateFunc != nil {
		return m.UpdateFunc(ctx, order)
	}
	return nil
}

func (m *MockOrderRepository) Delete(ctx context.Context, id string) error {
	m.DeleteCalls = append(m.DeleteCalls, struct{ Ctx context.Context; ID string }{ctx, id})
	if m.DeleteFunc != nil {
		return m.DeleteFunc(ctx, id)
	}
	return nil
}