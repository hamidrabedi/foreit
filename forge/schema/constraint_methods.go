package schema

// CheckConstraint creates a CHECK constraint.
func CheckConstraint(name string, condition string) Constraint {
	return Constraint{
		Name:      name,
		Type:      ConstraintTypeCheck,
		Condition: condition,
	}
}

// UniqueConstraint creates a UNIQUE constraint.
func UniqueConstraint(name string, fields ...string) Constraint {
	return Constraint{
		Name:   name,
		Type:   ConstraintTypeUnique,
		Fields: fields,
	}
}

// ExclusionConstraint creates an EXCLUDE constraint (PostgreSQL).
func ExclusionConstraint(name string, using string, operators ...string) Constraint {
	return Constraint{
		Name:         name,
		Type:         ConstraintTypeExclusion,
		ExcludeUsing: using,
		ExcludeWith:  operators,
	}
}

func (c Constraint) WithFields(fields ...string) Constraint {
	c.Fields = append(c.Fields, fields...)
	return c
}

func (c Constraint) WithCondition(condition string) Constraint {
	c.Condition = condition
	return c
}

func (c Constraint) WithDeferrable(deferrable DeferrableType) Constraint {
	c.Deferrable = deferrable
	return c
}

func (c Constraint) WithNotValid() Constraint {
	c.NotValid = true
	return c
}

func (c Constraint) WithExcludeUsing(using string) Constraint {
	c.ExcludeUsing = using
	return c
}

func (c Constraint) WithExcludeWith(operators ...string) Constraint {
	c.ExcludeWith = append(c.ExcludeWith, operators...)
	return c
}
