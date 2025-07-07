package database

import (
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// Config holds database configuration
type Config struct {
	Type     string
	Path     string // For SQLite
	Host     string // For PostgreSQL
	Port     string
	Name     string
	User     string
	Password string
	SSLMode  string
}

// NewConfig creates a database configuration from environment variables
func NewConfig() *Config {
	dbType := os.Getenv("DB_TYPE")
	if dbType == "" {
		dbType = "sqlite"
	}

	config := &Config{
		Type: dbType,
	}

	switch dbType {
	case "sqlite":
		config.Path = os.Getenv("DB_PATH")
		if config.Path == "" {
			config.Path = "./baby-tracker.db"
		}
	case "postgres":
		config.Host = getEnv("DB_HOST", "localhost")
		config.Port = getEnv("DB_PORT", "5432")
		config.Name = getEnv("DB_NAME", "baby_tracker")
		config.User = getEnv("DB_USER", "postgres")
		config.Password = os.Getenv("DB_PASSWORD")
		config.SSLMode = getEnv("DB_SSLMODE", "disable")
	}

	return config
}

// Connect establishes a database connection
func Connect(config *Config) error {
	var err error
	var dialector gorm.Dialector

	// Configure logger
	logLevel := logger.Silent
	if os.Getenv("ENV") == "development" {
		logLevel = logger.Info
	}

	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
	}

	switch config.Type {
	case "sqlite":
		dialector = sqlite.Open(config.Path)
		log.Printf("Connecting to SQLite database: %s", config.Path)
	case "postgres":
		dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
			config.Host, config.Port, config.User, config.Password, config.Name, config.SSLMode)
		dialector = postgres.Open(dsn)
		log.Printf("Connecting to PostgreSQL database: %s@%s:%s/%s", config.User, config.Host, config.Port, config.Name)
	default:
		return fmt.Errorf("unsupported database type: %s", config.Type)
	}

	DB, err = gorm.Open(dialector, gormConfig)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	// Get underlying SQL database to configure connection pool
	sqlDB, err := DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get database instance: %w", err)
	}

	// Set connection pool settings
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	// Test the connection
	if err := sqlDB.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	log.Println("Database connection established successfully")
	return nil
}

// Close closes the database connection
func Close() error {
	if DB != nil {
		sqlDB, err := DB.DB()
		if err != nil {
			return err
		}
		return sqlDB.Close()
	}
	return nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
