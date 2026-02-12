package utils

import "github.com/gin-gonic/gin"

// Cấu trúc chuẩn cho API response
type APIResponse struct {
	Status  int         `json:"status" example:"200"`
	Message string      `json:"message" example:"Success"`
	Data    interface{} `json:"data,omitempty"`
}

// Cấu trúc response với nhiều lỗi validation
type ValidationErrors struct {
	Status  int               `json:"status" example:"400"`
	Message string            `json:"message" example:"Validation failed"`
	Errors  map[string]string `json:"errors"`
}

// Trả về response thành công
func SuccessResponse(ctx *gin.Context, statusCode int, message string, data interface{}) {
	ctx.JSON(statusCode, APIResponse{
		Status:  statusCode,
		Message: message,
		Data:    data,
	})
}

// Trả về response lỗi
func ErrorResponse(ctx *gin.Context, statusCode int, message string) {
	ctx.JSON(statusCode, APIResponse{
		Status:  statusCode,
		Message: message,
		Data:    nil,
	})
}

// Trả về response lỗi validation với nhiều errors
func ValidationErrorResponse(ctx *gin.Context, errors map[string]string) {
	ctx.JSON(400, ValidationErrors{
		Status:  400,
		Message: "Validation failed",
		Errors:  errors,
	})
}
