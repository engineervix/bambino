package database

import (
	"github.com/engineervix/baby-tracker/internal/config"
	"github.com/engineervix/baby-tracker/internal/models"
	"gorm.io/gorm"
)

// RunMigrations runs all database migrations
func RunMigrations(db *gorm.DB, cfg *config.Config) error {
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

	// TODO: Add any custom migrations here

	return nil
}
