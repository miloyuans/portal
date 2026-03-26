package repository

import (
	"context"
	"log/slog"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"portal/internal/model"
)

type SettingsRepository struct {
	collection *mongo.Collection
	logger     *slog.Logger
}

func NewSettingsRepository(db *mongo.Database, logger *slog.Logger) *SettingsRepository {
	return &SettingsRepository{
		collection: db.Collection(SettingsCollection),
		logger:     logger,
	}
}

func (r *SettingsRepository) GetByRealm(ctx context.Context, realm string, defaultIdleTimeoutMinutes int) (model.PortalSettings, error) {
	var out model.PortalSettings
	err := r.collection.FindOne(ctx, bson.M{"realm": realm}).Decode(&out)
	if err == mongo.ErrNoDocuments {
		return defaultSettings(realm, defaultIdleTimeoutMinutes), nil
	}
	return out, err
}

func (r *SettingsRepository) Upsert(ctx context.Context, settings model.PortalSettings) error {
	now := time.Now().UTC()
	if settings.CreatedAt.IsZero() {
		settings.CreatedAt = now
	}
	settings.UpdatedAt = now

	_, err := r.collection.UpdateOne(
		ctx,
		bson.M{"realm": settings.Realm},
		bson.M{
			"$set": settings,
			"$setOnInsert": bson.M{
				"realm":     settings.Realm,
				"createdAt": settings.CreatedAt,
			},
		},
		options.Update().SetUpsert(true),
	)
	return err
}
