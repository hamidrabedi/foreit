package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/forgego/forge/pkg/api"
	"github.com/forgego/forge/pkg/models"
)

type ResourceService[T any] interface {
	List(ctx context.Context, params ListParams) (*PaginationResult[T], error)
	Retrieve(ctx context.Context, id any) (T, error)
	Create(ctx context.Context, input any) (T, error)
	Update(ctx context.Context, id any, input any) (T, error)
	Delete(ctx context.Context, id any) error
}

type ServiceWithSerializer[T any] interface {
	ResourceService[T]
	Serializer() api.Serializer[T]
}

type ListParams struct {
	Predicates []models.Condition
	Search     string
	SearchFields []string
	Sort       []OrderBy
	Page       int
	PageSize   int
	Cursor     *string
	Limit      int
}

type OrderBy struct {
	Field     string
	Direction string
}

type PaginationResult[T any] struct {
	Items      []T
	Total      int64
	NextCursor *string
	HasMore    bool
}

type Permission interface {
	Allow(ctx context.Context, action string) bool
}

type Filter[T any] interface {
	Apply(ctx context.Context, qs models.QuerySet[T], params ListParams) models.QuerySet[T]
}

type Pagination[T any] interface {
	Paginate(ctx context.Context, qs models.QuerySet[T], params ListParams) (*PaginationResult[T], error)
}

type BaseResourceService[T any] struct {
	manager     models.Manager[T]
	serializer  api.Serializer[T]
	filters     []Filter[T]
	pagination  Pagination[T]
	permissions []Permission
	
	beforeValidate func(ctx context.Context, input any) error
	afterValidate  func(ctx context.Context, input any) error
	beforeCreate   func(ctx context.Context, input any) error
	afterCreate    func(ctx context.Context, obj T) error
	beforeUpdate   func(ctx context.Context, id any, input any) error
	afterUpdate    func(ctx context.Context, obj T) error
	beforeDelete   func(ctx context.Context, id any) error
	afterDelete    func(ctx context.Context, id any) error
}

func NewBaseResourceService[T any](
	manager models.Manager[T],
	serializer api.Serializer[T],
) *BaseResourceService[T] {
	return &BaseResourceService[T]{
		manager:     manager,
		serializer:  serializer,
		filters:     []Filter[T]{},
		pagination:  &DefaultPagination[T]{},
		permissions: []Permission{},
	}
}

func (s *BaseResourceService[T]) Serializer() api.Serializer[T] {
	return s.serializer
}

func (s *BaseResourceService[T]) WithFilters(filters ...Filter[T]) *BaseResourceService[T] {
	s.filters = append(s.filters, filters...)
	return s
}

func (s *BaseResourceService[T]) WithPagination(pagination Pagination[T]) *BaseResourceService[T] {
	s.pagination = pagination
	return s
}

func (s *BaseResourceService[T]) WithPermissions(permissions ...Permission) *BaseResourceService[T] {
	s.permissions = append(s.permissions, permissions...)
	return s
}

func (s *BaseResourceService[T]) WithBeforeValidate(hook func(ctx context.Context, input any) error) *BaseResourceService[T] {
	s.beforeValidate = hook
	return s
}

func (s *BaseResourceService[T]) WithAfterValidate(hook func(ctx context.Context, input any) error) *BaseResourceService[T] {
	s.afterValidate = hook
	return s
}

func (s *BaseResourceService[T]) WithBeforeCreate(hook func(ctx context.Context, input any) error) *BaseResourceService[T] {
	s.beforeCreate = hook
	return s
}

func (s *BaseResourceService[T]) WithAfterCreate(hook func(ctx context.Context, obj T) error) *BaseResourceService[T] {
	s.afterCreate = hook
	return s
}

func (s *BaseResourceService[T]) WithBeforeUpdate(hook func(ctx context.Context, id any, input any) error) *BaseResourceService[T] {
	s.beforeUpdate = hook
	return s
}

func (s *BaseResourceService[T]) WithAfterUpdate(hook func(ctx context.Context, obj T) error) *BaseResourceService[T] {
	s.afterUpdate = hook
	return s
}

func (s *BaseResourceService[T]) WithBeforeDelete(hook func(ctx context.Context, id any) error) *BaseResourceService[T] {
	s.beforeDelete = hook
	return s
}

