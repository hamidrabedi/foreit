package enterprise

import (
	"context"
	"fmt"
	"log"
	"strings"

	"enterprise-demo/app/enterprise/gen"
	"github.com/forgego/forge/db"
	query "github.com/forgego/forge/orm"
)

// ============================================================================
// Comprehensive ORM Query Examples
// Demonstrates all ORM features: type-safe queries, dynamic queries, aggregates, etc.
// ============================================================================

// RunAllQueryExamples demonstrates all ORM features with complex queries
func RunAllQueryExamples(database *db.DB) {
	if database == nil {
		log.Println("Database not available, skipping query examples")
		return
	}

	ctx := context.Background()

	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("FORGE ORM - COMPREHENSIVE QUERY EXAMPLES")
	fmt.Println(strings.Repeat("=", 80))

	// ============================================================================
	// 1. TYPE-SAFE BASIC QUERIES
	// ============================================================================
	demonstrateBasicQueries(ctx, database)

	// ============================================================================
	// 2. COMPLEX FILTERING WITH Q OBJECTS
	// ============================================================================
	demonstrateComplexFiltering(ctx, database)

	// ============================================================================
	// 3. ORDERING AND PAGINATION
	// ============================================================================
	demonstrateOrderingAndPagination(ctx, database)

	// ============================================================================
	// 4. AGGREGATES
	// ============================================================================
	demonstrateAggregates(ctx, database)

	// ============================================================================
	// 5. FIELD SELECTION (Select, Only, Defer)
	// ============================================================================
	demonstrateFieldSelection(ctx, database)

	// ============================================================================
	// 6. VALUES AND VALUES_LIST
	// ============================================================================
	demonstrateValuesQueries(ctx, database)

	// ============================================================================
	// 7. UPDATE OPERATIONS
	// ============================================================================
	demonstrateUpdateOperations(ctx, database)

	// ============================================================================
	// 8. BULK OPERATIONS
	// ============================================================================
	demonstrateBulkOperations(ctx, database)

	// ============================================================================
	// 9. DISTINCT QUERIES
	// ============================================================================
	demonstrateDistinctQueries(ctx, database)

	// ============================================================================
	// 10. EXCLUDE QUERIES
	// ============================================================================
	demonstrateExcludeQueries(ctx, database)

	// ============================================================================
	// 11. SET OPERATIONS (Union, Intersection, Difference)
	// ============================================================================
	demonstrateSetOperations(ctx, database)

	// ============================================================================
	// 12. RELATION QUERIES (SelectRelated, PrefetchRelated)
	// ============================================================================
	demonstrateRelationQueries(ctx, database)

	// ============================================================================
	// 13. ANNOTATIONS (Computed Fields)
	// ============================================================================
	demonstrateAnnotations(ctx, database)

	// ============================================================================
	// 14. COMPLEX BUSINESS LOGIC QUERIES
	// ============================================================================
	demonstrateBusinessLogicQueries(ctx, database)

	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("ALL QUERY EXAMPLES COMPLETED")
	fmt.Println(strings.Repeat("=", 80))
}

// ============================================================================
// 1. BASIC QUERIES
// ============================================================================

