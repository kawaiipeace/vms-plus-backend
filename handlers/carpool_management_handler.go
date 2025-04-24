package handlers

import (
	"fmt"
	"log"
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

type CarpoolManagementHandler struct {
	Role string
}

// SearchCarpools godoc
// @Summary Search carpool management
// @Description Search carpool management by criteria
// @Tags Carpool-management
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Param search query string false "Search query for carpool_name or emp_name"
// @Param is_active query string false "Filter by is_active status (comma-separated, e.g., '1,0')"
// @Param order_by query string false "Order by fields: carpool_name, number_of_drivers, number_of_vehicles"
// @Param order_dir query string false "Order direction: asc or desc"
// @Param page query int false "Page number (default: 1)"
// @Param limit query int false "Number of records per page (default: 10)"
// @Router /api/carpool-management/search [get]
func (h *CarpoolManagementHandler) SearchCarpools(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))    // Default: page 1
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10")) // Default: 10 items per page
	offset := (page - 1) * limit

	var carpools []models.VmsMasCarpoolList
	query := config.DB.Model(&models.VmsMasCarpoolList{})
	query = query.Table("vms_mas_carpool cp").Select(`cp.*,
		(select count(*) from vms_mas_carpool_driver cpd where is_deleted='0' and cpd.mas_carpool_uid=cp.mas_carpool_uid) number_of_drivers,
		(select count(*) from vms_mas_carpool_vehicle cpv where is_deleted='0' and cpv.mas_carpool_uid=cp.mas_carpool_uid) number_of_vehicles
	`)
	search := strings.ToUpper(c.Query("search"))
	if search != "" {
		query = query.Where("UPPER(carpool_name) LIKE ? OR UPPER(emp_name) LIKE ?", "%"+search+"%", "%"+search+"%")
	}
	if isActive := c.Query("is_active"); isActive != "" {
		isActiveList := strings.Split(isActive, ",")
		query = query.Where("is_active IN (?)", isActiveList)
	}
	orderBy := c.Query("order_by")
	orderDir := c.Query("order_dir")
	if orderDir != "desc" {
		orderDir = "asc"
	}
	switch orderBy {
	case "carpool_name":
		query = query.Order("carpool_name " + orderDir)
	case "number_of_drivers":
		query = query.Order("number_of_drivers " + orderDir)
	case "number_of_vehicles":
		query = query.Order("number_of_vehicles " + orderDir)
	default:
		query = query.Order("carpool_name " + orderDir) // Default ordering by carpool_name
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	query = query.Limit(limit).
		Offset(offset)

	if err := query.Find(&carpools).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	for i := range carpools {
		if carpools[i].IsActive == "1" {
			carpools[i].CarpoolStatus = "เปิด"
		} else if carpools[i].IsActive == "0" {
			carpools[i].CarpoolStatus = "ปิด"
		}
	}
	if len(carpools) == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "No carpools found",
			"pagination": gin.H{
				"page":       page,
				"limit":      limit,
				"totalPages": (total + int64(limit) - 1) / int64(limit), // Calculate total pages
				"carpools":   carpools,
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
			"carpools": carpools,
		})
	}
}

// CreateCarpool godoc
// @Summary Create a new carpool
// @Description Create a new carpool and save it to the database
// @Tags Carpool-management
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Param carpool body models.VmsMasCarpoolRequest true "VmsMasCarpoolRequest data"
// @Router /api/carpool-management/create [post]
func (h *CarpoolManagementHandler) CreateCarpool(c *gin.Context) {
	user := funcs.GetAuthenUser(c, h.Role)
	var carpool models.VmsMasCarpoolRequest

	if err := c.ShouldBindJSON(&carpool); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	carpool.MasCarpoolUID = uuid.New().String()
	carpool.CarpoolDeptSap = ""
	carpool.CarpoolType = ""
	carpool.IsHaveDriverForCarpool = "0"
	carpool.IsMustPassStatus30 = "0"
	carpool.IsMustPassStatus40 = "0"
	carpool.IsMustPassStatus50 = "0"
	carpool.IsActive = "1"
	carpool.CreatedAt = time.Now()
	carpool.CreatedBy = user.EmpID
	carpool.UpdatedAt = time.Now()
	carpool.UpdatedBy = user.EmpID
	if err := config.DB.Create(&carpool).Error; err != nil {
		log.Println("DB Error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":         "Carpool created successfully",
		"data":            carpool,
		"mas_carpool_uid": carpool.MasCarpoolUID,
	})
}

