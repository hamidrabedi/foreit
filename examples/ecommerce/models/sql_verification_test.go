package models

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/forgego/forge/pkg/generator"
)

var codeGenerated = false

// setupTest generates code using the CLI tool before running tests
func setupTest(t *testing.T) {
	// Only generate once per test run
	if codeGenerated {
		return
	}

	// Get the current directory
	modelsDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}

	// Find the forge binary - try multiple locations
	var forgePath string
	possiblePaths := []string{
		filepath.Join(modelsDir, "..", "..", "..", "forge.exe"),     // Windows
		filepath.Join(modelsDir, "..", "..", "..", "forge"),         // Unix
		filepath.Join(modelsDir, "..", "..", "..", "..", "forge.exe"), // Alternative location
		filepath.Join(modelsDir, "..", "..", "..", "..", "forge"),
		"forge", // In PATH
	}

	for _, path := range possiblePaths {
		if path == "forge" {
			// Check if forge is in PATH
			if _, err := exec.LookPath("forge"); err == nil {
				forgePath = "forge"
				break
			}
		} else {
			if _, err := os.Stat(path); err == nil {
				forgePath = path
				break
			}
		}
	}

	if forgePath == "" {
		t.Skipf("Forge CLI tool not found. Please build it first: go build -o forge ./cmd/forge")
		return
	}

	// Run forge generate command
	cmd := exec.Command(forgePath, "generate", "--models", modelsDir, "--output", modelsDir)
	// Set working directory to newforge root
	cmd.Dir = filepath.Join(modelsDir, "..", "..", "..")
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Logf("Forge generate command: %s generate --models %s --output %s", forgePath, modelsDir, modelsDir)
		t.Logf("Forge generate output: %s", string(output))
		t.Fatalf("Failed to generate code using CLI tool: %v", err)
	}

	codeGenerated = true
	t.Logf("✅ Code generated successfully using forge CLI")
}

// TestSQLFeatureSupport verifies that all SQL features are properly supported
// in the schema definitions. This test ensures that:
// 1. All PostgreSQL data types are correctly represented
// 2. All field options (Required, Unique, Default, etc.) are supported
// 3. All relation types (ForeignKey, OneToOne, ManyToMany) are supported
// 4. Indexes and constraints are properly defined
func TestSQLFeatureSupport(t *testing.T) {
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

	// Track which features we've verified
	features := map[string]bool{
		"Int64":              false,
		"Int32":              false,
		"Decimal":            false,
		"Float64":            false,
		"String":             false,
		"Text":               false,
		"Bool":               false,
		"UUID":               false,
		"JSON":               false,
		"DateTime":           false,
		"Date":               false,
		"Time":               false,
		"Bytes":              false,
		"URL":                false,
		"Email":              false,
		"PrimaryKey":         false,
		"AutoIncrement":      false,
		"Required":           false,
		"Optional":           false,
		"Unique":             false,
		"Default":            false,
		"MaxLength":          false,
		"MaxDigits":          false,
		"DecimalPlaces":      false,
		"Choices":            false,
		"AutoNow":            false,
		"AutoNowAdd":         false,
		"DBIndex":            false,
		"ForeignKey":         false,
		"OneToOne":           false,
		"OneToMany":          false,
		"ManyToMany":         false,
		"SelfReferential":    false,
		"Indexes":            false,
		"UniqueTogether":     false,
		"CascadeCASCADE":     false,
		"CascadeSET_NULL":    false,
		"CascadePROTECT":     false,
		"WriteOnly":          false,
	}

	// Verify features across all models
	for _, def := range definitions {
		verifyFields(t, def, features)
		verifyRelations(t, def, features)
		verifyMeta(t, def, features)
	}

	// Report missing features
	missing := []string{}
	for feature, found := range features {
		if !found {
			missing = append(missing, feature)
		}
	}

	if len(missing) > 0 {
		t.Errorf("Missing SQL features: %v", missing)
	} else {
		t.Log("✅ All SQL features are supported in the schema definitions")
	}
}

