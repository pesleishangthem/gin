package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

const HeaderRequestID = "X-Request-ID"

type LogFields struct {
	Method    string
	Path      string
	Status    int
	Latency   time.Duration
	ClientIP  string
	UserAgent string
	RequestID string
	Error     string
}

func LogMiddleware(log *zerolog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		reqID := c.GetHeader(RequestIDKey)
		c.Writer.Header().Set(HeaderRequestID, reqID)
		c.Next()
		LogRequest(log, LogFields{
			Method:    c.Request.Method,
			Path:      c.Request.URL.Path,
			Status:    c.Writer.Status(),
			Latency:   time.Since(start),
			ClientIP:  c.ClientIP(),
			UserAgent: c.Request.UserAgent(),
			RequestID: reqID,
		})
	}
}

func LogRequest(log *zerolog.Logger, f LogFields) {
	event := log.Info()
	if f.Status >= 500 {
		event = log.Error()
	} else if f.Status >= 400 {
		event = log.Warn()
	}
	event.
		Str("method", f.Method).
		Str("path", f.Path).
		Int("status", f.Status).
		Dur("latency", f.Latency).
		Str("client_ip", f.ClientIP).
		Str("user_agent", f.UserAgent).
		Str("request_id", f.RequestID).
		Msg("http_request")
}
