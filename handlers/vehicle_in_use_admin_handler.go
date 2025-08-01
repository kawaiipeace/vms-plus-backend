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

type VehicleInUseAdminHandler struct {
	Role string
}

var StatusNameMapVehicelInUseAdmin = map[string]string{
	"50":  "รอรับยานพาหนะ",
	"50e": "รับยานพาหนะล่าช้า",
	"51":  "รับยานพาหนะ",
	"60":  "อยู่ระหว่างเดินทาง",
	"60e": "คืนยานพาหนะล่าช้า",
	"70":  "ส่งคืนยานพาหนะ",
	"80":  "เสร็จสิ้น",
}

func (h *VehicleInUseAdminHandler) SetQueryRole(user *models.AuthenUserEmp, query *gorm.DB) *gorm.DB {
	query = funcs.SetQueryAdminRole(user, query)
	return query
}
func (h *VehicleInUseAdminHandler) SetQueryStatusCanUpdate(query *gorm.DB) *gorm.DB {
	return query.Where("ref_request_status_code in ('60','70','71') and is_deleted = '0'")
}

// SearchRequests godoc
// @Summary Search booking requests and get summary counts by request status code
// @Description Search for requests using a keyword and get the summary of counts grouped by request status code
// @Tags Vehicle-in-use-admin
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Param search query string false "Search keyword (matches request_no, vehicle_license_plate, vehicle_user_emp_name, or work_place)"
// @Param ref_request_status_code query string false "Filter by multiple request status codes (comma-separated, e.g., 'A,B,C')"
// @Param vehicle_owner_dept_sap query string false "Filter by vehicle owner department SAP"
// @Param startdate query string false "Filter by start datetime (YYYY-MM-DD format)"
// @Param enddate query string false "Filter by end datetime (YYYY-MM-DD format)"
// @Param order_by query string false "Order by request_no, start_datetime, ref_request_status_code"
// @Param order_dir query string false "Order direction: asc or desc"
// @Param page query int false "Page number (default: 1)"
// @Param limit query int false "Number of records per page (default: 10)"
// @Router /api/vehicle-in-use-admin/search-requests [get]
func (h *VehicleInUseAdminHandler) SearchRequests(c *gin.Context) {
	user := funcs.GetAuthenUser(c, h.Role)
	if c.IsAborted() {
		return
	}
	statusNameMap := StatusNameMapVehicelInUseAdmin
	var requests []models.VmsTrnRequestVehicleInUseList
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
			"(select max(parking_place) from vms_mas_vehicle_department d where d.mas_vehicle_uid = vms_trn_request.mas_vehicle_uid AND d.is_deleted = '0' AND d.is_active = '1') parking_place ").
		Joins("LEFT JOIN vms_mas_vehicle v on v.mas_vehicle_uid = vms_trn_request.mas_vehicle_uid").
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
		requests = []models.VmsTrnRequestVehicleInUseList{}
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
// @Tags Vehicle-in-use-admin
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Param trn_request_uid path string true "TrnRequestUID (trn_request_uid)"
// @Router /api/vehicle-in-use-admin/request/{trn_request_uid} [get]
func (h *VehicleInUseAdminHandler) GetRequest(c *gin.Context) {
	funcs.GetAuthenUser(c, h.Role)
	if c.IsAborted() {
		return
	}
	funcs.GetRequest(c, StatusNameMapVehicelInUseAdmin)

	funcs.GetAuthenUser(c, h.Role)
	if c.IsAborted() {
		return
	}

	request, err := funcs.GetRequestVehicelInUse(c, StatusNameMapUser)
	if err != nil {
		return
	}

	c.JSON(http.StatusOK, request)
}

