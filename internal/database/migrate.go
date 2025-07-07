package database

import (
	"github.com/engineervix/baby-tracker/internal/config"
	"github.com/engineervix/baby-tracker/internal/models"
	"gorm.io/gorm"
)

// RunMigrations runs all database migrations
func RunMigrations(db *gorm.DB, cfg *config.Config) error {
	// Enable foreign key constraints for SQLite
	if cfg.DBType == "sqlite" {
		if err := db.Exec("PRAGMA foreign_keys = ON").Error; err != nil {
			return err
		}
	}

	// Enable UUID extension for PostgreSQL
	if cfg.DBType == "postgres" {
		if err := db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"").Error; err != nil {
			return err
		}
	}

	// Auto-migrate all models
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
		return err
	}

	// Add indexes for better performance
	if err := createIndexes(db); err != nil {
		return err
	}

	return nil
}

func createIndexes(db *gorm.DB) error {
	// Activity indexes
	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_activities_baby_id ON activities(baby_id)").Error; err != nil {
		return err
	}
	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_activities_type ON activities(type)").Error; err != nil {
		return err
	}
	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_activities_start_time ON activities(start_time)").Error; err != nil {
		return err
	}

	// Baby indexes
	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_babies_user_id ON babies(user_id)").Error; err != nil {
		return err
	}

	// Session indexes
	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_sessions_expires_at ON sessions(expires_at)").Error; err != nil {
		return err
	}

	return nil
}
