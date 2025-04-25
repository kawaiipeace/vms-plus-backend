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
	"gorm.io/gorm"
)

type DriverManagementHandler struct {
	Role string
}

// SearchDrivers godoc
// @Summary Get drivers by name with pagination
// @Description Get a list of drivers filtered by name with pagination
// @Tags Drivers-management
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Param search query string false "driver_name,driver_nickname,driver_dept_sap_short_name_work to search"
// @Param driver_dept_sap_work query string false "Filter by driver department SAP"
// @Param work_type query string false "work type 1: ค้างคืน, 2: ไป-กลับ Filter by multiple work_type (comma-separated, e.g., '1,2')"
// @Param ref_driver_status_code query string false "Filter by driver status code (comma-separated, e.g., '1,2')"
// @Param is_active query string false "Filter by is_active status (comma-separated, e.g., '1,0')"
// @Param driver_license_end_date query string false "Filter by  driver license end date (YYYY-MM-DD)"
// @Param approved_job_driver_end_date query string false "Filter by approved job driver end date (YYYY-MM-DD)"
// @Param order_by query string false "Order by driver_name, driver_license_end_date, approved_job_driver_end_date,"
// @Param order_dir query string false "Order direction: asc or desc"
// @Param page query int false "Page number (default: 1)"
// @Param limit query int false "Number of records per page (default: 10)"
// @Router /api/driver-management/search [get]
func (h *DriverManagementHandler) SearchDrivers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))    // Default: page 1
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10")) // Default: 10 items per page
	offset := (page - 1) * limit

	var drivers []models.VmsMasDriverList
	query := config.DB.Model(&models.VmsMasDriverList{})
	query = query.Select("vms_mas_driver.*,(select max(driver_license_end_date) from vms_mas_driver_license where mas_driver_uid = vms_mas_driver.mas_driver_uid) as driver_license_end_date")
	query = query.Where("is_deleted = ?", "0")

	name := strings.ToUpper(c.Query("name"))
	if name != "" {
		query = query.Where("UPPER(driver_name) LIKE ? OR UPPER(driver_nickname) LIKE ? OR UPPER(driver_dept_sap_short_name_work) LIKE ?", "%"+name+"%", "%"+name+"%", "%"+name+"%")
	}

	if driverDeptSAP := c.Query("driver_dept_sap_work"); driverDeptSAP != "" {
		query = query.Where("driver_dept_sap_work = ?", driverDeptSAP)
	}

	if workType := c.Query("work_type"); workType != "" {
		workTypes := strings.Split(workType, ",")
		query = query.Where("work_type IN (?)", workTypes)
	}

	if statusCodes := c.Query("ref_driver_status_code"); statusCodes != "" {
		statusCodeList := strings.Split(statusCodes, ",")
		query = query.Where("ref_driver_status_code IN (?)", statusCodeList)
	}

	if isActive := c.Query("is_active"); isActive != "" {
		isActiveList := strings.Split(isActive, ",")
		query = query.Where("is_active IN (?)", isActiveList)
	}

	if licenseEndDate := c.Query("driver_license_end_date"); licenseEndDate != "" {
		query = query.Where("driver_license_end_date <= ?", licenseEndDate)
	}

	if approvedEndDate := c.Query("approved_job_driver_end_date"); approvedEndDate != "" {
		query = query.Where("approved_job_driver_end_date <= ?", approvedEndDate)
	}

	orderBy := c.Query("order_by")
	orderDir := c.Query("order_dir")
	if orderDir != "desc" {
		orderDir = "asc"
	}
	switch orderBy {
	case "driver_name":
		query = query.Order("driver_name " + orderDir)
	case "driver_license_end_date":
		query = query.Order("driver_license_end_date " + orderDir)
	case "approved_job_driver_end_date":
		query = query.Order("approved_job_driver_end_date " + orderDir)
	default:
		query = query.Order("driver_name " + orderDir) // Default ordering by name
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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

// CreateDriver godoc
// @Summary Create a new driver
// @Description This endpoint allows creating a new driver.
// @Tags Drivers-management
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Param data body models.VmsMasDriverRequest true "VmsMasDriverRequest data"
// @Router /api/driver-management/create-driver [post]
func (h *DriverManagementHandler) CreateDriver(c *gin.Context) {
	user := funcs.GetAuthenUser(c, h.Role)
	var driver models.VmsMasDriverRequest
	if err := c.ShouldBindJSON(&driver); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON input"})
		return
	}

	driver.MasDriverUID = uuid.New().String()
	driver.CreatedBy = user.EmpID
	driver.UpdatedBy = user.EmpID
	driver.CreatedAt = time.Now()
	driver.UpdatedAt = time.Now()
	driver.IsDeleted = "0"
	driver.IsActive = "1"
	driver.IsReplacement = "0"
	driver.RefOtherUseCode = "0"

	driver.DriverLicense.MasDriverUID = driver.MasDriverUID
	driver.DriverLicense.MasDriverLicenseUID = uuid.New().String()
	driver.DriverLicense.CreatedBy = user.EmpID
	driver.DriverLicense.UpdatedBy = user.EmpID
	driver.DriverLicense.CreatedAt = time.Now()
	driver.DriverLicense.UpdatedAt = time.Now()
	driver.DriverLicense.IsDeleted = "0"
	driver.DriverLicense.IsActive = "1"

	for i := range driver.DriverDocuments {
		driver.DriverDocuments[i].MasDriverUID = driver.MasDriverUID
		driver.DriverDocuments[i].MasDriverDocumentUID = uuid.New().String()
		driver.DriverDocuments[i].CreatedBy = user.EmpID
		driver.DriverDocuments[i].UpdatedBy = user.EmpID
		driver.DriverDocuments[i].CreatedAt = time.Now()
		driver.DriverDocuments[i].UpdatedAt = time.Now()
		driver.DriverDocuments[i].IsDeleted = "0"
	}

	driverLicense := driver.DriverLicense
	driverDocuments := driver.DriverDocuments

	driver.DriverLicense = models.VmsMasDriverLicenseRequest{}
	driver.DriverDocuments = []models.VmsMasDriverDocument{}

	if err := config.DB.Create(&driver).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create driver"})
		return
	}
	if err := config.DB.Create(&driverLicense).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create driver License"})
		return
	}

	if err := config.DB.Create(&driverDocuments).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to driver Certificate"})
		return
	}

	// Return success response
	c.JSON(http.StatusCreated, gin.H{"message": "Driver created successfully",
		"data":           driver,
		"mas_driver_uid": driver.MasDriverUID,
	})
}