// CreateVehicleTripDetail godoc
// @Summary Create Vehicle Travel List for a booking request
// @Description This endpoint allows to Create Vehicle Travel List.
// @Tags Vehicle-in-use-admin
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Param data body models.VmsTrnTripDetailRequest true "VmsTrnTripDetailRequest data"
// @Router /api/vehicle-in-use-admin/create-travel-detail [post]
func (h *VehicleInUseAdminHandler) CreateVehicleTripDetail(c *gin.Context) {
	user := funcs.GetAuthenUser(c, h.Role)
	if c.IsAborted() {
		return
	}

	var request models.VmsTrnTripDetail
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON input", "message": messages.ErrInvalidJSONInput.Error()})
		return
	}
	request.TrnTripDetailUID = uuid.New().String()
	request.CreatedBy = user.EmpID
	request.CreatedAt = time.Now()
	request.UpdatedBy = user.EmpID
	request.UpdatedAt = time.Now()

	var trnRequest struct {
		MasVehicleUID                   string `gorm:"column:mas_vehicle_uid"`
		VehicleLicensePlate             string `gorm:"column:vehicle_license_plate"`
		VehicleLicensePlateProvinceFull string `gorm:"column:vehicle_license_plate_province_full"`
		MasVehicleDepartmentUID         string `gorm:"column:mas_vehicle_department_uid"`
		MasCarpoolUID                   string `gorm:"column:mas_carpool_uid"`
		EmployeeOrDriverID              string `gorm:"column:driver_emp_id"`
	}
	query := h.SetQueryRole(user, config.DB)
	query = h.SetQueryStatusCanUpdate(query)
	if err := query.Table("public.vms_trn_request").Where("trn_request_uid = ?", request.TrnRequestUID).First(&trnRequest).Error; err != nil {
		c.JSON(http.StatusMethodNotAllowed, gin.H{"error": "Booking can not update", "message": messages.ErrBookingNotFound.Error()})
		return
	}
	request.MasVehicleUID = func() string {
		if trnRequest.MasVehicleUID == "" {
			return funcs.DefaultUUID()
		} else {
			return trnRequest.MasVehicleUID
		}
	}()
	request.MasCarpoolUID = func() string {
		if trnRequest.MasCarpoolUID == "" {
			return funcs.DefaultUUID()
		} else {
			return trnRequest.MasCarpoolUID
		}
	}()
	request.MasVehicleDepartmentUID = func() string {
		if trnRequest.MasVehicleDepartmentUID == "" {
			return funcs.DefaultUUID()
		} else {
			return trnRequest.MasVehicleDepartmentUID
		}
	}()
	request.VehicleLicensePlate = trnRequest.VehicleLicensePlate
	request.VehicleLicensePlateProvinceFull = trnRequest.VehicleLicensePlateProvinceFull
	request.EmployeeOrDriverID = trnRequest.EmployeeOrDriverID
	request.IsDeleted = "0"

	if err := config.DB.Create(&request).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create", "message": messages.ErrInternalServer.Error()})
		return
	}

	funcs.UpdateVehicleMileage(request.TrnRequestUID, request.TripEndMiles)
	// Return success response
	c.JSON(http.StatusCreated, gin.H{"message": "Created successfully", "data": request})
}

