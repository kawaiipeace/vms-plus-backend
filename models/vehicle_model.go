package models

import (
	"time"
)

type VmsMasVehicleList struct {
	MasVehicleUID                    string `gorm:"primaryKey;column:mas_vehicle_uid" json:"mas_vehicle_uid"`
	VehicleLicensePlate              string `gorm:"column:vehicle_license_plate" json:"vehicle_license_plate"`
	VehicleLicensePlateProvinceShort string `gorm:"column:vehicle_license_plate_province_short" json:"vehicle_license_plate_province_short"`
	VehicleLicensePlateProvinceFull  string `gorm:"column:vehicle_license_plate_province_full" json:"vehicle_license_plate_province_full"`
	VehicleBrandName                 string `gorm:"column:vehicle_brand_name" json:"vehicle_brand_name"`
	VehicleModelName                 string `gorm:"column:vehicle_model_name" json:"vehicle_model_name"`
	CarType                          string `gorm:"column:CarTypeDetail" json:"car_type"`
	VehiclePeaID                     string `gorm:"column:vehicle_pea_id" json:"vehicle_pea_id"`
	VehicleOwnerDeptSAP              string `gorm:"column:vehicle_owner_dept_sap" json:"vehicle_owner_dept_sap"`
	VehicleOwnerDeptShort            string `gorm:"column:vehicle_owner_dept_short" json:"vehicle_owner_dept_short"`
	VehicleImg                       string `gorm:"column:vehicle_img" json:"vehicle_img"` // Store image URL or file path
	VehicleColor                     string `gorm:"column:vehicle_color" json:"vehicle_color"`
	VehicleMileage                   int    `gorm:"column:vehicle_mileage" json:"vehicle_mileage"`
	LastMonthMileage                 int    `gorm:"column:last_month_mileage" json:"last_month_mileage"`
	Seat                             int    `gorm:"column:seat" json:"seat"`
	IsAdminChooseDriver              bool   `json:"is_admin_choose_driver"`
	CarpoolName                      string `gorm:"column:carpool_name" json:"-"`
	FleetCardNo                      string `gorm:"column:fleet_card_no" json:"fleet_card_no"`
}

func (VmsMasVehicleList) TableName() string {
	return "vms_mas_vehicle"
}

type VmsMasVehicleCarpoolList struct {
	MasVehicleUID                    string    `gorm:"primaryKey;column:mas_vehicle_uid" json:"mas_vehicle_uid"`
	VehicleLicensePlate              string    `gorm:"column:vehicle_license_plate" json:"vehicle_license_plate"`
	VehicleLicensePlateProvinceShort string    `gorm:"column:vehicle_license_plate_province_short" json:"vehicle_license_plate_province_short"`
	VehicleLicensePlateProvinceFull  string    `gorm:"column:vehicle_license_plate_province_full" json:"vehicle_license_plate_province_full"`
	VehicleBrandName                 string    `gorm:"column:vehicle_brand_name" json:"vehicle_brand_name"`
	VehicleModelName                 string    `gorm:"column:vehicle_model_name" json:"vehicle_model_name"`
	VehiclePeaID                     string    `gorm:"column:vehicle_pea_id" json:"vehicle_pea_id"`
	CarType                          string    `gorm:"column:CarTypeDetail" json:"car_type"`
	VehicleOwnerDeptSAP              string    `gorm:"column:vehicle_owner_dept_sap" json:"vehicle_owner_dept_sap"`
	VehicleOwnerDeptShort            string    `gorm:"column:vehicle_owner_dept_short" json:"vehicle_owner_dept_short"`
	VehicleImg                       string    `gorm:"column:vehicle_img" json:"vehicle_img"` // Store image URL or file path
	VehicleColor                     string    `gorm:"column:vehicle_color" json:"vehicle_color"`
	FuelTypeName                     string    `gorm:"column:fuel_type_name" json:"fuel_type_name"`
	FleetCardNo                      string    `gorm:"column:fleet_card_no" json:"fleet_card_no"`
	IsTaxCredit                      string    `gorm:"column:is_tax_credit" json:"is_tax_credit"`
	VehicleMileage                   string    `gorm:"column:vehicle_mileage" json:"vehicle_mileage"`
	Age                              string    `gorm:"column:age" json:"age"`
	RefVehicleStatusName             string    `gorm:"column:ref_vehicle_status_name" json:"ref_vehicle_status_name"`
	Seat                             int       `gorm:"column:seat" json:"seat"`
	VehicleRegistrationDate          time.Time `gorm:"column:vehicle_registration_date" json:"vehicle_registration_date"`
	IsAdminChooseDriver              bool      `json:"is_admin_choose_driver"`
}