// GetCarpool godoc
// @Summary Retrieve a specific carpool
// @Description This endpoint fetches details of a specific carpool using its unique identifier (MasCarpoolUID).
// @Tags Carpool-management
// @Accept json
// @Produce json
// @Param mas_carpool_uid path string true "MasCarpoolUID (mas_carpool_uid)"
// @Router /api/carpool-management/carpool/{mas_carpool_uid} [get]
func (h *CarpoolManagementHandler) GetCarpool(c *gin.Context) {
	masCarpoolUID := c.Param("mas_carpool_uid")
	var carpool models.VmsMasCarpoolResponse

	if err := config.DB.Where("mas_carpool_uid = ? AND is_deleted = ?", masCarpoolUID, "0").
		Preload("CarpoolChooseDriver").
		Preload("CarpoolChooseCar").
		First(&carpool).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Carpool not found"})
		return
	}

	c.JSON(http.StatusOK, carpool)
}

// UpdateCarpool godoc
// @Summary Update an existing carpool
// @Description Update an existing carpool's details in the database
// @Tags Carpool-management
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Param mas_carpool_uid path string true "MasCarpoolUID (mas_carpool_uid)"
// @Param carpool body models.VmsMasCarpoolRequest true "VmsMasCarpoolRequest data"
// @Router /api/carpool-management/update/{mas_carpool_uid} [put]
func (h *CarpoolManagementHandler) UpdateCarpool(c *gin.Context) {
	user := funcs.GetAuthenUser(c, h.Role)
	masCarpoolUID := c.Param("mas_carpool_uid")
	var request models.VmsMasCarpoolRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var existingCarpool models.VmsMasCarpoolRequest
	if err := config.DB.Where("mas_carpool_uid = ? AND is_deleted = ?", masCarpoolUID, "0").First(&existingCarpool).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Carpool not found"})
		return
	}
	request.MasCarpoolUID = existingCarpool.MasCarpoolUID
	request.UpdatedAt = time.Now()
	request.UpdatedBy = user.EmpID

	request.IsHaveDriverForCarpool = existingCarpool.IsHaveDriverForCarpool
	request.IsMustPassStatus30 = existingCarpool.IsMustPassStatus30
	request.IsMustPassStatus40 = existingCarpool.IsMustPassStatus40
	request.IsMustPassStatus50 = existingCarpool.IsMustPassStatus50
	request.IsActive = existingCarpool.IsActive
	request.CreatedAt = existingCarpool.CreatedAt
	request.CreatedBy = existingCarpool.CreatedBy

	if err := config.DB.Save(&request).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to update: %v", err)})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Carpool updated successfully", "data": request})
}

// DeleteCarpool godoc
// @Summary Delete a carpool
// @Description This endpoint deletes a carpool using its unique identifier (MasCarpoolUID).
// @Tags Carpool-management
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Param mas_carpool_uid path string true "MasCarpoolUID (mas_carpool_uid)"
// @Router /api/carpool-management/delete/{mas_carpool_uid} [delete]
func (h *CarpoolManagementHandler) DeleteCarpool(c *gin.Context) {
	user := funcs.GetAuthenUser(c, h.Role)
	masCarpoolUID := c.Param("mas_carpool_uid")

	var carpool models.VmsMasCarpoolRequest
	if err := config.DB.Where("mas_carpool_uid = ? AND is_deleted = ?", masCarpoolUID, "0").First(&carpool).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Carpool not found"})
		return
	}

	if err := config.DB.Model(&carpool).UpdateColumns(map[string]interface{}{
		"is_deleted": "1",
		"updated_by": user.EmpID,
		"updated_at": time.Now(),
	}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete carpool"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Carpool deleted successfully"})
}

