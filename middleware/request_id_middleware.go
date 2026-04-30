package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const RequestIDKey = "X-Request-ID"

func RequestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		// 1. Check if the client/load-balancer already sent a Request ID
		rid := c.GetHeader(RequestIDKey)
		if rid == "" {
			rid = uuid.New().String()
		}

		// 2. Set it in the Context and the Response Header
		c.Set(RequestIDKey, rid)
		c.Header(RequestIDKey, rid)

		c.Next()
	}
}
