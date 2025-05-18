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

// SearchCarpoolDriver godoc
// @Summary Search carpool drivers
// @Description Search carpool drivers with pagination and filters
// @Tags Carpool-management
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Param mas_carpool_uid path string true "MasCarpoolUID (mas_carpool_uid)"
// @Param search query string false "Search query for driver_name or driver_license_no"
// @Param is_active query string false "Filter by is_active status (comma-separated, e.g., '1,0')"
// @Param order_by query string false "Order by fields: driver_name"
// @Param order_dir query string false "Order direction: asc or desc"
// @Param page query int false "Page number (default: 1)"
// @Param limit query int false "Number of records per page (default: 10)"
// @Router /api/carpool-management/driver-search/{mas_carpool_uid} [get]
func (h *CarpoolManagementHandler) SearchCarpoolDriver(c *gin.Context) {
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

	var drivers []models.VmsMasCarpoolDriverList
	query := config.DB.Table("vms_mas_carpool_driver cpd").
		Select(
			`cpd.mas_carpool_driver_uid, cpd.mas_carpool_uid, d.mas_driver_uid,
			d.driver_image, d.driver_name, d.driver_nickname,
			d.driver_dept_sap_short_name_hire, d.driver_contact_number,
			dl.driver_license_end_date,
			d.approved_job_driver_end_date, d.driver_average_satisfaction_score,
			d.ref_driver_status_code, rds.ref_driver_status_desc AS driver_status_name,
			d.is_active`).
		Joins("LEFT JOIN vms_mas_driver d ON d.mas_driver_uid = cpd.mas_driver_uid").
		Joins("LEFT JOIN vms_mas_driver_license dl ON dl.mas_driver_uid = d.mas_driver_uid").
		Joins("LEFT JOIN vms_ref_driver_status rds ON rds.ref_driver_status_code = d.ref_driver_status_code").
		Where("cpd.mas_carpool_uid = ? AND cpd.is_deleted = ?", masCarpoolUID, "0").
		Group("cpd.mas_carpool_driver_uid, cpd.mas_carpool_uid, d.mas_driver_uid, dl.driver_license_end_date, rds.ref_driver_status_desc")

	search := strings.ToUpper(c.Query("search"))
	if search != "" {
		query = query.Where("UPPER(d.driver_name) ILIKE ? OR UPPER(d.driver_license_no) ILIKE ?", "%"+search+"%", "%"+search+"%")
	}

	if isActive := c.Query("is_active"); isActive != "" {
		isActiveList := strings.Split(isActive, ",")
		query = query.Where("cpd.is_active IN (?)", isActiveList)
	}

	orderBy := c.Query("order_by")
	orderDir := c.Query("order_dir")
	if orderDir != "desc" {
		orderDir = "asc"
	}
	switch orderBy {
	case "driver_name":
		query = query.Order("d.driver_name " + orderDir)
	default:
		query = query.Order("d.driver_name " + orderDir) // Default ordering by driver_name
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	query = query.Limit(limit).Offset(offset)
	if err := query.Find(&drivers).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	for i := range drivers {
		funcs.TrimStringFields(&drivers[i])
	}
	if len(drivers) == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "No carpool drivers found",
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
			"drivers": drivers,
		})
	}
}

