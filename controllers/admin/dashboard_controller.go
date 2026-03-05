package admin

import (
	"net/http"

	adminSvc "phikhanh/services/admin"
	"phikhanh/utils"

	"github.com/gin-gonic/gin"
)

type DashboardController struct {
	appService  *adminSvc.ApplicationAdminService
	svcService  *adminSvc.ServiceAdminService
	deptService *adminSvc.DepartmentService
	userRepo    interface{} // User repository
}

func NewDashboardController(
	appService *adminSvc.ApplicationAdminService,
	svcService *adminSvc.ServiceAdminService,
	deptService *adminSvc.DepartmentService,
) *DashboardController {
	return &DashboardController{
		appService:  appService,
		svcService:  svcService,
		deptService: deptService,
	}
}

func (c *DashboardController) ShowDashboard(ctx *gin.Context) {
	// Lấy dữ liệu từ các service
	// Applications stats
	appList, _ := c.appService.GetList("", nil, 1)
	var applicationCount int64
	var processingCount int64
	if appList != nil {
		applicationCount = appList.TotalItems
		// Count processing status
		for _, item := range appList.Items {
			if item.Status == "Processing" {
				processingCount++
			}
		}
	}

	// Services count
	services, _ := c.svcService.GetAll()
	serviceCount := len(services)

	// Departments count
	departments, _ := c.deptService.GetAll()
	departmentCount := len(departments)

	data := utils.GetAdminData(ctx, "Dashboard", "dashboard")
	data["ApplicationCount"] = applicationCount
	data["ProcessingCount"] = processingCount
	data["ServiceCount"] = serviceCount
	data["DepartmentCount"] = departmentCount
	data["CsrfToken"] = getCsrfToken(ctx)

	utils.RenderHTML(ctx, http.StatusOK, "admin/dashboard.html", data)
}
