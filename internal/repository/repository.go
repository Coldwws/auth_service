package repository

import (
	"authorization_service/internal/model"
	"context"
)

type AuthRepository interface {
	Login(ctx context.Context, login *model.Login) (string, error)
	GetRefreshToken(ctx context.Context, refreshToken string) (string, error)
	GetAccessToken(ctx context.Context, accessToken string) (string, error)
}

type AccessRepository interface {
	Check(ctx context.Context, endpointAdress string) error
}
