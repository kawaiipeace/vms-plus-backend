package handlers

import (
	"net/http"
	"sort"
	"strings"
	"time"
	"vms_plus_be/config"
	"vms_plus_be/funcs"
	"vms_plus_be/messages"
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
	user := funcs.GetAuthenUser(c, "*")
	var lists []models.MasUserEmp
	search := c.Query("search")

	query := config.DBu

	if search != "" {
		query = query.Where("emp_id ILIKE ? OR full_name ILIKE ?", "%"+search+"%", "%"+search+"%")
	}
	query = query.Where("bureau_dept_sap = ?", user.BureauDeptSap)
	query = query.Limit(100)
	if err := query.
		Find(&lists).Error; err != nil {
		c.JSON(http.StatusOK, []interface{}{})
		return
	}

	for i := range lists {
		lists[i].ImageUrl = funcs.GetEmpImage(lists[i].EmpID)
	}

	c.JSON(http.StatusOK, lists)
}

// ListReceivedKeyUser godoc
// @Summary Retrieve the ReceivedKey Users
// @Description This endpoint allows a user to retrieve ReceivedKey Users.
// @Tags MAS
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Param trn_request_uid query string true "TrnReuestUID"
// @Param search query string false "Search by Employee ID or Full Name"
// @Router /api/mas/user-received-key-users [get]
func (h *MasHandler) ListReceivedKeyUser(c *gin.Context) {
	user := funcs.GetAuthenUser(c, "*")
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
		c.JSON(http.StatusNotFound, gin.H{"error": "Request not found", "message": messages.ErrNotfound.Error()})
		return
	}
	query := config.DBu
	query = query.Where("bureau_dept_sap = ?", user.BureauDeptSap)
	// Apply search filter if provided
	if search != "" {
		query = query.Where("emp_id ILIKE ? OR full_name ILIKE ?", "%"+search+"%", "%"+search+"%")
	}
	query = query.Limit(100)
	query = query.Order(gorm.Expr(
		"CASE emp_id WHEN ? THEN 1 WHEN ? THEN 2 WHEN ? THEN 3 ELSE 4 END",
		request.DriverEmpID, request.VehicleUserEmpID, request.CreatedRequestEmpID))

	query = query.Order("CASE emp_id WHEN '" + request.VehicleUserEmpID + "' THEN 1 WHEN '" + request.DriverEmpID + "' THEN 2 WHEN '" + request.CreatedRequestEmpID + "' THEN 3 ELSE 4 END")
	// Execute query
	if err := query.
		Find(&lists).Error; err != nil {
		c.JSON(http.StatusOK, []interface{}{})
		return
	}
	// Loop to modify or set AnnualDriver
	for i := range lists {
		lists[i].ImageUrl = funcs.GetEmpImage(lists[i].EmpID)
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
	user := funcs.GetAuthenUser(c, "*")
	var lists []models.MasUserDriver
	search := c.Query("search")

	query := config.DB

	// Apply search filter if provided
	if search != "" {
		query = query.Where("emp_id ILIKE ? OR full_name ILIKE ?", "%"+search+"%", "%"+search+"%")
	}
	query = query.Where("bureau_dept_sap = ?", user.BureauDeptSap)
	query = query.Limit(100)

	// Execute query
	if err := query.Preload("AnnualDriver").
		Find(&lists).Error; err != nil {
		c.JSON(http.StatusOK, []interface{}{})
		return
	}

	// Loop to modify or set AnnualDriver
	for i := range lists {
		lists[i].ImageUrl = funcs.GetEmpImage(lists[i].EmpID)
		lists[i].AnnualDriver.AnnualYYYY = 2568
		lists[i].AnnualDriver.DriverLicenseNo = "A123456"
		lists[i].AnnualDriver.RequestAnnualDriverNo = "B00001"
		lists[i].AnnualDriver.DriverLicenseExpireDate = time.Now().AddDate(1, 0, 0)
		lists[i].AnnualDriver.RequestIssueDate = time.Now().AddDate(0, -6, 0)
		lists[i].AnnualDriver.RequestExpireDate = time.Now().AddDate(0, 6, 0)
	}

	c.JSON(http.StatusOK, lists)
}

