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
	return  nil
}