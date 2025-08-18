package funcs

import (
	"encoding/csv"
	"fmt"
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
	"github.com/tealeg/xlsx"
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
		Preload("MasDriver.DriverStatus").
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

	if request.MasDriver.DriverImage != "" {
		request.DriverImageURL = request.MasDriver.DriverImage
	} else {
		request.DriverImageURL = GetEmpImage(request.DriverEmpID)
	}
	request.CanCancelRequest = true
	request.IsUseDriver = request.MasCarpoolDriverUID != ""
	request.RefRequestStatusName = StatusNameMap[request.RefRequestStatusCode]

	request.VehicleLicensePlate = request.MasVehicle.VehicleLicensePlate
	request.VehicleLicensePlateProvinceShort = request.MasVehicle.VehicleLicensePlateProvinceShort
	request.VehicleLicensePlateProvinceFull = request.MasVehicle.VehicleLicensePlateProvinceFull
	request.MasVehicle.VehicleDepartment.VehicleOwnerDeptShort = GetDeptSAPShort(request.MasVehicle.VehicleDepartment.VehicleOwnerDeptSap)
	request.MasVehicle.Age = CalculateAgeInt(request.MasVehicle.VehicleRegistrationDate)
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
	request.MasDriver.VendorName = request.MasDriver.DriverDeptSAPShort
	if err := config.DB.Table("vms_mas_carpool mc").
		Joins("INNER JOIN vms_mas_carpool_driver mcv ON mcv.mas_carpool_uid = mc.mas_carpool_uid AND mcv.mas_driver_uid = ? AND mcv.is_deleted = '0' AND mcv.is_active = '1'", request.MasDriver.MasDriverUID).
		Select("mc.carpool_name").
		Scan(&driverCarpoolName).Error; err == nil {
		if driverCarpoolName != "" {
			request.MasDriver.VendorName = driverCarpoolName
		}
	}

	if vehicleCarpoolName != "" {
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
	} else {
		userList, err := userhub.GetUserList(userhub.ServiceListUserRequest{
			ServiceCode: "vms",
			Role:        "admin-department-main",
			DeptSaps:    request.MasVehicle.VehicleDepartment.BureauDeptSap,
			Limit:       100,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "message": messages.ErrInternalServer.Error()})
			return request, err
		}
		if len(userList) > 0 {
			request.MasVehicle.VehicleDepartment.VehicleUser = userList[0]
		} else {
			userList, err := userhub.GetUserList(userhub.ServiceListUserRequest{
				ServiceCode: "vms",
				Role:        "admin-department",
				DeptSaps:    request.MasVehicle.VehicleDepartment.BureauDeptSap,
				Limit:       100,
			})
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "message": messages.ErrInternalServer.Error()})
				return request, err
			}
			if len(userList) > 0 {
				request.MasVehicle.VehicleDepartment.VehicleUser = userList[0]
			}
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
		Preload("MasDriver.DriverLicense.DriverLicenseType").
		Preload("MasDriver.DriverStatus").
		Preload("RefRequestStatus").
		Preload("RequestVehicleType").
		Preload("VehicleImagesReceived").
		Preload("VehicleImagesReturned").
		Preload("VehicleImageInspect").
		Preload("ReceiverKeyTypeDetail").
		Preload("RefTripType").
		Preload("SatisfactionSurveyAnswers.SatisfactionSurveyQuestions").
		Select("vms_trn_request.*,k.receiver_type,k.ref_vehicle_key_type_code ref_vehicle_key_type_code,k.receiver_personal_id,k.receiver_fullname,k.receiver_dept_sap,"+
			"k.appointment_start appointment_key_handover_start_datetime,k.appointment_end appointment_key_handover_end_datetime,k.appointment_location appointment_key_handover_place,"+
			"k.receiver_dept_name_short,k.receiver_dept_name_full,k.receiver_desk_phone,k.receiver_mobile_phone,k.receiver_position,k.remark receiver_remark,k.actual_receive_time received_key_datetime").
		Joins("LEFT JOIN vms_trn_vehicle_key_handover k ON k.trn_request_uid = vms_trn_request.trn_request_uid").
		First(&request, "vms_trn_request.trn_request_uid = ?", trnRequestUID).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{})
		return request, err
	}
	if request.MasDriver.DriverBirthdate != (time.Time{}) {
		request.MasDriver.Age = request.MasDriver.CalculateAgeInYearsMonths()
	}
	request.ParkingPlace = request.MasVehicle.VehicleDepartment.ParkingPlace

	request.VehicleUserImageUrl = GetEmpImage(request.VehicleUserEmpID)

	if request.MasDriver.DriverImage != "" {
		request.DriverImageURL = request.MasDriver.DriverImage
	} else {
		request.DriverImageURL = GetEmpImage(request.DriverEmpID)
	}
	switch request.ReceiverKeyType {
	case 1:
		request.ReceivedKeyImageURL = request.MasDriver.DriverImage
	case 2:
		request.ReceivedKeyImageURL = GetEmpImage(request.ReceivedKeyEmpID)
	default:
		request.ReceivedKeyImageURL = config.DefaultAvatarURL
	}

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
	request.MasVehicle.Age = CalculateAgeInt(request.MasVehicle.VehicleRegistrationDate)
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
	} else {
		userList, err := userhub.GetUserList(userhub.ServiceListUserRequest{
			ServiceCode: "vms",
			Role:        "admin-department-main",
			DeptSaps:    request.MasVehicle.VehicleDepartment.BureauDeptSap,
			Limit:       100,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "message": messages.ErrInternalServer.Error()})
			return request, err
		}
		if len(userList) > 0 {
			request.MasVehicle.VehicleDepartment.VehicleUser = userList[0]
		} else {
			userList, err := userhub.GetUserList(userhub.ServiceListUserRequest{
				ServiceCode: "vms",
				Role:        "admin-department",
				DeptSaps:    request.MasVehicle.VehicleDepartment.BureauDeptSap,
				Limit:       100,
			})
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "message": messages.ErrInternalServer.Error()})
				return request, err
			}
			if len(userList) > 0 {
				request.MasVehicle.VehicleDepartment.VehicleUser = userList[0]
			}
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

	if time.Now().After(request.ReserveEndDatetime.Time) {
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
	if IsAllowScoreButton(request.TrnRequestUID) {
		if request.TripDetailsCount > 0 {
			request.CanScoreButton = true
			request.CanPickupButton = false
			request.CanTravelCardButton = false
		} else {
			request.CanScoreButton = false
			request.CanPickupButton = false
			request.CanTravelCardButton = true
		}
	} else if IsAllowPickupButton(request.TrnRequestUID) {
		request.CanScoreButton = false
		request.CanPickupButton = true
		request.CanTravelCardButton = false
	} else {
		request.CanScoreButton = false
		request.CanPickupButton = false
		request.CanTravelCardButton = false
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

		if result.MasVehicleUID != "" && result.MasVehicleUID != DefaultUUID() {
			if err := config.DB.Table("vms_mas_vehicle_department").
				Select("bureau_dept_sap").
				Where("mas_vehicle_uid = ? AND is_deleted = '0' AND is_active = '1'", result.MasVehicleUID).
				Scan(&bureauDeptSap).Error; err != nil {
				return nil, err
			}
		}
		request := userhub.ServiceListUserRequest{
			ServiceCode:   "vms",
			Search:        "",
			Role:          "admin-department-main",
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
			Order("is_main_approver DESC").
			Pluck("approver_emp_no", &empIDs).Error; err != nil {
			return nil, err
		}
		return empIDs, nil
	} else {
		var approvedRequestEmpID string
		if err := config.DB.Table("vms_trn_request").
			Select("approved_request_emp_id").
			Where("trn_request_uid = ?", trnRequestUID).
			Scan(&approvedRequestEmpID).Error; err != nil {
			return nil, err
		}
		fmt.Println("approvedRequestEmpID", approvedRequestEmpID)
		empIDs = append(empIDs, approvedRequestEmpID)
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
		if err := config.DB.Table("vms_trn_request").
			Where("trn_request_uid = ?", trnRequestUID).
			Update("ref_request_status_code", "30").Error; err != nil {
			return
		}
		var confirmedRequestEmpID string
		if err := config.DB.Table("vms_trn_request").
			Where("trn_request_uid = ?", trnRequestUID).
			Select("confirmed_request_emp_id").
			Scan(&confirmedRequestEmpID).Error; err != nil {
			return
		}
		/*CreateTrnRequestActionLog(trnRequestUID,
			"30",
			"รอผู้ดูแลยานพาหนะตรวจสอบ",
			confirmedRequestEmpID,
			"level1-approval",
			"",
		)*/
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

		var confirmedRequestEmpID string
		if err := config.DB.Table("vms_trn_request").
			Where("trn_request_uid = ?", trnRequestUID).
			Select("confirmed_request_emp_id").
			Scan(&confirmedRequestEmpID).Error; err != nil {
			return
		}
		CreateTrnRequestActionLog(trnRequestUID,
			"30",
			"รอผู้ดูแลยานพาหนะตรวจสอบ",
			confirmedRequestEmpID,
			"level1-approval",
			"",
		)
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
		var adminEmpNo string
		if err := config.DB.Table("vms_trn_request").
			Joins("INNER JOIN vms_mas_carpool_admin ON vms_mas_carpool_admin.mas_carpool_uid = vms_trn_request.mas_carpool_uid AND vms_mas_carpool_admin.is_deleted = '0' AND vms_mas_carpool_admin.is_active = '1' AND is_main_admin = '1'").
			Where("trn_request_uid = ?", trnRequestUID).
			Select("admin_emp_no").
			Scan(&adminEmpNo).Error; err != nil {
			return
		}
		CreateTrnRequestActionLog(trnRequestUID,
			"40",
			"รออนุมัติ จากเจ้าของยานพาหนะ",
			adminEmpNo,
			"admin-department",
			"",
		)
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
		approvedEmpID := UpdateApproverRequest(trnRequestUID)
		UpdateRecievedKeyUser(trnRequestUID)

		var receivedKey models.VmsTrnRequestApprovedWithRecieiveKey
		if err := config.DB.First(&receivedKey, "trn_request_uid = ?", trnRequestUID).Error; err != nil {
			return
		}
		CreateTrnRequestActionLog(trnRequestUID,
			"50",
			GetDateBuddhistYear(receivedKey.ReceivedKeyStartDatetime.Time)+" สถานที่ "+receivedKey.ReceivedKeyPlace+" นัดหมายรับกุญแจ",
			approvedEmpID,
			"approval-department",
			"",
		)
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

		date_8_00 := time.Date(date.Year(), date.Month(), date.Day(), 1, 0, 0, 0, time.UTC)
		date_12_00 := time.Date(date.Year(), date.Month(), date.Day(), 5, 0, 0, 0, time.UTC)
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

		yesterday_8_00 := time.Date(yesterday.Year(), yesterday.Month(), yesterday.Day(), 1, 0, 0, 0, time.UTC)
		yesterday_12_00 := time.Date(yesterday.Year(), yesterday.Month(), yesterday.Day(), 5, 0, 0, 0, time.UTC)
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

func UpdateApproverRequest(trnRequestUID string) string {
	empIDs, err := GetFinalApprovalEmpIDs(trnRequestUID)
	if err != nil {
		return ""
	}
	if len(empIDs) > 0 {
		empUser := GetUserEmpInfo(empIDs[0])
		var request models.VmsTrnRequestApproved
		request.TrnRequestUID = trnRequestUID
		request.ApprovedRequestEmpID = empUser.EmpID
		request.ApprovedRequestEmpName = empUser.FullName
		request.ApprovedRequestDeptSAP = empUser.DeptSAP
		request.ApprovedRequestDeptNameShort = empUser.DeptSAPShort
		request.ApprovedRequestDeptNameFull = empUser.DeptSAPFull
		request.ApprovedRequestDeskPhone = empUser.TelInternal
		request.ApprovedRequestMobilePhone = empUser.TelMobile
		request.ApprovedRequestPosition = empUser.Position
		request.ApprovedRequestDatetime = models.TimeWithZone{Time: time.Now()}
		request.UpdatedAt = time.Now()
		request.UpdatedBy = "system"
		request.RefRequestStatusCode = "50"

		if err := config.DB.Save(&request).Error; err != nil {
			return ""
		}
		return empIDs[0]
	}
	return ""
}

func UpdateRecievedKeyUser(trnRequestUID string) {
	var trnRequest models.VmsTrnRequestResponse
	if err := config.DB.First(&trnRequest, "trn_request_uid = ?", trnRequestUID).Error; err != nil {
		return
	}
	var request = models.VmsTrnReceivedKeyPEA{}
	request.TrnRequestUID = trnRequestUID
	if trnRequest.DriverEmpID[:1] == "D" {
		request.ReceiverType = 1 // Driver
		request.ReceiverPersonalId = trnRequest.DriverEmpID
		request.ReceiverFullname = trnRequest.DriverEmpName
		request.ReceiverDeptSAP = trnRequest.DriverEmpDeptSAP
		request.ReceiverDeptNameShort = trnRequest.DriverEmpDeptNameShort
		request.ReceiverDeptNameFull = trnRequest.DriverEmpDeptNameFull
		request.ReceiverPosition = trnRequest.DriverEmpPosition
		request.ReceiverMobilePhone = trnRequest.DriverMobileContact
		request.ReceiverDeskPhone = trnRequest.DriverInternalContact
	} else {
		request.ReceiverType = 2 // PEA
		empUser := GetUserEmpInfo(trnRequest.VehicleUserEmpID)
		request.ReceiverPersonalId = empUser.EmpID
		request.ReceiverFullname = empUser.FullName
		request.ReceiverDeptSAP = empUser.DeptSAP
		request.ReceiverDeptNameShort = empUser.DeptSAPShort
		request.ReceiverDeptNameFull = empUser.DeptSAPFull
		request.ReceiverPosition = empUser.Position
		request.ReceiverMobilePhone = empUser.TelMobile
		request.ReceiverDeskPhone = empUser.TelInternal
	}
	if err := config.DB.Save(&request).Error; err != nil {
		return
	}
}
func ExportRequests(c *gin.Context, user *models.AuthenUserEmp, query *gorm.DB, statusNameMap map[string]string) {
	if c.Query("format") == "csv" {
		ExportRequestsCSV(c, user, query, statusNameMap)
	} else {
		ExportRequestsXLSX(c, user, query, statusNameMap)
	}
}

func ExportRequestsXLSX(c *gin.Context, user *models.AuthenUserEmp, query *gorm.DB, statusNameMap map[string]string) {
	var requests []models.VmsTrnRequestList

	// Use the keys from statusNameMap as the list of valid status codes
	statusCodes := make([]string, 0, len(statusNameMap))
	for code := range statusNameMap {
		statusCodes = append(statusCodes, code)
	}

	query = query.Table("public.vms_trn_request").
		Select(`vms_trn_request.*, v.vehicle_license_plate,v.vehicle_license_plate_province_short,v.vehicle_license_plate_province_full,
			fn_get_long_short_dept_name_by_dept_sap(d.vehicle_owner_dept_sap) vehicle_department_dept_sap_short,ref_trip_type_code,       
			(select max(mc.carpool_name) from vms_mas_carpool mc where mc.mas_carpool_uid=vms_trn_request.mas_carpool_uid) vehicle_carpool_name,
			(select log.action_detail from vms_log_request_action log where log.trn_request_uid=vms_trn_request.trn_request_uid order by log.log_request_action_datetime desc limit 1) action_detail
		`).
		Joins("LEFT JOIN vms_mas_vehicle v on v.mas_vehicle_uid = vms_trn_request.mas_vehicle_uid").
		Joins("LEFT JOIN vms_mas_vehicle_department d on d.mas_vehicle_department_uid=vms_trn_request.mas_vehicle_department_uid").
		Where("vms_trn_request.ref_request_status_code IN (?)", statusCodes)
	query = query.Where("vms_trn_request.is_deleted = ?", "0")
	// Apply additional filters (search, date range, etc.)
	if search := c.Query("search"); search != "" {
		query = query.Where("vms_trn_request.request_no ILIKE ? OR v.vehicle_license_plate ILIKE ? OR vms_trn_request.vehicle_user_emp_name ILIKE ? OR vms_trn_request.work_place ILIKE ?", "%"+search+"%", "%"+search+"%", "%"+search+"%", "%"+search+"%")
	}
	if startDate := c.Query("startdate"); startDate != "" {
		query = query.Where("vms_trn_request.reserve_end_datetime >= ?", startDate)
	}
	if endDate := c.Query("enddate"); endDate != "" {
		query = query.Where("vms_trn_request.reserve_start_datetime <= ?", endDate)
	}
	if refRequestStatusCodes := c.Query("ref_request_status_code"); refRequestStatusCodes != "" {
		// Split the comma-separated codes into a slice
		codes := strings.Split(refRequestStatusCodes, ",")
		// Include additional keys with the same text in StatusNameMapUser
		additionalCodes := make(map[string]bool)
		for _, code := range codes {
			if name, exists := statusNameMap[code]; exists {
				for key, value := range statusNameMap {
					if value == name {
						additionalCodes[key] = true
					}
				}
			}
		}
		// Merge the original codes with the additional codes
		for key := range additionalCodes {
			codes = append(codes, key)
		}
		//fmt.Println("codes", codes)
		query = query.Where("vms_trn_request.ref_request_status_code IN (?)", codes)
	}

	// Ordering
	orderBy := c.Query("order_by")
	orderDir := c.Query("order_dir")
	if orderDir != "desc" {
		orderDir = "asc"
	}
	switch orderBy {
	case "request_no":
		query = query.Order("vms_trn_request.request_no " + orderDir)
	case "start_datetime":
		query = query.Order("vms_trn_request.start_datetime " + orderDir)
	case "ref_request_status_code":
		query = query.Order("vms_trn_request.ref_request_status_code " + orderDir)
	}

	// Execute the main query
	if err := query.Scan(&requests).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "message": messages.ErrInternalServer.Error()})
		return
	}
	for i := range requests {
		requests[i].RefRequestStatusName = statusNameMap[requests[i].RefRequestStatusCode]
		switch requests[i].TripType {
		case 0:
			requests[i].TripTypeName = "ไป-กลับ"
		case 1:
			requests[i].TripTypeName = "ค้างแรม"
		}
	}

	// Set headers
	headers := []string{
		"เลขที่คำขอ",
		"ผู้ใช้ยานพาหนะ",
		"รหัสพนักงาน",
		"หน่วยงานที่สังกัด",
		"ยานพาหนะ",
		"สังกัดยานพาหนะ",
		"สถานที่ปฏิบัติงาน",
		"วันที่เดินทางเริ่มต้น",
		"วันที่เดินทางสิ้นสุด",
		"ประเภทการเดินทาง",
		"รายละเอียด",
		"สถานะคำขอ",
	}
	file := xlsx.NewFile()
	sheet, err := file.AddSheet("Booking Requests")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create Excel sheet", "message": err.Error()})
		return
	}
	headerRow := sheet.AddRow()
	for _, header := range headers {
		cell := headerRow.AddCell()
		cell.Value = header
	}

	for _, request := range requests {
		row := sheet.AddRow()
		row.AddCell().Value = request.RequestNo
		row.AddCell().Value = request.VehicleUserEmpName
		row.AddCell().Value = request.VehicleUserEmpID
		row.AddCell().Value = request.VehicleUserDeptNameShort
		row.AddCell().Value = request.VehicleLicensePlate + " " + request.VehicleLicensePlateProvinceFull
		row.AddCell().Value = request.VehicleDepartmentDeptSapShort
		row.AddCell().Value = request.WorkPlace

		row.AddCell().Value = GetDateWithZone(request.ReserveStartDatetime.Time)
		row.AddCell().Value = GetDateWithZone(request.ReserveEndDatetime.Time)
		row.AddCell().Value = request.TripTypeName
		row.AddCell().Value = request.ActionDetail
		row.AddCell().Value = request.RefRequestStatusName
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
	c.Header("Content-Disposition", "attachment; filename=requests.xlsx")
	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("File-Name", fmt.Sprintf("requests_%s.xlsx", time.Now().Format("2006-01-02")))
	c.Header("Content-Transfer-Encoding", "binary")
	if err := file.Write(c.Writer); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to write Excel file", "message": err.Error()})
		return
	}
}

