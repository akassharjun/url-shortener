package cmd

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/newrelic/go-agent/v3/newrelic"
	"go.uber.org/zap"
	"url-shortener/src/constants"
)

func SetRequestId(c *gin.Context) {
	var requestId string

	if c.GetHeader(string(constants.RequestId)) != "" {
		requestId = c.GetHeader("x-request-id")
	} else {
		// generate one
		requestId = uuid.New().String()
	}

	c.Set(string(constants.RequestId), requestId)
	c.Next()
}

func SetLoggerInContext(logger *zap.SugaredLogger) gin.HandlerFunc {
	return func(c *gin.Context) {
		requestId, exists := c.Get(string(constants.RequestId))

		if exists {
			logger = logger.With(zap.Any("x-request-id", requestId))
		}

		ctx := c.Request.Context()
		ctx = context.WithValue(ctx, "logger", logger)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

func UseJson(c *gin.Context) {
	c.Writer.Header().Set("Content-Type", "application/json")
	c.Next()
}

func SetNewRelicInContext() gin.HandlerFunc {
	return func(c *gin.Context) {
		//Setup context
		ctx := c.Request.Context()

		//Set newrelic context
		var txn newrelic.Transaction
		//newRelicTransaction is the key populated by nrgin Middleware
		value, exists := c.Get("newRelicTransaction")
		if exists {
			if v, ok := value.(newrelic.Transaction); ok {
				txn = v
			}
			ctx = context.WithValue(ctx, constants.NewRelicTransaction, txn)
		}
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}