// ListConfirmerUser godoc
// @Summary Retrieve the Confirmer Users
// @Description This endpoint allows a user to retrieve Confirmer Users.
// @Tags MAS
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Param search query string false "Search by Employee ID or Full Name"
// @Router /api/mas/user-confirmer-users [get]
func (h *MasHandler) ListConfirmerUser(c *gin.Context) {
	user := funcs.GetAuthenUser(c, "*")
	var lists []models.MasUserEmp
	search := c.Query("search")

	query := config.DBu
	if search != "" {
		query = query.Where("emp_id ILIKE ? OR full_name ILIKE ?", "%"+search+"%", "%"+search+"%")
	}
	query = query.Where("bureau_dept_sap = ? AND level_code in ('M1','M2','M3')", user.BureauDeptSap)
	query = query.Limit(100)
	query = query.Order("level_code")

	// Execute query
	if err := query.
		Find(&lists).Error; err != nil {
		c.JSON(http.StatusOK, []interface{}{})
		return
	}

	// For loop to set Image_url for each element in the slice
	for i := range lists {
		lists[i].ImageUrl = funcs.GetEmpImage(lists[i].EmpID)
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
// @Param trn_request_uid query string true "TrnReuestUID"
// @Param search query string false "Search by Employee ID or Full Name"
// @Router /api/mas/user-admin-approval-users [get]
func (h *MasHandler) ListAdminApprovalUser(c *gin.Context) {
	var lists []models.MasUserEmp
	trnRequestUID := c.Query("trn_request_uid")
	if trnRequestUID == "" {
		h.ListConfirmerLicenseUser(c)
		return
	}
	search := c.Query("search")
	var result struct {
		MasCarpoolUID string
		MasVehicleUID string
	}
	if err := config.DB.Table("vms_trn_request").
		Select("mas_carpool_uid, mas_vehicle_uid").
		Where("trn_request_uid = ?", trnRequestUID).
		Scan(&result).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Request not found", "message": messages.ErrNotfound.Error()})
		return
	}
	query := config.DB

	if result.MasCarpoolUID != "" {
		query = query.Where("emp_id in (select emp_uid from vms_mas_carpool_admin ca where ca.mas_carpool_uid = ? AND ca.is_deleted='0' AND ca.is_active='1')", result.MasCarpoolUID)
	} else {
		var bureauDeptSap string
		if err := config.DB.Table("vms_mas_vehicle_department").
			Select("bureau_dept_sap").
			Where("mas_vehicle_uid = ?", result.MasVehicleUID).
			Scan(&bureauDeptSap).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Request not found", "message": messages.ErrNotfound.Error()})
			return
		}
		query = query.Where("emp_id in (select da.emp_id from vms_mas_department_admin da where da.bureau_dept_sap = ?)", bureauDeptSap)
	}

	// Apply search filter if provided
	if search != "" {
		query = query.Where("emp_id ILIKE ? OR full_name ILIKE ?", "%"+search+"%", "%"+search+"%")
	}
	query = query.Limit(100)

	// Execute query
	if err := query.
		Find(&lists).Error; err != nil {
		c.JSON(http.StatusOK, []interface{}{})
		return
	}

	// For loop to set Image_url for each element in the slice
	for i := range lists {
		lists[i].ImageUrl = funcs.GetEmpImage(lists[i].EmpID)
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
// @Param trn_request_uid query string true "TrnReuestUID"
// @Param search query string false "Search by Employee ID or Full Name"
// @Router /api/mas/user-final-approval-users [get]
func (h *MasHandler) ListFinalApprovalUser(c *gin.Context) {
	var lists []models.MasUserEmp
	trnRequestUID := c.Query("trn_request_uid")
	if trnRequestUID == "" {
		h.ListConfirmerLicenseUser(c)
		return
	}
	search := c.Query("search")
	var result struct {
		MasCarpoolUID string
		MasVehicleUID string
	}
	if err := config.DB.Table("vms_trn_request").
		Select("mas_carpool_uid, mas_vehicle_uid").
		Where("trn_request_uid = ?", trnRequestUID).
		Scan(&result).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Request not found", "message": messages.ErrNotfound.Error()})
		return
	}
	query := config.DB

	if result.MasCarpoolUID != "" {
		query = query.Where("emp_id in (select emp_uid from vms_mas_carpool_approver ca where ca.mas_carpool_uid = ? AND ca.is_deleted='0' AND ca.is_active='1')", result.MasCarpoolUID)
	} else {
		var bureauDeptSap string
		if err := config.DB.Table("vms_mas_vehicle_department").
			Select("bureau_dept_sap").
			Where("mas_vehicle_uid = ?", result.MasVehicleUID).
			Scan(&bureauDeptSap).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Request not found", "message": messages.ErrNotfound.Error()})
			return
		}
		query = query.Where("emp_id in (select da.emp_id from vms_mas_department_approver da where da.bureau_dept_sap = ?)", bureauDeptSap)
	}

	if search != "" {
		query = query.Where("emp_id ILIKE ? OR full_name ILIKE ?", "%"+search+"%", "%"+search+"%")
	}
	query = query.Limit(100)

	// Execute query
	if err := query.
		Find(&lists).Error; err != nil {
		c.JSON(http.StatusOK, []interface{}{})
		return
	}

	// For loop to set Image_url for each element in the slice
	for i := range lists {
		lists[i].ImageUrl = funcs.GetEmpImage(lists[i].EmpID)
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
	if c.IsAborted() {
		return
	}
	EmpID := c.Param("emp_id")

	var userEmp models.MasUserEmp
	if err := config.DBu.
		First(&userEmp, "emp_id = ?", EmpID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found", "message": messages.ErrNotfound.Error()})
		return
	}
	userEmp.ImageUrl = funcs.GetEmpImage(userEmp.EmpID)
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
	if c.IsAborted() {
		return
	}

	var list []models.VmsMasSatisfactionSurveyQuestions
	if err := config.DB.
		Order("question_no").
		Find(&list, "is_deleted = ? AND is_active = ?", "0", "1").Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Questions not found", "message": messages.ErrNotfound.Error()})
		return
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve vehicle departments", "message": messages.ErrInternalServer.Error()})
		return
	}

	// Second Query
	if err := config.DB.Table("vms_mas_carpool").
		Select("CAST(mas_carpool_uid AS TEXT) AS vehicle_owner_dept_sap, carpool_name AS dept_sap_short, carpool_name AS dept_sap_full, 'Car pool' AS dept_type").
		Where("is_deleted = ? AND is_active = ?", "0", "1").
		Find(&carpools).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve carpools", "message": messages.ErrInternalServer.Error()})
		return
	}
	results := append(vehicleDepts, carpools...)

	c.JSON(http.StatusOK, results)
}

