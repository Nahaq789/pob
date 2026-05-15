package middleware

import (
	"context"
	"pob/pkg/tracing"

	"github.com/gin-gonic/gin"
)

func TraceMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		contextTrace := tracing.NewContextTrace()

		ctx := c.Request.Context()
		ctx = context.WithValue(ctx, tracing.ContextTraceKey, contextTrace.GenerateID())

		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}
