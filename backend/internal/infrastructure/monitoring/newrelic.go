package monitoring

import (
	"context"
	"log"
	"os"

	"github.com/newrelic/go-agent/v3/newrelic"
)

type NewRelicClient struct {
	app *newrelic.Application
}

func NewNewRelicClient() (*NewRelicClient, error) {
	apiKey := os.Getenv("NEW_RELIC_API_KEY")
	appName := os.Getenv("NEW_RELIC_APP_NAME")

	if apiKey == "" {
		log.Println("NEW_RELIC_API_KEY not set, New Relic monitoring disabled")
		return &NewRelicClient{app: nil}, nil
	}

	if appName == "" {
		appName = "slm-handson-api"
	}

	app, err := newrelic.NewApplication(
		newrelic.ConfigAppName(appName),
		newrelic.ConfigLicense(apiKey),
		newrelic.ConfigDistributedTracerEnabled(true),
		newrelic.ConfigEnabled(true),
	)

	if err != nil {
		log.Printf("Warning: Failed to initialize New Relic: %v. Continuing without New Relic monitoring.", err)
		return &NewRelicClient{app: nil}, nil
	}

	log.Printf("New Relic APM initialized for app: %s", appName)

	return &NewRelicClient{app: app}, nil
}

func (nr *NewRelicClient) StartTransaction(name string) *newrelic.Transaction {
	if nr.app == nil {
		return nil
	}
	return nr.app.StartTransaction(name)
}

func (nr *NewRelicClient) StartWebTransaction(name string) *newrelic.Transaction {
	if nr.app == nil {
		return nil
	}
	return nr.app.StartTransaction(name)
}

func (nr *NewRelicClient) RecordCustomEvent(eventType string, params map[string]interface{}) {
	if nr.app == nil {
		return
	}
	nr.app.RecordCustomEvent(eventType, params)
}

func (nr *NewRelicClient) RecordCustomMetric(name string, value float64) {
	if nr.app == nil {
		return
	}
	nr.app.RecordCustomMetric(name, value)
}

func (nr *NewRelicClient) NoticeError(err error) {
	if nr.app == nil {
		return
	}
	// New Relic v3では、エラーはトランザクション経由で報告される
	// ここでは単純にログ出力のみ行う
	log.Printf("New Relic Error: %v", err)
}

func (nr *NewRelicClient) GetApplication() *newrelic.Application {
	return nr.app
}

// ビジネスメトリクス記録用のヘルパー関数
func (nr *NewRelicClient) RecordProductView(productID string, userID string) {
	nr.RecordCustomEvent("ProductView", map[string]interface{}{
		"productId": productID,
		"userId":    userID,
	})
}

func (nr *NewRelicClient) RecordAddToCart(productID string, quantity int, userID string) {
	nr.RecordCustomEvent("AddToCart", map[string]interface{}{
		"productId": productID,
		"quantity":  quantity,
		"userId":    userID,
	})
}

func (nr *NewRelicClient) RecordPurchase(orderID string, amount float64, itemCount int, userID string) {
	nr.RecordCustomEvent("Purchase", map[string]interface{}{
		"orderId":   orderID,
		"amount":    amount,
		"itemCount": itemCount,
		"userId":    userID,
	})

	// 売上メトリクスも記録
	nr.RecordCustomMetric("Custom/Revenue", amount)
	nr.RecordCustomMetric("Custom/OrderCount", 1)
}

func (nr *NewRelicClient) RecordError(errorType string, message string, context map[string]interface{}) {
	nr.RecordCustomEvent("ApplicationError", map[string]interface{}{
		"errorType": errorType,
		"message":   message,
		"context":   context,
	})
}

// トランザクションにカスタム属性を追加するヘルパー
func AddCustomAttributes(txn *newrelic.Transaction, attributes map[string]interface{}) {
	if txn == nil {
		return
	}

	for key, value := range attributes {
		txn.AddAttribute(key, value)
	}
}

// データベースセグメント作成のヘルパー
func StartDatastoreSegment(txn *newrelic.Transaction, product, collection, operation string) *newrelic.DatastoreSegment {
	if txn == nil {
		return nil
	}

	return &newrelic.DatastoreSegment{
		StartTime:  txn.StartSegmentNow(),
		Product:    newrelic.DatastoreProduct(product),
		Collection: collection,
		Operation:  operation,
	}
}

// 外部サービス呼び出しセグメント作成のヘルパー
func StartExternalSegment(txn *newrelic.Transaction, url string) *newrelic.ExternalSegment {
	if txn == nil {
		return nil
	}

	return &newrelic.ExternalSegment{
		StartTime: txn.StartSegmentNow(),
		URL:       url,
	}
}

// コンテキストから New Relic トランザクションを取得
func GetTransactionFromContext(ctx context.Context) *newrelic.Transaction {
	return newrelic.FromContext(ctx)
}

// コンテキストに New Relic トランザクションを設定
func SetTransactionToContext(ctx context.Context, txn *newrelic.Transaction) context.Context {
	if txn == nil {
		return ctx
	}
	return newrelic.NewContext(ctx, txn)
}