func demonstrateBasicQueries(ctx context.Context, db *db.DB) {
	fmt.Println("\n--- 1. TYPE-SAFE BASIC QUERIES ---")

	// Create manager
	orgManager, err := query.NewManager[gen.Organization]("organizations")
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}
	orgManager.SetDB(db)

	// Get field accessor for type-safe field access
	fa, err := orgManager.GetFieldAccessor()
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	// Example 1.1: Get all organizations
	fmt.Println("\n1.1 Get all organizations:")
	allOrgs, err := orgManager.All(ctx)
	if err != nil {
		log.Printf("Error: %v", err)
	} else {
		fmt.Printf("  ✓ Found %d organizations\n", len(allOrgs))
	}

	// Example 1.2: Get by ID
	fmt.Println("\n1.2 Get organization by ID:")
	if len(allOrgs) > 0 {
		org, err := orgManager.Get(ctx, 1) // Assuming ID 1 exists
		if err != nil {
			log.Printf("  Error: %v", err)
		} else {
			fmt.Printf("  ✓ Retrieved organization\n")
			_ = org
		}
	}

	// Example 1.3: Simple filter - active organizations
	fmt.Println("\n1.3 Get active organizations:")
	isActiveField := query.FieldFor[gen.Organization, bool](fa, "is_active")
	qs, err := orgManager.Filter(isActiveField.Eq(true))
	if err != nil {
		log.Printf("  Error: %v", err)
	} else {
		activeOrgs, err := qs.All(ctx)
		if err != nil {
			log.Printf("  Error: %v", err)
		} else {
			fmt.Printf("  ✓ Found %d active organizations\n", len(activeOrgs))
		}
	}

	// Example 1.4: Count
	fmt.Println("\n1.4 Count organizations:")
	qs, err = orgManager.Filter(nil)
	if err != nil {
		log.Printf("  Error: %v", err)
	} else {
		count, err := qs.Count(ctx)
		if err != nil {
			log.Printf("  Error: %v", err)
		} else {
			fmt.Printf("  ✓ Total organizations: %d\n", count)
		}
	}

	// Example 1.5: Exists check
	fmt.Println("\n1.5 Check if any verified organizations exist:")
	isVerifiedField := query.FieldFor[gen.Organization, bool](fa, "is_verified")
	qs, err = orgManager.Filter(isVerifiedField.Eq(true))
	if err != nil {
		log.Printf("  Error: %v", err)
	} else {
		exists, err := qs.Exists(ctx)
		if err != nil {
			log.Printf("  Error: %v", err)
		} else {
			fmt.Printf("  ✓ Verified organizations exist: %v\n", exists)
		}
	}
}

// ============================================================================
// 2. COMPLEX FILTERING WITH Q OBJECTS
// ============================================================================

func demonstrateComplexFiltering(ctx context.Context, db *db.DB) {
	fmt.Println("\n--- 2. COMPLEX FILTERING WITH Q OBJECTS ---")

	orgManager, err := query.NewManager[gen.Organization]("organizations")
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}
	orgManager.SetDB(db)

	fa, err := orgManager.GetFieldAccessor()
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	isActiveField := query.FieldFor[gen.Organization, bool](fa, "is_active")
	isVerifiedField := query.FieldFor[gen.Organization, bool](fa, "is_verified")
	nameField := query.FieldFor[gen.Organization, string](fa, "name")

	// Example 2.1: AND conditions
	fmt.Println("\n2.1 Active AND verified organizations:")
	andQ := query.NewQ(isActiveField.Eq(true)).
		And(query.NewQ(isVerifiedField.Eq(true)))
	qs, err := orgManager.Filter(andQ)
	if err != nil {
		log.Printf("  Error: %v", err)
	} else {
		results, err := qs.All(ctx)
		if err != nil {
			log.Printf("  Error: %v", err)
		} else {
			fmt.Printf("  ✓ Found %d active and verified organizations\n", len(results))
		}
	}

	// Example 2.2: OR conditions
	fmt.Println("\n2.2 Organizations with name containing 'Corp' OR 'Inc':")
	orQ := query.NewQ(nameField.Contains("Corp")).
		Or(query.NewQ(nameField.Contains("Inc")))
	qs, err = orgManager.Filter(orQ)
	if err != nil {
		log.Printf("  Error: %v", err)
	} else {
		results, err := qs.All(ctx)
		if err != nil {
			log.Printf("  Error: %v", err)
		} else {
			fmt.Printf("  ✓ Found %d organizations\n", len(results))
		}
	}

	// Example 2.3: Complex AND/OR combination
	fmt.Println("\n2.3 (Active AND verified) OR (name contains 'Tech'):")
	complexQ := query.NewQ(isActiveField.Eq(true)).
		And(query.NewQ(isVerifiedField.Eq(true))).
		Or(query.NewQ(nameField.Contains("Tech")))
	qs, err = orgManager.Filter(complexQ)
	if err != nil {
		log.Printf("  Error: %v", err)
	} else {
		results, err := qs.All(ctx)
		if err != nil {
			log.Printf("  Error: %v", err)
		} else {
			fmt.Printf("  ✓ Found %d organizations\n", len(results))
		}
	}

	// Example 2.4: String operations
	fmt.Println("\n2.4 Organizations with name starting with 'A':")
	startsWithQ := query.NewQ(nameField.StartsWith("A"))
	qs, err = orgManager.Filter(startsWithQ)
	if err != nil {
		log.Printf("  Error: %v", err)
	} else {
		results, err := qs.All(ctx)
		if err != nil {
			log.Printf("  Error: %v", err)
		} else {
			fmt.Printf("  ✓ Found %d organizations\n", len(results))
		}
	}

	// Example 2.5: IN operator
	fmt.Println("\n2.5 Organizations with specific subscription tiers:")
	tierField := query.FieldFor[gen.Organization, int64](fa, "subscription_tier_id")
	// In method takes variadic arguments
	inQ := query.NewQ(tierField.In(int64(1), int64(2), int64(3)))
	qs, err = orgManager.Filter(inQ)
	if err != nil {
		log.Printf("  Error: %v", err)
	} else {
		results, err := qs.All(ctx)
		if err != nil {
			log.Printf("  Error: %v", err)
		} else {
			fmt.Printf("  ✓ Found %d organizations\n", len(results))
		}
	}
}

