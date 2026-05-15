package model

import (
	"time"

	"github.com/google/uuid"
)

type RefreshToken struct {
	RefreshTokenId uuid.UUID
	UserId         uuid.UUID
	TokenHash      string
	ExpiredAt      time.Time
	CreatedAt      time.Time
}

func NewRefreshToken(userId uuid.UUID, tokenHash string) *RefreshToken {
	now := time.Now()
	return &RefreshToken{
		RefreshTokenId: uuid.New(),
		UserId:         userId,
		TokenHash:      tokenHash,
		ExpiredAt:      now.Add(7 * 24 * time.Hour),
		CreatedAt:      now,
	}
}

func FromRefreshToken(id, userId uuid.UUID, tokenHash string, expiredAt, createdAt time.Time) *RefreshToken {
	return &RefreshToken{
		RefreshTokenId: id,
		UserId:         userId,
		TokenHash:      tokenHash,
		ExpiredAt:      expiredAt,
		CreatedAt:      createdAt,
	}
}
