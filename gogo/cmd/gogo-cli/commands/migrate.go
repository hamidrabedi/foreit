package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/gogo/pkg/orm"
	"github.com/spf13/cobra"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// RegisterMigrateCommand registers the migrate command
func RegisterMigrateCommand(rootCmd *cobra.Command) {
	migrateCmd := &cobra.Command{
		Use:   "migrate",
		Short: "Run database migrations",
		Long:  "Manage versioned database migrations using golang-migrate",
	}

	migrateUpCmd := &cobra.Command{
		Use:   "up",
		Short: "Apply all pending migrations",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runMigrate("up", nil)
		},
	}

	migrateDownCmd := &cobra.Command{
		Use:   "down",
		Short: "Rollback the last migration",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runMigrate("down", nil)
		},
	}

	migrateStepsCmd := &cobra.Command{
		Use:   "steps [n]",
		Short: "Apply or rollback N migrations",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			n, err := strconv.Atoi(args[0])
			if err != nil {
				return fmt.Errorf("invalid number of steps: %w", err)
			}
			return runMigrate("steps", &n)
		},
	}

	migrateVersionCmd := &cobra.Command{
		Use:   "version",
		Short: "Show current migration version",
		RunE: func(cmd *cobra.Command, args []string) error {
			return showVersion()
		},
	}

	migrateCmd.AddCommand(migrateUpCmd)
	migrateCmd.AddCommand(migrateDownCmd)
	migrateCmd.AddCommand(migrateStepsCmd)
	migrateCmd.AddCommand(migrateVersionCmd)

	rootCmd.AddCommand(migrateCmd)
}

func runMigrate(action string, steps *int) error {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		return fmt.Errorf("DATABASE_URL environment variable not set")
	}

	driver := os.Getenv("DATABASE_DRIVER")
	if driver == "" {
		driver = "postgres"
	}

	var db *gorm.DB
	var err error

	switch driver {
	case "postgres":
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	case "sqlite", "sqlite3":
		db, err = gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	default:
		return fmt.Errorf("unsupported driver: %s", driver)
	}

	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	migrationsPath := os.Getenv("MIGRATIONS_PATH")
	if migrationsPath == "" {
		cwd, _ := os.Getwd()
		migrationsPath = filepath.Join(cwd, "migrations")
	}

	client := &orm.Client{db: db, driver: driver, dsn: dsn}
	migrator, err := client.Migrator(migrationsPath)
	if err != nil {
		return fmt.Errorf("failed to create migrator: %w", err)
	}
	defer migrator.Close()

	switch action {
	case "up":
		if err := migrator.Up(); err != nil {
			return err
		}
		fmt.Println("Migrations applied successfully")
	case "down":
		if err := migrator.Down(); err != nil {
			return err
		}
		fmt.Println("Migration rolled back successfully")
	case "steps":
		if err := migrator.Steps(*steps); err != nil {
			return err
		}
		fmt.Printf("Applied %d migration steps\n", *steps)
	}

	return nil
}

func showVersion() error {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		return fmt.Errorf("DATABASE_URL environment variable not set")
	}

	driver := os.Getenv("DATABASE_DRIVER")
	if driver == "" {
		driver = "postgres"
	}

	var db *gorm.DB
	var err error

	switch driver {
	case "postgres":
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	case "sqlite", "sqlite3":
		db, err = gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	default:
		return fmt.Errorf("unsupported driver: %s", driver)
	}

	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	migrationsPath := os.Getenv("MIGRATIONS_PATH")
	if migrationsPath == "" {
		cwd, _ := os.Getwd()
		migrationsPath = filepath.Join(cwd, "migrations")
	}

	client := &orm.Client{db: db, driver: driver, dsn: dsn}
	migrator, err := client.Migrator(migrationsPath)
	if err != nil {
		return fmt.Errorf("failed to create migrator: %w", err)
	}
	defer migrator.Close()

	version, dirty, err := migrator.Version()
	if err != nil {
		return fmt.Errorf("failed to get version: %w", err)
	}

	if dirty {
		fmt.Printf("Current version: %d (dirty)\n", version)
	} else {
		fmt.Printf("Current version: %d\n", version)
	}

	return nil
}
