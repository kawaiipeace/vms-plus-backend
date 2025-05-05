package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
	"vms_plus_be/config"
	"vms_plus_be/funcs"
	"vms_plus_be/models"

	"github.com/gin-gonic/gin"
)

type VehicleManagementHandler struct {
	Role string
}

// SearchVehicles godoc
// @Summary Get vehicles by license plate, brand, or model with pagination
// @Description Get a list of vehicles filtered by license plate, brand, or model with pagination
// @Tags Vehicle-management
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Param search query string false "vehicle_license_plate,vehicle_brand_name,vehicle_model_name to search"
// @Param vehicle_owner_dept_sap query string false "Filter by vehicle owner department SAP"
// @Param ref_vehicle_category_code query string false "Filter by Car type"
// @Param ref_vehicle_status_code query string false "Filter by vehicle status code (comma-separated, e.g., '1,2')"
// @Param ref_fuel_type_id query string false "Filter by ref_fuel_type_id"
// @Param order_by query string false "Order by vehicle_license_plate, vehicle_mileage, age,is_active"
// @Param order_dir query string false "Order direction: asc or desc"
// @Param page query int false "Page number (default: 1)"
// @Param limit query int false "Number of records per page (default: 10)"
// @Router /api/vehicle-management/search [get]
func (h *VehicleManagementHandler) SearchVehicles(c *gin.Context) {
	funcs.GetAuthenUser(c, h.Role)
	if c.IsAborted() {
		return
	}
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))    // Default: page 1
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10")) // Default: 10 items per page
	offset := (page - 1) * limit

	var vehicles []models.VmsMasVehicleManagementList

	query := config.DB.Table("public.vms_mas_vehicle AS v").
		Select(`v.mas_vehicle_uid,v.vehicle_license_plate,v.vehicle_brand_name,v.vehicle_model_name,v.ref_vehicle_type_code,
				(select max(ref_vehicle_type_name) from vms_ref_vehicle_type s where s.ref_vehicle_type_code=v.ref_vehicle_type_code) ref_vehicle_type_name,
				(select max(s.dept_short) from vms_mas_department s where s.dept_sap=d.vehicle_owner_dept_sap) vehicle_owner_dept_short,
				v.ref_vehicle_type_code,d.fleet_card_no,'1' is_tax_credit,d.vehicle_mileage,
				d.vehicle_get_date,d.ref_vehicle_status_code,v.ref_fuel_type_id,d.is_active,
				(select max(mc.carpool_name) from vms_mas_carpool mc,vms_mas_carpool_vehicle mcv where mc.mas_carpool_uid=mc.mas_carpool_uid and mcv.mas_vehicle_uid=v.mas_vehicle_uid) vehicle_carpool_name
            `).
		Joins("INNER JOIN public.vms_mas_vehicle_department AS d ON v.mas_vehicle_uid = d.mas_vehicle_uid")

	query = query.Where("v.is_deleted = ?", "0")

	search := strings.ToUpper(c.Query("search"))
	if search != "" {
		query = query.Where("v.vehicle_license_plate ILIKE ? OR vehicle_brand_name ILIKE ? OR vehicle_model_name ILIKE ?", "%"+search+"%", "%"+search+"%", "%"+search+"%")
	}

	if vehicleOwnerDeptSAP := c.Query("vehicle_owner_dept_sap"); vehicleOwnerDeptSAP != "" {
		query = query.Where("vehicle_owner_dept_sap = ?", vehicleOwnerDeptSAP)
	}

	if categoryCode := c.Query("ref_vehicle_category_code"); categoryCode != "" {
		query = query.Where("ref_vehicle_category_code = ?", categoryCode)
	}

	if statusCodes := c.Query("ref_vehicle_status_code"); statusCodes != "" {
		statusCodeList := strings.Split(statusCodes, ",")
		query = query.Where("ref_vehicle_status_code IN (?)", statusCodeList)
	}

	if fuelTypeID := c.Query("ref_fuel_type_id"); fuelTypeID != "" {
		query = query.Where("ref_fuel_type_id = ?", fuelTypeID)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	orderBy := c.Query("order_by")
	orderDir := c.Query("order_dir")
	if orderDir != "desc" {
		orderDir = "asc"
	}
	switch orderBy {
	case "vehicle_license_plate":
		query = query.Order("vehicle_license_plate " + orderDir)
	case "vehicle_mileage":
		query = query.Order("vehicle_mileage " + orderDir)
	case "age":
		query = query.Order("age " + orderDir)
	case "is_active":
		query = query.Order("is_active " + orderDir)
	default:
		query = query.Order("vehicle_license_plate " + orderDir) // Default ordering by license plate
	}

	query = query.Limit(limit).
		Offset(offset)

	if err := query.
		Preload("VmsRefFuelType").
		Find(&vehicles).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	for i := range vehicles {
		vehicles[i].Age = funcs.CalculateAge(vehicles[i].VehicleGetDate)
		funcs.TrimStringFields(&vehicles[i])
	}

	if len(vehicles) == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "No vehicles found",
			"pagination": gin.H{
				"page":       page,
				"limit":      limit,
				"totalPages": (total + int64(limit) - 1) / int64(limit), // Calculate total pages
				"vehicles":   vehicles,
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
			"vehicles": vehicles,
		})
	}
}

// UpdateVehicleIsActive godoc
// @Summary Update vehicle active status
// @Description This endpoint updates the active status of a vehicle.
// @Tags Vehicle-management
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Param data body models.VmsMasVehicleIsActiveUpdate true "VmsMasVehicleIsActiveUpdate data"
// @Router /api/vehicle-management/update-vehicle-is-active [put]
func (h *VehicleManagementHandler) UpdateVehicleIsActive(c *gin.Context) {
	user := funcs.GetAuthenUser(c, h.Role)
	if c.IsAborted() {
		return
	}
	var request, vehicle, result models.VmsMasVehicleIsActiveUpdate

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := config.DB.First(&vehicle, "mas_vehicle_uid = ? and is_deleted = ?", request.MasVehicleUID, "0").Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Vehicle not found"})
		return
	}
	request.UpdatedAt = time.Now()
	request.UpdatedBy = user.EmpID

	if err := config.DB.Model(&vehicle).Update("is_active", request.IsActive).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to update: %v", err)})
		return
	}

	if err := config.DB.First(&result, "mas_vehicle_uid = ?", request.MasVehicleUID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Vehicle not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Updated successfully", "result": result})
}
