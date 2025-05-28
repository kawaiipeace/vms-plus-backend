package models

import (
	"strings"
	"time"
)

type VmsMasVehicleList struct {
	MasVehicleUID                    string `gorm:"primaryKey;column:mas_vehicle_uid" json:"mas_vehicle_uid"`
	VehicleLicensePlate              string `gorm:"column:vehicle_license_plate" json:"vehicle_license_plate"`
	VehicleLicensePlateProvinceShort string `gorm:"column:vehicle_license_plate_province_short" json:"vehicle_license_plate_province_short"`
	VehicleLicensePlateProvinceFull  string `gorm:"column:vehicle_license_plate_province_full" json:"vehicle_license_plate_province_full"`
	VehicleBrandName                 string `gorm:"column:vehicle_brand_name" json:"vehicle_brand_name"`
	VehicleModelName                 string `gorm:"column:vehicle_model_name" json:"vehicle_model_name"`
	CarType                          string `gorm:"column:car_type" json:"car_type"`
	VehicleOwnerDeptSAP              string `gorm:"column:dept_short" json:"vehicle_owner_dept_sap"`
	VehicleImg                       string `gorm:"column:vehicle_img" json:"vehicle_img"` // Store image URL or file path
	VehicleColor                     string `gorm:"column:vehicle_color" json:"vehicle_color"`
	Seat                             int    `gorm:"column:Seat" json:"seat"`
	IsAdminChooseDriver              bool   `json:"is_admin_choose_driver"`
}

func (VmsMasVehicleList) TableName() string {
	return "vms_mas_vehicle"
}
func AssignVehicleImageFromIndex(vehicles []VmsMasVehicleList) []VmsMasVehicleList {
	// List of random URLs
	imageUrls := []string{
		"http://pntdev.ddns.net:28089/VMS_PLUS/PIX/cars/Vehicle-1.svg",
		"http://pntdev.ddns.net:28089/VMS_PLUS/PIX/cars/Vehicle-2.svg",
		"http://pntdev.ddns.net:28089/VMS_PLUS/PIX/cars/Vehicle-3.svg",
		"http://pntdev.ddns.net:28089/VMS_PLUS/PIX/cars/Vehicle-4.svg",
	}

	// Seed the random generator
	for i := range vehicles {
		vehicles[i].VehicleImg = imageUrls[i%len(imageUrls)]
		if strings.TrimSpace(vehicles[i].VehicleLicensePlate) == "7กษ 4377" {
			vehicles[i].IsAdminChooseDriver = true
		}
	}
	return vehicles
}

// VmsMasCarpoolList
type VmsMasCarpoolCarBooking struct {
	MasCarpoolUID         string                 `gorm:"primaryKey;column:mas_carpool_uid" json:"mas_carpool_uid"`
	CarpoolName           string                 `gorm:"column:carpool_name" json:"carpool_name"`
	RefCarpoolChooseCarID int                    `gorm:"column:ref_carpool_choose_car_id" json:"ref_carpool_choose_car_id" example:"1"`
	RefCarpoolChooseCar   VmsRefCarpoolChooseCar `gorm:"foreignKey:RefCarpoolChooseCarID;references:RefCarpoolChooseCarID" json:"ref_carpool_choose_car"`
}

func (VmsMasCarpoolCarBooking) TableName() string {
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
	// List of random URLs
	imageUrls := []string{
		"http://pntdev.ddns.net:28089/VMS_PLUS/PIX/cartype/EV.svg",
		"http://pntdev.ddns.net:28089/VMS_PLUS/PIX/cartype/VAN.svg",
		"http://pntdev.ddns.net:28089/VMS_PLUS/PIX/cartype/SUV.svg",
	}

	// Seed the random generator
	for i := range vehicle_types {
		vehicle_types[i].VehicleTypeImage = imageUrls[i%len(imageUrls)]
	}

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
	Seat                             int                     `gorm:"column:Seat" json:"seat"`
	RefFuelType                      VmsRefFuelType          `gorm:"foreignKey:RefFuelTypeID;references:RefFuelTypeID" json:"ref_fuel_type"`
	VehicleGetDate                   time.Time               `gorm:"column:vehicle_get_date" json:"vehicle_get_date"`
	Age                              int                     `json:"age"`
	VehicleDepartment                VmsMasVehicleDepartment `gorm:"foreignKey:MasVehicleUID;references:MasVehicleUID" json:"vehicle_department"`
	IsAdminChooseDriver              bool                    `json:"is_admin_choose_driver"`
}

