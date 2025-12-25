package admin

import (
	"context"
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
)

// BulkCreateRequest represents a bulk create request
type BulkCreateRequest struct {
	Items []map[string]interface{} `json:"items"`
}

// BulkUpdateRequest represents a bulk update request
type BulkUpdateRequest struct {
	Items []BulkUpdateItem `json:"items"`
}

// BulkUpdateItem represents a single item in bulk update
type BulkUpdateItem struct {
	ID      interface{}            `json:"id"`
	Updates map[string]interface{} `json:"updates"`
}

// BulkDeleteRequest represents a bulk delete request
type BulkDeleteRequest struct {
	IDs []interface{} `json:"ids"`
}

// BulkCreate handles POST /admin/api/{model}/bulk - Bulk create records
func (h *CRUDHandler) BulkCreate(c *fiber.Ctx) error {
	// Check permissions
	permCtx := &Context{
		User:    GetUserFromContext(c),
		Request: c,
		Model:   h.modelMeta,
		Action:  "create",
	}

	allowed, err := h.permissionChecker.CheckPermission(permCtx)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Permission check failed",
			"details": err.Error(),
		})
	}
	if !allowed {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "You don't have permission to create this resource",
		})
	}

	// Parse request body
	var req BulkCreateRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
			"details": err.Error(),
		})
	}

	if len(req.Items) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "No items provided",
		})
	}

	// Limit batch size
	maxBatchSize := 100
	if len(req.Items) > maxBatchSize {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": fmt.Sprintf("Batch size exceeds maximum of %d", maxBatchSize),
		})
	}

	ctx := context.Background()
	
	// Check if transaction is requested
	useTransaction := c.Query("transaction", "true") == "true"
	
	var created []interface{}
	var errors []map[string]interface{}
	
	if useTransaction {
		// Execute in transaction
		result, errs := h.bulkCreateInTransaction(ctx, req.Items, c)
		created = result
		errors = errs
	} else {
		// Execute without transaction (partial success allowed)
		created = make([]interface{}, 0, len(req.Items))
		errors = make([]map[string]interface{}, 0)
		
		// Create each item
		for i, item := range req.Items {
		// Execute before hook
		hookCtx := &HookContext{
			Model:   h.modelMeta,
			Action:  "create",
			User:    GetUserFromContext(c),
			Request: c,
			Data:    item,
		}
		if err := h.hookRegistry.Execute(h.modelMeta.Name, HookBeforeCreate, hookCtx); err != nil {
			errors = append(errors, map[string]interface{}{
				"index": i,
				"error": err.Error(),
			})
			continue
		}

		// Validate data
		if err := h.validateData(item); err != nil {
			errors = append(errors, map[string]interface{}{
				"index": i,
				"error": fmt.Sprintf("Validation failed: %s", err.Error()),
			})
			continue
		}

		// Create record
		result, err := h.entHelper.Create(ctx, h.modelMeta.Name, item)
		if err != nil {
			errors = append(errors, map[string]interface{}{
				"index": i,
				"error": err.Error(),
			})
			continue
		}

		// Execute after hook
		hookCtx.Result = result
		if err := h.hookRegistry.Execute(h.modelMeta.Name, HookAfterCreate, hookCtx); err != nil {
			// Log error but don't fail the creation
		}

			created = append(created, h.modelToMap(result))
		}
	}

	response := fiber.Map{
		"data":   created,
		"count":  len(created),
		"errors": errors,
	}

	if len(errors) > 0 && len(created) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	return c.Status(fiber.StatusCreated).JSON(response)
}

