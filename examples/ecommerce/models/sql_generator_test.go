package models

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/forgego/forge/codegen"
)

// TestGenerateSQL verifies that SQL can be generated from all models
// This test generates SQL CREATE TABLE statements and verifies they are correct
func TestGenerateSQL(t *testing.T) {
	// Generate code first using CLI tool
	setupTest(t)

	parser := generator.NewASTParser()
	definitions, err := parser.ParseDirectory(".")
	if err != nil {
		t.Fatalf("Failed to parse models: %v", err)
	}

	if len(definitions) == 0 {
		t.Fatal("No model definitions found")
	}

	// Generate SQL for each model
	sqlStatements := []string{}
	for _, def := range definitions {
		sql, err := generateSQLForModel(def, "postgres")
		if err != nil {
			t.Errorf("Failed to generate SQL for %s: %v", def.Name, err)
			continue
		}
		sqlStatements = append(sqlStatements, sql)
		t.Logf("✅ Generated SQL for %s", def.Name)
	}

	if len(sqlStatements) != len(definitions) {
		t.Errorf("Expected %d SQL statements, got %d", len(definitions), len(sqlStatements))
	}

	// Verify SQL contains expected elements
	for i, sql := range sqlStatements {
		def := definitions[i]
		verifySQL(t, def, sql)
	}

	t.Logf("✅ Successfully generated SQL for all %d models", len(definitions))
}

