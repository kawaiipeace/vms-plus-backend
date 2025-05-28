package funcs

import (
	"net/http"
	"sort"
	"strings"
	"time"
	"vms_plus_be/config"
	"vms_plus_be/messages"
	"vms_plus_be/models"

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
		var count int64
		query = query.Debug()
		if err := query.Table("vms_trn_request").Where("ref_request_status_code IN ?", statusCodes).Count(&count).Error; err != nil {
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
		Preload("MasVehicle.RefFuelType").
		Preload("MasVehicle.VehicleDepartment").
		Preload("RefCostType").
		Preload("MasDriver").
		Preload("RefRequestStatus").
		Preload("RefTripType").
		Preload("RefCostType").
		First(&request, "trn_request_uid = ?", trnRequestUID).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{})
		return request, err
	}
	if request.MasDriver.DriverBirthdate != (time.Time{}) {
		request.MasDriver.Age = request.MasDriver.CalculateAgeInYearsMonths()
	}
	request.NumberOfAvailableDrivers = 2
	request.DriverImageURL = config.DefaultAvatarURL
	request.CanCancelRequest = true
	request.IsUseDriver = request.MasCarpoolDriverUID != ""
	request.RefRequestStatusName = StatusNameMap[request.RefRequestStatusCode]
	request.VehicleLicensePlate = request.MasVehicle.VehicleLicensePlate
	request.VehicleLicensePlateProvinceShort = request.MasVehicle.VehicleLicensePlateProvinceShort
	request.VehicleLicensePlateProvinceFull = request.MasVehicle.VehicleLicensePlateProvinceFull

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
		Preload("VehicleImagesReturned").
		Preload("VehicleImageInspect").
		Preload("ReceiverKeyTypeDetail").
		Preload("RefTripType").
		Preload("SatisfactionSurveyAnswers.SatisfactionSurveyQuestions").
		Select("vms_trn_request.*,k.receiver_personal_id,k.receiver_fullname,k.receiver_dept_sap,"+
			"k.receiver_dept_name_short,k.receiver_dept_name_full,k.receiver_desk_phone,k.receiver_mobile_phone,k.receiver_position,k.remark receiver_remark").
		Joins("LEFT JOIN vms_trn_vehicle_key_handover k ON k.trn_request_uid = vms_trn_request.trn_request_uid").
		First(&request, "vms_trn_request.trn_request_uid = ?", trnRequestUID).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{})
		return request, err
	}
	if request.MasDriver.DriverBirthdate != (time.Time{}) {
		request.MasDriver.Age = request.MasDriver.CalculateAgeInYearsMonths()
	}
	request.NumberOfAvailableDrivers = 2
	request.DriverImageURL = config.DefaultAvatarURL
	request.ReceivedKeyImageURL = config.DefaultAvatarURL
	request.CanCancelRequest = true
	request.IsUseDriver = request.MasCarpoolDriverUID != ""
	request.RefRequestStatusName = StatusNameMap[request.RefRequestStatusCode]
	request.FleetCardNo = request.MasVehicle.VehicleDepartment.FleetCardNo
	request.VehicleLicensePlate = request.MasVehicle.VehicleLicensePlate
	request.VehicleLicensePlateProvinceShort = request.MasVehicle.VehicleLicensePlateProvinceShort
	request.VehicleLicensePlateProvinceFull = request.MasVehicle.VehicleLicensePlateProvinceFull

	if err := config.DB.
		Preload("RefTripType").
		Where("trn_request_uid <> ?", request.TrnRequestUID).
		Order("created_at DESC").
		First(&request.NextRequest).Error; err == nil {
		request.NextRequest.RefRequestStatusName = StatusNameMap[request.NextRequest.RefRequestStatusCode]
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
