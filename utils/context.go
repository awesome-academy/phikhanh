package utils

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// ExtractUserID - Lấy và parse user_id từ Gin context (user API)
func ExtractUserID(ctx *gin.Context) (uuid.UUID, *ServiceError) {
	userIDVal, exists := ctx.Get("user_id")
	if !exists {
		return uuid.Nil, NewUnauthorizedError(MsgUnauthorized)
	}

	userIDStr, ok := userIDVal.(string)
	if !ok {
		return uuid.Nil, NewUnauthorizedError(MsgInvalidUserID)
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return uuid.Nil, NewBadRequestError(MsgInvalidUserID + ": " + MsgInvalidUUIDFormat)
	}

	return userID, nil
}

// ExtractAdminID - Lấy và parse admin_id từ Gin context (admin SSR)
func ExtractAdminID(ctx *gin.Context) (uuid.UUID, *ServiceError) {
	adminIDVal, exists := ctx.Get("admin_id")
	if !exists {
		return uuid.Nil, NewUnauthorizedError(MsgUnauthorized)
	}

	adminIDStr, ok := adminIDVal.(string)
	if !ok {
		return uuid.Nil, NewUnauthorizedError(MsgInvalidUserID)
	}

	adminID, err := uuid.Parse(adminIDStr)
	if err != nil {
		return uuid.Nil, NewBadRequestError(MsgInvalidUserID + ": " + MsgInvalidUUIDFormat)
	}

	return adminID, nil
}
