package schema

// NewIndex creates a simple index on the given fields (name is optional).
func NewIndex(fields ...string) Index {
	return Index{
		Fields: fields,
		Type:   IndexTypeBTree,
	}
}

// IndexNamed creates a simple index with an explicit name.
func IndexNamed(name string, fields ...string) Index {
	return Index{
		Name:   name,
		Fields: fields,
		Type:   IndexTypeBTree,
	}
}

func (i Index) WithName(name string) Index {
	i.Name = name
	return i
}

func (i Index) WithFields(fields ...string) Index {
	i.Fields = append(i.Fields, fields...)
	return i
}

func (i Index) WithField(name string, order IndexOrder) Index {
	i.FieldSpecs = append(i.FieldSpecs, IndexField{
		Name:  name,
		Order: order,
	})
	i.Fields = append(i.Fields, name)
	return i
}

func (i Index) WithExpression(expr string) Index {
	i.Expressions = append(i.Expressions, expr)
	return i
}

func (i Index) WithUnique() Index {
	i.Unique = true
	return i
}

func (i Index) WithType(indexType IndexType) Index {
	i.Type = indexType
	return i
}

func (i Index) WithCondition(condition string) Index {
	i.Condition = condition
	return i
}

func (i Index) WithInclude(columns ...string) Index {
	i.Include = append(i.Include, columns...)
	return i
}

func (i Index) WithOpClasses(classes ...string) Index {
	i.OpClasses = append(i.OpClasses, classes...)
	return i
}

func (i Index) WithTablespace(tablespace string) Index {
	i.Tablespace = tablespace
	return i
}

func (i Index) WithOption(key, value string) Index {
	if i.Options == nil {
		i.Options = make(map[string]string)
	}
	i.Options[key] = value
	return i
}
