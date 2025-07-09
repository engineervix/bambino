package cmd

import (
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/engineervix/baby-tracker/internal/config"
	"github.com/engineervix/baby-tracker/internal/database"
	"github.com/engineervix/baby-tracker/internal/models"
)

func setupTestDatabase(t *testing.T) (*gorm.DB, func()) {
	// Setup test database
	tmpfile, err := os.CreateTemp("", "test-*.db")
	require.NoError(t, err)
	tmpfile.Close()

	// Create test config
	cfg := &config.Config{
		DBType: "sqlite",
		DBPath: tmpfile.Name(),
		Env:    "test",
	}

	// Open database
	db, err := gorm.Open(sqlite.Open(tmpfile.Name()), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	require.NoError(t, err)

	// Run migrations
	err = database.RunMigrations(db, cfg)
	require.NoError(t, err)

	cleanup := func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
		os.Remove(tmpfile.Name())
	}

	return db, cleanup
}

func TestCreateUser_Integration(t *testing.T) {
	db, cleanup := setupTestDatabase(t)
	defer cleanup()

	// Test user creation
	username := "testparent"
	password := "testpassword123"
	babyName := "Test Baby"
	birthDate := time.Now().AddDate(0, 0, -7)

	// Create user directly (simulating what createUser does)
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	require.NoError(t, err)

	user := &models.User{
		Username:     username,
		PasswordHash: string(hash),
	}
	err = db.Create(user).Error
	require.NoError(t, err)

	baby := &models.Baby{
		UserID:    user.ID,
		Name:      babyName,
		BirthDate: birthDate,
	}
	err = db.Create(baby).Error
	require.NoError(t, err)

	// Verify user was created
	var loadedUser models.User
	err = db.Where("username = ?", username).First(&loadedUser).Error
	require.NoError(t, err)
	assert.Equal(t, username, loadedUser.Username)

	// Verify password hash
	err = bcrypt.CompareHashAndPassword([]byte(loadedUser.PasswordHash), []byte(password))
	assert.NoError(t, err)

	// Verify baby was created
	var loadedBaby models.Baby
	err = db.Where("user_id = ?", loadedUser.ID).First(&loadedBaby).Error
	require.NoError(t, err)
	assert.Equal(t, babyName, loadedBaby.Name)
	assert.Equal(t, loadedUser.ID, loadedBaby.UserID)
	assert.WithinDuration(t, birthDate, loadedBaby.BirthDate, 24*time.Hour)

	// Test duplicate username
	t.Run("duplicate username", func(t *testing.T) {
		duplicateUser := &models.User{
			Username:     username,
			PasswordHash: "anotherhash",
		}
		err = db.Create(duplicateUser).Error
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "UNIQUE")
	})
}

func TestUserBabyRelationship_Integration(t *testing.T) {
	db, cleanup := setupTestDatabase(t)
	defer cleanup()

	// Create user
	user := &models.User{
		Username:     "parent",
		PasswordHash: "$2a$10$test",
	}
	err := db.Create(user).Error
	require.NoError(t, err)

	// Create multiple babies
	babies := []models.Baby{
		{
			UserID:    user.ID,
			Name:      "First Baby",
			BirthDate: time.Now().AddDate(-1, 0, 0),
		},
		{
			UserID:    user.ID,
			Name:      "Second Baby",
			BirthDate: time.Now().AddDate(0, -6, 0),
		},
	}

	for i := range babies {
		babies[i].ID = uuid.New()
		err = db.Create(&babies[i]).Error
		require.NoError(t, err)
	}

	// Test loading user with babies
	var loadedUser models.User
	err = db.Preload("Babies").First(&loadedUser, user.ID).Error
	require.NoError(t, err)
	assert.Len(t, loadedUser.Babies, 2)

	// Test cascade delete
	err = db.Delete(&user).Error
	require.NoError(t, err)

	// Verify babies are deleted
	var babyCount int64
	err = db.Model(&models.Baby{}).Where("user_id = ?", user.ID).Count(&babyCount).Error
	require.NoError(t, err)
	assert.Equal(t, int64(0), babyCount)
}

func TestDatabaseMigration_Integration(t *testing.T) {
	// This tests that migrations can run multiple times without error
	tmpfile, err := os.CreateTemp("", "test-*.db")
	require.NoError(t, err)
	defer os.Remove(tmpfile.Name())
	tmpfile.Close()

	cfg := &config.Config{
		DBType: "sqlite",
		DBPath: tmpfile.Name(),
		Env:    "test",
	}

	db, err := database.Connect(cfg)
	require.NoError(t, err)
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}()

	// Run migrations multiple times
	for i := 0; i < 3; i++ {
		err = database.RunMigrations(db, cfg)
		assert.NoError(t, err, "Migration run %d should succeed", i+1)
	}

	// Verify all tables exist
	tables := []string{
		"users", "babies", "activities",
		"feed_activities", "pump_activities", "diaper_activities",
		"sleep_activities", "growth_measurements", "health_records", "milestones",
	}

	for _, table := range tables {
		var exists bool
		err = db.Raw("SELECT EXISTS (SELECT 1 FROM sqlite_master WHERE type='table' AND name=?)", table).Scan(&exists).Error
		require.NoError(t, err)
		assert.True(t, exists, "Table %s should exist", table)
	}
}
