package handlers

import (
	"fmt"
	"net/http"
	"sort"
	"strings"
	"time"
	"vms_plus_be/config"
	"vms_plus_be/funcs"
	"vms_plus_be/messages"
	"vms_plus_be/models"
	"vms_plus_be/userhub"

	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type BookingUserHandler struct {
	Role string
}

var MenuNameMapUser = map[string]string{
	"20,21,30,31,40,41,50,51,60,70,71": "กำลังดำเนินการ",
	"80":                               "เสร็จสิ้น",
	"90":                               "ยกเลิกคำขอ",
}

var StatusNameMapUser = map[string]string{
	"20": "รออนุมัติ",
	"21": "ถูกตีกลับ",
	"30": "รอตรวจสอบ",
	"31": "ถูกตีกลับ",
	"40": "รออนุมัติ",
	"41": "ถูกตีกลับ",
	"50": "รอรับกุญแจ",
	"51": "รอรับยานพาหนะ",
	"60": "เดินทาง",
	"70": "รอตรวจสอบ",
	"71": "คืนยานพาหนะไม่สำเร็จ",
	"80": "เสร็จสิ้น",
	"90": "ยกเลิกคำขอ",
}

func (h *BookingUserHandler) SetQueryRole(user *models.AuthenUserEmp, query *gorm.DB) *gorm.DB {
	return query.Where("created_request_emp_id = ?", user.EmpID)
}

func (h *BookingUserHandler) SetQueryStatusCanUpdate(query *gorm.DB) *gorm.DB {
	return query.Where("ref_request_status_code in ('21','31','41') and is_deleted = '0'")
}
func (h *BookingUserHandler) SetQueryStatusCanCancel(query *gorm.DB) *gorm.DB {
	return query.Where("ref_request_status_code in ('20','21','31','30','40','41') and is_deleted = '0'")
}

// CreateRequest godoc
// @Summary Create a new booking request
// @Description This endpoint allows a booking user to create a new request.
// @Tags Booking-user
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Param data body models.VmsTrnRequestRequest true "VmsTrnRequestRequest data"
// @Router /api/booking-user/create-request [post]
func (h *BookingUserHandler) CreateRequest(c *gin.Context) {
	user := funcs.GetAuthenUser(c, h.Role)
	if c.IsAborted() {
		return
	}

	var request models.VmsTrnRequestRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON input", "message": messages.ErrInvalidJSONInput.Error()})
		return
	}
	request.TrnRequestUID = uuid.New().String()
	request.CreatedAt = time.Now()
	request.CreatedBy = user.EmpID
	request.UpdatedAt = time.Now()
	request.UpdatedBy = user.EmpID
	request.IsDeleted = "0"

	createUser := funcs.GetUserEmpInfo(user.EmpID)
	request.CreatedRequestEmpID = createUser.EmpID
	request.CreatedRequestEmpName = createUser.FullName
	request.CreatedRequestDeptSAP = createUser.DeptSAP
	request.CreatedRequestDeptNameShort = createUser.DeptSAPShort
	request.CreatedRequestDeptNameFull = createUser.DeptSAPFull
	request.CreatedRequestDeskPhone = createUser.TelInternal
	request.CreatedRequestMobilePhone = createUser.TelMobile
	request.CreatedRequestPosition = createUser.Position
	request.CreatedRequestDatetime = models.TimeWithZone{Time: time.Now()}

	vehicleUser := funcs.GetUserEmpInfo(request.VehicleUserEmpID)
	request.VehicleUserEmpID = vehicleUser.EmpID
	request.VehicleUserEmpName = vehicleUser.FullName
	request.VehicleUserDeptSAP = vehicleUser.DeptSAP
	request.VehicleUserDeptNameShort = vehicleUser.DeptSAPShort
	request.VehicleUserDeptNameFull = vehicleUser.DeptSAPFull
	//request.VehicleUserDeskPhone = vehicleUser.TelInternal
	//request.VehicleUserMobilePhone = vehicleUser.TelMobile
	request.VehicleUserPosition = vehicleUser.Position

	confirmUser := funcs.GetUserEmpInfo(request.ConfirmedRequestEmpID)
	request.ConfirmedRequestEmpID = confirmUser.EmpID
	request.ConfirmedRequestEmpName = confirmUser.FullName
	request.ConfirmedRequestDeptSAP = confirmUser.DeptSAP
	request.ConfirmedRequestDeptNameShort = confirmUser.DeptSAPShort
	request.ConfirmedRequestDeptNameFull = confirmUser.DeptSAPFull
	request.ConfirmedRequestDeskPhone = confirmUser.TelInternal
	request.ConfirmedRequestMobilePhone = confirmUser.TelMobile
	request.ConfirmedRequestPosition = confirmUser.Position

	request.IsAdminChooseDriver = "0"
	request.RefRequestTypeCode = 1
	request.IsHaveSubRequest = "0"
	request.MasVehicleEvUID = ""

	if request.MasVehicleUID != nil && *request.MasVehicleUID != "" {
		var vehicle models.VmsMasVehicleDepartment
		if err := config.DB.First(&vehicle, "mas_vehicle_uid = ? AND is_deleted = '0'", request.MasVehicleUID).Error; err == nil {
			request.MasVehicleDepartmentUID = &vehicle.MasVehicleDepartmentUID
		}
		var carpool models.VmsMasCarpoolVehicle
		if err := config.DB.First(&carpool, "mas_vehicle_uid = ? AND is_deleted = '0'", request.MasVehicleUID).Error; err == nil {
			request.MasCarpoolUID = &carpool.MasCarpoolUID
		}
	}

	if request.MasCarpoolUID == nil || *request.MasCarpoolUID == "" {
		request.MasCarpoolUID = nil
	} else {
		var carpool models.VmsMasCarpoolRequest
		if err := config.DB.First(&carpool, "mas_carpool_uid = ? AND is_deleted = '0'", request.MasCarpoolUID).Error; err == nil {
			request.MasCarpoolUID = &carpool.MasCarpoolUID
			vehicleUser, _ := userhub.GetUserInfo(request.VehicleUserEmpID)

			if carpool.RefCarpoolChooseCarID == 3 {
				query := config.DB.Raw(`SELECT mas_vehicle_uid FROM fn_get_available_vehicles_view (?, ?, ?, ?) where mas_carpool_uid = ? and "CarTypeDetail" = ?`,
					request.ReserveStartDatetime, request.ReserveEndDatetime, vehicleUser.BureauDeptSap, vehicleUser.BusinessArea, request.MasCarpoolUID, request.RequestedVehicleType)

				var vehicle models.VmsMasVehicle
				if err := query.First(&vehicle).
					Select("mas_vehicle_uid").
					Order("vehicle_mileage").
					Limit(1).Error; err == nil && vehicle.MasVehicleUID != "" {
					request.MasVehicleUID = &vehicle.MasVehicleUID
				}
			}
			if carpool.RefCarpoolChooseDriverID == 3 {
				query := config.DB.Raw(`SELECT mas_driver_uid FROM fn_get_available_drivers_view (?, ?, ?, ?,?) where mas_carpool_uid = ?`,
					request.ReserveStartDatetime,
					request.ReserveEndDatetime,
					vehicleUser.BureauDeptSap,
					vehicleUser.BusinessArea,
					request.RefTripTypeCode,
					request.MasCarpoolUID)

				var driver models.VmsMasDriver
				if err := query.Scan(&driver).
					Select("mas_driver_uid, w_thismth.job_count, w_thismth.total_days").
					Joins("LEFT JOIN public.vms_trn_driver_monthly_workload AS w_thismth ON w_thismth.workload_year = ? AND w_thismth.workload_month = ? AND w_thismth.driver_emp_id = d.driver_id AND w_thismth.is_deleted = ?", request.ReserveStartDatetime.Year(), request.ReserveStartDatetime.Month(), "0").
					Order("total_days, job_count").
					Limit(1).Error; err == nil && driver.MasDriverUID != "" {
					request.MasCarPoolDriverUID = &driver.MasDriverUID
				}
			}
		}
	}

	if request.MasVehicleUID != nil && *request.MasVehicleUID != "" {
		var vehicle models.VmsMasVehicleDepartment
		if err := config.DB.First(&vehicle, "mas_vehicle_uid = ? AND is_deleted = '0'", request.MasVehicleUID).Error; err == nil {
			request.MasVehicleDepartmentUID = &vehicle.MasVehicleDepartmentUID
		}
	}

	if request.IsPEAEmployeeDriver == "1" && request.DriverEmpID != "" {
		driverUser := funcs.GetUserEmpInfo(request.DriverEmpID)
		request.DriverEmpID = driverUser.EmpID
		request.DriverEmpName = driverUser.FullName
		request.DriverEmpDeptSAP = driverUser.DeptSAP
		request.DriverEmpPosition = driverUser.Position
		request.DriverEmpDeptNameShort = funcs.GetDeptSAPShort(driverUser.DeptSAP)
		request.DriverEmpDeptNameFull = funcs.GetDeptSAPFull(driverUser.DeptSAP)
		//request.DriverEmpDeskPhone = driverUser.TelInternal
		//request.DriverEmpMobilePhone = driverUser.TelMobile

	} else if request.MasCarPoolDriverUID != nil && *request.MasCarPoolDriverUID != "" {
		var driver models.VmsMasDriver
		if err := config.DB.First(&driver, "mas_driver_uid = ? AND is_deleted = '0'", request.MasCarPoolDriverUID).Error; err == nil {
			request.DriverEmpID = driver.DriverID
			request.DriverEmpName = driver.DriverName
			request.DriverEmpDeptSAP = driver.DriverDeptSAP
			request.DriverEmpDeptNameShort = funcs.GetDeptSAPShort(driver.DriverDeptSAP)
			request.DriverEmpDeptNameFull = funcs.GetDeptSAPFull(driver.DriverDeptSAP)
		}
	}

	if request.MasVehicleDepartmentUID == nil || *request.MasVehicleDepartmentUID == "" {
		request.MasVehicleDepartmentUID = nil
	}
	if request.MasVehicleUID == nil || *request.MasVehicleUID == "" {
		request.MasVehicleUID = nil
	}
	if request.MasCarPoolDriverUID == nil || *request.MasCarPoolDriverUID == "" {
		request.MasCarPoolDriverUID = nil
	}

	YY := time.Now().Year() + 543
	BCode := vehicleUser.BusinessArea[0:1]
	var running int
	err := config.DB.Raw("SELECT nextval('vehicle_request_seq_' || lower(?))", BCode).Scan(&running).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate running number", "message": messages.ErrInternalServer.Error()})
		return
	}
	//'V' + BCode + 'YY' + 'RA' + Running 6 หลัก เช่น VZ68RA000001
	request.RequestNo = "V" + BCode + fmt.Sprintf("%02d", YY%100) + "RA" + fmt.Sprintf("%06d", running)
	request.RefRequestStatusCode = "20" // รออนุมัติจากต้นสังกัด

	if err := config.DB.Create(&request).
		Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create request", "message": messages.ErrCreateRequest.Error()})
		return
	}

	var result struct {
		models.VmsTrnRequestRequest
		RequestNo string `gorm:"column:request_no" json:"request_no"`
	}
	if err := config.DB.First(&result, "trn_request_uid = ? and is_deleted = ?", request.TrnRequestUID, "0").Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Booking not found", "message": messages.ErrBookingNotFound.Error()})
		return
	}
	funcs.CreateTrnRequestActionLog(result.TrnRequestUID,
		result.RefRequestStatusCode,
		"สร้างคำขอ",
		user.EmpID,
		"vehicle-user",
		"",
	)

	funcs.CheckMustPassStatus(request.TrnRequestUID)

	c.JSON(http.StatusCreated, gin.H{"message": "Request created successfully",
		"data":            result,
		"request_no":      request.RequestNo,
		"trn_request_uid": request.TrnRequestUID,
	})
}

