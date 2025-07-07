package database

import (
	"fmt"

	"gorm.io/gorm"

	"github.com/engineervix/baby-tracker/internal/config"
	"github.com/engineervix/baby-tracker/internal/models"
)

// RunTestMigrations runs AutoMigrate for tests
func RunTestMigrations(db *gorm.DB, cfg *config.Config) error {
	// Debug: Ensure we're in test mode
	if cfg.Env != "test" {
		return fmt.Errorf("RunTestMigrations called with non-test environment: %s", cfg.Env)
	}

	// Enable foreign key constraints for SQLite
	if cfg.DBType == "sqlite" {
		if err := db.Exec("PRAGMA foreign_keys = ON").Error; err != nil {
			return err
		}
	}

	// Auto-migrate all models for tests
	err := db.AutoMigrate(
		&models.User{},
		&models.Baby{},
		&models.Session{},
		&models.Activity{},
		&models.FeedActivity{},
		&models.PumpActivity{},
		&models.DiaperActivity{},
		&models.SleepActivity{},
		&models.GrowthMeasurement{},
		&models.HealthRecord{},
		&models.Milestone{},
	)

	if err != nil {
		return fmt.Errorf("failed to run test migrations: %w", err)
	}

	return nil
}
