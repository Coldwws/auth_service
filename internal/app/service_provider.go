package app

import (
	apiAuth "authorization_service/internal/api"
	"authorization_service/internal/client/db"
	"authorization_service/internal/client/db/pg"
	"authorization_service/internal/closer"
	"authorization_service/internal/config"
	"authorization_service/internal/repository"
	"authorization_service/internal/service"
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
)

type serviceProvider struct {
	config   *config.Config
	pgPool   *pgxpool.Pool
	dbClient db.Client

	authRepository repository.AuthRepository
	authService    service.AuthService

	authApi *apiAuth.Server
}

func NewServiceProvider(cfg *config.Config) *serviceProvider {
	return &serviceProvider{
		config: cfg,
	}
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

func (s *serviceProvider) AuthService(ctx context.Context) service.AuthService {
	if s.authService == nil {
		s.authService = authService.NewUserService(s.AuthRepository(ctx))
	}
	return s.authService
}

func (s *serviceProvider) UserAPI() *apiAuth.Server {

	if s.authApi == nil {
		s.authApi = apiAuth.NewServer(s.AuthService(ctx))
	}

	return s.authApi
}
