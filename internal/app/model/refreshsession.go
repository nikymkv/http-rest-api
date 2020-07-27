package model

import (
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// RefreshSession ...
type RefreshSession struct {
	ID                    primitive.ObjectID `json:"session_id" bson:"_id"`
	UserID                string             `json:"user_id" bson:"user_id"`
	RefreshToken          string             `json:"refresh_token" bson:"refresh_token"`
	EncryptedRefreshToken string             `json:"encrypted_refresh_token" bson:"encrypted_refresh_token"`
	Fingerprint           string             `json:"fingerpint" bson:"fingerprint"`
	ExpiresIn             string             `json:"expires_in" bson:"expires_in"`
	CreatedAt             string             `json:"created_at" bson:"created_at"`
}

// CreateRefreshSession ...
func (rs *RefreshSession) CreateRefreshSession(userID string, rt string, fingerprint string) error {
	expiresIn := time.Now().Add(time.Hour * 24).Unix()
	createdAt := time.Now().Unix()

	b, err := encryptToken(rt)
	if err != nil {
		return err
	}

	rs.UserID = userID
	rs.RefreshToken = rt
	rs.EncryptedRefreshToken = b
	rs.Fingerprint = fingerprint
	rs.ExpiresIn = fmt.Sprintf("%d", expiresIn)
	rs.CreatedAt = fmt.Sprintf("%d", createdAt)

	return nil
}

func encryptToken(t string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(t), bcrypt.MinCost)
	if err != nil {
		return "", err
	}

	return string(b), nil
}

// CompareRefreshToken ...
func (rs *RefreshSession) CompareRefreshToken(t string) bool {
	return bcrypt.CompareHashAndPassword([]byte(rs.RefreshToken), []byte(t)) == nil
}
