package schema

// BaseFieldBuilder provides common field building methods that all field builders share.
// This reduces code duplication and provides a consistent API across all field types.
type BaseFieldBuilder struct {
	field Field
}

// DBCollation sets the character set collation for the field (PostgreSQL, MySQL)
func (b *BaseFieldBuilder) DBCollation(collation string) *BaseFieldBuilder {
	b.field.DBCollation = collation
	return b
}

// DBComment sets the column comment/description
func (b *BaseFieldBuilder) DBComment(comment string) *BaseFieldBuilder {
	b.field.DBComment = comment
	return b
}

// DBTablespace sets the tablespace for field storage (PostgreSQL)
func (b *BaseFieldBuilder) DBTablespace(tablespace string) *BaseFieldBuilder {
	b.field.DBTablespace = tablespace
	return b
}

// DBType sets an explicit database column type override
func (b *BaseFieldBuilder) DBType(dbType string) *BaseFieldBuilder {
	b.field.DBType = dbType
	return b
}

// DBDefault sets a database-level default (SQL expression)
func (b *BaseFieldBuilder) DBDefault(expr string) *BaseFieldBuilder {
	b.field.DBDefault = expr
	return b
}

// ErrorMessages sets custom error messages per validation type
func (b *BaseFieldBuilder) ErrorMessages(messages map[string]string) *BaseFieldBuilder {
	b.field.ErrorMessages = messages
	return b
}

// HelpText sets the help text for the field
func (b *BaseFieldBuilder) HelpText(text string) *BaseFieldBuilder {
	b.field.HelpText = text
	return b
}

// VerboseName sets the human-readable name for the field
func (b *BaseFieldBuilder) VerboseName(name string) *BaseFieldBuilder {
	b.field.VerboseName = name
	return b
}

// DBColumn sets a custom database column name
func (b *BaseFieldBuilder) DBColumn(name string) *BaseFieldBuilder {
	b.field.DBColumn = name
	return b
}

// DBIndex creates a database index on this field
func (b *BaseFieldBuilder) DBIndex() *BaseFieldBuilder {
	b.field.DBIndex = true
	return b
}

// Required marks the field as required (NOT NULL)
func (b *BaseFieldBuilder) Required() *BaseFieldBuilder {
	b.field.Required = true
	return b
}

// Optional marks the field as optional (allows NULL)
func (b *BaseFieldBuilder) Optional() *BaseFieldBuilder {
	b.field.Required = false
	return b
}

// Unique marks the field as unique
func (b *BaseFieldBuilder) Unique() *BaseFieldBuilder {
	b.field.Unique = true
	return b
}

// Primary marks the field as primary key
func (b *BaseFieldBuilder) Primary() *BaseFieldBuilder {
	b.field.PrimaryKey = true
	return b
}

// Blank allows the field to be blank (for strings, allows empty string)
func (b *BaseFieldBuilder) Blank() *BaseFieldBuilder {
	b.field.Blank = true
	return b
}

// Editable controls whether the field is editable in forms/admin
func (b *BaseFieldBuilder) Editable(editable bool) *BaseFieldBuilder {
	b.field.Editable = editable
	return b
}

// Serialize controls whether the field is serialized in forms/API
func (b *BaseFieldBuilder) Serialize(serialize bool) *BaseFieldBuilder {
	b.field.Serialize = serialize
	return b
}

// Validators adds one or more validators to the field
func (b *BaseFieldBuilder) Validators(validators ...Validator) *BaseFieldBuilder {
	if b.field.Validators == nil {
		b.field.Validators = []Validator{}
	}
	b.field.Validators = append(b.field.Validators, validators...)
	return b
}

// ValidationTag sets a validation tag (e.g., for struct tags)
func (b *BaseFieldBuilder) ValidationTag(tag string) *BaseFieldBuilder {
	b.field.ValidationTag = tag
	return b
}

// GeneratedColumn marks the field as a generated column with the given expression
// If stored is true, the column is STORED; otherwise it's VIRTUAL
func (b *BaseFieldBuilder) GeneratedColumn(expression string, stored bool) *BaseFieldBuilder {
	b.field.Generated = true
	b.field.GeneratedExpr = expression
	b.field.IsStored = stored
	return b
}

// UniqueForDate makes the field unique for each date value of the specified date field
func (b *BaseFieldBuilder) UniqueForDate(dateField string) *BaseFieldBuilder {
	b.field.UniqueForDate = dateField
	return b
}

// UniqueForMonth makes the field unique for each month value of the specified date field
func (b *BaseFieldBuilder) UniqueForMonth(dateField string) *BaseFieldBuilder {
	b.field.UniqueForMonth = dateField
	return b
}

// UniqueForYear makes the field unique for each year value of the specified date field
func (b *BaseFieldBuilder) UniqueForYear(dateField string) *BaseFieldBuilder {
	b.field.UniqueForYear = dateField
	return b
}

// GetField returns the underlying field (for internal use)
func (b *BaseFieldBuilder) GetField() *Field {
	return &b.field
}

