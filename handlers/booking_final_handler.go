package handlers

import (
	"fmt"
	"net/http"
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

type BookingFinalHandler struct {
	Role string
}

var MenuNameMapFinal = map[string]string{
	"40,41": "ตรวจสอบคำขอ",
	"50,51": "ให้กุญแจ",
	"60":    "เดินทาง",
	"70,71": "ตรวจสอบการคืนยานพาหนะ",
	"80":    "เสร็จสิ้น",
	"90":    "ยกเลิกคำขอ",
}
var StatusNameMapFinal = map[string]string{
	"40": "รออนุมัติ",
	"41": "ตีกลับ",
	"50": "อนุมัติแล้ว",
	"90": "ยกเลิกคำขอ",
}

func (h *BookingFinalHandler) SetQueryRole(user *models.AuthenUserEmp, query *gorm.DB) *gorm.DB {
	query = funcs.SetQueryApproverRole(user, query)
	return query
}
func (h *BookingFinalHandler) SetQueryStatusCanUpdate(query *gorm.DB) *gorm.DB {
	return query.Where("ref_request_status_code in ('40') and is_deleted = '0'")
}

// MenuRequests godoc
// @Summary Summary booking requests by request status code
// @Description Summary booking requests, counts grouped by request status code
// @Tags Booking-final
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Router /api/booking-final/menu-requests [get]
func (h *BookingFinalHandler) MenuRequests(c *gin.Context) {
	// Get authenticated user role if needed
	user := funcs.GetAuthenUser(c, h.Role)
	if c.IsAborted() {
		return
	}
	statusMenuMap := MenuNameMapFinal
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
// @Tags Booking-final
// @Accept  json
// @Produce  json
// @Param search query string false "Search keyword (matches request_no, vehicle_license_plate, vehicle_user_emp_name, or work_place)"
// @Param ref_request_status_code query string false "Filter by multiple request status codes (comma-separated, e.g., 'A,B,C')"
// @Param startdate query string false "Filter by start datetime (YYYY-MM-DD format)"
// @Param enddate query string false "Filter by end datetime (YYYY-MM-DD format)"
// @Param order_by query string false "Order by request_no, start_datetime, ref_request_status_code"
// @Param order_dir query string false "Order direction: asc or desc"
// @Param page query int false "Page number (default: 1)"
// @Param limit query int false "Number of records per page (default: 10)"
// @Router /api/booking-final/search-requests [get]
func (h *BookingFinalHandler) SearchRequests(c *gin.Context) {
	user := funcs.GetAuthenUser(c, h.Role)
	if c.IsAborted() {
		return
	}
	statusNameMap := StatusNameMapFinal
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
		v.vehicle_license_plate,v.vehicle_license_plate_province_short,v.vehicle_license_plate_province_full,v."CarTypeDetail" car_type,
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
		//use switch case
		switch requests[i].TripType {
		case 0:
			requests[i].TripTypeName = "ไป-กลับ"
		case 1:
			requests[i].TripTypeName = "ค้างแรม"
		}

		if requests[i].IsPEAEmployeeDriver != 1 {
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
	summaryQuery = summaryQuery.Order("vms_trn_request.ref_request_status_code")

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
// @Tags Booking-final
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Param trn_request_uid path string true "TrnRequestUID (trn_request_uid)"
// @Router /api/booking-final/request/{trn_request_uid} [get]
func (h *BookingFinalHandler) GetRequest(c *gin.Context) {
	user := funcs.GetAuthenUser(c, h.Role)
	if c.IsAborted() {
		return
	}
	request, err := funcs.GetRequest(c, StatusNameMapFinal)
	if err != nil {
		return
	}
	if request.RefRequestStatusCode == "40" {
		request.ProgressRequestStatus = []models.ProgressRequestStatus{
			{ProgressIcon: "3", ProgressName: "อนุมัติจากต้นสังกัด"},
			{ProgressIcon: "3", ProgressName: "อนุมัติจากผู้ดูแลยานพาหนะ"},
			{ProgressIcon: "1", ProgressName: "รออนุมัติให้ใช้ยานพาหนะ"},
		}
		empUser := funcs.GetUserEmpInfo(user.EmpID)
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
	if request.MasCarpoolUID == nil || *request.MasCarpoolUID == "" {
		request.ProgressRequestStatus = append(request.ProgressRequestStatus[:0], request.ProgressRequestStatus[1:]...)
	}
	c.JSON(http.StatusOK, request)
}

// UpdateRejected godoc
// @Summary Update rejected status for an item
// @Description This endpoint allows users to update the rejected status of an item.
// @Tags Booking-final
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Param data body models.VmsTrnRequestRejected true "VmsTrnRequestRejected data"
// @Router /api/booking-final/update-rejected [put]
func (h *BookingFinalHandler) UpdateRejected(c *gin.Context) {
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

	request.RefRequestStatusCode = "41"
	rejectUser := funcs.GetUserEmpInfo(user.EmpID)
	request.RejectedRequestEmpID = rejectUser.EmpID
	request.RejectedRequestEmpName = rejectUser.FullName
	request.RejectedRequestDeptSAP = rejectUser.DeptSAP
	request.RejectedRequestDeptNameShort = rejectUser.DeptSAPShort
	request.RejectedRequestDeptNameFull = rejectUser.DeptSAPFull
	request.RejectedRequestDeskPhone = rejectUser.TelInternal
	request.RejectedRequestMobilePhone = rejectUser.TelMobile
	request.RejectedRequestPosition = rejectUser.Position
	request.RejectedRequestDatetime = models.TimeWithZone{Time: time.Now()}
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

	funcs.CreateTrnRequestActionLog(request.TrnRequestUID,
		request.RefRequestStatusCode,
		"ถูกตีกลับ จากเจ้าของยานพาหนะ",
		user.EmpID,
		"approval-department",
		request.RejectedRequestReason,
	)

	c.JSON(http.StatusOK, gin.H{"message": "Updated successfully", "result": result})
}

// UpdateApproved godoc
// @Summary Update sended back status for an item
// @Description This endpoint allows users to update the sended back status of an item.
// @Tags Booking-final
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Param data body models.VmsTrnRequestApproved true "VmsTrnRequestApproved data"
// @Router /api/booking-final/update-approved [put]
func (h *BookingFinalHandler) UpdateApproved(c *gin.Context) {
	user := funcs.GetAuthenUser(c, h.Role)
	if c.IsAborted() {
		return
	}
	var request, trnRequest models.VmsTrnRequestApproved
	var result struct {
		models.VmsTrnRequestApproved
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
	request.RefRequestStatusCode = "50"
	empUser := funcs.GetUserEmpInfo(user.EmpID)
	request.ApprovedRequestEmpID = empUser.EmpID
	request.ApprovedRequestEmpName = empUser.FullName
	request.ApprovedRequestDeptSAP = empUser.DeptSAP
	request.ApprovedRequestDeptNameShort = empUser.DeptSAPShort
	request.ApprovedRequestDeptNameFull = empUser.DeptSAPFull
	request.ApprovedRequestDeskPhone = empUser.TelInternal
	request.ApprovedRequestMobilePhone = empUser.TelMobile
	request.ApprovedRequestPosition = empUser.Position
	request.ApprovedRequestDatetime = models.TimeWithZone{Time: time.Now()}
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

	// Check if a record exists in vms_trn_vehicle_key_handover

	var approvedWithReceiveKey models.VmsTrnRequestApprovedWithRecieiveKey
	if err := config.DB.First(&approvedWithReceiveKey, "trn_request_uid = ?", request.TrnRequestUID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Booking not found", "message": messages.ErrBookingNotFound.Error()})
		return
	}
	var keyHandover models.VmsTrnReceivedKeyPEA
	if err := config.DB.First(&keyHandover, "trn_request_uid = ?", request.TrnRequestUID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			// Create a new record if it doesn't exist
			newKeyHandover := models.VmsTrnVehicleKeyHandover{
				HandoverUid:      uuid.New().String(),
				TrnRequestUID:    request.TrnRequestUID,
				AppointmentStart: approvedWithReceiveKey.ReceivedKeyStartDatetime,
				AppointmentEnd:   approvedWithReceiveKey.ReceivedKeyEndDatetime,
				AppointmentPlace: approvedWithReceiveKey.ReceivedKeyPlace,
				ReceiverType:     0,
				CreatedAt:        time.Now(),
				CreatedBy:        user.EmpID,
				UpdatedAt:        time.Now(),
				UpdatedBy:        user.EmpID,
			}

			if err := config.DB.Create(&newKeyHandover).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to create key handover record: %v", err), "message": messages.ErrInternalServer.Error()})
				return
			}
			//update vms_trn_request set appointment_key_handover_place,appointment_key_handover_start_datetime,appointment_key_handover_end_datetime
			if err := config.DB.Table("vms_trn_request").
				Where("trn_request_uid = ?", request.TrnRequestUID).
				Update("appointment_key_handover_place", approvedWithReceiveKey.ReceivedKeyPlace).
				Update("appointment_key_handover_start_datetime", approvedWithReceiveKey.ReceivedKeyStartDatetime).
				Update("appointment_key_handover_end_datetime", approvedWithReceiveKey.ReceivedKeyEndDatetime).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to update : %v", err), "message": messages.ErrInternalServer.Error()})
				return
			}

		} else {
			// Handle other errors
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to query key handover record: %v", err), "message": messages.ErrInternalServer.Error()})
			return
		}
	}

	funcs.UpdateRecievedKeyUser(request.TrnRequestUID)

	var receivedKey models.VmsTrnRequestApprovedWithRecieiveKey
	if err := config.DB.First(&receivedKey, "trn_request_uid = ?", request.TrnRequestUID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Booking not found", "message": messages.ErrBookingNotFound.Error()})
		return
	}
	funcs.CreateTrnRequestActionLog(request.TrnRequestUID,
		request.RefRequestStatusCode,
		funcs.GetDateTime2BuddhistYear(receivedKey.ReceivedKeyStartDatetime.TimeWithZoneToTime(), receivedKey.ReceivedKeyEndDatetime.TimeWithZoneToTime())+" สถานที่ "+receivedKey.ReceivedKeyPlace+" นัดหมายรับกุญแจ",
		user.EmpID,
		"approval-department",
		"",
	)

	c.JSON(http.StatusOK, gin.H{"message": "Updated successfully", "result": result})
}

// UpdateCanceled godoc
// @Summary Update cancel status for an item
// @Description This endpoint allows users to update the cancel status of an item.
// @Tags Booking-final
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Param data body models.VmsTrnRequestCanceled true "VmsTrnRequestCanceled data"
// @Router /api/booking-final/update-canceled [put]
func (h *BookingFinalHandler) UpdateCanceled(c *gin.Context) {
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
	request.RefRequestStatusCode = "90"
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to update : %v", err)})
		return
	}

	if err := config.DB.First(&result, "trn_request_uid = ?", request.TrnRequestUID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Booking not found", "message": messages.ErrBookingNotFound.Error()})
		return
	}

	funcs.CreateTrnRequestActionLog(request.TrnRequestUID,
		request.RefRequestStatusCode,
		"ยกเลิกคำขอ",
		user.EmpID,
		"approval-department",
		request.CanceledRequestReason,
	)

	c.JSON(http.StatusOK, gin.H{"message": "Updated successfully", "result": result})
}

// ExportRequests godoc
// @Summary Export booking requests
// @Description Export booking requests by criteria
// @Tags Booking-final
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Param search query string false "Search keyword (matches request_no, vehicle_license_plate, vehicle_user_emp_name, or work_place)"
// @Param ref_request_status_code query string false "Filter by multiple request status codes (comma-separated, e.g., 'A,B,C')"
// @Param startdate query string false "Filter by start datetime (YYYY-MM-DD format)"
// @Param enddate query string false "Filter by end datetime (YYYY-MM-DD format)"
// @Param order_by query string false "Order by request_no, start_datetime, ref_request_status_code"
// @Param order_dir query string false "Order direction: asc or desc"
// @Router /api/booking-final/export-requests [get]
func (h *BookingFinalHandler) ExportRequests(c *gin.Context) {
	user := funcs.GetAuthenUser(c, h.Role)
	if c.IsAborted() {
		return
	}
	query := h.SetQueryRole(user, config.DB)
	funcs.ExportRequests(c, user, query, StatusNameMapFinal)
}
