package service

import (
	"context"
	"log/slog"
	"pob/user/internal/model"
	"pob/user/internal/repository"
	"pob/user/internal/service/dto/auth"
	"pob/user/internal/shared"
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
	l := shared.FromContext(ctx)
	l.InfoContext(ctx, "login start", slog.String("username", d.UserName))

	au := model.NewAuth(d.UserName, d.Password)
	token, err := a.repository.Login(ctx, au)
	if err != nil {
		// service層ではエラーの詳細はrepository層で出力済みなのでwarnに留める
		l.WarnContext(ctx, "login failed", slog.String("username", d.UserName), slog.Any("error", err))
		return auth.TokenResponse{}, err
	}

	l.InfoContext(ctx, "login success", slog.String("username", d.UserName))
	return auth.TokenResponse{
		AccessToken:  token.Token,
		RefreshToken: token.RefreshToken,
	}, nil
}
