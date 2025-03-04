package models

import (
	"time"
)

type VmsTrnRequest_Request struct {
	//Step1
	VehicleUserEmpID       string    `gorm:"column:vehicle_user_emp_id" json:"vehicle_user_emp_id" example:"E12345"`
	VehicleUserEmpName     string    `gorm:"column:vehicle_user_emp_name" json:"vehicle_user_emp_name" example:"John Doe"`
	VehicleUserDeptSAP     string    `gorm:"column:vehicle_user_dept_sap" json:"vehicle_user_dept_sap" example:"DPT001"`
	CarUserInternalContact string    `gorm:"column:car_user_internal_contact_number" json:"car_user_internal_contact_number" example:"1234567890"`
	CarUserMobileContact   string    `gorm:"column:car_user_mobile_contact_number" json:"car_user_mobile_contact_number" example:"0987654321"`
	StartDatetime          time.Time `gorm:"column:start_datetime" json:"start_datetime" example:"2025-01-01T08:00:00Z"`
	EndDatetime            time.Time `gorm:"column:end_datetime" json:"end_datetime" example:"2025-01-01T10:00:00Z"`
	TripType               int       `gorm:"column:trip_type" json:"trip_type" example:"1"`
	ReservedTimeType       string    `gorm:"column:reserved_time_type" json:"reserved_time_type" example:"1"`
	WorkPlace              string    `gorm:"column:work_place" json:"work_place" example:"Head Office"`
	Objective              string    `gorm:"column:objective" json:"objective" example:"Business Meeting"`
	NumberOfPassengers     int       `gorm:"column:number_of_passengers" json:"number_of_passengers" example:"3"`
	Remark                 string    `gorm:"column:remark" json:"remark" example:"Urgent request"`
	ReferenceNumber        string    `gorm:"column:reference_number" json:"reference_number" example:"REF123456"`
	AttachedDocument       string    `gorm:"column:attached_document" json:"attached_document" example:"document.pdf"`
	RefCostTypeCode        int       `gorm:"column:ref_cost_type_code" json:"ref_cost_type_code" example:"101"`
	CostNo                 string    `gorm:"column:cost_no" json:"cost_no" example:"COST2024001"`

	//Step 2
	MasVehicleUID        string `gorm:"column:mas_vehicle_uid" json:"mas_vehicle_uid" example:"389b0f63-4195-4ece-bf35-0011c2f5f28c"`
	IsAdminChooseVehicle string `gorm:"column:is_admin_choose_vehicle;type:bit(1)" json:"is_admin_choose_vehicle" example:"0"`
	RequestVehicleTypeID int    `gorm:"column:requested_vehicle_type_id" json:"requested_vehicle_type_id" example:"1"`

	//Step 3
	IsDriverNeed        string `gorm:"column:is_driver_need" json:"is_driver_need" example:"1"`
	MasCarPoolDriverUID string `gorm:"column:mas_carpool_driver_uid" json:"mas_carpool_driver_uid" example:"a6c8a34b-9245-49c8-a12b-45fae77a4e7d"`
	IsPEAEmployeeDriver string `gorm:"column:is_pea_employee_driver" json:"is_pea_employee_driver" example:"1"`
	IsAdminChooseDriver string `gorm:"column:is_admin_choose_driver" json:"is_admin_choose_driver" example:"0"`

	PickupPlace    string    `gorm:"column:pickup_place" json:"pickup_place" example:"Main Office"`
	PickupDateTime time.Time `gorm:"column:pickup_datetime" json:"pickup_datetime" example:"2025-02-16T08:30:00Z"`

	//Step 4

	//
	RefRequestTypeCode int    `gorm:"column:ref_request_type_code" json:"ref_request_type_code" example:"1"`
	IsHaveSubRequest   string `gorm:"column:is_have_sub_request" json:"is_have_sub_request" example:"1"`
}

