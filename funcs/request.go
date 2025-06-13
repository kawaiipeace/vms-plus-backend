package funcs

import (
	"net/http"
	"sort"
	"strings"
	"time"
	"vms_plus_be/config"
	"vms_plus_be/messages"
	"vms_plus_be/models"
	"vms_plus_be/userhub"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

var StatusNameMap = map[string]string{
	"20": "รออนุมัติ",
	"21": "ถูกตีกลับ",
	"30": "รอตรวจสอบ",
	"31": "ถูกตีกลับ",
	"40": "รออนุมัติ",
	"41": "ถูกตีกลับ",
	"50": "รอรับกุญแจ",
	"51": "รอรับยานพาหนะ",
	"60": "เดินทาง",
	"70": "รอตรวจสอบ",
	"71": "คืนยานพาหนะไม่สำเร็จ",
	"80": "เสร็จสิ้น",
	"90": "ยกเลิกคำขอ",
}

func GetRequestNo(empID string) (string, error) {
	requestNo := empID

	return requestNo, nil
}

func MenuRequests(statusMenuMap map[string]string, query *gorm.DB) ([]models.VmsTrnRequestSummary, error) {
	var summary []models.VmsTrnRequestSummary

	// Group the request counts by statusMenuMap
	groupedSummary := make(map[string]int)
	for key := range statusMenuMap {
		statusCodes := strings.Split(key, ",")
		querySummary := query.Table("vms_trn_request").Session(&gorm.Session{})
		var count int64
		if err := querySummary.Where("vms_trn_request.ref_request_status_code IN ?", statusCodes).Count(&count).Error; err != nil {
			return nil, err
		}
		groupedSummary[key] += int(count)
	}

	// Build the summary from the grouped data
	for key, count := range groupedSummary {
		summary = append(summary, models.VmsTrnRequestSummary{
			RefRequestStatusCode: key,
			RefRequestStatusName: statusMenuMap[key],
			Count:                count,
		})
	}
	// Sort the summary by RefRequestStatusCode
	sort.Slice(summary, func(i, j int) bool {
		return summary[i].RefRequestStatusCode < summary[j].RefRequestStatusCode
	})
	return summary, nil
}

func GetRequest(c *gin.Context, statusNameMap map[string]string) (models.VmsTrnRequestResponse, error) {
	id := c.Param("trn_request_uid")
	var request models.VmsTrnRequestResponse
	trnRequestUID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid TrnRequestUID", "message": messages.ErrInvalidUID.Error()})
		return request, err
	}

	if err := config.DB.
		Select("vms_trn_request.*,k.receiver_personal_id,k.receiver_fullname,k.receiver_dept_sap,"+
			"k.appointment_start appointment_key_handover_start_datetime,k.appointment_end appointment_key_handover_end_datetime,k.appointment_location appointment_key_handover_place,"+
			"k.receiver_dept_name_short,k.receiver_dept_name_full,k.receiver_desk_phone,k.receiver_mobile_phone,k.receiver_position,k.remark receiver_remark").
		Joins("LEFT JOIN vms_trn_vehicle_key_handover k ON k.trn_request_uid = vms_trn_request.trn_request_uid").
		Preload("MasVehicle.RefFuelType").
		Preload("MasVehicle.VehicleDepartment",
			func(db *gorm.DB) *gorm.DB {
				return db.Select("*, fn_get_vehicle_distance_two_months(mas_vehicle_uid, ?) AS vehicle_distance,public.fn_get_oil_station_eng_by_fleetcard(fleet_card_no) as fleet_card_oil_stations", time.Now())
			},
		).
		Preload("RefCostType").
		Preload("MasDriver").
		Preload("MasDriver.DriverLicense.DriverLicenseType").
		Preload("RefRequestStatus").
		Preload("RefTripType").
		Preload("RefCostType").
		First(&request, "vms_trn_request.trn_request_uid = ?", trnRequestUID).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{})
		return request, err
	}
	if request.MasDriver.DriverBirthdate != (time.Time{}) {
		request.MasDriver.Age = request.MasDriver.CalculateAgeInYearsMonths()
	}
	var carpool models.VmsMasCarpoolCarBookingResponse
	if err := config.DB.First(&carpool, "mas_carpool_uid = ?", request.MasCarpoolUID).Error; err == nil {
		request.CarpoolName = carpool.CarpoolName
		if carpool.RefCarpoolChooseCarID == 2 {
			request.IsAdminChooseVehicle = "1"
		}
		if carpool.RefCarpoolChooseDriverID == 2 {
			request.IsAdminChooseDriver = "1"
			var count int64
			vehicleUser, _ := userhub.GetUserInfo(request.VehicleUserEmpID)
			query := config.DB.Raw(`SELECT count(mas_driver_uid) FROM fn_get_available_drivers_view (?, ?, ?, ?, ?) where mas_carpool_uid = ?`,
				request.ReserveStartDatetime, request.ReserveEndDatetime, vehicleUser.BureauDeptSap, vehicleUser.BusinessArea, request.RefTripType.RefTripTypeCode, request.MasCarpoolUID)
			if err := query.Scan(&count).Error; err == nil {
				request.NumberOfAvailableDrivers = int(count)
			}
		}
		if carpool.RefCarpoolChooseDriverID == 2 && request.IsPEAEmployeeDriver == "0" && request.MasCarpoolDriverUID == "" {
			request.CanChooseDriver = true
		}
		if carpool.RefCarpoolChooseCarID == 2 && request.MasVehicleUID == "" {
			request.CanChooseVehicle = true
		}

	}

	request.VehicleUserImageUrl = GetEmpImage(request.VehicleUserEmpID)
	request.ConfirmedRequestImageUrl = GetEmpImage(request.ConfirmedRequestEmpID)
	request.DriverEmpImageUrl = GetEmpImage(request.DriverEmpID)

	request.DriverImageURL = config.DefaultAvatarURL
	request.CanCancelRequest = true
	request.IsUseDriver = request.MasCarpoolDriverUID != ""
	request.RefRequestStatusName = StatusNameMap[request.RefRequestStatusCode]

	request.VehicleLicensePlate = request.MasVehicle.VehicleLicensePlate
	request.VehicleLicensePlateProvinceShort = request.MasVehicle.VehicleLicensePlateProvinceShort
	request.VehicleLicensePlateProvinceFull = request.MasVehicle.VehicleLicensePlateProvinceFull
	request.MasVehicle.VehicleDepartment.VehicleOwnerDeptShort = GetDeptSAPShort(request.MasVehicle.VehicleDepartment.VehicleOwnerDeptSap)

	var vehicleImgs []models.VmsMasVehicleImg
	if err := config.DB.Where("mas_vehicle_uid = ?", request.MasVehicle.MasVehicleUID).Find(&vehicleImgs).Error; err == nil {
		request.MasVehicle.VehicleImgs = make([]string, 0)
		for _, img := range vehicleImgs {
			request.MasVehicle.VehicleImgs = append(request.MasVehicle.VehicleImgs, img.VehicleImgFile)
		}
	}

	if request.RefRequestStatusCode == "90" {
		// Check VmsLogRequest
		var logRequest models.VmsLogRequest
		if err := config.DB.
			Where("trn_request_uid = ? AND is_deleted = '0' AND ref_request_status_code = ?", request.TrnRequestUID, request.RefRequestStatusCode).
			First(&logRequest).Error; err == nil {
			request.CanceledRequestRole = logRequest.ActionByRole
		}
	}

	return request, nil
}

