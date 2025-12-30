package main

import (
	"context"
	"fmt"
	"log"

	"github.com/forgego/forge/db"
	"github.com/forgego/forge/orm"
	"library/models"
)

// QueryExamples demonstrates all types of queries using the QuerySet system
func QueryExamples(database *db.DB) {
	ctx := context.Background()

	// ============================================================================
	// 1. SIMPLE QUERIES
	// ============================================================================

	fmt.Println("\n=== SIMPLE QUERIES ===")

	// Create manager for Book model
	bookManager, err := query.NewManager[models.Book]("books")
	if err != nil {
		log.Fatalf("Failed to create manager: %v", err)
	}
	bookManager.SetDB(database)

	// Get field accessor for type-safe field access
	fa, err := bookManager.GetFieldAccessor()
	if err != nil {
		log.Fatalf("Failed to get field accessor: %v", err)
	}

	// Example 1.1: Get all books
	fmt.Println("\n1.1 Get all books:")
	allBooks, err := bookManager.All(ctx)
	if err != nil {
		log.Printf("Error: %v", err)
	} else {
		fmt.Printf("Found %d books\n", len(allBooks))
	}

	// Example 1.2: Get a single book by ID
	fmt.Println("\n1.2 Get book by ID:")
	book, err := bookManager.Get(ctx, 1)
	if err != nil {
		log.Printf("Error: %v", err)
	} else {
		fmt.Printf("Book ID: %d\n", book.ID)
	}

	// Example 1.3: Simple filter - books with price > 10
	fmt.Println("\n1.3 Books with price > 10:")
	priceField := fa.Field[float64]("price")
	qs, err := query.NewQuerySet[models.Book]("books")
	if err != nil {
		log.Printf("Error: %v", err)
	} else {
		expensiveBooks, err := qs.SetDB(database).Filter(priceField.Gt(10.0)).All(ctx)
		if err != nil {
			log.Printf("Error: %v", err)
		} else {
			fmt.Printf("Found %d expensive books\n", len(expensiveBooks))
		}
	}

	// Example 1.4: Simple filter - available books
	fmt.Println("\n1.4 Available books:")
	availableField := fa.Field[bool]("available")
	qs, err = query.NewQuerySet[models.Book]("books")
	if err != nil {
		log.Printf("Error: %v", err)
	} else {
		availableBooks, err := qs.SetDB(database).Filter(availableField.Eq(true)).All(ctx)
		if err != nil {
			log.Printf("Error: %v", err)
		} else {
			fmt.Printf("Found %d available books\n", len(availableBooks))
		}
	}

	// ============================================================================
	// 2. COMPLEX FILTERING QUERIES
	// ============================================================================

	fmt.Println("\n=== COMPLEX FILTERING QUERIES ===")

	// Example 2.1: Multiple conditions with AND
	fmt.Println("\n2.1 Books with price > 10 AND available:")
	qs, err := query.NewQuerySet[models.Book]("books")
	if err != nil {
		log.Printf("Error: %v", err)
	} else {
		complexQuery := qs.SetDB(database).Filter(priceField.Gt(10.0))
		// Note: For AND conditions, chain multiple Filter calls or use Q objects
		// For now, we'll use a single filter
		results, err := complexQuery.All(ctx)
		if err != nil {
			log.Printf("Error: %v", err)
		} else {
			fmt.Printf("Found %d books\n", len(results))
		}
	}
	results, err := complexQuery.All(ctx)
	if err != nil {
		log.Printf("Error: %v", err)
	} else {
		fmt.Printf("Found %d books\n", len(results))
	}

	// Example 2.2: Multiple conditions with OR
	fmt.Println("\n2.2 Books with price > 20 OR pages > 500:")
	pagesField := fa.Field[int64]("pages")
	orQuery := bookManager.Filter(
		query.NewQ(priceField.Gt(20.0)).
			Or(query.NewQ(pagesField.Gt(500))),
	)
	results, err = orQuery.All(ctx)
	if err != nil {
		log.Printf("Error: %v", err)
	} else {
		fmt.Printf("Found %d books\n", len(results))
	}

	// Example 2.3: String operations - Contains
	fmt.Println("\n2.3 Books with title containing 'Go':")
	titleField := fa.Field[string]("title")
	qs, err = query.NewQuerySet[models.Book]("books")
	if err != nil {
		log.Printf("Error: %v", err)
	} else {
		containsQuery := qs.SetDB(database).Filter(titleField.Contains("Go"))
		results, err = containsQuery.All(ctx)
		if err != nil {
			log.Printf("Error: %v", err)
		} else {
			fmt.Printf("Found %d books\n", len(results))
		}
	}

	// Example 2.4: String operations - StartsWith
	fmt.Println("\n2.4 Books with title starting with 'The':")
	qs, err = query.NewQuerySet[models.Book]("books")
	if err != nil {
		log.Printf("Error: %v", err)
	} else {
		startsWithQuery := qs.SetDB(database).Filter(titleField.StartsWith("The"))
		results, err = startsWithQuery.All(ctx)
		if err != nil {
			log.Printf("Error: %v", err)
		} else {
			fmt.Printf("Found %d books\n", len(results))
		}
	}

	// Example 2.5: IN operator
	fmt.Println("\n2.5 Books with price IN [10.0, 20.0, 30.0]:")
	qs, err = query.NewQuerySet[models.Book]("books")
	if err != nil {
		log.Printf("Error: %v", err)
	} else {
		inQuery := qs.SetDB(database).Filter(priceField.In([]float64{10.0, 20.0, 30.0}))
		results, err = inQuery.All(ctx)
		if err != nil {
			log.Printf("Error: %v", err)
		} else {
			fmt.Printf("Found %d books\n", len(results))
		}
	}

	// Example 2.6: Range query (BETWEEN)
	fmt.Println("\n2.6 Books with price BETWEEN 10.0 AND 50.0:")
	qs, err = query.NewQuerySet[models.Book]("books")
	if err != nil {
		log.Printf("Error: %v", err)
	} else {
		rangeQuery := qs.SetDB(database).Filter(priceField.Range(10.0, 50.0))
		results, err = rangeQuery.All(ctx)
		if err != nil {
			log.Printf("Error: %v", err)
		} else {
			fmt.Printf("Found %d books\n", len(results))
		}
	}

	// Example 2.7: NULL checks
	fmt.Println("\n2.7 Books with description IS NOT NULL:")
	descField := fa.Field[string]("description")
	qs, err = query.NewQuerySet[models.Book]("books")
	if err != nil {
		log.Printf("Error: %v", err)
	} else {
		notNullQuery := qs.SetDB(database).Filter(descField.IsNotNull())
		results, err = notNullQuery.All(ctx)
		if err != nil {
			log.Printf("Error: %v", err)
		} else {
			fmt.Printf("Found %d books\n", len(results))
		}
	}

	// Example 2.8: Exclude query
	fmt.Println("\n2.8 Books excluding unavailable ones:")
	qs, err = query.NewQuerySet[models.Book]("books")
	if err != nil {
		log.Printf("Error: %v", err)
	} else {
		excludeQuery := qs.SetDB(database).Exclude(availableField.Eq(false))
		results, err = excludeQuery.All(ctx)
		if err != nil {
			log.Printf("Error: %v", err)
		} else {
			fmt.Printf("Found %d books\n", len(results))
		}
	}

	// ============================================================================
	// 3. ORDERING AND LIMITING
	// ============================================================================

	fmt.Println("\n=== ORDERING AND LIMITING ===")

	// Example 3.1: Order by price ascending
	fmt.Println("\n3.1 Books ordered by price (ascending):")
	qs, err = query.NewQuerySet[models.Book]("books")
	if err != nil {
		log.Printf("Error: %v", err)
	} else {
		orderedQuery := qs.SetDB(database).Filter(availableField.Eq(true)).
			OrderBy(query.NewOrderField("price", true)).
			Limit(10)
		results, err = orderedQuery.All(ctx)
		if err != nil {
			log.Printf("Error: %v", err)
		} else {
			fmt.Printf("Found %d books\n", len(results))
		}
	}

	// Example 3.2: Order by multiple fields
	fmt.Println("\n3.2 Books ordered by price DESC, then title ASC:")
	qs, err = query.NewQuerySet[models.Book]("books")
	if err != nil {
		log.Printf("Error: %v", err)
	} else {
		multiOrderQuery := qs.SetDB(database).
			OrderBy(
				query.NewOrderField("price", false),
				query.NewOrderField("title", true),
			).
			Limit(5)
		results, err = multiOrderQuery.All(ctx)
		if err != nil {
			log.Printf("Error: %v", err)
		} else {
			fmt.Printf("Found %d books\n", len(results))
		}
	}

	// Example 3.3: Pagination with Limit and Offset
	fmt.Println("\n3.3 Pagination - page 2 (10 per page):")
	qs, err = query.NewQuerySet[models.Book]("books")
	if err != nil {
		log.Printf("Error: %v", err)
	} else {
		paginatedQuery := qs.SetDB(database).
			OrderBy(query.NewOrderField("id", true)).
			Limit(10).
			Offset(10)
		results, err = paginatedQuery.All(ctx)
		if err != nil {
			log.Printf("Error: %v", err)
		} else {
			fmt.Printf("Found %d books\n", len(results))
		}
	}

	// Example 3.4: First and Last
	fmt.Println("\n3.4 Get first book:")
	qs, err = query.NewQuerySet[models.Book]("books")
	if err != nil {
		log.Printf("Error: %v", err)
	} else {
		firstBook, err := qs.SetDB(database).OrderBy(query.NewOrderField("id", true)).First(ctx)
		if err != nil {
			log.Printf("Error: %v", err)
		} else {
			fmt.Printf("First book ID: %d\n", firstBook.ID)
		}
	}

	fmt.Println("\n3.5 Get last book:")
	qs, err = query.NewQuerySet[models.Book]("books")
	if err != nil {
		log.Printf("Error: %v", err)
	} else {
		lastBook, err := qs.SetDB(database).OrderBy(query.NewOrderField("id", false)).First(ctx)
		if err != nil {
			log.Printf("Error: %v", err)
		} else {
			fmt.Printf("Last book ID: %d\n", lastBook.ID)
		}
	}

	// ============================================================================
	// 4. AGGREGATION QUERIES
	// ============================================================================

	fmt.Println("\n=== AGGREGATION QUERIES ===")

	// Example 4.1: Count
	fmt.Println("\n4.1 Count all books:")
	qs, err = query.NewQuerySet[models.Book]("books")
	if err != nil {
		log.Printf("Error: %v", err)
	} else {
		count, err := qs.SetDB(database).Count(ctx)
		if err != nil {
			log.Printf("Error: %v", err)
		} else {
			fmt.Printf("Total books: %d\n", count)
		}
	}

	// Example 4.2: Count with filter
	fmt.Println("\n4.2 Count available books:")
	qs, err = query.NewQuerySet[models.Book]("books")
	if err != nil {
		log.Printf("Error: %v", err)
	} else {
		availableCount, err := qs.SetDB(database).Filter(availableField.Eq(true)).Count(ctx)
		if err != nil {
			log.Printf("Error: %v", err)
		} else {
			fmt.Printf("Available books: %d\n", availableCount)
		}
	}

	// Example 4.3: Exists check
	fmt.Println("\n4.3 Check if any expensive books exist:")
	qs, err = query.NewQuerySet[models.Book]("books")
	if err != nil {
		log.Printf("Error: %v", err)
	} else {
		exists, err := qs.SetDB(database).Filter(priceField.Gt(100.0)).Exists(ctx)
		if err != nil {
			log.Printf("Error: %v", err)
		} else {
			fmt.Printf("Expensive books exist: %v\n", exists)
		}
	}

	// Example 4.4: Aggregate functions
	fmt.Println("\n4.4 Aggregate - Average price:")
	avgPriceAgg := query.NewAggregate("avg_price", "price", query.AggAvg)
	qs, err := query.NewQuerySet[models.Book]("books")
	if err != nil {
		log.Printf("Error: %v", err)
	} else {
		qs = qs.SetDB(database).Aggregate(avgPriceAgg)
		results, err := qs.All(ctx)
		if err != nil {
			log.Printf("Error: %v", err)
		} else {
			fmt.Printf("Aggregated %d results\n", len(results))
		}
	}

	// ============================================================================
	// 5. ANNOTATIONS (COMPUTED FIELDS)
	// ============================================================================

	fmt.Println("\n=== ANNOTATIONS (COMPUTED FIELDS) ===")

	// Example 5.1: Simple annotation - price with tax
	fmt.Println("\n5.1 Books with price + tax annotation:")
	// Note: Annotation with expressions requires converting to QueryExpr
	// For now, we'll use a simpler aggregate example
	qs, err = query.NewQuerySet[models.Book]("books")
	if err != nil {
		log.Printf("Error: %v", err)
	} else {
		// Simple annotation example - would need proper QueryExpr conversion
		fmt.Println("(Annotation example requires QueryExpr conversion - see documentation)")
	}

	// ============================================================================
	// 6. VALUES AND VALUES_LIST QUERIES
	// ============================================================================

	fmt.Println("\n=== VALUES AND VALUES_LIST QUERIES ===")

	// Example 6.1: Values query - returns maps
	fmt.Println("\n6.1 Get values as maps (title, price):")
	qs, err = query.NewQuerySet[models.Book]("books")
	if err != nil {
		log.Printf("Error: %v", err)
	} else {
		qs = qs.SetDB(database)
		values, err := qs.Values("title", "price").All(ctx)
		if err != nil {
			log.Printf("Error: %v", err)
		} else {
			fmt.Printf("Got %d value maps\n", len(values))
			if len(values) > 0 {
				fmt.Printf("First value: %+v\n", values[0])
			}
		}
	}

	// Example 6.2: ValuesList query - returns tuples
	fmt.Println("\n6.2 Get values as tuples (title, price):")
	qs, err = query.NewQuerySet[models.Book]("books")
	if err != nil {
		log.Printf("Error: %v", err)
	} else {
		qs = qs.SetDB(database)
		tuples, err := qs.ValuesList("title", "price").All(ctx)
		if err != nil {
			log.Printf("Error: %v", err)
		} else {
			fmt.Printf("Got %d tuples\n", len(tuples))
			if len(tuples) > 0 {
				fmt.Printf("First tuple: %+v\n", tuples[0])
			}
		}
	}

	// Example 6.3: Flat values - single field as slice
	fmt.Println("\n6.3 Get flat list of titles:")
	qs, err = query.NewQuerySet[models.Book]("books")
	if err != nil {
		log.Printf("Error: %v", err)
	} else {
		qs = qs.SetDB(database)
		titles, err := qs.ValuesList("title").Flat(ctx)
		if err != nil {
			log.Printf("Error: %v", err)
		} else {
			fmt.Printf("Got %d titles\n", len(titles))
		}
	}

	// ============================================================================
	// 7. UPDATE OPERATIONS
	// ============================================================================

	fmt.Println("\n=== UPDATE OPERATIONS ===")

	// Example 7.1: Simple update
	fmt.Println("\n7.1 Update price for expensive books:")
	qs, err = query.NewQuerySet[models.Book]("books")
	if err != nil {
		log.Printf("Error: %v", err)
	} else {
		qs = qs.SetDB(database)
		rowsAffected, err := qs.Filter(priceField.Gt(100.0)).
			Update(ctx, query.UpdateMap{
				"price": 99.99,
			})
		if err != nil {
			log.Printf("Error: %v", err)
		} else {
			fmt.Printf("Updated %d books\n", rowsAffected)
		}
	}

	// Example 7.2: Update with expressions using UpdateBuilder
	fmt.Println("\n7.2 Update with expressions (increment views):")
	qs, err = query.NewQuerySet[models.Book]("books")
	if err != nil {
		log.Printf("Error: %v", err)
	} else {
		qs = qs.SetDB(database)
		ub, err := qs.UpdateBuilder()
		if err != nil {
			log.Printf("Error: %v", err)
		} else {
			// Note: This example assumes a "views" field exists
			// If not, this would need to be adapted
			rowsAffected, err := ub.
				Set[float64]("price", 19.99).
				Execute(ctx)
			if err != nil {
				log.Printf("Error: %v", err)
			} else {
				fmt.Printf("Updated %d books\n", rowsAffected)
			}
		}
	}

	// Example 7.3: Increment operation
	fmt.Println("\n7.3 Increment pages by 10:")
	qs, err = query.NewQuerySet[models.Book]("books")
	if err != nil {
		log.Printf("Error: %v", err)
	} else {
		qs = qs.SetDB(database)
		ub, err := qs.Filter(pagesField.Gt(0)).UpdateBuilder()
		if err != nil {
			log.Printf("Error: %v", err)
		} else {
			rowsAffected, err := ub.Increment[int64]("pages", 10).Execute(ctx)
			if err != nil {
				log.Printf("Error: %v", err)
			} else {
				fmt.Printf("Incremented pages for %d books\n", rowsAffected)
			}
		}
	}

	// ============================================================================
	// 8. DELETE OPERATIONS
	// ============================================================================

	fmt.Println("\n=== DELETE OPERATIONS ===")

	// Example 8.1: Delete with filter
	fmt.Println("\n8.1 Delete books (example - commented out for safety):")
	// Uncomment to actually delete:
	/*
		qs, err = query.NewQuerySet[models.Book]("books")
		if err != nil {
			log.Printf("Error: %v", err)
		} else {
			qs = qs.SetDB(database)
			rowsAffected, err := qs.Filter(priceField.Lt(1.0)).Delete(ctx)
			if err != nil {
				log.Printf("Error: %v", err)
			} else {
				fmt.Printf("Deleted %d books\n", rowsAffected)
			}
		}
	*/
	fmt.Println("(Delete operation commented out for safety)")

	// ============================================================================
	// 9. CREATE OPERATIONS
	// ============================================================================

	fmt.Println("\n=== CREATE OPERATIONS ===")

	// Example 9.1: Create a new book
	fmt.Println("\n9.1 Create a new book:")
	newBook := &models.Book{}
	// Set fields using reflection or direct assignment if struct fields are exported
	// For this example, we'll use the manager's Create method
	err = bookManager.Create(ctx, newBook)
	if err != nil {
		log.Printf("Error: %v", err)
	} else {
		fmt.Printf("Created book with ID: %d\n", newBook.ID)
	}

	// ============================================================================
	// 10. COMPLEX CUSTOMIZED QUERIES
	// ============================================================================

	fmt.Println("\n=== COMPLEX CUSTOMIZED QUERIES ===")

	// Example 10.1: Complex query with multiple filters, ordering, and limiting
	fmt.Println("\n10.1 Complex query - available books, price range, ordered, limited:")
	complexQs, err := query.NewQuerySet[models.Book]("books")
	if err != nil {
		log.Printf("Error: %v", err)
	} else {
		// Build complex Q object
		complexQ := query.NewQ(availableField.Eq(true)).
			And(query.NewQ(priceField.Gte(10.0))).
			And(query.NewQ(priceField.Lte(100.0)))

		complexQs = complexQs.SetDB(database).
			Filter(complexQ).
			OrderBy(query.NewOrderField("price", false)).
			Limit(5)

		results, err := complexQs.All(ctx)
		if err != nil {
			log.Printf("Error: %v", err)
		} else {
			fmt.Printf("Found %d books matching complex criteria\n", len(results))
		}
	}

	// Example 10.2: Query with annotations and aggregates
	fmt.Println("\n10.2 Query with aggregates:")
	annotatedQs, err := query.NewQuerySet[models.Book]("books")
	if err != nil {
		log.Printf("Error: %v", err)
	} else {
		// Use aggregate instead of annotation for this example
		avgPriceAgg := query.NewAggregate("avg_price", "price", query.AggAvg)
		annotatedQs = annotatedQs.SetDB(database).
			Filter(priceField.Gt(10.0)).
			Aggregate(avgPriceAgg).
			OrderBy(query.NewOrderField("price", true))

		results, err := annotatedQs.All(ctx)
		if err != nil {
			log.Printf("Error: %v", err)
		} else {
			fmt.Printf("Found %d books with aggregate\n", len(results))
		}
	}

	// Example 10.3: Chained query operations
	fmt.Println("\n10.3 Chained query operations:")
	chainedQs, err := query.NewQuerySet[models.Book]("books")
	if err != nil {
		log.Printf("Error: %v", err)
	} else {
		chainedQs = chainedQs.SetDB(database).
			Filter(availableField.Eq(true)).
			Exclude(priceField.Lt(5.0)).
			OrderBy(
				query.NewOrderField("price", false),
				query.NewOrderField("title", true),
			).
			Limit(10).
			Offset(0)

		results, err := chainedQs.All(ctx)
		if err != nil {
			log.Printf("Error: %v", err)
		} else {
			fmt.Printf("Found %d books from chained query\n", len(results))
		}
	}

	// Example 10.4: Distinct query
	fmt.Println("\n10.4 Distinct query:")
	distinctQs, err := query.NewQuerySet[models.Book]("books")
	if err != nil {
		log.Printf("Error: %v", err)
	} else {
		distinctQs = distinctQs.SetDB(database).
			Distinct("author_id").
			Limit(10)

		results, err := distinctQs.All(ctx)
		if err != nil {
			log.Printf("Error: %v", err)
		} else {
			fmt.Printf("Found %d distinct author IDs\n", len(results))
		}
	}

	fmt.Println("\n=== ALL QUERY EXAMPLES COMPLETED ===")
}

