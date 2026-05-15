package repository

import (
	"context"
	"errors"
	"log/slog"
	"pob/user/internal/model"
	"pob/user/internal/shared"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
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

func (r *RefreshTokenRepository) FindByUserId(ctx context.Context, u string) (*model.RefreshToken, error) {
	var id, userId uuid.UUID
	var tokenHash string
	var expiresAt, createdAt time.Time
	query := `select id, user_id, token_hash, expires_at, created_at from refresh_tokens where user_id = $1`

	err := r.db.GetClient().QueryRow(ctx, query, u).Scan(&id, &userId, &tokenHash, &expiresAt, &createdAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			slog.DebugContext(ctx, "refresh token not found", slog.String("user_id", u))
			return nil, err
		}
		slog.ErrorContext(ctx, "db error on FindByUserId", slog.String("user_id", u), slog.Any("error", err))
		return nil, err
	}

	result := model.FromRefreshToken(id, userId, tokenHash, expiresAt, createdAt)
	return result, nil
}
