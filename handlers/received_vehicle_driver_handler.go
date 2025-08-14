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

type ReceivedVehicleDriverHandler struct {
	Role string
}

func (h *ReceivedVehicleDriverHandler) SetQueryRole(user *models.AuthenUserEmp, query *gorm.DB) *gorm.DB {
	return query.Where("driver_emp_id = ?", user.EmpID)
}
func (h *ReceivedVehicleDriverHandler) SetQueryStatusCanUpdate(query *gorm.DB) *gorm.DB {
	return query.Where("ref_request_status_code in ('51') and is_deleted = '0'")
}

// SearchRequests godoc
// @Summary Search booking requests and get summary counts by request status code
// @Description Search for requests using a keyword and get the summary of counts grouped by request status code
// @Tags Received-vehicle-driver
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Param search query string false "Search keyword (matches request_no, vehicle_license_plate, vehicle_user_emp_name, or work_place)"
// @Param ref_request_status_code query string false "Filter by multiple request status codes (comma-separated, e.g., 'A,B,C')"
// @Param startdate query string false "Filter by start datetime (YYYY-MM-DD format)"
// @Param enddate query string false "Filter by end datetime (YYYY-MM-DD format)"
// @Param received_key_start_datetime query string false "Filter by received key start datetime (YYYY-MM-DD format)"
// @Param received_key_end_datetime query string false "Filter by received key end datetime (YYYY-MM-DD format)"
// @Param order_by query string false "Order by request_no, start_datetime, ref_request_status_code"
// @Param order_dir query string false "Order direction: asc or desc"
// @Param page query int false "Page number (default: 1)"
// @Param limit query int false "Number of records per page (default: 10)"
// @Router /api/received-vehicle-driver/search-requests [get]
func (h *ReceivedVehicleDriverHandler) SearchRequests(c *gin.Context) {
	user := funcs.GetAuthenUser(c, h.Role)
	if c.IsAborted() {
		return
	}
	var statusNameMap = StatusNameMapReceivedKeyUser
	var requests []models.VmsTrnRequestList
	var summary []models.VmsTrnRequestSummary

	// Use the keys from statusNameMap as the list of valid status codes
	statusCodes := make([]string, 0, len(statusNameMap))
	for code := range statusNameMap {
		statusCodes = append(statusCodes, code)
	}

	// Build the main query
	query := h.SetQueryRole(user, config.DB)
	query = query.Table("public.vms_trn_request AS req").
		Select("req.*, v.vehicle_license_plate,v.vehicle_license_plate_province_short,v.vehicle_license_plate_province_full,"+
			"(select max(parking_place) from vms_mas_vehicle_department d where d.mas_vehicle_uid = req.mas_vehicle_uid AND d.is_deleted = '0' AND d.is_active = '1') parking_place ").
		Joins("LEFT JOIN vms_mas_vehicle v on v.mas_vehicle_uid = req.mas_vehicle_uid").
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

	if receivedKeyStartDatetime := c.Query("received_key_start_datetime"); receivedKeyStartDatetime != "" {
		query = query.Where("req.received_key_start_datetime >= ?", receivedKeyStartDatetime)
	}
	if receivedKeyEndDatetime := c.Query("received_key_end_datetime"); receivedKeyEndDatetime != "" {
		query = query.Where("req.received_key_end_datetime <= ?", receivedKeyEndDatetime)
	}
	if refRequestStatusCodes := c.Query("ref_request_status_code"); refRequestStatusCodes != "" {
		// Split the comma-separated codes into a slice
		has51e := false
		codes := strings.Split(refRequestStatusCodes, ",")
		for i := range codes {
			if codes[i] == "51e" {
				has51e = true
				codes[i] = "51"
			}
		}
		query = query.Where("req.ref_request_status_code IN (?)", codes)
		if has51e {
			query = query.Where("req.ref_request_status_code = '51' AND reserve_start_datetime < NOW()")
		}
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
	default:
		query = query.Order("req.request_no desc")
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
	}

	// Build the summary query
	summaryQuery := h.SetQueryRole(user, config.DB)
	summaryQuery = summaryQuery.Table("public.vms_trn_request AS req").
		Select(`CASE 
			WHEN req.ref_request_status_code = '51' AND reserve_start_datetime < NOW() THEN '51e'
			ELSE req.ref_request_status_code
		END as ref_request_status_code, COUNT(*) as count`).
		Group(`CASE 
			WHEN req.ref_request_status_code = '51' AND reserve_start_datetime < NOW() THEN '51e'
			ELSE req.ref_request_status_code
		END`)

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
		requests = []models.VmsTrnRequestList{}
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
// @Tags Received-vehicle-driver
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Param trn_request_uid path string true "TrnRequestUID (trn_request_uid)"
// @Router /api/received-vehicle-driver/request/{trn_request_uid} [get]
func (h *ReceivedVehicleDriverHandler) GetRequest(c *gin.Context) {
	funcs.GetAuthenUser(c, h.Role)
	if c.IsAborted() {
		return
	}
	request, err := funcs.GetRequestVehicelInUse(c, StatusNameMapReceivedKeyUser)
	if err != nil {
		return
	}
	c.JSON(http.StatusOK, request)
}

// ReceivedVehicle godoc
// @Summary Update vehicle pickup for a booking request
// @Description This endpoint allows to update vehicle pickup.
// @Tags Received-vehicle-driver
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Param data body models.VmsTrnReceivedVehicle true "VmsTrnReceivedVehicle data"
// @Router /api/received-vehicle-driver/received-vehicle [put]
func (h *ReceivedVehicleDriverHandler) ReceivedVehicle(c *gin.Context) {
	user := funcs.GetAuthenUser(c, h.Role)
	if c.IsAborted() {
		return
	}
	var request, trnRequest models.VmsTrnReceivedVehicle
	var result struct {
		models.VmsTrnReceivedVehicle
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

	request.RefRequestStatusCode = "60"
	request.UpdatedAt = time.Now()
	request.UpdatedBy = user.EmpID

	for i := range request.VehicleImages {
		request.VehicleImages[i].TrnVehicleImgReceivedUID = uuid.New().String()
		request.VehicleImages[i].TrnRequestUID = request.TrnRequestUID
		request.VehicleImages[i].CreatedAt = time.Now()
		request.VehicleImages[i].CreatedBy = user.EmpID
		request.VehicleImages[i].UpdatedAt = time.Now()
		request.VehicleImages[i].UpdatedBy = user.EmpID
		request.VehicleImages[i].IsDeleted = "0"
	}

	if len(request.VehicleImages) > 0 {
		if err := config.DB.Where("trn_request_uid = ?", request.TrnRequestUID).Delete(&models.VehicleImageReceived{}).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update vehicle images", "message": messages.ErrInternalServer.Error()})
			return
		}
	}

	if err := config.DB.Save(&request).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to update : %v", err), "message": messages.ErrInternalServer.Error()})
		return
	}

	if err := config.DB.
		Preload("VehicleImages").
		First(&result, "trn_request_uid = ?", request.TrnRequestUID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Booking not found", "message": messages.ErrBookingNotFound.Error()})
		return
	}

	if result.RefRequestStatusCode == request.RefRequestStatusCode {
		funcs.CreateTrnRequestActionLog(result.TrnRequestUID,
			result.RefRequestStatusCode,
			"กรุณาบันทึกเลขไมล์และการเติมเชื้อเพลิง",
			user.EmpID,
			"driver",
			"",
		)
	}
	c.JSON(http.StatusOK, gin.H{"message": "Updated successfully", "result": result})
}

