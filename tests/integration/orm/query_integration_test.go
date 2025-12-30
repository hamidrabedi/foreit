package orm

import (
	"context"
	"testing"

	"github.com/forgego/forge/orm"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Example integration test for QuerySet
func TestQuerySet_Integration(t *testing.T) {
	withTestDB(t, func(testDB *sql.DB) {
		// Convert to forge DB type
		forgeDB := &db.DB{DB: testDB}

		// Load fixtures
		LoadFixtures(t, testDB, StandardFixtures())

		ctx := context.Background()

		t.Run("Filter and All", func(t *testing.T) {
			manager, err := query.NewManager[Book]("books")
			require.NoError(t, err)
			manager.SetDB(forgeDB)

			fa, err := manager.GetFieldAccessor()
			require.NoError(t, err)

			priceField := fa.Field[float64]("price")
			qs, err := query.NewQuerySet[Book]("books")
			require.NoError(t, err)

			books, err := qs.SetDB(forgeDB).Filter(priceField.Gt(20.0)).All(ctx)
			require.NoError(t, err)
			assert.Greater(t, len(books), 0)
		})

		t.Run("Count", func(t *testing.T) {
			manager, err := query.NewManager[Book]("books")
			require.NoError(t, err)
			manager.SetDB(forgeDB)

			fa, err := manager.GetFieldAccessor()
			require.NoError(t, err)

			availableField := fa.Field[bool]("available")
			qs, err := query.NewQuerySet[Book]("books")
			require.NoError(t, err)

			count, err := qs.SetDB(forgeDB).Filter(availableField.Eq(true)).Count(ctx)
			require.NoError(t, err)
			assert.Greater(t, count, int64(0))
		})

		t.Run("Create", func(t *testing.T) {
			manager, err := query.NewManager[Book]("books")
			require.NoError(t, err)
			manager.SetDB(forgeDB)

			newBook := &Book{
				Title:     "Test Book",
				ISBN:      "978-0-123456-99-9",
				AuthorID:  1,
				Price:     15.99,
				Available: true,
			}

			err = manager.Create(ctx, newBook)
			require.NoError(t, err)
			assert.Greater(t, newBook.ID, int64(0))

			// Verify it was created
			AssertRecordExists(t, testDB, "books", map[string]interface{}{
				"id":    newBook.ID,
				"title": "Test Book",
			})
		})

		t.Run("Update", func(t *testing.T) {
			manager, err := query.NewManager[Book]("books")
			require.NoError(t, err)
			manager.SetDB(forgeDB)

			fa, err := manager.GetFieldAccessor()
			require.NoError(t, err)

			priceField := fa.Field[float64]("price")
			qs, err := query.NewQuerySet[Book]("books")
			require.NoError(t, err)

			rowsAffected, err := qs.SetDB(forgeDB).
				Filter(priceField.Eq(29.99)).
				Update(ctx, query.UpdateMap{
					"price": 30.00,
				})
			require.NoError(t, err)
			assert.Greater(t, rowsAffected, int64(0))
		})
	})
}

// Book is a simple test model
type Book struct {
	ID        int64   `db:"id"`
	Title     string  `db:"title"`
	ISBN      string  `db:"isbn"`
	AuthorID  int64   `db:"author_id"`
	Price     float64 `db:"price"`
	Available bool    `db:"available"`
	Pages     int64   `db:"pages"`
}
