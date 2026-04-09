package access

import (
	"authorization_service/internal/repository"
	"authorization_service/internal/service"
	"context"
)

type accessService struct {
	accessRepository repository.AccessRepository
}

func NewAccessService(accessRepository repository.AccessRepository) service.AccessService {
	return &accessService{
		accessRepository: accessRepository,
	}
}

func (s *accessService) Check(ctx context.Context, endpointAdress string) error {
	return nil
}
