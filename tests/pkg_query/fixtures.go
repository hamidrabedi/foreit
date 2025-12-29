package query

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"testing"
)

// Fixture represents test data
type Fixture struct {
	Model string
	Data  []map[string]interface{}
}

// LoadFixtures loads test data into database
func LoadFixtures(t *testing.T, db *sql.DB, fixtures []Fixture) {
	ctx := context.Background()

	for _, fixture := range fixtures {
		for _, row := range fixture.Data {
			// Build INSERT query
			var columns []string
			var placeholders []string
			var values []interface{}
			var valueIndex int

			for col, val := range row {
				columns = append(columns, col)
				valueIndex++
				placeholders = append(placeholders, fmt.Sprintf("$%d", valueIndex))
				values = append(values, val)
			}

			query := fmt.Sprintf(
				"INSERT INTO %s (%s) VALUES (%s)",
				fixture.Model,
				joinStrings(columns, ", "),
				joinStrings(placeholders, ", "),
			)

			_, err := db.ExecContext(ctx, query, values...)
			if err != nil {
				t.Fatalf("failed to insert fixture into %s: %v", fixture.Model, err)
			}
		}
	}
}

// StandardFixtures returns common test fixtures
func StandardFixtures() []Fixture {
	return []Fixture{
		{
			Model: "authors",
			Data: []map[string]interface{}{
				{"id": 1, "name": "John Doe", "email": "john@example.com", "is_active": true},
				{"id": 2, "name": "Jane Smith", "email": "jane@example.com", "is_active": true},
				{"id": 3, "name": "Bob Johnson", "email": "bob@example.com", "is_active": false},
			},
		},
		{
			Model: "books",
			Data: []map[string]interface{}{
				{"id": 1, "title": "Go Programming", "isbn": "978-0-123456-78-9", "author_id": 1, "price": 29.99, "available": true, "pages": 300},
				{"id": 2, "title": "Advanced Go", "isbn": "978-0-123456-79-6", "author_id": 1, "price": 39.99, "available": true, "pages": 450},
				{"id": 3, "title": "Go Patterns", "isbn": "978-0-123456-80-2", "author_id": 2, "price": 19.99, "available": false, "pages": 200},
				{"id": 4, "title": "The Go Way", "isbn": "978-0-123456-81-9", "author_id": 2, "price": 49.99, "available": true, "pages": 500},
				{"id": 5, "title": "Go Essentials", "isbn": "978-0-123456-82-6", "author_id": 3, "price": 24.99, "available": true, "pages": 350},
			},
		},
		{
			Model: "categories",
			Data: []map[string]interface{}{
				{"id": 1, "name": "Programming", "slug": "programming"},
				{"id": 2, "name": "Web Development", "slug": "web-development"},
				{"id": 3, "name": "Database", "slug": "database"},
			},
		},
	}
}

// joinStrings joins a slice of strings with a separator
func joinStrings(strs []string, sep string) string {
	if len(strs) == 0 {
		return ""
	}
	if len(strs) == 1 {
		return strs[0]
	}
	result := strs[0]
	for _, s := range strs[1:] {
		result += sep + s
	}
	return result
}
