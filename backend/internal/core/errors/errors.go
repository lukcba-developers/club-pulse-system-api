package errors

import (
	"fmt"
	"net/http"
)

// ErrorType identifies the kind of error
type ErrorType string

const (
	ErrorTypeValidation   ErrorType = "VALIDATION_ERROR"
	ErrorTypeNotFound     ErrorType = "NOT_FOUND"
	ErrorTypeUnauthorized ErrorType = "UNAUTHORIZED"
	ErrorTypeForbidden    ErrorType = "FORBIDDEN"
	ErrorTypeInternal     ErrorType = "INTERNAL_SERVER_ERROR"
	ErrorTypeConflict     ErrorType = "CONFLICT"
	ErrorTypeTooManyReqs  ErrorType = "TOO_MANY_REQUESTS"
)

// Common error messages (centralized for consistency)
const (
	MsgInvalidInput    = "Invalid input"
	MsgUnauthorized    = "Authorization required"
	MsgInvalidToken    = "Invalid or expired token"
	MsgNotFound        = "Resource not found"
	MsgForbidden       = "Access denied"
	MsgInternalError   = "Internal server error"
	MsgInvalidID       = "Invalid ID format"
	MsgConflict        = "Resource conflict"
	MsgTooManyRequests = "Too many requests, please try again later"
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

// Common error constructors
func NewValidation(msg string) *AppError {
	return New(ErrorTypeValidation, msg)
}

func NewNotFound(resource string) *AppError {
	return New(ErrorTypeNotFound, fmt.Sprintf("%s not found", resource))
}

func NewUnauthorized() *AppError {
	return New(ErrorTypeUnauthorized, MsgUnauthorized)
}

func NewForbidden() *AppError {
	return New(ErrorTypeForbidden, MsgForbidden)
}

func NewInvalidID(idType string) *AppError {
	return New(ErrorTypeValidation, fmt.Sprintf("invalid %s ID", idType))
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
	case ErrorTypeTooManyReqs:
		return http.StatusTooManyRequests
	default:
		return http.StatusInternalServerError
	}
}
