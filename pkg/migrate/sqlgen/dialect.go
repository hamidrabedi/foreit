package sqlgen

// Dialect represents a SQL database dialect
type Dialect int

const (
	DialectPostgres Dialect = iota
	DialectSQLite
	DialectMySQL
)

// String returns the string representation of the dialect
func (d Dialect) String() string {
	switch d {
	case DialectPostgres:
		return "postgres"
	case DialectSQLite:
		return "sqlite"
	case DialectMySQL:
		return "mysql"
	default:
		return "unknown"
	}
}

// SupportsDollarQuotes returns true if the dialect supports dollar-quoted strings
func (d Dialect) SupportsDollarQuotes() bool {
	return d == DialectPostgres
}

// SupportsHashComments returns true if the dialect supports hash comments
func (d Dialect) SupportsHashComments() bool {
	return d == DialectMySQL
}

// TODO: Add more dialect-specific features:
// - MySQL: AUTO_INCREMENT syntax differences
// - SQLite: Limited ALTER TABLE support
// - PostgreSQL: Advanced features (arrays, JSONB, custom types)
// - Dialect-specific type mappings
