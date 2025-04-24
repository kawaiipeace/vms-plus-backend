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
	updateQuery := `
		UPDATE public.vms_trn_request
		SET vehicle_department_dept_sap = d.vehicle_owner_dept_sap
		FROM vms_mas_vehicle_department d
		WHERE d.mas_vehicle_uid::character varying = vms_trn_request.mas_vehicle_uid
		AND trn_request_uid = ?`

	if err := config.DB.Exec(updateQuery, trnRequestUID).Error; err != nil {
		return err
	}

	updateQuery = `
		UPDATE public.vms_trn_request
		SET vehicle_department_dept_sap_short = d.dept_short,
		    vehicle_department_dept_sap_full = d.dept_full
		FROM vms_mas_department d
		WHERE d.dept_sap = vms_trn_request.vehicle_department_dept_sap
		AND trn_request_uid = ?`

	if err := config.DB.Exec(updateQuery, trnRequestUID).Error; err != nil {
		return err
	}

	updateQuery = `
		UPDATE public.vms_trn_request
		SET vehicle_license_plate = v.vehicle_license_plate,
			vehicle_license_plate_province_short = v.vehicle_license_plate_province_short,
			vehicle_license_plate_province_full = v.vehicle_license_plate_province_full
		FROM vms_mas_vehicle v
		WHERE v.mas_vehicle_uid::text = vms_trn_request.mas_vehicle_uid
		AND trn_request_uid = ?`

	if err := config.DB.Exec(updateQuery, trnRequestUID).Error; err != nil {
		return err
	}
	return nil
}
