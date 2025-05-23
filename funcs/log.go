package funcs

import (
	"encoding/csv"
	"errors"
	"io"
	"log"
	"time"
	"vms_plus_be/config"
	"vms_plus_be/models"

	"github.com/google/uuid"
)

func CreateTrnLog(trnRequestUID, refStatusCode, logRemark, createdBy string) error {
	return nil
}

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

	return nil
}

// ParseCSV parses a CSV file and returns a slice of maps where each map represents a row with column names as keys.
func ParseCSV(reader io.Reader) ([]map[string]string, error) {
	csvReader := csv.NewReader(reader)
	headers, err := csvReader.Read()
	if err != nil {
		return nil, errors.New("failed to read CSV headers")
	}

	var records []map[string]string
	for {
		row, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, errors.New("failed to read CSV row")
		}

		record := make(map[string]string)
		for i, header := range headers {
			record[header] = row[i]
		}
		records = append(records, record)
	}

	return records, nil
}
