package repository

import (
	"context"
	"crypto/rsa"
	"errors"
	"log/slog"
	"pob/user/internal/model"
	"pob/user/internal/model/apperror"
	"pob/user/internal/shared"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type AuthRepository struct {
	db         *shared.DBClient
	privateKey *rsa.PrivateKey
}

func NewAuthRepository(db *shared.DBClient, p *rsa.PrivateKey) *AuthRepository {
	return &AuthRepository{
		db:         db,
		privateKey: p,
	}
}

func (a *AuthRepository) Login(ctx context.Context, auth model.Auth) (model.Jwt, error) {
	user, err := a.findByUserName(ctx, auth.UserName)
	if err != nil {
		if errors.Is(err, apperror.ErrUserNotFound) {
			// ユーザーが存在しない場合は認証失敗として扱う（ユーザー存在有無を隠蔽）
			slog.WarnContext(ctx, "user not found on login", slog.String("username", auth.UserName))
			return model.Jwt{}, apperror.ErrInvalidCredentials
		}
		slog.ErrorContext(ctx, "failed to find user", slog.String("username", auth.UserName), slog.Any("error", err))
		return model.Jwt{}, err
	}

	if !Compare(user.PasswordHash, auth.PasswordPlane) {
		slog.WarnContext(ctx, "invalid password",
			slog.String("username", auth.UserName),
			slog.String("user_id", user.UserId.String()),
		)
		return model.Jwt{}, apperror.ErrInvalidCredentials
	}

	token, err := a.genJwt(user.UserId.String())
	if err != nil {
		slog.ErrorContext(ctx, "failed to generate jwt",
			slog.String("username", auth.UserName),
			slog.String("user_id", user.UserId.String()),
			slog.Any("error", err),
		)
		return model.Jwt{}, err
	}

	slog.InfoContext(ctx, "jwt generated", slog.String("user_id", user.UserId.String()))
	return token, nil
}

func (a *AuthRepository) findByUserName(ctx context.Context, u string) (model.User, error) {
	var id uuid.UUID
	var username, passwordHash string
	var createdAt, updatedAt time.Time
	query := `select id, username, password_hash, created_at, updated_at from users where username = $1`

	err := a.db.GetClient().QueryRow(ctx, query, u).Scan(&id, &username, &passwordHash, &createdAt, &updatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			slog.DebugContext(ctx, "user not found", slog.String("username", u))
			return model.User{}, apperror.ErrUserNotFound
		}
		slog.ErrorContext(ctx, "db error on findByUserName", slog.String("username", u), slog.Any("error", err))
		return model.User{}, err
	}

	return model.User{
		UserId:       id,
		UserName:     username,
		PasswordHash: passwordHash,
		CreatedAt:    createdAt,
		UpdatedAt:    updatedAt,
	}, nil
}

func (a *AuthRepository) genJwt(userId string) (model.Jwt, error) {
	accessClaims := model.Claims{
		UserID: userId,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodRS256, accessClaims)
	signedAccessToken, err := accessToken.SignedString(a.privateKey)
	if err != nil {
		return model.Jwt{}, err
	}

	refreshClaims := model.Claims{
		UserID: userId,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodRS256, refreshClaims)
	signedRefreshToken, err := refreshToken.SignedString(a.privateKey)
	if err != nil {
		return model.Jwt{}, err
	}

	return model.Jwt{
		Token:        signedAccessToken,
		RefreshToken: signedRefreshToken,
	}, nil
}
