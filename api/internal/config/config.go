package config

import (
	"os"
)

type Config struct {
	Port      string
	JWTSecret string
	AppEnv    string
}

func New() *Config {
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "your-secret-key"
	}

	appEnv := os.Getenv("APP_ENV")
	if appEnv == "" {
		appEnv = "development"
	}

	return &Config{
		Port:      port,
		JWTSecret: jwtSecret,
		AppEnv:    appEnv,
	}
}
