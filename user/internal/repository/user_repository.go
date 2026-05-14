package repository

import (
	"context"
	"log/slog"
	"pob/user/internal/model"
	"pob/user/internal/shared"
)

type UserRepository struct {
	db *shared.DBClient
}

func NewUserRepository(db *shared.DBClient) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (u *UserRepository) Register(ctx context.Context, user model.User) error {
	query := `insert into users (id, username, password_hash, created_at, updated_at) values ($1, $2, $3, $4, $5)`

	_, err := u.db.GetClient().Exec(ctx, query, user.UserId, user.UserName, user.PasswordHash, user.CreatedAt, user.UpdatedAt)
	if err != nil {
		slog.ErrorContext(ctx, "failed to register user",
			slog.String("user_id", user.UserId.String()),
			slog.String("username", user.UserName),
			slog.Any("error", err),
		)
		return err
	}
	return nil
}
