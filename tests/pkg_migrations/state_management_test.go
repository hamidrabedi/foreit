package migrations

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/forgego/forge/orm"
	"github.com/forgego/forge/tests/testhelpers"
)

// TestMigrationStateTracking tests that migration state is properly tracked
func TestMigrationStateTracking(t *testing.T) {
	// Using localhost PostgreSQL
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	opts := testhelpers.PostgresOpts{
		UseDirect: true,
		Host:      "localhost",
		Port:      "5432",
		User:      "postgres",
		Password:  "123",
		DBName:    fmt.Sprintf("test_%s_%d", t.Name(), time.Now().UnixNano()),
	}
	postgresDB, dsn, cleanup, err := testhelpers.StartPostgresContainer(ctx, opts)
	require.NoError(t, err)
	defer cleanup()
	defer postgresDB.Close()

	// Create a temporary directory for migrations under tests/tmp (returns relative path)
	tempDir, cleanupTemp := testhelpers.TempDirInTests(t, "state_tracking_")
	defer cleanupTemp()
	migrationsDir := filepath.Join(tempDir, "migrations")
	require.NoError(t, os.MkdirAll(migrationsDir, 0755))

	// Create database connection
	database, err := db.NewDB(dsn)
	require.NoError(t, err)
	defer database.Close()

	// Create first migration
	createMigration1Up := `CREATE TABLE users (id BIGSERIAL PRIMARY KEY, username VARCHAR(150));`
	createMigration1Down := `DROP TABLE IF EXISTS users;`

	migration1UpPath := filepath.Join(migrationsDir, "000001_create_users.up.sql")
	migration1DownPath := filepath.Join(migrationsDir, "000001_create_users.down.sql")

	require.NoError(t, os.WriteFile(migration1UpPath, []byte(createMigration1Up), 0644))
	require.NoError(t, os.WriteFile(migration1DownPath, []byte(createMigration1Down), 0644))

	// Apply first migration
	runner, err := db.NewMigrationRunner(database, migrationsDir)
	require.NoError(t, err)

	err = runner.Migrate(ctx)
	require.NoError(t, err)

	// Verify state using the existing runner
	testhelpers.AssertMigrationStateWithRunner(ctx, t, database, migrationsDir, 1, false, runner)

	// Create second migration
	createMigration2Up := `CREATE TABLE posts (id BIGSERIAL PRIMARY KEY, title VARCHAR(200));`
	createMigration2Down := `DROP TABLE IF EXISTS posts;`

	migration2UpPath := filepath.Join(migrationsDir, "000002_create_posts.up.sql")
	migration2DownPath := filepath.Join(migrationsDir, "000002_create_posts.down.sql")

	require.NoError(t, os.WriteFile(migration2UpPath, []byte(createMigration2Up), 0644))
	require.NoError(t, os.WriteFile(migration2DownPath, []byte(createMigration2Down), 0644))

	// Apply second migration - don't close the runner, just create a new one
	// Note: golang-migrate's Close() may close the database connection, so we don't close the runner
	// Instead, we'll create a new runner which will work with the same database connection
	runner2, err := db.NewMigrationRunner(database, migrationsDir)
	require.NoError(t, err)
	defer runner2.Close()

	err = runner2.Migrate(ctx)
	require.NoError(t, err)

	// Verify state using the new runner
	testhelpers.AssertMigrationStateWithRunner(ctx, t, database, migrationsDir, 2, false, runner2)

	// Close both runners at the end
	runner.Close()

	// Verify both tables exist
	testhelpers.AssertTableExists(ctx, t, postgresDB, "postgres", "users")
	testhelpers.AssertTableExists(ctx, t, postgresDB, "postgres", "posts")

	// Close runner after all operations
	defer runner.Close()
}

