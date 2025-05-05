package handlers

import (
	"fmt"
	"net/http"
	"sort"
	"strings"
	"time"
	"vms_plus_be/config"
	"vms_plus_be/funcs"
	"vms_plus_be/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type VehicleInspectionAdminHandler struct {
	Role string
}

var StatusNameMapVehicelInspectionAdmin = map[string]string{
	"70": "รอตรวจสอบ",
	"71": "ตีกลับยานพาหนะ",
}

// SearchRequests godoc
// @Summary Search booking requests and get summary counts by request status code
// @Description Search for requests using a keyword and get the summary of counts grouped by request status code
// @Tags Vehicle-inspection-admin
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
// @Router /api/vehicle-inspection-admin/search-requests [get]
func (h *VehicleInspectionAdminHandler) SearchRequests(c *gin.Context) {
	funcs.GetAuthenUser(c, h.Role)
	if c.IsAborted() {
		return
	}
	statusNameMap := StatusNameMapVehicelInspectionAdmin
	var requests []models.VmsTrnRequestVehicleInUseList
	var summary []models.VmsTrnRequestSummary

	// Use the keys from statusNameMap as the list of valid status codes
	statusCodes := make([]string, 0, len(statusNameMap))
	for code := range statusNameMap {
		statusCodes = append(statusCodes, code)
	}

	// Build the main query
	query := config.DB.Table("public.vms_trn_request AS req").
		Select("req.*, status.ref_request_status_desc,"+
			"(select parking_place from vms_mas_vehicle_department d where d.mas_vehicle_uid::text = req.mas_vehicle_uid) parking_place ").
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
// @Tags Vehicle-inspection-admin
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Param trn_request_uid path string true "TrnRequestUID (trn_request_uid)"
// @Router /api/vehicle-inspection-admin/request/{trn_request_uid} [get]
func (h *VehicleInspectionAdminHandler) GetRequest(c *gin.Context) {
	funcs.GetAuthenUser(c, h.Role)
	if c.IsAborted() {
		return
	}
	request, err := funcs.GetRequestVehicelInUse(c, StatusNameMapVehicelInspectionAdmin)
	if err != nil {
		return
	}
	c.JSON(http.StatusOK, request)
}

// CreateVehicleTripDetail godoc
// @Summary Create Vehicle Travel List for a booking request
// @Description This endpoint allows to Create Vehicle Travel List.
// @Tags Vehicle-inspection-admin
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Param data body models.VmsTrnTripDetailRequest true "VmsTrnTripDetailRequest data"
// @Router /api/vehicle-inspection-admin/create-travel-detail [post]
func (h *VehicleInspectionAdminHandler) CreateVehicleTripDetail(c *gin.Context) {
	user := funcs.GetAuthenUser(c, h.Role)
	if c.IsAborted() {
		return
	}

	var request models.VmsTrnTripDetail
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON input"})
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
	if err := config.DB.Table("public.vms_trn_request").Where("trn_request_uid = ?", request.TrnRequestUID).First(&trnRequest).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Booking not found"})
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create"})
		return
	}

	// Return success response
	c.JSON(http.StatusCreated, gin.H{"message": "Created successfully", "data": request})
}

// UpdateVehicleTripDetail godoc
// @Summary Update Vehicle Travel List for a booking request
// @Description This endpoint allows to Update Vehicle Travel List.
// @Tags Vehicle-inspection-admin
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Param trn_trip_detail_uid path string true "TrnTripDetailUID"
// @Param data body models.VmsTrnTripDetailRequest true "VmsTrnTripDetailRequest data"
// @Router /api/vehicle-inspection-admin/update-travel-detail/{trn_trip_detail_uid} [put]
func (h *VehicleInspectionAdminHandler) UpdateVehicleTripDetail(c *gin.Context) {
	user := funcs.GetAuthenUser(c, h.Role)
	if c.IsAborted() {
		return
	}
	uid := c.Param("trn_trip_detail_uid")
	trnTripDetailUid, err := uuid.Parse(uid)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Uid"})
		return
	}
	var existing models.VmsTrnTripDetail
	if err := config.DB.Where("trn_trip_detail_uid = ? AND is_deleted = ?", trnTripDetailUid, "0").First(&existing).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Trip not found"})
		return
	}
	var request models.VmsTrnTripDetailRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	existing.VmsTrnTripDetailRequest = request
	existing.UpdatedBy = user.EmpID
	existing.UpdatedAt = time.Now()

	if err := config.DB.Save(&existing).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to update : %v", err)})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Updated successfully", "data": request})
}