// CreateCarpoolDriver godoc
// @Summary Create a new carpool driver
// @Description Create a new carpool driver and save it to the database
// @Tags Carpool-management
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Param driver body []models.VmsMasCarpoolDriver true "VmsMasCarpoolDriver array"
// @Router /api/carpool-management/driver-create [post]
func (h *CarpoolManagementHandler) CreateCarpoolDriver(c *gin.Context) {
	user := funcs.GetAuthenUser(c, h.Role)
	if c.IsAborted() {
		return
	}

	var requests []models.VmsMasCarpoolDriver
	if err := c.ShouldBindJSON(&requests); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if len(requests) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No carpool vehicle data provided", "message": messages.ErrInvalidJSONInput.Error()})
		return

	}
	for i := range requests {
		var existingCarpool models.VmsMasCarpoolRequest
		queryRole := h.SetQueryRole(user, config.DB)
		if err := queryRole.Where("mas_carpool_uid = ? AND is_deleted = ?", requests[i].MasCarpoolUID, "0").First(&existingCarpool).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Carpool not found", "message": messages.ErrNotfound.Error()})
			return
		}
	}
	for i := range requests {
		var existingDriver models.VmsMasCarpoolDriver
		if err := config.DB.Where("mas_driver_uid = ? AND is_deleted = ?", requests[i].MasDriverUID, "0").First(&existingDriver).Error; err == nil {
			c.JSON(http.StatusConflict, gin.H{
				"error": fmt.Sprintf("Driver with MasCarpoolUID %s and MasDriverUID %s already exists", requests[i].MasCarpoolUID, requests[i].MasDriverUID),
			})
			return
		}

		requests[i].MasCarpoolDriverUID = uuid.New().String()
		requests[i].CreatedAt = time.Now()
		requests[i].CreatedBy = user.EmpID
		requests[i].UpdatedAt = time.Now()
		requests[i].UpdatedBy = user.EmpID
		requests[i].IsDeleted = "0"
		requests[i].IsActive = "1"

		requests[i].StartDate = time.Now()
		requests[i].EndDate = time.Now().AddDate(1, 0, 0)
	}

	if err := config.DB.Create(&requests).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":      "Carpool drivers created successfully",
		"data":         requests,
		"carpool_name": GetCarpoolName(requests[0].MasCarpoolUID),
	})
}

// DeleteCarpoolDriver godoc
// @Summary Delete a carpool driver
// @Description This endpoint deletes a carpool driver using its unique identifier (MasCarpoolDriverUID).
// @Tags Carpool-management
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Param mas_carpool_driver_uid path string true "MasCarpoolDriverUID (mas_carpool_driver_uid)"
// @Router /api/carpool-management/driver-delete/{mas_carpool_driver_uid} [delete]
func (h *CarpoolManagementHandler) DeleteCarpoolDriver(c *gin.Context) {
	user := funcs.GetAuthenUser(c, h.Role)
	if c.IsAborted() {
		return
	}
	masCarpoolDriverUID := c.Param("mas_carpool_driver_uid")

	var driver models.VmsMasCarpoolDriver
	if err := config.DB.Where("mas_carpool_driver_uid = ? AND is_deleted = ?", masCarpoolDriverUID, "0").First(&driver).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Carpool driver not found"})
		return
	}
	var existingCarpool models.VmsMasCarpoolRequest
	queryRole := h.SetQueryRole(user, config.DB)
	if err := queryRole.Where("mas_carpool_uid = ? AND is_deleted = ?", driver.MasCarpoolUID, "0").First(&existingCarpool).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Carpool not found", "message": messages.ErrNotfound.Error()})
		return
	}
	if err := config.DB.Model(&driver).UpdateColumns(map[string]interface{}{
		"is_deleted": "1",
		"updated_by": user.EmpID,
		"updated_at": time.Now(),
	}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete carpool driver"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Carpool driver deleted successfully", "carpool_name": GetCarpoolName(driver.MasCarpoolUID)})
}

// SearchMasDrivers godoc
// @Summary Get drivers by name with pagination
// @Description Get a list of drivers filtered by name with pagination
// @Tags Carpool-management
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Param name query string false "Driver name to search"
// @Param page query int false "Page number (default: 1)"
// @Param limit query int false "Number of records per page (default: 10)"
// @Router /api/carpool-management/driver-mas-search [get]
func (h *CarpoolManagementHandler) SearchMasDrivers(c *gin.Context) {
	funcs.GetAuthenUser(c, h.Role)
	if c.IsAborted() {
		return
	}
	name := strings.ToUpper(c.Query("name"))
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))    // Default: page 1
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10")) // Default: 10 items per page
	offset := (page - 1) * limit

	var drivers []models.VmsMasDriver
	query := h.SetQueryRoleDept(funcs.GetAuthenUser(c, h.Role), config.DB)
	query = query.Model(&models.VmsMasDriver{})
	query = query.Where("is_deleted = ?", "0")
	// Apply search filter
	if name != "" {
		searchTerm := "%" + name + "%"
		query = query.Where(`
            driver_name ILIKE ? OR 
            driver_id ILIKE ?`,
			searchTerm, searchTerm)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "message": messages.ErrInternalServer.Error()})
		return
	}
	query = query.Limit(limit).
		Offset(offset)

	if err := query.
		Preload("DriverStatus").
		Find(&drivers).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "message": messages.ErrInternalServer.Error()})
		return
	}

	if len(drivers) == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "No drivers found",
			"pagination": gin.H{
				"page":       page,
				"limit":      limit,
				"totalPages": (total + int64(limit) - 1) / int64(limit), // Calculate total pages
				"drivers":    drivers,
			},
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"pagination": gin.H{
				"total":      total,
				"page":       page,
				"limit":      limit,
				"totalPages": (total + int64(limit) - 1) / int64(limit), // Calculate total pages
			},
			"drivers": drivers,
		})
	}
}

