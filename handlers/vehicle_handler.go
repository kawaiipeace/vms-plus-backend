package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"vms_plus_be/config"
	"vms_plus_be/funcs"
	"vms_plus_be/messages"
	"vms_plus_be/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type VehicleHandler struct {
	Role string
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
// @Param ref_trip_type_code query int false "Filter by Trip Type Code (0: Round Trip, 1: Overnight)"
// @Param page query int false "Page number (default: 1)"
// @Param limit query int false "Number of records per page (default: 10)"
// @Router /api/vehicle/search [get]
func (h *VehicleHandler) SearchVehicles(c *gin.Context) {
	user := funcs.GetAuthenUser(c, h.Role)
	if c.IsAborted() {
		return
	}
	searchText := c.Query("search")            // Text search for brand name & license plate
	ownerDept := c.Query("vehicle_owner_dept") // Filter by vehicle owner department
	carType := c.Query("car_type")             // Filter by car type
	categoryCode := c.Query("category_code")   // Filter by car type
	ref_trip_type_code, _ := strconv.Atoi(c.Query("ref_trip_type_code"))

	//	user.BusinessArea = "J000"
	//carpool
	var carpools []models.VmsMasCarpoolCarBooking
	queryCarpool := config.DB
	queryCarpool = queryCarpool.Model(&models.VmsMasCarpoolCarBooking{})
	queryCarpool = queryCarpool.Where("is_deleted = '0' AND is_active = '1'")
	queryCarpool = queryCarpool.Where("ref_carpool_choose_car_id IN (2, 3)")
	queryCarpool = queryCarpool.Where("carpool_main_business_area= ?", user.BusinessArea)
	//queryCarpool = queryCarpool.Where("(select count(*) from vms_mas_carpool_vehicle cpv where is_deleted='0' and cpv.mas_carpool_uid=cp.mas_carpool_uid) > 0 AND " +
	//	"(select count(*) from vms_mas_carpool_approver cpa where is_deleted='0' and cpa.mas_carpool_uid=cp.mas_carpool_uid) > 0")
	queryCarpool.
		Preload("RefCarpoolChooseCar").
		Table("vms_mas_carpool cp").
		Find(&carpools)

	totalGroups := len(carpools)
	// Pagination parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))        // Default page = 1
	pageLimit, _ := strconv.Atoi(c.DefaultQuery("limit", "10")) // Default limit = 10
	offset := (page - 1) * pageLimit                            // Calculate offset
	limit := pageLimit

	if page == 1 {
		limit = limit - totalGroups
	} else {
		offset = offset - totalGroups
		carpools = make([]models.VmsMasCarpoolCarBooking, 0)
	}

	var vehicles []models.VmsMasVehicleList
	var total int64

	query := config.DB.Table("vms_mas_vehicle v").Select("*")
	query = query.Joins("LEFT JOIN vms_mas_vehicle_department vd ON v.mas_vehicle_uid = vd.mas_vehicle_uid")
	query = query.Where("v.is_deleted = '0'")
	query = query.Where("vd.ref_vehicle_status_code = '0' and vd.is_deleted = '0' and vd.is_active = '1'")
	query = query.Where("(vd.bureau_dept_sap = ?) OR (bureau_ba like ? AND (ref_other_use_code = 2 OR ref_other_use_code= 1 AND ? = 0))",
		user.BureauDeptSap, user.BusinessArea[:1]+"%", ref_trip_type_code)
	//ref_other_use_code = 2 -> ref_trip_type_code=1 ค้างแรม
	//ref_other_use_code = 1 -> ref_trip_type_code=0 ไปกลับ

	// Apply text search (VehicleBrandName OR VehicleLicensePlate)
	if searchText != "" {
		query = query.Where("vehicle_brand_name ILIKE ? OR vehicle_model_name ILIKE ? OR v.vehicle_license_plate ILIKE ?", "%"+searchText+"%", "%"+searchText+"%", "%"+searchText+"%")
	}

	// Apply filters if provided
	if ownerDept != "" {
		query = query.Where("vehicle_owner_dept_sap = ?", ownerDept)
	}
	if carType != "" {
		query = query.Where("car_type = ?", carType)
	}
	if categoryCode != "" {
		query = query.Where("ref_vehicle_type_code = ?", categoryCode)
	}

	// Count total records
	query.Count(&total)
	query = query.Select("v.*")
	// Execute query with pagination
	query.Offset(offset).Limit(limit).Find(&vehicles)
	vehicles = models.AssignVehicleImageFromIndex(vehicles)
	for i := range vehicles {
		funcs.TrimStringFields(&vehicles[i])
	}
	// Respond with JSON
	c.JSON(http.StatusOK, gin.H{
		"pagination": gin.H{
			"total":       total,
			"totalGroups": totalGroups,
			"page":        page,
			"limit":       pageLimit,
			"totalPages":  (total + int64(limit) - 1) / int64(limit), // Calculate total pages
		},
		"vehicles": vehicles,
		"carpools": carpools,
	})
}

