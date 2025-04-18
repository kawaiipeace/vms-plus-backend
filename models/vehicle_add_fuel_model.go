package models

import "time"

type VmsTrnAddFuel_List struct {
	TrnAddFuelUid        string    `gorm:"column:trn_add_fuel_uid" json:"trn_add_fuel_uid" example:"123e4567-e89b-12d3-a456-426614174000"`
	TrnRequestUID        string    `gorm:"column:trn_request_uid" json:"trn_request_uid" example:"456e4567-e89b-12d3-a456-426614174001"`
	TripStartDatetime    time.Time `gorm:"column:trip_start_datetime" json:"trip_start_datetime" example:"2025-03-26T08:00:00"`
	TripEndDatetime      time.Time `gorm:"column:trip_end_datetime" json:"trip_end_datetime" example:"2025-03-26T10:00:00"`
	TripDeparturePlace   string    `gorm:"column:trip_departure_place" json:"trip_departure_place" example:"Changi Airport"`
	TripDestinationPlace string    `gorm:"column:trip_destination_place" json:"trip_destination_place" example:"Marina Bay Sands"`
	TripStartMiles       int       `gorm:"column:trip_start_miles" json:"trip_start_miles" example:"5000"`
	TripEndMiles         int       `gorm:"column:trip_end_miles" json:"trip_end_miles" example:"5050"`
	TripDetail           string    `gorm:"column:trip_detail" json:"trip_detail" example:"Routine transport between airport and hotel."`
}

type VmsTrnAddFuel_Request struct {
	TrnRequestUID  string    `gorm:"column:trn_request_uid" json:"trn_request_uid" example:"456e4567-e89b-12d3-a456-426614174001"`
	PricePerLiter  float64   `gorm:"column:price_per_liter;type:numeric(10,2)" json:"price_per_liter"`
	SumLiter       float64   `gorm:"column:sum_liter;type:numeric(10,2)" json:"sum_liter"`
	BeforeVatPrice float64   `gorm:"column:before_vat_price;type:numeric(10,2)" json:"before_vat_price"`
	Vat            float64   `gorm:"column:vat;type:numeric(10,2)" json:"vat"`
	SumPrice       float64   `gorm:"column:sum_price;type:numeric(10,2)" json:"sum_price"`
	ReceiptImg     string    `gorm:"column:receipt_img;type:varchar(100)" json:"receipt_img"`
	TaxInvoiceNo   string    `gorm:"column:tax_invoice_no;type:varchar(20)" json:"tax_invoice_no"`
	TaxInvoiceDate time.Time `gorm:"column:tax_invoice_date;type:timestamp" json:"tax_invoice_date"`
}

type VmsTrnAddFuel struct {
	TrnAddFuelUID                    string `gorm:"column:trn_add_fuel_uid" json:"trn_add_fuel_uid" example:"123e4567-e89b-12d3-a456-426614174000"`
	MasVehicleUID                    string `gorm:"column:mas_vehicle_uid" json:"mas_vehicle_uid" example:"789e4567-e89b-12d3-a456-426614174002"`
	VehicleLicensePlate              string `gorm:"column:vehicle_license_plate" json:"vehicle_license_plate" example:"SGP1234"`
	VehicleLicensePlateProvinceShort string `gorm:"column:vehicle_license_plate_province_short" json:"vehicle_license_plate_province_short" example:"SG"`
	VehicleLicensePlateProvinceFull  string `gorm:"column:vehicle_license_plate_province_full" json:"vehicle_license_plate_province_full" example:"Singapore"`
	VmsTrnAddFuel_Request
	CreatedAt               time.Time `gorm:"column:created_at" json:"-"`
	CreatedBy               string    `gorm:"column:created_by" json:"-"`
	UpdatedAt               time.Time `gorm:"column:updated_at" json:"-"`
	UpdatedBy               string    `gorm:"column:updated_by" json:"-"`
	MasVehicleDepartmentUID string    `gorm:"column:mas_vehicle_department_uid" json:"mas_vehicle_department_uid" example:"abc12345-6789-1234-5678-abcdef012345"`
}

func (VmsTrnAddFuel) TableName() string {
	return "public.vms_trn_add_fuel"
}
