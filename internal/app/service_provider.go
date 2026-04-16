package app

import (
	"context"
	apiAuth "github.com/Coldwws/auth_service/internal/api"
	"github.com/Coldwws/auth_service/internal/client/db"
	"github.com/Coldwws/auth_service/internal/client/db/pg"
	"github.com/Coldwws/auth_service/internal/closer"
	"github.com/Coldwws/auth_service/internal/config"
	"github.com/Coldwws/auth_service/internal/repository"
	accessRepo "github.com/Coldwws/auth_service/internal/repository/access"
	authRepo "github.com/Coldwws/auth_service/internal/repository/auth"
	"github.com/Coldwws/auth_service/internal/service"
	accessServ "github.com/Coldwws/auth_service/internal/service/access"
	authServ "github.com/Coldwws/auth_service/internal/service/auth"
	"log"

	"github.com/jackc/pgx/v4/pgxpool"
)

type serviceProvider struct {
	config   *config.Config
	pgPool   *pgxpool.Pool
	dbClient db.Client

	authRepository   repository.AuthRepository
	accessRepository repository.AccessRepository

	authService   service.AuthService
	accessService service.AccessService

	authApi *apiAuth.Server
}

func NewServiceProvider(cfg *config.Config) *serviceProvider {
	return &serviceProvider{config: cfg}
}

var ctx = context.Background()

func (s *serviceProvider) PGConfig() config.PGConfig {
	if s.config == nil {
		log.Fatal("config is nil")
	}
	return s.config.PG
}

func (s *serviceProvider) PGPool() *pgxpool.Pool {
	if s.pgPool == nil {
		pool, err := pgxpool.Connect(ctx, s.PGConfig().DSN())
		if err != nil {
			log.Fatalf("Failed to connect database: %v", err)
		}
		err = pool.Ping(ctx)
		if err != nil {
			log.Fatalf("Failed to ping database: %v", err)
		}
		closer.Add(func() error {
			pool.Close()
			return nil
		})
		s.pgPool = pool
	}
	return s.pgPool
}

func (s *serviceProvider) DBClient(ctx context.Context) db.Client {
	if s.dbClient == nil {
		cl, err := pg.New(ctx, s.PGConfig().DSN())
		if err != nil {
			log.Fatalf("Failed to init pg client %v", err)
		}
		err = cl.DB().Ping(ctx)
		if err != nil {
			log.Fatalf("ping error: %v", err.Error())
		}
		closer.Add(cl.Close)
		s.dbClient = cl
	}
	return s.dbClient
}

func (s *serviceProvider) AuthRepository(ctx context.Context) repository.AuthRepository {
	if s.authRepository == nil {
		s.authRepository = authRepo.NewRepository(s.DBClient(ctx))
	}
	return s.authRepository
}

func (s *serviceProvider) AccessRepository(ctx context.Context) repository.AccessRepository {
	if s.accessRepository == nil {
		s.accessRepository = accessRepo.NewAccessRepository(s.DBClient(ctx))
	}
	return s.accessRepository
}

func (s *serviceProvider) AuthService(ctx context.Context) service.AuthService {
	if s.authService == nil {
		tokenCfg := s.config.Token
		s.authService = authServ.NewAuthService(
			s.AuthRepository(ctx),
			[]byte(tokenCfg.AccessSecretKey()),
			[]byte(tokenCfg.RefreshSecretKey()),
			tokenCfg.AccessTTL(),
			tokenCfg.RefreshTTL(),
		)
	}
	return s.authService
}

func (s *serviceProvider) AccessService(ctx context.Context) service.AccessService {
	if s.accessService == nil {
		s.accessService = accessServ.NewAccessService(
			s.AccessRepository(ctx),
			s.config.Token.AccessSecretKey(),
		)
	}
	return s.accessService
}

func (s *serviceProvider) AuthAPI() *apiAuth.Server {
	if s.authApi == nil {
		s.authApi = apiAuth.NewServer(
			s.AuthService(ctx),
			s.AccessService(ctx),
		)
	}
	return s.authApi
}