type VmsTrnRequest_Create struct {
	VmsTrnRequest_Request
	TrnRequestUID              string    `gorm:"column:trn_request_uid;type:uuid;" json:"trn_request_uid"`
	RefRequestTypeCode         int       `gorm:"column:ref_request_type_code" json:"ref_request_type_code"`
	RefRequestStatusCode       string    `gorm:"column:ref_request_status_code;default:'20'" json:"ref_request_status_code"`
	CreatedRequestDatetime     time.Time `gorm:"column:created_request_datetime;autoCreateTime" json:"created_request_datetime"`
	CreatedRequestEmpID        string    `gorm:"column:created_request_emp_id" json:"created_request_emp_id"`
	CreatedRequestEmpName      string    `gorm:"column:created_request_emp_name" json:"created_request_emp_name"`
	CreatedRequestDeptSAP      string    `gorm:"column:created_request_dept_sap" json:"created_request_dept_sap"`
	CreatedRequestDeptSAPShort string    `gorm:"column:created_request_dept_sap_name_short" json:"created_request_dept_sap_short"`
	CreatedRequestDeptSAPFull  string    `gorm:"column:created_request_dept_sap_name_full" json:"created_request_dept_sap_full"`
	LogCreate
}

func (VmsTrnRequest_Create) TableName() string {
	return "public.vms_trn_request"
}

type VmsTrnRequest_Response struct {
	TrnRequestUID              string    `gorm:"column:trn_request_uid;type:uuid;" json:"trn_request_uid"`
	VehicleUserEmpName         string    `gorm:"column:vehicle_user_emp_name" json:"vehicle_user_emp_name" example:"John Smith"`
	VehicleUserDeptSAP         string    `gorm:"column:vehicle_user_dept_sap" json:"vehicle_user_dept_sap" example:"HR"`
	VehicleUserEmpID           string    `gorm:"column:vehicle_user_emp_id" json:"vehicle_user_emp_id" example:"700001"`
	CarUserMobileContactNumber string    `gorm:"column:car_user_mobile_contact_number" json:"car_user_mobile_contact_number" example:"9876543210"`
	VehicleLicensePlate        string    `gorm:"column:vehicle_license_plate" json:"vehicle_license_plate" example:"ABC1234"`
	ApprovedRequestEmpID       string    `gorm:"column:approved_request_emp_id" json:"approved_request_emp_id" example:"EMP67890"`
	ApprovedRequestEmpName     string    `gorm:"column:approved_request_emp_name" json:"approved_request_emp_name" example:"Jane Doe"`
	ApprovedRequestDeptSAP     string    `gorm:"column:approved_request_dept_sap" json:"approved_request_dept_sap" example:"Finance"`
	StartDateTime              time.Time `gorm:"column:start_datetime" json:"start_datetime" example:"2025-02-16T08:30:00Z"`
	EndDateTime                time.Time `gorm:"column:end_datetime" json:"end_datetime" example:"2025-02-16T09:30:00Z"`
	DateRange                  string    `gorm:"column:date_range" json:"date_range" example:"2025-02-16 to 2025-02-17"`
	TripType                   int       `gorm:"column:trip_type" json:"trip_type" example:"1"`
	WorkPlace                  string    `gorm:"column:work_place" json:"work_place" example:"Office"`
	Objective                  string    `gorm:"column:objective" json:"objective" example:"Project meeting"`
	Remark                     string    `gorm:"column:remark" json:"remark" example:"Special request for parking spot"`
	NumberOfPassengers         int       `gorm:"column:number_of_passengers" json:"number_of_passengers" example:"4"`
	PickupPlace                string    `gorm:"column:pickup_place" json:"pickup_place" example:"Main Office"`
	PickupDateTime             time.Time `gorm:"column:pickup_datetime" json:"pickup_datetime" example:"2025-02-16T08:00:00Z"`
	ReferenceNumber            string    `gorm:"column:reference_number" json:"reference_number" example:"REF123456"`
	AttachedDocument           string    `gorm:"column:attached_document" json:"attached_document" example:"document.pdf"`
	IsPEAEmployeeDriver        string    `gorm:"column:is_pea_employee_driver" json:"is_pea_employee_driver" example:"1"`
	IsAdminChooseDriver        string    `gorm:"column:is_admin_choose_driver" json:"is_admin_choose_driver" example:"1"`
	RefCostTypeCode            string    `gorm:"column:ref_cost_type_code" json:"ref_cost_type_code" example:"COST123"`
	CostNo                     string    `gorm:"column:cost_no" json:"cost_no" example:"COSTNO123"`

	MasCarpoolDriverUID string       `gorm:"column:mas_carpool_driver_uid;type:uuid" json:"mas_carpool_driver_uid"`
	VMSMasDriver        VmsMasDriver `gorm:"foreignKey:MasCarpoolDriverUID;references:MasDriverUID" json:"driver"`

	MasVehicleUID string        `gorm:"column:mas_vehicle_uid;type:uuid" json:"mas_vehicle_uid"`
	VmsMasVehicle VmsMasVehicle `gorm:"foreignKey:MasVehicleUID;references:MasVehicleUID" json:"vehicle"`

	ReceivedKeyPlace         string    `gorm:"column:received_key_place" json:"received_key_place"`
	ReceivedKeyStartDatetime time.Time `gorm:"column:received_key_start_datetime" json:"received_key_start_datetime"`
	ReceivedKeyEndDatetime   time.Time `gorm:"column:received_key_end_datetime" json:"received_key_end_datetime"`
}

