package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "bambino",
	Short: "Baby Tracker - Self-hosted baby activity tracking",
	Long: `Baby Tracker is a privacy-focused, self-hosted application 
for tracking baby activities like feeding, sleeping, and diaper changes.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func init() {
	// Add global flags here if needed
}