// SearchCarpoolAdmin godoc
// @Summary Search admin carpools
// @Description Search admin carpools with pagination and filters
// @Tags Carpool-management
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Param mas_carpool_uid path string true "MasCarpoolUID (mas_carpool_uid)"
// @Param search query string false "Search query for admin_emp_no or admin_name"
// @Param is_active query string false "Filter by is_active status (comma-separated, e.g., '1,0')"
// @Param order_by query string false "Order by fields: admin_emp_no, admin_name"
// @Param order_dir query string false "Order direction: asc or desc"
// @Param page query int false "Page number (default: 1)"
// @Param limit query int false "Number of records per page (default: 10)"
// @Router /api/carpool-management/admin-search/{mas_carpool_uid} [get]
func (h *CarpoolManagementHandler) SearchCarpoolAdmin(c *gin.Context) {
	masCarpoolUID := c.Param("mas_carpool_uid")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))    // Default: page 1
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10")) // Default: 10 items per page
	offset := (page - 1) * limit

	var admins []models.VmsMasCarpoolAdminList
	query := config.DB.Table("vms_mas_carpool_admin cpa").
		Select("cpa.*,dept.dept_short admin_dept_sap_short,emp.full_name admin_emp_name").
		Joins("LEFT JOIN vms_mas_department dept ON dept.dept_sap = cpa.admin_dept_sap").
		Joins("LEFT JOIN vms_user.mas_employee emp ON emp.emp_id = cpa.admin_emp_no").
		Where("mas_carpool_uid = ? AND cpa.is_deleted = ?", masCarpoolUID, "0")

	search := strings.ToUpper(c.Query("search"))
	if search != "" {
		query = query.Where("UPPER(admin_emp_no) LIKE ? OR UPPER(admin_name) LIKE ?", "%"+search+"%", "%"+search+"%")
	}

	if isActive := c.Query("is_active"); isActive != "" {
		isActiveList := strings.Split(isActive, ",")
		query = query.Where("is_active IN (?)", isActiveList)
	}

	orderBy := c.Query("order_by")
	orderDir := c.Query("order_dir")
	if orderDir != "desc" {
		orderDir = "asc"
	}
	switch orderBy {
	case "admin_emp_no":
		query = query.Order("admin_emp_no " + orderDir)
	case "admin_name":
		query = query.Order("admin_name " + orderDir)
	default:
		query = query.Order("admin_emp_no " + orderDir) // Default ordering by admin_emp_no
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	query = query.Limit(limit).Offset(offset)
	if err := query.Find(&admins).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if len(admins) == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "No admin carpools found",
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
			"admins": admins,
		})
	}
}

// GetCarpoolAdmin godoc
// @Summary Retrieve a specific admin carpool
// @Description This endpoint fetches details of a specific admin carpool using its unique identifier (MasCarpoolAdminUID).
// @Tags Carpool-management
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Param mas_carpool_admin_uid path string true "MasCarpoolAdminUID (mas_carpool_admin_uid)"
// @Router /api/carpool-management/admin-detail/{mas_carpool_admin_uid} [get]
func (h *CarpoolManagementHandler) GetCarpoolAdmin(c *gin.Context) {
	masCarpoolAdminUID := c.Param("mas_carpool_admin_uid")

	var admin models.VmsMasCarpoolAdminList
	query := config.DB.Table("vms_mas_carpool_admin cpa").
		Select("cpa.*, dept.dept_short admin_dept_sap_short, emp.full_name admin_emp_name").
		Joins("LEFT JOIN vms_mas_department dept ON dept.dept_sap = cpa.admin_dept_sap").
		Joins("LEFT JOIN vms_user.mas_employee emp ON emp.emp_id = cpa.admin_emp_no").
		Where("cpa.mas_carpool_admin_uid = ? AND cpa.is_deleted = ?", masCarpoolAdminUID, "0")

	if err := query.First(&admin).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Admin carpool not found"})
		return
	}

	c.JSON(http.StatusOK, admin)
}

// CreateCarpoolAdmin godoc
// @Summary Create a new admin carpool
// @Description Create a new admin carpool and save it to the database
// @Tags Carpool-management
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Param carpool body []models.VmsMasCarpoolAdmin true "VmsMasCarpoolAdmin array"
// @Router /api/carpool-management/admin-create [post]
func (h *CarpoolManagementHandler) CreateCarpoolAdmin(c *gin.Context) {
	user := funcs.GetAuthenUser(c, h.Role)

	var requests []models.VmsMasCarpoolAdmin
	if err := c.ShouldBindJSON(&requests); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	for i := range requests {
		var existingAdmin models.VmsMasCarpoolAdmin
		if err := config.DB.Where("mas_carpool_uid = ? AND admin_emp_no = ? AND is_deleted = ?", requests[i].MasCarpoolUID, requests[i].AdminEmpNo, "0").First(&existingAdmin).Error; err == nil {
			c.JSON(http.StatusConflict, gin.H{
				"error": fmt.Sprintf("Admin with MasCarpoolUID %s and AdminEmpNo %s already exists", requests[i].MasCarpoolUID, requests[i].AdminEmpNo),
			})
			return
		}

		requests[i].MasCarpoolAdminUID = uuid.New().String()
		requests[i].CreatedAt = time.Now()
		requests[i].CreatedBy = user.EmpID
		requests[i].UpdatedAt = time.Now()
		requests[i].UpdatedBy = user.EmpID
		requests[i].IsDeleted = "0"
		requests[i].IsActive = "1"
		requests[i].IsMainAdmin = "0"

		empUser := funcs.GetUserEmpInfo(requests[i].AdminEmpNo)
		requests[i].AdminDeptSap = empUser.DeptSAP
	}

	if err := config.DB.Create(&requests).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Admin carpools created successfully",
		"data":    requests,
	})
}

