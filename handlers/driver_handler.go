package handlers

import (
	"net/http"
	"strconv"
	"strings"
	"vms_plus_be/config"
	"vms_plus_be/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type DriverHandler struct {
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
// @Param page query int false "Page number (default: 1)"
// @Param limit query int false "Number of records per page (default: 10)"
// @Router /api/driver/search [get]
func (h *DriverHandler) GetDrivers(c *gin.Context) {
	name := c.Query("name")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))    // Default: page 1
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10")) // Default: 10 items per page
	offset := (page - 1) * limit

	var drivers []models.VmsMasDriver

	// Get the total count of drivers (for pagination)
	var total int64
	config.DB.Model(&models.VmsMasDriver{}).Where("driver_name LIKE ?", "%"+name+"%").Count(&total)

	// Fetch the drivers with pagination
	result := config.DB.Where("driver_name LIKE ?", "%"+name+"%").Limit(limit).Offset(offset).Find(&drivers)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	for i := range drivers {
		drivers[i].Age = drivers[i].CalculateAgeInYearsMonths()
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
	name := c.Query("name")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))    // Default: page 1
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10")) // Default: 10 items per page
	offset := (page - 1) * limit

	var drivers []models.VmsMasDriver

	// Get the total count of drivers (for pagination)
	var total int64
	config.DB.Model(&models.VmsMasDriver{}).Where("driver_name LIKE ?", "%"+name+"%").Count(&total)

	// Fetch the drivers with pagination
	result := config.DB.Where("driver_name LIKE ?", "%"+name+"%").Limit(limit).Offset(offset).Find(&drivers)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	for i := range drivers {
		drivers[i].Age = drivers[i].CalculateAgeInYearsMonths()
		drivers[i].Status = "ว่าง"
		if strings.HasSuffix(drivers[i].DriverID, "1") {
			drivers[i].Status = "ไม่ว่าง"
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

// GetDriver godoc
// @Summary Retrieve a specific driver
// @Description This endpoint fetches details of a driver using its unique identifier (MasDriverUID).
// @Tags Drivers
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Param id path string true "MasDriverUID (MasDriverUID)"
// @Router /api/driver/{id} [get]
func (h *DriverHandler) GetDriver(c *gin.Context) {
	//funcs.GetAuthenUser(c, h.Role)
	masDriverUID := c.Param("id")
	parsedID, err := uuid.Parse(masDriverUID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid MasDriverUID"})
		return
	}
	var driver models.VmsMasDriver
	if err := config.DB.
		First(&driver, "mas_driver_uid = ?", parsedID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Request not found"})
		return
	}
	driver.Age = driver.CalculateAgeInYearsMonths()
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
