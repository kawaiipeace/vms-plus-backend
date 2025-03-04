package funcs

import (
	"log"
	"time"
	"vms_plus_be/config"
	"vms_plus_be/models"
)

func CreateTrnLog(trnRequestUID, refStatusCode, logRemark, createdBy string) error {
	logReq := models.VmsLogRequest{
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
