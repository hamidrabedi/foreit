package api

import (
	"reflect"

	"github.com/forgego/forge/schema"
)

// NonSerializableFields returns the names of fields on model whose Serialize
// flag is false (write-only values such as passwords). Both the field Name
// and its DB column name are returned when they differ. It returns nil when
// model does not implement schema.Schema.
func NonSerializableFields(model interface{}) []string {
	if model == nil {
		return nil
	}
	if s, ok := model.(schema.Schema); ok {
		return nonSerializableFromFields(s.Fields())
	}
	// Handle a non-pointer model whose pointer implements schema.Schema.
	v := reflect.ValueOf(model)
	if v.IsValid() && v.Kind() != reflect.Ptr {
		ptr := reflect.New(v.Type())
		ptr.Elem().Set(v)
		if s, ok := ptr.Interface().(schema.Schema); ok {
			return nonSerializableFromFields(s.Fields())
		}
	}
	return nil
}

func nonSerializableFromFields(fields []schema.Field) []string {
	var out []string
	for _, f := range fields {
		if f.Serialize {
			continue
		}
		out = append(out, f.Name)
		if f.DBColumn != "" && f.DBColumn != f.Name {
			out = append(out, f.DBColumn)
		}
	}
	return out
}

// stripExcludedFields removes ExcludeResponseFields keys from a serialized map.
func (vs *BaseViewSet) stripExcludedFields(m map[string]interface{}) map[string]interface{} {
	if m == nil || len(vs.ExcludeResponseFields) == 0 {
		return m
	}
	for _, k := range vs.ExcludeResponseFields {
		delete(m, k)
	}
	return m
}

// stripExcludedFieldsMany removes ExcludeResponseFields keys from each map.
func (vs *BaseViewSet) stripExcludedFieldsMany(items []map[string]interface{}) []map[string]interface{} {
	if len(vs.ExcludeResponseFields) == 0 {
		return items
	}
	for _, m := range items {
		vs.stripExcludedFields(m)
	}
	return items
}
