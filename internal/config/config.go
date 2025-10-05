package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env      string `env:"ENV" env-required:"true"`
	Server   `                 env-required:"true" env-prefix:"SERVER_"`
	Logger   `                 env-required:"true" env-prefix:"LOGGER_"`
	Postgres `                 env-required:"true" env-prefix:"POSTGRES_"`
	Redis    `                 env-required:"true" env-prefix:"REDIS_"`
	JWT      `                 env-required:"true" env-prefix:"JWT_"`
}

type Server struct {
	Host         string `env:"HOST"          env-required:"true"`
	Port         int    `env:"PORT"          env-required:"true"`
	ReadTimeout  int    `env:"READ_TIMEOUT"  env-required:"true"`
	WriteTimeout int    `env:"WRITE_TIMEOUT" env-required:"true"`
	IdleTimeout  int    `env:"IDLE_TIMEOUT"  env-required:"true"`
}

type Logger struct {
	Level  string `env:"LEVEL"  env-required:"true"`
	Output string `env:"OUTPUT" env-required:"true"`
	Format string `env:"FORMAT" env-required:"true"`
}

type Postgres struct {
	User     string `env:"USER"     env-required:"true"`
	Password string `env:"PASSWORD" env-required:"true"`
	Host     string `env:"HOST"     env-required:"true"`
	Port     int    `env:"PORT"     env-required:"true"`
	DB       string `env:"DB"       env-required:"true"`
	URL      string `env:"URL"      env-required:"true"`
}

type Redis struct {
	Host string `env:"HOST" env-required:"true"`
	Port int    `env:"PORT" env-required:"true"`
	DB   int    `env:"DB" env-required:"true"`
	URL  string `env:"URL"  env-required:"true"`
}

type JWT struct {
	Secret string `env:"SECRET" env-required:"true"`
}

func Load() (*Config, error) {
	c := new(Config)

	err := cleanenv.ReadEnv(c)
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	return c, nil
}

func MustLoad() *Config {
	c, err := Load()
	if err != nil {
		panic(err)
	}

	return c
}
