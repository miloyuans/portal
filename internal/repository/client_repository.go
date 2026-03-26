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

type ClientRepository struct {
	collection *mongo.Collection
	logger     *slog.Logger
}

func NewClientRepository(db *mongo.Database, logger *slog.Logger) *ClientRepository {
	return &ClientRepository{
		collection: db.Collection(ClientsCollection),
		logger:     logger,
	}
}

func (r *ClientRepository) UpsertMany(ctx context.Context, clients []model.ClientProjection) error {
	if len(clients) == 0 {
		return nil
	}

	now := time.Now().UTC()
	models := make([]mongo.WriteModel, 0, len(clients))
	for _, client := range clients {
		client.UpdatedAt = now
		models = append(models, mongo.NewUpdateOneModel().
			SetFilter(bson.M{"realm": client.Realm, "clientId": client.ClientID}).
			SetUpdate(bson.M{
				"$set": client,
				"$setOnInsert": bson.M{
					"realm":    client.Realm,
					"clientId": client.ClientID,
				},
			}).
			SetUpsert(true))
	}

	_, err := r.collection.BulkWrite(ctx, models, options.BulkWrite().SetOrdered(false))
	return err
}

func (r *ClientRepository) ListByRealm(ctx context.Context, realm string) ([]model.ClientProjection, error) {
	cursor, err := r.collection.Find(ctx, bson.M{"realm": realm}, options.Find().SetSort(bson.D{{Key: "clientId", Value: 1}}))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var out []model.ClientProjection
	for cursor.Next(ctx) {
		var item model.ClientProjection
		if err := cursor.Decode(&item); err != nil {
			return nil, err
		}
		out = append(out, item)
	}
	return out, cursor.Err()
}

func (r *ClientRepository) GetByRealmAndClientID(ctx context.Context, realm, clientID string) (model.ClientProjection, error) {
	var out model.ClientProjection
	err := r.collection.FindOne(ctx, bson.M{"realm": realm, "clientId": clientID}).Decode(&out)
	return out, err
}