// SearchBookingVehicles godoc
// @Summary Search vehicles by brand, license plate, and filters
// @Description Retrieves vehicles based on search text, department, and car type filters
// @Tags Vehicle
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Param emp_id query string false "Employee ID (emp_id) default(700001)"
// @Param start_date query string false "Start Date (YYYY-MM-DD HH:mm:ss)" default(2025-05-30 08:00:00)
// @Param end_date query string false "End Date (YYYY-MM-DD HH:mm:ss)" default(2025-05-30 16:00:00)
// @Param search query string false "Search text (Vehicle Brand Name or License Plate)"
// @Param vehicle_owner_dept query string false "Filter by Vehicle Owner Department"
// @Param car_type query string false "Filter by Car Type"
// @Param category_code query string false "Filter by Vehicle Category Code"
// @Param page query int false "Page number (default: 1)"
// @Param limit query int false "Number of records per page (default: 10)"
// @Router /api/vehicle/search-booking [get]
func (h *VehicleHandler) SearchBookingVehicles(c *gin.Context) {
	user := funcs.GetAuthenUser(c, h.Role)
	if c.IsAborted() {
		return
	}

	empID := c.Query("emp_id")
	bureauDeptSap := user.BureauDeptSap
	businessArea := user.BusinessArea

	if empID != "" {
		empUser := funcs.GetUserEmpInfo(empID)
		bureauDeptSap = empUser.BureauDeptSap
		businessArea = empUser.BusinessArea
	}

	startDate := c.Query("start_date")
	endDate := c.Query("end_date")
	if startDate == "" || endDate == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Start Date and End Date are required", "message": messages.ErrInvalidRequest.Error()})
		return
	}
	StartTimeWithZone, err1 := models.GetTimeWithZone(startDate)
	EndTimeWithZone, err2 := models.GetTimeWithZone(endDate)
	if err1 != nil || err2 != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get start time with zone", "message": messages.ErrInternalServer.Error()})
		return
	}

	fmt.Println(StartTimeWithZone)
	fmt.Println(EndTimeWithZone)

	var vehicleCanBookings []models.VmsMasVehicleCanBooking

	queryCanBooking := config.DB.Raw(`SELECT * FROM fn_get_available_vehicles_view (?, ?, ?, ?)`,
		StartTimeWithZone, EndTimeWithZone, bureauDeptSap, businessArea)
	err := queryCanBooking.Scan(&vehicleCanBookings).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch available vehicles", "message": messages.ErrInternalServer.Error()})
		return
	}

	searchText := c.Query("search")            // Text search for brand name & license plate
	ownerDept := c.Query("vehicle_owner_dept") // Filter by vehicle owner department
	carType := c.Query("car_type")             // Filter by car type
	categoryCode := c.Query("category_code")   // Filter by car type
	masVehicleUIDs := make([]string, 0)
	masCarpoolUIDs := make([]string, 0)
	adminChooseDriverMasCarpoolUIDs := make([]string, 0)
	systemChooseDriverMasCarpoolUIDs := make([]string, 0)
	adminChooseDriverMasVehicleUIDs := make([]string, 0)
	systemChooseDriverMasVehicleUIDs := make([]string, 0)
	for _, vehicleCanBooking := range vehicleCanBookings {

		if vehicleCanBooking.RefCarpoolChooseCarID == 2 || vehicleCanBooking.RefCarpoolChooseCarID == 3 {
			masCarpoolUIDs = append(masCarpoolUIDs, vehicleCanBooking.MasCarpoolUID)
		} else {
			masVehicleUIDs = append(masVehicleUIDs, vehicleCanBooking.MasVehicleUID)
		}
		if vehicleCanBooking.RefCarpoolChooseDriverID == 2 {
			adminChooseDriverMasCarpoolUIDs = append(adminChooseDriverMasCarpoolUIDs, vehicleCanBooking.MasCarpoolUID)
			adminChooseDriverMasVehicleUIDs = append(adminChooseDriverMasVehicleUIDs, vehicleCanBooking.MasVehicleUID)
		}
		if vehicleCanBooking.RefCarpoolChooseDriverID == 3 {
			systemChooseDriverMasCarpoolUIDs = append(systemChooseDriverMasCarpoolUIDs, vehicleCanBooking.MasCarpoolUID)
			systemChooseDriverMasVehicleUIDs = append(systemChooseDriverMasVehicleUIDs, vehicleCanBooking.MasVehicleUID)
		}
	}

	//carpool
	var carpools []models.VmsMasCarpoolCarBooking
	queryCarpool := config.DB.Model(&models.VmsMasCarpoolCarBooking{})
	queryCarpool = queryCarpool.Where("mas_carpool_uid IN (?) AND is_deleted = '0' AND is_active = '1'", masCarpoolUIDs)
	queryCarpool.Preload("RefCarpoolChooseCar").
		Table("vms_mas_carpool cp").
		Find(&carpools)

	for i := range carpools {
		if funcs.Contains(adminChooseDriverMasCarpoolUIDs, carpools[i].MasCarpoolUID) {
			carpools[i].IsAdminChooseDriver = true
		} else {
			carpools[i].IsAdminChooseDriver = false
		}
		if funcs.Contains(systemChooseDriverMasCarpoolUIDs, carpools[i].MasCarpoolUID) {
			carpools[i].IsSystemChooseDriver = true
		} else {
			carpools[i].IsSystemChooseDriver = false
		}
	}
	if len(ownerDept) > 10 {
		filteredCarpools := make([]models.VmsMasCarpoolCarBooking, 0, len(carpools))
		for _, carpool := range carpools {
			if carpool.MasCarpoolUID == ownerDept {
				filteredCarpools = append(filteredCarpools, carpool)
			}
		}
		carpools = filteredCarpools
	}
	totalGroups := len(carpools)
	// Pagination parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))        // Default page = 1
	pageLimit, _ := strconv.Atoi(c.DefaultQuery("limit", "10")) // Default limit = 10
	offset := (page - 1) * pageLimit                            // Calculate offset
	limit := pageLimit
	if page == 1 {
		limit = pageLimit - totalGroups
	} else {
		offset = offset - totalGroups
		carpools = make([]models.VmsMasCarpoolCarBooking, 0)
	}

	var vehicles []models.VmsMasVehicleList
	var total int64

	var carpoolReadyUIDs []string
	if err := config.DB.Table("vms_mas_carpool cp").Where("cp.is_active = '1' AND cp.is_deleted = '0'").
		Where("(select count(*) from vms_mas_carpool_vehicle cpv where cpv.is_deleted='0' and cpv.mas_carpool_uid=cp.mas_carpool_uid) > 0 AND "+
			"(select count(*) from vms_mas_carpool_approver cpa where cpa.is_deleted='0' and cpa.mas_carpool_uid=cp.mas_carpool_uid) > 0").
		Pluck("mas_carpool_uid", &carpoolReadyUIDs).Error; err != nil {
		fmt.Println("ready carpool error", err)
		//return
	}

	query := config.DB.Table("vms_mas_vehicle v").Select("v.*,vd.vehicle_owner_dept_sap,vd.vehicle_pea_id,vd.fleet_card_no,cpv.mas_carpool_uid as carpool_uid,vi.vehicle_img_file vehicle_img,cp.carpool_name")
	query = query.Joins("LEFT JOIN (SELECT DISTINCT ON (mas_vehicle_uid) * FROM vms_mas_vehicle_department WHERE is_deleted = '0' AND is_active = '1' ORDER BY mas_vehicle_uid, created_at DESC) vd ON v.mas_vehicle_uid = vd.mas_vehicle_uid")
	query = query.Joins("LEFT JOIN (SELECT DISTINCT ON (mas_vehicle_uid) * FROM vms_mas_carpool_vehicle WHERE is_deleted = '0' AND is_active = '1' ORDER BY mas_vehicle_uid, created_at DESC) cpv ON cpv.mas_vehicle_uid = v.mas_vehicle_uid")
	query = query.Joins("LEFT JOIN vms_mas_carpool cp ON cp.mas_carpool_uid = cpv.mas_carpool_uid")
	query = query.Joins("LEFT JOIN (SELECT DISTINCT ON (mas_vehicle_uid) * FROM vms_mas_vehicle_img WHERE ref_vehicle_img_side_code = 1 ORDER BY mas_vehicle_uid, ref_vehicle_img_side_code) vi ON vi.mas_vehicle_uid = v.mas_vehicle_uid")
	query = query.Where("v.mas_vehicle_uid IN (?) AND v.is_deleted = '0' AND v.is_active = '1'", masVehicleUIDs)
	query = query.Where("cpv.mas_carpool_uid is null OR (cp.is_active = '1' AND cp.is_deleted = '0' AND cp.mas_carpool_uid IN (?))", carpoolReadyUIDs)

	if searchText != "" {
		query = query.Where("vehicle_brand_name ILIKE ? OR vehicle_model_name ILIKE ? OR v.vehicle_license_plate ILIKE ?", "%"+searchText+"%", "%"+searchText+"%", "%"+searchText+"%")
	}

	// Apply filters if provided
	if ownerDept != "" {
		if len(ownerDept) > 10 {
			//search by mas_carpool_uid
			query = query.Where("cpv.mas_carpool_uid = ?", ownerDept)
		} else {
			query = query.Where("vehicle_owner_dept_sap = ?", ownerDept)
		}
	}
	if carType != "" {
		query = query.Where("\"CarTypeDetail\" = ?", carType)
	}
	if categoryCode != "" {
		query = query.Where("\"CarTypeDetail\" = ?", categoryCode)
	}

	// Count total records
	query.Count(&total)
	query = query.Model(&models.VmsMasVehicleList{})
	// Execute query with pagination
	query = query.Preload("RefFuelType")
	query.Offset(offset).Limit(limit).Find(&vehicles)

	for i := range vehicles {
		if funcs.Contains(adminChooseDriverMasVehicleUIDs, vehicles[i].MasVehicleUID) {
			vehicles[i].IsAdminChooseDriver = true
		} else {
			vehicles[i].IsAdminChooseDriver = false
		}
		if funcs.Contains(systemChooseDriverMasVehicleUIDs, vehicles[i].MasVehicleUID) {
			vehicles[i].IsSystemChooseDriver = true
		} else {
			vehicles[i].IsSystemChooseDriver = false
		}
		funcs.TrimStringFields(&vehicles[i])
		if vehicles[i].CarpoolName != "" {
			vehicles[i].VehicleOwnerDeptSAP = ""
			vehicles[i].VehicleOwnerDeptShort = vehicles[i].CarpoolName + "(carpool)"
		} else {
			vehicles[i].VehicleOwnerDeptShort = funcs.GetDeptSAPShort(vehicles[i].VehicleOwnerDeptSAP)
		}
	}
	// Respond with JSON
	c.JSON(http.StatusOK, gin.H{
		"pagination": gin.H{
			"total":       total,
			"totalGroups": totalGroups,
			"page":        page,
			"limit":       pageLimit,
			"totalPages":  (total + int64(limit) - 1) / int64(limit), // Calculate total pages
		},
		"vehicles": vehicles,
		"carpools": carpools,
	})
}

