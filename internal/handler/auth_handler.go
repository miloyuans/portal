package handler

import (
	"log/slog"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"

	"portal/internal/auth"
	"portal/internal/config"
	"portal/internal/model"
	"portal/internal/repository"
	sessionpkg "portal/internal/session"
	syncsvc "portal/internal/sync"
)

// AuthHandler serves OIDC login, callback, logout and auth/me.
type AuthHandler struct {
	cfg      config.Config
	oidc     *auth.OIDCClient
	sync     *syncsvc.Service
	sessions *sessionpkg.Manager
	repos    *repository.Repositories
	logger   *slog.Logger
}

// NewAuthHandler creates an AuthHandler.
func NewAuthHandler(cfg config.Config, oidc *auth.OIDCClient, syncService *syncsvc.Service, sessions *sessionpkg.Manager, repos *repository.Repositories, logger *slog.Logger) *AuthHandler {
	return &AuthHandler{
		cfg:      cfg,
		oidc:     oidc,
		sync:     syncService,
		sessions: sessions,
		repos:    repos,
		logger:   logger,
	}
}

// Login starts the OIDC browser login flow.
func (h *AuthHandler) Login(c *gin.Context) {
	loginURL, err := h.prepareLogin(c)
	if err != nil {
		JSONError(c, http.StatusInternalServerError, "LOGIN_PREPARE_FAILED", "failed to prepare oidc login", err.Error())
		return
	}
	c.Redirect(http.StatusFound, loginURL)
}

// LoginURL prepares the OIDC login flow and returns the Keycloak authorization URL.
func (h *AuthHandler) LoginURL(c *gin.Context) {
	loginURL, err := h.prepareLogin(c)
	if err != nil {
		JSONError(c, http.StatusInternalServerError, "LOGIN_PREPARE_FAILED", "failed to prepare oidc login", err.Error())
		return
	}
	JSONSuccess(c, http.StatusOK, gin.H{
		"loginUrl": loginURL,
	})
}

func (h *AuthHandler) prepareLogin(c *gin.Context) (string, error) {
	state, err := auth.NewStateValue()
	if err != nil {
		return "", err
	}
	nonce, err := auth.NewStateValue()
	if err != nil {
		return "", err
	}

	auth.SetTransientCookie(c, h.cfg, h.cfg.Session.StateCookieName, state, h.cfg.Session.StateCookieMaxAge)
	auth.SetTransientCookie(c, h.cfg, h.cfg.Session.NonceCookieName, nonce, h.cfg.Session.StateCookieMaxAge)
	return h.oidc.AuthCodeURL(state, nonce), nil
}

// Callback handles the OIDC authorization code callback.
func (h *AuthHandler) Callback(c *gin.Context) {
	state := c.Query("state")
	if state == "" {
		JSONError(c, http.StatusBadRequest, "STATE_REQUIRED", "missing oidc state", nil)
		return
	}

	stateCookie, err := c.Cookie(h.cfg.Session.StateCookieName)
	if err != nil || stateCookie != state {
		JSONError(c, http.StatusBadRequest, "STATE_MISMATCH", "oidc state validation failed", nil)
		return
	}

	nonce, err := c.Cookie(h.cfg.Session.NonceCookieName)
	if err != nil {
		JSONError(c, http.StatusBadRequest, "NONCE_MISSING", "oidc nonce validation failed", nil)
		return
	}

	auth.ClearCookie(c, h.cfg, h.cfg.Session.StateCookieName)
	auth.ClearCookie(c, h.cfg, h.cfg.Session.NonceCookieName)

	code := c.Query("code")
	if code == "" {
		JSONError(c, http.StatusBadRequest, "CODE_REQUIRED", "missing authorization code", nil)
		return
	}

	tokenBundle, err := h.oidc.Exchange(c.Request.Context(), code, nonce)
	if err != nil {
		JSONError(c, http.StatusUnauthorized, "TOKEN_EXCHANGE_FAILED", "oidc token exchange failed", err.Error())
		return
	}
	if tokenBundle.Claims.Subject == "" {
		JSONError(c, http.StatusUnauthorized, "SUBJECT_MISSING", "oidc token is missing subject claim", nil)
		return
	}

	syncResult, err := h.sync.SyncCurrentUser(c.Request.Context(), tokenBundle.Claims.Subject)
	if err != nil {
		JSONError(c, http.StatusBadGateway, "SYNC_FAILED", "failed to synchronize current user", err.Error())
		return
	}

	settings, err := h.repos.Settings.GetGlobal(c.Request.Context(), h.cfg.Session.IdleTimeoutMinutes)
	if err != nil {
		JSONError(c, http.StatusInternalServerError, "SETTINGS_LOOKUP_FAILED", "failed to load portal settings", err.Error())
		return
	}

	displayName := strings.TrimSpace(syncResult.User.FirstName + " " + syncResult.User.LastName)
	if displayName == "" {
		displayName = syncResult.User.Username
	}

	sessionRecord, err := h.sessions.Create(c.Request.Context(), model.PortalSession{
		RealmID:     syncResult.Realm.RealmID,
		UserID:      syncResult.User.UserID,
		Username:    syncResult.User.Username,
		DisplayName: displayName,
		RealmRoles:  syncResult.User.RealmRoles,
		ClientRoles: syncResult.User.ClientRoles,
		IDToken:     tokenBundle.IDToken,
	}, settings.IdleTimeoutMinutes)
	if err != nil {
		JSONError(c, http.StatusInternalServerError, "SESSION_CREATE_FAILED", "failed to create portal session", err.Error())
		return
	}

	h.sessions.SetCookie(c, sessionRecord.SessionID, sessionRecord.ExpiresAt)
	h.logger.Info("portal session established",
		slog.String("realmId", sessionRecord.RealmID),
		slog.String("userId", sessionRecord.UserID),
	)

	c.Redirect(http.StatusFound, h.cfg.Server.PublicWebURL+"/portal")
}

// Logout deletes the portal session and returns the Keycloak logout URL.
func (h *AuthHandler) Logout(c *gin.Context) {
	session, err := h.sessions.GetByRequest(c.Request.Context(), c.Request)
	if err == nil {
		if deleteErr := h.sessions.Delete(c.Request.Context(), session.SessionID); deleteErr != nil && deleteErr != mongo.ErrNoDocuments {
			JSONError(c, http.StatusInternalServerError, "SESSION_DELETE_FAILED", "failed to delete portal session", deleteErr.Error())
			return
		}
	}

	h.sessions.ClearCookie(c)
	redirectURI := h.cfg.Server.PublicWebURL + "/login"
	if c.Query("reason") == "expired" {
		redirectURI = h.cfg.Server.PublicWebURL + "/session-expired"
	}

	idTokenHint := ""
	if err == nil {
		idTokenHint = session.IDToken
	}
	JSONSuccess(c, http.StatusOK, gin.H{
		"logoutUrl": h.oidc.LogoutURL(idTokenHint, redirectURI),
	})
}

// Me returns the current session summary.
func (h *AuthHandler) Me(c *gin.Context) {
	sessionValue, ok := c.Get("portalSession")
	if !ok {
		JSONError(c, http.StatusUnauthorized, "AUTH_REQUIRED", "portal session not found", nil)
		return
	}
	session, _ := sessionValue.(model.PortalSession)
	JSONSuccess(c, http.StatusOK, session.View())
}
