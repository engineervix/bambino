package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoad(t *testing.T) {
	// Save current env and restore after test
	originalEnv := os.Environ()
	defer func() {
		os.Clearenv()
		for _, env := range originalEnv {
			pair := splitEnvVar(env)
			os.Setenv(pair[0], pair[1])
		}
	}()

	tests := []struct {
		name     string
		envVars  map[string]string
		expected *Config
	}{
		{
			name:    "default values",
			envVars: map[string]string{},
			expected: &Config{
				Port:                   "8080",
				Env:                    "development",
				DBType:                 "sqlite",
				DBPath:                 "./bambino.db",
				DBHost:                 "localhost",
				DBPort:                 "5432",
				DBName:                 "baby",
				DBUser:                 "postgres",
				DBPassword:             "",
				DBSSLMode:              "disable",
				SessionSecret:          "change-me",
				SessionMaxAge:          86400,
				AllowedOrigins:         "http://localhost:5173",
				SentryDSN:              "",
				SentryTracesSampleRate: 0.1,
			},
		},
		{
			name: "custom values",
			envVars: map[string]string{
				"PORT":            "3000",
				"ENV":             "production",
				"DB_TYPE":         "postgres",
				"DB_HOST":         "db.example.com",
				"DB_PORT":         "5433",
				"DB_PASSWORD":     "secret",
				"DB_SSLMODE":      "require",
				"SESSION_SECRET":  "very-secret-key",
				"SESSION_MAX_AGE": "3600",
			},
			expected: &Config{
				Port:                   "3000",
				Env:                    "production",
				DBType:                 "postgres",
				DBPath:                 "./bambino.db",
				DBHost:                 "db.example.com",
				DBPort:                 "5433",
				DBName:                 "baby",
				DBUser:                 "postgres",
				DBPassword:             "secret",
				DBSSLMode:              "require",
				SessionSecret:          "very-secret-key",
				SessionMaxAge:          3600,
				AllowedOrigins:         "http://localhost:5173",
				SentryDSN:              "",
				SentryTracesSampleRate: 0.1,
			},
		},
		{
			name: "sentry configuration",
			envVars: map[string]string{
				"SENTRY_DSN":                "https://example@sentry.io/123",
				"SENTRY_TRACES_SAMPLE_RATE": "0.5",
			},
			expected: &Config{
				Port:                   "8080",
				Env:                    "development",
				DBType:                 "sqlite",
				DBPath:                 "./bambino.db",
				DBHost:                 "localhost",
				DBPort:                 "5432",
				DBName:                 "baby",
				DBUser:                 "postgres",
				DBPassword:             "",
				DBSSLMode:              "disable",
				SessionSecret:          "change-me",
				SessionMaxAge:          86400,
				AllowedOrigins:         "http://localhost:5173",
				SentryDSN:              "https://example@sentry.io/123",
				SentryTracesSampleRate: 0.5,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clear env
			os.Clearenv()

			// Set test env vars
			for k, v := range tt.envVars {
				os.Setenv(k, v)
			}

			// Load config
			cfg := Load()

			// Assert
			assert.Equal(t, tt.expected, cfg)
		})
	}
}

func TestValidate(t *testing.T) {
	tests := []struct {
		name    string
		config  *Config
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid development config",
			config: &Config{
				Env:           "development",
				DBType:        "sqlite",
				SessionSecret: "change-me",
			},
			wantErr: false,
		},
		{
			name: "invalid session secret in production",
			config: &Config{
				Env:           "production",
				DBType:        "sqlite",
				SessionSecret: "change-me",
			},
			wantErr: true,
			errMsg:  "SESSION_SECRET must be set to a secure value in production",
		},
		{
			name: "missing postgres password in production",
			config: &Config{
				Env:           "production",
				DBType:        "postgres",
				DBPassword:    "",
				SessionSecret: "secure-secret",
			},
			wantErr: true,
			errMsg:  "DB_PASSWORD must be set for PostgreSQL in production",
		},
		{
			name: "invalid SSL mode",
			config: &Config{
				Env:           "development",
				DBType:        "postgres",
				DBSSLMode:     "invalid",
				SessionSecret: "secret",
			},
			wantErr: true,
			errMsg:  "DB_SSLMODE must be one of",
		},
		{
			name: "valid postgres config",
			config: &Config{
				Env:           "production",
				DBType:        "postgres",
				DBPassword:    "password",
				DBSSLMode:     "require",
				SessionSecret: "secure-secret",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestGetEnv(t *testing.T) {
	// Save and restore env
	originalValue := os.Getenv("TEST_VAR")
	defer os.Setenv("TEST_VAR", originalValue)

	// Test with env var set
	os.Setenv("TEST_VAR", "custom")
	assert.Equal(t, "custom", getEnv("TEST_VAR", "default"))

	// Test with env var unset
	os.Unsetenv("TEST_VAR")
	assert.Equal(t, "default", getEnv("TEST_VAR", "default"))
}

// Helper function to split environment variable
func splitEnvVar(env string) []string {
	for i := 0; i < len(env); i++ {
		if env[i] == '=' {
			return []string{env[:i], env[i+1:]}
		}
	}
	return []string{env, ""}
}
