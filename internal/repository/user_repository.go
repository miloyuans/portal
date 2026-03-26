package repository

import (
	"context"
	"log/slog"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"portal/internal/model"
)

// UserRepository persists user projections.
type UserRepository struct {
	collection *mongo.Collection
	logger     *slog.Logger
}

// NewUserRepository creates a UserRepository.
func NewUserRepository(db *mongo.Database, logger *slog.Logger) *UserRepository {
	return &UserRepository{
		collection: db.Collection(UsersCollection),
		logger:     logger,
	}
}

// Upsert stores a user projection.
func (r *UserRepository) Upsert(ctx context.Context, user model.UserProjection) error {
	_, err := r.collection.UpdateOne(
		ctx,
		bson.M{"realmId": user.RealmID, "userId": user.UserID},
		bson.M{
			"$set": user,
			"$setOnInsert": bson.M{
				"realmId": user.RealmID,
				"userId":  user.UserID,
			},
		},
		options.Update().SetUpsert(true),
	)
	return err
}

// GetByRealmAndUserID loads a user projection by user ID.
func (r *UserRepository) GetByRealmAndUserID(ctx context.Context, realmID, userID string) (model.UserProjection, error) {
	var out model.UserProjection
	err := r.collection.FindOne(ctx, bson.M{"realmId": realmID, "userId": userID}).Decode(&out)
	return out, err
}

// GetByRealmAndUsername loads a user projection by username.
func (r *UserRepository) GetByRealmAndUsername(ctx context.Context, realmID, username string) (model.UserProjection, error) {
	var out model.UserProjection
	err := r.collection.FindOne(ctx, bson.M{"realmId": realmID, "username": username}).Decode(&out)
	return out, err
}
