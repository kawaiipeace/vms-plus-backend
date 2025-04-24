package models

import (
	"time"
)

// VmsMasVehicleManagementList
type VmsMasVehicleManagementList struct {
	VehicleLicensePlate  string    `gorm:"column:vehicle_license_plate" json:"vehicle_license_plate"`
	VehicleBrandName     string    `gorm:"column:vehicle_brand_name" json:"vehicle_brand_name"`
	VehicleModelName     string    `gorm:"column:vehicle_model_name" json:"vehicle_model_name"`
	RefVehicleTypeCode   string    `gorm:"column:ref_vehicle_type_code" json:"ref_vehicle_type_code"`
	RefVehicleTypeName   string    `gorm:"column:ref_vehicle_type_name" json:"ref_vehicle_type_name"`
	VehicleOwnerDeptSAP  string    `gorm:"column:vehicle_owner_dept_short" json:"vehicle_owner_dept_short"`
	FleetCardNo          string    `gorm:"column:fleet_card_no" json:"fleet_card_no"`
	IsTaxCredit          bool      `gorm:"column:is_tax_credit" json:"is_tax_credit"`
	VehicleMileage       float64   `gorm:"column:vehicle_mileage" json:"vehicle_mileage"`
	VehicleGetDate       time.Time `gorm:"column:vehicle_get_date" json:"vehicle_get_date"` // Changed to time.Time
	RefVehicleStatusCode string    `gorm:"column:ref_vehicle_status_code" json:"ref_vehicle_status_code"`
	VehicleCarpoolName   string    `gorm:"column:vehicle_carpool_name" json:"vehicle_carpool_name"`
	Age                  string    `json:"age"`
}
