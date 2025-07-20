package configs

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DB   DbConfig
	Auth AuthConfig
}

type DbConfig struct {
	DSN string
}

type AuthConfig struct {
	SecretKey             string
	RefreshTokenSecretKey string
}

func LoadConfig() (*Config, error) {
	godotenv.Load()

	return &Config{
		DB: DbConfig{
			DSN: os.Getenv("DB_DSN"),
		},
		Auth: AuthConfig{
			SecretKey:             os.Getenv("SECRET_KEY"),
			RefreshTokenSecretKey: os.Getenv("REFRESH_SECRET_KEY"),
		},
	}, nil
}
