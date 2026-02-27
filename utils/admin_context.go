package utils

import "github.com/gin-gonic/gin"

// GetAdminData - Lấy thông tin admin từ context/cookie cho SSR pages
func GetAdminData(ctx *gin.Context, pageTitle, activeMenu string) gin.H {
	adminName, _ := ctx.Cookie("admin_name")
	if adminName == "" {
		adminName = "Admin"
	}

	roleStr := ""
	if role, ok := ctx.Get("admin_role"); ok {
		if r, isStr := role.(string); isStr {
			roleStr = r
		}
	}

	return gin.H{
		"PageTitle":  pageTitle,
		"ActiveMenu": activeMenu,
		"AdminName":  adminName,
		"AdminRole":  roleStr,
	}
}
