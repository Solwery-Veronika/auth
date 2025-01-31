package config

import (
	"log"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Postgres Postgres
	Platform Platform
	Service  Service
}

type Postgres struct {
	User     string `env:"POSTGRES_USER"`
	Password string `env:"POSTGRES_PASSWORD"`
	Host     string `env:"POSTGRES_HOST"`
	Port     string `env:"POSTGRES_PORT"`
	Database string `env:"POSTGRES_DBNAME"`
}

type Platform struct {
	Secret string `env:"PLATFORM_SECRET"`
}

type Service struct {
	Port string `env:"AUTH_SERVICE_PORT"`
}

func MustLoad() *Config {
	cfg := &Config{}
	err := cleanenv.ReadEnv(cfg)
	if err != nil {
		log.Fatalf("can not read env variables: %s", err)
	}
	return cfg
}
