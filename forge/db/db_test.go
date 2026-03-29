package db

import (
	"context"
	"strings"
	"testing"
	"time"
)

func TestNewDBFromConfig_NilConfig(t *testing.T) {
	database, err := NewDBFromConfig(nil)
	if err == nil {
		t.Fatal("expected error for nil config, got nil")
	}
	if database != nil {
		t.Fatalf("expected nil database for nil config, got %#v", database)
	}
	if !strings.Contains(err.Error(), "config is nil") {
		t.Fatalf("expected nil-config error message, got %q", err.Error())
	}
}

func TestPing_NilDB(t *testing.T) {
	var database DB
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := database.Ping(ctx)
	if err == nil {
		t.Fatal("expected error for nil database, got nil")
	}
	if !strings.Contains(err.Error(), "not initialized") {
		t.Fatalf("expected not initialized error message, got %q", err.Error())
	}
}

func TestIsConnected_NilDB(t *testing.T) {
	var database DB
	if database.IsConnected() {
		t.Fatal("expected IsConnected to return false for nil database")
	}
}

func TestDB_Dialect_Nil(t *testing.T) {
	var db *DB = nil
	if db.Dialect() != nil {
		t.Fatal("expected nil dialect for nil db")
	}
}

func TestDB_SetDialect_Nil(t *testing.T) {
	var db *DB = nil
	// Should not panic
	db.SetDialect(nil)
}

func TestDB_PoolStats_Nil(t *testing.T) {
	var db *DB = nil
	if db.PoolStats() != nil {
		t.Fatal("expected nil stats for nil db")
	}
}

func TestDB_RebindPlaceholders_Nil(t *testing.T) {
	var db *DB = nil
	query := "SELECT * FROM users WHERE id = ?"
	if db.RebindPlaceholders(query) != query {
		t.Fatal("expected unmodified query for nil db")
	}
}

func TestDB_Rebind_Nil(t *testing.T) {
	var db *DB = nil
	query := "SELECT * FROM users WHERE id = ?"
	if db.Rebind(query) != query {
		t.Fatal("expected unmodified query for nil db")
	}
}

func TestDB_Ping_Nil(t *testing.T) {
	var db *DB = nil
	err := db.Ping(context.Background())
	if err == nil {
		t.Fatal("expected error for nil db")
	}
	if err.Error() != "database connection is nil" {
		t.Fatalf("expected 'database connection is nil', got: %v", err)
	}
}

func TestDB_IsConnected_NilReceiver(t *testing.T) {
	var db *DB = nil
	if db.IsConnected() {
		t.Fatal("expected false for nil db")
	}
}

func TestNewDBWithValidation_InvalidDSN(t *testing.T) {
	// Test with an invalid DSN that will fail to connect
	database, err := NewDBWithValidation("invalid:dsn@tcp(localhost:3306)/nonexistent")
	if err == nil {
		t.Fatal("expected error for invalid DSN, got nil")
	}
	if database != nil {
		t.Fatalf("expected nil database for invalid DSN, got %#v", database)
	}
	// The error could be either "failed to connect" or "failed to ping database"
	// depending on which driver is attempted first
	errStr := err.Error()
	if !strings.Contains(errStr, "failed to connect") && !strings.Contains(errStr, "failed to ping") {
		t.Fatalf("expected connection-related error message, got %q", errStr)
	}
}

