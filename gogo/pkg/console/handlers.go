package console

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func listModel(c *fiber.Ctx) error {
	modelName := c.Params("model")

	consoleInfo, ok := globalRegistry.consoles[modelName]
	if !ok {
		return c.Status(404).JSON(fiber.Map{
			"error": "Model not found",
		})
	}

	page, _ := strconv.Atoi(c.Query("page", "1"))
	pageSize, _ := strconv.Atoi(c.Query("page_size", "20"))
	search := c.Query("search", "")

	response := fiber.Map{
		"model":        modelName,
		"page":         page,
		"page_size":    pageSize,
		"search":       search,
		"list_display": consoleInfo.Console.ListDisplay(),
		"filters":      consoleInfo.Console.ListFilters(),
		"actions":      consoleInfo.Console.Actions(),
		"items":        []interface{}{},
	}

	return c.JSON(response)
}

func showModel(c *fiber.Ctx) error {
	modelName := c.Params("model")
	id := c.Params("id")

	_, ok := globalRegistry.consoles[modelName]
	if !ok {
		return c.Status(404).JSON(fiber.Map{
			"error": "Model not found",
		})
	}

	response := fiber.Map{
		"model": modelName,
		"id":    id,
		"item":  nil,
	}

	return c.JSON(response)
}

func newModel(c *fiber.Ctx) error {
	modelName := c.Params("model")

	_, ok := globalRegistry.consoles[modelName]
	if !ok {
		return c.Status(404).JSON(fiber.Map{
			"error": "Model not found",
		})
	}

	response := fiber.Map{
		"model": modelName,
		"form": fiber.Map{
			"fields": []interface{}{},
		},
	}

	return c.JSON(response)
}

func createModel(c *fiber.Ctx) error {
	modelName := c.Params("model")

	_, ok := globalRegistry.consoles[modelName]
	if !ok {
		return c.Status(404).JSON(fiber.Map{
			"error": "Model not found",
		})
	}

	var data map[string]interface{}
	if err := c.BodyParser(&data); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid request data",
		})
	}

	response := fiber.Map{
		"model":   modelName,
		"item":    data,
		"message": "Created successfully",
	}

	return c.Status(201).JSON(response)
}

func editModel(c *fiber.Ctx) error {
	modelName := c.Params("model")
	id := c.Params("id")

	_, ok := globalRegistry.consoles[modelName]
	if !ok {
		return c.Status(404).JSON(fiber.Map{
			"error": "Model not found",
		})
	}

	response := fiber.Map{
		"model": modelName,
		"id":    id,
		"form": fiber.Map{
			"fields": []interface{}{},
			"item":   nil,
		},
	}

	return c.JSON(response)
}

func updateModel(c *fiber.Ctx) error {
	modelName := c.Params("model")
	id := c.Params("id")

	_, ok := globalRegistry.consoles[modelName]
	if !ok {
		return c.Status(404).JSON(fiber.Map{
			"error": "Model not found",
		})
	}

	var data map[string]interface{}
	if err := c.BodyParser(&data); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid request data",
		})
	}

	response := fiber.Map{
		"model":   modelName,
		"id":      id,
		"item":    data,
		"message": "Updated successfully",
	}

	return c.JSON(response)
}

func deleteModel(c *fiber.Ctx) error {
	modelName := c.Params("model")
	id := c.Params("id")

	_, ok := globalRegistry.consoles[modelName]
	if !ok {
		return c.Status(404).JSON(fiber.Map{
			"error": "Model not found",
		})
	}

	response := fiber.Map{
		"model":   modelName,
		"id":      id,
		"message": "Deleted successfully",
	}

	return c.JSON(response)
}
