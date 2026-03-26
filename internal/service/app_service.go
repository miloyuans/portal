package service

import (
	"context"

	"portal/internal/config"
	"portal/internal/model"
	"portal/internal/permission"
	"portal/internal/repository"
)

type AppService struct {
	permissions *permission.Service
	repos       *repository.Repositories
	cfg         config.Config
}

func NewAppService(permissions *permission.Service, repos *repository.Repositories, cfg config.Config) *AppService {
	return &AppService{
		permissions: permissions,
		repos:       repos,
		cfg:         cfg,
	}
}

func (s *AppService) Me(ctx context.Context, session model.PortalSession) (model.CurrentUserProfile, error) {
	settings, err := s.repos.Settings.GetByRealm(ctx, session.Realm, s.cfg.Session.DefaultIdleTimeoutMinutes)
	if err != nil {
		return model.CurrentUserProfile{}, err
	}
	return model.CurrentUserProfile{
		Realm:    session.Realm,
		User:     session.View(),
		IsAdmin:  s.permissions.IsAdmin(session),
		Settings: settings,
	}, nil
}

func (s *AppService) Apps(ctx context.Context, session model.PortalSession) ([]model.PortalApp, error) {
	return s.permissions.BuildVisibleApps(ctx, session)
}
