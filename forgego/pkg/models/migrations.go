package models

import (
	"context"
	"fmt"
	"reflect"
)

// CreateTables creates tables from registered schemas using BUN.
func CreateTables(ctx context.Context, db *DB, schemas ...interface{}) error {
	for _, schema := range schemas {
		// Get the model definition for this schema
		schemaType := reflect.TypeOf(schema)
		if schemaType.Kind() == reflect.Ptr {
			schemaType = schemaType.Elem()
		}

		schemaMutex.RLock()
		modelDef, ok := schemaRegistry[schemaType]
		schemaMutex.RUnlock()

		if !ok {
			return fmt.Errorf("schema %s not registered, call RegisterSchema first", schemaType.Name())
		}

		// Create table using BUN
		tableName := modelDef.tableName
		if tableName == "" {
			tableName = toSnakeCase(schemaType.Name()) + "s"
		}

		// Create a new instance of the schema type for BUN
		schemaPtr := reflect.New(schemaType).Interface()
		
		// Use BUN's CreateTable
		_, err := db.NewCreateTable().
			Model(schemaPtr).
			IfNotExists().
			Exec(ctx)
		if err != nil {
			return fmt.Errorf("failed to create table for %s: %w", schemaType.Name(), err)
		}
	}

	return nil
}