// bulkCreateInTransaction creates multiple records in a transaction
func (h *CRUDHandler) bulkCreateInTransaction(ctx context.Context, items []map[string]interface{}, c *fiber.Ctx) ([]interface{}, []map[string]interface{}) {
	// Get Ent client
	clientVal := reflect.ValueOf(h.client)
	if clientVal.Kind() == reflect.Ptr {
		clientVal = clientVal.Elem()
	}
	
	// Get Tx() method for transaction
	// Ent's Tx() signature: Tx(ctx context.Context) (*Tx, error)
	txMethod := clientVal.MethodByName("Tx")
	if !txMethod.IsValid() {
		// Fallback to non-transactional if Tx not available
		return h.bulkCreateWithoutTransaction(ctx, items, c)
	}
	
	// Start transaction
	txResults := txMethod.Call([]reflect.Value{reflect.ValueOf(ctx)})
	if len(txResults) < 2 {
		return nil, []map[string]interface{}{{"error": "failed to start transaction"}}
	}
	
	if !txResults[1].IsNil() {
		err := txResults[1].Interface().(error)
		return nil, []map[string]interface{}{{"error": fmt.Sprintf("failed to start transaction: %v", err)}}
	}
	
	txClient := txResults[0].Interface()
	txHelper := NewEntClientHelper(txClient)
	
	created := make([]interface{}, 0, len(items))
	
	// Create all items in transaction
	for i, item := range items {
		// Execute before hook
		hookCtx := &HookContext{
			Model:   h.modelMeta,
			Action:  "create",
			User:    GetUserFromContext(c),
			Request: c,
			Data:    item,
		}
		if err := h.hookRegistry.Execute(h.modelMeta.Name, HookBeforeCreate, hookCtx); err != nil {
			// Rollback transaction
			h.rollbackTransaction(txClient)
			return nil, []map[string]interface{}{{
				"index": i,
				"error": fmt.Sprintf("hook failed: %v", err),
			}}
		}
		
		// Validate
		if err := h.validateData(item); err != nil {
			// Rollback transaction
			h.rollbackTransaction(txClient)
			return nil, []map[string]interface{}{{
				"index": i,
				"error": fmt.Sprintf("validation failed: %v", err),
			}}
		}
		
		// Create using transaction client
		result, err := txHelper.Create(ctx, h.modelMeta.Name, item)
		if err != nil {
			// Rollback transaction
			h.rollbackTransaction(txClient)
			return nil, []map[string]interface{}{{
				"index": i,
				"error": fmt.Sprintf("create failed: %v", err),
			}}
		}
		
		created = append(created, result)
	}
	
	// Commit transaction
	if err := h.commitTransaction(txClient); err != nil {
		return nil, []map[string]interface{}{{"error": fmt.Sprintf("commit failed: %v", err)}}
	}
	
	// Execute after hooks (outside transaction)
	for i, item := range created {
		hookCtx := &HookContext{
			Model:   h.modelMeta,
			Action:  "create",
			User:    GetUserFromContext(c),
			Request: c,
			Result:  item,
		}
		// Execute but don't fail if hook errors
		h.hookRegistry.Execute(h.modelMeta.Name, HookAfterCreate, hookCtx)
	}
	
	// Convert results to maps
	resultMaps := make([]interface{}, len(created))
	for i, item := range created {
		resultMaps[i] = h.modelToMap(item)
	}
	
	return resultMaps, nil
}

// commitTransaction commits an Ent transaction
func (h *CRUDHandler) commitTransaction(txClient interface{}) error {
	txVal := reflect.ValueOf(txClient)
	if txVal.Kind() == reflect.Ptr {
		txVal = txVal.Elem()
	}
	
	commitMethod := txVal.MethodByName("Commit")
	if !commitMethod.IsValid() {
		return fmt.Errorf("Commit() method not found")
	}
	
	results := commitMethod.Call(nil)
	if len(results) > 0 && !results[0].IsNil() {
		err := results[0].Interface().(error)
		return err
	}
	
	return nil
}

// rollbackTransaction rolls back an Ent transaction
func (h *CRUDHandler) rollbackTransaction(txClient interface{}) error {
	txVal := reflect.ValueOf(txClient)
	if txVal.Kind() == reflect.Ptr {
		txVal = txVal.Elem()
	}
	
	rollbackMethod := txVal.MethodByName("Rollback")
	if !rollbackMethod.IsValid() {
		return fmt.Errorf("Rollback() method not found")
	}
	
	results := rollbackMethod.Call(nil)
	if len(results) > 0 && !results[0].IsNil() {
		err := results[0].Interface().(error)
		return err
	}
	
	return nil
}

