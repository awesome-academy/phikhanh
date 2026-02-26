package middlewares

import (
	"net/http"
	"strings"

	"phikhanh/utils"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware - Middleware xác thực JWT token
func AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			utils.ErrorResponse(ctx, http.StatusUnauthorized, "Authorization header required")
			ctx.Abort()
			return
		}

		var token string

		if strings.HasPrefix(authHeader, "Bearer ") {
			token = strings.TrimPrefix(authHeader, "Bearer ")
		} else {
			token = authHeader
		}

		token = strings.TrimSpace(token)
		if token == "" {
			utils.ErrorResponse(ctx, http.StatusUnauthorized, "Token is required")
			ctx.Abort()
			return
		}

		claims, err := utils.VerifyToken(token)
		if err != nil {
			utils.ErrorResponse(ctx, http.StatusUnauthorized, "Invalid or expired token")
			ctx.Abort()
			return
		}

		ctx.Set("user_id", claims.UserID)
		ctx.Set("role", claims.Role)
		ctx.Next()
	}
}
