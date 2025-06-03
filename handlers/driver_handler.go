package handlers

import (
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

// DriverHandlerInfo godoc
// @Summary Driver handler information
// @Description This endpoint allows a user to get driver handler information.
// @Tags Drivers
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Router /api/01-02-driver [get]
func (h *VehicleHandler) DriverHandlerInfo(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Driver handler information",
	})
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
		Preload("DriverVendor").
		Find(&drivers).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "message": messages.ErrInternalServer.Error()})
		return
	}
	for i := range drivers {
		drivers[i].Age = drivers[i].CalculateAgeInYearsMonths()
		drivers[i].Status = "ว่าง"
		if strings.HasSuffix(drivers[i].DriverID, "1") {
			drivers[i].Status = "ไม่ว่าง"
		}
		if drivers[i].WorkType == 1 {
			drivers[i].WorkTypeName = "ค้างคืน"
		} else if drivers[i].WorkType == 2 {
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

// GetBookingDrivers godoc
// @Summary Get drivers by name with pagination
// @Description Get a list of drivers filtered by name with pagination
// @Tags Drivers
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Param emp_id path string true "Employee ID (emp_id) default(700001)"
// @Param start_date query string false "Start Date (YYYY-MM-DD HH:mm:ss)" default(2025-05-30 08:00:00)
// @Param end_date query string false "End Date (YYYY-MM-DD HH:mm:ss)" default(2025-05-30 16:00:00)
// @Param name query string false "Driver name to search"
// @Param work_type query string false "work type to search (0: ไป-กลับ,1: ค้างคืน)" default(0)
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

	empID := c.Param("emp_id")
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

	var driverCanBookings []models.VmsMasDriverCanBooking
	tripTypeCode := 0
	if workType == "1" {
		tripTypeCode = 1
	}
	queryCanBooking := config.DB.Raw(`SELECT * FROM fn_get_available_drivers_view (?, ?, ?, ?, ?)`,
		startDate, endDate, bureauDeptSap, businessArea, tripTypeCode)
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
	query = query.Select("vms_mas_driver.*, w_thismth.job_count, w_thismth.total_days")
	query = query.Joins("LEFT JOIN public.vms_trn_driver_monthly_workload AS w_thismth ON w_thismth.workload_year = ? AND w_thismth.workload_month = ? AND w_thismth.driver_emp_id = vms_mas_driver.driver_id AND w_thismth.is_deleted = ?", date.Year(), date.Month(), "0")
	query = query.Where("vms_mas_driver.is_deleted = ? AND vms_mas_driver.is_replacement = ?", "0", "0")
	query = query.Where("vms_mas_driver.mas_driver_uid IN (?)", masDriverUIDs)
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
		if strings.HasSuffix(drivers[i].DriverID, "1") {
			drivers[i].Status = "ไม่ว่าง"
		}
		if drivers[i].WorkType == 1 {
			drivers[i].WorkTypeName = "ค้างคืน"
		} else if drivers[i].WorkType == 2 {
			drivers[i].WorkTypeName = "ไป-กลับ"
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
		if strings.HasSuffix(drivers[i].DriverID, "1") {
			drivers[i].Status = "ไม่ว่าง"
		}
		if drivers[i].WorkType == 1 {
			drivers[i].WorkTypeName = "ค้างคืน"
		} else if drivers[i].WorkType == 2 {
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
	var driver models.VmsMasDriver
	if err := config.DB.
		Preload("DriverLicense", func(db *gorm.DB) *gorm.DB {
			return db.Order("driver_license_end_date DESC").Limit(1)
		}).
		Preload("DriverLicense.DriverLicenseType").
		Preload("DriverStatus").
		Preload("DriverVendor").
		First(&driver, "mas_driver_uid = ?", parsedID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Request not found", "message": messages.ErrNotfound.Error()})
		return
	}
	driver.Age = driver.CalculateAgeInYearsMonths()
	if driver.WorkType == 1 {
		driver.WorkTypeName = "ค้างคืน"
	} else if driver.WorkType == 2 {
		driver.WorkTypeName = "ไป-กลับ"
	}
	driver.DriverDeptSAPShort = funcs.GetDeptSAPShort(driver.DriverDeptSAP)
	driver.WorkCount = 4
	driver.WorkDays = 16
	driver.Status = "ว่าง"
	if strings.HasSuffix(driver.DriverID, "1") {
		driver.Status = "ไม่ว่าง"
		vehicleUser1 := models.MasUserEmp{
			EmpID:        "E123",
			FullName:     "Somchai Prasert",
			DeptSAP:      "D01",
			DeptSAPShort: "Admin",
			DeptSAPFull:  "Administration Department",
			TelMobile:    "0812345678",
			TelInternal:  "5678",
		}

		vehicleUser2 := models.MasUserEmp{
			EmpID:        "E456",
			FullName:     "Nidnoi Chaiyaphum",
			DeptSAP:      "D02",
			DeptSAPShort: "HR",
			DeptSAPFull:  "Human Resources Department",
			TelMobile:    "0818765432",
			TelInternal:  "4321",
		}

		// Create two trip detail instances
		tripDetail1 := models.VmsDriverTripDetail{
			TrnRequestUID: "456e4567-e89b-12d3-a456-426614174001",
			RequestNo:     "REQ12345",
			WorkPlace:     "Bangkok",
			StartDatetime: "2025-03-29T08:00:00",
			EndDatetime:   "2025-03-29T18:00:00",
			VehicleUser:   vehicleUser1,
		}

		tripDetail2 := models.VmsDriverTripDetail{
			TrnRequestUID: "456e4567-e89b-12d3-a456-426614174002",
			RequestNo:     "REQ67890",
			WorkPlace:     "Chiang Mai",
			StartDatetime: "2025-03-30T09:00:00",
			EndDatetime:   "2025-03-30T19:00:00",
			VehicleUser:   vehicleUser2,
		}

		// Append the trip details to the DriverTripDetails slice
		driver.DriverTripDetails = append(driver.DriverTripDetails, tripDetail1, tripDetail2)

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
