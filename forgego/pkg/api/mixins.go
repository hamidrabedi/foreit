package api

import (
	"database/sql"
	"errors"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/forgego/forge/pkg/models"
)

type ListMixin[T any] struct{}

func (m *ListMixin[T]) List(c *fiber.Ctx, manager models.Manager[T], serializer Serializer[T], filters []FilterBackend[T], pagination PaginationBackend[T]) error {
	ctx := c.UserContext()
	qs := manager.All()

	for _, filter := range filters {
		qs = filter.Apply(c, qs)
	}

	var results []*T
	var total int64
	var err error

	if pagination != nil {
		results, total, err = pagination.Paginate(ctx, qs, c)
	} else {
		results, err = qs.All(ctx)
		if err == nil {
			total = int64(len(results))
		}
	}

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to list records",
			"details": err.Error(),
		})
	}

	serialized := make([]interface{}, 0, len(results))
	for _, obj := range results {
		rep, err := serializer.ToRepresentation(obj)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to serialize result",
				"details": err.Error(),
			})
		}
		serialized = append(serialized, rep)
	}

	response := fiber.Map{
		"results": serialized,
	}
	if pagination != nil {
		response["count"] = total
	} else {
		response["count"] = len(serialized)
	}

	return c.JSON(response)
}

type RetrieveMixin[T any] struct{}

func (m *RetrieveMixin[T]) Retrieve(c *fiber.Ctx, manager models.Manager[T], serializer Serializer[T]) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "ID is required",
		})
	}

	ctx := c.UserContext()
	condition := parseIDCondition(id)
	obj, err := manager.Get(ctx, condition)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Record not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve record",
			"details": err.Error(),
		})
	}

	rep, err := serializer.ToRepresentation(obj)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to serialize result",
			"details": err.Error(),
		})
	}

	return c.JSON(rep)
}

type CreateMixin[T any] struct{}

func (m *CreateMixin[T]) Create(c *fiber.Ctx, manager models.Manager[T], serializer Serializer[T]) error {
	body := c.Body()
	obj, err := serializer.FromCreate(body)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to parse request",
			"details": err.Error(),
		})
	}

	ctx := c.UserContext()
	if err := manager.Create(ctx, obj); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to create record",
			"details": err.Error(),
		})
	}

	rep, err := serializer.ToRepresentation(obj)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to serialize result",
			"details": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(rep)
}

type UpdateMixin[T any] struct{}

func (m *UpdateMixin[T]) Update(c *fiber.Ctx, manager models.Manager[T], serializer Serializer[T]) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "ID is required",
		})
	}

	ctx := c.UserContext()
	obj, err := manager.Get(ctx, models.NewStringCondition("id", "=", id))

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Record not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve record",
			"details": err.Error(),
		})
	}

	body := c.Body()
	if err := serializer.FromUpdate(obj, body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to parse request",
			"details": err.Error(),
		})
	}

	if err := manager.Update(ctx, obj); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to update record",
			"details": err.Error(),
		})
	}

	rep, err := serializer.ToRepresentation(obj)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to serialize result",
			"details": err.Error(),
		})
	}

	return c.JSON(rep)
}

type DestroyMixin[T any] struct{}

func (m *DestroyMixin[T]) Destroy(c *fiber.Ctx, manager models.Manager[T]) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "ID is required",
		})
	}

	ctx := c.UserContext()
	condition := parseIDCondition(id)
	obj, err := manager.Get(ctx, condition)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Record not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve record",
			"details": err.Error(),
		})
	}

	if err := manager.Delete(ctx, obj); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete record",
			"details": err.Error(),
		})
	}

	return c.SendStatus(fiber.StatusNoContent)
}

// parseIDCondition parses an ID parameter and returns the appropriate condition.
// Tries to parse as int64 first, falls back to string if parsing fails.
func parseIDCondition(id string) models.Condition {
	if intID, err := strconv.ParseInt(id, 10, 64); err == nil {
		return models.NewIntCondition("id", "=", intID)
	}
	return models.NewStringCondition("id", "=", id)
}

