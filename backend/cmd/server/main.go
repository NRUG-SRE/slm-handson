package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/NRUG-SRE/slm-handson/backend/internal/domain/repository"
	"github.com/NRUG-SRE/slm-handson/backend/internal/infrastructure/monitoring"
	"github.com/NRUG-SRE/slm-handson/backend/internal/infrastructure/persistence/memory"
	"github.com/NRUG-SRE/slm-handson/backend/internal/interface/api"
	"github.com/NRUG-SRE/slm-handson/backend/internal/usecase"
	"github.com/NRUG-SRE/slm-handson/backend/pkg/config"
)

func main() {
	// 設定読み込み
	cfg := config.Load()
	if err := cfg.Validate(); err != nil {
		log.Fatalf("Configuration validation failed: %v", err)
	}

	// New Relic クライアント初期化
	nrClient, err := monitoring.NewNewRelicClient()
	if err != nil {
		log.Fatalf("Failed to initialize New Relic: %v", err)
	}

	// リポジトリ初期化
	var (
		productRepo repository.ProductRepository = memory.NewProductRepository()
		cartRepo    repository.CartRepository    = memory.NewCartRepository()
		orderRepo   repository.OrderRepository   = memory.NewOrderRepository()
	)

	// ユースケース初期化
	var (
		productUseCase = usecase.NewProductUseCase(productRepo)
		cartUseCase    = usecase.NewCartUseCase(cartRepo, productRepo)
		orderUseCase   = usecase.NewOrderUseCase(orderRepo, cartRepo, productRepo)
	)

	// ルーター初期化
	router := api.NewRouter(productUseCase, cartUseCase, orderUseCase, nrClient)
	ginEngine := router.SetupRoutes()

	// HTTPサーバー設定
	server := &http.Server{
		Addr:         fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port),
		Handler:      ginEngine,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// サーバー起動（ゴルーチンで実行）
	go func() {
		log.Printf("Starting server on %s", server.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// グレースフルシャットダウンの設定
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// シャットダウンシグナル待機
	<-quit
	log.Println("Shutting down server...")

	// グレースフルシャットダウン実行
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}