package config

import (
	"context"
	"fmt"

	"github.com/joho/godotenv"
	"github.com/sethvargo/go-envconfig"
)

type (
	// Config of service
	Config struct {
		InstanceConfig
		DBConfig
		RedisConfig
	}

	InstanceConfig struct {
		ServiceName string `env:"SERVICE_NAME,required"`
		HTTPPort    string `env:"HTTP_PORT, default=8081"`
	}

	DBConfig struct {
		PGURL          string `env:"PG_URL,required"`
		MigrationsPath string `env:"MIGRATIONS_PATH, default=migrations/"`
		UseMigrates    bool   `env:"USE_MIGRATES,default=false"`
	}

	RedisConfig struct {
		RedisAddr string `env:"REDIS_ADDR,required"`
	}
)

// Get
func Get() (Config, error) {
	return parseConfig()
}

func parseConfig() (cfg Config, err error) {
	godotenv.Load()

	err = envconfig.Process(context.Background(), &cfg)
	if err != nil {
		return cfg, fmt.Errorf("fill config: %w", err)
	}

	return cfg, nil
}
