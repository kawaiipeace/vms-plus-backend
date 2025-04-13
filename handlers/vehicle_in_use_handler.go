package handlers

import (
	"fmt"
	"net/http"
	"time"
	"vms_plus_be/config"
	"vms_plus_be/funcs"
	"vms_plus_be/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type VehicleInUseHandler struct {
	Role string
}

// SearchRequests godoc
// @Summary Search booking requests and get summary counts by request status code
// @Description Search for requests using a keyword and get the summary of counts grouped by request status code
// @Tags Vehicle-in-use
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
// @Param page_size query int false "Number of records per page (default: 10)"
// @Router /api/vehicle-in-use/search-requests [get]
func (h *VehicleInUseHandler) SearchRequests(c *gin.Context) {
	//funcs.GetAuthenUser(c, h.Role)
	statusNameMap := map[string]string{
		"51": "รับยานพาหนะ",
		"60": "เดินทาง",
	}

	var requests []models.VmsTrnRequest_List
	var summary []models.VmsTrnRequest_Summary

	// Use the keys from statusNameMap as the list of valid status codes
	statusCodes := make([]string, 0, len(statusNameMap))
	for code := range statusNameMap {
		statusCodes = append(statusCodes, code)
	}

	// Build the main query
	query := config.DB.Table("public.vms_trn_request AS req").
		Select("req.*, status.ref_request_status_desc").
		Joins("LEFT JOIN public.vms_ref_request_status AS status ON req.ref_request_status_code = status.ref_request_status_code").
		Where("req.ref_request_status_code IN (?)", statusCodes)

	// Apply additional filters (search, date range, etc.)
	if search := c.Query("search"); search != "" {
		query = query.Where("req.request_no LIKE ? OR req.vehicle_license_plate LIKE ? OR req.vehicle_user_emp_name LIKE ? OR req.work_place LIKE ?", "%"+search+"%", "%"+search+"%", "%"+search+"%", "%"+search+"%")
	}
	if startDate := c.Query("startdate"); startDate != "" {
		query = query.Where("req.start_datetime >= ?", startDate)
	}
	if endDate := c.Query("enddate"); endDate != "" {
		query = query.Where("req.start_datetime <= ?", endDate)
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
		summary = append(summary, models.VmsTrnRequest_Summary{
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
// @Tags Vehicle-in-use
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Param trn_request_uid path string true "TrnRequestUID (trn_request_uid)"
// @Router /api/vehicle-in-use/request/{trn_request_uid} [get]
func (h *VehicleInUseHandler) GetRequest(c *gin.Context) {
	funcs.GetAuthenUser(c, h.Role)
	statusNameMap := map[string]string{
		"51": "รับยานพาหนะ",
		"60": "เดินทาง",
	}
	funcs.GetRequest(c, statusNameMap)
}

// CreateVehicleTripDetail godoc
// @Summary Create Vehicle Travel List for a booking request
// @Description This endpoint allows to Create Vehicle Travel List.
// @Tags Vehicle-in-use
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Param data body models.VmsTrnTripDetail_Request true "VmsTrnTripDetail_Request data"
// @Router /api/vehicle-in-use/create-travel-detail [post]
func (h *VehicleInUseHandler) CreateVehicleTripDetail(c *gin.Context) {
	user := funcs.GetAuthenUser(c, h.Role)

	var req models.VmsTrnTripDetail
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON input"})
		return
	}
	req.TrnTripDetailUID = uuid.New().String()
	req.CreatedBy = user.EmpID
	req.CreatedAt = time.Now()
	req.UpdatedBy = user.EmpID
	req.UpdatedAt = time.Now()

	var existingReq struct {
		MasVehicleUID                   string `gorm:"column:mas_vehicle_uid"`
		VehicleLicensePlate             string `gorm:"column:vehicle_license_plate"`
		VehicleLicensePlateProvinceFull string `gorm:"column:vehicle_license_plate_province_full"`
		MasVehicleDepartmentUID         string `gorm:"column:mas_vehicle_department_uid"`
		MasCarpoolUID                   string `gorm:"column:mas_carpool_uid"`
		EmployeeOrDriverID              string `gorm:"column:driver_emp_id"`
	}
	if err := config.DB.Table("public.vms_trn_request").Where("trn_request_uid = ?", req.TrnRequestUID).First(&existingReq).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Booking not found"})
		return
	}

	req.MasCarpoolUID = existingReq.MasCarpoolUID
	req.MasVehicleDepartmentUID = existingReq.MasVehicleDepartmentUID
	req.VehicleLicensePlate = existingReq.VehicleLicensePlate
	req.VehicleLicensePlateProvinceFull = existingReq.VehicleLicensePlateProvinceFull
	req.EmployeeOrDriverID = existingReq.EmployeeOrDriverID

	if err := config.DB.Create(&req).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create"})
		return
	}

	// Return success response
	c.JSON(http.StatusCreated, gin.H{"message": "Created successfully", "data": req})
}

// UpdateVehicleTripDetail godoc
// @Summary Update Vehicle Travel List for a booking request
// @Description This endpoint allows to Update Vehicle Travel List.
// @Tags Vehicle-in-use
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Param trn_trip_detail_uid path string true "TrnTripDetailUID"
// @Param data body models.VmsTrnTripDetail_Request true "VmsTrnTripDetail_Request data"
// @Router /api/vehicle-in-use/update-travel-detail/{trn_trip_detail_uid} [put]
func (h *VehicleInUseHandler) UpdateVehicleTripDetail(c *gin.Context) {
	user := funcs.GetAuthenUser(c, h.Role)
	id := c.Param("trn_trip_detail_uid")
	uid, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	var existing models.VmsTrnTripDetail
	if err := config.DB.Where("trn_trip_detail_uid = ?", uid).First(&existing).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Trip not found"})
		return
	}
	var req models.VmsTrnTripDetail_Request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	existing.VmsTrnTripDetail_Request = req
	existing.UpdatedBy = user.EmpID
	existing.UpdatedAt = time.Now()

	if err := config.DB.Save(&existing).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to update : %v", err)})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Updated successfully", "data": req})
}

