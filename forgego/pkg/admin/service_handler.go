package admin

import (
	"github.com/gofiber/fiber/v2"
	"github.com/forgego/forge/pkg/api"
	"github.com/forgego/forge/pkg/service"
)

type ServiceCRUDHandler[T any] struct {
	modelMeta *ModelMeta
	service   service.ResourceService[T]
	registry  *Registry
}

func NewServiceCRUDHandler[T any](meta *ModelMeta, svc service.ResourceService[T], registry *Registry) *ServiceCRUDHandler[T] {
	return &ServiceCRUDHandler[T]{
		modelMeta: meta,
		service:   svc,
		registry:  registry,
	}
}

func (h *ServiceCRUDHandler[T]) List(c *fiber.Ctx) error {
	ctx := c.UserContext()
	params := service.ParseListParams(c)
	
	result, err := h.service.List(ctx, params)
	if err != nil {
		if err == service.ErrPermissionDenied {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "You don't have permission to list this resource",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to list records",
			"details": err.Error(),
		})
	}
	
	serializer := h.getSerializer()
	serialized := make([]interface{}, 0, len(result.Items))
	for _, obj := range result.Items {
		objPtr := &obj
		rep, err := serializer.ToRepresentation(objPtr)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to serialize result",
				"details": err.Error(),
			})
		}
		serialized = append(serialized, rep)
	}
	
	return c.JSON(fiber.Map{
		"count":   result.Total,
		"results": serialized,
	})
}

func (h *ServiceCRUDHandler[T]) Get(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "ID is required",
		})
	}
	
	ctx := c.UserContext()
	obj, err := h.service.Retrieve(ctx, id)
	if err != nil {
		if err == service.ErrPermissionDenied {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "You don't have permission to retrieve this resource",
			})
		}
		if err == service.ErrNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Record not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve record",
			"details": err.Error(),
		})
	}
	
	serializer := h.getSerializer()
	objPtr := &obj
	rep, err := serializer.ToRepresentation(objPtr)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to serialize result",
			"details": err.Error(),
		})
	}
	
	return c.JSON(rep)
}

func (h *ServiceCRUDHandler[T]) Create(c *fiber.Ctx) error {
	ctx := c.UserContext()
	body := c.Body()
	
	obj, err := h.service.Create(ctx, body)
	if err != nil {
		if err == service.ErrPermissionDenied {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "You don't have permission to create this resource",
			})
		}
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to create record",
			"details": err.Error(),
		})
	}
	
	serializer := h.getSerializer()
	objPtr := &obj
	rep, err := serializer.ToRepresentation(objPtr)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to serialize result",
			"details": err.Error(),
		})
	}
	
	return c.Status(fiber.StatusCreated).JSON(rep)
}

func (h *ServiceCRUDHandler[T]) Update(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "ID is required",
		})
	}
	
	ctx := c.UserContext()
	body := c.Body()
	
	obj, err := h.service.Update(ctx, id, body)
	if err != nil {
		if err == service.ErrPermissionDenied {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "You don't have permission to update this resource",
			})
		}
		if err == service.ErrNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Record not found",
			})
		}
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to update record",
			"details": err.Error(),
		})
	}
	
	serializer := h.getSerializer()
	objPtr := &obj
	rep, err := serializer.ToRepresentation(objPtr)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to serialize result",
			"details": err.Error(),
		})
	}
	
	return c.JSON(rep)
}

func (h *ServiceCRUDHandler[T]) Delete(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "ID is required",
		})
	}
	
	ctx := c.UserContext()
	err := h.service.Delete(ctx, id)
	if err != nil {
		if err == service.ErrPermissionDenied {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "You don't have permission to delete this resource",
			})
		}
		if err == service.ErrNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Record not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete record",
			"details": err.Error(),
		})
	}
	
	return c.SendStatus(fiber.StatusNoContent)
}

func (h *ServiceCRUDHandler[T]) getSerializer() api.Serializer[T] {
	if svcWithSerializer, ok := h.service.(service.ServiceWithSerializer[T]); ok {
		return svcWithSerializer.Serializer()
	}
	return nil
}

