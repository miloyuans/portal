package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"portal/internal/config"
)

func CORS(cfg config.Config) gin.HandlerFunc {
	allowMap := make(map[string]struct{}, len(cfg.Server.AllowedOrigins))
	for _, origin := range cfg.Server.AllowedOrigins {
		allowMap[strings.TrimSpace(origin)] = struct{}{}
	}

	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")
		if _, ok := allowMap[origin]; ok {
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
			c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
			c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Request-ID")
			c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		}

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}
