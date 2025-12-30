package parse

import (
	"strings"

	"github.com/forgego/forge/db/migrate/core"
)

// SQLParser parses SQL migration files to extract schema changes
// Uses lexer-based statement splitting and limited DDL parsing
type SQLParser struct {
	lexer      *Lexer
	classifier *Classifier
	ddlParser  *DDLParser
}

// NewSQLParser creates a new SQL parser using lexer-based approach
func NewSQLParser() *SQLParser {
	return &SQLParser{
		lexer:      NewLexer(),
		classifier: NewClassifier(),
		ddlParser:  NewDDLParser(),
	}
}

// ParseUpSQL parses up migration SQL and extracts changes
// Uses lexer-based statement splitting and limited DDL parsing
// Returns UnknownChange for unparseable statements (fail-soft)
func (p *SQLParser) ParseUpSQL(sql string) ([]core.Change, error) {
	// 1. Split SQL into statements using lexer
	statements, err := p.lexer.Scan(sql)
	if err != nil {
		// If lexer fails, try to continue with best-effort
		// This shouldn't happen often, but fail softly
		return []core.Change{&core.UnknownChange{SQL: sql}}, nil
	}

	// 2. Parse DDL statements only
	var changes []core.Change
	for _, stmt := range statements {
		// Skip empty statements
		if strings.TrimSpace(stmt.Text) == "" {
			continue
		}

		// Parse statement using DDL parser
		parsed, err := p.ddlParser.ParseStatement(stmt)
		if err != nil {
			// Fail softly: add as UnknownChange
			// TODO: Add optional logging for UnknownChange statements to help identify parsing gaps
			changes = append(changes, &core.UnknownChange{SQL: stmt.Text})
			continue
		}

		changes = append(changes, parsed...)
	}

	return changes, nil
}
