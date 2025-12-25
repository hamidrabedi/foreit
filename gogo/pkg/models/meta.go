package models

// Meta represents model metadata options (like Django's Meta class)
type Meta struct {
	// TableName is the database table name
	TableName string
	
	// Ordering specifies default ordering
	Ordering []string
	
	// Indexes specifies database indexes
	Indexes []Index
	
	// UniqueTogether specifies unique constraints
	UniqueTogether [][]string
	
	// VerboseName is the human-readable name
	VerboseName string
	
	// VerboseNamePlural is the plural human-readable name
	VerboseNamePlural string
	
	// DBTable is an alias for TableName
	DBTable string
	
	// Managed controls whether migrations manage this model
	Managed bool
	
	// Abstract indicates if this is an abstract base model
	Abstract bool
	
	// Proxy indicates if this is a proxy model
	Proxy bool
	
	// DefaultSelectRelated specifies default fields to eagerly load
	DefaultSelectRelated []string
	
	// DefaultPrefetchRelated specifies default relationships to prefetch
	DefaultPrefetchRelated []string
	
	// Constraints specifies database constraints
	Constraints []Constraint
	
	// Permissions specifies model-level permissions
	Permissions []string
	
	// DefaultManagerName specifies the default manager name
	DefaultManagerName string
}

// Index represents a database index
type Index struct {
	Fields []string
	Name   string
	Unique bool
	Type   string // BTREE, HASH, GIN, etc.
}

// Constraint represents a database constraint
type Constraint struct {
	Name      string
	Type      string // CHECK, UNIQUE, FOREIGN KEY, etc.
	Fields    []string
	Condition string // For CHECK constraints
}

// ModelMeta interface for models that have metadata
type ModelMeta interface {
	Meta() *Meta
}

// GetMeta retrieves metadata from a model
func GetMeta(model Model) *Meta {
	if metaModel, ok := model.(ModelMeta); ok {
		return metaModel.Meta()
	}
	return &Meta{
		Managed: true,
	}
}

// DefaultMeta returns default metadata
func DefaultMeta() *Meta {
	return &Meta{
		Managed: true,
		Ordering: []string{"-id"},
	}
}

