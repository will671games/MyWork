package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type DBConfig struct {
	Host     string `env:"DB_HOST" env-default:"sample"`
	Port     string `env:"DB_PORT" env-default:"5432"`
	User     string `env:"DB_USER" env-default:"sample"`
	Password string `env:"DB_PASSWORD" env-default:"sample"`
	Database string `env:"DB_NAME" env-default:"sample"`
	SSLMode  string `env:"DB_SSL_MODE" env-default:"disable"`
	Migrate  bool   `env:"DB_MIGRATE" env-default:"true"`
}

type HttpConfig struct {
	Port string `env:"HTTP_PORT" env-default:"8080"`
}

func LoadEnv() error {
	err := godotenv.Load("config.env")
	if err != nil {
		log.Printf("Warning: config.env file not found, using system environment variables: %v", err)
	}
	return nil
}

func NewDBConfig() DBConfig {
	LoadEnv()

	return DBConfig{
		Host:     getEnv("DB_HOST", "localhost"),
		Port:     getEnv("DB_PORT", "5432"),
		User:     getEnv("DB_USER", "postgres"),
		Password: getEnv("DB_PASSWORD", "postgres"),
		Database: getEnv("DB_NAME", "walletdb"),
		SSLMode:  getEnv("DB_SSL_MODE", "disable"),
		Migrate:  getBoolEnv("DB_MIGRATE", true),
	}
}

func NewHttpConfig() HttpConfig {
	LoadEnv()

	return HttpConfig{
		Port: getEnv("HTTP_PORT", "8080"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getBoolEnv(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if parsed, err := strconv.ParseBool(value); err == nil {
			return parsed
		}
	}
	return defaultValue
}
