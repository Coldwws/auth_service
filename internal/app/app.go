package app

import (
	"authorization_service/internal/closer"
	"authorization_service/internal/config"
	"authorization_service/pkg/access_v1"
	"authorization_service/pkg/auth_v1"
	"context"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
)

type App struct {
	config          *config.Config
	serviceProvider *serviceProvider
	grpcServer      *grpc.Server
}

func NewApp(ctx context.Context) (*App, error) {
	a := &App{}
	err := a.InitDeps(ctx)
	if err != nil {
		return nil, err
	}
	return a, nil
}

func (a *App) InitDeps(ctx context.Context) error {
	inits := []func(context.Context) error{
		a.InitConfig,
		a.initServiceProvider,
		a.initGRPCServer,
	}
	for _, f := range inits {
		if err := f(ctx); err != nil {
			return err
		}
	}
	return nil
}

func (a *App) InitConfig(_ context.Context) error {
	if err := godotenv.Load("docker.env"); err != nil {
		log.Println("Warning: local.env not found, using system env")
	}

	cfg := config.LoadConfig()
	a.config = &cfg

	return nil
}

func (a *App) initServiceProvider(_ context.Context) error {
	a.serviceProvider = NewServiceProvider(a.config)
	return nil
}

func (a *App) initGRPCServer(ctx context.Context) error {
	a.grpcServer = grpc.NewServer()
	grpc.Creds(insecure.NewCredentials())

	reflection.Register(a.grpcServer)
	auth_v1.RegisterAuthV1Server(a.grpcServer, a.serviceProvider.AuthAPI())
	access_v1.RegisterAccessV1Server(a.grpcServer, a.serviceProvider.AuthAPI())

	return nil

}

func (a *App) runGRPCServer() error {

	log.Println("GRPC server is running on:", a.serviceProvider.config.GRPC.Address())

	list, err := net.Listen("tcp", a.serviceProvider.config.GRPC.Address())
	if err != nil {
		return err
	}

	err = a.grpcServer.Serve(list)
	if err != nil {
		return err
	}

	return nil

}

func (a *App) Run() error {
	defer func() {
		closer.CloseAll()
		closer.Wait()
	}()

	err := a.runGRPCServer()
	if err != nil {
		log.Printf("failed to run grpc server: %v", err)
	}

	return nil
}
