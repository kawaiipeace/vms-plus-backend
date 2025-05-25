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
	"vms_plus_be/messages"
	"vms_plus_be/models"
	"vms_plus_be/userhub"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/tealeg/xlsx"
	"gorm.io/gorm"
)

type CarpoolManagementHandler struct {
	Role string
}

func (h *CarpoolManagementHandler) SetQueryRole(user *models.AuthenUserEmp, query *gorm.DB) *gorm.DB {
	if user.EmpID == "" {
		return query
	}
	return query
}

func (h *CarpoolManagementHandler) SetQueryRoleDept(user *models.AuthenUserEmp, query *gorm.DB) *gorm.DB {
	if user.EmpID == "" {
		return query
	}
	return query
}

func GetCarpoolTypeName(carpoolType string) string {
	switch carpoolType {
	case "01":
		return "Car Pool สำนักงานใหญ่"
	case "02":
		return "Car Pool เขต"
	case "03":
		return "Car Pool หน่วยงาน"
	default:
		return ""
	}
}

func GetCarpoolName(MasCarpoolUID string) string {
	var carpool models.VmsMasCarpoolList
	if err := config.DB.Where("mas_carpool_uid = ?", MasCarpoolUID).First(&carpool).Error; err != nil {
		log.Printf("Error fetching carpool name: %v", err)
		return ""
	}
	return carpool.CarpoolName
}

func DoSearchCarpools(c *gin.Context, h *CarpoolManagementHandler, isLimit bool) ([]models.VmsMasCarpoolList, int64, error) {
	user := funcs.GetAuthenUser(c, h.Role)
	if c.IsAborted() {
		return nil, 0, nil
	}
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))    // Default: page 1
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10")) // Default: 10 items per page
	offset := (page - 1) * limit
	if !isLimit {
		limit = 1000000
	}
	var carpools []models.VmsMasCarpoolList
	query := h.SetQueryRole(user, config.DB)
	query = query.Model(&models.VmsMasCarpoolList{})
	query = query.Table("vms_mas_carpool cp").Select(`cp.*,
		(select count(*) from vms_mas_carpool_driver cpd where is_deleted='0' and cpd.mas_carpool_uid=cp.mas_carpool_uid) number_of_drivers,
		(select count(*) from vms_mas_carpool_vehicle cpv where is_deleted='0' and cpv.mas_carpool_uid=cp.mas_carpool_uid) number_of_vehicles,
		(select count(*) from vms_mas_carpool_approver cpa where is_deleted='0' and cpa.mas_carpool_uid=cp.mas_carpool_uid) number_of_approvers,
		(select admin_emp_name from vms_mas_carpool_admin cpa where is_deleted='0' and is_main_admin='1' and cpa.mas_carpool_uid=cp.mas_carpool_uid) carpool_admin_emp_name,
		(select dept_short from vms_mas_carpool_admin cpa,vms_mas_department md where cpa.is_deleted='0' and is_main_admin='1' and md.dept_sap=cpa.admin_dept_sap and cpa.mas_carpool_uid=cp.mas_carpool_uid) carpool_admin_dept_sap,
		case carpool_type when '01' then 'สำนักงานใหญ่'
			when '02' then (select dept_short from vms_mas_department md where md.dept_sap=cp.carpool_dept_sap)
			when '03' then 
				case when (select count(*) from vms_mas_carpool_authorized_dept cad where cad.mas_carpool_uid=cp.mas_carpool_uid and cad.is_deleted='0') > 1 
					then cast((select count(*) from vms_mas_carpool_authorized_dept cad where cad.mas_carpool_uid=cp.mas_carpool_uid and cad.is_deleted='0') as text)||' หน่วยงาน'
					else (select md.dept_short from vms_mas_department md,vms_mas_carpool_authorized_dept cad where md.dept_sap=cad.dept_sap and cad.mas_carpool_uid=cp.mas_carpool_uid and cad.is_deleted='0')
				end
		end as carpool_authorized_depts

	`)
	search := strings.ToUpper(c.Query("search"))
	if search != "" {
		query = query.Where("UPPER(cp.carpool_name) ILIKE ? OR EXISTS (SELECT 1 FROM vms_mas_carpool_admin cpa WHERE cpa.mas_carpool_uid = cp.mas_carpool_uid AND UPPER(cpa.admin_emp_name) ILIKE ?)", "%"+search+"%", "%"+search+"%")
	}
	if deptSap := c.Query("dept_sap"); deptSap != "" {
		deptSapList := strings.Split(deptSap, ",")
		query = query.Where("EXISTS (SELECT 1 FROM vms_mas_carpool_authorized_dept cad WHERE cad.mas_carpool_uid = cp.mas_carpool_uid AND cad.dept_sap IN (?))", deptSapList)
	}
	if isActive := c.Query("is_active"); isActive != "" {
		isActiveList := strings.Split(isActive, ",")
		conditions := []string{}
		args := []interface{}{}

		for _, status := range isActiveList {
			if status == "2" {
				conditions = append(conditions, "((SELECT COUNT(*) FROM vms_mas_carpool_vehicle cpv WHERE is_deleted='0' AND cpv.mas_carpool_uid=cp.mas_carpool_uid) = 0 OR (SELECT COUNT(*) FROM vms_mas_carpool_approver cpa WHERE is_deleted='0' AND cpa.mas_carpool_uid=cp.mas_carpool_uid) = 0)")
			} else {
				conditions = append(conditions, "cp.is_active = ?")
				args = append(args, status)
			}
		}

		query = query.Where(strings.Join(conditions, " OR "), args...)
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "message": messages.ErrInternalServer.Error()})
		return nil, 0, err
	}
	query = query.Limit(limit).
		Offset(offset)

	if err := query.Find(&carpools).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "message": messages.ErrInternalServer.Error()})
		return nil, 0, err
	}

	for i := range carpools {
		carpools[i].CarpoolTypeName = GetCarpoolTypeName(carpools[i].CarpoolType)
		if carpools[i].NumberOfVehicles == 0 || carpools[i].NumberOfApprovers == 0 {
			carpools[i].CarpoolStatus = "ไม่พร้อมใช้งาน"
		} else if carpools[i].IsActive == "1" {
			carpools[i].CarpoolStatus = "เปิด"
		} else if carpools[i].IsActive == "0" {
			carpools[i].CarpoolStatus = "ปิด"
		}
	}
	return carpools, total, nil
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
// @Param is_active query string false "Filter by is_active status (comma-separated, e.g., '2,1,0') 2=ไม่พร้อมใช้งาน 1=เปิด 0=ปิด"
// @Param dept_sap query string false "Filter by dept_sap"
// @Param order_by query string false "Order by fields: carpool_name, number_of_drivers, number_of_vehicles"
// @Param order_dir query string false "Order direction: asc or desc"
// @Param page query int false "Page number (default: 1)"
// @Param limit query int false "Number of records per page (default: 10)"
// @Router /api/carpool-management/search [get]
func (h *CarpoolManagementHandler) SearchCarpools(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))    // Default: page 1
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10")) // Default: 10 items per page

	carpools, total, err := DoSearchCarpools(c, h, true)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "message": messages.ErrInternalServer.Error()})
		return
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

