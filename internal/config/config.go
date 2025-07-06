package config

import (
	"os"
	"strconv"
)

type Config struct {
	Port           string
	Env            string
	DBType         string
	DBPath         string
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
		DBPath:         getEnv("DB_PATH", "./baby-tracker.db"),
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