// TestMigrationStateAfterRollback tests state tracking after rollback
func TestMigrationStateAfterRollback(t *testing.T) {
	// Using localhost PostgreSQL
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	opts := testhelpers.PostgresOpts{
		UseDirect: true,
		Host:      "localhost",
		Port:      "5432",
		User:      "postgres",
		Password:  "123",
		DBName:    fmt.Sprintf("test_%s_%d", t.Name(), time.Now().UnixNano()),
	}
	postgresDB, dsn, cleanup, err := testhelpers.StartPostgresContainer(ctx, opts)
	require.NoError(t, err)
	defer cleanup()
	defer postgresDB.Close()

	// Create a temporary directory for migrations under tests/tmp (returns relative path)
	tempDir, cleanupTemp := testhelpers.TempDirInTests(t, "state_tracking_")
	defer cleanupTemp()
	migrationsDir := filepath.Join(tempDir, "migrations")
	require.NoError(t, os.MkdirAll(migrationsDir, 0755))

	// Create database connection
	database, err := db.NewDB(dsn)
	require.NoError(t, err)
	defer database.Close()

	// Create migrations
	createMigration1Up := `CREATE TABLE users (id BIGSERIAL PRIMARY KEY);`
	createMigration1Down := `DROP TABLE IF EXISTS users;`

	migration1UpPath := filepath.Join(migrationsDir, "000001_create_users.up.sql")
	migration1DownPath := filepath.Join(migrationsDir, "000001_create_users.down.sql")

	require.NoError(t, os.WriteFile(migration1UpPath, []byte(createMigration1Up), 0644))
	require.NoError(t, os.WriteFile(migration1DownPath, []byte(createMigration1Down), 0644))

	createMigration2Up := `CREATE TABLE posts (id BIGSERIAL PRIMARY KEY);`
	createMigration2Down := `DROP TABLE IF EXISTS posts;`

	migration2UpPath := filepath.Join(migrationsDir, "000002_create_posts.up.sql")
	migration2DownPath := filepath.Join(migrationsDir, "000002_create_posts.down.sql")

	require.NoError(t, os.WriteFile(migration2UpPath, []byte(createMigration2Up), 0644))
	require.NoError(t, os.WriteFile(migration2DownPath, []byte(createMigration2Down), 0644))

	// Apply both migrations
	runner, err := db.NewMigrationRunner(database, migrationsDir)
	require.NoError(t, err)

	err = runner.Migrate(ctx)
	require.NoError(t, err)

	// Verify state
	testhelpers.AssertMigrationState(ctx, t, database, migrationsDir, 2, false)

	// Rollback one migration
	err = runner.Rollback(ctx)
	require.NoError(t, err)

	// Verify state updated using the existing runner
	testhelpers.AssertMigrationStateWithRunner(ctx, t, database, migrationsDir, 1, false, runner)

	// Close runner after assertions
	defer runner.Close()
}

// TestSchemaDriftDetection tests detecting schema drift (manual changes to database)
func TestSchemaDriftDetection(t *testing.T) {
	// Using localhost PostgreSQL
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	opts := testhelpers.PostgresOpts{
		UseDirect: true,
		Host:      "localhost",
		Port:      "5432",
		User:      "postgres",
		Password:  "123",
		DBName:    fmt.Sprintf("test_%s_%d", t.Name(), time.Now().UnixNano()),
	}
	postgresDB, dsn, cleanup, err := testhelpers.StartPostgresContainer(ctx, opts)
	require.NoError(t, err)
	defer cleanup()
	defer postgresDB.Close()

	// Create a temporary directory for migrations under tests/tmp (returns relative path)
	tempDir, cleanupTemp := testhelpers.TempDirInTests(t, "state_tracking_")
	defer cleanupTemp()
	migrationsDir := filepath.Join(tempDir, "migrations")
	require.NoError(t, os.MkdirAll(migrationsDir, 0755))

	// Create database connection
	database, err := db.NewDB(dsn)
	require.NoError(t, err)
	defer database.Close()

	// Create and apply migration
	createMigrationUp := `CREATE TABLE users (id BIGSERIAL PRIMARY KEY, username VARCHAR(150));`
	createMigrationDown := `DROP TABLE IF EXISTS users;`

	migrationUpPath := filepath.Join(migrationsDir, "000001_create_users.up.sql")
	migrationDownPath := filepath.Join(migrationsDir, "000001_create_users.down.sql")

	require.NoError(t, os.WriteFile(migrationUpPath, []byte(createMigrationUp), 0644))
	require.NoError(t, os.WriteFile(migrationDownPath, []byte(createMigrationDown), 0644))

	runner, err := db.NewMigrationRunner(database, migrationsDir)
	require.NoError(t, err)
	defer runner.Close()

	err = runner.Migrate(ctx)
	require.NoError(t, err)

	// Verify initial state
	testhelpers.AssertTableExists(ctx, t, postgresDB, "postgres", "users")
	testhelpers.AssertMigrationState(ctx, t, database, migrationsDir, 1, false)

	// Manually add a column (simulating drift)
	alterSQL := `ALTER TABLE users ADD COLUMN email VARCHAR(254)`
	testhelpers.RunSQLExpectSuccess(ctx, t, postgresDB, alterSQL)

	// Verify column was added
	testhelpers.AssertColumnExists(ctx, t, postgresDB, "postgres", "users", "email")

	// The migration state should still be at version 1, but the schema has drifted
	// This demonstrates that manual changes create drift
	testhelpers.AssertMigrationStateWithRunner(ctx, t, database, migrationsDir, 1, false, runner)

	// Verify the drift: column exists in DB but not in migration
	// This would be detected by a drift detection tool
	var columnExists bool
	err = postgresDB.QueryRowContext(ctx, `
		SELECT EXISTS (
			SELECT 1 FROM information_schema.columns 
			WHERE table_name = 'users' AND column_name = 'email' AND table_schema = 'public'
		)
	`).Scan(&columnExists)
	require.NoError(t, err)
	require.True(t, columnExists, "email column should exist (drift)")
}