// UpdateVehicleTripDetail godoc
// @Summary Update Vehicle Travel List for a booking request
// @Description This endpoint allows to Update Vehicle Travel List.
// @Tags Vehicle-in-use-admin
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Param trn_trip_detail_uid path string true "TrnTripDetailUID"
// @Param data body models.VmsTrnTripDetailRequest true "VmsTrnTripDetailRequest data"
// @Router /api/vehicle-in-use-admin/update-travel-detail/{trn_trip_detail_uid} [put]
func (h *VehicleInUseAdminHandler) UpdateVehicleTripDetail(c *gin.Context) {
	user := funcs.GetAuthenUser(c, h.Role)
	if c.IsAborted() {
		return
	}
	uid := c.Param("trn_trip_detail_uid")
	trnTripDetailUid, err := uuid.Parse(uid)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Uid", "message": messages.ErrInvalidUID.Error()})
		return
	}
	var existing models.VmsTrnTripDetail
	if err := config.DB.Where("trn_trip_detail_uid = ? AND is_deleted = ?", trnTripDetailUid, "0").First(&existing).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Trip not found", "message": messages.ErrNotfound.Error()})
		return
	}
	var request models.VmsTrnTripDetailRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "message": messages.ErrInvalidJSONInput.Error()})
		return
	}

	var trnRequest models.VmsTrnRequestList
	query := h.SetQueryRole(user, config.DB)
	query = h.SetQueryStatusCanUpdate(query)
	if err := query.Table("public.vms_trn_request").Where("trn_request_uid = ?", existing.TrnRequestUID).First(&trnRequest).Error; err != nil {
		c.JSON(http.StatusMethodNotAllowed, gin.H{"error": "Booking can not update", "message": messages.ErrBookingNotFound.Error()})
		return
	}

	existing.VmsTrnTripDetailRequest = request
	existing.UpdatedBy = user.EmpID
	existing.UpdatedAt = time.Now()

	if err := config.DB.Save(&existing).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to update : %v", err), "message": messages.ErrInternalServer.Error()})
		return
	}

	//MaxMileage
	var maxMileage int
	if err := config.DB.Table("vms_trn_trip_detail").
		Where("trn_request_uid = ?", existing.TrnRequestUID).
		Select("MAX(trip_end_miles)").Scan(&maxMileage).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get max mileage", "message": messages.ErrInternalServer.Error()})
		return
	}
	funcs.UpdateVehicleMileage(existing.TrnRequestUID, maxMileage)

	c.JSON(http.StatusOK, gin.H{"message": "Updated successfully", "data": request})
}

// DeleteVehicleTripDetail godoc
// @Summary Update Vehicle Travel List for a booking request
// @Description This endpoint allows to Update Vehicle Travel List.
// @Tags Vehicle-in-use-admin
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Param trn_trip_detail_uid path string true "TrnTripDetailUID"
// @Router /api/vehicle-in-use-admin/delete-travel-detail/{trn_trip_detail_uid} [delete]
func (h *VehicleInUseAdminHandler) DeleteVehicleTripDetail(c *gin.Context) {
	user := funcs.GetAuthenUser(c, h.Role)
	if c.IsAborted() {
		return
	}
	uid := c.Param("trn_trip_detail_uid")
	trnTripDetailUid, err := uuid.Parse(uid)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Uid", "message": messages.ErrInvalidUID.Error()})
		return
	}

	var existing models.VmsTrnTripDetail
	if err := config.DB.Where("trn_trip_detail_uid = ? AND is_deleted = ?", trnTripDetailUid, "0").First(&existing).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Trip not found"})
		return
	}
	var trnRequest models.VmsTrnRequestList
	query := h.SetQueryRole(user, config.DB)
	query = h.SetQueryStatusCanUpdate(query)
	if err := query.Table("public.vms_trn_request").Where("trn_request_uid = ?", existing.TrnRequestUID).First(&trnRequest).Error; err != nil {
		c.JSON(http.StatusMethodNotAllowed, gin.H{"error": "Booking can not update", "message": messages.ErrBookingNotFound.Error()})
		return
	}
	if err := config.DB.Model(&existing).UpdateColumns(map[string]interface{}{
		"is_deleted": "1",
		"updated_by": user.EmpID,
		"updated_at": time.Now(),
	}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete"})
	}

	//max tripendmiles
	var maxMiles int
	if err := query.Table("public.vms_trn_trip_detail").Where("trn_request_uid = ? AND is_deleted = ?", existing.TrnRequestUID, "0").
		Select("MAX(trip_end_miles) as max_trip_end_miles").First(&maxMiles).Error; err == nil {
		funcs.UpdateVehicleMileage(existing.TrnRequestUID, maxMiles)
	}
	c.JSON(http.StatusOK, gin.H{"message": "Deleted successfully"})
}

