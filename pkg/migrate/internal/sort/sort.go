package sort

import (
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/forgego/forge/pkg/migrate/types"
)

// SortMigrationFiles sorts migration files by version number
func SortMigrationFiles(files []string) []string {
	// Parse version from filename and sort numerically
	// Filename format: VERSION_NAME.up.sql (e.g., 000001_create_users.up.sql)
	sort.Slice(files, func(i, j int) bool {
		vi := extractVersionFromFilename(files[i])
		vj := extractVersionFromFilename(files[j])
		return vi < vj
	})
	return files
}

// extractVersionFromFilename extracts the version number from a migration filename
// Returns 0 if version cannot be parsed
func extractVersionFromFilename(filename string) uint64 {
	basename := filepath.Base(filename)
	// Extract version prefix (e.g., "000001" from "000001_create_users.up.sql")
	parts := strings.Split(basename, "_")
	if len(parts) == 0 {
		return 0
	}

	versionStr := parts[0]
	version, err := strconv.ParseUint(versionStr, 10, 64)
	if err != nil {
		return 0
	}
	return version
}

// SortChangesByType sorts changes to ensure CREATE TABLE comes before other changes
func SortChangesByType(changes []types.Change) []types.Change {
	// Create separate slices for different change types
	var createTables []types.Change
	var addColumns []types.Change
	var addForeignKeys []types.Change
	var otherChanges []types.Change

	for _, change := range changes {
		switch change.Type() {
		case types.ChangeTypeCreateTable:
			createTables = append(createTables, change)
		case types.ChangeTypeAddColumn:
			addColumns = append(addColumns, change)
		case types.ChangeTypeAddForeignKey:
			addForeignKeys = append(addForeignKeys, change)
		default:
			otherChanges = append(otherChanges, change)
		}
	}

	// Combine in order: CREATE TABLE, ADD COLUMN, ADD FOREIGN KEY, others
	result := make([]types.Change, 0, len(changes))
	result = append(result, createTables...)
	result = append(result, addColumns...)
	result = append(result, addForeignKeys...)
	result = append(result, otherChanges...)

	return result
}
