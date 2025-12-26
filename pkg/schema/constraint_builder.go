package schema

// ConstraintBuilder provides a fluent interface for building database constraints
type ConstraintBuilder struct {
	constraint Constraint
}

// NewConstraint creates a new constraint builder with the given name and type
func NewConstraint(name string, constraintType ConstraintType) *ConstraintBuilder {
	return &ConstraintBuilder{
		constraint: Constraint{
			Name: name,
			Type: constraintType,
		},
	}
}

// CheckConstraint creates a CHECK constraint builder
func CheckConstraint(name string, condition string) *ConstraintBuilder {
	return &ConstraintBuilder{
		constraint: Constraint{
			Name:      name,
			Type:      ConstraintTypeCheck,
			Condition: condition,
		},
	}
}

// UniqueConstraint creates a UNIQUE constraint builder
func UniqueConstraint(name string, fields ...string) *ConstraintBuilder {
	return &ConstraintBuilder{
		constraint: Constraint{
			Name:   name,
			Type:   ConstraintTypeUnique,
			Fields: fields,
		},
	}
}

// ExclusionConstraint creates an EXCLUDE constraint builder (PostgreSQL)
func ExclusionConstraint(name string, using string, operators ...string) *ConstraintBuilder {
	return &ConstraintBuilder{
		constraint: Constraint{
			Name:         name,
			Type:         ConstraintTypeExclusion,
			ExcludeUsing: using,
			ExcludeWith:  operators,
		},
	}
}

// Fields sets the fields involved in the constraint
func (b *ConstraintBuilder) Fields(fields ...string) *ConstraintBuilder {
	b.constraint.Fields = append(b.constraint.Fields, fields...)
	return b
}

// Condition sets the SQL expression for CHECK constraints
func (b *ConstraintBuilder) Condition(condition string) *ConstraintBuilder {
	b.constraint.Condition = condition
	return b
}

// Deferrable makes the constraint deferrable
func (b *ConstraintBuilder) Deferrable(deferrable DeferrableType) *ConstraintBuilder {
	b.constraint.Deferrable = deferrable
	return b
}

// NotValid marks the constraint as NOT VALID (PostgreSQL) - constraint not validated immediately
func (b *ConstraintBuilder) NotValid() *ConstraintBuilder {
	b.constraint.NotValid = true
	return b
}

// ExcludeUsing sets the index method for EXCLUDE constraints (e.g., "gist", "btree")
func (b *ConstraintBuilder) ExcludeUsing(using string) *ConstraintBuilder {
	b.constraint.ExcludeUsing = using
	return b
}

// ExcludeWith sets the operators for EXCLUDE constraints (PostgreSQL)
func (b *ConstraintBuilder) ExcludeWith(operators ...string) *ConstraintBuilder {
	b.constraint.ExcludeWith = append(b.constraint.ExcludeWith, operators...)
	return b
}

// Build returns the built constraint
func (b *ConstraintBuilder) Build() Constraint {
	return b.constraint
}

// SimpleCheck creates a simple CHECK constraint (backward compatibility helper)
func SimpleCheck(name string, condition string) Constraint {
	return Constraint{
		Name:      name,
		Type:      ConstraintTypeCheck,
		Condition: condition,
	}
}

// SimpleUnique creates a simple UNIQUE constraint
func SimpleUnique(name string, fields ...string) Constraint {
	return Constraint{
		Name:   name,
		Type:   ConstraintTypeUnique,
		Fields: fields,
	}
}

