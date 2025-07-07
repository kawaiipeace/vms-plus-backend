package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
	"vms_plus_be/config"
	"vms_plus_be/funcs"
	"vms_plus_be/messages"
	"vms_plus_be/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// SearchCarpoolVehicle godoc
// @Summary Search carpool vehicles
// @Description Search carpool vehicles with pagination and filters
// @Tags Carpool-management
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Param mas_carpool_uid path string true "MasCarpoolUID (mas_carpool_uid)"
// @Param search query string false "Search query for vehicle_no or vehicle_owner_dept_short"
// @Param is_active query string false "Filter by is_active status (comma-separated, e.g., '1,0')"
// @Param order_by query string false "Order by fields: vehicle_license_plate"
// @Param order_dir query string false "Order direction: asc or desc"
// @Param page query int false "Page number (default: 1)"
// @Param limit query int false "Number of records per page (default: 10)"
// @Router /api/carpool-management/vehicle-search/{mas_carpool_uid} [get]
func (h *CarpoolManagementHandler) SearchCarpoolVehicle(c *gin.Context) {
	user := funcs.GetAuthenUser(c, h.Role)
	if c.IsAborted() {
		return
	}
	masCarpoolUID := c.Param("mas_carpool_uid")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))    // Default: page 1
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10")) // Default: 10 items per page
	offset := (page - 1) * limit

	var existingCarpool models.VmsMasCarpoolRequest
	queryRole := h.SetQueryRole(user, config.DB)
	if err := queryRole.Where("mas_carpool_uid = ? AND is_deleted = ?", masCarpoolUID, "0").First(&existingCarpool).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Carpool not found", "message": messages.ErrNotfound.Error()})
		return
	}

	var vehicles []models.VmsMasCarpoolVehicleList
	query := config.DB.Table("vms_mas_carpool_vehicle cpv").
		Select(
			`cpv.mas_carpool_vehicle_uid,
			cpv.mas_carpool_uid,
			cpv.mas_vehicle_uid,
			v.vehicle_license_plate,
			v.vehicle_license_plate_province_short,
			v.vehicle_brand_name,
			v.vehicle_model_name,
			v.ref_vehicle_type_code,
			(select max(ref_vehicle_type_name) from vms_ref_vehicle_type s where s.ref_vehicle_type_code=v.ref_vehicle_type_code) ref_vehicle_type_name,
			(select max(s.dept_long_short) from vms_mas_department s where s.dept_sap=d.vehicle_owner_dept_sap) vehicle_owner_dept_short,
			v.ref_vehicle_type_code,
			d.fleet_card_no,
			v.is_tax_credit,
			d.vehicle_mileage,
			d.vehicle_get_date,
			v.seat,
			v.vehicle_color,
			v.vehicle_gear,
			d.vehicle_pea_id,
			d.ref_vehicle_status_code,
			(select max(s.ref_vehicle_status_short_name) from vms_ref_vehicle_status s where s.ref_vehicle_status_code=d.ref_vehicle_status_code) ref_vehicle_status_name,
			v.ref_fuel_type_id,
			(select max(s.ref_fuel_type_name_th) from vms_ref_fuel_type s where s.ref_fuel_type_id=v.ref_fuel_type_id) fuel_type_name,
			cpv.is_active,
			d.parking_place
		`).
		Joins("LEFT JOIN vms_mas_vehicle v ON v.mas_vehicle_uid = cpv.mas_vehicle_uid").
		Joins("LEFT JOIN (SELECT DISTINCT ON (mas_vehicle_uid) * FROM vms_mas_vehicle_department WHERE is_deleted = '0' AND is_active = '1' ORDER BY mas_vehicle_uid, created_at DESC) d ON v.mas_vehicle_uid = d.mas_vehicle_uid").
		Where("cpv.mas_carpool_uid = ? AND cpv.is_deleted = ?", masCarpoolUID, "0")

	search := strings.ToUpper(c.Query("search"))
	if search != "" {
		query = query.Where("UPPER(v.vehicle_no) ILIKE ? OR UPPER(v.vehicle_name) ILIKE ?", "%"+search+"%", "%"+search+"%")
	}

	if isActive := c.Query("is_active"); isActive != "" {
		isActiveList := strings.Split(isActive, ",")
		query = query.Where("cpv.is_active IN (?)", isActiveList)
	}

	orderBy := c.Query("order_by")
	orderDir := c.Query("order_dir")
	if orderDir != "desc" {
		orderDir = "asc"
	}
	switch orderBy {
	case "vehicle_no":
		query = query.Order("v.vehicle_license_plate " + orderDir)
	case "vehicle_owner_dept_short":
		query = query.Order("v.vehicle_owner_dept_short " + orderDir)
	default:
		query = query.Order("v.vehicle_license_plate " + orderDir) // Default ordering by vehicle_no
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "message": messages.ErrInternalServer.Error()})
		return
	}

	query = query.Limit(limit).Offset(offset)
	if err := query.Find(&vehicles).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "message": messages.ErrInternalServer.Error()})
		return
	}
	for i := range vehicles {
		vehicles[i].Age = funcs.CalculateAge(vehicles[i].VehicleGetDate)
		funcs.TrimStringFields(&vehicles[i])
		var vehicleImgs []models.VmsMasVehicleImg
		if err := config.DB.Where("mas_vehicle_uid = ?", vehicles[i].MasVehicleUID).Find(&vehicleImgs).Error; err == nil {
			vehicles[i].VehicleImgs = make([]string, 0)
			for _, img := range vehicleImgs {
				vehicles[i].VehicleImgs = append(vehicles[i].VehicleImgs, img.VehicleImgFile)
			}
		}
	}
	if len(vehicles) == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "No carpool vehicles found",
			"pagination": gin.H{
				"page":       page,
				"limit":      limit,
				"totalPages": (total + int64(limit) - 1) / int64(limit),
			},
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"pagination": gin.H{
				"total":      total,
				"page":       page,
				"limit":      limit,
				"totalPages": (total + int64(limit) - 1) / int64(limit),
			},
			"vehicles": vehicles,
		})
	}
}

