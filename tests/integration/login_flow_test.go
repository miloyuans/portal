//go:build integration

package integration

import (
	"os"
	"testing"
)

func TestLoginFlowSkeleton(t *testing.T) {
	baseURL := os.Getenv("PORTAL_BASE_URL")
	if baseURL == "" {
		t.Skip("set PORTAL_BASE_URL to run the login flow integration skeleton")
	}

	t.Skip("integration skeleton: start docker compose, navigate browser or automation through OIDC login, then assert /api/v1/me returns the synced session profile")
}
