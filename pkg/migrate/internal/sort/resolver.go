package sort

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/forgego/forge/pkg/migrate/types"
)

// Resolver resolves migration dependencies and orders migrations correctly
type Resolver struct {
	migrationsDir string
}

// NewResolver creates a new dependency resolver
func NewResolver(migrationsDir string) *Resolver {
	return &Resolver{
		migrationsDir: migrationsDir,
	}
}

// MigrationWithDeps represents a migration with its dependencies
type MigrationWithDeps struct {
	Version      string
	Name         string
	Path         string
	Dependencies []types.Dependency
}

// ResolveDependencies resolves migration dependencies and returns migrations in correct order
func (r *Resolver) ResolveDependencies() ([]MigrationWithDeps, error) {
	// Find all migration files
	pattern := filepath.Join(r.migrationsDir, "*_*.up.sql")
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return nil, fmt.Errorf("failed to scan migrations directory: %w", err)
	}

	// Parse migrations and their dependencies
	migrations := make(map[string]*MigrationWithDeps)
	for _, match := range matches {
		mig, err := r.parseMigration(match)
		if err != nil {
			continue // Skip files that can't be parsed
		}
		migrations[mig.Version] = mig
	}

	// Resolve order using topological sort
	ordered, err := r.topologicalSort(migrations)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve dependencies: %w", err)
	}

	return ordered, nil
}

// parseMigration parses a migration file and extracts dependencies from comments
func (r *Resolver) parseMigration(filePath string) (*MigrationWithDeps, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	basename := filepath.Base(filePath)
	// Extract version and name from filename like "000001_create_users.up.sql"
	parts := strings.Split(basename, "_")
	if len(parts) < 2 {
		return nil, fmt.Errorf("invalid migration filename format")
	}
	version := parts[0]
	name := strings.TrimSuffix(strings.Join(parts[1:], "_"), ".up.sql")

	mig := &MigrationWithDeps{
		Version:      version,
		Name:         name,
		Path:         filePath,
		Dependencies: []types.Dependency{},
	}

	// Parse dependencies from comments
	// Format: -- DEPENDS: app:version or -- DEPENDS: version
	dependencyPattern := regexp.MustCompile(`(?i)--\s*DEPENDS:\s*([^\n]+)`)
	matches := dependencyPattern.FindAllStringSubmatch(string(content), -1)
	for _, match := range matches {
		if len(match) < 2 {
			continue
		}
		depStr := strings.TrimSpace(match[1])
		dep := r.parseDependency(depStr)
		if dep != nil {
			mig.Dependencies = append(mig.Dependencies, *dep)
		}
	}

	return mig, nil
}

// parseDependency parses a dependency string into a Dependency struct
func (r *Resolver) parseDependency(depStr string) *types.Dependency {
	// Format: "app:version" or just "version"
	parts := strings.Split(depStr, ":")
	if len(parts) == 2 {
		return &types.Dependency{
			App:     strings.TrimSpace(parts[0]),
			Version: strings.TrimSpace(parts[1]),
		}
	} else if len(parts) == 1 {
		return &types.Dependency{
			Version: strings.TrimSpace(parts[0]),
		}
	}
	return nil
}

// topologicalSort performs topological sort on migrations based on dependencies
func (r *Resolver) topologicalSort(migrations map[string]*MigrationWithDeps) ([]MigrationWithDeps, error) {
	// Build dependency graph
	inDegree := make(map[string]int)
	for version := range migrations {
		inDegree[version] = 0
	}

	// Calculate in-degrees
	for _, mig := range migrations {
		for _, dep := range mig.Dependencies {
			// For same-app dependencies, just use version
			if dep.App == "" {
				if _, exists := migrations[dep.Version]; exists {
					inDegree[mig.Version]++
				}
			}
		}
	}

	// Find migrations with no dependencies
	queue := []string{}
	for version, degree := range inDegree {
		if degree == 0 {
			queue = append(queue, version)
		}
	}

	var result []MigrationWithDeps
	visited := make(map[string]bool)

	// Process queue
	for len(queue) > 0 {
		// Sort queue by version for deterministic ordering
		for i := 0; i < len(queue)-1; i++ {
			for j := i + 1; j < len(queue); j++ {
				vi, _ := strconv.ParseUint(queue[i], 10, 64)
				vj, _ := strconv.ParseUint(queue[j], 10, 64)
				if vi > vj {
					queue[i], queue[j] = queue[j], queue[i]
				}
			}
		}

		current := queue[0]
		queue = queue[1:]

		if visited[current] {
			continue
		}
		visited[current] = true

		if mig, exists := migrations[current]; exists {
			result = append(result, *mig)
		}

		// Update in-degrees of dependent migrations
		for _, mig := range migrations {
			for _, dep := range mig.Dependencies {
				if dep.Version == current {
					inDegree[mig.Version]--
					if inDegree[mig.Version] == 0 && !visited[mig.Version] {
						queue = append(queue, mig.Version)
					}
				}
			}
		}
	}

	// Check for cycles (unvisited migrations)
	if len(result) < len(migrations) {
		var unvisited []string
		for version := range migrations {
			if !visited[version] {
				unvisited = append(unvisited, version)
			}
		}
		return nil, fmt.Errorf("circular dependency detected in migrations: %v", unvisited)
	}

	return result, nil
}

// FormatDependencyComment formats a dependency as a comment for inclusion in migration files
func FormatDependencyComment(dep types.Dependency) string {
	if dep.App != "" {
		return fmt.Sprintf("-- DEPENDS: %s:%s", dep.App, dep.Version)
	}
	return fmt.Sprintf("-- DEPENDS: %s", dep.Version)
}
