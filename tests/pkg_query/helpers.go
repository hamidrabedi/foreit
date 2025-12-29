package query

import (
	"database/sql"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// QueryCapture captures SQL queries executed
type QueryCapture struct {
	Queries []string
	db      *sql.DB
}

// NewQueryCapture creates a new query capture
func NewQueryCapture(db *sql.DB) *QueryCapture {
	return &QueryCapture{
		Queries: make([]string, 0),
		db:      db,
	}
}

// CaptureQueries captures all SQL queries executed in a function
func CaptureQueries(t *testing.T, db *sql.DB, fn func()) []string {
	// This is a simplified version - in production, you'd use a query logger
	// or database driver that supports query interception
	capture := NewQueryCapture(db)
	fn()
	return capture.Queries
}

// AssertSQLContains checks if generated SQL contains expected parts
func AssertSQLContains(t *testing.T, sql string, expected string) {
	assert.Contains(t, strings.ToLower(sql), strings.ToLower(expected),
		"Expected SQL to contain '%s', but got: %s", expected, sql)
}

// AssertSQLNotContains checks if generated SQL does not contain parts
func AssertSQLNotContains(t *testing.T, sql string, notExpected string) {
	assert.NotContains(t, strings.ToLower(sql), strings.ToLower(notExpected),
		"Expected SQL to not contain '%s', but got: %s", notExpected, sql)
}

// AssertQueryEqual compares two queries for equivalence
// This is a simplified version - full implementation would compare SQL and parameters
func AssertQueryEqual(t *testing.T, expected, actual interface{}, msgAndArgs ...interface{}) {
	assert.Equal(t, expected, actual, msgAndArgs...)
}

// RequireNoError requires no error or fails test
func RequireNoError(t *testing.T, err error, msgAndArgs ...interface{}) {
	require.NoError(t, err, msgAndArgs...)
}

// AssertCount asserts the count of records
func AssertCount(t *testing.T, db *sql.DB, table string, expected int) {
	var count int
	err := db.QueryRow(fmt.Sprintf("SELECT COUNT(*) FROM %s", table)).Scan(&count)
	require.NoError(t, err)
	assert.Equal(t, expected, count, "Expected %d records in %s, got %d", expected, table, count)
}

// AssertRecordExists asserts a record exists with given conditions
func AssertRecordExists(t *testing.T, db *sql.DB, table string, conditions map[string]interface{}) {
	var count int
	whereParts := make([]string, 0, len(conditions))
	args := make([]interface{}, 0, len(conditions))
	argIndex := 1

	for col, val := range conditions {
		whereParts = append(whereParts, fmt.Sprintf("%s = $%d", col, argIndex))
		args = append(args, val)
		argIndex++
	}

	query := fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE %s", table, strings.Join(whereParts, " AND "))
	err := db.QueryRow(query, args...).Scan(&count)
	require.NoError(t, err)
	assert.Greater(t, count, 0, "Expected record to exist in %s with conditions %v", table, conditions)
}

// AssertRecordNotExists asserts a record does not exist
func AssertRecordNotExists(t *testing.T, db *sql.DB, table string, conditions map[string]interface{}) {
	var count int
	whereParts := make([]string, 0, len(conditions))
	args := make([]interface{}, 0, len(conditions))
	argIndex := 1

	for col, val := range conditions {
		whereParts = append(whereParts, fmt.Sprintf("%s = $%d", col, argIndex))
		args = append(args, val)
		argIndex++
	}

	query := fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE %s", table, strings.Join(whereParts, " AND "))
	err := db.QueryRow(query, args...).Scan(&count)
	require.NoError(t, err)
	assert.Equal(t, 0, count, "Expected record to not exist in %s with conditions %v", table, conditions)
}
