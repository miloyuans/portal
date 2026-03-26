package middleware

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
)

func Logger(logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		started := time.Now()
		c.Next()
		requestID, _ := c.Get(RequestIDKey)
		logger.Info("request completed",
			slog.String("requestID", fmt.Sprint(requestID)),
			slog.String("method", c.Request.Method),
			slog.String("path", c.Request.URL.Path),
			slog.Int("status", c.Writer.Status()),
			slog.Duration("duration", time.Since(started)),
			slog.String("ip", c.ClientIP()),
		)
	}
}
