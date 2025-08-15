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
