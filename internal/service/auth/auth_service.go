package service

import (
	"context"
	"github.com/Coldwws/auth_service/internal/model"
	"github.com/Coldwws/auth_service/internal/repository"
	"github.com/Coldwws/auth_service/internal/service"
	"github.com/Coldwws/auth_service/internal/utils"
	"golang.org/x/crypto/bcrypt"
	"time"

	"github.com/pkg/errors"
)

type authService struct {
	authRepository   repository.AuthRepository
	accessSecretKey  []byte
	refreshSecretKey []byte
	accessTTL        time.Duration
	refreshTTL       time.Duration
}

func NewAuthService(
	authRepository repository.AuthRepository,
	accessSecretKey []byte,
	refreshSecretKey []byte,
	accessTTL time.Duration,
	refreshTTL time.Duration,
) service.AuthService {
	return &authService{
		authRepository:   authRepository,
		accessSecretKey:  accessSecretKey,
		refreshSecretKey: refreshSecretKey,
		accessTTL:        accessTTL,
		refreshTTL:       refreshTTL,
	}
}
func (s *authService) Register(ctx context.Context, email, password, role string) (int64, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return 0, errors.Wrap(err, "failed to hash password")
	}

	id, err := s.authRepository.Register(ctx, email, string(hash), role)
	if err != nil {
		return 0, errors.Wrap(err, "failed to register user")
	}

	return id, nil
}

func (s *authService) Login(ctx context.Context, login *model.Login) (string, error) {
	// Получаем пользователя из репозитория (возвращает хэш пароля и роль)
	userInfo, err := s.authRepository.Login(ctx, login)
	if err != nil {
		return "", errors.Wrap(err, "failed to find user")
	}

	// Верифицируем пароль
	if !utils.VerifyPassword(userInfo.PasswordHash, login.Password) {
		return "", errors.New("invalid credentials")
	}

	// Генерируем refresh токен
	refreshToken, err := utils.GenerateToken(
		model.UserInfo{Email: userInfo.Email, Role: userInfo.Role},
		s.refreshSecretKey,
		s.refreshTTL,
	)
	if err != nil {
		return "", errors.Wrap(err, "failed to generate refresh token")
	}

	// Сохраняем refresh токен в БД
	err = s.authRepository.SaveRefreshToken(ctx, userInfo.ID, refreshToken, time.Now().Add(s.refreshTTL))
	if err != nil {
		return "", errors.Wrap(err, "failed to save refresh token")
	}

	return refreshToken, nil
}

func (s *authService) GetRefreshToken(ctx context.Context, oldRefreshToken string) (string, error) {
	// Верифицируем старый refresh токен
	claims, err := utils.VerifyToken(oldRefreshToken, s.refreshSecretKey)
	if err != nil {
		return "", errors.Wrap(err, "invalid refresh token")
	}

	// Проверяем, что токен есть в БД (не отозван)
	err = s.authRepository.CheckRefreshToken(ctx, oldRefreshToken)
	if err != nil {
		return "", errors.Wrap(err, "refresh token not found or revoked")
	}

	// Генерируем новый refresh токен
	newRefreshToken, err := utils.GenerateToken(
		model.UserInfo{Email: claims.Email, Role: claims.Role},
		s.refreshSecretKey,
		s.refreshTTL,
	)
	if err != nil {
		return "", errors.Wrap(err, "failed to generate new refresh token")
	}

	// Заменяем старый токен новым
	err = s.authRepository.ReplaceRefreshToken(ctx, oldRefreshToken, newRefreshToken, time.Now().Add(s.refreshTTL))
	if err != nil {
		return "", errors.Wrap(err, "failed to replace refresh token")
	}

	return newRefreshToken, nil
}

func (s *authService) GetAccessToken(ctx context.Context, refreshToken string) (string, error) {
	// Верифицируем refresh токен
	claims, err := utils.VerifyToken(refreshToken, s.refreshSecretKey)
	if err != nil {
		return "", errors.Wrap(err, "invalid refresh token")
	}

	// Проверяем наличие в БД
	err = s.authRepository.CheckRefreshToken(ctx, refreshToken)
	if err != nil {
		return "", errors.Wrap(err, "refresh token not found or revoked")
	}

	// Генерируем access токен
	accessToken, err := utils.GenerateToken(
		model.UserInfo{Email: claims.Email, Role: claims.Role},
		s.accessSecretKey,
		s.accessTTL,
	)
	if err != nil {
		return "", errors.Wrap(err, "failed to generate access token")
	}

	return accessToken, nil
}
