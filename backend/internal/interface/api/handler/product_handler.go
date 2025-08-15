package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/newrelic/go-agent/v3/newrelic"

	"github.com/NRUG-SRE/slm-handson/backend/internal/domain/entity"
	"github.com/NRUG-SRE/slm-handson/backend/internal/infrastructure/monitoring"
	"github.com/NRUG-SRE/slm-handson/backend/internal/interface/api/presenter"
	"github.com/NRUG-SRE/slm-handson/backend/internal/usecase"
)

type ProductHandler struct {
	productUseCase *usecase.ProductUseCase
	nrClient       *monitoring.NewRelicClient
}

func NewProductHandler(productUseCase *usecase.ProductUseCase, nrClient *monitoring.NewRelicClient) *ProductHandler {
	return &ProductHandler{
		productUseCase: productUseCase,
		nrClient:       nrClient,
	}
}

func (h *ProductHandler) GetProducts(c *gin.Context) {
	ctx := c.Request.Context()

	// New Relic トランザクションにカスタム属性を追加
	if txn := newrelic.FromContext(ctx); txn != nil {
		txn.AddAttribute("handler", "GetProducts")
	}

	products, err := h.productUseCase.GetAllProducts(ctx)
	if err != nil {
		h.nrClient.NoticeError(err)
		presenter.InternalServerErrorResponse(c, "Failed to get products")
		return
	}

	// New Relic カスタムイベント記録
	h.nrClient.RecordCustomEvent("ProductListView", map[string]interface{}{
		"productCount": len(products),
		"userAgent":    c.GetHeader("User-Agent"),
	})

	presenter.SuccessResponse(c, http.StatusOK, products)
}

func (h *ProductHandler) GetProduct(c *gin.Context) {
	ctx := c.Request.Context()
	productID := c.Param("id")

	// New Relic トランザクションにカスタム属性を追加
	if txn := newrelic.FromContext(ctx); txn != nil {
		txn.AddAttribute("handler", "GetProduct")
		txn.AddAttribute("product.id", productID)
	}

	if productID == "" {
		presenter.BadRequestResponse(c, "Product ID is required")
		return
	}

	product, err := h.productUseCase.GetProductByID(ctx, productID)
	if err != nil {
		if err == entity.ErrProductNotFound {
			presenter.NotFoundResponse(c, "Product not found")
			return
		}

		h.nrClient.NoticeError(err)
		presenter.InternalServerErrorResponse(c, "Failed to get product")
		return
	}

	// New Relic ビジネスメトリクス記録
	userID := c.GetHeader("X-User-ID") // 実際のアプリではセッションから取得
	if userID == "" {
		userID = "anonymous"
	}
	h.nrClient.RecordProductView(productID, userID)

	presenter.SuccessResponse(c, http.StatusOK, product)
}

// SLMデモ用のエラー発生エンドポイント
func (h *ProductHandler) TriggerError(c *gin.Context) {
	ctx := c.Request.Context()

	// New Relic トランザクションにカスタム属性を追加
	if txn := newrelic.FromContext(ctx); txn != nil {
		txn.AddAttribute("handler", "TriggerError")
		txn.AddAttribute("demo", "error_simulation")
	}

	// 意図的にエラーを発生させる
	err := fmt.Errorf("simulated error for SLM demonstration")
	h.nrClient.NoticeError(err)

	h.nrClient.RecordError("DemoError", "Intentional error for SLM testing", map[string]interface{}{
		"endpoint":  "/api/v1/error",
		"userAgent": c.GetHeader("User-Agent"),
	})

	presenter.InternalServerErrorResponse(c, "This is a simulated error for SLM demonstration")
}
