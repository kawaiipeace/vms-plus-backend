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
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BookingAdminHandler struct {
	Role string
}

var MenuNameMapAdmin = map[string]string{
	"30,31,40": "ตรวจสอบคำขอ",
	"50,51":    "ให้กุญแจ",
	"60":       "เดินทาง",
	"70,71":    "ตรวจสอบยานพาหนะ",
	"80":       "เสร็จสิ้น",
	"90":       "ยกเลิก",
}

var StatusNameMapAdmin = map[string]string{
	"30": "รอตรวจสอบ",
	"31": "ตีกลับ",
	"40": "รออนุมัติ",
	"90": "ยกเลิกคำขอ",
}

var StatusNameMapAdminDetail = map[string]string{
	"30": "รอตรวจสอบ",
	"31": "ตีกลับ",
	"40": "รออนุมัติ",
	"70": "ตรวจสอบยานพาหนะ",
	"71": "ตรวจสอบยานพาหนะไม่ผ่าน",
	"80": "เสร็จสิ้น",
	"90": "ยกเลิกคำขอ",
}

func (h *BookingAdminHandler) SetQueryRole(user *models.AuthenUserEmp, query *gorm.DB) *gorm.DB {
	return funcs.SetQueryAdminRole(user, query)
}

func (h *BookingAdminHandler) SetQueryStatusCanUpdate(query *gorm.DB) *gorm.DB {
	return query.Where("ref_request_status_code in ('30') and is_deleted = '0'")
}

