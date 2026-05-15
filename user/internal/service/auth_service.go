package service

import (
	"context"
	"log/slog"
	"pob/user/internal/model"
	"pob/user/internal/repository"
	"pob/user/internal/service/dto/auth"
)

type AuthService struct {
	ar *repository.AuthRepository
	rr *repository.RefreshTokenRepository
}

func NewAuthService(ar *repository.AuthRepository, rr *repository.RefreshTokenRepository) *AuthService {
	return &AuthService{
		ar: ar,
		rr: rr,
	}
}

func (a *AuthService) Login(ctx context.Context, d auth.Login) (auth.TokenResponse, error) {
	slog.InfoContext(ctx, "login start", slog.String("username", d.UserName))

	au := model.NewAuth(d.UserName, d.Password)
	token, err := a.ar.Login(ctx, au)
	if err != nil {
		slog.WarnContext(ctx, "login failed", slog.String("username", d.UserName), slog.Any("error", err))
		return auth.TokenResponse{}, err
	}

	hashed, err := repository.Hash(token.RefreshToken)
	if err != nil {
		slog.ErrorContext(ctx, "failed to hash", "error", err)
		return auth.TokenResponse{}, err
	}
	refreshToken := model.NewRefreshToken(token.UserId, hashed)
	if err := a.rr.Save(ctx, *refreshToken); err != nil {
		slog.ErrorContext(ctx, "failed to save refresh token",
			slog.String("username", d.UserName),
			slog.String("user_id", token.UserId.String()),
			slog.Any("error", err),
		)
		return auth.TokenResponse{}, err
	}

	slog.InfoContext(ctx, "login success", slog.String("username", d.UserName))
	return auth.TokenResponse{
		AccessToken:  token.Token,
		RefreshToken: token.RefreshToken,
	}, nil
}

func (a *AuthService) Refresh(ctx context.Context) error {
	return nil
}
