package cmd

import (
	"fmt"
	"log"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
	"golang.org/x/term"

	"github.com/engineervix/bambino/internal/config"
	"github.com/engineervix/bambino/internal/database"
	"github.com/engineervix/bambino/internal/models"
	"github.com/engineervix/bambino/internal/utils"
)

var createUserCmd = &cobra.Command{
	Use:   "create-user",
	Short: "Create a new user",
	Long:  `Creates a new user account with the specified username.`,
	Run: func(cmd *cobra.Command, args []string) {
		username, _ := cmd.Flags().GetString("username")
		babyName, _ := cmd.Flags().GetString("baby-name")
		birthDateStr, _ := cmd.Flags().GetString("birth-date")

		if username == "" {
			log.Fatal("Username is required")
		}

		// Parse birth date if provided
		var birthDate time.Time
		if birthDateStr != "" {
			var err error
			birthDate, err = time.Parse("2006-01-02", birthDateStr)
			if err != nil {
				log.Fatalf("Invalid birth date format. Use YYYY-MM-DD: %v", err)
			}
		} else {
			// Default to 1 week ago
			birthDate = time.Now().AddDate(0, 0, -7)
			fmt.Printf("No birth date provided, using default: %s\n", birthDate.Format("2006-01-02"))
		}

		// Prompt for password
		fmt.Print("Enter password: ")
		password, err := term.ReadPassword(int(syscall.Stdin))
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println()

		fmt.Print("Confirm password: ")
		confirmPassword, err := term.ReadPassword(int(syscall.Stdin))
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println()

		if string(password) != string(confirmPassword) {
			log.Fatal("Passwords do not match")
		}

		createUser(username, string(password), babyName, birthDate)
	},
}

func init() {
	rootCmd.AddCommand(createUserCmd)

	createUserCmd.Flags().StringP("username", "u", "", "Username for the new user (required)")
	createUserCmd.Flags().StringP("baby-name", "b", "Baby", "Name of the baby")
	createUserCmd.Flags().StringP("birth-date", "d", "", "Birth date (YYYY-MM-DD). Defaults to 1 week ago")
	createUserCmd.MarkFlagRequired("username")
}

func createUser(username, password, babyName string, birthDate time.Time) {
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

	// Hash password using Argon2
	fmt.Println("Hashing password...")
	hash, err := utils.HashPassword(password)
	if err != nil {
		log.Fatalf("Failed to hash password: %v", err)
	}

	// Create user
	user := models.User{
		Username:     username,
		PasswordHash: hash,
	}

	if err := db.Create(&user).Error; err != nil {
		log.Fatalf("Failed to create user: %v", err)
	}

	// Create baby profile
	baby := models.Baby{
		UserID:    user.ID,
		Name:      babyName,
		BirthDate: birthDate,
	}

	if err := db.Create(&baby).Error; err != nil {
		log.Fatalf("Failed to create baby profile: %v", err)
	}

	fmt.Printf("✅ User '%s' created successfully!\n", username)
	fmt.Printf("   Baby profile '%s' created.\n", babyName)
	fmt.Printf("   Birth date: %s (age: %d days)\n",
		birthDate.Format("2006-01-02"),
		int(time.Since(birthDate).Hours()/24))
	fmt.Printf("   Password hashed with Argon2id\n")
}