// CreateCarpoolVehicle godoc
// @Summary Create a new carpool vehicle
// @Description Create a new carpool vehicle and save it to the database
// @Tags Carpool-management
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Param vehicle body []models.VmsMasCarpoolVehicle true "VmsMasCarpoolVehicle array"
// @Router /api/carpool-management/vehicle-create [post]
func (h *CarpoolManagementHandler) CreateCarpoolVehicle(c *gin.Context) {
	user := funcs.GetAuthenUser(c, h.Role)
	if c.IsAborted() {
		return
	}

	var requests []models.VmsMasCarpoolVehicle
	if err := c.ShouldBindJSON(&requests); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "message": messages.ErrInvalidJSONInput.Error()})
		return
	}
	if len(requests) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No carpool vehicle data provided", "message": messages.ErrInvalidJSONInput.Error()})
		return

	}
	var existingCarpool models.VmsMasCarpoolRequest
	for i := range requests {
		queryRole := h.SetQueryRole(user, config.DB)
		if err := queryRole.Where("mas_carpool_uid = ? AND is_deleted = ?", requests[i].MasCarpoolUID, "0").First(&existingCarpool).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Carpool not found", "message": messages.ErrNotfound.Error()})
			return
		}
	}
	//check vehicle is exist
	for i := range requests {
		var existingVehicle models.VmsMasVehicle
		if err := config.DB.Where("mas_vehicle_uid = ? AND is_deleted = ?", requests[i].MasVehicleUID, "0").First(&existingVehicle).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   fmt.Sprintf("Vehicle with MasVehicleUID %s not found", requests[i].MasVehicleUID),
				"message": "ข้อมูลรถที่ระบุไม่มีอยู่ในระบบ",
			})
			return
		}
	}

	//check vehicle is not duplicate in another carpool
	for i := range requests {
		var existingVehicle models.VmsMasCarpoolVehicle
		if err := config.DB.Where("mas_vehicle_uid = ? AND is_deleted = ?", requests[i].MasVehicleUID, requests[i].MasCarpoolUID, "0").First(&existingVehicle).Error; err == nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   fmt.Sprintf("Vehicle with MasVehicleUID %s already exists in another carpool", requests[i].MasVehicleUID),
				"message": "ข้อมูลรถที่ระบุมีอยู่ในกลุ่มรถอื่นแล้ว",
			})
			return
		}

		requests[i].MasCarpoolVehicleUID = uuid.New().String()
		requests[i].CreatedAt = time.Now()
		requests[i].CreatedBy = user.EmpID
		requests[i].UpdatedAt = time.Now()
		requests[i].UpdatedBy = user.EmpID
		requests[i].IsDeleted = "0"
		requests[i].IsActive = existingCarpool.IsActive
		requests[i].StartDate = time.Now()
		requests[i].EndDate = time.Now().AddDate(10, 0, 0)
		var masVehicleDepartmentUID string
		if err := config.DB.Table("vms_mas_vehicle_department").Where("mas_vehicle_uid = ?", requests[i].MasVehicleUID).Pluck("mas_vehicle_department_uid", &masVehicleDepartmentUID).Error; err == nil {
			requests[i].MasVehicleDepartmentUID = masVehicleDepartmentUID
		}
	}

	if err := config.DB.Create(&requests).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "message": messages.ErrInternalServer.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":      "Carpool vehicles created successfully",
		"data":         requests,
		"carpool_name": GetCarpoolName(requests[0].MasCarpoolUID),
	})
}

