package syncsvc

import (
	"context"
	"log/slog"
	"time"

	"portal/internal/config"
	"portal/internal/kcadmin"
	"portal/internal/model"
	"portal/internal/repository"
)

// Service synchronizes the current user's related Keycloak data into MongoDB.
type Service struct {
	admin  *kcadmin.Client
	repos  *repository.Repositories
	cfg    config.Config
	logger *slog.Logger
}

// Result stores the current-user sync result.
type Result struct {
	Realm        model.RealmProjection
	Clients      []model.ClientProjection
	User         model.UserProjection
	RoleMappings kcadmin.MappingsRepresentation
}

// NewService creates a sync service.
func NewService(admin *kcadmin.Client, repos *repository.Repositories, cfg config.Config, logger *slog.Logger) *Service {
	return &Service{
		admin:  admin,
		repos:  repos,
		cfg:    cfg,
		logger: logger,
	}
}

// SyncCurrentUser synchronizes the current realm, clients and current user snapshot.
func (s *Service) SyncCurrentUser(ctx context.Context, userID string) (Result, error) {
	syncCtx, cancel := context.WithTimeout(ctx, time.Duration(s.cfg.Sync.TimeoutSeconds)*time.Second)
	defer cancel()

	realmDTO, err := s.admin.GetRealm(syncCtx, s.cfg.Keycloak.Realm)
	if err != nil {
		return Result{}, err
	}

	clientDTOs, err := s.admin.ListClients(syncCtx, s.cfg.Keycloak.Realm, kcadmin.ListClientsOptions{Max: 500})
	if err != nil {
		return Result{}, err
	}

	userDTO, err := s.admin.GetUserByID(syncCtx, s.cfg.Keycloak.Realm, userID)
	if err != nil {
		return Result{}, err
	}

	roleMappings, err := s.admin.GetUserRoleMappings(syncCtx, s.cfg.Keycloak.Realm, userID)
	if err != nil {
		return Result{}, err
	}

	effectiveRealmRoles, err := s.admin.GetUserEffectiveRealmRoles(syncCtx, s.cfg.Keycloak.Realm, userID, true)
	if err != nil {
		return Result{}, err
	}

	realmID := realmDTO.ID
	if realmID == "" {
		realmID = realmDTO.Realm
	}

	now := time.Now().UTC()
	realmProjection := model.RealmProjection{
		RealmID:     realmID,
		RealmName:   realmDTO.Realm,
		DisplayName: realmDTO.DisplayName,
		Enabled:     realmDTO.Enabled,
		Attributes:  realmDTO.Attributes,
		SyncedAt:    now,
	}

	clientProjections := make([]model.ClientProjection, 0, len(clientDTOs))
	clientRoles := make(map[string][]string)
	for _, clientDTO := range clientDTOs {
		clientProjections = append(clientProjections, model.ClientProjection{
			RealmID:    realmID,
			ClientUUID: clientDTO.ID,
			ClientID:   clientDTO.ClientID,
			Name:       firstNonEmpty(clientDTO.Name, clientDTO.ClientID),
			Enabled:    clientDTO.Enabled,
			BaseURL:    clientDTO.BaseURL,
			RootURL:    clientDTO.RootURL,
			Protocol:   clientDTO.Protocol,
			Attributes: clientDTO.Attributes,
			SyncedAt:   now,
		})

		roles, err := s.admin.GetUserEffectiveClientRoles(syncCtx, s.cfg.Keycloak.Realm, userID, clientDTO.ID)
		if err != nil {
			return Result{}, err
		}
		if len(roles) == 0 {
			continue
		}

		names := make([]string, 0, len(roles))
		for _, role := range roles {
			names = append(names, role.Name)
		}
		clientRoles[clientDTO.ClientID] = names
	}

	realmRoleNames := make([]string, 0, len(effectiveRealmRoles))
	for _, role := range effectiveRealmRoles {
		realmRoleNames = append(realmRoleNames, role.Name)
	}

	userProjection := model.UserProjection{
		RealmID:     realmID,
		UserID:      userDTO.ID,
		Username:    userDTO.Username,
		Email:       userDTO.Email,
		Enabled:     userDTO.Enabled,
		FirstName:   userDTO.FirstName,
		LastName:    userDTO.LastName,
		Attributes:  userDTO.Attributes,
		RealmRoles:  realmRoleNames,
		ClientRoles: clientRoles,
		SyncedAt:    now,
	}

	if err := s.repos.Realms.Upsert(syncCtx, realmProjection); err != nil {
		return Result{}, err
	}
	if err := s.repos.Clients.UpsertMany(syncCtx, clientProjections); err != nil {
		return Result{}, err
	}
	if err := s.repos.Users.Upsert(syncCtx, userProjection); err != nil {
		return Result{}, err
	}

	s.logger.Info("login sync completed",
		slog.String("realmId", realmID),
		slog.String("userId", userProjection.UserID),
		slog.Int("clients", len(clientProjections)),
		slog.Int("realmRoles", len(realmRoleNames)),
	)

	return Result{
		Realm:        realmProjection,
		Clients:      clientProjections,
		User:         userProjection,
		RoleMappings: roleMappings,
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
