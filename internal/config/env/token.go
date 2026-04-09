package env

import (
	"authorization_service/internal/config"
	"errors"
	"os"
	"time"
)

const accessSecretKeyEnvName = "ACCESS_TOKEN_SECRET_KEY"
const refreshSecretKeyEnvName = "REFRESH_TOKEN_SECRET_KEY"
const accessTTLEnvName = "ACCESS_TTL"
const refreshTTLEnvName = "REFRESH_TTL"

type tokenConfig struct {
	accessSecretKey  string
	refreshSecretKey string
	accessTTL        time.Duration
	refreshTTL       time.Duration
}

func NewTokenConfig() (config.TokenConfig, error) {
	accessKey := os.Getenv(accessSecretKeyEnvName)
	if len(accessKey) == 0 {
		return nil, errors.New(accessKey)
	}

	refreshKey := os.Getenv(refreshSecretKeyEnvName)
	if len(refreshKey) == 0 {
		return nil, errors.New(refreshKey)
	}

	accessTTL := os.Getenv(accessTTLEnvName)
	if len(accessTTL) == 0 {
		return nil, errors.New(accessTTL)
	}

	accessDurationTTL, err := time.ParseDuration(accessTTL)
	if err != nil {
		return nil, err
	}

	refreshTTL := os.Getenv(refreshTTLEnvName)
	if len(refreshTTL) == 0 {
		return nil, errors.New(refreshTTL)
	}

	refreshDurationTTL, err := time.ParseDuration(refreshTTL)
	if err != nil {
		return nil, err
	}

	return &tokenConfig{
		accessSecretKey:  accessKey,
		refreshSecretKey: refreshKey,
		accessTTL:        accessDurationTTL,
		refreshTTL:       refreshDurationTTL,
	}, nil
}

func (c *tokenConfig) AccessSecretKey() string {
	return c.accessSecretKey
}

func (c *tokenConfig) RefreshSecretKey() string {
	return c.refreshSecretKey
}

func (c *tokenConfig) AccessTTL() time.Duration {
	return c.accessTTL
}

func (c *tokenConfig) RefreshTTL() time.Duration {
	return c.refreshTTL
}