// DeleteCarpoolVehicle godoc
// @Summary Delete a carpool vehicle
// @Description This endpoint deletes a carpool vehicle using its unique identifier (MasCarpoolVehicleUID).
// @Tags Carpool-management
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Param mas_carpool_vehicle_uid path string true "MasCarpoolVehicleUID (mas_carpool_vehicle_uid)"
// @Router /api/carpool-management/vehicle-delete/{mas_carpool_vehicle_uid} [delete]
func (h *CarpoolManagementHandler) DeleteCarpoolVehicle(c *gin.Context) {
	user := funcs.GetAuthenUser(c, h.Role)
	if c.IsAborted() {
		return
	}
	masCarpoolVehicleUID := c.Param("mas_carpool_vehicle_uid")

	var vehicle models.VmsMasCarpoolVehicle
	if err := config.DB.Where("mas_carpool_vehicle_uid = ? AND is_deleted = ?", masCarpoolVehicleUID, "0").First(&vehicle).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Carpool vehicle not found", "message": messages.ErrNotfound.Error()})
		return
	}
	var existingCarpool models.VmsMasCarpoolRequest
	queryRole := h.SetQueryRole(user, config.DB)
	if err := queryRole.Where("mas_carpool_uid = ? AND is_deleted = ?", vehicle.MasCarpoolUID, "0").First(&existingCarpool).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Carpool not found", "message": messages.ErrNotfound.Error()})
		return
	}
	if err := config.DB.Model(&vehicle).UpdateColumns(map[string]interface{}{
		"is_active":  "0",
		"is_deleted": "1",
		"updated_by": user.EmpID,
		"updated_at": time.Now(),
	}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete carpool vehicle", "message": messages.ErrInternalServer.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Carpool vehicle deleted successfully", "carpool_name": GetCarpoolName(vehicle.MasCarpoolUID)})
}

