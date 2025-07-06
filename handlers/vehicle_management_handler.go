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

type VehicleManagementHandler struct {
	Role string
}

func (h *VehicleManagementHandler) SetQueryRole(user *models.AuthenUserEmp, query *gorm.DB) *gorm.DB {
	return query
}

func (h *VehicleManagementHandler) SetQueryRoleDept(user *models.AuthenUserEmp, query *gorm.DB) *gorm.DB {
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
// @Param is_tax_credit query string false "Filter by is_tax_credit"
// @Param order_by query string false "Order by vehicle_license_plate, vehicle_mileage, age,is_active"
// @Param order_dir query string false "Order direction: asc or desc"
// @Param page query int false "Page number (default: 1)"
// @Param limit query int false "Number of records per page (default: 10)"
// @Router /api/vehicle-management/search [get]
func (h *VehicleManagementHandler) SearchVehicles(c *gin.Context) {
	user := funcs.GetAuthenUser(c, h.Role)
	if c.IsAborted() {
		return
	}
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))    // Default: page 1
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10")) // Default: 10 items per page
	offset := (page - 1) * limit

	var vehicles []models.VmsMasVehicleManagementList

	query := h.SetQueryRole(user, config.DB)
	query = h.SetQueryRoleDept(user, query)
	query = query.Table("public.vms_mas_vehicle AS v").
		Select(`
		v.mas_vehicle_uid, v.vehicle_license_plate,v.vehicle_license_plate_province_short,v.vehicle_license_plate_province_full, v.vehicle_brand_name, v.vehicle_model_name, v.ref_vehicle_type_code,
		"CarTypeDetail" AS ref_vehicle_type_name, md.dept_short AS vehicle_owner_dept_short, d.fleet_card_no, v.is_tax_credit, d.vehicle_mileage,
		v.vehicle_registration_date, d.ref_vehicle_status_code, v.ref_fuel_type_id, v.is_active, mc.carpool_name vehicle_carpool_name, vs.ref_vehicle_status_short_name
	`).
		Joins("INNER JOIN public.vms_mas_vehicle_department AS d ON v.mas_vehicle_uid = d.mas_vehicle_uid AND d.is_deleted = '0' AND d.is_active = '1'").
		Joins("LEFT JOIN vms_ref_vehicle_type AS rt ON rt.ref_vehicle_type_code = v.ref_vehicle_type_code").
		Joins("LEFT JOIN vms_mas_department AS md ON md.dept_sap = d.vehicle_owner_dept_sap").
		Joins("LEFT JOIN vms_mas_carpool_vehicle AS mcv ON mcv.mas_vehicle_uid = v.mas_vehicle_uid AND mcv.is_deleted = '0'").
		Joins("LEFT JOIN vms_mas_carpool AS mc ON mc.mas_carpool_uid = mcv.mas_carpool_uid AND mc.is_deleted = '0'").
		Joins("LEFT JOIN vms_ref_vehicle_status AS vs ON vs.ref_vehicle_status_code = d.ref_vehicle_status_code").
		Where("v.is_deleted = ?", "0")

	search := strings.ToUpper(c.Query("search"))
	if search != "" {
		query = query.Where("v.vehicle_license_plate ILIKE ? OR vehicle_brand_name ILIKE ? OR vehicle_model_name ILIKE ?", "%"+search+"%", "%"+search+"%", "%"+search+"%")
	}

	if vehicleOwnerDeptSAP := c.Query("vehicle_owner_dept_sap"); vehicleOwnerDeptSAP != "" {
		query = query.Where("vehicle_owner_dept_sap = ?", vehicleOwnerDeptSAP)
	}

	if categoryCode := c.Query("ref_vehicle_category_code"); categoryCode != "" {
		query = query.Where("v.\"CarTypeDetail\" = ?", categoryCode)
	}

	if statusCodes := c.Query("ref_vehicle_status_code"); statusCodes != "" {
		statusCodeList := strings.Split(statusCodes, ",")
		query = query.Where("d.ref_vehicle_status_code IN (?)", statusCodeList)
	}

	if fuelTypeID := c.Query("ref_fuel_type_id"); fuelTypeID != "" {
		query = query.Where("ref_fuel_type_id = ?", fuelTypeID)
	}

	if isTaxCredit := c.Query("is_tax_credit"); isTaxCredit != "" {
		query = query.Where("is_tax_credit = ?", isTaxCredit)
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
		query = query.Order("v.is_active " + orderDir)
	default:
		query = query.Order("vehicle_license_plate " + orderDir) // Default ordering by license plate
	}

	query = query.Limit(limit).
		Offset(offset)

	if err := query.
		Preload("RefFuelType").
		Preload("RefVehicleStatus").
		Find(&vehicles).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "message": messages.ErrInternalServer.Error()})
		return
	}

	for i := range vehicles {
		vehicles[i].Age = funcs.CalculateAge(vehicles[i].VehicleGetDate)
		vehicles[i].RefVehicleStatus.RefVehicleStatusName = vehicles[i].RefVehicleStatusShortName
		funcs.TrimStringFields(&vehicles[i])
	}

	if len(vehicles) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"message": "No vehicles found",
			"pagination": gin.H{
				"page":       page,
				"limit":      limit,
				"totalPages": (total + int64(limit) - 1) / int64(limit), // Calculate total pages
				"vehicles":   vehicles,
			},
			"vehicles": []models.VmsMasVehicleManagementList{},
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

	query := h.SetQueryRole(user, config.DB)
	query = h.SetQueryRoleDept(user, query)
	if err := query.First(&vehicle, "mas_vehicle_uid = ? and is_deleted = ?", request.MasVehicleUID, "0").Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Vehicle not found", "message": messages.ErrNotfound.Error()})
		return
	}
	request.UpdatedAt = time.Now()
	request.UpdatedBy = user.EmpID

	if err := config.DB.Model(&vehicle).Update("is_active", request.IsActive).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to update: %v", err), "message": messages.ErrInternalServer.Error()})
		return
	}

	if err := config.DB.First(&result, "mas_vehicle_uid = ?", request.MasVehicleUID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Vehicle not found", "message": messages.ErrNotfound.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Updated successfully", "result": result})
}