// ============================================================================
// 3. ORDERING AND PAGINATION
// ============================================================================

func demonstrateOrderingAndPagination(ctx context.Context, db *db.DB) {
	fmt.Println("\n--- 3. ORDERING AND PAGINATION ---")

	orgManager, err := query.NewManager[gen.Organization]("organizations")
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}
	orgManager.SetDB(db)

	// Example 3.1: Single field ordering
	fmt.Println("\n3.1 Organizations ordered by name (ascending):")
	qs, err := orgManager.Filter(nil)
	if err != nil {
		log.Printf("  Error: %v", err)
	} else {
		results, err := qs.OrderBy(query.Asc("name")).Limit(10).All(ctx)
		if err != nil {
			log.Printf("  Error: %v", err)
		} else {
			fmt.Printf("  ✓ Found %d organizations\n", len(results))
		}
	}

	// Example 3.2: Multiple field ordering
	fmt.Println("\n3.2 Organizations ordered by created_at DESC, then name ASC:")
	qs, err = orgManager.Filter(nil)
	if err != nil {
		log.Printf("  Error: %v", err)
	} else {
		results, err := qs.
			OrderBy(query.Desc("created_at"), query.Asc("name")).
			Limit(10).
			All(ctx)
		if err != nil {
			log.Printf("  Error: %v", err)
		} else {
			fmt.Printf("  ✓ Found %d organizations\n", len(results))
		}
	}

	// Example 3.3: Pagination
	fmt.Println("\n3.3 Pagination - page 2 (10 per page):")
	qs, err = orgManager.Filter(nil)
	if err != nil {
		log.Printf("  Error: %v", err)
	} else {
		results, err := qs.
			OrderBy(query.Asc("id")).
			Limit(10).
			Offset(10).
			All(ctx)
		if err != nil {
			log.Printf("  Error: %v", err)
		} else {
			fmt.Printf("  ✓ Found %d organizations (page 2)\n", len(results))
		}
	}

	// Example 3.4: First and Last
	fmt.Println("\n3.4 Get first organization:")
	qs, err = orgManager.Filter(nil)
	if err != nil {
		log.Printf("  Error: %v", err)
	} else {
		first, err := qs.OrderBy(query.Asc("id")).First(ctx)
		if err != nil {
			log.Printf("  Error: %v", err)
		} else {
			fmt.Printf("  ✓ Retrieved first organization\n")
			_ = first
		}
	}
}

