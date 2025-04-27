package handlers

import (
	"net/http"
	"strings"
	"time"
	"vms_plus_be/config"
	"vms_plus_be/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type MasHandler struct {
}

// ListVehicleUser godoc
// @Summary Retrieve the Vehicle Users
// @Description This endpoint allows a user to retrieve Vehicle Users.
// @Tags MAS
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Param search query string false "Search by Employee ID or Full Name"
// @Router /api/mas/user-vehicle-users [get]
func (h *MasHandler) ListVehicleUser(c *gin.Context) {
	var lists []models.MasUserDriver
	search := c.Query("search")

	query := config.DB

	// Apply search filter if provided
	if search != "" {
		query = query.Where("emp_id = ? OR full_name ILIKE ?", search, "%"+search+"%")
	}

	// Execute query
	if err := query.
		Find(&lists).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Not found"})
		return
	}

	// Loop to modify or set AnnualDriver
	for i := range lists {
		lists[i].ImageURL = config.DefaultAvatarURL
	}

	c.JSON(http.StatusOK, lists)
}

// ListVehicleUser godoc
// @Summary Retrieve the ReceivedKey Users
// @Description This endpoint allows a user to retrieve ReceivedKey Users.
// @Tags MAS
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Param search query string false "Search by Employee ID or Full Name"
// @Param trn_request_uid query string false "TrnReuestUID"
// @Router /api/mas/user-received-key-users [get]
func (h *MasHandler) ListReceivedKeyUser(c *gin.Context) {
	var lists []models.MasUserDriver
	search := c.Query("search")
	trnRequestUID := c.Query("trn_request_uid")
	var request struct {
		DriverEmpID         string `gorm:"column:driver_emp_id" json:"driver_emp_id" example:"700001"`
		VehicleUserEmpID    string `gorm:"column:vehicle_user_emp_id" json:"vehicle_user_emp_id" example:"990001"`
		CreatedRequestEmpID string `gorm:"column:created_request_emp_id" json:"created_request_emp_id" example:"700001"`
	}
	if err := config.DB.Table("vms_trn_request").
		First(&request, "trn_request_uid = ?", trnRequestUID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Request not found"})
		return
	}
	query := config.DB

	// Apply search filter if provided
	if search != "" {
		query = query.Where("emp_id = ? OR full_name ILIKE ?", search, "%"+search+"%")
	}
	query = query.Order(gorm.Expr(
		"CASE emp_id WHEN ? THEN 1 WHEN ? THEN 2 WHEN ? THEN 3 ELSE 4 END",
		request.DriverEmpID, request.VehicleUserEmpID, request.CreatedRequestEmpID))

	query = query.Order("CASE emp_id WHEN '" + request.VehicleUserEmpID + "' THEN 1 WHEN '" + request.DriverEmpID + "' THEN 2 WHEN '" + request.CreatedRequestEmpID + "' THEN 3 ELSE 4 END")
	// Execute query
	if err := query.
		Find(&lists).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Not found"})
		return
	}
	// Loop to modify or set AnnualDriver
	for i := range lists {
		lists[i].ImageURL = config.DefaultAvatarURL
		var roles []string
		if lists[i].EmpID == request.DriverEmpID {
			roles = append(roles, "ผู้ขับขี่")
		}
		if lists[i].EmpID == request.VehicleUserEmpID {
			roles = append(roles, "ผู้ใช้ยานพาหนะ")
		}
		if lists[i].EmpID == request.CreatedRequestEmpID {
			roles = append(roles, "ผู้ใช้สร้างคำขอ")
		}
		if len(roles) > 0 {
			lists[i].FullName = lists[i].FullName + " (" + strings.Join(roles, ", ") + ")"
		}
	}

	c.JSON(http.StatusOK, lists)
}

