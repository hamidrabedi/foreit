package admin

import (
	"context"
	"fmt"
	"reflect"
	"strings"
)

// EntClientHelper provides helper methods to work with Ent client using reflection
// This allows us to work with any Ent model type dynamically
type EntClientHelper struct {
	client interface{} // *ent.Client
}

// NewEntClientHelper creates a new Ent client helper
func NewEntClientHelper(client interface{}) *EntClientHelper {
	return &EntClientHelper{
		client: client,
	}
}

// GetModelClient gets the model-specific client (e.g., client.User, client.Article)
// This uses reflection to call methods like client.User() on the Ent client
func (h *EntClientHelper) GetModelClient(modelName string) (reflect.Value, error) {
	if h.client == nil {
		return reflect.Value{}, fmt.Errorf("ent client is nil")
	}

	clientVal := reflect.ValueOf(h.client)
	if clientVal.Kind() == reflect.Ptr {
		clientVal = clientVal.Elem()
	}

	// Ent generates methods like User(), Article(), etc. on the client
	// The method name matches the model name (capitalized)
	methodName := modelName
	method := clientVal.MethodByName(methodName)
	if !method.IsValid() {
		// Try with lowercase first letter (e.g., user instead of User)
		methodName = strings.ToLower(modelName[:1]) + modelName[1:]
		method = clientVal.MethodByName(methodName)
		if !method.IsValid() {
			return reflect.Value{}, fmt.Errorf("model client method %s() not found on ent client", modelName)
		}
	}

	// Call the method to get the model client (e.g., client.User())
	result := method.Call(nil)
	if len(result) == 0 {
		return reflect.Value{}, fmt.Errorf("model client method %s() returned no value", modelName)
	}

	return result[0], nil
}

// Create creates a new record using Ent's Create builder
func (h *EntClientHelper) Create(ctx context.Context, modelName string, data map[string]interface{}) (interface{}, error) {
	modelClient, err := h.GetModelClient(modelName)
	if err != nil {
		return nil, err
	}

	// Get the Create method (e.g., client.User().Create())
	createMethod := modelClient.MethodByName("Create")
	if !createMethod.IsValid() {
		return nil, fmt.Errorf("Create() method not found for model %s", modelName)
	}

	// Call Create() to get the Create builder
	createBuilder := createMethod.Call(nil)[0]

	// Apply field setters to the builder
	// Ent generates methods like SetName(), SetEmail(), etc.
	if err := h.applyFieldSetters(createBuilder, data); err != nil {
		return nil, fmt.Errorf("failed to set fields: %w", err)
	}

	// Call Save() on the builder to execute the create
	saveMethod := createBuilder.MethodByName("Save")
	if !saveMethod.IsValid() {
		return nil, fmt.Errorf("Save() method not found for model %s", modelName)
	}

	// Save() takes a context
	results := saveMethod.Call([]reflect.Value{reflect.ValueOf(ctx)})
	if len(results) < 1 {
		return nil, fmt.Errorf("Save() returned no value")
	}

	// Check for error (Save returns (Model, error))
	if len(results) > 1 && !results[1].IsNil() {
		errVal := results[1].Interface()
		if err, ok := errVal.(error); ok {
			return nil, err
		}
	}

	return results[0].Interface(), nil
}

// Update updates an existing record using Ent's Update builder
func (h *EntClientHelper) Update(ctx context.Context, modelName string, id interface{}, data map[string]interface{}) (interface{}, error) {
	modelClient, err := h.GetModelClient(modelName)
	if err != nil {
		return nil, err
	}

	// Get the Update method (e.g., client.User().Update())
	updateMethod := modelClient.MethodByName("Update")
	if !updateMethod.IsValid() {
		return nil, fmt.Errorf("Update() method not found for model %s", modelName)
	}

	// Call Update() to get the Update builder
	updateBuilder := updateMethod.Call(nil)[0]

	// Apply ID filter (e.g., WhereID(id))
	if err := h.applyIDFilter(updateBuilder, id); err != nil {
		return nil, fmt.Errorf("failed to set ID filter: %w", err)
	}

	// Apply field setters
	if err := h.applyFieldSetters(updateBuilder, data); err != nil {
		return nil, fmt.Errorf("failed to set fields: %w", err)
	}

	// Call Save() to execute the update
	saveMethod := updateBuilder.MethodByName("Save")
	if !saveMethod.IsValid() {
		return nil, fmt.Errorf("Save() method not found for model %s", modelName)
	}

	results := saveMethod.Call([]reflect.Value{reflect.ValueOf(ctx)})
	if len(results) < 1 {
		return nil, fmt.Errorf("Save() returned no value")
	}

	if len(results) > 1 && !results[1].IsNil() {
		errVal := results[1].Interface()
		if err, ok := errVal.(error); ok {
			return nil, err
		}
	}

	return results[0].Interface(), nil
}

