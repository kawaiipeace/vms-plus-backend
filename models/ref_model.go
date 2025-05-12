package models

//VmsRefRequestStatus
type VmsRefRequestStatus struct {
	RefRequestStatusCode string `gorm:"column:ref_request_status_code" json:"ref_request_status_code"`
	RefRequestStatusDesc string `gorm:"column:ref_request_status_desc" json:"ref_request_status_desc"`
}

func (VmsRefRequestStatus) TableName() string {
	return "vms_ref_request_status"
}

// VmsRefFuelType
type VmsRefFuelType struct {
	RefFuelTypeID     int    `gorm:"primaryKey;column:ref_fuel_type_id" json:"ref_fuel_type_id"`
	RefFuelTypeNameTh string `gorm:"column:ref_fuel_type_name_th" json:"ref_fuel_type_name_th"`
	RefFuelTypeNameEn string `gorm:"column:ref_fuel_type_name_en" json:"ref_fuel_type_name_en"`
}

func (VmsRefFuelType) TableName() string {
	return "vms_ref_fuel_type"
}

// VmsRefCostType
type VmsRefCostType struct {
	RefCostTypeCode string `gorm:"column:ref_cost_type_code" json:"ref_cost_type_code"`
	RefCostTypeName string `gorm:"column:ref_cost_type_name" json:"ref_cost_type_name"`
	RefCostNo       string `gorm:"column:ref_cost_no" json:"ref_cost_no"`
}

func (VmsRefCostType) TableName() string {
	return "vms_ref_cost_type"
}

// VmsRefOilStationBrand
type VmsRefOilStationBrand struct {
	RefOilStationBrandId       int    `gorm:"primaryKey;column:ref_oil_station_brand_id" json:"ref_oil_station_brand_id"`
	RefOilStationBrandNameTh   string `gorm:"column:ref_oil_station_brand_name_th" json:"ref_oil_station_brand_name_th"`
	RefOilStationBrandNameEn   string `gorm:"column:ref_oil_station_brand_name_en" json:"ref_oil_station_brand_name_en"`
	RefOilStationBrandNameFull string `gorm:"column:ref_oil_station_brand_name_full" json:"ref_oil_station_brand_name_full"`
	RefOilStationBrandImg      string `gorm:"column:ref_oil_station_brand_img" json:"ref_oil_station_brand_img"`
}

func (VmsRefOilStationBrand) TableName() string {
	return "vms_ref_oil_station_brand"
}

// VmsRefVehicleImgSide
type VmsRefVehicleImgSide struct {
	RefVehicleImgSideCode int    `gorm:"primaryKey;column:ref_vehicle_img_side_code" json:"ref_vehicle_img_side_code"`
	VehicleImgDescription string `gorm:"column:vehicle_img_description" json:"vehicle_img_description"`
}

func (VmsRefVehicleImgSide) TableName() string {
	return "vms_ref_vehicle_img_side"
}

// VmsRefPaymentType
type VmsRefPaymentType struct {
	RefPaymentTypeCode int    `gorm:"column:ref_payment_type_code;primarykey" json:"ref_payment_type_code"`
	RefPaymentTypeName string `gorm:"column:ref_payment_type_name" json:"ref_payment_type_name"`
}

func (VmsRefPaymentType) TableName() string {
	return "vms_ref_payment_type"
}

// VmsRefOtherUse
type VmsRefOtherUse struct {
	RefOtherUseCode int    `gorm:"column:ref_other_use_code;primarykey" json:"ref_other_use_code"`
	RefOtherUseDesc string `gorm:"column:ref_other_use_desc" json:"ref_other_use_desc"`
}

func (VmsRefOtherUse) TableName() string {
	return "vms_ref_other_use"
}