func GetRequestVehicelInUse(c *gin.Context, statusNameMap map[string]string) (models.VmsTrnRequestVehicleInUseResponse, error) {
	id := c.Param("trn_request_uid")
	var request models.VmsTrnRequestVehicleInUseResponse
	trnRequestUID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid TrnRequestUID", "message": messages.ErrInvalidUID.Error()})
		return request, err
	}

	if err := config.DB.
		Preload("MasVehicle.RefFuelType").
		Preload("MasVehicle.VehicleDepartment").
		Preload("RefCostType").
		Preload("MasDriver").
		Preload("RefRequestStatus").
		Preload("RequestVehicleType").
		Preload("VehicleImagesReceived").
		Preload("VehicleImagesReturned").
		Preload("VehicleImageInspect").
		Preload("ReceiverKeyTypeDetail").
		Preload("RefTripType").
		Preload("SatisfactionSurveyAnswers.SatisfactionSurveyQuestions").
		Select("vms_trn_request.*,k.ref_vehicle_key_type_code ref_vehicle_key_type_code,k.receiver_personal_id,k.receiver_fullname,k.receiver_dept_sap,"+
			"k.appointment_start appointment_key_handover_start_datetime,k.appointment_end appointment_key_handover_end_datetime,k.appointment_location appointment_key_handover_place,"+
			"k.receiver_dept_name_short,k.receiver_dept_name_full,k.receiver_desk_phone,k.receiver_mobile_phone,k.receiver_position,k.remark receiver_remark").
		Joins("LEFT JOIN vms_trn_vehicle_key_handover k ON k.trn_request_uid = vms_trn_request.trn_request_uid").
		First(&request, "vms_trn_request.trn_request_uid = ?", trnRequestUID).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{})
		return request, err
	}
	if request.MasDriver.DriverBirthdate != (time.Time{}) {
		request.MasDriver.Age = request.MasDriver.CalculateAgeInYearsMonths()
	}
	request.ParkingPlace = request.MasVehicle.VehicleDepartment.ParkingPlace
	request.DriverImageURL = config.DefaultAvatarURL
	request.ReceivedKeyImageURL = config.DefaultAvatarURL
	request.CanCancelRequest = true
	request.IsUseDriver = request.MasCarpoolDriverUID != ""
	request.RefRequestStatusName = StatusNameMap[request.RefRequestStatusCode]
	request.FleetCardNo = request.MasVehicle.VehicleDepartment.FleetCardNo
	request.VehicleLicensePlate = request.MasVehicle.VehicleLicensePlate
	request.VehicleLicensePlateProvinceShort = request.MasVehicle.VehicleLicensePlateProvinceShort
	request.VehicleLicensePlateProvinceFull = request.MasVehicle.VehicleLicensePlateProvinceFull

	// Get vehicle department details
	if err := config.DB.Where("mas_vehicle_uid = ?", request.MasVehicle.MasVehicleUID).
		Select("*,public.fn_get_oil_station_eng_by_fleetcard(fleet_card_no) as fleet_card_oil_stations").
		Where("is_deleted = '0' AND is_active = '1'").
		First(&request.MasVehicle.VehicleDepartment).Error; err != nil {

	}
	request.MasVehicle.VehicleDepartment.VehicleOwnerDeptShort = GetDeptSAPShort(request.MasVehicle.VehicleDepartment.VehicleOwnerDeptSap)
	//check if vehicle is carpool
	var carpoolVehicle models.VmsMasCarpoolVehicle
	masCarpoolUID := ""
	if err := config.DB.Where("mas_vehicle_uid = ? AND is_deleted = '0' AND is_active = '1'", request.MasVehicle.MasVehicleUID).First(&carpoolVehicle).Error; err == nil {
		masCarpoolUID = carpoolVehicle.MasCarpoolUID
	}
	if masCarpoolUID != "" {
		var carpoolAdmin models.VmsMasCarpoolAdmin
		if err := config.DB.Where("mas_carpool_uid = ? AND is_deleted = '0' AND is_active = '1'", masCarpoolUID).
			Select("admin_emp_no,admin_emp_name,admin_dept_sap,admin_position,mobile_contact_number,internal_contact_number").
			Order("is_main_admin DESC").
			First(&carpoolAdmin).Error; err == nil {
			request.MasVehicle.VehicleDepartment.VehicleUser.EmpID = carpoolAdmin.AdminEmpNo
			request.MasVehicle.VehicleDepartment.VehicleUser.FullName = carpoolAdmin.AdminEmpName
			request.MasVehicle.VehicleDepartment.VehicleUser.DeptSAP = carpoolAdmin.AdminDeptSap
			request.MasVehicle.VehicleDepartment.VehicleUser.DeptSAPFull = GetDeptSAPFull(carpoolAdmin.AdminDeptSap)
			request.MasVehicle.VehicleDepartment.VehicleUser.DeptSAPShort = GetDeptSAPShort(carpoolAdmin.AdminDeptSap)
			request.MasVehicle.VehicleDepartment.VehicleUser.ImageUrl = GetEmpImage(carpoolAdmin.AdminEmpNo)
			request.MasVehicle.VehicleDepartment.VehicleUser.Position = carpoolAdmin.AdminPosition
			request.MasVehicle.VehicleDepartment.VehicleUser.TelMobile = carpoolAdmin.MobileContactNumber
			request.MasVehicle.VehicleDepartment.VehicleUser.TelInternal = carpoolAdmin.InternalContactNumber
			request.MasVehicle.VehicleDepartment.VehicleUser.IsEmployee = true
		}
	}

	request.MileUsed = request.MileEnd - request.MileStart
	if err := config.DB.
		Table("vms_trn_add_fuel").
		Where("trn_request_uid = ? AND is_deleted = '0'", request.TrnRequestUID).
		Count(&request.AddFuelsCount).Error; err != nil {
		request.AddFuelsCount = 0
	}
	if err := config.DB.
		Table("vms_trn_trip_detail").
		Where("trn_request_uid = ? AND is_deleted = '0'", request.TrnRequestUID).
		Count(&request.TripDetailsCount).Error; err != nil {
		request.TripDetailsCount = 0
	}
	request.IsReturnOverDue = false
	if time.Now().Truncate(24 * time.Hour).After(request.ReserveEndDatetime.Truncate(24 * time.Hour)) {
		request.IsReturnOverDue = true
	}

	//c.JSON(http.StatusOK, request)
	return request, nil
}

