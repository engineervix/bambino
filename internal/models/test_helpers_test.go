package models_test

import (
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/engineervix/bambino/internal/config"
	"github.com/engineervix/bambino/internal/database"
	"github.com/engineervix/bambino/internal/models"
)

// TestDB holds the test database instance
type TestDB struct {
	DB      *gorm.DB
	Config  *config.Config
	Cleanup func()
}

// setupTestDB creates a test database
func setupTestDB(t *testing.T) *TestDB {
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

	// Run migrations
	err = database.RunMigrations(db, cfg)
	require.NoError(t, err)

	return &TestDB{
		DB:     db,
		Config: cfg,
		Cleanup: func() {
			sqlDB, _ := db.DB()
			sqlDB.Close()
			os.Remove(tmpfile.Name())
		},
	}
}

// createTestUser creates a user for testing
func createTestUser(t *testing.T, db *gorm.DB) *models.User {
	t.Helper()

	user := &models.User{
		ID:           uuid.New(),
		Username:     "testuser_" + uuid.New().String()[:8],
		PasswordHash: "$2a$10$test.hash", // Not a real hash, just for testing
	}

	err := db.Create(user).Error
	require.NoError(t, err)

	return user
}

// createTestBaby creates a baby for testing
func createTestBaby(t *testing.T, db *gorm.DB, userID uuid.UUID) *models.Baby {
	t.Helper()

	baby := &models.Baby{
		ID:        uuid.New(),
		UserID:    userID,
		Name:      "Test Baby",
		BirthDate: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
	}

	err := db.Create(baby).Error
	require.NoError(t, err)

	return baby
}

// Helper functions for pointers
func floatPtr(f float64) *float64 {
	return &f
}

func intPtr(i int) *int {
	return &i
}

func timePtr(t time.Time) *time.Time {
	return &t
}
