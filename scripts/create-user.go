package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"

	"github.com/engineervix/baby-tracker/internal/database"
	"github.com/engineervix/baby-tracker/internal/models"
)

func main() {
	var username, password, babyName string
	var babyBirthDate string

	flag.StringVar(&username, "username", "", "Username for the new user")
	flag.StringVar(&password, "password", "", "Password for the new user")
	flag.StringVar(&babyName, "baby-name", "", "Baby's name (optional)")
	flag.StringVar(&babyBirthDate, "baby-birth-date", "", "Baby's birth date YYYY-MM-DD (optional)")
	flag.Parse()

	if username == "" || password == "" {
		fmt.Println("Usage: go run create-user.go --username=<username> --password=<password> [--baby-name=<name>] [--baby-birth-date=YYYY-MM-DD]")
		fmt.Println("\nExample:")
		fmt.Println("  go run create-user.go --username=parent --password=secure123 --baby-name=\"Emma\" --baby-birth-date=2024-01-15")
		os.Exit(1)
	}

	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	// Connect to database
	dbConfig := database.NewConfig()
	if err := database.Connect(dbConfig); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close()

	// Run migrations
	if err := database.AutoMigrate(); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Check if user already exists
	var existingUser models.User
	if err := database.DB.Where("username = ?", username).First(&existingUser).Error; err == nil {
		log.Fatalf("User with username '%s' already exists", username)
	}

	// Hash password
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("Failed to hash password: %v", err)
	}

	// Create user
	user := &models.User{
		Username:     username,
		PasswordHash: string(hash),
	}

	if err := database.DB.Create(user).Error; err != nil {
		log.Fatalf("Failed to create user: %v", err)
	}

	fmt.Printf("✅ Created user: %s (ID: %s)\n", user.Username, user.ID)

	// Create baby profile if requested
	if babyName != "" {
		baby := &models.Baby{
			UserID: user.ID,
			Name:   babyName,
		}

		// Parse birth date if provided
		if babyBirthDate != "" {
			birthDate, err := time.Parse("2006-01-02", babyBirthDate)
			if err != nil {
				log.Printf("Warning: Invalid birth date format '%s', using today's date", babyBirthDate)
				baby.BirthDate = time.Now()
			} else {
				baby.BirthDate = birthDate
			}
		} else {
			baby.BirthDate = time.Now()
		}

		if err := database.DB.Create(baby).Error; err != nil {
			log.Printf("Warning: Failed to create baby profile: %v", err)
		} else {
			fmt.Printf("✅ Created baby profile: %s (Born: %s)\n", baby.Name, baby.BirthDate.Format("2006-01-02"))
		}
	}

	fmt.Println("\n🎉 Setup complete! You can now log in with your credentials.")
}