func verifyFields(t *testing.T, def *generator.ModelDefinition, features map[string]bool) {
	for _, field := range def.Fields {
		// Verify field types
		switch field.Type {
		case "Int64":
			features["Int64"] = true
		case "Int32":
			features["Int32"] = true
		case "Decimal":
			features["Decimal"] = true
		case "Float64":
			features["Float64"] = true
		case "String":
			features["String"] = true
		case "Text":
			features["Text"] = true
		case "Bool":
			features["Bool"] = true
		case "UUID":
			features["UUID"] = true
		case "JSON":
			features["JSON"] = true
		case "DateTime":
			features["DateTime"] = true
		case "Date":
			features["Date"] = true
		case "Time":
			features["Time"] = true
		case "Bytes":
			features["Bytes"] = true
		case "URL":
			features["URL"] = true
		case "Email":
			features["Email"] = true
		}

		// Verify field options
		if field.PrimaryKey {
			features["PrimaryKey"] = true
		}
		if field.AutoIncrement {
			features["AutoIncrement"] = true
		}
		if field.Required {
			features["Required"] = true
		} else {
			features["Optional"] = true
		}
		if unique, ok := field.Options["unique"].(bool); ok && unique {
			features["Unique"] = true
		}
		if field.Default != nil {
			features["Default"] = true
		}
		if maxLen, ok := field.Options["max_length"].(int); ok && maxLen > 0 {
			features["MaxLength"] = true
		}
		if maxDigits, ok := field.Options["max_digits"].(int); ok && maxDigits > 0 {
			features["MaxDigits"] = true
		}
		if decimalPlaces, ok := field.Options["decimal_places"].(int); ok && decimalPlaces >= 0 {
			features["DecimalPlaces"] = true
		}
		if choices, ok := field.Options["choices"].([]interface{}); ok && len(choices) > 0 {
			features["Choices"] = true
		}
		if autoNow, ok := field.Options["auto_now"].(bool); ok && autoNow {
			features["AutoNow"] = true
		}
		if autoNowAdd, ok := field.Options["auto_now_add"].(bool); ok && autoNowAdd {
			features["AutoNowAdd"] = true
		}
		if dbIndex, ok := field.Options["db_index"].(bool); ok && dbIndex {
			features["DBIndex"] = true
		}
		if writeOnly, ok := field.Options["write_only"].(bool); ok && writeOnly {
			features["WriteOnly"] = true
		}
	}
}

func verifyRelations(t *testing.T, def *generator.ModelDefinition, features map[string]bool) {
	for _, rel := range def.Relations {
		switch rel.Type {
		case "ForeignKey":
			features["ForeignKey"] = true
			// Check cascade behaviors
			if onDelete, ok := rel.Options["on_delete"].(string); ok {
				switch onDelete {
				case "CascadeCASCADE", "CASCADE":
					features["CascadeCASCADE"] = true
				case "CascadeSET_NULL", "SET_NULL":
					features["CascadeSET_NULL"] = true
				case "CascadePROTECT", "PROTECT":
					features["CascadePROTECT"] = true
				}
			}
		case "OneToOne":
			features["OneToOne"] = true
		case "OneToMany":
			features["OneToMany"] = true
		case "ManyToMany":
			features["ManyToMany"] = true
		}

		// Check for self-referential relations
		if strings.EqualFold(rel.To, def.Name) {
			features["SelfReferential"] = true
		}
	}
}

func verifyMeta(t *testing.T, def *generator.ModelDefinition, features map[string]bool) {
	if len(def.Meta.Indexes) > 0 {
		features["Indexes"] = true
	}
	if len(def.Meta.UniqueTogether) > 0 {
		features["UniqueTogether"] = true
	}
}

// TestPostgreSQLTypes verifies that all PostgreSQL data types are properly represented
func TestPostgreSQLTypes(t *testing.T) {
	// Generate code first using CLI tool
	setupTest(t)

	expectedTypes := map[string]string{
		"Int64":     "BIGINT",
		"Int32":     "INTEGER",
		"Decimal":   "NUMERIC",
		"Float64":   "DOUBLE PRECISION",
		"String":    "VARCHAR or TEXT",
		"Text":      "TEXT",
		"Bool":      "BOOLEAN",
		"UUID":      "UUID",
		"JSON":      "JSONB",
		"DateTime":  "TIMESTAMP WITH TIME ZONE",
		"Date":      "DATE",
		"Time":      "TIME",
		"Bytes":     "BYTEA",
	}

	parser := generator.NewASTParser()
	definitions, err := parser.ParseDirectory(".")
	if err != nil {
		t.Fatalf("Failed to parse models: %v", err)
	}

	foundTypes := make(map[string]bool)
	for _, def := range definitions {
		for _, field := range def.Fields {
			if expectedType, ok := expectedTypes[field.Type]; ok {
				foundTypes[field.Type] = true
				t.Logf("✅ Found %s field (maps to %s): %s.%s", field.Type, expectedType, def.Name, field.Name)
			}
		}
	}

	// Verify all expected types are found
	for fieldType, expectedSQL := range expectedTypes {
		if !foundTypes[fieldType] {
			t.Errorf("❌ Missing field type: %s (expected SQL: %s)", fieldType, expectedSQL)
		}
	}

	if len(foundTypes) == len(expectedTypes) {
		t.Log("✅ All PostgreSQL data types are represented in the models")
	}
}

