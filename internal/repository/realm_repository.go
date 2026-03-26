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

type RealmRepository struct {
	collection *mongo.Collection
	logger     *slog.Logger
}

func NewRealmRepository(db *mongo.Database, logger *slog.Logger) *RealmRepository {
	return &RealmRepository{
		collection: db.Collection(RealmsCollection),
		logger:     logger,
	}
}

func (r *RealmRepository) Upsert(ctx context.Context, realm model.RealmProjection) error {
	realm.UpdatedAt = time.Now().UTC()
	_, err := r.collection.UpdateOne(
		ctx,
		bson.M{"realm": realm.Realm},
		bson.M{
			"$set": realm,
			"$setOnInsert": bson.M{
				"realm": realm.Realm,
			},
		},
		options.Update().SetUpsert(true),
	)
	return err
}

func (r *RealmRepository) GetByRealm(ctx context.Context, realm string) (model.RealmProjection, error) {
	var out model.RealmProjection
	err := r.collection.FindOne(ctx, bson.M{"realm": realm}).Decode(&out)
	return out, err
}
