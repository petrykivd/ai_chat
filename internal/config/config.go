package config

import (
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

var (
	Postgres  PostgresConfig
	Anthropic AnthropicConfig
	Server    ServerConfig
)

func Load() error {
	_ = godotenv.Load()

	if err := envconfig.Process("", &Postgres); err != nil {
		return err
	}
	if err := envconfig.Process("", &Anthropic); err != nil {
		return err
	}
	if err := envconfig.Process("", &Server); err != nil {
		return err
	}
	return nil
}