// MenuRequests godoc
// @Summary Summary booking requests by request status code
// @Description Summary booking requests, counts grouped by request status code
// @Tags Booking-user
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Router /api/booking-user/menu-requests [get]
func (h *BookingUserHandler) MenuRequests(c *gin.Context) {
	user := funcs.GetAuthenUser(c, h.Role)
	if c.IsAborted() {
		return
	}
	query := h.SetQueryRole(user, config.DB)
	statusMenuMap := MenuNameMapUser
	summary, err := funcs.MenuRequests(statusMenuMap, query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "message": messages.ErrInternalServer.Error()})
		return
	}

	c.JSON(http.StatusOK, summary)
}

// SearchRequests godoc
// @Summary Search booking requests and get summary counts by request status code
// @Description Search for requests using a keyword and get the summary of counts grouped by request status code
// @Tags Booking-user
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Param search query string false "Search keyword (matches request_no, vehicle_license_plate, vehicle_user_emp_name, or work_place)"
// @Param ref_request_status_code query string false "Filter by multiple request status codes (comma-separated, e.g., 'A,B,C')"
// @Param startdate query string false "Filter by start datetime (YYYY-MM-DD format)"
// @Param enddate query string false "Filter by end datetime (YYYY-MM-DD format)"
// @Param order_by query string false "Order by request_no, start_datetime, ref_request_status_code"
// @Param order_dir query string false "Order direction: asc or desc"
// @Param page query int false "Page number (default: 1)"
// @Param limit query int false "Number of records per page (default: 10)"
// @Router /api/booking-user/search-requests [get]
func (h *BookingUserHandler) SearchRequests(c *gin.Context) {
	user := funcs.GetAuthenUser(c, h.Role)
	if c.IsAborted() {
		return
	}
	statusNameMap := StatusNameMapUser

	var requests []models.VmsTrnRequestList
	var summary []models.VmsTrnRequestSummary

	// Use the keys from statusNameMap as the list of valid status codes
	statusCodes := make([]string, 0, len(statusNameMap))
	for code := range statusNameMap {
		statusCodes = append(statusCodes, code)
	}

	query := h.SetQueryRole(user, config.DB)
	query = query.Table("public.vms_trn_request AS req").
		Select(`req.*, v.vehicle_license_plate,v.vehicle_license_plate_province_short,v.vehicle_license_plate_province_full,
			fn_get_long_short_dept_name_by_dept_sap(d.vehicle_owner_dept_sap) vehicle_department_dept_sap_short,       
			(select max(mc.carpool_name) from vms_mas_carpool mc where mc.mas_carpool_uid=req.mas_carpool_uid) vehicle_carpool_name,
			(select log.action_detail from vms_log_request_action log where log.trn_request_uid=req.trn_request_uid order by log.log_request_action_datetime desc limit 1) action_detail
		`).
		Joins("LEFT JOIN vms_mas_vehicle v on v.mas_vehicle_uid = req.mas_vehicle_uid").
		Joins("LEFT JOIN vms_mas_vehicle_department d on d.mas_vehicle_department_uid=req.mas_vehicle_department_uid").
		Where("req.ref_request_status_code IN (?)", statusCodes)
	query = query.Where("req.is_deleted = ?", "0")
	// Apply additional filters (search, date range, etc.)
	if search := c.Query("search"); search != "" {
		query = query.Where("req.request_no ILIKE ? OR v.vehicle_license_plate ILIKE ? OR req.vehicle_user_emp_name ILIKE ? OR req.work_place ILIKE ?", "%"+search+"%", "%"+search+"%", "%"+search+"%", "%"+search+"%")
	}
	if startDate := c.Query("startdate"); startDate != "" {
		query = query.Where("req.reserve_end_datetime >= ?", startDate)
	}
	if endDate := c.Query("enddate"); endDate != "" {
		query = query.Where("req.reserve_start_datetime <= ?", endDate)
	}
	if refRequestStatusCodes := c.Query("ref_request_status_code"); refRequestStatusCodes != "" {
		// Split the comma-separated codes into a slice
		codes := strings.Split(refRequestStatusCodes, ",")
		// Include additional keys with the same text in StatusNameMapUser
		additionalCodes := make(map[string]bool)
		for _, code := range codes {
			if name, exists := StatusNameMapUser[code]; exists {
				for key, value := range StatusNameMapUser {
					if value == name {
						additionalCodes[key] = true
					}
				}
			}
		}
		// Merge the original codes with the additional codes
		for key := range additionalCodes {
			codes = append(codes, key)
		}
		//fmt.Println("codes", codes)
		query = query.Where("req.ref_request_status_code IN (?)", codes)
	}

	// Ordering
	orderBy := c.Query("order_by")
	orderDir := c.Query("order_dir")
	if orderDir != "desc" {
		orderDir = "asc"
	}
	switch orderBy {
	case "request_no":
		query = query.Order("req.request_no " + orderDir)
	case "start_datetime":
		query = query.Order("req.start_datetime " + orderDir)
	case "ref_request_status_code":
		query = query.Order("req.ref_request_status_code " + orderDir)
	}

	// Pagination
	page := c.DefaultQuery("page", "1")
	limit := c.DefaultQuery("limit", "10")
	var pageInt, pageSizeInt int
	fmt.Sscanf(page, "%d", &pageInt)
	fmt.Sscanf(limit, "%d", &pageSizeInt)
	if pageInt < 1 {
		pageInt = 1
	}
	if pageSizeInt < 1 {
		pageSizeInt = 10
	}
	offset := (pageInt - 1) * pageSizeInt
	var total int64
	if err := query.Count(&total).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "message": messages.ErrInternalServer.Error()})
		return
	}

	query = query.Offset(offset).Limit(pageSizeInt)

	// Execute the main query
	if err := query.Scan(&requests).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "message": messages.ErrInternalServer.Error()})
		return
	}
	for i := range requests {
		requests[i].RefRequestStatusName = statusNameMap[requests[i].RefRequestStatusCode]
	}

	// Build the summary query
	summaryQuery := h.SetQueryRole(user, config.DB)
	summaryQuery = summaryQuery.Table("public.vms_trn_request AS req").
		Select("req.ref_request_status_code, COUNT(*) as count").
		Where("req.ref_request_status_code IN (?)", statusCodes).
		Group("req.ref_request_status_code")

	// Execute the summary query
	dbSummary := []struct {
		RefRequestStatusCode string `gorm:"column:ref_request_status_code"`
		Count                int    `gorm:"column:count"`
	}{}
	if err := summaryQuery.Scan(&dbSummary).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "message": messages.ErrInternalServer.Error()})
		return
	}

	// Create a complete summary with all statuses from statusNameMap
	groupedSummary := make(map[string]struct {
		Count   int
		MinCode string
	})

	// Aggregate counts and find the minimum RefRequestStatusCode for each RefRequestStatusName
	for _, dbItem := range dbSummary {
		name := statusNameMap[dbItem.RefRequestStatusCode]
		if data, exists := groupedSummary[name]; exists {
			groupedSummary[name] = struct {
				Count   int
				MinCode string
			}{
				Count:   data.Count + dbItem.Count,
				MinCode: min(data.MinCode, dbItem.RefRequestStatusCode),
			}
		} else {
			groupedSummary[name] = struct {
				Count   int
				MinCode string
			}{
				Count:   dbItem.Count,
				MinCode: dbItem.RefRequestStatusCode,
			}
		}
	}

	// Build the summary from the grouped data
	for name, data := range groupedSummary {
		summary = append(summary, models.VmsTrnRequestSummary{
			RefRequestStatusName: name,
			RefRequestStatusCode: data.MinCode,
			Count:                data.Count,
		})
	}
	// Sort the summary by RefRequestStatusCode
	sort.Slice(summary, func(i, j int) bool {
		return summary[i].RefRequestStatusCode < summary[j].RefRequestStatusCode
	})
	if requests == nil {
		requests = []models.VmsTrnRequestList{}
		summary = []models.VmsTrnRequestSummary{}
	}
	// Return both the filtered requests and the complete summary
	c.JSON(http.StatusOK, gin.H{
		"pagination": gin.H{
			"total":      total,
			"page":       page,
			"limit":      pageSizeInt,
			"totalPages": (total + int64(pageSizeInt) - 1) / int64(pageSizeInt), // Calculate total pages
		},
		"requests": requests,
		"summary":  summary,
	})
}

