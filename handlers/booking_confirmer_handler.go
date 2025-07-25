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

type BookingConfirmerHandler struct {
	Role string
}

var MenuNameMapConfirmer = map[string]string{
	"20,21,30": "คำขอใช้ยานพาหนะ",
	"00":       "คำขออนุมัติทำหน้าที่ขับรถยนต์",
}

var StatusNameMapConfirmer = map[string]string{
	"20": "รออนุมัติ",
	"21": "ตีกลับ",
	"30": "อนุมัติแล้ว",
	"90": "ยกเลิกคำขอ",
}

func (h *BookingConfirmerHandler) SetQueryRole(user *models.AuthenUserEmp, query *gorm.DB) *gorm.DB {
	return query.Where("confirmed_request_emp_id = ? ", user.EmpID)
}

func (h *BookingConfirmerHandler) SetQueryStatusCanUpdate(query *gorm.DB) *gorm.DB {
	return query.Where("ref_request_status_code in ('20') and is_deleted = '0'")
}

// MenuRequests godoc
// @Summary Summary booking requests by request status code
// @Description Summary booking requests, counts grouped by request status code
// @Tags Booking-confirmer
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Router /api/booking-confirmer/menu-requests [get]
func (h *BookingConfirmerHandler) MenuRequests(c *gin.Context) {
	user := funcs.GetAuthenUser(c, h.Role)
	if c.IsAborted() {
		return
	}
	statusMenuMap := MenuNameMapConfirmer
	query := h.SetQueryRole(user, config.DB)
	summary, err := funcs.MenuRequests(statusMenuMap, query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "message": messages.ErrInternalServer.Error()})
		return
	}
	for i := range summary {
		if summary[i].RefRequestStatusCode == "00" {
			//get count from vms_trn_request_annual_driver
			var count int64
			query := h.SetQueryRole(user, config.DB)
			query = query.Table("vms_trn_request_annual_driver").Where("is_deleted = ?", "0")
			if err := query.Count(&count).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "message": messages.ErrInternalServer.Error()})
				return
			}
			summary[i].Count = int(count)
			summary[i].RefRequestStatusName = "คำขออนุมัติทำหน้าที่ขับรถยนต์"
		}
	}
	sort.Slice(summary, func(i, j int) bool {
		return summary[i].RefRequestStatusCode > summary[j].RefRequestStatusCode
	})
	c.JSON(http.StatusOK, summary)
}