// GetDriver godoc
// @Summary Retrieve a specific driver
// @Description This endpoint fetches details of a specific driver using its unique identifier (MasDriverUID).
// @Tags Drivers-management
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Param mas_driver_uid path string true "MasDriverUID (mas_driver_uid)"
// @Router /api/driver-management/driver/{mas_driver_uid} [get]
func (h *DriverManagementHandler) GetDriver(c *gin.Context) {
	masDriverUID := c.Param("mas_driver_uid")
	var driver models.VmsMasDriverResponse

	if err := config.DB.Where("mas_driver_uid = ? AND is_deleted = ?", masDriverUID, "0").
		Preload("DriverStatus").
		Preload("DriverLicense", func(db *gorm.DB) *gorm.DB {
			return db.Order("driver_license_end_date DESC").Limit(1)
		}).
		Preload("DriverLicense.DriverLicenseType").
		Preload("DriverCertificate").
		First(&driver).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Driver not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"driver": driver})
}

// UpdateDriverDetail godoc
// @Summary Update driver details
// @Description This endpoint updates the details of a driver using its unique identifier (MasDriverUID).
// @Tags Drivers-management
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Param data body models.VmsMasDriverDetailUpdate true "VmsMasDriverDetailUpdate data"
// @Router /api/driver-management/update-driver-detail [put]
func (h *DriverManagementHandler) UpdateDriverDetail(c *gin.Context) {
	user := funcs.GetAuthenUser(c, h.Role)
	var request, driver, result models.VmsMasDriverDetailUpdate

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := config.DB.First(&driver, "mas_driver_uid = ? and is_deleted = ?", request.MasDriverUID, "0").Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Driver not found"})
		return
	}
	request.UpdatedAt = time.Now()
	request.UpdatedBy = user.EmpID

	if err := config.DB.Save(&request).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to update : %v", err)})
		return
	}

	if err := config.DB.First(&result, "mas_driver_uid = ?", request.MasDriverUID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Driver not found"})
		return

	}
	c.JSON(http.StatusOK, gin.H{"message": "Updated successfully", "result": result})
}

