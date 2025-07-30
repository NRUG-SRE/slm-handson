package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/newrelic/go-agent/v3/newrelic"
	
	"github.com/NRUG-SRE/slm-handson/backend/internal/domain/entity"
	"github.com/NRUG-SRE/slm-handson/backend/internal/infrastructure/monitoring"
	"github.com/NRUG-SRE/slm-handson/backend/internal/interface/api/presenter"
	"github.com/NRUG-SRE/slm-handson/backend/internal/usecase"
)

type OrderHandler struct {
	orderUseCase *usecase.OrderUseCase
	nrClient     *monitoring.NewRelicClient
}

type CreateOrderRequest struct {
	Items []CreateOrderItem `json:"items" binding:"required"`
}

type CreateOrderItem struct {
	ProductID string `json:"productId" binding:"required"`
	Quantity  int    `json:"quantity" binding:"required,min=1"`
}

// cart_handlerと同じDefaultCartIDを使用

func NewOrderHandler(orderUseCase *usecase.OrderUseCase, nrClient *monitoring.NewRelicClient) *OrderHandler {
	return &OrderHandler{
		orderUseCase: orderUseCase,
		nrClient:     nrClient,
	}
}

func (h *OrderHandler) CreateOrder(c *gin.Context) {
	ctx := c.Request.Context()
	cartID := DefaultCartID // 実際のアプリではユーザーセッションから取得
	
	// New Relic トランザクションにカスタム属性を追加
	if txn := newrelic.FromContext(ctx); txn != nil {
		txn.AddAttribute("handler", "CreateOrder")
		txn.AddAttribute("cart.id", cartID)
	}

	order, err := h.orderUseCase.CreateOrder(ctx, cartID)
	if err != nil {
		if err == entity.ErrEmptyCart {
			presenter.UnprocessableEntityResponse(c, "Cart is empty")
			return
		}
		
		h.nrClient.NoticeError(err)
		presenter.InternalServerErrorResponse(c, "Failed to create order")
		return
	}

	// New Relic ビジネスメトリクス記録
	userID := c.GetHeader("X-User-ID")
	if userID == "" {
		userID = "anonymous"
	}
	h.nrClient.RecordPurchase(
		order.ID,
		float64(order.TotalAmount),
		order.GetItemCount(),
		userID,
	)

	presenter.SuccessResponse(c, http.StatusCreated, order)
}

func (h *OrderHandler) GetOrder(c *gin.Context) {
	ctx := c.Request.Context()
	orderID := c.Param("id")
	
	// New Relic トランザクションにカスタム属性を追加
	if txn := newrelic.FromContext(ctx); txn != nil {
		txn.AddAttribute("handler", "GetOrder")
		txn.AddAttribute("order.id", orderID)
	}

	if orderID == "" {
		presenter.BadRequestResponse(c, "Order ID is required")
		return
	}

	order, err := h.orderUseCase.GetOrder(ctx, orderID)
	if err != nil {
		if err == entity.ErrOrderNotFound {
			presenter.NotFoundResponse(c, "Order not found")
			return
		}
		
		h.nrClient.NoticeError(err)
		presenter.InternalServerErrorResponse(c, "Failed to get order")
		return
	}

	presenter.SuccessResponse(c, http.StatusOK, order)
}

func (h *OrderHandler) GetOrders(c *gin.Context) {
	ctx := c.Request.Context()
	
	// New Relic トランザクションにカスタム属性を追加
	if txn := newrelic.FromContext(ctx); txn != nil {
		txn.AddAttribute("handler", "GetOrders")
	}

	orders, err := h.orderUseCase.GetAllOrders(ctx)
	if err != nil {
		h.nrClient.NoticeError(err)
		presenter.InternalServerErrorResponse(c, "Failed to get orders")
		return
	}

	presenter.SuccessResponse(c, http.StatusOK, orders)
}