package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type TokenProvider interface {
	GenerateToken(id string, name string, email string) (string, error)
	VerifyToken(token string) (*User, error)
}

type simpleTokenProvider struct {
	secret string
}

func NewTokenProvider(secret string) TokenProvider {
	return &simpleTokenProvider{secret}
}

func (r *simpleTokenProvider) GenerateToken(id string, name string, email string) (string, error) {
	signingToken := []byte(r.secret)
	claims := jwt.MapClaims{
		"id":       id,
		"name":     name,
		"email":    email,
		"exp":      time.Now().Add(7 * 24 * time.Hour).Unix(),
		"issuedAt": time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)
	return token.SignedString(signingToken)
}

func (r *simpleTokenProvider) VerifyToken(tokenStr string) (*User, error) {
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		return []byte(r.secret), nil
	})

	if err != nil || !token.Valid {
		return nil, errors.New("invalid token")
	}

	claims := token.Claims.(jwt.MapClaims)

	idStr := claims["id"].(string)
	name := claims["name"].(string)
	email := claims["email"].(string)

	uid, err := uuid.Parse(idStr)
	if err != nil {
		return nil, errors.New("invalid uuid in token")
	}

	return &User{
		Id:    uid,
		Name:  name,
		Email: email,
	}, nil
}
