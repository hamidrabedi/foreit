package schema

// IndexBuilder provides a fluent interface for building database indexes
type IndexBuilder struct {
	index Index
}

// NewIndex creates a new index builder with the given name
func NewIndex(name string) *IndexBuilder {
	return &IndexBuilder{
		index: Index{
			Name:   name,
			Fields: []string{},
			Type:   IndexTypeBTree,
		},
	}
}

// Fields adds one or more field names to the index
func (b *IndexBuilder) Fields(fields ...string) *IndexBuilder {
	b.index.Fields = append(b.index.Fields, fields...)
	return b
}

// Field adds a single field with optional ordering
func (b *IndexBuilder) Field(name string, order IndexOrder) *IndexBuilder {
	b.index.FieldSpecs = append(b.index.FieldSpecs, IndexField{
		Name:  name,
		Order: order,
	})
	// Also add to Fields
	b.index.Fields = append(b.index.Fields, name)
	return b
}

// Expression adds a functional index expression (e.g., "LOWER(title)")
func (b *IndexBuilder) Expression(expr string) *IndexBuilder {
	b.index.Expressions = append(b.index.Expressions, expr)
	return b
}

// Unique makes the index unique
func (b *IndexBuilder) Unique() *IndexBuilder {
	b.index.Unique = true
	return b
}

// Type sets the index type (btree, hash, gin, gist, brin)
func (b *IndexBuilder) Type(indexType IndexType) *IndexBuilder {
	b.index.Type = indexType
	return b
}

// Condition sets a WHERE condition for partial indexes
func (b *IndexBuilder) Condition(condition string) *IndexBuilder {
	b.index.Condition = condition
	return b
}

// Include adds columns to include in a covering index (PostgreSQL)
func (b *IndexBuilder) Include(columns ...string) *IndexBuilder {
	b.index.Include = append(b.index.Include, columns...)
	return b
}

// OpClasses sets operator classes for PostgreSQL indexes
func (b *IndexBuilder) OpClasses(classes ...string) *IndexBuilder {
	b.index.OpClasses = append(b.index.OpClasses, classes...)
	return b
}

// Tablespace sets the tablespace for the index (PostgreSQL)
func (b *IndexBuilder) Tablespace(tablespace string) *IndexBuilder {
	b.index.Tablespace = tablespace
	return b
}

// Option sets an index option (e.g., "fillfactor", "50")
func (b *IndexBuilder) Option(key, value string) *IndexBuilder {
	if b.index.Options == nil {
		b.index.Options = make(map[string]string)
	}
	b.index.Options[key] = value
	return b
}

// Build returns the built index
func (b *IndexBuilder) Build() Index {
	return b.index
}

// SimpleIndex creates a simple index with just field names
func SimpleIndex(name string, fields ...string) Index {
	return Index{
		Name:   name,
		Fields: fields,
		Type:   IndexTypeBTree,
	}
}

// UniqueIndex creates a unique index with field names
func UniqueIndex(name string, fields ...string) Index {
	return Index{
		Name:   name,
		Fields: fields,
		Unique: true,
		Type:   IndexTypeBTree,
	}
}

// PartialIndex creates a partial index with a condition
func PartialIndex(name string, condition string, fields ...string) Index {
	return Index{
		Name:      name,
		Fields:    fields,
		Condition: condition,
		Type:      IndexTypeBTree,
	}
}