func (VmsMasVehicleCarpoolList) TableName() string {
	return "vms_mas_vehicle"
}

func AssignVehicleImageFromIndex(vehicles []VmsMasVehicleList) []VmsMasVehicleList {
	return vehicles
}

// VmsMasCarpoolList
type VmsMasCarpoolCarBooking struct {
	MasCarpoolUID         string                 `gorm:"primaryKey;column:mas_carpool_uid" json:"mas_carpool_uid"`
	CarpoolName           string                 `gorm:"column:carpool_name" json:"carpool_name"`
	RefCarpoolChooseCarID int                    `gorm:"column:ref_carpool_choose_car_id" json:"ref_carpool_choose_car_id" example:"1"`
	RefCarpoolChooseCar   VmsRefCarpoolChooseCar `gorm:"foreignKey:RefCarpoolChooseCarID;references:RefCarpoolChooseCarID" json:"ref_carpool_choose_car"`
	IsAdminChooseDriver   bool                   `json:"is_admin_choose_driver"`
}

func (VmsMasCarpoolCarBooking) TableName() string {
	return "vms_mas_carpool"
}

// VmsMasCarpoolList
type VmsMasCarpoolCarBookingResponse struct {
	MasCarpoolUID            string                    `gorm:"primaryKey;column:mas_carpool_uid" json:"mas_carpool_uid"`
	CarpoolName              string                    `gorm:"column:carpool_name" json:"carpool_name"`
	RefCarpoolChooseCarID    int                       `gorm:"column:ref_carpool_choose_car_id" json:"ref_carpool_choose_car_id" example:"1"`
	RefCarpoolChooseCar      VmsRefCarpoolChooseCar    `gorm:"foreignKey:RefCarpoolChooseCarID;references:RefCarpoolChooseCarID" json:"ref_carpool_choose_car"`
	RefCarpoolChooseDriverID int                       `gorm:"column:ref_carpool_choose_driver_id" json:"ref_carpool_choose_driver_id" example:"1"`
	RefCarpoolChooseDriver   VmsRefCarpoolChooseDriver `gorm:"foreignKey:RefCarpoolChooseDriverID;references:RefCarpoolChooseDriverID" json:"ref_carpool_choose_driver"`
	IsAdminChooseDriver      bool                      `json:"is_admin_choose_driver"`
}

func (VmsMasCarpoolCarBookingResponse) TableName() string {
	return "vms_mas_carpool"
}

type VmsRefVehicleType struct {
	RefVehicleTypeCode int    `gorm:"column:ref_vehicle_type_code;primarykey" json:"ref_vehicle_type_code"`
	RefVehicleTypeName string `gorm:"column:ref_vehicle_type_name" json:"ref_vehicle_type_name"`
	AvailableUnits     int    `gorm:"column:available_units" json:"available_units"`
	VehicleTypeImage   string `json:"vehicle_type_image"`
}

func (VmsRefVehicleType) TableName() string {
	return "vms_ref_vehicle_type"
}

type VmsRefCarTypeDetail struct {
	CarTypeDetail string `gorm:"column:car_type_detail" json:"car_type_detail"`
}

