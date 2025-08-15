package entity

import (
	"time"

	"github.com/google/uuid"
)

type OrderStatus string

const (
	OrderStatusPending   OrderStatus = "pending"
	OrderStatusCompleted OrderStatus = "completed"
	OrderStatusFailed    OrderStatus = "failed"
	OrderStatusCanceled  OrderStatus = "canceled"
)

type OrderItem struct {
	ID        string    `json:"id"`
	ProductID string    `json:"productId"`
	Product   *Product  `json:"product"`
	Quantity  int       `json:"quantity"`
	Price     int       `json:"price"` // 注文時の価格を保存
	CreatedAt time.Time `json:"createdAt"`
}

type Order struct {
	ID          string       `json:"id"`
	Items       []*OrderItem `json:"items"`
	TotalAmount int          `json:"totalAmount"`
	Status      OrderStatus  `json:"status"`
	CreatedAt   time.Time    `json:"createdAt"`
	UpdatedAt   time.Time    `json:"updatedAt"`
}

func NewOrder(cart *Cart) (*Order, error) {
	if cart.IsEmpty() {
		return nil, ErrEmptyCart
	}

	now := time.Now()
	order := &Order{
		ID:          uuid.New().String(),
		Items:       make([]*OrderItem, 0, len(cart.Items)),
		TotalAmount: cart.TotalAmount,
		Status:      OrderStatusPending,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	// カートアイテムを注文アイテムに変換
	for _, cartItem := range cart.Items {
		orderItem := &OrderItem{
			ID:        uuid.New().String(),
			ProductID: cartItem.ProductID,
			Product:   cartItem.Product,
			Quantity:  cartItem.Quantity,
			Price:     cartItem.Product.Price, // 注文時の価格を記録
			CreatedAt: now,
		}
		order.Items = append(order.Items, orderItem)
	}

	return order, nil
}

func (o *Order) Complete() error {
	if o.Status != OrderStatusPending {
		return ErrInvalidOrderStatus
	}
	o.Status = OrderStatusCompleted
	o.UpdatedAt = time.Now()
	return nil
}

func (o *Order) Fail() error {
	if o.Status != OrderStatusPending {
		return ErrInvalidOrderStatus
	}
	o.Status = OrderStatusFailed
	o.UpdatedAt = time.Now()
	return nil
}

func (o *Order) Cancel() error {
	if o.Status == OrderStatusCompleted {
		return ErrInvalidOrderStatus
	}
	o.Status = OrderStatusCanceled
	o.UpdatedAt = time.Now()
	return nil
}

func (o *Order) GetItemCount() int {
	count := 0
	for _, item := range o.Items {
		count += item.Quantity
	}
	return count
}

func (o *Order) IsCompleted() bool {
	return o.Status == OrderStatusCompleted
}

func (o *Order) IsPending() bool {
	return o.Status == OrderStatusPending
}

func (o *Order) IsFailed() bool {
	return o.Status == OrderStatusFailed
}

func (o *Order) IsCanceled() bool {
	return o.Status == OrderStatusCanceled
}
