package state

import (
	"regexp"
	"strings"

	"github.com/forgego/forge/pkg/generator"
	"github.com/forgego/forge/pkg/migrations/core"
)

// SQLParser parses SQL migration files to extract schema changes
type SQLParser struct {
	// Regular expressions for parsing SQL
	createTableRegex *regexp.Regexp
	alterTableRegex  *regexp.Regexp
	createIndexRegex *regexp.Regexp
}

// NewSQLParser creates a new SQL parser
func NewSQLParser() *SQLParser {
	return &SQLParser{
		createTableRegex: regexp.MustCompile(`(?i)CREATE\s+TABLE\s+(?:IF\s+NOT\s+EXISTS\s+)?["']?(\w+)["']?`),
		alterTableRegex:  regexp.MustCompile(`(?i)ALTER\s+TABLE\s+["']?(\w+)["']?`),
		createIndexRegex: regexp.MustCompile(`(?i)CREATE\s+(?:UNIQUE\s+)?INDEX\s+(?:IF\s+NOT\s+EXISTS\s+)?(\w+)\s+ON\s+["']?(\w+)["']?`),
	}
}

// ParseUpSQL parses up migration SQL and extracts changes
// This is a simplified parser - a full implementation would need to parse
// the complete SQL structure including columns, types, constraints, etc.
func (p *SQLParser) ParseUpSQL(sql string) ([]core.Change, error) {
	var changes []core.Change

	// Split SQL into statements
	statements := splitSQLStatements(sql)

	for _, stmt := range statements {
		stmt = strings.TrimSpace(stmt)
		if stmt == "" || strings.HasPrefix(stmt, "--") {
			continue
		}

		// Try to match CREATE TABLE
		if matches := p.createTableRegex.FindStringSubmatch(stmt); len(matches) > 0 {
			tableParser := NewTableParser()
			createTable, err := tableParser.ParseCreateTable(stmt)
			if err == nil {
				changes = append(changes, createTable)
			}
		}

		// Try to match ALTER TABLE ADD COLUMN
		if strings.Contains(strings.ToUpper(stmt), "ALTER TABLE") && strings.Contains(strings.ToUpper(stmt), "ADD COLUMN") {
			// Parse column addition
			// This is complex and would require full SQL parsing
		}

		// Try to match CREATE INDEX
		if matches := p.createIndexRegex.FindStringSubmatch(stmt); len(matches) > 0 {
			indexName := matches[1]
			tableName := matches[2]
			isUnique := strings.Contains(strings.ToUpper(stmt), "UNIQUE")

			// Extract fields from index definition
			// This is simplified - a full parser would extract the field list
			fields := extractIndexFields(stmt)

			changes = append(changes, &core.AddIndex{
				Table: tableName,
				Index: generator.IndexDefinition{
					Name:   indexName,
					Fields: fields,
					Unique: isUnique,
				},
			})
		}
	}

	return changes, nil
}

// splitSQLStatements splits SQL into individual statements
func splitSQLStatements(sql string) []string {
	// Simple split by semicolon
	// A proper implementation would handle quoted strings, comments, etc.
	statements := strings.Split(sql, ";")
	var result []string
	for _, stmt := range statements {
		stmt = strings.TrimSpace(stmt)
		if stmt != "" {
			result = append(result, stmt)
		}
	}
	return result
}

// extractIndexFields extracts field names from an index definition
func extractIndexFields(sql string) []string {
	// This is a simplified extraction
	// Look for patterns like (field1, field2, ...)
	re := regexp.MustCompile(`\(([^)]+)\)`)
	matches := re.FindStringSubmatch(sql)
	if len(matches) > 1 {
		fieldsStr := matches[1]
		fields := strings.Split(fieldsStr, ",")
		var result []string
		for _, field := range fields {
			field = strings.TrimSpace(field)
			// Remove quotes
			field = strings.Trim(field, `"`)
			field = strings.Trim(field, `'`)
			if field != "" {
				result = append(result, field)
			}
		}
		return result
	}
	return []string{}
}

// Note: Full SQL parsing is complex and would require a proper SQL parser library
// For now, this provides a basic structure that can be extended
// The current implementation focuses on CREATE INDEX statements as an example
// A production implementation would need to parse:
// - CREATE TABLE with full column definitions
// - ALTER TABLE ADD/DROP/MODIFY COLUMN
// - CREATE/DROP INDEX
// - ALTER TABLE ADD/DROP CONSTRAINT
// - ALTER TABLE ADD/DROP FOREIGN KEY

