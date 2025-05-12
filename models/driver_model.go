package models

import (
	"fmt"
	"time"
)

type VmsMasDriverShort struct {
	MasDriverUID   string `gorm:"primaryKey;column:mas_driver_uid;type:uuid" json:"mas_driver_uid"`
	DriverID       string `gorm:"column:driver_id" json:"driver_id"`
	DriverName     string `gorm:"column:driver_name" json:"driver_name"`
	DriverImage    string `gorm:"column:driver_image" json:"driver_image"`
	DriverNickname string `gorm:"column:driver_nickname" json:"driver_nickname"`
}

func (VmsMasDriverShort) TableName() string {
	return "vms_mas_driver"
}

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
	DriverTotalSatisfactionReview  int                   `gorm:"column:driver_total_satisfaction_review" json:"driver_total_satisfaction_review"`
	DriverBirthdate                time.Time             `gorm:"column:driver_birthdate" json:"driver_birthdate"`
	WorkType                       int                   `gorm:"column:work_type" json:"work_type"`
	WorkTypeName                   string                `gorm:"column:work_type_name" json:"work_type_name"`
	ContractNo                     string                `gorm:"column:contract_no" json:"contract_no"`
	EndDate                        time.Time             `gorm:"column:end_date" json:"contract_end_date"`
	Age                            string                `json:"age"`
	Status                         string                `gorm:"-" json:"status"`
	RefDriverStatusCode            int                   `gorm:"column:ref_driver_status_code" json:"-"`
	DriverStatus                   VmsRefDriverStatus    `gorm:"foreignKey:RefDriverStatusCode;references:RefDriverStatusCode" json:"driver_status"`
	WorkDays                       int                   `gorm:"-" json:"work_days"`
	WorkCount                      int                   `gorm:"-" json:"work_count"`
	DriverTripDetails              []VmsDriverTripDetail `gorm:"-" json:"trip_Details"`
	DriverLicense                  VmsMasDriverLicense   `gorm:"foreignKey:MasDriverUID;references:MasDriverUID" json:"driver_license"`
}

func (VmsMasDriver) TableName() string {
	return "vms_mas_driver"
}

// VmsMasDriverLicense
type VmsMasDriverLicense struct {
	MasDriverLicenseUID      string                  `gorm:"column:mas_driver_license_uid;primaryKey" json:"mas_driver_license_uid"`
	MasDriverUID             string                  `gorm:"column:mas_driver_uid;type:uuid" json:"mas_driver_uid"`
	RefDriverLicenseTypeCode string                  `gorm:"column:ref_driver_license_type_code;type:varchar(2)" json:"ref_driver_license_type_code"`
	DriverLicenseNo          string                  `gorm:"column:driver_license_no;type:varchar(10)" json:"driver_license_no"`
	DriverLicenseStartDate   time.Time               `gorm:"column:driver_license_start_date" json:"driver_license_start_date"`
	DriverLicenseEndDate     time.Time               `gorm:"column:driver_license_end_date" json:"driver_license_end_date"`
	DriverLicenseType        VmsRefDriverLicenseType `gorm:"foreignKey:RefDriverLicenseTypeCode;references:RefDriverLicenseTypeCode" json:"driver_license_type"`
}

func (VmsMasDriverLicense) TableName() string {
	return "vms_mas_driver_license"
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

// VmsRefDriverStatus
type VmsRefDriverStatus struct {
	RefDriverStatusCode int    `gorm:"column:ref_driver_status_code;primaryKey" json:"ref_driver_status_code"`
	RefDriverStatusDesc string `gorm:"column:ref_driver_status_desc" json:"ref_driver_status_desc"`
}

func (VmsRefDriverStatus) TableName() string {
	return "vms_ref_driver_status"
}
