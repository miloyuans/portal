package handler

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"portal/internal/config"
)

func externalBaseURL(c *gin.Context, cfg config.Config) string {
	scheme := forwardedHeaderValue(c.Request.Header, "X-Forwarded-Proto")
	if scheme == "" {
		if c.Request.TLS != nil {
			scheme = "https"
		} else {
			scheme = "http"
		}
	}

	host := forwardedHeaderValue(c.Request.Header, "X-Forwarded-Host")
	if host == "" {
		host = c.Request.Host
	}

	if host == "" {
		return strings.TrimRight(cfg.Server.PublicWebURL, "/")
	}
	return strings.TrimRight(scheme+"://"+host, "/")
}

func callbackURL(c *gin.Context, cfg config.Config) string {
	return externalBaseURL(c, cfg) + "/api/auth/callback"
}

func forwardedHeaderValue(header http.Header, key string) string {
	value := strings.TrimSpace(header.Get(key))
	if value == "" {
		return ""
	}
	if index := strings.Index(value, ","); index >= 0 {
		value = value[:index]
	}
	return strings.TrimSpace(value)
}