// SearchBookingVehiclesCarpool godoc
// @Summary Search vehicel from carpools by brand, license plate, and filters
// @Description Retrieves vehicles from carpools based on search text, department, and car type filters
// @Tags Vehicle
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Param emp_id query string false "Employee ID (emp_id) default(700001)"
// @Param start_date query string false "Start Date (YYYY-MM-DD HH:mm:ss)" default(2025-05-30 08:00:00)
// @Param end_date query string false "End Date (YYYY-MM-DD HH:mm:ss)" default(2025-05-30 16:00:00)
// @Param search query string false "Search text (Vehicle Brand Name or License Plate)"
// @Param vehicle_owner_dept query string false "Filter by Vehicle Owner Department"
// @Param mas_carpool_uid query string false "Filter by MasCarpoolUID"
// @Param car_type query string false "Filter by Car Type"
// @Param page query int false "Page number (default: 1)"
// @Param limit query int false "Number of records per page (default: 10)"
// @Router /api/vehicle/search-booking-carpool [get]
func (h *VehicleHandler) SearchBookingVehiclesCarpool(c *gin.Context) {
	user := funcs.GetAuthenUser(c, h.Role)
	if c.IsAborted() {
		return
	}

	empID := c.Query("emp_id")
	bureauDeptSap := user.BureauDeptSap
	businessArea := user.BusinessArea

	if empID != "" {
		empUser := funcs.GetUserEmpInfo(empID)
		bureauDeptSap = empUser.BureauDeptSap
		businessArea = empUser.BusinessArea
	}

	startDate := c.Query("start_date")
	endDate := c.Query("end_date")
	if startDate == "" || endDate == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Start Date and End Date are required", "message": messages.ErrInvalidRequest.Error()})
		return
	}
	StartTimeWithZone, err1 := models.GetTimeWithZone(startDate)
	EndTimeWithZone, err2 := models.GetTimeWithZone(endDate)
	if err1 != nil || err2 != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get start time with zone", "message": messages.ErrInternalServer.Error()})
		return
	}

	fmt.Println(StartTimeWithZone)
	fmt.Println(EndTimeWithZone)

	var vehicleCanBookings []models.VmsMasVehicleCanBooking
	masCarpoolUID := c.Query("mas_carpool_uid")
	queryCanBooking := config.DB.Raw(`SELECT * FROM fn_get_available_vehicles_view (?, ?, ?, ?) where mas_carpool_uid = ?`,
		StartTimeWithZone, EndTimeWithZone, bureauDeptSap, businessArea, masCarpoolUID)
	err := queryCanBooking.Scan(&vehicleCanBookings).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch available vehicles", "message": messages.ErrInternalServer.Error()})
		return
	}
	masVehicleUIDs := make([]string, 0)
	for _, vehicleCanBooking := range vehicleCanBookings {
		masVehicleUIDs = append(masVehicleUIDs, vehicleCanBooking.MasVehicleUID)
	}
	searchText := c.Query("search")            // Text search for brand name & license plate
	ownerDept := c.Query("vehicle_owner_dept") // Filter by vehicle owner department
	carType := c.Query("car_type")             // Filter by car type

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))        // Default page = 1
	pageLimit, _ := strconv.Atoi(c.DefaultQuery("limit", "10")) // Default limit = 10
	offset := (page - 1) * pageLimit                            // Calculate offset
	limit := pageLimit

	var vehicles []models.VmsMasVehicleList
	var total int64

	query := config.DB.Table("vms_mas_vehicle v").Select("v.*,vd.vehicle_owner_dept_sap,vd.vehicle_mileage,0 as last_month_mileage,vd.vehicle_pea_id,cpv.mas_carpool_uid as carpool_uid,vi.vehicle_img_file vehicle_img,cp.carpool_name")
	query = query.Joins("LEFT JOIN (SELECT DISTINCT ON (mas_vehicle_uid) * FROM vms_mas_vehicle_department WHERE is_deleted = '0' AND is_active = '1' ORDER BY mas_vehicle_uid, created_at DESC) vd ON v.mas_vehicle_uid = vd.mas_vehicle_uid")
	query = query.Joins("LEFT JOIN (SELECT DISTINCT ON (mas_vehicle_uid) * FROM vms_mas_carpool_vehicle WHERE is_deleted = '0' AND is_active = '1' ORDER BY mas_vehicle_uid, created_at DESC) cpv ON cpv.mas_vehicle_uid = v.mas_vehicle_uid")
	query = query.Joins("LEFT JOIN vms_mas_carpool cp ON cp.mas_carpool_uid = cpv.mas_carpool_uid")
	query = query.Joins("LEFT JOIN (SELECT DISTINCT ON (mas_vehicle_uid) * FROM vms_mas_vehicle_img WHERE ref_vehicle_img_side_code = 1 ORDER BY mas_vehicle_uid, ref_vehicle_img_side_code) vi ON vi.mas_vehicle_uid = v.mas_vehicle_uid")
	query = query.Where("v.mas_vehicle_uid IN (?) AND v.is_deleted = '0'", masVehicleUIDs)
	if searchText != "" {
		query = query.Where("vehicle_brand_name ILIKE ? OR vehicle_model_name ILIKE ? OR v.vehicle_license_plate ILIKE ?", "%"+searchText+"%", "%"+searchText+"%", "%"+searchText+"%")
	}

	// Apply filters if provided
	if ownerDept != "" {
		query = query.Where("vehicle_owner_dept_sap = ?", ownerDept)
	}
	if carType != "" {
		query = query.Where("\"CarTypeDetail\" = ?", carType)
	}

	// Count total records
	query.Count(&total)
	query = query.Model(&models.VmsMasVehicleList{})
	// Execute query with pagination
	query.Offset(offset).Limit(limit).Find(&vehicles)

	for i := range vehicles {
		funcs.TrimStringFields(&vehicles[i])
		if vehicles[i].CarpoolName != "" {
			vehicles[i].VehicleOwnerDeptSAP = ""
			vehicles[i].VehicleOwnerDeptShort = vehicles[i].CarpoolName
		} else {
			vehicles[i].VehicleOwnerDeptShort = funcs.GetDeptSAPShort(vehicles[i].VehicleOwnerDeptSAP)
		}
	}
	// Respond with JSON
	c.JSON(http.StatusOK, gin.H{
		"pagination": gin.H{
			"total":      total,
			"page":       page,
			"limit":      pageLimit,
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
// @Param mas_vehicle_uid path string true "MasVehicleUID (mas_vehicle_uid)"
// @Router /api/vehicle/{mas_vehicle_uid} [get]
func (h *VehicleHandler) GetVehicle(c *gin.Context) {
	funcs.GetAuthenUser(c, h.Role)
	if c.IsAborted() {
		return
	}
	vehicleID := c.Param("mas_vehicle_uid")

	// Parse the string ID to uuid.UUID
	parsedID, err := uuid.Parse(vehicleID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid vehicle ID", "message": messages.ErrInvalidUID.Error()})
		return
	}

	// Fetch the vehicle record from the database
	var vehicle models.VmsMasVehicle
	if err := config.DB.Preload("RefFuelType").
		First(&vehicle, "mas_vehicle_uid = ? AND is_deleted = '0'", parsedID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Vehicle not found", "message": messages.ErrNotfound.Error()})
		return
	}
	vehicle.Age = funcs.CalculateAgeInt(vehicle.VehicleRegistrationDate)
	var vehicleImgs []models.VmsMasVehicleImg
	if err := config.DB.Where("mas_vehicle_uid = ?", parsedID).Find(&vehicleImgs).Error; err == nil {
		vehicle.VehicleImgs = make([]string, 0)
		for _, img := range vehicleImgs {
			vehicle.VehicleImgs = append(vehicle.VehicleImgs, img.VehicleImgFile)
		}
	}

	// Get vehicle department details
	if err := config.DB.Where("mas_vehicle_uid = ?", parsedID).
		Select("*,public.fn_get_oil_station_eng_by_fleetcard(fleet_card_no) as fleet_card_oil_stations").
		Where("is_deleted = '0' AND is_active = '1'").
		First(&vehicle.VehicleDepartment).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Vehicle department not found", "message": messages.ErrNotfound.Error()})
		return
	}
	vehicle.VehicleDepartment.VehicleOwnerDeptShort = funcs.GetDeptSAPShort(vehicle.VehicleDepartment.VehicleOwnerDeptSap)
	//check if vehicle is carpool
	var carpoolVehicle models.VmsMasCarpoolVehicle
	masCarpoolUID := ""
	if err := config.DB.Where("mas_vehicle_uid = ? AND is_deleted = '0' AND is_active = '1'", parsedID).First(&carpoolVehicle).Error; err == nil {
		masCarpoolUID = carpoolVehicle.MasCarpoolUID
	}
	if masCarpoolUID != "" {
		var carpoolAdmin models.VmsMasCarpoolAdmin
		if err := config.DB.Where("mas_carpool_uid = ? AND is_deleted = '0' AND is_active = '1'", masCarpoolUID).
			Select("admin_emp_no,admin_emp_name,admin_dept_sap,admin_position,mobile_contact_number,internal_contact_number").
			Order("is_main_admin DESC").
			First(&carpoolAdmin).Error; err == nil {
			vehicle.VehicleDepartment.VehicleUser.EmpID = carpoolAdmin.AdminEmpNo
			vehicle.VehicleDepartment.VehicleUser.FullName = carpoolAdmin.AdminEmpName
			vehicle.VehicleDepartment.VehicleUser.DeptSAP = carpoolAdmin.AdminDeptSap
			vehicle.VehicleDepartment.VehicleUser.DeptSAPFull = funcs.GetDeptSAPFull(carpoolAdmin.AdminDeptSap)
			vehicle.VehicleDepartment.VehicleUser.DeptSAPShort = funcs.GetDeptSAPShort(carpoolAdmin.AdminDeptSap)
			vehicle.VehicleDepartment.VehicleUser.ImageUrl = funcs.GetEmpImage(carpoolAdmin.AdminEmpNo)
			vehicle.VehicleDepartment.VehicleUser.Position = carpoolAdmin.AdminPosition
			vehicle.VehicleDepartment.VehicleUser.TelMobile = carpoolAdmin.MobileContactNumber
			vehicle.VehicleDepartment.VehicleUser.TelInternal = carpoolAdmin.InternalContactNumber
			vehicle.VehicleDepartment.VehicleUser.IsEmployee = true
		}
		var carpool models.VmsMasCarpoolList
		if err := config.DB.Table("vms_mas_carpool").Where("mas_carpool_uid = ? AND is_deleted = '0'", masCarpoolUID).First(&carpool).Error; err == nil {
			vehicle.VehicleDepartment.VehicleOwnerDeptSap = ""
			vehicle.VehicleDepartment.VehicleOwnerDeptShort = carpool.CarpoolName
		}

	}

	funcs.TrimStringFields(&vehicle)

	c.JSON(http.StatusOK, vehicle)
}

// GetTypes godoc
// @Summary Get vehicle types
// @Description Fetches all vehicle types, optionally filtered by name
// @Tags Vehicle
// @Accept json
// @Produce json
// @Param emp_id query string false "Vehicle User EmpID (emp_id)" default(700001)
// @Param start_date query string false "Start Date (YYYY-MM-DD HH:mm:ss)" default(2025-05-30 08:00:00)
// @Param end_date query string false "End Date (YYYY-MM-DD HH:mm:ss)" default(2025-05-30 16:00:00)
// @Param name query string false "Filter by vehicle type name (partial match)"
// @Param mas_carpool_uid query string false "Filter by mas_carpool_uid"
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Router /api/vehicle/types [get]
func (h *VehicleHandler) GetTypes(c *gin.Context) {
	user := funcs.GetAuthenUser(c, "*")
	var vehicleTypes []models.VmsRefVehicleType
	name := c.Query("name") // Get the 'name' query parameter

	empID := c.Query("emp_id")
	bureauDeptSap := user.BureauDeptSap
	businessArea := user.BusinessArea

	if empID != "" {
		empUser := funcs.GetUserEmpInfo(empID)
		bureauDeptSap = empUser.BureauDeptSap
		businessArea = empUser.BusinessArea
	}

	startDate := c.Query("start_date")
	endDate := c.Query("end_date")
	if startDate == "" || endDate == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Start Date and End Date are required", "message": messages.ErrInvalidRequest.Error()})
		return
	}
	masCarpoolUID := c.Query("mas_carpool_uid")
	query := config.DB.Raw(`SELECT "CarTypeDetail" ref_vehicle_type_name,count(*) AS available_units FROM fn_get_available_vehicles_view (?, ?, ?, ?)
		WHERE "CarTypeDetail" ILIKE ? AND (? = '' OR mas_carpool_uid::text = ?)
		group by "CarTypeDetail"`,
		startDate, endDate, bureauDeptSap, businessArea, "%"+name+"%", masCarpoolUID, masCarpoolUID)

	if masCarpoolUID != "" {
		query = query.Where("mas_carpool_uid = ?", masCarpoolUID)
	}
	err := query.Scan(&vehicleTypes).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch available vehicles", "message": messages.ErrInternalServer.Error()})
		return
	}

	vehicleTypes = models.AssignTypeImageFromIndex(vehicleTypes)
	if len(vehicleTypes) == 0 {
		vehicleTypes = []models.VmsRefVehicleType{}
	}
	// Respond with JSON
	c.JSON(http.StatusOK, vehicleTypes)
}

