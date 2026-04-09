package access

import (
	"authorization_service/internal/client/db"
	"authorization_service/internal/repository"
	"context"

	sq "github.com/Masterminds/squirrel"
)

var (
	psq = sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
)

type repo struct {
	db db.Client
}

func NewAccessRepository(db db.Client) repository.AccessRepository {
	return &repo{
		db: db,
	}
}

func (s *repo) Check(ctx context.Context, endpointAdress string) error {
	return nil
}