// Delete deletes a record using Ent's Delete builder
func (h *EntClientHelper) Delete(ctx context.Context, modelName string, id interface{}) error {
	modelClient, err := h.GetModelClient(modelName)
	if err != nil {
		return err
	}

	// Get the Delete method (e.g., client.User().Delete())
	deleteMethod := modelClient.MethodByName("Delete")
	if !deleteMethod.IsValid() {
		return fmt.Errorf("Delete() method not found for model %s", modelName)
	}

	// Call Delete() to get the Delete builder
	deleteBuilder := deleteMethod.Call(nil)[0]

	// Apply ID filter
	if err := h.applyIDFilter(deleteBuilder, id); err != nil {
		return fmt.Errorf("failed to set ID filter: %w", err)
	}

	// Call Exec() to execute the delete
	execMethod := deleteBuilder.MethodByName("Exec")
	if !execMethod.IsValid() {
		return fmt.Errorf("Exec() method not found for model %s", modelName)
	}

	results := execMethod.Call([]reflect.Value{reflect.ValueOf(ctx)})
	if len(results) > 0 && !results[0].IsNil() {
		errVal := results[0].Interface()
		if err, ok := errVal.(error); ok {
			return err
		}
	}

	return nil
}

// Get retrieves a single record by ID
func (h *EntClientHelper) Get(ctx context.Context, modelName string, id interface{}) (interface{}, error) {
	modelClient, err := h.GetModelClient(modelName)
	if err != nil {
		return nil, err
	}

	// Get the Query method (e.g., client.User().Query())
	queryMethod := modelClient.MethodByName("Query")
	if !queryMethod.IsValid() {
		return nil, fmt.Errorf("Query() method not found for model %s", modelName)
	}

	// Call Query() to get the query builder
	queryBuilder := queryMethod.Call(nil)[0]

	// Apply ID filter
	if err := h.applyIDFilter(queryBuilder, id); err != nil {
		return nil, fmt.Errorf("failed to set ID filter: %w", err)
	}

	// Call Only() to get only the ID field, then First() to get the record
	// Actually, we want all fields, so we'll call First() directly
	firstMethod := queryBuilder.MethodByName("First")
	if !firstMethod.IsValid() {
		return nil, fmt.Errorf("First() method not found for model %s", modelName)
	}

	results := firstMethod.Call([]reflect.Value{reflect.ValueOf(ctx)})
	if len(results) < 1 {
		return nil, fmt.Errorf("First() returned no value")
	}

	if len(results) > 1 && !results[1].IsNil() {
		errVal := results[1].Interface()
		if err, ok := errVal.(error); ok {
			return nil, err
		}
	}

	return results[0].Interface(), nil
}

