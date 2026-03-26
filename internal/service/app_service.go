package service

import (
	"context"

	"portal/internal/model"
	"portal/internal/permission"
	"portal/internal/repository"
)

// AppService serves portal-facing application data.
type AppService struct {
	permissions *permission.Service
	repos       *repository.Repositories
}

// NewAppService creates an AppService.
func NewAppService(permissions *permission.Service, repos *repository.Repositories) *AppService {
	return &AppService{
		permissions: permissions,
		repos:       repos,
	}
}

// Me returns the current session summary.
func (s *AppService) Me(session model.PortalSession) model.SessionView {
	return session.View()
}

// Apps resolves the visible portal apps for the current session.
func (s *AppService) Apps(ctx context.Context, session model.PortalSession) ([]model.PortalAppView, error) {
	return s.permissions.ResolveApps(ctx, session)
}

// Profile returns the current user's projected profile.
func (s *AppService) Profile(ctx context.Context, session model.PortalSession, defaultIdleTimeout int) (model.CurrentUserProfile, error) {
	user, err := s.repos.Users.GetByRealmAndUserID(ctx, session.RealmID, session.UserID)
	if err != nil {
		return model.CurrentUserProfile{}, err
	}
	realm, err := s.repos.Realms.GetByRealmID(ctx, session.RealmID)
	if err != nil {
		return model.CurrentUserProfile{}, err
	}
	settings, err := s.repos.Settings.GetGlobal(ctx, defaultIdleTimeout)
	if err != nil {
		return model.CurrentUserProfile{}, err
	}
	return model.CurrentUserProfile{
		Session:  session.View(),
		User:     user,
		Realm:    realm,
		Settings: settings,
	}, nil
}

// Realms returns the projected realm list.
func (s *AppService) Realms(ctx context.Context) ([]model.RealmProjection, error) {
	return s.repos.Realms.List(ctx)
}

// SyncStatus returns a compact sync summary for the current session.
func (s *AppService) SyncStatus(ctx context.Context, session model.PortalSession, defaultIdleTimeout int) (model.SyncStatus, error) {
	realm, err := s.repos.Realms.GetByRealmID(ctx, session.RealmID)
	if err != nil {
		return model.SyncStatus{}, err
	}
	user, err := s.repos.Users.GetByRealmAndUserID(ctx, session.RealmID, session.UserID)
	if err != nil {
		return model.SyncStatus{}, err
	}
	clients, err := s.repos.Clients.ListByRealm(ctx, session.RealmID)
	if err != nil {
		return model.SyncStatus{}, err
	}
	settings, err := s.repos.Settings.GetGlobal(ctx, defaultIdleTimeout)
	if err != nil {
		return model.SyncStatus{}, err
	}
	return model.SyncStatus{
		RealmID:           realm.RealmID,
		RealmSyncedAt:     realm.SyncedAt,
		UserSyncedAt:      user.SyncedAt,
		ClientCount:       len(clients),
		SettingsUpdatedAt: settings.UpdatedAt,
	}, nil
}
