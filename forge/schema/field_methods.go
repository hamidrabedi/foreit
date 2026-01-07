package schema

import "github.com/google/uuid"

func newField(name string, fieldType FieldType) Field {
	return Field{
		Name:      name,
		Type:      fieldType,
		Editable:  true,
		Serialize: true,
	}
}

func (f Field) WithDefault(value interface{}) Field {
	f.Default = value
	return f
}

func (f Field) WithDefaultUUID(value uuid.UUID) Field {
	f.Default = value.String()
	return f
}

func (f Field) WithDefaultNewUUID() Field {
	f.Default = uuid.NewString()
	return f
}

func (f Field) WithRequired() Field {
	f.Required = true
	return f
}

func (f Field) WithOptional() Field {
	f.Required = false
	return f
}

func (f Field) WithUnique() Field {
	f.Unique = true
	return f
}

func (f Field) WithPrimary() Field {
	f.PrimaryKey = true
	return f
}

func (f Field) WithBlank() Field {
	f.Blank = true
	return f
}

func (f Field) WithEditable(editable bool) Field {
	f.Editable = editable
	return f
}

func (f Field) WithSerialize(serialize bool) Field {
	f.Serialize = serialize
	return f
}

func (f Field) WithWriteOnly() Field {
	f.Serialize = false
	return f
}

func (f Field) WithValidators(validators ...Validator) Field {
	f.Validators = append(f.Validators, validators...)
	return f
}

func (f Field) WithValidationTag(tag string) Field {
	f.ValidationTag = tag
	return f
}

func (f Field) WithErrorMessages(messages map[string]string) Field {
	f.ErrorMessages = messages
	return f
}

func (f Field) WithHelpText(text string) Field {
	f.HelpText = text
	return f
}

func (f Field) WithVerboseName(name string) Field {
	f.VerboseName = name
	return f
}

func (f Field) WithUniqueForDate(dateField string) Field {
	f.UniqueForDate = dateField
	return f
}

func (f Field) WithUniqueForMonth(dateField string) Field {
	f.UniqueForMonth = dateField
	return f
}

func (f Field) WithUniqueForYear(dateField string) Field {
	f.UniqueForYear = dateField
	return f
}

func (f Field) WithMaxLength(n int) Field {
	f.MaxLength = intPtr(n)
	return f
}

func (f Field) WithMinLength(n int) Field {
	f.MinLength = intPtr(n)
	return f
}

func (f Field) WithMaxValue(val float64) Field {
	f.MaxValue = float64Ptr(val)
	return f
}

func (f Field) WithMinValue(val float64) Field {
	f.MinValue = float64Ptr(val)
	return f
}

func (f Field) WithMaxDigits(digits int) Field {
	f.MaxDigits = intPtr(digits)
	return f
}

func (f Field) WithDecimalPlaces(places int) Field {
	f.DecimalPlaces = intPtr(places)
	return f
}

func (f Field) WithChoices(choices ...Choice) Field {
	f.Choices = choices
	return f
}

func (f Field) WithChoicesFromPairs(pairs ...string) Field {
	f.Choices = WithChoices(pairs...)
	return f
}

func (f Field) WithAutoIncrement() Field {
	f.AutoIncrement = true
	return f
}

func (f Field) WithAutoNow() Field {
	f.AutoNow = true
	return f
}

func (f Field) WithAutoNowAdd() Field {
	f.AutoNowAdd = true
	return f
}

func (f Field) WithDBColumn(name string) Field {
	f.DBColumn = name
	return f
}

func (f Field) WithDBType(dbType string) Field {
	f.DBType = dbType
	return f
}

func (f Field) WithDBCollation(collation string) Field {
	f.DBCollation = collation
	return f
}

func (f Field) WithDBComment(comment string) Field {
	f.DBComment = comment
	return f
}

func (f Field) WithDBTablespace(tablespace string) Field {
	f.DBTablespace = tablespace
	return f
}

func (f Field) WithDBDefault(expr string) Field {
	f.DBDefault = expr
	return f
}

func (f Field) WithDBIndex() Field {
	f.DBIndex = true
	return f
}

func (f Field) WithGeneratedColumn(expression string, stored bool) Field {
	f.Generated = true
	f.GeneratedExpr = expression
	f.IsStored = stored
	return f
}

func intPtr(value int) *int {
	v := value
	return &v
}

func float64Ptr(value float64) *float64 {
	v := value
	return &v
}
