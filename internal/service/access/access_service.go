package access

import (
	"authorization_service/internal/repository"
	"authorization_service/internal/service"
	"authorization_service/internal/utils"
	"context"
	"errors"
	"google.golang.org/grpc/metadata"
	"strings"
)

type accessService struct {
	accessRepository repository.AccessRepository
	accessSecretKey  string
}

func NewAccessService(accessRepository repository.AccessRepository, accessSecretKey string) service.AccessService {
	return &accessService{
		accessRepository: accessRepository,
		accessSecretKey:  accessSecretKey,
	}
}

func (s *accessService) Check(ctx context.Context, endpointAddress string) error {
	// достаём токен из входящего контекста
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return errors.New("metadata not provided")
	}

	authHeader, ok := md["authorization"]
	if !ok || len(authHeader) == 0 {
		return errors.New("authorization token not provided")
	}

	token := strings.TrimPrefix(authHeader[0], "Bearer ")

	// верифицируем access токен
	_, err := utils.VerifyToken(token, []byte(s.accessSecretKey))
	if err != nil {
		return errors.New("invalid access token")
	}

	// проверяем доступ к эндпоинту
	return s.accessRepository.Check(ctx, endpointAddress)
}
