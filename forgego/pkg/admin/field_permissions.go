package admin

import (
	"context"
)

type FieldPermissionChecker struct {
	permissionChecker *PermissionChecker
}

func NewFieldPermissionChecker(permissionChecker *PermissionChecker) *FieldPermissionChecker {
	return &FieldPermissionChecker{
		permissionChecker: permissionChecker,
	}
}

func (fpc *FieldPermissionChecker) CanViewField(ctx context.Context, user interface{}, modelMeta *ModelMeta, fieldName string, obj interface{}) bool {
	permCtx := &Context{
		User:     user,
		Request:  nil,
		Model:    modelMeta,
		Action:   "view_field",
		Resource: obj,
		Field:    fieldName,
	}
	
	allowed, err := fpc.permissionChecker.CheckPermission(permCtx)
	if err != nil {
		return false
	}
	
	return allowed
}

func (fpc *FieldPermissionChecker) CanEditField(ctx context.Context, user interface{}, modelMeta *ModelMeta, fieldName string, obj interface{}) bool {
	permCtx := &Context{
		User:     user,
		Request:  nil,
		Model:    modelMeta,
		Action:   "edit_field",
		Resource: obj,
		Field:    fieldName,
	}
	
	allowed, err := fpc.permissionChecker.CheckPermission(permCtx)
	if err != nil {
		return false
	}
	
	return allowed
}

func (fpc *FieldPermissionChecker) GetReadOnlyFields(ctx context.Context, user interface{}, modelMeta *ModelMeta, obj interface{}) []string {
	readOnlyFields := make([]string, 0)
	
	for _, field := range modelMeta.Fields {
		if field.ReadOnly {
			readOnlyFields = append(readOnlyFields, field.Name)
			continue
		}
		
		if !fpc.CanEditField(ctx, user, modelMeta, field.Name, obj) {
			readOnlyFields = append(readOnlyFields, field.Name)
		}
	}
	
	return readOnlyFields
}

func (fpc *FieldPermissionChecker) GetVisibleFields(ctx context.Context, user interface{}, modelMeta *ModelMeta, obj interface{}) []string {
	visibleFields := make([]string, 0)
	
	for _, field := range modelMeta.Fields {
		if fpc.CanViewField(ctx, user, modelMeta, field.Name, obj) {
			visibleFields = append(visibleFields, field.Name)
		}
	}
	
	return visibleFields
}

func (fpc *FieldPermissionChecker) FilterFields(ctx context.Context, user interface{}, modelMeta *ModelMeta, obj interface{}, fields []string) []string {
	filtered := make([]string, 0)
	
	for _, fieldName := range fields {
		if fpc.CanViewField(ctx, user, modelMeta, fieldName, obj) {
			filtered = append(filtered, fieldName)
		}
	}
	
	return filtered
}