// DeleteVehicleTripDetail godoc
// @Summary Update Vehicle Travel List for a booking request
// @Description This endpoint allows to Update Vehicle Travel List.
// @Tags Vehicle-inspection-admin
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Param trn_trip_detail_uid path string true "TrnTripDetailUID"
// @Router /api/vehicle-inspection-admin/delete-travel-detail/{trn_trip_detail_uid} [delete]
func (h *VehicleInspectionAdminHandler) DeleteVehicleTripDetail(c *gin.Context) {
	user := funcs.GetAuthenUser(c, h.Role)
	if c.IsAborted() {
		return
	}
	uid := c.Param("trn_trip_detail_uid")
	trnTripDetailUid, err := uuid.Parse(uid)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Uid"})
		return
	}

	var existing models.VmsTrnTripDetail
	if err := config.DB.Where("trn_trip_detail_uid = ? AND is_deleted = ?", trnTripDetailUid, "0").First(&existing).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Trip not found"})
		return
	}

	if err := config.DB.Model(&existing).UpdateColumns(map[string]interface{}{
		"is_deleted": "1",
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
// @Tags Vehicle-inspection-admin
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Param trn_request_uid path string true "TrnRequestUID (trn_request_uid)"
// @Param search query string false "Search keyword (matches place)"
// @Router /api/vehicle-inspection-admin/travel-details/{trn_request_uid} [get]
func (h *VehicleInspectionAdminHandler) GetVehicleTripDetails(c *gin.Context) {
	funcs.GetAuthenUser(c, h.Role)
	if c.IsAborted() {
		return
	}
	uid := c.Param("trn_request_uid")
	trnRequestUid, err := uuid.Parse(uid)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Uid"})
		return
	}
	query := config.DB
	if search := c.Query("search"); search != "" {
		query = query.Where("trip_departure_place ILIKE ? OR trip_destination_place ILIKE ?", "%"+search+"%", "%"+search+"%")
	}

	// Fetch the vehicle record from the database
	var trips []models.VmsTrnTripDetailList
	if err := query.Find(&trips, "trn_request_uid = ? AND is_deleted = ?", trnRequestUid, "0").Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Trip not found"})
		return
	}

	// Return the vehicle data as a JSON response
	c.JSON(http.StatusOK, trips)
}

// GetVehicleTripDetail godoc
// @Summary Retrieve details of a specific trip detail
// @Description Fetch detailed information about a trip detail using their unique TrnTripDetailUID.
// @Tags Vehicle-inspection-admin
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Param trn_trip_detail_uid path string true "TrnTripDetailUID"
// @Router /api/vehicle-inspection/travel-detail/{trn_trip_detail_uid} [get]
func (h *VehicleInspectionAdminHandler) GetVehicleTripDetail(c *gin.Context) {
	funcs.GetAuthenUser(c, h.Role)
	if c.IsAborted() {
		return
	}
	uid := c.Param("trn_trip_detail_uid")
	trnTripDetailUid, err := uuid.Parse(uid)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Uid"})
		return
	}
	// Fetch the vehicle record from the database
	var trip models.VmsTrnTripDetail
	if err := config.DB.First(&trip, "trn_trip_detail_uid = ? AND is_deleted = ?", trnTripDetailUid, "0").Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Trip not found"})
		return
	}

	// Return the vehicle data as a JSON response
	c.JSON(http.StatusOK, trip)
}

// CreateVehicleAddFuel godoc
// @Summary Create Vehicle Add Fuel entry
// @Description This endpoint allows to create a new Vehicle Add Fuel entry.
// @Tags Vehicle-inspection-admin
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Param data body models.VmsTrnAddFuelRequest true "VmsTrnAddFuelRequest data"
// @Router /api/vehicle-inspection-admin/create-add-fuel [post]
func (h *VehicleInspectionAdminHandler) CreateVehicleAddFuel(c *gin.Context) {
	user := funcs.GetAuthenUser(c, h.Role)
	if c.IsAborted() {
		return
	}

	var request models.VmsTrnAddFuel
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON input"})
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
	if err := config.DB.Table("public.vms_trn_request").Where("trn_request_uid = ?", request.TrnRequestUID).First(&trnRequest).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Booking not found"})
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Created successfully", "data": request})
}