// ExportCarpools godoc
// @Summary Export carpool management
// @Description Export carpool management by criteria
// @Tags Carpool-management
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Param search query string false "Search query for carpool_name or emp_name"
// @Param is_active query string false "Filter by is_active status (comma-separated, e.g., '2,1,0') 2=ไม่พร้อมใช้งาน 1=เปิด 0=ปิด"
// @Param dept_sap query string false "Filter by dept_sap"
// @Param order_by query string false "Order by fields: carpool_name, number_of_drivers, number_of_vehicles"
// @Param order_dir query string false "Order direction: asc or desc"
// @Router /api/carpool-management/export [get]
func (h *CarpoolManagementHandler) ExportCarpools(c *gin.Context) {
	carpools, _, err := DoSearchCarpools(c, h, false)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "message": messages.ErrInternalServer.Error()})
		return
	}
	// Create Excel file
	file := xlsx.NewFile()
	sheet, err := file.AddSheet("Carpools")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create Excel sheet", "message": err.Error()})
		return
	}

	// Set headers
	headers := []string{
		"ชื่อ Car Pool",
		"ประเภท Car Pool",
		"หน่วยงานที่รับผิดชอบ",
		"ผู้ดูแล Car Pool",
		"สังกัดผู้ดูแล",
		"จำนวนยานพาหนะ",
		"จำนวนพนักงานขับรถ",
		"จำนวนผู้อนุมัติ",
		"สถานะ",
	}

	row := sheet.AddRow()
	for _, header := range headers {
		cell := row.AddCell()
		cell.Value = header
	}

	for _, carpool := range carpools {
		row := sheet.AddRow()
		row.AddCell().Value = carpool.CarpoolName
		row.AddCell().Value = carpool.CarpoolTypeName
		row.AddCell().Value = carpool.CarpoolAuthorizedDepts
		row.AddCell().Value = carpool.CarpoolAdminEmpName
		row.AddCell().Value = carpool.CarpoolAdminDeptSap
		row.AddCell().Value = strconv.Itoa(carpool.NumberOfVehicles) + " คัน"
		row.AddCell().Value = strconv.Itoa(carpool.NumberOfDrivers) + " คน"
		row.AddCell().Value = strconv.Itoa(carpool.NumberOfApprovers) + " คน"
		row.AddCell().Value = carpool.CarpoolStatus
	}

	c.Header("Content-Disposition", "attachment; filename=carpools.xlsx")
	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("File-Name", fmt.Sprintf("carpools_%s.xlsx", time.Now().Format("2006-01-02")))
	c.Header("Content-Transfer-Encoding", "binary")
	if err := file.Write(c.Writer); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to write Excel file", "message": err.Error()})
		return
	}
}

