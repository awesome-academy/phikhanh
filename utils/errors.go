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
	StatusCode int // dùng cho user controllers: svcErr.StatusCode
	Code       int // dùng cho admin controllers: appErr.Code
	Message    string
}

func (e *AppError) Error() string {
	return e.Message
}

// ServiceError - Alias của AppError để tương thích với user controllers
type ServiceError = AppError

func NewBadRequestError(message string) *AppError {
	return &AppError{Code: http.StatusBadRequest, StatusCode: http.StatusBadRequest, Message: message}
}

func NewNotFoundError(message string) *AppError {
	return &AppError{Code: http.StatusNotFound, StatusCode: http.StatusNotFound, Message: message}
}

func NewUnauthorizedError(message string) *AppError {
	return &AppError{Code: http.StatusUnauthorized, StatusCode: http.StatusUnauthorized, Message: message}
}

func NewForbiddenError(message string) *AppError {
	return &AppError{Code: http.StatusForbidden, StatusCode: http.StatusForbidden, Message: message}
}

func NewInternalServerError(err error) *AppError {
	return &AppError{
		Code:       http.StatusInternalServerError,
		StatusCode: http.StatusInternalServerError,
		Message:    "An internal error occurred. Please try again later.",
	}
}

func NewTooManyRequestsError(message string) *AppError {
	return &AppError{Code: http.StatusTooManyRequests, StatusCode: http.StatusTooManyRequests, Message: message}
}

// Aliases cho user services
func NewServiceBadRequest(message string) *AppError   { return NewBadRequestError(message) }
func NewServiceNotFound(message string) *AppError     { return NewNotFoundError(message) }
func NewServiceUnauthorized(message string) *AppError { return NewUnauthorizedError(message) }
func NewServiceInternalError(err error) *AppError     { return NewInternalServerError(err) }
