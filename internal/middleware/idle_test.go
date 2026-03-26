package middleware

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"

	"portal/internal/config"
	"portal/internal/model"
	sessionpkg "portal/internal/session"
)

func TestIdleTimeoutMiddlewareReturnsSessionExpired(t *testing.T) {
	gin.SetMode(gin.TestMode)
	manager := sessionpkg.NewManager(nil, config.Config{
		Session: config.SessionConfig{CookieName: "portal_session"},
	})

	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set(SessionKey, model.PortalSession{
			SessionID:         "session-1",
			ExpiresAt:         time.Now().UTC().Add(-time.Minute),
			AbsoluteExpiresAt: time.Now().UTC().Add(time.Hour),
		})
		c.Next()
	})
	router.Use(IdleTimeout(manager, nil, 15))
	router.GET("/protected", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	request := httptest.NewRequest(http.MethodGet, "/protected", nil)
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	if response.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", response.Code)
	}

	var payload map[string]any
	if err := json.Unmarshal(response.Body.Bytes(), &payload); err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}
}