// SearchMasVehicles godoc
// @Summary Search vehicles for add to Carpool vehicle
// @Description Search vehicles for add to Carpool vehicle
// @Tags Carpool-management
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Param search query string false "Search text (Vehicle Brand Name or License Plate)"
// @Param page query int false "Page number (default: 1)"
// @Param limit query int false "Number of records per page (default: 10)"
// @Router /api/carpool-management/vehicle-mas-search [get]
func (h *CarpoolManagementHandler) SearchMasVehicles(c *gin.Context) {
	user := funcs.GetAuthenUser(c, h.Role)
	if c.IsAborted() {
		return
	}
	searchText := c.Query("search") // Text search for brand name & license plate

	// Pagination parameters
	//page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))    // Default page = 1
	//limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10")) // Default limit = 10
	limit := 100
	page := 1

	offset := (page - 1) * limit // Calculate offset

	var vehicles []models.VmsMasVehicleCarpoolList
	var total int64
	query := h.SetQueryRoleDept(funcs.GetAuthenUser(c, h.Role), config.DB)
	query = query.Table("vms_mas_vehicle v")
	query = query.Select("v.*,d.vehicle_owner_dept_sap vehicle_owner_dept_sap,fn_get_long_short_dept_name_by_dept_sap(d.vehicle_owner_dept_sap) vehicle_owner_dept_short" +
		",d.fleet_card_no,d.vehicle_mileage vehicle_mileage" +
		",(select max(s.ref_vehicle_status_short_name) from vms_ref_vehicle_status s where s.ref_vehicle_status_code=d.ref_vehicle_status_code) ref_vehicle_status_name" +
		",(select max(s.ref_fuel_type_name_th) from vms_ref_fuel_type s where s.ref_fuel_type_id=v.ref_fuel_type_id) fuel_type_name")
	query = query.Model(&models.VmsMasVehicleCarpoolList{})
	query = query.Where("v.is_deleted = '0'")
	query = query.Joins("INNER JOIN public.vms_mas_vehicle_department AS d ON v.mas_vehicle_uid = d.mas_vehicle_uid")
	query = query.Where("not exists (select 1 from vms_mas_carpool_vehicle cv where cv.mas_vehicle_uid = v.mas_vehicle_uid and cv.is_deleted = '0')")
	// Apply text search (VehicleBrandName OR VehicleLicensePlate)
	if searchText != "" {
		query = query.Where("v.vehicle_brand_name ILIKE ? OR v.vehicle_license_plate ILIKE ? OR fn_get_long_short_dept_name_by_dept_sap(d.vehicle_owner_dept_sap) ILIKE ?",
			"%"+searchText+"%", "%"+searchText+"%", "%"+searchText+"%")
	}

	// Count total records
	query.Count(&total)
	deptSAPWork := user.DeptSAP
	//order by vehicle_owner_dept_sap=deptSAPWork first
	query = query.Order(fmt.Sprintf("CASE WHEN d.vehicle_owner_dept_sap = '%s' THEN 0 ELSE 1 END", deptSAPWork))

	// Execute query with pagination
	query.Offset(offset).Limit(limit).Find(&vehicles)

	for i := range vehicles {
		funcs.TrimStringFields(&vehicles[i])
		vehicles[i].Age = funcs.CalculateAge(vehicles[i].VehicleRegistrationDate)
	}
	// Respond with JSON
	if len(vehicles) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"pagination": gin.H{
				"total":      total,
				"page":       page,
				"limit":      limit,
				"totalPages": (total + int64(limit) - 1) / int64(limit), // Calculate total pages
			},
			"vehicles": []models.VmsMasVehicleList{},
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"pagination": gin.H{
			"total":      total,
			"page":       page,
			"limit":      limit,
			"totalPages": (total + int64(limit) - 1) / int64(limit), // Calculate total pages
		},
		"vehicles": vehicles,
	})
}

// GetMasVehicleDetail godoc
// @Summary Retrieve a specific carpool vehicle
// @Description This endpoint fetches details of a specific carpool vehicle using its unique identifier (MasVehicleUID).
// @Tags Carpool-management
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Param active body []models.VmsMasVehicleArray true "array of VmsMasVehicleArray"
// @Router /api/carpool-management/vehicle-mas-details [post]
func (h *CarpoolManagementHandler) GetMasVehicleDetail(c *gin.Context) {
	funcs.GetAuthenUser(c, h.Role)
	if c.IsAborted() {
		return
	}
	var requests []models.VmsMasVehicleArray
	if err := c.ShouldBindJSON(&requests); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "message": messages.ErrInternalServer.Error()})
		return
	}
	var vehicles []models.VmsMasCarpoolVehicleDetail
	masVehicleUIDs := make([]string, len(requests))
	for i := range requests {
		masVehicleUIDs[i] = requests[i].MasVehicleUID
	}

	query := config.DB.Table("vms_mas_vehicle v").
		Select(
			`v.mas_vehicle_uid,
			v.vehicle_license_plate,
			v.vehicle_brand_name,
			v.vehicle_model_name,
			v.ref_vehicle_type_code,
			(select max(ref_vehicle_type_name) from vms_ref_vehicle_type s where s.ref_vehicle_type_code=v.ref_vehicle_type_code) ref_vehicle_type_name,
			(select max(s.dept_long_short) from vms_mas_department s where s.dept_sap=d.vehicle_owner_dept_sap) vehicle_owner_dept_short,
			v.ref_vehicle_type_code,
			d.fleet_card_no,
			v.is_tax_credit,
			d.vehicle_mileage,
			d.vehicle_get_date,
			d.ref_vehicle_status_code,
			(select max(s.ref_vehicle_status_short_name) from vms_ref_vehicle_status s where s.ref_vehicle_status_code=d.ref_vehicle_status_code) ref_vehicle_status_name,
			d.is_active,
			v.seat,
			v.vehicle_color,
			v.vehicle_gear,
			v.ref_fuel_type_id,
			(select max(s.ref_fuel_type_name_th) from vms_ref_fuel_type s where s.ref_fuel_type_id=v.ref_fuel_type_id) fuel_type_name,
			d.vehicle_pea_id,
			d.parking_place
		`).
		Joins("LEFT JOIN (SELECT DISTINCT ON (mas_vehicle_uid) * FROM vms_mas_vehicle_department WHERE is_deleted = '0' AND is_active = '1' ORDER BY mas_vehicle_uid, created_at DESC) d ON v.mas_vehicle_uid = d.mas_vehicle_uid").
		Where("v.mas_vehicle_uid IN (?) AND v.is_deleted = ?", masVehicleUIDs, "0")

	if err := query.Find(&vehicles).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Vehicle not found", "message": messages.ErrNotfound.Error()})
		return
	}

	for i := range vehicles {
		vehicles[i].Age = funcs.CalculateAge(vehicles[i].VehicleGetDate)
		funcs.TrimStringFields(&vehicles[i])
	}

	c.JSON(http.StatusOK, vehicles)
}

