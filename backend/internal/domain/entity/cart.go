package entity

import (
	"time"

	"github.com/google/uuid"
)

type CartItem struct {
	ID        string    `json:"id"`
	ProductID string    `json:"productId"`
	Product   *Product  `json:"product"`
	Quantity  int       `json:"quantity"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type Cart struct {
	ID          string      `json:"id"`
	Items       []*CartItem `json:"items"`
	TotalAmount int         `json:"totalAmount"`
	CreatedAt   time.Time   `json:"createdAt"`
	UpdatedAt   time.Time   `json:"updatedAt"`
}

func NewCart() *Cart {
	now := time.Now()
	return &Cart{
		ID:          uuid.New().String(),
		Items:       make([]*CartItem, 0),
		TotalAmount: 0,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

func NewCartItem(productID string, product *Product, quantity int) *CartItem {
	now := time.Now()
	return &CartItem{
		ID:        uuid.New().String(),
		ProductID: productID,
		Product:   product,
		Quantity:  quantity,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

func (c *Cart) AddItem(product *Product, quantity int) error {
	if !product.IsAvailable(quantity) {
		return ErrInsufficientStock
	}

	// 既存のアイテムがあるかチェック
	for _, item := range c.Items {
		if item.ProductID == product.ID {
			totalQuantity := item.Quantity + quantity
			if !product.IsAvailable(totalQuantity) {
				return ErrInsufficientStock
			}
			item.Quantity = totalQuantity
			item.UpdatedAt = time.Now()
			c.calculateTotal()
			c.UpdatedAt = time.Now()
			return nil
		}
	}

	// 新しいアイテムを追加
	newItem := NewCartItem(product.ID, product, quantity)
	c.Items = append(c.Items, newItem)
	c.calculateTotal()
	c.UpdatedAt = time.Now()
	return nil
}

func (c *Cart) UpdateItemQuantity(itemID string, quantity int) error {
	for _, item := range c.Items {
		if item.ID == itemID {
			if quantity <= 0 {
				return c.RemoveItem(itemID)
			}
			if !item.Product.IsAvailable(quantity) {
				return ErrInsufficientStock
			}
			item.Quantity = quantity
			item.UpdatedAt = time.Now()
			c.calculateTotal()
			c.UpdatedAt = time.Now()
			return nil
		}
	}
	return ErrItemNotFound
}

func (c *Cart) RemoveItem(itemID string) error {
	for i, item := range c.Items {
		if item.ID == itemID {
			c.Items = append(c.Items[:i], c.Items[i+1:]...)
			c.calculateTotal()
			c.UpdatedAt = time.Now()
			return nil
		}
	}
	return ErrItemNotFound
}

func (c *Cart) Clear() {
	c.Items = make([]*CartItem, 0)
	c.TotalAmount = 0
	c.UpdatedAt = time.Now()
}

func (c *Cart) GetItemCount() int {
	count := 0
	for _, item := range c.Items {
		count += item.Quantity
	}
	return count
}

func (c *Cart) IsEmpty() bool {
	return len(c.Items) == 0
}

func (c *Cart) calculateTotal() {
	total := 0
	for _, item := range c.Items {
		total += item.Product.Price * item.Quantity
	}
	c.TotalAmount = total
}