// UpdateCarpoolAdmin godoc
// @Summary Update an admin carpool
// @Description Update an admin carpool's details in the database
// @Tags Carpool-management
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Param mas_carpool_admin_uid path string true "MasCarpoolAdminUID (mas_carpool_admin_uid)"
// @Param admin body models.VmsMasCarpoolAdmin true "VmsMasCarpoolAdmin data"
// @Router /api/carpool-management/admin-update/{mas_carpool_admin_uid} [put]
func (h *CarpoolManagementHandler) UpdateCarpoolAdmin(c *gin.Context) {
	user := funcs.GetAuthenUser(c, h.Role)
	masCarpoolAdminUID := c.Param("mas_carpool_admin_uid")

	var request models.VmsMasCarpoolAdmin
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var existingAdmin models.VmsMasCarpoolAdmin
	if err := config.DB.Where("mas_carpool_admin_uid = ? AND is_deleted = ?", masCarpoolAdminUID, "0").First(&existingAdmin).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Admin carpool not found"})
		return
	}

	request.MasCarpoolAdminUID = existingAdmin.MasCarpoolAdminUID
	request.CreatedAt = existingAdmin.CreatedAt
	request.CreatedBy = existingAdmin.CreatedBy
	request.UpdatedAt = time.Now()
	request.UpdatedBy = user.EmpID
	request.IsDeleted = existingAdmin.IsDeleted
	request.IsMainAdmin = existingAdmin.IsMainAdmin
	request.IsActive = existingAdmin.IsActive
	empUser := funcs.GetUserEmpInfo(request.AdminEmpNo)
	request.AdminDeptSap = empUser.DeptSAP

	if err := config.DB.Save(&request).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to update: %v", err)})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Admin carpool updated successfully", "data": request})
}

// DeleteCarpoolAdmin godoc
// @Summary Delete an admin carpool
// @Description This endpoint deletes an admin carpool using its unique identifier (MasCarpoolAdminUID).
// @Tags Carpool-management
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Param mas_carpool_admin_uid path string true "MasCarpoolAdminUID (mas_carpool_admin_uid)"
// @Router /api/carpool-management/admin-delete/{mas_carpool_admin_uid} [delete]
func (h *CarpoolManagementHandler) DeleteCarpoolAdmin(c *gin.Context) {
	user := funcs.GetAuthenUser(c, h.Role)
	masCarpoolAdminUID := c.Param("mas_carpool_admin_uid")

	var adminCarpool models.VmsMasCarpoolAdmin
	if err := config.DB.Where("mas_carpool_admin_uid = ? AND is_deleted = ?", masCarpoolAdminUID, "0").First(&adminCarpool).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Admin carpool not found"})
		return
	}

	if err := config.DB.Model(&adminCarpool).UpdateColumns(map[string]interface{}{
		"is_deleted": "1",
		"updated_by": user.EmpID,
		"updated_at": time.Now(),
	}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete admin carpool"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Admin carpool deleted successfully"})
}

// SearchCarpoolApprover godoc
// @Summary Search carpool approvers
// @Description Search carpool approvers with pagination and filters
// @Tags Carpool-management
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Param mas_carpool_uid path string true "MasCarpoolUID (mas_carpool_uid)"
// @Param search query string false "Search query for approver_emp_no or approver_name"
// @Param is_active query string false "Filter by is_active status (comma-separated, e.g., '1,0')"
// @Param order_by query string false "Order by fields: approver_emp_no, approver_name"
// @Param order_dir query string false "Order direction: asc or desc"
// @Param page query int false "Page number (default: 1)"
// @Param limit query int false "Number of records per page (default: 10)"
// @Router /api/carpool-management/approver-search/{mas_carpool_uid} [get]
func (h *CarpoolManagementHandler) SearchCarpoolApprover(c *gin.Context) {
	masCarpoolUID := c.Param("mas_carpool_uid")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))    // Default: page 1
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10")) // Default: 10 items per page
	offset := (page - 1) * limit

	var approvers []models.VmsMasCarpoolApproverList
	query := config.DB.Table("vms_mas_carpool_approver cpa").
		Select("cpa.*, dept.dept_short approver_dept_sap_short, emp.full_name approver_emp_name").
		Joins("LEFT JOIN vms_mas_department dept ON dept.dept_sap = cpa.approver_dept_sap").
		Joins("LEFT JOIN vms_user.mas_employee emp ON emp.emp_id = cpa.approver_emp_no").
		Where("mas_carpool_uid = ? AND cpa.is_deleted = ?", masCarpoolUID, "0")

	search := strings.ToUpper(c.Query("search"))
	if search != "" {
		query = query.Where("UPPER(approver_emp_no) LIKE ? OR UPPER(approver_name) LIKE ?", "%"+search+"%", "%"+search+"%")
	}

	if isActive := c.Query("is_active"); isActive != "" {
		isActiveList := strings.Split(isActive, ",")
		query = query.Where("is_active IN (?)", isActiveList)
	}

	orderBy := c.Query("order_by")
	orderDir := c.Query("order_dir")
	if orderDir != "desc" {
		orderDir = "asc"
	}
	switch orderBy {
	case "approver_emp_no":
		query = query.Order("approver_emp_no " + orderDir)
	case "approver_name":
		query = query.Order("approver_name " + orderDir)
	default:
		query = query.Order("approver_emp_no " + orderDir) // Default ordering by approver_emp_no
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	query = query.Limit(limit).Offset(offset)
	if err := query.Find(&approvers).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if len(approvers) == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "No approver carpools found",
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
			"approvers": approvers,
		})
	}
}

