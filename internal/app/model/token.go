package model

import (
	"encoding/base64"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// Token ...
type Token struct {
	UserID      string `json:"user_id"`
	Fingerprint string `json:"fingerprint"`
	jwt.StandardClaims
}

// CreatePairTokens ...
func (t *Token) CreatePairTokens(userID string, fingerprint string) (map[string]string, error) {
	at := &Token{
		UserID:      userID,
		Fingerprint: fingerprint,
	}

	rt := &Token{
		UserID:      userID,
		Fingerprint: fingerprint,
	}

	accessToken, err := at.CreateAccess()
	if err != nil {
		return nil, err
	}

	refreshToken, err := rt.CreateRefresh()
	if err != nil {
		return nil, err
	}

	tokens := map[string]string{
		"accessToken":  accessToken,
		"refreshToken": refreshToken,
	}

	return tokens, nil
}

// CreateAccess ...
func (t *Token) CreateAccess() (string, error) {

	claims := Token{
		t.UserID,
		t.Fingerprint,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Minute * 15).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)

	tokenString, err := token.SignedString([]byte(os.Getenv("token_password")))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// CreateRefresh ...
func (t *Token) CreateRefresh() (string, error) {

	claims := Token{
		t.UserID,
		t.Fingerprint,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)

	tokenString, err := token.SignedString([]byte(os.Getenv("token_password")))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ParseToken ...
func (t *Token) ParseToken(tk string) error {
	token, err := jwt.ParseWithClaims(tk, t, func(*jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("token_password")), nil
	})
	if err != nil {
		return err
	}

	if !token.Valid {
		return ErrInvalidToken
	}

	return nil
}

// DecodeFromBase64 ...
func (t *Token) DecodeFromBase64(encodedToken string) (string, error) {
	byteDecodedRefreshToken, err := base64.StdEncoding.DecodeString(encodedToken)
	if err != nil {
		return "", err
	}

	return string(byteDecodedRefreshToken), nil
}
