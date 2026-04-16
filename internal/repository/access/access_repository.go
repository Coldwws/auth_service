package access

import (
	"context"
	"github.com/Coldwws/auth_service/internal/client/db"
	"github.com/Coldwws/auth_service/internal/repository"
	sq "github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
)

type repo struct {
	db db.Client
}

func NewAccessRepository(db db.Client) repository.AccessRepository {
	return &repo{
		db: db,
	}
}

func (s *repo) Check(ctx context.Context, endpointAddress string) error {
	builder := sq.Select("1").
		PlaceholderFormat(sq.Dollar).
		From("access_rules").
		Where(sq.Eq{"endpoint_address": endpointAddress}).
		Limit(1)

	query, args, err := builder.ToSql()
	if err != nil {
		return errors.Wrap(err, "build query")
	}

	q := db.Query{
		Name:     "access_repository.Check",
		QueryRaw: query,
	}

	var dummy int
	err = s.db.DB().ScanOneContext(ctx, &dummy, q, args...)
	if err != nil {
		return errors.Errorf("access denied for endpoint: %s", endpointAddress)
	}

	return nil
}
