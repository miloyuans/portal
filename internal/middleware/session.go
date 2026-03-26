package middleware

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"

	"portal/internal/model"
	sessionpkg "portal/internal/session"
)

const SessionKey = "portalSession"

func Session(manager *sessionpkg.Manager) gin.HandlerFunc {
	return func(c *gin.Context) {
		session, err := manager.GetByRequest(c.Request.Context(), c.Request)
		if err != nil {
			if errors.Is(err, http.ErrNoCookie) || errors.Is(err, mongo.ErrNoDocuments) {
				abortJSON(c, http.StatusUnauthorized, "AUTH_REQUIRED", "portal session not found", nil)
				return
			}
			abortJSON(c, http.StatusInternalServerError, "SESSION_LOOKUP_FAILED", "failed to load portal session", err.Error())
			return
		}
		c.Set(SessionKey, session)
		c.Next()
	}
}

func CurrentSession(c *gin.Context) model.PortalSession {
	value, _ := c.Get(SessionKey)
	session, _ := value.(model.PortalSession)
	return session
}
