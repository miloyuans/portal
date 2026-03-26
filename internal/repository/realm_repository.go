package repository

import (
	"context"
	"log/slog"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"portal/internal/model"
)

// RealmRepository persists realm projections.
type RealmRepository struct {
	collection *mongo.Collection
	logger     *slog.Logger
}

// NewRealmRepository creates a RealmRepository.
func NewRealmRepository(db *mongo.Database, logger *slog.Logger) *RealmRepository {
	return &RealmRepository{
		collection: db.Collection(RealmsCollection),
		logger:     logger,
	}
}

// Upsert stores a realm projection.
func (r *RealmRepository) Upsert(ctx context.Context, realm model.RealmProjection) error {
	_, err := r.collection.UpdateOne(
		ctx,
		bson.M{"realmId": realm.RealmID},
		bson.M{
			"$set": realm,
			"$setOnInsert": bson.M{"realmId": realm.RealmID},
		},
		options.Update().SetUpsert(true),
	)
	return err
}

// GetByRealmID loads a realm projection by realm ID.
func (r *RealmRepository) GetByRealmID(ctx context.Context, realmID string) (model.RealmProjection, error) {
	var out model.RealmProjection
	err := r.collection.FindOne(ctx, bson.M{"realmId": realmID}).Decode(&out)
	return out, err
}

// List returns all realm projections.
func (r *RealmRepository) List(ctx context.Context) ([]model.RealmProjection, error) {
	cursor, err := r.collection.Find(ctx, bson.M{}, options.Find().SetSort(bson.D{{Key: "realmName", Value: 1}}))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var out []model.RealmProjection
	for cursor.Next(ctx) {
		var item model.RealmProjection
		if err := cursor.Decode(&item); err != nil {
			return nil, err
		}
		out = append(out, item)
	}
	return out, cursor.Err()
}