// GetCarTypeDetails godoc
// @Summary Get car type details
// @Description Fetches details of car types including their names and descriptions
// @Tags Vehicle
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Router /api/vehicle/car-types-by-detail [get]
func (h *VehicleHandler) GetCarTypeDetails(c *gin.Context) {
	var carTypeDetails []models.VmsRefCarTypeDetail

	query := `
		SELECT 
			DISTINCT trim("CarTypeDetail") AS car_type_detail
		FROM 
			vms_mas_vehicle
		WHERE 
			is_deleted = '0' AND "CarTypeDetail">''
		GROUP BY 
			"CarTypeDetail"
	`

	// Execute the query
	if err := config.DB.Raw(query).Scan(&carTypeDetails).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch car type details", "message": messages.ErrInternalServer.Error()})
		return
	}

	// Respond with the result
	c.JSON(http.StatusOK, carTypeDetails)
}

// GetDepartments godoc
// @Summary Get department list
// @Description Fetches a list of departments grouped by dept_sap, including dept_short and dept_full
// @Tags Vehicle
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Param emp_id query string false "Employee ID (emp_id) default(700001)"
// @Param start_date query string false "Start Date (YYYY-MM-DD HH:mm:ss)" default(2025-05-30 08:00:00)
// @Param end_date query string false "End Date (YYYY-MM-DD HH:mm:ss)" default(2025-05-30 16:00:00)
// @Router /api/vehicle/departments [get]
func (h *VehicleHandler) GetDepartments(c *gin.Context) {
	user := funcs.GetAuthenUser(c, "*")
	var departments []struct {
		DeptSap   string `gorm:"column:dept_sap" json:"dept_sap"`
		DeptShort string `gorm:"column:dept_short" json:"dept_short"`
		DeptFull  string `gorm:"column:dept_full" json:"dept_full"`
	}
	empID := c.Query("emp_id")
	bureauDeptSap := user.BureauDeptSap
	businessArea := user.BusinessArea

	if empID != "" {
		empUser := funcs.GetUserEmpInfo(empID)
		bureauDeptSap = empUser.BureauDeptSap
		businessArea = empUser.BusinessArea
	}

	startDate := c.Query("start_date")
	endDate := c.Query("end_date")
	if startDate == "" || endDate == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Start Date and End Date are required", "message": messages.ErrInvalidRequest.Error()})
		return
	}

	query := config.DB.Raw(`SELECT CASE WHEN carpool_name!='' THEN mas_carpool_uid::text ELSE vehicle_owner_dept_sap END AS dept_sap,
		max(CASE WHEN carpool_name!='' THEN carpool_name ELSE fn_get_long_short_dept_name_by_dept_sap(vehicle_owner_dept_sap) END) AS dept_short,
		max(CASE WHEN carpool_name!='' THEN carpool_name ELSE fn_get_long_full_dept_name_by_dept_sap(vehicle_owner_dept_sap) END) AS dept_full
	 FROM fn_get_available_vehicles_view (?, ?, ?, ?) group by (CASE WHEN carpool_name!='' THEN mas_carpool_uid::text ELSE vehicle_owner_dept_sap END)`,
		startDate, endDate, bureauDeptSap, businessArea)

	err := query.Scan(&departments).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch available vehicles", "message": messages.ErrInternalServer.Error()})
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
// @Param mas_vehicle_uid query string false "MasVehicleUID (mas_vehicle_uid)"
// @Param mas_carpool_uid query string false "Filter by MasCarpoolUID"
// @Param emp_id query string false "Employee ID (emp_id) default(700001)"
// @Param start_date query string false "Start Date (YYYY-MM-DD HH:mm:ss)" default(2025-05-30 08:00:00)
// @Param end_date query string false "End Date (YYYY-MM-DD HH:mm:ss)" default(2025-05-30 16:00:00)
// @Param work_type query string false "work type to search (0: ไป-กลับ,1: ค้างคืน)" default(0)
// @Param search query string false "Search text (Vehicle Brand Name or License Plate)"
// @Param vehicle_owner_dept query string false "Filter by Vehicle Owner Department"
// @Router /api/vehicle-info [get]
func (h *VehicleHandler) GetVehicleInfo(c *gin.Context) {
	user := funcs.GetAuthenUser(c, "*")
	if c.IsAborted() {
		return
	}
	masVehicleUID := c.Query("mas_vehicle_uid")
	masCarpoolUID := c.Query("mas_carpool_uid")
	empID := c.Query("emp_id")
	bureauDeptSap := user.BureauDeptSap
	businessArea := user.BusinessArea

	if empID != "" {
		empUser := funcs.GetUserEmpInfo(empID)
		bureauDeptSap = empUser.BureauDeptSap
		businessArea = empUser.BusinessArea
	}
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")
	workType := c.Query("work_type")
	tripTypeCode := 0
	if workType == "1" {
		tripTypeCode = 1
	} else {
		tripTypeCode = 0
	}
	if masCarpoolUID == "" {
		var carpoolVehicle models.VmsMasCarpoolVehicle
		if err := config.DB.Where("mas_vehicle_uid = ? AND is_deleted = '0' AND is_active = '1'", masVehicleUID).First(&carpoolVehicle).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Vehicle not found", "message": messages.ErrNotfound.Error()})
			return
		}
		masCarpoolUID = carpoolVehicle.MasCarpoolUID
	}
	var count int64
	query := config.DB.Raw(`SELECT count(mas_driver_uid) FROM fn_get_available_drivers_view (?, ?, ?, ?, ?) where mas_carpool_uid = ?`,
		startDate, endDate, bureauDeptSap, businessArea, tripTypeCode, masCarpoolUID)
	if err := query.Scan(&count).Error; err == nil {
		c.JSON(http.StatusOK, gin.H{"number_of_available_drivers": int(count)})
		return
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "message": messages.ErrInternalServer.Error()})
		return
	}

}
