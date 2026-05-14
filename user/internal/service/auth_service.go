package service

import (
	"context"
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
	au := model.NewAuth(d.UserName, d.Password)
	token, err := a.repository.Login(ctx, au)
	if err != nil {
		return auth.TokenResponse{}, err
	}

	return auth.TokenResponse{
		JwtToken: token.Token,
	}, nil
}