func (VmsTrnRequest_Response) TableName() string {
	return "public.vms_trn_request"
}

type VmsTrnRequest_Update_VehicleUser struct {
	TrnRequestUID          string `gorm:"column:trn_request_uid;type:uuid;" json:"trn_request_uid" example:"a7de5318-1e05-4511-abe7-8c1c6374ab29"`
	VehicleUserEmpID       string `gorm:"column:vehicle_user_emp_id" json:"vehicle_user_emp_id" example:"700001"`
	VehicleUserEmpName     string `gorm:"column:vehicle_user_emp_name" json:"vehicle_user_emp_name" example:"John Doe"`
	VehicleUserDeptSAP     string `gorm:"column:vehicle_user_dept_sap" json:"vehicle_user_dept_sap" example:"DPT001"`
	CarUserInternalContact string `gorm:"column:car_user_internal_contact_number" json:"car_user_internal_contact_number" example:"1234567890"`
	CarUserMobileContact   string `gorm:"column:car_user_mobile_contact_number" json:"car_user_mobile_contact_number" example:"0987654321"`
}

func (VmsTrnRequest_Update_VehicleUser) TableName() string {
	return "public.vms_trn_request"
}

type VmsTrnRequest_Update_Trip struct {
	TrnRequestUID      string    `gorm:"column:trn_request_uid;type:uuid;" json:"trn_request_uid" example:"a7de5318-1e05-4511-abe7-8c1c6374ab29"`
	StartDatetime      time.Time `gorm:"column:start_datetime" json:"start_datetime" example:"2025-01-01T08:00:00Z"`
	EndDatetime        time.Time `gorm:"column:end_datetime" json:"end_datetime" example:"2025-01-01T10:00:00Z"`
	TripType           int       `gorm:"column:trip_type" json:"trip_type" example:"1"`
	ReservedTimeType   string    `gorm:"column:reserved_time_type" json:"reserved_time_type" example:"1"`
	WorkPlace          string    `gorm:"column:work_place" json:"work_place" example:"Head Office"`
	Objective          string    `gorm:"column:objective" json:"objective" example:"Business Meeting"`
	NumberOfPassengers int       `gorm:"column:number_of_passengers" json:"number_of_passengers" example:"3"`
}

func (VmsTrnRequest_Update_Trip) TableName() string {
	return "public.vms_trn_request"
}

