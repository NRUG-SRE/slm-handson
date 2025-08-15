package repository

import (
	"context"

	"github.com/NRUG-SRE/slm-handson/backend/internal/domain/entity"
)

type ProductRepository interface {
	GetAll(ctx context.Context) ([]*entity.Product, error)
	GetByID(ctx context.Context, id string) (*entity.Product, error)
	Create(ctx context.Context, product *entity.Product) error
	Update(ctx context.Context, product *entity.Product) error
	Delete(ctx context.Context, id string) error
	UpdateStock(ctx context.Context, id string, newStock int) error
	DecreaseStock(ctx context.Context, id string, quantity int) error
	IncreaseStock(ctx context.Context, id string, quantity int) error
}
