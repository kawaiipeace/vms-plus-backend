package funcs

import (
	"log"
	"time"
	"vms_plus_be/config"
	"vms_plus_be/models"

	"github.com/google/uuid"
)

func CreateTrnRequestActionLog(trnRequestUID, refStatusCode, actionDetail, actionByPersonalID, actionByRole, remark string) error {
	var user models.MasUserEmp
	if actionByRole == "driver" {
		//user = GetUserEmpInfo(actionByPersonalID)
	} else {
		user = GetUserEmpInfo(actionByPersonalID)
	}
	logReq := models.VmsLogRequest{
		LogRequestActionUID:      uuid.New().String(),
		TrnRequestUID:            trnRequestUID,
		RefRequestStatusCode:     refStatusCode,
		LogRequestActionDatetime: time.Now(),
		ActionByPersonalID:       actionByPersonalID,
		ActionByRole:             actionByRole,
		ActionByFullname:         user.FullName,
		ActionByPosition:         user.Position,
		ActionByDepartment:       user.DeptSAP,
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

func CreateTrnLog1(trnRequestUID, refStatusCode, logRemark, createdBy string) error {
	return nil
}
