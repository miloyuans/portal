package session

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/mongo"

	"portal/internal/config"
	"portal/internal/model"
	"portal/internal/repository"
)

var (
	// ErrSessionExpired indicates the portal session has expired.
	ErrSessionExpired = errors.New("session expired")
	// ErrInvalidSessionCookie indicates the session cookie signature is invalid.
	ErrInvalidSessionCookie = errors.New("invalid session cookie")
)

// Manager manages portal sessions and signed session cookies.
type Manager struct {
	repo *repository.SessionRepository
	cfg  config.Config
}

// NewManager creates a session manager.
func NewManager(repo *repository.SessionRepository, cfg config.Config) *Manager {
	return &Manager{
		repo: repo,
		cfg:  cfg,
	}
}

// Create creates and stores a new portal session.
func (m *Manager) Create(ctx context.Context, session model.PortalSession, idleTimeoutMinutes int) (model.PortalSession, error) {
	now := time.Now().UTC()
	session.SessionID = uuid.NewString()
	session.CreatedAt = now
	session.LastActiveAt = now
	session.AbsoluteExpiresAt = now.Add(time.Duration(m.cfg.Session.AbsoluteTimeoutMinutes) * time.Minute)
	if idleTimeoutMinutes <= 0 {
		idleTimeoutMinutes = m.cfg.Session.IdleTimeoutMinutes
	}
	session.IdleTimeoutMinutes = idleTimeoutMinutes
	session.ExpiresAt = now.Add(time.Duration(idleTimeoutMinutes) * time.Minute)
	if session.ExpiresAt.After(session.AbsoluteExpiresAt) {
		session.ExpiresAt = session.AbsoluteExpiresAt
	}

	if err := m.repo.Create(ctx, session); err != nil {
		return model.PortalSession{}, err
	}
	return session, nil
}

// SetCookie writes the signed session cookie.
func (m *Manager) SetCookie(c *gin.Context, sessionID string, expiresAt time.Time) {
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     m.cfg.Session.CookieName,
		Value:    m.signSessionID(sessionID),
		Path:     "/",
		HttpOnly: m.cfg.Session.HTTPOnly,
		Secure:   m.cfg.Session.Secure,
		SameSite: sameSite(m.cfg.Session.SameSite),
		Expires:  expiresAt,
	})
}

// ClearCookie removes the session cookie.
func (m *Manager) ClearCookie(c *gin.Context) {
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     m.cfg.Session.CookieName,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: m.cfg.Session.HTTPOnly,
		Secure:   m.cfg.Session.Secure,
		SameSite: sameSite(m.cfg.Session.SameSite),
		Expires:  time.Unix(0, 0),
	})
}

// GetByRequest loads and validates the signed session cookie, then loads the session document.
func (m *Manager) GetByRequest(ctx context.Context, request *http.Request) (model.PortalSession, error) {
	cookie, err := request.Cookie(m.cfg.Session.CookieName)
	if err != nil {
		return model.PortalSession{}, err
	}

	sessionID, err := m.verifyCookieValue(cookie.Value)
	if err != nil {
		return model.PortalSession{}, err
	}
	return m.repo.GetByID(ctx, sessionID)
}

// Validate validates the portal session against idle and absolute timeout.
func (m *Manager) Validate(session model.PortalSession) error {
	now := time.Now().UTC()
	if now.After(session.ExpiresAt) || now.After(session.AbsoluteExpiresAt) {
		return ErrSessionExpired
	}
	return nil
}

// Touch refreshes the session idle deadline.
func (m *Manager) Touch(ctx context.Context, session model.PortalSession, idleTimeoutMinutes int) (model.PortalSession, error) {
	now := time.Now().UTC()
	expiresAt := now.Add(time.Duration(idleTimeoutMinutes) * time.Minute)
	if expiresAt.After(session.AbsoluteExpiresAt) {
		expiresAt = session.AbsoluteExpiresAt
	}
	if err := m.repo.Touch(ctx, session.SessionID, now, expiresAt); err != nil {
		return model.PortalSession{}, err
	}
	session.IdleTimeoutMinutes = idleTimeoutMinutes
	session.LastActiveAt = now
	session.ExpiresAt = expiresAt
	return session, nil
}

// Delete removes a session by session ID.
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

func (m *Manager) signSessionID(sessionID string) string {
	mac := hmac.New(sha256.New, []byte(m.cfg.Session.SigningKey))
	mac.Write([]byte(sessionID))
	signature := base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
	return sessionID + "." + signature
}

func (m *Manager) verifyCookieValue(value string) (string, error) {
	parts := strings.Split(value, ".")
	if len(parts) != 2 {
		return "", ErrInvalidSessionCookie
	}

	sessionID := parts[0]
	expected := m.signSessionID(sessionID)
	if !hmac.Equal([]byte(expected), []byte(value)) {
		return "", ErrInvalidSessionCookie
	}
	return sessionID, nil
}

func sameSite(value string) http.SameSite {
	switch strings.ToLower(value) {
	case "strict":
		return http.SameSiteStrictMode
	case "none":
		return http.SameSiteNoneMode
	default:
		return http.SameSiteLaxMode
	}
}
