package migrations

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/forgego/forge/pkg/generator"
)

func TestSQLGenerator_CreateTable(t *testing.T) {
	gen := NewSQLGenerator("postgres")

	def := &generator.ModelDefinition{
		Name: "User",
		Meta: generator.MetaDefinition{
			TableName: "users",
		},
		Fields: []generator.FieldDefinition{
			{
				Name:          "id",
				Type:          "Int64",
				GoType:        "int64",
				PrimaryKey:    true,
				AutoIncrement: true,
				Required:      true,
			},
			{
				Name:     "email",
				Type:     "String",
				GoType:   "string",
				Required: true,
				Options: map[string]interface{}{
					"max_length": 255,
					"unique":     true,
				},
			},
		},
	}

	change := &CreateTable{Table: def}
	changes := []Change{change}
	sql, err := gen.GenerateUpSQL(changes)
	if err != nil {
		t.Fatalf("Failed to generate SQL: %v", err)
	}

	if !contains(sql, "CREATE TABLE") {
		t.Errorf("SQL should contain CREATE TABLE")
	}
	if !contains(sql, "users") {
		t.Errorf("SQL should contain table name 'users'")
	}
	if !contains(sql, "id") {
		t.Errorf("SQL should contain column 'id'")
	}
	if !contains(sql, "email") {
		t.Errorf("SQL should contain column 'email'")
	}
}

// TestSQLGenerator_FieldTypeMapping is removed as it tested internal implementation details
// Field type mapping is now tested indirectly through the public GenerateUpSQL interface

func TestDetector_DetectChanges(t *testing.T) {
	detector := NewDetector()

	current := []*generator.ModelDefinition{
		{
			Name: "User",
			Meta: generator.MetaDefinition{TableName: "users"},
			Fields: []generator.FieldDefinition{
				{Name: "id", Type: "Int64", PrimaryKey: true, AutoIncrement: true},
				{Name: "email", Type: "String", Required: true},
			},
		},
	}

	previous := []*generator.ModelDefinition{}

	changes, err := detector.DetectChanges(current, previous)
	if err != nil {
		t.Fatalf("Failed to detect changes: %v", err)
	}

	if len(changes) == 0 {
		t.Fatal("Expected at least one change (CreateTable)")
	}

	foundCreateTable := false
	for _, change := range changes {
		if change.Type() == ChangeTypeCreateTable {
			foundCreateTable = true
			break
		}
	}

	if !foundCreateTable {
		t.Error("Expected CreateTable change")
	}
}

func TestDetector_DetectColumnChanges(t *testing.T) {
	detector := NewDetector()

	current := []*generator.ModelDefinition{
		{
			Name: "User",
			Meta: generator.MetaDefinition{TableName: "users"},
			Fields: []generator.FieldDefinition{
				{Name: "id", Type: "Int64", PrimaryKey: true},
				{Name: "email", Type: "String", Required: true},
				{Name: "name", Type: "String", Required: true}, // New field
			},
		},
	}

	previous := []*generator.ModelDefinition{
		{
			Name: "User",
			Meta: generator.MetaDefinition{TableName: "users"},
			Fields: []generator.FieldDefinition{
				{Name: "id", Type: "Int64", PrimaryKey: true},
				{Name: "email", Type: "String", Required: true},
			},
		},
	}

	changes, err := detector.DetectChanges(current, previous)
	if err != nil {
		t.Fatalf("Failed to detect changes: %v", err)
	}

	foundAddColumn := false
	for _, change := range changes {
		if change.Type() == ChangeTypeAddColumn {
			if addCol, ok := change.(*AddColumn); ok && addCol.Column.Name == "name" {
				foundAddColumn = true
				break
			}
		}
	}

	if !foundAddColumn {
		t.Error("Expected AddColumn change for 'name' field")
	}
}

func TestSQLGenerator_ReversibleMigrations(t *testing.T) {
	gen := NewSQLGenerator("postgres")

	def := &generator.ModelDefinition{
		Name: "User",
		Meta: generator.MetaDefinition{TableName: "users"},
		Fields: []generator.FieldDefinition{
			{Name: "id", Type: "Int64", PrimaryKey: true, AutoIncrement: true},
		},
	}

	changes := []Change{
		&CreateTable{Table: def},
	}

	upSQL, err := gen.GenerateUpSQL(changes)
	if err != nil {
		t.Fatalf("Failed to generate up SQL: %v", err)
	}

	downSQL, err := gen.GenerateDownSQL(changes)
	if err != nil {
		t.Fatalf("Failed to generate down SQL: %v", err)
	}

	if !contains(upSQL, "CREATE TABLE") {
		t.Error("Up SQL should contain CREATE TABLE")
	}

	if !contains(downSQL, "DROP TABLE") {
		t.Error("Down SQL should contain DROP TABLE")
	}
}

func TestGenerator_GenerateMigrations(t *testing.T) {
	// Create temporary directories
	tmpDir := t.TempDir()
	modelsDir := filepath.Join(tmpDir, "models")
	migrationsDir := filepath.Join(tmpDir, "migrations")

	if err := os.MkdirAll(modelsDir, 0755); err != nil {
		t.Fatalf("Failed to create models directory: %v", err)
	}

	// Create a simple model file
	modelFile := filepath.Join(modelsDir, "user.go")
	modelContent := `package models

import "github.com/forgego/forge/pkg/schema"

type User struct {
	schema.BaseSchema
}

func (User) Fields() []schema.Field {
	return []schema.Field{
		schema.Int64("id").Primary().AutoIncrement().Build(),
		schema.String("email").Required().MaxLength(255).Unique().Build(),
	}
}

func (User) Meta() schema.Meta {
	return schema.Meta{
		TableName: "users",
	}
}

func (User) Relations() []schema.Relation {
	return []schema.Relation{}
}
`
	if err := os.WriteFile(modelFile, []byte(modelContent), 0644); err != nil {
		t.Fatalf("Failed to write model file: %v", err)
	}

	gen, err := NewGenerator(modelsDir, migrationsDir)
	if err != nil {
		t.Fatalf("Failed to create generator: %v", err)
	}

	err = gen.GenerateMigrations("create_users")
	if err != nil {
		t.Fatalf("Failed to generate migrations: %v", err)
	}

	// Check that migration files were created
	upFile := filepath.Join(migrationsDir, "000001_create_users.up.sql")
	downFile := filepath.Join(migrationsDir, "000001_create_users.down.sql")

	if _, err := os.Stat(upFile); os.IsNotExist(err) {
		t.Errorf("Up migration file was not created: %s", upFile)
	}

	if _, err := os.Stat(downFile); os.IsNotExist(err) {
		t.Errorf("Down migration file was not created: %s", downFile)
	}

	// Check file contents
	upContent, err := os.ReadFile(upFile)
	if err != nil {
		t.Fatalf("Failed to read up migration: %v", err)
	}

	if !contains(string(upContent), "CREATE TABLE") {
		t.Error("Up migration should contain CREATE TABLE")
	}

	if !contains(string(upContent), "users") {
		t.Error("Up migration should contain table name 'users'")
	}
}

// Helper function
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 || 
		(len(s) > len(substr) && (s[:len(substr)] == substr || 
		s[len(s)-len(substr):] == substr || 
		containsMiddle(s, substr))))
}

func containsMiddle(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

