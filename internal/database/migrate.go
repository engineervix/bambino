package database

import (
	"database/sql"
	"embed"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"gorm.io/gorm"

	"github.com/engineervix/bambino/internal/config"
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

// RunMigrations runs all pending migrations
func RunMigrations(db *gorm.DB, cfg *config.Config) error {
	// For tests, always use AutoMigrate as fallback
	if cfg.Env == "test" {
		return RunTestMigrations(db, cfg)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to get database instance: %w", err)
	}

	m, err := createMigrator(sqlDB, cfg)
	if err != nil {
		return fmt.Errorf("failed to create migrator: %w", err)
	}
	// Don't close the migrator as it closes the shared database connection
	// defer m.Close()

	// Run all pending migrations
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	return nil
}

// MigrateDown rolls back one migration
func MigrateDown(db *gorm.DB, cfg *config.Config) error {
	// For tests, this operation doesn't make sense with AutoMigrate
	if cfg.Env == "test" {
		return fmt.Errorf("migrate down not supported in test mode")
	}

	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to get database instance: %w", err)
	}

	m, err := createMigrator(sqlDB, cfg)
	if err != nil {
		return fmt.Errorf("failed to create migrator: %w", err)
	}
	// Don't close the migrator as it closes the shared database connection
	// defer m.Close()

	// Roll back one migration
	if err := m.Steps(-1); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to rollback migration: %w", err)
	}

	return nil
}

// MigrateStatus shows current migration status
func MigrateStatus(db *gorm.DB, cfg *config.Config) (uint, bool, error) {
	// For tests, return fake status
	if cfg.Env == "test" {
		return 1, false, nil
	}

	sqlDB, err := db.DB()
	if err != nil {
		return 0, false, fmt.Errorf("failed to get database instance: %w", err)
	}

	m, err := createMigrator(sqlDB, cfg)
	if err != nil {
		return 0, false, fmt.Errorf("failed to create migrator: %w", err)
	}
	// Don't close the migrator as it closes the shared database connection
	// defer m.Close()

	version, dirty, err := m.Version()
	if err != nil {
		return 0, false, fmt.Errorf("failed to get migration version: %w", err)
	}

	return version, dirty, nil
}

// createMigrator creates a new migrator instance
func createMigrator(sqlDB *sql.DB, cfg *config.Config) (*migrate.Migrate, error) {
	// Create source from embedded filesystem
	source, err := iofs.New(migrationsFS, "migrations")
	if err != nil {
		return nil, fmt.Errorf("failed to create migration source: %w", err)
	}

	// Create database driver
	var driver database.Driver
	switch cfg.DBType {
	case "sqlite":
		driver, err = sqlite3.WithInstance(sqlDB, &sqlite3.Config{})
	case "postgres":
		driver, err = postgres.WithInstance(sqlDB, &postgres.Config{})
	default:
		return nil, fmt.Errorf("unsupported database type: %s", cfg.DBType)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to create database driver: %w", err)
	}

	// Create migrator
	m, err := migrate.NewWithInstance("iofs", source, cfg.DBType, driver)
	if err != nil {
		return nil, fmt.Errorf("failed to create migrator: %w", err)
	}

	return m, nil
}
