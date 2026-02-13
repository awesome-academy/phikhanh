package utils

import "errors"

// ServiceError - Struct chứa status code và message để trả về từ service
type ServiceError struct {
	StatusCode int
	Message    string
	Err        error
}

func (e *ServiceError) Error() string {
	return e.Message
}

// NewServiceError - Tạo ServiceError mới
func NewServiceError(statusCode int, message string, err error) *ServiceError {
	return &ServiceError{
		StatusCode: statusCode,
		Message:    message,
		Err:        err,
	}
}

// Custom errors
var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrCitizenIDExists    = errors.New("citizen_id already exists")
	ErrEmailExists        = errors.New("email already exists")
	ErrUnauthorized       = errors.New("unauthorized")
	ErrUserNotFound       = errors.New("user not found")
	ErrServiceNotFound    = errors.New("service not found")
	ErrInvalidInput       = errors.New("invalid input")
	ErrInternalServer     = errors.New("internal server error")
)

// Predefined ServiceErrors
func ErrInvalidCredentialsResponse() *ServiceError {
	return NewServiceError(401, "Invalid citizen ID or password", ErrInvalidCredentials)
}

func ErrCitizenIDExistsResponse() *ServiceError {
	return NewServiceError(400, "Citizen ID already exists", ErrCitizenIDExists)
}

func ErrEmailExistsResponse() *ServiceError {
	return NewServiceError(400, "Email already exists", ErrEmailExists)
}

func ErrUnauthorizedResponse() *ServiceError {
	return NewServiceError(401, "Unauthorized access", ErrUnauthorized)
}

func ErrUserNotFoundResponse() *ServiceError {
	return NewServiceError(404, "User not found", ErrUserNotFound)
}

func ErrServiceNotFoundResponse() *ServiceError {
	return NewServiceError(404, "Service not found", ErrServiceNotFound)
}

func ErrInvalidInputResponse() *ServiceError {
	return NewServiceError(400, "Invalid input data", ErrInvalidInput)
}

func ErrInternalServerResponse() *ServiceError {
	return NewServiceError(500, "An error occurred, please try again later", ErrInternalServer)
}
