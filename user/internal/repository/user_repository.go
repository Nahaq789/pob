package repository

import (
	"context"
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
	l := shared.FromContext(ctx)
	query := `insert`
	c := u.db.GetClient()

	_, err := c.Exec(ctx, query)
	if err != nil {
		l.ErrorContext(ctx, "failed to register user", "error", err)
		return err
	}
	return nil
}

