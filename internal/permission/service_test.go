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