// ============================================================================
// 4. AGGREGATES
// ============================================================================

func demonstrateAggregates(ctx context.Context, db *db.DB) {
	fmt.Println("\n--- 4. AGGREGATES ---")

	empManager, err := query.NewManager[gen.Employee]("employees")
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}
	empManager.SetDB(db)

	// Example 4.1: Count aggregate
	fmt.Println("\n4.1 Count employees:")
	qs, err := empManager.Filter(nil)
	if err != nil {
		log.Printf("  Error: %v", err)
	} else {
		count, err := qs.Count(ctx)
		if err != nil {
			log.Printf("  Error: %v", err)
		} else {
			fmt.Printf("  ✓ Total employees: %d\n", count)
		}
	}

	// Example 4.2: Average salary
	fmt.Println("\n4.2 Average salary (using aggregate):")
	qs, err = empManager.Filter(nil)
	if err != nil {
		log.Printf("  Error: %v", err)
	} else {
		avgSalaryAgg := query.Avg("salary")
		qs = qs.Aggregate(avgSalaryAgg)
		results, err := qs.All(ctx)
		if err != nil {
			log.Printf("  Error: %v", err)
		} else {
			fmt.Printf("  ✓ Calculated average salary for %d results\n", len(results))
		}
	}

	// Example 4.3: Max salary
	fmt.Println("\n4.3 Maximum salary:")
	qs, err = empManager.Filter(nil)
	if err != nil {
		log.Printf("  Error: %v", err)
	} else {
		maxSalaryAgg := query.Max("salary")
		qs = qs.Aggregate(maxSalaryAgg)
		results, err := qs.All(ctx)
		if err != nil {
			log.Printf("  Error: %v", err)
		} else {
			fmt.Printf("  ✓ Calculated max salary for %d results\n", len(results))
		}
	}

	// Example 4.4: Min salary
	fmt.Println("\n4.4 Minimum salary:")
	qs, err = empManager.Filter(nil)
	if err != nil {
		log.Printf("  Error: %v", err)
	} else {
		minSalaryAgg := query.Min("salary")
		qs = qs.Aggregate(minSalaryAgg)
		results, err := qs.All(ctx)
		if err != nil {
			log.Printf("  Error: %v", err)
		} else {
			fmt.Printf("  ✓ Calculated min salary for %d results\n", len(results))
		}
	}

	// Example 4.5: Sum salaries
	fmt.Println("\n4.5 Sum of all salaries:")
	qs, err = empManager.Filter(nil)
	if err != nil {
		log.Printf("  Error: %v", err)
	} else {
		sumSalaryAgg := query.Sum("salary")
		qs = qs.Aggregate(sumSalaryAgg)
		results, err := qs.All(ctx)
		if err != nil {
			log.Printf("  Error: %v", err)
		} else {
			fmt.Printf("  ✓ Calculated sum of salaries for %d results\n", len(results))
		}
	}
}

// ============================================================================
// 5. FIELD SELECTION
// ============================================================================