// List retrieves multiple records with pagination and filters
func (h *EntClientHelper) List(ctx context.Context, modelName string, params *QueryParams) ([]interface{}, int64, error) {
	modelClient, err := h.GetModelClient(modelName)
	if err != nil {
		return nil, 0, err
	}

	// Get the Query method
	queryMethod := modelClient.MethodByName("Query")
	if !queryMethod.IsValid() {
		return nil, 0, fmt.Errorf("Query() method not found for model %s", modelName)
	}

	// Call Query() to get the query builder
	queryBuilder := queryMethod.Call(nil)[0]

	// Apply filters (including relationship filters)
	if err := h.applyFilters(queryBuilder, params.Filters, modelName); err != nil {
		// Some filters might fail (e.g., relationship filters not yet fully supported)
		// Log but continue for now
		// return nil, 0, fmt.Errorf("failed to apply filters: %w", err)
	}

	// Apply search (if any)
	if params.Search != "" {
		searchFields := params.SearchFields
		if len(searchFields) == 0 {
			// Default search fields - could be configured per model
			searchFields = []string{"name", "title", "email"}
		}
		if err := h.applySearch(queryBuilder, params.Search, searchFields); err != nil {
			// Search might not be supported for all models, so we log but don't fail
			// return nil, 0, fmt.Errorf("failed to apply search: %w", err)
		}
	}

	// Count total records (before pagination)
	countMethod := queryBuilder.MethodByName("Count")
	if !countMethod.IsValid() {
		return nil, 0, fmt.Errorf("Count() method not found for model %s", modelName)
	}
	countResults := countMethod.Call([]reflect.Value{reflect.ValueOf(ctx)})
	var total int64
	if len(countResults) > 0 {
		if countVal, ok := countResults[0].Interface().(int); ok {
			total = int64(countVal)
		} else if countVal, ok := countResults[0].Interface().(int64); ok {
			total = countVal
		}
	}

	// Apply sorting
	if err := h.applySorting(queryBuilder, params.SortBy, params.SortOrder); err != nil {
		return nil, 0, fmt.Errorf("failed to apply sorting: %w", err)
	}

	// Apply pagination
	if err := h.applyPagination(queryBuilder, params.Page, params.PageSize); err != nil {
		return nil, 0, fmt.Errorf("failed to apply pagination: %w", err)
	}

	// Call All() to get all records
	allMethod := queryBuilder.MethodByName("All")
	if !allMethod.IsValid() {
		return nil, 0, fmt.Errorf("All() method not found for model %s", modelName)
	}

	results := allMethod.Call([]reflect.Value{reflect.ValueOf(ctx)})
	if len(results) < 1 {
		return nil, 0, fmt.Errorf("All() returned no value")
	}

	if len(results) > 1 && !results[1].IsNil() {
		errVal := results[1].Interface()
		if err, ok := errVal.(error); ok {
			return nil, 0, err
		}
	}

	// Convert slice to []interface{}
	sliceVal := results[0]
	if sliceVal.Kind() != reflect.Slice {
		return nil, 0, fmt.Errorf("All() did not return a slice")
	}

	items := make([]interface{}, sliceVal.Len())
	for i := 0; i < sliceVal.Len(); i++ {
		items[i] = sliceVal.Index(i).Interface()
	}

	return items, total, nil
}

// applyFieldSetters applies field setters to a builder (Create or Update)
func (h *EntClientHelper) applyFieldSetters(builder reflect.Value, data map[string]interface{}) error {
	for fieldName, value := range data {
		// Ent generates setters like SetName(), SetEmail(), etc.
		setterName := "Set" + strings.ToUpper(fieldName[:1]) + fieldName[1:]
		setter := builder.MethodByName(setterName)
		if !setter.IsValid() {
			// Field might not exist or might be read-only, skip it
			continue
		}

		// Convert value to appropriate type
		valueVal := reflect.ValueOf(value)
		setterType := setter.Type()

		// Setter takes one parameter
		if setterType.NumIn() != 1 {
			continue
		}

		paramType := setterType.In(0)
		convertedVal := h.convertValue(valueVal, paramType)
		if !convertedVal.IsValid() {
			continue
		}

		// Call the setter
		setter.Call([]reflect.Value{convertedVal})
	}

	return nil
}

// applyIDFilter applies an ID filter to a query/update/delete builder
func (h *EntClientHelper) applyIDFilter(builder reflect.Value, id interface{}) error {
	// Ent generates methods like WhereID() or ID()
	methods := []string{"WhereID", "ID"}
	for _, methodName := range methods {
		method := builder.MethodByName(methodName)
		if method.IsValid() {
			idVal := reflect.ValueOf(id)
			methodType := method.Type()
			if methodType.NumIn() == 1 {
				paramType := methodType.In(0)
				convertedID := h.convertValue(idVal, paramType)
				if convertedID.IsValid() {
					method.Call([]reflect.Value{convertedID})
					return nil
				}
			}
		}
	}

	return fmt.Errorf("could not find ID filter method (WhereID or ID)")
}