// UpdateVehicleAddFuel godoc
// @Summary Update Vehicle Add Fuel entry
// @Description This endpoint allows to update an existing Vehicle Add Fuel entry.
// @Tags Vehicle-inspection-admin
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Param trn_add_fuel_uid path string true "TrnAddFuelUID"
// @Param data body models.VmsTrnAddFuel true "VmsTrnAddFuel data"
// @Router /api/vehicle-inspection-admin/update-add-fuel/{trn_add_fuel_uid} [put]
func (h *VehicleInspectionAdminHandler) UpdateVehicleAddFuel(c *gin.Context) {
	user := funcs.GetAuthenUser(c, h.Role)
	if c.IsAborted() {
		return
	}
	uid := c.Param("trn_add_fuel_uid")
	trnAddFuelUid, err := uuid.Parse(uid)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Uid"})
		return
	}

	var existing models.VmsTrnAddFuel
	if err := config.DB.Where("trn_add_fuel_uid = ? AND is_deleted = ?", trnAddFuelUid, "0").First(&existing).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Add Fuel entry not found"})
		return
	}

	var request models.VmsTrnAddFuelRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON input"})
		return
	}

	existing.VmsTrnAddFuelRequest = request
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
// @Tags Vehicle-inspection-admin
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Param trn_add_fuel_uid path string true "TrnAddFuelUID"
// @Router /api/vehicle-inspection-admin/delete-add-fuel/{trn_add_fuel_uid} [delete]
func (h *VehicleInspectionAdminHandler) DeleteVehicleAddFuel(c *gin.Context) {
	user := funcs.GetAuthenUser(c, h.Role)
	if c.IsAborted() {
		return
	}
	uid := c.Param("trn_add_fuel_uid")
	trnAddFuelUid, err := uuid.Parse(uid)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Uid"})
		return
	}
	var existing models.VmsTrnAddFuel
	if err := config.DB.Where("trn_add_fuel_uid = ? AND is_deleted = ?", trnAddFuelUid, "0").First(&existing).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Add Fuel entry not found"})
		return
	}

	if err := config.DB.Model(&existing).UpdateColumns(map[string]interface{}{
		"is_deleted": "1",
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
// @Tags Vehicle-inspection-admin
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Param trn_request_uid path string true "TrnRequestUID"
// @Param search query string false "Search keyword (matches tax_invoice_no)"
// @Router /api/vehicle-inspection-admin/add-fuel-details/{trn_request_uid} [get]
func (h *VehicleInspectionAdminHandler) GetVehicleAddFuelDetails(c *gin.Context) {
	funcs.GetAuthenUser(c, h.Role)
	if c.IsAborted() {
		return
	}
	uid := c.Param("trn_request_uid")
	trnRequestUid, err := uuid.Parse(uid)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Uid"})
		return
	}

	query := config.DB.Where("trn_request_uid = ? AND is_deleted = ?", trnRequestUid, "0")
	if search := c.Query("search"); search != "" {
		query = query.Where("tax_invoice_no ILIKE ?", "%"+search+"%")
	}

	var fuels []models.VmsTrnAddFuel
	query = query.
		Preload("RefContType").
		Preload("RefOilStationBrand").
		Preload("RefFuelType").
		Preload("RefPaymentType")
	if err := query.Find(&fuels).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Add Fuel entries not found"})
		return
	}

	c.JSON(http.StatusOK, fuels)
}

// GetVehicleAddFuelDetail godoc
// @Summary Retrieve details of a specific Add Fuel entry
// @Description Fetch detailed information about an Add Fuel entry using its unique TrnAddFuelUID.
// @Tags Vehicle-inspection-admin
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Param trn_add_fuel_uid path string true "TrnAddFuelUID"
// @Router /api/vehicle-inspection-admin/add-fuel-detail/{trn_add_fuel_uid} [get]
func (h *VehicleInspectionAdminHandler) GetVehicleAddFuelDetail(c *gin.Context) {
	funcs.GetAuthenUser(c, h.Role)
	if c.IsAborted() {
		return
	}
	uid := c.Param("trn_add_fuel_uid")
	trnAddFuelUid, err := uuid.Parse(uid)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Uid"})
		return
	}
	var fuel models.VmsTrnAddFuel
	if err := config.DB.
		Preload("RefContType").
		Preload("RefOilStationBrand").
		Preload("RefFuelType").
		Preload("RefPaymentType").
		Where("trn_add_fuel_uid = ? AND is_deleted = ?", trnAddFuelUid, "0").First(&fuel).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Add Fuel entry not found"})
		return
	}

	c.JSON(http.StatusOK, fuel)
}

