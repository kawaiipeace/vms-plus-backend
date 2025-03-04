package handlers

import (
	"net/http"
	"strconv"
	"vms_plus_be/config"
	"vms_plus_be/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type VehicleHandler struct {
}

// SearchVehicles godoc
// @Summary Search vehicles by brand, license plate, and filters
// @Description Retrieves vehicles based on search text, department, and car type filters
// @Tags Vehicle
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Param search query string false "Search text (Vehicle Brand Name or License Plate)"
// @Param vehicle_owner_dept query string false "Filter by icle Owner Department"
// @Param car_type query string false "Filter by Car Type"
// @Param category_code query string false "Filter by Vehicle Category Code"
// @Param page query int false "Page number (default: 1)"
// @Param limit query int false "Number of records per page (default: 10)"
// @Router /api/vehicle/search [get]
func (h *VehicleHandler) SearchVehicles(c *gin.Context) {
	searchText := c.Query("search")            // Text search for brand name & license plate
	ownerDept := c.Query("vehicle_owner_dept") // Filter by vehicle owner department
	carType := c.Query("car_type")             // Filter by car type
	categoryCode := c.Query("category_code")   // Filter by car type

	// Pagination parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))    // Default page = 1
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10")) // Default limit = 10
	offset := (page - 1) * limit                            // Calculate offset

	var vehicles []models.VmsMasVehicle_List
	var total int64

	query := config.DB.Model(&models.VmsMasVehicle_List{})
	query = query.Where("is_deleted = '0'")
	// Apply text search (VehicleBrandName OR VehicleLicensePlate)
	if searchText != "" {
		query = query.Where("vehicle_brand_name LIKE ? OR vehicle_license_plate LIKE ?", "%"+searchText+"%", "%"+searchText+"%")
	}

	// Apply filters if provided
	if ownerDept != "" {
		query = query.Where("vehicle_owner_dept_sap = ?", ownerDept)
	}
	if carType != "" {
		query = query.Where("car_type = ?", carType)
	}
	if categoryCode != "" {
		query = query.Where("ref_vehicle_category_code = ?", categoryCode)
	}

	// Count total records
	query.Count(&total)

	// Execute query with pagination
	query.Offset(offset).Limit(limit).Find(&vehicles)

	// Respond with JSON
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

// GetVehicle godoc
// @Summary Retrieve details of a specific vehicle
// @Description This endpoint allows a user to retrieve the details of a specific vehicle associated with a booking request.
// @Tags Vehicle
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Param id path string true "Vehicle ID (mas_vehicle_uid)"
// @Router /api/vehicle/{id} [get]
func (h *VehicleHandler) GetVehicle(c *gin.Context) {
	vehicleID := c.Param("id")

	// Parse the string ID to uuid.UUID
	parsedID, err := uuid.Parse(vehicleID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid vehicle ID"})
		return
	}

	// Fetch the vehicle record from the database
	var vehicle models.VmsMasVehicle
	if err := config.DB.Preload("RefFuelType").First(&vehicle, "mas_vehicle_uid = ? AND is_deleted = '0'", parsedID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Vehicle not found"})
		return
	}

	// Return the vehicle data as a JSON response
	c.JSON(http.StatusOK, vehicle)
}

// GetCategory godoc
// @Summary Get vehicle categories
// @Description Fetches vehicle categories with optional filtering and pagination
// @Tags Vehicle
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Param page query int false "Page number (default: 1)"
// @Param limit query int false "Number of records per page (default: 10)"
// @Router /api/vehicle/category [get]
func (h *VehicleHandler) GetCategory(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))    // Default: page 1
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10")) // Default: 10 items per page
	offset := (page - 1) * limit

	var categories []models.VmsRefCategory
	var total int64

	query := config.DB.Model(&models.VmsRefCategory{})

	// Count total records
	query.Count(&total)

	// Fetch categories with pagination
	query.Offset(offset).Limit(limit).Find(&categories)

	// Respond with JSON
	c.JSON(http.StatusOK, gin.H{
		"pagination": gin.H{
			"total":      total,
			"page":       page,
			"limit":      limit,
			"totalPages": (total + int64(limit) - 1) / int64(limit), // Calculate total pages
		},
		"categories": categories,
	})
}
