package permission

import (
	"testing"

	"portal/internal/config"
	"portal/internal/model"
)

func TestIsAdmin(t *testing.T) {
	service := NewService(nil, config.Config{
		Permission: config.PermissionConfig{
			AdminRealmRoles: []string{"portal-admin"},
		},
	})

	if !service.IsAdmin(model.PortalSession{RealmRoles: []string{"viewer", "portal-admin"}}) {
		t.Fatalf("expected portal-admin realm role to grant admin access")
	}
}

func TestClientVisibleToUser(t *testing.T) {
	service := NewService(nil, config.Config{})
	client := model.ClientProjection{ClientID: "sales-app"}
	meta := model.PortalClientMeta{
		ClientID:            "sales-app",
		RequiredRealmRoles:  []string{"finance-admin"},
		RequiredClientRoles: []string{"viewer"},
	}

	session := model.PortalSession{
		RealmRoles: []string{"employee"},
		ClientRoles: map[string][]string{
			"sales-app": {"viewer"},
		},
	}

	if !service.clientVisibleToUser(client, meta, true, session) {
		t.Fatalf("expected client role intersection to grant visibility")
	}
}