// GetRequest godoc
// @Summary Retrieve a specific booking request
// @Description This endpoint fetches details of a specific booking request using its unique identifier (TrnRequestUID).
// @Tags Booking-user
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Param trn_request_uid path string true "TrnRequestUID (trn_request_uid)"
// @Router /api/booking-user/request/{trn_request_uid} [get]
func (h *BookingUserHandler) GetRequest(c *gin.Context) {
	funcs.GetAuthenUser(c, h.Role)
	if c.IsAborted() {
		return
	}

	request, err := funcs.GetRequest(c, StatusNameMapUser)
	if err != nil {
		return
	}
	if request.RefRequestStatusCode == "20" {
		empUser := funcs.GetUserEmpInfo(request.ConfirmedRequestEmpID)
		request.ProgressRequestStatusEmp = models.ProgressRequestStatusEmp{
			ActionRole:   "ผู้อนุมัติต้นสังกัด",
			EmpID:        empUser.EmpID,
			EmpName:      empUser.FullName,
			EmpPosition:  empUser.Position,
			DeptSAP:      empUser.DeptSAP,
			DeptSAPShort: empUser.DeptSAPShort,
			DeptSAPFull:  empUser.DeptSAPFull,
			PhoneNumber:  empUser.TelInternal,
			MobileNumber: empUser.TelMobile,
		}
		request.ProgressRequestStatus = []models.ProgressRequestStatus{
			{ProgressIcon: "1", ProgressName: "รออนุมัติจากต้นสังกัด"},
			{ProgressIcon: "0", ProgressName: "รอผู้ดูแลยานพาหนะตรวจสอบ"},
			{ProgressIcon: "0", ProgressName: "รออนุมัติให้ใช้ยานพาหนะ"},
		}
	}

	if request.RefRequestStatusCode == "21" {
		request.ProgressRequestStatusEmp = funcs.GetProgressRequestStatusEmp(request.TrnRequestUID, "21", "ผู้อนุมัติต้นสังกัด")
		request.ProgressRequestStatus = []models.ProgressRequestStatus{
			{ProgressIcon: "2", ProgressName: "ถูกตีกลับจากต้นสังกัด"},
			{ProgressIcon: "0", ProgressName: "รอผู้ดูแลยานพาหนะตรวจสอบ"},
			{ProgressIcon: "0", ProgressName: "รออนุมัติให้ใช้ยานพาหนะ"},
		}
	}
	if request.RefRequestStatusCode == "30" {
		request.ProgressRequestStatus = []models.ProgressRequestStatus{
			{ProgressIcon: "3", ProgressName: "อนุมัติจากต้นสังกัด"},
			{ProgressIcon: "1", ProgressName: "รอผู้ดูแลยานพาหนะตรวจสอบ"},
			{ProgressIcon: "0", ProgressName: "รออนุมัติให้ใช้ยานพาหนะ"},
		}
		empIDs, err := funcs.GetAdminApprovalEmpIDs(request.TrnRequestUID)
		if err == nil && len(empIDs) > 0 {
			empUser := funcs.GetUserEmpInfo(empIDs[0])
			request.ProgressRequestStatusEmp = models.ProgressRequestStatusEmp{
				ActionRole:   "ผู้ดูแลยานพาหนะ",
				EmpID:        empUser.EmpID,
				EmpName:      empUser.FullName,
				EmpPosition:  empUser.Position,
				DeptSAP:      empUser.DeptSAP,
				DeptSAPShort: empUser.DeptSAPShort,
				DeptSAPFull:  empUser.DeptSAPFull,
				PhoneNumber:  empUser.TelInternal,
				MobileNumber: empUser.TelMobile,
			}
		}
	}
	if request.RefRequestStatusCode == "31" {
		request.ProgressRequestStatusEmp = funcs.GetProgressRequestStatusEmp(request.TrnRequestUID, "31", "ผู้ดูแลยานพาหนะ")
		request.ProgressRequestStatus = []models.ProgressRequestStatus{
			{ProgressIcon: "3", ProgressName: "อนุมัติจากต้นสังกัด"},
			{ProgressIcon: "2", ProgressName: "ถูกตีกลับจากผู้ดูแลยานพาหนะ"},
			{ProgressIcon: "0", ProgressName: "รออนุมัติให้ใช้ยานพาหนะ"},
		}
	}
	if request.RefRequestStatusCode == "40" {
		request.ProgressRequestStatusEmp = funcs.GetProgressRequestStatusEmp(request.TrnRequestUID, "40", "ผู้อนุมัติให้ใช้ยานพาหนะ")
		request.ProgressRequestStatus = []models.ProgressRequestStatus{
			{ProgressIcon: "3", ProgressName: "อนุมัติจากต้นสังกัด"},
			{ProgressIcon: "3", ProgressName: "อนุมัติจากผู้ดูแลยานพาหนะ"},
			{ProgressIcon: "1", ProgressName: "รออนุมัติให้ใช้ยานพาหนะ"},
		}
		empIDs, err := funcs.GetFinalApprovalEmpIDs(request.TrnRequestUID)
		if err == nil && len(empIDs) > 0 {
			empUser := funcs.GetUserEmpInfo(empIDs[0])
			request.ProgressRequestStatusEmp = models.ProgressRequestStatusEmp{
				ActionRole:   "ผู้อนุมัติให้ใช้ยานพาหนะ",
				EmpID:        empUser.EmpID,
				EmpName:      empUser.FullName,
				EmpPosition:  empUser.Position,
				DeptSAP:      empUser.DeptSAP,
				DeptSAPShort: empUser.DeptSAPShort,
				DeptSAPFull:  empUser.DeptSAPFull,
				PhoneNumber:  empUser.TelInternal,
				MobileNumber: empUser.TelMobile,
			}
		}
	}
	if request.RefRequestStatusCode == "41" {
		request.ProgressRequestStatusEmp = funcs.GetProgressRequestStatusEmp(request.TrnRequestUID, "41", "ผู้อนุมัติให้ใช้ยานพาหนะ")
		request.ProgressRequestStatus = []models.ProgressRequestStatus{
			{ProgressIcon: "3", ProgressName: "อนุมัติจากต้นสังกัด"},
			{ProgressIcon: "3", ProgressName: "อนุมัติจากผู้ดูแลยานพาหนะ"},
			{ProgressIcon: "2", ProgressName: "ถูกตีกลับจากเจ้าของยานพาหนะ"},
		}
	}
	if request.RefRequestStatusCode >= "50" && request.RefRequestStatusCode < "90" { //
		request.ProgressRequestStatusEmp = funcs.GetProgressRequestStatusEmp(request.TrnRequestUID, "50", "ผู้อนุมัติให้ใช้ยานพาหนะ")
		request.ProgressRequestStatus = []models.ProgressRequestStatus{
			{ProgressIcon: "3", ProgressName: "อนุมัติจากต้นสังกัด"},
			{ProgressIcon: "3", ProgressName: "อนุมัติจากผู้ดูแลยานพาหนะ"},
			{ProgressIcon: "3", ProgressName: "อนุมัติให้ใช้ยานพาหนะ"},
		}

	}
	if request.RefRequestStatusCode == "90" {
		request.ProgressRequestStatusEmp = funcs.GetProgressRequestStatusEmp(request.TrnRequestUID, "90", "ผู้ยกเลิกคำขอ")
		if request.CanceledRequestRole == "vehicle-user" {
			request.ProgressRequestStatus = []models.ProgressRequestStatus{
				{ProgressIcon: "2", ProgressName: "ยกเลิกจากผู้ขอใช้ยานพาหนะ"},
			}
		}
		if request.CanceledRequestRole == "level1-approval" {
			request.ProgressRequestStatus = []models.ProgressRequestStatus{
				{ProgressIcon: "2", ProgressName: "ยกเลิกจากต้นสังกัด"},
			}
		}
		if request.CanceledRequestRole == "admin-department" || request.CanceledRequestRole == "admin-carpool" {
			request.ProgressRequestStatus = []models.ProgressRequestStatus{
				{ProgressIcon: "3", ProgressName: "อนุมัติจากต้นสังกัด"},
				{ProgressIcon: "2", ProgressName: "ยกเลิกจากผู้ดูแลยานพาหนะ"},
			}
		}
		if request.CanceledRequestRole == "approval-department" || request.CanceledRequestRole == "approval-carpool" {
			request.ProgressRequestStatus = []models.ProgressRequestStatus{
				{ProgressIcon: "3", ProgressName: "อนุมัติจากต้นสังกัด"},
				{ProgressIcon: "3", ProgressName: "อนุมัติจากผู้ดูแลยานพาหนะ"},
				{ProgressIcon: "2", ProgressName: "ยกเลิกจากผู้ให้ใช้ยานพาหนะ"},
			}
		}
	}
	c.JSON(http.StatusOK, request)
}

