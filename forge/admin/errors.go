package admin

import (
	"fmt"
)

// AdminError represents an admin-specific error
type AdminError struct {
	Code    string
	Message string
	Err     error
}

// Error implements error interface
func (e *AdminError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %s (%v)", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// Unwrap returns the underlying error
func (e *AdminError) Unwrap() error {
	return e.Err
}

// Error codes
const (
	ErrModelNotFound    = "MODEL_NOT_FOUND"
	ErrInstanceNotFound = "INSTANCE_NOT_FOUND"
	ErrPermissionDenied = "PERMISSION_DENIED"
	ErrValidationFailed = "VALIDATION_FAILED"
	ErrInvalidOperation = "INVALID_OPERATION"
	ErrManagerNotSet    = "MANAGER_NOT_SET"
)

// NewAdminError creates a new admin error
func NewAdminError(code, message string, err error) *AdminError {
	return &AdminError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

// IsModelNotFound checks if error is model not found
func IsModelNotFound(err error) bool {
	if adminErr, ok := err.(*AdminError); ok {
		return adminErr.Code == ErrModelNotFound
	}
	return false
}

// IsInstanceNotFound checks if error is instance not found
func IsInstanceNotFound(err error) bool {
	if adminErr, ok := err.(*AdminError); ok {
		return adminErr.Code == ErrInstanceNotFound
	}
	return false
}

// IsPermissionDenied checks if error is permission denied
func IsPermissionDenied(err error) bool {
	if adminErr, ok := err.(*AdminError); ok {
		return adminErr.Code == ErrPermissionDenied
	}
	return false
}
