package commons

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	AppName string
	AppEnv  string
	AppPort string

	PostgresHost     string
	PostgresPort     string
	PostgresUser     string
	PostgresPassword string
	PostgresDB       string
	PostgresSSLMode  string

	JWTSecret           string
	JWTAccessTTLMinutes string
}

var AppConfig Config

func LoadConfig() {
	_ = godotenv.Load()

	AppConfig = Config{
		AppName: getEnv("APP_NAME", "reserveflow-v1"),
		AppEnv:  getEnv("APP_ENV", "development"),
		AppPort: getEnv("APP_PORT", "8083"),

		PostgresHost:     getEnv("POSTGRES_HOST", "localhost"),
		PostgresPort:     getEnv("POSTGRES_PORT", "5433"),
		PostgresUser:     getEnv("POSTGRES_USER", "postgres"),
		PostgresPassword: getEnv("POSTGRES_PASSWORD", "postgres"),
		PostgresDB:       getEnv("POSTGRES_DB", "reserveflow_v1"),
		PostgresSSLMode:  getEnv("POSTGRES_SSLMODE", "disable"),

		JWTSecret:           getEnv("JWT_SECRET", "change-this-secret"),
		JWTAccessTTLMinutes: getEnv("JWT_ACCESS_TTL_MINUTES", "15"),
	}
}

func (c Config) PostgresDSN() string {
	return fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		c.PostgresHost,
		c.PostgresUser,
		c.PostgresPassword,
		c.PostgresDB,
		c.PostgresPort,
		c.PostgresSSLMode,
	)
}

func getEnv(key string, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	return value
}
