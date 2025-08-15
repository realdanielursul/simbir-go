package config

import (
	"fmt"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type (
	Config struct {
		App      `yaml:"app"`
		HTTP     `yaml:"http"`
		Postgres `yaml:"postgres"`
		JWT      `yaml:"jwt"`
		Hasher   `yaml:"hasher"`
	}

	App struct {
		Name    string `yaml:"name"`
		Version string `yaml:"version"`
	}

	HTTP struct {
		Port string `yaml:"port"`
	}

	Postgres struct {
		Host     string `yaml:"host"`
		Port     string `yaml:"port"`
		Username string `yaml:"username"`
		Password string `yaml:"password" env:"DB_PASSWORD"`
		Database string `yaml:"database"`
		SSLMode  string `yaml:"ssl_mode"`
	}

	JWT struct {
		SignKey  string        `yaml:"sign_key" env:"SIGN_KEY"`
		TokenTTL time.Duration `yaml:"token_ttl"`
	}

	Hasher struct {
		Salt string `yaml:"salt" env:"SALT"`
	}
)

func NewConfig(configPath string) (*Config, error) {
	cfg := &Config{}

	if err := cleanenv.ReadConfig(configPath, cfg); err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	if err := cleanenv.UpdateEnv(cfg); err != nil {
		return nil, fmt.Errorf("error updating env file: %w", err)
	}

	return cfg, nil
}
