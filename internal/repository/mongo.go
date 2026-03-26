package repository

import (
	"context"
	"log/slog"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"portal/internal/config"
	"portal/internal/model"
)

const (
	// RealmsCollection stores realm projections.
	RealmsCollection = "kc_realms"
	// ClientsCollection stores client projections.
	ClientsCollection = "kc_clients"
	// ClientMetaCollection stores portal-only client metadata.
	ClientMetaCollection = "portal_client_meta"
	// UsersCollection stores user projections.
	UsersCollection = "kc_users"
	// SessionsCollection stores portal sessions.
	SessionsCollection = "portal_sessions"
	// SettingsCollection stores global portal settings.
	SettingsCollection = "portal_settings"
)

// Mongo wraps a MongoDB client and database.
type Mongo struct {
	Client   *mongo.Client
	Database *mongo.Database
	logger   *slog.Logger
}

// Repositories groups all Mongo repositories.
type Repositories struct {
	Realms      *RealmRepository
	Clients     *ClientRepository
	ClientMetas *ClientMetaRepository
	Users       *UserRepository
	Sessions    *SessionRepository
	Settings    *SettingsRepository
}

// NewMongo connects to MongoDB and ensures indexes.
func NewMongo(ctx context.Context, cfg config.Config, logger *slog.Logger) (*Mongo, error) {
	connectCtx, cancel := context.WithTimeout(ctx, cfg.Mongo.ConnectTimeout)
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

	if err := mongoDB.Ping(connectCtx); err != nil {
		_ = client.Disconnect(context.Background())
		return nil, err
	}

	if err := mongoDB.EnsureIndexes(ctx); err != nil {
		return nil, err
	}
	return mongoDB, nil
}

// EnsureIndexes creates the required Mongo indexes.
func (m *Mongo) EnsureIndexes(ctx context.Context) error {
	indexes := []struct {
		collection string
		models     []mongo.IndexModel
	}{
		{
			collection: RealmsCollection,
			models: []mongo.IndexModel{
				{Keys: bson.D{{Key: "realmId", Value: 1}}, Options: options.Index().SetUnique(true).SetName("ux_realm_id")},
			},
		},
		{
			collection: ClientsCollection,
			models: []mongo.IndexModel{
				{Keys: bson.D{{Key: "realmId", Value: 1}, {Key: "clientId", Value: 1}}, Options: options.Index().SetUnique(true).SetName("ux_realm_client_id")},
				{Keys: bson.D{{Key: "realmId", Value: 1}, {Key: "clientUuid", Value: 1}}, Options: options.Index().SetUnique(true).SetName("ux_realm_client_uuid")},
			},
		},
		{
			collection: ClientMetaCollection,
			models: []mongo.IndexModel{
				{Keys: bson.D{{Key: "realmId", Value: 1}, {Key: "clientId", Value: 1}}, Options: options.Index().SetUnique(true).SetName("ux_realm_client_meta")},
			},
		},
		{
			collection: UsersCollection,
			models: []mongo.IndexModel{
				{Keys: bson.D{{Key: "realmId", Value: 1}, {Key: "userId", Value: 1}}, Options: options.Index().SetUnique(true).SetName("ux_realm_user_id")},
				{Keys: bson.D{{Key: "realmId", Value: 1}, {Key: "username", Value: 1}}, Options: options.Index().SetUnique(true).SetName("ux_realm_username")},
			},
		},
		{
			collection: SessionsCollection,
			models: []mongo.IndexModel{
				{Keys: bson.D{{Key: "sessionId", Value: 1}}, Options: options.Index().SetUnique(true).SetName("ux_session_id")},
				{Keys: bson.D{{Key: "expiresAt", Value: 1}}, Options: options.Index().SetExpireAfterSeconds(0).SetName("ttl_expires_at")},
			},
		},
		{
			collection: SettingsCollection,
			models: []mongo.IndexModel{
				{Keys: bson.D{{Key: "_id", Value: 1}}, Options: options.Index().SetName("ix_settings_id")},
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

// Ping checks MongoDB readiness.
func (m *Mongo) Ping(ctx context.Context) error {
	return m.Client.Ping(ctx, nil)
}

// Close disconnects the MongoDB client.
func (m *Mongo) Close(ctx context.Context) error {
	return m.Client.Disconnect(ctx)
}

// NewRepositories creates all repository instances.
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

func defaultSettings(idleTimeoutMinutes int) model.PortalSettings {
	return model.PortalSettings{
		ID:                 "global",
		IdleTimeoutMinutes: idleTimeoutMinutes,
		IdleWarnSeconds:    60,
	}
}
