package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// RequireRole - Middleware kiểm tra role của user
func RequireRole(allowedRoles ...string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Lấy role từ context (được set bởi auth middleware)
		role, exists := ctx.Get("admin_role")
		if !exists {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		roleStr, ok := role.(string)
		if !ok {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid role format"})
			return
		}

		// Kiểm tra role có trong danh sách allowed roles
		allowed := false
		for _, ar := range allowedRoles {
			if roleStr == ar {
				allowed = true
				break
			}
		}

		if !allowed {
			ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Access denied"})
			return
		}

		ctx.Next()
	}
}
