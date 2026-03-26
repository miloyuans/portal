package middleware

import (
	"net/http"
	"slices"

	"github.com/gin-gonic/gin"

	"portal/internal/config"
)

func RequireAdmin(cfg config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		session := CurrentSession(c)
		for _, role := range session.RealmRoles {
			if slices.Contains(cfg.Permission.AdminRealmRoles, role) {
				c.Next()
				return
			}
		}
		abortJSON(c, http.StatusForbidden, "FORBIDDEN", "admin role required", nil)
	}
}
