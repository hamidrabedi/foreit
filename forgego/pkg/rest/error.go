package rest

import (
	"net/http"
	
	"github.com/gofiber/fiber/v2"
)

type Error struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Status  int    `json:"-"`
}

func NewError(status int, code, message string) *Error {
	return &Error{
		Code:    code,
		Message: message,
		Status:  status,
	}
}

func (e *Error) Error() string {
	return e.Message
}

func HandleError(c *fiber.Ctx, err error) error {
	if apiErr, ok := err.(*Error); ok {
		return c.Status(apiErr.Status).JSON(fiber.Map{
			"error": fiber.Map{
				"code":    apiErr.Code,
				"message": apiErr.Message,
			},
		})
	}
	
	return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
		"error": fiber.Map{
			"code":    "internal_error",
			"message": err.Error(),
		},
	})
}

var (
	ErrNotFound      = NewError(http.StatusNotFound, "not_found", "Resource not found")
	ErrBadRequest    = NewError(http.StatusBadRequest, "bad_request", "Invalid request")
	ErrUnauthorized  = NewError(http.StatusUnauthorized, "unauthorized", "Unauthorized")
	ErrForbidden     = NewError(http.StatusForbidden, "forbidden", "Forbidden")
	ErrInternalError = NewError(http.StatusInternalServerError, "internal_error", "Internal server error")
)

