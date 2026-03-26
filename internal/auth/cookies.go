package auth

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"portal/internal/config"
)

func NewStateValue() (string, error) {
	buffer := make([]byte, 32)
	if _, err := rand.Read(buffer); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(buffer), nil
}

func SetTransientCookie(c *gin.Context, cfg config.Config, name, value string, maxAgeSeconds int) {
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     name,
		Value:    value,
		Path:     "/",
		MaxAge:   maxAgeSeconds,
		HttpOnly: true,
		Secure:   cfg.Server.CookieSecure,
		Domain:   cfg.Server.CookieDomain,
		SameSite: sameSite(cfg.Server.CookieSameSite),
		Expires:  time.Now().Add(time.Duration(maxAgeSeconds) * time.Second),
	})
}

func ClearCookie(c *gin.Context, cfg config.Config, name string) {
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     name,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   cfg.Server.CookieSecure,
		Domain:   cfg.Server.CookieDomain,
		SameSite: sameSite(cfg.Server.CookieSameSite),
		Expires:  time.Unix(0, 0),
	})
}

func sameSite(value string) http.SameSite {
	switch value {
	case "Strict":
		return http.SameSiteStrictMode
	case "None":
		return http.SameSiteNoneMode
	default:
		return http.SameSiteLaxMode
	}
}
