package sqlgen

import (
	"github.com/forgego/forge/pkg/migrate/errors"
	"github.com/forgego/forge/pkg/migrate/types"
	"fmt"
	"strings"

	"github.com/forgego/forge/pkg/generator"
)

// baseBuilder contains common SQL building logic
type baseBuilder struct {
	isSQLite   bool
	isPostgres bool
}

// BuildColumnDefinition builds SQL column definition from field
func (b *baseBuilder) BuildColumnDefinition(field generator.FieldDefinition) (string, error) {
	columnName := field.Name
	sqlType := mapFieldTypeToSQL(field, b.isSQLite, b.isPostgres)

	var parts []string
	parts = append(parts, fmt.Sprintf(`"%s"`, columnName))
	parts = append(parts, sqlType)

	// Handle primary key and auto-increment
	if b.isSQLite {
		if field.PrimaryKey && field.AutoIncrement {
			parts = []string{fmt.Sprintf(`"%s"`, columnName), "INTEGER", "PRIMARY KEY", "AUTOINCREMENT"}
		} else if field.PrimaryKey {
			parts = append(parts, "PRIMARY KEY")
		}
	} else {
		if field.PrimaryKey {
			parts = append(parts, "PRIMARY KEY")
		}
		if field.AutoIncrement {
			parts = append(parts, "GENERATED ALWAYS AS IDENTITY")
		}
	}

	// Handle required/optional
	if field.Required && !field.PrimaryKey {
		parts = append(parts, "NOT NULL")
	}

	// Handle default values
	// Check for AutoNowAdd or AutoNow options first (these imply DEFAULT now())
	if autoNowAdd, ok := field.Options["auto_now_add"].(bool); ok && autoNowAdd {
		parts = append(parts, "DEFAULT now()")
		// Make created_at NOT NULL when AutoNowAdd is set
		if !field.Required && !field.PrimaryKey {
			// Check if NOT NULL is already in parts
			hasNotNull := false
			for _, part := range parts {
				if part == "NOT NULL" {
					hasNotNull = true
					break
				}
			}
			if !hasNotNull {
				parts = append(parts, "NOT NULL")
			}
		}
	} else if autoNow, ok := field.Options["auto_now"].(bool); ok && autoNow {
		parts = append(parts, "DEFAULT now()")
	} else if field.Default != nil {
		defaultVal := formatDefaultValue(field.Default, field.GoType, field.Type, field.Options, b.isSQLite)
		parts = append(parts, fmt.Sprintf("DEFAULT %s", defaultVal))
	}

	// Handle unique constraint
	if unique, ok := field.Options["unique"].(bool); ok && unique {
		parts = append(parts, "UNIQUE")
	}

	return strings.Join(parts, " "), nil
}

