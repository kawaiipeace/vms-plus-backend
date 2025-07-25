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

type ReceivedVehicleUserHandler struct {
	Role string
}

var StatusNameMapReceivedVehicleUser = map[string]string{
	"51":  "รับยานพาหนะ",
	"60":  "เดินทาง",
	"60e": "เกินวันเดินทาง",
}

func (h *ReceivedVehicleUserHandler) SetQueryRole(user *models.AuthenUserEmp, query *gorm.DB) *gorm.DB {
	return query.Where("created_request_emp_id = ? OR vehicle_user_emp_id = ?", user.EmpID, user.EmpID)
}
func (h *ReceivedVehicleUserHandler) SetQueryStatusCanUpdate(query *gorm.DB) *gorm.DB {
	return query.Where("ref_request_status_code in ('51','60','71') and is_deleted = '0'")
}

// SearchRequests godoc
// @Summary Search booking requests and get summary counts by request status code
// @Description Search for requests using a keyword and get the summary of counts grouped by request status code
// @Tags Received-vehicle-user
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
// @Router /api/received-vehicle-user/search-requests [get]
func (h *ReceivedVehicleUserHandler) SearchRequests(c *gin.Context) {
	user := funcs.GetAuthenUser(c, h.Role)
	if c.IsAborted() {
		return
	}
	statusNameMap := StatusNameMapReceivedVehicleUser
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
			"(select max(parking_place) from vms_mas_vehicle_department d where d.mas_vehicle_uid = req.mas_vehicle_uid and d.is_deleted = '0' and d.is_active = '1') parking_place,"+
			"k.receiver_personal_id,k.receiver_fullname,k.receiver_dept_sap,"+
			"k.appointment_start appointment_key_handover_start_datetime,k.appointment_end appointment_key_handover_end_datetime,k.appointment_location appointment_key_handover_place,"+
			"k.receiver_dept_name_short,k.receiver_dept_name_full,k.receiver_desk_phone,k.receiver_mobile_phone,k.receiver_position,k.remark receiver_remark").
		Joins("LEFT JOIN vms_trn_vehicle_key_handover k ON k.trn_request_uid = req.trn_request_uid").
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
		requests[i].CanScoreButton = funcs.IsAllowScoreButton(requests[i].TrnRequestUid)
		if requests[i].CanScoreButton {
			requests[i].CanPickupButton = false
		} else {
			requests[i].CanPickupButton = funcs.IsAllowPickupButton(requests[i].TrnRequestUid)
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
// @Tags Received-vehicle-user
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Param trn_request_uid path string true "TrnRequestUID (trn_request_uid)"
// @Router /api/received-vehicle-user/request/{trn_request_uid} [get]
func (h *ReceivedVehicleUserHandler) GetRequest(c *gin.Context) {
	funcs.GetAuthenUser(c, h.Role)
	if c.IsAborted() {
		return
	}
	request, err := funcs.GetRequestVehicelInUse(c, StatusNameMapReceivedVehicleUser)
	if err != nil {
		return
	}
	c.JSON(http.StatusOK, request)
}

// ReceivedVehicle godoc
// @Summary Update vehicle pickup for a booking request
// @Description This endpoint allows to update vehicle pickup.
// @Tags Received-vehicle-user
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Param data body models.VmsTrnReceivedVehicle true "VmsTrnReceivedVehicle data"
// @Router /api/received-vehicle-user/received-vehicle [put]
func (h *ReceivedVehicleUserHandler) ReceivedVehicle(c *gin.Context) {
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

	empUser := funcs.GetUserEmpInfo(user.EmpID)
	request.ReceivedVehicleEmpID = empUser.EmpID
	request.ReceivedVehicleEmpName = empUser.FullName
	request.ReceivedVehicleDeptSAP = empUser.DeptSAP
	request.ReceivedVehicleDeptSAPShort = funcs.GetDeptSAPShort(empUser.DeptSAP)
	request.ReceivedVehicleDeptSAPFull = funcs.GetDeptSAPFull(empUser.DeptSAP)
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
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update vehicle images"})
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
			"นำยานพาหนะออกไปปฎิบัติงานแล้ว",
			user.EmpID,
			"vehicle_user",
			"",
		)
	}

	c.JSON(http.StatusOK, gin.H{"message": "Updated successfully", "result": result})
}

// GetTravelCard godoc
// @Summary Retrieve a travel-card of pecific booking request
// @Description This endpoint fetches a travel-card of pecific booking request using its unique identifier (TrnRequestUID).
// @Tags Received-vehicle-user
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Param trn_request_uid path string true "TrnRequestUID (trn_request_uid)"
// @Router /api/received-vehicle-user/travel-card/{trn_request_uid} [get]
func (h *ReceivedVehicleUserHandler) GetTravelCard(c *gin.Context) {
	user := funcs.GetAuthenUser(c, h.Role)
	if c.IsAborted() {
		return
	}
	id := c.Param("trn_request_uid")
	trnRequestUID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid TrnRequestUID", "message": messages.ErrInvalidUID.Error()})
		return
	}
	var request models.VmsTrnTravelCard
	query := h.SetQueryRole(user, config.DB)
	query = query.Table("public.vms_trn_request AS req").
		Select("req.*, v.vehicle_license_plate,v.vehicle_license_plate_province_short,v.vehicle_license_plate_province_full").
		Joins("LEFT JOIN vms_mas_vehicle v on v.mas_vehicle_uid = req.mas_vehicle_uid")

	if err := query.
		First(&request, "trn_request_uid = ?", trnRequestUID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Booking not found", "message": messages.ErrBookingNotFound.Error()})
		return
	}
	request.VehicleUserImageURL = funcs.GetEmpImage(request.VehicleUserEmpID)
	request.ApprovedRequestDeptSAPShort = request.ApprovedRequestPosition + " " + request.ApprovedRequestDeptSAPShort

	c.JSON(http.StatusOK, request)

}