// GetTravelCard godoc
// @Summary Retrieve a travel-card of pecific booking request
// @Description This endpoint fetches a travel-card of pecific booking request using its unique identifier (TrnRequestUID).
// @Tags Received-vehicle-driver
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Param trn_request_uid path string true "TrnRequestUID (trn_request_uid)"
// @Router /api/received-vehicle-driver/travel-card/{trn_request_uid} [get]
func (h *ReceivedVehicleDriverHandler) GetTravelCard(c *gin.Context) {
	user := funcs.GetAuthenUser(c, h.Role)
	if c.IsAborted() {
		return
	}
	id := c.Param("trn_request_uid")
	trnRequestUID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid TrnRequestUID"})
		return
	}
	var request models.VmsTrnTravelCard
	query := h.SetQueryRole(user, config.DB)
	query = query.Table("public.vms_trn_request AS req").
		Select("req.*, v.vehicle_license_plate,v.vehicle_license_plate_province_short,v.vehicle_license_plate_province_full").
		Joins("LEFT JOIN vms_mas_vehicle v on v.mas_vehicle_uid = req.mas_vehicle_uid")

	if err := query.
		First(&request, "trn_request_uid = ?", trnRequestUID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Booking not found"})
		return
	}
	request.VehicleUserImageURL = funcs.GetEmpImage(request.VehicleUserEmpID)
	request.VehicleUserDeptSAPShort = request.VehicleUserPosition + " " + request.VehicleUserDeptSAPShort
	request.ApprovedRequestDeptSAPShort = request.ApprovedRequestPosition + " " + request.ApprovedRequestDeptSAPShort

	c.JSON(http.StatusOK, request)

}
