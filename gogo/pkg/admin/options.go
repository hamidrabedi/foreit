package admin

import "github.com/gogo/internal/admin"

// Re-export options for convenience
var (
	WithPermissions = admin.WithPermissions
	WithFields      = admin.WithFields
	WithFilters     = admin.WithFilters
	WithSorting     = admin.WithSorting
	WithSearch      = admin.WithSearch
	WithTableName   = admin.WithTableName
	WithRule        = admin.WithRule
	WithActions     = admin.WithActions
)

// Re-export types for convenience
type (
	Permissions = admin.Permissions
	FieldsConfig = admin.FieldsConfig
	Action      = admin.Action
	ModelMeta   = admin.ModelMeta
	FieldMeta   = admin.FieldMeta
	FieldType   = admin.FieldType
	Choice      = admin.Choice
	OpenAPISpec = admin.OpenAPISpec
	HookType    = admin.HookType
	HookFunc    = admin.HookFunc
	HookContext = admin.HookContext
)

// Hook type constants
const (
	HookBeforeCreate = admin.HookBeforeCreate
	HookAfterCreate  = admin.HookAfterCreate
	HookBeforeUpdate = admin.HookBeforeUpdate
	HookAfterUpdate  = admin.HookAfterUpdate
	HookBeforeDelete = admin.HookBeforeDelete
	HookAfterDelete  = admin.HookAfterDelete
	HookBeforeList   = admin.HookBeforeList
	HookAfterList    = admin.HookAfterList
	HookBeforeGet    = admin.HookBeforeGet
	HookAfterGet     = admin.HookAfterGet
)

