package admin

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	redirectDepartments = "/admin/departments"
	redirectServices    = "/admin/services"
)

// parseUUID - Parse UUID từ path param, redirect với error nếu không hợp lệ
func parseUUID(ctx *gin.Context, redirectTo string) (uuid.UUID, bool) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.Redirect(http.StatusFound, redirectTo+"?error=Invalid+ID+format")
		return uuid.Nil, false
	}
	return id, true
}

// parseDepartmentID - Parse UUID cho department routes
func parseDepartmentID(ctx *gin.Context) (uuid.UUID, bool) {
	return parseUUID(ctx, redirectDepartments)
}

// parseServiceID - Parse UUID cho service routes
func parseServiceID(ctx *gin.Context) (uuid.UUID, bool) {
	return parseUUID(ctx, redirectServices)
}