// GetTravelCard godoc
// @Summary Retrieve a travel-card of pecific booking request
// @Description This endpoint fetches a travel-card of pecific booking request using its unique identifier (TrnRequestUID).
// @Tags Vehicle-inspection-admin
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Param trn_request_uid path string true "TrnRequestUID (trn_request_uid)"
// @Router /api/vehicle-inspection-admin/travel-card/{trn_request_uid} [get]
func (h *VehicleInspectionAdminHandler) GetTravelCard(c *gin.Context) {
	funcs.GetAuthenUser(c, h.Role)
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
	if err := config.DB.
		First(&request, "trn_request_uid = ?", trnRequestUid).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Booking not found"})
		return
	}
	request.VehicleUserImageURL = config.DefaultAvatarURL
	c.JSON(http.StatusOK, request)
}

// UpdateReturnedVehicle godoc
// @Summary Update returned vehicle for a booking request
// @Description This endpoint allows to update the returned vehicle details.
// @Tags Vehicle-inspection-admin
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Param data body models.VmsTrnReturnedVehicleNoImage true "VmsTrnReturnedVehicleNoImage data"
// @Router /api/vehicle-inspection-admin/update-returned-vehicle [put]
func (h *VehicleInspectionAdminHandler) UpdateReturnedVehicle(c *gin.Context) {
	user := funcs.GetAuthenUser(c, h.Role)
	if c.IsAborted() {
		return
	}
	var request, trnRequest models.VmsTrnReturnedVehicleNoImage
	var result struct {
		models.VmsTrnReturnedVehicleNoImage
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

	request.UpdatedAt = time.Now()
	request.UpdatedBy = user.EmpID

	if err := config.DB.Save(&request).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to update : %v", err)})
		return
	}

	if err := config.DB.
		First(&result, "trn_request_uid = ?", request.TrnRequestUID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Booking not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Updated successfully", "result": result})
}

// UpdateReturnedVehicleImages godoc
// @Summary Update vehicle return for a booking request
// @Description This endpoint allows to update vehicle return.
// @Tags Vehicle-inspection-admin
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Param data body models.VmsTrnReturnedVehicleImages true "VmsTrnReturnedVehicleImages data"
// @Router /api/vehicle-inspection-admin/update-returned-vehicle-images [put]
func (h *VehicleInspectionAdminHandler) UpdateReturnedVehicleImages(c *gin.Context) {
	user := funcs.GetAuthenUser(c, h.Role)
	if c.IsAborted() {
		return
	}
	var request, trnRequest models.VmsTrnReturnedVehicleImages
	var result struct {
		models.VmsTrnReturnedVehicleImages
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

	request.UpdatedAt = time.Now()
	request.UpdatedBy = user.EmpID

	for i := range request.VehicleImages {
		request.VehicleImages[i].TrnVehicleImgReturnedUID = uuid.New().String()
		request.VehicleImages[i].TrnRequestUID = request.TrnRequestUID
	}

	if len(request.VehicleImages) > 0 {
		if err := config.DB.Where("trn_request_uid = ?", request.TrnRequestUID).Delete(&models.VehicleImageReturned{}).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update vehicle images"})
			return
		}
	}

	if err := config.DB.Save(&request).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to update : %v", err)})
		return
	}

	if err := config.DB.
		Preload("VehicleImages").
		First(&result, "trn_request_uid = ?", request.TrnRequestUID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Booking not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Updated successfully", "result": result})
}

// GetSatisfactionSurvey godoc
// @Summary Retrieve Satisfaction Survey for a booking request
// @Description This endpoint fetches Satisfaction Survey for a booking request.
// @Tags Vehicle-inspection-admin
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Param trn_request_uid path string true "TrnRequestUID"
// @Router /api/vehicle-inspection-admin/satisfaction-survey/{trn_request_uid} [get]
func (h *VehicleInspectionAdminHandler) GetSatisfactionSurvey(c *gin.Context) {
	funcs.GetAuthenUser(c, h.Role)
	if c.IsAborted() {
		return
	}
	uid := c.Param("trn_request_uid")
	trnRequestUID, err := uuid.Parse(uid)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Uid"})
		return
	}

	var list []models.VmsTrnSatisfactionSurveyAnswersResponse
	if err := config.DB.
		Preload("SatisfactionSurveyQuestions").
		Where("trn_request_uid = ?", trnRequestUID).Find(&list).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Satisfaction survey not found"})
		return
	}
	for i := range list {
		desc := list[i].SatisfactionSurveyQuestions.MasSatisfactionSurveyQuestionsDesc
		parts := strings.SplitN(desc, ":", 2)
		list[i].SatisfactionSurveyQuestions.MasSatisfactionSurveyQuestionsTitle = parts[0] // Title before colon
		if len(parts) > 1 {
			list[i].SatisfactionSurveyQuestions.MasSatisfactionSurveyQuestionsDesc = parts[1] // Remaining description after colon
		} else {
			list[i].SatisfactionSurveyQuestions.MasSatisfactionSurveyQuestionsDesc = "" // Empty if no colon found
		}
	}
	c.JSON(http.StatusOK, list)
}

