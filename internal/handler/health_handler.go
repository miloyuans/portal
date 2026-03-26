package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"portal/internal/auth"
	"portal/internal/repository"
)

type HealthHandler struct {
	db   *repository.Mongo
	oidc *auth.OIDCClient
}

func NewHealthHandler(db *repository.Mongo, oidc *auth.OIDCClient) *HealthHandler {
	return &HealthHandler{
		db:   db,
		oidc: oidc,
	}
}

func (h *HealthHandler) Healthz(c *gin.Context) {
	JSONSuccess(c, http.StatusOK, gin.H{"status": "ok"})
}

func (h *HealthHandler) Readyz(c *gin.Context) {
	if err := h.db.Ping(c.Request.Context()); err != nil {
		JSONError(c, http.StatusServiceUnavailable, "MONGO_NOT_READY", "mongo is not ready", err.Error())
		return
	}
	if !h.oidc.Ready() {
		JSONError(c, http.StatusServiceUnavailable, "OIDC_NOT_READY", "oidc provider is not ready", nil)
		return
	}
	JSONSuccess(c, http.StatusOK, gin.H{"status": "ready"})
}