// GetCarpoolApprover godoc
// @Summary Retrieve a specific carpool approver
// @Description This endpoint fetches details of a specific carpool approver using its unique identifier (MasCarpoolApproverUID).
// @Tags Carpool-management
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Param mas_carpool_approver_uid path string true "MasCarpoolApproverUID (mas_carpool_approver_uid)"
// @Router /api/carpool-management/approver-detail/{mas_carpool_approver_uid} [get]
func (h *CarpoolManagementHandler) GetCarpoolApprover(c *gin.Context) {
	masCarpoolApproverUID := c.Param("mas_carpool_approver_uid")

	var approver models.VmsMasCarpoolApproverList
	query := config.DB.Table("vms_mas_carpool_approver cpa").
		Select("cpa.*, dept.dept_short approver_dept_sap_short, emp.full_name approver_emp_name").
		Joins("LEFT JOIN vms_mas_department dept ON dept.dept_sap = cpa.approver_dept_sap").
		Joins("LEFT JOIN vms_user.mas_employee emp ON emp.emp_id = cpa.approver_emp_no").
		Where("cpa.mas_carpool_approver_uid = ? AND cpa.is_deleted = ?", masCarpoolApproverUID, "0")

	if err := query.First(&approver).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Carpool approver not found"})
		return
	}

	c.JSON(http.StatusOK, approver)
}

// CreateCarpoolApprover godoc
// @Summary Create a new carpool approver
// @Description Create a new carpool approver and save it to the database
// @Tags Carpool-management
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Param approver body []models.VmsMasCarpoolApprover true "VmsMasCarpoolApprover array"
// @Router /api/carpool-management/approver-create [post]
func (h *CarpoolManagementHandler) CreateCarpoolApprover(c *gin.Context) {
	user := funcs.GetAuthenUser(c, h.Role)

	var requests []models.VmsMasCarpoolApprover
	if err := c.ShouldBindJSON(&requests); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	for i := range requests {
		var existingApprover models.VmsMasCarpoolApprover
		if err := config.DB.Where("mas_carpool_uid = ? AND approver_emp_no = ? AND is_deleted = ?", requests[i].MasCarpoolUID, requests[i].ApproverEmpNo, "0").First(&existingApprover).Error; err == nil {
			c.JSON(http.StatusConflict, gin.H{
				"error": fmt.Sprintf("Approver with MasCarpoolUID %s and ApproverEmpNo %s already exists", requests[i].MasCarpoolUID, requests[i].ApproverEmpNo),
			})
			return
		}

		requests[i].MasCarpoolApproverUID = uuid.New().String()
		requests[i].CreatedAt = time.Now()
		requests[i].CreatedBy = user.EmpID
		requests[i].UpdatedAt = time.Now()
		requests[i].UpdatedBy = user.EmpID
		requests[i].IsDeleted = "0"
		requests[i].IsActive = "1"
		requests[i].IsMainApprover = "0"

		empUser := funcs.GetUserEmpInfo(requests[i].ApproverEmpNo)
		requests[i].ApproverDeptSap = empUser.DeptSAP
	}

	if err := config.DB.Create(&requests).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Approver carpools created successfully",
		"data":    requests,
	})
}

