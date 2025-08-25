package middleware

import (
	"fmt"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/newrelic/go-agent/v3/integrations/nrgin"
	"github.com/newrelic/go-agent/v3/newrelic"

	"github.com/NRUG-SRE/slm-handson/backend/internal/infrastructure/monitoring"
)

func NewRelicMiddleware(nrClient *monitoring.NewRelicClient) gin.HandlerFunc {
	if nrClient.GetApplication() == nil {
		// New Relic が無効な場合は何もしないミドルウェアを返す
		return func(c *gin.Context) {
			c.Next()
		}
	}

	return nrgin.Middleware(nrClient.GetApplication())
}

// Distributed Tracingヘッダーを処理するミドルウェア
func DistributedTracingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// New Relic v3では、distributed tracingはW3C Trace Contextとして処理される
		// カスタムヘッダーは属性として記録
		if tracePayload := c.GetHeader("newrelic"); tracePayload != "" {
			if txn := newrelic.FromContext(c.Request.Context()); txn != nil {
				// カスタム属性として記録
				txn.AddAttribute("distributedTrace.payload", tracePayload)
			}
		}

		// フロントエンドからのNew Relicトレース情報
		if newRelicTrace := c.GetHeader("X-NewRelic-Trace"); newRelicTrace != "" {
			if txn := newrelic.FromContext(c.Request.Context()); txn != nil {
				txn.AddAttribute("frontend.trace.id", newRelicTrace)
			}
		}

		// ブラウザからのリクエストであることを記録
		if browserFlag := c.GetHeader("X-NewRelic-Browser"); browserFlag == "true" {
			if txn := newrelic.FromContext(c.Request.Context()); txn != nil {
				txn.AddAttribute("request.source", "browser")
			}
		}

		// W3C Trace Context ヘッダーの処理（New Relic v3で自動サポート）
		if traceParent := c.GetHeader("traceparent"); traceParent != "" {
			if txn := newrelic.FromContext(c.Request.Context()); txn != nil {
				// 明示的にカスタム属性として記録
				txn.AddAttribute("trace.parent", traceParent)
			}
		}

		// トレース状態ヘッダーも処理
		if traceState := c.GetHeader("tracestate"); traceState != "" {
			if txn := newrelic.FromContext(c.Request.Context()); txn != nil {
				txn.AddAttribute("trace.state", traceState)
			}
		}

		// セッションIDを受け取ってNew Relicに記録
		if sessionId := c.GetHeader("X-Session-ID"); sessionId != "" {
			if txn := newrelic.FromContext(c.Request.Context()); txn != nil {
				txn.AddAttribute("session.id", sessionId)
				// ユーザーIDとしても設定（RUMとの連携用）
				txn.AddAttribute("user.id", sessionId)
			}
			// Ginコンテキストにも保存（ハンドラーで利用可能に）
			c.Set("SessionID", sessionId)
		}

		c.Next()
	}
}

func LoggingMiddleware() gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		return fmt.Sprintf("%s - [%s] \"%s %s %s %d %s \"%s\" %s\"\n",
			param.ClientIP,
			param.TimeStamp.Format(time.RFC1123),
			param.Method,
			param.Path,
			param.Request.Proto,
			param.StatusCode,
			param.Latency,
			param.Request.UserAgent(),
			param.ErrorMessage,
		)
	})
}

func RequestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = generateRequestID()
		}

		c.Header("X-Request-ID", requestID)
		c.Set("RequestID", requestID)

		// New Relic トランザクションに属性を追加
		if txn := newrelic.FromContext(c.Request.Context()); txn != nil {
			txn.AddAttribute("request.id", requestID)
		}

		c.Next()
	}
}

func generateRequestID() string {
	// 簡単なリクエストID生成（実際にはより堅牢な実装が必要）
	return fmt.Sprintf("%d", time.Now().UnixNano())
}

func RecoveryMiddleware(nrClient *monitoring.NewRelicClient) gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		if err, ok := recovered.(string); ok {
			log.Printf("Panic recovered: %s", err)

			// New Relic にエラーを報告
			if nrClient != nil {
				nrClient.NoticeError(fmt.Errorf("panic: %s", err))
			}
		}

		c.AbortWithStatus(500)
	})
}
