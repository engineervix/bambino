package cmd

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/spf13/cobra"

	"github.com/engineervix/bambino/internal/config"
	"github.com/engineervix/bambino/internal/database"
	"github.com/engineervix/bambino/internal/models"
	"github.com/engineervix/bambino/internal/utils"
)

var seedCmd = &cobra.Command{
	Use:   "seed",
	Short: "Seed the database with development data",
	Long: `Seeds the database with minimal test data for development purposes.
This command is intended for development use only and should not be run in production.`,
	Run: func(cmd *cobra.Command, args []string) {
		seedDatabase()
	},
}

func init() {
	rootCmd.AddCommand(seedCmd)
}

func seedDatabase() {
	// Check if we're in development
	if os.Getenv("ENV") == "production" {
		log.Fatal("❌ Seed command is not allowed in production environment")
	}

	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
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

	fmt.Println("🌱 Seeding database with development data...")

	// Create test user
	hashedPassword, err := utils.HashPassword("password123")
	if err != nil {
		log.Fatalf("Failed to hash password: %v", err)
	}

	testUser := &models.User{
		Username:     "testuser",
		PasswordHash: hashedPassword,
	}

	// Check if user already exists
	var existingUser models.User
	result := db.Where("username = ?", testUser.Username).First(&existingUser)
	if result.Error == nil {
		fmt.Println("⚠️  Test user already exists, skipping user creation")
		testUser = &existingUser
	} else {
		if err := db.Create(testUser).Error; err != nil {
			log.Fatalf("Failed to create test user: %v", err)
		}
		fmt.Printf("✅ Created test user: %s\n", testUser.Username)
	}

	// Create test baby
	birthDate := time.Now().AddDate(0, -6, 0) // 6 months ago
	birthWeight := 3.2                        // kg
	birthHeight := 50.0                       // cm

	testBaby := &models.Baby{
		UserID:      testUser.ID,
		Name:        "Test Baby",
		BirthDate:   birthDate,
		TrackSleep:  true,
		BirthWeight: &birthWeight,
		BirthHeight: &birthHeight,
	}

	// Check if baby already exists
	var existingBaby models.Baby
	result = db.Where("user_id = ? AND name = ?", testUser.ID, testBaby.Name).First(&existingBaby)
	if result.Error == nil {
		fmt.Println("⚠️  Test baby already exists, skipping baby creation")
		testBaby = &existingBaby
	} else {
		if err := db.Create(testBaby).Error; err != nil {
			log.Fatalf("Failed to create test baby: %v", err)
		}
		fmt.Printf("✅ Created test baby: %s\n", testBaby.Name)
	}

	// Create sample activities
	now := time.Now()

	// Feed activity
	feedActivity := &models.Activity{
		BabyID:    testBaby.ID,
		Type:      models.ActivityTypeFeed,
		StartTime: now.Add(-2 * time.Hour),
		EndTime:   &[]time.Time{now.Add(-90 * time.Minute)}[0],
		Notes:     "Morning feeding",
	}

	if err := db.Create(feedActivity).Error; err != nil {
		log.Printf("Warning: Failed to create feed activity: %v", err)
	} else {
		// Create feed details
		feedDetails := &models.FeedActivity{
			ActivityID:      feedActivity.ID,
			FeedType:        models.FeedTypeBottle,
			AmountML:        &[]float64{120.0}[0],
			DurationMinutes: &[]int{30}[0],
		}
		if err := db.Create(feedDetails).Error; err != nil {
			log.Printf("Warning: Failed to create feed details: %v", err)
		} else {
			fmt.Println("✅ Created sample feed activity")
		}
	}

	// Diaper activity
	diaperActivity := &models.Activity{
		BabyID:    testBaby.ID,
		Type:      models.ActivityTypeDiaper,
		StartTime: now.Add(-3 * time.Hour),
		Notes:     "Wet diaper",
	}

	if err := db.Create(diaperActivity).Error; err != nil {
		log.Printf("Warning: Failed to create diaper activity: %v", err)
	} else {
		// Create diaper details
		diaperDetails := &models.DiaperActivity{
			ActivityID: diaperActivity.ID,
			Wet:        true,
			Dirty:      false,
		}
		if err := db.Create(diaperDetails).Error; err != nil {
			log.Printf("Warning: Failed to create diaper details: %v", err)
		} else {
			fmt.Println("✅ Created sample diaper activity")
		}
	}

	// Sleep activity
	sleepActivity := &models.Activity{
		BabyID:    testBaby.ID,
		Type:      models.ActivityTypeSleep,
		StartTime: now.Add(-5 * time.Hour),
		EndTime:   &[]time.Time{now.Add(-3 * time.Hour)}[0],
		Notes:     "Afternoon nap",
	}

	if err := db.Create(sleepActivity).Error; err != nil {
		log.Printf("Warning: Failed to create sleep activity: %v", err)
	} else {
		// Create sleep details
		sleepDetails := &models.SleepActivity{
			ActivityID: sleepActivity.ID,
			Location:   "crib",
			Quality:    &[]int{4}[0], // Good quality (1-5 scale)
		}
		if err := db.Create(sleepDetails).Error; err != nil {
			log.Printf("Warning: Failed to create sleep details: %v", err)
		} else {
			fmt.Println("✅ Created sample sleep activity")
		}
	}

	fmt.Println("\n🎉 Database seeding completed!")
	fmt.Println("\nTest credentials:")
	fmt.Printf("  Username: %s\n", testUser.Username)
	fmt.Println("  Password: password123")
	fmt.Printf("  Baby: %s (born %s)\n", testBaby.Name, testBaby.BirthDate.Format("2006-01-02"))
}