func (s *BaseResourceService[T]) WithAfterDelete(hook func(ctx context.Context, id any) error) *BaseResourceService[T] {
	s.afterDelete = hook
	return s
}

// List implements ResourceService.List.
// Permission checks happen first to fail fast before any database queries.
func (s *BaseResourceService[T]) List(ctx context.Context, params ListParams) (*PaginationResult[T], error) {
	for _, p := range s.permissions {
		if !p.Allow(ctx, "list") {
			return nil, ErrPermissionDenied
		}
	}

	qs := s.manager.All()

	if len(params.Predicates) > 0 {
		qs = qs.Filter(params.Predicates...)
	}

	for _, f := range s.filters {
		qs = f.Apply(ctx, qs, params)
	}

	if params.Search != "" {
		qs = s.applySearch(qs, params.Search, params.SearchFields)
	}

	if len(params.Sort) > 0 {
		orders := make([]models.OrderBy, len(params.Sort))
		for i, sort := range params.Sort {
			direction := models.OrderAsc
			if sort.Direction == "DESC" || sort.Direction == "desc" {
				direction = models.OrderDesc
			}
			orders[i] = models.OrderBy{
				Field:     sort.Field,
				Direction: direction,
			}
		}
		qs = qs.OrderBy(orders...)
	}

	if params.Cursor != nil {
		return s.listWithCursor(ctx, qs, params)
	}
	return s.listWithOffset(ctx, qs, params)
}

// Retrieve implements ResourceService.Retrieve.
// Assumes primary key field is named "id" - this should be made configurable.
func (s *BaseResourceService[T]) Retrieve(ctx context.Context, id any) (T, error) {
	var zero T

	for _, p := range s.permissions {
		if !p.Allow(ctx, "retrieve") {
			return zero, ErrPermissionDenied
		}
	}

	idCondition := newStringCondition("id", "=", id)
	result, err := s.manager.Get(ctx, idCondition)
	if err != nil {
		return zero, fmt.Errorf("failed to retrieve: %w", err)
	}

	return *result, nil
}

// Create implements ResourceService.Create.
// Lifecycle: BeforeValidate → Validate (Serializer) → AfterValidate → BeforeCreate → Create → AfterCreate
// All validation happens in Service - ViewSet never validates.
func (s *BaseResourceService[T]) Create(ctx context.Context, input any) (T, error) {
	var zero T

	for _, p := range s.permissions {
		if !p.Allow(ctx, "create") {
			return zero, ErrPermissionDenied
		}
	}

	if s.beforeValidate != nil {
		if err := s.beforeValidate(ctx, input); err != nil {
			return zero, fmt.Errorf("before validate failed: %w", err)
		}
	}

	var body []byte
	switch v := input.(type) {
	case []byte:
		body = v
	case string:
		body = []byte(v)
	default:
		return zero, fmt.Errorf("invalid input type: expected []byte or string")
	}

	obj, err := s.serializer.FromCreate(body)
	if err != nil {
		return zero, fmt.Errorf("validation failed: %w", err)
	}

	if s.afterValidate != nil {
		if err := s.afterValidate(ctx, input); err != nil {
			return zero, fmt.Errorf("after validate failed: %w", err)
		}
	}

	if s.beforeCreate != nil {
		if err := s.beforeCreate(ctx, input); err != nil {
			return zero, fmt.Errorf("before create failed: %w", err)
		}
	}

	err = s.manager.Create(ctx, obj)
	if err != nil {
		return zero, fmt.Errorf("failed to create: %w", err)
	}

	result := *obj

	if s.afterCreate != nil {
		if err := s.afterCreate(ctx, result); err != nil {
			// AfterCreate errors don't fail the request (already persisted)
		}
	}

	return result, nil
}

