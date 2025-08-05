package handlers

import (
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"
	"vms_plus_be/config"
	"vms_plus_be/funcs"
	"vms_plus_be/messages"
	"vms_plus_be/models"
	"vms_plus_be/userhub"

	"github.com/gin-gonic/gin"
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
	search := c.Query("search")

	request := userhub.ServiceListUserRequest{
		ServiceCode:   "vms",
		Search:        search,
		BureauDeptSap: user.BureauDeptSap,
		Limit:         100,
	}
	lists, err := userhub.GetUserList(request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
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
	search := c.Query("search")
	trnRequestUID := c.Query("trn_request_uid")
	var trnRequest struct {
		DriverEmpID         string `gorm:"column:driver_emp_id" json:"driver_emp_id" example:"700001"`
		VehicleUserEmpID    string `gorm:"column:vehicle_user_emp_id" json:"vehicle_user_emp_id" example:"990001"`
		CreatedRequestEmpID string `gorm:"column:created_request_emp_id" json:"created_request_emp_id" example:"700001"`
	}
	if err := config.DB.Table("vms_trn_request").
		First(&trnRequest, "trn_request_uid = ?", trnRequestUID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Request not found", "message": messages.ErrNotfound.Error()})
		return
	}

	request := userhub.ServiceListUserRequest{
		ServiceCode:   "vms",
		Search:        search,
		BureauDeptSap: user.BureauDeptSap,
		Limit:         100,
	}

	lists, err := userhub.GetUserList(request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// Sort lists based on role priority
	sort.Slice(lists, func(i, j int) bool {
		getPriority := func(empID string) int {
			switch empID {
			case trnRequest.VehicleUserEmpID:
				return 1
			case trnRequest.DriverEmpID:
				return 2
			case trnRequest.CreatedRequestEmpID:
				return 3
			default:
				return 4
			}
		}
		return getPriority(lists[i].EmpID) < getPriority(lists[j].EmpID)
	})
	// Loop to modify or set AnnualDriver
	for i := range lists {
		var roles []string
		if lists[i].EmpID == trnRequest.DriverEmpID {
			roles = append(roles, "ผู้ขับขี่")
		}
		if lists[i].EmpID == trnRequest.VehicleUserEmpID {
			roles = append(roles, "ผู้ใช้ยานพาหนะ")
		}
		if lists[i].EmpID == trnRequest.CreatedRequestEmpID {
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

	request := userhub.ServiceListUserRequest{
		ServiceCode:   "vms",
		Search:        search,
		BureauDeptSap: user.BureauDeptSap,
		Limit:         100,
	}
	users, err := userhub.GetUserList(request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// Convert users to MasUserDriver format and add to lists

	annual_yyyy := time.Now().Year() + 543
	empIDs := make([]string, len(users))
	for i, user := range users {
		empIDs[i] = user.EmpID
	}

	var annualDrivers []models.VmsTrnAnnualDriver
	config.DB.Where("created_request_emp_id IN (?) AND is_deleted = ? AND annual_yyyy = ? AND ref_request_annual_driver_status_code = ?",
		empIDs, "0", annual_yyyy, "30").
		Select("created_request_emp_id", "annual_yyyy", "driver_license_no", "request_annual_driver_no", "driver_license_expire_date", "request_issue_date", "request_expire_date").
		Find(&annualDrivers)

	annualDriverMap := make(map[string]models.VmsTrnAnnualDriver)
	for _, ad := range annualDrivers {
		annualDriverMap[ad.CreatedRequestEmpId] = ad
	}

	lists = make([]models.MasUserDriver, len(users))
	for i, user := range users {
		lists[i] = models.MasUserDriver{
			EmpID:        user.EmpID,
			FullName:     user.FullName,
			DeptSAP:      user.DeptSAP,
			DeptSAPShort: user.DeptSAPShort,
			DeptSAPFull:  user.DeptSAPFull,
			TelMobile:    user.TelMobile,
			TelInternal:  user.TelInternal,
			ImageUrl:     user.ImageUrl,
			AnnualDriver: annualDriverMap[user.EmpID],
			Position:     user.Position,
		}
	}
	// Sort lists to put the current user's emp_id first
	sort.SliceStable(lists, func(i, j int) bool {
		if lists[i].EmpID == user.EmpID {
			return true
		}
		if lists[j].EmpID == user.EmpID {
			return false
		}
		return false
	})
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
// @Param emp_id query string true "EmpID"
// @Router /api/mas/user-confirmer-users [get]
func (h *MasHandler) ListConfirmerUser(c *gin.Context) {
	user := funcs.GetAuthenUser(c, "*")
	empID := c.Query("emp_id")
	if empID == "" {
		empID = user.EmpID
	}
	userInfo := funcs.GetUserEmpInfo(empID)
	managers := funcs.GetUserManager(userInfo.DeptSAP)
	list := []models.MasUserEmp{}
	listA := []models.MasUserEmp{}
	for _, manager := range managers {
		if manager.Type == "L" && manager.LevelCode >= "M5" {
			list = append(list, models.MasUserEmp{
				EmpID:        strconv.Itoa(manager.EmpIDLeader),
				FullName:     manager.EmpName,
				Position:     manager.PlansTextShort,
				DeptSAP:      strconv.Itoa(manager.DeptSAP),
				DeptSAPShort: funcs.GetDeptSAPShort(strconv.Itoa(manager.DeptSAP)),
				DeptSAPFull:  funcs.GetDeptSAPFull(strconv.Itoa(manager.DeptSAP)),
				ImageUrl:     funcs.GetEmpImage(strconv.Itoa(manager.EmpIDLeader)),
				IsEmployee:   true,
			})
		}
		if manager.Type == "L" && manager.LevelCode >= "M3" && manager.LevelCode < "M5" {
			listA = append(listA, models.MasUserEmp{
				EmpID:        strconv.Itoa(manager.EmpIDLeader),
				FullName:     manager.EmpName,
				Position:     manager.PlansTextShort,
				DeptSAP:      strconv.Itoa(manager.DeptSAP),
				DeptSAPShort: funcs.GetDeptSAPShort(strconv.Itoa(manager.DeptSAP)),
				DeptSAPFull:  funcs.GetDeptSAPFull(strconv.Itoa(manager.DeptSAP)),
				ImageUrl:     funcs.GetEmpImage(strconv.Itoa(manager.EmpIDLeader)),
				IsEmployee:   true,
			})
		}
		if manager.Type == "A" && manager.LevelCode >= "M3" {
			listA = append(listA, models.MasUserEmp{
				EmpID:        strconv.Itoa(manager.EmpIDLeader),
				FullName:     manager.EmpName,
				Position:     manager.PlansTextShort,
				DeptSAP:      strconv.Itoa(manager.DeptSAP),
				DeptSAPShort: funcs.GetDeptSAPShort(strconv.Itoa(manager.DeptSAP)),
				DeptSAPFull:  funcs.GetDeptSAPFull(strconv.Itoa(manager.DeptSAP)),
				ImageUrl:     funcs.GetEmpImage(strconv.Itoa(manager.EmpIDLeader)),
				IsEmployee:   true,
			})
		}
	}

	if len(list) == 0 {
		for _, manager := range managers {
			if manager.Type == "U" {
				list = append(list, models.MasUserEmp{
					EmpID:        strconv.Itoa(manager.EmpIDLeader),
					FullName:     manager.EmpName,
					Position:     manager.PlansTextShort,
					DeptSAP:      strconv.Itoa(manager.DeptSAP),
					DeptSAPShort: funcs.GetDeptSAPShort(strconv.Itoa(manager.DeptSAP)),
					DeptSAPFull:  funcs.GetDeptSAPFull(strconv.Itoa(manager.DeptSAP)),
					ImageUrl:     funcs.GetEmpImage(strconv.Itoa(manager.EmpIDLeader)),
					IsEmployee:   true,
				})
			}
		}
	}

	if len(list) == 0 && len(managers) > 0 {
		upperManagers := funcs.GetUserManager(strconv.Itoa(managers[0].DeptUpper))
		for _, manager := range upperManagers {
			if manager.Type == "L" && manager.LevelCode >= "M5" {
				list = append(list, models.MasUserEmp{
					EmpID:        strconv.Itoa(manager.EmpIDLeader),
					FullName:     manager.EmpName,
					Position:     manager.PlansTextShort,
					DeptSAP:      strconv.Itoa(manager.DeptSAP),
					DeptSAPShort: funcs.GetDeptSAPShort(strconv.Itoa(manager.DeptSAP)),
					DeptSAPFull:  funcs.GetDeptSAPFull(strconv.Itoa(manager.DeptSAP)),
					ImageUrl:     funcs.GetEmpImage(strconv.Itoa(manager.EmpIDLeader)),
					IsEmployee:   true,
				})
			}
		}

		if len(list) == 0 {
			for _, manager := range upperManagers {
				if manager.Type == "U" {
					list = append(list, models.MasUserEmp{
						EmpID:        strconv.Itoa(manager.EmpIDLeader),
						FullName:     manager.EmpName,
						Position:     manager.PlansTextShort,
						DeptSAP:      strconv.Itoa(manager.DeptSAP),
						DeptSAPShort: funcs.GetDeptSAPShort(strconv.Itoa(manager.DeptSAP)),
						DeptSAPFull:  funcs.GetDeptSAPFull(strconv.Itoa(manager.DeptSAP)),
						ImageUrl:     funcs.GetEmpImage(strconv.Itoa(manager.EmpIDLeader)),
						IsEmployee:   true,
					})
				}
			}
		}
	}
	list = append(list, listA...)
	for i := range list {
		empInfo := funcs.GetUserEmpInfo(list[i].EmpID)
		list[i].TelMobile = empInfo.TelMobile
		list[i].TelInternal = empInfo.TelInternal
	}

	//if empID=list[i].EmpID move to list[0]
	for i := range list {
		if list[i].EmpID == empID {
			list[0], list[i] = list[i], list[0]
			break
		}
	}

	c.JSON(http.StatusOK, list)
}

// ListAdminDepartmentOrCarpoolUser godoc
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
func (h *MasHandler) ListAdminDepartmentOrCarpoolUser(c *gin.Context) {
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

	var empIDs []string
	if result.MasCarpoolUID != "" && result.MasCarpoolUID != funcs.DefaultUUID() {
		if err := config.DB.Table("vms_mas_carpool_admin").
			Select("admin_emp_no").
			Where("mas_carpool_uid = ? AND is_deleted = '0' AND is_active = '1'", result.MasCarpoolUID).
			Pluck("admin_emp_no", &empIDs).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch carpool admins", "message": err.Error()})
			return
		}
		request := userhub.ServiceListUserRequest{
			ServiceCode: "vms",
			Search:      search,
			EmpIDs:      strings.Join(empIDs, ","),
			Limit:       100,
		}
		lists, err := userhub.GetUserList(request)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, lists)
	} else {
		var bureauDeptSap string
		request := userhub.ServiceListUserRequest{
			ServiceCode:   "vms",
			Search:        search,
			Role:          "admin-department",
			BureauDeptSap: bureauDeptSap,
			Limit:         100,
		}
		lists, err := userhub.GetUserList(request)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, lists)
	}
	c.JSON(http.StatusOK, []interface{}{})

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
	trnRequestUID := c.Query("trn_request_uid")
	if trnRequestUID == "" {
		h.ListConfirmerLicenseUser(c)
		return
	}
	search := c.Query("search")
	var result struct {
		MasCarpoolUID    string
		MasVehicleUID    string
		VehicleUserEmpID string
	}
	if err := config.DB.Table("vms_trn_request").
		Select("mas_carpool_uid, mas_vehicle_uid,vehicle_user_emp_id").
		Where("trn_request_uid = ?", trnRequestUID).
		Scan(&result).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Request not found", "message": messages.ErrNotfound.Error()})
		return
	}

	var empIDs []string
	if result.MasCarpoolUID != "" && result.MasCarpoolUID != funcs.DefaultUUID() {
		if err := config.DB.Table("vms_mas_carpool_approver").
			Select("approver_emp_no").
			Where("mas_carpool_uid = ? AND is_deleted = '0' AND is_active = '1'", result.MasCarpoolUID).
			Pluck("approver_emp_no", &empIDs).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch carpool admins", "message": err.Error()})
			return
		}

		request := userhub.ServiceListUserRequest{
			ServiceCode: "vms",
			Search:      search,
			EmpIDs:      strings.Join(empIDs, ","),
			Limit:       100,
		}
		lists, err := userhub.GetUserList(request)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, lists)
		return
	} else {
		userInfo := funcs.GetUserEmpInfo(result.VehicleUserEmpID)
		managers := funcs.GetUserManager(userInfo.DeptSAP)
		list := []models.MasUserEmp{}
		for _, manager := range managers {
			if manager.LevelCode == "S1" {
				list = append(list, models.MasUserEmp{
					EmpID:        strconv.Itoa(manager.EmpIDLeader),
					FullName:     manager.EmpName,
					Position:     manager.PlansTextShort,
					DeptSAP:      strconv.Itoa(manager.DeptSAP),
					DeptSAPShort: funcs.GetDeptSAPShort(strconv.Itoa(manager.DeptSAP)),
					DeptSAPFull:  funcs.GetDeptSAPFull(strconv.Itoa(manager.DeptSAP)),
					ImageUrl:     funcs.GetEmpImage(strconv.Itoa(manager.EmpIDLeader)),
					IsEmployee:   true,
				})
			}
		}

		for _, manager := range managers {
			if manager.LevelCode == "M6" {
				list = append(list, models.MasUserEmp{
					EmpID:        strconv.Itoa(manager.EmpIDLeader),
					FullName:     manager.EmpName,
					Position:     manager.PlansTextShort,
					DeptSAP:      strconv.Itoa(manager.DeptSAP),
					DeptSAPShort: funcs.GetDeptSAPShort(strconv.Itoa(manager.DeptSAP)),
					DeptSAPFull:  funcs.GetDeptSAPFull(strconv.Itoa(manager.DeptSAP)),
					ImageUrl:     funcs.GetEmpImage(strconv.Itoa(manager.EmpIDLeader)),
					IsEmployee:   true,
				})
			}
		}

		if len(list) == 0 && len(managers) > 0 {
			upperManagers := funcs.GetUserManager(strconv.Itoa(managers[0].DeptUpper))
			for _, manager := range upperManagers {
				if manager.LevelCode == "S1" || manager.LevelCode == "M6" {
					list = append(list, models.MasUserEmp{
						EmpID:        strconv.Itoa(manager.EmpIDLeader),
						FullName:     manager.EmpName,
						Position:     manager.PlansTextShort,
						DeptSAP:      strconv.Itoa(manager.DeptSAP),
						DeptSAPShort: funcs.GetDeptSAPShort(strconv.Itoa(manager.DeptSAP)),
						DeptSAPFull:  funcs.GetDeptSAPFull(strconv.Itoa(manager.DeptSAP)),
						ImageUrl:     funcs.GetEmpImage(strconv.Itoa(manager.EmpIDLeader)),
						IsEmployee:   true,
					})
				}
			}
		}
		for i := range list {
			empInfo := funcs.GetUserEmpInfo(list[i].EmpID)
			list[i].TelMobile = empInfo.TelMobile
			list[i].TelInternal = empInfo.TelInternal
		}
		c.JSON(http.StatusOK, list)
		return
	}
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

	userEmp, err := userhub.GetUserInfo(EmpID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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
		Select("vd.vehicle_owner_dept_sap, MAX(d.dept_long_short) AS dept_sap_short, MAX(d.dept_full) AS dept_sap_full, 'PEA' AS dept_type").
		Joins("INNER JOIN vms_mas_department d ON d.dept_sap = vd.vehicle_owner_dept_sap").
		Where("vd.is_deleted = ? AND vd.is_active = ? AND d.is_deleted = ?", "0", "1", "0").
		Where("d.dept_long_short IS NOT NULL").
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
	query = query.Where("dept_long_short IS NOT NULL")
	if search != "" {
		query = query.Where("dept_sap ILIKE ? OR dept_long_short ILIKE ? OR dept_full ILIKE ?", "%"+search+"%", "%"+search+"%", "%"+search+"%")
	}
	if err := query.
		Order("dept_sap").
		Limit(100).
		Find(&driverDepts).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve driver departments", "message": messages.ErrInternalServer.Error()})
		return
	}
	var carpools []models.VmsMasDepartment
	if err := config.DB.Table("vms_mas_carpool").
		Select("CAST(mas_carpool_uid AS TEXT) AS dept_sap, carpool_name AS dept_long_short, carpool_name AS dept_full, 'Car pool' AS dept_type").
		Where("is_deleted = ? AND is_active = ?", "0", "1").
		Find(&carpools).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve carpools", "message": messages.ErrInternalServer.Error()})
		return
	}
	results := append(driverDepts, carpools...)
	c.JSON(http.StatusOK, results)
}

