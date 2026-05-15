package service

import (
	"context"
	"crypto/rsa"
	"log/slog"
	"pob/pkg/jwt"
	"pob/user/internal/model"
	"pob/user/internal/model/apperror"
	"pob/user/internal/repository"
	"pob/user/internal/service/dto/auth"
)

type AuthService struct {
	ar        *repository.AuthRepository
	rr        *repository.RefreshTokenRepository
	publicKey *rsa.PublicKey
}

func NewAuthService(ar *repository.AuthRepository, rr *repository.RefreshTokenRepository, publicKey *rsa.PublicKey) *AuthService {
	return &AuthService{
		ar:        ar,
		rr:        rr,
		publicKey: publicKey,
	}
}

func (a *AuthService) Login(ctx context.Context, d auth.LoginRequest) (auth.TokenResponse, error) {
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

func (a *AuthService) Refresh(ctx context.Context, rt string) (auth.TokenResponse, error) {
	slog.InfoContext(ctx, "refresh start")

	claims, err := jwt.VerifyToken(rt, a.publicKey)
	if err != nil {
		slog.WarnContext(ctx, "failed to verify refresh token", slog.String("error", err.Error()))
		return auth.TokenResponse{}, err
	}
	userId := claims.UserID
	token, err := a.rr.FindByUserId(ctx, userId)
	if err != nil {
		slog.ErrorContext(ctx, "failed to find refresh token", slog.String("user_id", userId), slog.String("error", err.Error()))
		return auth.TokenResponse{}, err
	}

	ok := repository.Compare(token.TokenHash, rt)
	if !ok {
		slog.WarnContext(ctx, "refresh token mismatch", slog.String("user_id", userId))
		return auth.TokenResponse{}, apperror.ErrInvalidCredentials
	}

	newToken, err := a.ar.GenJwt(token.UserId)
	if err != nil {
		slog.ErrorContext(ctx, "failed to generate jwt", slog.String("user_id", userId), slog.String("error", err.Error()))
		return auth.TokenResponse{}, err
	}

	hashedRefresh, err := repository.Hash(newToken.RefreshToken)
	if err != nil {
		slog.ErrorContext(ctx, "failed to hash refresh token", slog.String("user_id", userId), slog.String("error", err.Error()))
		return auth.TokenResponse{}, err
	}
	newRefresh := model.NewRefreshToken(token.UserId, hashedRefresh)
	if err := a.rr.Save(ctx, *newRefresh); err != nil {
		slog.ErrorContext(ctx, "failed to save refresh token", slog.String("user_id", userId), slog.String("error", err.Error()))
		return auth.TokenResponse{}, err
	}

	slog.InfoContext(ctx, "refresh success", slog.String("user_id", userId))
	return auth.TokenResponse{
		AccessToken:  newToken.Token,
		RefreshToken: newToken.RefreshToken,
	}, nil
}

func (a *AuthService) Logout(ctx context.Context, at string) error {
	slog.InfoContext(ctx, "logout start")

	claims, err := jwt.VerifyToken(at, a.publicKey)
	if err != nil {
		slog.WarnContext(ctx, "failed to verify access token", slog.String("error", err.Error()))
		return apperror.ErrInvalidCredentials
	}

	userId := claims.UserID
	if err := a.rr.DeleteByUserId(ctx, userId); err != nil {
		return err
	}

	slog.InfoContext(ctx, "logout success", slog.String("user_id", userId))
	return nil
}