// DeleteVehicleTripDetail godoc
// @Summary Update Vehicle Travel List for a booking request
// @Description This endpoint allows to Update Vehicle Travel List.
// @Tags Vehicle-in-use
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Param trn_trip_detail_uid path string true "TrnTripDetailUID"
// @Router /api/vehicle-in-use/delete-travel-detail/{trn_trip_detail_uid} [delete]
func (h *VehicleInUseHandler) DeleteVehicleTripDetail(c *gin.Context) {
	user := funcs.GetAuthenUser(c, h.Role)
	id := c.Param("trn_trip_detail_uid")
	uid, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var existing models.VmsTrnTripDetail
	if err := config.DB.Where("trn_trip_detail_uid = ?", uid).First(&existing).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Trip not found"})
		return
	}

	if err := config.DB.Model(&existing).UpdateColumns(map[string]interface{}{
		"is_deleted": true,
		"updated_by": user.EmpID,
		"updated_at": time.Now(),
	}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete"})
	}

	c.JSON(http.StatusOK, gin.H{"message": "Deleted successfully"})
}

// GetVehicleTripDetails godoc
// @Summary Retrieve list of trip detail
// @Description Retrieve a list of trip detail in TrnRequestUID
// @Tags Vehicle-in-use
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Param trn_request_uid path string true "TrnRequestUID (trn_request_uid)"
// @Router /api/vehicle-in-use/travel-details/{trn_request_uid} [get]
func (h *VehicleInUseHandler) GetVehicleTripDetails(c *gin.Context) {
	id := c.Param("trn_request_uid")
	uid, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	// Fetch the vehicle record from the database
	var trips []models.VmsTrnTripDetail_List
	if err := config.DB.Find(&trips, "trn_trip_detail_uid = ? AND is_deleted = false", uid).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Trip not found"})
		return
	}

	// Return the vehicle data as a JSON response
	c.JSON(http.StatusOK, trips)
}

// GetVehicleTripDetail godoc
// @Summary Retrieve details of a specific trip detail
// @Description Fetch detailed information about a trip detail using their unique TrnTripDetailUID.
// @Tags Vehicle-in-use
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Param trn_trip_detail_uid path string true "TrnTripDetailUID"
// @Router /api/vehicle-in-use/travel-detail/{trn_trip_detail_uid} [get]
func (h *VehicleInUseHandler) GetVehicleTripDetail(c *gin.Context) {
	id := c.Param("trn_trip_detail_uid")
	uid, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	// Fetch the vehicle record from the database
	var trip models.VmsTrnTripDetail
	if err := config.DB.First(&trip, "trn_trip_detail_uid = ? AND is_deleted = false", uid).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Trip not found"})
		return
	}

	// Return the vehicle data as a JSON response
	c.JSON(http.StatusOK, trip)
}