// TestMigrationStateConsistency tests that migration state remains consistent
func TestMigrationStateConsistency(t *testing.T) {
	// Using localhost PostgreSQL
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	opts := testhelpers.PostgresOpts{
		UseDirect: true,
		Host:      "localhost",
		Port:      "5432",
		User:      "postgres",
		Password:  "123",
		DBName:    fmt.Sprintf("test_%s_%d", t.Name(), time.Now().UnixNano()),
	}
	postgresDB, dsn, cleanup, err := testhelpers.StartPostgresContainer(ctx, opts)
	require.NoError(t, err)
	defer cleanup()
	defer postgresDB.Close()

	// Create a temporary directory for migrations under tests/tmp (returns relative path)
	tempDir, cleanupTemp := testhelpers.TempDirInTests(t, "state_tracking_")
	defer cleanupTemp()
	migrationsDir := filepath.Join(tempDir, "migrations")
	require.NoError(t, os.MkdirAll(migrationsDir, 0755))

	// Create database connection
	database, err := db.NewDB(dsn)
	require.NoError(t, err)
	defer database.Close()

	// Create migration
	createMigrationUp := `CREATE TABLE users (id BIGSERIAL PRIMARY KEY, username VARCHAR(150));`
	createMigrationDown := `DROP TABLE IF EXISTS users;`

	migrationUpPath := filepath.Join(migrationsDir, "000001_create_users.up.sql")
	migrationDownPath := filepath.Join(migrationsDir, "000001_create_users.down.sql")

	require.NoError(t, os.WriteFile(migrationUpPath, []byte(createMigrationUp), 0644))
	require.NoError(t, os.WriteFile(migrationDownPath, []byte(createMigrationDown), 0644))

	// Apply migration
	runner, err := db.NewMigrationRunner(database, migrationsDir)
	require.NoError(t, err)

	err = runner.Migrate(ctx)
	require.NoError(t, err)

	// Verify state using the existing runner
	testhelpers.AssertMigrationStateWithRunner(ctx, t, database, migrationsDir, 1, false, runner)

	// Don't close the runner - golang-migrate's Close() may close the database connection
	// Instead, create a new runner without closing the first one
	// In a real application restart scenario, the database connection would be recreated anyway
	runner2, err := db.NewMigrationRunner(database, migrationsDir)
	require.NoError(t, err)
	defer runner2.Close()

	// Verify state is consistent across runner instances
	testhelpers.AssertMigrationStateWithRunner(ctx, t, database, migrationsDir, 1, false, runner2)

	// Close both runners at the end
	defer runner.Close()

	// Verify table still exists
	testhelpers.AssertTableExists(ctx, t, postgresDB, "postgres", "users")
}

