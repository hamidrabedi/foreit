package db

import (
	"context"
	"crypto/sha256"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/forgego/forge/internal/testutils"
	"github.com/stretchr/testify/require"
)

func TestMigrationFileChecksums(t *testing.T) {
	dir := t.TempDir()
	lf := []byte("SELECT 1;\n")
	crlf := []byte("SELECT 1;\r\n")
	require.NoError(t, os.WriteFile(filepath.Join(dir, "000007_first.up.sql"), lf, 0600))
	require.NoError(t, os.WriteFile(filepath.Join(dir, "42_second.up.sql"), crlf, 0600))
	require.NoError(t, os.WriteFile(filepath.Join(dir, "42_second.down.sql"), lf, 0600))
	up, down, err := migrationFileChecksums(dir, 7)
	require.NoError(t, err)
	require.Equal(t, fmt.Sprintf("%x", sha256.Sum256(lf)), up)
	require.Nil(t, down)
	other, down, err := migrationFileChecksums(dir, 42)
	require.NoError(t, err)
	require.NotEqual(t, up, other)
	require.Equal(t, fmt.Sprintf("%x", sha256.Sum256(crlf)), other)
	require.NotNil(t, down)
	require.Equal(t, up, *down)
	_, _, err = migrationFileChecksums(dir, 8)
	require.ErrorContains(t, err, "version 8")
	require.NoError(t, os.Remove(filepath.Join(dir, "000007_first.up.sql")))
	require.NoError(t, os.Symlink(filepath.Join(dir, "absent"), filepath.Join(dir, "000007_first.up.sql")))
	_, _, err = migrationFileChecksums(dir, 7)
	require.ErrorContains(t, err, "version 7")
}

func TestChecksumBaseline_RecordsEachAppliedVersion(t *testing.T) {
	sqlDB := testutils.SetupTestDB(t)
	t.Cleanup(func() { _ = sqlDB.Close() })
	// Use an isolated schema without altering other tests' migration state.
	schema := fmt.Sprintf("checksum_baseline_%d", time.Now().UnixNano())
	_, err := sqlDB.Exec("CREATE SCHEMA " + schema)
	require.NoError(t, err)
	t.Cleanup(func() { _, _ = sqlDB.Exec("DROP SCHEMA " + schema + " CASCADE") })
	// SET is session-local, so initialize every connection used by this runner.
	sqlDB.SetMaxOpenConns(3)
	sqlDB.SetMaxIdleConns(3)
	forConnections := make([]interface{ Close() error }, 0, 3)
	for i := 0; i < 3; i++ {
		conn, err := sqlDB.Conn(context.Background())
		require.NoError(t, err)
		_, err = conn.ExecContext(context.Background(), "SET search_path TO "+schema)
		require.NoError(t, err)
		forConnections = append(forConnections, conn)
	}
	for _, conn := range forConnections {
		require.NoError(t, conn.Close())
	}
	testChecksumBaseline(t, &DB{DB: sqlDB, Driver: "postgres"})
}

func TestChecksumBaseline_SQLiteRecordsEachAppliedVersion(t *testing.T) {
	database, err := NewDBWithDriver("sqlite3", filepath.Join(t.TempDir(), "test.db"))
	require.NoError(t, err)
	t.Cleanup(func() { _ = database.DB.Close() })
	testChecksumBaseline(t, database)
}

func testChecksumBaseline(t *testing.T, database *DB) {
	t.Helper()
	ctx := context.Background()
	dir := t.TempDir()
	expected := map[uint]string{}
	for i := 1; i <= 3; i++ {
		contents := []byte(fmt.Sprintf("CREATE TABLE checksum_test_%d (id INTEGER PRIMARY KEY);\n", i))
		require.NoError(t, os.WriteFile(filepath.Join(dir, fmt.Sprintf("%06d_test.up.sql", i)), contents, 0600))
		require.NoError(t, os.WriteFile(filepath.Join(dir, fmt.Sprintf("%06d_test.down.sql", i)), []byte(fmt.Sprintf("DROP TABLE checksum_test_%d;\n", i)), 0600))
		expected[uint(i)] = fmt.Sprintf("%x", sha256.Sum256(contents))
	}
	runner, err := NewMigrationRunner(database, dir)
	require.NoError(t, err)
	assertRows := func(count int) {
		t.Helper()
		rows, err := database.Query("SELECT version, up_sha256, down_sha256, applied_at FROM forge_migration_checksums ORDER BY version")
		require.NoError(t, err)
		defer rows.Close()
		n := 0
		for rows.Next() {
			var version uint
			var up, down string
			var applied any
			require.NoError(t, rows.Scan(&version, &up, &down, &applied))
			require.Equal(t, expected[version], up)
			require.Equal(t, fmt.Sprintf("%x", sha256.Sum256([]byte(fmt.Sprintf("DROP TABLE checksum_test_%d;\n", version)))), down)
			require.NotNil(t, applied)
			n++
		}
		require.NoError(t, rows.Err())
		require.Equal(t, count, n)
	}
	require.NoError(t, runner.Up(ctx))
	assertRows(3)
	require.NoError(t, runner.Rollback(ctx))
	assertRows(2)
	require.NoError(t, runner.Up(ctx))
	assertRows(3)
	require.NoError(t, runner.RollbackTo(ctx, 1))
	assertRows(1)
	require.NoError(t, runner.MigrateTo(ctx, 3))
	assertRows(3)
	require.NoError(t, runner.RollbackSteps(ctx, 3))
	assertRows(0)
}
