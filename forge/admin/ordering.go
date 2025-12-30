package admin

// Ordering represents an ordering specification
type Ordering[T any] struct {
	field     FieldExpr[T, interface{}]
	direction OrderDirection
}

// OrderDirection specifies sort direction
type OrderDirection string

const (
	OrderAsc  OrderDirection = "asc"
	OrderDesc OrderDirection = "desc"
)

// OrderBy creates an ordering specification
func OrderBy[T any](field FieldExpr[T, interface{}]) *OrderingBuilder[T] {
	return &OrderingBuilder[T]{
		field: field,
	}
}

// OrderingBuilder helps build ordering specifications
type OrderingBuilder[T any] struct {
	field FieldExpr[T, interface{}]
}

// Asc sets ascending order
func (b *OrderingBuilder[T]) Asc() Ordering[T] {
	return Ordering[T]{
		field:     b.field,
		direction: OrderAsc,
	}
}

// Desc sets descending order
func (b *OrderingBuilder[T]) Desc() Ordering[T] {
	return Ordering[T]{
		field:     b.field,
		direction: OrderDesc,
	}
}

// Field returns the field expression
func (o Ordering[T]) Field() FieldExpr[T, interface{}] {
	return o.field
}

// Direction returns the sort direction
func (o Ordering[T]) Direction() OrderDirection {
	return o.direction
}

// IsAscending returns true if ascending
func (o Ordering[T]) IsAscending() bool {
	return o.direction == OrderAsc
}

// IsDescending returns true if descending
func (o Ordering[T]) IsDescending() bool {
	return o.direction == OrderDesc
}
