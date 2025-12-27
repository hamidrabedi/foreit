package testhelpers

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"

	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

// PostgresOpts contains options for starting a Postgres container
type PostgresOpts struct {
	Version     string
	User        string
	Password    string
	DBName      string
	Env         []string
	WaitTimeout time.Duration
}

// DefaultPostgresOpts returns sensible defaults
func DefaultPostgresOpts() PostgresOpts {
	return PostgresOpts{
		Version:     "15",
		User:        "testuser",
		Password:    "testpass",
		DBName:      "testdb",
		WaitTimeout: 30 * time.Second,
	}
}

// StartPostgresContainer starts an ephemeral Postgres container using Dockertest
func StartPostgresContainer(ctx context.Context, opts PostgresOpts) (*sql.DB, string, func() error, error) {
	if opts.Version == "" {
		opts.Version = "15"
	}
	if opts.WaitTimeout == 0 {
		opts.WaitTimeout = 30 * time.Second
	}

	pool, err := dockertest.NewPool("")
	if err != nil {
		return nil, "", nil, fmt.Errorf("could not connect to docker: %w", err)
	}

	// Set max wait time
	pool.MaxWait = opts.WaitTimeout

	env := []string{
		fmt.Sprintf("POSTGRES_USER=%s", opts.User),
		fmt.Sprintf("POSTGRES_PASSWORD=%s", opts.Password),
		fmt.Sprintf("POSTGRES_DB=%s", opts.DBName),
	}
	env = append(env, opts.Env...)

	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "postgres",
		Tag:        opts.Version,
		Env:        env,
		ExposedPorts: []string{"5432"},
	}, func(config *docker.HostConfig) {
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{Name: "no"}
	})
	if err != nil {
		return nil, "", nil, fmt.Errorf("could not start resource: %w", err)
	}

	// Get port
	port := resource.GetPort("5432/tcp")
	dsn := fmt.Sprintf("postgres://%s:%s@localhost:%s/%s?sslmode=disable",
		opts.User, opts.Password, port, opts.DBName)

	cleanup := func() error {
		return pool.Purge(resource)
	}

	// Retry connection
	var db *sql.DB
	retries := 0
	maxRetries := 30
	for retries < maxRetries {
		db, err = sql.Open("postgres", dsn)
		if err == nil {
			if err = db.PingContext(ctx); err == nil {
				break
			}
			db.Close()
		}
		retries++
		time.Sleep(100 * time.Millisecond)
	}

	if retries >= maxRetries {
		cleanup()
		return nil, "", nil, fmt.Errorf("could not connect to postgres after %d retries", maxRetries)
	}

	return db, dsn, cleanup, nil
}

// StartSQLiteMemory returns an in-memory SQLite connection for tests
func StartSQLiteMemory(dsn string) (*sql.DB, error) {
	if dsn == "" {
		dsn = "file::memory:?cache=shared"
	}
	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to create sqlite DB: %w", err)
	}

	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to ping sqlite: %w", err)
	}

	return db, nil
}

// WaitForDBReady polls the DB connection until ready or timeout
func WaitForDBReady(ctx context.Context, db *sql.DB, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			if err := db.PingContext(ctx); err == nil {
				return nil
			}
			if time.Now().After(deadline) {
				return fmt.Errorf("database not ready after %v", timeout)
			}
		}
	}
}