// UpdateCarpoolApprover godoc
// @Summary Update a carpool approver
// @Description Update a carpool approver's details in the database
// @Tags Carpool-management
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Param mas_carpool_approver_uid path string true "MasCarpoolApproverUID (mas_carpool_approver_uid)"
// @Param approver body models.VmsMasCarpoolApprover true "VmsMasCarpoolApprover data"
// @Router /api/carpool-management/approver-update/{mas_carpool_approver_uid} [put]
func (h *CarpoolManagementHandler) UpdateCarpoolApprover(c *gin.Context) {
	user := funcs.GetAuthenUser(c, h.Role)
	masCarpoolApproverUID := c.Param("mas_carpool_approver_uid")

	var request models.VmsMasCarpoolApprover
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var existingApprover models.VmsMasCarpoolApprover
	if err := config.DB.Where("mas_carpool_approver_uid = ? AND is_deleted = ?", masCarpoolApproverUID, "0").First(&existingApprover).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Carpool approver not found"})
		return
	}

	request.MasCarpoolApproverUID = existingApprover.MasCarpoolApproverUID
	request.CreatedAt = existingApprover.CreatedAt
	request.CreatedBy = existingApprover.CreatedBy
	request.UpdatedAt = time.Now()
	request.UpdatedBy = user.EmpID
	request.IsDeleted = existingApprover.IsDeleted
	request.IsMainApprover = existingApprover.IsMainApprover
	request.IsActive = existingApprover.IsActive
	empUser := funcs.GetUserEmpInfo(request.ApproverEmpNo)
	request.ApproverDeptSap = empUser.DeptSAP

	if err := config.DB.Save(&request).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to update: %v", err)})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Carpool approver updated successfully", "data": request})
}

// DeleteCarpoolApprover godoc
// @Summary Delete a carpool approver
// @Description This endpoint deletes a carpool approver using its unique identifier (MasCarpoolApproverUID).
// @Tags Carpool-management
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Param mas_carpool_approver_uid path string true "MasCarpoolApproverUID (mas_carpool_approver_uid)"
// @Router /api/carpool-management/approver-delete/{mas_carpool_approver_uid} [delete]
func (h *CarpoolManagementHandler) DeleteCarpoolApprover(c *gin.Context) {
	user := funcs.GetAuthenUser(c, h.Role)
	masCarpoolApproverUID := c.Param("mas_carpool_approver_uid")

	var approver models.VmsMasCarpoolApprover
	if err := config.DB.Where("mas_carpool_approver_uid = ? AND is_deleted = ?", masCarpoolApproverUID, "0").First(&approver).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Carpool approver not found"})
		return
	}

	if err := config.DB.Model(&approver).UpdateColumns(map[string]interface{}{
		"is_deleted": "1",
		"updated_by": user.EmpID,
		"updated_at": time.Now(),
	}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete carpool approver"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Carpool approver deleted successfully"})
}

// SearchCarpoolVehicle godoc
// @Summary Search carpool vehicles
// @Description Search carpool vehicles with pagination and filters
// @Tags Carpool-management
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Param mas_carpool_uid path string true "MasCarpoolUID (mas_carpool_uid)"
// @Param search query string false "Search query for vehicle_no or vehicle_owner_dept_short"
// @Param is_active query string false "Filter by is_active status (comma-separated, e.g., '1,0')"
// @Param order_by query string false "Order by fields: vehicle_license_plate"
// @Param order_dir query string false "Order direction: asc or desc"
// @Param page query int false "Page number (default: 1)"
// @Param limit query int false "Number of records per page (default: 10)"
// @Router /api/carpool-management/vehicle-search/{mas_carpool_uid} [get]
func (h *CarpoolManagementHandler) SearchCarpoolVehicle(c *gin.Context) {
	masCarpoolUID := c.Param("mas_carpool_uid")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))    // Default: page 1
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10")) // Default: 10 items per page
	offset := (page - 1) * limit

	var vehicles []models.VmsMasCarpoolVehicleList
	query := config.DB.Table("vms_mas_carpool_vehicle cpv").
		Select(
			`cpv.mas_carpool_vehicle_uid,cpv.mas_carpool_uid,cpv.mas_vehicle_uid,v.vehicle_license_plate,v.vehicle_brand_name,v.vehicle_model_name,v.ref_vehicle_type_code,
				(select max(ref_vehicle_type_name) from vms_ref_vehicle_type s where s.ref_vehicle_type_code=v.ref_vehicle_type_code) ref_vehicle_type_name,
				(select max(s.dept_short) from vms_mas_department s where s.dept_sap=d.vehicle_owner_dept_sap) vehicle_owner_dept_short,
				v.ref_vehicle_type_code,d.fleet_card_no,'1' is_tax_credit,d.vehicle_mileage,
				d.vehicle_get_date,d.ref_vehicle_status_code
		`).
		Joins("LEFT JOIN vms_mas_vehicle v ON v.mas_vehicle_uid = cpv.mas_vehicle_uid").
		Joins("INNER JOIN public.vms_mas_vehicle_department AS d ON v.mas_vehicle_uid = d.mas_vehicle_uid").
		Where("cpv.mas_carpool_uid = ? AND cpv.is_deleted = ?", masCarpoolUID, "0")

	search := strings.ToUpper(c.Query("search"))
	if search != "" {
		query = query.Where("UPPER(v.vehicle_no) LIKE ? OR UPPER(v.vehicle_name) LIKE ?", "%"+search+"%", "%"+search+"%")
	}

	if isActive := c.Query("is_active"); isActive != "" {
		isActiveList := strings.Split(isActive, ",")
		query = query.Where("cpv.is_active IN (?)", isActiveList)
	}

	orderBy := c.Query("order_by")
	orderDir := c.Query("order_dir")
	if orderDir != "desc" {
		orderDir = "asc"
	}
	switch orderBy {
	case "vehicle_no":
		query = query.Order("v.vehicle_license_plate " + orderDir)
	case "vehicle_owner_dept_short":
		query = query.Order("v.vehicle_owner_dept_short " + orderDir)
	default:
		query = query.Order("v.vehicle_license_plate " + orderDir) // Default ordering by vehicle_no
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	query = query.Limit(limit).Offset(offset)
	if err := query.Find(&vehicles).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	for i := range vehicles {
		vehicles[i].Age = funcs.CalculateAge(vehicles[i].VehicleGetDate)
		funcs.TrimStringFields(&vehicles[i])
	}
	if len(vehicles) == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "No carpool vehicles found",
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
			"vehicles": vehicles,
		})
	}
}

