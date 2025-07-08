package models

import "time"

//VmsTrnReceivedVehicle
type VmsTrnReceivedVehicle struct {
	TrnRequestUID         string       `gorm:"column:trn_request_uid;primaryKey" json:"trn_request_uid" example:"0b07440c-ab04-49d0-8730-d62ce0a9bab9"`
	PickupDatetime        TimeWithZone `gorm:"column:pickup_datetime" json:"pickup_datetime" example:"2025-03-26T14:30:00+07:00"`
	MileStart             int          `gorm:"column:mile_start" json:"mile_start" example:"10000"`
	FuelStart             int          `gorm:"column:fuel_start" json:"fuel_start" example:"50"`
	ReceivedVehicleRemark string       `gorm:"column:received_vehicle_remark" json:"received_vehicle_remark" example:"Minor scratch on bumper"`

	VehicleImages               []VehicleImageReceived `gorm:"foreignKey:TrnRequestUID;references:TrnRequestUID" json:"vehicle_images"`
	ReceivedVehicleEmpID        string                 `gorm:"column:received_vehicle_emp_id" json:"-"`
	ReceivedVehicleEmpName      string                 `gorm:"column:received_vehicle_emp_name" json:"-"`
	ReceivedVehicleDeptSAP      string                 `gorm:"column:received_vehicle_dept_sap" json:"-"`
	ReceivedVehicleDeptSAPShort string                 `gorm:"column:received_vehicle_dept_name_short" json:"-"`
	ReceivedVehicleDeptSAPFull  string                 `gorm:"column:received_vehicle_dept_name_full" json:"-"`
	RefRequestStatusCode        string                 `gorm:"column:ref_request_status_code" json:"-"`
	UpdatedAt                   time.Time              `gorm:"column:updated_at" json:"-"`
	UpdatedBy                   string                 `gorm:"column:updated_by" json:"-"`
}

func (VmsTrnReceivedVehicle) TableName() string {
	return "public.vms_trn_request"
}

// VehicleImageReceived
type VehicleImageReceived struct {
	TrnVehicleImgReceivedUID string    `gorm:"column:trn_vehicle_img_received_uid;primaryKey" json:"-"`
	TrnRequestUID            string    `gorm:"column:trn_request_uid;" json:"-"`
	RefVehicleImgSideCode    int       `gorm:"column:ref_vehicle_img_side_code" json:"ref_vehicle_img_side_code" example:"1"`
	VehicleImgFile           string    `gorm:"column:vehicle_img_file" json:"vehicle_img_file" example:"http://vms.pea.co.th/side_image.jpg"`
	CreatedAt                time.Time `gorm:"column:created_at" json:"-"`
	CreatedBy                string    `gorm:"column:created_by" json:"-"`
	UpdatedAt                time.Time `gorm:"column:updated_at" json:"-"`
	UpdatedBy                string    `gorm:"column:updated_by" json:"-"`
	IsDeleted                string    `gorm:"column:is_deleted" json:"-"`
}

func (VehicleImageReceived) TableName() string {
	return "public.vms_trn_vehicle_img_received"
}

//VmsTrnTravelCard
type VmsTrnTravelCard struct {
	TrnRequestUID string       `gorm:"column:trn_request_uid;primaryKey;" json:"trn_request_uid" example:"a7de5318-1e05-4511-abe7-8c1c6374ab29"`
	StartDateTime TimeWithZone `gorm:"column:start_datetime" json:"start_datetime" example:"2025-02-16T08:30:00+07:00"`
	EndDateTime   TimeWithZone `gorm:"column:end_datetime" json:"end_datetime" example:"2025-02-16T09:30:00+07:00"`

	VehicleLicensePlate              string `gorm:"column:vehicle_license_plate" json:"vehicle_license_plate"`
	VehicleLicensePlateProvinceShort string `gorm:"column:vehicle_license_plate_province_short" json:"vehicle_license_plate_province_short"`
	VehicleLicensePlateProvinceFull  string `gorm:"column:vehicle_license_plate_province_full" json:"vehicle_license_plate_province_full"`
	WorkPlace                        string `gorm:"column:work_place" json:"work_place" example:"Office"`

	VehicleUserEmpID        string `gorm:"column:vehicle_user_emp_id" json:"vehicle_user_emp_id" example:"700001"`
	VehicleUserEmpName      string `gorm:"column:vehicle_user_emp_name" json:"vehicle_user_emp_name" example:"John Smith"`
	VehicleUserDeptSAP      string `gorm:"column:vehicle_user_dept_sap" json:"vehicle_user_dept_sap" example:"HR"`
	VehicleUserDeptSAPShort string `gorm:"column:vehicle_user_dept_name_short" json:"vehicle_user_dept_sap_short" example:"Finance"`
	VehicleUserDeptSAPFull  string `gorm:"column:vehicle_user_dept_name_full" json:"vehicle_user_dept_sap_full" example:"Finance"`
	VehicleUserImageURL     string `gorm:"-" json:"vehicle_user_image_url"`

	ApprovedRequestEmpID        string `gorm:"column:approved_request_emp_id" json:"approved_request_emp_id" example:"EMP67890"`
	ApprovedRequestEmpName      string `gorm:"column:approved_request_emp_name" json:"approved_request_emp_name" example:"Jane Doe"`
	ApprovedRequestDeptSAP      string `gorm:"column:approved_request_dept_sap" json:"approved_request_dept_sap" example:"Finance"`
	ApprovedRequestDeptSAPShort string `gorm:"column:approved_request_dept_sap_short" json:"approved_request_dept_sap_short" example:"Finance"`
	ApprovedRequestDeptSAPFull  string `gorm:"column:approved_request_dept_sap_full" json:"approved_request_dept_sap_full" example:"Finance"`
}

func (VmsTrnTravelCard) TableName() string {
	return "public.vms_trn_request"
}
