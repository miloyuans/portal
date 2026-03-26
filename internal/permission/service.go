package permission

import (
	"context"
	"slices"
	"sort"

	"portal/internal/model"
	"portal/internal/repository"
)

// Service resolves portal visibility and admin access rules.
type Service struct {
	repos *repository.Repositories
}

// NewService creates a PermissionService.
func NewService(repos *repository.Repositories) *Service {
	return &Service{repos: repos}
}

// IsPortalAdmin reports whether the current session has portal admin access.
func (s *Service) IsPortalAdmin(session model.PortalSession) bool {
	return slices.Contains(session.RealmRoles, "portal_admin")
}

// ResolveApps returns all visible portal apps for the current user.
func (s *Service) ResolveApps(ctx context.Context, session model.PortalSession) ([]model.PortalAppView, error) {
	clients, err := s.repos.Clients.ListByRealm(ctx, session.RealmID)
	if err != nil {
		return nil, err
	}
	metas, err := s.repos.ClientMetas.ListByRealm(ctx, session.RealmID)
	if err != nil {
		return nil, err
	}

	metaByClientID := make(map[string]model.PortalClientMeta, len(metas))
	for _, meta := range metas {
		metaByClientID[meta.ClientID] = meta
	}

	type sortableApp struct {
		view model.PortalAppView
		sort int
	}

	var sortable []sortableApp
	for _, client := range clients {
		meta, ok := metaByClientID[client.ClientID]
		if !ok || !client.Enabled || !meta.Visible {
			continue
		}

		canView := s.canView(client, meta, session)
		if !canView {
			continue
		}

		sortable = append(sortable, sortableApp{
			sort: meta.Sort,
			view: model.PortalAppView{
				ClientID:    client.ClientID,
				DisplayName: firstNonEmpty(meta.DisplayName, client.Name, client.ClientID),
				Category:    meta.Category,
				Icon:        meta.Icon,
				LaunchURL:   firstNonEmpty(meta.LaunchURL, client.BaseURL, client.RootURL),
				CanView:     true,
				CanLaunch:   true,
				CanAdmin:    s.canAdmin(meta, session),
			},
		})
	}

	sort.Slice(sortable, func(i, j int) bool {
		if sortable[i].sort != sortable[j].sort {
			return sortable[i].sort < sortable[j].sort
		}
		if sortable[i].view.DisplayName == sortable[j].view.DisplayName {
			return sortable[i].view.ClientID < sortable[j].view.ClientID
		}
		return sortable[i].view.DisplayName < sortable[j].view.DisplayName
	})

	apps := make([]model.PortalAppView, 0, len(sortable))
	for _, item := range sortable {
		apps = append(apps, item.view)
	}
	return apps, nil
}

func (s *Service) canView(client model.ClientProjection, meta model.PortalClientMeta, session model.PortalSession) bool {
	if !client.Enabled || !meta.Visible {
		return false
	}
	if len(meta.AccessRules.AnyRealmRoles) == 0 && len(meta.AccessRules.AnyClientRoles) == 0 {
		return s.IsPortalAdmin(session)
	}
	if hasIntersection(session.RealmRoles, meta.AccessRules.AnyRealmRoles) {
		return true
	}
	return hasIntersection(session.ClientRoles[client.ClientID], meta.AccessRules.AnyClientRoles)
}

func (s *Service) canAdmin(meta model.PortalClientMeta, session model.PortalSession) bool {
	if s.IsPortalAdmin(session) {
		return true
	}
	return hasIntersection(session.RealmRoles, meta.AccessRules.AdminRealmRoles)
}

func hasIntersection(left, right []string) bool {
	if len(left) == 0 || len(right) == 0 {
		return false
	}
	for _, candidate := range left {
		if slices.Contains(right, candidate) {
			return true
		}
	}
	return false
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if value != "" {
			return value
		}
	}
	return ""
}