// GetVehicleTimeLine godoc
// @Summary Get vehicle timeline
// @Description Get vehicle timeline by date range
// @Tags Vehicle-management
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Param start_date query string true "Start date (YYYY-MM-DD)" default(2025-05-01)
// @Param end_date query string true "End date (YYYY-MM-DD)" default(2025-06-30)
// @Param search query string false "Search by vehicle license plate, brand, or model"
// @Param vehicle_owner_dept_sap query string false "Filter by vehicle owner department SAP"
// @Param vehicel_car_type_detail query string false "Filter by Car type"
// @Param ref_timeline_status_id query string false "Filter by Timeline status (comma-separated, e.g., '1,2')"
// @Param page query int false "Page number (default: 1)"
// @Param limit query int false "Number of records per page (default: 10)"
// @Router /api/vehicle-management/timeline [get]
func (h *VehicleManagementHandler) GetVehicleTimeLine(c *gin.Context) {
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

	if startDate.After(endDate) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Start date must be before end date", "message": messages.ErrInvalidDate.Error()})
		return
	}
	// check if endDate - startDate > 3 months
	if endDate.Sub(startDate) > 90*24*time.Hour {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Date range cannot exceed 3 months", "message": messages.ErrInvalidDate.Error()})
		return
	}

	var vehicles []models.VehicleTimeLine
	query := h.SetQueryRole(user, config.DB)
	query = h.SetQueryRoleDept(user, query)
	query = query.Table("public.vms_mas_vehicle AS v").
		Select(`v.mas_vehicle_uid, v.vehicle_license_plate, v.vehicle_license_plate_province_short, 
				v.vehicle_license_plate_province_full, d.county, d.vehicle_get_date, d.vehicle_pea_id, 
				d.vehicle_license_plate_province_short, d.vehicle_license_plate_province_full, 
				public.fn_get_long_short_dept_name_by_dept_sap(d.vehicle_owner_dept_sap)  AS vehicle_dept_name,
				(select max(cp.carpool_name) from vms_mas_carpool cp, vms_mas_carpool_vehicle cpv where cpv.is_deleted = '0' and cpv.is_active = '1' and cpv.mas_carpool_uid = cp.mas_carpool_uid and cpv.mas_vehicle_uid = v.mas_vehicle_uid) AS vehicle_carpool_name, 
				v."CarTypeDetail" AS vehicle_car_type_detail, 0 AS vehicle_mileage,
				v.vehicle_brand_name,v.vehicle_model_name,
				public.fn_get_vehicle_distance_two_months(v.mas_vehicle_uid, ?) AS vehicle_distance`, startDate).
		Joins("LEFT JOIN (SELECT DISTINCT ON (mas_vehicle_uid) * FROM vms_mas_vehicle_department WHERE is_deleted = '0' AND is_active = '1' ORDER BY mas_vehicle_uid, created_at DESC) d ON v.mas_vehicle_uid = d.mas_vehicle_uid").
		Where("v.is_deleted = ?", "0")
	if refTimelineStatusID := c.Query("ref_timeline_status_id"); refTimelineStatusID != "" {
		refTimelineStatusIDList := strings.Split(refTimelineStatusID, ",")
		query = query.Where(`exists (select 1 from vms_trn_request r where r.mas_vehicle_uid = v.mas_vehicle_uid AND r.ref_request_status_code != '90' AND (
					('1' in (?) AND r.ref_request_status_code < '50') OR
					('2' in (?) AND r.ref_request_status_code >= '50' AND r.ref_request_status_code < '80' AND r.ref_trip_type_code = 0) OR 
					('3' in (?) AND r.ref_request_status_code >= '50' AND r.ref_request_status_code < '80' AND r.ref_trip_type_code = 1) OR
					('4' in (?) AND r.ref_request_status_code = '80')
				) AND
				 (reserve_start_datetime BETWEEN ? AND ? OR reserve_end_datetime BETWEEN ? AND ?)	
			)`, refTimelineStatusIDList, refTimelineStatusIDList, refTimelineStatusIDList, refTimelineStatusIDList, startDate, endDate, startDate, endDate)
	}

	if search := c.Query("search"); search != "" {
		query = query.Where("v.vehicle_license_plate ILIKE ? OR v.vehicle_brand_name ILIKE ? OR v.vehicle_model_name ILIKE ?", "%"+search+"%", "%"+search+"%", "%"+search+"%")
	}
	if vehicleOwnerDeptSAP := c.Query("vehicle_owner_dept_sap"); vehicleOwnerDeptSAP != "" {
		query = query.Where("d.vehicle_owner_dept_sap = ?", vehicleOwnerDeptSAP)
	}

	if vehicleCarTypeDetail := c.Query("vehicel_car_type_detail"); vehicleCarTypeDetail != "" {
		query = query.Where("v.\"CarTypeDetail\" = ?", vehicleCarTypeDetail)
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

	if err := query.Find(&vehicles).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "message": messages.ErrInternalServer.Error()})
		return
	}

	for i := range vehicles {
		parts := strings.Split(vehicles[i].VehicleDistance, ",")
		if len(parts) > 3 {
			vehicles[i].VehicleMileage = parts[3]
		}
		// Preload the vehicle requests for each vehicle
		query := config.DB.Table("vms_trn_request r").
			Preload("MasDriver").
			Where("mas_vehicle_uid = ? AND is_deleted = ? AND (reserve_start_datetime BETWEEN ? AND ? OR reserve_end_datetime BETWEEN ? AND ?)", vehicles[i].MasVehicleUID, "0", startDate, endDate, startDate, endDate).
			Where("ref_request_status_code != '90'")
		if refTimelineStatusID := c.Query("ref_timeline_status_id"); refTimelineStatusID != "" {
			refTimelineStatusIDList := strings.Split(refTimelineStatusID, ",")
			query = query.Where(`
					('1' in (?) AND r.ref_request_status_code < '50') OR
					('2' in (?) AND r.ref_request_status_code >= '50' AND r.ref_request_status_code < '80' AND r.ref_trip_type_code = 0) OR 
					('3' in (?) AND r.ref_request_status_code >= '50' AND r.ref_request_status_code < '80' AND r.ref_trip_type_code = 1) OR
					('4' in (?) AND r.ref_request_status_code = '80')
				`, refTimelineStatusIDList, refTimelineStatusIDList, refTimelineStatusIDList, refTimelineStatusIDList)
		}
		if err := query.Find(&vehicles[i].VehicleTrnRequests).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "message": messages.ErrInternalServer.Error()})
			return
		}

		for j := range vehicles[i].VehicleTrnRequests {
			vehicles[i].VehicleTrnRequests[j].TripDetails = []models.VmsTrnTripDetail{
				{
					TrnTripDetailUID: uuid.New().String(),
					VmsTrnTripDetailRequest: models.VmsTrnTripDetailRequest{
						TrnRequestUID:        vehicles[i].VehicleTrnRequests[j].TrnRequestUID,
						TripStartDatetime:    vehicles[i].VehicleTrnRequests[j].ReserveStartDatetime,
						TripEndDatetime:      vehicles[i].VehicleTrnRequests[j].ReserveEndDatetime,
						TripDeparturePlace:   vehicles[i].VehicleTrnRequests[j].WorkPlace,
						TripDestinationPlace: vehicles[i].VehicleTrnRequests[j].WorkPlace,
						TripStartMiles:       0,
						TripEndMiles:         0,
					},
				},
			}
			if vehicles[i].VehicleTrnRequests[j].RefRequestStatusCode == "80" {
				vehicles[i].VehicleTrnRequests[j].RefTimelineStatusID = "4"
				vehicles[i].VehicleTrnRequests[j].TimeLineStatus = "เสร็จสิ้น"
			} else if vehicles[i].VehicleTrnRequests[j].RefRequestStatusCode < "50" {
				vehicles[i].VehicleTrnRequests[j].RefTimelineStatusID = "1"
				vehicles[i].VehicleTrnRequests[j].TimeLineStatus = "รออนุมัติ"
			} else if vehicles[i].VehicleTrnRequests[j].RefTripTypeCode == 0 {
				vehicles[i].VehicleTrnRequests[j].RefTimelineStatusID = "2"
				vehicles[i].VehicleTrnRequests[j].TimeLineStatus = "ไป-กลับ"
			} else if vehicles[i].VehicleTrnRequests[j].RefTripTypeCode == 1 {
				vehicles[i].VehicleTrnRequests[j].RefTimelineStatusID = "3"
				vehicles[i].VehicleTrnRequests[j].TimeLineStatus = "ค้างแรม"
			}
			vehicles[i].VehicleTrnRequests[j].RefRequestStatusName = StatusNameMapUser[vehicles[i].VehicleTrnRequests[j].RefRequestStatusCode]
		}
	}
	thaiMonths := []string{"ม.ค.", "ก.พ.", "มี.ค.", "เม.ย.", "พ.ค.", "มิ.ย.", "ก.ค.", "ส.ค.", "ก.ย.", "ต.ค.", "พ.ย.", "ธ.ค."}
	lastMonthDate := time.Date(startDate.Year(), startDate.Month()-1, 1, 0, 0, 0, 0, startDate.Location())
	lastMonth := fmt.Sprintf("%s%02d", thaiMonths[lastMonthDate.Month()-1], (lastMonthDate.Year()+543)%100)

	c.JSON(http.StatusOK, gin.H{
		"pagination": gin.H{
			"total":      total,
			"page":       page,
			"limit":      pageSizeInt,
			"totalPages": (total + int64(pageSizeInt) - 1) / int64(pageSizeInt), // Calculate total pages
		},
		"last_month": lastMonth,
		"vehicles":   vehicles,
	})
}

