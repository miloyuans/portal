package middleware

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"portal/internal/repository"
	sessionpkg "portal/internal/session"
)

// IdleTimeout enforces portal idle timeout and refreshes lastActiveAt.
func IdleTimeout(manager *sessionpkg.Manager, settingsRepo *repository.SettingsRepository, defaultIdleTimeout int) gin.HandlerFunc {
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

		idleTimeoutMinutes := defaultIdleTimeout
		if settingsRepo != nil {
			settings, err := settingsRepo.GetGlobal(c.Request.Context(), defaultIdleTimeout)
			if err != nil {
				abortJSON(c, http.StatusInternalServerError, "SETTINGS_LOOKUP_FAILED", "failed to load session settings", err.Error())
				return
			}
			idleTimeoutMinutes = settings.IdleTimeoutMinutes
		}

		refreshed, err := manager.Touch(c.Request.Context(), session, idleTimeoutMinutes)
		if err != nil {
			abortJSON(c, http.StatusInternalServerError, "SESSION_TOUCH_FAILED", "failed to refresh portal session", err.Error())
			return
		}
		c.Set(SessionKey, refreshed)
		c.Next()
	}
}
