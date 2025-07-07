package database

import (
	"fmt"
	"log"

	"github.com/engineervix/baby-tracker/internal/models"
)

// AutoMigrate runs automatic migrations for all models
func AutoMigrate() error {
	if DB == nil {
		return fmt.Errorf("database connection not initialized")
	}

	log.Println("Starting database migration...")

	// List of models to migrate in order
	modelsToMigrate := []interface{}{
		&models.User{},
		&models.Baby{},
		&models.Activity{},
		&models.FeedActivity{},
		&models.DiaperActivity{},
		&models.SleepActivity{},
		&models.PumpActivity{},
		&models.GrowthMeasurement{},
		&models.HealthRecord{},
		&models.Milestone{},
		&models.Session{},
	}

	// Run auto migration
	if err := DB.AutoMigrate(modelsToMigrate...); err != nil {
		return fmt.Errorf("failed to run auto migration: %w", err)
	}

	// Add any custom indexes or constraints here
	if err := addCustomConstraints(); err != nil {
		return fmt.Errorf("failed to add custom constraints: %w", err)
	}

	log.Println("Database migration completed successfully")
	return nil
}

// addCustomConstraints adds any custom database constraints not handled by GORM
func addCustomConstraints() error {
	// Add composite indexes for better query performance
	indexes := []struct {
		model  interface{}
		name   string
		fields []string
	}{
		{
			model:  &models.Activity{},
			name:   "idx_activities_baby_type",
			fields: []string{"baby_id", "type"},
		},
		{
			model:  &models.Activity{},
			name:   "idx_activities_baby_start",
			fields: []string{"baby_id", "start_time"},
		},
	}

	for _, idx := range indexes {
		for _, field := range idx.fields {
			if err := DB.Migrator().CreateIndex(idx.model, field); err != nil {
				// Ignore error if index already exists
				log.Printf("Note: Index on %s for %s might already exist: %v", field, idx.name, err)
			}
		}
	}

	// For PostgreSQL, add check constraints (SQLite doesn't support ALTER TABLE ADD CONSTRAINT)
	if DB.Dialector.Name() == "postgres" {
		constraints := []string{
			"ALTER TABLE sleep_activities ADD CONSTRAINT quality_check CHECK (quality >= 1 AND quality <= 5)",
		}

		for _, constraint := range constraints {
			if err := DB.Exec(constraint).Error; err != nil {
				// Ignore if constraint already exists
				log.Printf("Note: Constraint might already exist: %v", err)
			}
		}
	}

	return nil
}

// DropAllTables drops all tables (use with caution!)
func DropAllTables() error {
	if DB == nil {
		return fmt.Errorf("database connection not initialized")
	}

	log.Println("WARNING: Dropping all tables...")

	// Drop tables in reverse order to avoid foreign key constraints
	tablesToDrop := []interface{}{
		&models.Session{},
		&models.Milestone{},
		&models.HealthRecord{},
		&models.GrowthMeasurement{},
		&models.PumpActivity{},
		&models.SleepActivity{},
		&models.DiaperActivity{},
		&models.FeedActivity{},
		&models.Activity{},
		&models.Baby{},
		&models.User{},
	}

	for _, model := range tablesToDrop {
		if err := DB.Migrator().DropTable(model); err != nil {
			return fmt.Errorf("failed to drop table: %w", err)
		}
	}

	log.Println("All tables dropped successfully")
	return nil
}

// ResetDatabase drops all tables and runs migrations (use with caution!)
func ResetDatabase() error {
	if err := DropAllTables(); err != nil {
		return err
	}
	return AutoMigrate()
}
