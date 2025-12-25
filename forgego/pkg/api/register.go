package api

import (
	"reflect"

	"github.com/gofiber/fiber/v2"
	"github.com/forgego/forge/pkg/models"
)

func RegisterModel[T any](
	app *fiber.App,
	path string,
	db *models.DB,
	modelDef *models.ModelDefinition[T],
) {
	if path == "" {
		path = "/api/v1"
	}
	if path[len(path)-1] == '/' {
		path = path[:len(path)-1]
	}
	
	manager := models.NewModelManager[T](modelDef, db)
	serializer := NewModelSerializerWithDefinition[T](modelDef)
	
	// Mark primary key and timestamp fields as read-only
	fields := modelDef.GetFields()
	for _, field := range fields {
		fieldName := field.GetName()
		
		// Check if field is primary key by checking the underlying FieldDescriptor
		if isPrimaryKeyField(field) {
			serializer = serializer.ReadOnly(fieldName)
		}
		
		// Mark timestamp fields as read-only
		if fieldName == "CreatedAt" || fieldName == "UpdatedAt" {
			serializer = serializer.ReadOnly(fieldName)
		}
	}
	
	viewSet := NewBaseViewSet[T](manager, serializer)
	
	RegisterViewSet[T](app, path, viewSet)
}

// isPrimaryKeyField checks if a FieldDefinition represents a primary key field.
// It uses reflection to access the underlying FieldDescriptor's options.
func isPrimaryKeyField(field models.FieldDefinition) bool {
	// Use reflection to check if field is a fieldDefinitionAdapter
	fieldValue := reflect.ValueOf(field)
	if fieldValue.Kind() == reflect.Ptr {
		fieldValue = fieldValue.Elem()
	}
	
	// Check if it has a descriptor field (fieldDefinitionAdapter)
	descriptorField := fieldValue.FieldByName("descriptor")
	if descriptorField.IsValid() {
		// Get the descriptor's options
		optionsMethod := descriptorField.MethodByName("GetOptions")
		if optionsMethod.IsValid() {
			results := optionsMethod.Call(nil)
			if len(results) > 0 {
				options := results[0].Interface()
				optionsValue := reflect.ValueOf(options)
				primaryKeyField := optionsValue.FieldByName("PrimaryKey")
				if primaryKeyField.IsValid() && primaryKeyField.Kind() == reflect.Bool {
					return primaryKeyField.Bool()
				}
			}
		}
	}
	
	// Fallback: check by field name convention
	fieldName := field.GetName()
	return fieldName == "ID" || fieldName == "Id"
}

