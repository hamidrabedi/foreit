package admin

import (
	"fmt"
	"reflect"
	"strings"
	"time"

	"entgo.io/ent"
)

// SchemaIntrospector extracts metadata from Ent schemas at runtime
type SchemaIntrospector struct {
	registry *Registry
}

// NewSchemaIntrospector creates a new schema introspector
func NewSchemaIntrospector(registry *Registry) *SchemaIntrospector {
	return &SchemaIntrospector{
		registry: registry,
	}
}

// Introspect extracts metadata from an Ent schema
func (s *SchemaIntrospector) Introspect(schema ent.Schema) (*ModelMeta, error) {
	// Get the schema type
	schemaType := reflect.TypeOf(schema)
	if schemaType.Kind() == reflect.Ptr {
		schemaType = schemaType.Elem()
	}

	modelName := schemaType.Name()
	meta, err := s.registry.GetModel(modelName)
	if err != nil {
		return nil, fmt.Errorf("model %s not found in registry: %w", modelName, err)
	}

	// Extract fields from schema
	fields := schema.Fields()
	fieldMetas := make([]FieldMeta, 0, len(fields))

	// Extract field metadata
	for i, f := range fields {
		fieldMeta := s.extractFieldMeta(f, i)
		fieldMetas = append(fieldMetas, fieldMeta)
	}

	meta.Fields = fieldMetas
	meta.Schema = schema

	// Extract edges (relationships) from schema
	edges := schema.Edges()
	relationshipMetas := make([]RelationshipMeta, 0, len(edges))
	
	for _, edge := range edges {
		relMeta := s.extractRelationshipMeta(edge)
		if relMeta != nil {
			relationshipMetas = append(relationshipMetas, *relMeta)
		}
	}
	
	meta.Relationships = relationshipMetas

	// Set default table name if not set
	if meta.TableName == "" {
		meta.TableName = s.inferTableName(modelName)
	}

	return meta, nil
}

// extractFieldMeta extracts metadata from an Ent field
func (s *SchemaIntrospector) extractFieldMeta(f ent.Field, index int) FieldMeta {
	// Try to get field descriptor using reflection
	// Ent fields are typically descriptors that we can inspect
	fieldValue := reflect.ValueOf(f)
	
	// Default field name
	fieldName := fmt.Sprintf("field_%d", index)
	
	// Try to extract field name from the descriptor
	// Ent field descriptors typically have a Name() method or Name field
	if fieldValue.IsValid() {
		// Try to call Name() method
		nameMethod := fieldValue.MethodByName("Name")
		if nameMethod.IsValid() && nameMethod.Type().NumIn() == 0 {
			if results := nameMethod.Call(nil); len(results) > 0 {
				if name, ok := results[0].Interface().(string); ok {
					fieldName = name
				}
			}
		}
		
		// Try to get Name field directly
		if fieldName == fmt.Sprintf("field_%d", index) {
			nameField := fieldValue.FieldByName("Name")
			if nameField.IsValid() && nameField.Kind() == reflect.String {
				fieldName = nameField.String()
			}
		}
	}
	
	// Extract field type information
	fieldTypeInfo := s.extractFieldTypeInfo(f, fieldValue)
	
	// Build field metadata
	meta := FieldMeta{
		Name:       fieldName,
		Label:      s.formatLabel(fieldName),
		Type:       fieldTypeInfo.AdminType,
		DBType:     fieldTypeInfo.DBType,
		Required:   fieldTypeInfo.Required,
		Unique:     fieldTypeInfo.Unique,
		Filterable: true, // Default to filterable
		Sortable:   true, // Default to sortable
		Searchable: fieldTypeInfo.AdminType == FieldTypeText || fieldTypeInfo.AdminType == FieldTypeTextarea,
	}
	
	// Extract default value if available
	if fieldTypeInfo.Default != nil {
		meta.Default = fieldTypeInfo.Default
	}
	
	// Extract enum choices if available
	if len(fieldTypeInfo.Choices) > 0 {
		meta.Choices = fieldTypeInfo.Choices
		if meta.Type == FieldTypeText {
			meta.Type = FieldTypeSelect
		}
	}
	
	// Check if field should be read-only (e.g., created_at, updated_at, id)
	if strings.Contains(fieldName, "_at") || strings.Contains(fieldName, "_on") || fieldName == "id" {
		meta.ReadOnly = true
	}
	
	return meta
}

// FieldTypeInfo contains extracted type information
type FieldTypeInfo struct {
	AdminType FieldType
	DBType    string
	Required  bool
	Unique    bool
	Default   interface{}
	Choices   []Choice
}

