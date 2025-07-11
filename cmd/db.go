package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/spf13/cobra"

	"github.com/engineervix/bambino/internal/config"
	"github.com/engineervix/bambino/internal/database"
)

var dbCmd = &cobra.Command{
	Use:   "db",
	Short: "Database management commands",
	Long:  `Commands for managing the database.`,
}

var testDbCmd = &cobra.Command{
	Use:   "test",
	Short: "Test database connection",
	Long:  `Tests the database connection and displays connection info.`,
	Run: func(cmd *cobra.Command, args []string) {
		testDatabaseConnection()
	},
}

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Run database migrations",
	Long:  `Runs all pending database migrations.`,
	Run: func(cmd *cobra.Command, args []string) {
		runMigrations()
	},
}

var migrateDownCmd = &cobra.Command{
	Use:   "migrate-down",
	Short: "Rollback one migration",
	Long:  `Rolls back the most recent migration.`,
	Run: func(cmd *cobra.Command, args []string) {
		rollbackMigration()
	},
}

var migrateStatusCmd = &cobra.Command{
	Use:   "migrate-status",
	Short: "Show migration status",
	Long:  `Shows current migration version and status.`,
	Run: func(cmd *cobra.Command, args []string) {
		showMigrationStatus()
	},
}

func init() {
	rootCmd.AddCommand(dbCmd)
	dbCmd.AddCommand(testDbCmd)
	dbCmd.AddCommand(migrateCmd)
	dbCmd.AddCommand(migrateDownCmd)
	dbCmd.AddCommand(migrateStatusCmd)
}

func testDatabaseConnection() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		// Only log in development - production uses Docker env vars
		if os.Getenv("ENV") != "production" {
			log.Println("No .env file found")
		}
	}

	// Load configuration
	cfg := config.Load()

	// Validate configuration
	if err := cfg.Validate(); err != nil {
		log.Fatalf("Configuration error: %v", err)
	}

	fmt.Printf("Testing database connection...\n")
	fmt.Printf("Database Type: %s\n", cfg.DBType)

	if cfg.DBType == "sqlite" {
		fmt.Printf("Database Path: %s\n", cfg.DBPath)
	} else {
		fmt.Printf("Database Host: %s:%s\n", cfg.DBHost, cfg.DBPort)
		fmt.Printf("Database Name: %s\n", cfg.DBName)
		fmt.Printf("Database User: %s\n", cfg.DBUser)
		fmt.Printf("SSL Mode: %s\n", cfg.DBSSLMode)
	}

	// Connect to database
	db, err := database.Connect(cfg)
	if err != nil {
		log.Fatalf("❌ Failed to connect to database: %v", err)
	}

	// Test connection
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("❌ Failed to get database instance: %v", err)
	}

	if err := sqlDB.Ping(); err != nil {
		log.Fatalf("❌ Failed to ping database: %v", err)
	}

	fmt.Println("✅ Database connection successful!")

	// Show table information
	var tables []string
	if cfg.DBType == "sqlite" {
		db.Raw("SELECT name FROM sqlite_master WHERE type='table' AND name NOT LIKE 'sqlite_%' AND name != 'schema_migrations'").Pluck("name", &tables)
	} else {
		db.Raw("SELECT table_name FROM information_schema.tables WHERE table_schema = 'public' AND table_name != 'schema_migrations'").Pluck("table_name", &tables)
	}

	if len(tables) > 0 {
		fmt.Printf("\nExisting tables:\n")
		for _, table := range tables {
			fmt.Printf("  - %s\n", table)
		}
	} else {
		fmt.Println("\nNo tables found. Run 'bambino db migrate' to create tables.")
	}
}

func runMigrations() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		// Only log in development - production uses Docker env vars
		if os.Getenv("ENV") != "production" {
			log.Println("No .env file found")
		}
	}

	// Load configuration
	cfg := config.Load()

	// Validate configuration
	if err := cfg.Validate(); err != nil {
		log.Fatalf("Configuration error: %v", err)
	}

	// Connect to database
	db, err := database.Connect(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	fmt.Println("Running database migrations...")

	// Run migrations
	if err := database.RunMigrations(db, cfg); err != nil {
		log.Fatalf("❌ Failed to run migrations: %v", err)
	}

	fmt.Println("✅ Migrations completed successfully!")
}

func rollbackMigration() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		// Only log in development - production uses Docker env vars
		if os.Getenv("ENV") != "production" {
			log.Println("No .env file found")
		}
	}

	// Load configuration
	cfg := config.Load()

	// Validate configuration
	if err := cfg.Validate(); err != nil {
		log.Fatalf("Configuration error: %v", err)
	}

	// Connect to database
	db, err := database.Connect(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	fmt.Println("Rolling back migration...")

	// Rollback migration
	if err := database.MigrateDown(db, cfg); err != nil {
		log.Fatalf("❌ Failed to rollback migration: %v", err)
	}

	fmt.Println("✅ Migration rolled back successfully!")
}

func showMigrationStatus() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		// Only log in development - production uses Docker env vars
		if os.Getenv("ENV") != "production" {
			log.Println("No .env file found")
		}
	}

	// Load configuration
	cfg := config.Load()

	// Validate configuration
	if err := cfg.Validate(); err != nil {
		log.Fatalf("Configuration error: %v", err)
	}

	// Connect to database
	db, err := database.Connect(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Get migration status
	version, dirty, err := database.MigrateStatus(db, cfg)
	if err != nil {
		log.Fatalf("❌ Failed to get migration status: %v", err)
	}

	fmt.Printf("Migration Status:\n")
	fmt.Printf("  Current Version: %d\n", version)
	fmt.Printf("  Dirty: %t\n", dirty)

	if dirty {
		fmt.Println("⚠️  Warning: Database is in dirty state. Manual intervention may be required.")
	} else {
		fmt.Println("✅ Database is clean")
	}
}