// GetMasDepartment godoc
// @Summary Retrieve the Pea Department
// @Description This endpoint allows a user to Pea Department.
// @Tags Carpool-management
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Param carpool_type path string true "CarpoolType (trn_request_uid)"
// @Param search query string false "Search by DeptSap Code or DetpSap Name"
// @Router /api/carpool-management/mas-department/{carpool_type} [get]
func (h *CarpoolManagementHandler) GetMasDepartment(c *gin.Context) {
	user := funcs.GetAuthenUser(c, h.Role)
	if c.IsAborted() {
		return
	}
	carpoolType := c.Param("carpool_type")
	if carpoolType == "01" {
		c.JSON(http.StatusOK, []interface{}{})
		return
	}
	var lists []models.VmsMasDepartment

	query := h.SetQueryRoleDept(user, config.DB)
	if carpoolType == "02" {
		query = query.Where("resource_name = 'การไฟฟ้าเขต'")
	}
	search := c.Query("search")
	query = query.Where("is_deleted = ? AND is_active = ?", "0", "1").Limit(100)
	if search != "" {
		query = query.Where("dept_sap ILIKE ? OR dept_short ILIKE ? OR dept_full ILIKE ?", "%"+search+"%", "%"+search+"%", "%"+search+"%")
	}
	query = query.Where("is_deleted = ? AND is_active = ?", "0", "1")
	if err := query.
		Find(&lists).Error; err != nil {
		c.JSON(http.StatusOK, []interface{}{})
		return
	}

	c.JSON(http.StatusOK, lists)
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
	if c.IsAborted() {
		return
	}
	var carpool models.VmsMasCarpoolRequest

	if err := c.ShouldBindJSON(&carpool); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "message": messages.ErrInvalidJSONInput.Error()})
		return
	}

	carpool.MasCarpoolUID = uuid.New().String()
	carpool.IsHaveDriverForCarpool = "0"

	carpool.IsActive = "1"
	carpool.CreatedAt = time.Now()
	carpool.CreatedBy = user.EmpID
	carpool.UpdatedAt = time.Now()
	carpool.UpdatedBy = user.EmpID

	for i := range carpool.CarpoolAuthorizedDepts {
		carpool.CarpoolAuthorizedDepts[i].MasCarpoolAuthorizedDeptUID = uuid.New().String()
		carpool.CarpoolAuthorizedDepts[i].MasCarpoolUID = carpool.MasCarpoolUID
		carpool.CarpoolAuthorizedDepts[i].CreatedAt = time.Now()
		carpool.CarpoolAuthorizedDepts[i].CreatedBy = user.EmpID
		carpool.CarpoolAuthorizedDepts[i].UpdatedAt = time.Now()
		carpool.CarpoolAuthorizedDepts[i].UpdatedBy = user.EmpID
		carpool.CarpoolAuthorizedDepts[i].IsDeleted = "0"
	}

	if err := config.DB.Create(&carpool).Error; err != nil {
		log.Println("DB Error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "message": messages.ErrInternalServer.Error()})
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
	user := funcs.GetAuthenUser(c, h.Role)
	if c.IsAborted() {
		return
	}
	masCarpoolUID := c.Param("mas_carpool_uid")
	var carpool models.VmsMasCarpoolResponse
	query := h.SetQueryRole(user, config.DB)
	if err := query.Where("mas_carpool_uid = ? AND is_deleted = ?", masCarpoolUID, "0").
		Preload("CarpoolChooseDriver").
		Preload("CarpoolChooseCar").
		Preload("CarpoolAuthorizedDepts.MasDepartment").
		Preload("CarpoolAdmins").
		Preload("CarpoolApprovers").
		Preload("CarpoolVehicles").
		Preload("CarpoolDrivers").
		First(&carpool).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Carpool not found", "message": messages.ErrNotfound.Error()})
		return
	}
	if len(carpool.CarpoolVehicles) > 0 && len(carpool.CarpoolApprovers) > 0 {
		carpool.IsCarpoolReady = true
	}
	if len(carpool.CarpoolDrivers) > 0 {
		carpool.IsCarpoolChooseDriver = true
	}
	carpool.CarpoolTypeName = GetCarpoolTypeName(carpool.CarpoolType)

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
	if c.IsAborted() {
		return
	}
	masCarpoolUID := c.Param("mas_carpool_uid")
	var request models.VmsMasCarpoolRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var existingCarpool models.VmsMasCarpoolRequest
	query := h.SetQueryRole(user, config.DB)
	if err := query.Where("mas_carpool_uid = ? AND is_deleted = ?", masCarpoolUID, "0").First(&existingCarpool).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Carpool not found", "message": messages.ErrNotfound.Error()})
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

	if err := config.DB.Where("mas_carpool_uid = ?", masCarpoolUID).Delete(&models.VmsMasCarpoolAuthorizedDept{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to delete authorized departments: %v", err), "message": messages.ErrInternalServer.Error()})
		return
	}

	for i := range request.CarpoolAuthorizedDepts {
		request.CarpoolAuthorizedDepts[i].MasCarpoolAuthorizedDeptUID = uuid.New().String()
		request.CarpoolAuthorizedDepts[i].MasCarpoolUID = request.MasCarpoolUID
		request.CarpoolAuthorizedDepts[i].CreatedAt = time.Now()
		request.CarpoolAuthorizedDepts[i].CreatedBy = user.EmpID
		request.CarpoolAuthorizedDepts[i].IsDeleted = "0"
	}

	if err := config.DB.Save(&request).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to update: %v", err), "message": messages.ErrInternalServer.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Carpool updated successfully", "data": request, "carpool_name": GetCarpoolName(masCarpoolUID)})
}

