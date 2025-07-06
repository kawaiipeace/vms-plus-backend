package handlers

import (
	"fmt"
	"net/http"
	"slices"
	"strconv"
	"strings"
	"time"
	"vms_plus_be/config"
	"vms_plus_be/funcs"
	"vms_plus_be/messages"
	"vms_plus_be/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/tealeg/xlsx"
	"gorm.io/gorm"
)

type DriverManagementHandler struct {
	Role string
}

func (h *DriverManagementHandler) SetQueryRole(user *models.AuthenUserEmp, query *gorm.DB) *gorm.DB {
	if user.EmpID == "" {
		return query
	}
	return query
}

func (h *DriverManagementHandler) SetQueryRoleDept(user *models.AuthenUserEmp, query *gorm.DB) *gorm.DB {
	if slices.Contains(user.Roles, "admin-super") {
		return query
	}
	if slices.Contains(user.Roles, "admin-region") {
		return query.Where("d.bureau_ba = ?", user.BusinessArea)
	}
	if slices.Contains(user.Roles, "admin-approval") {
		return query.Where("d.bureau_dept_sap = ?", user.BureauDeptSap)
	}
	return nil
}

func (h *DriverManagementHandler) CheckRoleDriverOwner(user *models.AuthenUserEmp, masDriverUID string) bool {
	query := config.DB.Table("vms_mas_driver_department").Where("mas_driver_uid = ?", masDriverUID)
	query = h.SetQueryRoleDept(user, query)

	//if exist return true
	var driverOwnerDeptSAP string
	if err := query.First(&driverOwnerDeptSAP).Error; err != nil {
		return false
	}
	return true
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
	user := funcs.GetAuthenUser(c, h.Role)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))    // Default: page 1
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10")) // Default: 10 items per page
	offset := (page - 1) * limit

	var drivers []models.VmsMasDriverList
	query := h.SetQueryRole(user, config.DB)
	query = h.SetQueryRoleDept(user, query).
		Model(&models.VmsMasDriverList{}).
		Select("vms_mas_driver.*, (select max(driver_license_end_date) from vms_mas_driver_license where vms_mas_driver.mas_driver_uid = vms_mas_driver_license.mas_driver_uid and vms_mas_driver_license.is_deleted = '0') as driver_license_end_date").
		Where("vms_mas_driver.is_deleted = ?", "0").Debug()

	if search := strings.ToUpper(c.Query("search")); search != "" {
		query = query.Where("UPPER(vms_mas_driver.driver_name) ILIKE ? OR UPPER(vms_mas_driver.driver_nickname) ILIKE ? OR UPPER(vms_mas_driver.driver_dept_sap_short_work) ILIKE ?", "%"+search+"%", "%"+search+"%", "%"+search+"%")
	}

	if driverDeptSAP := c.Query("driver_dept_sap_work"); driverDeptSAP != "" {
		query = query.Where("vms_mas_driver.driver_dept_sap_work = ?", driverDeptSAP)
	}

	if workType := c.Query("work_type"); workType != "" {
		query = query.Where("vms_mas_driver.work_type IN (?)", strings.Split(workType, ","))
	}

	if statusCodes := c.Query("ref_driver_status_code"); statusCodes != "" {
		query = query.Where("vms_mas_driver.ref_driver_status_code IN (?)", strings.Split(statusCodes, ","))
	}

	if isActive := c.Query("is_active"); isActive != "" {
		query = query.Where("vms_mas_driver.is_active IN (?)", strings.Split(isActive, ","))
	}

	if licenseEndDate := c.Query("driver_license_end_date"); licenseEndDate != "" {
		query = query.Where("dl.driver_license_end_date <= ?", licenseEndDate)
	}

	if approvedEndDate := c.Query("approved_job_driver_end_date"); approvedEndDate != "" {
		query = query.Where("vms_mas_driver.approved_job_driver_end_date <= ?", approvedEndDate)
	}

	orderBy := c.Query("order_by")
	orderDir := c.DefaultQuery("order_dir", "asc")

	switch orderBy {
	case "driver_name", "driver_license_end_date", "approved_job_driver_end_date":
		query = query.Order(orderBy + " " + orderDir)
	default:
		query = query.Order("driver_name " + orderDir)
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
		c.JSON(http.StatusOK, gin.H{
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

// GetDriverRunningNumber retrieves the next sequence number for a given sequence name
func GetDriverRunningNumber(sequenceName string) int {
	var nextVal int
	err := config.DB.Raw(fmt.Sprintf("SELECT nextval('%s')", sequenceName)).Scan(&nextVal).Error
	if err != nil {
		panic(fmt.Sprintf("Failed to get next sequence value for %s: %v", sequenceName, err))
	}
	return nextVal
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
	if c.IsAborted() {
		return
	}
	var driver models.VmsMasDriverRequest
	if err := c.ShouldBindJSON(&driver); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON input", "message": messages.ErrNotfound.Error()})
		return
	}

	driver.MasDriverUID = uuid.New().String()
	driver.CreatedBy = user.EmpID
	driver.UpdatedBy = user.EmpID
	driver.CreatedAt = time.Now()
	driver.UpdatedAt = time.Now()
	driver.IsDeleted = "0"
	driver.IsActive = "1"

	driver.DriverLicense.MasDriverUID = driver.MasDriverUID
	driver.DriverLicense.MasDriverLicenseUID = uuid.New().String()
	driver.DriverLicense.CreatedBy = user.EmpID
	driver.DriverLicense.UpdatedBy = user.EmpID
	driver.DriverLicense.CreatedAt = time.Now()
	driver.DriverLicense.UpdatedAt = time.Now()
	driver.DriverLicense.IsDeleted = "0"
	driver.DriverLicense.IsActive = "1"

	BCode := user.BusinessArea[0:1]
	driver.DriverID = fmt.Sprintf("DB%06d", GetDriverRunningNumber("vehicle_driver_seq_"+BCode))

	var departmentHire struct {
		DeptShort string `gorm:"column:dept_short" json:"dept_short"`
		DeptFull  string `gorm:"column:dept_full" json:"dept_full"`
	}

	if err := config.DB.Model(&models.VmsMasDepartment{}).
		Where("dept_sap = ?", driver.DriverDeptSapHire).
		Select("dept_short, dept_full").
		Find(&departmentHire).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to retrieve department details: %v", err), "message": messages.ErrInternalServer.Error()})
		return
	}
	driver.DriverDeptSapShortNameHire = departmentHire.DeptShort

	var departmentWork struct {
		DeptShort string `gorm:"column:dept_short" json:"dept_short"`
		DeptFull  string `gorm:"column:dept_full" json:"dept_full"`
	}

	if err := config.DB.Model(&models.VmsMasDepartment{}).
		Where("dept_sap = ?", driver.DriverDeptSapWork).
		Select("dept_short, dept_full").
		Find(&departmentWork).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to retrieve department details: %v", err), "message": messages.ErrInternalServer.Error()})
		return
	}

	driver.DriverDeptSapShortWork = departmentWork.DeptShort
	driver.DriverDeptSapFullWork = departmentWork.DeptFull

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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create driver", "message": messages.ErrInternalServer.Error()})
		return
	}
	if err := config.DB.Create(&driverLicense).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create driver License", "message": messages.ErrInternalServer.Error()})
		return
	}

	if err := config.DB.Create(&driverDocuments).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to driver Certificate", "message": messages.ErrInternalServer.Error()})
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
	user := funcs.GetAuthenUser(c, h.Role)
	if c.IsAborted() {
		return
	}

	masDriverUID := c.Param("mas_driver_uid")
	var driver models.VmsMasDriverResponse
	query := h.SetQueryRole(user, config.DB)
	query = h.SetQueryRoleDept(user, query)
	if err := query.Where("mas_driver_uid = ? AND is_deleted = ?", masDriverUID, "0").
		Preload("DriverStatus").
		Preload("DriverLicense", func(db *gorm.DB) *gorm.DB {
			return db.Order("driver_license_end_date DESC").Limit(1)
		}).
		Preload("DriverLicense.DriverLicenseType").
		Preload("DriverDocuments").
		First(&driver).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Driver not found", "message": messages.ErrNotfound.Error()})
		return
	}
	//
	driver.AlertDriverStatus = "ปฏิบัติงานปกติ"
	driver.AlertDriverStatusDesc = "เข้าปฏิบัติงานตามปกติ"
	driver.DriverTotalSatisfactionReview = 200

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
	if c.IsAborted() {
		return
	}
	var request, driver, result models.VmsMasDriverDetailUpdate

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	queryRole := h.SetQueryRole(user, config.DB)
	queryRole = h.SetQueryRoleDept(user, queryRole)
	if err := queryRole.First(&driver, "mas_driver_uid = ? and is_deleted = ?", request.MasDriverUID, "0").Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Driver not found", "message": messages.ErrNotfound.Error()})
		return
	}
	request.UpdatedAt = time.Now()
	request.UpdatedBy = user.EmpID

	if err := config.DB.Save(&request).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to update : %v", err), "message": messages.ErrInternalServer.Error()})
		return
	}

	if err := config.DB.First(&result, "mas_driver_uid = ?", request.MasDriverUID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Driver not found", "message": messages.ErrNotfound.Error()})
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
	if c.IsAborted() {
		return
	}
	var request, driver models.VmsMasDriverContractUpdate
	var result struct {
		models.VmsMasDriverContractUpdate
		DriverID       string `gorm:"column:driver_id" json:"driver_id"`
		DriverName     string `gorm:"column:driver_name" json:"driver_name" example:"John Doe"`
		DriverNickname string `gorm:"column:driver_nickname" json:"driver_nickname" example:"Johnny"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "message": messages.ErrInvalidJSONInput.Error()})
		return
	}

	var departmentHire struct {
		DeptShort string `gorm:"column:dept_short" json:"dept_short"`
		DeptFull  string `gorm:"column:dept_full" json:"dept_full"`
	}

	if err := config.DB.Model(&models.VmsMasDepartment{}).
		Where("dept_sap = ?", request.DriverDeptSapHire).
		Select("dept_short, dept_full").
		Find(&departmentHire).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to retrieve department details: %v", err), "message": messages.ErrInternalServer.Error()})
		return
	}
	request.DriverDeptSapShortNameHire = departmentHire.DeptShort

	var departmentWork struct {
		DeptShort string `gorm:"column:dept_short" json:"dept_short"`
		DeptFull  string `gorm:"column:dept_full" json:"dept_full"`
	}

	if err := config.DB.Model(&models.VmsMasDepartment{}).
		Where("dept_sap = ?", request.DriverDeptSapWork).
		Select("dept_short, dept_full").
		Find(&departmentWork).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to retrieve department details: %v", err), "message": messages.ErrInternalServer.Error()})
		return
	}

	request.DriverDeptSapShortWork = departmentWork.DeptShort
	request.DriverDeptSapFullWork = departmentWork.DeptFull

	queryRole := h.SetQueryRole(user, config.DB)
	queryRole = h.SetQueryRoleDept(user, queryRole)
	if err := queryRole.First(&driver, "mas_driver_uid = ? and is_deleted = ?", request.MasDriverUID, "0").Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Driver not found", "message": messages.ErrNotfound.Error()})
		return
	}
	request.UpdatedAt = time.Now()
	request.UpdatedBy = user.EmpID

	if err := config.DB.Save(&request).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to update : %v", err), "message": messages.ErrInternalServer.Error()})
		return
	}

	if err := config.DB.First(&result, "mas_driver_uid = ?", request.MasDriverUID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Driver not found", "message": messages.ErrNotfound.Error()})
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
	if c.IsAborted() {
		return
	}
	var driver models.VmsMasDriver
	var request, driverLicense models.VmsMasDriverLicenseUpdate
	var result struct {
		models.VmsMasDriverLicenseUpdate
		DriverID       string `gorm:"column:driver_id" json:"driver_id"`
		DriverName     string `gorm:"column:driver_name" json:"driver_name" example:"John Doe"`
		DriverNickname string `gorm:"column:driver_nickname" json:"driver_nickname" example:"Johnny"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "message": messages.ErrInvalidJSONInput.Error()})
		return
	}

	queryRole := h.SetQueryRole(user, config.DB)
	queryRole = h.SetQueryRoleDept(user, queryRole)
	if err := queryRole.First(&driver, "mas_driver_uid = ? and is_deleted = ?", request.MasDriverUID, "0").Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Driver not found", "message": messages.ErrNotfound.Error()})
		return
	}
	if err := config.DB.First(&driverLicense, "mas_driver_uid = ?", request.MasDriverUID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Driver license not found", "message": messages.ErrNotfound.Error()})
		return
	}
	request.UpdatedAt = time.Now()
	request.UpdatedBy = user.EmpID
	request.MasDriverUID = driver.MasDriverUID
	request.MasDriverLicenseUID = driverLicense.MasDriverLicenseUID
	if err := config.DB.Save(&request).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to update: %v", err), "message": messages.ErrInternalServer.Error()})
		return
	}
	/*
		//if exists driver_license_no update else create
		if err := config.DB.Where("driver_license_no = ?", request.DriverLicenseNo).First(&models.VmsMasDriverLicense{}).Error; err != nil {
			if err := config.DB.Save(&request).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to update: %v", err), "message": messages.ErrInternalServer.Error()})
				return
			}
		} else {
			createRequest := models.VmsMasDriverLicenseRequest{
				MasDriverLicenseUID:      uuid.New().String(),
				CreatedAt:                time.Now(),
				CreatedBy:                user.EmpID,
				UpdatedAt:                time.Now(),
				UpdatedBy:                user.EmpID,
				IsDeleted:                "0",
				IsActive:                 "1",
				MasDriverUID:             driver.MasDriverUID,
				DriverLicenseNo:          request.DriverLicenseNo,
				DriverLicenseStartDate:   request.DriverLicenseStartDate,
				DriverLicenseEndDate:     request.DriverLicenseEndDate,
				RefDriverLicenseTypeCode: request.RefDriverLicenseTypeCode,
			}
			if err := config.DB.Save(&createRequest).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to update: %v", err), "message": messages.ErrInternalServer.Error()})
				return
			}
		}
	*/

	if err := config.DB.First(&result).
		Order("updated_at desc").Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to fetch updated documents: %v", err), "message": messages.ErrInternalServer.Error()})
		return
	}
	result.DriverID = driver.DriverID
	result.DriverName = driver.DriverName
	result.DriverNickname = driver.DriverNickname
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
// @Param data body models.VmsMasDriverDocumentUpdate true "VmsMasDriverDocumentUpdate"
// @Router /api/driver-management/update-driver-documents [put]
func (h *DriverManagementHandler) UpdateDriverDocuments(c *gin.Context) {
	user := funcs.GetAuthenUser(c, h.Role)
	if c.IsAborted() {
		return
	}
	var driver models.VmsMasDriver
	var request models.VmsMasDriverDocumentUpdate
	var result struct {
		models.VmsMasDriverDocumentUpdate
		DriverID       string `gorm:"column:driver_id" json:"driver_id"`
		DriverName     string `gorm:"column:driver_name" json:"driver_name" example:"John Doe"`
		DriverNickname string `gorm:"column:driver_nickname" json:"driver_nickname" example:"Johnny"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "message": messages.ErrInvalidJSONInput.Error()})
		return
	}

	queryRole := h.SetQueryRole(user, config.DB)
	queryRole = h.SetQueryRoleDept(user, queryRole)
	if err := queryRole.First(&driver, "mas_driver_uid = ? and is_deleted = ?", request.MasDriverUID, "0").Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Driver not found", "message": messages.ErrNotfound.Error()})
		return
	}

	if err := config.DB.Where("mas_driver_uid = ? AND is_deleted = ?", request.MasDriverUID, "0").
		Delete(&models.VmsMasDriverDocument{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to delete existing documents: %v", err), "message": messages.ErrInternalServer.Error()})
		return
	}
	var masDriverLicenseUID string
	if err := config.DB.Model(&models.VmsMasDriverLicense{}).
		Where("mas_driver_uid = ? AND is_deleted = ?", request.MasDriverUID, "0").
		Pluck("mas_driver_license_uid", &masDriverLicenseUID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to retrieve mas_driver_license_uid: %v", err), "message": messages.ErrInternalServer.Error()})
		return
	}
	if request.DriverLicense.DriverDocumentFile != "" {
		if err := config.DB.Model(&models.VmsMasDriverLicense{}).
			Where("mas_driver_license_uid = ? AND is_deleted = ?", masDriverLicenseUID, "0").
			Update("driver_license_image", request.DriverLicense.DriverDocumentFile).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to update driver license image: %v", err), "message": messages.ErrInternalServer.Error()})
			return
		}
	}

	// Add new documents
	for i := range request.DriverDocuments {
		request.DriverDocuments[i].MasDriverUID = request.MasDriverUID
		request.DriverDocuments[i].MasDriverDocumentUID = uuid.New().String()
		request.DriverDocuments[i].CreatedBy = user.EmpID
		request.DriverDocuments[i].UpdatedBy = user.EmpID
		request.DriverDocuments[i].CreatedAt = time.Now()
		request.DriverDocuments[i].UpdatedAt = time.Now()
		request.DriverDocuments[i].IsDeleted = "0"
	}

	if err := config.DB.Create(&request.DriverDocuments).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update driver Documents", "message": messages.ErrInternalServer.Error()})
		return
	}

	if err := config.DB.
		Preload("DriverDocuments").
		Preload("DriverLicense").
		First(&result, "mas_driver_uid = ?", request.MasDriverUID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Driver not found", "message": messages.ErrNotfound.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Updated successfully", "result": result})
}

// UpdateDriverLayoffStatus godoc
// @Summary Update driver leave status
// @Description This endpoint updates the leave status of a driver using its unique identifier (MasDriverUID).
// @Tags Drivers-management
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Param data body models.VmsMasDriverLeaveStatusUpdate true "VmsMasDriverLeaveStatusUpdate data"
// @Router /api/driver-management/update-driver-leave-status [put]
func (h *DriverManagementHandler) UpdateDriverLeaveStatus(c *gin.Context) {
	user := funcs.GetAuthenUser(c, h.Role)
	if c.IsAborted() {
		return
	}
	var driver models.VmsMasDriver
	var request models.VmsMasDriverLeaveStatusUpdate
	var result struct {
		models.VmsMasDriverLeaveStatusUpdate
		DriverID       string `gorm:"column:driver_id" json:"driver_id"`
		DriverName     string `gorm:"column:driver_name" json:"driver_name" example:"John Doe"`
		DriverNickname string `gorm:"column:driver_nickname" json:"driver_nickname" example:"Johnny"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "message": messages.ErrInvalidJSONInput.Error()})
		return
	}

	queryRole := h.SetQueryRole(user, config.DB)
	queryRole = h.SetQueryRoleDept(user, queryRole)
	if err := queryRole.First(&driver, "mas_driver_uid = ? and is_deleted = ?", request.MasDriverUID, "0").Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Driver not found", "message": messages.ErrNotfound.Error()})
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
		c.JSON(http.StatusNotFound, gin.H{"error": "Driver not found", "message": messages.ErrNotfound.Error()})
		return

	}
	result.DriverID = driver.DriverID
	result.DriverName = driver.DriverName
	result.DriverNickname = driver.DriverNickname

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
	if c.IsAborted() {
		return
	}
	var driver models.VmsMasDriver
	var request models.VmsMasDriverIsActiveUpdate
	var result struct {
		models.VmsMasDriverIsActiveUpdate
		DriverID       string `gorm:"column:driver_id" json:"driver_id"`
		DriverName     string `gorm:"column:driver_name" json:"driver_name" example:"John Doe"`
		DriverNickname string `gorm:"column:driver_nickname" json:"driver_nickname" example:"Johnny"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "message": messages.ErrInvalidJSONInput.Error()})
		return
	}
	queryRole := h.SetQueryRole(user, config.DB)
	queryRole = h.SetQueryRoleDept(user, queryRole)
	if err := queryRole.First(&driver, "mas_driver_uid = ? and is_deleted = ?", request.MasDriverUID, "0").Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Driver not found", "message": messages.ErrNotfound.Error()})
		return
	}
	request.UpdatedAt = time.Now()
	request.UpdatedBy = user.EmpID

	if err := config.DB.Model(&driver).Update("is_active", request.IsActive).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to update: %v", err), "message": messages.ErrInternalServer.Error()})
		return
	}

	if err := config.DB.First(&result, "mas_driver_uid = ?", request.MasDriverUID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Driver not found", "message": messages.ErrNotfound.Error()})
		return
	}
	result.DriverID = driver.DriverID
	result.DriverName = driver.DriverName
	result.DriverNickname = driver.DriverNickname
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
	if c.IsAborted() {
		return
	}
	var request, driver models.VmsMasDriverDelete

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "message": messages.ErrInvalidJSONInput.Error()})
		return
	}

	queryRole := h.SetQueryRole(user, config.DB)
	queryRole = h.SetQueryRoleDept(user, queryRole)
	if err := queryRole.First(&driver, "mas_driver_uid = ? and is_deleted = ?", request.MasDriverUID, "0").Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Driver not found", "message": messages.ErrNotfound.Error()})
		return
	}

	if err := config.DB.Model(&driver).UpdateColumns(map[string]interface{}{
		"is_deleted": "1",
		"updated_by": user.EmpID,
		"updated_at": time.Now(),
	}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete", "message": messages.ErrInternalServer.Error()})
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
	if c.IsAborted() {
		return
	}
	var driver models.VmsMasDriver
	var request models.VmsMasDriverLayoffStatusUpdate
	var result struct {
		models.VmsMasDriverLayoffStatusUpdate
		DriverID       string `gorm:"column:driver_id" json:"driver_id"`
		DriverName     string `gorm:"column:driver_name" json:"driver_name" example:"John Doe"`
		DriverNickname string `gorm:"column:driver_nickname" json:"driver_nickname" example:"Johnny"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "message": messages.ErrInvalidJSONInput.Error()})
		return
	}
	queryRole := h.SetQueryRole(user, config.DB)
	queryRole = h.SetQueryRoleDept(user, queryRole)
	if err := queryRole.First(&driver, "mas_driver_uid = ? and is_deleted = ?", request.MasDriverUID, "0").Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Driver not found", "message": messages.ErrNotfound.Error()})
		return
	}
	request.UpdatedAt = time.Now()
	request.UpdatedBy = user.EmpID
	request.RefDriverStatusCode = 4

	if err := config.DB.Save(&request).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to update: %v", err), "message": messages.ErrInternalServer.Error()})
		return
	}
	// update ReplacedMasDriverUID is_replacement to 0
	if request.ReplacedMasDriverUID != "" {
		if err := config.DB.Model(&models.VmsMasDriver{}).Where("mas_driver_uid = ?", request.ReplacedMasDriverUID).
			Update("is_replacement", "0").Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to update: %v", err), "message": messages.ErrInternalServer.Error()})
			return
		}
	}
	if request.RefDriverStatusCode == 4 {
		if err := config.DB.Model(&driver).Update("is_active", "0").Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to update: %v", err), "message": messages.ErrInternalServer.Error()})
			return
		}
	}
	if err := config.DB.First(&result, "mas_driver_uid = ?", request.MasDriverUID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Driver not found", "message": messages.ErrNotfound.Error()})
		return
	}
	result.DriverID = driver.DriverID
	result.DriverName = driver.DriverName
	result.DriverNickname = driver.DriverNickname
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
	if c.IsAborted() {
		return
	}
	var driver models.VmsMasDriver
	var request models.VmsMasDriverResignStatusUpdate
	var result struct {
		models.VmsMasDriverLayoffStatusUpdate
		DriverID       string `gorm:"column:driver_id" json:"driver_id"`
		DriverName     string `gorm:"column:driver_name" json:"driver_name" example:"John Doe"`
		DriverNickname string `gorm:"column:driver_nickname" json:"driver_nickname" example:"Johnny"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "message": messages.ErrInvalidJSONInput.Error()})
		return
	}

	queryRole := h.SetQueryRole(user, config.DB)
	queryRole = h.SetQueryRoleDept(user, queryRole)
	if err := queryRole.First(&driver, "mas_driver_uid = ? and is_deleted = ?", request.MasDriverUID, "0").Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Driver not found", "message": messages.ErrNotfound.Error()})
		return
	}
	request.UpdatedAt = time.Now()
	request.UpdatedBy = user.EmpID
	request.RefDriverStatusCode = 3

	if err := config.DB.Save(&request).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to update: %v", err), "message": messages.ErrInternalServer.Error()})
		return
	}

	// update ReplacedMasDriverUID is_replacement to 0
	if request.ReplacedMasDriverUID != "" {
		if err := config.DB.Model(&models.VmsMasDriver{}).Where("mas_driver_uid = ?", request.ReplacedMasDriverUID).
			Update("is_replacement", "0").Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to update: %v", err), "message": messages.ErrInternalServer.Error()})
			return
		}
	}
	if err := config.DB.First(&result, "mas_driver_uid = ?", request.MasDriverUID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Driver not found", "message": messages.ErrNotfound.Error()})
		return
	}
	result.DriverID = driver.DriverID
	result.DriverName = driver.DriverName
	result.DriverNickname = driver.DriverNickname
	c.JSON(http.StatusOK, gin.H{"message": "Updated successfully", "result": result})
}

// GetReplacementDrivers godoc
// @Summary Get replacement drivers by name
// @Description Get a list of replacement drivers filtered by name
// @Tags Drivers-management
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Param name query string false "Driver name to search"
// @Router /api/driver-management/replacement-drivers [get]
func (h *DriverManagementHandler) GetReplacementDrivers(c *gin.Context) {
	name := strings.ToUpper(c.Query("name"))
	var drivers []models.VmsMasDriver
	user := funcs.GetAuthenUser(c, h.Role)
	query := h.SetQueryRole(user, config.DB)
	query = h.SetQueryRoleDept(user, query)
	query = query.Model(&models.VmsMasDriver{})
	query = query.Where("is_deleted = ? AND is_replacement = ?", "0", "1")
	// Apply search filter
	if name != "" {
		searchTerm := "%" + name + "%"
		query = query.Where(`
            driver_name ILIKE ? OR 
            driver_id ILIKE ?`,
			searchTerm, searchTerm)
	}

	if err := query.
		Preload("DriverStatus").
		Find(&drivers).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "message": messages.ErrInternalServer.Error()})
		return
	}
	c.JSON(http.StatusOK, drivers)
}

// GetDriverTimeLine godoc
// @Summary Get driver timeline
// @Description Get driver timeline by date range
// @Tags Drivers-management
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Param search query string false "driver_name,driver_nickname,driver_dept_sap_short_name_work to search"
// @Param start_date query string true "Start date (YYYY-MM-DD)" default(2025-05-01)
// @Param end_date query string true "End date (YYYY-MM-DD)" default(2025-05-31)
// @Param work_type query string false "work type 1: ค้างคืน, 2: ไป-กลับ Filter by multiple work_type (comma-separated, e.g., '1,2')"
// @Param ref_driver_status_code query string false "Filter by driver status code (comma-separated, e.g., '1,2')"
// @Param ref_timeline_status_id query string false "Filter by timeline status (comma-separated, e.g., '1,2')"
// @Param is_active query string false "Filter by is_active status (comma-separated, e.g., '1,0')"
// @Param page query int false "Page number (default: 1)"
// @Param limit query int false "Number of records per page (default: 10)"
// @Router /api/driver-management/timeline [get]
func (h *DriverManagementHandler) GetDriverTimeLine(c *gin.Context) {
	user := funcs.GetAuthenUser(c, h.Role)
	if c.IsAborted() {
		return
	}
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
	lastMonthDate := time.Date(startDate.Year(), startDate.Month()-1, 1, 0, 0, 0, 0, startDate.Location())

	query := h.SetQueryRole(user, config.DB)
	query = h.SetQueryRoleDept(user, query)
	query = query.Table("public.vms_mas_driver AS d").
		Select("d.*, w_thismth.job_count job_count_this_month, w_thismth.total_days total_day_this_month, w_lastmth.job_count job_count_last_month, w_lastmth.total_days total_day_last_month").
		Joins("LEFT JOIN public.vms_trn_driver_monthly_workload AS w_thismth ON w_thismth.workload_year = ? AND w_thismth.workload_month = ? AND w_thismth.driver_emp_id = d.driver_id AND w_thismth.is_deleted = ?", startDate.Year(), startDate.Month(), "0").
		Joins("LEFT JOIN public.vms_trn_driver_monthly_workload AS w_lastmth ON w_lastmth.workload_year = ? AND w_lastmth.workload_month = ? AND w_lastmth.driver_emp_id = d.driver_id AND w_lastmth.is_deleted = ?", lastMonthDate.Year(), lastMonthDate.Month(), "0").
		Where("d.is_deleted = ?", "0")

	if refTimelineStatusID := c.Query("ref_timeline_status_id"); refTimelineStatusID != "" {
		refTimelineStatusIDList := strings.Split(refTimelineStatusID, ",")
		query = query.Where(`exists (select 1 from vms_trn_request r where r.mas_carpool_driver_uid = d.mas_driver_uid AND r.ref_request_status_code != '90' AND (
						('1' in (?) AND r.ref_request_status_code < '50') OR
						('2' in (?) AND r.ref_request_status_code >= '50' AND r.ref_request_status_code < '80' AND r.ref_trip_type_code = 0) OR 
						('3' in (?) AND r.ref_request_status_code >= '50' AND r.ref_request_status_code < '80' AND r.ref_trip_type_code = 1) OR
						('4' in (?) AND r.ref_request_status_code = '80')
					)AND
						 (reserve_start_datetime BETWEEN ? AND ? OR reserve_end_datetime BETWEEN ? AND ?)
				)`, refTimelineStatusIDList, refTimelineStatusIDList, refTimelineStatusIDList, refTimelineStatusIDList, startDate, endDate, startDate, endDate)
	}

	search := strings.ToUpper(c.Query("search"))
	if search != "" {
		query = query.Where("UPPER(driver_name) ILIKE ? OR UPPER(driver_nickname) ILIKE ? OR UPPER(driver_dept_sap_short_work) ILIKE ?", "%"+search+"%", "%"+search+"%", "%"+search+"%")
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
		drivers[i].WorkLastMonth = fmt.Sprintf("%d วัน/%d งาน", drivers[i].TotalDayLastMonth, drivers[i].JobCountLastMonth)
		drivers[i].WorkThisMonth = fmt.Sprintf("%d วัน/%d งาน", drivers[i].TotalDayThisMonth, drivers[i].JobCountThisMonth)
		// Preload the driver requests for each driver
		query := config.DB.Table("vms_trn_request r").
			Select("r.*, v.vehicle_license_plate, v.vehicle_license_plate_province_short").
			Joins("LEFT JOIN vms_mas_vehicle v ON v.mas_vehicle_uid = r.mas_vehicle_uid AND v.is_deleted = ?", "0").
			Where("r.mas_carpool_driver_uid = ? AND r.is_deleted = ? AND r.ref_request_status_code != '90' AND (r.reserve_start_datetime BETWEEN ? AND ? OR r.reserve_end_datetime BETWEEN ? AND ?)", drivers[i].MasDriverUID, "0", startDate, endDate, startDate, endDate)

		if refTimelineStatusID := c.Query("ref_timeline_status_id"); refTimelineStatusID != "" {
			refTimelineStatusIDList := strings.Split(refTimelineStatusID, ",")
			query = query.Where(`
							('1' in (?) AND r.ref_request_status_code < '50') OR
							('2' in (?) AND r.ref_request_status_code >= '50' AND r.ref_request_status_code < '80' AND r.ref_trip_type_code = 0) OR
							('3' in (?) AND r.ref_request_status_code >= '50' AND r.ref_request_status_code < '80' AND r.ref_trip_type_code = 1) OR
							('4' in (?) AND r.ref_request_status_code = '80')
						`, refTimelineStatusIDList, refTimelineStatusIDList, refTimelineStatusIDList, refTimelineStatusIDList)
		}

		if err := query.Find(&drivers[i].DriverTrnRequests).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "message": messages.ErrInternalServer.Error()})
			return
		}

		// Preload the driver status for each driver

		for j := range drivers[i].DriverTrnRequests {
			drivers[i].DriverTrnRequests[j].TripDetails = []models.VmsTrnTripDetail{
				{
					TrnTripDetailUID: uuid.New().String(),
					VmsTrnTripDetailRequest: models.VmsTrnTripDetailRequest{
						TrnRequestUID:        drivers[i].DriverTrnRequests[j].TrnRequestUID,
						TripStartDatetime:    drivers[i].DriverTrnRequests[j].ReserveStartDatetime,
						TripEndDatetime:      drivers[i].DriverTrnRequests[j].ReserveEndDatetime,
						TripDeparturePlace:   drivers[i].DriverTrnRequests[j].WorkPlace,
						TripDestinationPlace: drivers[i].DriverTrnRequests[j].WorkPlace,
						TripStartMiles:       0,
						TripEndMiles:         0,
					},
				},
			}

			if drivers[i].DriverTrnRequests[j].RefRequestStatusCode == "80" {
				drivers[i].DriverTrnRequests[j].RefTimelineStatusID = "4"
				drivers[i].DriverTrnRequests[j].TimeLineStatus = "เสร็จสิ้น"
			} else if drivers[i].DriverTrnRequests[j].RefRequestStatusCode < "50" {
				drivers[i].DriverTrnRequests[j].RefTimelineStatusID = "1"
				drivers[i].DriverTrnRequests[j].TimeLineStatus = "รออนุมัติ"
			} else if drivers[i].DriverTrnRequests[j].RefTripTypeCode == 0 {
				drivers[i].DriverTrnRequests[j].RefTimelineStatusID = "2"
				drivers[i].DriverTrnRequests[j].TimeLineStatus = "ไป-กลับ"
			} else if drivers[i].DriverTrnRequests[j].RefTripTypeCode == 1 {
				drivers[i].DriverTrnRequests[j].RefTimelineStatusID = "3"
				drivers[i].DriverTrnRequests[j].TimeLineStatus = "ค้างแรม"
			}
			drivers[i].DriverTrnRequests[j].RefRequestStatusName = StatusNameMapUser[drivers[i].DriverTrnRequests[j].RefRequestStatusCode]
		}
	}
	thaiMonths := []string{"ม.ค.", "ก.พ.", "มี.ค.", "เม.ย.", "พ.ค.", "มิ.ย.", "ก.ค.", "ส.ค.", "ก.ย.", "ต.ค.", "พ.ย.", "ธ.ค."}
	lastMonth := fmt.Sprintf("%s%02d", thaiMonths[lastMonthDate.Month()-1], (lastMonthDate.Year()+543)%100)
	c.JSON(http.StatusOK, gin.H{
		"drivers":    drivers,
		"last_month": lastMonth,
		"pagination": gin.H{
			"total":      total,
			"page":       pageInt,
			"limit":      pageSizeInt,
			"totalPages": (total + int64(pageSizeInt) - 1) / int64(pageSizeInt), // Calculate total pages
		},
	})
}

// ImportDriver godoc
// @Summary Import driver data from CSV
// @Description This endpoint allows importing driver data from a CSV file.
// @Tags Drivers-management
// @Accept multipart/form-data
// @Produce json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Param file formData file true "CSV file containing driver data"
// @Router /api/driver-management/import-driver [post]
func (h *DriverManagementHandler) ImportDriver(c *gin.Context) {
	user := funcs.GetAuthenUser(c, h.Role)
	if c.IsAborted() {
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File is required", "message": messages.ErrInvalidFileType.Error()})
		return
	}

	src, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open file", "message": messages.ErrInternalServer.Error()})
		return
	}
	defer src.Close()

	// Parse CSV file
	records, err := funcs.ParseCSV(src)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid CSV format", "message": messages.ErrInvalidFileType.Error()})
		return
	}

	var drivers []models.VmsMasDriverImport
	for _, record := range records {
		masDriverUID := uuid.New().String()
		driver := models.VmsMasDriverImport{
			MasDriverUID:   masDriverUID,
			DriverName:     record["driver_name"],
			DriverNickname: record["driver_nickname"],
			DriverID:       fmt.Sprintf("DB%06d", GetDriverRunningNumber("vehicle_driver_seq_b")),
			WorkType: func() int {
				workType, _ := strconv.Atoi(record["work_type"])
				return workType
			}(),

			DriverIdentificationNo: record["driver_identification_no"],
			DriverBirthdate: func() time.Time {
				birthdate, _ := time.Parse("2006-01-02 15:04:05", record["driver_birthdate"])
				return birthdate
			}(),
			ContractNo:                 record["contract_no"],
			VendorName:                 record["vendor_name"],
			DriverDeptSapHire:          record["driver_dept_sap_hire"],
			DriverDeptSapShortNameHire: record["driver_dept_sap_short_name_hire"],
			DriverDeptSapWork:          record["driver_dept_sap_work"],
			DriverDeptSapShortNameWork: record["driver_dept_sap_short_work"],
			ApprovedJobDriverStartDate: func() time.Time {
				startDate, _ := time.Parse("2006-01-02", record["approved_job_driver_start_date"])
				return startDate
			}(),
			ApprovedJobDriverEndDate: func() time.Time {
				endDate, _ := time.Parse("2006-01-02", record["approved_job_driver_end_date"])
				return endDate
			}(),
			RefOtherUseCode: "0",
			DriverLicense: models.VmsMasDriverLicenseRequest{
				MasDriverLicenseUID:      uuid.New().String(),
				MasDriverUID:             masDriverUID,
				DriverLicenseNo:          record["driver_license_no"],
				RefDriverLicenseTypeCode: record["ref_driver_license_type_code"],
				DriverLicenseStartDate: func() time.Time {
					startDate, _ := time.Parse("2006-01-02", record["driver_license_start_date"])
					return startDate
				}(),
				DriverLicenseEndDate: func() time.Time {
					endDate, _ := time.Parse("2006-01-02", record["driver_license_end_date"])
					return endDate
				}(),
				CreatedBy: user.EmpID,
				UpdatedBy: user.EmpID,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
				IsDeleted: "0",
				IsActive:  "1",
			},
			IsReplacement: "0",
			CreatedBy:     user.EmpID,
			UpdatedBy:     user.EmpID,
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
			IsDeleted:     "0",
			IsActive:      "1",
		}
		fmt.Println(driver.DriverIdentificationNo, driver.ApprovedJobDriverStartDate)
		//check if driver already exists
		if err := config.DB.Where("driver_identification_no = ? AND approved_job_driver_start_date = ? AND is_deleted = ?", driver.DriverIdentificationNo, driver.ApprovedJobDriverStartDate, "0").First(&models.VmsMasDriver{}).Error; err != nil {
			drivers = append(drivers, driver)
		}
	}
	if len(drivers) > 0 {
		if err := config.DB.Create(&drivers).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to import drivers", "message": messages.ErrInternalServer.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Drivers imported successfully", "count": len(drivers)})
}

// ReportDriverWork godoc
// @Summary Get driver work report
// @Description Get driver work report by date range
// @Tags Drivers-management
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Param start_date query string true "Start date (YYYY-MM-DD)"
// @Param end_date query string true "End date (YYYY-MM-DD)"
// @Param show_all query string false "Show all vehicles (1 for true, 0 for false)"
// @Param mas_driver_uid body []string true "Array of driver mas_driver_uid"
// @Router /api/driver-management/work-report [post]
func (h *DriverManagementHandler) GetDriverWorkReport(c *gin.Context) {
	user := funcs.GetAuthenUser(c, h.Role)
	if c.IsAborted() {
		return
	}
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
	var masDriverUIDs []string
	if err := c.ShouldBindJSON(&masDriverUIDs); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid mas_driver_uid format", "message": messages.ErrInvalidJSONInput.Error()})
		return
	}

	var driverWorkReports []models.DriverWorkReport

	query := h.SetQueryRole(user, config.DB)
	query = h.SetQueryRoleDept(user, query)
	query = query.Table("public.vms_mas_driver AS d").
		Select(`d.mas_driver_uid, d.driver_name, d.driver_nickname, d.driver_id, 
				d.driver_dept_sap_short_work, d.driver_dept_sap_full_work, 
				r.reserve_start_datetime, r.reserve_end_datetime,
				v.vehicle_license_plate, v.vehicle_license_plate_province_short, 
				v.vehicle_license_plate_province_full, v."CarTypeDetail" AS vehicle_car_type_detail,
				td.trip_start_datetime, td.trip_end_datetime, td.trip_departure_place, td.trip_destination_place,
				td.trip_start_miles, td.trip_end_miles, td.trip_detail`).
		Joins("INNER JOIN vms_trn_request r ON r.mas_carpool_driver_uid = d.mas_driver_uid").
		Joins("INNER JOIN vms_mas_vehicle v ON v.mas_vehicle_uid = r.mas_vehicle_uid").
		Joins("LEFT JOIN vms_trn_trip_detail td ON td.trn_request_uid = r.trn_request_uid").
		Where("d.is_deleted = ? AND d.is_active = ?", "0", "1").
		Where("r.reserve_start_datetime BETWEEN ? AND ? OR r.reserve_end_datetime BETWEEN ? AND ? OR ? BETWEEN r.reserve_start_datetime AND r.reserve_end_datetime OR ? BETWEEN r.reserve_start_datetime AND r.reserve_end_datetime",
			startDate, endDate, startDate, endDate, startDate, endDate)

	query = query.Where("d.mas_driver_uid::Text IN (?)", masDriverUIDs)

	if err := query.Find(&driverWorkReports).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "message": "Internal server error"})
		return
	}

	file := xlsx.NewFile()
	sheet, err := file.AddSheet("Driver Work Reports")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create Excel sheet", "message": err.Error()})
		return
	}

	// Add header row
	headerRow := sheet.AddRow()
	headers := []string{
		"รหัสพนักงานขับรถ", "ชื่อพนักงานขับรถ", "ชื่อเล่น",
		"หน่วยงาน (ย่อ)", "หน่วยงาน (เต็ม)", "วันเวลาเริ่มงาน", "วันเวลาสิ้นสุดงาน",
		"ทะเบียนรถ", "จังหวัด (ย่อ)", "จังหวัด (เต็ม)", "ประเภทรถ",
		"วันเวลาเริ่มต้นการเดินทาง", "วันเวลาสิ้นสุดการเดินทาง", "สถานที่ออกเดินทาง", "สถานที่ถึง",
		"ระยะทางเริ่มต้น", "ระยะทางสิ้นสุด", "รายละเอียดการเดินทาง",
	}
	for _, header := range headers {
		cell := headerRow.AddCell()
		cell.Value = header
	}

	// Add data rows
	for _, report := range driverWorkReports {
		row := sheet.AddRow()
		row.AddCell().Value = report.DriverID
		row.AddCell().Value = report.DriverName
		row.AddCell().Value = report.DriverNickname
		row.AddCell().Value = report.DriverDeptSapShortWork
		row.AddCell().Value = report.DriverDeptSapFullWork
		row.AddCell().Value = report.ReserveStartDatetime.Format("2006-01-02 15:04:05")
		row.AddCell().Value = report.ReserveEndDatetime.Format("2006-01-02 15:04:05")
		row.AddCell().Value = report.VehicleLicensePlate
		row.AddCell().Value = report.VehicleLicensePlateProvinceShort
		row.AddCell().Value = report.VehicleLicensePlateProvinceFull
		row.AddCell().Value = report.VehicleCarTypeDetail
		row.AddCell().Value = report.TripStartDatetime.Format("2006-01-02 15:04:05")
		row.AddCell().Value = report.TripEndDatetime.Format("2006-01-02 15:04:05")
		row.AddCell().Value = report.TripDeparturePlace
		row.AddCell().Value = report.TripDestinationPlace
		row.AddCell().Value = fmt.Sprintf("%f", report.TripStartMiles)
		row.AddCell().Value = fmt.Sprintf("%f", report.TripEndMiles)
	}

	// Add style to the header row (bold, background color)
	headerStyle := xlsx.NewStyle()
	font := xlsx.DefaultFont()
	font.Bold = true
	headerStyle.Font = *font
	headerStyle.ApplyFont = true
	headerStyle.Font.Color = "FFFFFF"
	headerStyle.Fill = *xlsx.NewFill("solid", "4F81BD", "4F81BD")
	headerStyle.ApplyFill = true
	headerStyle.Alignment.Horizontal = "center"
	headerStyle.Alignment.Vertical = "center"
	headerStyle.ApplyAlignment = true
	headerStyle.Border = xlsx.Border{
		Left:   "thin",
		Top:    "thin",
		Bottom: "thin",
		Right:  "thin",
	}
	headerStyle.ApplyBorder = true

	// Apply style and auto-size columns for header row
	for i, cell := range headerRow.Cells {
		cell.SetStyle(headerStyle)
		// Auto-size columns (set a default width)
		col := sheet.Col(i)
		if col != nil {
			col.Width = 20
		}
	}

	// Write the file to response
	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Disposition", "attachment; filename=driver_work_reports.xlsx")
	c.Header("File-Name", fmt.Sprintf("driver_work_reports_%s_to_%s.xlsx", startDate.Format("2006-01-02"), endDate.Format("2006-01-02")))
	c.Header("Content-Transfer-Encoding", "binary")
	if err := file.Write(c.Writer); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to write Excel file", "message": err.Error()})
		return
	}
}