// ListDepartment godoc
// @Summary Retrieve the Department Tree
// @Description This endpoint allows a user to retrieve the Department Tree.
// @Tags MAS
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Param dept_upper query string false "Filter by Department Upper"
// @Router /api/mas/department-tree [get]
func (h *MasHandler) GetDepartmentTree(c *gin.Context) {
	deptUpper := c.Query("dept_upper")

	var departments []models.VmsMasDepartmentTree
	if err := config.DB.
		Where("dept_upper = ? AND is_active = ?", deptUpper, "1").
		Find(&departments).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Departments not found"})
		return
	}

	for i := range departments {
		departments[i].DeptFull = departments[i].DeptFull + " (" + departments[i].DeptSAP + ")" + "-" + departments[i].ResourceName
		populateSubDepartments(&departments[i], 8)
	}
	c.JSON(http.StatusOK, departments)
}

func populateSubDepartments(department *models.VmsMasDepartmentTree, levels int) {
	if levels <= 0 {
		return
	}
	var subDepartments []models.VmsMasDepartmentTree
	if err := config.DB.
		Where("dept_upper = ? AND is_active = ?", department.DeptSAP, "1").
		Find(&subDepartments).Error; err == nil {
		department.DeptUnder = subDepartments
		for i := range subDepartments {
			subDepartments[i].DeptFull = subDepartments[i].DeptFull + " (" + subDepartments[i].DeptSAP + ")" + "-" + subDepartments[i].ResourceName
			populateSubDepartments(&subDepartments[i], levels-1)
		}
	}
}

// ListDriverVendor godoc
// @Summary Retrieve the Driver Vendors
// @Description This endpoint allows a user to retrieve Driver Vendors.
// @Tags MAS
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Param search query string false "Search by Vendor Code or Vendor Name"
// @Router /api/mas/driver-vendors [get]
func (h *MasHandler) ListDriverVendor(c *gin.Context) {
	search := c.Query("search")
	var vendors []models.VmsMasDriverVendor

	query := config.DB
	query = query.Where("is_deleted = ?", "0")
	// Apply search filter if provided
	if search != "" {
		query = query.Where("mas_vendor_code ILIKE ? OR mas_vendor_name ILIKE ?", "%"+search+"%", "%"+search+"%")
	}

	// Execute query
	if err := query.Order("mas_vendor_code").Find(&vendors).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve driver vendors", "message": messages.ErrInternalServer.Error()})
		return
	}

	c.JSON(http.StatusOK, vendors)
}

// ListDriverDepartment godoc
// @Summary Retrieve the Driver Departments
// @Description This endpoint allows a user to retrieve Driver Departments.
// @Tags MAS
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Param search query string false "Search by DeptSap Code or DetpSap Name"
// @Router /api/mas/driver-departments [get]
func (h *MasHandler) ListDriverDepartment(c *gin.Context) {
	var driverDepts []models.VmsMasDepartment
	search := c.Query("search")
	query := config.DB
	query = query.Where("is_deleted = ? AND is_active = ?", "0", "1")
	if search != "" {
		query = query.Where("dept_sap ILIKE ? OR dept_short ILIKE ? OR dept_full ILIKE ?", "%"+search+"%", "%"+search+"%")
	}
	if err := query.
		Order("dept_sap").
		Limit(100).
		Find(&driverDepts).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve driver departments", "message": messages.ErrInternalServer.Error()})
		return
	}

	c.JSON(http.StatusOK, driverDepts)
}

