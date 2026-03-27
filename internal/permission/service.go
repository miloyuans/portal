package permission

import (
	"context"
	"errors"
	"slices"
	"sort"

	"go.mongodb.org/mongo-driver/mongo"

	"portal/internal/model"
	"portal/internal/repository"
)

var (
	// ErrAppNotVisible indicates the current session cannot see the requested app.
	ErrAppNotVisible = errors.New("portal app not visible")
	// ErrLaunchDisabled indicates the app is visible but intentionally not launchable.
	ErrLaunchDisabled = errors.New("portal app launch disabled")
	// ErrLaunchTargetMissing indicates the app is visible but has no launch target.
	ErrLaunchTargetMissing = errors.New("portal app launch target missing")
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
		if !ok {
			continue
		}

		view, visible := s.buildAppView(client, meta, session)
		if !visible {
			continue
		}

		sortable = append(sortable, sortableApp{
			sort: meta.Sort,
			view: view,
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

// ResolveLaunch returns the final app launch target using only projected Mongo data.
func (s *Service) ResolveLaunch(ctx context.Context, session model.PortalSession, clientID string) (model.PortalLaunchView, error) {
	client, err := s.repos.Clients.GetByRealmAndClientID(ctx, session.RealmID, clientID)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return model.PortalLaunchView{}, ErrAppNotVisible
		}
		return model.PortalLaunchView{}, err
	}

	meta, err := s.repos.ClientMetas.GetByRealmAndClientID(ctx, session.RealmID, clientID)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return model.PortalLaunchView{}, ErrAppNotVisible
		}
		return model.PortalLaunchView{}, err
	}

	view, visible := s.buildAppView(client, meta, session)
	if !visible {
		return model.PortalLaunchView{}, ErrAppNotVisible
	}
	if !view.CanLaunch {
		if model.NormalizeLaunchMode(meta.LaunchMode) == model.LaunchModeDisabled {
			return model.PortalLaunchView{}, ErrLaunchDisabled
		}
		return model.PortalLaunchView{}, ErrLaunchTargetMissing
	}

	return model.PortalLaunchView{
		ClientID:    view.ClientID,
		DisplayName: view.DisplayName,
		LaunchMode:  view.LaunchMode,
		LaunchURL:   view.LaunchURL,
	}, nil
}

func (s *Service) buildAppView(client model.ClientProjection, meta model.PortalClientMeta, session model.PortalSession) (model.PortalAppView, bool) {
	if !s.canView(client, meta, session) {
		return model.PortalAppView{}, false
	}

	launchMode := model.NormalizeLaunchMode(meta.LaunchMode)
	launchURL := resolveLaunchURL(client, meta)
	canLaunch := launchMode != model.LaunchModeDisabled && launchURL != ""

	return model.PortalAppView{
		ClientID:    client.ClientID,
		DisplayName: firstNonEmpty(meta.DisplayName, client.Name, client.ClientID),
		Category:    meta.Category,
		Icon:        meta.Icon,
		LaunchMode:  launchMode,
		LaunchURL:   launchURL,
		CanView:     true,
		CanLaunch:   canLaunch,
		CanAdmin:    s.canAdmin(meta, session),
	}, true
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

func resolveLaunchURL(client model.ClientProjection, meta model.PortalClientMeta) string {
	switch model.NormalizeLaunchMode(meta.LaunchMode) {
	case model.LaunchModeDisabled:
		return ""
	case model.LaunchModeDirect, model.LaunchModeSPInitiated:
		return firstNonEmpty(meta.LaunchURL, meta.LaunchConfig["launchUrl"], client.BaseURL, client.RootURL)
	default:
		return ""
	}
}
