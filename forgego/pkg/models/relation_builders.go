package models

// ManyToOneBuilder builds a ManyToOne relation descriptor.
// Note: Field types are type-erased to avoid import cycles.
type ManyToOneBuilder[T, R any, FK comparable] struct {
	name      string
	fromField interface{} // field.Field[T, FK] - type-erased to avoid import cycle
	toField   interface{} // field.Field[R, FK] - type-erased to avoid import cycle
	onDelete  string
	onUpdate  string
	backRef   string
}

// ManyToOne creates a new ManyToOne builder.
func ManyToOne[T, R any, FK comparable](name string) *ManyToOneBuilder[T, R, FK] {
	return &ManyToOneBuilder[T, R, FK]{
		name: name,
	}
}

// From sets the foreign key field on the owner model.
// Accepts any type that implements the Field interface (type-erased to avoid import cycles).
func (b *ManyToOneBuilder[T, R, FK]) From(field interface{}) *ManyToOneBuilder[T, R, FK] {
	b.fromField = field
	return b
}

// To sets the primary key field on the related model.
// Accepts any type that implements the Field interface (type-erased to avoid import cycles).
func (b *ManyToOneBuilder[T, R, FK]) To(field interface{}) *ManyToOneBuilder[T, R, FK] {
	b.toField = field
	return b
}

// OnDelete sets the ON DELETE action.
func (b *ManyToOneBuilder[T, R, FK]) OnDelete(action string) *ManyToOneBuilder[T, R, FK] {
	b.onDelete = action
	return b
}

// OnUpdate sets the ON UPDATE action.
func (b *ManyToOneBuilder[T, R, FK]) OnUpdate(action string) *ManyToOneBuilder[T, R, FK] {
	b.onUpdate = action
	return b
}

// BackRef sets the back reference name.
func (b *ManyToOneBuilder[T, R, FK]) BackRef(name string) *ManyToOneBuilder[T, R, FK] {
	b.backRef = name
	return b
}

// Build creates the RelationDescriptor.
func (b *ManyToOneBuilder[T, R, FK]) Build() RelationDescriptor {
	var zero R
	relatedModelName := getTypeName(zero)
	return &manyToOneDescriptor{
		name:         b.name,
		fromField:    b.fromField,
		toField:      b.toField,
		onDelete:     b.onDelete,
		onUpdate:     b.onUpdate,
		backRef:      b.backRef,
		relatedModel: relatedModelName,
	}
}

// manyToOneDescriptor implements RelationDescriptor for ManyToOne.
type manyToOneDescriptor struct {
	name         string
	fromField    interface{} // field.Field[T, FK] - type-erased
	toField      interface{} // field.Field[R, FK] - type-erased
	onDelete     string
	onUpdate     string
	backRef      string
	relatedModel string // Store related model name for easier access
}

func (d *manyToOneDescriptor) GetName() string {
	return d.name
}

func (d *manyToOneDescriptor) GetType() RelationType {
	return RelationTypeForeignKey // Maps to "foreignkey" in existing system
}

// OneToManyBuilder builds a OneToMany relation descriptor.
// Note: Field types are type-erased to avoid import cycles.
type OneToManyBuilder[T, R any, FK comparable] struct {
	name      string
	fromField interface{} // field.Field[T, FK] - type-erased to avoid import cycle
	toField   interface{} // field.Field[R, FK] - type-erased to avoid import cycle
	onDelete  string
	backRef   string
}

// OneToManyRelation creates a new OneToMany builder.
// Note: This is the Ent-inspired builder version. The old OneToMany in relationships.go is kept for compatibility.
func OneToManyRelation[T, R any, FK comparable](name string) *OneToManyBuilder[T, R, FK] {
	return &OneToManyBuilder[T, R, FK]{
		name: name,
	}
}

// From sets the primary key field on the owner model.
// Accepts any type that implements the Field interface (type-erased to avoid import cycles).
func (b *OneToManyBuilder[T, R, FK]) From(field interface{}) *OneToManyBuilder[T, R, FK] {
	b.fromField = field
	return b
}

// To sets the foreign key field on the related model.
// Accepts any type that implements the Field interface (type-erased to avoid import cycles).
func (b *OneToManyBuilder[T, R, FK]) To(field interface{}) *OneToManyBuilder[T, R, FK] {
	b.toField = field
	return b
}

// OnDelete sets the ON DELETE action.
func (b *OneToManyBuilder[T, R, FK]) OnDelete(action string) *OneToManyBuilder[T, R, FK] {
	b.onDelete = action
	return b
}

// BackRef sets the back reference name.
func (b *OneToManyBuilder[T, R, FK]) BackRef(name string) *OneToManyBuilder[T, R, FK] {
	b.backRef = name
	return b
}

// Build creates the RelationDescriptor.
func (b *OneToManyBuilder[T, R, FK]) Build() RelationDescriptor {
	var zero R
	relatedModelName := getTypeName(zero)
	return &oneToManyDescriptor{
		name:         b.name,
		fromField:    b.fromField,
		toField:      b.toField,
		onDelete:     b.onDelete,
		backRef:      b.backRef,
		relatedModel: relatedModelName,
	}
}

// oneToManyDescriptor implements RelationDescriptor for OneToMany.
type oneToManyDescriptor struct {
	name         string
	fromField    interface{} // field.Field[T, FK] - type-erased
	toField      interface{} // field.Field[R, FK] - type-erased
	onDelete     string
	backRef      string
	relatedModel string // Store related model name for easier access
}

func (d *oneToManyDescriptor) GetName() string {
	return d.name
}

func (d *oneToManyDescriptor) GetType() RelationType {
	return RelationTypeOneToMany
}

