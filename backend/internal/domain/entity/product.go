package entity

import (
	"time"

	"github.com/google/uuid"
)

type Product struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       int       `json:"price"` // 価格は円単位で整数で管理
	ImageURL    string    `json:"imageUrl"`
	Stock       int       `json:"stock"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

func NewProduct(name, description string, price int, imageURL string, stock int) *Product {
	now := time.Now()
	return &Product{
		ID:          uuid.New().String(),
		Name:        name,
		Description: description,
		Price:       price,
		ImageURL:    imageURL,
		Stock:       stock,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

func (p *Product) UpdateStock(newStock int) {
	p.Stock = newStock
	p.UpdatedAt = time.Now()
}

func (p *Product) DecreaseStock(quantity int) error {
	if p.Stock < quantity {
		return ErrInsufficientStock
	}
	p.Stock -= quantity
	p.UpdatedAt = time.Now()
	return nil
}

func (p *Product) IncreaseStock(quantity int) {
	p.Stock += quantity
	p.UpdatedAt = time.Now()
}

func (p *Product) IsInStock() bool {
	return p.Stock > 0
}

func (p *Product) IsAvailable(quantity int) bool {
	return p.Stock >= quantity
}