// UpdateVehicleUser godoc
// @Summary Update vehicle information for a booking user
// @Description This endpoint allows a booking user to update the vehicle details associated with their request.
// @Tags Booking-user
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Param data body models.VmsTrnRequestVehicleUser true "VmsTrnRequestVehicleUser data"
// @Router /api/booking-user/update-vehicle-user [put]
func (h *BookingUserHandler) UpdateVehicleUser(c *gin.Context) {
	user := funcs.GetAuthenUser(c, h.Role)
	if c.IsAborted() {
		return
	}

	var request, trnRequest models.VmsTrnRequestVehicleUser
	var result struct {
		models.VmsTrnRequestVehicleUser
		models.VmsTrnRequestRequestNo
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "message": messages.ErrInvalidJSONInput.Error()})
		return
	}
	query := h.SetQueryRole(user, config.DB)
	query = h.SetQueryStatusCanUpdate(query)
	if err := query.First(&trnRequest, "trn_request_uid = ?", request.TrnRequestUID).Error; err != nil {
		c.JSON(http.StatusMethodNotAllowed, gin.H{"error": "Booking can not update", "message": messages.ErrBookingCannotUpdate.Error()})
		return
	}

	vehicleUser := funcs.GetUserEmpInfo(request.VehicleUserEmpID)
	request.VehicleUserEmpID = vehicleUser.EmpID
	request.VehicleUserEmpName = vehicleUser.FullName
	request.VehicleUserDeptSAP = vehicleUser.DeptSAP
	request.VehicleUserDeptNameShort = vehicleUser.DeptSAPShort
	request.VehicleUserDeptNameFull = vehicleUser.DeptSAPFull
	//request.VehicleUserDeskPhone = vehicleUser.TelInternal
	//request.VehicleUserMobilePhone = vehicleUser.TelMobile
	request.VehicleUserPosition = vehicleUser.Position

	request.UpdatedAt = time.Now()
	request.UpdatedBy = user.EmpID

	if err := config.DB.Save(&request).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to update : %v", err), "message": messages.ErrInternalServer.Error()})
		return
	}

	if err := config.DB.First(&result, "trn_request_uid = ?", request.TrnRequestUID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Booking not found", "message": messages.ErrBookingNotFound.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Updated successfully", "result": result})
}

