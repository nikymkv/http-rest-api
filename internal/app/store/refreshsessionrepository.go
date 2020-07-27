package store

import (
	"context"

	"github.com/http-rest-api/internal/app/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// RefreshSessionRepository ...
type RefreshSessionRepository struct {
	store                    *Store
	refreshSessionCollection *mongo.Collection
}

// CreateNewSession ...
func (rsr *RefreshSessionRepository) CreateNewSession(ctx context.Context, userID string, rt string, fingerprint string) error {
	refreshSession := &model.RefreshSession{}
	refreshSession.CreateRefreshSession(userID, rt, fingerprint)

	_, err := rsr.refreshSessionCollection.InsertOne(ctx, bson.D{
		{Key: "user_id", Value: refreshSession.UserID},
		{Key: "refresh_token", Value: refreshSession.EncryptedRefreshToken},
		{Key: "fingerprint", Value: refreshSession.Fingerprint},
		{Key: "expires_in", Value: refreshSession.ExpiresIn},
		{Key: "created_at", Value: refreshSession.CreatedAt},
	})
	if err != nil {
		return err
	}

	return nil
}

// CheckRefreshSession ...
func (rsr *RefreshSessionRepository) CheckRefreshSession(ctx context.Context, userID string, rt string, fingerprint string) (string, error) {
	refreshSession := model.RefreshSession{}
	filter := bson.M{
		"user_id":     userID,
		"fingerprint": fingerprint,
	}

	err := rsr.refreshSessionCollection.FindOne(ctx, filter).Decode(&refreshSession)
	if err != nil {
		return "", err
	}

	if refreshSession.CompareRefreshToken(rt) {
		return refreshSession.ID.Hex(), nil
	}

	return "", ErrInvalidToken
}

// DeleteRefreshSession ...
func (rsr *RefreshSessionRepository) DeleteRefreshSession(ctx context.Context, sessionID string) error {
	sessionObjectID, err := primitive.ObjectIDFromHex(sessionID)
	if err != nil {
		return err
	}
	filter := bson.M{"_id": sessionObjectID}

	session, err := rsr.store.client.StartSession()
	if err != nil{
		return err
	}

	err = session.StartTransaction()
	if err != nil {
		return err
	}

	if err = mongo.WithSession(ctx, session, func(sc mongo.SessionContext) error {
		result, err := rsr.refreshSessionCollection.DeleteOne(ctx, filter)
		if err != nil {
			return err
		}

		if result.DeletedCount == 0 {
			return ErrDocumentNotFound
		}

		err = session.CommitTransaction(sc)
		if err != nil {
			return err
		}

		return nil
	}); err != nil {
		return err
	}
	session.EndSession(ctx)

	_, err = rsr.refreshSessionCollection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}

	return nil
}

// DeleteAllRefreshSessions ...
func (rsr *RefreshSessionRepository) DeleteAllRefreshSessions(ctx context.Context, userID string) error {
	filter := bson.M{"user_id": userID}

	session, err := rsr.store.client.StartSession()
	if err != nil{
		return err
	}

	err = session.StartTransaction()
	if err != nil {
		return err
	}

	if err = mongo.WithSession(ctx, session, func(sc mongo.SessionContext) error {
		res, err := rsr.refreshSessionCollection.DeleteMany(ctx, filter)
		if err != nil {
			return err
		}

		if res.DeletedCount == 0 {
			return ErrDocumentNotFound
		}

		err = session.CommitTransaction(sc)
		if err != nil {
			return err
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}