// CreateCarpoolVehicle godoc
// @Summary Create a new carpool vehicle
// @Description Create a new carpool vehicle and save it to the database
// @Tags Carpool-management
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Param vehicle body []models.VmsMasCarpoolVehicle true "VmsMasCarpoolVehicle array"
// @Router /api/carpool-management/vehicle-create [post]
func (h *CarpoolManagementHandler) CreateCarpoolVehicle(c *gin.Context) {
	user := funcs.GetAuthenUser(c, h.Role)

	var requests []models.VmsMasCarpoolVehicle
	if err := c.ShouldBindJSON(&requests); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	for i := range requests {
		var existingVehicle models.VmsMasCarpoolVehicle
		if err := config.DB.Where("mas_carpool_uid = ? AND mas_vehicle_uid = ? AND is_deleted = ?", requests[i].MasCarpoolUID, requests[i].MasVehicleUID, "0").First(&existingVehicle).Error; err == nil {
			c.JSON(http.StatusConflict, gin.H{
				"error": fmt.Sprintf("Vehicle with MasCarpoolUID %s and MasVehicleUID %s already exists", requests[i].MasCarpoolUID, requests[i].MasVehicleUID),
			})
			return
		}

		requests[i].MasCarpoolVehicleUID = uuid.New().String()
		requests[i].CreatedAt = time.Now()
		requests[i].CreatedBy = user.EmpID
		requests[i].UpdatedAt = time.Now()
		requests[i].UpdatedBy = user.EmpID
		requests[i].IsDeleted = "0"
		requests[i].IsActive = "1"
	}

	if err := config.DB.Create(&requests).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Carpool vehicles created successfully",
		"data":    requests,
	})
}

// UpdateCarpoolVehicle godoc
// @Summary Update a carpool vehicle
// @Description Update a carpool vehicle's details in the database
// @Tags Carpool-management
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Param mas_carpool_vehicle_uid path string true "MasCarpoolVehicleUID (mas_carpool_vehicle_uid)"
// @Param vehicle body models.VmsMasCarpoolVehicle true "VmsMasCarpoolVehicle data"
// @Router /api/carpool-management/vehicle-update/{mas_carpool_vehicle_uid} [put]
func (h *CarpoolManagementHandler) UpdateCarpoolVehicle(c *gin.Context) {
	user := funcs.GetAuthenUser(c, h.Role)
	masCarpoolVehicleUID := c.Param("mas_carpool_vehicle_uid")

	var request models.VmsMasCarpoolVehicle
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var existingVehicle models.VmsMasCarpoolVehicle
	if err := config.DB.Where("mas_carpool_vehicle_uid = ? AND is_deleted = ?", masCarpoolVehicleUID, "0").First(&existingVehicle).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Carpool vehicle not found"})
		return
	}

	request.MasCarpoolVehicleUID = existingVehicle.MasCarpoolVehicleUID
	request.CreatedAt = existingVehicle.CreatedAt
	request.CreatedBy = existingVehicle.CreatedBy
	request.UpdatedAt = time.Now()
	request.UpdatedBy = user.EmpID
	request.IsDeleted = existingVehicle.IsDeleted
	request.IsActive = existingVehicle.IsActive

	if err := config.DB.Save(&request).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to update: %v", err)})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Carpool vehicle updated successfully", "data": request})
}