// UpdateTrip godoc
// @Summary Update trip details for a booking request
// @Description This endpoint allows a booking user to update the details of an existing trip associated with their request.
// @Tags Booking-user
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Param data body models.VmsTrnRequestTrip true "VmsTrnRequestTrip data"
// @Router /api/booking-user/update-trip [put]
func (h *BookingUserHandler) UpdateTrip(c *gin.Context) {
	user := funcs.GetAuthenUser(c, h.Role)
	if c.IsAborted() {
		return
	}

	var request, trnRequest models.VmsTrnRequestTrip
	var result struct {
		models.VmsTrnRequestTrip
		models.VmsTrnRequestRequestNo
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "message": messages.ErrInvalidJSONInput.Error()})
		return
	}
	query := h.SetQueryRole(user, config.DB)
	query = h.SetQueryStatusCanUpdate(query)
	if err := query.First(&trnRequest, "trn_request_uid = ?", request.TrnRequestUID).Error; err != nil {
		c.JSON(http.StatusMethodNotAllowed, gin.H{"error": "Booking can not update", "message": messages.ErrBookingCannotUpdate.Error()})
		return
	}

	request.UpdatedAt = time.Now()
	request.UpdatedBy = user.EmpID

	if err := config.DB.Save(&request).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to update : %v", err), "message": messages.ErrInternalServer.Error()})
		return
	}

	if err := config.DB.First(&result, "trn_request_uid = ?", request.TrnRequestUID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Booking not found", "message": messages.ErrBookingNotFound.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Updated successfully", "result": result})
}