// GetVehicleTripDetails godoc
// @Summary Retrieve list of trip detail
// @Description Retrieve a list of trip detail in TrnRequestUID
// @Tags Vehicle-in-use-admin
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Param trn_request_uid path string true "TrnRequestUID (trn_request_uid)"
// @Param search query string false "Search keyword (matches place)"
// @Router /api/vehicle-in-use-admin/travel-details/{trn_request_uid} [get]
func (h *VehicleInUseAdminHandler) GetVehicleTripDetails(c *gin.Context) {
	user := funcs.GetAuthenUser(c, h.Role)
	if c.IsAborted() {
		return
	}
	uid := c.Param("trn_request_uid")
	trnRequestUid, err := uuid.Parse(uid)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Uid"})
		return
	}
	var trnRequest models.VmsTrnRequestList
	query := h.SetQueryRole(user, config.DB)
	if err := query.Table("public.vms_trn_request").Where("trn_request_uid = ?", trnRequestUid).First(&trnRequest).Error; err != nil {
		c.JSON(http.StatusMethodNotAllowed, gin.H{"error": "Booking can not update", "message": messages.ErrBookingNotFound.Error()})
		return
	}

	queryTrip := config.DB
	if search := c.Query("search"); search != "" {
		queryTrip = query.Where("trip_departure_place ILIKE ? OR trip_destination_place ILIKE ?", "%"+search+"%", "%"+search+"%")
	}

	// Fetch the vehicle record from the database
	var trips []models.VmsTrnTripDetailList
	if err := queryTrip.Find(&trips, "trn_request_uid = ? AND is_deleted = ?", trnRequestUid, "0").Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Trip not found", "message": messages.ErrNotfound.Error()})
		return
	}

	// Return the vehicle data as a JSON response
	c.JSON(http.StatusOK, trips)
}

// GetVehicleTripDetail godoc
// @Summary Retrieve details of a specific trip detail
// @Description Fetch detailed information about a trip detail using their unique TrnTripDetailUID.
// @Tags Vehicle-in-use-admin
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Param trn_trip_detail_uid path string true "TrnTripDetailUID"
// @Router /api/vehicle-in-use-admin/travel-detail/{trn_trip_detail_uid} [get]
func (h *VehicleInUseAdminHandler) GetVehicleTripDetail(c *gin.Context) {
	user := funcs.GetAuthenUser(c, h.Role)
	if c.IsAborted() {
		return
	}
	uid := c.Param("trn_trip_detail_uid")
	trnTripDetailUid, err := uuid.Parse(uid)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Uid", "message": messages.ErrInvalidUID.Error()})
		return
	}
	// Fetch the vehicle record from the database
	var trip models.VmsTrnTripDetail
	if err := config.DB.First(&trip, "trn_trip_detail_uid = ? AND is_deleted = ?", trnTripDetailUid, "0").Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Trip not found", "message": messages.ErrNotfound.Error()})
		return
	}
	var trnRequest models.VmsTrnRequestList
	query := h.SetQueryRole(user, config.DB)
	if err := query.Table("public.vms_trn_request").Where("trn_request_uid = ?", trip.TrnRequestUID).First(&trnRequest).Error; err != nil {
		c.JSON(http.StatusMethodNotAllowed, gin.H{"error": "Booking can not update", "message": messages.ErrBookingNotFound.Error()})
		return
	}
	// Return the vehicle data as a JSON response
	c.JSON(http.StatusOK, trip)
}