func GetAdminApprovalEmpIDs(trnRequestUID string) ([]string, error) {
	var result struct {
		MasCarpoolUID string
		MasVehicleUID string
	}
	if err := config.DB.Table("vms_trn_request").
		Select("mas_carpool_uid, mas_vehicle_uid").
		Where("trn_request_uid = ?", trnRequestUID).
		Scan(&result).Error; err != nil {
		return nil, err
	}

	var empIDs []string
	if result.MasCarpoolUID != "" && result.MasCarpoolUID != DefaultUUID() {
		if err := config.DB.Table("vms_mas_carpool_admin").
			Select("admin_emp_no").
			Where("mas_carpool_uid = ? AND is_deleted = '0' AND is_active = '1'", result.MasCarpoolUID).
			Pluck("admin_emp_no", &empIDs).Error; err != nil {
			return nil, err
		}
		return empIDs, nil
	} else {
		var bureauDeptSap string
		request := userhub.ServiceListUserRequest{
			ServiceCode:   "vms",
			Search:        "",
			Role:          "admin_approval",
			BureauDeptSap: bureauDeptSap,
			Limit:         100,
		}
		lists, err := userhub.GetUserList(request)
		if err != nil {
			return nil, err
		}
		for _, list := range lists {
			empIDs = append(empIDs, list.EmpID)
		}
		return empIDs, nil
	}
}

