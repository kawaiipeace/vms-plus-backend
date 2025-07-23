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

	//replace vehicle_owner_dept_short with carpool_name if in carpool
	vehicleCarpoolName := ""
	if err := config.DB.Table("vms_mas_carpool mc").
		Joins("INNER JOIN vms_mas_carpool_vehicle mcv ON mcv.mas_carpool_uid = mc.mas_carpool_uid AND mcv.mas_vehicle_uid = ? AND mcv.is_deleted = '0' AND mcv.is_active = '1'", request.MasVehicle.MasVehicleUID).
		Select("mc.carpool_name").
		Scan(&vehicleCarpoolName).Error; err == nil {
		if vehicleCarpoolName != "" {
			request.MasVehicle.VehicleDepartment.VehicleOwnerDeptShort = vehicleCarpoolName
			request.VehicleDepartmentDeptSapShort = vehicleCarpoolName
			request.VehicleDepartmentDeptSapFull = vehicleCarpoolName
		}
	}
	//replace vehicle_owner_dept_short with carpool_name if in carpool
	driverCarpoolName := ""
	if err := config.DB.Table("vms_mas_carpool mc").
		Joins("INNER JOIN vms_mas_carpool_driver mcv ON mcv.mas_carpool_uid = mc.mas_carpool_uid AND mcv.mas_driver_uid = ? AND mcv.is_deleted = '0' AND mcv.is_active = '1'", request.MasDriver.MasDriverUID).
		Select("mc.carpool_name").
		Scan(&driverCarpoolName).Error; err == nil {
		if driverCarpoolName != "" {
			request.MasDriver.VendorName = driverCarpoolName
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

	if request.MasCarpoolUID != "" {
		var carpoolAdmin models.VmsMasCarpoolAdmin
		if err := config.DB.Where("mas_carpool_uid = ? AND is_deleted = '0' AND is_active = '1'", request.MasCarpoolUID).
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

	//VmsTrnSatisfactionSurveyAnswersResponse

	var satisfactionSurveyAnswers []models.VmsTrnSatisfactionSurveyAnswersResponse
	if err := config.DB.Table("vms_trn_satisfaction_survey_answers").
		Where("trn_request_uid = ?", request.TrnRequestUID).
		Find(&satisfactionSurveyAnswers).Error; err == nil {
		for i := range request.SatisfactionSurveyAnswers {
			for j := range satisfactionSurveyAnswers {
				if request.SatisfactionSurveyAnswers[i].MasSatisfactionSurveyQuestionsUID == satisfactionSurveyAnswers[j].MasSatisfactionSurveyQuestionsUID {
					request.SatisfactionSurveyAnswers[i].SurveyAnswer = satisfactionSurveyAnswers[j].SurveyAnswer
				}
			}

		}
	}
	//replace vehicle_owner_dept_short with carpool_name if in carpool
	vehicleCarpoolName := ""
	if err := config.DB.Table("vms_mas_carpool mc").
		Joins("INNER JOIN vms_mas_carpool_vehicle mcv ON mcv.mas_carpool_uid = mc.mas_carpool_uid AND mcv.mas_vehicle_uid = ? AND mcv.is_deleted = '0' AND mcv.is_active = '1'", request.MasVehicle.MasVehicleUID).
		Select("mc.carpool_name").
		Scan(&vehicleCarpoolName).Error; err == nil {
		if vehicleCarpoolName != "" {
			request.MasVehicle.VehicleDepartment.VehicleOwnerDeptShort = vehicleCarpoolName
		}
	}
	//replace vehicle_owner_dept_short with carpool_name if in carpool
	driverCarpoolName := ""
	request.MasDriver.VendorName = request.MasDriver.DriverDeptSAPShort
	if err := config.DB.Table("vms_mas_carpool mc").
		Joins("INNER JOIN vms_mas_carpool_driver mcv ON mcv.mas_carpool_uid = mc.mas_carpool_uid AND mcv.mas_driver_uid = ? AND mcv.is_deleted = '0' AND mcv.is_active = '1'", request.MasDriver.MasDriverUID).
		Select("mc.carpool_name").
		Scan(&driverCarpoolName).Error; err == nil {
		if driverCarpoolName != "" {
			request.MasDriver.VendorName = driverCarpoolName
		}
	}
	request.CanScoreButton = IsAllowScoreButton(request.TrnRequestUID)
	if request.CanScoreButton {
		request.CanPickupButton = false
	} else {
		request.CanPickupButton = IsAllowPickupButton(request.TrnRequestUID)
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

func UpdateVehicleMileage(trnRequestUID string, mileage int) error {
	var masVehicleUID string
	if err := config.DB.Table("vms_trn_request").
		Where("trn_request_uid = ?", trnRequestUID).
		Select("mas_vehicle_uid").
		Scan(&masVehicleUID).Error; err != nil {
		return err
	}

	if err := config.DB.Table("vms_mas_vehicle_department").
		Where("mas_vehicle_uid = ?", masVehicleUID).
		Update("vehicle_mileage", mileage).Error; err != nil {
		return err
	}

	return nil
}

func UpdateVehicleParkingPlace(trnRequestUID string, parkingPlace string) error {
	var masVehicleUID string
	if err := config.DB.Table("vms_trn_request").
		Where("trn_request_uid = ?", trnRequestUID).
		Select("mas_vehicle_uid").
		Scan(&masVehicleUID).Error; err != nil {
		return err
	}

	if err := config.DB.Table("vms_mas_vehicle_department").
		Where("mas_vehicle_uid = ?", masVehicleUID).
		Update("parking_place", parkingPlace).Error; err != nil {
		return err
	}

	return nil
}
func CheckMustPassStatus30Department(trnRequestUID string) {
	var exists bool
	err := config.DB.
		Table("vms_trn_request").
		Select("1").
		Where(`
			vms_trn_request.ref_request_status_code = '20' AND
			vms_trn_request.mas_carpool_uid is null AND
			vms_trn_request.trn_request_uid = ?`, trnRequestUID).
		Limit(1).
		Scan(&exists).Error
	if err != nil {
		return
	} else if exists {
		//update vms_trn_request set ref_request_status_code='30'
		if err := config.DB.Table("vms_trn_request").
			Where("trn_request_uid = ?", trnRequestUID).
			Update("ref_request_status_code", "30").Error; err != nil {
			return
		}
	}
}

func CheckMustPassStatus30(trnRequestUID string) {
	CheckMustPassStatus30Department(trnRequestUID)

	var exists bool
	err := config.DB.
		Table("vms_mas_carpool").
		Select("1").
		Joins("INNER JOIN vms_trn_request ON vms_trn_request.mas_carpool_uid = vms_mas_carpool.mas_carpool_uid").
		Where(`
        vms_mas_carpool.is_must_pass_status_30 = '1' AND
        vms_trn_request.ref_request_status_code = '20' AND
        vms_trn_request.trn_request_uid = ?
    `, trnRequestUID).
		Limit(1).
		Scan(&exists).Error

	if err != nil {
		return
	} else if exists {
		//update vms_trn_request set ref_request_status_code='30'
		if err := config.DB.Table("vms_trn_request").
			Where("trn_request_uid = ?", trnRequestUID).
			Update("ref_request_status_code", "30").Error; err != nil {
			return
		}
	}
}

func CheckMustPassStatus40(trnRequestUID string) {
	var exists bool
	err := config.DB.
		Table("vms_mas_carpool").
		Select("1").
		Joins("INNER JOIN vms_trn_request ON vms_trn_request.mas_carpool_uid = vms_mas_carpool.mas_carpool_uid").
		Where(`
        vms_mas_carpool.is_must_pass_status_40 = '1' AND
        vms_trn_request.ref_request_status_code = '30' AND
        vms_trn_request.trn_request_uid = ?
    `, trnRequestUID).
		Limit(1).
		Scan(&exists).Error

	if err != nil {
		return
	} else if exists {
		//update vms_trn_request set ref_request_status_code='40'
		if err := config.DB.Table("vms_trn_request").
			Where("trn_request_uid = ?", trnRequestUID).
			Update("ref_request_status_code", "40").Error; err != nil {
			return
		}
		SetReceivedKey(trnRequestUID, "")
	}
}

func CheckMustPassStatus50(trnRequestUID string) {
	var exists bool
	err := config.DB.
		Table("vms_mas_carpool").
		Select("1").
		Joins("INNER JOIN vms_trn_request ON vms_trn_request.mas_carpool_uid = vms_mas_carpool.mas_carpool_uid").
		Where(`
        vms_mas_carpool.is_must_pass_status_50 = '1' AND
        vms_trn_request.ref_request_status_code = '40' AND
        vms_trn_request.trn_request_uid = ?
    `, trnRequestUID).
		Limit(1).
		Scan(&exists).Error

	if err != nil {
		return
	} else if exists {
		//update vms_trn_request set ref_request_status_code='50'
		if err := config.DB.Table("vms_trn_request").
			Where("trn_request_uid = ?", trnRequestUID).
			Update("ref_request_status_code", "50").Error; err != nil {
			return
		}
	}
}

func CheckMustPassStatus(trnRequestUID string) {
	CheckMustPassStatus30(trnRequestUID)
	CheckMustPassStatus40(trnRequestUID)
	CheckMustPassStatus50(trnRequestUID)
}

func IsAllowPickupButton(trnRequestUID string) bool {
	var exists bool
	err := config.DB.
		Table("vms_trn_request").
		Select("1").
		Where("trn_request_uid = ? AND ref_request_status_code < '60'", trnRequestUID).
		Limit(1).
		Scan(&exists).Error

	if err != nil {
		return false
	}
	if exists {
		return true
	}
	return false
}

func IsAllowScoreButton(trnRequestUID string) bool {
	var exists bool
	err := config.DB.
		Table("vms_trn_request").
		Select("1").
		Where("trn_request_uid = ? AND ref_request_status_code >= '60' AND ref_request_status_code < '70' AND date(reserve_end_datetime) >= date(?)", trnRequestUID, time.Now()).
		Limit(1).
		Scan(&exists).Error

	if err != nil {
		return false
	}
	if exists {
		return true
	}
	return false
}

// Test
func SetReceivedKey(trnRequestUID string, handoverUID string) {
	if handoverUID == "" {
		handoverUID = uuid.New().String()
	}
	request := models.VmsTrnRequestApprovedWithRecieiveKey{
		HandoverUID:              handoverUID,
		TrnRequestUID:            trnRequestUID,
		ReceiverType:             0,
		CreatedBy:                "system",
		CreatedAt:                time.Now(),
		UpdatedBy:                "system",
		UpdatedAt:                time.Now(),
		ReceivedKeyStartDatetime: models.TimeWithZone{Time: time.Now()},
		ReceivedKeyEndDatetime:   models.TimeWithZone{Time: time.Now()},
		ReceivedKeyPlace:         "-",
	}
	var requestDetail struct {
		MasCarpoolUID        string
		ReserveEndDatetime   time.Time
		ReserveStartDatetime time.Time
	}

	if err := config.DB.Table("vms_trn_request").
		Where("trn_request_uid = ?", trnRequestUID).Select("mas_carpool_uid, reserve_end_datetime, reserve_start_datetime").Scan(&requestDetail).Error; err != nil {
		return
	}
	//ReceivedKeyPlace = carpool_contact_place
	var carpoolContactPlace string
	if err := config.DB.Table("vms_mas_carpool").
		Select("carpool_contact_place").
		Where("mas_carpool_uid = ?", requestDetail.MasCarpoolUID).
		Scan(&carpoolContactPlace).Error; err == nil {
		request.ReceivedKeyPlace = carpoolContactPlace
	}
	if requestDetail.ReserveStartDatetime.Hour() >= 12 {
		date := requestDetail.ReserveStartDatetime.Truncate(24 * time.Hour)
		// convert to 8:00 at Bangkok
		bangkokLoc, err := time.LoadLocation("Asia/Bangkok")
		if err != nil {
			bangkokLoc = time.UTC // fallback to UTC if Bangkok location fails to load
		}
		date_8_00 := time.Date(date.Year(), date.Month(), date.Day(), 8, 0, 0, 0, bangkokLoc)
		date_12_00 := time.Date(date.Year(), date.Month(), date.Day(), 12, 0, 0, 0, bangkokLoc)
		request.ReceivedKeyStartDatetime = models.TimeWithZone{Time: date_8_00}
		request.ReceivedKeyEndDatetime = models.TimeWithZone{Time: date_12_00}
	} else {
		date := requestDetail.ReserveStartDatetime.Truncate(24 * time.Hour)
		var holidays []models.VmsMasHolidays
		if err := config.DB.Table("vms_mas_holidays").
			Select("mas_holidays_date").
			Find(&holidays).Error; err != nil {
			return
		}
		//find yesterday with not sunday,saturday,holiday
		yesterday := date.AddDate(0, 0, -1)
		for IsHoliday(yesterday, holidays) {
			yesterday = yesterday.AddDate(0, 0, -1)
		}

		//settime yesterday to 8:00:00
		bangkokLoc, err := time.LoadLocation("Asia/Bangkok")
		if err != nil {
			bangkokLoc = time.UTC // fallback to UTC if Bangkok location fails to load
		}
		yesterday_8_00 := time.Date(yesterday.Year(), yesterday.Month(), yesterday.Day(), 8, 0, 0, 0, bangkokLoc)
		yesterday_12_00 := time.Date(yesterday.Year(), yesterday.Month(), yesterday.Day(), 12, 0, 0, 0, bangkokLoc)
		request.ReceivedKeyStartDatetime = models.TimeWithZone{Time: yesterday_8_00}
		request.ReceivedKeyEndDatetime = models.TimeWithZone{Time: yesterday_12_00}
	}
	if err := config.DB.Save(&request).Error; err != nil {
		return
	}

	//update vms_trn_request set appointment_key_handover_place,appointment_key_handover_start_datetime,appointment_key_handover_end_datetime
	if err := config.DB.Table("vms_trn_request").
		Where("trn_request_uid = ?", trnRequestUID).
		Update("appointment_key_handover_place", request.ReceivedKeyPlace).
		Update("appointment_key_handover_start_datetime", request.ReceivedKeyStartDatetime).
		Update("appointment_key_handover_end_datetime", request.ReceivedKeyEndDatetime).Error; err != nil {
		return
	}
}
