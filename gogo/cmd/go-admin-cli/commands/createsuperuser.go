package commands

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

// NewCreateSuperuserCommand creates a new createsuperuser command
func NewCreateSuperuserCommand() *cobra.Command {
	var (
		username string
		email    string
		password string
	)

	cmd := &cobra.Command{
		Use:   "createsuperuser",
		Short: "Create a superuser account",
		Long:  "Creates a superuser account for admin access",
		RunE: func(cmd *cobra.Command, args []string) error {
			return createSuperuser(username, email, password)
		},
	}

	cmd.Flags().StringVarP(&username, "username", "u", "", "Username for superuser")
	cmd.Flags().StringVarP(&email, "email", "e", "", "Email for superuser")
	cmd.Flags().StringVarP(&password, "password", "p", "", "Password for superuser")

	return cmd
}

func createSuperuser(username, email, password string) error {
	// Check if we're in a project directory
	wd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}

	// Check for go.mod
	if _, err := os.Stat(filepath.Join(wd, "go.mod")); os.IsNotExist(err) {
		return fmt.Errorf("not in a Go project directory (go.mod not found)")
	}

	// Prompt for missing fields
	if username == "" {
		fmt.Print("Username: ")
		fmt.Scanln(&username)
	}

	if email == "" {
		fmt.Print("Email: ")
		fmt.Scanln(&email)
	}

	if password == "" {
		fmt.Print("Password: ")
		fmt.Scanln(&password)
	}

	// Validate inputs
	if username == "" || email == "" || password == "" {
		return fmt.Errorf("username, email, and password are required")
	}

	fmt.Println("Creating superuser...")
	fmt.Printf("  Username: %s\n", username)
	fmt.Printf("  Email: %s\n", email)
	fmt.Println("  Password: [hidden]")

	// TODO: Implement actual superuser creation
	// - Connect to database
	// - Create user with admin privileges
	// - Hash password
	// - Save to database

	fmt.Println("✅ Superuser created successfully!")
	fmt.Println("Note: This is a placeholder. In production, this would create the user in the database.")

	return nil
}