// ListDriverUser godoc
// @Summary Retrieve the Driver Users
// @Description This endpoint allows a user to retrieve Driver Users.
// @Tags MAS
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Param search query string false "Search by Employee ID or Full Name"
// @Router /api/mas/user-driver-users [get]
func (h *MasHandler) ListDriverUser(c *gin.Context) {
	var lists []models.MasUserDriver
	search := c.Query("search")

	query := config.DB

	// Apply search filter if provided
	if search != "" {
		query = query.Where("emp_id = ? OR full_name ILIKE ?", search, "%"+search+"%")
	}

	// Execute query
	if err := query.Preload("AnnualDriver").
		Find(&lists).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Not found"})
		return
	}

	// Loop to modify or set AnnualDriver
	for i := range lists {
		lists[i].ImageURL = config.DefaultAvatarURL
		lists[i].AnnualDriver.AnnualYYYY = 2568
		lists[i].AnnualDriver.DriverLicenseNo = "A123456"
		lists[i].AnnualDriver.RequestAnnualDriverNo = "B00001"
		lists[i].AnnualDriver.DriverLicenseExpireDate = time.Now().AddDate(1, 0, 0)
		lists[i].AnnualDriver.RequestIssueDate = time.Now().AddDate(0, -6, 0)
		lists[i].AnnualDriver.RequestExpireDate = time.Now().AddDate(0, 6, 0)
	}

	c.JSON(http.StatusOK, lists)
}

// ListApprovalUser godoc
// @Summary Retrieve the Approval Users
// @Description This endpoint allows a user to retrieve Approval Users.
// @Tags MAS
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Param search query string false "Search by Employee ID or Full Name"
// @Router /api/mas/user-approval-users [get]
func (h *MasHandler) ListApprovalUser(c *gin.Context) {
	var lists []models.MasUserEmp
	search := c.Query("search")

	query := config.DB

	// Apply search filter if provided
	if search != "" {
		query = query.Where("emp_id = ? OR full_name ILIKE ?", search, "%"+search+"%")
	}

	// Execute query
	if err := query.
		Find(&lists).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Not found"})
		return
	}

	defaultURL := "http://pntdev.ddns.net:28089/VMS_PLUS/PIX/user-avatar.jpg"
	// For loop to set Image_url for each element in the slice
	for i := range lists {
		lists[i].Image_url = defaultURL
	}

	c.JSON(http.StatusOK, lists)
}

// ListAdminApprovalUser godoc
// @Summary Retrieve the Admin Approval Users
// @Description This endpoint allows a user to retrieve Admin Approval Users.
// @Tags MAS
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Param search query string false "Search by Employee ID or Full Name"
// @Router /api/mas/user-admin-approval-users [get]
func (h *MasHandler) ListAdminApprovalUser(c *gin.Context) {
	var lists []models.MasUserEmp
	search := c.Query("search")

	query := config.DB

	// Apply search filter if provided
	if search != "" {
		query = query.Where("emp_id = ? OR full_name ILIKE ?", search, "%"+search+"%")
	}

	// Execute query
	if err := query.
		Find(&lists).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Not found"})
		return
	}

	defaultURL := "http://pntdev.ddns.net:28089/VMS_PLUS/PIX/user-avatar.jpg"
	// For loop to set Image_url for each element in the slice
	for i := range lists {
		lists[i].Image_url = defaultURL
	}

	c.JSON(http.StatusOK, lists)
}

// ListFinalApprovalUser godoc
// @Summary Retrieve the Final Approval User
// @Description This endpoint allows a user to retrieve Final Approval User.
// @Tags MAS
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Param search query string false "Search by Employee ID or Full Name"
// @Router /api/mas/user-final-approval-users [get]
func (h *MasHandler) ListFinalApprovalUser(c *gin.Context) {
	var lists []models.MasUserEmp
	search := c.Query("search")

	query := config.DB

	// Apply search filter if provided
	if search != "" {
		query = query.Where("emp_id = ? OR full_name ILIKE ?", search, "%"+search+"%")
	}

	// Execute query
	if err := query.
		Find(&lists).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Not found"})
		return
	}

	defaultURL := "http://pntdev.ddns.net:28089/VMS_PLUS/PIX/user-avatar.jpg"
	// For loop to set Image_url for each element in the slice
	for i := range lists {
		lists[i].Image_url = defaultURL
	}

	c.JSON(http.StatusOK, lists)
}