// extractFieldTypeInfo extracts type information from an Ent field
func (s *SchemaIntrospector) extractFieldTypeInfo(f ent.Field, fieldValue reflect.Value) FieldTypeInfo {
	info := FieldTypeInfo{
		AdminType: FieldTypeText,
		DBType:    "TEXT",
		Required:  true,
		Unique:    false,
	}
	
	// Try to get descriptor information
	// Ent fields implement the Field interface which has Descriptor() method
	descriptorMethod := fieldValue.MethodByName("Descriptor")
	if descriptorMethod.IsValid() {
		results := descriptorMethod.Call(nil)
		if len(results) > 0 {
			descriptor := results[0]
			if descriptor.IsValid() {
				info = s.extractFromDescriptor(descriptor)
			}
		}
	}
	
	// Fallback: Try to infer from field type
	if info.AdminType == FieldTypeText {
		fieldType := reflect.TypeOf(f)
		info = s.inferTypeFromReflection(fieldType, info)
	}
	
	return info
}

// extractFromDescriptor extracts information from an Ent field descriptor
func (s *SchemaIntrospector) extractFromDescriptor(descriptor reflect.Value) FieldTypeInfo {
	info := FieldTypeInfo{
		AdminType: FieldTypeText,
		DBType:    "TEXT",
		Required:  true,
		Unique:    false,
	}
	
	if !descriptor.IsValid() {
		return info
	}
	
	// Try to get Info field (field.Info)
	infoField := descriptor.FieldByName("Info")
	if infoField.IsValid() {
		infoType := infoField.FieldByName("Type")
		if infoType.IsValid() {
			info = s.mapReflectionTypeToFieldType(infoType, info)
		}
	}
	
	// Try to get Optional field
	optionalField := descriptor.FieldByName("Optional")
	if optionalField.IsValid() && optionalField.Kind() == reflect.Bool {
		info.Required = !optionalField.Bool()
	}
	
	// Try to get Unique field
	uniqueField := descriptor.FieldByName("Unique")
	if uniqueField.IsValid() && uniqueField.Kind() == reflect.Bool {
		info.Unique = uniqueField.Bool()
	}
	
	// Try to get Default field
	defaultField := descriptor.FieldByName("Default")
	if defaultField.IsValid() && defaultField.IsValid() {
		if !defaultField.IsNil() {
			info.Default = defaultField.Interface()
		}
	}
	
	// Try to get Enum field
	enumField := descriptor.FieldByName("Enum")
	if enumField.IsValid() {
		if enumField.Kind() == reflect.Slice {
			choices := make([]Choice, 0, enumField.Len())
			for i := 0; i < enumField.Len(); i++ {
				elem := enumField.Index(i)
				if elem.Kind() == reflect.String {
					val := elem.String()
					choices = append(choices, Choice{
						Value: val,
						Label: s.formatLabel(val),
					})
				}
			}
			info.Choices = choices
		}
	}
	
	// Try to get Size field for strings
	sizeField := descriptor.FieldByName("Size")
	if sizeField.IsValid() && sizeField.Kind() == reflect.Int {
		size := int(sizeField.Int())
		if size > 0 && size <= 255 {
			info.DBType = fmt.Sprintf("VARCHAR(%d)", size)
		} else if size > 255 {
			info.AdminType = FieldTypeTextarea
		}
	}
	
	return info
}

// inferTypeFromReflection infers field type from reflection
func (s *SchemaIntrospector) inferTypeFromReflection(fieldType reflect.Type, info FieldTypeInfo) FieldTypeInfo {
	// Check if it's a pointer to a field descriptor
	if fieldType.Kind() == reflect.Ptr {
		fieldType = fieldType.Elem()
	}
	
	// Try to match common Ent field types
	typeName := fieldType.String()
	
	// Check for common patterns
	if strings.Contains(typeName, "String") {
		info.AdminType = FieldTypeText
		info.DBType = "TEXT"
	} else if strings.Contains(typeName, "Int") {
		info.AdminType = FieldTypeNumber
		info.DBType = "INTEGER"
	} else if strings.Contains(typeName, "Float") {
		info.AdminType = FieldTypeNumber
		info.DBType = "REAL"
	} else if strings.Contains(typeName, "Bool") {
		info.AdminType = FieldTypeBoolean
		info.DBType = "BOOLEAN"
	} else if strings.Contains(typeName, "Time") {
		info.AdminType = FieldTypeDateTime
		info.DBType = "TIMESTAMP"
	}
	
	return info
}

