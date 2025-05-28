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

type DriverLicenseApproverHandler struct {
	Role string
}

var LicenseStatusNameMapApprover = map[string]string{
	"10": "รอตรวจสอบ",
	"11": "ตีกลับคำขอ",
	"20": "รออนุมัติ",
	"30": "อนุมัติ",
	"90": "ยกเลิกคำขอ",
}

func (h *DriverLicenseApproverHandler) SetQueryRole(user *models.AuthenUserEmp, query *gorm.DB) *gorm.DB {
	if user.EmpID == "" {
		return query
	}
	return query
}

func (h *DriverLicenseApproverHandler) SetQueryRoleDept(user *models.AuthenUserEmp, query *gorm.DB) *gorm.DB {
	if user.EmpID == "" {
		return query
	}
	return query
}
func (h *DriverLicenseApproverHandler) SetQueryStatusCanUpdate(query *gorm.DB) *gorm.DB {
	return query.Where("ref_request_annual_driver_status_code in ('20') and is_deleted = '0'")
}

// SearchRequests godoc
// @Summary Search driver license annual requests and get summary counts by request status code
// @Description Search for annual driver license requests using a keyword and get the summary of counts grouped by request status code
// @Tags Driver-license-approver
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Param search query string false "Search keyword (matches request_annual_driver_no, created_request_emp_name)"
// @Param ref_request_annual_driver_status_code query string false "Filter by multiple request status codes (comma-separated, e.g., '10,11,20')"
// @Param ref_driver_license_type_code query string false "Filter by multiple license type codes (comma-separated, e.g., '1,2,3')"
// @Param annual_yyyy query string false "Filter by annual yyyy"
// @Param start_created_request_datetime query string false "Filter by start created datetime (YYYY-MM-DD format)"
// @Param end_created_request_datetime query string false "Filter by end created datetime (YYYY-MM-DD format)"
// @Param start_driver_license_expire_date query string false "Filter by start license expire datetime (YYYY-MM-DD format)"
// @Param end_driver_license_expire_date query string false "Filter by end license expire datetime (YYYY-MM-DD format)"
// @Param order_by query string false "Order by request_annual_driver_no, driver_license_expire_date"
// @Param order_dir query string false "Order direction: asc or desc"
// @Param page query int false "Page number (default: 1)"
// @Param limit query int false "Number of records per page (default: 10)"
// @Router /api/driver-license-approver/search-requests [get]
func (h *DriverLicenseApproverHandler) SearchRequests(c *gin.Context) {
	user := funcs.GetAuthenUser(c, h.Role)
	if c.IsAborted() {
		return
	}
	statusNameMap := LicenseStatusNameMapApprover

	var requests []models.VmsDriverLicenseAnnualList
	var summary []models.VmsTrnRequestAnnualDriverSummary

	// Use the keys from statusNameMap as the list of valid status codes
	statusCodes := make([]string, 0, len(statusNameMap))
	for code := range statusNameMap {
		statusCodes = append(statusCodes, code)
	}

	// Build the main query
	query := h.SetQueryRole(user, config.DB)
	query = query.Table("public.vms_trn_request_annual_driver AS req").
		Select("req.*, rcode.ref_driver_license_type_name").
		Joins("LEFT JOIN public.vms_ref_driver_license_type AS rcode ON req.ref_driver_license_type_code = rcode.ref_driver_license_type_code").
		Where("req.ref_request_annual_driver_status_code IN (?)", statusCodes)

	// Apply additional filters (search, date range, etc.)
	if search := c.Query("search"); search != "" {
		query = query.Where("req.request_annual_driver_no ILIKE ? OR req.created_request_emp_name ILIKE ?", "%"+search+"%", "%"+search+"%")
	}
	if startDate := c.Query("start_created_request_datetime"); startDate != "" {
		query = query.Where("req.created_request_datetime >= ?", startDate)
	}
	if endDate := c.Query("end_created_request_datetime"); endDate != "" {
		query = query.Where("req.created_request_datetime <= ?", endDate)
	}
	//ref_driver_license_type_code
	if refDriverLicenseTypeCode := c.Query("ref_driver_license_type_code"); refDriverLicenseTypeCode != "" {
		codes := strings.Split(refDriverLicenseTypeCode, ",")
		query = query.Where("req.ref_driver_license_type_code IN (?)", codes)
	}

	if refRequestStatusCodes := c.Query("ref_request_annual_driver_status_code"); refRequestStatusCodes != "" {
		// Split the comma-separated codes into a slice
		codes := strings.Split(refRequestStatusCodes, ",")
		query = query.Where("req.ref_request_annual_driver_status_code  IN (?)", codes)
	}
	// Ordering
	orderBy := c.Query("order_by")
	orderDir := c.Query("order_dir")
	if orderDir != "desc" {
		orderDir = "asc"
	}
	switch orderBy {
	case "request_annual_driver_no":
		query = query.Order("req.request_annual_driver_no " + orderDir)
	case "driver_license_expire_date":
		query = query.Order("req.driver_license_expire_date " + orderDir)
	case "created_request_datetime":
		query = query.Order("req.created_request_datetime " + orderDir)
	default:
		query = query.Order("req.created_request_datetime desc")
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
	if err := query.
		Scan(&requests).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "message": messages.ErrInternalServer.Error()})
		return
	}
	for i := range requests {
		requests[i].RefRequestAnnualDriverStatusName = statusNameMap[requests[i].RefRequestAnnualDriverStatusCode]
	}

	// Build the summary query
	summaryQuery := h.SetQueryRole(user, config.DB)
	summaryQuery = summaryQuery.Table("public.vms_trn_request_annual_driver AS req").
		Select("req.ref_request_annual_driver_status_code, COUNT(*) as count").
		Where("req.ref_request_annual_driver_status_code IN (?)", statusCodes).
		Group("req.ref_request_annual_driver_status_code")

	// Execute the summary query
	dbSummary := []struct {
		RefRequestAnnualDriverStatusCode string `gorm:"column:ref_request_annual_driver_status_code"`
		Count                            int    `gorm:"column:count"`
	}{}
	if err := summaryQuery.Scan(&dbSummary).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "message": messages.ErrInternalServer.Error()})
		return
	}

	// Create a complete summary with all statuses from statusNameMap
	for code, name := range statusNameMap {
		count := 0
		for _, dbItem := range dbSummary {
			if dbItem.RefRequestAnnualDriverStatusCode == code {
				count = dbItem.Count
				break
			}
		}
		summary = append(summary, models.VmsTrnRequestAnnualDriverSummary{
			RefRequestAnnualDriverStatusCode: code,
			RefRequestAnnualDriverStatusName: name,
			Count:                            count,
		})
	}
	// Sort the summary by RefRequestStatusCode
	sort.Slice(summary, func(i, j int) bool {
		return summary[i].RefRequestAnnualDriverStatusCode < summary[j].RefRequestAnnualDriverStatusCode
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

// GetDriverLicenseAnnual godoc
// @Summary Retrieve a specific driver license annual record
// @Description Get the details of a driver license annual record by its unique identifier
// @Tags Driver-license-approver
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Param trn_request_annual_driver_uid path string true "trnRequestAnnualDriverUID (trn_request_annual_driver_uid)"
// @Router /api/driver-license-approver/license-annual/{trn_request_annual_driver_uid} [get]
func (h *DriverLicenseApproverHandler) GetDriverLicenseAnnual(c *gin.Context) {
	user := funcs.GetAuthenUser(c, h.Role)
	if c.IsAborted() {
		return
	}
	trnRequestAnnualDriverUID := c.Param("trn_request_annual_driver_uid")
	var request models.VmsDriverLicenseAnnualResponse
	query := h.SetQueryRole(user, config.DB)
	if err := query.
		Preload("DriverLicenseType").
		Preload("DriverCertificateType").
		First(&request, "trn_request_annual_driver_uid = ? and is_deleted = ?", trnRequestAnnualDriverUID, "0").Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "annual not found", "message": messages.ErrNotfound.Error()})
		return
	}
	if request.RefRequestAnnualDriverStatusCode == "10" {
		request.ProgressRequestStatus = []models.ProgressRequestStatus{
			{ProgressIcon: "3", ProgressName: "ขออนุมัติ", ProgressDatetime: request.CreatedRequestDatetime},
			{ProgressIcon: "1", ProgressName: "รอต้นสังกัดตรวจสอบ", ProgressDatetime: request.ConfirmedRequestDatetime},
			{ProgressIcon: "0", ProgressName: "รออนุมัติให้ทำหน้าที่ขับรถประจำปี", ProgressDatetime: request.ApprovedRequestDatetime},
		}
		request.ProgressRequestStatusEmp = models.ProgressRequestStatusEmp{
			ActionRole:   "ผู้ขออนุมัติ",
			EmpID:        request.CreatedRequestEmpID,
			EmpName:      request.CreatedRequestEmpName,
			EmpPosition:  request.CreatedRequestEmpPosition,
			DeptSAP:      request.CreatedRequestDeptSap,
			DeptSAPShort: request.CreatedRequestDeptSapNameShort,
			DeptSAPFull:  request.CreatedRequestDeptSapNameFull,
			PhoneNumber:  request.CreatedRequestPhoneNumber,
			MobileNumber: request.CreatedRequestMobileNumber,
		}
	}
	if request.RefRequestAnnualDriverStatusCode == "11" {
		request.ProgressRequestStatus = []models.ProgressRequestStatus{
			{ProgressIcon: "3", ProgressName: "ขออนุมัติ", ProgressDatetime: request.CreatedRequestDatetime},
			{ProgressIcon: "2", ProgressName: "ตีกลับจากต้นสังกัด", ProgressDatetime: request.ConfirmedRequestDatetime},
			{ProgressIcon: "0", ProgressName: "รออนุมัติให้ทำหน้าที่ขับรถประจำปี", ProgressDatetime: request.ApprovedRequestDatetime},
		}
		request.ProgressRequestStatusEmp = models.ProgressRequestStatusEmp{
			ActionRole:   "ผู้อนุมัติต้นสังกัด",
			EmpID:        request.ConfirmedRequestEmpID,
			EmpName:      request.ConfirmedRequestEmpName,
			EmpPosition:  request.ConfirmedRequestEmpPosition,
			DeptSAP:      request.ConfirmedRequestDeptSap,
			DeptSAPShort: request.ConfirmedRequestDeptSapShort,
			DeptSAPFull:  request.ConfirmedRequestDeptSapFull,
			PhoneNumber:  request.ConfirmedRequestPhoneNumber,
			MobileNumber: request.ConfirmedRequestMobileNumber,
		}
	}
	if request.RefRequestAnnualDriverStatusCode == "20" {
		request.ProgressRequestStatus = []models.ProgressRequestStatus{
			{ProgressIcon: "3", ProgressName: "ขออนุมัติ", ProgressDatetime: request.CreatedRequestDatetime},
			{ProgressIcon: "3", ProgressName: "ต้นสังกัดตรวจสอบ", ProgressDatetime: request.ConfirmedRequestDatetime},
			{ProgressIcon: "1", ProgressName: "รออนุมัติให้ทำหน้าที่ขับรถประจำปี", ProgressDatetime: request.ApprovedRequestDatetime},
		}
		request.ProgressRequestStatusEmp = models.ProgressRequestStatusEmp{
			ActionRole:   "ผู้อนุมัติต้นสังกัด",
			EmpID:        request.ConfirmedRequestEmpID,
			EmpName:      request.ConfirmedRequestEmpName,
			EmpPosition:  request.ConfirmedRequestEmpPosition,
			DeptSAP:      request.ConfirmedRequestDeptSap,
			DeptSAPShort: request.ConfirmedRequestDeptSapShort,
			DeptSAPFull:  request.ConfirmedRequestDeptSapFull,
			PhoneNumber:  request.ConfirmedRequestPhoneNumber,
			MobileNumber: request.ConfirmedRequestMobileNumber,
		}
	}
	if request.RefRequestAnnualDriverStatusCode == "21" {
		request.ProgressRequestStatus = []models.ProgressRequestStatus{
			{ProgressIcon: "3", ProgressName: "ขออนุมัติ", ProgressDatetime: request.CreatedRequestDatetime},
			{ProgressIcon: "3", ProgressName: "ต้นสังกัดตรวจสอบ", ProgressDatetime: request.ConfirmedRequestDatetime},
			{ProgressIcon: "2", ProgressName: "ตีกลับจากผู้อนุมัติ", ProgressDatetime: request.RejectedRequestDatetime},
		}
		request.ProgressRequestStatusEmp = models.ProgressRequestStatusEmp{
			ActionRole:   "ผู้อนุมัติต้นสังกัด",
			EmpID:        request.ConfirmedRequestEmpID,
			EmpName:      request.ConfirmedRequestEmpName,
			EmpPosition:  request.ConfirmedRequestEmpPosition,
			DeptSAP:      request.ConfirmedRequestDeptSap,
			DeptSAPShort: request.ConfirmedRequestDeptSapShort,
			DeptSAPFull:  request.ConfirmedRequestDeptSapFull,
			PhoneNumber:  request.ConfirmedRequestPhoneNumber,
			MobileNumber: request.ConfirmedRequestMobileNumber,
		}
	}
	if request.RefRequestAnnualDriverStatusCode == "30" {
		request.ProgressRequestStatus = []models.ProgressRequestStatus{
			{ProgressIcon: "3", ProgressName: "ขออนุมัติ", ProgressDatetime: request.CreatedRequestDatetime},
			{ProgressIcon: "3", ProgressName: "ต้นสังกัดตรวจสอบ", ProgressDatetime: request.ConfirmedRequestDatetime},
			{ProgressIcon: "3", ProgressName: "อนุมัติให้ทำหน้าที่ขับรถประจำปี", ProgressDatetime: request.ApprovedRequestDatetime},
		}
		request.ProgressRequestStatusEmp = models.ProgressRequestStatusEmp{
			ActionRole:   "ผู้อนุมัติให้ทำหน้าที่ขับรถประจำปี",
			EmpID:        request.ApprovedRequestEmpID,
			EmpName:      request.ApprovedRequestEmpName,
			EmpPosition:  request.ApprovedRequestEmpPosition,
			DeptSAP:      request.ApprovedRequestDeptSap,
			DeptSAPShort: request.ApprovedRequestDeptSapShort,
			DeptSAPFull:  request.ApprovedRequestDeptSapFull,
			PhoneNumber:  request.ApprovedRequestPhoneNumber,
			MobileNumber: request.ApprovedRequestMobileNumber,
		}
	}
	if request.RefRequestAnnualDriverStatusCode == "90" && request.CanceledRequestEmpID == request.CreatedRequestEmpID {
		request.ProgressRequestStatus = []models.ProgressRequestStatus{
			{ProgressIcon: "2", ProgressName: "ยกเลิก", ProgressDatetime: request.CanceledRequestDatetime},
		}
		request.ProgressRequestStatusEmp = models.ProgressRequestStatusEmp{
			ActionRole:   "ผู้ขออนุมัติ",
			EmpID:        request.CreatedRequestEmpID,
			EmpName:      request.CreatedRequestEmpName,
			EmpPosition:  request.CreatedRequestEmpPosition,
			DeptSAP:      request.CreatedRequestDeptSap,
			DeptSAPShort: request.CreatedRequestDeptSapNameShort,
			DeptSAPFull:  request.CreatedRequestDeptSapNameFull,
			PhoneNumber:  request.CreatedRequestPhoneNumber,
			MobileNumber: request.CreatedRequestMobileNumber,
		}
	}
	if request.RefRequestAnnualDriverStatusCode == "90" && request.CanceledRequestEmpID == request.ConfirmedRequestEmpID {

		request.ProgressRequestStatus = []models.ProgressRequestStatus{
			{ProgressIcon: "2", ProgressName: "ยกเลิกจากต้นสังกัด", ProgressDatetime: request.CanceledRequestDatetime},
		}
		request.ProgressRequestStatusEmp = models.ProgressRequestStatusEmp{
			ActionRole:   "ผู้อนุมัติต้นสังกัด",
			EmpID:        request.ConfirmedRequestEmpID,
			EmpName:      request.ConfirmedRequestEmpName,
			EmpPosition:  request.ConfirmedRequestEmpPosition,
			DeptSAP:      request.ConfirmedRequestDeptSap,
			DeptSAPShort: request.ConfirmedRequestDeptSapShort,
			DeptSAPFull:  request.ConfirmedRequestDeptSapFull,
			PhoneNumber:  request.ConfirmedRequestPhoneNumber,
			MobileNumber: request.ConfirmedRequestMobileNumber,
		}
	}

	if request.RefRequestAnnualDriverStatusCode == "93" && request.CanceledRequestEmpID == request.ApprovedRequestEmpID {
		request.ProgressRequestStatus = []models.ProgressRequestStatus{
			{ProgressIcon: "3", ProgressName: "อนุมัติจากต้นสังกัด", ProgressDatetime: request.ApprovedRequestDatetime},
			{ProgressIcon: "2", ProgressName: "ยกเลิกจากผู้อนุมัติ", ProgressDatetime: request.CanceledRequestDatetime},
		}
		request.ProgressRequestStatusEmp = models.ProgressRequestStatusEmp{
			ActionRole:   "ผู้อนุมัติให้ทำหน้าที่ขับรถประจำปี",
			EmpID:        request.ApprovedRequestEmpID,
			EmpName:      request.ApprovedRequestEmpName,
			EmpPosition:  request.ApprovedRequestEmpPosition,
			DeptSAP:      request.ApprovedRequestDeptSap,
			DeptSAPShort: request.ApprovedRequestDeptSapShort,
			DeptSAPFull:  request.ApprovedRequestDeptSapFull,
			PhoneNumber:  request.ApprovedRequestPhoneNumber,
			MobileNumber: request.ApprovedRequestMobileNumber,
		}
	}
	request.RefRequestAnnualDriverStatusName = LicenseStatusNameMapApprover[request.RefRequestAnnualDriverStatusCode]
	request.CreatedRequestImageUrl = funcs.GetEmpImage(request.CreatedRequestEmpID)
	request.ConfirmedRequestImageUrl = funcs.GetEmpImage(request.ConfirmedRequestEmpID)
	request.ApprovedRequestImageUrl = funcs.GetEmpImage(request.ApprovedRequestEmpID)
	request.ProgressRequestHistory = GetProgressRequestHistory(request)
	c.JSON(http.StatusOK, request)
}