func (VmsMasVehicle) TableName() string {
	return "vms_mas_vehicle"
}
func (v *VmsMasVehicle) CalculateAge() int {
	now := time.Now()
	// Subtract the registration year from the current year
	age := now.Year() - v.VehicleGetDate.Year()

	// Adjust if the current date is before the registration date in the year
	if now.YearDay() < v.VehicleGetDate.YearDay() {
		age--
	}
	return age
}

type VmsMasVehicleDepartment struct {
	MasVehicleUID                    string     `gorm:"column:mas_vehicle_uid;primaryKey" json:"-"`
	County                           string     `gorm:"column:county" json:"county"`
	VehicleGetDate                   time.Time  `gorm:"column:vehicle_get_date" json:"vehicle_get_date"`
	VehiclePeaID                     string     `gorm:"column:vehicle_pea_id" json:"vehicle_pea_id"`
	VehicleLicensePlate              string     `gorm:"column:vehicle_license_plate" json:"vehicle_license_plate"`
	VehicleLicensePlateProvinceShort string     `gorm:"column:vehicle_license_plate_province_short" json:"vehicle_license_plate_province_short"`
	VehicleLicensePlateProvinceFull  string     `gorm:"column:vehicle_license_plate_province_full" json:"vehicle_license_plate_province_full"`
	VehicleAssetNo                   string     `gorm:"column:vehicle_asset_no" json:"vehicle_asset_no"`
	AssetClass                       string     `gorm:"column:asset_class" json:"asset_class"`
	AssetSubcategory                 string     `gorm:"column:asset_subcategory" json:"asset_subcategory"`
	RefPeaOfficialVehicleTypeCode    int        `gorm:"column:ref_pea_official_vehicle_type_code" json:"ref_pea_official_vehicle_type_code"`
	VehicleCondition                 int        `gorm:"column:vehicle_condition" json:"vehicle_condition"`
	VehicleMileage                   int        `gorm:"column:vehicle_mileage" json:"vehicle_mileage"`
	VehicleOwnerDeptSap              string     `gorm:"column:vehicle_owner_dept_sap" json:"vehicle_owner_dept_sap"`
	VehicleCostCenter                string     `gorm:"column:vehicle_cost_center" json:"vehicle_cost_center"`
	OwnerDeptName                    string     `gorm:"column:owner_dept_name" json:"owner_dept_name"`
	VehicleImg                       string     `gorm:"column:vehicle_img" json:"vehicle_img"`
	VehicleUserEmpID                 string     `gorm:"column:vehicle_user_emp_id" json:"vehicle_user_emp_id"`
	VehicleUserEmpName               string     `gorm:"column:vehicle_user_emp_name" json:"vehicle_user_emp_name"`
	VehicleAdminEmpID                string     `gorm:"column:vehicle_admin_emp_id" json:"vehicle_admin_emp_id"`
	VehicleAdminEmpName              string     `gorm:"column:vehicle_admin_emp_name" json:"vehicle_admin_emp_name"`
	ParkingPlace                     string     `gorm:"column:parking_place" json:"parking_place"`
	FleetCardNo                      string     `gorm:"column:fleet_card_no" json:"fleet_card_no"`
	IsInCarpool                      []byte     `gorm:"column:is_in_carpool;type:bit(1)" json:"is_in_carpool"`
	Remark                           string     `gorm:"column:remark" json:"remark"`
	RefVehicleStatusCode             int        `gorm:"column:ref_vehicle_status_code" json:"ref_vehicle_status_code"`
	RefOtherUseCode                  *int       `gorm:"column:ref_other_use_code" json:"ref_other_use_code"`
	VehicleUser                      MasUserEmp `gorm:"foreignKey:VehicleUserEmpID;references:EmpID" json:"vehicle_user"`
}

func (VmsMasVehicleDepartment) TableName() string {
	return "vms_mas_vehicle_department"
}
