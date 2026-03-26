package middleware

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	sessionpkg "portal/internal/session"
)

func IdleTimeout(manager *sessionpkg.Manager) gin.HandlerFunc {
	return func(c *gin.Context) {
		session := CurrentSession(c)
		if err := manager.Validate(session); err != nil {
			if errors.Is(err, sessionpkg.ErrSessionExpired) {
				_ = manager.Delete(c.Request.Context(), session.SessionID)
				manager.ClearCookie(c)
				abortJSON(c, http.StatusUnauthorized, "SESSION_EXPIRED", "portal session expired", nil)
				return
			}
			abortJSON(c, http.StatusUnauthorized, "INVALID_SESSION", "portal session is invalid", nil)
			return
		}

		if err := manager.Touch(c.Request.Context(), session); err != nil {
			abortJSON(c, http.StatusInternalServerError, "SESSION_TOUCH_FAILED", "failed to refresh portal session", err.Error())
			return
		}
		c.Next()
	}
}