// UpdatePickup godoc
// @Summary Update pickup details for a booking request
// @Description This endpoint allows a booking user to update the pickup information for an existing booking.
// @Tags Booking-user
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Param data body models.VmsTrnRequestPickup true "VmsTrnRequestPickup data"
// @Router /api/booking-user/update-pickup [put]
func (h *BookingUserHandler) UpdatePickup(c *gin.Context) {
	user := funcs.GetAuthenUser(c, h.Role)
	if c.IsAborted() {
		return
	}
	var request, trnRequest models.VmsTrnRequestPickup
	var result struct {
		models.VmsTrnRequestPickup
		models.VmsTrnRequestRequestNo
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "message": messages.ErrInvalidJSONInput.Error()})
		return
	}

	query := h.SetQueryRole(user, config.DB)
	query = h.SetQueryStatusCanUpdate(query)
	if err := query.First(&trnRequest, "trn_request_uid = ?", request.TrnRequestUID).Error; err != nil {
		c.JSON(http.StatusMethodNotAllowed, gin.H{"error": "Booking can not update", "message": messages.ErrBookingCannotUpdate.Error()})
		return
	}
	request.UpdatedAt = time.Now()
	request.UpdatedBy = user.EmpID

	if err := config.DB.Save(&request).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to update : %v", err), "message": messages.ErrInternalServer.Error()})
		return
	}

	if err := config.DB.First(&result, "trn_request_uid = ?", request.TrnRequestUID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Booking not found", "message": messages.ErrBookingNotFound.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Updated successfully", "result": result})
}

