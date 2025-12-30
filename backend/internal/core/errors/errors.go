package errors

import (
	"fmt"
	"net/http"
)

// ErrorType identifies the kind of error
type ErrorType string

const (
	ErrorTypeValidation    ErrorType = "VALIDATION_ERROR"
	ErrorTypeNotFound      ErrorType = "NOT_FOUND"
	ErrorTypeUnauthorized  ErrorType = "UNAUTHORIZED"
	ErrorTypeForbidden     ErrorType = "FORBIDDEN"
	ErrorTypeInternal      ErrorType = "INTERNAL_SERVER_ERROR"
	ErrorTypeConflict      ErrorType = "CONFLICT"
)

// AppError is the standard error struct for the application
type AppError struct {
	Type    ErrorType `json:"type"`
	Message string    `json:"message"`
	Code    int       `json:"code"`
	Err     error     `json:"-"` // Internal error for logging
}

func (e *AppError) Error() string {
	return fmt.Sprintf("[%s] %s", e.Type, e.Message)
}

func New(errType ErrorType, msg string) *AppError {
	return &AppError{
		Type:    errType,
		Message: msg,
		Code:    mapTypeToStatusCode(errType),
	}
}

func Wrap(errType ErrorType, msg string, err error) *AppError {
	return &AppError{
		Type:    errType,
		Message: msg,
		Code:    mapTypeToStatusCode(errType),
		Err:     err,
	}
}

// Validation helper (e.g. for creating validation error details)
func NewValidation(msg string) *AppError {
	return New(ErrorTypeValidation, msg)
}

func mapTypeToStatusCode(t ErrorType) int {
	switch t {
	case ErrorTypeValidation:
		return http.StatusBadRequest
	case ErrorTypeNotFound:
		return http.StatusNotFound
	case ErrorTypeUnauthorized:
		return http.StatusUnauthorized
	case ErrorTypeForbidden:
		return http.StatusForbidden
	case ErrorTypeConflict:
		return http.StatusConflict
	default:
		return http.StatusInternalServerError
	}
}