// UpdateDriverContract godoc
// @Summary Update driver contract details
// @Description This endpoint updates the contract details of a driver using its unique identifier (MasDriverUID).
// @Tags Drivers-management
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Param data body models.VmsMasDriverContractUpdate true "VmsMasDriverContractUpdate data"
// @Router /api/driver-management/update-driver-contract [put]
func (h *DriverManagementHandler) UpdateDriverContract(c *gin.Context) {
	user := funcs.GetAuthenUser(c, h.Role)
	var request, driver, result models.VmsMasDriverContractUpdate

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := config.DB.First(&driver, "mas_driver_uid = ? and is_deleted = ?", request.MasDriverUID, "0").Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Driver not found"})
		return
	}
	request.UpdatedAt = time.Now()
	request.UpdatedBy = user.EmpID

	if err := config.DB.Save(&request).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to update : %v", err)})
		return
	}

	if err := config.DB.First(&result, "mas_driver_uid = ?", request.MasDriverUID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Driver not found"})
		return

	}
	c.JSON(http.StatusOK, gin.H{"message": "Updated successfully", "result": result})
}

// UpdateDriverLicense godoc
// @Summary Update driver license details
// @Description This endpoint updates the license details of a driver using its unique identifier (MasDriverUID).
// @Tags Drivers-management
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Param data body models.VmsMasDriverLicenseUpdate true "VmsMasDriverLicenseUpdate data"
// @Router /api/driver-management/update-driver-license [put]
func (h *DriverManagementHandler) UpdateDriverLicense(c *gin.Context) {
	user := funcs.GetAuthenUser(c, h.Role)
	var request, driver, result models.VmsMasDriverLicenseUpdate

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := config.DB.First(&driver, "mas_driver_license_uid = ?", request.MasDriverLicenseUID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Driver license not found"})
		return
	}
	request.UpdatedAt = time.Now()
	request.UpdatedBy = user.EmpID
	request.MasDriverUID = driver.MasDriverUID
	if err := config.DB.Save(&request).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to update: %v", err)})
		return
	}

	if err := config.DB.First(&result, "mas_driver_license_uid = ?", request.MasDriverLicenseUID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Driver license not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Updated successfully", "result": result})
}

// UpdateDriverDocuments godoc
// @Summary Update driver document details
// @Description This endpoint updates the document details of a driver using its unique identifier (MasDriverUID).
// @Tags Drivers-management
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Param mas_driver_uid path string true "MasDriverUID (mas_driver_uid)"
// @Param data body []models.VmsMasDriverDocument true "VmsMasDriverDocument array"
// @Router /api/driver-management/update-driver-documents/{mas_driver_uid} [put]
func (h *DriverManagementHandler) UpdateDriverDocuments(c *gin.Context) {
	user := funcs.GetAuthenUser(c, h.Role)
	var request, result []models.VmsMasDriverDocument
	var driver models.VmsMasDriver
	masDriverUID := c.Param("mas_driver_uid")

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := config.DB.First(&driver, "mas_driver_uid = ? and is_deleted = ?", masDriverUID, "0").Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Driver document not found"})
		return
	}

	if err := config.DB.Where("mas_driver_uid = ? AND is_deleted = ?", masDriverUID, "0").
		Delete(&models.VmsMasDriverDocument{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to delete existing documents: %v", err)})
		return
	}

	// Add new documents
	for i := range request {
		request[i].MasDriverUID = masDriverUID
		request[i].MasDriverDocumentUID = uuid.New().String()
		request[i].CreatedBy = user.EmpID
		request[i].UpdatedBy = user.EmpID
		request[i].CreatedAt = time.Now()
		request[i].UpdatedAt = time.Now()
		request[i].IsDeleted = "0"
	}

	if err := config.DB.Create(&request).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to update: %v", err)})
		return
	}

	if err := config.DB.Where("mas_driver_uid = ? AND is_deleted = ?", masDriverUID, "0").Find(&result).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to fetch updated documents: %v", err)})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Updated successfully", "result": result})
}

func (h *DriverManagementHandler) UpdateDriverLeaveStatus(c *gin.Context) {
	user := funcs.GetAuthenUser(c, h.Role)
	var driver models.VmsMasDriver
	var request, result models.VmsMasDriverLeaveStatusUpdate

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := config.DB.First(&driver, "mas_driver_uid = ? and is_deleted = ?", request.MasDriverUID, "0").Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Driver not found"})
		return
	}
	request.TrnDriverLeaveUID = uuid.NewString()
	request.CreatedAt = time.Now()
	request.CreatedBy = user.EmpID
	request.UpdatedAt = time.Now()
	request.UpdatedBy = user.EmpID
	request.RefDriverStatusCode = 2
	request.IsDeleted = "0"

	if err := config.DB.Save(&request).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to update : %v", err)})
		return
	}

	if err := config.DB.First(&result, "trn_driver_leave_uid = ?", request.TrnDriverLeaveUID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Driver not found"})
		return

	}
	c.JSON(http.StatusOK, gin.H{"message": "Updated successfully", "result": result})
}