// UpdateDriverLicenseAnnualCanceled godoc
// @Summary Update cancel status for a driver license annual record
// @Description This endpoint allows users to update the cancel status of a driver license annual record.
// @Tags Driver-license-approver
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Param data body models.VmsDriverLicenseAnnualCanceled true "VmsDriverLicenseAnnualCanceled data"
// @Router /api/driver-license-approver/update-license-annual-canceled [put]
func (h *DriverLicenseApproverHandler) UpdateDriverLicenseAnnualCanceled(c *gin.Context) {
	user := funcs.GetAuthenUser(c, h.Role)
	if c.IsAborted() {
		return
	}
	var request, driverLicenseAnnual models.VmsDriverLicenseAnnualCanceled
	var result struct {
		models.VmsDriverLicenseAnnualCanceled
		models.VmsTrnRequestAnnualDriverNo
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "message": messages.ErrInternalServer.Error()})
		return
	}
	query := h.SetQueryRole(user, config.DB)
	query = h.SetQueryStatusCanUpdate(query)
	if err := query.First(&driverLicenseAnnual, "trn_request_annual_driver_uid = ? AND is_deleted = ?", request.TrnRequestAnnualDriverUID, "0").Error; err != nil {
		c.JSON(http.StatusMethodNotAllowed, gin.H{"error": "Driver license annual can not update", "message": messages.ErrAnnualCannotUpdate.Error()})
		return
	}
	request.RefRequestAnnualDriverStatusCode = "90"
	request.UpdatedAt = time.Now()
	request.UpdatedBy = user.EmpID

	empUser := funcs.GetUserEmpInfo(user.EmpID)
	request.CanceledRequestEmpID = empUser.EmpID
	request.CanceledRequestEmpName = empUser.FullName
	request.CanceledRequestDeptSAP = empUser.DeptSAP
	request.CanceledRequestDeptSAPShort = empUser.DeptSAPShort
	request.CanceledRequestDeptSAPFull = empUser.DeptSAPFull
	request.CanceledRequestDatetime = time.Now()

	if err := config.DB.Save(&request).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to update: %v", err), "message": messages.ErrInternalServer.Error()})
		return
	}

	if err := config.DB.First(&result, "trn_request_annual_driver_uid = ? AND is_deleted = ?", request.TrnRequestAnnualDriverUID, "0").Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Driver license annual record not found", "message": messages.ErrNotfound.Error()})
		return
	}
	funcs.CreateRequestAnnualLicenseNotification(request.TrnRequestAnnualDriverUID)
	c.JSON(http.StatusOK, gin.H{"message": "Updated successfully", "result": result})
}