// CreateVehicleAddFuel godoc
// @Summary Create Vehicle Add Fuel entry
// @Description This endpoint allows to create a new Vehicle Add Fuel entry.
// @Tags Vehicle-in-use
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Param data body models.VmsTrnAddFuel true "VmsTrnAddFuel data"
// @Router /api/vehicle-in-use/create-add-fuel [post]
func (h *VehicleInUseHandler) CreateVehicleAddFuel(c *gin.Context) {
	user := funcs.GetAuthenUser(c, h.Role)

	var req models.VmsTrnAddFuel
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON input"})
		return
	}

	req.TrnAddFuelUID = uuid.New().String()
	req.CreatedBy = user.EmpID
	req.CreatedAt = time.Now()
	req.UpdatedBy = user.EmpID
	req.UpdatedAt = time.Now()

	if err := config.DB.Create(&req).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Created successfully", "data": req})
}

// UpdateVehicleAddFuel godoc
// @Summary Update Vehicle Add Fuel entry
// @Description This endpoint allows to update an existing Vehicle Add Fuel entry.
// @Tags Vehicle-in-use
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Param trn_add_fuel_uid path string true "TrnAddFuelUID"
// @Param data body models.VmsTrnAddFuel true "VmsTrnAddFuel data"
// @Router /api/vehicle-in-use/update-add-fuel/{trn_add_fuel_uid} [put]
func (h *VehicleInUseHandler) UpdateVehicleAddFuel(c *gin.Context) {
	user := funcs.GetAuthenUser(c, h.Role)
	id := c.Param("trn_add_fuel_uid")

	var existing models.VmsTrnAddFuel
	if err := config.DB.Where("trn_add_fuel_uid = ?", id).First(&existing).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Add Fuel entry not found"})
		return
	}

	var req models.VmsTrnAddFuel
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON input"})
		return
	}

	existing = req
	existing.UpdatedBy = user.EmpID
	existing.UpdatedAt = time.Now()

	if err := config.DB.Save(&existing).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Updated successfully", "data": existing})
}

// DeleteVehicleAddFuel godoc
// @Summary Delete Vehicle Add Fuel entry
// @Description This endpoint allows to mark a Vehicle Add Fuel entry as deleted.
// @Tags Vehicle-in-use
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Param trn_add_fuel_uid path string true "TrnAddFuelUID"
// @Router /api/vehicle-in-use/delete-add-fuel/{trn_add_fuel_uid} [delete]
func (h *VehicleInUseHandler) DeleteVehicleAddFuel(c *gin.Context) {
	user := funcs.GetAuthenUser(c, h.Role)
	id := c.Param("trn_add_fuel_uid")
	uid, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var existing models.VmsTrnTripDetail
	if err := config.DB.Where("trn_add_fuel_uid = ?", uid).First(&existing).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Trip not found"})
		return
	}

	if err := config.DB.Model(&existing).UpdateColumns(map[string]interface{}{
		"is_deleted": true,
		"updated_by": user.EmpID,
		"updated_at": time.Now(),
	}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete"})
	}

	c.JSON(http.StatusOK, gin.H{"message": "Deleted successfully"})
}

// GetVehicleAddFuelDetails godoc
// @Summary Retrieve a list of Add Fuel entries
// @Description Fetch a list of Add Fuel entries in TrnRequestUID.
// @Tags Vehicle-in-use
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Param trn_request_uid path string true "TrnRequestUID"
// @Router /api/vehicle-in-use/add-fuel-details/{trn_request_uid} [get]
func (h *VehicleInUseHandler) GetVehicleAddFuelDetails(c *gin.Context) {
	id := c.Param("trn_request_uid")

	var fuels []models.VmsTrnAddFuel
	if err := config.DB.Where("trn_request_uid = ? AND is_deleted = false", id).Find(&fuels).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Add Fuel entries not found"})
		return
	}

	c.JSON(http.StatusOK, fuels)
}