// GetUserEmp godoc
// @Summary Retrieve a specific user
// @Description This endpoint fetches details of a user.
// @Tags MAS
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Param emp_id path string true "EmpID (emp_id)"
// @Router /api/mas/user/{emp_id} [get]
func (h *MasHandler) GetUserEmp(c *gin.Context) {
	//funcs.GetAuthenUser(c, h.Role)
	EmpID := c.Param("emp_id")

	var userEmp models.MasUserEmp
	if err := config.DB.
		First(&userEmp, "emp_id = ?", EmpID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	c.JSON(http.StatusOK, userEmp)
}

// GetVmsMasSatisfactionSurveyQuestions godoc
// @Summary Retrieve a specific user
// @Description This endpoint fetches details of a user.
// @Tags MAS
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Router /api/mas/satisfaction_survey_questions [get]
func (h *MasHandler) ListVmsMasSatisfactionSurveyQuestions(c *gin.Context) {
	//funcs.GetAuthenUser(c, h.Role)

	var list []models.VmsMasSatisfactionSurveyQuestions
	if err := config.DB.
		Order("ordering").
		Find(&list, "is_deleted = ? AND is_active = ?", "0", "1").Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Questions not found"})
		return
	}
	// Process the splitting
	for i := range list {
		desc := list[i].MasSatisfactionSurveyQuestionsDesc
		parts := strings.SplitN(desc, ":", 2)
		list[i].MasSatisfactionSurveyQuestionsTitle = parts[0] // Title before colon
		if len(parts) > 1 {
			list[i].MasSatisfactionSurveyQuestionsDesc = parts[1] // Remaining description after colon
		} else {
			list[i].MasSatisfactionSurveyQuestionsDesc = "" // Empty if no colon found
		}
	}

	c.JSON(http.StatusOK, list)
}

// ListVehicleDepartment godoc
// @Summary Retrieve the Vehicle Departments
// @Description This endpoint allows a user to retrieve Vehicle Departments.
// @Tags MAS
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Router /api/mas/vehicle-departments [get]
func (h *MasHandler) ListVehicleDepartment(c *gin.Context) {

	var vehicleDepts []models.VmsMasVehicleDepartmentList
	var carpools []models.VmsMasVehicleDepartmentList

	// First Query
	if err := config.DB.Table("vms_mas_vehicle_department AS vd").
		Select("vd.vehicle_owner_dept_sap, MAX(d.dept_short) AS dept_sap_short, MAX(d.dept_full) AS dept_sap_full, 'PEA' AS dept_type").
		Joins("INNER JOIN vms_mas_department d ON d.dept_sap = vd.vehicle_owner_dept_sap").
		Where("vd.is_deleted = ? AND vd.is_active = ? AND d.is_deleted = ?", "0", "1", "0").
		Group("vd.vehicle_owner_dept_sap").
		Order("vd.vehicle_owner_dept_sap").
		Find(&vehicleDepts).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve vehicle departments"})
		return
	}

	// Second Query
	if err := config.DB.Table("vms_mas_carpool").
		Select("CAST(mas_carpool_uid AS TEXT) AS vehicle_owner_dept_sap, carpool_name AS dept_sap_short, carpool_name AS dept_sap_full, 'Car pool' AS dept_type").
		Where("is_deleted = ? AND is_active = ?", "0", "1").
		Find(&carpools).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve carpools"})
		return
	}
	results := append(vehicleDepts, carpools...)

	c.JSON(http.StatusOK, results)
}
