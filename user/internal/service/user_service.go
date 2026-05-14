package service

import (
	"context"
	"log/slog"
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
	slog.InfoContext(ctx, "registration start", slog.String("username", r.UserName))

	hashedPassword, err := repository.Hash(r.Password)
	if err != nil {
		slog.ErrorContext(ctx, "failed to hash password", slog.String("username", r.UserName), slog.Any("error", err))
		return err
	}

	user := model.NewUser(r.UserName, hashedPassword)
	if err := u.repository.Register(ctx, user); err != nil {
		slog.WarnContext(ctx, "registration failed", slog.String("username", r.UserName), slog.Any("error", err))
		return err
	}

	slog.InfoContext(ctx, "registration success",
		slog.String("user_id", user.UserId.String()),
		slog.String("username", user.UserName),
	)
	return nil
}