// SearchRequests godoc
// @Summary Search booking requests and get summary counts by request status code
// @Description Search for requests using a keyword and get the summary of counts grouped by request status code
// @Tags Booking-confirmer
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
// @Router /api/booking-confirmer/search-requests [get]
func (h *BookingConfirmerHandler) SearchRequests(c *gin.Context) {
	user := funcs.GetAuthenUser(c, h.Role)
	if c.IsAborted() {
		return
	}

	var requests []models.VmsTrnRequestList
	var summary []models.VmsTrnRequestSummary

	statusNameMap := StatusNameMapConfirmer
	statusCodes := make([]string, 0, len(statusNameMap))
	for code := range statusNameMap {
		statusCodes = append(statusCodes, code)
	}

	// Build the main query
	query := h.SetQueryRole(user, config.DB)
	query = query.Table("public.vms_trn_request AS req").
		Select(`req.*, v.vehicle_license_plate,v.vehicle_license_plate_province_short,v.vehicle_license_plate_province_full,
		fn_get_long_short_dept_name_by_dept_sap(d.vehicle_owner_dept_sap) vehicle_department_dept_sap_short,       
		(select max(mc.carpool_name) from vms_mas_carpool mc where mc.mas_carpool_uid=req.mas_carpool_uid) vehicle_carpool_name
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
// @Tags Booking-confirmer
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Param trn_request_uid path string true "TrnRequestUID (trn_request_uid)"
// @Router /api/booking-confirmer/request/{trn_request_uid} [get]
func (h *BookingConfirmerHandler) GetRequest(c *gin.Context) {
	user := funcs.GetAuthenUser(c, h.Role)
	if c.IsAborted() {
		return
	}
	request, err := funcs.GetRequest(c, StatusNameMapConfirmer)
	if err != nil {
		return
	}
	if request.RefRequestStatusCode == "20" {
		request.ProgressRequestStatus = []models.ProgressRequestStatus{
			{ProgressIcon: "1", ProgressName: "รออนุมัติจากต้นสังกัด"},
			{ProgressIcon: "0", ProgressName: "รอผู้ดูแลยานพาหนะตรวจสอบ"},
			{ProgressIcon: "0", ProgressName: "รออนุมัติให้ใช้ยานพาหนะ"},
		}
		empUser := funcs.GetUserEmpInfo(user.EmpID)
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
	}

	if request.RefRequestStatusCode == "21" {
		request.ProgressRequestStatus = []models.ProgressRequestStatus{
			{ProgressIcon: "2", ProgressName: "ถูกตีกลับจากต้นสังกัด"},
			{ProgressIcon: "0", ProgressName: "รอผู้ดูแลยานพาหนะตรวจสอบ"},
			{ProgressIcon: "0", ProgressName: "รออนุมัติให้ใช้ยานพาหนะ"},
		}
		request.ProgressRequestStatusEmp = funcs.GetProgressRequestStatusEmp(request.TrnRequestUID, "21", "ผู้อนุมัติต้นสังกัด")
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

// UpdateRejected godoc
// @Summary Update sended back status for an item
// @Description This endpoint allows users to update the sended back status of an item.
// @Tags Booking-confirmer
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Param data body models.VmsTrnRequestRejected true "VmsTrnRequestRejected data"
// @Router /api/booking-confirmer/update-rejected [put]
func (h *BookingConfirmerHandler) UpdateRejected(c *gin.Context) {
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
	request.RefRequestStatusCode = "21"
	request.UpdatedAt = time.Now()
	request.UpdatedBy = user.EmpID

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

	if err := config.DB.Save(&request).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to update : %v", err), "message": messages.ErrInternalServer.Error()})
		return
	}
	if err := config.DB.First(&result, "trn_request_uid = ?", request.TrnRequestUID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Booking not found"})
		return
	}
	funcs.CreateTrnRequestActionLog(request.TrnRequestUID,
		request.RefRequestStatusCode,
		"ถูกตีกลับ จากต้นสังกัด",
		user.EmpID,
		"level1-approval",
		request.RejectedRequestReason,
	)

	c.JSON(http.StatusOK, gin.H{"message": "Updated successfully", "result": result})
}

// UpdateApproved godoc
// @Summary Update sended back status for an item
// @Description This endpoint allows users to update the sended back status of an item.
// @Tags Booking-confirmer
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Param data body models.VmsTrnRequestConfirmed true "VmsTrnRequestConfirmed data"
// @Router /api/booking-confirmer/update-approved [put]
func (h *BookingConfirmerHandler) UpdateApproved(c *gin.Context) {
	user := funcs.GetAuthenUser(c, h.Role)
	if c.IsAborted() {
		return
	}
	var request, trnRequest models.VmsTrnRequestConfirmed
	var result struct {
		models.VmsTrnRequestConfirmed
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

	request.RefRequestStatusCode = "30" // ยืนยันคำขอแล้ว รอตรวจสอบคำขอ
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
		"รอผู้ดูแลยานพาหนะตรวจสอบ",
		user.EmpID,
		"level1-approval",
		"",
	)
	funcs.CheckMustPassStatus(request.TrnRequestUID)
	c.JSON(http.StatusOK, gin.H{"message": "Updated successfully", "result": result})
}

// UpdateCanceled godoc
// @Summary Update cancel status for an item
// @Description This endpoint allows users to update the cancel status of an item.
// @Tags Booking-confirmer
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Param data body models.VmsTrnRequestCanceled true "VmsTrnRequestCanceled data"
// @Router /api/booking-confirmer/update-canceled [put]
func (h *BookingConfirmerHandler) UpdateCanceled(c *gin.Context) {
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

// ExportRequests godoc
// @Summary Export booking requests
// @Description Export booking requests by criteria
// @Tags Booking-confirmer
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
// @Router /api/booking-confirmer/export-requests [get]
func (h *BookingConfirmerHandler) ExportRequests(c *gin.Context) {
	user := funcs.GetAuthenUser(c, h.Role)
	if c.IsAborted() {
		return
	}
	query := h.SetQueryRole(user, config.DB)
	funcs.ExportRequests(c, user, query, StatusNameMapConfirmer)
}
