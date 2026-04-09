package repository

import (
	"authorization_service/internal/model"
	"context"
	"time"
)

type AuthRepository interface {
	Register(ctx context.Context, email, passwordHash, role string) (int64, error)
	Login(ctx context.Context, login *model.Login) (*model.UserDB, error)
	SaveRefreshToken(ctx context.Context, userID int64, token string, expiresAt time.Time) error
	CheckRefreshToken(ctx context.Context, token string) error
	ReplaceRefreshToken(ctx context.Context, oldToken, newToken string, expiresAt time.Time) error
}

type AccessRepository interface {
	Check(ctx context.Context, endpointAdress string) error
}
