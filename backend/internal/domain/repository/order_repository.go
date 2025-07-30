package repository

import (
	"context"

	"github.com/NRUG-SRE/slm-handson/backend/internal/domain/entity"
)

type OrderRepository interface {
	GetAll(ctx context.Context) ([]*entity.Order, error)
	GetByID(ctx context.Context, id string) (*entity.Order, error)
	Create(ctx context.Context, order *entity.Order) error
	Update(ctx context.Context, order *entity.Order) error
	Delete(ctx context.Context, id string) error
}