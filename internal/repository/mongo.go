package repository

import (
	"context"
	"log/slog"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"portal/internal/config"
	"portal/internal/model"
)

const (
	RealmsCollection     = "kc_realms"
	ClientsCollection    = "kc_clients"
	ClientMetaCollection = "portal_client_meta"
	UsersCollection      = "kc_users"
	SessionsCollection   = "portal_sessions"
	SettingsCollection   = "portal_settings"
)

type Mongo struct {
	Client   *mongo.Client
	Database *mongo.Database
	logger   *slog.Logger
}

type Repositories struct {
	Realms      *RealmRepository
	Clients     *ClientRepository
	ClientMetas *ClientMetaRepository
	Users       *UserRepository
	Sessions    *SessionRepository
	Settings    *SettingsRepository
}

func NewMongo(ctx context.Context, cfg config.Config, logger *slog.Logger) (*Mongo, error) {
	connectCtx, cancel := context.WithTimeout(ctx, time.Duration(cfg.Mongo.ConnectTimeout)*time.Second)
	defer cancel()

	client, err := mongo.Connect(connectCtx, options.Client().ApplyURI(cfg.Mongo.URI))
	if err != nil {
		return nil, err
	}

	mongoDB := &Mongo{
		Client:   client,
		Database: client.Database(cfg.Mongo.Database),
		logger:   logger,
	}

	if err := mongoDB.EnsureIndexes(ctx); err != nil {
		return nil, err
	}
	return mongoDB, nil
}

func (m *Mongo) EnsureIndexes(ctx context.Context) error {
	indexes := []struct {
		collection string
		models     []mongo.IndexModel
	}{
		{
			collection: RealmsCollection,
			models: []mongo.IndexModel{
				{Keys: bson.D{{Key: "realm", Value: 1}}, Options: options.Index().SetUnique(true).SetName("ux_realm")},
			},
		},
		{
			collection: ClientsCollection,
			models: []mongo.IndexModel{
				{Keys: bson.D{{Key: "realm", Value: 1}, {Key: "clientId", Value: 1}}, Options: options.Index().SetUnique(true).SetName("ux_realm_client_id")},
				{Keys: bson.D{{Key: "realm", Value: 1}, {Key: "clientUuid", Value: 1}}, Options: options.Index().SetUnique(true).SetName("ux_realm_client_uuid")},
			},
		},
		{
			collection: ClientMetaCollection,
			models: []mongo.IndexModel{
				{Keys: bson.D{{Key: "realm", Value: 1}, {Key: "clientId", Value: 1}}, Options: options.Index().SetUnique(true).SetName("ux_realm_client_meta")},
			},
		},
		{
			collection: UsersCollection,
			models: []mongo.IndexModel{
				{Keys: bson.D{{Key: "realm", Value: 1}, {Key: "userId", Value: 1}}, Options: options.Index().SetUnique(true).SetName("ux_realm_user")},
			},
		},
		{
			collection: SessionsCollection,
			models: []mongo.IndexModel{
				{Keys: bson.D{{Key: "sessionId", Value: 1}}, Options: options.Index().SetUnique(true).SetName("ux_session_id")},
				{Keys: bson.D{{Key: "expiresAt", Value: 1}}, Options: options.Index().SetExpireAfterSeconds(0).SetName("ttl_expires_at")},
				{Keys: bson.D{{Key: "realm", Value: 1}, {Key: "userId", Value: 1}}, Options: options.Index().SetName("ix_realm_user")},
			},
		},
		{
			collection: SettingsCollection,
			models: []mongo.IndexModel{
				{Keys: bson.D{{Key: "realm", Value: 1}}, Options: options.Index().SetUnique(true).SetName("ux_settings_realm")},
			},
		},
	}

	for _, entry := range indexes {
		if _, err := m.Database.Collection(entry.collection).Indexes().CreateMany(ctx, entry.models); err != nil {
			return err
		}
	}
	m.logger.Info("mongo indexes ensured")
	return nil
}

func (m *Mongo) Ping(ctx context.Context) error {
	return m.Client.Ping(ctx, nil)
}

func (m *Mongo) Close(ctx context.Context) error {
	return m.Client.Disconnect(ctx)
}

func NewRepositories(db *mongo.Database, logger *slog.Logger) *Repositories {
	return &Repositories{
		Realms:      NewRealmRepository(db, logger),
		Clients:     NewClientRepository(db, logger),
		ClientMetas: NewClientMetaRepository(db, logger),
		Users:       NewUserRepository(db, logger),
		Sessions:    NewSessionRepository(db, logger),
		Settings:    NewSettingsRepository(db, logger),
	}
}

func defaultSettings(realm string, idleTimeoutMinutes int) model.PortalSettings {
	now := time.Now().UTC()
	return model.PortalSettings{
		Realm:              realm,
		IdleTimeoutMinutes: idleTimeoutMinutes,
		CreatedAt:          now,
		UpdatedAt:          now,
	}
}
