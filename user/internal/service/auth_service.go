package service

import (
	"context"
	"log/slog"
	"pob/user/internal/model"
	"pob/user/internal/repository"
	"pob/user/internal/service/dto/auth"
)

type AuthService struct {
	repository *repository.AuthRepository
}

func NewAuthService(r *repository.AuthRepository) *AuthService {
	return &AuthService{
		repository: r,
	}
}

func (a *AuthService) Login(ctx context.Context, d auth.Login) (auth.TokenResponse, error) {
	slog.InfoContext(ctx, "login start", slog.String("username", d.UserName))

	au := model.NewAuth(d.UserName, d.Password)
	token, err := a.repository.Login(ctx, au)
	if err != nil {
		slog.WarnContext(ctx, "login failed", slog.String("username", d.UserName), slog.Any("error", err))
		return auth.TokenResponse{}, err
	}

	slog.InfoContext(ctx, "login success", slog.String("username", d.UserName))
	return auth.TokenResponse{
		AccessToken:  token.Token,
		RefreshToken: token.RefreshToken,
	}, nil
}