// CreateVehicleAddFuel godoc
// @Summary Create Vehicle Add Fuel entry
// @Description This endpoint allows to create a new Vehicle Add Fuel entry.
// @Tags Vehicle-in-use-admin
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Param data body models.VmsTrnAddFuelRequest true "VmsTrnAddFuelRequest data"
// @Router /api/vehicle-in-use-admin/create-add-fuel [post]
func (h *VehicleInUseAdminHandler) CreateVehicleAddFuel(c *gin.Context) {
	user := funcs.GetAuthenUser(c, h.Role)
	if c.IsAborted() {
		return
	}

	var request models.VmsTrnAddFuel
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON input", "message": messages.ErrInvalidJSONInput.Error()})
		return
	}

	var trnRequest struct {
		MasVehicleUID                    string `gorm:"column:mas_vehicle_uid"`
		VehicleLicensePlate              string `gorm:"column:vehicle_license_plate"`
		VehicleLicensePlateProvinceShort string `gorm:"column:vehicle_license_plate_province_short"`
		VehicleLicensePlateProvinceFull  string `gorm:"column:vehicle_license_plate_province_full"`
		MasVehicleDepartmentUID          string `gorm:"column:mas_vehicle_department_uid"`
		MasCarpoolUID                    string `gorm:"column:mas_carpool_uid"`
		RefCostTypeCode                  int    `gorm:"column:ref_cost_type_code"`
	}
	query := h.SetQueryRole(user, config.DB)
	query = h.SetQueryStatusCanUpdate(query)
	if err := query.Table("public.vms_trn_request").Where("trn_request_uid = ?", request.TrnRequestUID).First(&trnRequest).Error; err != nil {
		c.JSON(http.StatusMethodNotAllowed, gin.H{"error": "Booking can not update", "message": messages.ErrBookingNotFound.Error()})
		return
	}
	request.RefCostTypeCode = trnRequest.RefCostTypeCode
	request.VehicleLicensePlate = trnRequest.VehicleLicensePlate
	request.VehicleLicensePlateProvinceShort = trnRequest.VehicleLicensePlateProvinceShort
	request.VehicleLicensePlateProvinceFull = trnRequest.VehicleLicensePlateProvinceFull
	request.MasVehicleUID = func() string {
		if trnRequest.MasVehicleUID == "" {
			return funcs.DefaultUUID()
		} else {
			return trnRequest.MasVehicleUID
		}
	}()
	request.MasVehicleDepartmentUID = func() string {
		if trnRequest.MasVehicleDepartmentUID == "" {
			return funcs.DefaultUUID()
		} else {
			return trnRequest.MasVehicleDepartmentUID
		}
	}()
	request.AddFuelDateTime = time.Now()
	request.TrnAddFuelUID = uuid.New().String()
	request.CreatedBy = user.EmpID
	request.CreatedAt = time.Now()
	request.UpdatedBy = user.EmpID
	request.UpdatedAt = time.Now()
	request.IsDeleted = "0"

	if err := config.DB.Create(&request).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create", "message": messages.ErrInternalServer.Error()})
		return
	}

	funcs.UpdateVehicleMileage(request.TrnRequestUID, request.Mile)

	c.JSON(http.StatusCreated, gin.H{"message": "Created successfully", "data": request})
}

// UpdateVehicleAddFuel godoc
// @Summary Update Vehicle Add Fuel entry
// @Description This endpoint allows to update an existing Vehicle Add Fuel entry.
// @Tags Vehicle-in-use-admin
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Param trn_add_fuel_uid path string true "TrnAddFuelUID"
// @Param data body models.VmsTrnAddFuel true "VmsTrnAddFuel data"
// @Router /api/vehicle-in-use-admin/update-add-fuel/{trn_add_fuel_uid} [put]
func (h *VehicleInUseAdminHandler) UpdateVehicleAddFuel(c *gin.Context) {
	user := funcs.GetAuthenUser(c, h.Role)
	if c.IsAborted() {
		return
	}
	uid := c.Param("trn_add_fuel_uid")
	trnAddFuelUid, err := uuid.Parse(uid)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Uid", "message": messages.ErrInvalidUID.Error()})
		return
	}

	var existing models.VmsTrnAddFuel
	if err := config.DB.Where("trn_add_fuel_uid = ? AND is_deleted = ?", trnAddFuelUid, "0").First(&existing).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Add Fuel entry not found", "message": messages.ErrNotfound.Error()})
		return
	}
	var trnRequest models.VmsTrnRequestList
	query := h.SetQueryRole(user, config.DB)
	query = h.SetQueryStatusCanUpdate(query)
	if err := query.Table("public.vms_trn_request").Where("trn_request_uid = ?", existing.TrnRequestUID).First(&trnRequest).Error; err != nil {
		c.JSON(http.StatusMethodNotAllowed, gin.H{"error": "Booking can not update", "message": messages.ErrBookingNotFound.Error()})
		return
	}
	var request models.VmsTrnAddFuelRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON input", "message": messages.ErrInvalidJSONInput.Error()})
		return
	}

	existing.VmsTrnAddFuelRequest = request
	existing.UpdatedBy = user.EmpID
	existing.UpdatedAt = time.Now()

	if err := config.DB.Save(&existing).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update", "message": messages.ErrInternalServer.Error()})
		return
	}

	funcs.UpdateVehicleMileage(request.TrnRequestUID, request.Mile)
	c.JSON(http.StatusOK, gin.H{"message": "Updated successfully", "data": existing})
}

