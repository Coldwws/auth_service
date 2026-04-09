package env

import (
	"errors"
	"os"
	"time"
)

type TokenConfig struct {
	accessSecretKey  string
	refreshSecretKey string
	accessTTL        time.Duration
	refreshTTL       time.Duration
}

func NewTokenConfig() (*TokenConfig, error) {
	accessKey := os.Getenv("ACCESS_TOKEN_SECRET_KEY")
	if len(accessKey) == 0 {
		return nil, errors.New("ACCESS_TOKEN_SECRET_KEY is empty")
	}

	refreshKey := os.Getenv("REFRESH_TOKEN_SECRET_KEY")
	if len(refreshKey) == 0 {
		return nil, errors.New("REFRESH_TOKEN_SECRET_KEY is empty")
	}

	accessTTL := os.Getenv("ACCESS_TTL")
	if len(accessTTL) == 0 {
		return nil, errors.New("ACCESS_TTL is empty")
	}
	accessDuration, err := time.ParseDuration(accessTTL)
	if err != nil {
		return nil, err
	}

	refreshTTL := os.Getenv("REFRESH_TTL")
	if len(refreshTTL) == 0 {
		return nil, errors.New("REFRESH_TTL is empty")
	}
	refreshDuration, err := time.ParseDuration(refreshTTL)
	if err != nil {
		return nil, err
	}

	return &TokenConfig{
		accessSecretKey:  accessKey,
		refreshSecretKey: refreshKey,
		accessTTL:        accessDuration,
		refreshTTL:       refreshDuration,
	}, nil
}

func (c *TokenConfig) AccessSecretKey() string   { return c.accessSecretKey }
func (c *TokenConfig) RefreshSecretKey() string  { return c.refreshSecretKey }
func (c *TokenConfig) AccessTTL() time.Duration  { return c.accessTTL }
func (c *TokenConfig) RefreshTTL() time.Duration { return c.refreshTTL }
