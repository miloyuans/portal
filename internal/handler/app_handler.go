package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"portal/internal/middleware"
	"portal/internal/permission"
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

// Launch returns the final launch target for a visible app.
func (h *PortalHandler) Launch(c *gin.Context) {
	clientID := c.Param("clientId")
	if clientID == "" {
		JSONError(c, http.StatusBadRequest, "CLIENT_ID_REQUIRED", "clientId is required", nil)
		return
	}

	launch, err := h.service.Launch(c.Request.Context(), middleware.CurrentSession(c), clientID)
	if err != nil {
		switch {
		case errors.Is(err, permission.ErrAppNotVisible):
			JSONError(c, http.StatusForbidden, "APP_NOT_VISIBLE", "the requested app is not visible for the current session", nil)
		case errors.Is(err, permission.ErrLaunchDisabled):
			JSONError(c, http.StatusConflict, "APP_LAUNCH_DISABLED", "the requested app is visible but launch is disabled", nil)
		case errors.Is(err, permission.ErrLaunchTargetMissing):
			JSONError(c, http.StatusConflict, "APP_LAUNCH_TARGET_MISSING", "the requested app has no launch target configured", nil)
		default:
			JSONError(c, http.StatusInternalServerError, "APP_LAUNCH_FAILED", "failed to resolve app launch target", err.Error())
		}
		return
	}
	JSONSuccess(c, http.StatusOK, launch)
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