// DeleteVehicleAddFuel godoc
// @Summary Delete Vehicle Add Fuel entry
// @Description This endpoint allows to mark a Vehicle Add Fuel entry as deleted.
// @Tags Vehicle-in-use-admin
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Param trn_add_fuel_uid path string true "TrnAddFuelUID"
// @Router /api/vehicle-in-use-admin/delete-add-fuel/{trn_add_fuel_uid} [delete]
func (h *VehicleInUseAdminHandler) DeleteVehicleAddFuel(c *gin.Context) {
	user := funcs.GetAuthenUser(c, h.Role)
	if c.IsAborted() {
		return
	}
	uid := c.Param("trn_add_fuel_uid")
	trnAddFuelUid, err := uuid.Parse(uid)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Uid", "message": messages.ErrInvalidUID.Error()})
		return
	}
	var existing models.VmsTrnAddFuel
	if err := config.DB.Where("trn_add_fuel_uid = ? AND is_deleted = ?", trnAddFuelUid, "0").First(&existing).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Add Fuel entry not found", "message": messages.ErrNotfound.Error()})
		return
	}
	var trnRequest models.VmsTrnRequestList
	query := h.SetQueryRole(user, config.DB)
	query = h.SetQueryStatusCanUpdate(query)
	if err := query.Table("public.vms_trn_request").Where("trn_request_uid = ?", existing.TrnRequestUID).First(&trnRequest).Error; err != nil {
		c.JSON(http.StatusMethodNotAllowed, gin.H{"error": "Booking can not update", "message": messages.ErrBookingNotFound.Error()})
		return
	}

	if err := config.DB.Model(&existing).UpdateColumns(map[string]interface{}{
		"is_deleted": "1",
		"updated_by": user.EmpID,
		"updated_at": time.Now(),
	}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete", "message": messages.ErrInternalServer.Error()})
	}

	var maxMileage int
	if err := config.DB.Table("vms_trn_add_fuel").
		Where("trn_request_uid = ?", existing.TrnRequestUID).
		Select("MAX(mile)").Scan(&maxMileage).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get max mileage", "message": messages.ErrInternalServer.Error()})
		return
	}
	funcs.UpdateVehicleMileage(existing.TrnRequestUID, maxMileage)
	c.JSON(http.StatusOK, gin.H{"message": "Deleted successfully"})
}

// GetVehicleAddFuelDetails godoc
// @Summary Retrieve a list of Add Fuel entries
// @Description Fetch a list of Add Fuel entries in TrnRequestUID.
// @Tags Vehicle-in-use-admin
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Param trn_request_uid path string true "TrnRequestUID"
// @Param search query string false "Search keyword (matches tax_invoice_no)"
// @Router /api/vehicle-in-use-admin/add-fuel-details/{trn_request_uid} [get]
func (h *VehicleInUseAdminHandler) GetVehicleAddFuelDetails(c *gin.Context) {
	user := funcs.GetAuthenUser(c, h.Role)
	if c.IsAborted() {
		return
	}
	uid := c.Param("trn_request_uid")
	trnRequestUid, err := uuid.Parse(uid)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Uid", "message": messages.ErrInvalidUID.Error()})
		return
	}

	var trnRequest models.VmsTrnRequestList
	query := h.SetQueryRole(user, config.DB)
	if err := query.Table("public.vms_trn_request").Where("trn_request_uid = ?", trnRequestUid).First(&trnRequest).Error; err != nil {
		c.JSON(http.StatusMethodNotAllowed, gin.H{"error": "Booking can not update", "message": messages.ErrBookingNotFound.Error()})
		return
	}

	queryTrip := config.DB
	queryTrip = queryTrip.Where("trn_request_uid = ? AND is_deleted = ?", trnRequestUid, "0")
	if search := c.Query("search"); search != "" {
		queryTrip = queryTrip.Where("tax_invoice_no ILIKE ?", "%"+search+"%")
	}

	var fuels []models.VmsTrnAddFuel
	queryTrip = queryTrip.
		Preload("RefCostType").
		Preload("RefOilStationBrand").
		Preload("RefFuelType").
		Preload("RefPaymentType")
	if err := queryTrip.Find(&fuels).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Add Fuel entries not found", "message": messages.ErrNotfound.Error()})
		return
	}

	c.JSON(http.StatusOK, fuels)
}

