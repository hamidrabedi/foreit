package migrations

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/forgego/forge/pkg/db"
	"github.com/forgego/forge/pkg/migrate"
	"github.com/forgego/forge/tests/testhelpers"
)

// TestMigrationApplySQLite tests migrations against SQLite using the migration system
// NOTE: Currently skipped because migration generator uses config which defaults to postgres
// TODO: Add support for specifying driver in migration generator
func TestMigrationApplySQLite(t *testing.T) {
	t.Skip("Skipping SQLite test - migration generator needs driver configuration support")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Create temporary directories
	tempDir, cleanupTemp := testhelpers.TempDirInTests(t, "sqlite_migration_")
	defer cleanupTemp()
	modelsDir := filepath.Join(tempDir, "models")
	migrationsDir := filepath.Join(tempDir, "migrations")
	require.NoError(t, os.MkdirAll(modelsDir, 0755))
	require.NoError(t, os.MkdirAll(migrationsDir, 0755))

	// Create a model file
	modelContent := `package models

import "github.com/forgego/forge/pkg/schema"

type User struct {
	schema.BaseSchema
}

func (User) Fields() []schema.Field {
	return []schema.Field{
		schema.Int64("id").Primary().AutoIncrement().Build(),
		schema.String("username").Required().MaxLength(150).Build(),
		schema.String("email").Required().MaxLength(255).Unique().Build(),
		schema.Bool("is_active").Default(true).Build(),
		schema.Timestamp("created_at").Default("CURRENT_TIMESTAMP").Build(),
	}
}

func (User) Meta() schema.Meta {
	return schema.Meta{TableName: "users"}
}

func (User) Relations() []schema.Relation {
	return []schema.Relation{}
}
`
	modelFile := filepath.Join(modelsDir, "user.go")
	require.NoError(t, os.WriteFile(modelFile, []byte(modelContent), 0644))

	// Generate migrations using the migration system
	gen, err := migrate.NewGenerator(modelsDir, migrationsDir)
	require.NoError(t, err)

	err = gen.GenerateMigrations("create_users")
	require.NoError(t, err)

	// Create SQLite database
	sqliteDB, err := testhelpers.StartSQLiteMemory("file::memory:?cache=shared")
	require.NoError(t, err)
	defer sqliteDB.Close()

	// Get DSN for db.NewDB
	dsn := "file::memory:?cache=shared"
	database, err := db.NewDBWithDriver("sqlite3", dsn, nil)
	require.NoError(t, err)
	defer database.Close()

	// Apply migrations using the migration runner
	runner, err := db.NewMigrationRunner(database, migrationsDir)
	require.NoError(t, err)
	defer runner.Close()

	err = runner.Migrate(ctx)
	require.NoError(t, err)

	// Verify table exists
	testhelpers.AssertTableExists(ctx, t, sqliteDB, "sqlite", "users")

	// Verify columns
	testhelpers.AssertColumnExists(ctx, t, sqliteDB, "sqlite", "users", "id")
	testhelpers.AssertColumnExists(ctx, t, sqliteDB, "sqlite", "users", "username")
	testhelpers.AssertColumnExists(ctx, t, sqliteDB, "sqlite", "users", "email")
}

// TestMigrationApplyPostgres tests migrations against Postgres using the migration system
func TestMigrationApplyPostgres(t *testing.T) {
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

	t.Logf("Connected to Postgres: %s", dsn)

	// Create temporary directories
	tempDir, cleanupTemp := testhelpers.TempDirInTests(t, "postgres_migration_")
	defer cleanupTemp()
	modelsDir := filepath.Join(tempDir, "models")
	migrationsDir := filepath.Join(tempDir, "migrations")
	require.NoError(t, os.MkdirAll(modelsDir, 0755))
	require.NoError(t, os.MkdirAll(migrationsDir, 0755))

	// Create User model
	userModel := `package models

import "github.com/forgego/forge/pkg/schema"

type User struct {
	schema.BaseSchema
}

func (User) Fields() []schema.Field {
	return []schema.Field{
		schema.Int64("id").Primary().AutoIncrement().Build(),
		schema.String("username").Required().MaxLength(150).Unique().Build(),
		schema.String("email").Required().MaxLength(255).Unique().Build(),
		schema.Bool("is_active").Default(true).Build(),
		schema.Timestamp("created_at").Default("NOW()").Build(),
	}
}

func (User) Meta() schema.Meta {
	return schema.Meta{TableName: "users"}
}

func (User) Relations() []schema.Relation {
	return []schema.Relation{}
}
`
	require.NoError(t, os.WriteFile(filepath.Join(modelsDir, "user.go"), []byte(userModel), 0644))

	// Create Post model with FK
	postModel := `package models

import "github.com/forgego/forge/pkg/schema"

type Post struct {
	schema.BaseSchema
}

func (Post) Fields() []schema.Field {
	return []schema.Field{
		schema.Int64("id").Primary().AutoIncrement().Build(),
		schema.String("title").Required().MaxLength(200).Build(),
		schema.Text("content").Build(),
		schema.Int64("user_id").Required().Build(),
		schema.Timestamp("created_at").Default("NOW()").Build(),
	}
}

func (Post) Meta() schema.Meta {
	return schema.Meta{TableName: "posts"}
}

func (Post) Relations() []schema.Relation {
	return []schema.Relation{
		schema.ForeignKey("user_id", "User").Required().OnDelete(schema.CascadeCASCADE).RelatedName("posts").Build(),
	}
}
`
	require.NoError(t, os.WriteFile(filepath.Join(modelsDir, "post.go"), []byte(postModel), 0644))

	// Generate migrations using the migration system
	gen, err := migrate.NewGenerator(modelsDir, migrationsDir)
	require.NoError(t, err)

	err = gen.GenerateMigrations("create_users_and_posts")
	require.NoError(t, err)

	// Apply migrations using the migration runner
	database, err := db.NewDB(dsn)
	require.NoError(t, err)
	defer database.Close()

	runner, err := db.NewMigrationRunner(database, migrationsDir)
	require.NoError(t, err)
	defer runner.Close()

	err = runner.Migrate(ctx)
	require.NoError(t, err)

	// Verify tables and FK
	testhelpers.AssertTableExists(ctx, t, postgresDB, "postgres", "users")
	testhelpers.AssertTableExists(ctx, t, postgresDB, "postgres", "posts")
	testhelpers.AssertForeignKeyExists(ctx, t, postgresDB, "postgres", "posts", "user_id")

	// Test: FK constraint on delete
	insertUserSQL := `INSERT INTO users (username, email) VALUES ('testuser', 'test@example.com')`
	testhelpers.RunSQLExpectSuccess(ctx, t, postgresDB, insertUserSQL)

	insertPostSQL := `INSERT INTO posts (title, user_id) VALUES ('Test Post', 1)`
	testhelpers.RunSQLExpectSuccess(ctx, t, postgresDB, insertPostSQL)

	// Verify rows exist
	testhelpers.AssertRowCount(ctx, t, postgresDB, "users", 1)
	testhelpers.AssertRowCount(ctx, t, postgresDB, "posts", 1)

	// Delete user and verify cascade
	deleteUserSQL := `DELETE FROM users WHERE id = 1`
	testhelpers.RunSQLExpectSuccess(ctx, t, postgresDB, deleteUserSQL)

	// Posts should be deleted due to CASCADE
	testhelpers.AssertRowCount(ctx, t, postgresDB, "posts", 0)
}

