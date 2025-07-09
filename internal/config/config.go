package config

import (
	"errors"
	"os"
	"strconv"
)

type Config struct {
	Port           string
	Env            string
	DBType         string
	DBPath         string
	DBHost         string
	DBPort         string
	DBName         string
	DBUser         string
	DBPassword     string
	DBSSLMode      string
	SessionSecret  string
	SessionMaxAge  int
	AllowedOrigins string
}

func Load() *Config {
	maxAge, _ := strconv.Atoi(getEnv("SESSION_MAX_AGE", "86400"))

	return &Config{
		Port:           getEnv("PORT", "8080"),
		Env:            getEnv("ENV", "development"),
		DBType:         getEnv("DB_TYPE", "sqlite"),
		DBPath:         getEnv("DB_PATH", "./bambino.db"),
		DBHost:         getEnv("DB_HOST", "localhost"),
		DBPort:         getEnv("DB_PORT", "5432"),
		DBName:         getEnv("DB_NAME", "baby_tracker"),
		DBUser:         getEnv("DB_USER", "postgres"),
		DBPassword:     getEnv("DB_PASSWORD", ""),
		DBSSLMode:      getEnv("DB_SSLMODE", "disable"),
		SessionSecret:  getEnv("SESSION_SECRET", "change-me"),
		SessionMaxAge:  maxAge,
		AllowedOrigins: getEnv("ALLOWED_ORIGINS", "http://localhost:5173"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// Validate ensures required configuration values are properly set
func (c *Config) Validate() error {
	if c.SessionSecret == "change-me" && c.Env == "production" {
		return errors.New("SESSION_SECRET must be set to a secure value in production")
	}

	if c.DBType == "postgres" && c.DBPassword == "" && c.Env == "production" {
		return errors.New("DB_PASSWORD must be set for PostgreSQL in production")
	}

	// Validate SSL mode values
	validSSLModes := map[string]bool{
		"disable":     true,
		"require":     true,
		"verify-ca":   true,
		"verify-full": true,
	}
	if c.DBType == "postgres" && !validSSLModes[c.DBSSLMode] {
		return errors.New("DB_SSLMODE must be one of: disable, require, verify-ca, verify-full")
	}

	return nil
}