// GetMasDriverDetails godoc
// @Summary Retrieve a specific driver
// @Description This endpoint fetches details of a specific driver using its unique identifier (MasDriverUID).
// @Tags Carpool-management
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Param active body []models.VmsMasDriverArray true "array of VmsMasDriverArray"
// @Router /api/carpool-management/driver-mas-details [post]
func (h *CarpoolManagementHandler) GetMasDriverDetails(c *gin.Context) {
	funcs.GetAuthenUser(c, h.Role)
	if c.IsAborted() {
		return
	}
	var request []models.VmsMasDriverArray
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var drivers []models.VmsMasCarpoolDriverDetail
	masDriverUIDs := make([]string, len(request))
	for i := range request {
		masDriverUIDs[i] = request[i].MasDriverUID
	}

	query := config.DB.Table("vms_mas_driver d").
		Select(
			`d.mas_driver_uid,
			d.driver_image,
			d.driver_name,
			d.driver_nickname,
			d.driver_dept_sap_short_name_hire,
			d.driver_contact_number,
			(select driver_license_end_date from vms_mas_driver_license s where s.mas_driver_uid=d.mas_driver_uid) driver_license_end_date,
			d.approved_job_driver_end_date,
			d.driver_average_satisfaction_score,
			d.ref_driver_status_code,
			(select max(s.ref_driver_status_desc) from vms_ref_driver_status s WHERE s.ref_driver_status_code = d.ref_driver_status_code) AS driver_status_name,
			d.is_active
	`).
		Where("d.mas_driver_uid in (?) AND d.is_deleted = ?", masDriverUIDs, "0")

	if err := query.Find(&drivers).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Driver not found", "message": messages.ErrNotfound.Error()})
		return
	}

	c.JSON(http.StatusOK, drivers)
}

// SetActiveCarpoolDriver godoc
// @Summary Set active status for a carpool driver
// @Description Update the active status of a carpool driver using its unique identifier (MasCarpoolDriverUID).
// @Tags Carpool-management
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Param active body models.VmsMasCarpoolDriverActive true "VmsMasCarpoolDriverActive data"
// @Router /api/carpool-management/driver-set-active [put]
func (h *CarpoolManagementHandler) SetActiveCarpoolDriver(c *gin.Context) {
	user := funcs.GetAuthenUser(c, h.Role)
	if c.IsAborted() {
		return
	}

	var request models.VmsMasCarpoolDriverActive
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var driver models.VmsMasCarpoolDriver
	if err := config.DB.Where("mas_carpool_driver_uid = ? AND is_deleted = ?", request.MasCarpoolDriverUID, "0").First(&driver).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Carpool driver not found", "message": messages.ErrNotfound.Error()})
		return
	}

	var existingCarpool models.VmsMasCarpoolRequest
	queryRole := h.SetQueryRole(user, config.DB)
	if err := queryRole.Where("mas_carpool_uid = ? AND is_deleted = ?", driver.MasCarpoolUID, "0").First(&existingCarpool).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Carpool not found", "message": messages.ErrNotfound.Error()})
		return
	}

	driver.IsActive = request.IsActive
	driver.UpdatedAt = time.Now()
	driver.UpdatedBy = user.EmpID

	if err := config.DB.Save(&driver).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to update active status: %v", err), "message": messages.ErrInternalServer.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Carpool driver active status updated successfully", "data": request, "carpool_name": GetCarpoolName(driver.MasCarpoolUID)})
}