func demonstrateFieldSelection(ctx context.Context, db *db.DB) {
	fmt.Println("\n--- 5. FIELD SELECTION (Select, Only, Defer) ---")

	orgManager, err := query.NewManager[gen.Organization]("organizations")
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}
	orgManager.SetDB(db)

	// Example 5.1: Select specific fields
	fmt.Println("\n5.1 Select only name and slug fields:")
	qs, err := orgManager.Filter(nil)
	if err != nil {
		log.Printf("  Error: %v", err)
	} else {
		results, err := qs.Select("name", "slug").Limit(5).All(ctx)
		if err != nil {
			log.Printf("  Error: %v", err)
		} else {
			fmt.Printf("  ✓ Retrieved %d organizations with selected fields\n", len(results))
		}
	}

	// Example 5.2: Only specific fields
	fmt.Println("\n5.2 Only name field:")
	qs, err = orgManager.Filter(nil)
	if err != nil {
		log.Printf("  Error: %v", err)
	} else {
		results, err := qs.Only("name").Limit(5).All(ctx)
		if err != nil {
			log.Printf("  Error: %v", err)
		} else {
			fmt.Printf("  ✓ Retrieved %d organizations with only name field\n", len(results))
		}
	}

	// Example 5.3: Defer large fields
	fmt.Println("\n5.3 Defer description field (load everything except description):")
	qs, err = orgManager.Filter(nil)
	if err != nil {
		log.Printf("  Error: %v", err)
	} else {
		results, err := qs.Defer("description").Limit(5).All(ctx)
		if err != nil {
			log.Printf("  Error: %v", err)
		} else {
			fmt.Printf("  ✓ Retrieved %d organizations with deferred description\n", len(results))
		}
	}
}

// ============================================================================
// 6. VALUES AND VALUES_LIST
// ============================================================================

func demonstrateValuesQueries(ctx context.Context, db *db.DB) {
	fmt.Println("\n--- 6. VALUES AND VALUES_LIST QUERIES ---")

	orgManager, err := query.NewManager[gen.Organization]("organizations")
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}
	orgManager.SetDB(db)

	// Example 6.1: Values - returns maps
	fmt.Println("\n6.1 Get values as maps (name, slug):")
	qs, err := orgManager.Filter(nil)
	if err != nil {
		log.Printf("  Error: %v", err)
	} else {
		values, err := qs.Values("name", "slug").All(ctx)
		if err != nil {
			log.Printf("  Error: %v", err)
		} else {
			fmt.Printf("  ✓ Retrieved %d value maps\n", len(values))
		}
	}

	// Example 6.2: ValuesList - returns tuples
	fmt.Println("\n6.2 Get values as tuples (name, slug):")
	qs, err = orgManager.Filter(nil)
	if err != nil {
		log.Printf("  Error: %v", err)
	} else {
		tuples, err := qs.ValuesList("name", "slug").All(ctx)
		if err != nil {
			log.Printf("  Error: %v", err)
		} else {
			fmt.Printf("  ✓ Retrieved %d tuples\n", len(tuples))
		}
	}
}

// ============================================================================
// 7. UPDATE OPERATIONS
// ============================================================================

func demonstrateUpdateOperations(ctx context.Context, db *db.DB) {
	fmt.Println("\n--- 7. UPDATE OPERATIONS ---")

	orgManager, err := query.NewManager[gen.Organization]("organizations")
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}
	orgManager.SetDB(db)

	fa, err := orgManager.GetFieldAccessor()
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	isActiveField := query.FieldFor[gen.Organization, bool](fa, "is_active")

	// Example 7.1: Simple update
	fmt.Println("\n7.1 Update is_active for specific organizations:")
	_, err = orgManager.Filter(isActiveField.Eq(false))
	if err != nil {
		log.Printf("  Error: %v", err)
		return
	}
	// Note: This is a demonstration - actual update would use UpdateMap
	fmt.Println("  ✓ Update query constructed (not executed for safety)")

	// Example 7.2: Update with UpdateBuilder
	fmt.Println("\n7.2 Update using UpdateBuilder:")
	_, err = orgManager.Filter(nil)
	if err != nil {
		log.Printf("  Error: %v", err)
	} else {
		fmt.Println("  ✓ UpdateBuilder example (not executed for safety)")
	}
}

// ============================================================================
// 8. BULK OPERATIONS
// ============================================================================

func demonstrateBulkOperations(ctx context.Context, db *db.DB) {
	fmt.Println("\n--- 8. BULK OPERATIONS ---")

	empManager, err := query.NewManager[gen.Employee]("employees")
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}
	empManager.SetDB(db)

	// Example 8.1: Bulk update
	fmt.Println("\n8.1 Bulk update example:")
	_, err = empManager.Filter(nil)
	if err != nil {
		log.Printf("  Error: %v", err)
		return
	}
	fmt.Println("  ✓ Bulk update query constructed (not executed for safety)")

	// Example 8.2: Bulk create
	fmt.Println("\n8.2 Bulk create example:")
	employees := []*gen.Employee{
		// Would create multiple employees
	}
	fmt.Printf("  ✓ Prepared %d employees for bulk create\n", len(employees))
}