// generateSQLForModel generates a CREATE TABLE statement for a model
func generateSQLForModel(def *generator.ModelDefinition, driver string) (string, error) {
	var sql strings.Builder

	tableName := def.Meta.TableName
	if tableName == "" {
		tableName = strings.ToLower(def.Name) + "s"
	}

	sql.WriteString(fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (\n", tableName))

	// Generate column definitions
	columns := []string{}
	for _, field := range def.Fields {
		colDef, err := buildColumnSQL(field, driver)
		if err != nil {
			return "", fmt.Errorf("failed to build column %s: %w", field.Name, err)
		}
		columns = append(columns, "    "+colDef)
	}

	sql.WriteString(strings.Join(columns, ",\n"))
	sql.WriteString("\n);\n")

	// Generate indexes
	for _, idx := range def.Meta.Indexes {
		if len(idx.Fields) == 0 {
			continue
		}
		indexType := "INDEX"
		if idx.Unique {
			indexType = "UNIQUE INDEX"
		}
		indexName := idx.Name
		if indexName == "" {
			indexName = fmt.Sprintf("idx_%s_%s", tableName, strings.Join(idx.Fields, "_"))
		}
		escapedFields := make([]string, len(idx.Fields))
		for i, field := range idx.Fields {
			escapedFields[i] = fmt.Sprintf(`"%s"`, field)
		}
		sql.WriteString(fmt.Sprintf("\nCREATE %s IF NOT EXISTS %s ON %s (%s);\n",
			indexType, indexName, tableName, strings.Join(escapedFields, ", ")))
	}

	return sql.String(), nil
}

// buildColumnSQL builds SQL column definition from field definition
func buildColumnSQL(field generator.FieldDefinition, driver string) (string, error) {
	columnName := field.Name
	sqlType := mapFieldTypeToSQL(field.Type, field.GoType, field.Options, driver)

	var parts []string
	parts = append(parts, fmt.Sprintf(`"%s"`, columnName))
	parts = append(parts, sqlType)

	// Handle primary key and auto-increment
	if field.PrimaryKey && field.AutoIncrement {
		if driver == "sqlite" || driver == "sqlite3" {
			parts = []string{fmt.Sprintf(`"%s"`, columnName), "INTEGER", "PRIMARY KEY", "AUTOINCREMENT"}
		} else {
			parts = []string{fmt.Sprintf(`"%s"`, columnName), sqlType, "PRIMARY KEY", "GENERATED ALWAYS AS IDENTITY"}
		}
	} else if field.PrimaryKey {
		parts = append(parts, "PRIMARY KEY")
	} else if field.AutoIncrement {
		if driver != "sqlite" && driver != "sqlite3" {
			parts = append(parts, "GENERATED ALWAYS AS IDENTITY")
		}
	}

	// Handle required/optional
	if field.Required && !field.PrimaryKey {
		parts = append(parts, "NOT NULL")
	}

	// Handle default values
	if field.Default != nil {
		defaultVal := formatDefaultValue(field.Default, field.GoType, driver)
		parts = append(parts, fmt.Sprintf("DEFAULT %s", defaultVal))
	}

	// Handle unique constraint
	if unique, ok := field.Options["unique"].(bool); ok && unique {
		parts = append(parts, "UNIQUE")
	}

	return strings.Join(parts, " "), nil
}

// mapFieldTypeToSQL maps field types to SQL types
func mapFieldTypeToSQL(fieldType, goType string, options map[string]interface{}, driver string) string {
	isPostgres := driver == "postgres" || driver == "postgresql"
	isSQLite := driver == "sqlite" || driver == "sqlite3"

	switch fieldType {
	case "Decimal":
		maxDigits := 10
		decimalPlaces := 2
		if md, ok := options["max_digits"].(int); ok && md > 0 {
			maxDigits = md
		}
		if dp, ok := options["decimal_places"].(int); ok && dp >= 0 {
			decimalPlaces = dp
		}
		return fmt.Sprintf("NUMERIC(%d, %d)", maxDigits, decimalPlaces)

	case "JSON":
		if isPostgres {
			return "JSONB"
		}
		if isSQLite {
			return "TEXT"
		}
		return "JSON"

	case "Date":
		return "DATE"

	case "DateTime":
		if isPostgres {
			return "TIMESTAMP WITH TIME ZONE"
		}
		return "TIMESTAMP"

	case "Time":
		return "TIME"

	case "String":
		if maxLen, ok := options["max_length"].(int); ok && maxLen > 0 {
			return fmt.Sprintf("VARCHAR(%d)", maxLen)
		}
		return "TEXT"

	case "Text":
		return "TEXT"

	case "UUID":
		if isPostgres {
			return "UUID"
		}
		return "CHAR(36)"

	case "Bytes":
		if isSQLite {
			return "BLOB"
		}
		if isPostgres {
			return "BYTEA"
		}
		return "BLOB"

	case "Float64":
		if isSQLite {
			return "REAL"
		}
		return "DOUBLE PRECISION"

	case "Int64", "Int":
		if isSQLite {
			return "INTEGER"
		}
		return "BIGINT"

	case "Int32":
		return "INTEGER"

	case "Bool":
		if isSQLite {
			return "INTEGER"
		}
		return "BOOLEAN"

	case "Email", "URL":
		if maxLen, ok := options["max_length"].(int); ok && maxLen > 0 {
			return fmt.Sprintf("VARCHAR(%d)", maxLen)
		}
		return "TEXT"
	}

	// Fallback to Go type
	return mapGoTypeToSQL(goType, driver)
}

// mapGoTypeToSQL maps Go types to SQL types (fallback)
func mapGoTypeToSQL(goType string, driver string) string {
	isSQLite := driver == "sqlite" || driver == "sqlite3"

	switch goType {
	case "int64", "int":
		if isSQLite {
			return "INTEGER"
		}
		return "BIGINT"
	case "int32":
		return "INTEGER"
	case "string":
		return "TEXT"
	case "bool":
		if isSQLite {
			return "INTEGER"
		}
		return "BOOLEAN"
	case "time.Time", "*time.Time":
		return "TIMESTAMP"
	case "float64":
		if isSQLite {
			return "REAL"
		}
		return "DOUBLE PRECISION"
	case "float32":
		return "REAL"
	case "[]byte":
		if isSQLite {
			return "BLOB"
		}
		return "BYTEA"
	default:
		return "TEXT"
	}
}

// formatDefaultValue formats default value for SQL
func formatDefaultValue(value interface{}, goType string, driver string) string {
	isSQLite := driver == "sqlite" || driver == "sqlite3"

	switch v := value.(type) {
	case string:
		return fmt.Sprintf("'%s'", strings.ReplaceAll(v, "'", "''"))
	case bool:
		if isSQLite {
			if v {
				return "1"
			}
			return "0"
		}
		if v {
			return "TRUE"
		}
		return "FALSE"
	case int, int32, int64:
		return fmt.Sprintf("%d", v)
	case float32, float64:
		return fmt.Sprintf("%f", v)
	default:
		return fmt.Sprintf("'%v'", v)
	}
}

// verifySQL verifies that generated SQL is correct
func verifySQL(t *testing.T, def *generator.ModelDefinition, sql string) {
	// Check that SQL contains CREATE TABLE
	if !strings.Contains(sql, "CREATE TABLE") {
		t.Errorf("SQL for %s does not contain CREATE TABLE", def.Name)
	}

	// Check that table name is present
	tableName := def.Meta.TableName
	if tableName == "" {
		tableName = strings.ToLower(def.Name) + "s"
	}
	if !strings.Contains(sql, tableName) {
		t.Errorf("SQL for %s does not contain table name %s", def.Name, tableName)
	}

	// Check that primary key field is present
	for _, field := range def.Fields {
		if field.PrimaryKey {
			if !strings.Contains(sql, field.Name) {
				t.Errorf("SQL for %s does not contain primary key field %s", def.Name, field.Name)
			}
			if !strings.Contains(sql, "PRIMARY KEY") {
				t.Errorf("SQL for %s does not contain PRIMARY KEY constraint", def.Name)
			}
			break
		}
	}

	// Check that required fields have NOT NULL
	for _, field := range def.Fields {
		if field.Required && !field.PrimaryKey {
			// Check if field name appears with NOT NULL nearby
			fieldIndex := strings.Index(sql, fmt.Sprintf(`"%s"`, field.Name))
			if fieldIndex >= 0 {
				fieldContext := sql[fieldIndex:min(fieldIndex+200, len(sql))]
				if !strings.Contains(fieldContext, "NOT NULL") {
					t.Logf("Warning: Required field %s.%s may not have NOT NULL constraint", def.Name, field.Name)
				}
			}
		}
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// TestGenerateSQLForAllModels generates SQL for all models and writes to a file for inspection
func TestGenerateSQLForAllModels(t *testing.T) {
	// Generate code first using CLI tool
	setupTest(t)

	parser := generator.NewASTParser()
	definitions, err := parser.ParseDirectory(".")
	if err != nil {
		t.Fatalf("Failed to parse models: %v", err)
	}

	// Create output directory
	outputDir := filepath.Join("..", "generated_sql")
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		t.Fatalf("Failed to create output directory: %v", err)
	}

	// Generate SQL for each model
	for _, def := range definitions {
		sql, err := generateSQLForModel(def, "postgres")
		if err != nil {
			t.Errorf("Failed to generate SQL for %s: %v", def.Name, err)
			continue
		}

		// Write to file
		filename := filepath.Join(outputDir, fmt.Sprintf("%s.sql", strings.ToLower(def.Name)))
		if err := os.WriteFile(filename, []byte(sql), 0644); err != nil {
			t.Errorf("Failed to write SQL file for %s: %v", def.Name, err)
			continue
		}

		t.Logf("✅ Generated SQL for %s -> %s", def.Name, filename)
	}

	t.Logf("✅ Generated SQL files for all %d models in %s", len(definitions), outputDir)
}
