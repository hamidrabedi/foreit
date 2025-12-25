package models

import (
	"fmt"
	"reflect"
	"sync"
)

var (
	schemaRegistry = make(map[reflect.Type]*ModelDefinition[interface{}])
	schemaMutex    sync.RWMutex
)

// RegisterSchema registers a model schema and builds ModelDefinition.
func RegisterSchema[T any](schema T) error {
	return registerSchemaReflection(schema)
}

// RegisterSchemas registers multiple schemas.
func RegisterSchemas(schemas ...interface{}) error {
	for _, schema := range schemas {
		err := registerSchemaReflection(schema)
		if err != nil {
			return fmt.Errorf("failed to register schema %T: %w", schema, err)
		}
	}
	return nil
}

// registerSchemaReflection uses pure reflection to register a schema.
func registerSchemaReflection(schema interface{}) error {
	schemaValue := reflect.ValueOf(schema)
	schemaType := schemaValue.Type()
	if schemaType.Kind() == reflect.Ptr {
		schemaType = schemaType.Elem()
		schemaValue = schemaValue.Elem()
	}

	// Check if it implements SchemaInterface
	schemaInterfaceType := reflect.TypeOf((*SchemaInterface)(nil)).Elem()
	if !schemaType.Implements(schemaInterfaceType) {
		return fmt.Errorf("type %T does not implement SchemaInterface", schema)
	}

	// Call Fields() method
	fieldsMethod := schemaValue.MethodByName("Fields")
	if !fieldsMethod.IsValid() {
		if schemaValue.CanAddr() {
			fieldsMethod = schemaValue.Addr().MethodByName("Fields")
		}
	}

	if !fieldsMethod.IsValid() {
		return fmt.Errorf("type %T does not have Fields() method", schema)
	}

	fieldsResult := fieldsMethod.Call(nil)
	if len(fieldsResult) == 0 {
		return fmt.Errorf("Fields() method returned no results")
	}

	fieldDescriptors := fieldsResult[0].Interface().([]FieldDescriptor)

	// Call Relations() method
	relationsMethod := schemaValue.MethodByName("Relations")
	if !relationsMethod.IsValid() {
		if schemaValue.CanAddr() {
			relationsMethod = schemaValue.Addr().MethodByName("Relations")
		}
	}

	var relationDescriptors []RelationDescriptor
	if relationsMethod.IsValid() {
		relationsResult := relationsMethod.Call(nil)
		if len(relationsResult) > 0 {
			if descs, ok := relationsResult[0].Interface().([]RelationDescriptor); ok {
				relationDescriptors = descs
			}
		}
	}

	// Build Field definitions from descriptors
	fields := make([]FieldDefinition, 0, len(fieldDescriptors))
	for _, fd := range fieldDescriptors {
		fieldDef := fieldDescriptorToFieldDefinition(fd)
		if fieldDef != nil {
			fields = append(fields, fieldDef)
		}
	}

	// Build Relation definitions from descriptors
	relationships := make([]RelationDefinition, 0, len(relationDescriptors))
	for _, rd := range relationDescriptors {
		relDef := relationDescriptorToRelationDefinition(rd)
		if relDef.Name != "" {
			relationships = append(relationships, relDef)
		}
	}

	// Get table name if AdvancedSchemaInterface is implemented
	var tableName string
	advancedInterfaceType := reflect.TypeOf((*AdvancedSchemaInterface)(nil)).Elem()
	if schemaType.Implements(advancedInterfaceType) {
		tableNameMethod := schemaValue.MethodByName("TableName")
		if tableNameMethod.IsValid() {
			tableNameResult := tableNameMethod.Call(nil)
			if len(tableNameResult) > 0 {
				tableName = tableNameResult[0].String()
			}
		}
	}

	// Create ModelDefinition (type-erased)
	modelDef := &ModelDefinition[interface{}]{
		tableName:     tableName,
		fields:        fields,
		relationships: relationships,
		meta: ModelDefinitionMeta{
			TableName: tableName,
		},
	}

	schemaMutex.Lock()
	schemaRegistry[schemaType] = modelDef
	schemaMutex.Unlock()

	return nil
}