// bulkCreateWithoutTransaction creates records without transaction
func (h *CRUDHandler) bulkCreateWithoutTransaction(ctx context.Context, items []map[string]interface{}, c *fiber.Ctx) ([]interface{}, []map[string]interface{}) {
	created := make([]interface{}, 0, len(items))
	errors := make([]map[string]interface{}, 0)
	
	for i, item := range items {
		// Execute before hook
		hookCtx := &HookContext{
			Model:   h.modelMeta,
			Action:  "create",
			User:    GetUserFromContext(c),
			Request: c,
			Data:    item,
		}
		if err := h.hookRegistry.Execute(h.modelMeta.Name, HookBeforeCreate, hookCtx); err != nil {
			errors = append(errors, map[string]interface{}{
				"index": i,
				"error": err.Error(),
			})
			continue
		}
		
		// Validate
		if err := h.validateData(item); err != nil {
			errors = append(errors, map[string]interface{}{
				"index": i,
				"error": fmt.Sprintf("Validation failed: %s", err.Error()),
			})
			continue
		}
		
		// Create
		result, err := h.entHelper.Create(ctx, h.modelMeta.Name, item)
		if err != nil {
			errors = append(errors, map[string]interface{}{
				"index": i,
				"error": err.Error(),
			})
			continue
		}
		
		// Execute after hook
		hookCtx.Result = result
		if err := h.hookRegistry.Execute(h.modelMeta.Name, HookAfterCreate, hookCtx); err != nil {
			// Log but don't fail
		}
		
		created = append(created, h.modelToMap(result))
	}
	
	return created, errors
}

// BulkUpdate handles PUT /admin/api/{model}/bulk - Bulk update records
func (h *CRUDHandler) BulkUpdate(c *fiber.Ctx) error {
	// Check permissions
	permCtx := &Context{
		User:    GetUserFromContext(c),
		Request: c,
		Model:   h.modelMeta,
		Action:  "update",
	}

	allowed, err := h.permissionChecker.CheckPermission(permCtx)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Permission check failed",
			"details": err.Error(),
		})
	}
	if !allowed {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "You don't have permission to update this resource",
		})
	}

	// Parse request body
	var req BulkUpdateRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
			"details": err.Error(),
		})
	}

	if len(req.Items) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "No items provided",
		})
	}

	// Limit batch size
	maxBatchSize := 100
	if len(req.Items) > maxBatchSize {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": fmt.Sprintf("Batch size exceeds maximum of %d", maxBatchSize),
		})
	}

	ctx := context.Background()
	updated := make([]interface{}, 0, len(req.Items))
	errors := make([]map[string]interface{}, 0)

	// Update each item
	for i, item := range req.Items {
		// Get existing record for permission check
		existing, err := h.entHelper.Get(ctx, h.modelMeta.Name, item.ID)
		if err != nil {
			errors = append(errors, map[string]interface{}{
				"index": i,
				"id":    item.ID,
				"error": fmt.Sprintf("Record not found: %v", err),
			})
			continue
		}

		// Check rule-based permissions with the existing resource
		permCtx.Resource = existing
		allowed, err = h.permissionChecker.CheckPermission(permCtx)
		if err != nil || !allowed {
			errors = append(errors, map[string]interface{}{
				"index": i,
				"id":    item.ID,
				"error": "Permission denied",
			})
			continue
		}

		// Execute before hook
		hookCtx := &HookContext{
			Model:    h.modelMeta,
			Action:   "update",
			User:     GetUserFromContext(c),
			Request:  c,
			Data:     item.Updates,
			Resource: item.ID,
		}
		if err := h.hookRegistry.Execute(h.modelMeta.Name, HookBeforeUpdate, hookCtx); err != nil {
			errors = append(errors, map[string]interface{}{
				"index": i,
				"id":    item.ID,
				"error": err.Error(),
			})
			continue
		}

		// Validate data
		if err := h.validateData(item.Updates); err != nil {
			errors = append(errors, map[string]interface{}{
				"index": i,
				"id":    item.ID,
				"error": fmt.Sprintf("Validation failed: %s", err.Error()),
			})
			continue
		}

		// Update record
		result, err := h.entHelper.Update(ctx, h.modelMeta.Name, item.ID, item.Updates)
		if err != nil {
			errors = append(errors, map[string]interface{}{
				"index": i,
				"id":    item.ID,
				"error": err.Error(),
			})
			continue
		}

		// Execute after hook
		hookCtx.Result = result
		if err := h.hookRegistry.Execute(h.modelMeta.Name, HookAfterUpdate, hookCtx); err != nil {
			// Log error but don't fail the update
		}

		updated = append(updated, h.modelToMap(result))
	}

	response := fiber.Map{
		"data":   updated,
		"count":  len(updated),
		"errors": errors,
	}

	if len(errors) > 0 && len(updated) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	return c.JSON(response)
}

