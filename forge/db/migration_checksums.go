package db

import (
	"context"
	"crypto/sha256"
	"database/sql/driver"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/golang-migrate/migrate/v4"
)

// migrationChecksumLockKey identifies Forge's migration-and-checksum critical section.
// It is distinct from golang-migrate's per-step advisory lock.
const migrationChecksumLockKey int64 = 0x466f7267654d6967

const postgresMigrationChecksumsDDL = `CREATE TABLE IF NOT EXISTS forge_migration_checksums (version BIGINT PRIMARY KEY, up_sha256 CHAR(64) NOT NULL, down_sha256 CHAR(64), applied_at TIMESTAMPTZ NOT NULL DEFAULT now())`
const sqliteMigrationChecksumsDDL = `CREATE TABLE IF NOT EXISTS forge_migration_checksums (version INTEGER PRIMARY KEY, up_sha256 TEXT NOT NULL, down_sha256 TEXT, applied_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP)`

func (mr *MigrationRunner) withMigrationChecksums(ctx context.Context, run func() error) (err error) {
	ddl := sqliteMigrationChecksumsDDL
	if mr.db.Driver == "postgres" || mr.db.Driver == "postgresql" {
		ddl = postgresMigrationChecksumsDDL
		conn, connErr := mr.db.DB.Conn(ctx)
		if connErr != nil {
			return connErr
		}
		defer conn.Close()
		if _, err = conn.ExecContext(ctx, "SELECT pg_advisory_lock($1)", migrationChecksumLockKey); err != nil {
			return err
		}
		defer func() {
			unlockCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			_, unlockErr := conn.ExecContext(unlockCtx, "SELECT pg_advisory_unlock($1)", migrationChecksumLockKey)
			if unlockErr != nil {
				// Discard the session so a failed unlock cannot leak a lock into the pool.
				_ = conn.Raw(func(any) error { return driver.ErrBadConn })
				err = errors.Join(err, fmt.Errorf("release migration checksum lock: %w", unlockErr))
			}
		}()
	}
	// SQLite serializes writers itself; no additional advisory lock is needed.
	if _, err := mr.db.ExecContext(ctx, ddl); err != nil {
		return fmt.Errorf("create migration checksum table: %w", err)
	}
	return run()
}

func migrationFileChecksums(path string, version uint) (string, *string, error) {
	entries, err := os.ReadDir(path)
	if err != nil {
		return "", nil, fmt.Errorf("read migration version %d: %w", version, err)
	}
	for _, entry := range entries {
		name := entry.Name()
		prefix, _, ok := strings.Cut(name, "_")
		n, parseErr := strconv.ParseUint(prefix, 10, 64)
		if entry.IsDir() || !ok || parseErr != nil || n != uint64(version) || !strings.HasSuffix(name, ".up.sql") {
			continue
		}
		up, err := os.ReadFile(filepath.Join(path, name))
		if err != nil {
			return "", nil, fmt.Errorf("read up migration version %d: %w", version, err)
		}
		upHash := fmt.Sprintf("%x", sha256.Sum256(up))
		down, err := os.ReadFile(filepath.Join(path, strings.TrimSuffix(name, ".up.sql")+".down.sql"))
		if errors.Is(err, os.ErrNotExist) {
			return upHash, nil, nil
		}
		if err != nil {
			return "", nil, fmt.Errorf("read down migration version %d: %w", version, err)
		}
		downHash := fmt.Sprintf("%x", sha256.Sum256(down))
		return upHash, &downHash, nil
	}
	return "", nil, fmt.Errorf("up migration file for version %d not found", version)
}

func (mr *MigrationRunner) applyChecksumSteps(ctx context.Context, target *uint) error {
	if target != nil {
		if _, _, err := migrationFileChecksums(mr.migrationsPath, *target); err != nil {
			return err
		}
	}
	for {
		if err := ctx.Err(); err != nil {
			return err
		}
		version, _, err := mr.migrate.Version()
		if err != nil && err != migrate.ErrNilVersion {
			return err
		}
		if target != nil && err == nil && version == *target {
			return nil
		}
		err = mr.migrate.Steps(1)
		var short migrate.ErrShortLimit
		// Steps(1) returns bare os.ErrNotExist when the source has no next version.
		// File read failures are wrapped and must still be reported.
		if err == os.ErrNotExist || errors.Is(err, migrate.ErrNoChange) || errors.As(err, &short) {
			return nil
		}
		if err != nil {
			return fmt.Errorf("apply migration: %w", err)
		}
		version, _, err = mr.migrate.Version()
		if err != nil {
			return err
		}
		up, down, err := migrationFileChecksums(mr.migrationsPath, version)
		if err != nil {
			return err
		}
		_, err = mr.db.ExecContext(ctx, `INSERT INTO forge_migration_checksums (version, up_sha256, down_sha256) VALUES ($1, $2, $3) ON CONFLICT (version) DO UPDATE SET up_sha256 = excluded.up_sha256, down_sha256 = excluded.down_sha256, applied_at = CURRENT_TIMESTAMP`, version, up, down)
		if err != nil {
			return fmt.Errorf("record migration version %d: %w", version, err)
		}
	}
}

func (mr *MigrationRunner) rollbackChecksumStep(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if err := mr.migrate.Steps(-1); err != nil {
		return err
	}
	version, _, err := mr.migrate.Version()
	if err == migrate.ErrNilVersion {
		_, err = mr.db.ExecContext(ctx, "DELETE FROM forge_migration_checksums")
	} else if err == nil {
		_, err = mr.db.ExecContext(ctx, "DELETE FROM forge_migration_checksums WHERE version > $1", version)
	}
	return err
}
