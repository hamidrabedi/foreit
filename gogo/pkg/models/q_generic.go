package models

// Q represents a type-safe query condition (Django-style)
// Using map[string]interface{} but with type-safe field references
type Q map[string]interface{}

// And combines two Q objects with AND
func (q Q) And(other Q) Q {
	result := make(Q)
	for k, v := range q {
		result[k] = v
	}
	for k, v := range other {
		result[k] = v
	}
	return result
}

// Or creates an OR condition (requires special handling in query builder)
func (q Q) Or(other Q) Q {
	// OR conditions need special handling
	// For now, return a combined Q (implementation would need OR support)
	return q.And(other)
}

// Not negates a condition
func (q Q) Not() Q {
	result := make(Q)
	for k, v := range q {
		// Add __not lookup
		if field, lookup := parseLookup(k); lookup == "exact" {
			result[field+"__ne"] = v
		} else {
			result[k+"__not"] = v
		}
	}
	return result
}

// Type-safe Q building with field references
// Example:
//   Email := NewFieldRef[string]("email")
//   q := Email.Contains("@example.com").And(Email.IsNotNull())