func ExportRequestsCSV(c *gin.Context, user *models.AuthenUserEmp, query *gorm.DB, statusNameMap map[string]string) {
	var requests []models.VmsTrnRequestList

	// Use the keys from statusNameMap as the list of valid status codes
	statusCodes := make([]string, 0, len(statusNameMap))
	for code := range statusNameMap {
		statusCodes = append(statusCodes, code)
	}

	query = query.Table("public.vms_trn_request").
		Select(`vms_trn_request.*, v.vehicle_license_plate,v.vehicle_license_plate_province_short,v.vehicle_license_plate_province_full,
			fn_get_long_short_dept_name_by_dept_sap(d.vehicle_owner_dept_sap) vehicle_department_dept_sap_short,ref_trip_type_code,       
			(select max(mc.carpool_name) from vms_mas_carpool mc where mc.mas_carpool_uid=vms_trn_request.mas_carpool_uid) vehicle_carpool_name,
			(select log.action_detail from vms_log_request_action log where log.trn_request_uid=vms_trn_request.trn_request_uid order by log.log_request_action_datetime desc limit 1) action_detail
		`).
		Joins("LEFT JOIN vms_mas_vehicle v on v.mas_vehicle_uid = vms_trn_request.mas_vehicle_uid").
		Joins("LEFT JOIN vms_mas_vehicle_department d on d.mas_vehicle_department_uid=vms_trn_request.mas_vehicle_department_uid").
		Where("vms_trn_request.ref_request_status_code IN (?)", statusCodes)
	query = query.Where("vms_trn_request.is_deleted = ?", "0")
	// Apply additional filters (search, date range, etc.)
	if search := c.Query("search"); search != "" {
		query = query.Where("vms_trn_request.request_no ILIKE ? OR v.vehicle_license_plate ILIKE ? OR vms_trn_request.vehicle_user_emp_name ILIKE ? OR vms_trn_request.work_place ILIKE ?", "%"+search+"%", "%"+search+"%", "%"+search+"%", "%"+search+"%")
	}
	if startDate := c.Query("startdate"); startDate != "" {
		query = query.Where("vms_trn_request.reserve_end_datetime >= ?", startDate)
	}
	if endDate := c.Query("enddate"); endDate != "" {
		query = query.Where("vms_trn_request.reserve_start_datetime <= ?", endDate)
	}
	if refRequestStatusCodes := c.Query("ref_request_status_code"); refRequestStatusCodes != "" {
		// Split the comma-separated codes into a slice
		codes := strings.Split(refRequestStatusCodes, ",")
		// Include additional keys with the same text in StatusNameMapUser
		additionalCodes := make(map[string]bool)
		for _, code := range codes {
			if name, exists := statusNameMap[code]; exists {
				for key, value := range statusNameMap {
					if value == name {
						additionalCodes[key] = true
					}
				}
			}
		}
		// Merge the original codes with the additional codes
		for key := range additionalCodes {
			codes = append(codes, key)
		}
		//fmt.Println("codes", codes)
		query = query.Where("vms_trn_request.ref_request_status_code IN (?)", codes)
	}

	// Ordering
	orderBy := c.Query("order_by")
	orderDir := c.Query("order_dir")
	if orderDir != "desc" {
		orderDir = "asc"
	}
	switch orderBy {
	case "request_no":
		query = query.Order("vms_trn_request.request_no " + orderDir)
	case "start_datetime":
		query = query.Order("vms_trn_request.start_datetime " + orderDir)
	case "ref_request_status_code":
		query = query.Order("vms_trn_request.ref_request_status_code " + orderDir)
	}

	// Execute the main query
	if err := query.Scan(&requests).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "message": messages.ErrInternalServer.Error()})
		return
	}
	for i := range requests {
		requests[i].RefRequestStatusName = statusNameMap[requests[i].RefRequestStatusCode]
		switch requests[i].TripType {
		case 0:
			requests[i].TripTypeName = "ไป-กลับ"
		case 1:
			requests[i].TripTypeName = "ค้างแรม"
		}
	}

	// Set headers
	headers := []string{
		"เลขที่คำขอ",
		"ผู้ใช้ยานพาหนะ",
		"รหัสพนักงาน",
		"หน่วยงานที่สังกัด",
		"ยานพาหนะ",
		"สังกัดยานพาหนะ",
		"สถานที่ปฏิบัติงาน",
		"วันที่เดินทางเริ่มต้น",
		"วันที่เดินทางสิ้นสุด",
		"ประเภทการเดินทาง",
		"รายละเอียด",
		"สถานะคำขอ",
	}

	// Set CSV headers
	c.Header("Content-Disposition", "attachment; filename=requests.csv")
	c.Header("Content-Type", "text/csv; charset=utf-8")
	c.Header("File-Name", fmt.Sprintf("requests_%s.csv", time.Now().Format("2006-01-02")))
	c.Header("Content-Transfer-Encoding", "binary")

	writer := csv.NewWriter(c.Writer)
	defer writer.Flush()

	// Write header row
	if err := writer.Write(headers); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to write CSV header", "message": err.Error()})
		return
	}

	// Write data rows
	for _, request := range requests {
		row := []string{
			request.RequestNo,
			request.VehicleUserEmpName,
			request.VehicleUserEmpID,
			request.VehicleUserDeptNameShort,
			request.VehicleLicensePlate + " " + request.VehicleLicensePlateProvinceFull,
			request.VehicleDepartmentDeptSapShort,
			request.WorkPlace,
			request.ReserveStartDatetime.Format("2006-01-02 15:04:05"),
			request.ReserveEndDatetime.Format("2006-01-02 15:04:05"),
			request.TripTypeName,
			request.ActionDetail,
			request.RefRequestStatusName,
		}
		if err := writer.Write(row); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to write CSV row", "message": err.Error()})
			return
		}
	}
}

func UpdateDriverAvgScore(driverID string) {
	//fn_driver_avg_score('DZ000229')
	var driverScore struct {
		AvgScore        float64 `json:"avg_score"`
		TotalEvaluation int     `json:"total_evaluation"`
	}
	query := fmt.Sprintf("SELECT * FROM fn_driver_avg_score('%s')", driverID)
	if err := config.DB.Raw(query).Scan(&driverScore).Error; err != nil {
		fmt.Println("Error getting driver score:", err)
		return
	}
	//update to vms_mas_driver
	config.DB.Model(&models.VmsMasDriver{}).Where("driver_id = ?", driverID).
		Updates(map[string]interface{}{
			"driver_average_satisfaction_score": driverScore.AvgScore,
			"driver_total_satisfaction_review":  driverScore.TotalEvaluation,
		})

	fmt.Println("Driver score:", driverScore)
}
