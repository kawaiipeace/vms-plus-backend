package models

type VmsMasVehicle_List struct {
	MasVehicleUID       string `gorm:"primaryKey;column:mas_vehicle_uid" json:"mas_vehicle_uid"`
	VehicleLicensePlate string `gorm:"column:vehicle_license_plate;uniqueIndex" json:"vehicle_license_plate"`
	VehicleBrandName    string `gorm:"column:vehicle_brand_name" json:"vehicle_brand_name"`
	VehicleModelName    string `gorm:"column:vehicle_model_name" json:"vehicle_model_name"`
	CarType             string `gorm:"column:car_type" json:"car_type"`
	VehicleOwnerDeptSAP string `gorm:"column:vehicle_owner_dept_sap" json:"vehicle_owner_dept_sap"`
	VehicleImg          string `gorm:"column:vehicle_img" json:"vehicle_img"` // Store image URL or file path
	Seat                int    `gorm:"column:Seat" json:"seat"`
}

func (VmsMasVehicle_List) TableName() string {
	return "vms_mas_vehicle"
}

type VmsRefCategory struct {
	RefVehicleCategoryCode string `gorm:"column:ref_vehicle_category_code" json:"ref_vehicle_category_code"`
	RefVehicleCategoryName string `gorm:"column:ref_vehicle_category_name" json:"ref_vehicle_category_name"`
	AvailableUnits         int    `gorm:"column:available_units" json:"available_units"`
}

func (VmsRefCategory) TableName() string {
	return "vms_ref_vehicle_category"
}

type VmsMasVehicle struct {
	MasVehicleUID         string `gorm:"primaryKey;column:mas_vehicle_uid;type:uuid" json:"mas_vehicle_uid"`
	VehicleBrandName      string `gorm:"column:vehicle_brand_name" json:"vehicle_brand_name"`
	VehicleModelName      string `gorm:"column:vehicle_model_name" json:"vehicle_model_name"`
	VehicleLicensePlate   string `gorm:"column:vehicle_license_plate" json:"vehicle_license_plate"`
	VehicleImg            string `gorm:"column:vehicle_img" json:"vehicle_img"`
	CarType               string `gorm:"column:CarType" json:"CarType"`
	VehicleOwnerDeptSap   string `gorm:"column:vehicle_owner_dept_sap" json:"vehicle_owner_dept_sap"`
	IsHasFleetCard        byte   `gorm:"column:is_has_fleet_card" json:"is_has_fleet_card"`
	VehicleGear           string `gorm:"column:vehicle_gear" json:"vehicle_gear"`
	RefVehicleSubtypeCode int    `gorm:"column:ref_vehicle_subtype_code" json:"ref_vehicle_subtype_code"`
	VehicleUserEmpID      string `gorm:"column:vehicle_user_emp_id" json:"vehicle_user_emp_id"`
	RefFuelTypeID         int    `gorm:"column:ref_fuel_type_id" json:"ref_fuel_type_id"`
	Seat                  int    `gorm:"column:Seat" json:"seat"`
	// Foreign key relation (Optional if you want to use GORM's automatic foreign key mapping)
	RefFuelType VmsRefFuelType `gorm:"foreignKey:RefFuelTypeID;references:RefFuelTypeID" json:"ref_fuel_type"`
}

func (VmsMasVehicle) TableName() string {
	return "vms_mas_vehicle"
}