// MenuRequests godoc
// @Summary Summary booking requests by request status code
// @Description Summary booking requests, counts grouped by request status code
// @Tags Booking-admin
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Router /api/booking-admin/menu-requests [get]
func (h *BookingAdminHandler) MenuRequests(c *gin.Context) {
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
// @Tags Booking-admin
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Param search query string false "Search keyword (matches request_no, vehicle_license_plate, vehicle_user_emp_name, or work_place)"
// @Param ref_request_status_code query string false "Filter by multiple request status codes (comma-separated, e.g., 'A,B,C')"
// @Param startdate query string false "Filter by start datetime (YYYY-MM-DD format)"
// @Param enddate query string false "Filter by end datetime (YYYY-MM-DD format)"
// @Param ref_request_status_code query string false "Filter by multiple request status codes (comma-separated, e.g., 'A,B,C')"
// @Param vehicle_owner_dept_sap query string false "Filter by vehicle owner department SAP"
// @Param order_by query string false "Order by request_no, start_datetime, ref_request_status_code"
// @Param order_dir query string false "Order direction: asc or desc"
// @Param page query int false "Page number (default: 1)"
// @Param limit query int false "Number of records per page (default: 10)"
// @Router /api/booking-admin/search-requests [get]
func (h *BookingAdminHandler) SearchRequests(c *gin.Context) {
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
	query = query.Table("public.vms_trn_request").
		Select(
			`vms_trn_request.*,
			v.vehicle_license_plate,v.vehicle_license_plate_province_short,v.vehicle_license_plate_province_full,
			case vms_trn_request.is_pea_employee_driver when '1' then vms_trn_request.driver_emp_name else (select driver_name from vms_mas_driver d where d.mas_driver_uid=vms_trn_request.mas_carpool_driver_uid) end driver_name,
			case vms_trn_request.is_pea_employee_driver when '1' then vms_trn_request.driver_emp_dept_name_short else (select driver_dept_sap_short_work from vms_mas_driver d where d.mas_driver_uid=vms_trn_request.mas_carpool_driver_uid) end driver_dept_name,
			fn_get_long_short_dept_name_by_dept_sap(d.vehicle_owner_dept_sap) vehicle_department_dept_sap_short,       
			mc.carpool_name vehicle_carpool_name,ref_carpool_choose_car_id,ref_carpool_choose_driver_id,
			(select max(md.vendor_name) from vms_mas_driver md where md.mas_driver_uid=vms_trn_request.mas_carpool_driver_uid) driver_vendor_name,
			(select max(mc.carpool_name) from vms_mas_carpool mc,vms_mas_carpool_driver mcd where mc.mas_carpool_uid=mcd.mas_carpool_uid and mcd.mas_driver_uid=vms_trn_request.mas_carpool_driver_uid) driver_carpool_name
		`).
		Joins("LEFT JOIN vms_mas_vehicle_department d on d.mas_vehicle_department_uid=vms_trn_request.mas_vehicle_department_uid").
		Joins("LEFT JOIN vms_mas_vehicle AS v ON v.mas_vehicle_uid = vms_trn_request.mas_vehicle_uid").
		Joins("LEFT JOIN public.vms_ref_request_status AS status ON vms_trn_request.ref_request_status_code = status.ref_request_status_code").
		Joins("LEFT JOIN vms_mas_carpool mc ON mc.mas_carpool_uid = vms_trn_request.mas_carpool_uid").
		Where("vms_trn_request.ref_request_status_code IN (?)", statusCodes)

	query = query.Where("vms_trn_request.is_deleted = ?", "0")

	if search := c.Query("search"); search != "" {
		query = query.Where("vms_trn_request.request_no ILIKE ? OR v.vehicle_license_plate ILIKE ? OR vms_trn_request.vehicle_user_emp_name ILIKE ? OR vms_trn_request.work_place ILIKE ?", "%"+search+"%", "%"+search+"%", "%"+search+"%", "%"+search+"%")
	}
	if startDate := c.Query("startdate"); startDate != "" {
		query = query.Where("vms_trn_request.reserve_end_datetime >= ?", startDate)
	}
	if endDate := c.Query("enddate"); endDate != "" {
		query = query.Where("vms_trn_request.reserve_start_datetime <= ?", endDate)
	}
	if refRequestStatusCodes := c.Query("ref_request_status_code"); refRequestStatusCodes != "" {
		// Split the comma-separated codes into a slice
		codes := strings.Split(refRequestStatusCodes, ",")
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
		if requests[i].RefCarpoolChooseDriverID == 2 && requests[i].IsPEAEmployeeDriver == 0 && requests[i].MasCarpoolDriverUID == nil {
			requests[i].CanChooseDriver = true
		}
		if requests[i].RefCarpoolChooseCarID == 2 && requests[i].MasVehicleUID == nil {
			requests[i].CanChooseVehicle = true
		}
		if requests[i].TripType == 0 {
			requests[i].TripTypeName = "ไป-กลับ"
		} else if requests[i].TripType == 1 {
			requests[i].TripTypeName = "ค้างแรม"
		}
		if requests[i].IsPEAEmployeeDriver != 1 && requests[i].DriverCarpoolName != "" {
			requests[i].DriverDeptName = requests[i].DriverCarpoolName
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
// @Tags Booking-admin
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Param trn_request_uid path string true "TrnRequestUID (trn_request_uid)"
// @Router /api/booking-admin/request/{trn_request_uid} [get]
func (h *BookingAdminHandler) GetRequest(c *gin.Context) {
	user := funcs.GetAuthenUser(c, h.Role)
	if c.IsAborted() {
		return
	}
	request, err := funcs.GetRequest(c, StatusNameMapAdminDetail)
	if err != nil {
		return
	}
	if request.RefRequestStatusCode == "30" {
		request.ProgressRequestStatus = []models.ProgressRequestStatus{
			{ProgressIcon: "3", ProgressName: "อนุมัติจากต้นสังกัด"},
			{ProgressIcon: "1", ProgressName: "รอผู้ดูแลยานพาหนะตรวจสอบ"},
			{ProgressIcon: "0", ProgressName: "รออนุมัติให้ใช้ยานพาหนะ"},
		}
		empUser := funcs.GetUserEmpInfo(user.EmpID)
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
	if request.RefRequestStatusCode == "31" {
		request.ProgressRequestStatus = []models.ProgressRequestStatus{
			{ProgressIcon: "3", ProgressName: "อนุมัติจากต้นสังกัด"},
			{ProgressIcon: "2", ProgressName: "ถูกตีกลับจากผู้ดูแลยานพาหนะ"},
			{ProgressIcon: "0", ProgressName: "รออนุมัติให้ใช้ยานพาหนะ"},
		}
		request.ProgressRequestStatusEmp = funcs.GetProgressRequestStatusEmp(request.TrnRequestUID, "31", "ผู้ดูแลยานพาหนะ")

	}
	if request.RefRequestStatusCode == "40" {
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
		request.ProgressRequestStatus = []models.ProgressRequestStatus{
			{ProgressIcon: "3", ProgressName: "อนุมัติจากต้นสังกัด"},
			{ProgressIcon: "3", ProgressName: "อนุมัติจากผู้ดูแลยานพาหนะ"},
			{ProgressIcon: "2", ProgressName: "ถูกตีกลับจากเจ้าของยานพาหนะ"},
		}
		request.ProgressRequestStatusEmp = funcs.GetProgressRequestStatusEmp(request.TrnRequestUID, "41", "ผู้อนุมัติให้ใช้ยานพาหนะ")
	}

	if request.RefRequestStatusCode >= "50" && request.RefRequestStatusCode < "90" { //
		request.ProgressRequestStatus = []models.ProgressRequestStatus{
			{ProgressIcon: "3", ProgressName: "อนุมัติจากต้นสังกัด"},
			{ProgressIcon: "3", ProgressName: "อนุมัติจากผู้ดูแลยานพาหนะ"},
			{ProgressIcon: "3", ProgressName: "อนุมัติให้ใช้ยานพาหนะ"},
		}
		request.ProgressRequestStatusEmp = funcs.GetProgressRequestStatusEmp(request.TrnRequestUID, "50", "ผู้อนุมัติให้ใช้ยานพาหนะ")
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
		if request.CanceledRequestRole == "admin-approval" {
			request.ProgressRequestStatus = []models.ProgressRequestStatus{
				{ProgressIcon: "3", ProgressName: "อนุมัติจากต้นสังกัด"},
				{ProgressIcon: "2", ProgressName: "ยกเลิกจากผู้ดูแลยานพาหนะ"},
			}
		}
		if request.CanceledRequestRole == "final-approval" {
			request.ProgressRequestStatus = []models.ProgressRequestStatus{
				{ProgressIcon: "3", ProgressName: "อนุมัติจากต้นสังกัด"},
				{ProgressIcon: "3", ProgressName: "อนุมัติจากผู้ดูแลยานพาหนะ"},
				{ProgressIcon: "2", ProgressName: "ยกเลิกจากผู้ให้ใช้ยานพาหนะ"},
			}
		}

	}
	c.JSON(http.StatusOK, request)
}

// UpdateRejected godoc
// @Summary Update rejected status for an item
// @Description This endpoint allows users to update the rejected status of an item.
// @Tags Booking-admin
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Param data body models.VmsTrnRequestRejected true "VmsTrnRequestRejected data"
// @Router /api/booking-admin/update-rejected [put]
func (h *BookingAdminHandler) UpdateRejected(c *gin.Context) {
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
		c.JSON(http.StatusMethodNotAllowed, gin.H{"error": "Booking not found", "message": messages.ErrBookingNotFound.Error()})
		return
	}
	funcs.CreateTrnRequestActionLog(request.TrnRequestUID,
		request.RefRequestStatusCode,
		"ผู้ดูแลยานพาหนะตีกลับคำขอ",
		user.EmpID,
		"admin-approval",
		request.RejectedRequestReason,
	)

	c.JSON(http.StatusOK, gin.H{"message": "Updated successfully", "result": result})
}

// UpdateApproved godoc
// @Summary Update Approved status for an item
// @Description This endpoint allows users to update the sended back status of an item.
// @Tags Booking-admin
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Param data body models.VmsTrnRequestApprovedWithRecieiveKey true "VmsTrnRequestApprovedWithRecieiveKey data"
// @Router /api/booking-admin/update-approved [put]
func (h *BookingAdminHandler) UpdateApproved(c *gin.Context) {
	user := funcs.GetAuthenUser(c, h.Role)
	if c.IsAborted() {
		return
	}
	var request models.VmsTrnRequestApprovedWithRecieiveKey
	var trnRequest models.VmsTrnRequestList
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
	request.HandoverUID = uuid.New().String()
	request.ReceiverType = 0
	request.CreatedBy = user.EmpID
	request.CreatedAt = time.Now()
	request.UpdatedBy = user.EmpID
	request.UpdatedAt = time.Now()
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

	var trnRequestList models.VmsTrnRequestList
	if err := config.DB.First(&trnRequestList, "trn_request_uid = ?", request.TrnRequestUID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Booking not found", "message": messages.ErrBookingNotFound.Error()})
		return
	}
	result.RequestNo = trnRequestList.RequestNo

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
// @Tags Booking-admin
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Param data body models.VmsTrnRequestCanceled true "VmsTrnRequestCanceled data"
// @Router /api/booking-admin/update-canceled [put]
func (h *BookingAdminHandler) UpdateCanceled(c *gin.Context) {
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
	if err := query.First(&trnRequest, "trn_request_uid = ? AND is_deleted = ?  AND ref_request_status_code='30'", request.TrnRequestUID, "0").Error; err != nil {
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
			"admin-approval",
			request.CanceledRequestReason,
		)
	}
	c.JSON(http.StatusOK, gin.H{"message": "Updated successfully", "result": result})
}

// UpdateVehicleUser godoc
// @Summary Update vehicle information for a booking user
// @Description This endpoint allows a booking user to update the vehicle details associated with their request.
// @Tags Booking-admin
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Param data body models.VmsTrnRequestVehicleUser true "VmsTrnRequestVehicleUser data"
// @Router /api/booking-admin/update-vehicle-user [put]
func (h *BookingAdminHandler) UpdateVehicleUser(c *gin.Context) {
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
	request.VehicleUserDeptNameShort = funcs.GetDeptSAPShort(vehicleUser.DeptSAP)
	request.VehicleUserDeptNameFull = funcs.GetDeptSAPFull(vehicleUser.DeptSAP)
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
// @Tags Booking-admin
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Param data body models.VmsTrnRequestTrip true "VmsTrnRequestTrip data"
// @Router /api/booking-admin/update-trip [put]
func (h *BookingAdminHandler) UpdateTrip(c *gin.Context) {
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
// @Tags Booking-admin
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Param data body models.VmsTrnRequestPickup true "VmsTrnRequestPickup data"
// @Router /api/booking-admin/update-pickup [put]
func (h *BookingAdminHandler) UpdatePickup(c *gin.Context) {
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
// @Tags Booking-admin
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Param data body models.VmsTrnRequestDocument true "VmsTrnRequestDocument data"
// @Router /api/booking-admin/update-document [put]
func (h *BookingAdminHandler) UpdateDocument(c *gin.Context) {
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
// @Tags Booking-admin
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Param data body models.VmsTrnRequestCost true "VmsTrnRequestCost data"
// @Router /api/booking-admin/update-cost [put]
func (h *BookingAdminHandler) UpdateCost(c *gin.Context) {
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
// @Tags Booking-admin
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Param data body models.VmsTrnRequestDriver true "VmsTrnRequestDriver data"
// @Router /api/booking-admin/update-driver [put]
func (h *BookingAdminHandler) UpdateDriver(c *gin.Context) {
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
	var driver models.VmsMasDriver
	if err := config.DB.First(&driver, "mas_driver_uid = ? AND is_deleted = '0'", request.MasCarPoolDriverUID).Error; err == nil {
		request.DriverEmpID = driver.DriverID
		request.DriverEmpName = driver.DriverName
		request.DriverEmpDeptSAP = driver.DriverDeptSAP
		request.DriverEmpDeptNameShort = funcs.GetDeptSAPShort(driver.DriverDeptSAP)
		request.DriverEmpDeptNameFull = funcs.GetDeptSAPFull(driver.DriverDeptSAP)
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

// UpdateVehicle godoc
// @Summary Update vehicle for a booking request
// @Description This endpoint allows a booking user to update the cost information for an existing booking request.
// @Tags Booking-admin
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Param data body models.VmsTrnRequestVehicle true "VmsTrnRequestVehicle data"
// @Router /api/booking-admin/update-vehicle [put]
func (h *BookingAdminHandler) UpdateVehicle(c *gin.Context) {
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

	var vehicle models.VmsMasVehicleDepartment
	if err := config.DB.First(&vehicle, "mas_vehicle_uid = ? AND is_deleted = '0'", request.MasVehicleUID).Error; err == nil {
		request.MasVehicleDepartmentUID = vehicle.MasVehicleDepartmentUID
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