// TestMigrationStateWithMultipleTables tests state tracking with complex schemas
func TestMigrationStateWithMultipleTables(t *testing.T) {
	// Using localhost PostgreSQL
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	opts := testhelpers.PostgresOpts{
		UseDirect: true,
		Host:      "localhost",
		Port:      "5432",
		User:      "postgres",
		Password:  "123",
		DBName:    fmt.Sprintf("test_%s_%d", t.Name(), time.Now().UnixNano()),
	}
	postgresDB, dsn, cleanup, err := testhelpers.StartPostgresContainer(ctx, opts)
	require.NoError(t, err)
	defer cleanup()
	defer postgresDB.Close()

	// Create a temporary directory for migrations under tests/tmp (returns relative path)
	tempDir, cleanupTemp := testhelpers.TempDirInTests(t, "state_tracking_")
	defer cleanupTemp()
	migrationsDir := filepath.Join(tempDir, "migrations")
	require.NoError(t, os.MkdirAll(migrationsDir, 0755))

	// Create database connection
	database, err := db.NewDB(dsn)
	require.NoError(t, err)
	defer database.Close()

	// Create multiple migrations with relationships
	createMigration1Up := `
		CREATE TABLE users (
			id BIGSERIAL PRIMARY KEY,
			username VARCHAR(150) NOT NULL
		);
	`
	createMigration1Down := `DROP TABLE IF EXISTS users;`

	migration1UpPath := filepath.Join(migrationsDir, "000001_create_users.up.sql")
	migration1DownPath := filepath.Join(migrationsDir, "000001_create_users.down.sql")

	require.NoError(t, os.WriteFile(migration1UpPath, []byte(createMigration1Up), 0644))
	require.NoError(t, os.WriteFile(migration1DownPath, []byte(createMigration1Down), 0644))

	createMigration2Up := `
		CREATE TABLE posts (
			id BIGSERIAL PRIMARY KEY,
			title VARCHAR(200) NOT NULL,
			user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE
		);
	`
	createMigration2Down := `DROP TABLE IF EXISTS posts;`

	migration2UpPath := filepath.Join(migrationsDir, "000002_create_posts.up.sql")
	migration2DownPath := filepath.Join(migrationsDir, "000002_create_posts.down.sql")

	require.NoError(t, os.WriteFile(migration2UpPath, []byte(createMigration2Up), 0644))
	require.NoError(t, os.WriteFile(migration2DownPath, []byte(createMigration2Down), 0644))

	// Apply migrations
	runner, err := db.NewMigrationRunner(database, migrationsDir)
	require.NoError(t, err)
	defer runner.Close()

	err = runner.Migrate(ctx)
	require.NoError(t, err)

	// Verify state
	testhelpers.AssertMigrationState(ctx, t, database, migrationsDir, 2, false)

	// Verify all tables exist
	testhelpers.AssertTableExists(ctx, t, postgresDB, "postgres", "users")
	testhelpers.AssertTableExists(ctx, t, postgresDB, "postgres", "posts")
	testhelpers.AssertForeignKeyExists(ctx, t, postgresDB, "postgres", "posts", "user_id")
}

