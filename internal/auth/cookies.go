package auth

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"portal/internal/config"
)

// NewStateValue returns a random URL-safe state or nonce value.
func NewStateValue() (string, error) {
	buffer := make([]byte, 32)
	if _, err := rand.Read(buffer); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(buffer), nil
}

// SetTransientCookie writes an OIDC state/nonce cookie.
func SetTransientCookie(c *gin.Context, cfg config.Config, name, value string, maxAge time.Duration) {
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     name,
		Value:    value,
		Path:     "/",
		MaxAge:   int(maxAge.Seconds()),
		HttpOnly: cfg.Session.HTTPOnly,
		Secure:   cfg.Session.Secure,
		SameSite: sameSite(cfg.Session.SameSite),
		Expires:  time.Now().Add(maxAge),
	})
}

// ClearCookie clears an OIDC cookie.
func ClearCookie(c *gin.Context, cfg config.Config, name string) {
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     name,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: cfg.Session.HTTPOnly,
		Secure:   cfg.Session.Secure,
		SameSite: sameSite(cfg.Session.SameSite),
		Expires:  time.Unix(0, 0),
	})
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
