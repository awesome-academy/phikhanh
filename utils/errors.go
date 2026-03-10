package utils

import "net/http"

// Common error messages
const (
	MsgUnauthorized       = "Unauthorized"
	MsgInvalidCredentials = "Invalid email or password"
	MsgInternalError      = "An error occurred, please try again later"
	MsgInvalidUserID      = "Invalid user ID"
	MsgInvalidUUIDFormat  = "Invalid UUID format"
)

// AppError - Unified error type cho cả admin và user
type AppError struct {
	StatusCode int
	Message    string
}

func (e *AppError) Error() string {
	return e.Message
}

// ServiceError - Alias của AppError để tương thích với user controllers
type ServiceError = AppError

func NewBadRequestError(message string) *AppError {
	return &AppError{StatusCode: http.StatusBadRequest, Message: message}
}

func NewNotFoundError(message string) *AppError {
	return &AppError{StatusCode: http.StatusNotFound, Message: message}
}

func NewUnauthorizedError(message string) *AppError {
	return &AppError{StatusCode: http.StatusUnauthorized, Message: message}
}

func NewForbiddenError(message string) *AppError {
	return &AppError{StatusCode: http.StatusForbidden, Message: message}
}

func NewInternalServerError(err error) *AppError {
	return &AppError{
		StatusCode: http.StatusInternalServerError,
		Message:    "An internal error occurred. Please try again later.",
	}
}

func NewTooManyRequestsError(message string) *AppError {
	return &AppError{StatusCode: http.StatusTooManyRequests, Message: message}
}
