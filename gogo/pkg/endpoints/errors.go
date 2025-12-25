package endpoints

import (
	"github.com/gofiber/fiber/v2"
)

// Error represents an API error
type Error struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details interface{} `json:"details,omitempty"`
}

// Error implements error interface
func (e *Error) Error() string {
	return e.Message
}

// Common errors
var (
	ErrNotFound     = &Error{Code: "not_found", Message: "Resource not found"}
	ErrUnauthorized = &Error{Code: "unauthorized", Message: "Unauthorized"}
	ErrForbidden    = &Error{Code: "forbidden", Message: "Forbidden"}
	ErrBadRequest   = &Error{Code: "bad_request", Message: "Bad request"}
	ErrValidation   = &Error{Code: "validation_error", Message: "Validation failed"}
	ErrInternal     = &Error{Code: "internal_error", Message: "Internal server error"}
)

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error Error `json:"error"`
}

// HandleError handles errors and returns appropriate HTTP response
func HandleError(c *fiber.Ctx, err error) error {
	var apiErr *Error
	
	switch e := err.(type) {
	case *Error:
		apiErr = e
	default:
		apiErr = ErrInternal
		apiErr.Details = err.Error()
	}
	
	statusCode := getStatusCode(apiErr)
	
	return c.Status(statusCode).JSON(ErrorResponse{
		Error: *apiErr,
	})
}

// getStatusCode returns HTTP status code for an error
func getStatusCode(err *Error) int {
	switch err.Code {
	case "not_found":
		return fiber.StatusNotFound
	case "unauthorized":
		return fiber.StatusUnauthorized
	case "forbidden":
		return fiber.StatusForbidden
	case "bad_request", "validation_error":
		return fiber.StatusBadRequest
	default:
		return fiber.StatusInternalServerError
	}
}

// NewError creates a new error
func NewError(code, message string) *Error {
	return &Error{
		Code:    code,
		Message: message,
	}
}

// NewValidationError creates a validation error with details
func NewValidationError(details interface{}) *Error {
	return &Error{
		Code:    "validation_error",
		Message: "Validation failed",
		Details: details,
	}
}

