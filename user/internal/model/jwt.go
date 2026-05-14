package model

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type Jwt struct {
	UserId       uuid.UUID
	Token        string
	RefreshToken string
}

type Claims struct {
	UserID string
	jwt.RegisteredClaims
}

func NewClaims(userID string, duration time.Duration) *Claims {
	return &Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
}
