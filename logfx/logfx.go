package logfx

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Ctx returns a logger with the Request ID pre-attached as a field
func Ctx(c *gin.Context, logger *zap.Logger) *zap.Logger {
	if rid, exists := c.Get("X-Request-ID"); exists {
		return logger.With(zap.String("request_id", rid.(string)))
	}
	return logger
}
