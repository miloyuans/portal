package permission

import (
	"context"
	"slices"
	"sort"
	"strings"

	"portal/internal/config"
	"portal/internal/model"
	"portal/internal/repository"
)

type Service struct {
	repos *repository.Repositories
	cfg   config.Config
}

func NewService(repos *repository.Repositories, cfg config.Config) *Service {
	return &Service{
		repos: repos,
		cfg:   cfg,
	}
}

func (s *Service) IsAdmin(session model.PortalSession) bool {
	for _, role := range session.RealmRoles {
		if slices.Contains(s.cfg.Permission.AdminRealmRoles, role) {
			return true
		}
	}
	return false
}

func (s *Service) BuildVisibleApps(ctx context.Context, session model.PortalSession) ([]model.PortalApp, error) {
	clients, err := s.repos.Clients.ListByRealm(ctx, session.Realm)
	if err != nil {
		return nil, err
	}
	metas, err := s.repos.ClientMetas.ListByRealm(ctx, session.Realm)
	if err != nil {
		return nil, err
	}

	metaByClientID := make(map[string]model.PortalClientMeta, len(metas))
	for _, meta := range metas {
		metaByClientID[meta.ClientID] = meta
	}

	isAdmin := s.IsAdmin(session)
	apps := make([]model.PortalApp, 0)
	for _, client := range clients {
		meta, hasMeta := metaByClientID[client.ClientID]
		targetURL := firstNonEmpty(meta.TargetURL, client.RootURL, client.BaseURL)
		if !client.Enabled || strings.TrimSpace(targetURL) == "" {
			continue
		}

		enabled := client.Enabled
		showInPortal := true
		if hasMeta {
			enabled = meta.Enabled
			showInPortal = meta.ShowInPortal
		}
		if !enabled || !showInPortal {
			continue
		}

		if !isAdmin && !s.clientVisibleToUser(client, meta, hasMeta, session) {
			continue
		}

		apps = append(apps, model.PortalApp{
			ClientID:    client.ClientID,
			DisplayName: firstNonEmpty(meta.DisplayName, client.Name, client.ClientID),
			Description: firstNonEmpty(meta.Description, client.Description),
			TargetURL:   targetURL,
			Icon:        meta.Icon,
			Category:    meta.Category,
			Tags:        meta.Tags,
			SortOrder:   meta.SortOrder,
		})
	}

	sort.Slice(apps, func(i, j int) bool {
		if apps[i].SortOrder == apps[j].SortOrder {
			return apps[i].DisplayName < apps[j].DisplayName
		}
		return apps[i].SortOrder < apps[j].SortOrder
	})
	return apps, nil
}

func (s *Service) clientVisibleToUser(client model.ClientProjection, meta model.PortalClientMeta, hasMeta bool, session model.PortalSession) bool {
	userClientRoles := session.ClientRoles[client.ClientID]
	if hasMeta && (len(meta.RequiredRealmRoles) > 0 || len(meta.RequiredClientRoles) > 0) {
		return hasIntersection(session.RealmRoles, meta.RequiredRealmRoles) || hasIntersection(userClientRoles, meta.RequiredClientRoles)
	}
	return len(userClientRoles) > 0
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
		if strings.TrimSpace(value) != "" {
			return value
		}
	}
	return ""
}