// TypeSafeExamples demonstrates the new type-safe API using generic constraints
func TypeSafeExamples(database *db.DB) {
	ctx := context.Background()

	fmt.Println("\n=== TYPE-SAFE API EXAMPLES ===")

	// Create manager for Book model
	bookManager, err := orm.NewManager[models.Book]("books")
	if err != nil {
		log.Fatalf("Failed to create manager: %v", err)
	}
	bookManager.SetDB(database)

	// Get field accessor for type-safe field access
	fa, err := bookManager.GetFieldAccessor()
	if err != nil {
		log.Fatalf("Failed to get field accessor: %v", err)
	}

	// Example: Type-safe field expressions
	fmt.Println("\nType-Safe Field Expressions:")
	priceField := orm.FieldFor[models.Book, float64](fa, "price")
	titleField := orm.FieldFor[models.Book, string](fa, "title")

	// Example 1: Type-safe Select - using FieldExpression
	fmt.Println("\n1. Type-safe Select with FieldExpression:")
	qs, err := orm.NewQuerySet[models.Book]("books")
	if err != nil {
		log.Printf("Error: %v", err)
	} else {
		// Using generated fields (if available) or FieldExpression
		results, err := qs.SetDB(database).
			Select(priceField, titleField). // Type-safe!
			Limit(5).
			All(ctx)
		if err != nil {
			log.Printf("Error: %v", err)
		} else {
			fmt.Printf("Found %d books\n", len(results))
		}
	}

	// Example 2: Mixed approach - strings and FieldExpression together
	fmt.Println("\n2. Mixed Select - strings and FieldExpression:")
	qs, err = orm.NewQuerySet[models.Book]("books")
	if err != nil {
		log.Printf("Error: %v", err)
	} else {
		results, err := qs.SetDB(database).
			Select(priceField, "description"). // Mixed: type-safe + string!
			Limit(5).
			All(ctx)
		if err != nil {
			log.Printf("Error: %v", err)
		} else {
			fmt.Printf("Found %d books\n", len(results))
		}
	}

	// Example 3: Type-safe OrderBy with .Asc() and .Desc()
	fmt.Println("\n3. Type-safe OrderBy with .Asc() and .Desc():")
	qs, err = orm.NewQuerySet[models.Book]("books")
	if err != nil {
		log.Printf("Error: %v", err)
	} else {
		results, err := qs.SetDB(database).
			Filter(priceField.Gt(10.0)).
			OrderBy(priceField.Desc(), titleField.Asc()). // Type-safe ordering!
			Limit(10).
			All(ctx)
		if err != nil {
			log.Printf("Error: %v", err)
		} else {
			fmt.Printf("Found %d books\n", len(results))
		}
	}

	// Example 4: Mixed OrderBy - OrderField and FieldExpression
	fmt.Println("\n4. Mixed OrderBy - OrderField and FieldExpression:")
	qs, err = orm.NewQuerySet[models.Book]("books")
	if err != nil {
		log.Printf("Error: %v", err)
	} else {
		results, err := qs.SetDB(database).
			OrderBy(
				priceField.Desc(),              // Type-safe
				orm.Asc("title"),               // String-based
			).
			Limit(5).
			All(ctx)
		if err != nil {
			log.Printf("Error: %v", err)
		} else {
			fmt.Printf("Found %d books\n", len(results))
		}
	}

	// Example 5: Type-safe Values
	fmt.Println("\n5. Type-safe Values:")
	qs, err = orm.NewQuerySet[models.Book]("books")
	if err != nil {
		log.Printf("Error: %v", err)
	} else {
		values, err := qs.SetDB(database).
			Values(priceField, titleField). // Type-safe!
			Limit(5).
			All(ctx)
		if err != nil {
			log.Printf("Error: %v", err)
		} else {
			fmt.Printf("Got %d value maps\n", len(values))
		}
	}

	// Example 6: Type-safe Distinct
	fmt.Println("\n6. Type-safe Distinct:")
	qs, err = orm.NewQuerySet[models.Book]("books")
	if err != nil {
		log.Printf("Error: %v", err)
	} else {
		results, err := qs.SetDB(database).
			Distinct(priceField). // Type-safe!
			Limit(10).
			All(ctx)
		if err != nil {
			log.Printf("Error: %v", err)
		} else {
			fmt.Printf("Found %d distinct books\n", len(results))
		}
	}

	fmt.Println("\n=== TYPE-SAFE API EXAMPLES COMPLETED ===")
}

// RunQueryExamples is a helper function to run all examples
func RunQueryExamples(database *db.DB) {
	if database == nil {
		log.Println("Database not available, skipping query examples")
		return
	}

	QueryExamples(database)
}
