package models

import "time"

type VmsTrnReceivedVehicle struct {
	TrnRequestUID             string    `gorm:"column:trn_request_uid;type:uuid;" json:"trn_request_uid" example:"a7de5318-1e05-4511-abe7-8c1c6374ab29"`
	PickupDatetime            time.Time `gorm:"column:pickup_datetime" json:"pickup_datetime" example:"2025-03-26T14:30:00"`
	VehicleImgOutsideFront    string    `gorm:"column:vehicle_img_outside_front" json:"vehicle_img_outside_front" example:"http://vms.pea.co.th/front_image.jpg"`
	VehicleImgOutsideBehind   string    `gorm:"column:vehicle_img_outside_behind" json:"vehicle_img_outside_behind" example:"http://vms.pea.co.th/behind_image.jpg"`
	VehicleImgOutsideLeft     string    `gorm:"column:vehicle_img_outside_left" json:"vehicle_img_outside_left" example:"http://vms.pea.co.th/left_image.jpg"`
	VehicleImgOutsideRight    string    `gorm:"column:vehicle_img_outside_right" json:"vehicle_img_outside_right" example:"http://vms.pea.co.th/right_image.jpg"`
	VehicleImgInsideFrontseat string    `gorm:"column:vehicle_img_inside_frontseat" json:"vehicle_img_inside_frontseat" example:"http://vms.pea.co.th/frontseat_image.jpg"`
	VehicleImgInsideBackseat  string    `gorm:"column:vehicle_img_inside_backseat" json:"vehicle_img_inside_backseat" example:"http://vms.pea.co.th/backseat_image.jpg"`
	MileStart                 int       `gorm:"column:mile_start" json:"mile_start" example:"10000"`
	FuelStart                 int       `gorm:"column:fuel_start" json:"fuel_start" example:"50"`
	ReceivedVehicleRemark     string    `gorm:"column:received_vehicle_remark" json:"received_vehicle_remark" example:"Minor scratch on bumper"`
}

func (VmsTrnReceivedVehicle) TableName() string {
	return "public.vms_trn_request"
}

//VmsTrnTravelCard
type VmsTrnTravelCard struct {
	TrnRequestUID string    `gorm:"column:trn_request_uid;type:uuid;" json:"trn_request_uid" example:"a7de5318-1e05-4511-abe7-8c1c6374ab29"`
	StartDateTime time.Time `gorm:"column:start_datetime" json:"start_datetime" example:"2025-02-16T08:30:00Z"`
	EndDateTime   time.Time `gorm:"column:end_datetime" json:"end_datetime" example:"2025-02-16T09:30:00Z"`

	VehicleLicensePlate              string `gorm:"column:vehicle_license_plate" json:"vehicle_license_plate"`
	VehicleLicensePlateProvinceShort string `gorm:"column:vehicle_license_plate_province_short" json:"vehicle_license_plate_province_short"`
	VehicleLicensePlateProvinceFull  string `gorm:"column:vehicle_license_plate_province_full" json:"vehicle_license_plate_province_full"`
	WorkPlace                        string `gorm:"column:work_place" json:"work_place" example:"Office"`

	VehicleUserEmpID        string `gorm:"column:vehicle_user_emp_id" json:"vehicle_user_emp_id" example:"700001"`
	VehicleUserEmpName      string `gorm:"column:vehicle_user_emp_name" json:"vehicle_user_emp_name" example:"John Smith"`
	VehicleUserDeptSAP      string `gorm:"column:vehicle_user_dept_sap" json:"vehicle_user_dept_sap" example:"HR"`
	VehicleUserDeptSAPShort string `gorm:"column:vehicle_user_dept_sap_name_short" json:"vehicle_user_dept_sap_short" example:"Finance"`
	VehicleUserDeptSAPFull  string `gorm:"column:vehicle_user_dept_sap_name_full" json:"vehicle_user_dept_sap_full" example:"Finance"`
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
