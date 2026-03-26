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

type ClientMetaRepository struct {
	collection *mongo.Collection
	logger     *slog.Logger
}

func NewClientMetaRepository(db *mongo.Database, logger *slog.Logger) *ClientMetaRepository {
	return &ClientMetaRepository{
		collection: db.Collection(ClientMetaCollection),
		logger:     logger,
	}
}

func (r *ClientMetaRepository) ListByRealm(ctx context.Context, realm string) ([]model.PortalClientMeta, error) {
	cursor, err := r.collection.Find(ctx, bson.M{"realm": realm}, options.Find().SetSort(bson.D{{Key: "sortOrder", Value: 1}, {Key: "clientId", Value: 1}}))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var out []model.PortalClientMeta
	for cursor.Next(ctx) {
		var item model.PortalClientMeta
		if err := cursor.Decode(&item); err != nil {
			return nil, err
		}
		out = append(out, item)
	}
	return out, cursor.Err()
}

func (r *ClientMetaRepository) Upsert(ctx context.Context, meta model.PortalClientMeta) error {
	now := time.Now().UTC()
	if meta.CreatedAt.IsZero() {
		meta.CreatedAt = now
	}
	meta.UpdatedAt = now

	_, err := r.collection.UpdateOne(
		ctx,
		bson.M{"realm": meta.Realm, "clientId": meta.ClientID},
		bson.M{
			"$set": meta,
			"$setOnInsert": bson.M{
				"realm":     meta.Realm,
				"clientId":  meta.ClientID,
				"createdAt": meta.CreatedAt,
			},
		},
		options.Update().SetUpsert(true),
	)
	return err
}

func (r *ClientMetaRepository) SeedDefaults(ctx context.Context, metas []model.PortalClientMeta) error {
	if len(metas) == 0 {
		return nil
	}

	models := make([]mongo.WriteModel, 0, len(metas))
	now := time.Now().UTC()
	for _, meta := range metas {
		if meta.CreatedAt.IsZero() {
			meta.CreatedAt = now
		}
		meta.UpdatedAt = now
		models = append(models, mongo.NewUpdateOneModel().
			SetFilter(bson.M{"realm": meta.Realm, "clientId": meta.ClientID}).
			SetUpdate(bson.M{
				"$setOnInsert": meta,
			}).
			SetUpsert(true))
	}

	_, err := r.collection.BulkWrite(ctx, models, options.BulkWrite().SetOrdered(false))
	return err
}

func (r *ClientMetaRepository) Delete(ctx context.Context, realm, clientID string) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"realm": realm, "clientId": clientID})
	return err
}