// GetVehicleAddFuelDetail godoc
// @Summary Retrieve details of a specific Add Fuel entry
// @Description Fetch detailed information about an Add Fuel entry using its unique TrnAddFuelUID.
// @Tags Vehicle-in-use-admin
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Param trn_add_fuel_uid path string true "TrnAddFuelUID"
// @Router /api/vehicle-in-use-admin/add-fuel-detail/{trn_add_fuel_uid} [get]
func (h *VehicleInUseAdminHandler) GetVehicleAddFuelDetail(c *gin.Context) {
	user := funcs.GetAuthenUser(c, h.Role)
	if c.IsAborted() {
		return
	}
	uid := c.Param("trn_add_fuel_uid")
	trnAddFuelUid, err := uuid.Parse(uid)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Uid", "message": messages.ErrInvalidUID.Error()})
		return
	}
	var fuel models.VmsTrnAddFuel
	if err := config.DB.
		Preload("RefCostType").
		Preload("RefOilStationBrand").
		Preload("RefFuelType").
		Preload("RefPaymentType").
		Where("trn_add_fuel_uid = ? AND is_deleted = ?", trnAddFuelUid, "0").First(&fuel).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Add Fuel entry not found", "message": messages.ErrNotfound.Error()})
		return
	}
	var trnRequest models.VmsTrnRequestList
	query := h.SetQueryRole(user, config.DB)
	if err := query.Table("public.vms_trn_request").Where("trn_request_uid = ?", fuel.TrnRequestUID).First(&trnRequest).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Add Fuel entries not found", "message": messages.ErrNotfound.Error()})
		return
	}
	c.JSON(http.StatusOK, fuel)
}

// GetTravelCard godoc
// @Summary Retrieve a travel-card of pecific booking request
// @Description This endpoint fetches a travel-card of pecific booking request using its unique identifier (TrnRequestUID).
// @Tags Vehicle-in-use-admin
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Param trn_request_uid path string true "TrnRequestUID (trn_request_uid)"
// @Router /api/vehicle-in-use-admin/travel-card/{trn_request_uid} [get]
func (h *VehicleInUseAdminHandler) GetTravelCard(c *gin.Context) {
	user := funcs.GetAuthenUser(c, h.Role)
	if c.IsAborted() {
		return
	}
	id := c.Param("trn_request_uid")
	trnRequestUid, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Uid"})
		return
	}
	var request models.VmsTrnTravelCard
	query := h.SetQueryRole(user, config.DB)
	query = query.Table("public.vms_trn_request AS req").
		Select("req.*, v.vehicle_license_plate,v.vehicle_license_plate_province_short,v.vehicle_license_plate_province_full").
		Joins("LEFT JOIN vms_mas_vehicle v on v.mas_vehicle_uid = req.mas_vehicle_uid")

	if err := query.
		First(&request, "trn_request_uid = ?", trnRequestUid).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Booking not found"})
		return
	}
	request.VehicleUserImageURL = funcs.GetEmpImage(request.VehicleUserEmpID)
	request.VehicleUserDeptSAPShort = request.VehicleUserPosition + " " + request.VehicleUserDeptSAPShort
	request.ApprovedRequestDeptSAPShort = request.ApprovedRequestPosition + " " + request.ApprovedRequestDeptSAPShort

	c.JSON(http.StatusOK, request)
}

