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

type CartHandler struct {
	cartUseCase *usecase.CartUseCase
	nrClient    *monitoring.NewRelicClient
}

type AddToCartRequest struct {
	ProductID string `json:"productId" binding:"required"`
	Quantity  int    `json:"quantity" binding:"required,min=1"`
}

type UpdateCartItemRequest struct {
	Quantity int `json:"quantity"`
}

// DefaultCartIDはconstants.goで定義

func NewCartHandler(cartUseCase *usecase.CartUseCase, nrClient *monitoring.NewRelicClient) *CartHandler {
	return &CartHandler{
		cartUseCase: cartUseCase,
		nrClient:    nrClient,
	}
}

func (h *CartHandler) GetCart(c *gin.Context) {
	ctx := c.Request.Context()
	cartID := DefaultCartID // 実際のアプリではユーザーセッションから取得
	
	// New Relic トランザクションにカスタム属性を追加
	if txn := newrelic.FromContext(ctx); txn != nil {
		txn.AddAttribute("handler", "GetCart")
		txn.AddAttribute("cart.id", cartID)
	}

	cart, err := h.cartUseCase.GetCart(ctx, cartID)
	if err != nil {
		h.nrClient.NoticeError(err)
		presenter.InternalServerErrorResponse(c, "Failed to get cart")
		return
	}

	presenter.SuccessResponse(c, http.StatusOK, cart)
}

func (h *CartHandler) AddToCart(c *gin.Context) {
	ctx := c.Request.Context()
	cartID := DefaultCartID
	
	var req AddToCartRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		presenter.BadRequestResponse(c, "Invalid request body")
		return
	}
	
	// New Relic トランザクションにカスタム属性を追加
	if txn := newrelic.FromContext(ctx); txn != nil {
		txn.AddAttribute("handler", "AddToCart")
		txn.AddAttribute("cart.id", cartID)
		txn.AddAttribute("product.id", req.ProductID)
		txn.AddAttribute("quantity", req.Quantity)
	}

	cart, err := h.cartUseCase.AddToCart(ctx, cartID, req.ProductID, req.Quantity)
	if err != nil {
		if err == entity.ErrProductNotFound {
			presenter.NotFoundResponse(c, "Product not found")
			return
		}
		if err == entity.ErrInsufficientStock {
			presenter.UnprocessableEntityResponse(c, "Insufficient stock")
			return
		}
		
		h.nrClient.NoticeError(err)
		presenter.InternalServerErrorResponse(c, "Failed to add item to cart")
		return
	}

	// New Relic ビジネスメトリクス記録
	userID := c.GetHeader("X-User-ID")
	if userID == "" {
		userID = "anonymous"
	}
	h.nrClient.RecordAddToCart(req.ProductID, req.Quantity, userID)

	presenter.SuccessResponse(c, http.StatusOK, cart)
}

func (h *CartHandler) UpdateCartItem(c *gin.Context) {
	ctx := c.Request.Context()
	cartID := DefaultCartID
	itemID := c.Param("id")
	
	var req UpdateCartItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		presenter.BadRequestResponse(c, "Invalid request body")
		return
	}
	
	// 数量のバリデーション（負の値のみ拒否、0は削除として許可）
	if req.Quantity < 0 {
		presenter.BadRequestResponse(c, "Quantity cannot be negative")
		return
	}
	
	// New Relic トランザクションにカスタム属性を追加
	if txn := newrelic.FromContext(ctx); txn != nil {
		txn.AddAttribute("handler", "UpdateCartItem")
		txn.AddAttribute("cart.id", cartID)
		txn.AddAttribute("item.id", itemID)
		txn.AddAttribute("quantity", req.Quantity)
	}

	cart, err := h.cartUseCase.UpdateCartItem(ctx, cartID, itemID, req.Quantity)
	if err != nil {
		if err == entity.ErrItemNotFound {
			presenter.NotFoundResponse(c, "Cart item not found")
			return
		}
		if err == entity.ErrInsufficientStock {
			presenter.UnprocessableEntityResponse(c, "Insufficient stock")
			return
		}
		
		h.nrClient.NoticeError(err)
		presenter.InternalServerErrorResponse(c, "Failed to update cart item")
		return
	}

	presenter.SuccessResponse(c, http.StatusOK, cart)
}