// UpdateDriverIsActive godoc
// @Summary Update driver active status
// @Description This endpoint updates the active status of a driver using its unique identifier (MasDriverUID).
// @Tags Drivers-management
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Param data body models.VmsMasDriverIsActiveUpdate true "VmsMasDriverIsActiveUpdate data"
// @Router /api/driver-management/update-driver-is-active [put]
func (h *DriverManagementHandler) UpdateDriverIsActive(c *gin.Context) {
	user := funcs.GetAuthenUser(c, h.Role)
	var request, driver, result models.VmsMasDriverIsActiveUpdate

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := config.DB.First(&driver, "mas_driver_uid = ? and is_deleted = ?", request.MasDriverUID, "0").Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Driver not found"})
		return
	}
	request.UpdatedAt = time.Now()
	request.UpdatedBy = user.EmpID

	if err := config.DB.Model(&driver).Update("is_active", request.IsActive).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to update: %v", err)})
		return
	}

	if err := config.DB.First(&result, "mas_driver_uid = ?", request.MasDriverUID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Driver not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Updated successfully", "result": result})
}

// DeleteDriver godoc
// @Summary Delete a driver
// @Description This endpoint deletes a driver using its unique identifier (MasDriverUID).
// @Tags Drivers-management
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Param data body models.VmsMasDriverDelete true "VmsMasDriverDelete data"
// @Router /api/driver-management/delete-driver [delete]
func (h *DriverManagementHandler) DeleteDriver(c *gin.Context) {
	user := funcs.GetAuthenUser(c, h.Role)
	var request, driver models.VmsMasDriverDelete

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := config.DB.First(&driver, "mas_driver_uid = ? and is_deleted = ? and driver_name = ?", request.MasDriverUID, "0", request.DriverName).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Driver not found"})
		return
	}

	if err := config.DB.Model(&driver).UpdateColumns(map[string]interface{}{
		"is_deleted": "1",
		"updated_by": user.EmpID,
		"updated_at": time.Now(),
	}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete"})
	}

	c.JSON(http.StatusOK, gin.H{"message": "Deleted successfully"})
}

// UpdateDriverLayoffStatus godoc
// @Summary Update driver layoff status
// @Description This endpoint updates the layoff status of a driver using its unique identifier (MasDriverUID).
// @Tags Drivers-management
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Param data body models.VmsMasDriverLayoffStatusUpdate true "VmsMasDriverLayoffStatusUpdate data"
// @Router /api/driver-management/update-driver-layoff-status [put]
func (h *DriverManagementHandler) UpdateDriverLayoffStatus(c *gin.Context) {
	user := funcs.GetAuthenUser(c, h.Role)
	var request, driver, result models.VmsMasDriverLayoffStatusUpdate

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := config.DB.First(&driver, "mas_driver_uid = ? and is_deleted = ?", request.MasDriverUID, "0").Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Driver not found"})
		return
	}
	request.UpdatedAt = time.Now()
	request.UpdatedBy = user.EmpID
	request.RefDriverStatusCode = 4

	if err := config.DB.Save(&request).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to update: %v", err)})
		return
	}

	if err := config.DB.First(&result, "mas_driver_uid = ?", request.MasDriverUID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Driver not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Updated successfully", "result": result})
}

// UpdateDriverResignStatus godoc
// @Summary Update driver resign status
// @Description This endpoint updates the resign status of a driver using its unique identifier (MasDriverUID).
// @Tags Drivers-management
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Param data body models.VmsMasDriverResignStatusUpdate true "VmsMasDriverResignStatusUpdate data"
// @Router /api/driver-management/update-driver-resign-status [put]
func (h *DriverManagementHandler) UpdateDriverResignStatus(c *gin.Context) {
	user := funcs.GetAuthenUser(c, h.Role)
	var request, driver, result models.VmsMasDriverResignStatusUpdate

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := config.DB.First(&driver, "mas_driver_uid = ? and is_deleted = ?", request.MasDriverUID, "0").Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Driver not found"})
		return
	}
	request.UpdatedAt = time.Now()
	request.UpdatedBy = user.EmpID
	request.RefDriverStatusCode = 3

	if err := config.DB.Save(&request).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to update: %v", err)})
		return
	}

	if err := config.DB.First(&result, "mas_driver_uid = ?", request.MasDriverUID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Driver not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Updated successfully", "result": result})
}