// BuildCreateTable generates CREATE TABLE statement
func (b *baseBuilder) BuildCreateTable(c *types.CreateTable) (string, error) {
	tableName := c.TableName()
	var parts []string

	parts = append(parts, fmt.Sprintf("-- Create table: %s", tableName))
	parts = append(parts, fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (", tableName))

	var columnDefs []string
	for _, field := range c.Table.Fields {
		colDef, err := b.BuildColumnDefinition(field)
		if err != nil {
			return "", fmt.Errorf("failed to build column %s: %w", field.Name, err)
		}
		columnDefs = append(columnDefs, "    "+colDef)
	}

	parts = append(parts, strings.Join(columnDefs, ",\n"))
	parts = append(parts, ");")

	return strings.Join(parts, "\n"), nil
}

// BuildAddColumn generates ALTER TABLE ADD COLUMN statement
func (b *baseBuilder) BuildAddColumn(c *types.AddColumn) (string, error) {
	colDef, err := b.BuildColumnDefinition(c.Column)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("ALTER TABLE %s ADD COLUMN %s;", c.Table, colDef), nil
}

// BuildAddIndex generates CREATE INDEX statement
func (b *baseBuilder) BuildAddIndex(c *types.AddIndex) (string, error) {
	if len(c.Index.Fields) == 0 {
		return "", nil
	}

	escapedFields := make([]string, len(c.Index.Fields))
	for i, field := range c.Index.Fields {
		escapedFields[i] = fmt.Sprintf(`"%s"`, field)
	}

	indexType := "INDEX"
	if c.Index.Unique {
		indexType = "UNIQUE INDEX"
	}

	indexName := c.Index.Name
	if indexName == "" {
		indexName = fmt.Sprintf("idx_%s_%s", c.Table, strings.Join(c.Index.Fields, "_"))
	}

	return fmt.Sprintf("CREATE %s IF NOT EXISTS %s ON %s (%s);",
		indexType, indexName, c.Table, strings.Join(escapedFields, ", ")), nil
}

// BuildModifyIndex generates DROP and CREATE INDEX statements
func (b *baseBuilder) BuildModifyIndex(c *types.ModifyIndex) (string, error) {
	oldIndexName := c.OldIndex.Name
	if oldIndexName == "" {
		oldIndexName = fmt.Sprintf("idx_%s_%s", c.Table, strings.Join(c.OldIndex.Fields, "_"))
	}

	var statements []string
	statements = append(statements, fmt.Sprintf("DROP INDEX IF EXISTS %s;", oldIndexName))

	newIndex := &types.AddIndex{
		Table: c.Table,
		Index: c.NewIndex,
	}
	newIndexSQL, err := b.BuildAddIndex(newIndex)
	if err != nil {
		return "", err
	}
	statements = append(statements, newIndexSQL)

	return strings.Join(statements, "\n"), nil
}

// BuildAddForeignKey generates ALTER TABLE ADD CONSTRAINT FOREIGN KEY statement
func (b *baseBuilder) BuildAddForeignKey(c *types.AddForeignKey) (string, error) {
	// Validate target table is not empty
	if c.TargetTable == "" {
		return "", errors.NewMigrationError(
			errors.ErrInvalidChange,
			fmt.Sprintf("foreign key %s.%s has empty target table", c.Table, c.Relation.Name),
			nil,
		)
	}

	onDelete := "NO ACTION"
	if onDeleteVal, ok := c.Relation.Options["on_delete"].(string); ok {
		onDelete = mapCascadeType(onDeleteVal)
	}

	onUpdate := "NO ACTION"
	if onUpdateVal, ok := c.Relation.Options["on_update"].(string); ok {
		onUpdate = mapCascadeType(onUpdateVal)
	}

	fkName := fmt.Sprintf("fk_%s_%s", c.Table, c.Relation.Name)

	// Use DO block to check if constraint exists before adding (PostgreSQL-specific)
	// This prevents errors when constraint already exists (e.g., in incremental migrations)
	return fmt.Sprintf(`DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_constraint 
        WHERE conname = '%s' 
        AND conrelid = '%s'::regclass
    ) THEN
        ALTER TABLE %s ADD CONSTRAINT %s FOREIGN KEY ("%s") REFERENCES %s (id) ON DELETE %s ON UPDATE %s;
    END IF;
END $$;`,
		fkName, c.Table, c.Table, fkName, c.Relation.Name, c.TargetTable, onDelete, onUpdate), nil
}

// BuildModifyForeignKey generates DROP and ADD FOREIGN KEY statements
func (b *baseBuilder) BuildModifyForeignKey(c *types.ModifyForeignKey) (string, error) {
	oldFKName := fmt.Sprintf("fk_%s_%s", c.Table, c.OldFK.Name)

	var statements []string
	statements = append(statements, fmt.Sprintf("ALTER TABLE %s DROP CONSTRAINT IF EXISTS %s;", c.Table, oldFKName))

	newFK := &types.AddForeignKey{
		Table:       c.Table,
		Relation:    c.NewFK,
		TargetTable: c.TargetTable,
	}
	newFKSQL, err := b.BuildAddForeignKey(newFK)
	if err != nil {
		return "", err
	}
	statements = append(statements, newFKSQL)

	return strings.Join(statements, "\n"), nil
}

// BuildAddConstraint generates ALTER TABLE ADD CONSTRAINT statement
func (b *baseBuilder) BuildAddConstraint(c *types.AddConstraint) (string, error) {
	var constraintSQL string
	switch strings.ToUpper(c.Constraint.Type) {
	case "CHECK":
		if c.Constraint.Condition != "" {
			constraintSQL = fmt.Sprintf("CHECK (%s)", c.Constraint.Condition)
		} else if len(c.Constraint.Fields) > 0 {
			escapedFields := make([]string, len(c.Constraint.Fields))
			for i, field := range c.Constraint.Fields {
				escapedFields[i] = fmt.Sprintf(`"%s"`, field)
			}
			constraintSQL = fmt.Sprintf("CHECK (%s)", strings.Join(escapedFields, ", "))
		} else {
			return "", errors.NewMigrationError(
				errors.ErrInvalidChange,
				"CHECK constraint requires condition or fields",
				nil,
			)
		}
	case "UNIQUE":
		if len(c.Constraint.Fields) > 0 {
			escapedFields := make([]string, len(c.Constraint.Fields))
			for i, field := range c.Constraint.Fields {
				escapedFields[i] = fmt.Sprintf(`"%s"`, field)
			}
			constraintSQL = fmt.Sprintf("UNIQUE (%s)", strings.Join(escapedFields, ", "))
		} else {
			return "", errors.NewMigrationError(
				errors.ErrInvalidChange,
				"UNIQUE constraint requires fields",
				nil,
			)
		}
	default:
		if c.Constraint.Condition != "" {
			constraintSQL = c.Constraint.Condition
		} else {
			return "", errors.NewMigrationError(
				errors.ErrInvalidChange,
				fmt.Sprintf("constraint type %s requires condition", c.Constraint.Type),
				nil,
			)
		}
	}

	return fmt.Sprintf("ALTER TABLE %s ADD CONSTRAINT %s %s;",
		c.Table, c.Constraint.Name, constraintSQL), nil
}

// orderChanges orders changes by dependency
func orderChanges(changes []types.Change) []types.Change {
	ordered := make([]types.Change, 0, len(changes))
	
	// First pass: CreateTable
	for _, change := range changes {
		if change.Type() == types.ChangeTypeCreateTable {
			ordered = append(ordered, change)
		}
	}
	
	// Second pass: AddColumn, RenameColumn, ModifyColumn
	for _, change := range changes {
		ct := change.Type()
		if ct == types.ChangeTypeAddColumn || ct == types.ChangeTypeRenameColumn || ct == types.ChangeTypeModifyColumn {
			ordered = append(ordered, change)
		}
	}
	
	// Third pass: AddForeignKey, ModifyForeignKey
	for _, change := range changes {
		if change.Type() == types.ChangeTypeAddForeignKey || change.Type() == types.ChangeTypeModifyForeignKey {
			ordered = append(ordered, change)
		}
	}
	
	// Fourth pass: AddIndex, ModifyIndex, AddConstraint
	for _, change := range changes {
		ct := change.Type()
		if ct == types.ChangeTypeAddIndex || ct == types.ChangeTypeModifyIndex || ct == types.ChangeTypeAddConstraint {
			ordered = append(ordered, change)
		}
	}
	
	// Fifth pass: Everything else
	for _, change := range changes {
		ct := change.Type()
		if ct != types.ChangeTypeCreateTable && ct != types.ChangeTypeAddColumn && ct != types.ChangeTypeRenameColumn &&
			ct != types.ChangeTypeModifyColumn && ct != types.ChangeTypeAddForeignKey && ct != types.ChangeTypeModifyForeignKey &&
			ct != types.ChangeTypeAddIndex && ct != types.ChangeTypeModifyIndex && ct != types.ChangeTypeAddConstraint {
			ordered = append(ordered, change)
		}
	}
	
	return ordered
}

// mapCascadeType maps cascade type strings to SQL
func mapCascadeType(cascade string) string {
	switch cascade {
	case "CASCADE":
		return "CASCADE"
	case "SET_NULL", "SET NULL":
		return "SET NULL"
	case "PROTECT":
		return "RESTRICT"
	case "SET_DEFAULT", "SET DEFAULT":
		return "SET DEFAULT"
	case "DO_NOTHING", "NO ACTION":
		return "NO ACTION"
	default:
		return "NO ACTION"
	}
}

// formatDefaultValue formats default value for SQL
func formatDefaultValue(value interface{}, goType string, fieldType string, fieldOptions map[string]interface{}, isSQLite bool) string {
	// Check if value is a SQL expression (string that looks like a function call)
	if str, ok := value.(string); ok {
		// Common SQL functions that should not be quoted
		sqlFunctions := []string{"now()", "CURRENT_TIMESTAMP", "CURRENT_DATE", "CURRENT_TIME",
			"uuid_generate_v4()", "gen_random_uuid()", "random()"}
		for _, fn := range sqlFunctions {
			if strings.EqualFold(str, fn) {
				return str // Return unquoted SQL expression
			}
		}
		// Check if it looks like a function call (contains parentheses)
		if strings.Contains(str, "(") && strings.Contains(str, ")") {
			return str // Likely a SQL function, return unquoted
		}
		// Regular string, quote it
		return fmt.Sprintf("'%s'", strings.ReplaceAll(str, "'", "''"))
	}

	switch v := value.(type) {
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
		// For Decimal fields, format to match decimal places
		if fieldType == "Decimal" {
			decimalPlaces := 2
			if dp, ok := fieldOptions["decimal_places"].(int); ok && dp >= 0 {
				decimalPlaces = dp
			}
			return fmt.Sprintf("%.*f", decimalPlaces, v)
		}
		return fmt.Sprintf("%f", v)
	default:
		return fmt.Sprintf("'%v'", v)
	}
}

// mapFieldTypeToSQL maps field types to SQL types
func mapFieldTypeToSQL(field generator.FieldDefinition, isSQLite, isPostgres bool) string {
	switch field.Type {
	case "Decimal":
		maxDigits := 10
		decimalPlaces := 2
		if md, ok := field.Options["max_digits"].(int); ok && md > 0 {
			maxDigits = md
		}
		if dp, ok := field.Options["decimal_places"].(int); ok && dp >= 0 {
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
		if maxLen, ok := field.Options["max_length"].(int); ok && maxLen > 0 {
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
		if maxLen, ok := field.Options["max_length"].(int); ok && maxLen > 0 {
			return fmt.Sprintf("VARCHAR(%d)", maxLen)
		}
		return "TEXT"
	}

	// Fallback to Go type
	return mapGoTypeToSQL(field.GoType, isSQLite, isPostgres)
}

// mapGoTypeToSQL maps Go types to SQL types (fallback)
func mapGoTypeToSQL(goType string, isSQLite, isPostgres bool) string {
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