// GetVehicleAddFuelDetail godoc
// @Summary Retrieve details of a specific Add Fuel entry
// @Description Fetch detailed information about an Add Fuel entry using its unique TrnAddFuelUID.
// @Tags Vehicle-in-use
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Param trn_add_fuel_uid path string true "TrnAddFuelUID"
// @Router /api/vehicle-in-use/add-fuel-detail/{trn_add_fuel_uid} [get]
func (h *VehicleInUseHandler) GetVehicleAddFuelDetail(c *gin.Context) {
	id := c.Param("trn_add_fuel_uid")

	var fuel models.VmsTrnAddFuel
	if err := config.DB.Where("trn_add_fuel_uid = ? AND is_deleted = false", id).First(&fuel).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Add Fuel entry not found"})
		return
	}

	c.JSON(http.StatusOK, fuel)
}

// GetTravelCard godoc
// @Summary Retrieve a travel-card of pecific booking request
// @Description This endpoint fetches a travel-card of pecific booking request using its unique identifier (TrnRequestUID).
// @Tags Vehicle-in-use
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Param trn_request_uid path string true "TrnRequestUID (trn_request_uid)"
// @Router /api/vehicle-in-use/travel-card/{trn_request_uid} [get]
func (h *VehicleInUseHandler) GetTravelCard(c *gin.Context) {
	//funcs.GetAuthenUser(c, h.Role)
	id := c.Param("trn_request_uid")
	trnRequestUID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid TrnRequestUID"})
		return
	}
	var request models.VmsTrnTravelCard
	if err := config.DB.
		First(&request, "trn_request_uid = ?", trnRequestUID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Booking not found"})
		return
	}
	request.VehicleUserImageURL = config.DefaultURL
	c.JSON(http.StatusOK, request)
}

// UpdateVehicleTripDetail godoc
// @Summary Update Satisfaction Survey for a booking request
// @Description This endpoint allows to Update Satisfaction Survey for a booking request.
// @Tags Vehicle-in-use
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Param trn_request_uid path string true "TrnRequestUID"
// @Param data body []models.VmsTrnSatisfactionSurveyAnswers true "Array of VmsTrnSatisfactionSurveyAnswers data"
// @Router /api/vehicle-in-use/update_satisfaction_survey/{trn_request_uid} [put]
func (h *VehicleInUseHandler) UpdateSatisfactionSurvey(c *gin.Context) {
	user := funcs.GetAuthenUser(c, h.Role)
	id := c.Param("trn_request_uid")
	trnRequestUID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	var existing models.VmsTrnRequest_Response
	if err := config.DB.First(&existing, "trn_request_uid = ?", trnRequestUID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Booking not found"})
		return
	}

	var reqs []models.VmsTrnSatisfactionSurveyAnswers
	if err := c.ShouldBindJSON(&reqs); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	for _, req := range reqs {
		var existingAnswer models.VmsTrnSatisfactionSurveyAnswers
		if err := config.DB.First(&existingAnswer, "trn_request_uid = ? AND mas_satisfaction_survey_questions_code = ?", trnRequestUID, req.MasSatisfactionSurveyQuestionsCode).Error; err != nil {
			if gorm.ErrRecordNotFound == err {
				// Handle record not found logic (create new record)
				req.TrnSatisfactionSurveyAnswersUID = uuid.NewString()
				req.TrnRequestUID = trnRequestUID.String()
				req.SurveyAnswerDate = time.Now()
				req.SurveyAnswerEmpID = user.EmpID
				if err := config.DB.Create(&req).Error; err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to create record: %v", err)})
					return
				}
			} else {
				// Handle other errors
				c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to check existence: %v", err)})
				return
			}
		} else {
			// Handle record exists (update logic)
			existingAnswer.TrnRequestUID = trnRequestUID.String()
			existingAnswer.SurveyAnswer = req.SurveyAnswer
			existingAnswer.SurveyAnswerDate = time.Now()
			existingAnswer.SurveyAnswerEmpID = user.EmpID

			if err := config.DB.Save(&existingAnswer).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to update record: %v", err)})
				return
			}
		}
	}
	c.JSON(http.StatusOK, gin.H{"message": "Updated successfully", "data": reqs})
}