type VmsTrnRequest_Update_Pickup struct {
	TrnRequestUID  string    `gorm:"column:trn_request_uid;type:uuid;" json:"trn_request_uid" example:"a7de5318-1e05-4511-abe7-8c1c6374ab29"`
	PickupPlace    string    `gorm:"column:pickup_place" json:"pickup_place" example:"Main Office"`
	PickupDateTime time.Time `gorm:"column:pickup_datetime" json:"pickup_datetime" example:"2025-02-16T08:00:00Z"`
}

func (VmsTrnRequest_Update_Pickup) TableName() string {
	return "public.vms_trn_request"
}

type VmsTrnRequest_Update_Document struct {
	TrnRequestUID    string `gorm:"column:trn_request_uid;type:uuid;" json:"trn_request_uid" example:"a7de5318-1e05-4511-abe7-8c1c6374ab29"`
	ReferenceNumber  string `gorm:"column:reference_number" json:"reference_number" example:"REF123456"`
	AttachedDocument string `gorm:"column:attached_document" json:"attached_document" example:"document.pdf"`
}

func (VmsTrnRequest_Update_Document) TableName() string {
	return "public.vms_trn_request"
}

type VmsTrnRequest_Update_Cost struct {
	TrnRequestUID   string `gorm:"column:trn_request_uid;type:uuid;" json:"trn_request_uid" example:"a7de5318-1e05-4511-abe7-8c1c6374ab29"`
	RefCostTypeCode string `gorm:"column:ref_cost_type_code" json:"ref_cost_type_code" example:"COST123"`
	CostNo          string `gorm:"column:cost_no" json:"cost_no" example:"COSTNO123"`
}

func (VmsTrnRequest_Update_Cost) TableName() string {
	return "public.vms_trn_request"
}

type VmsTrnRequest_Update_VehicleType struct {
	TrnRequestUID        string `gorm:"column:trn_request_uid;type:uuid;" json:"trn_request_uid" example:"a7de5318-1e05-4511-abe7-8c1c6374ab29"`
	RequestVehicleTypeId int    `gorm:"column:requested_vehicle_type_id" json:"requested_vehicle_type_id" example:"1"`
}

func (VmsTrnRequest_Update_VehicleType) TableName() string {
	return "public.vms_trn_request"
}

type VmsTrnRequest_Update_Approver struct {
	TrnRequestUID        string `gorm:"column:trn_request_uid;type:uuid;" json:"trn_request_uid" example:"a7de5318-1e05-4511-abe7-8c1c6374ab29"`
	ApprovedRequestEmpId string `gorm:"column:approved_request_emp_id" json:"approved_request_emp_id" example:"700002"`
}

func (VmsTrnRequest_Update_Approver) TableName() string {
	return "public.vms_trn_request"
}

// Approver
type VmsTrnRequest_Approved struct {
	TrnRequestUID string `gorm:"column:trn_request_uid;type:uuid;" json:"trn_request_uid" example:"a7de5318-1e05-4511-abe7-8c1c6374ab29"`
}

func (VmsTrnRequest_Approved) TableName() string {
	return "public.vms_trn_request"
}

type VmsTrnRequest_Approved_Update struct {
	VmsTrnRequest_Approved
	RefRequestStatusCode        string `gorm:"column:ref_request_status_code" json:"ref_request_status_code"`
	ApprovedRequestEmpID        string `gorm:"column:approved_request_emp_id" json:"approved_request_emp_id"`
	ApprovedRequestEmpName      string `gorm:"column:approved_request_emp_name" json:"approved_request_emp_name"`
	ApprovedRequestDeptSAP      string `gorm:"column:approved_request_dept_sap" json:"approved_request_dept_sap"`
	ApprovedRequestDeptSAPShort string `gorm:"column:approved_request_dept_sap_short" json:"approved_request_dept_sap_short"`
	ApprovedRequestDeptSAPFull  string `gorm:"column:approved_request_dept_sap_full" json:"approved_request_dept_sap_full"`
	LogUpdate
}

