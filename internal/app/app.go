package app

import (
	"authorization_service/internal/closer"
	"authorization_service/internal/config"
	"authorization_service/internal/logger"
	"authorization_service/internal/tracing"
	"authorization_service/pkg/access_v1"
	"authorization_service/pkg/auth_v1"
	"context"
	"github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"
	"github.com/joho/godotenv"
	"github.com/natefinch/lumberjack"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"os"
)

var (
	serviceName = "auth_service"
)

type App struct {
	config          *config.Config
	serviceProvider *serviceProvider
	grpcServer      *grpc.Server
	logLevel        string
}

func NewApp(ctx context.Context, logLevel string) (*App, error) {
	a := &App{logLevel: logLevel}
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
		a.initLogger,
		a.initTracing,
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

func (a *App) initLogger(_ context.Context) error {
	logger.Init(a.getCore(a.getAtomicLevel()))

	return nil
}

func (a *App) initTracing(_ context.Context) error {
	tracing.Init(logger.Logger(), serviceName)
	return nil
}

func (a *App) getCore(level zap.AtomicLevel) zapcore.Core {
	stdout := zapcore.AddSync(os.Stdout)

	file := zapcore.AddSync(&lumberjack.Logger{
		Filename:   "logs/app.log",
		MaxSize:    10, // megabytes
		MaxBackups: 3,
		MaxAge:     7, // days
	})

	productionCfg := zap.NewProductionEncoderConfig()
	productionCfg.TimeKey = "timestamp"
	productionCfg.EncodeTime = zapcore.ISO8601TimeEncoder

	developmentCfg := zap.NewDevelopmentEncoderConfig()
	developmentCfg.EncodeLevel = zapcore.CapitalColorLevelEncoder

	consoleEncoder := zapcore.NewConsoleEncoder(developmentCfg)
	fileEncoder := zapcore.NewJSONEncoder(productionCfg)

	return zapcore.NewTee(
		zapcore.NewCore(consoleEncoder, stdout, level),
		zapcore.NewCore(fileEncoder, file, level),
	)
}

func (a *App) getAtomicLevel() zap.AtomicLevel {
	var level zapcore.Level
	if err := level.Set(a.logLevel); err != nil {
		log.Fatalf("failed to set log level: %w", err)
	}

	return zap.NewAtomicLevelAt(level)
}

func (a *App) initGRPCServer(ctx context.Context) error {
	a.grpcServer = grpc.NewServer(
		grpc.Creds(insecure.NewCredentials()),
		grpc.UnaryInterceptor(otgrpc.OpenTracingServerInterceptor(opentracing.GlobalTracer())),
	)
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
