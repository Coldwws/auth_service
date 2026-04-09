package service

import (
	"authorization_service/internal/model"
	"authorization_service/internal/repository"
	"authorization_service/internal/service"
	"context"
)

type authService struct {
	authRepository repository.AuthRepository
}

func NewAuthService(authRepository repository.AuthRepository) service.AuthService {
	return &authService{
		authRepository: authRepository,
	}
}

func (s *authService) Login(ctx context.Context, login *model.Login) (string, error) {
	return "", nil
}

func (s *authService) GetRefreshToken(ctx context.Context, refreshToken string) (string, error) {
	return "", nil
}

func (s *authService) GetAccessToken(ctx context.Context, accessToken string) (string, error) {
	return "", nil
}