// UpdateSendedBack godoc
// @Summary Update booking request to "Sended Back" status
// @Description This endpoint allows to update the booking request status to "Sended Back"
// @Tags Vehicle-inspection-admin
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Param data body models.VmsTrnRequestSendedBack true "VmsTrnRequestSendedBack data"
// @Router /api/vehicle-inspection-admin/update-sended-back [put]
func (h *VehicleInspectionAdminHandler) UpdateSendedBack(c *gin.Context) {
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

	request.RefRequestStatusCode = "71" //ตรวจสอบยานพาหนะไม่ผ่าน รอแก้ไขเพื่อส่งคืน
	empUser := funcs.GetUserEmpInfo(user.EmpID)
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

// UpdateAccepted godoc
// @Summary Update booking request to "Accepted" status
// @Description This endpoint allows to update the booking request status to "Accepted"
// @Tags Vehicle-inspection-admin
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Param data body models.VmsTrnRequestAccepted true "VmsTrnRequestAccepted data"
// @Router /api/vehicle-inspection-admin/update-accepted [put]
func (h *VehicleInspectionAdminHandler) UpdateAccepted(c *gin.Context) {
	user := funcs.GetAuthenUser(c, h.Role)
	if c.IsAborted() {
		return
	}
	var request, trnRequest models.VmsTrnRequestAccepted
	var result struct {
		models.VmsTrnRequestAccepted
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

	request.RefRequestStatusCode = "80" // ยืนยันการยอมรับการส่งคืน
	empUser := funcs.GetUserEmpInfo(user.EmpID)
	request.AcceptedVehicleEmpID = empUser.EmpID
	request.AcceptedVehicleEmpName = empUser.FullName
	request.AcceptedVehicleDeptSAP = empUser.DeptSAP
	request.AcceptedVehicleDeptSAPShort = empUser.DeptSAPShort
	request.AcceptedVehicleDeptSAPFull = empUser.DeptSAPFull
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
		"รับกุญแจและยานพาหนะคืนแล้ว จบคำขอ",
		user.EmpID)

	c.JSON(http.StatusOK, gin.H{"message": "Updated successfully", "result": result})
}

// UpdateInspectVehicleImages godoc
// @Summary Update vehicle inspect for a booking request
// @Description This endpoint allows to update vehicle inspect.
// @Tags Vehicle-inspection-admin
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Param data body models.VmsTrnInspectVehicleImages true "VmsTrnInspectVehicleImages data"
// @Router /api/vehicle-inspection-admin/update-inspect-vehicle-images [put]
func (h *VehicleInspectionAdminHandler) UpdateInspectVehicleImages(c *gin.Context) {
	user := funcs.GetAuthenUser(c, h.Role)
	if c.IsAborted() {
		return
	}
	var request, trnRequest models.VmsTrnInspectVehicleImages
	var result struct {
		models.VmsTrnInspectVehicleImages
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

	request.UpdatedAt = time.Now()
	request.UpdatedBy = user.EmpID

	for i := range request.VehicleImages {
		request.VehicleImages[i].TrnVehicleImgReturnedUID = uuid.New().String()
		request.VehicleImages[i].TrnRequestUID = request.TrnRequestUID
	}

	if len(request.VehicleImages) > 0 {
		if err := config.DB.Where("trn_request_uid = ?", request.TrnRequestUID).Delete(&models.VehicleImageReturned{}).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update vehicle images"})
			return
		}
	}

	if err := config.DB.Save(&request).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to update : %v", err)})
		return
	}

	if err := config.DB.
		Preload("VehicleImages").
		First(&result, "trn_request_uid = ?", request.TrnRequestUID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Booking not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Updated successfully", "result": result})
}