// DeleteCarpool godoc
// @Summary Delete a carpool
// @Description This endpoint deletes a carpool using its unique identifier (MasCarpoolUID).
// @Tags Carpool-management
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Param data body models.VmsMasCarpoolDelete true "VmsMasCarpoolDelete data"
// @Router /api/carpool-management/delete [delete]
func (h *CarpoolManagementHandler) DeleteCarpool(c *gin.Context) {
	user := funcs.GetAuthenUser(c, h.Role)
	if c.IsAborted() {
		return
	}
	var request, carpool models.VmsMasCarpoolDelete

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	queryRole := h.SetQueryRole(user, config.DB)
	if err := queryRole.Where("mas_carpool_uid = ? AND is_deleted = ? AND carpool_name = ?", request.MasCarpoolUID, "0", request.CarpoolName).First(&carpool).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Carpool not found", "message": messages.ErrNotfound.Error()})
		return
	}
	var count int64
	tables := []string{"vms_mas_carpool_admin", "vms_mas_carpool_approver", "vms_mas_carpool_vehicle", "vms_mas_carpool_driver"}
	for _, table := range tables {
		if err := config.DB.Table(table).
			Where("mas_carpool_uid = ? AND is_deleted = ?", request.MasCarpoolUID, "0").
			Count(&count).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to check dependencies in table %s: %v", table, err), "message": messages.ErrInternalServer.Error()})
			return
		}
		if count > 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   fmt.Sprintf("Cannot delete admin carpool as it has dependencies in table %s", table),
				"message": "การลบหน่วยงานที่สามารถใช้บริการกลุ่มยานพาหนะนี้ จำเป็นต้องลบยานพาหนะ, พนักงานขับรถ, ผู้ดูแลยานพาหนะ และผู้อนุมัติที่สังกัดหน่วยงานนั้น ออกจากกลุ่มก่อน",
			})
			return
		}
	}
	var requests int64
	if err := config.DB.Table("vms_trn_request").
		Where("mas_vehicle_uid IN (SELECT mas_vehicle_uid FROM vms_mas_carpool_vehicle WHERE mas_carpool_uid = ? AND is_deleted = ?)", request.MasCarpoolUID, "0").
		Count(&requests).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to check dependencies in vms_trn_request: %v", err), "message": messages.ErrInternalServer.Error()})
		return
	}
	if requests > 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Cannot delete admin carpool as it has dependencies in vms_trn_request",
			"message": "ไม่สามารถลบกลุ่มได้เนื่องจากมีคำขอใช้ยานพาหนะของกลุ่มที่ยังดำเนินการไม่เสร็จสิ้น",
		})
		return
	}
	if err := config.DB.Model(&carpool).UpdateColumns(map[string]interface{}{
		"is_deleted": "1",
		"updated_by": user.EmpID,
		"updated_at": time.Now(),
	}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete carpool", "message": messages.ErrInternalServer.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Carpool deleted successfully", "carpool_name": GetCarpoolName(request.MasCarpoolUID)})
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
	if c.IsAborted() {
		return
	}

	var request models.VmsMasCarpoolActive
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var carpool models.VmsMasCarpoolRequest
	queryRole := h.SetQueryRole(user, config.DB)
	if err := queryRole.Where("mas_carpool_uid = ? AND is_deleted = ?", request.MasCarpoolUID, "0").First(&carpool).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Carpool not found", "message": messages.ErrNotfound.Error()})
		return
	}

	carpool.IsActive = request.IsActive
	carpool.UpdatedAt = time.Now()
	carpool.UpdatedBy = user.EmpID

	if err := config.DB.Save(&carpool).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to update active status: %v", err), "message": messages.ErrInternalServer.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Carpool active status updated successfully", "data": request, "carpool_name": GetCarpoolName(request.MasCarpoolUID)})
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
	user := funcs.GetAuthenUser(c, h.Role)
	if c.IsAborted() {
		return
	}

	masCarpoolUID := c.Param("mas_carpool_uid")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))    // Default: page 1
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10")) // Default: 10 items per page
	offset := (page - 1) * limit
	var carpool models.VmsMasCarpoolList
	query := h.SetQueryRole(user, config.DB)
	if err := query.Where("mas_carpool_uid = ? AND is_deleted = ?", masCarpoolUID, "0").First(&carpool).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Carpool not found", "message": messages.ErrNotfound.Error()})
		return
	}
	var admins []models.VmsMasCarpoolAdminList
	query = config.DB.Table("vms_mas_carpool_admin cpa").
		Select("cpa.*,dept.dept_short admin_dept_sap_short").
		Joins("LEFT JOIN vms_mas_department dept ON dept.dept_sap = cpa.admin_dept_sap").
		Where("mas_carpool_uid = ? AND cpa.is_deleted = ?", masCarpoolUID, "0")

	search := strings.ToUpper(c.Query("search"))
	if search != "" {
		query = query.Where("UPPER(admin_emp_no) ILIKE ? OR UPPER(admin_name) ILIKE ?", "%"+search+"%", "%"+search+"%")
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "message": messages.ErrInternalServer.Error()})
		return
	}

	query = query.Limit(limit).Offset(offset)
	if err := query.Find(&admins).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "message": messages.ErrInternalServer.Error()})
		return
	}

	if len(admins) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"message": "No admin carpools found",
			"pagination": gin.H{
				"page":       page,
				"limit":      limit,
				"totalPages": (total + int64(limit) - 1) / int64(limit),
			},
			"admins": []models.VmsMasCarpoolAdminList{},
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
	user := funcs.GetAuthenUser(c, h.Role)
	masCarpoolAdminUID := c.Param("mas_carpool_admin_uid")

	var admin models.VmsMasCarpoolAdminList
	query := h.SetQueryRole(user, config.DB)
	query = query.Table("vms_mas_carpool_admin cpa").
		Select("cpa.*, dept.dept_short admin_dept_sap_short").
		Joins("LEFT JOIN vms_mas_department dept ON dept.dept_sap = cpa.admin_dept_sap").
		Where("cpa.mas_carpool_admin_uid = ? AND cpa.is_deleted = ?", masCarpoolAdminUID, "0")

	if err := query.First(&admin).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Admin carpool not found", "message": messages.ErrNotfound.Error()})
		return
	}
	var carpool models.VmsMasCarpoolList
	queryRole := h.SetQueryRole(user, config.DB)
	if err := queryRole.Where("mas_carpool_uid = ? AND is_deleted = ?", admin.MasCarpoolUID, "0").First(&carpool).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Carpool not found", "message": messages.ErrNotfound.Error()})
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
	if c.IsAborted() {
		return
	}

	var requests []models.VmsMasCarpoolAdmin
	if err := c.ShouldBindJSON(&requests); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "message": messages.ErrInvalidJSONInput.Error()})
		return
	}
	if len(requests) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No admin carpool data provided", "message": messages.ErrInvalidJSONInput.Error()})
		return
	}
	for i := range requests {
		var existingCarpool models.VmsMasCarpoolRequest
		queryRole := h.SetQueryRole(user, config.DB)
		if err := queryRole.Where("mas_carpool_uid = ? AND is_deleted = ?", requests[i].MasCarpoolUID, "0").First(&existingCarpool).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Carpool not found", "message": messages.ErrNotfound.Error()})
			return
		}
	}

	for i := range requests {
		var existingAdmin models.VmsMasCarpoolAdmin
		if err := config.DB.Where("mas_carpool_uid = ? AND admin_emp_no = ? AND is_deleted = ?", requests[i].MasCarpoolUID, requests[i].AdminEmpNo, "0").First(&existingAdmin).Error; err == nil {
			c.JSON(http.StatusConflict, gin.H{
				"error":   fmt.Sprintf("มีผู้ดูแลที่มี MasCarpoolUID %s และ AdminEmpNo %s อยู่แล้ว", requests[i].MasCarpoolUID, requests[i].AdminEmpNo),
				"message": "ไม่สามารถเพิ่มผู้ดูแลได้เนื่องจากมีผู้ดูแลที่มีรหัสพนักงานนี้อยู่ในกลุ่มยานพาหนะนี้แล้ว",
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
		requests[i].AdminEmpName = empUser.FullName
	}

	if err := config.DB.Create(&requests).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "message": messages.ErrInternalServer.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":      "Admin carpools created successfully",
		"data":         requests,
		"carpool_name": GetCarpoolName(requests[0].MasCarpoolUID),
	})
}

