package database

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/engineervix/baby-tracker/internal/config"
)

// createTestDB creates a test database without running migrations
func createTestDB(t *testing.T) (*gorm.DB, *config.Config, func()) {
	t.Helper()

	// Create temp file for SQLite
	tmpfile, err := os.CreateTemp("", "test-*.db")
	require.NoError(t, err)
	tmpfile.Close()

	// Create test config
	cfg := &config.Config{
		DBType: "sqlite",
		DBPath: tmpfile.Name(),
		Env:    "test",
	}

	// Open database with logging disabled for tests
	db, err := gorm.Open(sqlite.Open(tmpfile.Name()), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	require.NoError(t, err)

	cleanup := func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
		os.Remove(tmpfile.Name())
	}

	return db, cfg, cleanup
}
