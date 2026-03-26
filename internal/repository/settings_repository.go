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

// SettingsRepository persists global portal settings.
type SettingsRepository struct {
	collection *mongo.Collection
	logger     *slog.Logger
}

// NewSettingsRepository creates a SettingsRepository.
func NewSettingsRepository(db *mongo.Database, logger *slog.Logger) *SettingsRepository {
	return &SettingsRepository{
		collection: db.Collection(SettingsCollection),
		logger:     logger,
	}
}

// GetGlobal loads the global portal settings document.
func (r *SettingsRepository) GetGlobal(ctx context.Context, defaultIdleTimeoutMinutes int) (model.PortalSettings, error) {
	var out model.PortalSettings
	err := r.collection.FindOne(ctx, bson.M{"_id": "global"}).Decode(&out)
	if err == mongo.ErrNoDocuments {
		return defaultSettings(defaultIdleTimeoutMinutes), nil
	}
	return out, err
}

// UpsertGlobal stores the global portal settings document.
func (r *SettingsRepository) UpsertGlobal(ctx context.Context, settings model.PortalSettings) error {
	settings.ID = "global"
	settings.UpdatedAt = time.Now().UTC()
	_, err := r.collection.UpdateOne(
		ctx,
		bson.M{"_id": "global"},
		bson.M{
			"$set": settings,
			"$setOnInsert": bson.M{"_id": "global"},
		},
		options.Update().SetUpsert(true),
	)
	return err
}