func GetFinalApprovalEmpIDs(trnRequestUID string) ([]string, error) {
	var result struct {
		MasCarpoolUID string
		MasVehicleUID string
	}
	if err := config.DB.Table("vms_trn_request").
		Select("mas_carpool_uid, mas_vehicle_uid").
		Where("trn_request_uid = ?", trnRequestUID).
		Scan(&result).Error; err != nil {
		return nil, err
	}

	var empIDs []string
	if result.MasCarpoolUID != "" && result.MasCarpoolUID != DefaultUUID() {
		if err := config.DB.Table("vms_mas_carpool_approver").
			Select("approver_emp_no").
			Where("mas_carpool_uid = ? AND is_deleted = '0' AND is_active = '1'", result.MasCarpoolUID).
			Pluck("approver_emp_no", &empIDs).Error; err != nil {
			return nil, err
		}
		return empIDs, nil
	} else {
		var bureauDeptSap string
		request := userhub.ServiceListUserRequest{
			ServiceCode:   "vms",
			Search:        "",
			Role:          "final_approval",
			BureauDeptSap: bureauDeptSap,
			Limit:         100,
		}
		lists, err := userhub.GetUserList(request)
		if err != nil {
			return nil, err
		}
		for _, list := range lists {
			empIDs = append(empIDs, list.EmpID)
		}
		return empIDs, nil
	}
}

func GetProgressRequestStatusEmp(trnRequestUID, refRequestStatusCode, actionRoleName string) models.ProgressRequestStatusEmp {
	var logRequest models.VmsLogRequest
	if err := config.DB.
		Where("trn_request_uid = ? AND ref_request_status_code = ? AND is_deleted = '0'",
			trnRequestUID, refRequestStatusCode).
		First(&logRequest).Error; err != nil {
		return models.ProgressRequestStatusEmp{}
	}
	empUser := GetUserEmpInfo(logRequest.ActionByPersonalID)
	progressRequestStatusEmp := models.ProgressRequestStatusEmp{
		ActionRole:   actionRoleName,
		EmpID:        empUser.EmpID,
		EmpName:      empUser.FullName,
		EmpPosition:  empUser.Position,
		DeptSAP:      empUser.DeptSAP,
		DeptSAPShort: empUser.DeptSAPShort,
		DeptSAPFull:  empUser.DeptSAPFull,
		PhoneNumber:  empUser.TelInternal,
		MobileNumber: empUser.TelMobile,
	}

	return progressRequestStatusEmp
}
