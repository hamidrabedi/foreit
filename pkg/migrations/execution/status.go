package execution

import (
	"context"
	"fmt"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/golang-migrate/migrate/v4"
)

// DetailedStatus provides detailed migration status information
type DetailedStatus struct {
	Current    string
	Next       string
	Applied    []MigrationInfo
	Pending    []MigrationInfo
	OutOfOrder []MigrationInfo
	Status     string // "OK", "PENDING", "DIRTY"
	Error      string
}

// MigrationInfo contains information about a single migration
type MigrationInfo struct {
	Version     string
	Name        string
	Applied     bool
	Checksum    string
	AppliedAt   string
}

// StatusReporter collects detailed migration status
type StatusReporter struct {
	migrationsPath string
	migrate        *migrate.Migrate
}

// NewStatusReporter creates a new status reporter
func NewStatusReporter(migrationsPath string, m *migrate.Migrate) *StatusReporter {
	return &StatusReporter{
		migrationsPath: migrationsPath,
		migrate:        m,
	}
}

// GetDetailedStatus returns detailed migration status
func (r *StatusReporter) GetDetailedStatus(ctx context.Context) (*DetailedStatus, error) {
	status := &DetailedStatus{
		Applied:    []MigrationInfo{},
		Pending:    []MigrationInfo{},
		OutOfOrder: []MigrationInfo{},
	}

	// Get current version
	version, dirty, err := r.migrate.Version()
	if err != nil && err != migrate.ErrNilVersion {
		return nil, fmt.Errorf("failed to get migration version: %w", err)
	}

	if err == migrate.ErrNilVersion {
		status.Current = "No migration applied yet"
		status.Status = "PENDING"
	} else {
		status.Current = fmt.Sprintf("%d", version)
		if dirty {
			status.Status = "DIRTY"
			status.Error = "Migration is in a dirty state - manual intervention required"
		} else {
			status.Status = "OK"
		}
	}

	// Get all migration files
	allMigrations, err := r.getAllMigrations()
	if err != nil {
		return nil, fmt.Errorf("failed to get migration files: %w", err)
	}

	// Get applied migrations from database
	appliedVersions, err := r.getAppliedVersions(ctx)
	if err != nil {
		// If we can't get applied versions, assume none are applied
		appliedVersions = make(map[uint]bool)
	}

	// Categorize migrations
	for _, mig := range allMigrations {
		migVersion, _ := parseVersion(mig.Version)
		isApplied := appliedVersions[migVersion]

		info := MigrationInfo{
			Version: mig.Version,
			Name:    mig.Name,
			Applied: isApplied,
		}

		if isApplied {
			status.Applied = append(status.Applied, info)
		} else {
			// Check if this migration is out of order
			if len(status.Applied) > 0 && migVersion < version {
				status.OutOfOrder = append(status.OutOfOrder, info)
			} else {
				status.Pending = append(status.Pending, info)
			}
		}
	}

	// Sort applied migrations
	sort.Slice(status.Applied, func(i, j int) bool {
		vi, _ := parseVersion(status.Applied[i].Version)
		vj, _ := parseVersion(status.Applied[j].Version)
		return vi < vj
	})

	// Sort pending migrations
	sort.Slice(status.Pending, func(i, j int) bool {
		vi, _ := parseVersion(status.Pending[i].Version)
		vj, _ := parseVersion(status.Pending[j].Version)
		return vi < vj
	})

	// Set next migration
	if len(status.Pending) > 0 {
		status.Next = status.Pending[0].Version
	} else {
		status.Next = "Already at latest version"
	}

	return status, nil
}

// MigrationFile represents a migration file
type MigrationFile struct {
	Version string
	Name    string
	Path    string
}

// getAllMigrations gets all migration files
func (r *StatusReporter) getAllMigrations() ([]MigrationFile, error) {
	pattern := filepath.Join(r.migrationsPath, "*_*.up.sql")
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return nil, err
	}

	var migrations []MigrationFile
	for _, match := range matches {
		basename := filepath.Base(match)
		// Extract version and name from filename like "000001_create_users.up.sql"
		parts := strings.Split(basename, "_")
		if len(parts) < 2 {
			continue
		}
		version := parts[0]
		name := strings.TrimSuffix(strings.Join(parts[1:], "_"), ".up.sql")
		
		migrations = append(migrations, MigrationFile{
			Version: version,
			Name:    name,
			Path:    match,
		})
	}

	// Sort by version
	sort.Slice(migrations, func(i, j int) bool {
		vi, _ := parseVersion(migrations[i].Version)
		vj, _ := parseVersion(migrations[j].Version)
		return vi < vj
	})

	return migrations, nil
}

// getAppliedVersions gets applied migration versions from database
func (r *StatusReporter) getAppliedVersions(ctx context.Context) (map[uint]bool, error) {
	// This would query the schema_migrations table
	// For now, return empty map as golang-migrate handles this internally
	// A full implementation would query: SELECT version FROM schema_migrations
	applied := make(map[uint]bool)
	
	// Try to get version to see if any migrations are applied
	version, _, err := r.migrate.Version()
	if err == nil {
		// If we got a version, we know at least that one is applied
		// In a full implementation, we'd query all applied versions
		applied[version] = true
	}
	
	return applied, nil
}

// parseVersion parses a version string to uint
func parseVersion(versionStr string) (uint, error) {
	version, err := strconv.ParseUint(versionStr, 10, 64)
	if err != nil {
		return 0, err
	}
	return uint(version), nil
}

