package config

import (
	"github.com/Coldwws/auth_service/internal/config/env"
	"log"
	"os"
	"time"
)

type Config struct {
	Env   string
	GRPC  GRPCConfig
	PG    PGConfig
	Token TokenConfig
}
type TokenConfig interface {
	AccessSecretKey() string
	RefreshSecretKey() string
	AccessTTL() time.Duration
	RefreshTTL() time.Duration
}

func LoadConfig() Config {
	return Config{
		Env:   getEnv("ENV", "local"),
		GRPC:  loadGRPC(),
		PG:    loadPG(),
		Token: mustLoadTokenConfig(),
	}
}

func mustLoadTokenConfig() TokenConfig {
	cfg, err := env.NewTokenConfig()
	if err != nil {
		log.Fatalf("failed to load token config: %v", err)
	}
	return cfg
}

func mustEnv(key string) string {
	v := os.Getenv(key)
	if v == "" {
		log.Fatalf("missing env var :%s", key)
	}
	return v
}

func getEnv(key, def string) string {
	v := os.Getenv(key)
	if v == "" {
		return def
	}
	return v
}