// applyFilters applies filters to a query builder
func (h *EntClientHelper) applyFilters(queryBuilder reflect.Value, filters map[string]interface{}, modelName string) error {
	for key, filterValue := range filters {
		var fieldName, operator string
		var value interface{}
		
		// Check if it's a FilterSpec map
		if filterMap, ok := filterValue.(map[string]interface{}); ok {
			if field, ok := filterMap["field"].(string); ok {
				fieldName = field
			}
			if op, ok := filterMap["operator"].(string); ok {
				operator = op
			}
			value = filterMap["value"]
		} else {
			// Simple filter - check if key has operator
			parts := strings.Split(key, "__")
			fieldName = parts[0]
			if len(parts) > 1 {
				operator = parts[1]
			} else {
				operator = "eq"
			}
			value = filterValue
		}
		
		// Check if this is a relationship filter (e.g., "author__id", "author__email")
		if strings.Contains(fieldName, "__") {
			parts := strings.SplitN(fieldName, "__", 2)
			if len(parts) == 2 {
				relName := parts[0]
				relField := parts[1]
				// Try to apply relationship filter
				if err := h.applyRelationshipFilter(queryBuilder, relName, relField, operator, value); err == nil {
					continue // Successfully applied relationship filter
				}
				// If relationship filter fails, fall through to regular filter
			}
		}

		// Apply filter based on operator
		if err := h.applyFilterWithOperator(queryBuilder, fieldName, operator, value); err != nil {
			return fmt.Errorf("failed to apply filter %s %s %v: %w", fieldName, operator, value, err)
		}
	}
	return nil
}

// applyFilterWithOperator applies a filter with a specific operator
func (h *EntClientHelper) applyFilterWithOperator(queryBuilder reflect.Value, fieldName string, operator string, value interface{}) error {
	// Handle special operators
	switch operator {
	case "between":
		// Value should be a comma-separated string or array
		return h.applyBetweenFilter(queryBuilder, fieldName, value)
	case "isnull":
		// Check if field is null or not null
		return h.applyNullFilter(queryBuilder, fieldName, value)
	case "in":
		// Value should be a comma-separated string or array
		return h.applyInFilter(queryBuilder, fieldName, value)
	case "contains", "like":
		// String contains/like search
		return h.applyContainsFilter(queryBuilder, fieldName, value)
	}
	
	// For comparison operators, try to find Ent's predicate methods
	// Ent generates methods like WhereNameEQ, WhereNameNE, WhereNameGT, etc.
	predicateMethodName := h.buildPredicateMethodName(fieldName, operator)
	predicateMethod := queryBuilder.MethodByName(predicateMethodName)
	
	if predicateMethod.IsValid() {
		valueVal := reflect.ValueOf(value)
		predicateMethodType := predicateMethod.Type()
		if predicateMethodType.NumIn() == 1 {
			paramType := predicateMethodType.In(0)
			convertedVal := h.convertValue(valueVal, paramType)
			if convertedVal.IsValid() {
				predicateMethod.Call([]reflect.Value{convertedVal})
				return nil
			}
		}
	}
	
	// Fallback to simple Where method
	whereMethodName := "Where" + strings.ToUpper(fieldName[:1]) + fieldName[1:]
	whereMethod := queryBuilder.MethodByName(whereMethodName)
	if whereMethod.IsValid() {
		valueVal := reflect.ValueOf(value)
		whereMethodType := whereMethod.Type()
		if whereMethodType.NumIn() == 1 {
			paramType := whereMethodType.In(0)
			convertedVal := h.convertValue(valueVal, paramType)
			if convertedVal.IsValid() {
				whereMethod.Call([]reflect.Value{convertedVal})
				return nil
			}
		}
	}
	
	return fmt.Errorf("could not apply filter for field %s with operator %s", fieldName, operator)
}

// buildPredicateMethodName builds Ent predicate method name
func (h *EntClientHelper) buildPredicateMethodName(fieldName string, operator string) string {
	fieldCapitalized := strings.ToUpper(fieldName[:1]) + fieldName[1:]
	
	operatorMap := map[string]string{
		"eq":  "EQ",
		"ne":  "NEQ",
		"gt":  "GT",
		"gte": "GTE",
		"lt":  "LT",
		"lte": "LTE",
	}
	
	opSuffix, ok := operatorMap[operator]
	if !ok {
		opSuffix = "EQ"
	}
	
	return "Where" + fieldCapitalized + opSuffix
}

// applyBetweenFilter applies a between filter (for date ranges, number ranges)
func (h *EntClientHelper) applyBetweenFilter(queryBuilder reflect.Value, fieldName string, value interface{}) error {
	var min, max interface{}
	
	// Parse value - could be "min,max" string or array
	if str, ok := value.(string); ok {
		parts := strings.Split(str, ",")
		if len(parts) == 2 {
			min = strings.TrimSpace(parts[0])
			max = strings.TrimSpace(parts[1])
		} else {
			return fmt.Errorf("between filter requires two values separated by comma")
		}
	} else if arr, ok := value.([]interface{}); ok && len(arr) == 2 {
		min = arr[0]
		max = arr[1]
	} else {
		return fmt.Errorf("invalid between filter value format")
	}
	
	// Apply GTE for min value
	if err := h.applyFilterWithOperator(queryBuilder, fieldName, "gte", min); err != nil {
		return err
	}
	
	// Apply LTE for max value
	return h.applyFilterWithOperator(queryBuilder, fieldName, "lte", max)
}

