package api

import (
	"github.com/gin-gonic/gin"

	"github.com/NRUG-SRE/slm-handson/backend/internal/infrastructure/monitoring"
	"github.com/NRUG-SRE/slm-handson/backend/internal/interface/api/handler"
	"github.com/NRUG-SRE/slm-handson/backend/internal/interface/api/middleware"
	"github.com/NRUG-SRE/slm-handson/backend/internal/usecase"
)

type Router struct {
	healthHandler  *handler.HealthHandler
	productHandler *handler.ProductHandler
	cartHandler    *handler.CartHandler
	orderHandler   *handler.OrderHandler
	swaggerHandler *handler.SwaggerHandler
	nrClient       *monitoring.NewRelicClient
}

func NewRouter(
	productUseCase *usecase.ProductUseCase,
	cartUseCase *usecase.CartUseCase,
	orderUseCase *usecase.OrderUseCase,
	nrClient *monitoring.NewRelicClient,
) *Router {
	return &Router{
		healthHandler:  handler.NewHealthHandler(),
		productHandler: handler.NewProductHandler(productUseCase, nrClient),
		cartHandler:    handler.NewCartHandler(cartUseCase, nrClient),
		orderHandler:   handler.NewOrderHandler(orderUseCase, nrClient),
		swaggerHandler: handler.NewSwaggerHandler(),
		nrClient:       nrClient,
	}
}

func (r *Router) SetupRoutes() *gin.Engine {
	// プロダクション環境用の設定
	gin.SetMode(gin.ReleaseMode)

	router := gin.New()

	// ミドルウェア設定
	router.Use(middleware.LoggingMiddleware())
	router.Use(middleware.RecoveryMiddleware(r.nrClient))
	router.Use(middleware.CORS())
	router.Use(middleware.NewRelicMiddleware(r.nrClient))
	router.Use(middleware.RequestIDMiddleware())

	// ヘルスチェックエンドポイント
	router.GET("/health", r.healthHandler.HealthCheck)

	// APIルートグループ
	apiV1 := router.Group("/api")
	{
		// 商品関連エンドポイント
		apiV1.GET("/products", r.productHandler.GetProducts)
		apiV1.GET("/products/:id", r.productHandler.GetProduct)

		// カート関連エンドポイント
		apiV1.GET("/cart", r.cartHandler.GetCart)
		apiV1.POST("/cart/items", r.cartHandler.AddToCart)
		apiV1.PUT("/cart/items/:id", r.cartHandler.UpdateCartItem)

		// 注文関連エンドポイント
		apiV1.POST("/orders", r.orderHandler.CreateOrder)
		apiV1.GET("/orders/:id", r.orderHandler.GetOrder)
		apiV1.GET("/orders", r.orderHandler.GetOrders) // オプション: 全注文取得

		// SLMデモ用エンドポイント
		apiV1.GET("/v1/error", r.productHandler.TriggerError)

		// Swagger APIドキュメントエンドポイント
		docsGroup := apiV1.Group("/docs")
		{
			docsGroup.GET("", r.swaggerHandler.SwaggerUI)
			docsGroup.GET("/swagger.yaml", r.swaggerHandler.ServeSwaggerYAML)
		}
	}

	return router
}
