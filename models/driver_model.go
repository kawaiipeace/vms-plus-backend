package models

import (
	"fmt"
	"time"
)

type VmsMasDriver struct {
	MasDriverUID                   string                `gorm:"primaryKey;column:mas_driver_uid;type:uuid" json:"mas_driver_uid"`
	DriverID                       string                `gorm:"column:driver_id" json:"driver_id"`
	DriverName                     string                `gorm:"column:driver_name" json:"driver_name"`
	DriverImage                    string                `gorm:"column:driver_image" json:"driver_image"`
	DriverNickname                 string                `gorm:"column:driver_nickname" json:"driver_nickname"`
	DriverDeptSAP                  string                `gorm:"column:driver_dept_sap" json:"driver_dept_sap"`
	DriverIdentificationNo         string                `gorm:"column:driver_identification_no" json:"driver_identification_no"`
	DriverContactNumber            string                `gorm:"column:driver_contact_number" json:"driver_contact_number"`
	DriverAverageSatisfactionScore float64               `gorm:"column:driver_average_satisfaction_score" json:"driver_average_satisfaction_score"`
	DriverBirthdate                time.Time             `gorm:"column:driver_birthdate" json:"driver_birthdate"`
	Age                            string                `json:"age"`
	Status                         string                `json:"status"`
	DriverTripDetails              []VmsDriverTripDetail `gorm:"-" json:"trip_Details"`
}

func (VmsMasDriver) TableName() string {
	return "vms_mas_driver"
}

type VmsDriverTripDetail struct {
	TrnRequestUID string     `gorm:"column:trn_request_uid" json:"trn_request_uid" example:"456e4567-e89b-12d3-a456-426614174001"`
	RequestNo     string     `gorm:"column:request_no" json:"request_no"`
	WorkPlace     string     `gorm:"column:work_place" json:"work_place"`
	StartDatetime string     `gorm:"column:start_datetime" json:"start_datetime"`
	EndDatetime   string     `gorm:"column:end_datetime" json:"end_datetime"`
	VehicleUser   MasUserEmp `gorm:"foreignKey:VehicleUserEmpID;references:EmpID" json:"vehicle_user"`
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

// VmsTrnAnnualDriver
type VmsTrnAnnualDriver struct {
	TrnRequestAnnualDriverUid string    `gorm:"column:trn_request_annual_driver_uid" json:"-"`
	CreatedRequestEmpId       string    `gorm:"column:created_request_emp_id" json:"-"`
	RequestAnnualDriverNo     string    `gorm:"column:request_annual_driver_no" json:"request_annual_driver_no"`
	RequestIssueDate          time.Time `gorm:"column:request_issue_date" json:"request_issue_date"`
	RequestExpireDate         time.Time `gorm:"column:request_expire_date" json:"request_expire_date"`
	AnnualYYYY                int       `gorm:"column:annual_yyyy" json:"annual_yyyy"`
	DriverLicenseNo           string    `gorm:"column:driver_license_no" json:"driver_license_no"`
	DriverLicenseExpireDate   time.Time `gorm:"column:driver_license_expire_date" json:"driver_license_expire_date"`
}

func (VmsTrnAnnualDriver) TableName() string {
	return "vms_trn_request_annual_driver"
}
