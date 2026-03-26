package repository

import (
	"context"
	"log/slog"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"portal/internal/model"
)

// SessionRepository persists portal sessions.
type SessionRepository struct {
	collection *mongo.Collection
	logger     *slog.Logger
}

// NewSessionRepository creates a SessionRepository.
func NewSessionRepository(db *mongo.Database, logger *slog.Logger) *SessionRepository {
	return &SessionRepository{
		collection: db.Collection(SessionsCollection),
		logger:     logger,
	}
}

// Create stores a new portal session.
func (r *SessionRepository) Create(ctx context.Context, session model.PortalSession) error {
	_, err := r.collection.InsertOne(ctx, session)
	return err
}

// GetByID loads a session by session ID.
func (r *SessionRepository) GetByID(ctx context.Context, sessionID string) (model.PortalSession, error) {
	var out model.PortalSession
	err := r.collection.FindOne(ctx, bson.M{"sessionId": sessionID}).Decode(&out)
	return out, err
}

// Touch refreshes the idle session timestamps.
func (r *SessionRepository) Touch(ctx context.Context, sessionID string, lastActiveAt, expiresAt time.Time) error {
	result, err := r.collection.UpdateOne(
		ctx,
		bson.M{"sessionId": sessionID},
		bson.M{"$set": bson.M{
			"lastActiveAt": lastActiveAt,
			"expiresAt":    expiresAt,
		}},
	)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return mongo.ErrNoDocuments
	}
	return nil
}

// Delete removes a session.
func (r *SessionRepository) Delete(ctx context.Context, sessionID string) error {
	result, err := r.collection.DeleteOne(ctx, bson.M{"sessionId": sessionID})
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return mongo.ErrNoDocuments
	}
	return nil
}