func (VmsRefCarTypeDetail) TableName() string {
	return "vms_mas_vehicle"
}

func AssignTypeImageFromIndex(vehicle_types []VmsRefVehicleType) []VmsRefVehicleType {
	return vehicle_types
}

type VmsMasVehicle struct {
	MasVehicleUID                    string                  `gorm:"primaryKey;column:mas_vehicle_uid;type:uuid" json:"mas_vehicle_uid"`
	VehicleBrandName                 string                  `gorm:"column:vehicle_brand_name" json:"vehicle_brand_name"`
	VehicleModelName                 string                  `gorm:"column:vehicle_model_name" json:"vehicle_model_name"`
	VehicleLicensePlate              string                  `gorm:"column:vehicle_license_plate" json:"vehicle_license_plate"`
	VehicleLicensePlateProvinceShort string                  `gorm:"column:vehicle_license_plate_province_short" json:"vehicle_license_plate_province_short"`
	VehicleLicensePlateProvinceFull  string                  `gorm:"column:vehicle_license_plate_province_full" json:"vehicle_license_plate_province_full"`
	VehicleImgs                      []string                `gorm:"-" json:"vehicle_imgs"`
	CarType                          string                  `gorm:"column:CarType" json:"CarType"`
	VehicleOwnerDeptSap              string                  `gorm:"column:vehicle_owner_dept_sap" json:"vehicle_owner_dept_sap"`
	IsHasFleetCard                   byte                    `gorm:"column:is_has_fleet_card" json:"is_has_fleet_card"`
	VehicleGear                      string                  `gorm:"column:vehicle_gear" json:"vehicle_gear"`
	RefVehicleSubtypeCode            int                     `gorm:"column:ref_vehicle_subtype_code" json:"ref_vehicle_subtype_code"`
	VehicleUserEmpID                 string                  `gorm:"column:vehicle_user_emp_id" json:"vehicle_user_emp_id"`
	RefFuelTypeID                    int                     `gorm:"column:ref_fuel_type_id" json:"ref_fuel_type_id"`
	Seat                             int                     `gorm:"column:seat" json:"seat"`
	RefFuelType                      VmsRefFuelType          `gorm:"foreignKey:RefFuelTypeID;references:RefFuelTypeID" json:"ref_fuel_type"`
	VehicleRegistrationDate          time.Time               `gorm:"column:vehicle_registration_date" json:"vehicle_registration_date"`
	Age                              int                     `json:"age"`
	VehicleDepartment                VmsMasVehicleDepartment `gorm:"foreignKey:MasVehicleUID;references:MasVehicleUID" json:"vehicle_department"`
	IsAdminChooseDriver              string                  `json:"is_admin_choose_driver"`
}

func (VmsMasVehicle) TableName() string {
	return "vms_mas_vehicle"
}

