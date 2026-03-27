package repository

import (
	"context"
	"log/slog"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"portal/internal/model"
)

// ClientMetaRepository persists portal client metadata.
type ClientMetaRepository struct {
	collection *mongo.Collection
	logger     *slog.Logger
}

// NewClientMetaRepository creates a ClientMetaRepository.
func NewClientMetaRepository(db *mongo.Database, logger *slog.Logger) *ClientMetaRepository {
	return &ClientMetaRepository{
		collection: db.Collection(ClientMetaCollection),
		logger:     logger,
	}
}

// ListByRealm returns all portal client metadata for a realm.
func (r *ClientMetaRepository) ListByRealm(ctx context.Context, realmID string) ([]model.PortalClientMeta, error) {
	cursor, err := r.collection.Find(ctx, bson.M{"realmId": realmID}, options.Find().SetSort(bson.D{{Key: "sort", Value: 1}, {Key: "clientId", Value: 1}}))
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

// GetByRealmAndClientID returns portal metadata for a single client.
func (r *ClientMetaRepository) GetByRealmAndClientID(ctx context.Context, realmID, clientID string) (model.PortalClientMeta, error) {
	var out model.PortalClientMeta
	err := r.collection.FindOne(ctx, bson.M{"realmId": realmID, "clientId": clientID}).Decode(&out)
	return out, err
}

// Upsert stores portal client metadata.
func (r *ClientMetaRepository) Upsert(ctx context.Context, meta model.PortalClientMeta) error {
	_, err := r.collection.UpdateOne(
		ctx,
		bson.M{"realmId": meta.RealmID, "clientId": meta.ClientID},
		bson.M{
			"$set": meta,
			"$setOnInsert": bson.M{
				"realmId":  meta.RealmID,
				"clientId": meta.ClientID,
			},
		},
		options.Update().SetUpsert(true),
	)
	return err
}

// SeedDefaults inserts metadata only when it does not already exist.
func (r *ClientMetaRepository) SeedDefaults(ctx context.Context, metas []model.PortalClientMeta) error {
	if len(metas) == 0 {
		return nil
	}

	models := make([]mongo.WriteModel, 0, len(metas))
	for _, meta := range metas {
		models = append(models, mongo.NewUpdateOneModel().
			SetFilter(bson.M{"realmId": meta.RealmID, "clientId": meta.ClientID}).
			SetUpdate(bson.M{"$setOnInsert": meta}).
			SetUpsert(true))
	}

	_, err := r.collection.BulkWrite(ctx, models, options.BulkWrite().SetOrdered(false))
	return err
}

// Delete removes portal client metadata.
func (r *ClientMetaRepository) Delete(ctx context.Context, realmID, clientID string) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"realmId": realmID, "clientId": clientID})
	return err
}