// SetActiveCarpoolVehicle godoc
// @Summary Set active status for a carpool vehicle
// @Description Update the active status of a carpool vehicle using its unique identifier (MasCarpoolVehicleUID).
// @Tags Carpool-management
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Param active body models.VmsMasCarpoolVehicleActive true "VmsMasCarpoolVehicleActive data"
// @Router /api/carpool-management/vehicle-set-active [put]
func (h *CarpoolManagementHandler) SetActiveCarpoolVehicle(c *gin.Context) {
	user := funcs.GetAuthenUser(c, h.Role)
	if c.IsAborted() {
		return
	}

	var request models.VmsMasCarpoolVehicleActive
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "message": messages.ErrInternalServer.Error()})
		return
	}

	var vehicle models.VmsMasCarpoolVehicle
	if err := config.DB.Where("mas_carpool_vehicle_uid = ? AND is_deleted = ?", request.MasCarpoolVehicleUID, "0").First(&vehicle).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Carpool vehicle not found", "message": messages.ErrNotfound.Error()})
		return
	}

	var existingCarpool models.VmsMasCarpoolRequest
	queryRole := h.SetQueryRole(user, config.DB)
	if err := queryRole.Where("mas_carpool_uid = ? AND is_deleted = ?", vehicle.MasCarpoolUID, "0").First(&existingCarpool).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Carpool not found", "message": messages.ErrNotfound.Error()})
		return
	}

	//update is_active to 1 in carpool_vehicle
	if err := config.DB.Model(&models.VmsMasCarpoolVehicle{}).Where("mas_carpool_vehicle_uid = ?", vehicle.MasCarpoolVehicleUID).UpdateColumns(map[string]interface{}{
		"is_active":  request.IsActive,
		"updated_at": time.Now(),
		"updated_by": user.EmpID,
	}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update carpool vehicle", "message": messages.ErrInternalServer.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Carpool vehicle active status updated successfully", "data": request, "carpool_name": GetCarpoolName(vehicle.MasCarpoolUID)})
}

