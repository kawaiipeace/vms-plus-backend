package handlers

import (
	"net/http"
	"strconv"
	"vms_plus_be/config"
	"vms_plus_be/models"

	"github.com/gin-gonic/gin"
)

type DriverHandler struct {
}

// GetDriversByName godoc
// @Summary Get drivers by name with pagination
// @Description Get a list of drivers filtered by name with pagination
// @Tags Drivers
// @Accept json
// @Produce json
// @Param name query string true "Driver name to search"
// @Param page query int false "Page number (default: 1)"
// @Param limit query int false "Number of records per page (default: 10)"
// @Router /api/driver/search [get]
func (h *DriverHandler) GetDriversByName(c *gin.Context) {
	name := c.Query("name")
	if name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Name parameter is required"})
		return
	}

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

	// Convert Driver to DriverResponse (inherit model + add Age)
	var response []models.VmsMasDriver
	for _, driver := range drivers {
		response = append(response, models.VmsMasDriver{
			Age: driver.CalculateAgeInYearsMonths(), // Add formatted age
		})
	}

	if len(response) == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "No drivers found",
			"pagination": gin.H{
				"page":       page,
				"limit":      limit,
				"totalPages": (total + int64(limit) - 1) / int64(limit), // Calculate total pages
				"drivers":    response,
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
			"drivers": response,
		})
	}
}