// ListConfirmerLicenseUser godoc
// @Summary Retrieve the Confirmer License Users
// @Description	This endpoint allows a user to retrieve Confirmer License Users.
// @Tags MAS
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Param emp_id query string false "EmpID"
// @Router /api/mas/user-confirmer-license-users [get]
func (h *MasHandler) ListConfirmerLicenseUser(c *gin.Context) {
	h.ListConfirmerUser(c)
}

// ListApprovalLicenseUser godoc
// @Summary Retrieve the Approval License Users
// @Description	This endpoint allows a user to retrieve Approval License Users.
// @Tags MAS
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Param emp_id query string false "EmpID"
// @Router /api/mas/user-approval-license-users [get]
func (h *MasHandler) ListApprovalLicenseUser(c *gin.Context) {
	user := funcs.GetAuthenUser(c, "*")
	empID := c.Query("emp_id")
	if empID == "" {
		empID = user.EmpID
	}
	userInfo := funcs.GetUserEmpInfo(empID)
	managers := funcs.GetUserManager(userInfo.DeptSAP)
	list := []models.MasUserEmp{}
	for _, manager := range managers {
		if manager.LevelCode == "S1" {
			list = append(list, models.MasUserEmp{
				EmpID:        strconv.Itoa(manager.EmpIDLeader),
				FullName:     manager.EmpName,
				Position:     manager.PlansTextShort,
				DeptSAP:      strconv.Itoa(manager.DeptSAP),
				DeptSAPShort: funcs.GetDeptSAPShort(strconv.Itoa(manager.DeptSAP)),
				DeptSAPFull:  funcs.GetDeptSAPFull(strconv.Itoa(manager.DeptSAP)),
				ImageUrl:     funcs.GetEmpImage(strconv.Itoa(manager.EmpIDLeader)),
				IsEmployee:   true,
			})
		}
	}

	for _, manager := range managers {
		if manager.LevelCode == "M6" {
			list = append(list, models.MasUserEmp{
				EmpID:        strconv.Itoa(manager.EmpIDLeader),
				FullName:     manager.EmpName,
				Position:     manager.PlansTextShort,
				DeptSAP:      strconv.Itoa(manager.DeptSAP),
				DeptSAPShort: funcs.GetDeptSAPShort(strconv.Itoa(manager.DeptSAP)),
				DeptSAPFull:  funcs.GetDeptSAPFull(strconv.Itoa(manager.DeptSAP)),
				ImageUrl:     funcs.GetEmpImage(strconv.Itoa(manager.EmpIDLeader)),
				IsEmployee:   true,
			})
		}
	}

	if len(list) == 0 && len(managers) > 0 {
		upperManagers := funcs.GetUserManager(strconv.Itoa(managers[0].DeptUpper))
		for _, manager := range upperManagers {
			if manager.LevelCode == "S1" || manager.LevelCode == "M6" {
				list = append(list, models.MasUserEmp{
					EmpID:        strconv.Itoa(manager.EmpIDLeader),
					FullName:     manager.EmpName,
					Position:     manager.PlansTextShort,
					DeptSAP:      strconv.Itoa(manager.DeptSAP),
					DeptSAPShort: funcs.GetDeptSAPShort(strconv.Itoa(manager.DeptSAP)),
					DeptSAPFull:  funcs.GetDeptSAPFull(strconv.Itoa(manager.DeptSAP)),
					ImageUrl:     funcs.GetEmpImage(strconv.Itoa(manager.EmpIDLeader)),
					IsEmployee:   true,
				})
			}
		}
	}

	for i := range list {
		empInfo := funcs.GetUserEmpInfo(list[i].EmpID)
		list[i].TelMobile = empInfo.TelMobile
		list[i].TelInternal = empInfo.TelInternal
	}
	c.JSON(http.StatusOK, list)
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
							HolidaysDate: models.TimeWithZone{Time: d},
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
		return holidays[i].HolidaysDate.Before(holidays[j].HolidaysDate.Time)
	})
	c.JSON(http.StatusOK, holidays)
}
