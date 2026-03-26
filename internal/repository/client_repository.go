package repository

import (
	"context"
	"log/slog"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"portal/internal/model"
)

// ClientRepository persists client projections.
type ClientRepository struct {
	collection *mongo.Collection
	logger     *slog.Logger
}

// NewClientRepository creates a ClientRepository.
func NewClientRepository(db *mongo.Database, logger *slog.Logger) *ClientRepository {
	return &ClientRepository{
		collection: db.Collection(ClientsCollection),
		logger:     logger,
	}
}

// UpsertMany stores client projections.
func (r *ClientRepository) UpsertMany(ctx context.Context, clients []model.ClientProjection) error {
	if len(clients) == 0 {
		return nil
	}

	models := make([]mongo.WriteModel, 0, len(clients))
	for _, client := range clients {
		models = append(models, mongo.NewUpdateOneModel().
			SetFilter(bson.M{"realmId": client.RealmID, "clientId": client.ClientID}).
			SetUpdate(bson.M{
				"$set": client,
				"$setOnInsert": bson.M{
					"realmId":  client.RealmID,
					"clientId": client.ClientID,
				},
			}).
			SetUpsert(true))
	}

	_, err := r.collection.BulkWrite(ctx, models, options.BulkWrite().SetOrdered(false))
	return err
}

// ListByRealm returns all clients for a realm.
func (r *ClientRepository) ListByRealm(ctx context.Context, realmID string) ([]model.ClientProjection, error) {
	cursor, err := r.collection.Find(ctx, bson.M{"realmId": realmID}, options.Find().SetSort(bson.D{{Key: "clientId", Value: 1}}))
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

// GetByRealmAndClientID returns a client projection by client ID.
func (r *ClientRepository) GetByRealmAndClientID(ctx context.Context, realmID, clientID string) (model.ClientProjection, error) {
	var out model.ClientProjection
	err := r.collection.FindOne(ctx, bson.M{"realmId": realmID, "clientId": clientID}).Decode(&out)
	return out, err
}

// GetByRealmAndClientUUID returns a client projection by client UUID.
func (r *ClientRepository) GetByRealmAndClientUUID(ctx context.Context, realmID, clientUUID string) (model.ClientProjection, error) {
	var out model.ClientProjection
	err := r.collection.FindOne(ctx, bson.M{"realmId": realmID, "clientUuid": clientUUID}).Decode(&out)
	return out, err
}