// mapReflectionTypeToFieldType maps a reflection type to admin field type
func (s *SchemaIntrospector) mapReflectionTypeToFieldType(typeValue reflect.Value, info FieldTypeInfo) FieldTypeInfo {
	if !typeValue.IsValid() {
		return info
	}
	
	// Get the actual type
	var actualType reflect.Type
	if typeValue.Kind() == reflect.Interface {
		if typeValue.IsNil() {
			return info
		}
		actualType = typeValue.Elem().Type()
	} else {
		actualType = typeValue.Type()
	}
	
	// Handle time.Time specifically
	if actualType == reflect.TypeOf(time.Time{}) {
		info.AdminType = FieldTypeDateTime
		info.DBType = "TIMESTAMP"
		return info
	}
	
	// Map based on kind
	switch actualType.Kind() {
	case reflect.String:
		info.AdminType = FieldTypeText
		info.DBType = "TEXT"
		
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32:
		info.AdminType = FieldTypeNumber
		info.DBType = "INTEGER"
		
	case reflect.Int64:
		info.AdminType = FieldTypeNumber
		info.DBType = "BIGINT"
		
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32:
		info.AdminType = FieldTypeNumber
		info.DBType = "INTEGER"
		
	case reflect.Uint64:
		info.AdminType = FieldTypeNumber
		info.DBType = "BIGINT"
		
	case reflect.Float32:
		info.AdminType = FieldTypeNumber
		info.DBType = "REAL"
		
	case reflect.Float64:
		info.AdminType = FieldTypeNumber
		info.DBType = "DOUBLE PRECISION"
		
	case reflect.Bool:
		info.AdminType = FieldTypeBoolean
		info.DBType = "BOOLEAN"
		
	case reflect.Slice, reflect.Array:
		info.AdminType = FieldTypeJSON
		info.DBType = "JSONB"
		
	case reflect.Map:
		info.AdminType = FieldTypeJSON
		info.DBType = "JSONB"
		
	case reflect.Struct:
		info.AdminType = FieldTypeJSON
		info.DBType = "JSONB"
	}
	
	return info
}

// formatLabel converts a field name to a human-readable label
func (s *SchemaIntrospector) formatLabel(name string) string {
	// Convert snake_case or camelCase to Title Case
	parts := strings.FieldsFunc(name, func(r rune) bool {
		return r == '_' || r == '-'
	})

	var result strings.Builder
	for i, part := range parts {
		if i > 0 {
			result.WriteString(" ")
		}
		if len(part) > 0 {
			result.WriteString(strings.ToUpper(part[:1]))
			if len(part) > 1 {
				result.WriteString(strings.ToLower(part[1:]))
			}
		}
	}

	// Handle camelCase
	if result.Len() == 0 {
		// Try to split camelCase
		var words []string
		var current strings.Builder
		for _, r := range name {
			if r >= 'A' && r <= 'Z' {
				if current.Len() > 0 {
					words = append(words, current.String())
					current.Reset()
				}
				current.WriteRune(r)
			} else {
				current.WriteRune(r)
			}
		}
		if current.Len() > 0 {
			words = append(words, current.String())
		}
		if len(words) > 0 {
			return strings.Join(words, " ")
		}
		return name
	}

	return result.String()
}

// inferTableName infers a table name from a model name
func (s *SchemaIntrospector) inferTableName(modelName string) string {
	// Convert PascalCase to snake_case and pluralize
	var result strings.Builder
	for i, r := range modelName {
		if i > 0 && r >= 'A' && r <= 'Z' {
			result.WriteString("_")
		}
		result.WriteRune(r)
	}
	
	tableName := strings.ToLower(result.String())
	
	// Simple pluralization (add 's' or 'es')
	if strings.HasSuffix(tableName, "y") {
		tableName = strings.TrimSuffix(tableName, "y") + "ies"
	} else if strings.HasSuffix(tableName, "s") || strings.HasSuffix(tableName, "x") || 
		strings.HasSuffix(tableName, "z") || strings.HasSuffix(tableName, "ch") || 
		strings.HasSuffix(tableName, "sh") {
		tableName += "es"
	} else {
		tableName += "s"
	}
	
	return tableName
}