// ReturnedVehicle godoc
// @Summary Update returned vehicle for a booking request
// @Description This endpoint allows to update the returned vehicle details.
// @Tags Vehicle-in-use-admin
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Param data body models.VmsTrnReturnedVehicle true "VmsTrnReturnedVehicle data"
// @Router /api/vehicle-in-use-admin/returned-vehicle [put]
func (h *VehicleInUseAdminHandler) ReturnedVehicle(c *gin.Context) {
	user := funcs.GetAuthenUser(c, h.Role)
	if c.IsAborted() {
		return
	}
	var request, trnRequest models.VmsTrnReturnedVehicle
	var result struct {
		models.VmsTrnReturnedVehicle
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

	request.RefRequestStatusCode = "70"
	request.UpdatedAt = time.Now()
	request.UpdatedBy = user.EmpID

	empUser := funcs.GetUserEmpInfo(user.EmpID)
	request.ReturnedVehicleEmpID = empUser.EmpID
	request.ReturnedVehicleEmpName = empUser.FullName
	request.ReturnedVehicleDeptSAP = empUser.DeptSAP
	request.ReturnedVehicleDeptNameShort = empUser.DeptSAPShort
	request.ReturnedVehicleDeptNameFull = empUser.DeptSAPFull
	for i := range request.VehicleImages {
		request.VehicleImages[i].TrnVehicleImgReturnedUID = uuid.New().String()
		request.VehicleImages[i].TrnRequestUID = request.TrnRequestUID
		request.VehicleImages[i].CreatedAt = time.Now()
		request.VehicleImages[i].CreatedBy = user.EmpID
		request.VehicleImages[i].UpdatedAt = time.Now()
		request.VehicleImages[i].UpdatedBy = user.EmpID
		request.VehicleImages[i].IsDeleted = "0"
	}

	if len(request.VehicleImages) > 0 {
		if err := config.DB.Where("trn_request_uid = ?", request.TrnRequestUID).Delete(&models.VehicleImageReturned{}).Error; err != nil {
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
		funcs.CreateTrnRequestActionLog(request.TrnRequestUID,
			request.RefRequestStatusCode,
			"รอผู้ดูแลยานพาหนะตรวจสอบ",
			user.EmpID,
			"admin-department",
			"",
		)
	}
	funcs.UpdateVehicleMileage(request.TrnRequestUID, request.MileEnd)
	funcs.UpdateVehicleParkingPlace(request.TrnRequestUID, request.ReturnedParkingPlace)
	c.JSON(http.StatusOK, gin.H{"message": "Updated successfully", "result": result})
}

// UpdateReceivedVehicle godoc
// @Summary Update vehicle pickup for a booking request
// @Description This endpoint allows to update vehicle pickup.
// @Tags Vehicle-in-use-admin
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Param data body models.VmsTrnReceivedVehicleNoImgage true "VmsTrnReceivedVehicleNoImgage data"
// @Router /api/vehicle-in-use-admin/update-received-vehicle [put]
func (h *VehicleInUseAdminHandler) UpdateReceivedVehicle(c *gin.Context) {
	user := funcs.GetAuthenUser(c, h.Role)
	if c.IsAborted() {
		return
	}
	var request, trnRequest models.VmsTrnReceivedVehicleNoImgage
	var result struct {
		models.VmsTrnReceivedVehicleNoImgage
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to update : %v", err)})
		return
	}

	if err := config.DB.
		First(&result, "trn_request_uid = ?", request.TrnRequestUID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Booking not found", "message": messages.ErrBookingNotFound.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Updated successfully", "result": result})
}

// UpdateReceivedVehicleImages godoc
// @Summary Update vehicle pickup for a booking request
// @Description This endpoint allows to update vehicle pickup.
// @Tags Vehicle-in-use-admin
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Param data body models.VmsTrnReceivedVehicleImages true "VmsTrnReceivedVehicleImages data"
// @Router /api/vehicle-in-use-admin/update-received-vehicle-images [put]
func (h *VehicleInUseAdminHandler) UpdateReceivedVehicleImages(c *gin.Context) {
	user := funcs.GetAuthenUser(c, h.Role)
	if c.IsAborted() {
		return
	}
	var request, trnRequest models.VmsTrnReceivedVehicleImages
	var result struct {
		models.VmsTrnReceivedVehicleImages
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

	c.JSON(http.StatusOK, gin.H{"message": "Updated successfully", "result": result})
}
