package state

import (
	"testing"

	"github.com/forgego/forge/pkg/migrations/core"
)

func TestTableParser_ParseCreateTable(t *testing.T) {
	parser := NewTableParser()

	tests := []struct {
		name     string
		sql      string
		wantErr  bool
		checkFn  func(t *testing.T, change *core.CreateTable)
	}{
		{
			name: "simple table with id and email",
			sql: `CREATE TABLE users (
				id BIGINT PRIMARY KEY,
				email VARCHAR(255) NOT NULL
			)`,
			wantErr: false,
			checkFn: func(t *testing.T, change *core.CreateTable) {
				if change.Table.Meta.TableName != "users" {
					t.Errorf("Expected table name 'users', got '%s'", change.Table.Meta.TableName)
				}
				if len(change.Table.Fields) < 2 {
					t.Errorf("Expected at least 2 fields, got %d", len(change.Table.Fields))
				}
			},
		},
		{
			name: "table with autoincrement",
			sql: `CREATE TABLE posts (
				id INTEGER PRIMARY KEY AUTOINCREMENT,
				title TEXT NOT NULL
			)`,
			wantErr: false,
			checkFn: func(t *testing.T, change *core.CreateTable) {
				if change.Table.Meta.TableName != "posts" {
					t.Errorf("Expected table name 'posts', got '%s'", change.Table.Meta.TableName)
				}
				// Check for autoincrement field
				foundAutoInc := false
				for _, field := range change.Table.Fields {
					if field.AutoIncrement {
						foundAutoInc = true
						break
					}
				}
				if !foundAutoInc {
					t.Error("Expected to find autoincrement field")
				}
			},
		},
		{
			name: "table with unique constraint",
			sql: `CREATE TABLE products (
				id BIGINT PRIMARY KEY,
				sku VARCHAR(100) UNIQUE NOT NULL
			)`,
			wantErr: false,
			checkFn: func(t *testing.T, change *core.CreateTable) {
				if change.Table.Meta.TableName != "products" {
					t.Errorf("Expected table name 'products', got '%s'", change.Table.Meta.TableName)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			change, err := parser.ParseCreateTable(tt.sql)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseCreateTable() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && change != nil && tt.checkFn != nil {
				tt.checkFn(t, change)
			}
		})
	}
}


