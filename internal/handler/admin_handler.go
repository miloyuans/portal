package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"portal/internal/config"
	"portal/internal/middleware"
	"portal/internal/model"
	"portal/internal/repository"
)

type AdminHandler struct {
	repos *repository.Repositories
	cfg   config.Config
}

func NewAdminHandler(repos *repository.Repositories, cfg config.Config) *AdminHandler {
	return &AdminHandler{
		repos: repos,
		cfg:   cfg,
	}
}

type upsertClientMetaRequest struct {
	ClientID            string   `json:"clientId"`
	DisplayName         string   `json:"displayName"`
	Description         string   `json:"description"`
	TargetURL           string   `json:"targetUrl"`
	Icon                string   `json:"icon"`
	Category            string   `json:"category"`
	SortOrder           int      `json:"sortOrder"`
	Enabled             bool     `json:"enabled"`
	ShowInPortal        bool     `json:"showInPortal"`
	RequiredRealmRoles  []string `json:"requiredRealmRoles"`
	RequiredClientRoles []string `json:"requiredClientRoles"`
	Tags                []string `json:"tags"`
}

type updateSettingsRequest struct {
	IdleTimeoutMinutes int `json:"idleTimeoutMinutes"`
}

func (h *AdminHandler) ListClientMetas(c *gin.Context) {
	session := middleware.CurrentSession(c)
	metas, err := h.repos.ClientMetas.ListByRealm(c.Request.Context(), session.Realm)
	if err != nil {
		JSONError(c, http.StatusInternalServerError, "CLIENT_META_LIST_FAILED", "failed to load portal client metadata", err.Error())
		return
	}
	JSONSuccess(c, http.StatusOK, metas)
}

func (h *AdminHandler) UpsertClientMeta(c *gin.Context) {
	session := middleware.CurrentSession(c)
	var request upsertClientMetaRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		JSONError(c, http.StatusBadRequest, "INVALID_REQUEST", "invalid client meta payload", err.Error())
		return
	}

	clientID := request.ClientID
	if routeClientID := c.Param("clientId"); routeClientID != "" {
		clientID = routeClientID
	}
	if clientID == "" {
		JSONError(c, http.StatusBadRequest, "CLIENT_ID_REQUIRED", "clientId is required", nil)
		return
	}

	meta := model.PortalClientMeta{
		Realm:               session.Realm,
		ClientID:            clientID,
		DisplayName:         request.DisplayName,
		Description:         request.Description,
		TargetURL:           request.TargetURL,
		Icon:                request.Icon,
		Category:            request.Category,
		SortOrder:           request.SortOrder,
		Enabled:             request.Enabled,
		ShowInPortal:        request.ShowInPortal,
		RequiredRealmRoles:  request.RequiredRealmRoles,
		RequiredClientRoles: request.RequiredClientRoles,
		Tags:                request.Tags,
		UpdatedBy:           session.Username,
		CreatedAt:           time.Now().UTC(),
	}

	if err := h.repos.ClientMetas.Upsert(c.Request.Context(), meta); err != nil {
		JSONError(c, http.StatusInternalServerError, "CLIENT_META_SAVE_FAILED", "failed to save portal client meta", err.Error())
		return
	}
	JSONSuccess(c, http.StatusOK, meta)
}

func (h *AdminHandler) DeleteClientMeta(c *gin.Context) {
	session := middleware.CurrentSession(c)
	clientID := c.Param("clientId")
	if clientID == "" {
		JSONError(c, http.StatusBadRequest, "CLIENT_ID_REQUIRED", "clientId is required", nil)
		return
	}
	if err := h.repos.ClientMetas.Delete(c.Request.Context(), session.Realm, clientID); err != nil {
		JSONError(c, http.StatusInternalServerError, "CLIENT_META_DELETE_FAILED", "failed to delete portal client meta", err.Error())
		return
	}
	JSONSuccess(c, http.StatusOK, gin.H{"deleted": true, "clientId": clientID})
}

func (h *AdminHandler) GetSettings(c *gin.Context) {
	session := middleware.CurrentSession(c)
	settings, err := h.repos.Settings.GetByRealm(c.Request.Context(), session.Realm, h.cfg.Session.DefaultIdleTimeoutMinutes)
	if err != nil {
		JSONError(c, http.StatusInternalServerError, "SETTINGS_LOOKUP_FAILED", "failed to load portal settings", err.Error())
		return
	}
	JSONSuccess(c, http.StatusOK, settings)
}

func (h *AdminHandler) UpdateSettings(c *gin.Context) {
	session := middleware.CurrentSession(c)
	var request updateSettingsRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		JSONError(c, http.StatusBadRequest, "INVALID_REQUEST", "invalid settings payload", err.Error())
		return
	}
	if request.IdleTimeoutMinutes <= 0 {
		request.IdleTimeoutMinutes = h.cfg.Session.DefaultIdleTimeoutMinutes
	}

	settings := model.PortalSettings{
		Realm:              session.Realm,
		IdleTimeoutMinutes: request.IdleTimeoutMinutes,
		CreatedAt:          time.Now().UTC(),
	}
	if err := h.repos.Settings.Upsert(c.Request.Context(), settings); err != nil {
		JSONError(c, http.StatusInternalServerError, "SETTINGS_SAVE_FAILED", "failed to update portal settings", err.Error())
		return
	}
	JSONSuccess(c, http.StatusOK, settings)
}