// applyNullFilter applies a null/not null filter
func (h *EntClientHelper) applyNullFilter(queryBuilder reflect.Value, fieldName string, value interface{}) error {
	isNull := false
	if str, ok := value.(string); ok {
		isNull = strings.ToLower(str) == "true" || str == "1"
	} else if b, ok := value.(bool); ok {
		isNull = b
	}
	
	// Ent generates IsNull() and NotNull() methods
	methodName := "Where" + strings.ToUpper(fieldName[:1]) + fieldName[1:]
	if isNull {
		methodName += "IsNil"
	} else {
		methodName += "NotNil"
	}
	
	method := queryBuilder.MethodByName(methodName)
	if method.IsValid() {
		method.Call(nil)
		return nil
	}
	
	return fmt.Errorf("null filter method not found for field %s", fieldName)
}

// applyInFilter applies an "in" filter (value in array)
func (h *EntClientHelper) applyInFilter(queryBuilder reflect.Value, fieldName string, value interface{}) error {
	var values []interface{}
	
	// Parse value - could be comma-separated string or array
	if str, ok := value.(string); ok {
		parts := strings.Split(str, ",")
		values = make([]interface{}, len(parts))
		for i, part := range parts {
			values[i] = strings.TrimSpace(part)
		}
	} else if arr, ok := value.([]interface{}); ok {
		values = arr
	} else {
		return fmt.Errorf("invalid in filter value format")
	}
	
	// Ent generates WhereXIn() methods
	methodName := "Where" + strings.ToUpper(fieldName[:1]) + fieldName[1:] + "In"
	method := queryBuilder.MethodByName(methodName)
	if method.IsValid() {
		// Convert values to appropriate slice type
		valueSlice := reflect.MakeSlice(reflect.SliceOf(reflect.TypeOf(values[0])), len(values), len(values))
		for i, v := range values {
			valueSlice.Index(i).Set(reflect.ValueOf(v))
		}
		method.Call([]reflect.Value{valueSlice})
		return nil
	}
	
	return fmt.Errorf("in filter method not found for field %s", fieldName)
}

// applyRelationshipFilter applies a filter on a related model
func (h *EntClientHelper) applyRelationshipFilter(queryBuilder reflect.Value, relName string, relField string, operator string, value interface{}) error {
	// Ent generates methods like WhereHasAuthor() or WhereHasAuthorWith()
	// For filtering by related fields, we use WhereHasXWith() pattern
	
	hasMethodName := "WhereHas" + strings.ToUpper(relName[:1]) + relName[1:]
	hasMethod := queryBuilder.MethodByName(hasMethodName)
	
	if !hasMethod.IsValid() {
		// Try With suffix
		hasMethodName += "With"
		hasMethod = queryBuilder.MethodByName(hasMethodName)
	}
	
	if hasMethod.IsValid() {
		// WhereHasXWith() takes a predicate function
		// We need to create a predicate that filters by the related field
		// This is complex with reflection, so we'll use a simpler approach:
		// Try to find a method that takes the field name and value
		
		// For now, return an error indicating relationship filtering needs the predicate
		// In production, you'd use Ent's predicate builder
		return fmt.Errorf("relationship filtering requires predicate builder (not yet implemented)")
	}
	
	return fmt.Errorf("relationship %s not found", relName)
}

// applyContainsFilter applies a contains/like filter (for text search)
func (h *EntClientHelper) applyContainsFilter(queryBuilder reflect.Value, fieldName string, value interface{}) error {
	// Ent doesn't have built-in LIKE, but we can use Where with a predicate
	// For now, we'll use a simple equality check and note that full-text search
	// would require database-specific implementations
	
	// Try to find Contains method
	methodName := "Where" + strings.ToUpper(fieldName[:1]) + fieldName[1:] + "Contains"
	method := queryBuilder.MethodByName(methodName)
	if method.IsValid() {
		valueVal := reflect.ValueOf(value)
		method.Call([]reflect.Value{valueVal})
		return nil
	}
	
	// Fallback: use simple Where (exact match)
	// In production, this would use database-specific LIKE/ILIKE
	return h.applyFilterWithOperator(queryBuilder, fieldName, "eq", value)
}

