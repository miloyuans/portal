package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"portal/internal/middleware"
	"portal/internal/model"
	"portal/internal/repository"
	"portal/internal/service"
)

// AdminHandler serves /api/admin endpoints.
type AdminHandler struct {
	repos               *repository.Repositories
	service             *service.AppService
	defaultIdleTimeout  int
}

// NewAdminHandler creates an AdminHandler.
func NewAdminHandler(repos *repository.Repositories, service *service.AppService, defaultIdleTimeout int) *AdminHandler {
	return &AdminHandler{
		repos:              repos,
		service:            service,
		defaultIdleTimeout: defaultIdleTimeout,
	}
}

// ListRealms returns projected realms for admin use.
func (h *AdminHandler) ListRealms(c *gin.Context) {
	realms, err := h.repos.Realms.List(c.Request.Context())
	if err != nil {
		JSONError(c, http.StatusInternalServerError, "REALM_LIST_FAILED", "failed to load realms", err.Error())
		return
	}
	JSONSuccess(c, http.StatusOK, realms)
}

// ListClients returns projected clients merged with portal metadata.
func (h *AdminHandler) ListClients(c *gin.Context) {
	session := middleware.CurrentSession(c)
	clients, err := h.repos.Clients.ListByRealm(c.Request.Context(), session.RealmID)
	if err != nil {
		JSONError(c, http.StatusInternalServerError, "CLIENT_LIST_FAILED", "failed to load clients", err.Error())
		return
	}
	metas, err := h.repos.ClientMetas.ListByRealm(c.Request.Context(), session.RealmID)
	if err != nil {
		JSONError(c, http.StatusInternalServerError, "CLIENT_META_LIST_FAILED", "failed to load client metadata", err.Error())
		return
	}

	metaByClientID := make(map[string]model.PortalClientMeta, len(metas))
	for _, meta := range metas {
		metaByClientID[meta.ClientID] = meta
	}

	type adminClientRow struct {
		Client model.ClientProjection  `json:"client"`
		Meta   *model.PortalClientMeta `json:"meta,omitempty"`
	}

	rows := make([]adminClientRow, 0, len(clients))
	for _, client := range clients {
		row := adminClientRow{Client: client}
		if meta, ok := metaByClientID[client.ClientID]; ok {
			row.Meta = &meta
		}
		rows = append(rows, row)
	}

	JSONSuccess(c, http.StatusOK, rows)
}

// UpdateClientMeta updates portal client metadata.
func (h *AdminHandler) UpdateClientMeta(c *gin.Context) {
	session := middleware.CurrentSession(c)
	clientID := c.Param("clientId")
	if clientID == "" {
		JSONError(c, http.StatusBadRequest, "CLIENT_ID_REQUIRED", "clientId is required", nil)
		return
	}

	var payload model.PortalClientMeta
	if err := c.ShouldBindJSON(&payload); err != nil {
		JSONError(c, http.StatusBadRequest, "INVALID_REQUEST", "invalid client meta payload", err.Error())
		return
	}
	payload.RealmID = session.RealmID
	payload.ClientID = clientID

	if err := h.repos.ClientMetas.Upsert(c.Request.Context(), payload); err != nil {
		JSONError(c, http.StatusInternalServerError, "CLIENT_META_SAVE_FAILED", "failed to save portal client metadata", err.Error())
		return
	}
	JSONSuccess(c, http.StatusOK, payload)
}

// GetUser returns a projected user.
func (h *AdminHandler) GetUser(c *gin.Context) {
	session := middleware.CurrentSession(c)
	userID := c.Param("userId")
	user, err := h.repos.Users.GetByRealmAndUserID(c.Request.Context(), session.RealmID, userID)
	if err != nil {
		JSONError(c, http.StatusNotFound, "USER_NOT_FOUND", "failed to load projected user", err.Error())
		return
	}
	JSONSuccess(c, http.StatusOK, user)
}

// GetSessionSettings returns global session settings.
func (h *AdminHandler) GetSessionSettings(c *gin.Context) {
	settings, err := h.repos.Settings.GetGlobal(c.Request.Context(), h.defaultIdleTimeout)
	if err != nil {
		JSONError(c, http.StatusInternalServerError, "SETTINGS_LOOKUP_FAILED", "failed to load session settings", err.Error())
		return
	}
	JSONSuccess(c, http.StatusOK, settings)
}

// UpdateSessionSettings updates global session settings.
func (h *AdminHandler) UpdateSessionSettings(c *gin.Context) {
	var payload model.PortalSettings
	if err := c.ShouldBindJSON(&payload); err != nil {
		JSONError(c, http.StatusBadRequest, "INVALID_REQUEST", "invalid session settings payload", err.Error())
		return
	}
	if payload.IdleTimeoutMinutes <= 0 {
		payload.IdleTimeoutMinutes = h.defaultIdleTimeout
	}
	if payload.IdleWarnSeconds <= 0 {
		payload.IdleWarnSeconds = 60
	}

	if err := h.repos.Settings.UpsertGlobal(c.Request.Context(), payload); err != nil {
		JSONError(c, http.StatusInternalServerError, "SETTINGS_SAVE_FAILED", "failed to save session settings", err.Error())
		return
	}
	JSONSuccess(c, http.StatusOK, payload)
}

// SyncStatus returns the latest sync summary for the current admin session.
func (h *AdminHandler) SyncStatus(c *gin.Context) {
	status, err := h.service.SyncStatus(c.Request.Context(), middleware.CurrentSession(c), h.defaultIdleTimeout)
	if err != nil {
		JSONError(c, http.StatusInternalServerError, "SYNC_STATUS_FAILED", "failed to load sync status", err.Error())
		return
	}
	JSONSuccess(c, http.StatusOK, status)
}
