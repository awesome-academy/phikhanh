package utils

import "fmt"

// ServiceError - Struct chứa status code và message để trả về từ service
type ServiceError struct {
	StatusCode int
	Message    string
	Err        error // Lưu error gốc để log
}

func (e *ServiceError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

// NewServiceError - Tạo ServiceError với status code và message tùy chỉnh
func NewServiceError(statusCode int, message string) *ServiceError {
	return &ServiceError{
		StatusCode: statusCode,
		Message:    message,
	}
}

// Common error constructors
func NewBadRequestError(message string) *ServiceError {
	return NewServiceError(400, message)
}

func NewUnauthorizedError(message string) *ServiceError {
	return NewServiceError(401, message)
}

func NewNotFoundError(message string) *ServiceError {
	return NewServiceError(404, message)
}

// NewInternalServerError - Trả về error 500 với message cụ thể từ error
func NewInternalServerError(err error) *ServiceError {
	return &ServiceError{
		StatusCode: 500,
		Message:    "An error occurred, please try again later",
		Err:        err,
	}
}