// ReportTripDetail godoc
// @Summary Get vehicle report trip detail
// @Description Get vehicle report trip by date range
// @Tags Vehicle-management
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Param start_date query string true "Start date (YYYY-MM-DD)"
// @Param end_date query string true "End date (YYYY-MM-DD)"
// @Param show_all query string false "Show all vehicles (1 for true, 0 for false)"
// @Param mas_vehicle_uid body []string true "Array of vehicle mas_vehicle_uid"
// @Router /api/vehicle-management/report-trip-detail [post]
func (h *VehicleManagementHandler) ReportTripDetail(c *gin.Context) {
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
	var masVehicleUIDs []string
	if err := c.ShouldBindJSON(&masVehicleUIDs); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid mas_vehicle_uid format", "message": messages.ErrInvalidJSONInput.Error()})
		return
	}
	var tripReports []models.VehicleReportTripDetail

	query := h.SetQueryRole(user, config.DB)
	query = h.SetQueryRoleDept(user, query)
	query = query.Table("public.vms_mas_vehicle AS v").
		Select(`v.mas_vehicle_uid, v.vehicle_license_plate, v.vehicle_license_plate_province_short, 
				v.vehicle_license_plate_province_full, d.county, d.vehicle_get_date, d.vehicle_pea_id, 
				d.vehicle_license_plate_province_short, d.vehicle_license_plate_province_full, 
				public.fn_get_long_short_dept_name_by_dept_sap(d.vehicle_owner_dept_sap) AS vehicle_dept_name, 
				(select max(mc.carpool_name) from vms_mas_carpool mc, vms_mas_carpool_vehicle cpv where cpv.is_deleted = '0' and cpv.is_active = '1' and cpv.mas_carpool_uid = mc.mas_carpool_uid and cpv.mas_vehicle_uid = v.mas_vehicle_uid) AS vehicle_carpool_name, 
				v."CarTypeDetail" AS vehicle_car_type_detail,
				v.vehicle_brand_name,v.vehicle_model_name,
				td.trip_start_datetime, td.trip_end_datetime,td.trip_departure_place,td.trip_destination_place,td.trip_start_miles,td.trip_end_miles,td.trip_detail`).
		Joins("INNER JOIN public.vms_mas_vehicle_department AS d ON v.mas_vehicle_uid = d.mas_vehicle_uid AND d.is_deleted = '0' AND d.is_active = '1'").
		Joins("LEFT JOIN vms_trn_request r ON r.mas_vehicle_uid = v.mas_vehicle_uid AND ("+
			"r.reserve_start_datetime BETWEEN ? AND ? "+
			"OR r.reserve_end_datetime BETWEEN ? AND ? "+
			"OR ? BETWEEN r.reserve_start_datetime AND r.reserve_end_datetime "+
			"OR ? BETWEEN r.reserve_start_datetime AND r.reserve_end_datetime)",
			startDate, endDate, startDate, endDate, startDate, endDate).
		Joins("LEFT JOIN vms_trn_trip_detail td ON td.trn_request_uid = r.trn_request_uid AND td.is_deleted = '0'").
		Where("v.is_deleted = ? AND d.is_deleted = ? AND d.is_active = ?", "0", "0", "1")

	query = query.Where("v.mas_vehicle_uid::Text IN (?)", masVehicleUIDs).Debug()

	if err := query.Find(&tripReports).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "message": "Internal server error"})
		return
	}

	file := xlsx.NewFile()
	sheet, err := file.AddSheet("Trip Reports")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create Excel sheet", "message": err.Error()})
		return
	}

	// Add header row
	headerRow := sheet.AddRow()
	headers := []string{
		"รหัสยานพาหนะ", "ป้ายทะเบียน", "จังหวัด (ย่อ)", "จังหวัด (เต็ม)",
		"ชื่อหน่วยงาน", "ชื่อคาร์พูล", "รายละเอียดประเภทรถ", "เวลาเริ่มต้นการเดินทาง", "เวลาสิ้นสุดการเดินทาง",
		"สถานที่ออกเดินทาง", "สถานที่ปลายทาง", "ระยะไมล์เริ่มต้น", "ระยะไมล์สิ้นสุด", "รายละเอียดการเดินทาง",
	}
	for _, header := range headers {
		cell := headerRow.AddCell()
		cell.Value = header
	}

	// Add data rows
	for _, report := range tripReports {
		row := sheet.AddRow()
		row.AddCell().Value = report.VehiclePEAID
		row.AddCell().Value = report.VehicleLicensePlate
		row.AddCell().Value = report.VehicleLicensePlateProvinceShort
		row.AddCell().Value = report.VehicleLicensePlateProvinceFull
		row.AddCell().Value = report.VehicleDeptName
		row.AddCell().Value = report.CarpoolName
		row.AddCell().Value = report.VehicleCarTypeDetail
		row.AddCell().Value = report.TripStartDatetime.Format("2006-01-02 15:04:05")
		row.AddCell().Value = report.TripEndDatetime.Format("2006-01-02 15:04:05")
		row.AddCell().Value = report.TripDeparturePlace
		row.AddCell().Value = report.TripDestinationPlace
		row.AddCell().Value = strconv.FormatFloat(float64(report.TripStartMiles), 'f', 2, 64)
		row.AddCell().Value = strconv.FormatFloat(float64(report.TripEndMiles), 'f', 2, 64)
		row.AddCell().Value = report.TripDetail
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
	c.Header("Content-Disposition", "attachment; filename=trip_reports.xlsx")
	c.Header("File-Name", fmt.Sprintf("trip_reports_%s_to_%s.xlsx", startDate.Format("2006-01-02"), endDate.Format("2006-01-02")))
	c.Header("Content-Transfer-Encoding", "binary")
	if err := file.Write(c.Writer); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to write Excel file", "message": err.Error()})
		return
	}
}