// extractRelationshipMeta extracts metadata from an Ent edge
func (s *SchemaIntrospector) extractRelationshipMeta(edge ent.Edge) *RelationshipMeta {
	edgeValue := reflect.ValueOf(edge)
	if !edgeValue.IsValid() {
		return nil
	}

	meta := &RelationshipMeta{}

	// Extract edge name
	nameMethod := edgeValue.MethodByName("Name")
	if nameMethod.IsValid() && nameMethod.Type().NumIn() == 0 {
		if results := nameMethod.Call(nil); len(results) > 0 {
			if name, ok := results[0].Interface().(string); ok {
				meta.Name = name
			}
		}
	}

	if meta.Name == "" {
		// Try to get name from field
		nameField := edgeValue.FieldByName("Name")
		if nameField.IsValid() && nameField.Kind() == reflect.String {
			meta.Name = nameField.String()
		}
	}

	if meta.Name == "" {
		return nil // Can't extract edge without name
	}

	// Extract target type
	tagMethod := edgeValue.MethodByName("Tag")
	if tagMethod.IsValid() && tagMethod.Type().NumIn() == 0 {
		if results := tagMethod.Call(nil); len(results) > 0 {
			if tag, ok := results[0].Interface().(string); ok {
				// Tag typically contains the target type information
				// Parse it to get the model name
				meta.TargetModel = s.parseTargetModelFromTag(tag)
			}
		}
	}

	// Try to get descriptor
	descriptorMethod := edgeValue.MethodByName("Descriptor")
	if descriptorMethod.IsValid() {
		results := descriptorMethod.Call(nil)
		if len(results) > 0 {
			descriptor := results[0]
			if descriptor.IsValid() {
				s.extractRelationshipFromDescriptor(descriptor, meta)
			}
		}
	}

	// Infer relationship type from edge properties
	if meta.Type == "" {
		meta.Type = s.inferRelationshipType(edgeValue, meta)
	}

	return meta
}

// extractRelationshipFromDescriptor extracts relationship info from edge descriptor
func (s *SchemaIntrospector) extractRelationshipFromDescriptor(descriptor reflect.Value, meta *RelationshipMeta) {
	if !descriptor.IsValid() {
		return
	}

	// Try to get Type field
	typeField := descriptor.FieldByName("Type")
	if typeField.IsValid() {
		// Ent edge types: O2O, O2M, M2O, M2M
		typeStr := typeField.String()
		switch {
		case strings.Contains(typeStr, "O2O"):
			meta.Type = RelationshipTypeOneToOne
		case strings.Contains(typeStr, "O2M"):
			meta.Type = RelationshipTypeOneToMany
		case strings.Contains(typeStr, "M2O"):
			meta.Type = RelationshipTypeManyToOne
		case strings.Contains(typeStr, "M2M"):
			meta.Type = RelationshipTypeManyToMany
		}
	}

	// Try to get Required field
	requiredField := descriptor.FieldByName("Required")
	if requiredField.IsValid() && requiredField.Kind() == reflect.Bool {
		meta.Required = requiredField.Bool()
	}

	// Try to get Unique field
	uniqueField := descriptor.FieldByName("Unique")
	if uniqueField.IsValid() && uniqueField.Kind() == reflect.Bool {
		meta.Unique = uniqueField.Bool()
	}

	// Try to get Inverse field
	inverseField := descriptor.FieldByName("Inverse")
	if inverseField.IsValid() && inverseField.Kind() == reflect.Bool {
		meta.Inverse = inverseField.Bool()
	}
}

// inferRelationshipType infers relationship type from edge properties
func (s *SchemaIntrospector) inferRelationshipType(edgeValue reflect.Value, meta *RelationshipMeta) RelationshipType {
	// Check if it's unique (OneToOne)
	if meta.Unique {
		return RelationshipTypeOneToOne
	}

	// Check for ManyToMany (typically has a through table or junction table)
	// This is a simplified check - in production, you'd check for M2M edge type
	tagMethod := edgeValue.MethodByName("Tag")
	if tagMethod.IsValid() {
		results := tagMethod.Call(nil)
		if len(results) > 0 {
			if tag, ok := results[0].Interface().(string); ok {
				if strings.Contains(tag, "M2M") || strings.Contains(tag, "many_to_many") {
					return RelationshipTypeManyToMany
				}
			}
		}
	}

	// Default: assume ManyToOne (foreign key relationship)
	// In Ent, edges are typically defined from the "many" side
	return RelationshipTypeManyToOne
}

// parseTargetModelFromTag parses the target model name from an edge tag
func (s *SchemaIntrospector) parseTargetModelFromTag(tag string) string {
	// Ent tags typically contain type information
	// Example: "ent.Edge" or "User" or "ent.Schema"
	// This is a simplified parser - in production, you'd use Ent's actual tag format
	
	// Try to extract model name from tag
	// Remove common prefixes
	tag = strings.TrimPrefix(tag, "ent.")
	tag = strings.TrimPrefix(tag, "schema.")
	
	// If it's a simple name, return it
	if tag != "" && !strings.Contains(tag, ".") {
		return tag
	}
	
	return ""
}
