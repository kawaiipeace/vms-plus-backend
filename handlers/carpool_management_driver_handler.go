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
	masCarpoolUID := c.Param("mas_carpool_uid")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))    // Default: page 1
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10")) // Default: 10 items per page
	offset := (page - 1) * limit

	var drivers []models.VmsMasCarpoolDriverList
	query := config.DB.Table("vms_mas_carpool_driver cpd").
		Select(
			`cpd.mas_carpool_driver_uid,
			 cpd.mas_carpool_uid,
			 d.mas_driver_uid,
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
		Joins("LEFT JOIN vms_mas_driver d ON d.mas_driver_uid = cpd.mas_driver_uid").
		Where("cpd.mas_carpool_uid = ? AND cpd.is_deleted = ?", masCarpoolUID, "0")

	search := strings.ToUpper(c.Query("search"))
	if search != "" {
		query = query.Where("UPPER(d.driver_name) LIKE ? OR UPPER(d.driver_license_no) LIKE ?", "%"+search+"%", "%"+search+"%")
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

	var requests []models.VmsMasCarpoolDriver
	if err := c.ShouldBindJSON(&requests); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
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
		"message": "Carpool drivers created successfully",
		"data":    requests,
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
	masCarpoolDriverUID := c.Param("mas_carpool_driver_uid")

	var driver models.VmsMasCarpoolDriver
	if err := config.DB.Where("mas_carpool_driver_uid = ? AND is_deleted = ?", masCarpoolDriverUID, "0").First(&driver).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Carpool driver not found"})
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

	c.JSON(http.StatusOK, gin.H{"message": "Carpool driver deleted successfully"})
}

// GetMasDriverDetail godoc
// @Summary Retrieve a specific driver
// @Description This endpoint fetches details of a specific driver using its unique identifier (MasDriverUID).
// @Tags Carpool-management
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Param mas_driver_uid path string true "MasDriverUID (mas_driver_uid)"
// @Router /api/carpool-management/driver-mas-detail/{mas_driver_uid} [get]
func (h *CarpoolManagementHandler) GetMasDriverDetail(c *gin.Context) {
	masDriverUID := c.Param("mas_driver_uid")

	var driver models.VmsMasCarpoolDriverDetail
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
		Where("d.mas_driver_uid = ? AND d.is_deleted = ?", masDriverUID, "0")

	if err := query.First(&driver).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Driver not found"})
		return
	}

	funcs.TrimStringFields(&driver)

	c.JSON(http.StatusOK, driver)
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

	var request models.VmsMasCarpoolDriverActive
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var driver models.VmsMasCarpoolDriver
	if err := config.DB.Where("mas_carpool_driver_uid = ? AND is_deleted = ?", request.MasCarpoolDriverUID, "0").First(&driver).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Carpool driver not found"})
		return
	}

	driver.IsActive = request.IsActive
	driver.UpdatedAt = time.Now()
	driver.UpdatedBy = user.EmpID

	if err := config.DB.Save(&driver).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to update active status: %v", err)})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Carpool driver active status updated successfully", "data": request})
}
