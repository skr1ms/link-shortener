package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DB   DbConfig
	Auth AuthConfig
}

type DbConfig struct {
	URL string
}

type AuthConfig struct {
	SecretKey             string
	RefreshTokenSecretKey string
}

func LoadConfig() (*Config, error) {
	godotenv.Load()

	return &Config{
		DB: DbConfig{
			URL: os.Getenv("DB_URL"),
		},
		Auth: AuthConfig{
			SecretKey:             os.Getenv("SECRET_KEY"),
			RefreshTokenSecretKey: os.Getenv("REFRESH_SECRET_KEY"),
		},
	}, nil
}
