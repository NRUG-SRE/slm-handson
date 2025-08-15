package mocks

import (
	"context"

	"github.com/NRUG-SRE/slm-handson/backend/internal/domain/entity"
)

// MockCartRepository はCartRepositoryのモック実装
type MockCartRepository struct {
	GetByIDFunc     func(ctx context.Context, id string) (*entity.Cart, error)
	GetOrCreateFunc func(ctx context.Context, id string) (*entity.Cart, error)
	SaveFunc        func(ctx context.Context, cart *entity.Cart) error
	DeleteFunc      func(ctx context.Context, id string) error
	ClearFunc       func(ctx context.Context, id string) error

	// 呼び出し記録用
	GetByIDCalls []struct {
		Ctx context.Context
		ID  string
	}
	GetOrCreateCalls []struct {
		Ctx context.Context
		ID  string
	}
	SaveCalls []struct {
		Ctx  context.Context
		Cart *entity.Cart
	}
	DeleteCalls []struct {
		Ctx context.Context
		ID  string
	}
	ClearCalls []struct {
		Ctx context.Context
		ID  string
	}
}

func (m *MockCartRepository) GetByID(ctx context.Context, id string) (*entity.Cart, error) {
	m.GetByIDCalls = append(m.GetByIDCalls, struct {
		Ctx context.Context
		ID  string
	}{ctx, id})
	if m.GetByIDFunc != nil {
		return m.GetByIDFunc(ctx, id)
	}
	return nil, nil
}

func (m *MockCartRepository) GetOrCreate(ctx context.Context, id string) (*entity.Cart, error) {
	m.GetOrCreateCalls = append(m.GetOrCreateCalls, struct {
		Ctx context.Context
		ID  string
	}{ctx, id})
	if m.GetOrCreateFunc != nil {
		return m.GetOrCreateFunc(ctx, id)
	}
	return nil, nil
}

func (m *MockCartRepository) Save(ctx context.Context, cart *entity.Cart) error {
	m.SaveCalls = append(m.SaveCalls, struct {
		Ctx  context.Context
		Cart *entity.Cart
	}{ctx, cart})
	if m.SaveFunc != nil {
		return m.SaveFunc(ctx, cart)
	}
	return nil
}

func (m *MockCartRepository) Delete(ctx context.Context, id string) error {
	m.DeleteCalls = append(m.DeleteCalls, struct {
		Ctx context.Context
		ID  string
	}{ctx, id})
	if m.DeleteFunc != nil {
		return m.DeleteFunc(ctx, id)
	}
	return nil
}

func (m *MockCartRepository) Clear(ctx context.Context, id string) error {
	m.ClearCalls = append(m.ClearCalls, struct {
		Ctx context.Context
		ID  string
	}{ctx, id})
	if m.ClearFunc != nil {
		return m.ClearFunc(ctx, id)
	}
	return nil
}
