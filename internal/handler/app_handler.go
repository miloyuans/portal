package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"portal/internal/middleware"
	"portal/internal/service"
)

type AppHandler struct {
	service *service.AppService
}

func NewAppHandler(service *service.AppService) *AppHandler {
	return &AppHandler{
		service: service,
	}
}

func (h *AppHandler) Me(c *gin.Context) {
	profile, err := h.service.Me(c.Request.Context(), middleware.CurrentSession(c))
	if err != nil {
		JSONError(c, http.StatusInternalServerError, "ME_LOOKUP_FAILED", "failed to resolve current session", err.Error())
		return
	}
	JSONSuccess(c, http.StatusOK, profile)
}

func (h *AppHandler) Apps(c *gin.Context) {
	apps, err := h.service.Apps(c.Request.Context(), middleware.CurrentSession(c))
	if err != nil {
		JSONError(c, http.StatusInternalServerError, "APPS_LOOKUP_FAILED", "failed to resolve visible apps", err.Error())
		return
	}
	JSONSuccess(c, http.StatusOK, apps)
}
