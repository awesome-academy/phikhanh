package utils

import "github.com/gin-gonic/gin"

// Cấu trúc chuẩn cho API response
type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// Trả về response thành công
func SuccessResponse(ctx *gin.Context, statusCode int, message string, data interface{}) {
	ctx.JSON(statusCode, APIResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// Trả về response lỗi
func ErrorResponse(ctx *gin.Context, statusCode int, message string) {
	ctx.JSON(statusCode, APIResponse{
		Success: false,
		Message: message,
		Data:    nil,
	})
}
