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

type UserRepository struct {
	collection *mongo.Collection
	logger     *slog.Logger
}

func NewUserRepository(db *mongo.Database, logger *slog.Logger) *UserRepository {
	return &UserRepository{
		collection: db.Collection(UsersCollection),
		logger:     logger,
	}
}

func (r *UserRepository) Upsert(ctx context.Context, user model.UserProjection) error {
	user.UpdatedAt = time.Now().UTC()
	_, err := r.collection.UpdateOne(
		ctx,
		bson.M{"realm": user.Realm, "userId": user.UserID},
		bson.M{
			"$set": user,
			"$setOnInsert": bson.M{
				"realm":  user.Realm,
				"userId": user.UserID,
			},
		},
		options.Update().SetUpsert(true),
	)
	return err
}

func (r *UserRepository) GetByRealmAndUserID(ctx context.Context, realm, userID string) (model.UserProjection, error) {
	var out model.UserProjection
	err := r.collection.FindOne(ctx, bson.M{"realm": realm, "userId": userID}).Decode(&out)
	return out, err
}
