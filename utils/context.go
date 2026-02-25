package utils

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// ExtractUserID - Lấy và parse user_id từ Gin context an toàn
func ExtractUserID(ctx *gin.Context) (uuid.UUID, *ServiceError) {
	userIDVal, exists := ctx.Get("user_id")
	if !exists {
		return uuid.Nil, NewUnauthorizedError("Unauthorized")
	}

	userIDStr, ok := userIDVal.(string)
	if !ok {
		return uuid.Nil, NewUnauthorizedError("Invalid user ID")
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return uuid.Nil, NewBadRequestError("Invalid user ID format")
	}

	return userID, nil
}