// UpdateCarpoolAdmin godoc
// @Summary Update an admin carpool to main
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
	if c.IsAborted() {
		return
	}
	masCarpoolAdminUID := c.Param("mas_carpool_admin_uid")

	var request models.VmsMasCarpoolAdmin
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "message": messages.ErrInvalidJSONInput.Error()})
		return
	}

	var carpool models.VmsMasCarpoolList
	queryRole := h.SetQueryRole(user, config.DB)
	if err := queryRole.Where("mas_carpool_uid = ? AND is_deleted = ?", masCarpoolAdminUID, "0").First(&carpool).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Carpool not found", "message": messages.ErrNotfound.Error()})
		return
	}

	var existingAdmin models.VmsMasCarpoolAdmin
	if err := config.DB.Where("mas_carpool_admin_uid = ? AND is_deleted = ?", masCarpoolAdminUID, "0").First(&existingAdmin).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Admin carpool not found", "message": messages.ErrNotfound.Error()})
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to update: %v", err), "message": messages.ErrInternalServer.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Admin carpool updated successfully", "data": request, "carpool_name": GetCarpoolName(request.MasCarpoolUID)})
}

// UpdateCarpoolMainAdmin godoc
// @Summary Update an admin carpool to main admin
// @Description Update an admin carpool to main admin
// @Tags Carpool-management
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Param mas_carpool_admin_uid path string true "MasCarpoolAdminUID (mas_carpool_admin_uid)"
// @Router /api/carpool-management/admin-update-main-admin/{mas_carpool_admin_uid} [put]
func (h *CarpoolManagementHandler) UpdateCarpoolMainAdmin(c *gin.Context) {
	user := funcs.GetAuthenUser(c, h.Role)
	if c.IsAborted() {
		return
	}

	masCarpoolAdminUID := c.Param("mas_carpool_admin_uid")

	var existingAdmin models.VmsMasCarpoolAdmin
	queryRole := h.SetQueryRole(user, config.DB)
	if err := queryRole.Where("mas_carpool_admin_uid = ? AND is_deleted = ?", masCarpoolAdminUID, "0").First(&existingAdmin).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Admin carpool not found", "message": messages.ErrNotfound.Error()})
		return
	}

	if err := config.DB.Model(&models.VmsMasCarpoolAdmin{}).
		Where("mas_carpool_uid = ?", existingAdmin.MasCarpoolUID).
		Update("is_main_admin", "0").Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to update is_main_admin for all: %v", err), "message": messages.ErrInternalServer.Error()})
		return
	}

	if err := config.DB.Model(&existingAdmin).Updates(map[string]interface{}{
		"is_main_admin": "1",
		"updated_by":    user.EmpID,
		"updated_at":    time.Now(),
	}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to update is_main_admin: %v", err), "message": messages.ErrInternalServer.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Admin carpool updated successfully", "carpool_name": GetCarpoolName(existingAdmin.MasCarpoolUID)})
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
	if c.IsAborted() {
		return
	}
	masCarpoolAdminUID := c.Param("mas_carpool_admin_uid")

	var adminCarpool models.VmsMasCarpoolAdmin
	if err := config.DB.Where("mas_carpool_admin_uid = ? AND is_deleted = ?", masCarpoolAdminUID, "0").First(&adminCarpool).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Admin carpool not found", "message": messages.ErrNotfound.Error()})
		return
	}

	var carpool models.VmsMasCarpoolList
	queryRole := h.SetQueryRole(user, config.DB)
	if err := queryRole.Where("mas_carpool_uid = ? AND is_deleted = ?", adminCarpool.MasCarpoolUID, "0").First(&carpool).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Carpool not found", "message": messages.ErrNotfound.Error()})
		return
	}

	if err := config.DB.Model(&adminCarpool).UpdateColumns(map[string]interface{}{
		"is_deleted": "1",
		"updated_by": user.EmpID,
		"updated_at": time.Now(),
	}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete admin carpool", "message": messages.ErrInternalServer.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Admin carpool deleted successfully", "carpool_name": GetCarpoolName(adminCarpool.MasCarpoolUID)})
}