type VmsTrnRequest_Canceled struct {
	TrnRequestUID         string `gorm:"column:trn_request_uid;type:uuid;" json:"trn_request_uid" example:"a7de5318-1e05-4511-abe7-8c1c6374ab29"`
	CanceledRequestReason string `gorm:"column:canceled_request_reason;" json:"canceled_request_reason" example:"Test Cancel"`
}

func (VmsTrnRequest_Canceled) TableName() string {
	return "public.vms_trn_request"
}

type VmsTrnRequest_Canceled_Update struct {
	VmsTrnRequest_Canceled
	RefRequestStatusCode        string `gorm:"column:ref_request_status_code" json:"ref_request_status_code"`
	CanceledRequestEmpID        string `gorm:"column:canceled_request_emp_id" json:"canceled_request_emp_id"`
	CanceledRequestEmpName      string `gorm:"column:canceled_request_emp_name" json:"canceled_request_emp_name"`
	CanceledRequestDeptSAP      string `gorm:"column:canceled_request_dept_sap" json:"canceled_request_dept_sap"`
	CanceledRequestDeptSAPShort string `gorm:"column:canceled_request_dept_sap_short" json:"canceled_request_dept_sap_short"`
	CanceledRequestDeptSAPFull  string `gorm:"column:canceled_request_dept_sap_full" json:"canceled_request_dept_sap_full"`
	LogUpdate
}

type VmsTrnRequest_SendedBack struct {
	TrnRequestUID           string `gorm:"column:trn_request_uid;type:uuid;" json:"trn_request_uid" example:"a7de5318-1e05-4511-abe7-8c1c6374ab29"`
	SendedBackRequestReason string `gorm:"column:sended_back_request_reason;" json:"sended_back_request_reason" example:"Test Send Back"`
}

func (VmsTrnRequest_SendedBack) TableName() string {
	return "public.vms_trn_request"
}

type VmsTrnRequest_SendedBack_Update struct {
	VmsTrnRequest_SendedBack
	RefRequestStatusCode          string `gorm:"column:ref_request_status_code" json:"ref_request_status_code"`
	SendedBackRequestEmpID        string `gorm:"column:sended_back_request_emp_id" json:"sended_back_request_emp_id"`
	SendedBackRequestEmpName      string `gorm:"column:sended_back_request_emp_name" json:"sended_back_request_emp_name"`
	SendedBackRequestDeptSAP      string `gorm:"column:sended_back_request_dept_sap" json:"sended_back_request_dept_sap"`
	SendedBackRequestDeptSAPShort string `gorm:"column:sended_back_request_dept_sap_short" json:"sended_back_request_dept_sap_short"`
	SendedBackRequestDeptSAPFull  string `gorm:"column:sended_back_request_dept_sap_full" json:"sended_back_request_dept_sap_full"`
	LogUpdate
}

type VmsTrnRequest_Update_Driver struct {
	TrnRequestUID       string `gorm:"column:trn_request_uid;type:uuid;" json:"trn_request_uid" example:"a7de5318-1e05-4511-abe7-8c1c6374ab29"`
	MasCarPoolDriverUID string `gorm:"column:mas_carpool_driver_uid" json:"mas_carpool_driver_uid" example:"a6c8a34b-9245-49c8-a12b-45fae77a4e7d"`
}

func (VmsTrnRequest_Update_Driver) TableName() string {
	return "public.vms_trn_request"
}

type VmsTrnRequest_Update_Vehicle struct {
	TrnRequestUID string `gorm:"column:trn_request_uid;type:uuid;" json:"trn_request_uid" example:"a7de5318-1e05-4511-abe7-8c1c6374ab29"`
	MasVehicleUID string `gorm:"column:mas_vehicle_uid;type:uuid" json:"mas_vehicle_uid"  example:"a6c8a34b-9245-49c8-a12b-45fae77a4e7d"`
}

func (VmsTrnRequest_Update_Vehicle) TableName() string {
	return "public.vms_trn_request"
}
