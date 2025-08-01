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
	"gorm.io/gorm"
)

type DriverHandler struct {
	Role string
}

// GetDrivers godoc
// @Summary Get drivers by name with pagination
// @Description Get a list of drivers filtered by name with pagination
// @Tags Drivers
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Param name query string false "Driver name to search"
// @Param work_type query string false "work type to search (0: ไป-กลับ ,1: ค้างคืน)"
// @Param page query int false "Page number (default: 1)"
// @Param limit query int false "Number of records per page (default: 10)"
// @Router /api/driver/search [get]
func (h *DriverHandler) GetDrivers(c *gin.Context) {
	funcs.GetAuthenUser(c, h.Role)
	if c.IsAborted() {
		return
	}
	name := strings.ToUpper(c.Query("name"))
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))    // Default: page 1
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10")) // Default: 10 items per page
	offset := (page - 1) * limit

	var drivers []models.VmsMasDriver
	query := config.DB.Model(&models.VmsMasDriver{})
	query = query.Where("is_deleted = ? AND is_replacement = ?", "0", "0")
	// Apply search filter
	if name != "" {
		searchTerm := "%" + name + "%"
		query = query.Where(`
            driver_name ILIKE ? OR 
            driver_id ILIKE ?`,
			searchTerm, searchTerm)
	}
	if workType := c.Query("work_type"); workType != "" {
		query = query.Where("work_type = ?", workType)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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
	for i := range drivers {
		drivers[i].Age = drivers[i].CalculateAgeInYearsMonths()
		drivers[i].Status = "ว่าง"

		switch drivers[i].WorkType {
		case 1:
			drivers[i].WorkTypeName = "ค้างคืน"
		case 2:
			drivers[i].WorkTypeName = "ไป-กลับ"
		}
		drivers[i].WorkCount = 0
		drivers[i].WorkDays = 0
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

// GetBookingDrivers godoc
// @Summary Get drivers by name with pagination
// @Description Get a list of drivers filtered by name with pagination
// @Tags Drivers
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Param emp_id query string false "Employee ID (emp_id) default(700001)"
// @Param start_date query string false "Start Date (YYYY-MM-DD HH:mm:ss)" default(2025-05-30 08:00:00)
// @Param end_date query string false "End Date (YYYY-MM-DD HH:mm:ss)" default(2025-05-30 16:00:00)
// @Param name query string false "Driver name to search"
// @Param work_type query string false "work type to search (0: ไป-กลับ,1: ค้างคืน)" default(0)
// @Param mas_carpool_uid query string false "MasCarpoolUID (mas_carpool_uid)"
// @Param mas_vehicle_uid query string false "MasVehicleUID (mas_vehicle_uid)"
// @Param page query int false "Page number (default: 1)"
// @Param limit query int false "Number of records per page (default: 10)"
// @Router /api/driver/search-booking [get]
func (h *DriverHandler) GetBookingDrivers(c *gin.Context) {
	user := funcs.GetAuthenUser(c, h.Role)
	if c.IsAborted() {
		return
	}
	name := strings.ToUpper(c.Query("name"))
	workType := c.Query("work_type")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))    // Default: page 1
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10")) // Default: 10 items per page
	offset := (page - 1) * limit

	empID := c.Query("emp_id")
	bureauDeptSap := user.BureauDeptSap
	businessArea := user.BusinessArea

	if empID != "" {
		empUser := funcs.GetUserEmpInfo(empID)
		bureauDeptSap = empUser.BureauDeptSap
		businessArea = empUser.BusinessArea
	}

	startDate := c.Query("start_date")
	endDate := c.Query("end_date")
	if startDate == "" || endDate == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Start Date and End Date are required", "message": messages.ErrInvalidRequest.Error()})
		return
	}

	StartTimeWithZone, err1 := models.GetTimeWithZone(startDate)
	EndTimeWithZone, err2 := models.GetTimeWithZone(endDate)
	if err1 != nil || err2 != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get start time with zone", "message": messages.ErrInternalServer.Error()})
		return
	}

	fmt.Println(StartTimeWithZone)
	fmt.Println(EndTimeWithZone)

	var driverCanBookings []models.VmsMasDriverCanBooking
	tripTypeCode := 0
	if workType == "1" {
		tripTypeCode = 1
	} else {
		tripTypeCode = 0
	}
	queryCanBooking := config.DB.Raw(`SELECT * FROM fn_get_available_drivers_view (?, ?, ?, ?, ?)`,
		StartTimeWithZone, EndTimeWithZone, bureauDeptSap, businessArea, tripTypeCode)
	err := queryCanBooking.Scan(&driverCanBookings).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch available vehicles", "message": messages.ErrInternalServer.Error()})
		return
	}
	masDriverUIDs := make([]string, 0)
	for _, driverCanBooking := range driverCanBookings {
		masDriverUIDs = append(masDriverUIDs, driverCanBooking.MasDriverUID)
	}
	date, _ := time.Parse("2006-01-02 15:04:05", startDate)
	var drivers []models.VmsMasDriver
	query := config.DB.Model(&models.VmsMasDriver{})
	query = query.Select("vms_mas_driver.*, w_thismth.job_count, w_thismth.total_days,mc.carpool_name")
	query = query.Joins("LEFT JOIN vms_mas_carpool mc ON mc.mas_carpool_uid = vms_mas_driver.mas_carpool_uid")
	query = query.Joins("LEFT JOIN public.vms_trn_driver_monthly_workload AS w_thismth ON w_thismth.workload_year = ? AND w_thismth.workload_month = ? AND w_thismth.driver_emp_id = vms_mas_driver.driver_id AND w_thismth.is_deleted = ?", date.Year(), date.Month(), "0")
	query = query.Where("vms_mas_driver.is_deleted = ? AND vms_mas_driver.is_replacement = ?", "0", "0")
	query = query.Where("vms_mas_driver.mas_driver_uid IN (?)", masDriverUIDs)

	isDepartment := false
	masVehicleUID := c.Query("mas_vehicle_uid")
	masCarpoolUID := c.Query("mas_carpool_uid")
	var vehicleCarpoolOrDeptSap struct {
		MasCarpoolUID string `gorm:"column:mas_carpool_uid"`
		BureauDeptSap string `gorm:"column:bureau_dept_sap"`
	}
	if masVehicleUID != "" && masCarpoolUID == "" {

		if err := config.DB.Table("vms_v_info_vehicle_all_active").Where("mas_vehicle_uid = ?", masVehicleUID).
			Select("mas_carpool_uid,bureau_dept_sap").
			First(&vehicleCarpoolOrDeptSap).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "MasVehicleUID not found", "message": messages.ErrNotfound.Error()})
			return
		}
		if vehicleCarpoolOrDeptSap.BureauDeptSap != "" && vehicleCarpoolOrDeptSap.MasCarpoolUID == "" {
			isDepartment = true
		}
		if vehicleCarpoolOrDeptSap.MasCarpoolUID != "" {
			query = query.Where("exists (select 1 from vms_mas_carpool_driver where mas_carpool_uid = ? AND mas_driver_uid = vms_mas_driver.mas_driver_uid AND is_deleted = '0')",
				vehicleCarpoolOrDeptSap.MasCarpoolUID)
		}
	}
	if masCarpoolUID != "" {
		query = query.Where("exists (select 1 from vms_mas_carpool_driver where mas_carpool_uid = ? AND mas_driver_uid = vms_mas_driver.mas_driver_uid AND is_deleted = '0')", masCarpoolUID)
	}
	if isDepartment {
		query = query.Where("vms_mas_driver.bureau_dept_sap = ?", vehicleCarpoolOrDeptSap.BureauDeptSap)
	}

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
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	query = query.Limit(limit).
		Offset(offset)

	if err := query.
		Preload("DriverStatus").
		Preload("DriverLicense", func(db *gorm.DB) *gorm.DB {
			// Use COALESCE to treat NULL driver_license_end_date as the earliest possible date for ordering
			return db.Order("COALESCE(driver_license_end_date, '1900-01-01') DESC")
		}).
		Preload("DriverLicense.DriverLicenseType").
		Find(&drivers).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "message": messages.ErrInternalServer.Error()})
		return
	}
	for i := range drivers {
		drivers[i].Age = drivers[i].CalculateAgeInYearsMonths()

		//อายุเกิน 60 ปี (60 years old)
		if drivers[i].DriverBirthdate.Before(time.Now().AddDate(-60, 0, 0)) {
			drivers[i].Status = "อายุเกิน 60 ปี"
			drivers[i].CanSelect = false
		} else {
			drivers[i].Status = "ว่าง"
			drivers[i].CanSelect = true
		}
		switch drivers[i].WorkType {
		case 1:
			drivers[i].WorkTypeName = "ค้างคืน"
		case 2:
			drivers[i].WorkTypeName = "ไป-กลับ"
		}
		drivers[i].VendorName = drivers[i].DriverDeptSAPShort
		if drivers[i].CarpoolName != "" {
			drivers[i].VendorName = drivers[i].CarpoolName
		}
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

// GetDriversOtherDept godoc
// @Summary Get drivers by name with pagination from other department
// @Description Get a list of drivers filtered by name with pagination from other department
// @Tags Drivers
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Param name query string false "Driver name to search"
// @Param page query int false "Page number (default: 1)"
// @Param limit query int false "Number of records per page (default: 10)"
// @Router /api/driver/search-other-dept [get]
func (h *DriverHandler) GetDriversOtherDept(c *gin.Context) {
	funcs.GetAuthenUser(c, h.Role)
	if c.IsAborted() {
		return
	}
	name := c.Query("name")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))    // Default: page 1
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10")) // Default: 10 items per page
	offset := (page - 1) * limit

	var drivers []models.VmsMasDriver

	// Get the total count of drivers (for pagination)
	var total int64
	config.DB.Model(&models.VmsMasDriver{}).Where("driver_name ILIKE ?", "%"+name+"%").Count(&total)

	// Fetch the drivers with pagination
	result := config.DB.Where("driver_name ILIKE ?", "%"+name+"%").Limit(limit).Offset(offset).
		Preload("DriverStatus").
		Find(&drivers)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error(), "message": messages.ErrInternalServer.Error()})
		return
	}

	for i := range drivers {
		drivers[i].Age = drivers[i].CalculateAgeInYearsMonths()
		drivers[i].Status = "ว่าง"

		switch drivers[i].WorkType {
		case 1:
			drivers[i].WorkTypeName = "ค้างคืน"
		case 2:
			drivers[i].WorkTypeName = "ไป-กลับ"
		}
		drivers[i].WorkCount = 4
		drivers[i].WorkDays = 16
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

// GetDriver godoc
// @Summary Retrieve a specific driver
// @Description This endpoint fetches details of a driver using its unique identifier (MasDriverUID).
// @Tags Drivers
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Param mas_driver_uid path string true "MasDriverUID (mas_driver_uid)"
// @Router /api/driver/{mas_driver_uid} [get]
func (h *DriverHandler) GetDriver(c *gin.Context) {
	funcs.GetAuthenUser(c, h.Role)
	if c.IsAborted() {
		return
	}
	masDriverUID := c.Param("mas_driver_uid")
	parsedID, err := uuid.Parse(masDriverUID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid MasDriverUID", "message": messages.ErrInvalidUID.Error()})
		return
	}
	date := time.Now()
	var driver models.VmsMasDriver
	query := config.DB.Model(&models.VmsMasDriver{})
	query = query.Select("vms_mas_driver.*, w_thismth.job_count, w_thismth.total_days,mc.carpool_name")
	query = query.Joins("LEFT JOIN vms_mas_carpool mc ON mc.mas_carpool_uid = vms_mas_driver.mas_carpool_uid")
	query = query.Joins("LEFT JOIN public.vms_trn_driver_monthly_workload AS w_thismth ON w_thismth.workload_year = ? AND w_thismth.workload_month = ? AND w_thismth.driver_emp_id = vms_mas_driver.driver_id AND w_thismth.is_deleted = ?", date.Year(), date.Month(), "0")

	if err := query.
		Preload("DriverLicense", func(db *gorm.DB) *gorm.DB {
			// Use COALESCE to treat NULL driver_license_end_date as the earliest possible date for ordering
			return db.Order("COALESCE(driver_license_end_date, '1900-01-01') DESC")
		}).
		Preload("DriverLicense.DriverLicenseType").
		Preload("DriverStatus").
		First(&driver, "mas_driver_uid = ?", parsedID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Request not found", "message": messages.ErrNotfound.Error()})
		return
	}
	driver.Age = driver.CalculateAgeInYearsMonths()
	switch driver.WorkType {
	case 1:
		driver.WorkTypeName = "ค้างคืน"
	case 2:
		driver.WorkTypeName = "ไป-กลับ"
	default:
		driver.WorkTypeName = "ไม่ระบุ"
	}
	driver.DriverDeptSAPShort = funcs.GetDeptSAPShort(driver.DriverDeptSAP)
	driver.VendorName = driver.DriverDeptSAPShort
	if driver.CarpoolName != "" {
		driver.VendorName = driver.CarpoolName
	}

	c.JSON(http.StatusOK, driver)
}

// GetDriverType godoc
// @Summary Retrieve driver types
// @Description This endpoint fetches all driver types from the database.
// @Tags Drivers
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Router /api/driver/work-type [get]
func (h *DriverHandler) GetWorkType(c *gin.Context) {
	workTypes := []gin.H{
		{"type": 1, "description": "ค้างคืน"},
		{"type": 2, "description": "ไป-กลับ"},
	}
	c.JSON(http.StatusOK, workTypes)
}
