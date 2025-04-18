package models

import "time"

type VmsTrnTripDetail_List struct {
	TrnTripDetailUID     string    `gorm:"column:trn_trip_detail_uid" json:"trn_trip_detail_uid" example:"123e4567-e89b-12d3-a456-426614174000"`
	TrnRequestUID        string    `gorm:"column:trn_request_uid" json:"trn_request_uid" example:"456e4567-e89b-12d3-a456-426614174001"`
	TripStartDatetime    time.Time `gorm:"column:trip_start_datetime" json:"trip_start_datetime" example:"2025-03-26T08:00:00"`
	TripEndDatetime      time.Time `gorm:"column:trip_end_datetime" json:"trip_end_datetime" example:"2025-03-26T10:00:00"`
	TripDeparturePlace   string    `gorm:"column:trip_departure_place" json:"trip_departure_place" example:"Changi Airport"`
	TripDestinationPlace string    `gorm:"column:trip_destination_place" json:"trip_destination_place" example:"Marina Bay Sands"`
	TripStartMiles       int       `gorm:"column:trip_start_miles" json:"trip_start_miles" example:"5000"`
	TripEndMiles         int       `gorm:"column:trip_end_miles" json:"trip_end_miles" example:"5050"`
	TripDetail           string    `gorm:"column:trip_detail" json:"trip_detail" example:"Routine transport between airport and hotel."`
}

type VmsTrnTripDetail_Request struct {
	TrnRequestUID        string    `gorm:"column:trn_request_uid" json:"trn_request_uid" example:"456e4567-e89b-12d3-a456-426614174001"`
	TripStartDatetime    time.Time `gorm:"column:trip_start_datetime" json:"trip_start_datetime" example:"2025-03-26T08:00:00"`
	TripEndDatetime      time.Time `gorm:"column:trip_end_datetime" json:"trip_end_datetime" example:"2025-03-26T10:00:00"`
	TripDeparturePlace   string    `gorm:"column:trip_departure_place" json:"trip_departure_place" example:"Changi Airport"`
	TripDestinationPlace string    `gorm:"column:trip_destination_place" json:"trip_destination_place" example:"Marina Bay Sands"`
	TripStartMiles       int       `gorm:"column:trip_start_miles" json:"trip_start_miles" example:"5000"`
	TripEndMiles         int       `gorm:"column:trip_end_miles" json:"trip_end_miles" example:"5050"`
	TripDetail           string    `gorm:"column:trip_detail" json:"trip_detail" example:"Routine transport between airport and hotel."`
}

type VmsTrnTripDetail struct {
	TrnTripDetailUID                 string `gorm:"column:trn_trip_detail_uid" json:"trn_trip_detail_uid" example:"123e4567-e89b-12d3-a456-426614174000"`
	MasVehicleUID                    string `gorm:"column:mas_vehicle_uid" json:"mas_vehicle_uid" example:"789e4567-e89b-12d3-a456-426614174002"`
	VehicleLicensePlate              string `gorm:"column:vehicle_license_plate" json:"vehicle_license_plate" example:"SGP1234"`
	VehicleLicensePlateProvinceShort string `gorm:"column:vehicle_license_plate_province_short" json:"vehicle_license_plate_province_short" example:"SG"`
	VehicleLicensePlateProvinceFull  string `gorm:"column:vehicle_license_plate_province_full" json:"vehicle_license_plate_province_full" example:"Singapore"`
	VmsTrnTripDetail_Request
	CreatedAt time.Time `gorm:"column:created_at" json:"-"`
	CreatedBy string    `gorm:"column:created_by" json:"-"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"-"`
	UpdatedBy string    `gorm:"column:updated_by" json:"-"`

	MasVehicleDepartmentUID string `gorm:"column:mas_vehicle_department_uid" json:"mas_vehicle_department_uid" example:"abc12345-6789-1234-5678-abcdef012345"`
	MasCarpoolUID           string `gorm:"column:mas_carpool_uid" json:"mas_carpool_uid" example:"xyz12345-6789-1234-5678-abcdef012345"`
	EmployeeOrDriverID      string `gorm:"column:employee_or_driver_id" json:"employee_or_driver_id" example:"driver001"`
}

func (VmsTrnTripDetail) TableName() string {
	return "public.vms_trn_trip_detail"
}
