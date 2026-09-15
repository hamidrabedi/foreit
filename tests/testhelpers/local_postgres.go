package testhelpers

import (
	"fmt"
	"net/url"
	"os"
	"testing"
	"time"
)

func testDatabaseURL(defaultURL string) string {
	if u := os.Getenv("FORGE_TEST_DATABASE_URL"); u != "" {
		return u
	}
	return defaultURL
}

func requireDB() bool {
	return os.Getenv("FORGE_REQUIRE_DB") == "1"
}

func dbUnavailableAction() string {
	if requireDB() {
		return "fatal"
	}
	return "skip"
}

func skipOrFailNoDB(t testing.TB, format string, args ...any) {
	t.Helper()
	msg := fmt.Sprintf(format, args...)
	if requireDB() {
		t.Fatalf("%s (FORGE_REQUIRE_DB=1 is set)", msg)
	}
	t.Skipf("%s (set FORGE_REQUIRE_DB=1 to turn into failure)", msg)
}

func applyDatabaseURLToOpts(opts *PostgresOpts) {
	raw := os.Getenv("FORGE_TEST_DATABASE_URL")
	if raw == "" {
		return
	}
	u, err := url.Parse(raw)
	if err != nil {
		return
	}
	if h := u.Hostname(); h != "" {
		opts.Host = h
	}
	if p := u.Port(); p != "" {
		opts.Port = p
	}
	if u.User != nil {
		opts.User = u.User.Username()
		if pass, ok := u.User.Password(); ok {
			opts.Password = pass
		}
	}
	if u.RawQuery != "" {
		opts.RawQuery = u.RawQuery
	}
}

// LocalPostgresOpts returns PostgreSQL options for local database
// Uses localhost with user "postgres" and password "123"
// Creates a unique database name for each test
func LocalPostgresOpts(testName string) PostgresOpts {
	opts := PostgresOpts{
		UseDirect: true,
		Host:      "localhost",
		Port:      "5432",
		User:      "postgres",
		Password:  "123",
		DBName:    fmt.Sprintf("test_%s_%d", sanitizeTestName(testName), time.Now().UnixNano()),
	}
	applyDatabaseURLToOpts(&opts)
	return opts
}