// GetModel returns the ModelDefinition for a type.
func GetModel[T any]() (*ModelDefinition[T], error) {
	var zero T
	t := reflect.TypeOf(zero)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	schemaMutex.RLock()
	defer schemaMutex.RUnlock()

	modelDef, ok := schemaRegistry[t]
	if !ok {
		return nil, fmt.Errorf("model %s not registered", t.Name())
	}

	// Convert from type-erased to typed ModelDefinition
	// This is a simplified version - in production, you'd need proper type conversion
	typedModelDef := &ModelDefinition[T]{
		tableName:     modelDef.tableName,
		fields:        modelDef.fields,
		relationships: modelDef.relationships,
		meta:          modelDef.meta,
	}

	return typedModelDef, nil
}

// fieldDescriptorToFieldDefinition converts a FieldDescriptor to FieldDefinition.
func fieldDescriptorToFieldDefinition(fd FieldDescriptor) FieldDefinition {
	return &fieldDefinitionAdapter{
		descriptor: fd,
	}
}

// fieldDefinitionAdapter adapts FieldDescriptor to FieldDefinition interface.
type fieldDefinitionAdapter struct {
	descriptor FieldDescriptor
}

func (f *fieldDefinitionAdapter) GetName() string {
	return f.descriptor.GetName()
}

func (f *fieldDefinitionAdapter) GetColumn() string {
	return f.descriptor.GetColumn()
}

// GetValidators implements FieldDefinitionWithValidators interface.
func (f *fieldDefinitionAdapter) GetValidators() []interface{} {
	options := f.descriptor.GetOptions()
	validators := []interface{}{}
	
	if options.Required {
		validators = append(validators, "required")
	}
	if options.MaxLength != nil {
		validators = append(validators, map[string]interface{}{
			"type": "max_length",
			"value": *options.MaxLength,
		})
	}
	if options.MinLength != nil {
		validators = append(validators, map[string]interface{}{
			"type": "min_length",
			"value": *options.MinLength,
		})
	}
	if options.Max != nil {
		validators = append(validators, map[string]interface{}{
			"type": "max",
			"value": *options.Max,
		})
	}
	if options.Min != nil {
		validators = append(validators, map[string]interface{}{
			"type": "min",
			"value": *options.Min,
		})
	}
	
	return validators
}

// relationDescriptorToRelationDefinition converts a RelationDescriptor to RelationDefinition.
func relationDescriptorToRelationDefinition(rd RelationDescriptor) RelationDefinition {
	relDef := RelationDefinition{
		Name:       rd.GetName(),
		Type:       rd.GetType(),
		RelatedModel: "",
		ForeignKey:   "",
		BackRef:      "",
		OnDelete:     "",
		OnUpdate:     "",
	}

	// Use reflection to extract field information from the descriptor
	rdValue := reflect.ValueOf(rd)
	if rdValue.Kind() == reflect.Ptr {
		rdValue = rdValue.Elem()
	}

	// Try to extract fromField from the descriptor
	// This is type-erased, so we need to use reflection
	fromField := rdValue.FieldByName("fromField")

	// Extract field column names if they implement FieldDefinition interface
	if fromField.IsValid() && !fromField.IsNil() {
		if fieldDef, ok := fromField.Interface().(FieldDefinition); ok {
			relDef.ForeignKey = fieldDef.GetColumn()
		} else {
			// Try to get column via reflection if it's a field.Field
			fieldValue := reflect.ValueOf(fromField.Interface())
			if fieldValue.Kind() == reflect.Struct {
				columnField := fieldValue.FieldByName("Column")
				if columnField.IsValid() && columnField.Kind() == reflect.String {
					relDef.ForeignKey = columnField.String()
				}
			}
		}
	}

	// Extract related model name from descriptor if stored
	relatedModelField := rdValue.FieldByName("relatedModel")
	if relatedModelField.IsValid() && relatedModelField.Kind() == reflect.String {
		relDef.RelatedModel = relatedModelField.String()
	}

	// Extract OnDelete and OnUpdate if present
	onDeleteField := rdValue.FieldByName("onDelete")
	if onDeleteField.IsValid() && onDeleteField.Kind() == reflect.String {
		relDef.OnDelete = onDeleteField.String()
	}

	onUpdateField := rdValue.FieldByName("onUpdate")
	if onUpdateField.IsValid() && onUpdateField.Kind() == reflect.String {
		relDef.OnUpdate = onUpdateField.String()
	}

	// Extract BackRef if present
	backRefField := rdValue.FieldByName("backRef")
	if backRefField.IsValid() && backRefField.Kind() == reflect.String {
		relDef.BackRef = backRefField.String()
	}

	return relDef
}