// UpdateDriverLicenseAnnualRejected godoc
// @Summary Update reject status for a driver license annual record
// @Description This endpoint allows users to update the reject status of a driver license annual record.
// @Tags Driver-license-approver
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Param data body models.VmsDriverLicenseAnnualRejected true "VmsDriverLicenseAnnualRejected data"
// @Router /api/driver-license-approver/update-license-annual-rejected [put]
func (h *DriverLicenseApproverHandler) UpdateDriverLicenseAnnualRejected(c *gin.Context) {
	user := funcs.GetAuthenUser(c, h.Role)
	if c.IsAborted() {
		return
	}
	var request, driverLicenseAnnual models.VmsDriverLicenseAnnualRejected
	var result struct {
		models.VmsDriverLicenseAnnualRejected
		models.VmsTrnRequestAnnualDriverNo
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "message": messages.ErrInternalServer.Error()})
		return
	}
	query := h.SetQueryRole(user, config.DB)
	query = h.SetQueryStatusCanUpdate(query)
	if err := query.First(&driverLicenseAnnual, "trn_request_annual_driver_uid = ? AND is_deleted = ?", request.TrnRequestAnnualDriverUID, "0").Error; err != nil {
		c.JSON(http.StatusMethodNotAllowed, gin.H{"error": "Driver license annual can not update", "message": messages.ErrAnnualCannotUpdate.Error()})
		return
	}
	request.RefRequestAnnualDriverStatusCode = "21"
	request.UpdatedAt = time.Now()
	request.UpdatedBy = user.EmpID

	empUser := funcs.GetUserEmpInfo(user.EmpID)
	request.RejectedRequestEmpID = empUser.EmpID
	request.RejectedRequestEmpName = empUser.FullName
	request.RejectedRequestDeptSAP = empUser.DeptSAP
	request.RejectedRequestDeptSAPShort = empUser.DeptSAPShort
	request.RejectedRequestDeptSAPFull = empUser.DeptSAPFull
	request.RejectedRequestDatetime = time.Now()

	if err := config.DB.Save(&request).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to update: %v", err), "message": messages.ErrInternalServer.Error()})
		return
	}

	if err := config.DB.First(&result, "trn_request_annual_driver_uid = ? AND is_deleted = ?", request.TrnRequestAnnualDriverUID, "0").Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Driver license annual record not found", "message": messages.ErrNotfound.Error()})
		return
	}
	funcs.CreateRequestAnnualLicenseNotification(request.TrnRequestAnnualDriverUID)
	c.JSON(http.StatusOK, gin.H{"message": "Updated successfully", "result": result})
}

