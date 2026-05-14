package repository

import (
	"context"
	"crypto/rsa"
	"errors"
	"fmt"
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
	l := shared.FromContext(ctx)
	user, err := a.findByUserName(ctx, auth.UserName)
	if err != nil {
		l.ErrorContext(ctx, "failed to login user", "error", err)
		return model.Jwt{}, err
	}

	// パスワード検証
	hashed := user.PasswordHash
	plane := auth.PasswordPlane

	if !Compare(hashed, plane) {
		l.ErrorContext(ctx, "invalid credentials", "username", auth.UserName)
		return model.Jwt{}, apperror.ErrInvalidCredentials
	}

	token, err := a.genJwt(user.UserId.String())
	if err != nil {
		l.ErrorContext(ctx, "failed to generate jwt token", "error", err)
		return model.Jwt{}, err
	}
	return token, err
}

func (a *AuthRepository) findByUserName(ctx context.Context, u string) (model.User, error) {
	l := shared.FromContext(ctx)

	var id uuid.UUID
	var username, passwordHash string
	var createdAt, updatedAt time.Time
	query := `select id, username, password_hash, created_at, updated_at from users where username = $1`

	c := a.db.GetClient()

	err := c.QueryRow(ctx, query, u).Scan(&id, &username, &passwordHash, &createdAt, &updatedAt)
	if err != nil {
		l.ErrorContext(ctx, "failed to select user", "error", err)
		if errors.Is(err, pgx.ErrNoRows) {
			return model.User{}, fmt.Errorf("user not found: %w", err)
		}
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