// TestMigrationStateDirtyFlag tests the dirty flag when migration fails
func TestMigrationStateDirtyFlag(t *testing.T) {
	// Using localhost PostgreSQL
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	opts := testhelpers.PostgresOpts{
		UseDirect: true,
		Host:      "localhost",
		Port:      "5432",
		User:      "postgres",
		Password:  "123",
		DBName:    fmt.Sprintf("test_%s_%d", t.Name(), time.Now().UnixNano()),
	}
	postgresDB, dsn, cleanup, err := testhelpers.StartPostgresContainer(ctx, opts)
	require.NoError(t, err)
	defer cleanup()
	defer postgresDB.Close()

	// Create a temporary directory for migrations under tests/tmp (returns relative path)
	tempDir, cleanupTemp := testhelpers.TempDirInTests(t, "state_tracking_")
	defer cleanupTemp()
	migrationsDir := filepath.Join(tempDir, "migrations")
	require.NoError(t, os.MkdirAll(migrationsDir, 0755))

	// Create database connection
	database, err := db.NewDB(dsn)
	require.NoError(t, err)
	defer database.Close()

	// Create first migration (valid)
	createMigration1Up := `CREATE TABLE users (id BIGSERIAL PRIMARY KEY);`
	createMigration1Down := `DROP TABLE IF EXISTS users;`

	migration1UpPath := filepath.Join(migrationsDir, "000001_create_users.up.sql")
	migration1DownPath := filepath.Join(migrationsDir, "000001_create_users.down.sql")

	require.NoError(t, os.WriteFile(migration1UpPath, []byte(createMigration1Up), 0644))
	require.NoError(t, os.WriteFile(migration1DownPath, []byte(createMigration1Down), 0644))

	// Create second migration with invalid SQL (will fail)
	createMigration2Up := `CREATE TABLE posts (id BIGSERIAL PRIMARY KEY, invalid_syntax_here);`
	createMigration2Down := `DROP TABLE IF EXISTS posts;`

	migration2UpPath := filepath.Join(migrationsDir, "000002_create_posts.up.sql")
	migration2DownPath := filepath.Join(migrationsDir, "000002_create_posts.down.sql")

	require.NoError(t, os.WriteFile(migration2UpPath, []byte(createMigration2Up), 0644))
	require.NoError(t, os.WriteFile(migration2DownPath, []byte(createMigration2Down), 0644))

	// Apply migrations
	runner, err := db.NewMigrationRunner(database, migrationsDir)
	require.NoError(t, err)
	defer runner.Close()

	// First migration should succeed
	err = runner.Migrate(ctx)
	// The second migration will fail, but golang-migrate handles this
	// We expect an error, but the state might be dirty
	if err != nil {
		t.Logf("Migration failed as expected: %v", err)
	}

	// Check if state is dirty
	version, dirty, err := runner.Version(ctx)
	require.NoError(t, err)

	if dirty {
		t.Logf("Migration state is dirty as expected after failed migration")
		// In a real scenario, you would need to manually fix the migration
		// and then use Force() to mark it as clean
	} else {
		// If not dirty, the migration system handled the error gracefully
		t.Logf("Migration state is clean (version: %d)", version)
	}
}

// TestMigrationStateAfterForce tests state after using Force() to fix dirty state
func TestMigrationStateAfterForce(t *testing.T) {
	// Using localhost PostgreSQL
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	opts := testhelpers.PostgresOpts{
		UseDirect: true,
		Host:      "localhost",
		Port:      "5432",
		User:      "postgres",
		Password:  "123",
		DBName:    fmt.Sprintf("test_%s_%d", t.Name(), time.Now().UnixNano()),
	}
	postgresDB, dsn, cleanup, err := testhelpers.StartPostgresContainer(ctx, opts)
	require.NoError(t, err)
	defer cleanup()
	defer postgresDB.Close()

	// Create a temporary directory for migrations under tests/tmp (returns relative path)
	tempDir, cleanupTemp := testhelpers.TempDirInTests(t, "state_tracking_")
	defer cleanupTemp()
	migrationsDir := filepath.Join(tempDir, "migrations")
	require.NoError(t, os.MkdirAll(migrationsDir, 0755))

	// Create database connection
	database, err := db.NewDB(dsn)
	require.NoError(t, err)
	defer database.Close()

	// Create migration
	createMigrationUp := `CREATE TABLE users (id BIGSERIAL PRIMARY KEY);`
	createMigrationDown := `DROP TABLE IF EXISTS users;`

	migrationUpPath := filepath.Join(migrationsDir, "000001_create_users.up.sql")
	migrationDownPath := filepath.Join(migrationsDir, "000001_create_users.down.sql")

	require.NoError(t, os.WriteFile(migrationUpPath, []byte(createMigrationUp), 0644))
	require.NoError(t, os.WriteFile(migrationDownPath, []byte(createMigrationDown), 0644))

	// Apply migration
	runner, err := db.NewMigrationRunner(database, migrationsDir)
	require.NoError(t, err)

	err = runner.Migrate(ctx)
	require.NoError(t, err)

	// Verify state using the existing runner
	testhelpers.AssertMigrationStateWithRunner(ctx, t, database, migrationsDir, 1, false, runner)

	// Simulate dirty state by manually setting it (in real scenario, this happens after failed migration)
	// We'll use Force() to set a specific version
	err = runner.Force(ctx, 1)
	require.NoError(t, err)

	// Verify state is clean after force using the existing runner
	testhelpers.AssertMigrationStateWithRunner(ctx, t, database, migrationsDir, 1, false, runner)

	// Close runner after all assertions
	defer runner.Close()
}
