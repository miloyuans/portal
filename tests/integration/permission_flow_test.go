//go:build integration

package integration

import (
	"os"
	"testing"
)

func TestPermissionProjectionSkeleton(t *testing.T) {
	baseURL := os.Getenv("PORTAL_BASE_URL")
	if baseURL == "" {
		t.Skip("set PORTAL_BASE_URL to run the permission projection integration skeleton")
	}

	t.Skip("integration skeleton: login as seeded users, call /api/v1/apps, and assert client visibility matches synced Keycloak realm/client roles plus portal_client_meta rules")
}
