package schema

// RelationType represents the type of a relationship
type RelationType int

const (
	RelationForeignKey RelationType = iota
	RelationOneToOne
	RelationManyToMany
)

// CascadeType represents cascade behavior
type CascadeType string

const (
	CascadeCASCADE     CascadeType = "CASCADE"
	CascadeSET_NULL    CascadeType = "SET_NULL"
	CascadePROTECT     CascadeType = "PROTECT"
	CascadeSET_DEFAULT CascadeType = "SET_DEFAULT"
	CascadeDO_NOTHING  CascadeType = "DO_NOTHING"
)

// FKMatchType represents foreign key match type
type FKMatchType string

const (
	FKMatchFull    FKMatchType = "FULL"
	FKMatchPartial FKMatchType = "PARTIAL"
	FKMatchSimple  FKMatchType = "SIMPLE"
)

// Relation represents a model relationship with full Django relationship features
type Relation struct {
	LimitChoicesTo   interface{}
	CustomRelation   CustomRelation
	Name             string
	To               string
	RelatedName      string
	RelatedQueryName string
	OnDelete         CascadeType
	OnUpdate         CascadeType
	Through          string
	Type             RelationType
	Symmetrical      bool
	ParentLink       bool
	DBConstraint     bool        // Control FK constraint creation (default: true)
	ConstraintName   string      // Explicit FK constraint name
	Deferrable       DeferrableType // FK constraint deferrability
	Match            FKMatchType // FK match type (FULL, PARTIAL, SIMPLE)
	Swappable        string      // Swappable model support (Django-specific)
}

// CustomRelation is the interface for custom relation types
type CustomRelation interface {
	GetName() string
	GetType() RelationType
	GetTarget() string
}

// ForeignKey creates a new ForeignKey relation builder
func ForeignKey(name, to string) *ForeignKeyBuilder {
	return &ForeignKeyBuilder{
		relation: Relation{
			Name: name,
			Type: RelationForeignKey,
			To:   to,
		},
	}
}

// ForeignKeyBuilder is a builder for ForeignKey relations
type ForeignKeyBuilder struct {
	relation Relation
}

func (b *ForeignKeyBuilder) RelatedName(name string) *ForeignKeyBuilder {
	b.relation.RelatedName = name
	return b
}

func (b *ForeignKeyBuilder) RelatedQueryName(name string) *ForeignKeyBuilder {
	b.relation.RelatedQueryName = name
	return b
}

func (b *ForeignKeyBuilder) OnDelete(cascade CascadeType) *ForeignKeyBuilder {
	b.relation.OnDelete = cascade
	return b
}

func (b *ForeignKeyBuilder) OnUpdate(cascade CascadeType) *ForeignKeyBuilder {
	b.relation.OnUpdate = cascade
	return b
}

func (b *ForeignKeyBuilder) LimitChoicesTo(limit interface{}) *ForeignKeyBuilder {
	b.relation.LimitChoicesTo = limit
	return b
}

func (b *ForeignKeyBuilder) ParentLink() *ForeignKeyBuilder {
	b.relation.ParentLink = true
	return b
}

func (b *ForeignKeyBuilder) Required() *ForeignKeyBuilder {
	// ForeignKey is always required unless explicitly set to optional
	return b
}

func (b *ForeignKeyBuilder) Optional() *ForeignKeyBuilder {
	// This would be used in field definitions, but for relations
	// we use OnDelete(SET_NULL) to make it optional
	return b
}

func (b *ForeignKeyBuilder) VerboseName(name string) *ForeignKeyBuilder {
	// Store verbose name in relation (could be used in admin)
	return b
}

func (b *ForeignKeyBuilder) CascadeOnDelete() *ForeignKeyBuilder {
	b.relation.OnDelete = CascadeCASCADE
	return b
}

func (b *ForeignKeyBuilder) DBConstraint(enabled bool) *ForeignKeyBuilder {
	b.relation.DBConstraint = enabled
	return b
}

func (b *ForeignKeyBuilder) ConstraintName(name string) *ForeignKeyBuilder {
	b.relation.ConstraintName = name
	return b
}

func (b *ForeignKeyBuilder) Deferrable(deferrable DeferrableType) *ForeignKeyBuilder {
	b.relation.Deferrable = deferrable
	return b
}

