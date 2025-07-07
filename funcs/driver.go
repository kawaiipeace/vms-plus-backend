package funcs

import (
	"time"
	"vms_plus_be/config"
	"vms_plus_be/models"

	"gorm.io/gorm"
)

func CheckDriverIsActive(masDriverUID string) {
	var driver models.VmsMasDriverResponse
	err := config.DB.Where("mas_driver_uid = ? AND is_deleted = ?", masDriverUID, "0").
		Preload("DriverLicense", func(db *gorm.DB) *gorm.DB {
			return db.Order("driver_license_end_date DESC").Limit(1)
		}).
		First(&driver).Error
	if err != nil {
		return
	}
	isActive := "1"
	refDriverStatusCode := 1

	if driver.IsReplacement == "1" {
		refDriverStatusCode = 6
	}
	//not in (1,6,7)
	if driver.RefDriverStatusCode != 1 && driver.RefDriverStatusCode != 6 && driver.RefDriverStatusCode != 7 {
		isActive = "0"
	}

	if driver.ApprovedJobDriverStartDate.Before(time.Now()) {
		isActive = "0"
	}

	if driver.ApprovedJobDriverEndDate.After(time.Now()) {
		isActive = "0"
	}

	if driver.DriverLicense.DriverLicenseEndDate.Before(time.Now()) {
		isActive = "0"
	}

	//vms_trn_driver_leave
	query := config.DB.Table("vms_trn_driver_leave").Where("mas_driver_uid = ? AND is_deleted = ?", masDriverUID, "0")
	query = query.Where("leave_start_date <= ?", time.Now())
	query = query.Where("leave_end_date >= ?", time.Now())
	var driverLeave models.VmsMasDriverLeaveStatusUpdate
	if err := query.First(&driverLeave).Error; err == nil {
		if driverLeave.TrnDriverLeaveUID != "" {
			if driverLeave.ReplacementDriverUID != "" {
				isActive = "0"
				refDriverStatusCode = 2
				//update is_active to 1, is_replacement to 0
				if err := config.DB.Model(&models.VmsMasDriver{}).Where("mas_driver_uid = ?", driverLeave.ReplacementDriverUID).
					Update("is_active", "1").
					Update("ref_driver_status_code", 7).
					Update("is_replacement", "0").Error; err != nil {
					return
				}
			}
		}
	}

	//update is_active
	if err := config.DB.Table("vms_mas_driver").Where("mas_driver_uid = ?", masDriverUID).
		Where("is_active != ? or ref_driver_status_code != ?", isActive, refDriverStatusCode).
		Update("is_active", isActive).
		Update("ref_driver_status_code", refDriverStatusCode).Error; err != nil {
		return
	}
}

func JobDriversCheckActive() {
	//get all job drivers
	var jobDrivers []models.VmsMasDriverResponse
	err := config.DB.Where("is_deleted = ?", "0").
		Find(&jobDrivers).Error
	if err != nil {
		return
	}

	for _, driver := range jobDrivers {
		CheckDriverIsActive(driver.MasDriverUID)
	}
}