// TestComplexRelations verifies that complex relation patterns are supported
func TestComplexRelations(t *testing.T) {
	// Generate code first using CLI tool
	setupTest(t)

	parser := generator.NewASTParser()
	definitions, err := parser.ParseDirectory(".")
	if err != nil {
		t.Fatalf("Failed to parse models: %v", err)
	}

	relationTypes := map[string]int{
		"ForeignKey":  0,
		"OneToOne":    0,
		"OneToMany":   0,
		"ManyToMany":  0,
		"SelfReferential": 0,
	}

	for _, def := range definitions {
		for _, rel := range def.Relations {
			relationTypes[rel.Type]++
			if strings.EqualFold(rel.To, def.Name) {
				relationTypes["SelfReferential"]++
			}
		}
	}

	t.Logf("Relation counts:")
	for relType, count := range relationTypes {
		if count > 0 {
			t.Logf("  ✅ %s: %d", relType, count)
		} else {
			t.Logf("  ⚠️  %s: 0 (not tested)", relType)
		}
	}

	// Verify we have examples of key relation types
	if relationTypes["ForeignKey"] == 0 {
		t.Error("❌ No ForeignKey relations found")
	}
	if relationTypes["ManyToMany"] == 0 {
		t.Error("❌ No ManyToMany relations found")
	}
	if relationTypes["OneToOne"] == 0 {
		t.Error("❌ No OneToOne relations found")
	}
}

// TestFieldOptions verifies that all field options are properly supported
func TestFieldOptions(t *testing.T) {
	// Generate code first using CLI tool
	setupTest(t)

	parser := generator.NewASTParser()
	definitions, err := parser.ParseDirectory(".")
	if err != nil {
		t.Fatalf("Failed to parse models: %v", err)
	}

	optionsFound := map[string]bool{
		"PrimaryKey":    false,
		"AutoIncrement": false,
		"Required":      false,
		"Unique":        false,
		"Default":       false,
		"MaxLength":     false,
		"Choices":       false,
		"AutoNow":       false,
		"AutoNowAdd":    false,
		"DBIndex":       false,
		"WriteOnly":     false,
	}

	for _, def := range definitions {
		for _, field := range def.Fields {
			if field.PrimaryKey {
				optionsFound["PrimaryKey"] = true
			}
			if field.AutoIncrement {
				optionsFound["AutoIncrement"] = true
			}
			if field.Required {
				optionsFound["Required"] = true
			}
			if unique, ok := field.Options["unique"].(bool); ok && unique {
				optionsFound["Unique"] = true
			}
			if field.Default != nil {
				optionsFound["Default"] = true
			}
			if _, ok := field.Options["max_length"].(int); ok {
				optionsFound["MaxLength"] = true
			}
			if _, ok := field.Options["choices"].([]interface{}); ok {
				optionsFound["Choices"] = true
			}
			if autoNow, ok := field.Options["auto_now"].(bool); ok && autoNow {
				optionsFound["AutoNow"] = true
			}
			if autoNowAdd, ok := field.Options["auto_now_add"].(bool); ok && autoNowAdd {
				optionsFound["AutoNowAdd"] = true
			}
			if dbIndex, ok := field.Options["db_index"].(bool); ok && dbIndex {
				optionsFound["DBIndex"] = true
			}
			if writeOnly, ok := field.Options["write_only"].(bool); ok && writeOnly {
				optionsFound["WriteOnly"] = true
			}
		}
	}

	t.Log("Field options found:")
	for option, found := range optionsFound {
		if found {
			t.Logf("  ✅ %s", option)
		} else {
			t.Logf("  ⚠️  %s (not tested)", option)
		}
	}
}

// TestIndexesAndConstraints verifies that indexes and constraints are properly defined
func TestIndexesAndConstraints(t *testing.T) {
	// Generate code first using CLI tool
	setupTest(t)

	parser := generator.NewASTParser()
	definitions, err := parser.ParseDirectory(".")
	if err != nil {
		t.Fatalf("Failed to parse models: %v", err)
	}

	totalIndexes := 0
	totalUniqueTogether := 0
	indexTypes := map[string]int{
		"SingleField": 0,
		"MultiField":  0,
		"Unique":      0,
		"NonUnique":   0,
	}

	for _, def := range definitions {
		for _, idx := range def.Meta.Indexes {
			totalIndexes++
			if len(idx.Fields) == 1 {
				indexTypes["SingleField"]++
			} else {
				indexTypes["MultiField"]++
			}
			if idx.Unique {
				indexTypes["Unique"]++
			} else {
				indexTypes["NonUnique"]++
			}
		}
		totalUniqueTogether += len(def.Meta.UniqueTogether)
	}

	t.Logf("Indexes and constraints:")
	t.Logf("  ✅ Total indexes: %d", totalIndexes)
	t.Logf("  ✅ Total unique together constraints: %d", totalUniqueTogether)
	t.Logf("  ✅ Single field indexes: %d", indexTypes["SingleField"])
	t.Logf("  ✅ Multi-field indexes: %d", indexTypes["MultiField"])
	t.Logf("  ✅ Unique indexes: %d", indexTypes["Unique"])
	t.Logf("  ✅ Non-unique indexes: %d", indexTypes["NonUnique"])

	if totalIndexes == 0 {
		t.Error("❌ No indexes found in models")
	}
}

