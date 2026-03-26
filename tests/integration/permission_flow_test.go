//go:build integration

package integration

import "testing"

func TestPortalAppsReflectAccessRulesSkeleton(t *testing.T) {
	t.Skip("integration skeleton: seed portal_client_meta, login as multiple users, assert /api/portal/apps returns the expected canView/canLaunch/canAdmin set")
}
