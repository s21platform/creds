package config

import (
	"log"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Service  Service
	Metrics  Metrics
	Platform Platform
	Logger   Logger
	Postgres Postgres
}

type Logger struct {
	Host string `env:"LOGGER_SERVICE_HOST"`
	Port string `env:"LOGGER_SERVICE_PORT"`
}

type Service struct {
	Port string `env:"CREDS_SERVICE_PORT"`
	Name string `env:"CREDS_SERVICE_NAME"`
}

type Postgres struct {
	User     string `env:"CREDS_SERVICE_POSTGRES_USER"`
	Password string `env:"CREDS_SERVICE_POSTGRES_PASSWORD"`
	Database string `env:"CREDS_SERVICE_POSTGRES_DB"`
	Host     string `env:"CREDS_SERVICE_POSTGRES_HOST"`
	Port     string `env:"CREDS_SERVICE_POSTGRES_PORT"`
}

type Metrics struct {
	Host string `env:"GRAFANA_HOST"`
	Port int    `env:"GRAFANA_PORT"`
}

type Platform struct {
	Env string `env:"ENV"`
}

func MustLoad() *Config {
	cfg := &Config{}
	err := cleanenv.ReadEnv(cfg)
	if err != nil {
		log.Fatalf("failed to read env variables: %s", err)
	}
	return cfg
}