// DeleteCarpoolVehicle godoc
// @Summary Delete a carpool vehicle
// @Description This endpoint deletes a carpool vehicle using its unique identifier (MasCarpoolVehicleUID).
// @Tags Carpool-management
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Param mas_carpool_vehicle_uid path string true "MasCarpoolVehicleUID (mas_carpool_vehicle_uid)"
// @Router /api/carpool-management/vehicle-delete/{mas_carpool_vehicle_uid} [delete]
func (h *CarpoolManagementHandler) DeleteCarpoolVehicle(c *gin.Context) {
	user := funcs.GetAuthenUser(c, h.Role)
	masCarpoolVehicleUID := c.Param("mas_carpool_vehicle_uid")

	var vehicle models.VmsMasCarpoolVehicle
	if err := config.DB.Where("mas_carpool_vehicle_uid = ? AND is_deleted = ?", masCarpoolVehicleUID, "0").First(&vehicle).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Carpool vehicle not found"})
		return
	}

	if err := config.DB.Model(&vehicle).UpdateColumns(map[string]interface{}{
		"is_deleted": "1",
		"updated_by": user.EmpID,
		"updated_at": time.Now(),
	}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete carpool vehicle"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Carpool vehicle deleted successfully"})
}

// GetCarpoolVehicle godoc
// @Summary Retrieve a specific carpool vehicle
// @Description This endpoint fetches details of a specific carpool vehicle using its unique identifier (MasCarpoolVehicleUID).
// @Tags Carpool-management
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Param mas_carpool_vehicle_uid path string true "MasCarpoolVehicleUID (mas_carpool_vehicle_uid)"
// @Router /api/carpool-management/vehicle-detail/{mas_carpool_vehicle_uid} [get]
func (h *CarpoolManagementHandler) GetCarpoolVehicle(c *gin.Context) {
	masCarpoolVehicleUID := c.Param("mas_carpool_vehicle_uid")

	var vehicle models.VmsMasCarpoolVehicleList
	query := config.DB.Table("vms_mas_carpool_vehicle cpv").
		Select(
			`cpv.mas_carpool_vehicle_uid,cpv.mas_carpool_uid,cpv.mas_vehicle_uid,v.vehicle_license_plate,v.vehicle_brand_name,v.vehicle_model_name,v.ref_vehicle_type_code,
				(select max(ref_vehicle_type_name) from vms_ref_vehicle_type s where s.ref_vehicle_type_code=v.ref_vehicle_type_code) ref_vehicle_type_name,
				(select max(s.dept_short) from vms_mas_department s where s.dept_sap=d.vehicle_owner_dept_sap) vehicle_owner_dept_short,
				v.ref_vehicle_type_code,d.fleet_card_no,'1' is_tax_credit,d.vehicle_mileage,
				d.vehicle_get_date,d.ref_vehicle_status_code
		`).
		Joins("LEFT JOIN vms_mas_vehicle v ON v.mas_vehicle_uid = cpv.mas_vehicle_uid").
		Joins("INNER JOIN public.vms_mas_vehicle_department AS d ON v.mas_vehicle_uid = d.mas_vehicle_uid").
		Where("cpv.mas_carpool_vehicle_uid = ? AND cpv.is_deleted = ?", masCarpoolVehicleUID, "0")

	if err := query.First(&vehicle).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Carpool vehicle not found"})
		return
	}

	vehicle.Age = funcs.CalculateAge(vehicle.VehicleGetDate)
	funcs.TrimStringFields(&vehicle)

	c.JSON(http.StatusOK, vehicle)
}

// SetActiveCarpool godoc
// @Summary Set active status for a carpool
// @Description Update the active status of a carpool using its unique identifier (MasCarpoolUID).
// @Tags Carpool-management
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Param active body models.VmsMasCarpoolActive true "VmsMasCarpoolActive data"
// @Router /api/carpool-management/set-active [put]
func (h *CarpoolManagementHandler) SetActiveCarpool(c *gin.Context) {
	user := funcs.GetAuthenUser(c, h.Role)

	var request models.VmsMasCarpoolActive
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var carpool models.VmsMasCarpoolRequest
	if err := config.DB.Where("mas_carpool_uid = ? AND is_deleted = ?", request.MasCarpoolUID, "0").First(&carpool).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Carpool not found"})
		return
	}

	carpool.IsActive = request.IsActive
	carpool.UpdatedAt = time.Now()
	carpool.UpdatedBy = user.EmpID

	if err := config.DB.Save(&carpool).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to update active status: %v", err)})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Carpool active status updated successfully", "data": request})
}
