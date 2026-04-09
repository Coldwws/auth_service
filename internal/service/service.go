package service

import (
	"authorization_service/internal/model"
	"context"
)

type AuthService interface {
	Login(ctx context.Context, login *model.Login) (string, error)
	GetRefreshToken(ctx context.Context, refreshToken string) (string, error)
	GetAccessToken(ctx context.Context, accessToken string) (string, error)
}

type AccessService interface {
	Check(ctx context.Context, endpointAdress string) error
}
