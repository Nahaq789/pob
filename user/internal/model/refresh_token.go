package model

import (
	"time"

	"github.com/google/uuid"
)

type RefreshToken struct {
	RefreshTokenId uuid.UUID
	UserId uuid.UUID
	TokenHash string
	ExpiredAt time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewRefreshToken(userId uuid.UUID, tokenHash string, expiredAt time.Time) *RefreshToken {
	return &RefreshToken{
		RefreshTokenId: uuid.New(),
		UserId: userId,
		TokenHash: tokenHash,
		ExpiredAt: expiredAt,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}