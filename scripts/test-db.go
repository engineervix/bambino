package main

import (
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"

	"github.com/engineervix/baby-tracker/internal/database"
	"github.com/engineervix/baby-tracker/internal/models"
)

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
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

	log.Println("Database connection and migration successful!")

	// Test creating a user
	if err := testCreateUser(); err != nil {
		log.Printf("Error testing user creation: %v", err)
	}

	// Test creating activities
	if err := testCreateActivities(); err != nil {
		log.Printf("Error testing activity creation: %v", err)
	}

	log.Println("All tests completed!")
}

func testCreateUser() error {
	log.Println("Testing user creation...")

	// Create a test user
	password := "testpassword123"
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	user := &models.User{
		Username:     fmt.Sprintf("testuser_%d", time.Now().Unix()),
		PasswordHash: string(hash),
	}

	if err := database.DB.Create(user).Error; err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	log.Printf("Created user: %s (ID: %s)", user.Username, user.ID)

	// Create a baby for this user
	baby := &models.Baby{
		UserID:      user.ID,
		Name:        "Test Baby",
		BirthDate:   time.Now().AddDate(0, -6, 0), // 6 months old
		BirthWeight: ptr(3.5),                     // 3.5 kg
		BirthHeight: ptr(50.0),                    // 50 cm
	}

	if err := database.DB.Create(baby).Error; err != nil {
		return fmt.Errorf("failed to create baby: %w", err)
	}

	log.Printf("Created baby: %s (ID: %s)", baby.Name, baby.ID)

	// Test some activities
	return testActivitiesForBaby(baby.ID)
}

func testActivitiesForBaby(babyID uuid.UUID) error {
	log.Println("Testing activity creation...")

	// Create a feed activity
	feedActivity := &models.Activity{
		BabyID:    babyID,
		Type:      models.ActivityTypeFeed,
		StartTime: time.Now().Add(-30 * time.Minute),
		EndTime:   ptr(time.Now().Add(-10 * time.Minute)),
		Notes:     "Test feeding session",
	}

	if err := database.DB.Create(feedActivity).Error; err != nil {
		return fmt.Errorf("failed to create feed activity: %w", err)
	}

	// Create feed details
	feedDetails := &models.FeedActivity{
		ActivityID:      feedActivity.ID,
		FeedType:        models.FeedTypeBottle,
		AmountML:        ptr(120.0),
		DurationMinutes: ptr(20),
	}

	if err := database.DB.Create(feedDetails).Error; err != nil {
		return fmt.Errorf("failed to create feed details: %w", err)
	}

	log.Printf("Created feed activity: %s", feedActivity.ID)

	// Create a diaper activity
	diaperActivity := &models.Activity{
		BabyID:    babyID,
		Type:      models.ActivityTypeDiaper,
		StartTime: time.Now().Add(-5 * time.Minute),
		Notes:     "Diaper change",
	}

	if err := database.DB.Create(diaperActivity).Error; err != nil {
		return fmt.Errorf("failed to create diaper activity: %w", err)
	}

	diaperDetails := &models.DiaperActivity{
		ActivityID: diaperActivity.ID,
		Wet:        true,
		Dirty:      true,
		Color:      "yellow",
	}

	if err := database.DB.Create(diaperDetails).Error; err != nil {
		return fmt.Errorf("failed to create diaper details: %w", err)
	}

	log.Printf("Created diaper activity: %s", diaperActivity.ID)

	// Query activities
	var activities []models.Activity
	if err := database.DB.Where("baby_id = ?", babyID).Find(&activities).Error; err != nil {
		return fmt.Errorf("failed to query activities: %w", err)
	}

	log.Printf("Found %d activities for baby", len(activities))

	return nil
}

func testCreateActivities() error {
	// Get a user to work with
	var user models.User
	if err := database.DB.First(&user).Error; err != nil {
		log.Println("No existing user found, skipping activity tests")
		return nil
	}

	// Get their baby
	var baby models.Baby
	if err := database.DB.Where("user_id = ?", user.ID).First(&baby).Error; err != nil {
		log.Println("No baby found for user, skipping activity tests")
		return nil
	}

	log.Printf("Using existing baby: %s", baby.Name)

	// Create various activity types
	activities := []struct {
		activity interface{}
		details  interface{}
	}{
		{
			activity: &models.Activity{
				BabyID:    baby.ID,
				Type:      models.ActivityTypeSleep,
				StartTime: time.Now().Add(-2 * time.Hour),
				EndTime:   ptr(time.Now()),
				Notes:     "Afternoon nap",
			},
			details: &models.SleepActivity{
				Location: "crib",
				Quality:  ptr(4),
			},
		},
		{
			activity: &models.Activity{
				BabyID:    baby.ID,
				Type:      models.ActivityTypeGrowth,
				StartTime: time.Now(),
				Notes:     "Monthly measurement",
			},
			details: &models.GrowthMeasurement{
				WeightKG:            ptr(7.5),
				HeightCM:            ptr(65.0),
				HeadCircumferenceCM: ptr(42.0),
			},
		},
	}

	for _, item := range activities {
		activity := item.activity.(*models.Activity)

		if err := database.DB.Create(activity).Error; err != nil {
			return fmt.Errorf("failed to create %s activity: %w", activity.Type, err)
		}

		// Set the activity ID for the details
		switch details := item.details.(type) {
		case *models.SleepActivity:
			details.ActivityID = activity.ID
		case *models.GrowthMeasurement:
			details.ActivityID = activity.ID
		}

		if err := database.DB.Create(item.details).Error; err != nil {
			return fmt.Errorf("failed to create %s details: %w", activity.Type, err)
		}

		log.Printf("Created %s activity: %s", activity.Type, activity.ID)
	}

	return nil
}

// Helper function to get pointer to value
func ptr[T any](v T) *T {
	return &v
}