//VmsRefDriverLicenseType
type VmsRefDriverLicenseType struct {
	RefDriverLicenseTypeCode string `gorm:"column:ref_driver_license_type_code;primaryKey;type:varchar(2)" json:"ref_driver_license_type_code"`
	RefDriverLicenseTypeName string `gorm:"column:ref_driver_license_type_name;type:varchar(50)" json:"ref_driver_license_type_name"`
	RefDriverLicenseTypeDesc string `gorm:"column:ref_driver_license_type_desc;type:varchar(350)" json:"ref_driver_license_type_desc"`
}

func (VmsRefDriverLicenseType) TableName() string {
	return "vms_ref_driver_license_type"
}

//VmsRefDriverCertificateType
type VmsRefDriverCertificateType struct {
	RefDriverCertificateTypeCode int    `gorm:"column:ref_driver_certificate_type_code;primaryKey" json:"ref_driver_certificate_type_code"`
	RefDriverCertificateTypeName string `gorm:"column:ref_driver_certificate_type_name" json:"ref_driver_certificate_type_name"`
	RefDriverCertificateTypeDesc string `gorm:"column:ref_driver_certificate_type_desc" json:"ref_driver_certificate_type_desc"`
}

func (VmsRefDriverCertificateType) TableName() string {
	return "vms_ref_driver_certificate_type"
}

// VmsRefCarpoolChooseCar
type VmsRefCarpoolChooseCar struct {
	RefCarpoolChooseCarID int    `gorm:"primaryKey;column:ref_carpool_choose_car_id" json:"ref_carpool_choose_car_id"`
	TypeOfChooseCar       string `gorm:"column:type_of_choose_car" json:"type_of_choose_car"`
}

func (VmsRefCarpoolChooseCar) TableName() string {
	return "vms_ref_carpool_choose_car"
}

// VmsRefCarpoolChooseDriver
type VmsRefCarpoolChooseDriver struct {
	RefCarpoolChooseDriverID int    `gorm:"primaryKey;column:ref_carpool_choose_driver_id" json:"ref_carpool_choose_driver_id"`
	TypeOfChooseDriver       string `gorm:"column:type_of_choose_Driver" json:"type_of_choose_driver"`
}

func (VmsRefCarpoolChooseDriver) TableName() string {
	return "vms_ref_carpool_choose_driver"
}

// RefVehicleKeyType
type VmsRefVehicleKeyType struct {
	RefVehicleKeyTypeCode string `gorm:"column:ref_vehicle_key_type_code;primaryKey" json:"ref_vehicle_key_type_code"`
	RefVehicleKeyTypeName string `gorm:"column:ref_vehicle_key_type_name" json:"ref_vehicle_key_type_name"`
}

func (VmsRefVehicleKeyType) TableName() string {
	return "vms_ref_vehicle_key_type"
}

// VmsRefLeaveTimeType
type VmsRefLeaveTimeType struct {
	LeaveTimeTypeCode int    `gorm:"primaryKey;column:leave_time_type_code" json:"leave_time_type_code"`
	LeaveTimeTypeName string `gorm:"column:leave_time_type_name" json:"leave_time_type_name"`
}

func (VmsRefLeaveTimeType) TableName() string {
	return "vms_ref_leave_time_type"
}

// VmsRefVehicleStatus
type VmsRefVehicleStatus struct {
	RefVehicleStatusCode int    `gorm:"primaryKey;column:ref_vehicle_status_code" json:"ref_vehicle_status_code"`
	RefVehicleStatusName string `gorm:"column:ref_vehicle_status_short_name" json:"ref_vehicle_status_name"`
}

func (VmsRefVehicleStatus) TableName() string {
	return "vms_ref_vehicle_status"
}

//VmsRefTripType
type VmsRefTripType struct {
	RefTripTypeCode int    `gorm:"column:ref_trip_type_code;primaryKey" json:"ref_trip_type_code"`
	RefTripTypeName string `gorm:"column:ref_trip_type_name" json:"ref_trip_type_name"`
}

func (VmsRefTripType) TableName() string {
	return "vms_ref_trip_type"
}