// ReportAddFuel godoc
// @Summary Get vehicle report add fuel detail
// @Description Get vehicle report add fuel by date range
// @Tags Vehicle-management
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Param start_date query string true "Start date (YYYY-MM-DD)"
// @Param end_date query string true "End date (YYYY-MM-DD)"
// @Param show_all query string false "Show all vehicles (1 for true, 0 for false)"
// @Param mas_vehicle_uid body []string true "Array of vehicle mas_vehicle_uid"
// @Router /api/vehicle-management/report-add-fuel [post]
func (h *VehicleManagementHandler) ReportAddFuel(c *gin.Context) {
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
	var masVehicleUIDs []string
	if err := c.ShouldBindJSON(&masVehicleUIDs); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid mas_vehicle_uid format", "message": messages.ErrInvalidJSONInput.Error()})
		return
	}
	var fuelReports []models.VehicleReportAddFuel
	query := h.SetQueryRole(user, config.DB)
	query = h.SetQueryRoleDept(user, query)
	query = query.Table("public.vms_mas_vehicle AS v").
		Select(`v.mas_vehicle_uid, v.vehicle_license_plate, v.vehicle_license_plate_province_short, 
				v.vehicle_license_plate_province_full, d.county, d.vehicle_get_date, d.vehicle_pea_id, 
				public.fn_get_long_short_dept_name_by_dept_sap(d.vehicle_owner_dept_sap) AS vehicle_dept_name, 
				(select max(mc.carpool_name) from vms_mas_carpool mc, vms_mas_carpool_vehicle cpv where cpv.is_deleted = '0' and cpv.is_active = '1' and cpv.mas_carpool_uid = mc.mas_carpool_uid and cpv.mas_vehicle_uid = v.mas_vehicle_uid) AS vehicle_carpool_name, 
				v.vehicle_brand_name,v.vehicle_model_name,
				af.add_fuel_date_time, af.mile, af.tax_invoice_date, af.tax_invoice_no,af.price_per_liter,af.sum_liter,af.sum_price,
				(select ref_cost_type_name from vms_ref_cost_type where ref_cost_type_code = af.ref_cost_type_code) as cost_type_name,
				(select ref_oil_station_brand_name_th from vms_ref_oil_station_brand where ref_oil_station_brand_id = af.ref_oil_station_brand_id) as oil_station_brand_name,	
				(select ref_fuel_type_name_th from vms_ref_fuel_type where ref_fuel_type_id = af.ref_fuel_type_id) as fuel_type_name,
				(select ref_payment_type_name from vms_ref_payment_type where ref_payment_type_code = af.ref_payment_type_code) as payment_type_name
				`).
		Joins("INNER JOIN public.vms_mas_vehicle_department AS d ON v.mas_vehicle_uid = d.mas_vehicle_uid AND d.is_deleted = '0' AND d.is_active = '1'").
		Joins("LEFT JOIN vms_trn_request r ON r.mas_vehicle_uid = v.mas_vehicle_uid AND ("+
			"r.reserve_start_datetime BETWEEN ? AND ? "+
			"OR r.reserve_end_datetime BETWEEN ? AND ? "+
			"OR ? BETWEEN r.reserve_start_datetime AND r.reserve_end_datetime "+
			"OR ? BETWEEN r.reserve_start_datetime AND r.reserve_end_datetime)",
			startDate, endDate, startDate, endDate, startDate, endDate).
		Joins("LEFT JOIN vms_trn_add_fuel af ON af.is_deleted = '0' AND af.trn_request_uid = r.trn_request_uid").
		Where("v.is_deleted = ? AND d.is_deleted = ? AND d.is_active = ?", "0", "0", "1")

	query = query.Where("v.mas_vehicle_uid::Text IN (?)", masVehicleUIDs)

	if err := query.Find(&fuelReports).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "message": messages.ErrInternalServer.Error()})
		return
	}

	file := xlsx.NewFile()
	sheet, err := file.AddSheet("Fuel Reports")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create Excel sheet", "message": err.Error()})
		return
	}

	// Add header row
	headerRow := sheet.AddRow()
	headers := []string{
		"รหัสยานพาหนะ", "ป้ายทะเบียน", "จังหวัด (ย่อ)", "จังหวัด (เต็ม)", "ชื่อหน่วยงาน", "ชื่อคาร์พูล",
		"วันที่เติมน้ำมัน", "เลขไมล์", "วันที่ใบกำกับภาษี", "เลขที่ใบกำกับภาษี", "ราคาต่อลิตร", "จำนวนลิตร",
		"ราคารวม", "ประเภทค่าใช้จ่าย", "แบรนด์สถานีบริการน้ำมัน", "ประเภทน้ำมัน", "ประเภทการชำระเงิน",
	}
	for _, header := range headers {
		cell := headerRow.AddCell()
		cell.Value = header
	}

	// Add data rows
	for _, report := range fuelReports {
		row := sheet.AddRow()
		row.AddCell().Value = report.VehiclePEAID
		row.AddCell().Value = report.VehicleLicensePlate
		row.AddCell().Value = report.VehicleLicensePlateProvinceShort
		row.AddCell().Value = report.VehicleLicensePlateProvinceFull
		row.AddCell().Value = report.VehicleDeptName
		row.AddCell().Value = report.CarpoolName
		row.AddCell().Value = report.AddFuelDateTime.Format("2006-01-02 15:04:05")
		row.AddCell().Value = strconv.FormatFloat(float64(report.Mile), 'f', 2, 64)
		row.AddCell().Value = report.TaxInvoiceDate.Format("2006-01-02")
		row.AddCell().Value = report.TaxInvoiceNo
		row.AddCell().Value = strconv.FormatFloat(float64(report.PricePerLiter), 'f', 2, 64)
		row.AddCell().Value = strconv.FormatFloat(float64(report.SumLiter), 'f', 2, 64)
		row.AddCell().Value = strconv.FormatFloat(float64(report.SumPrice), 'f', 2, 64)
		row.AddCell().Value = report.RefCostType
		row.AddCell().Value = report.RefOilStationBrand
		row.AddCell().Value = report.RefFuelType
		row.AddCell().Value = report.RefPaymentType
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
	c.Header("Content-Disposition", "attachment; filename=fuel_reports.xlsx")
	c.Header("File-Name", fmt.Sprintf("fuel_reports_%s_to_%s.xlsx", startDate.Format("2006-01-02"), endDate.Format("2006-01-02")))
	c.Header("Content-Transfer-Encoding", "binary")
	if err := file.Write(c.Writer); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to write Excel file", "message": err.Error()})
		return
	}
}
