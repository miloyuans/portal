//go:build integration

package integration

import "testing"

func TestOIDCLoginCallbackFlowSkeleton(t *testing.T) {
	t.Skip("integration skeleton: stand up Keycloak, drive /api/auth/login -> /api/auth/callback, assert portal session cookie and redirect to /portal")
}

func TestLoginSyncsCurrentUserSnapshotSkeleton(t *testing.T) {
	t.Skip("integration skeleton: after callback, assert kc_realms, kc_clients and kc_users are upserted for the current user only")
}