// SearchMasAdminUser godoc
// @Summary Retrieve the Admin Users
// @Description This endpoint allows a user to retrieve Admin Users.
// @Tags Carpool-management
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Param search query string false "Search by Employee ID or Full Name"
// @Router /api/carpool-management/admin-mas-search [get]
func (h *CarpoolManagementHandler) SearchMasAdminUser(c *gin.Context) {
	funcs.GetAuthenUser(c, h.Role)
	if c.IsAborted() {
		return
	}
	search := c.Query("search")

	request := userhub.ServiceListUserRequest{
		ServiceCode: "vms",
		Search:      search,
		Limit:       100,
	}
	lists, err := userhub.GetUserList(request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, lists)
}

// SearchMasApprovalUser godoc
// @Summary Retrieve the Admin Users
// @Description This endpoint allows a user to retrieve Admin Users.
// @Tags Carpool-management
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Param search query string false "Search by Employee ID or Full Name"
// @Router /api/carpool-management/approver-mas-search [get]
func (h *CarpoolManagementHandler) SearchMasApprovalUser(c *gin.Context) {
	funcs.GetAuthenUser(c, h.Role)
	if c.IsAborted() {
		return
	}
	search := c.Query("search")

	request := userhub.ServiceListUserRequest{
		ServiceCode: "vms",
		Search:      search,
		Limit:       100,
	}
	lists, err := userhub.GetUserList(request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, lists)
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
	user := funcs.GetAuthenUser(c, h.Role)
	masCarpoolUID := c.Param("mas_carpool_uid")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))    // Default: page 1
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10")) // Default: 10 items per page
	offset := (page - 1) * limit

	var carpool models.VmsMasCarpoolList
	queryRole := h.SetQueryRole(user, config.DB)
	if err := queryRole.Where("mas_carpool_uid = ? AND is_deleted = ?", masCarpoolUID, "0").First(&carpool).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Carpool not found", "message": messages.ErrNotfound.Error()})
		return
	}

	var approvers []models.VmsMasCarpoolApproverList
	query := config.DB.Table("vms_mas_carpool_approver cpa").
		Select("cpa.*, dept.dept_short approver_dept_sap_short").
		Joins("LEFT JOIN vms_mas_department dept ON dept.dept_sap = cpa.approver_dept_sap").
		Where("mas_carpool_uid = ? AND cpa.is_deleted = ?", masCarpoolUID, "0")

	search := strings.ToUpper(c.Query("search"))
	if search != "" {
		query = query.Where("UPPER(approver_emp_no) ILIKE ? OR UPPER(approver_name) ILIKE ?", "%"+search+"%", "%"+search+"%")
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "message": messages.ErrInternalServer.Error()})
		return
	}

	query = query.Limit(limit).Offset(offset)
	if err := query.Find(&approvers).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "message": messages.ErrInternalServer.Error()})
		return
	}

	if len(approvers) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"message": "No approver carpools found",
			"pagination": gin.H{
				"page":       page,
				"limit":      limit,
				"totalPages": (total + int64(limit) - 1) / int64(limit),
			},
			"approvers": []models.VmsMasCarpoolApproverList{},
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
	user := funcs.GetAuthenUser(c, h.Role)
	if c.IsAborted() {
		return
	}
	masCarpoolApproverUID := c.Param("mas_carpool_approver_uid")

	var approver models.VmsMasCarpoolApproverList
	query := config.DB.Table("vms_mas_carpool_approver cpa").
		Select("cpa.*, dept.dept_short approver_dept_sap_short").
		Joins("LEFT JOIN vms_mas_department dept ON dept.dept_sap = cpa.approver_dept_sap").
		Where("cpa.mas_carpool_approver_uid = ? AND cpa.is_deleted = ?", masCarpoolApproverUID, "0")

	if err := query.First(&approver).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Carpool approver not found", "message": messages.ErrNotfound.Error()})
		return
	}

	var carpool models.VmsMasCarpoolList
	queryRole := h.SetQueryRole(user, config.DB)
	if err := queryRole.Where("mas_carpool_uid = ? AND is_deleted = ?", approver.MasCarpoolUID, "0").First(&carpool).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Carpool not found", "message": messages.ErrNotfound.Error()})
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
	if c.IsAborted() {
		return
	}

	var requests []models.VmsMasCarpoolApprover
	if err := c.ShouldBindJSON(&requests); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "message": messages.ErrInvalidJSONInput.Error()})
		return
	}
	for i := range requests {
		var existingCarpool models.VmsMasCarpoolRequest
		queryRole := h.SetQueryRole(user, config.DB)
		if err := queryRole.Where("mas_carpool_uid = ? AND is_deleted = ?", requests[i].MasCarpoolUID, "0").First(&existingCarpool).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Carpool not found", "message": messages.ErrNotfound.Error()})
			return
		}
	}
	for i := range requests {
		var existingApprover models.VmsMasCarpoolApprover
		if err := config.DB.Where("mas_carpool_uid = ? AND approver_emp_no = ? AND is_deleted = ?", requests[i].MasCarpoolUID, requests[i].ApproverEmpNo, "0").First(&existingApprover).Error; err == nil {
			c.JSON(http.StatusConflict, gin.H{
				"error":   fmt.Sprintf("Approver with MasCarpoolUID %s and ApproverEmpNo %s already exists", requests[i].MasCarpoolUID, requests[i].ApproverEmpNo),
				"Message": "ไม่สามารถเพิ่มผู้อนุมัติได้เนื่องจากมีผู้อนุมัติที่มีรหัสพนักงานนี้อยู่ในกลุ่มยานพาหนะนี้แล้ว",
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
		"message":      "Approver carpools created successfully",
		"data":         requests,
		"carpool_name": GetCarpoolName(requests[0].MasCarpoolUID),
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
	if c.IsAborted() {
		return
	}
	masCarpoolApproverUID := c.Param("mas_carpool_approver_uid")

	var request models.VmsMasCarpoolApprover
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "message": messages.ErrInvalidJSONInput.Error()})
		return
	}

	var carpool models.VmsMasCarpoolList
	queryRole := h.SetQueryRole(user, config.DB)
	if err := queryRole.Where("mas_carpool_uid = ? AND is_deleted = ?", request.MasCarpoolUID, "0").First(&carpool).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Carpool not found", "message": messages.ErrNotfound.Error()})
		return
	}

	var existingApprover models.VmsMasCarpoolApprover
	if err := config.DB.Where("mas_carpool_approver_uid = ? AND is_deleted = ?", masCarpoolApproverUID, "0").First(&existingApprover).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Carpool approver not found", "message": messages.ErrNotfound.Error()})
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
	request.ApproverEmpName = empUser.FullName

	if err := config.DB.Save(&request).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to update: %v", err), "message": messages.ErrInternalServer.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Carpool approver updated successfully", "data": request, "carpool_name": GetCarpoolName(existingApprover.MasCarpoolUID)})
}