// UpdateDriverLicenseAnnualApproved godoc
// @Summary Update approved status for a driver license annual record
// @Description This endpoint allows users to update the approved status of a driver license annual record.
// @Tags Driver-license-approver
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Param data body models.VmsDriverLicenseAnnualApproved true "VmsDriverLicenseAnnualApproved data"
// @Router /api/driver-license-approver/update-license-annual-approved [put]
func (h *DriverLicenseApproverHandler) UpdateDriverLicenseAnnualApproved(c *gin.Context) {
	user := funcs.GetAuthenUser(c, h.Role)
	if c.IsAborted() {
		return
	}
	var request, driverLicenseAnnual models.VmsDriverLicenseAnnualApproved
	var result struct {
		models.VmsDriverLicenseAnnualApproved
		models.VmsTrnRequestAnnualDriverNo
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "message": messages.ErrInternalServer.Error()})
		return
	}

	query := h.SetQueryRole(user, config.DB)
	query = h.SetQueryStatusCanUpdate(query)
	if err := query.First(&driverLicenseAnnual, "trn_request_annual_driver_uid = ? AND is_deleted = ?", request.TrnRequestAnnualDriverUID, "0").Error; err != nil {
		c.JSON(http.StatusMethodNotAllowed, gin.H{"error": "Driver license annual can not update", "message": messages.ErrAnnualCannotUpdate.Error()})
		return
	}
	request.RefRequestAnnualDriverStatusCode = "30"
	request.UpdatedAt = time.Now()
	request.UpdatedBy = user.EmpID

	empUser := funcs.GetUserEmpInfo(user.EmpID)
	request.ApprovedRequestEmpID = empUser.EmpID
	request.ApprovedRequestEmpName = empUser.FullName
	request.ApprovedRequestDeptSAP = empUser.DeptSAP
	request.ApprovedRequestDeptSAPShort = empUser.DeptSAPShort
	request.ApprovedRequestDeptSAPFull = empUser.DeptSAPFull
	request.ApprovedRequestDatetime = time.Now()
	request.AnnualYYYY = driverLicenseAnnual.AnnualYYYY
	request.DriverLicenseExpireDate = driverLicenseAnnual.DriverLicenseExpireDate
	request.RequestIssueDate = request.ApprovedRequestDatetime
	annualYearEnd := time.Date(driverLicenseAnnual.AnnualYYYY-543, 12, 31, 23, 59, 59, 0, time.UTC)
	if driverLicenseAnnual.DriverLicenseExpireDate.Before(annualYearEnd) {
		request.RequestExpireDate = driverLicenseAnnual.DriverLicenseExpireDate
	} else {
		request.RequestExpireDate = annualYearEnd
	}

	if err := config.DB.Save(&request).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to update: %v", err), "message": messages.ErrInternalServer.Error()})
		return
	}

	if err := config.DB.First(&result, "trn_request_annual_driver_uid = ? AND is_deleted = ?", request.TrnRequestAnnualDriverUID, "0").Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Driver license annual record not found", "message": messages.ErrNotfound.Error()})
		return
	}
	funcs.CreateRequestAnnualLicenseNotification(request.TrnRequestAnnualDriverUID)
	c.JSON(http.StatusOK, gin.H{"message": "Updated successfully", "result": result})
}
