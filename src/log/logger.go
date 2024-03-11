package log

import (
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"github.com/newrelic/go-agent/v3/newrelic"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"time"
)

func NewZapGinRequestLogger(logger *zap.Logger) gin.HandlerFunc {
	return ginzap.GinzapWithConfig(logger, &ginzap.Config{
		TimeFormat: time.RFC3339,
		UTC:        true,
		SkipPaths:  nil,
		TraceID:    false,
		Context: func(c *gin.Context) []zapcore.Field {
			var fields []zapcore.Field

			fields = append(fields, zap.Any("header", c.Request.Header))
			fields = append(fields, zap.Any("body", c.Request.Body))
			fields = append(fields, zap.Any("query", c.Request.URL.Query()))
			fields = append(fields, zap.Any("params", c.Params))
			fields = append(fields, zap.Any("response", c.Request.Response))
			fields = append(fields, zap.String("x-request-id", c.GetString("x-request-id")))

			txn := newrelic.FromContext(c)
			if txn != nil {
				fields = append(fields, zap.String("nrTraceId", txn.GetTraceMetadata().TraceID))
				fields = append(fields, zap.String("nrSpanId", txn.GetTraceMetadata().SpanID))
			}

			return fields
		},
	})
}
