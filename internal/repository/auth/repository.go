package auth

import (
	"authorization_service/internal/client/db"
	"authorization_service/internal/model"
	"authorization_service/internal/repository"
	"context"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
)

const (
	usersTable         = "users"
	refreshTokensTable = "refresh_tokens"

	colID           = "id"
	colEmail        = "email"
	colPasswordHash = "password_hash"
	colRole         = "role"
	colUserID       = "user_id"
	colToken        = "token"
	colExpiresAt    = "expires_at"
)

var psq = sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

type repo struct {
	db db.Client
}

func NewRepository(db db.Client) repository.AuthRepository {
	return &repo{db: db}
}

func (r *repo) Register(ctx context.Context, email, passwordHash, role string) (int64, error) {
	builder := psq.Insert(usersTable).
		Columns(colEmail, colPasswordHash, colRole).
		Values(email, passwordHash, role).
		Suffix("RETURNING id")

	query, args, err := builder.ToSql()
	if err != nil {
		return 0, errors.Wrap(err, "build query")
	}

	q := db.Query{
		Name:     "auth_repository.Register",
		QueryRaw: query,
	}

	var id int64
	err = r.db.DB().ScanOneContext(ctx, &id, q, args...)
	if err != nil {
		return 0, errors.Wrap(err, "exec")
	}

	return id, nil
}

func (r *repo) Login(ctx context.Context, login *model.Login) (*model.UserDB, error) {
	builder := psq.Select(colID, colEmail, colPasswordHash).
		Column("COALESCE(role, 'user') AS role").
		From(usersTable).
		Where(sq.Eq{colEmail: login.Email}).
		Limit(1)

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "build query")
	}

	q := db.Query{
		Name:     "auth_repository.Login",
		QueryRaw: query,
	}

	var user model.UserDB
	err = r.db.DB().ScanOneContext(ctx, &user, q, args...)
	if err != nil {
		return nil, errors.Wrap(err, "user not found")
	}

	return &user, nil
}

func (r *repo) SaveRefreshToken(ctx context.Context, userID int64, token string, expiresAt time.Time) error {
	builder := psq.Insert(refreshTokensTable).
		Columns(colUserID, colToken, colExpiresAt).
		Values(userID, token, expiresAt)

	query, args, err := builder.ToSql()
	if err != nil {
		return errors.Wrap(err, "build query")
	}

	q := db.Query{
		Name:     "auth_repository.SaveRefreshToken",
		QueryRaw: query,
	}

	_, err = r.db.DB().ExecContext(ctx, q, args...)
	return errors.Wrap(err, "exec")
}

func (r *repo) CheckRefreshToken(ctx context.Context, token string) error {
	builder := psq.Select("1").
		From(refreshTokensTable).
		Where(sq.Eq{colToken: token}).
		Where(sq.Gt{colExpiresAt: time.Now()}).
		Limit(1)

	query, args, err := builder.ToSql()
	if err != nil {
		return errors.Wrap(err, "build query")
	}

	q := db.Query{
		Name:     "auth_repository.CheckRefreshToken",
		QueryRaw: query,
	}

	var dummy int
	err = r.db.DB().ScanOneContext(ctx, &dummy, q, args...)
	if err != nil {
		return errors.New("refresh token not found or expired")
	}

	return nil
}

func (r *repo) ReplaceRefreshToken(ctx context.Context, oldToken, newToken string, expiresAt time.Time) error {
	builder := psq.Update(refreshTokensTable).
		Set(colToken, newToken).
		Set(colExpiresAt, expiresAt).
		Where(sq.Eq{colToken: oldToken})

	query, args, err := builder.ToSql()
	if err != nil {
		return errors.Wrap(err, "build query")
	}

	q := db.Query{
		Name:     "auth_repository.ReplaceRefreshToken",
		QueryRaw: query,
	}

	_, err = r.db.DB().ExecContext(ctx, q, args...)
	return errors.Wrap(err, "exec")
}