// UpdateCarpoolMainApprover godoc
// @Summary Update an approver carpool to main approver
// @Description Update an approver carpool's to main approver
// @Tags Carpool-management
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Param mas_carpool_approver_uid path string true "UpdateCarpoolMainApprover (mas_carpool_approver_uid)"
// @Router /api/carpool-management/approver-update-main-approver/{mas_carpool_approver_uid} [put]
func (h *CarpoolManagementHandler) UpdateCarpoolMainApprover(c *gin.Context) {
	user := funcs.GetAuthenUser(c, h.Role)
	if c.IsAborted() {
		return
	}
	masCarpoolApproverUID := c.Param("mas_carpool_approver_uid")

	var existingApprover models.VmsMasCarpoolApprover
	if err := config.DB.Where("mas_carpool_approver_uid = ? AND is_deleted = ?", masCarpoolApproverUID, "0").First(&existingApprover).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Approver carpool not found", "message": messages.ErrNotfound.Error()})
		return
	}

	var carpool models.VmsMasCarpoolList
	queryRole := h.SetQueryRole(user, config.DB)
	if err := queryRole.Where("mas_carpool_uid = ? AND is_deleted = ?", existingApprover.MasCarpoolUID, "0").First(&carpool).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Carpool not found", "message": messages.ErrNotfound.Error()})
		return
	}

	if err := config.DB.Model(&models.VmsMasCarpoolApprover{}).
		Where("mas_carpool_uid = ?", existingApprover.MasCarpoolUID).
		Update("is_main_approver", "0").Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to update is_main_approver for all: %v", err), "message": messages.ErrInternalServer.Error()})
		return
	}

	if err := config.DB.Model(&existingApprover).Updates(map[string]interface{}{
		"is_main_approver": "1",
		"updated_by":       user.EmpID,
		"updated_at":       time.Now(),
	}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to update is_main_approver: %v", err), "message": messages.ErrInternalServer.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Approver carpool updated successfully", "carpool_name": GetCarpoolName(existingApprover.MasCarpoolUID)})
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
	if c.IsAborted() {
		return
	}
	masCarpoolApproverUID := c.Param("mas_carpool_approver_uid")

	var approver models.VmsMasCarpoolApprover
	if err := config.DB.Where("mas_carpool_approver_uid = ? AND is_deleted = ?", masCarpoolApproverUID, "0").First(&approver).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Carpool approver not found", "message": messages.ErrNotfound.Error()})
		return
	}

	var carpool models.VmsMasCarpoolList
	queryRole := h.SetQueryRole(user, config.DB)
	if err := queryRole.Where("mas_carpool_uid = ? AND is_deleted = ?", approver.MasCarpoolUID, "0").First(&carpool).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Carpool not found", "message": messages.ErrNotfound.Error()})
		return
	}

	if err := config.DB.Model(&approver).UpdateColumns(map[string]interface{}{
		"is_deleted": "1",
		"updated_by": user.EmpID,
		"updated_at": time.Now(),
	}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete carpool approver", "message": messages.ErrInternalServer.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Carpool approver deleted successfully", "carpool_name": GetCarpoolName(approver.MasCarpoolUID)})
}
