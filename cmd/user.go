package cmd

import (
	"fmt"
	"log"
	"syscall"

	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/term"

	"github.com/engineervix/baby-tracker/internal/config"
	"github.com/engineervix/baby-tracker/internal/database"
	"github.com/engineervix/baby-tracker/internal/models"
)

var createUserCmd = &cobra.Command{
	Use:   "create-user",
	Short: "Create a new user",
	Long:  `Creates a new user account with the specified username.`,
	Run: func(cmd *cobra.Command, args []string) {
		username, _ := cmd.Flags().GetString("username")
		babyName, _ := cmd.Flags().GetString("baby-name")

		if username == "" {
			log.Fatal("Username is required")
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

		createUser(username, string(password), babyName)
	},
}

func init() {
	rootCmd.AddCommand(createUserCmd)

	createUserCmd.Flags().StringP("username", "u", "", "Username for the new user (required)")
	createUserCmd.Flags().StringP("baby-name", "b", "Baby", "Name of the baby")
	createUserCmd.MarkFlagRequired("username")
}

func createUser(username, password, babyName string) {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// Load configuration
	cfg := config.Load()

	// Connect to database
	db, err := database.Connect(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Hash password
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal(err)
	}

	// Create user
	user := models.User{
		Username:     username,
		PasswordHash: string(hash),
	}

	if err := db.Create(&user).Error; err != nil {
		log.Fatalf("Failed to create user: %v", err)
	}

	// Create baby profile
	baby := models.Baby{
		UserID: user.ID,
		Name:   babyName,
	}

	if err := db.Create(&baby).Error; err != nil {
		log.Fatalf("Failed to create baby profile: %v", err)
	}

	fmt.Printf("✅ User '%s' created successfully!\n", username)
	fmt.Printf("   Baby profile '%s' created.\n", babyName)
}
