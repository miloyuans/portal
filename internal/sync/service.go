package syncsvc

import (
	"context"
	"log/slog"
	"time"

	"portal/internal/kcadmin"
	"portal/internal/model"
	"portal/internal/repository"
)

type Service struct {
	kc     *kcadmin.Client
	repos  *repository.Repositories
	logger *slog.Logger
}

type Result struct {
	Realm   model.RealmProjection
	Clients []model.ClientProjection
	User    model.UserProjection
}

func NewService(kc *kcadmin.Client, repos *repository.Repositories, logger *slog.Logger) *Service {
	return &Service{
		kc:     kc,
		repos:  repos,
		logger: logger,
	}
}

func (s *Service) SyncCurrentUser(ctx context.Context, userID string) (Result, error) {
	realm, clients, user, err := s.kc.SyncData(ctx, userID)
	if err != nil {
		return Result{}, err
	}

	if err := s.repos.Realms.Upsert(ctx, realm); err != nil {
		return Result{}, err
	}
	if err := s.repos.Clients.UpsertMany(ctx, clients); err != nil {
		return Result{}, err
	}
	defaultMetas := make([]model.PortalClientMeta, 0, len(clients))
	now := time.Now().UTC()
	for index, client := range clients {
		defaultMetas = append(defaultMetas, model.PortalClientMeta{
			Realm:        client.Realm,
			ClientID:     client.ClientID,
			DisplayName:  firstNonEmpty(client.Name, client.ClientID),
			Description:  client.Description,
			TargetURL:    firstNonEmpty(client.RootURL, client.BaseURL),
			SortOrder:    (index + 1) * 10,
			Enabled:      client.Enabled,
			ShowInPortal: true,
			CreatedAt:    now,
			UpdatedAt:    now,
		})
	}
	if err := s.repos.ClientMetas.SeedDefaults(ctx, defaultMetas); err != nil {
		return Result{}, err
	}
	if err := s.repos.Users.Upsert(ctx, user); err != nil {
		return Result{}, err
	}

	s.logger.Info("user sync completed",
		slog.String("realm", realm.Realm),
		slog.String("userId", user.UserID),
		slog.Int("clients", len(clients)),
		slog.Int("realmRoles", len(user.RealmRoles)),
	)

	return Result{
		Realm:   realm,
		Clients: clients,
		User:    user,
	}, nil
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if value != "" {
			return value
		}
	}
	return ""
}
