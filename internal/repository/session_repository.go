package repository

import (
	"context"
	"log/slog"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"portal/internal/model"
)

type SessionRepository struct {
	collection *mongo.Collection
	logger     *slog.Logger
}

func NewSessionRepository(db *mongo.Database, logger *slog.Logger) *SessionRepository {
	return &SessionRepository{
		collection: db.Collection(SessionsCollection),
		logger:     logger,
	}
}

func (r *SessionRepository) Create(ctx context.Context, session model.PortalSession) error {
	_, err := r.collection.InsertOne(ctx, session)
	return err
}

func (r *SessionRepository) GetByID(ctx context.Context, sessionID string) (model.PortalSession, error) {
	var out model.PortalSession
	err := r.collection.FindOne(ctx, bson.M{"sessionId": sessionID}).Decode(&out)
	return out, err
}

func (r *SessionRepository) Touch(ctx context.Context, sessionID string, lastSeenAt, expiresAt time.Time) error {
	_, err := r.collection.UpdateOne(
		ctx,
		bson.M{"sessionId": sessionID},
		bson.M{
			"$set": bson.M{
				"lastSeenAt": lastSeenAt,
				"expiresAt":  expiresAt,
				"updatedAt":  time.Now().UTC(),
			},
		},
	)
	return err
}

func (r *SessionRepository) Delete(ctx context.Context, sessionID string) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"sessionId": sessionID})
	return err
}
