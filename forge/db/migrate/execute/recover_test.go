package execute

import (
	"context"
	"database/sql"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	_ "github.com/mattn/go-sqlite3"
)

func setupTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("sqlite3", ":memory:")
	require.NoError(t, err)

	// Create schema_migrations table
	_, err = db.Exec(`
		CREATE TABLE schema_migrations (
			version INTEGER PRIMARY KEY,
			dirty BOOLEAN NOT NULL DEFAULT 0
		)
	`)
	require.NoError(t, err)

	return db
}

func TestRecovery_GetDirtyMigrationInfo(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	recovery := NewRecovery(db)
	ctx := context.Background()

	t.Run("No migrations", func(t *testing.T) {
		info, err := recovery.GetDirtyMigrationInfo(ctx)
		require.NoError(t, err)
		assert.Equal(t, uint(0), info.Version)
		assert.False(t, info.Dirty)
	})

	t.Run("Clean migration", func(t *testing.T) {
		_, err := db.Exec("INSERT INTO schema_migrations (version, dirty) VALUES (1, 0)")
		require.NoError(t, err)

		info, err := recovery.GetDirtyMigrationInfo(ctx)
		require.NoError(t, err)
		assert.Equal(t, uint(1), info.Version)
		assert.False(t, info.Dirty)
	})

	t.Run("Dirty migration", func(t *testing.T) {
		_, err := db.Exec("UPDATE schema_migrations SET dirty = 1 WHERE version = 1")
		require.NoError(t, err)

		info, err := recovery.GetDirtyMigrationInfo(ctx)
		require.NoError(t, err)
		assert.Equal(t, uint(1), info.Version)
		assert.True(t, info.Dirty)
	})
}

func TestRecovery_RecoverDirtyState(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	recovery := NewRecovery(db)
	ctx := context.Background()

	t.Run("No dirty migrations", func(t *testing.T) {
		_, err := db.Exec("INSERT INTO schema_migrations (version, dirty) VALUES (1, 0)")
		require.NoError(t, err)

		dirtyMig, err := recovery.RecoverDirtyState(ctx, "")
		require.NoError(t, err)
		assert.Nil(t, dirtyMig)
	})

	t.Run("Dirty migration found", func(t *testing.T) {
		_, err := db.Exec("INSERT INTO schema_migrations (version, dirty) VALUES (2, 1)")
		require.NoError(t, err)

		dirtyMig, err := recovery.RecoverDirtyState(ctx, "")
		require.NoError(t, err)
		require.NotNil(t, dirtyMig)
		assert.Equal(t, uint(2), dirtyMig.Version)
	})
}

func TestRecovery_MarkMigrationClean(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	recovery := NewRecovery(db)
	ctx := context.Background()

	// Insert dirty migration
	_, err := db.Exec("INSERT INTO schema_migrations (version, dirty) VALUES (1, 1)")
	require.NoError(t, err)

	// Mark as clean
	err = recovery.MarkMigrationClean(ctx, 1)
	require.NoError(t, err)

	// Verify
	var dirty bool
	err = db.QueryRow("SELECT dirty FROM schema_migrations WHERE version = 1").Scan(&dirty)
	require.NoError(t, err)
	assert.False(t, dirty)
}

func TestRecovery_GetAppliedMigrations(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	recovery := NewRecovery(db)
	ctx := context.Background()

	// Insert test migrations
	migrations := []struct {
		version uint
		dirty   bool
	}{
		{1, false},
		{2, false},
		{3, true},
		{4, false},
	}

	for _, m := range migrations {
		_, err := db.Exec("INSERT INTO schema_migrations (version, dirty) VALUES (?, ?)", m.version, m.dirty)
		require.NoError(t, err)
	}

	// Get all migrations
	applied, err := recovery.GetAppliedMigrations(ctx)
	require.NoError(t, err)
	assert.Len(t, applied, 4)

	// Verify order (should be ascending)
	for i := 0; i < len(applied)-1; i++ {
		assert.Less(t, applied[i].Version, applied[i+1].Version)
	}

	// Verify dirty flag
	assert.False(t, applied[0].Dirty)
	assert.True(t, applied[2].Dirty) // Version 3
}


func TestRecovery_ForceCleanState(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	recovery := NewRecovery(db)
	ctx := context.Background()

	// Insert dirty migration
	_, err := db.Exec("INSERT INTO schema_migrations (version, dirty) VALUES (1, 1)")
	require.NoError(t, err)

	// Force clean
	err = recovery.ForceCleanState(ctx, 1)
	require.NoError(t, err)

	// Verify
	var dirty bool
	err = db.QueryRow("SELECT dirty FROM schema_migrations WHERE version = 1").Scan(&dirty)
	require.NoError(t, err)
	assert.False(t, dirty)
}