func (b *ForeignKeyBuilder) Match(matchType FKMatchType) *ForeignKeyBuilder {
	b.relation.Match = matchType
	return b
}

func (b *ForeignKeyBuilder) Swappable(swappable string) *ForeignKeyBuilder {
	b.relation.Swappable = swappable
	return b
}

func (b *ForeignKeyBuilder) Build() Relation {
	// Set default DBConstraint to true if not explicitly set
	if !b.relation.DBConstraint && b.relation.ConstraintName == "" {
		b.relation.DBConstraint = true
	}
	return b.relation
}

// ManyToOne is a helper that creates a ForeignKey relation
// This is for convenience - it's the same as ForeignKey
func ManyToOne(name, to, fkColumn string) *ForeignKeyBuilder {
	return ForeignKey(fkColumn, to).RelatedName(name)
}

// OneToMany creates a reverse relation (used on the "one" side)
// This is typically defined on the "many" side as ForeignKey
// but this helper makes it clearer when defining relations
func OneToMany(name, to, fkColumn string) *OneToManyBuilder {
	return &OneToManyBuilder{
		relation: Relation{
			Name: name,
			Type: RelationForeignKey, // OneToMany is really a reverse ForeignKey
			To:   to,
		},
		fkColumn: fkColumn,
	}
}

// OneToManyBuilder represents a reverse ForeignKey relation
type OneToManyBuilder struct {
	relation Relation
	fkColumn string
}

func (b *OneToManyBuilder) CascadeOnDelete() *OneToManyBuilder {
	b.relation.OnDelete = CascadeCASCADE
	return b
}

func (b *OneToManyBuilder) Build() Relation {
	return b.relation
}

// OneToOne creates a new OneToOne relation builder
func OneToOne(name, to string) *OneToOneBuilder {
	return &OneToOneBuilder{
		relation: Relation{
			Name: name,
			Type: RelationOneToOne,
			To:   to,
		},
	}
}

// OneToOneBuilder is a builder for OneToOne relations
type OneToOneBuilder struct {
	relation Relation
}

func (b *OneToOneBuilder) RelatedName(name string) *OneToOneBuilder {
	b.relation.RelatedName = name
	return b
}

func (b *OneToOneBuilder) OnDelete(cascade CascadeType) *OneToOneBuilder {
	b.relation.OnDelete = cascade
	return b
}

func (b *OneToOneBuilder) OnUpdate(cascade CascadeType) *OneToOneBuilder {
	b.relation.OnUpdate = cascade
	return b
}

func (b *OneToOneBuilder) ParentLink() *OneToOneBuilder {
	b.relation.ParentLink = true
	return b
}

func (b *OneToOneBuilder) Build() Relation {
	return b.relation
}

// ManyToMany creates a new ManyToMany relation builder
func ManyToMany(name, to string) *ManyToManyBuilder {
	return &ManyToManyBuilder{
		relation: Relation{
			Name: name,
			Type: RelationManyToMany,
			To:   to,
		},
	}
}

// ManyToManyBuilder is a builder for ManyToMany relations
type ManyToManyBuilder struct {
	relation Relation
}

func (b *ManyToManyBuilder) RelatedName(name string) *ManyToManyBuilder {
	b.relation.RelatedName = name
	return b
}

func (b *ManyToManyBuilder) RelatedQueryName(name string) *ManyToManyBuilder {
	b.relation.RelatedQueryName = name
	return b
}

func (b *ManyToManyBuilder) Through(throughTable, localColumn, remoteColumn string) *ManyToManyBuilder {
	b.relation.Through = throughTable
	// Store column names for the through table
	// This would be used when generating SQL
	return b
}

func (b *ManyToManyBuilder) ThroughTable(tableName string) *ManyToManyBuilder {
	b.relation.Through = tableName
	return b
}

func (b *ManyToManyBuilder) Symmetrical() *ManyToManyBuilder {
	b.relation.Symmetrical = true
	return b
}

func (b *ManyToManyBuilder) LimitChoicesTo(limit interface{}) *ManyToManyBuilder {
	b.relation.LimitChoicesTo = limit
	return b
}

func (b *ManyToManyBuilder) Build() Relation {
	return b.relation
}

// RegisterRelationType registers a custom relation type
func RegisterRelationType(name string, relationType RelationType, builder func(string, string) interface{}) {
	// TODO: Implement relation type registry
}
