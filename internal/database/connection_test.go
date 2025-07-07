package database

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/engineervix/baby-tracker/internal/config"
)

func TestConnect_SQLite(t *testing.T) {
	// Create temp file for SQLite
	tmpfile, err := os.CreateTemp("", "test-*.db")
	require.NoError(t, err)
	defer os.Remove(tmpfile.Name())
	tmpfile.Close()

	cfg := &config.Config{
		DBType: "sqlite",
		DBPath: tmpfile.Name(),
	}

	db, err := Connect(cfg)
	require.NoError(t, err)
	assert.NotNil(t, db)

	// Test connection
	sqlDB, err := db.DB()
	require.NoError(t, err)

	err = sqlDB.Ping()
	assert.NoError(t, err)

	// Cleanup
	sqlDB.Close()
}

func TestConnect_InvalidDBType(t *testing.T) {
	cfg := &config.Config{
		DBType: "invalid",
	}

	db, err := Connect(cfg)
	assert.Nil(t, db)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unsupported database type")
}

func TestConnect_PostgreSQL(t *testing.T) {
	// Skip if not in integration test mode
	if os.Getenv("INTEGRATION_TEST") != "true" {
		t.Skip("Skipping PostgreSQL integration test")
	}

	cfg := &config.Config{
		DBType:     "postgres",
		DBHost:     os.Getenv("TEST_DB_HOST"),
		DBPort:     os.Getenv("TEST_DB_PORT"),
		DBName:     os.Getenv("TEST_DB_NAME"),
		DBUser:     os.Getenv("TEST_DB_USER"),
		DBPassword: os.Getenv("TEST_DB_PASSWORD"),
		DBSSLMode:  "disable",
	}

	db, err := Connect(cfg)
	require.NoError(t, err)
	assert.NotNil(t, db)

	// Test connection
	sqlDB, err := db.DB()
	require.NoError(t, err)

	err = sqlDB.Ping()
	assert.NoError(t, err)

	// Cleanup
	sqlDB.Close()
}
