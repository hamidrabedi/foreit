package schema

// CommonRelationMethods provides chainable methods that work for all relation builders
// Using generics, we implement common methods once for all relation types
type CommonRelationMethods[T any] struct {
	relation *Relation
	chainFn  func() T // Function that returns the concrete builder for chaining
}

// initCommonRelationMethods initializes CommonRelationMethods with the chain function
func initCommonRelationMethods[T any](cm *CommonRelationMethods[T], relation *Relation, chainFn func() T) {
	cm.relation = relation
	cm.chainFn = chainFn
}

// RelatedName sets the related name for the relation
func (c *CommonRelationMethods[T]) RelatedName(name string) T {
	c.relation.RelatedName = name
	return c.chainFn()
}

// RelatedQueryName sets the related query name for the relation
func (c *CommonRelationMethods[T]) RelatedQueryName(name string) T {
	c.relation.RelatedQueryName = name
	return c.chainFn()
}

// OnDelete sets the cascade behavior on delete
func (c *CommonRelationMethods[T]) OnDelete(cascade CascadeType) T {
	c.relation.OnDelete = cascade
	return c.chainFn()
}

// OnUpdate sets the cascade behavior on update
func (c *CommonRelationMethods[T]) OnUpdate(cascade CascadeType) T {
	c.relation.OnUpdate = cascade
	return c.chainFn()
}

// LimitChoicesTo limits the choices for the relation
func (c *CommonRelationMethods[T]) LimitChoicesTo(limit interface{}) T {
	c.relation.LimitChoicesTo = limit
	return c.chainFn()
}

// ParentLink marks the relation as a parent link
func (c *CommonRelationMethods[T]) ParentLink() T {
	c.relation.ParentLink = true
	return c.chainFn()
}

// CascadeOnDelete is a convenience method for OnDelete(CascadeCASCADE)
func (c *CommonRelationMethods[T]) CascadeOnDelete() T {
	c.relation.OnDelete = CascadeCASCADE
	return c.chainFn()
}