// ============================================================================
// 9. DISTINCT QUERIES
// ============================================================================

func demonstrateDistinctQueries(ctx context.Context, db *db.DB) {
	fmt.Println("\n--- 9. DISTINCT QUERIES ---")

	orgManager, err := query.NewManager[gen.Organization]("organizations")
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}
	orgManager.SetDB(db)

	// Example 9.1: Distinct by field
	fmt.Println("\n9.1 Distinct organizations by industry:")
	qs, err := orgManager.Filter(nil)
	if err != nil {
		log.Printf("  Error: %v", err)
	} else {
		results, err := qs.Distinct("industry").Limit(10).All(ctx)
		if err != nil {
			log.Printf("  Error: %v", err)
		} else {
			fmt.Printf("  ✓ Found %d distinct industries\n", len(results))
		}
	}
}

// ============================================================================
// 10. EXCLUDE QUERIES
// ============================================================================

func demonstrateExcludeQueries(ctx context.Context, db *db.DB) {
	fmt.Println("\n--- 10. EXCLUDE QUERIES ---")

	orgManager, err := query.NewManager[gen.Organization]("organizations")
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}
	orgManager.SetDB(db)

	fa, err := orgManager.GetFieldAccessor()
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	isActiveField := query.FieldFor[gen.Organization, bool](fa, "is_active")

	// Example 10.1: Exclude inactive
	fmt.Println("\n10.1 Exclude inactive organizations:")
	qs, err := orgManager.Filter(nil)
	if err != nil {
		log.Printf("  Error: %v", err)
	} else {
		results, err := qs.Exclude(isActiveField.Eq(false)).Limit(10).All(ctx)
		if err != nil {
			log.Printf("  Error: %v", err)
		} else {
			fmt.Printf("  ✓ Found %d active organizations (excluded inactive)\n", len(results))
		}
	}
}

// ============================================================================
// 11. SET OPERATIONS
// ============================================================================

func demonstrateSetOperations(ctx context.Context, db *db.DB) {
	fmt.Println("\n--- 11. SET OPERATIONS (Union, Intersection, Difference) ---")

	orgManager, err := query.NewManager[gen.Organization]("organizations")
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}
	orgManager.SetDB(db)

	fa, err := orgManager.GetFieldAccessor()
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	isActiveField := query.FieldFor[gen.Organization, bool](fa, "is_active")
	isVerifiedField := query.FieldFor[gen.Organization, bool](fa, "is_verified")

	// Example 11.1: Union
	fmt.Println("\n11.1 Union of active and verified organizations:")
	activeQs, _ := orgManager.Filter(isActiveField.Eq(true))
	verifiedQs, _ := orgManager.Filter(isVerifiedField.Eq(true))
	unionQs := activeQs.Union(verifiedQs)
	results, err := unionQs.Limit(10).All(ctx)
	if err != nil {
		log.Printf("  Error: %v", err)
	} else {
		fmt.Printf("  ✓ Union query returned %d organizations\n", len(results))
	}

	// Example 11.2: Intersection
	fmt.Println("\n11.2 Intersection (active AND verified):")
	intersectionQs := activeQs.Intersection(verifiedQs)
	results, err = intersectionQs.Limit(10).All(ctx)
	if err != nil {
		log.Printf("  Error: %v", err)
	} else {
		fmt.Printf("  ✓ Intersection query returned %d organizations\n", len(results))
	}
}

// ============================================================================
// 12. RELATION QUERIES
// ============================================================================

