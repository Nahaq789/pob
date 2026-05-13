package repository

import (
	"context"
	"pob/user/internal/model"
	"pob/user/internal/shared"
)

type RefreshTokenRepository struct {
	db *shared.DBClient
}

func NewRefreshTokenRepository(db *shared.DBClient) *RefreshTokenRepository {
	return &RefreshTokenRepository{
		db: db,
	}
}

func (r *RefreshTokenRepository) Save(ctx context.Context, token model.RefreshToken) error {
	l := shared.FromContext(ctx)
	query := `insert into refresh_tokens (id, user_id, token_hash, expires_at, created_at, updated_at) values ($1, $2, $3, $4, $5, %6)`
	c := r.db.GetClient()

	_, err := c.Exec(ctx, query, token.RefreshTokenId, token.UserId, token.TokenHash, token.ExpiredAt, token.CreatedAt, token.UpdatedAt)
	if err != nil {
		l.ErrorContext(ctx, "failed to save refresh token", "error", err)
		return err
	}
	return nil
}

