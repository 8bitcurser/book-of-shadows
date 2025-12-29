package errors

import (
	"errors"
	"fmt"
)

// Common errors
var (
	// Storage errors
	ErrNotFound        = errors.New("resource not found")
	ErrAlreadyExists   = errors.New("resource already exists")
	ErrDatabaseClosed  = errors.New("database connection is closed")
	ErrInvalidData     = errors.New("invalid data provided")

	// Cookie errors
	ErrCookieNotFound  = errors.New("cookie not found")
	ErrCookieTooLarge  = errors.New("cookie size exceeds limit")
	ErrInvalidCookie   = errors.New("invalid cookie format")

	// Investigator errors
	ErrInvalidAttribute = errors.New("invalid attribute")
	ErrInvalidSkill    = errors.New("invalid skill")
	ErrInvalidMode     = errors.New("invalid game mode")
)

// ValidationError represents a validation error with field information
type ValidationError struct {
	Field   string
	Message string
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("validation error on field '%s': %s", e.Field, e.Message)
}

// NewValidationError creates a new ValidationError
func NewValidationError(field, message string) ValidationError {
	return ValidationError{
		Field:   field,
		Message: message,
	}
}

// HTTPError represents an HTTP error with status code
type HTTPError struct {
	Code    int
	Message string
	Err     error
}

func (e HTTPError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("HTTP %d: %s: %v", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("HTTP %d: %s", e.Code, e.Message)
}

func (e HTTPError) Unwrap() error {
	return e.Err
}

// NewHTTPError creates a new HTTPError
func NewHTTPError(code int, message string, err error) HTTPError {
	return HTTPError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}