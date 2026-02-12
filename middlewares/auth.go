package middlewares

import (
	"net/http"
	"strings"

	"phikhanh/utils"

	"github.com/gin-gonic/gin"
)

// Middleware xác thực JWT token
func AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			utils.ErrorResponse(ctx, http.StatusUnauthorized, "Authorization header required")
			ctx.Abort()
			return
		}

		// Bearer token
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			utils.ErrorResponse(ctx, http.StatusUnauthorized, "Invalid authorization header format")
			ctx.Abort()
			return
		}

		token := parts[1]
		claims, err := utils.VerifyToken(token)
		if err != nil {
			utils.ErrorResponse(ctx, http.StatusUnauthorized, "Invalid or expired token")
			ctx.Abort()
			return
		}

		// Lưu thông tin user vào context
		ctx.Set("user_id", claims.UserID)
		ctx.Set("role", claims.Role)
		ctx.Next()
	}
}
