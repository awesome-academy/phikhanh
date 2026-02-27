package admin

import (
	"net/http"
	"phikhanh/utils"

	"github.com/gin-gonic/gin"
)

type DepartmentController struct{}

func NewDepartmentController() *DepartmentController {
	return &DepartmentController{}
}

func (c *DepartmentController) ShowDepartments(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "admin/departments.html", utils.GetAdminData(ctx, "Departments", "departments"))
}
