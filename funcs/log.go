package funcs

import (
	"log"
	"time"
	"vms_plus_be/config"
	"vms_plus_be/models"

	"github.com/google/uuid"
)

func CreateTrnLog(trnRequestUID, refStatusCode, logRemark, createdBy string) error {
	logReq := models.VmsLogRequest{
		LogRequestUID: uuid.New().String(),
		TrnRequestUID: trnRequestUID,
		RefStatusCode: refStatusCode,
		LogRemark:     logRemark,
		CreatedAt:     time.Now(),
		CreatedBy:     createdBy,
	}

	// Insert into database
	if err := config.DB.Create(&logReq).Error; err != nil {
		log.Println("Error inserting log:", err)
		return err
	}

	return nil
}

func UpdateTrnRequestData(trnRequestUID string) error {
	// First update query
	updateQuery1 := `
		UPDATE public.vms_trn_request
		SET vehicle_department_dept_sap = d.vehicle_owner_dept_sap
		FROM vms_mas_vehicle_department d
		WHERE d.mas_vehicle_uid::character varying = vms_trn_request.mas_vehicle_uid
		AND trn_request_uid = ?`

	if err := config.DB.Exec(updateQuery1, trnRequestUID).Error; err != nil {
		return err
	}

	// Second update query
	updateQuery2 := `
		UPDATE public.vms_trn_request
		SET vehicle_department_dept_sap_short = d.dept_short,
		    vehicle_department_dept_sap_full = d.dept_full
		FROM vms_mas_department d
		WHERE d.dept_sap = vms_trn_request.vehicle_department_dept_sap
		AND trn_request_uid = ?`

	if err := config.DB.Exec(updateQuery2, trnRequestUID).Error; err != nil {
		return err
	}

	return nil
}