// GetCarpoolVehicleTimeLine godoc
// @Summary Get Carpool vehicle timeline
// @Description Get Carpool vehicle timeline by date range
// @Tags Carpool-management
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Param mas_carpool_uid path string true "MasCarpoolUID (mas_carpool_uid)"
// @Param start_date query string true "Start date (YYYY-MM-DD)"
// @Param end_date query string true "End date (YYYY-MM-DD)"
// @Param search query string false "Search by vehicle license plate, brand, or model"
// @Param vehicel_car_type_detail query string false "Filter by Car type"
// @Param is_active query string false "Filter by is_active status (comma-separated, e.g., '1,0')"
// @Param ref_vehicle_status_code query string false "Filter by vehicle status code (comma-separated, e.g., '1,2')"
// @Router /api/carpool-management/vehicle-timeline/{mas_carpool_uid} [get]
func (h *CarpoolManagementHandler) GetCarpoolVehicleTimeLine(c *gin.Context) {
	user := funcs.GetAuthenUser(c, h.Role)
	if c.IsAborted() {
		return
	}
	masCarpoolUID := c.Param("mas_carpool_uid")
	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")

	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start date format", "message": messages.ErrInvalidDate.Error()})
		return
	}

	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid end date format", "message": messages.ErrInvalidDate.Error()})
		return
	}

	if startDate.After(endDate) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Start date must be before end date", "message": messages.ErrInvalidDate.Error()})
		return
	}
	// check if endDate - startDate > 3 months
	if endDate.Sub(startDate) > 90*24*time.Hour {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Date range cannot exceed 3 months", "message": messages.ErrInvalidDate.Error()})
		return
	}

	var vehicles []models.VehicleTimeLine

	query := h.SetQueryRole(user, config.DB)
	query = query.Table("public.vms_mas_vehicle AS v").
		Select(`v.mas_vehicle_uid, v.vehicle_license_plate, v.vehicle_license_plate_province_short, 
				v.vehicle_license_plate_province_full, d.county, d.vehicle_get_date, d.vehicle_pea_id, 
				d.vehicle_license_plate_province_short, d.vehicle_license_plate_province_full, 
				public.fn_get_long_short_dept_name_by_dept_sap(d.vehicle_owner_dept_sap)  AS vehicle_dept_name,
				(select max(cp.carpool_name) from vms_mas_carpool cp, vms_mas_carpool_vehicle cpv where cpv.is_deleted = '0' and cpv.is_active = '1' and cpv.mas_carpool_uid = cp.mas_carpool_uid and cpv.mas_vehicle_uid = v.mas_vehicle_uid) AS vehicle_carpool_name, 
				v."CarTypeDetail" AS vehicle_car_type_detail, 0 AS vehicle_mileage,
				v.vehicle_brand_name,v.vehicle_model_name,
				public.fn_get_vehicle_distance_two_months(v.mas_vehicle_uid, ?) AS vehicle_distance`, startDate).
		Joins("LEFT JOIN (SELECT DISTINCT ON (mas_vehicle_uid) * FROM vms_mas_vehicle_department WHERE is_deleted = '0' AND is_active = '1' ORDER BY mas_vehicle_uid, created_at DESC) d ON v.mas_vehicle_uid = d.mas_vehicle_uid").
		Where("v.is_deleted = ?", "0").
		Where("exists (select 1 from vms_mas_carpool_vehicle cv where cv.mas_vehicle_uid = v.mas_vehicle_uid and cv.is_deleted = '0' and cv.mas_carpool_uid = ?)", masCarpoolUID)
	if refTimelineStatusID := c.Query("ref_timeline_status_id"); refTimelineStatusID != "" {
		refTimelineStatusIDList := strings.Split(refTimelineStatusID, ",")
		query = query.Where(`exists (select 1 from vms_trn_request r where r.mas_vehicle_uid = v.mas_vehicle_uid AND r.ref_request_status_code != '90' AND (
						('1' in (?) AND r.ref_request_status_code < '50') OR
						('2' in (?) AND r.ref_request_status_code >= '50' AND r.ref_request_status_code < '80' AND r.ref_trip_type_code = 0) OR 
						('3' in (?) AND r.ref_request_status_code >= '50' AND r.ref_request_status_code < '80' AND r.ref_trip_type_code = 1) OR
						('4' in (?) AND r.ref_request_status_code = '80')
					) AND
						 (reserve_start_datetime BETWEEN ? AND ? OR reserve_end_datetime BETWEEN ? AND ?)
				)`, refTimelineStatusIDList, refTimelineStatusIDList, refTimelineStatusIDList, refTimelineStatusIDList, startDate, endDate, startDate, endDate)
	}

	if search := c.Query("search"); search != "" {
		query = query.Where("v.vehicle_license_plate ILIKE ? OR v.vehicle_brand_name ILIKE ? OR v.vehicle_model_name ILIKE ?", "%"+search+"%", "%"+search+"%", "%"+search+"%")
	}
	if vehicleOwnerDeptSAP := c.Query("vehicle_owner_dept_sap"); vehicleOwnerDeptSAP != "" {
		query = query.Where("d.vehicle_owner_dept_sap = ?", vehicleOwnerDeptSAP)
	}

	if vehicleCarTypeDetail := c.Query("vehicel_car_type_detail"); vehicleCarTypeDetail != "" {
		query = query.Where("v.\"CarTypeDetail\" = ?", vehicleCarTypeDetail)
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

	if err := query.Find(&vehicles).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "message": messages.ErrInternalServer.Error()})
		return
	}

	for i := range vehicles {
		parts := strings.Split(vehicles[i].VehicleDistance, ",")
		if len(parts) > 3 {
			vehicles[i].VehicleMileage = parts[3]
		}
		// Preload the vehicle requests for each vehicle
		query := config.DB.Table("vms_trn_request r").
			Preload("MasDriver").
			Where("mas_vehicle_uid = ? AND is_deleted = ? AND (reserve_start_datetime BETWEEN ? AND ? OR reserve_end_datetime BETWEEN ? AND ?)", vehicles[i].MasVehicleUID, "0", startDate, endDate, startDate, endDate).
			Where("ref_request_status_code != '90'")
		if refTimelineStatusID := c.Query("ref_timeline_status_id"); refTimelineStatusID != "" {
			refTimelineStatusIDList := strings.Split(refTimelineStatusID, ",")
			query = query.Where(`
					('1' in (?) AND r.ref_request_status_code < '50') OR
					('2' in (?) AND r.ref_request_status_code >= '50' AND r.ref_request_status_code < '80' AND r.ref_trip_type_code = 0) OR 
					('3' in (?) AND r.ref_request_status_code >= '50' AND r.ref_request_status_code < '80' AND r.ref_trip_type_code = 1) OR
					('4' in (?) AND r.ref_request_status_code = '80')
				`, refTimelineStatusIDList, refTimelineStatusIDList, refTimelineStatusIDList, refTimelineStatusIDList)
		}
		if err := query.Find(&vehicles[i].VehicleTrnRequests).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "message": messages.ErrInternalServer.Error()})
			return
		}

		for j := range vehicles[i].VehicleTrnRequests {
			vehicles[i].VehicleTrnRequests[j].TripDetails = []models.VmsTrnTripDetail{
				{
					TrnTripDetailUID: uuid.New().String(),
					VmsTrnTripDetailRequest: models.VmsTrnTripDetailRequest{
						TrnRequestUID:        vehicles[i].VehicleTrnRequests[j].TrnRequestUID,
						TripStartDatetime:    vehicles[i].VehicleTrnRequests[j].ReserveStartDatetime,
						TripEndDatetime:      vehicles[i].VehicleTrnRequests[j].ReserveEndDatetime,
						TripDeparturePlace:   vehicles[i].VehicleTrnRequests[j].WorkPlace,
						TripDestinationPlace: vehicles[i].VehicleTrnRequests[j].WorkPlace,
						TripStartMiles:       0,
						TripEndMiles:         0,
					},
				},
			}
			if vehicles[i].VehicleTrnRequests[j].RefRequestStatusCode == "80" {
				vehicles[i].VehicleTrnRequests[j].RefTimelineStatusID = "4"
				vehicles[i].VehicleTrnRequests[j].TimeLineStatus = "เสร็จสิ้น"
			} else if vehicles[i].VehicleTrnRequests[j].RefRequestStatusCode < "50" {
				vehicles[i].VehicleTrnRequests[j].RefTimelineStatusID = "1"
				vehicles[i].VehicleTrnRequests[j].TimeLineStatus = "รออนุมัติ"
			} else if vehicles[i].VehicleTrnRequests[j].RefTripTypeCode == 0 {
				vehicles[i].VehicleTrnRequests[j].RefTimelineStatusID = "2"
				vehicles[i].VehicleTrnRequests[j].TimeLineStatus = "ไป-กลับ"
			} else if vehicles[i].VehicleTrnRequests[j].RefTripTypeCode == 1 {
				vehicles[i].VehicleTrnRequests[j].RefTimelineStatusID = "3"
				vehicles[i].VehicleTrnRequests[j].TimeLineStatus = "ค้างแรม"
			}
			vehicles[i].VehicleTrnRequests[j].RefRequestStatusName = StatusNameMapUser[vehicles[i].VehicleTrnRequests[j].RefRequestStatusCode]
		}
	}
	thaiMonths := []string{"ม.ค.", "ก.พ.", "มี.ค.", "เม.ย.", "พ.ค.", "มิ.ย.", "ก.ค.", "ส.ค.", "ก.ย.", "ต.ค.", "พ.ย.", "ธ.ค."}
	lastMonthDate := time.Date(startDate.Year(), startDate.Month()-1, 1, 0, 0, 0, 0, startDate.Location())
	lastMonth := fmt.Sprintf("%s%02d", thaiMonths[lastMonthDate.Month()-1], (lastMonthDate.Year()+543)%100)

	c.JSON(http.StatusOK, gin.H{
		"pagination": gin.H{
			"total":      total,
			"page":       page,
			"limit":      pageSizeInt,
			"totalPages": (total + int64(pageSizeInt) - 1) / int64(pageSizeInt), // Calculate total pages
		},
		"last_month": lastMonth,
		"vehicles":   vehicles,
	})
}
