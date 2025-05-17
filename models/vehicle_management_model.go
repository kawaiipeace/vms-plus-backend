package models

import (
	"time"
)

// VmsMasVehicleManagementList
type VmsMasVehicleManagementList struct {
	MasVihicleUID             string              `gorm:"primaryKey;column:mas_vehicle_uid" json:"mas_vehicle_uid"`
	VehicleLicensePlate       string              `gorm:"column:vehicle_license_plate" json:"vehicle_license_plate"`
	VehicleBrandName          string              `gorm:"column:vehicle_brand_name" json:"vehicle_brand_name"`
	VehicleModelName          string              `gorm:"column:vehicle_model_name" json:"vehicle_model_name"`
	RefVehicleTypeCode        string              `gorm:"column:ref_vehicle_type_code" json:"ref_vehicle_type_code"`
	RefVehicleTypeName        string              `gorm:"column:ref_vehicle_type_name" json:"ref_vehicle_type_name"`
	VehicleOwnerDeptSAP       string              `gorm:"column:vehicle_owner_dept_short" json:"vehicle_owner_dept_short"`
	FleetCardNo               string              `gorm:"column:fleet_card_no" json:"fleet_card_no"`
	IsTaxCredit               bool                `gorm:"column:is_tax_credit" json:"is_tax_credit"`
	VehicleMileage            float64             `gorm:"column:vehicle_mileage" json:"vehicle_mileage"`
	VehicleGetDate            time.Time           `gorm:"column:vehicle_get_date" json:"vehicle_get_date"` // Changed to time.Time
	RefVehicleStatusCode      int                 `gorm:"column:ref_vehicle_status_code" json:"ref_vehicle_status_code"`
	RefVehicleStatus          VmsRefVehicleStatus `gorm:"foreignKey:RefVehicleStatusCode;references:RefVehicleStatusCode" json:"vms_ref_vehicle_status"`
	RefVehicleStatusShortName string              `gorm:"column:ref_vehicle_status_short_name" json:"-"`
	VehicleCarpoolName        string              `gorm:"column:vehicle_carpool_name" json:"vehicle_carpool_name"`
	IsAcictive                string              `gorm:"column:is_active" json:"is_active"`
	RefFuelTypeID             int                 `gorm:"column:ref_fuel_type_id" json:"ref_fuel_type_id"`
	RefFuelType               VmsRefFuelType      `gorm:"foreignKey:RefFuelTypeID;references:RefFuelTypeID" json:"vms_ref_fuel_type"`
	Age                       string              `json:"age"`
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

type VehicleTimeLine struct {
	MasVehicleUID                    string              `gorm:"column:mas_vehicle_uid" json:"mas_vehicle_uid"`
	VehicleLicensePlate              string              `gorm:"column:vehicle_license_plate" json:"vehicle_license_plate"`
	VehicleLicensePlateProvinceShort string              `gorm:"column:vehicle_license_plate_province_short" json:"vehicle_license_plate_province_short"`
	VehicleLicensePlateProvinceFull  string              `gorm:"column:vehicle_license_plate_province_full" json:"vehicle_license_plate_province_full"`
	VehicleDeptName                  string              `gorm:"column:vehicle_dept_name" json:"vehicle_dept_name"`
	CarpoolName                      string              `gorm:"column:vehicle_carpool_name" json:"vehicle_carpool_name"`
	VehicleCarTypeDetail             string              `gorm:"column:vehicle_car_type_detail" json:"vehicle_car_type_detail"`
	VehicleMileage                   string              `gorm:"column:vehicle_mileage" json:"vehicle_mileage"`
	VehicleTrnRequests               []VehicleTrnRequest `gorm:"foreignKey:MasVehicleUID;references:MasVehicleUID" json:"vehicle_t_requests"`
}

type VehicleTrnRequest struct {
	MasVehicleUID          string             `gorm:"column:mas_vehicle_uid" json:"mas_vehicle_uid"`
	TrnRequestUID          string             `gorm:"column:trn_request_uid" json:"trn_request_uid"`
	RequestNo              string             `gorm:"column:request_no" json:"request_no"`
	ReserveStartDatetime   string             `gorm:"column:reserve_start_datetime" json:"start_datetime"`
	ReserveEndDatetime     string             `gorm:"column:reserve_end_datetime" json:"end_datetime"`
	RefRequestStatusCode   string             `gorm:"column:ref_request_status_code" json:"ref_request_status_code"`
	RefRequestStatusName   string             `json:"ref_request_status_name"`
	VehicleUserEmpID       string             `gorm:"column:vehicle_user_emp_id" json:"vehicle_user_emp_id" example:"990001"`
	VehicleUserEmpName     string             `gorm:"column:vehicle_user_emp_name" json:"vehicle_user_emp_name"`
	VehicleUserDeptSAP     string             `gorm:"column:vehicle_user_dept_sap" json:"vehicle_user_dept_sap"`
	VehicleUserDeskPhone   string             `gorm:"column:vehicle_user_desk_phone" json:"car_user_internal_contact_number" example:"1122"`
	VehicleUserMobilePhone string             `gorm:"column:vehicle_user_mobile_phone" json:"car_user_mobile_contact_number" example:"0987654321"`
	IsPEAEmployeeDriver    string             `gorm:"column:is_pea_employee_driver" json:"is_pea_employee_driver" example:"1"`
	MasCarpoolDriverUID    string             `gorm:"column:mas_carpool_driver_uid;type:uuid" json:"mas_carpool_driver_uid"`
	DriverEmpID            string             `gorm:"column:driver_emp_id" json:"driver_emp_id" example:"700001"`
	DriverEmpName          string             `gorm:"column:driver_emp_name" json:"driver_emp_name" example:"John Doe"`
	DriverDeptSAP          string             `gorm:"column:driver_emp_dept_sap" json:"driver_emp_dept_sap" example:"DPT001"`
	DriverInternalContact  string             `gorm:"column:driver_internal_contact_number" json:"driver_internal_contact_number" example:"1234567890"`
	DriverMobileContact    string             `gorm:"column:driver_mobile_contact_number" json:"driver_mobile_contact_number" example:"0987654321"`
	MasDriver              VmsMasDriverShort  `gorm:"foreignKey:MasCarpoolDriverUID;references:MasDriverUID" json:"driver"`
	TripDetails            []VmsTrnTripDetail `gorm:"foreignKey:TrnRequestUID;references:TrnRequestUID" json:"trip_details"`
}

func (VehicleTrnRequest) TableName() string {
	return "public.vms_trn_request"
}

type VehicleReportTripDetail struct {
	MasVehicleUID                    string    `gorm:"column:mas_vehicle_uid" json:"mas_vehicle_uid"`
	VehicleLicensePlate              string    `gorm:"column:vehicle_license_plate" json:"vehicle_license_plate"`
	VehicleLicensePlateProvinceShort string    `gorm:"column:vehicle_license_plate_province_short" json:"vehicle_license_plate_province_short"`
	VehicleLicensePlateProvinceFull  string    `gorm:"column:vehicle_license_plate_province_full" json:"vehicle_license_plate_province_full"`
	VehicleDeptName                  string    `gorm:"column:vehicle_dept_name" json:"vehicle_dept_name"`
	CarpoolName                      string    `gorm:"column:vehicle_carpool_name" json:"vehicle_carpool_name"`
	VehicleCarTypeDetail             string    `gorm:"column:vehicle_car_type_detail" json:"vehicle_car_type_detail"`
	VehicleMileage                   string    `gorm:"column:vehicle_mileage" json:"vehicle_mileage"`
	TripStartDatetime                time.Time `gorm:"column:trip_start_datetime" json:"trip_start_datetime" example:"2025-03-26T08:00:00Z"`
	TripEndDatetime                  time.Time `gorm:"column:trip_end_datetime" json:"trip_end_datetime" example:"2025-03-26T10:00:00Z"`
	TripDeparturePlace               string    `gorm:"column:trip_departure_place" json:"trip_departure_place" example:"Changi Airport"`
	TripDestinationPlace             string    `gorm:"column:trip_destination_place" json:"trip_destination_place" example:"Marina Bay Sands"`
	TripStartMiles                   int       `gorm:"column:trip_start_miles" json:"trip_start_miles" example:"5000"`
	TripEndMiles                     int       `gorm:"column:trip_end_miles" json:"trip_end_miles" example:"5050"`
	TripDetail                       string    `gorm:"column:trip_detail" json:"trip_detail" example:"Routine transport between airport and hotel."`
}

type VehicleReportAddFuel struct {
	MasVehicleUID                    string    `gorm:"column:mas_vehicle_uid" json:"mas_vehicle_uid"`
	VehicleLicensePlate              string    `gorm:"column:vehicle_license_plate" json:"vehicle_license_plate"`
	VehicleLicensePlateProvinceShort string    `gorm:"column:vehicle_license_plate_province_short" json:"vehicle_license_plate_province_short"`
	VehicleLicensePlateProvinceFull  string    `gorm:"column:vehicle_license_plate_province_full" json:"vehicle_license_plate_province_full"`
	VehicleDeptName                  string    `gorm:"column:vehicle_dept_name" json:"vehicle_dept_name"`
	CarpoolName                      string    `gorm:"column:vehicle_carpool_name" json:"vehicle_carpool_name"`
	VehicleCarTypeDetail             string    `gorm:"column:vehicle_car_type_detail" json:"vehicle_car_type_detail"`
	Mile                             int       `gorm:"column:mile" json:"mile" example:"12000"`
	TaxInvoiceDate                   time.Time `gorm:"column:tax_invoice_date;type:timestamp" json:"tax_invoice_date" example:"2025-03-26T08:00:00Z"`
	TaxInvoiceNo                     string    `gorm:"column:tax_invoice_no;type:varchar(20)" json:"tax_invoice_no" example:"INV1234567890"`
	PricePerLiter                    float64   `gorm:"column:price_per_liter;type:numeric(10,2)" json:"price_per_liter" example:"35.50"`
	SumLiter                         float64   `gorm:"column:sum_liter;type:numeric(10,2)" json:"sum_liter" example:"50.00"`
	SumPrice                         float64   `gorm:"column:sum_price;type:numeric(10,2)" json:"sum_price" example:"1872.50"`
	ReceiptImg                       string    `gorm:"column:receipt_img;type:varchar(100)" json:"receipt_img" example:"http://vms.pea.co.th/receipt.jpg"`
	AddFuelDateTime                  time.Time `gorm:"column:add_fuel_date_time" json:"add_fuel_date_time" example:"2025-03-26T08:00:00Z"`
	RefCostTypeCode                  int       `gorm:"column:ref_cost_type_code" json:"ref_cost_type_code" example:"1"`
	RefCostType                      string    `gorm:"column:ref_cost_type" json:"ref_cost_type" example:"Fuel"`
	RefOilStationBrand               string    `gorm:"column:ref_oil_station_brand" json:"ref_oil_station_brand" example:"PTT"`
	RefFuelType                      string    `gorm:"column:ref_fuel_type" json:"ref_fuel_type" example:"Gasoline"`
	RefPaymentType                   string    `gorm:"column:ref_payment_type" json:"ref_payment_type" example:"Credit Card"`
}