// GetCarpoolDriverTimeLine godoc
// @Summary Get driver timeline
// @Description Get driver timeline by date range
// @Tags Carpool-management
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Param mas_carpool_uid path string true "MasCarpoolUID (mas_carpool_uid)"
// @Param start_date query string true "Start date (YYYY-MM-DD)"
// @Param end_date query string true "End date (YYYY-MM-DD)"
// @Param search query string false "driver_name,driver_nickname,driver_dept_sap_short_name_work to search"
// @Param work_type query string false "work type 1: ค้างคืน, 2: ไป-กลับ Filter by multiple work_type (comma-separated, e.g., '1,2')"
// @Param ref_driver_status_code query string false "Filter by driver status code (comma-separated, e.g., '1,2')"
// @Param is_active query string false "Filter by is_active status (comma-separated, e.g., '1,0')"
// @Param page query int false "Page number (default: 1)"
// @Param limit query int false "Number of records per page (default: 10)"
// @Router /api/carpool-management/driver-timeline/{mas_carpool_uid} [get]
func (h *CarpoolManagementHandler) GetCarpoolDriverTimeLine(c *gin.Context) {
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

	var drivers []models.DriverTimeLine

	query := h.SetQueryRole(user, config.DB).
		Table("public.vms_mas_driver AS d").
		Select("d.*").
		Where("d.is_deleted = ? AND d.is_active = ?", "0", "1").
		Joins("INNER JOIN vms_mas_carpool_driver cd ON cd.mas_driver_uid = d.mas_driver_uid AND cd.mas_carpool_uid = ? AND cd.is_deleted = ?", masCarpoolUID, "0").
		Joins(`INNER JOIN vms_trn_request r 
		   ON r.mas_carpool_driver_uid = d.mas_driver_uid 
		   AND r.is_pea_employee_driver = ? 
		   AND (r.reserve_start_datetime BETWEEN ? AND ? 
		   OR r.reserve_end_datetime BETWEEN ? AND ? 
		   OR ? BETWEEN r.reserve_start_datetime AND r.reserve_end_datetime 
		   OR ? BETWEEN r.reserve_start_datetime AND r.reserve_end_datetime)`,
			"0", startDate, endDate, startDate, endDate, startDate, endDate)

	name := strings.ToUpper(c.Query("name"))
	if name != "" {
		query = query.Where("UPPER(driver_name) ILIKE ? OR UPPER(driver_nickname) ILIKE ? OR UPPER(driver_dept_sap_short_name_work) ILIKE ?", "%"+name+"%", "%"+name+"%", "%"+name+"%")
	}
	if workType := c.Query("work_type"); workType != "" {
		workTypes := strings.Split(workType, ",")
		query = query.Where("work_type IN (?)", workTypes)
	}
	if refDriverStatusCode := c.Query("ref_driver_status_code"); refDriverStatusCode != "" {
		statusCodes := strings.Split(refDriverStatusCode, ",")
		query = query.Where("ref_driver_status_code IN (?)", statusCodes)
	}
	if isActive := c.Query("is_active"); isActive != "" {
		isActiveValues := strings.Split(isActive, ",")
		query = query.Where("is_active IN (?)", isActiveValues)
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
	if err := query.Find(&drivers).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "message": messages.ErrInternalServer.Error()})
		return
	}
	for i := range drivers {
		drivers[i].WorkLastMonth = fmt.Sprintf("%d วัน/%d งาน", 22, 3)
		drivers[i].WorkThisMonth = fmt.Sprintf("%d วัน/%d งาน", 16, 2)

		// Preload the driver requests for each driver
		if err := config.DB.Table("vms_trn_request").
			Preload("TripDetails").
			Where("mas_carpool_driver_uid = ? AND is_pea_employee_driver = ? AND is_deleted = ? AND (reserve_start_datetime BETWEEN ? AND ? OR reserve_end_datetime BETWEEN ? AND ?)", drivers[i].MasDriverUID, "0", "0", startDate, endDate, startDate, endDate).
			Find(&drivers[i].DriverTrnRequests).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "message": messages.ErrInternalServer.Error()})
			return
		}
		// Preload the driver status for each driver

		for j := range drivers[i].DriverTrnRequests {
			drivers[i].DriverTrnRequests[j].RefRequestStatusName = StatusNameMapUser[drivers[i].DriverTrnRequests[j].RefRequestStatusCode]
		}
	}
	c.JSON(http.StatusOK, gin.H{
		"drivers": drivers,
		"pagination": gin.H{
			"total":      total,
			"page":       pageInt,
			"limit":      pageSizeInt,
			"totalPages": (total + int64(pageSizeInt) - 1) / int64(pageSizeInt), // Calculate total pages
		},
	})
}
