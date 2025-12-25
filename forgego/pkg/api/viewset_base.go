package api

import (
	"database/sql"
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/forgego/forge/pkg/models"
)

func isNotFoundError(err error) bool {
	return errors.Is(err, sql.ErrNoRows)
}

type BaseViewSet[T any] struct {
	manager     models.Manager[T]
	serializer  Serializer[T]
	filters     []FilterBackend[T]
	pagination  PaginationBackend[T]
	permissions []Permission
	listMixin   *ListMixin[T]
	retrieveMixin *RetrieveMixin[T]
	createMixin *CreateMixin[T]
	updateMixin *UpdateMixin[T]
	destroyMixin *DestroyMixin[T]
}

func NewBaseViewSet[T any](manager models.Manager[T], serializer Serializer[T]) *BaseViewSet[T] {
	return &BaseViewSet[T]{
		manager:      manager,
		serializer:   serializer,
		filters:      []FilterBackend[T]{},
		pagination:   nil,
		permissions:  []Permission{},
		listMixin:    &ListMixin[T]{},
		retrieveMixin: &RetrieveMixin[T]{},
		createMixin:  &CreateMixin[T]{},
		updateMixin:  &UpdateMixin[T]{},
		destroyMixin: &DestroyMixin[T]{},
	}
}

func (vs *BaseViewSet[T]) WithFilters(filters ...FilterBackend[T]) *BaseViewSet[T] {
	vs.filters = append(vs.filters, filters...)
	return vs
}

func (vs *BaseViewSet[T]) WithPagination(pagination PaginationBackend[T]) *BaseViewSet[T] {
	vs.pagination = pagination
	return vs
}

func (vs *BaseViewSet[T]) WithPermissions(permissions ...Permission) *BaseViewSet[T] {
	vs.permissions = append(vs.permissions, permissions...)
	return vs
}

func (vs *BaseViewSet[T]) checkPermissions(c *fiber.Ctx) error {
	for _, perm := range vs.permissions {
		if !perm.HasPermission(c) {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "You don't have permission to perform this action",
			})
		}
	}
	return nil
}

func (vs *BaseViewSet[T]) List(c *fiber.Ctx) error {
	if err := vs.checkPermissions(c); err != nil {
		return err
	}
	return vs.listMixin.List(c, vs.manager, vs.serializer, vs.filters, vs.pagination)
}

func (vs *BaseViewSet[T]) Retrieve(c *fiber.Ctx) error {
	if err := vs.checkPermissions(c); err != nil {
		return err
	}
	return vs.retrieveMixin.Retrieve(c, vs.manager, vs.serializer)
}

func (vs *BaseViewSet[T]) Create(c *fiber.Ctx) error {
	if err := vs.checkPermissions(c); err != nil {
		return err
	}
	return vs.createMixin.Create(c, vs.manager, vs.serializer)
}

func (vs *BaseViewSet[T]) Update(c *fiber.Ctx) error {
	if err := vs.checkPermissions(c); err != nil {
		return err
	}
	return vs.updateMixin.Update(c, vs.manager, vs.serializer)
}

func (vs *BaseViewSet[T]) Destroy(c *fiber.Ctx) error {
	if err := vs.checkPermissions(c); err != nil {
		return err
	}
	return vs.destroyMixin.Destroy(c, vs.manager)
}

func (vs *BaseViewSet[T]) GetManager() models.Manager[T] {
	return vs.manager
}

func (vs *BaseViewSet[T]) GetSerializer() Serializer[T] {
	return vs.serializer
}