func TestRecovery_ValidateMigrationIntegrity(t *testing.T) {
	recovery := NewRecovery(nil) // No DB needed for file operations

	// Create temp directory with test migrations
	tmpDir := t.TempDir()

	testMigrations := []struct {
		filename string
		content  string
	}{
		{"20240101000001_create_users.up.sql", "CREATE TABLE users (id INT);"},
		{"20240101000002_add_index.up.sql", "CREATE INDEX idx_users ON users(id);"},
	}

	for _, tm := range testMigrations {
		path := filepath.Join(tmpDir, tm.filename)
		err := os.WriteFile(path, []byte(tm.content), 0644)
		require.NoError(t, err)
	}

	// Validate integrity
	checksums, err := recovery.ValidateMigrationIntegrity(tmpDir)
	require.NoError(t, err)
	assert.Len(t, checksums, 2)

	// Verify checksums are consistent
	checksums2, err := recovery.ValidateMigrationIntegrity(tmpDir)
	require.NoError(t, err)
	assert.Equal(t, checksums, checksums2)
}

func TestRecovery_CompareChecksums(t *testing.T) {
	recovery := NewRecovery(nil)

	// Create temp directory with test migrations
	tmpDir := t.TempDir()
	migrationFile := filepath.Join(tmpDir, "20240101000001_test.up.sql")
	
	// Initial content
	err := os.WriteFile(migrationFile, []byte("CREATE TABLE test;"), 0644)
	require.NoError(t, err)

	// Get initial checksums
	initialChecksums, err := recovery.ValidateMigrationIntegrity(tmpDir)
	require.NoError(t, err)

	t.Run("No modifications", func(t *testing.T) {
		modified, err := recovery.CompareChecksums(tmpDir, initialChecksums)
		require.NoError(t, err)
		assert.Empty(t, modified)
	})

	t.Run("Modified migration", func(t *testing.T) {
		// Modify the file
		err := os.WriteFile(migrationFile, []byte("CREATE TABLE test_modified;"), 0644)
		require.NoError(t, err)

		modified, err := recovery.CompareChecksums(tmpDir, initialChecksums)
		require.NoError(t, err)
		assert.Contains(t, modified, uint(20240101000001))
	})
}

func TestRecovery_RollbackPartialMigration(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	recovery := NewRecovery(db)
	ctx := context.Background()

	// Create a test table
	_, err := db.Exec("CREATE TABLE test_table (id INTEGER PRIMARY KEY)")
	require.NoError(t, err)

	// Insert migration record
	_, err = db.Exec("INSERT INTO schema_migrations (version, dirty) VALUES (1, 1)")
	require.NoError(t, err)

	// Rollback with down SQL
	downSQL := "DROP TABLE test_table;"
	err = recovery.RollbackPartialMigration(ctx, 1, downSQL)
	require.NoError(t, err)

	// Verify table is dropped
	var tableName string
	err = db.QueryRow("SELECT name FROM sqlite_master WHERE type='table' AND name='test_table'").Scan(&tableName)
	assert.Error(t, err) // Should be no rows

	// Verify migration record is removed
	var version uint
	err = db.QueryRow("SELECT version FROM schema_migrations WHERE version = 1").Scan(&version)
	assert.Error(t, err) // Should be no rows
}

func TestRecovery_SplitSQL(t *testing.T) {
	tests := []struct {
		name     string
		sql      string
		expected int
	}{
		{
			name:     "Single statement",
			sql:      "CREATE TABLE users (id INT);",
			expected: 1,
		},
		{
			name:     "Multiple statements",
			sql:      "CREATE TABLE users (id INT); CREATE INDEX idx ON users(id);",
			expected: 2,
		},
		{
			name:     "Statements with whitespace",
			sql:      "  CREATE TABLE t1;  \n\n  CREATE TABLE t2;  ",
			expected: 2,
		},
		{
			name:     "Empty statements",
			sql:      ";;;",
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			statements := splitSQL(tt.sql)
			assert.Len(t, statements, tt.expected)
		})
	}
}

func TestRecovery_GetRecoverySteps(t *testing.T) {
	recovery := NewRecovery(nil)

	t.Run("Generic error", func(t *testing.T) {
		steps := recovery.GetRecoverySteps("001", "some error")
		assert.NotEmpty(t, steps)
		assert.Contains(t, steps[0], "error message")
	})

	t.Run("Constraint error", func(t *testing.T) {
		steps := recovery.GetRecoverySteps("001", "constraint violation")
		found := false
		for _, step := range steps {
			if containsIgnoreCase(step, "foreign key") || containsIgnoreCase(step, "constraint") {
				found = true
				break
			}
		}
		assert.True(t, found, "Should mention constraint handling")
	})

	t.Run("Column error", func(t *testing.T) {
		steps := recovery.GetRecoverySteps("001", "column already exists")
		found := false
		for _, step := range steps {
			if containsIgnoreCase(step, "column") {
				found = true
				break
			}
		}
		assert.True(t, found, "Should mention column handling")
	})
}

func containsIgnoreCase(s, substr string) bool {
	return contains(toLower(s), toLower(substr))
}

func toLower(s string) string {
	// Simple toLower implementation
	result := make([]byte, len(s))
	for i := 0; i < len(s); i++ {
		if s[i] >= 'A' && s[i] <= 'Z' {
			result[i] = s[i] + 32
		} else {
			result[i] = s[i]
		}
	}
	return string(result)
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && findSubstr(s, substr) >= 0
}

func findSubstr(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}

