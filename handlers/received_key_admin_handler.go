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

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ReceivedKeyAdminHandler struct {
	Role string
}

var StatusNameMapReceivedKeyAdmin = map[string]string{
	"50":  "รอรับกุญแจ",
	"50e": "เกินวันนัดหมาย",
	"51":  "รับยานพาหนะ",
}

func (h *ReceivedKeyAdminHandler) SetQueryRole(user *models.AuthenUserEmp, query *gorm.DB) *gorm.DB {
	query = funcs.SetQueryAdminRole(user, query)
	return query
}
func (h *ReceivedKeyAdminHandler) SetQueryStatusCanUpdate(query *gorm.DB) *gorm.DB {
	return query.Where("ref_request_status_code in ('50','60','70') and is_deleted = '0'")
}

// SearchRequests godoc
// @Summary Search booking requests and get summary counts by request status code
// @Description Search for requests using a keyword and get the summary of counts grouped by request status code
// @Tags Received-key-admin
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Param search query string false "Search keyword (matches request_no, vehicle_license_plate, vehicle_user_emp_name, or work_place)"
// @Param ref_request_status_code query string false "Filter by multiple request status codes (comma-separated, e.g., 'A,B,C')"
// @Param startdate query string false "Filter by start datetime (YYYY-MM-DD format)"
// @Param enddate query string false "Filter by end datetime (YYYY-MM-DD format)"
// @Param vehicle_owner_dept_sap query string false "Filter by vehicle owner department SAP"
// @Param received_key_start_datetime query string false "Filter by received key start datetime (YYYY-MM-DD format)"
// @Param received_key_end_datetime query string false "Filter by received key end datetime (YYYY-MM-DD format)"
// @Param order_by query string false "Order by request_no, start_datetime, ref_request_status_code"
// @Param order_dir query string false "Order direction: asc or desc"
// @Param page query int false "Page number (default: 1)"
// @Param limit query int false "Number of records per page (default: 10)"
// @Router /api/received-key-admin/search-requests [get]
func (h *ReceivedKeyAdminHandler) SearchRequests(c *gin.Context) {
	user := funcs.GetAuthenUser(c, h.Role)
	if c.IsAborted() {
		return
	}
	var statusNameMap = StatusNameMapReceivedKeyUser
	var requests []models.VmsTrnRequestAdminList
	var summary []models.VmsTrnRequestSummary

	// Use the keys from statusNameMap as the list of valid status codes
	statusCodes := make([]string, 0, len(statusNameMap))
	for code := range statusNameMap {
		statusCodes = append(statusCodes, code)
	}

	// Build the main query
	query := h.SetQueryRole(user, config.DB)
	query = query.Table("public.vms_trn_request").
		Select("vms_trn_request.*, v.vehicle_license_plate,v.vehicle_license_plate_province_short,v.vehicle_license_plate_province_full,"+
			"case vms_trn_request.is_pea_employee_driver when '1' then vms_trn_request.driver_emp_name else (select driver_name from vms_mas_driver d where d.mas_driver_uid=vms_trn_request.mas_carpool_driver_uid) end driver_name,"+
			"case vms_trn_request.is_pea_employee_driver when '1' then vms_trn_request.driver_emp_dept_name_short else (select driver_dept_sap_short_work from vms_mas_driver d where d.mas_driver_uid=vms_trn_request.mas_carpool_driver_uid) end driver_dept_name,"+
			"fn_get_long_short_dept_name_by_dept_sap(d.vehicle_owner_dept_sap) vehicle_department_dept_sap_short,"+
			"mc.carpool_name vehicle_carpool_name,"+
			"(select Max(parking_place) from vms_mas_vehicle_department d where d.mas_vehicle_uid = vms_trn_request.mas_vehicle_uid and d.is_deleted = '0' and d.is_active = '1') parking_place, "+
			"k.receiver_personal_id,k.receiver_fullname,k.receiver_dept_sap,"+
			"k.appointment_start appointment_key_handover_start_datetime,k.appointment_end appointment_key_handover_end_datetime,k.appointment_location appointment_key_handover_place,"+
			"k.receiver_dept_name_short,k.receiver_dept_name_full,k.receiver_desk_phone,k.receiver_mobile_phone,k.receiver_position,k.remark receiver_remark").
		Joins("LEFT JOIN vms_trn_vehicle_key_handover k ON k.trn_request_uid = vms_trn_request.trn_request_uid").
		Joins("LEFT JOIN vms_mas_vehicle v on v.mas_vehicle_uid = vms_trn_request.mas_vehicle_uid").
		Joins("LEFT JOIN vms_mas_vehicle_department d on d.mas_vehicle_department_uid=vms_trn_request.mas_vehicle_department_uid").
		Joins("LEFT JOIN vms_mas_carpool mc ON mc.mas_carpool_uid = vms_trn_request.mas_carpool_uid").
		Where("vms_trn_request.ref_request_status_code IN (?)", statusCodes)
	query = query.Where("vms_trn_request.is_deleted = ?", "0")

	// Apply additional filters (search, date range, etc.)
	if search := c.Query("search"); search != "" {
		query = query.Where("vms_trn_request.request_no ILIKE ? OR v.vehicle_license_plate ILIKE ? OR vms_trn_request.vehicle_user_emp_name ILIKE ? OR vms_trn_request.work_place ILIKE ?", "%"+search+"%", "%"+search+"%", "%"+search+"%", "%"+search+"%")
	}
	if startDate := c.Query("startdate"); startDate != "" {
		query = query.Where("vms_trn_request.reserve_end_datetime >= ?", startDate)
	}
	if endDate := c.Query("enddate"); endDate != "" {
		query = query.Where("vms_trn_request.reserve_start_datetime <= ?", endDate)
	}

	if receivedKeyStartDatetime := c.Query("received_key_start_datetime"); receivedKeyStartDatetime != "" {
		query = query.Where("vms_trn_request.received_key_start_datetime >= ?", receivedKeyStartDatetime)
	}
	if receivedKeyEndDatetime := c.Query("received_key_end_datetime"); receivedKeyEndDatetime != "" {
		query = query.Where("vms_trn_request.received_key_end_datetime <= ?", receivedKeyEndDatetime)
	}
	if refRequestStatusCodes := c.Query("ref_request_status_code"); refRequestStatusCodes != "" {
		// Split the comma-separated codes into a slice
		codes := strings.Split(refRequestStatusCodes, ",")
		// Include additional keys with the same text in StatusNameMapUser
		additionalCodes := make(map[string]bool)
		for _, code := range codes {
			if name, exists := statusNameMap[code]; exists {
				for key, value := range statusNameMap {
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
		fmt.Println("codes", codes)
		query = query.Where("vms_trn_request.ref_request_status_code IN (?)", codes)
	}
	// Ordering
	orderBy := c.Query("order_by")
	orderDir := c.Query("order_dir")
	if orderDir != "desc" {
		orderDir = "asc"
	}
	switch orderBy {
	case "request_no":
		query = query.Order("vms_trn_request.request_no " + orderDir)
	case "start_datetime":
		query = query.Order("vms_trn_request.start_datetime " + orderDir)
	case "ref_request_status_code":
		query = query.Order("vms_trn_request.ref_request_status_code " + orderDir)
	default:
		query = query.Order("vms_trn_request.request_no desc")
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	query = query.Offset(offset).Limit(pageSizeInt)

	// Execute the main query
	if err := query.Scan(&requests).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	for i := range requests {
		requests[i].RefRequestStatusName = statusNameMap[requests[i].RefRequestStatusCode]

		switch requests[i].TripType {
		case 0:
			requests[i].TripTypeName = "ไป-กลับ"
		case 1:
			requests[i].TripTypeName = "ค้างแรม"
		}
		if requests[i].IsPEAEmployeeDriver != 1 && requests[i].DriverCarpoolName != "" {
			requests[i].DriverDeptName = requests[i].DriverCarpoolName
		}
		if requests[i].VehicleCarpoolName != "" {
			requests[i].VehicleDepartmentDeptSapShort = requests[i].VehicleCarpoolName
			requests[i].VehicleDeptName = requests[i].VehicleCarpoolName
			requests[i].VehicleCarpoolText = "Carpool"
			requests[i].VehicleCarpoolName = "Carpool"
		} else {
			requests[i].VehicleDeptName = requests[i].VehicleDepartmentDeptSapShort
			requests[i].VehicleCarpoolName = ""
			requests[i].VehicleCarpoolText = ""

		}
	}
	// Build the summary query
	summaryQuery := h.SetQueryRole(user, config.DB)
	summaryQuery = summaryQuery.Table("public.vms_trn_request").
		Select("vms_trn_request.ref_request_status_code, COUNT(*) as count").
		Where("vms_trn_request.ref_request_status_code IN (?)", statusCodes).
		Group("vms_trn_request.ref_request_status_code")

	// Execute the summary query
	dbSummary := []struct {
		RefRequestStatusCode string `gorm:"column:ref_request_status_code"`
		Count                int    `gorm:"column:count"`
	}{}
	if err := summaryQuery.Scan(&dbSummary).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Create a complete summary with all statuses from statusNameMap
	for code, name := range statusNameMap {
		count := 0
		for _, dbItem := range dbSummary {
			if dbItem.RefRequestStatusCode == code {
				count = dbItem.Count
				break
			}
		}
		summary = append(summary, models.VmsTrnRequestSummary{
			RefRequestStatusCode: code,
			RefRequestStatusName: name,
			Count:                count,
		})
	}
	if requests == nil {
		requests = []models.VmsTrnRequestAdminList{}
		summary = []models.VmsTrnRequestSummary{}
	}
	// Sort the summary by RefRequestStatusCode
	sort.Slice(summary, func(i, j int) bool {
		return summary[i].RefRequestStatusCode < summary[j].RefRequestStatusCode
	})
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
// @Tags Received-key-admin
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Param trn_request_uid path string true "TrnRequestUID (trn_request_uid)"
// @Router /api/received-key-admin/request/{trn_request_uid} [get]
func (h *ReceivedKeyAdminHandler) GetRequest(c *gin.Context) {
	funcs.GetAuthenUser(c, h.Role)
	if c.IsAborted() {
		return
	}
	request, err := funcs.GetRequestVehicelInUse(c, StatusNameMapReceivedKeyAdmin)
	if err != nil {
		return
	}
	c.JSON(http.StatusOK, request)
}

// UpdateRecieivedKey godoc
// @Summary Update key pickup driver for a booking request
// @Description This endpoint allows to update key pickup driver.
// @Tags Received-key-admin
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Param data body models.VmsTrnRequestUpdateRecieivedKey true "VmsTrnRequestUpdateRecieivedKey data"
// @Router /api/received-key-admin/update-recieived-key [put]
func (h *ReceivedKeyAdminHandler) UpdateRecieivedKey(c *gin.Context) {
	user := funcs.GetAuthenUser(c, h.Role)
	if c.IsAborted() {
		return
	}
	var request, trnRequest models.VmsTrnRequestUpdateRecieivedKey
	var result struct {
		models.VmsTrnRequestUpdateRecieivedKey
		models.VmsTrnRequestRequestNo
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "message": messages.ErrInvalidJSONInput.Error()})
		return
	}

	query := h.SetQueryRole(user, config.DB).Table("public.vms_trn_request")
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

	//update vms_trn_request set appointment_key_handover_place,appointment_key_handover_start_datetime,appointment_key_handover_end_datetime
	if err := config.DB.Table("vms_trn_request").
		Where("trn_request_uid = ?", request.TrnRequestUID).
		Update("appointment_key_handover_place", request.ReceivedKeyPlace).
		Update("appointment_key_handover_start_datetime", request.ReceivedKeyStartDatetime).
		Update("appointment_key_handover_end_datetime", request.ReceivedKeyEndDatetime).Error; err != nil {
		return
	}

	if err := config.DB.
		Table("public.vms_trn_request").
		Select("trn_request_uid,request_no,appointment_key_handover_place appointment_location,appointment_key_handover_start_datetime appointment_start,appointment_key_handover_end_datetime appointment_end").
		First(&result, "trn_request_uid = ?", request.TrnRequestUID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Booking not found", "message": messages.ErrBookingNotFound.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Updated successfully", "result": result})
}

// UpdateRecieivedKeyDetail godoc
// @Summary Confirm the update of key pickup driver for a booking request
// @Description This endpoint allows confirming the update of the key pickup driver for a specific booking request.
// @Tags Received-key-admin
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Param data body models.VmsTrnRequestUpdateRecieivedKeyDetail true "VmsTrnRequestUpdateRecieivedKeyDetail data"
// @Router /api/received-key-admin/update-recieived-key-detail [put]
func (h *ReceivedKeyAdminHandler) UpdateRecieivedKeyDetail(c *gin.Context) {
	user := funcs.GetAuthenUser(c, h.Role)
	if c.IsAborted() {
		return
	}
	var request models.VmsTrnRequestUpdateRecieivedKeyConfirmed
	var trnRequest models.VmsTrnRequestList
	var result struct {
		models.VmsTrnRequestUpdateRecieivedKeyDetail
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
	result.RequestNo = trnRequest.RequestNo
	c.JSON(http.StatusOK, gin.H{"message": "Updated successfully", "result": result})
}

// UpdateKeyPickupDriver godoc
// @Summary Update key pickup driver for a booking request
// @Description This endpoint allows to update key pickup driver.
// @Tags Received-key-admin
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Param data body models.VmsTrnReceivedKeyDriver true "VmsTrnReceivedKeyDriver data"
// @Router /api/received-key-admin/update-key-pickup-driver [put]
func (h *ReceivedKeyAdminHandler) UpdateKeyPickupDriver(c *gin.Context) {
	user := funcs.GetAuthenUser(c, h.Role)
	if c.IsAborted() {
		return
	}
	var request models.VmsTrnReceivedKeyDriver
	var trnRequest models.VmsTrnRequestResponse
	var result struct {
		models.VmsTrnReceivedKeyDriver
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
	request.ReceiverType = 1 // Driver
	request.ReceiverPersonalId = trnRequest.DriverEmpID
	request.ReceiverFullname = trnRequest.DriverEmpName
	request.ReceiverDeptSAP = trnRequest.DriverEmpDeptSAP
	request.ReceiverDeptNameShort = trnRequest.DriverEmpDeptNameShort
	request.ReceiverDeptNameFull = trnRequest.DriverEmpDeptNameFull
	request.ReceiverPosition = trnRequest.DriverEmpPosition
	request.ReceiverMobilePhone = trnRequest.DriverMobileContact
	request.ReceiverDeskPhone = trnRequest.DriverInternalContact
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
	requestStatus := models.VmsTrnRequestUpdateRecieivedKeyStatus{
		RefRequestStatusCode: "51", // "รับกุญแจยานพาหนะแล้ว รอบันทึกข้อมูลเมื่อนำยานพาหนะออกปฎิบัติงาน"
		TrnRequestUID:        request.TrnRequestUID,
		UpdatedAt:            time.Now(),
		UpdatedBy:            user.EmpID,
	}
	if err := config.DB.Save(&requestStatus).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to update : %v", err), "message": messages.ErrInternalServer.Error()})
		return
	}

	var parkingPlace string
	if err := config.DB.Table("public.vms_trn_request AS req").
		Joins("LEFT JOIN vms_mas_vehicle_department d on d.mas_vehicle_uid = req.mas_vehicle_uid AND d.is_deleted = '0' AND d.is_active = '1'").
		Select("d.parking_place").
		Where("req.trn_request_uid = ?", request.TrnRequestUID).
		First(&parkingPlace).Error; err != nil {
		parkingPlace = ""
	}

	funcs.CreateTrnRequestActionLog(result.TrnRequestUID,
		requestStatus.RefRequestStatusCode,
		"สถานที่ "+parkingPlace+" สถานที่จอดรถ",
		user.EmpID,
		"admin-department",
		"",
	)
	result.RequestNo = trnRequest.RequestNo
	c.JSON(http.StatusOK, gin.H{"message": "Updated successfully", "result": result})
}

// UpdateKeyPickupPEA godoc
// @Summary Update key pickup emp user for a booking request
// @Description This endpoint allows to update key pickup emp user.
// @Tags Received-key-admin
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Param data body models.VmsTrnReceivedKeyPEA true "VmsTrnReceivedKeyPEA data"
// @Router /api/received-key-admin/update-key-pickup-pea [put]
func (h *ReceivedKeyAdminHandler) UpdateKeyPickupPEA(c *gin.Context) {
	user := funcs.GetAuthenUser(c, h.Role)
	if c.IsAborted() {
		return
	}
	var request models.VmsTrnReceivedKeyPEA
	var trnRequest models.VmsTrnRequestList
	var result struct {
		models.VmsTrnReceivedKeyPEA
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
	request.ReceiverType = 2 // PEA
	empUser := funcs.GetUserEmpInfo(request.ReceiverPersonalId)
	request.ReceiverFullname = empUser.FullName
	request.ReceiverDeptSAP = empUser.DeptSAP
	request.ReceiverDeptNameShort = empUser.DeptSAPShort
	request.ReceiverDeptNameFull = empUser.DeptSAPFull
	request.UpdatedAt = time.Now()
	request.UpdatedBy = user.EmpID

	if err := config.DB.Save(&request).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to update : %v", err)})
		return
	}

	if err := config.DB.First(&result, "trn_request_uid = ?", request.TrnRequestUID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Booking not found"})
		return
	}
	requestStatus := models.VmsTrnRequestUpdateRecieivedKeyStatus{
		RefRequestStatusCode: "51", // "รับกุญแจยานพาหนะแล้ว รอบันทึกข้อมูลเมื่อนำยานพาหนะออกปฎิบัติงาน"
		TrnRequestUID:        request.TrnRequestUID,
		UpdatedAt:            time.Now(),
		UpdatedBy:            user.EmpID,
	}
	if err := config.DB.Save(&requestStatus).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to update : %v", err), "message": messages.ErrInternalServer.Error()})
		return
	}
	var parkingPlace string
	if err := config.DB.Table("public.vms_trn_request AS req").
		Joins("LEFT JOIN vms_mas_vehicle_department d on d.mas_vehicle_uid = req.mas_vehicle_uid AND d.is_deleted = '0' AND d.is_active = '1'").
		Select("d.parking_place").
		Where("req.trn_request_uid = ?", request.TrnRequestUID).
		First(&parkingPlace).Error; err != nil {
		parkingPlace = ""
	}

	funcs.CreateTrnRequestActionLog(result.TrnRequestUID,
		requestStatus.RefRequestStatusCode,
		"สถานที่ "+parkingPlace+" สถานที่จอดรถ",
		user.EmpID,
		"admin-department",
		"",
	)
	result.RequestNo = trnRequest.RequestNo
	c.JSON(http.StatusOK, gin.H{"message": "Updated successfully", "result": result})
}

// UpdateKeyPickupOutSider godoc
// @Summary Update key pickup outsource for a booking request
// @Description This endpoint allows to update key pickup outsource.
// @Tags Received-key-admin
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Param data body models.VmsTrnReceivedKeyOutSider true "VmsTrnReceivedKeyOutSider data"
// @Router /api/received-key-admin/update-key-pickup-outsider [put]
func (h *ReceivedKeyAdminHandler) UpdateKeyPickupOutSider(c *gin.Context) {
	user := funcs.GetAuthenUser(c, h.Role)
	if c.IsAborted() {
		return
	}
	var request models.VmsTrnReceivedKeyOutSider
	var trnRequest models.VmsTrnRequestList
	var result struct {
		models.VmsTrnReceivedKeyOutSider
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
	request.ReceiverType = 3 // Outsider
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
	requestStatus := models.VmsTrnRequestUpdateRecieivedKeyStatus{
		RefRequestStatusCode: "51", // "รับกุญแจยานพาหนะแล้ว รอบันทึกข้อมูลเมื่อนำยานพาหนะออกปฎิบัติงาน"
		TrnRequestUID:        request.TrnRequestUID,
		UpdatedAt:            time.Now(),
		UpdatedBy:            user.EmpID,
	}
	if err := config.DB.Save(&requestStatus).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to update : %v", err), "message": messages.ErrInternalServer.Error()})
		return
	}
	var parkingPlace string
	if err := config.DB.Table("public.vms_trn_request AS req").
		Joins("LEFT JOIN vms_mas_vehicle_department d on d.mas_vehicle_uid = req.mas_vehicle_uid AND d.is_deleted = '0' AND d.is_active = '1'").
		Select("d.parking_place").
		Where("req.trn_request_uid = ?", request.TrnRequestUID).
		First(&parkingPlace).Error; err != nil {
		parkingPlace = ""
	}

	funcs.CreateTrnRequestActionLog(result.TrnRequestUID,
		requestStatus.RefRequestStatusCode,
		"สถานที่ "+parkingPlace+" สถานที่จอดรถ",
		user.EmpID,
		"admin-department",
		"",
	)
	result.RequestNo = trnRequest.RequestNo
	c.JSON(http.StatusOK, gin.H{"message": "Updated successfully", "result": result})
}

// UpdateCanceled godoc
// @Summary Update cancel status for an item
// @Description This endpoint allows users to update the cancel status of an item.
// @Tags Received-key-admin
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Param data body models.VmsTrnRequestCanceled true "VmsTrnRequestCanceled data"
// @Router /api/received-key-admin/update-canceled [put]
func (h *ReceivedKeyAdminHandler) UpdateCanceled(c *gin.Context) {
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
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	query := h.SetQueryRole(user, config.DB)
	query = h.SetQueryStatusCanUpdate(query)
	if err := query.First(&trnRequest, "trn_request_uid = ?", request.TrnRequestUID).Error; err != nil {
		c.JSON(http.StatusMethodNotAllowed, gin.H{"error": "Booking can not update", "message": messages.ErrBookingCannotUpdate.Error()})
		return
	}

	request.RefRequestStatusCode = "90" // ยกเลิกคำขอ
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
	request.UpdatedAt = time.Now()
	request.UpdatedBy = user.EmpID

	if err := config.DB.Save(&request).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to update : %v", err)})
		return
	}

	if err := config.DB.First(&result, "trn_request_uid = ?", request.TrnRequestUID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Booking not found"})
		return
	}
	funcs.CreateTrnRequestActionLog(request.TrnRequestUID,
		request.RefRequestStatusCode,
		"ยกเลิกคำขอ",
		user.EmpID,
		"admin-department",
		request.CanceledRequestReason,
	)

	c.JSON(http.StatusOK, gin.H{"message": "Updated successfully", "result": result})
}