// applySorting applies sorting to a query builder
func (h *EntClientHelper) applySorting(queryBuilder reflect.Value, sortBy string, sortOrder string) error {
	if sortBy == "" {
		return nil
	}

	// Ent generates methods like OrderByID(), OrderByName(), etc.
	orderMethodName := "Order"
	if sortOrder == "desc" {
		orderMethodName = "OrderDesc"
	}

	// Try to find a field-specific order method
	fieldOrderMethodName := orderMethodName + strings.ToUpper(sortBy[:1]) + sortBy[1:]
	fieldOrderMethod := queryBuilder.MethodByName(fieldOrderMethodName)
	if fieldOrderMethod.IsValid() {
		fieldOrderMethod.Call(nil)
		return nil
	}

	// Fallback to generic Order() method if available
	orderMethod := queryBuilder.MethodByName(orderMethodName)
	if orderMethod.IsValid() {
		// Order() might take a field parameter
		if orderMethod.Type().NumIn() > 0 {
			// Try to pass the field name or a field enum
			// This is model-specific, so we'll skip for now
		} else {
			orderMethod.Call(nil)
		}
	}

	return nil
}

// applyPagination applies pagination to a query builder
func (h *EntClientHelper) applyPagination(queryBuilder reflect.Value, page int, pageSize int) error {
	if pageSize <= 0 {
		return nil
	}

	// Ent uses Limit() and Offset() methods
	limitMethod := queryBuilder.MethodByName("Limit")
	if limitMethod.IsValid() {
		limitMethod.Call([]reflect.Value{reflect.ValueOf(pageSize)})
	}

	offset := (page - 1) * pageSize
	offsetMethod := queryBuilder.MethodByName("Offset")
	if offsetMethod.IsValid() {
		offsetMethod.Call([]reflect.Value{reflect.ValueOf(offset)})
	}

	return nil
}

// applySearch applies search to a query builder
// Searches across multiple text fields using OR conditions
func (h *EntClientHelper) applySearch(queryBuilder reflect.Value, search string, searchFields []string) error {
	if search == "" {
		return nil
	}
	
	// If no search fields specified, try to find text fields from model metadata
	// For now, we'll use a simple approach: search in common text fields
	if len(searchFields) == 0 {
		searchFields = []string{"name", "title", "email", "description", "content"}
	}
	
	// Ent doesn't have built-in OR conditions easily accessible via reflection
	// For a proper implementation, we would need to:
	// 1. Use Ent's predicate system with Or()
	// 2. Or use raw SQL with database-specific full-text search
	// 
	// For now, we'll apply search to the first available field
	// In production, this should use Ent's predicate builder or raw SQL
	
	for _, fieldName := range searchFields {
		// Try Contains method first
		methodName := "Where" + strings.ToUpper(fieldName[:1]) + fieldName[1:] + "Contains"
		method := queryBuilder.MethodByName(methodName)
		if method.IsValid() {
			method.Call([]reflect.Value{reflect.ValueOf(search)})
			return nil // Found and applied, return
		}
	}
	
	// If no Contains method found, try simple Where on first field
	if len(searchFields) > 0 {
		fieldName := searchFields[0]
		whereMethodName := "Where" + strings.ToUpper(fieldName[:1]) + fieldName[1:]
		whereMethod := queryBuilder.MethodByName(whereMethodName)
		if whereMethod.IsValid() {
			whereMethod.Call([]reflect.Value{reflect.ValueOf(search)})
			return nil
		}
	}
	
	// Search not supported for this model
	return fmt.Errorf("search not supported - no searchable fields found")
}

// convertValue converts a value to the target type
func (h *EntClientHelper) convertValue(value reflect.Value, targetType reflect.Type) reflect.Value {
	if !value.IsValid() {
		return reflect.Value{}
	}

	// If types match, return as-is
	if value.Type().AssignableTo(targetType) {
		return value
	}

	// Try to convert
	if value.Type().ConvertibleTo(targetType) {
		return value.Convert(targetType)
	}

	// Handle string to other types
	if value.Kind() == reflect.String {
		strVal := value.String()
		switch targetType.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			// Try to parse as int
			// This is simplified - in production, use strconv
			return reflect.Value{}
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			return reflect.Value{}
		case reflect.Float32, reflect.Float64:
			return reflect.Value{}
		case reflect.Bool:
			return reflect.Value{}
		}
	}

	return reflect.Value{}
}

