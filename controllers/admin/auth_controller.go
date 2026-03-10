package admin

import (
	"net/http"

	adminDto "phikhanh/dto/admin"
	"phikhanh/models"
	adminSvc "phikhanh/services/admin"
	"phikhanh/utils"

	"github.com/gin-gonic/gin"
)

const (
	cookieMaxAge = 3600 // 1 hour
	cookiePath   = "/admin"
)

type AuthController struct {
	service        *adminSvc.AuthService
	activityLogSvc *adminSvc.ActivityLogService
}

func NewAuthController(service *adminSvc.AuthService, activityLogSvc *adminSvc.ActivityLogService) *AuthController {
	return &AuthController{service: service, activityLogSvc: activityLogSvc}
}

// ShowLogin - Hiển thị trang login admin
func (c *AuthController) ShowLogin(ctx *gin.Context) {
	if token, err := ctx.Cookie("admin_token"); err == nil && token != "" {
		if _, err := utils.VerifyToken(token); err == nil {
			ctx.Redirect(http.StatusFound, "/admin/dashboard")
			return
		}
	}

	utils.RenderHTML(ctx, http.StatusOK, "admin/auth/login.html", gin.H{
		"Error": ctx.Query("error"),
	})
}

// ProcessLogin - Xử lý login admin
func (c *AuthController) ProcessLogin(ctx *gin.Context) {
	var req adminDto.LoginRequest

	if err := ctx.ShouldBind(&req); err != nil {
		redirectWithError(ctx, formatErrorMessage(err))
		return
	}

	user, token, err := c.service.Login(req.Email, req.Password)
	if err != nil {
		redirectWithError(ctx, formatErrorMessage(err))
		return
	}

	ctx.SetCookie("admin_token", token, cookieMaxAge, cookiePath, "", false, true)
	ctx.SetCookie("admin_name", user.Name, cookieMaxAge, cookiePath, "", false, true)
	ctx.SetCookie("admin_role", string(user.Role), cookieMaxAge, cookiePath, "", false, true)

	// Record LOGIN activity
	c.activityLogSvc.RecordActivity(
		user.ID.String(),
		models.ActionLogin,
		user.ID.String(),
		"Admin logged in: "+user.Email,
		ctx.ClientIP(),
	)

	ctx.Redirect(http.StatusFound, "/admin/dashboard")
}

// ProcessLogout - Xử lý logout admin
func (c *AuthController) ProcessLogout(ctx *gin.Context) {
	// Record LOGOUT before clearing cookies
	if adminID, err := utils.ExtractAdminID(ctx); err == nil {
		c.activityLogSvc.RecordActivity(
			adminID.String(),
			models.ActionLogout,
			adminID.String(),
			"Admin logged out",
			ctx.ClientIP(),
		)
	}

	ctx.SetCookie("admin_token", "", -1, cookiePath, "", false, true)
	ctx.SetCookie("admin_name", "", -1, cookiePath, "", false, true)
	ctx.SetCookie("admin_role", "", -1, cookiePath, "", false, true)
	ctx.Redirect(http.StatusFound, "/admin/login")
}

// redirectWithError - Redirect về login với error message
func redirectWithError(ctx *gin.Context, message string) {
	setFlashError(ctx, message, "/admin/login")
}
