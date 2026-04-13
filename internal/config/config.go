package config

import (
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	ServerPort string

	DBHost     string
	DBPort     string
	DBName     string
	DBUser     string
	DBPassword string
	DBSSLMode  string

	JWTSecret string
	JWTTTL    time.Duration
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	cfg := &Config{
		ServerPort: getEnv("SERVER_PORT", "8080"),

		DBHost:     getEnv("DB_HOST", "db"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBName:     getEnv("DB_NAME", "wishlistdb"),
		DBUser:     getEnv("DB_USER", "wishlistuser"),
		DBPassword: getEnv("DB_PASSWORD", "wishlistpwd"),
		DBSSLMode:  getEnv("DB_SSLMODE", "disable"),

		JWTSecret: getEnv("JWT_SECRET", "change_me"),
		JWTTTL:    24 * time.Hour,
	}

	return cfg, nil
}

func (c *Config) DBConnString() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		c.DBUser,
		c.DBPassword,
		c.DBHost,
		c.DBPort,
		c.DBName,
		c.DBSSLMode,
	)
}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	return value
}