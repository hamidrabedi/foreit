package schema

// RelationType represents the type of a relationship
type RelationType int

const (
	RelationForeignKey RelationType = iota
	RelationOneToOne
	RelationManyToMany
)

func (t RelationType) String() string {
	switch t {
	case RelationForeignKey:
		return "foreign_key"
	case RelationOneToOne:
		return "one_to_one"
	case RelationManyToMany:
		return "many_to_many"
	default:
		return "unknown"
	}
}

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
	DBConstraint     bool           // Control FK constraint creation (default: true)
	ConstraintName   string         // Explicit FK constraint name
	Deferrable       DeferrableType // FK constraint deferrability
	Match            FKMatchType    // FK match type (FULL, PARTIAL, SIMPLE)
	Swappable        string         // Swappable model support (Django-specific)
}

// CustomRelation is the interface for custom relation types
type CustomRelation interface {
	GetName() string
	GetType() RelationType
	GetTarget() string
}

func newRelation(name, to string, relationType RelationType) Relation {
	return Relation{
		Name:         name,
		To:           to,
		Type:         relationType,
		DBConstraint: true,
	}
}

// ForeignKey creates a new ForeignKey relation
func ForeignKey(name, to string) Relation {
	return newRelation(name, to, RelationForeignKey)
}

// OneToOne creates a new OneToOne relation
func OneToOne(name, to string) Relation {
	return newRelation(name, to, RelationOneToOne)
}

// ManyToMany creates a new ManyToMany relation
func ManyToMany(name, to string) Relation {
	return newRelation(name, to, RelationManyToMany)
}

// ManyToOne is a helper that creates a ForeignKey relation
func ManyToOne(name, to, fkColumn string) Relation {
	return ForeignKey(fkColumn, to).WithRelatedName(name)
}

// OneToMany creates a reverse ForeignKey relation
func OneToMany(name, to, fkColumn string) Relation {
	_ = fkColumn
	return newRelation(name, to, RelationForeignKey)
}

// Relation chain methods

func (r Relation) WithRelatedName(name string) Relation {
	r.RelatedName = name
	return r
}

func (r Relation) WithRelatedQueryName(name string) Relation {
	r.RelatedQueryName = name
	return r
}

func (r Relation) WithOnDelete(cascade CascadeType) Relation {
	r.OnDelete = cascade
	return r
}

func (r Relation) WithOnUpdate(cascade CascadeType) Relation {
	r.OnUpdate = cascade
	return r
}

func (r Relation) WithLimitChoicesTo(limit interface{}) Relation {
	r.LimitChoicesTo = limit
	return r
}

func (r Relation) WithParentLink() Relation {
	r.ParentLink = true
	return r
}

func (r Relation) WithCascadeOnDelete() Relation {
	r.OnDelete = CascadeCASCADE
	return r
}

func (r Relation) WithRequired() Relation {
	return r
}

func (r Relation) WithOptional() Relation {
	return r
}

func (r Relation) WithDBConstraint(enabled bool) Relation {
	r.DBConstraint = enabled
	return r
}

func (r Relation) WithConstraintName(name string) Relation {
	r.ConstraintName = name
	return r
}

func (r Relation) WithDeferrable(deferrable DeferrableType) Relation {
	r.Deferrable = deferrable
	return r
}

func (r Relation) WithMatch(matchType FKMatchType) Relation {
	r.Match = matchType
	return r
}

func (r Relation) WithSwappable(swappable string) Relation {
	r.Swappable = swappable
	return r
}

func (r Relation) WithThrough(throughTable, localColumn, remoteColumn string) Relation {
	_ = localColumn
	_ = remoteColumn
	r.Through = throughTable
	return r
}

func (r Relation) WithThroughTable(tableName string) Relation {
	r.Through = tableName
	return r
}

func (r Relation) WithSymmetrical() Relation {
	r.Symmetrical = true
	return r
}
