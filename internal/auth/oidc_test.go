package auth

import (
	"strings"
	"testing"

	"portal/internal/config"
)

func TestLogoutURLUsesRequestedRedirect(t *testing.T) {
	client := &OIDCClient{
		cfg: config.Config{
			Keycloak: config.KeycloakConfig{
				BaseURL:               "http://localhost:8081",
				Realm:                 "portal",
				OIDCClientID:          "portal-api",
				PostLogoutRedirectURL: "http://localhost:5173/login",
			},
		},
	}

	url := client.LogoutURL("id-token", "http://localhost:5173/session-expired")
	if !strings.Contains(url, "session-expired") {
		t.Fatalf("expected logout url to include custom redirect, got %s", url)
	}
}
