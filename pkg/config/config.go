package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port       string
	JWTSecret  string
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string
}

func Load() (*Config, error) {
	_ = godotenv.Load(".env.local")
	_ = godotenv.Load()

	cfg := &Config{
		Port:       os.Getenv("PORT"),
		JWTSecret:  os.Getenv("JWTSecret"),
		DBHost:     os.Getenv("DB_HOST"),
		DBPort:     os.Getenv("DB_PORT"),
		DBUser:     os.Getenv("DB_USER"),
		DBPassword: os.Getenv("DB_PASSWORD"),
		DBName:     os.Getenv("DB_NAME"),
		DBSSLMode:  os.Getenv("DB_SSLMODE"),
	}
	if cfg.DBSSLMode == "" {
		cfg.DBSSLMode = os.Getenv("DBSSLMode")
	}
	if cfg.DBSSLMode == "" {
		cfg.DBSSLMode = "require"
	}
	if cfg.Port == "" {
		cfg.Port = "8080"
	}
	if cfg.JWTSecret == "" {
		return nil, fmt.Errorf("missing required JWTSecret environment variable")
	}
	if cfg.DBHost == "" || cfg.DBPort == "" || cfg.DBUser == "" || cfg.DBName == "" {
		return nil, fmt.Errorf("missing required database environment variables")
	}

	return cfg, nil
}
