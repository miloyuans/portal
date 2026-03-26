package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"portal/internal/middleware"
	"portal/internal/service"
)

// PortalHandler serves /api/portal endpoints.
type PortalHandler struct {
	service            *service.AppService
	defaultIdleTimeout int
}

// NewPortalHandler creates a PortalHandler.
func NewPortalHandler(service *service.AppService, defaultIdleTimeout int) *PortalHandler {
	return &PortalHandler{
		service:            service,
		defaultIdleTimeout: defaultIdleTimeout,
	}
}

// Apps returns visible portal apps.
func (h *PortalHandler) Apps(c *gin.Context) {
	apps, err := h.service.Apps(c.Request.Context(), middleware.CurrentSession(c))
	if err != nil {
		JSONError(c, http.StatusInternalServerError, "APPS_LOOKUP_FAILED", "failed to resolve visible apps", err.Error())
		return
	}
	JSONSuccess(c, http.StatusOK, apps)
}

// Realms returns projected realms.
func (h *PortalHandler) Realms(c *gin.Context) {
	realms, err := h.service.Realms(c.Request.Context())
	if err != nil {
		JSONError(c, http.StatusInternalServerError, "REALM_LIST_FAILED", "failed to load projected realms", err.Error())
		return
	}
	JSONSuccess(c, http.StatusOK, realms)
}

// Profile returns the current user's profile.
func (h *PortalHandler) Profile(c *gin.Context) {
	profile, err := h.service.Profile(c.Request.Context(), middleware.CurrentSession(c), h.defaultIdleTimeout)
	if err != nil {
		JSONError(c, http.StatusInternalServerError, "PROFILE_LOOKUP_FAILED", "failed to load current profile", err.Error())
		return
	}
	JSONSuccess(c, http.StatusOK, profile)
}
