package models

import (
	"time"
)

// VmsMasVehicleManagementList
type VmsMasVehicleManagementList struct {
	MasVihicleUID        string              `gorm:"primaryKey;column:mas_vehicle_uid" json:"mas_vehicle_uid"`
	VehicleLicensePlate  string              `gorm:"column:vehicle_license_plate" json:"vehicle_license_plate"`
	VehicleBrandName     string              `gorm:"column:vehicle_brand_name" json:"vehicle_brand_name"`
	VehicleModelName     string              `gorm:"column:vehicle_model_name" json:"vehicle_model_name"`
	RefVehicleTypeCode   string              `gorm:"column:ref_vehicle_type_code" json:"ref_vehicle_type_code"`
	RefVehicleTypeName   string              `gorm:"column:ref_vehicle_type_name" json:"ref_vehicle_type_name"`
	VehicleOwnerDeptSAP  string              `gorm:"column:vehicle_owner_dept_short" json:"vehicle_owner_dept_short"`
	FleetCardNo          string              `gorm:"column:fleet_card_no" json:"fleet_card_no"`
	IsTaxCredit          bool                `gorm:"column:is_tax_credit" json:"is_tax_credit"`
	VehicleMileage       float64             `gorm:"column:vehicle_mileage" json:"vehicle_mileage"`
	VehicleGetDate       time.Time           `gorm:"column:vehicle_get_date" json:"vehicle_get_date"` // Changed to time.Time
	RefVehicleStatusCode string              `gorm:"column:ref_vehicle_status_code" json:"ref_vehicle_status_code"`
	VmsRefVehicleStatus  VmsRefVehicleStatus `gorm:"foreignKey:RefVehicleStatusCode;references:RefVehicleStatusCode" json:"vms_ref_vehicle_status"`
	VehicleCarpoolName   string              `gorm:"column:vehicle_carpool_name" json:"vehicle_carpool_name"`
	IsAcictive           string              `gorm:"column:is_active" json:"is_active"`
	RefFuelTypeID        int                 `gorm:"column:ref_fuel_type_id" json:"ref_fuel_type_id"`
	VmsRefFuelType       VmsRefFuelType      `gorm:"foreignKey:RefFuelTypeID;references:RefFuelTypeID" json:"vms_ref_fuel_type"`
	Age                  string              `json:"age"`
}

func (VmsMasVehicleManagementList) TableName() string {
	return "vms_mas_vehicle"
}

// VmsMasVehicleIsActiveUpdate
type VmsMasVehicleIsActiveUpdate struct {
	MasVehicleUID string    `gorm:"primaryKey;column:mas_vehicle_uid" json:"mas_vehicle_uid" example:"f3b29096-140e-49dc-97ee-17fa9352aff6"`
	IsActive      string    `gorm:"column:is_active" json:"is_active" example:"1"`
	UpdatedAt     time.Time `gorm:"column:updated_at" json:"-"`
	UpdatedBy     string    `gorm:"column:updated_by" json:"-"`
}

func (VmsMasVehicleIsActiveUpdate) TableName() string {
	return "vms_mas_vehicle_department"
}
