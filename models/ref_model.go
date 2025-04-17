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
