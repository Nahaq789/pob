package repository

import (
	"context"
	"log/slog"
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
	query := `
		INSERT INTO refresh_tokens (id, user_id, token_hash, expires_at, created_at)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (user_id)
		DO UPDATE SET
  		token_hash = EXCLUDED.token_hash,
  		expires_at = EXCLUDED.expires_at
  	`
	c := r.db.GetClient()

	_, err := c.Exec(ctx, query, token.RefreshTokenId, token.UserId, token.TokenHash, token.ExpiredAt, token.CreatedAt)
	if err != nil {
		slog.ErrorContext(ctx, "failed to save refresh token", "error", err)
		return err
	}
	return nil
}