// UpdateDocument godoc
// @Summary Update document details for a booking request
// @Description This endpoint allows a booking user to update the document associated with their booking request.
// @Tags Booking-user
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Param data body models.VmsTrnRequestDocument true "VmsTrnRequestDocument data"
// @Router /api/booking-user/update-document [put]
func (h *BookingUserHandler) UpdateDocument(c *gin.Context) {
	user := funcs.GetAuthenUser(c, h.Role)
	if c.IsAborted() {
		return
	}
	var request, trnRequest models.VmsTrnRequestDocument
	var result struct {
		models.VmsTrnRequestDocument
		models.VmsTrnRequestRequestNo
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "message": messages.ErrInvalidJSONInput.Error()})
		return
	}
	query := h.SetQueryRole(user, config.DB)
	query = h.SetQueryStatusCanUpdate(query)
	if err := query.First(&trnRequest, "trn_request_uid = ?", request.TrnRequestUID).Error; err != nil {
		c.JSON(http.StatusMethodNotAllowed, gin.H{"error": "Booking can not update", "message": messages.ErrBookingCannotUpdate.Error()})
		return
	}
	request.UpdatedAt = time.Now()
	request.UpdatedBy = user.EmpID

	if err := config.DB.Save(&request).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to update : %v", err), "message": messages.ErrInternalServer.Error()})
		return
	}

	if err := config.DB.First(&result, "trn_request_uid = ?", request.TrnRequestUID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Booking not found", "message": messages.ErrBookingNotFound.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Updated successfully", "result": result})
}

// UpdateCost godoc
// @Summary Update cost details for a booking request
// @Description This endpoint allows a booking user to update the cost information for an existing booking request.
// @Tags Booking-user
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Param data body models.VmsTrnRequestCost true "VmsTrnRequestCost data"
// @Router /api/booking-user/update-cost [put]
func (h *BookingUserHandler) UpdateCost(c *gin.Context) {
	user := funcs.GetAuthenUser(c, h.Role)
	if c.IsAborted() {
		return
	}
	var request, trnRequest models.VmsTrnRequestCost
	var result struct {
		models.VmsTrnRequestCost
		models.VmsTrnRequestRequestNo
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "message": messages.ErrInvalidJSONInput.Error()})
		return
	}

	query := h.SetQueryRole(user, config.DB)
	query = h.SetQueryStatusCanUpdate(query)
	if err := query.First(&trnRequest, "trn_request_uid = ?", request.TrnRequestUID).Error; err != nil {
		c.JSON(http.StatusMethodNotAllowed, gin.H{"error": "Booking can not update", "message": messages.ErrBookingCannotUpdate.Error()})
		return
	}
	request.UpdatedAt = time.Now()
	request.UpdatedBy = user.EmpID

	if err := config.DB.Save(&request).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to update : %v", err), "message": messages.ErrInternalServer.Error()})
		return
	}

	if err := config.DB.First(&result, "trn_request_uid = ?", request.TrnRequestUID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Booking not found", "message": messages.ErrBookingNotFound.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Updated successfully", "result": result})
}

// UpdateVehicleType godoc
// @Summary Update vehicle type for a booking request
// @Description This endpoint allows a booking user to update the vehicle type associated with their booking request.
// @Tags Booking-user
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Param data body models.VmsTrnRequestVehicleType true "VmsTrnRequestVehicleType data"
// @Router /api/booking-user/update-vehicle-type [put]
func (h *BookingUserHandler) UpdateVehicleType(c *gin.Context) {
	user := funcs.GetAuthenUser(c, h.Role)
	if c.IsAborted() {
		return
	}

	var request, trnRequest models.VmsTrnRequestVehicleType
	var result struct {
		models.VmsTrnRequestVehicleType
		models.VmsTrnRequestRequestNo
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "message": messages.ErrInvalidJSONInput.Error()})
		return
	}

	query := h.SetQueryRole(user, config.DB)
	query = h.SetQueryStatusCanUpdate(query)
	if err := query.First(&trnRequest, "trn_request_uid = ?", request.TrnRequestUID).Error; err != nil {
		c.JSON(http.StatusMethodNotAllowed, gin.H{"error": "Booking can not update", "message": messages.ErrBookingCannotUpdate.Error()})
		return
	}
	request.UpdatedAt = time.Now()
	request.UpdatedBy = user.EmpID

	if err := config.DB.Save(&request).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to update : %v", err), "message": messages.ErrInternalServer.Error()})
		return
	}

	if err := config.DB.First(&result, "trn_request_uid = ?", request.TrnRequestUID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Booking not found", "message": messages.ErrBookingNotFound.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Updated successfully", "result": result})
}

