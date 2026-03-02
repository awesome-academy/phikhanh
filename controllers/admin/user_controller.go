package admin

import (
	"net/http"
	"phikhanh/utils"

	"github.com/gin-gonic/gin"
)

type UserController struct{}

func NewUserController() *UserController {
	return &UserController{}
}

func (c *UserController) ShowUsers(ctx *gin.Context) {
	utils.RenderHTML(ctx, http.StatusOK, "admin/users.html", utils.GetAdminData(ctx, "Users", "users"))
}
