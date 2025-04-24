package handlers

import (
	"net/http"
	"strconv"
	"strings"
	"vms_plus_be/config"
	"vms_plus_be/funcs"
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

	var vehicles []models.VmsMasVehicleList
	var total int64

	query := config.DB.Model(&models.VmsMasVehicleList{})
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
	vehicles = models.AssignVehicleImageFromIndex(vehicles)
	// Respond with JSON
	c.JSON(http.StatusOK, gin.H{
		"pagination": gin.H{
			"total":       total,
			"totalGroups": 2,
			"page":        page,
			"limit":       limit,
			"totalPages":  (total + int64(limit) - 1) / int64(limit), // Calculate total pages
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
// @Param mas_vehicle_uid path string true "MasVehicleUID (mas_vehicle_uid)"
// @Router /api/vehicle/{mas_vehicle_uid} [get]
func (h *VehicleHandler) GetVehicle(c *gin.Context) {
	vehicleID := c.Param("mas_vehicle_uid")

	// Parse the string ID to uuid.UUID
	parsedID, err := uuid.Parse(vehicleID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid vehicle ID"})
		return
	}

	// Fetch the vehicle record from the database
	var vehicle models.VmsMasVehicle
	if err := config.DB.Preload("RefFuelType").
		Preload("VehicleDepartment.VehicleUser").
		First(&vehicle, "mas_vehicle_uid = ? AND is_deleted = '0'", parsedID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Vehicle not found"})
		return
	}
	vehicle.Age = vehicle.CalculateAge()
	vehicle.VehicleImgs = []string{
		"http://pntdev.ddns.net:28089/VMS_PLUS/PIX/cars/Vehicle-1.svg",
		"http://pntdev.ddns.net:28089/VMS_PLUS/PIX/cars/Vehicle-2.svg",
		"http://pntdev.ddns.net:28089/VMS_PLUS/PIX/cars/Vehicle-3.svg",
	}
	vehicle.VehicleDepartment.VehicleUser.EmpID = vehicle.VehicleDepartment.VehicleUserEmpID
	vehicle.VehicleDepartment.VehicleUser.FullName = vehicle.VehicleDepartment.VehicleUserEmpName
	vehicle.VehicleDepartment.VehicleUser.DeptSAP = vehicle.VehicleDepartment.VehicleOwnerDeptSap
	vehicle.VehicleDepartment.VehicleUser.DeptSAPFull = vehicle.VehicleDepartment.OwnerDeptName
	vehicle.VehicleDepartment.VehicleUser.DeptSAPShort = "สฟฟ.มสด.4(ล)"
	if strings.TrimSpace(vehicle.VehicleLicensePlate) == "7กษ 4377" {
		vehicle.IsAdminChooseDriver = true
	}
	// Return the vehicle data as a JSON response
	funcs.TrimStringFields(&vehicle)
	c.JSON(http.StatusOK, vehicle)
}

// GetTypes godoc
// @Summary Get vehicle types
// @Description Fetches all vehicle types, optionally filtered by name
// @Tags Vehicle
// @Accept json
// @Produce json
// @Param name query string false "Filter by vehicle type name (partial match)"
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Router /api/vehicle/types [get]
func (h *VehicleHandler) GetTypes(c *gin.Context) {
	var types []models.VmsRefVehicleType
	name := c.Query("name") // Get the 'name' query parameter

	// Build the query
	query := config.DB
	if name != "" {
		query = query.Where("ref_vehicle_type_name LIKE ?", "%"+name+"%")
	}

	// Execute the query
	query.Find(&types)
	types = models.AssignTypeImageFromIndex(types)

	// Respond with JSON
	c.JSON(http.StatusOK, types)
}

// GetDepartments godoc
// @Summary Get department list
// @Description Fetches a list of departments grouped by dept_sap, including dept_short and dept_full
// @Tags Vehicle
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Router /api/vehicle/departments [get]
func (h *VehicleHandler) GetDepartments(c *gin.Context) {
	var departments []struct {
		DeptSap   string `json:"dept_sap"`
		DeptShort string `json:"dept_short"`
		DeptFull  string `json:"dept_full"`
	}

	// SQL query to group and join tables
	query := `
        SELECT 
            vmd.dept_sap,
            vmd.dept_short,
            vmd.dept_full
        FROM 
            vms_mas_vehicle_department vvd
        JOIN 
            vms_mas_department vmd
        ON 
            vvd.vehicle_owner_dept_sap = vmd.dept_sap
        GROUP BY 
            vmd.dept_sap, vmd.dept_short, vmd.dept_full
    `

	// Execute the query
	if err := config.DB.Raw(query).Scan(&departments).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch departments"})
		return
	}

	// Respond with the result
	c.JSON(http.StatusOK, departments)
}

// GetVehicleInfo godoc
// @Summary Get vehicle info
// @Description Get vehicle info
// @Tags Vehicle
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Param mas_vehicle_uid path string true "MasVehicleUID (mas_vehicle_uid)"
// @Router /api/vehicle-info/{mas_vehicle_uid} [get]
func (h *VehicleHandler) GetVehicleInfo(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"number_of_available_drivers": 2})
}
