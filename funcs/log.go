package funcs

import (
	"fmt"
	"log"
	"time"
	"vms_plus_be/config"
	"vms_plus_be/models"

	"github.com/google/uuid"
)

func CreateTrnRequestActionLog(trnRequestUID, refStatusCode, requestDetail, actionByPersonalID, actionByRole, requestRemark string) error {
	var user models.MasUserEmp
	if actionByRole == "driver" {
		//user = GetUserEmpInfo(actionByPersonalID)
	} else {
		user = GetUserEmpInfo(actionByPersonalID)
	}
	UpdateDetailToRequest(trnRequestUID, requestDetail)
	actionDetail, remark := GetActionDetail(trnRequestUID, refStatusCode, requestDetail, requestRemark)

	logReq := models.VmsLogRequest{
		LogRequestActionUID:      uuid.New().String(),
		TrnRequestUID:            trnRequestUID,
		RefRequestStatusCode:     refStatusCode,
		LogRequestActionDatetime: models.TimeWithZone{Time: time.Now()},
		ActionByPersonalID:       actionByPersonalID,
		ActionByRole:             actionByRole,
		ActionByFullname:         user.FullName,
		ActionByPosition:         user.Position,
		ActionByDepartment:       user.DeptSAPShort,
		ActionDetail:             actionDetail,
		Remark:                   remark,
		IsDeleted:                "0",
	}

	// Insert into database
	if err := config.DB.Create(&logReq).Error; err != nil {
		log.Println("Error inserting log:", err)
		return err
	}

	CreateRequestBookingNotification(trnRequestUID)

	return nil
}

func CreateTrnRequestAnnualLicenseActionLog(trnAnnualLicenseUID, refStatusCode, actionDetail, actionByPersonalID, actionByRole, remark string) error {
	CreateRequestAnnualLicenseNotification(trnAnnualLicenseUID)
	return nil
}

func UpdateDetailToRequest(trnRequestUID, action_detail string) error {
	//update detail to request
	config.DB.Table("vms_trn_request").
		Where("trn_request_uid = ?", trnRequestUID).
		Update("action_detail", action_detail)

	return nil
}

func GetActionDetail(trnRequestUID, refStatusCode, requestDetail, requestRemark string) (string, string) {
	actionDetail := ""
	remark := ""
	switch refStatusCode {
	case "20":
		actionDetail = "สร้างคำขอ"
	case "21":
		actionDetail = "ตีกลับคำขอ"
		remark = requestRemark
	case "30":
		actionDetail = "อนุมัติคำขอ"
	case "31":
		actionDetail = "ตีกลับคำขอ"
		remark = requestRemark
	case "40":
		actionDetail = "ผ่านการตรวจสอบ"
	case "41":
		actionDetail = "ตีกลับคำขอ"
		remark = requestRemark
	case "50":
		actionDetail = "อนุมัติให้ใช้ยานพาหนะ"
	case "51":
		actionDetail = "รับกุญแจ"

		var receivedKey struct {
			ActualReceiveTime     models.TimeWithZone `gorm:"column:actual_receive_time"`
			ReceiverFullname      string              `gorm:"column:receiver_fullname"`
			ReceiverPersonalID    string              `gorm:"column:receiver_personal_id"`
			RefVehicleKeyTypeName string              `gorm:"column:ref_vehicle_key_type_name"`
		}

		if err := config.DB.Table("vms_trn_vehicle_key_handover").
			Select("vms_trn_vehicle_key_handover.*, vms_ref_vehicle_key_type.ref_vehicle_key_type_name").
			Joins("LEFT JOIN vms_ref_vehicle_key_type ON vms_ref_vehicle_key_type.ref_vehicle_key_type_code = vms_trn_vehicle_key_handover.ref_vehicle_key_type_code").
			Where("trn_request_uid = ?", trnRequestUID).
			First(&receivedKey).Error; err != nil {
			remark = GetDateTimeBuddhistYear(time.Now())
		} else {
			remark = GetDateTimeBuddhistYear(receivedKey.ActualReceiveTime.Time) + " - " + receivedKey.ReceiverFullname + "(" + receivedKey.ReceiverPersonalID + ") เป็นผู้มารับ" + receivedKey.RefVehicleKeyTypeName

		}
	case "60":
		actionDetail = "รับยานพาหนะ"
	case "70":
		actionDetail = "คืนยานพาหนะ"
	case "71":
		actionDetail = "ตีกลับยานพาหนะ"
		remark = requestRemark
	case "80":
		actionDetail = "เสร็จสิ้น"
		var request struct {
			InspectVehicleDatetime models.TimeWithZone `gorm:"column:inspected_vehicle_datetime"`
		}
		if err := config.DB.Table("vms_trn_request").
			Where("trn_request_uid = ?", trnRequestUID).
			First(&request).Error; err != nil {
			remark = GetDateTimeBuddhistYear(time.Now()) + " - ยืนยันการคืนยานพาหนะ"
		} else {
			remark = GetDateTimeBuddhistYear(request.InspectVehicleDatetime.Time) + " - ยืนยันการคืนยานพาหนะ"
		}
	case "90":
		actionDetail = "ยกเลิกคำขอ"
		remark = requestRemark
	}

	return actionDetail, remark
}

func UpdateActionDetailToLogRequest() error {

	var logRequestAction []struct {
		LogRequestActionUID  string `gorm:"column:log_request_action_uid"`
		TrnRequestUID        string `gorm:"column:trn_request_uid"`
		RefRequestStatusCode string `gorm:"column:ref_request_status_code"`
		Remark               string `gorm:"column:remark"`
	}
	if err := config.DB.Table("vms_log_request_action").
		Find(&logRequestAction).Error; err != nil {
		return err
	}

	for _, action := range logRequestAction {
		actionDetail, remark := GetActionDetail(action.TrnRequestUID, action.RefRequestStatusCode, "", action.Remark)
		fmt.Println(action.LogRequestActionUID, action.TrnRequestUID, action.RefRequestStatusCode, action.Remark, actionDetail, remark)
		config.DB.Table("vms_log_request_action").
			Where("log_request_action_uid = ?", action.LogRequestActionUID).
			Updates(map[string]interface{}{
				"action_detail": actionDetail,
				"remark":        remark,
			})
	}

	return nil
}