// UpdateConfirmer godoc
// @Summary Update confirmer for a booking request
// @Description This endpoint allows a booking user to update the confirmer user of their booking request.
// @Tags Booking-user
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Param data body models.VmsTrnRequestConfirmer true "VmsTrnRequestConfirmer data"
// @Router /api/booking-user/update-confirmer [put]
func (h *BookingUserHandler) UpdateConfirmer(c *gin.Context) {
	user := funcs.GetAuthenUser(c, h.Role)
	if c.IsAborted() {
		return
	}
	var request, trnRequest models.VmsTrnRequestConfirmer
	var result struct {
		models.VmsTrnRequestConfirmer
		models.VmsTrnRequestRequestNo
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "message": messages.ErrInvalidJSONInput.Error()})
		return
	}

	query := h.SetQueryRole(user, config.DB)
	query = h.SetQueryStatusCanUpdate(query)
	if err := query.First(&trnRequest, "trn_request_uid = ?", request.TrnRequestUID).Error; err != nil {
		c.JSON(http.StatusMethodNotAllowed, gin.H{"error": "Booking can not update", "message": messages.ErrBookingCannotUpdate.Error()})
		return
	}
	confirmUser := funcs.GetUserEmpInfo(request.ConfirmedRequestEmpID)
	request.ConfirmedRequestEmpID = confirmUser.EmpID
	request.ConfirmedRequestEmpName = confirmUser.FullName
	request.ConfirmedRequestDeptSAP = confirmUser.DeptSAP
	request.ConfirmedRequestDeptNameShort = confirmUser.DeptSAPShort
	request.ConfirmedRequestDeptNameFull = confirmUser.DeptSAPFull
	request.ConfirmedRequestDeskPhone = confirmUser.TelInternal
	request.ConfirmedRequestMobilePhone = confirmUser.TelMobile
	request.ConfirmedRequestPosition = confirmUser.Position

	request.UpdatedAt = time.Now()
	request.UpdatedBy = user.EmpID

	if err := config.DB.Save(&request).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to update : %v", err), "message": messages.ErrInternalServer.Error()})
		return
	}

	if err := config.DB.First(&result, "trn_request_uid = ?", request.TrnRequestUID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Booking not found", "message": messages.ErrBookingNotFound.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Updated successfully", "result": result})
}

// UpdateResend godoc
// @Summary re-send reqeust
// @Description This endpoint allows users re-send reqeust.
// @Tags Booking-user
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Param data body models.VmsTrnRequestResend true "VmsTrnRequestResend data"
// @Router /api/booking-user/update-resend [put]
func (h *BookingUserHandler) UpdateResend(c *gin.Context) {
	user := funcs.GetAuthenUser(c, h.Role)
	if c.IsAborted() {
		return
	}
	var request, trnRequest models.VmsTrnRequestResend
	var result struct {
		models.VmsTrnRequestResend
		models.VmsTrnRequestRequestNo
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "message": messages.ErrInvalidJSONInput.Error()})
		return
	}

	query := h.SetQueryRole(user, config.DB)
	query = h.SetQueryStatusCanUpdate(query)
	if err := query.First(&trnRequest, "trn_request_uid = ?", request.TrnRequestUID).Error; err != nil {
		c.JSON(http.StatusMethodNotAllowed, gin.H{"error": "Booking can not update", "message": messages.ErrBookingCannotUpdate.Error()})
		return
	}
	request.RefRequestStatusCode = "20"
	request.UpdatedAt = time.Now()
	request.UpdatedBy = user.EmpID

	if err := config.DB.Save(&request).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to update : %v", err), "message": messages.ErrInternalServer.Error()})
		return
	}

	if err := config.DB.First(&result, "trn_request_uid = ?", request.TrnRequestUID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Booking not found", "message": messages.ErrBookingNotFound.Error()})
		return
	}
	funcs.CreateTrnRequestActionLog(result.TrnRequestUID,
		result.RefRequestStatusCode,
		"ส่งคำขออีกครั้ง",
		user.EmpID,
		"vehicle-user",
		"",
	)
	funcs.CheckMustPassStatus(request.TrnRequestUID)
	c.JSON(http.StatusOK, gin.H{"message": "Updated successfully", "result": result})
}

// UpdateCanceled godoc
// @Summary Update cancel status for an item
// @Description This endpoint allows users to update the cancel status of an item.
// @Tags Booking-user
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Param data body models.VmsTrnRequestCanceled true "VmsTrnRequestCanceled data"
// @Router /api/booking-user/update-canceled [put]
func (h *BookingUserHandler) UpdateCanceled(c *gin.Context) {

	user := funcs.GetAuthenUser(c, h.Role)
	if c.IsAborted() {
		return
	}
	var request, trnRequest models.VmsTrnRequestCanceled
	var result struct {
		models.VmsTrnRequestCanceled
		models.VmsTrnRequestRequestNo
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "message": messages.ErrInvalidJSONInput.Error()})
		return
	}

	query := h.SetQueryRole(user, config.DB)
	query = h.SetQueryStatusCanCancel(query)
	if err := query.First(&trnRequest, "trn_request_uid = ?", request.TrnRequestUID).Error; err != nil {
		c.JSON(http.StatusMethodNotAllowed, gin.H{"error": "Booking can not update", "message": messages.ErrBookingCannotUpdate.Error()})
		return
	}
	request.RefRequestStatusCode = "90"
	request.UpdatedAt = time.Now()
	request.UpdatedBy = user.EmpID

	cancelUser := funcs.GetUserEmpInfo(user.EmpID)
	request.CanceledRequestEmpID = cancelUser.EmpID
	request.CanceledRequestEmpName = cancelUser.FullName
	request.CanceledRequestDeptSAP = cancelUser.DeptSAP
	request.CanceledRequestDeptNameShort = cancelUser.DeptSAPShort
	request.CanceledRequestDeptNameFull = cancelUser.DeptSAPFull
	request.CanceledRequestDeskPhone = cancelUser.TelInternal
	request.CanceledRequestMobilePhone = cancelUser.TelMobile
	request.CanceledRequestPosition = cancelUser.Position
	request.CanceledRequestDatetime = models.TimeWithZone{Time: time.Now()}

	if err := config.DB.Save(&request).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to update : %v", err), "message": messages.ErrInternalServer.Error()})
		return
	}

	if err := config.DB.First(&result, "trn_request_uid = ?", request.TrnRequestUID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Booking not found", "message": messages.ErrBookingNotFound.Error()})
		return
	}
	funcs.CreateTrnRequestActionLog(result.TrnRequestUID,
		result.RefRequestStatusCode,
		"ยกเลิกคำขอ",
		user.EmpID,
		"vehicle-user",
		request.CanceledRequestReason,
	)

	c.JSON(http.StatusOK, gin.H{"message": "Updated successfully", "result": result})
}