type VmsMasVehicleDepartment struct {
	MasVehicleUID                    string    `gorm:"column:mas_vehicle_uid;primaryKey" json:"-"`
	MasVehicleDepartmentUID          string    `gorm:"column:mas_vehicle_department_uid" json:"-"`
	County                           string    `gorm:"column:county" json:"county"`
	VehicleGetDate                   time.Time `gorm:"column:vehicle_get_date" json:"vehicle_get_date"`
	VehiclePeaID                     string    `gorm:"column:vehicle_pea_id" json:"vehicle_pea_id"`
	VehicleLicensePlate              string    `gorm:"column:vehicle_license_plate" json:"vehicle_license_plate"`
	VehicleLicensePlateProvinceShort string    `gorm:"column:vehicle_license_plate_province_short" json:"vehicle_license_plate_province_short"`
	VehicleLicensePlateProvinceFull  string    `gorm:"column:vehicle_license_plate_province_full" json:"vehicle_license_plate_province_full"`
	VehicleAssetNo                   string    `gorm:"column:vehicle_asset_no" json:"vehicle_asset_no"`
	AssetClass                       string    `gorm:"column:asset_class" json:"asset_class"`
	AssetSubcategory                 string    `gorm:"column:asset_subcategory" json:"asset_subcategory"`
	RefPeaOfficialVehicleTypeCode    int       `gorm:"column:ref_pea_official_vehicle_type_code" json:"ref_pea_official_vehicle_type_code"`
	VehicleCondition                 int       `gorm:"column:vehicle_condition" json:"vehicle_condition"`
	VehicleMileage                   int       `gorm:"column:vehicle_mileage" json:"vehicle_mileage"`
	VehicleMileageLastMonth          int       `gorm:"vehicle_mileage_last_month" json:"vehicle_mileage_last_month"`
	VehicleOwnerDeptSap              string    `gorm:"column:vehicle_owner_dept_sap" json:"vehicle_owner_dept_sap"`
	VehicleOwnerDeptShort            string    `gorm:"-" json:"vehicle_owner_dept_short"`
	VehicleCostCenter                string    `gorm:"column:vehicle_cost_center" json:"vehicle_cost_center"`
	OwnerDeptName                    string    `gorm:"column:owner_dept_name" json:"owner_dept_name"`
	VehicleImg                       string    `gorm:"column:vehicle_img" json:"vehicle_img"`
	//VehicleUserEmpID                 string     `gorm:"column:vehicle_user_emp_id" json:"vehicle_user_emp_id"`
	//VehicleUserEmpName               string     `gorm:"column:vehicle_user_emp_name" json:"vehicle_user_emp_name"`
	//VehicleAdminEmpID                string     `gorm:"column:vehicle_admin_emp_id" json:"vehicle_admin_emp_id"`
	//VehicleAdminEmpName              string     `gorm:"column:vehicle_admin_emp_name" json:"vehicle_admin_emp_name"`
	ParkingPlace         string     `gorm:"column:parking_place" json:"parking_place"`
	FleetCardNo          string     `gorm:"column:fleet_card_no" json:"fleet_card_no"`
	FleetCardOilStations string     `gorm:"column:fleet_card_oil_stations" json:"fleet_card_oil_stations"`
	IsInCarpool          string     `gorm:"column:is_in_carpool" json:"is_in_carpool"`
	Remark               string     `gorm:"column:remark" json:"remark"`
	RefVehicleStatusCode int        `gorm:"column:ref_vehicle_status_code" json:"ref_vehicle_status_code"`
	RefOtherUseCode      *int       `gorm:"column:ref_other_use_code" json:"ref_other_use_code"`
	VehicleUser          MasUserEmp `gorm:"-" json:"vehicle_user"`
}

func (VmsMasVehicleDepartment) TableName() string {
	return "vms_mas_vehicle_department"
}

// VmsMasVehicleCanBooking
type VmsMasVehicleCanBooking struct {
	MasVehicleUID            string `gorm:"column:mas_vehicle_uid" json:"-"`
	MasCarpoolUID            string `gorm:"column:mas_carpool_uid" json:"-"`
	CarpoolName              string `gorm:"column:carpool_name" json:"carpool_name"`
	RefCarpoolChooseCarID    int    `gorm:"column:ref_carpool_choose_car_id" json:"ref_carpool_choose_car_id"`
	RefCarpoolChooseDriverID int    `gorm:"column:ref_carpool_choose_driver_id" json:"ref_carpool_choose_driver_id"`
}

func (VmsMasVehicleCanBooking) TableName() string {
	return "vms_mas_vehicle_can_booking"
}

type VmsMasVehicleImg struct {
	MasVehicleUID         string `gorm:"column:mas_vehicle_uid" json:"-"`
	RefVehicleImgSideCode int    `gorm:"column:ref_vehicle_img_side_code" json:"ref_vehicle_img_side_code"`
	VehicleImgFile        string `gorm:"column:vehicle_img_file" json:"vehicle_img_file"`
}

func (VmsMasVehicleImg) TableName() string {
	return "vms_mas_vehicle_img"
}
