package mocks

import (
	"context"

	"github.com/NRUG-SRE/slm-handson/backend/internal/domain/entity"
)

// MockProductRepository はProductRepositoryのモック実装
type MockProductRepository struct {
	GetAllFunc         func(ctx context.Context) ([]*entity.Product, error)
	GetByIDFunc        func(ctx context.Context, id string) (*entity.Product, error)
	CreateFunc         func(ctx context.Context, product *entity.Product) error
	UpdateFunc         func(ctx context.Context, product *entity.Product) error
	DeleteFunc         func(ctx context.Context, id string) error
	UpdateStockFunc    func(ctx context.Context, id string, newStock int) error
	DecreaseStockFunc  func(ctx context.Context, id string, quantity int) error
	IncreaseStockFunc  func(ctx context.Context, id string, quantity int) error

	// 呼び出し記録用
	GetAllCalls         []context.Context
	GetByIDCalls        []struct{ Ctx context.Context; ID string }
	CreateCalls         []struct{ Ctx context.Context; Product *entity.Product }
	UpdateCalls         []struct{ Ctx context.Context; Product *entity.Product }
	DeleteCalls         []struct{ Ctx context.Context; ID string }
	UpdateStockCalls    []struct{ Ctx context.Context; ID string; NewStock int }
	DecreaseStockCalls  []struct{ Ctx context.Context; ID string; Quantity int }
	IncreaseStockCalls  []struct{ Ctx context.Context; ID string; Quantity int }
}

func (m *MockProductRepository) GetAll(ctx context.Context) ([]*entity.Product, error) {
	m.GetAllCalls = append(m.GetAllCalls, ctx)
	if m.GetAllFunc != nil {
		return m.GetAllFunc(ctx)
	}
	return nil, nil
}

func (m *MockProductRepository) GetByID(ctx context.Context, id string) (*entity.Product, error) {
	m.GetByIDCalls = append(m.GetByIDCalls, struct{ Ctx context.Context; ID string }{ctx, id})
	if m.GetByIDFunc != nil {
		return m.GetByIDFunc(ctx, id)
	}
	return nil, nil
}

func (m *MockProductRepository) Create(ctx context.Context, product *entity.Product) error {
	m.CreateCalls = append(m.CreateCalls, struct{ Ctx context.Context; Product *entity.Product }{ctx, product})
	if m.CreateFunc != nil {
		return m.CreateFunc(ctx, product)
	}
	return nil
}

func (m *MockProductRepository) Update(ctx context.Context, product *entity.Product) error {
	m.UpdateCalls = append(m.UpdateCalls, struct{ Ctx context.Context; Product *entity.Product }{ctx, product})
	if m.UpdateFunc != nil {
		return m.UpdateFunc(ctx, product)
	}
	return nil
}

func (m *MockProductRepository) Delete(ctx context.Context, id string) error {
	m.DeleteCalls = append(m.DeleteCalls, struct{ Ctx context.Context; ID string }{ctx, id})
	if m.DeleteFunc != nil {
		return m.DeleteFunc(ctx, id)
	}
	return nil
}

func (m *MockProductRepository) UpdateStock(ctx context.Context, id string, newStock int) error {
	m.UpdateStockCalls = append(m.UpdateStockCalls, struct{ Ctx context.Context; ID string; NewStock int }{ctx, id, newStock})
	if m.UpdateStockFunc != nil {
		return m.UpdateStockFunc(ctx, id, newStock)
	}
	return nil
}

func (m *MockProductRepository) DecreaseStock(ctx context.Context, id string, quantity int) error {
	m.DecreaseStockCalls = append(m.DecreaseStockCalls, struct{ Ctx context.Context; ID string; Quantity int }{ctx, id, quantity})
	if m.DecreaseStockFunc != nil {
		return m.DecreaseStockFunc(ctx, id, quantity)
	}
	return nil
}

func (m *MockProductRepository) IncreaseStock(ctx context.Context, id string, quantity int) error {
	m.IncreaseStockCalls = append(m.IncreaseStockCalls, struct{ Ctx context.Context; ID string; Quantity int }{ctx, id, quantity})
	if m.IncreaseStockFunc != nil {
		return m.IncreaseStockFunc(ctx, id, quantity)
	}
	return nil
}