// ListConfirmerLicenseUser godoc
// @Summary Retrieve the Confirmer License Users
// @Description	This endpoint allows a user to retrieve Confirmer License Users.
// @Tags MAS
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Param search query string false "Search by Employee ID or Full Name"
// @Router /api/mas/user-confirmer-license-users [get]
func (h *MasHandler) ListConfirmerLicenseUser(c *gin.Context) {
	user := funcs.GetAuthenUser(c, "*")
	var lists []models.MasUserEmp
	search := c.Query("search")

	query := config.DBu
	if search != "" {
		query = query.Where("emp_id ILIKE ? OR full_name ILIKE ?", "%"+search+"%", "%"+search+"%")
	}
	query = query.Where("bureau_dept_sap = ? AND level_code in ('M1','M2','M3')", user.BureauDeptSap)
	query = query.Limit(100)
	query = query.Order("level_code")

	// Execute query
	if err := query.
		Find(&lists).Error; err != nil {
		c.JSON(http.StatusOK, []interface{}{})
		return
	}

	// For loop to set Image_url for each element in the slice
	for i := range lists {
		lists[i].ImageUrl = funcs.GetEmpImage(lists[i].EmpID)
	}

	c.JSON(http.StatusOK, lists)
}

// ListApprovalLicenseUser godoc
// @Summary Retrieve the Approval License Users
// @Description	This endpoint allows a user to retrieve Approval License Users.
// @Tags MAS
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Param search query string false "Search by Employee ID or Full Name"
// @Router /api/mas/user-approval-license-users [get]
func (h *MasHandler) ListApprovalLicenseUser(c *gin.Context) {
	user := funcs.GetAuthenUser(c, "*")
	var lists []models.MasUserEmp
	search := c.Query("search")

	query := config.DBu
	if search != "" {
		query = query.Where("emp_id ILIKE ? OR full_name ILIKE ?", "%"+search+"%", "%"+search+"%")
	}
	query = query.Where("bureau_dept_sap = ? AND level_code in ('M4','S1')", user.BureauDeptSap)
	query = query.Limit(100)
	query = query.Order("level_code")

	// Execute query
	if err := query.
		Find(&lists).Error; err != nil {
		c.JSON(http.StatusOK, []interface{}{})
		return
	}

	// For loop to set Image_url for each element in the slice
	for i := range lists {
		lists[i].ImageUrl = funcs.GetEmpImage(lists[i].EmpID)
	}

	c.JSON(http.StatusOK, lists)
}

// ListHoliday godoc
// @Summary Retrieve the Holidays
// @Description This endpoint allows a user to retrieve Holidays.
// @Tags MAS
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Param start_date query string false "Start date (YYYY-MM-DD)"
// @Param end_date query string false "End date (YYYY-MM-DD)"
// @Router /api/mas/holidays [get]
func (h *MasHandler) ListHoliday(c *gin.Context) {
	var holidays []models.VmsMasHolidays
	query := config.DB
	query = query.Where("is_deleted = ?", "0")
	query = query.Order("mas_holidays_date")
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")
	if startDate != "" {
		query = query.Where("mas_holidays_date >= ?", startDate)
	}
	if endDate != "" {
		query = query.Where("mas_holidays_date <= ?", endDate)
	}
	if err := query.Find(&holidays).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve holidays", "message": messages.ErrInternalServer.Error()})
		return
	}
	// Add weekends (Saturday and Sunday) to holidays
	if startDate != "" && endDate != "" {
		start, err := time.Parse("2006-01-02", startDate)
		if err == nil {
			end, err := time.Parse("2006-01-02", endDate)
			if err == nil {
				// Iterate through dates and add weekends
				for d := start; !d.After(end); d = d.AddDate(0, 0, 1) {
					if d.Weekday() == time.Saturday || d.Weekday() == time.Sunday {
						holidays = append(holidays, models.VmsMasHolidays{
							HolidaysDate: d,
							HolidaysDetail: map[time.Weekday]string{
								time.Saturday: "วันเสาร์",
								time.Sunday:   "วันอาทิตย์",
							}[d.Weekday()],
						})
					}
				}
			}
		}
	}
	// Sort holidays by date
	sort.Slice(holidays, func(i, j int) bool {
		return holidays[i].HolidaysDate.Before(holidays[j].HolidaysDate)
	})
	c.JSON(http.StatusOK, holidays)
}