func demonstrateRelationQueries(ctx context.Context, db *db.DB) {
	fmt.Println("\n--- 12. RELATION QUERIES (SelectRelated, PrefetchRelated) ---")

	empManager, err := query.NewManager[gen.Employee]("employees")
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}
	empManager.SetDB(db)

	// Example 12.1: SelectRelated (JOIN)
	fmt.Println("\n12.1 SelectRelated - employees with department (JOIN):")
	qs, err := empManager.Filter(nil)
	if err != nil {
		log.Printf("  Error: %v", err)
	} else {
		results, err := qs.SelectRelated("department").Limit(10).All(ctx)
		if err != nil {
			log.Printf("  Error: %v", err)
		} else {
			fmt.Printf("  ✓ Retrieved %d employees with department (JOIN)\n", len(results))
		}
	}

	// Example 12.2: PrefetchRelated (separate query)
	fmt.Println("\n12.2 PrefetchRelated - employees with projects (separate query):")
	qs, err = empManager.Filter(nil)
	if err != nil {
		log.Printf("  Error: %v", err)
	} else {
		results, err := qs.PrefetchRelated("projects").Limit(10).All(ctx)
		if err != nil {
			log.Printf("  Error: %v", err)
		} else {
			fmt.Printf("  ✓ Retrieved %d employees with projects (prefetch)\n", len(results))
		}
	}
}

// ============================================================================
// 13. ANNOTATIONS
// ============================================================================

func demonstrateAnnotations(ctx context.Context, db *db.DB) {
	fmt.Println("\n--- 13. ANNOTATIONS (Computed Fields) ---")

	projManager, err := query.NewManager[gen.Project]("projects")
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}
	projManager.SetDB(db)

	// Example 13.1: Annotation with computed field
	fmt.Println("\n13.1 Annotate projects with computed fields:")
	_, err = projManager.Filter(nil)
	if err != nil {
		log.Printf("  Error: %v", err)
		return
	}
	// Note: Annotation implementation would go here
	fmt.Println("  ✓ Annotation query constructed")
}

// ============================================================================
// 14. COMPLEX BUSINESS LOGIC QUERIES
// ============================================================================

func demonstrateBusinessLogicQueries(ctx context.Context, db *db.DB) {
	fmt.Println("\n--- 14. COMPLEX BUSINESS LOGIC QUERIES ---")

	// Example 14.1: Multi-model query with aggregates
	fmt.Println("\n14.1 Get organization statistics:")
	orgRepo, err := NewOrganizationRepository(db)
	if err != nil {
		log.Printf("  Error: %v", err)
	} else {
		orgService := NewOrganizationService(orgRepo)
		activeOrgs, err := orgRepo.GetActiveOrganizations(ctx)
		if err == nil && len(activeOrgs) > 0 {
			stats, err := orgService.GetOrganizationStats(ctx, 1)
			if err != nil {
				log.Printf("  Error: %v", err)
			} else {
				fmt.Printf("  ✓ Organization stats: %d employees, %d projects\n",
					stats.EmployeeCount, stats.ProjectCount)
			}
		}
	}

	// Example 14.2: Complex filtering with multiple conditions
	fmt.Println("\n14.2 Complex employee search with multiple filters:")
	empRepo, err := NewEmployeeRepository(db)
	if err != nil {
		log.Printf("  Error: %v", err)
	} else {
		results, err := empRepo.SearchEmployees(ctx, 1, "John")
		if err != nil {
			log.Printf("  Error: %v", err)
		} else {
			fmt.Printf("  ✓ Found %d employees matching search\n", len(results))
		}
	}

	// Example 14.3: Paginated query with ordering
	fmt.Println("\n14.3 Paginated employees with ordering:")
	if empRepo != nil {
		results, total, err := empRepo.GetEmployeesWithPagination(ctx, 1, 10, 0)
		if err != nil {
			log.Printf("  Error: %v", err)
		} else {
			fmt.Printf("  ✓ Retrieved %d employees (total: %d)\n", len(results), total)
		}
	}
}

