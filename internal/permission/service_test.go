package permission

import (
	"testing"

	"portal/internal/model"
)

func TestIsPortalAdmin(t *testing.T) {
	service := NewService(nil)
	if !service.IsPortalAdmin(model.PortalSession{RealmRoles: []string{"portal_admin"}}) {
		t.Fatalf("expected portal_admin role to grant admin access")
	}
}

func TestCanViewByClientRoles(t *testing.T) {
	service := NewService(nil)
	client := model.ClientProjection{ClientID: "sales-app", Enabled: true}
	meta := model.PortalClientMeta{
		ClientID: "sales-app",
		Visible:  true,
		AccessRules: model.AccessRules{
			AnyClientRoles: []string{"viewer"},
		},
	}
	session := model.PortalSession{
		ClientRoles: map[string][]string{
			"sales-app": {"viewer"},
		},
	}

	if !service.canView(client, meta, session) {
		t.Fatalf("expected matching client role to grant visibility")
	}
}

func TestBuildAppViewDefaultsToSPInitiated(t *testing.T) {
	service := NewService(nil)
	client := model.ClientProjection{
		ClientID: "finance-app",
		Name:     "Finance",
		Enabled:  true,
		BaseURL:  "https://finance.example.com",
	}
	meta := model.PortalClientMeta{
		ClientID:    "finance-app",
		DisplayName: "Finance Portal",
		Visible:     true,
	}
	session := model.PortalSession{
		RealmRoles: []string{"portal_admin"},
	}

	view, visible := service.buildAppView(client, meta, session)
	if !visible {
		t.Fatalf("expected app to be visible for portal admin")
	}
	if view.LaunchMode != model.LaunchModeSPInitiated {
		t.Fatalf("expected default launch mode to be sp_initiated, got %q", view.LaunchMode)
	}
	if !view.CanLaunch {
		t.Fatalf("expected default launch mode with base URL to be launchable")
	}
	if view.LaunchURL != "https://finance.example.com" {
		t.Fatalf("expected base URL to be used as fallback launch target, got %q", view.LaunchURL)
	}
}

func TestBuildAppViewDisabledLaunchMode(t *testing.T) {
	service := NewService(nil)
	client := model.ClientProjection{
		ClientID: "ops-app",
		Enabled:  true,
		BaseURL:  "https://ops.example.com",
	}
	meta := model.PortalClientMeta{
		ClientID:   "ops-app",
		Visible:    true,
		LaunchMode: model.LaunchModeDisabled,
		LaunchURL:  "https://ops.example.com",
	}
	session := model.PortalSession{
		RealmRoles: []string{"portal_admin"},
	}

	view, visible := service.buildAppView(client, meta, session)
	if !visible {
		t.Fatalf("expected app to remain visible for portal admin")
	}
	if view.CanLaunch {
		t.Fatalf("expected disabled launch mode to prevent launching")
	}
	if view.LaunchURL != "" {
		t.Fatalf("expected disabled launch mode to hide the launch URL, got %q", view.LaunchURL)
	}
}
