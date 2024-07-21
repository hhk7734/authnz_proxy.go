package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/hhk7734/zapx.go"
	"go.uber.org/zap"
)

func RequestID(generateIfNotExist bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader("X-Request-Id")
		if _, err := uuid.Parse(requestID); err != nil {
			if generateIfNotExist {
				requestID = uuid.New().String()
				c.Request.Header.Set("X-Request-Id", requestID)
			} else {
				return
			}
		}
		// Response
		c.Header("X-Request-Id", requestID)

		ctx := c.Request.Context()
		ctx = zapx.WithFields(ctx, zap.String("request_id", requestID))
		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}
