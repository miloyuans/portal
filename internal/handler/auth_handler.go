package handler

import (
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"

	"portal/internal/auth"
	"portal/internal/config"
	"portal/internal/model"
	"portal/internal/permission"
	"portal/internal/repository"
	sessionpkg "portal/internal/session"
	syncsvc "portal/internal/sync"
)

type AuthHandler struct {
	cfg         config.Config
	oidc        *auth.OIDCClient
	syncService *syncsvc.Service
	sessions    *sessionpkg.Manager
	repos       *repository.Repositories
	permissions *permission.Service
	logger      *slog.Logger
}

func NewAuthHandler(cfg config.Config, oidc *auth.OIDCClient, syncService *syncsvc.Service, sessions *sessionpkg.Manager, repos *repository.Repositories, permissions *permission.Service, logger *slog.Logger) *AuthHandler {
	return &AuthHandler{
		cfg:         cfg,
		oidc:        oidc,
		syncService: syncService,
		sessions:    sessions,
		repos:       repos,
		permissions: permissions,
		logger:      logger,
	}
}

func (h *AuthHandler) Login(c *gin.Context) {
	state, err := auth.NewStateValue()
	if err != nil {
		JSONError(c, http.StatusInternalServerError, "STATE_GENERATION_FAILED", "failed to generate oidc state", err.Error())
		return
	}
	nonce, err := auth.NewStateValue()
	if err != nil {
		JSONError(c, http.StatusInternalServerError, "NONCE_GENERATION_FAILED", "failed to generate oidc nonce", err.Error())
		return
	}

	auth.SetTransientCookie(c, h.cfg, h.cfg.Session.StateCookieName, state, h.cfg.Session.StateCookieMaxAgeSeconds)
	auth.SetTransientCookie(c, h.cfg, h.cfg.Session.NonceCookieName, nonce, h.cfg.Session.StateCookieMaxAgeSeconds)
	c.Redirect(http.StatusFound, h.oidc.AuthCodeURL(state, nonce))
}

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

	result, err := h.syncService.SyncCurrentUser(c.Request.Context(), tokenBundle.Claims.Subject)
	if err != nil {
		JSONError(c, http.StatusBadGateway, "SYNC_FAILED", "failed to sync current user from keycloak", err.Error())
		return
	}

	settings, err := h.repos.Settings.GetByRealm(c.Request.Context(), result.Realm.Realm, h.cfg.Session.DefaultIdleTimeoutMinutes)
	if err != nil {
		JSONError(c, http.StatusInternalServerError, "SETTINGS_LOOKUP_FAILED", "failed to load portal settings", err.Error())
		return
	}

	displayName := strings.TrimSpace(result.User.FirstName + " " + result.User.LastName)
	if displayName == "" {
		displayName = result.User.Username
	}

	sessionRecord, err := h.sessions.Create(c.Request.Context(), model.PortalSession{
		Realm:              result.Realm.Realm,
		UserID:             result.User.UserID,
		Username:           result.User.Username,
		Email:              result.User.Email,
		DisplayName:        displayName,
		RealmRoles:         result.User.RealmRoles,
		ClientRoles:        result.User.ClientRoles,
		AccessToken:        tokenBundle.AccessToken,
		RefreshToken:       tokenBundle.RefreshToken,
		IDToken:            tokenBundle.IDToken,
		IdleTimeoutMinutes: settings.IdleTimeoutMinutes,
		ExpiresAt:          time.Now().UTC().Add(time.Duration(h.cfg.Session.AbsoluteTTLHours) * time.Hour),
	})
	if err != nil {
		JSONError(c, http.StatusInternalServerError, "SESSION_CREATE_FAILED", "failed to create portal session", err.Error())
		return
	}

	h.sessions.SetCookie(c, sessionRecord.SessionID)
	h.logger.Info("portal session established",
		slog.String("realm", sessionRecord.Realm),
		slog.String("userId", sessionRecord.UserID),
		slog.Bool("isAdmin", h.permissions.IsAdmin(sessionRecord)),
	)

	c.Redirect(http.StatusFound, h.cfg.Server.PublicWebURL+"/auth/callback/success")
}

func (h *AuthHandler) Logout(c *gin.Context) {
	session, err := h.sessions.GetByRequest(c.Request.Context(), c.Request)
	if err == nil {
		if deleteErr := h.sessions.Delete(c.Request.Context(), session.SessionID); deleteErr != nil && deleteErr != mongo.ErrNoDocuments {
			JSONError(c, http.StatusInternalServerError, "SESSION_DELETE_FAILED", "failed to delete portal session", deleteErr.Error())
			return
		}
	}

	h.sessions.ClearCookie(c)
	idTokenHint := ""
	if err == nil {
		idTokenHint = session.IDToken
	}

	redirectURI := h.cfg.Server.PublicWebURL + "/login"
	if c.Query("reason") == "expired" {
		redirectURI = h.cfg.Server.PublicWebURL + "/session-expired"
	}
	c.Redirect(http.StatusFound, h.oidc.LogoutURL(idTokenHint, redirectURI))
}
