package service

import (
	"context"
	"pob/user/internal/model"
	"pob/user/internal/repository"
	"pob/user/internal/service/dto/user"
)

type UserService struct {
	repository *repository.UserRepository
}

func NewUserService(r *repository.UserRepository) *UserService {
	return &UserService{
		repository: r,
	}
}

func (u *UserService) Registration(ctx context.Context, r user.UserRegistration) error {
	hashedPassword, err := repository.HashPassword(r.Password)
	if err != nil {
		return err
	}
	user := model.NewUser(r.UserName, hashedPassword)
	dbErr := u.repository.Register(ctx, user)
	if dbErr != nil {
		return dbErr
	}
	return nil
}
