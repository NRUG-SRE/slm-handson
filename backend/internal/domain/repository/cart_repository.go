package repository

import (
	"context"

	"github.com/NRUG-SRE/slm-handson/backend/internal/domain/entity"
)

type CartRepository interface {
	GetByID(ctx context.Context, id string) (*entity.Cart, error)
	GetOrCreate(ctx context.Context, id string) (*entity.Cart, error)
	Save(ctx context.Context, cart *entity.Cart) error
	Delete(ctx context.Context, id string) error
	Clear(ctx context.Context, id string) error
}
