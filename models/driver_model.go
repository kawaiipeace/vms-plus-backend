package models

import (
	"fmt"
	"time"
)

type VmsMasDriver struct {
	MasDriverUID                   string    `gorm:"primaryKey;column:mas_driver_uid;type:uuid" json:"mas_driver_uid"`
	DriverName                     string    `gorm:"column:driver_name" json:"driver_name"`
	DriverImage                    string    `gorm:"column:driver_image" json:"driver_image"`
	DriverNickname                 string    `gorm:"column:driver_nickname" json:"driver_nickname"`
	DriverDeptSAP                  string    `gorm:"column:driver_dept_sap" json:"driver_dept_sap"`
	DriverIdentificationNo         string    `gorm:"column:driver_identification_no" json:"driver_identification_no"`
	DriverContactNumber            string    `gorm:"column:driver_contact_number" json:"driver_contact_number"`
	DriverAverageSatisfactionScore float64   `gorm:"column:driver_average_satisfaction_score" json:"driver_average_satisfaction_score"`
	DriverBirthdate                time.Time `gorm:"column:driver_birthdate" json:"driver_birthdate"`
	Age                            string    `json:"age"`
}

func (VmsMasDriver) TableName() string {
	return "vms_mas_driver"
}

func (d *VmsMasDriver) CalculateAgeInYearsMonths() string {
	if d.DriverBirthdate.IsZero() {
		return "ไม่ระบุ"
	}

	today := time.Now()
	years := today.Year() - d.DriverBirthdate.Year()
	months := today.Month() - d.DriverBirthdate.Month()

	// Adjust if birthday hasn't occurred yet this year
	if today.Day() < d.DriverBirthdate.Day() {
		months--
	}

	if months < 0 {
		years--
		months += 12
	}

	return fmt.Sprintf("%d ปี %d เดือน", years, months)
}
