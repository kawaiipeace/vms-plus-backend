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

type BookingAdminDeptHandler struct {
	Role string
}

func (h *BookingAdminDeptHandler) SetQueryRole(user *models.AuthenUserEmp, query *gorm.DB) *gorm.DB {
	if user.EmpID == "" {
		return query
	}
	return query.Where("created_request_emp_id = ?", user.EmpID)
}
func (h *BookingAdminDeptHandler) SetQueryStatusCanUpdate(query *gorm.DB) *gorm.DB {
	return query.Where("ref_request_status_code in ('30') and is_deleted = '0'")
}

// MenuRequests godoc
// @Summary Summary booking requests by request status code
// @Description Summary booking requests, counts grouped by request status code
// @Tags Booking-admin-dept
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Router /api/booking-admin-dept/menu-requests [get]
func (h *BookingAdminDeptHandler) MenuRequests(c *gin.Context) {
	// Get authenticated user role if needed
	user := funcs.GetAuthenUser(c, h.Role)
	if c.IsAborted() {
		return
	}

	statusMenuMap := MenuNameMapAdmin
	query := h.SetQueryRole(user, config.DB)
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
// @Tags Booking-admin-dept
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Param search query string false "Search keyword (matches request_no, vehicle_license_plate, vehicle_user_emp_name, or work_place)"
// @Param ref_request_status_code query string false "Filter by multiple request status codes (comma-separated, e.g., 'A,B,C')"
// @Param startdate query string false "Filter by start datetime (YYYY-MM-DD format)"
// @Param enddate query string false "Filter by end datetime (YYYY-MM-DD format)"
// @Param vehicle_owner_dept_sap query string false "Filter by vehicle owner department SAP"
// @Param order_by query string false "Order by request_no, start_datetime, ref_request_status_code"
// @Param order_dir query string false "Order direction: asc or desc"
// @Param page query int false "Page number (default: 1)"
// @Param limit query int false "Number of records per page (default: 10)"
// @Router /api/booking-admin-dept/search-requests [get]
func (h *BookingAdminDeptHandler) SearchRequests(c *gin.Context) {
	user := funcs.GetAuthenUser(c, h.Role)
	if c.IsAborted() {
		return
	}
	statusNameMap := StatusNameMapAdmin

	var requests []models.VmsTrnRequestAdminList
	var summary []models.VmsTrnRequestSummary

	// Use the keys from statusNameMap as the list of valid status codes
	statusCodes := make([]string, 0, len(statusNameMap))
	for code := range statusNameMap {
		statusCodes = append(statusCodes, code)
	}

	// Build the main query
	query := h.SetQueryRole(user, config.DB)
	query = query.Table("public.vms_trn_request AS req").
		Select(`
		req.*, 
		v.vehicle_license_plate, v.vehicle_license_plate_province_short, v.vehicle_license_plate_province_full,
		COALESCE(NULLIF(req.is_pea_employee_driver::VARCHAR, '1'), d.driver_name) AS driver_name,
		COALESCE(NULLIF(req.is_pea_employee_driver::VARCHAR, '1'), d.driver_dept_sap_hire) AS driver_dept_name,
		md.dept_short AS vehicle_dept_name,
		mc.carpool_name AS vehicle_carpool_name
	`).
		Joins("LEFT JOIN vms_mas_vehicle AS v ON v.mas_vehicle_uid = req.mas_vehicle_uid").
		Joins("LEFT JOIN vms_mas_driver AS d ON d.mas_driver_uid = req.mas_carpool_driver_uid").
		Joins("LEFT JOIN vms_mas_vehicle_department AS mvd ON mvd.mas_vehicle_uid = req.mas_vehicle_uid").
		Joins("LEFT JOIN vms_mas_department AS md ON md.dept_sap = mvd.vehicle_owner_dept_sap").
		Joins("LEFT JOIN vms_mas_carpool_vehicle AS mcv ON mcv.mas_vehicle_uid = req.mas_vehicle_uid").
		Joins("LEFT JOIN vms_mas_carpool AS mc ON mc.mas_carpool_uid = mcv.mas_carpool_uid").
		Where("req.ref_request_status_code IN (?)", statusCodes)

	query = query.Where("req.is_deleted = ?", "0")
	// Apply additional filters (search, date range, etc.)
	if search := c.Query("search"); search != "" {
		query = query.Where("req.request_no ILIKE ? OR req.vehicle_license_plate ILIKE ? OR req.vehicle_user_emp_name ILIKE ? OR req.work_place ILIKE ?", "%"+search+"%", "%"+search+"%", "%"+search+"%", "%"+search+"%")
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
		if requests[i].IsAdminChooseDriver == 1 && requests[i].IsPEAEmployeeDriver == 0 && (requests[i].MasCarpoolDriverUID == "" || requests[i].MasCarpoolDriverUID == funcs.DefaultUUID()) {
			requests[i].Can_Choose_Driver = true
		}
		if requests[i].IsAdminChooseVehicle == 1 && (requests[i].MasVehicleUID == "" || requests[i].MasVehicleUID == funcs.DefaultUUID()) {
			requests[i].Can_Choose_Vehicle = true
		}
		if requests[i].TripType == 1 {
			requests[i].TripTypeName = "ไป-กลับ"
		} else if requests[i].TripType == 2 {
			requests[i].TripTypeName = "ค้างแรม"
		}
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
	// Sort the summary by RefRequestStatusCode
	sort.Slice(summary, func(i, j int) bool {
		return summary[i].RefRequestStatusCode < summary[j].RefRequestStatusCode
	})
	if requests == nil {
		requests = []models.VmsTrnRequestAdminList{}
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
// @Tags Booking-admin-dept
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Param trn_request_uid path string true "TrnRequestUID (trn_request_uid)"
// @Router /api/booking-admin-dept/request/{trn_request_uid} [get]
func (h *BookingAdminDeptHandler) GetRequest(c *gin.Context) {
	funcs.GetAuthenUser(c, h.Role)
	if c.IsAborted() {
		return
	}
	request, err := funcs.GetRequest(c, StatusNameMapAdminDetail)
	if err != nil {
		return
	}
	if request.RefRequestStatusCode == "30" {
		request.ProgressRequestStatus = []models.ProgressRequestStatus{
			{ProgressIcon: "1", ProgressName: "รอผู้ดูแลยานพาหนะตรวจสอบ"},
			{ProgressIcon: "0", ProgressName: "รออนุมัติให้ใช้ยานพาหนะ"},
		}
	}
	if request.RefRequestStatusCode == "31" {
		request.ProgressRequestStatus = []models.ProgressRequestStatus{
			{ProgressIcon: "2", ProgressName: "ถูกตีกลับจากผู้ดูแลยานพาหนะ"},
			{ProgressIcon: "0", ProgressName: "รออนุมัติให้ใช้ยานพาหนะ"},
		}
	}
	if request.RefRequestStatusCode >= "40" {
		request.ProgressRequestStatus = []models.ProgressRequestStatus{
			{ProgressIcon: "3", ProgressName: "อนุมัติจากผู้ดูแลยานพาหนะ"},
			{ProgressIcon: "1", ProgressName: "รออนุมัติให้ใช้ยานพาหนะ"},
		}
	}
	if request.RefRequestStatusCode >= "41" {
		request.ProgressRequestStatus = []models.ProgressRequestStatus{
			{ProgressIcon: "3", ProgressName: "อนุมัติจากผู้ดูแลยานพาหนะ"},
			{ProgressIcon: "2", ProgressName: "ถูกตีกลับจากเจ้าของยานพาหนะ"},
		}
	}
	if request.RefRequestStatusCode >= "50" && request.RefRequestStatusCode < "90" { //
		request.ProgressRequestStatus = []models.ProgressRequestStatus{
			{ProgressIcon: "3", ProgressName: "อนุมัติจากผู้ดูแลยานพาหนะ"},
			{ProgressIcon: "3", ProgressName: "อนุมัติให้ใช้ยานพาหนะ"},
		}

	}
	if request.RefRequestStatusCode == "90" {
		request.ProgressRequestStatus = []models.ProgressRequestStatus{
			{ProgressIcon: "2", ProgressName: "ยกเลิก"},
		}
	}
	if request.RefRequestStatusCode == "91" {
		request.ProgressRequestStatus = []models.ProgressRequestStatus{
			{ProgressIcon: "2", ProgressName: "ยกเลิกจากผู้ขอใช้ยานพาหนะ"},
		}
	}
	if request.RefRequestStatusCode == "92" {
		request.ProgressRequestStatus = []models.ProgressRequestStatus{
			{ProgressIcon: "2", ProgressName: "ยกเลิกจากต้นสังกัด"},
		}
	}
	if request.RefRequestStatusCode == "93" {
		request.ProgressRequestStatus = []models.ProgressRequestStatus{
			{ProgressIcon: "3", ProgressName: "อนุมัติจากต้นสังกัด"},
			{ProgressIcon: "2", ProgressName: "ยกเลิกจากผู้ดูแลยานพาหนะ"},
		}
	}
	if request.RefRequestStatusCode == "94" {
		request.ProgressRequestStatus = []models.ProgressRequestStatus{
			{ProgressIcon: "3", ProgressName: "อนุมัติจากต้นสังกัด"},
			{ProgressIcon: "3", ProgressName: "อนุมัติจากผู้ดูแลยานพาหนะ"},
			{ProgressIcon: "2", ProgressName: "ยกเลิกจากผู้ให้ใช้ยานพาหนะ"},
		}
	}
	c.JSON(http.StatusOK, request)
}

// UpdateRejected godoc
// @Summary Update rejected status for an item
// @Description This endpoint allows users to update the rejected status of an item.
// @Tags Booking-admin-dept
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Param data body models.VmsTrnRequestRejected true "VmsTrnRequestRejected data"
// @Router /api/booking-admin-dept/update-rejected [put]
func (h *BookingAdminDeptHandler) UpdateRejected(c *gin.Context) {
	user := funcs.GetAuthenUser(c, h.Role)
	if c.IsAborted() {
		return
	}
	var request, trnRequest models.VmsTrnRequestRejected
	var result struct {
		models.VmsTrnRequestRejected
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
	request.RefRequestStatusCode = "31"
	rejectUser := funcs.GetUserEmpInfo(user.EmpID)
	request.RejectedRequestEmpID = rejectUser.EmpID
	request.RejectedRequestEmpName = rejectUser.FullName
	request.RejectedRequestDeptSAP = rejectUser.DeptSAP
	request.RejectedRequestDeptNameShort = rejectUser.DeptSAPShort
	request.RejectedRequestDeptNameFull = rejectUser.DeptSAPFull
	request.RejectedRequestDeskPhone = rejectUser.TelInternal
	request.RejectedRequestMobilePhone = rejectUser.TelMobile
	request.RejectedRequestPosition = rejectUser.Position
	request.RejectedRequestDatetime = time.Now()
	request.UpdatedAt = time.Now()
	request.UpdatedBy = user.EmpID

	if err := config.DB.Save(&request).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to update : %v", err), "message": messages.ErrInternalServer.Error()})
		return
	}
	if err := config.DB.First(&result, "trn_request_uid = ?", request.TrnRequestUID).Error; err != nil {
		c.JSON(http.StatusMethodNotAllowed, gin.H{"error": "Booking can not update", "message": messages.ErrBookingCannotUpdate.Error()})
		return
	}
	funcs.CreateTrnRequestActionLog(request.TrnRequestUID,
		request.RefRequestStatusCode,
		"ผู้ดูแลยานพาหนะตีกลับคำขอ",
		user.EmpID,
		"admin-dept-approval",
		request.RejectedRequestReason,
	)

	c.JSON(http.StatusOK, gin.H{"message": "Updated successfully", "result": result})
}

// UpdateApproved godoc
// @Summary Update Approved status for an item
// @Description This endpoint allows users to update the sended back status of an item.
// @Tags Booking-admin-dept
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Param data body models.VmsTrnRequestApprovedWithRecieiveKey true "VmsTrnRequestApprovedWithRecieiveKey data"
// @Router /api/booking-admin-dept/update-approved [put]
func (h *BookingAdminDeptHandler) UpdateApproved(c *gin.Context) {
	user := funcs.GetAuthenUser(c, h.Role)
	if c.IsAborted() {
		return
	}
	var request, trnRequest models.VmsTrnRequestApprovedWithRecieiveKey
	var result struct {
		models.VmsTrnRequestApprovedWithRecieiveKey
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
	if err := config.DB.Save(&request).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to update : %v", err), "message": messages.ErrInternalServer.Error()})
		return
	}

	requestStatus := models.VmsTrnRequestUpdateRecieivedKeyStatus{
		RefRequestStatusCode: "40", //
		TrnRequestUID:        request.TrnRequestUID,
		UpdatedAt:            time.Now(),
		UpdatedBy:            user.EmpID,
	}
	if err := config.DB.Save(&requestStatus).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to update : %v", err), "message": messages.ErrInternalServer.Error()})
		return
	}

	if err := config.DB.First(&result, "trn_request_uid = ?", request.TrnRequestUID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Booking not found", "message": messages.ErrBookingNotFound.Error()})
		return
	}

	funcs.CreateTrnRequestActionLog(request.TrnRequestUID,
		requestStatus.RefRequestStatusCode,
		"ผู้ดูแลยานพาหนะ ตรวจสอบผ่านแล้ว",
		user.EmpID,
		"admin-approval",
		"",
	)

	c.JSON(http.StatusOK, gin.H{"message": "Updated successfully", "result": result})
}

// UpdateCanceled godoc
// @Summary Update cancel status for an item
// @Description This endpoint allows users to update the cancel status of an item.
// @Tags Booking-admin-dept
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Param data body models.VmsTrnRequestCanceled true "VmsTrnRequestCanceled data"
// @Router /api/booking-admin-dept/update-canceled [put]
func (h *BookingAdminDeptHandler) UpdateCanceled(c *gin.Context) {
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
	query = h.SetQueryStatusCanUpdate(query)
	if err := query.First(&trnRequest, "trn_request_uid = ?", request.TrnRequestUID).Error; err != nil {
		c.JSON(http.StatusMethodNotAllowed, gin.H{"error": "Booking can not update", "message": messages.ErrBookingCannotUpdate.Error()})
		return
	}

	request.RefRequestStatusCode = "90" //
	cancelUser := funcs.GetUserEmpInfo(user.EmpID)
	request.CanceledRequestEmpID = cancelUser.EmpID
	request.CanceledRequestEmpName = cancelUser.FullName
	request.CanceledRequestDeptSAP = cancelUser.DeptSAP
	request.CanceledRequestDeptNameShort = cancelUser.DeptSAPShort
	request.CanceledRequestDeptNameFull = cancelUser.DeptSAPFull
	request.CanceledRequestDeskPhone = cancelUser.TelInternal
	request.CanceledRequestMobilePhone = cancelUser.TelMobile
	request.CanceledRequestPosition = cancelUser.Position
	request.CanceledRequestDatetime = time.Now()

	if err := config.DB.Save(&request).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to update : %v", err), "message": messages.ErrInternalServer.Error()})
		return
	}

	if err := config.DB.First(&result, "trn_request_uid = ?", request.TrnRequestUID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Booking not found", "message": messages.ErrBookingNotFound.Error()})
		return
	}
	if result.RefRequestStatusCode == request.RefRequestStatusCode {
		funcs.CreateTrnRequestActionLog(request.TrnRequestUID,
			request.RefRequestStatusCode,
			"ยกเลิกคำขอ",
			user.EmpID,
			"admin-dept-approval",
			request.CanceledRequestReason,
		)
	}
	c.JSON(http.StatusOK, gin.H{"message": "Updated successfully", "result": result})
}

// UpdateVehicleUser godoc
// @Summary Update vehicle information for a booking user
// @Description This endpoint allows a booking user to update the vehicle details associated with their request.
// @Tags Booking-admin-dept
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Param data body models.VmsTrnRequestVehicleUser true "VmsTrnRequestVehicleUser data"
// @Router /api/booking-admin-dept/update-vehicle-user [put]
func (h *BookingAdminDeptHandler) UpdateVehicleUser(c *gin.Context) {
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
	request.VehicleUserDeskPhone = vehicleUser.TelInternal
	request.VehicleUserMobilePhone = vehicleUser.TelMobile
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
// @Tags Booking-admin-dept
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Param data body models.VmsTrnRequestTrip true "VmsTrnRequestTrip data"
// @Router /api/booking-admin-dept/update-trip [put]
func (h *BookingAdminDeptHandler) UpdateTrip(c *gin.Context) {
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
// @Tags Booking-admin-dept
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Param data body models.VmsTrnRequestPickup true "VmsTrnRequestPickup data"
// @Router /api/booking-admin-dept/update-pickup [put]
func (h *BookingAdminDeptHandler) UpdatePickup(c *gin.Context) {
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
// @Tags Booking-admin-dept
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Param data body models.VmsTrnRequestDocument true "VmsTrnRequestDocument data"
// @Router /api/booking-admin-dept/update-document [put]
func (h *BookingAdminDeptHandler) UpdateDocument(c *gin.Context) {
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
// @Tags Booking-admin-dept
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Param data body models.VmsTrnRequestCost true "VmsTrnRequestCost data"
// @Router /api/booking-admin-dept/update-cost [put]
func (h *BookingAdminDeptHandler) UpdateCost(c *gin.Context) {
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

// UpdateDriver godoc
// @Summary Update driver for a booking request
// @Description This endpoint allows a booking user to update the cost information for an existing booking request.
// @Tags Booking-admin-dept
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Param data body models.VmsTrnRequestDriver true "VmsTrnRequestDriver data"
// @Router /api/booking-admin-dept/update-driver [put]
func (h *BookingAdminDeptHandler) UpdateDriver(c *gin.Context) {
	user := funcs.GetAuthenUser(c, h.Role)
	if c.IsAborted() {
		return
	}
	var request, trnRequest models.VmsTrnRequestDriver
	var result struct {
		models.VmsTrnRequestDriver
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

// UpdateVehicle godoc
// @Summary Update vehicle for a booking request
// @Description This endpoint allows a booking user to update the cost information for an existing booking request.
// @Tags Booking-admin-dept
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Param data body models.VmsTrnRequestVehicle true "VmsTrnRequestVehicle data"
// @Router /api/booking-admin-dept/update-vehicle [put]
func (h *BookingAdminDeptHandler) UpdateVehicle(c *gin.Context) {
	user := funcs.GetAuthenUser(c, h.Role)
	if c.IsAborted() {
		return
	}
	var request, trnRequest models.VmsTrnRequestVehicle
	var result struct {
		models.VmsTrnRequestVehicle
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
	var vehicle models.VmsMasVehicle
	if err := config.DB.First(&vehicle, "mas_vehicle_uid = ? AND is_deleted = '0'", request.MasVehicleUID).Error; err == nil {
		request.VehicleLicensePlate = strings.TrimSpace(vehicle.VehicleLicensePlate)
		request.VehicleLicensePlateProvinceShort = strings.TrimSpace(vehicle.VehicleLicensePlateProvinceShort)
		request.VehicleLicensePlateProvinceFull = strings.TrimSpace(vehicle.VehicleLicensePlateProvinceFull)
	}

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
