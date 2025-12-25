package admin

import (
	"context"
	"fmt"
	"reflect"
	"strings"
)

// RelationshipHelper provides utilities for working with Ent relationships
type RelationshipHelper struct {
	client interface{} // *ent.Client
}

// NewRelationshipHelper creates a new relationship helper
func NewRelationshipHelper(client interface{}) *RelationshipHelper {
	return &RelationshipHelper{
		client: client,
	}
}

// PrefetchRelationships prefetches related objects for a list of records
func (h *RelationshipHelper) PrefetchRelationships(ctx context.Context, modelName string, records []interface{}, relationships []string) error {
	if len(relationships) == 0 {
		return nil
	}

	// Get model client
	clientHelper := NewEntClientHelper(h.client)
	modelClient, err := clientHelper.GetModelClient(modelName)
	if err != nil {
		return err
	}

	// For each relationship, prefetch the related objects
	for _, relName := range relationships {
		if err := h.prefetchRelationship(ctx, modelClient, records, relName); err != nil {
			// Log error but continue with other relationships
			continue
		}
	}

	return nil
}

// prefetchRelationship prefetches a single relationship
func (h *RelationshipHelper) prefetchRelationship(ctx context.Context, modelClient reflect.Value, records []interface{}, relName string) error {
	// Ent generates methods like QueryX() for relationships
	// This is a placeholder - full implementation would:
	// 1. Get the edge metadata
	// 2. Use Ent's WithX() methods to eager load
	// 3. Populate the related objects in the records
	
	// For now, this is a placeholder
	return nil
}

// GetRelatedObjects gets related objects for a record
func (h *RelationshipHelper) GetRelatedObjects(ctx context.Context, modelName string, recordID interface{}, relationshipName string) ([]interface{}, error) {
	clientHelper := NewEntClientHelper(h.client)
	
	// Get the record first
	record, err := clientHelper.Get(ctx, modelName, recordID)
	if err != nil {
		return nil, err
	}

	// Get model client
	modelClient, err := clientHelper.GetModelClient(modelName)
	if err != nil {
		return nil, err
	}

	// Use reflection to call QueryX() method where X is the relationship name
	queryMethodName := "Query" + strings.ToUpper(relationshipName[:1]) + relationshipName[1:]
	queryMethod := modelClient.MethodByName(queryMethodName)
	
	if !queryMethod.IsValid() {
		return nil, fmt.Errorf("relationship %s not found on model %s", relationshipName, modelName)
	}

	// Call QueryX() to get the query builder
	queryBuilder := queryMethod.Call([]reflect.Value{reflect.ValueOf(record)})[0]

	// Call All() to get all related objects
	allMethod := queryBuilder.MethodByName("All")
	if !allMethod.IsValid() {
		return nil, fmt.Errorf("All() method not found for relationship %s", relationshipName)
	}

	results := allMethod.Call([]reflect.Value{reflect.ValueOf(ctx)})
	if len(results) < 1 {
		return nil, fmt.Errorf("All() returned no value")
	}

	if len(results) > 1 && !results[1].IsNil() {
		errVal := results[1].Interface()
		if err, ok := errVal.(error); ok {
			return nil, err
		}
	}

	// Convert slice to []interface{}
	sliceVal := results[0]
	if sliceVal.Kind() != reflect.Slice {
		return nil, fmt.Errorf("All() did not return a slice")
	}

	items := make([]interface{}, sliceVal.Len())
	for i := 0; i < sliceVal.Len(); i++ {
		items[i] = sliceVal.Index(i).Interface()
	}

	return items, nil
}

// SetRelatedObject sets a related object for a record (for ManyToOne, OneToOne)
func (h *RelationshipHelper) SetRelatedObject(ctx context.Context, modelName string, recordID interface{}, relationshipName string, relatedID interface{}) error {
	clientHelper := NewEntClientHelper(h.client)
	
	// Get the record
	record, err := clientHelper.Get(ctx, modelName, recordID)
	if err != nil {
		return err
	}

	// Get model client
	modelClient, err := clientHelper.GetModelClient(modelName)
	if err != nil {
		return err
	}

	// Get update builder
	updateMethod := modelClient.MethodByName("Update")
	if !updateMethod.IsValid() {
		return fmt.Errorf("Update() method not found for model %s", modelName)
	}

	updateBuilder := updateMethod.Call(nil)[0]

	// Apply ID filter
	if err := clientHelper.applyIDFilter(updateBuilder, recordID); err != nil {
		return err
	}

	// Set the relationship using SetX() method
	setMethodName := "Set" + strings.ToUpper(relationshipName[:1]) + relationshipName[1:]
	setMethod := updateBuilder.MethodByName(setMethodName)
	
	if !setMethod.IsValid() {
		return fmt.Errorf("Set%s() method not found", relationshipName)
	}

	// Get the related object
	relatedModelName := h.getRelatedModelName(modelName, relationshipName)
	relatedRecord, err := clientHelper.Get(ctx, relatedModelName, relatedID)
	if err != nil {
		return err
	}

	// Call SetX() with the related record
	setMethod.Call([]reflect.Value{reflect.ValueOf(relatedRecord)})

	// Save
	saveMethod := updateBuilder.MethodByName("Save")
	if !saveMethod.IsValid() {
		return fmt.Errorf("Save() method not found")
	}

	results := saveMethod.Call([]reflect.Value{reflect.ValueOf(ctx)})
	if len(results) > 1 && !results[1].IsNil() {
		errVal := results[1].Interface()
		if err, ok := errVal.(error); ok {
			return err
		}
	}

	return nil
}

// getRelatedModelName gets the related model name from relationship metadata
func (h *RelationshipHelper) getRelatedModelName(modelName string, relationshipName string) string {
	// This would need to look up the relationship metadata
	// For now, return a placeholder
	// In production, you'd query the registry for the model's relationships
	return ""
}