// BulkDelete handles DELETE /admin/api/{model}/bulk - Bulk delete records
func (h *CRUDHandler) BulkDelete(c *fiber.Ctx) error {
	// Check permissions
	permCtx := &Context{
		User:    GetUserFromContext(c),
		Request: c,
		Model:   h.modelMeta,
		Action:  "delete",
	}

	allowed, err := h.permissionChecker.CheckPermission(permCtx)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Permission check failed",
			"details": err.Error(),
		})
	}
	if !allowed {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "You don't have permission to delete this resource",
		})
	}

	// Parse request body
	var req BulkDeleteRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
			"details": err.Error(),
		})
	}

	if len(req.IDs) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "No IDs provided",
		})
	}

	// Limit batch size
	maxBatchSize := 100
	if len(req.IDs) > maxBatchSize {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": fmt.Sprintf("Batch size exceeds maximum of %d", maxBatchSize),
		})
	}

	ctx := context.Background()
	deleted := make([]interface{}, 0, len(req.IDs))
	errors := make([]map[string]interface{}, 0)

	// Delete each item
	for i, id := range req.IDs {
		// Get existing record for permission check
		existing, err := h.entHelper.Get(ctx, h.modelMeta.Name, id)
		if err != nil {
			errors = append(errors, map[string]interface{}{
				"index": i,
				"id":    id,
				"error": fmt.Sprintf("Record not found: %v", err),
			})
			continue
		}

		// Check rule-based permissions with the existing resource
		permCtx.Resource = existing
		allowed, err = h.permissionChecker.CheckPermission(permCtx)
		if err != nil || !allowed {
			errors = append(errors, map[string]interface{}{
				"index": i,
				"id":    id,
				"error": "Permission denied",
			})
			continue
		}

		// Execute before hook
		hookCtx := &HookContext{
			Model:    h.modelMeta,
			Action:   "delete",
			User:     GetUserFromContext(c),
			Request:  c,
			Resource: id,
		}
		if err := h.hookRegistry.Execute(h.modelMeta.Name, HookBeforeDelete, hookCtx); err != nil {
			errors = append(errors, map[string]interface{}{
				"index": i,
				"id":    id,
				"error": err.Error(),
			})
			continue
		}

		// Delete record
		err = h.entHelper.Delete(ctx, h.modelMeta.Name, id)
		if err != nil {
			errors = append(errors, map[string]interface{}{
				"index": i,
				"id":    id,
				"error": err.Error(),
			})
			continue
		}

		// Execute after hook
		if err := h.hookRegistry.Execute(h.modelMeta.Name, HookAfterDelete, hookCtx); err != nil {
			// Log error but don't fail the delete
		}

		deleted = append(deleted, id)
	}

	response := fiber.Map{
		"deleted": deleted,
		"count":    len(deleted),
		"errors":   errors,
	}

	if len(errors) > 0 && len(deleted) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	return c.JSON(response)
}