// Update implements ResourceService.Update.
// Lifecycle: BeforeValidate → Validate (Serializer) → AfterValidate → BeforeUpdate → Update → AfterUpdate
func (s *BaseResourceService[T]) Update(ctx context.Context, id any, input any) (T, error) {
	var zero T

	for _, p := range s.permissions {
		if !p.Allow(ctx, "update") {
			return zero, ErrPermissionDenied
		}
	}

	idCondition := newStringCondition("id", "=", id)
	existing, err := s.manager.Get(ctx, idCondition)
	if err != nil {
		return zero, fmt.Errorf("failed to retrieve existing record: %w", err)
	}

	if s.beforeValidate != nil {
		if err := s.beforeValidate(ctx, input); err != nil {
			return zero, fmt.Errorf("before validate failed: %w", err)
		}
	}

	var body []byte
	switch v := input.(type) {
	case []byte:
		body = v
	case string:
		body = []byte(v)
	default:
		return zero, fmt.Errorf("invalid input type: expected []byte or string")
	}

	err = s.serializer.FromUpdate(existing, body)
	if err != nil {
		return zero, fmt.Errorf("validation failed: %w", err)
	}

	if s.afterValidate != nil {
		if err := s.afterValidate(ctx, input); err != nil {
			return zero, fmt.Errorf("after validate failed: %w", err)
		}
	}

	if s.beforeUpdate != nil {
		if err := s.beforeUpdate(ctx, id, input); err != nil {
			return zero, fmt.Errorf("before update failed: %w", err)
		}
	}

	err = s.manager.Update(ctx, existing)
	if err != nil {
		return zero, fmt.Errorf("failed to update: %w", err)
	}

	result := *existing

	if s.afterUpdate != nil {
		if err := s.afterUpdate(ctx, result); err != nil {
			// AfterUpdate errors don't fail the request (already persisted)
		}
	}

	return result, nil
}

// Delete implements ResourceService.Delete.
// Lifecycle: BeforeDelete → Delete → AfterDelete
func (s *BaseResourceService[T]) Delete(ctx context.Context, id any) error {
	for _, p := range s.permissions {
		if !p.Allow(ctx, "delete") {
			return ErrPermissionDenied
		}
	}

	if s.beforeDelete != nil {
		if err := s.beforeDelete(ctx, id); err != nil {
			return fmt.Errorf("before delete failed: %w", err)
		}
	}

	idCondition := newStringCondition("id", "=", id)
	obj, err := s.manager.Get(ctx, idCondition)
	if err != nil {
		return fmt.Errorf("failed to retrieve record: %w", err)
	}

	err = s.manager.Delete(ctx, obj)
	if err != nil {
		return fmt.Errorf("failed to delete: %w", err)
	}

	if s.afterDelete != nil {
		if err := s.afterDelete(ctx, id); err != nil {
			// AfterDelete errors don't fail the request (already deleted)
		}
	}

	return nil
}

// Helper methods

func (s *BaseResourceService[T]) applySearch(qs models.QuerySet[T], search string, fields []string) models.QuerySet[T] {
	if len(fields) == 0 {
		return qs
	}

	var conditions []models.Condition
	for _, field := range fields {
		conditions = append(conditions, newStringCondition(field, "ILIKE", "%"+search+"%"))
	}

	if len(conditions) > 0 {
		qs = qs.Filter(conditions...)
	}

	return qs
}

func (s *BaseResourceService[T]) listWithOffset(ctx context.Context, qs models.QuerySet[T], params ListParams) (*PaginationResult[T], error) {
	total, err := qs.Count(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to count: %w", err)
	}

	offset := (params.Page - 1) * params.PageSize
	if params.PageSize <= 0 {
		params.PageSize = 20
	}
	if params.PageSize > 200 {
		params.PageSize = 200
	}

	results, err := qs.Limit(params.PageSize).Offset(offset).All(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list: %w", err)
	}

	items := make([]T, len(results))
	for i, r := range results {
		items[i] = *r
	}

	totalPages := int(total) / params.PageSize
	if int(total)%params.PageSize > 0 {
		totalPages++
	}
	hasMore := params.Page < totalPages

	return &PaginationResult[T]{
		Items:   items,
		Total:   int64(total),
		HasMore: hasMore,
	}, nil
}

func (s *BaseResourceService[T]) listWithCursor(ctx context.Context, qs models.QuerySet[T], params ListParams) (*PaginationResult[T], error) {
	// TODO: Implement proper cursor-based pagination
	return s.listWithOffset(ctx, qs, params)
}

var (
	ErrPermissionDenied = errors.New("permission denied")
	ErrNotFound         = errors.New("resource not found")
)

