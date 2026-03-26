package session

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/mongo"

	"portal/internal/auth"
	"portal/internal/config"
	"portal/internal/model"
	"portal/internal/repository"
)

var ErrSessionExpired = errors.New("session expired")

type Manager struct {
	repo *repository.SessionRepository
	cfg  config.Config
}

func NewManager(repo *repository.SessionRepository, _ any, cfg config.Config) *Manager {
	return &Manager{
		repo: repo,
		cfg:  cfg,
	}
}

func (m *Manager) Create(ctx context.Context, session model.PortalSession) (model.PortalSession, error) {
	now := time.Now().UTC()
	session.SessionID = uuid.NewString()
	session.CreatedAt = now
	session.UpdatedAt = now
	session.LastSeenAt = now
	if session.IdleTimeoutMinutes <= 0 {
		session.IdleTimeoutMinutes = m.cfg.Session.DefaultIdleTimeoutMinutes
	}
	if session.ExpiresAt.IsZero() {
		session.ExpiresAt = now.Add(time.Duration(m.cfg.Session.AbsoluteTTLHours) * time.Hour)
	}
	if err := m.repo.Create(ctx, session); err != nil {
		return model.PortalSession{}, err
	}
	return session, nil
}

func (m *Manager) SetCookie(c *gin.Context, sessionID string) {
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     m.cfg.Server.CookieName,
		Value:    sessionID,
		Path:     "/",
		HttpOnly: true,
		Secure:   m.cfg.Server.CookieSecure,
		Domain:   m.cfg.Server.CookieDomain,
		SameSite: authCookieSameSite(m.cfg.Server.CookieSameSite),
		Expires:  time.Now().Add(time.Duration(m.cfg.Session.AbsoluteTTLHours) * time.Hour),
	})
}

func (m *Manager) ClearCookie(c *gin.Context) {
	auth.ClearCookie(c, m.cfg, m.cfg.Server.CookieName)
}

func (m *Manager) GetByRequest(ctx context.Context, request *http.Request) (model.PortalSession, error) {
	cookie, err := request.Cookie(m.cfg.Server.CookieName)
	if err != nil {
		return model.PortalSession{}, err
	}
	return m.repo.GetByID(ctx, cookie.Value)
}

func (m *Manager) Validate(session model.PortalSession) error {
	now := time.Now().UTC()
	if now.After(session.ExpiresAt) {
		return ErrSessionExpired
	}
	if now.After(session.LastSeenAt.Add(time.Duration(session.IdleTimeoutMinutes) * time.Minute)) {
		return ErrSessionExpired
	}
	return nil
}

func (m *Manager) Touch(ctx context.Context, session model.PortalSession) error {
	return m.repo.Touch(ctx, session.SessionID, time.Now().UTC(), session.ExpiresAt)
}

func (m *Manager) Delete(ctx context.Context, sessionID string) error {
	if sessionID == "" {
		return nil
	}
	err := m.repo.Delete(ctx, sessionID)
	if err == mongo.ErrNoDocuments {
		return nil
	}
	return err
}

func authCookieSameSite(value string) http.SameSite {
	switch value {
	case "Strict":
		return http.SameSiteStrictMode
	case "None":
		return http.SameSiteNoneMode
	default:
		return http.SameSiteLaxMode
	}
}
