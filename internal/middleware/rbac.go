package middleware

import (
	"net/http"
	"slices"

	"github.com/gin-gonic/gin"
)

// RequirePortalAdmin restricts access to portal_admin users.
func RequirePortalAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := CurrentSession(c)
		if slices.Contains(session.RealmRoles, "portal_admin") {
			c.Next()
			return
		}
		abortJSON(c, http.StatusForbidden, "FORBIDDEN", "portal_admin role required", nil)
	}
}
