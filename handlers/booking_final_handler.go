package handlers

import (
	"fmt"
	"net/http"
	"strings"
	"time"
	"vms_plus_be/config"
	"vms_plus_be/funcs"
	"vms_plus_be/models"

	"github.com/gin-gonic/gin"
)

type BookingFinalHandler struct {
	Role string
}

var MenuNameMapFinal = map[string]string{
	"40,41": "ตรวจสอบคำขอ",
	"50,51": "ให้กุญแจ",
	"60":    "เดินทาง",
	"70,71": "ตรวจสอบยานพาหนะ",
	"80":    "เสร็จสิ้น",
	"90":    "ยกเลิกคำขอ",
}
var StatusNameMapFinal = map[string]string{
	"40": "รออนุมัติ",
	"41": "ตีกลับ",
	"50": "อนุมัติแล้ว",
	"90": "ยกเลิกคำขอ",
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
	funcs.GetAuthenUser(c, h.Role)
	if c.IsAborted() {
		return
	}

	statusMenuMap := MenuNameMapFinal

	summary, err := funcs.MenuRequests(statusMenuMap)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Respond with the grouped summary
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
	funcs.GetAuthenUser(c, h.Role)
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
	query := config.DB.Table("public.vms_trn_request AS req").
		Select(
			`req.*, status.ref_request_status_desc, case req.is_pea_employee_driver when '1' then req.driver_emp_name else (select driver_name from vms_mas_driver d where d.mas_driver_uid::text=req.mas_carpool_driver_uid) end driver_name,
			 case req.is_pea_employee_driver when '1' then req.driver_mas_dept_sap_name_short else (select driver_dept_sap_hire from vms_mas_driver d where d.mas_driver_uid::text=req.mas_carpool_driver_uid) end driver_dept_name,
			(select max(md.dept_short) from vms_mas_vehicle_department mvd,vms_mas_department md where md.dept_sap=mvd.vehicle_owner_dept_sap and mvd.mas_vehicle_uid::text=req.mas_vehicle_uid) vehicle_dept_name,
			(select mc.carpool_name from vms_mas_carpool mc,vms_mas_carpool_vehicle mcv where mc.mas_carpool_uid=mc.mas_carpool_uid and mcv.mas_vehicle_uid::text=req.mas_vehicle_uid) vehicle_carpool_name
			`).
		Joins("LEFT JOIN public.vms_ref_request_status AS status ON req.ref_request_status_code = status.ref_request_status_code").
		Where("req.ref_request_status_code IN (?)", statusCodes)

	// Apply additional filters (search, date range, etc.)
	if search := c.Query("search"); search != "" {
		query = query.Where("req.request_no ILIKE ? OR req.vehicle_license_plate ILIKE ? OR req.vehicle_user_emp_name ILIKE ? OR req.work_place ILIKE ?", "%"+search+"%", "%"+search+"%", "%"+search+"%", "%"+search+"%")
	}
	if startDate := c.Query("startdate"); startDate != "" {
		query = query.Where("req.start_datetime >= ?", startDate)
	}
	if endDate := c.Query("enddate"); endDate != "" {
		query = query.Where("req.start_datetime <= ?", endDate)
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
	summaryQuery := config.DB.Table("public.vms_trn_request AS req").
		Select("req.ref_request_status_code, COUNT(*) as count").
		Where("req.ref_request_status_code IN (?)", statusCodes).
		Group("req.ref_request_status_code")

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
	funcs.GetAuthenUser(c, h.Role)
	if c.IsAborted() {
		return
	}
	request, err := funcs.GetRequest(c, StatusNameMapFinal)
	if err != nil {
		return
	}
	if request.RefRequestStatusCode >= "40" {
		request.ProgressRequestStatus = []models.ProgressRequestStatus{
			{ProgressIcon: "3", ProgressName: "อนุมัติจากต้นสังกัด"},
			{ProgressIcon: "3", ProgressName: "อนุมัติจากผู้ดูแลยานพาหนะ"},
			{ProgressIcon: "1", ProgressName: "รออนุมัติให้ใช้ยานพาหนะ"},
		}
	}
	if request.RefRequestStatusCode >= "41" {
		request.ProgressRequestStatus = []models.ProgressRequestStatus{
			{ProgressIcon: "3", ProgressName: "อนุมัติจากต้นสังกัด"},
			{ProgressIcon: "3", ProgressName: "อนุมัติจากผู้ดูแลยานพาหนะ"},
			{ProgressIcon: "2", ProgressName: "ถูกตีกลับจากเจ้าของยานพาหนะ"},
		}
	}
	if request.RefRequestStatusCode >= "50" && request.RefRequestStatusCode < "90" { //
		request.ProgressRequestStatus = []models.ProgressRequestStatus{
			{ProgressIcon: "3", ProgressName: "อนุมัติจากต้นสังกัด"},
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

// UpdateSendedBack godoc
// @Summary Update sended back status for an item
// @Description This endpoint allows users to update the sended back status of an item.
// @Tags Booking-final
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Param data body models.VmsTrnRequestSendedBack true "VmsTrnRequestSendedBack data"
// @Router /api/booking-final/update-sended-back [put]
func (h *BookingFinalHandler) UpdateSendedBack(c *gin.Context) {
	user := funcs.GetAuthenUser(c, h.Role)
	if c.IsAborted() {
		return
	}
	var request, trnRequest models.VmsTrnRequestSendedBack
	var result struct {
		models.VmsTrnRequestSendedBack
		models.VmsTrnRequestRequestNo
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := config.DB.First(&trnRequest, "trn_request_uid = ?", request.TrnRequestUID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Booking not found"})
		return
	}

	empUser := funcs.GetUserEmpInfo(user.EmpID)
	request.RefRequestStatusCode = "41"
	request.SendedBackRequestEmpID = empUser.EmpID
	request.SendedBackRequestEmpName = empUser.FullName
	request.SendedBackRequestDeptSAP = empUser.DeptSAP
	request.SendedBackRequestDeptSAPShort = empUser.DeptSAPShort
	request.SendedBackRequestDeptSAPFull = empUser.DeptSAPFull
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
	funcs.CreateTrnLog(result.TrnRequestUID,
		result.RefRequestStatusCode,
		result.SendedBackRequestReason,
		user.EmpID)

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
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := config.DB.First(&trnRequest, "trn_request_uid = ?", request.TrnRequestUID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Booking not found"})
		return
	}

	empUser := funcs.GetUserEmpInfo(user.EmpID)
	request.RefRequestStatusCode = "50" //
	request.ApprovedRequestEmpID = empUser.EmpID
	request.ApprovedRequestEmpName = empUser.FullName
	request.ApprovedRequestDeptSAP = empUser.DeptSAP
	request.ApprovedRequestDeptSAPShort = empUser.DeptSAPShort
	request.ApprovedRequestDeptSAPFull = empUser.DeptSAPFull
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
	funcs.CreateTrnLog(result.TrnRequestUID,
		result.RefRequestStatusCode,
		"ผู้ดูแลยืนยันคำขอแล้ว",
		user.EmpID)

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

	if err := config.DB.First(&trnRequest, "trn_request_uid = ?", request.TrnRequestUID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Booking not found"})
		return
	}
	empUser := funcs.GetUserEmpInfo(user.EmpID)
	request.RefRequestStatusCode = "90"
	request.UpdatedAt = time.Now()
	request.UpdatedBy = user.EmpID
	request.CanceledRequestEmpID = empUser.EmpID
	request.CanceledRequestEmpName = empUser.FullName
	request.CanceledRequestDeptSAP = empUser.DeptSAP
	request.CanceledRequestDeptSAPShort = empUser.DeptSAPShort
	request.CanceledRequestDeptSAPFull = empUser.DeptSAPFull
	request.CanceledRequestDatetime = time.Now()
	if err := config.DB.Save(&request).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to update : %v", err)})
		return
	}

	if err := config.DB.First(&result, "trn_request_uid = ?", request.TrnRequestUID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Booking not found"})
		return
	}
	funcs.CreateTrnLog(result.TrnRequestUID,
		result.RefRequestStatusCode,
		result.CanceledRequestReason,
		user.EmpID)

	c.JSON(http.StatusOK, gin.H{"message": "Updated successfully", "result": result})
}
