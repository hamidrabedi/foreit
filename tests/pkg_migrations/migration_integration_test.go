package migrations

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/forgego/forge/tests/testhelpers"
)

// TestMigrationApplySQLite tests migrations against SQLite
func TestMigrationApplySQLite(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	db, err := testhelpers.StartSQLiteMemory("file::memory:?cache=shared")
	require.NoError(t, err)
	defer db.Close()

	err = testhelpers.WaitForDBReady(ctx, db, 5*time.Second)
	require.NoError(t, err)

	// Test: Create a basic table via migration
	createUserSQL := `
	CREATE TABLE users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username VARCHAR(150) NOT NULL,
		email VARCHAR(254) NOT NULL UNIQUE,
		is_active BOOLEAN DEFAULT 1,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	)
	`
	testhelpers.RunSQLExpectSuccess(ctx, t, db, createUserSQL)

	// Verify table exists
	testhelpers.AssertTableExists(ctx, t, db, "sqlite", "users")

	// Verify columns
	testhelpers.AssertColumnExists(ctx, t, db, "sqlite", "users", "id")
	testhelpers.AssertColumnExists(ctx, t, db, "sqlite", "users", "username")
	testhelpers.AssertColumnExists(ctx, t, db, "sqlite", "users", "email")
}

// TestMigrationApplyPostgres tests migrations against Postgres
func TestMigrationApplyPostgres(t *testing.T) {
	if os.Getenv("DATABASE_URL") == "" && os.Getenv("RUN_POSTGRES_TESTS") == "" {
		t.Skip("Postgres not available, skipping test")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	opts := testhelpers.DefaultPostgresOpts()
	db, dsn, cleanup, err := testhelpers.StartPostgresContainer(ctx, opts)
	require.NoError(t, err)
	defer cleanup()
	defer db.Close()

	t.Logf("Connected to Postgres: %s", dsn)

	// Test: Create users and posts tables with FK
	createUserSQL := `
	CREATE TABLE users (
		id BIGSERIAL PRIMARY KEY,
		username VARCHAR(150) NOT NULL UNIQUE,
		email VARCHAR(254) NOT NULL UNIQUE,
		is_active BOOLEAN DEFAULT true,
		created_at TIMESTAMP DEFAULT NOW()
	)
	`
	testhelpers.RunSQLExpectSuccess(ctx, t, db, createUserSQL)

	createPostSQL := `
	CREATE TABLE posts (
		id BIGSERIAL PRIMARY KEY,
		title VARCHAR(200) NOT NULL,
		content TEXT,
		user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
		created_at TIMESTAMP DEFAULT NOW()
	)
	`
	testhelpers.RunSQLExpectSuccess(ctx, t, db, createPostSQL)

	// Verify tables and FK
	testhelpers.AssertTableExists(ctx, t, db, "postgres", "users")
	testhelpers.AssertTableExists(ctx, t, db, "postgres", "posts")
	testhelpers.AssertForeignKeyExists(ctx, t, db, "postgres", "posts", "user_id")

	// Test: FK constraint on delete
	insertUserSQL := `INSERT INTO users (username, email) VALUES ('testuser', 'test@example.com')`
	testhelpers.RunSQLExpectSuccess(ctx, t, db, insertUserSQL)

	insertPostSQL := `INSERT INTO posts (title, user_id) VALUES ('Test Post', 1)`
	testhelpers.RunSQLExpectSuccess(ctx, t, db, insertPostSQL)

	// Verify rows exist
	testhelpers.AssertRowCount(ctx, t, db, "users", 1)
	testhelpers.AssertRowCount(ctx, t, db, "posts", 1)

	// Delete user and verify cascade
	deleteUserSQL := `DELETE FROM users WHERE id = 1`
	testhelpers.RunSQLExpectSuccess(ctx, t, db, deleteUserSQL)

	// Posts should be deleted due to CASCADE
	testhelpers.AssertRowCount(ctx, t, db, "posts", 0)
}

// TestIndexCreation tests that indexes are properly created
func TestIndexCreation(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	db, err := testhelpers.StartSQLiteMemory("")
	require.NoError(t, err)
	defer db.Close()


	// Create table (simpler format for SQLite)
	createTableSQL := `CREATE TABLE users (id INTEGER PRIMARY KEY, email VARCHAR(254) NOT NULL UNIQUE)`
	testhelpers.RunSQLExpectSuccess(ctx, t, db, createTableSQL)
	
	// Create index separately
	createIndexSQL := `CREATE INDEX idx_email ON users(email)`
	testhelpers.RunSQLExpectSuccess(ctx, t, db, createIndexSQL)
	testhelpers.AssertTableExists(ctx, t, db, "sqlite", "users")
	testhelpers.AssertIndexExists(ctx, t, db, "sqlite", "users", "idx_email")
}

// TestUniqueConstraintEnforcement tests that UNIQUE constraints are enforced
func TestUniqueConstraintEnforcement(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	db, err := testhelpers.StartSQLiteMemory("")
	require.NoError(t, err)
	defer db.Close()

	createSQL := `
	CREATE TABLE users (
		id INTEGER PRIMARY KEY,
		email VARCHAR(254) NOT NULL UNIQUE
	)
	`
	testhelpers.RunSQLExpectSuccess(ctx, t, db, createSQL)

	// Insert first user
	insert1 := `INSERT INTO users (email) VALUES ('test@example.com')`
	testhelpers.RunSQLExpectSuccess(ctx, t, db, insert1)

	// Insert duplicate should fail
	insert2 := `INSERT INTO users (email) VALUES ('test@example.com')`
	err = testhelpers.RunSQLExpectError(ctx, db, insert2)
	assert.Error(t, err, "duplicate email should fail")
}
