package models

import (
	"time"
)

type VmsTrnRequestList struct {
	TrnRequestUid                    string    `gorm:"column:trn_request_uid;type:uuid;" json:"trn_request_uid"`
	RequestNo                        string    `gorm:"column:request_no" json:"request_no"`
	VehicleUserEmpID                 string    `gorm:"column:vehicle_user_emp_id" json:"vehicle_user_emp_id"`
	VehicleUserEmpName               string    `gorm:"column:vehicle_user_emp_name" json:"vehicle_user_emp_name"`
	VehicleUserDeptSAPShort          string    `gorm:"column:vehicle_user_dept_sap_name_short" json:"vehicle_user_dept_sap_short" example:"Finance"`
	VehicleLicensePlate              string    `gorm:"column:vehicle_license_plate" json:"vehicle_license_plate"`
	VehicleLicensePlateProvinceShort string    `gorm:"column:vehicle_license_plate_province_short" json:"vehicle_license_plate_province_short"`
	VehicleLicensePlateProvinceFull  string    `gorm:"column:vehicle_license_plate_province_full" json:"vehicle_license_plate_province_full"`
	VehicleDepartmentDeptSapShort    string    `gorm:"column:vehicle_department_dept_sap_short" json:"vehicle_department_dept_sap_short"`
	WorkPlace                        string    `gorm:"column:work_place" json:"work_place"`
	StartDatetime                    string    `gorm:"column:start_datetime" json:"start_datetime"`
	EndDatetime                      string    `gorm:"column:end_datetime" json:"end_datetime"`
	RefRequestStatusCode             string    `gorm:"column:ref_request_status_code" json:"ref_request_status_code"`
	RefRequestStatusName             string    `json:"ref_request_status_name"`
	IsHaveSubRequest                 string    `gorm:"column:is_have_sub_request" json:"is_have_sub_request" example:"0"`
	ReceivedKeyPlace                 string    `gorm:"column:received_key_place" json:"received_key_place"`
	ReceivedKeyStartDatetime         time.Time `gorm:"column:received_key_start_datetime" json:"received_key_start_datetime"`
	ReceivedKeyEndDatetime           time.Time `gorm:"column:received_key_end_datetime" json:"received_key_end_datetime"`
}
type VmsTrnRequestSummary struct {
	RefRequestStatusCode string `gorm:"column:ref_request_status_code" json:"ref_request_status_code"`
	RefRequestStatusName string `json:"ref_request_status_name"`
	Count                int    `gorm:"column:count" json:"count"`
}

// VmsTrnRequestAdminList
type VmsTrnRequestAdminList struct {
	VmsTrnRequestList
	RefVehicleTypeName   string `gorm:"column:ref_vehicle_type_name" json:"ref_vehicle_type_name"`
	DriverEmpId          string `gorm:"column:driver_emp_id" json:"driver_emp_id"`
	MasVehicleUID        string `gorm:"column:mas_vehicle_uid" json:"mas_vehicle_uid"`
	MasCarpoolDriverUID  string `gorm:"column:mas_carpool_driver_uid" json:"mas_carpool_driver_uid"`
	DriverName           string `gorm:"column:driver_name" json:"driver_name"`
	DriverDeptName       string `gorm:"column:driver_dept_name" json:"driver_dept_name"`
	VehicleDeptName      string `gorm:"column:vehicle_dept_name" json:"vehicle_dept_name"`
	VehicleCarpoolName   string `gorm:"column:vehicle_carpool_name" json:"vehicle_carpool_name"`
	IsAdminChooseDriver  int    `gorm:"column:is_admin_choose_driver" json:"is_admin_choose_driver"`
	IsAdminChooseVehicle int    `gorm:"column:is_admin_choose_vehicle" json:"is_admin_choose_vehicle"`
	IsPEAEmployeeDriver  int    `gorm:"column:is_pea_employee_driver" json:"is_pea_employee_driver"`
	TripType             int    `gorm:"column:trip_type" json:"trip_type" example:"1"`
	TripTypeName         string `gorm:"-" json:"trip_type_name" example:"1"`
	Can_Choose_Vehicle   bool   `gorm:"-" json:"can_choose_vehicle"`
	Can_Choose_Driver    bool   `gorm:"-" json:"can_choose_driver"`
}

func (VmsTrnRequestAdminList) TableName() string {
	return "public.vms_trn_request"
}

type VmsTrnRequestRequest struct {
	//Step1
	VehicleUserEmpID             string    `gorm:"column:vehicle_user_emp_id" json:"vehicle_user_emp_id" example:"990001"`
	VehicleUserEmpName           string    `gorm:"column:vehicle_user_emp_name" json:"vehicle_user_emp_name" example:"John Doe"`
	VehicleUserDeptSAP           string    `gorm:"column:vehicle_user_dept_sap" json:"vehicle_user_dept_sap" example:"DPT001"`
	CarUserInternalContactNumber string    `gorm:"column:car_user_internal_contact_number" json:"car_user_internal_contact_number" example:"1234567890"`
	CarUserMobileContactNumber   string    `gorm:"column:car_user_mobile_contact_number" json:"car_user_mobile_contact_number" example:"0987654321"`
	StartDatetime                time.Time `gorm:"column:start_datetime" json:"start_datetime" example:"2025-01-01T08:00:00Z"`
	EndDatetime                  time.Time `gorm:"column:end_datetime" json:"end_datetime" example:"2025-01-01T10:00:00Z"`
	TripType                     int       `gorm:"column:trip_type" json:"trip_type" example:"1"`
	ReservedTimeType             string    `gorm:"column:reserved_time_type" json:"reserved_time_type" example:"1"`
	WorkPlace                    string    `gorm:"column:work_place" json:"work_place" example:"Head Office"`
	Objective                    string    `gorm:"column:objective" json:"objective" example:"Business Meeting"`
	NumberOfPassengers           int       `gorm:"column:number_of_passengers" json:"number_of_passengers" example:"3"`
	Remark                       string    `gorm:"column:remark" json:"remark" example:"Urgent request"`
	ReferenceNumber              string    `gorm:"column:reference_number" json:"reference_number" example:"REF123456"`
	AttachedDocument             string    `gorm:"column:attached_document" json:"attached_document" example:"document.pdf"`
	RefCostTypeCode              int       `gorm:"column:ref_cost_type_code" json:"ref_cost_type_code" example:"101"`
	CostNo                       string    `gorm:"column:cost_no" json:"cost_no" example:"COST2024001"`

	//Step 2
	MasVehicleUID         string `gorm:"column:mas_vehicle_uid" json:"mas_vehicle_uid" example:"389b0f63-4195-4ece-bf35-0011c2f5f28c"`
	IsAdminChooseVehicle  string `gorm:"column:is_admin_choose_vehicle;type:bit(1)" json:"is_admin_choose_vehicle" example:"0"`
	IsSystemChooseVehicle string `gorm:"-" json:"is_system_choose_vehicle" example:"0"`
	RequestVehicleTypeID  int    `gorm:"column:requested_vehicle_type_id" json:"requested_vehicle_type_id" example:"1"`

	//Step 3
	IsDriverNeed        string `gorm:"column:is_driver_need" json:"-" example:"1"`
	MasCarPoolDriverUID string `gorm:"column:mas_carpool_driver_uid" json:"mas_carpool_driver_uid" example:"a6c8a34b-9245-49c8-a12b-45fae77a4e7d"`
	IsPEAEmployeeDriver string `gorm:"column:is_pea_employee_driver" json:"is_pea_employee_driver" example:"1"`
	IsAdminChooseDriver string `gorm:"column:is_admin_choose_driver" json:"-" example:"0"`

	DriverEmpID           string `gorm:"column:driver_emp_id" json:"driver_emp_id" example:"700001"`
	DriverEmpName         string `gorm:"column:driver_emp_name" json:"driver_emp_name" example:"John Doe"`
	DriverDeptSAP         string `gorm:"column:driver_emp_dept_sap" json:"driver_emp_dept_sap" example:"DPT001"`
	DriverInternalContact string `gorm:"column:driver_internal_contact_number" json:"driver_internal_contact_number" example:"1234567890"`
	DriverMobileContact   string `gorm:"column:driver_mobile_contact_number" json:"driver_mobile_contact_number" example:"0987654321"`

	PickupPlace    string    `gorm:"column:pickup_place" json:"pickup_place" example:"Main Office"`
	PickupDateTime time.Time `gorm:"column:pickup_datetime" json:"pickup_datetime" example:"2025-02-16T08:30:00Z"`

	//Step 4
	ApprovedRequestEmpID        string `gorm:"column:approved_request_emp_id" json:"approved_request_emp_id" example:"EMP67890"`
	ApprovedRequestEmpName      string `gorm:"column:approved_request_emp_name" json:"approved_request_emp_name" example:"Jane Doe"`
	ApprovedRequestDeptSAP      string `gorm:"column:approved_request_dept_sap" json:"approved_request_dept_sap" example:"Finance"`
	ApprovedRequestDeptSAPShort string `gorm:"column:approved_request_dept_sap_short" json:"approved_request_dept_sap_short" example:"Finance"`
	ApprovedRequestDeptSAPFull  string `gorm:"column:approved_request_dept_sap_full" json:"approved_request_dept_sap_full" example:"Finance"`

	//
	RefRequestTypeCode int    `gorm:"column:ref_request_type_code" json:"-" example:"1"`
	IsHaveSubRequest   string `gorm:"column:is_have_sub_request" json:"-" example:"0"`
}

type VmsTrnRequestCreate struct {
	VmsTrnRequestRequest
	TrnRequestUID              string    `gorm:"column:trn_request_uid;type:uuid;" json:"trn_request_uid"`
	RequestNo                  string    `gorm:"column:request_no" json:"request_no"`
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

func (VmsTrnRequestCreate) TableName() string {
	return "public.vms_trn_request"
}

type VmsTrnRequestResponse struct {
	TrnRequestUID                    string    `gorm:"column:trn_request_uid;type:uuid;" json:"trn_request_uid"`
	RequestNo                        string    `gorm:"column:request_no" json:"request_no"`
	VehicleUserEmpName               string    `gorm:"column:vehicle_user_emp_name" json:"vehicle_user_emp_name" example:"John Smith"`
	VehicleUserDeptSAP               string    `gorm:"column:vehicle_user_dept_sap" json:"vehicle_user_dept_sap" example:"HR"`
	VehicleUserEmpID                 string    `gorm:"column:vehicle_user_emp_id" json:"vehicle_user_emp_id" example:"700001"`
	VehicleUserDeptSAPShort          string    `gorm:"column:vehicle_user_dept_sap_name_short" json:"vehicle_user_dept_sap_short" example:"Finance"`
	VehicleUserDeptSAPFull           string    `gorm:"column:vehicle_user_dept_sap_name_full" json:"vehicle_user_dept_sap_full" example:"Finance"`
	CarUserMobileContactNumber       string    `gorm:"column:car_user_mobile_contact_number" json:"car_user_mobile_contact_number" example:"9876543210"`
	CarUserInternalContactNumber     string    `gorm:"column:car_user_internal_contact_number" json:"car_user_internal_contact_number" example:"9876543210"`
	VehicleLicensePlate              string    `gorm:"column:vehicle_license_plate" json:"vehicle_license_plate" example:"ABC1234"`
	VehicleLicensePlateProvinceShort string    `gorm:"column:vehicle_license_plate_province_short" json:"vehicle_license_plate_province_short"`
	VehicleLicensePlateProvinceFull  string    `gorm:"column:vehicle_license_plate_province_full" json:"vehicle_license_plate_province_full"`
	ApprovedRequestEmpID             string    `gorm:"column:approved_request_emp_id" json:"approved_request_emp_id" example:"EMP67890"`
	ApprovedRequestEmpName           string    `gorm:"column:approved_request_emp_name" json:"approved_request_emp_name" example:"Jane Doe"`
	ApprovedRequestDeptSAP           string    `gorm:"column:approved_request_dept_sap" json:"approved_request_dept_sap" example:"Finance"`
	ApprovedRequestDeptSAPShort      string    `gorm:"column:approved_request_dept_sap_short" json:"approved_request_dept_sap_short" example:"Finance"`
	ApprovedRequestDeptSAPFull       string    `gorm:"column:approved_request_dept_sap_full" json:"approved_request_dept_sap_full" example:"Finance"`
	StartDateTime                    time.Time `gorm:"column:start_datetime" json:"start_datetime" example:"2025-02-16T08:30:00Z"`
	EndDateTime                      time.Time `gorm:"column:end_datetime" json:"end_datetime" example:"2025-02-16T09:30:00Z"`
	DateRange                        string    `gorm:"column:date_range" json:"date_range" example:"2025-02-16 to 2025-02-17"`
	TripType                         int       `gorm:"column:trip_type" json:"trip_type" example:"1"`
	WorkPlace                        string    `gorm:"column:work_place" json:"work_place" example:"Office"`
	Objective                        string    `gorm:"column:objective" json:"objective" example:"Project meeting"`
	Remark                           string    `gorm:"column:remark" json:"remark" example:"Special request for parking spot"`
	NumberOfPassengers               int       `gorm:"column:number_of_passengers" json:"number_of_passengers" example:"4"`
	PickupPlace                      string    `gorm:"column:pickup_place" json:"pickup_place" example:"Main Office"`
	PickupDateTime                   time.Time `gorm:"column:pickup_datetime" json:"pickup_datetime" example:"2025-02-16T08:00:00Z"`
	ReferenceNumber                  string    `gorm:"column:reference_number" json:"reference_number" example:"REF123456"`
	AttachedDocument                 string    `gorm:"column:attached_document" json:"attached_document" example:"document.pdf"`
	IsPEAEmployeeDriver              string    `gorm:"column:is_pea_employee_driver" json:"is_pea_employee_driver" example:"1"`
	IsAdminChooseDriver              string    `gorm:"column:is_admin_choose_driver" json:"is_admin_choose_driver" example:"1"`
	NumberOfAvailableDrivers         int       `gorm:"-" json:"number_of_available_drivers" example:"2"`
	RefCostTypeCode                  string    `gorm:"column:ref_cost_type_code" json:"ref_cost_type_code" example:"COST123"`
	CostNo                           string    `gorm:"column:cost_no" json:"cost_no" example:"COSTNO123"`

	MasCarpoolDriverUID  string            `gorm:"column:mas_carpool_driver_uid;type:uuid" json:"mas_carpool_driver_uid"`
	VMSMasDriver         VmsMasDriver      `gorm:"foreignKey:MasCarpoolDriverUID;references:MasDriverUID" json:"driver"`
	IsAdminChooseVehicle string            `gorm:"column:is_admin_choose_vehicle" json:"is_admin_choose_vehicle" example:"0"`
	RequestVehicleTypeID int               `gorm:"column:requested_vehicle_type_id" json:"requested_vehicle_type_id" example:"1"`
	RequestVehicleType   VmsRefVehicleType `gorm:"foreignKey:RequestVehicleTypeID;references:RefVehicleTypeCode" json:"request_vehicle_type"`

	DriverEmpID           string `gorm:"column:driver_emp_id" json:"driver_emp_id" example:"700001"`
	DriverEmpName         string `gorm:"column:driver_emp_name" json:"driver_emp_name" example:"John Doe"`
	DriverDeptSAP         string `gorm:"column:driver_emp_dept_sap" json:"driver_emp_dept_sap" example:"DPT001"`
	DriverInternalContact string `gorm:"column:driver_internal_contact_number" json:"driver_internal_contact_number" example:"1234567890"`
	DriverMobileContact   string `gorm:"column:driver_mobile_contact_number" json:"driver_mobile_contact_number" example:"0987654321"`
	DriverImageURL        string `gorm:"-" json:"driver_image_url"`

	MasVehicleUID                 string        `gorm:"column:mas_vehicle_uid;type:uuid" json:"mas_vehicle_uid"`
	VehicleDepartmentDeptSap      string        `gorm:"column:vehicle_department_dept_sap" json:"vehicle_department_dept_sap"`
	VehicleDepartmentDeptSapShort string        `gorm:"column:vehicle_department_dept_sap_short" json:"mas_vehicle_department_dept_sap_short"`
	VehicleDepartmentDeptSapFull  string        `gorm:"column:vehicle_department_dept_sap_full" json:"mas_vehicle_department_dept_sap_full"`
	VmsMasVehicle                 VmsMasVehicle `gorm:"foreignKey:MasVehicleUID;references:MasVehicleUID" json:"vehicle"`

	ReceivedKeyPlace         string    `gorm:"column:received_key_place" json:"received_key_place"`
	ReceivedKeyStartDatetime time.Time `gorm:"column:received_key_start_datetime" json:"received_key_start_datetime"`
	ReceivedKeyEndDatetime   time.Time `gorm:"column:received_key_end_datetime" json:"received_key_end_datetime"`

	CanCancelRequest        bool                    `gorm:"-" json:"can_cancel_request"`
	RefRequestStatusCode    string                  `gorm:"column:ref_request_status_code" json:"ref_request_status_code"`
	RefRequestStatus        VmsRefRequestStatus     `gorm:"foreignKey:RefRequestStatusCode;references:RefRequestStatusCode" json:"ref_request_status"`
	RefRequestStatusName    string                  `json:"ref_request_status_name"`
	SendedBackRequestReason string                  `gorm:"column:sended_back_request_reason;" json:"sended_back_request_reason" example:"Test Send Back"`
	CanceledRequestReason   string                  `gorm:"column:canceled_request_reason;" json:"canceled_request_reason" example:"Test Cancel"`
	ProgressRequestStatus   []ProgressRequestStatus `gorm:"-" json:"progress_request_status"`
}

func (VmsTrnRequestResponse) TableName() string {
	return "public.vms_trn_request"
}

type VmsTrnRequestRequestNo struct {
	RequestNo string `gorm:"column:request_no" json:"request_no"`
}

type ProgressRequestStatus struct {
	ProgressIcon string `gorm:"column:progress_icon" json:"progress_icon"`
	ProgressName string `gorm:"column:progress_name" json:"progress_name"`
}

// VmsTrnRequestVehicleUser
type VmsTrnRequestVehicleUser struct {
	TrnRequestUID               string    `gorm:"column:trn_request_uid;primarykey" json:"trn_request_uid" example:"8bd09808-61fa-42fd-8a03-bf961b5678cd"`
	VehicleUserEmpID            string    `gorm:"column:vehicle_user_emp_id" json:"vehicle_user_emp_id" example:"990001"`
	VehicleUserEmpName          string    `gorm:"column:vehicle_user_emp_name" json:"-"`
	VehicleUserDeptSAP          string    `gorm:"column:vehicle_user_dept_sap" json:"-"`
	VehicleUserDeptSAPNameShort string    `gorm:"column:vehicle_user_dept_sap_name_short" json:"-"`
	VehicleUserDeptSAPNameFull  string    `gorm:"column:vehicle_user_dept_sap_name_full" json:"-"`
	CarUserInternalContact      string    `gorm:"column:car_user_internal_contact_number" json:"car_user_internal_contact_number" example:"1234567890"`
	CarUserMobileContact        string    `gorm:"column:car_user_mobile_contact_number" json:"car_user_mobile_contact_number" example:"0987654321"`
	UpdatedAt                   time.Time `gorm:"column:updated_at" json:"-"`
	UpdatedBy                   string    `gorm:"column:updated_by" json:"-"`
}

func (VmsTrnRequestVehicleUser) TableName() string {
	return "public.vms_trn_request"
}

// VmsTrnRequestTrip
type VmsTrnRequestTrip struct {
	TrnRequestUID      string    `gorm:"column:trn_request_uid;primarykey" json:"trn_request_uid" example:"8bd09808-61fa-42fd-8a03-bf961b5678cd"`
	StartDatetime      time.Time `gorm:"column:start_datetime" json:"start_datetime" example:"2025-01-01T08:00:00Z"`
	EndDatetime        time.Time `gorm:"column:end_datetime" json:"end_datetime" example:"2025-01-01T10:00:00Z"`
	TripType           int       `gorm:"column:trip_type" json:"trip_type" example:"1"`
	ReservedTimeType   string    `gorm:"column:reserved_time_type" json:"reserved_time_type" example:"1"`
	WorkPlace          string    `gorm:"column:work_place" json:"work_place" example:"Head Office"`
	Objective          string    `gorm:"column:objective" json:"objective" example:"Business Meeting"`
	NumberOfPassengers int       `gorm:"column:number_of_passengers" json:"number_of_passengers" example:"3"`
	Remark             string    `gorm:"column:remark" json:"remark" example:"Special request for parking spot"`
	UpdatedAt          time.Time `gorm:"column:updated_at" json:"-"`
	UpdatedBy          string    `gorm:"column:updated_by" json:"-"`
}

func (VmsTrnRequestTrip) TableName() string {
	return "public.vms_trn_request"
}

// VmsTrnRequestPickup
type VmsTrnRequestPickup struct {
	TrnRequestUID  string    `gorm:"column:trn_request_uid;primarykey" json:"trn_request_uid" example:"8bd09808-61fa-42fd-8a03-bf961b5678cd"`
	PickupPlace    string    `gorm:"column:pickup_place" json:"pickup_place" example:"Main Office"`
	PickupDateTime time.Time `gorm:"column:pickup_datetime" json:"pickup_datetime" example:"2025-02-16T08:00:00Z"`
	UpdatedAt      time.Time `gorm:"column:updated_at" json:"-"`
	UpdatedBy      string    `gorm:"column:updated_by" json:"-"`
}

func (VmsTrnRequestPickup) TableName() string {
	return "public.vms_trn_request"
}

// VmsTrnRequestDocument
type VmsTrnRequestDocument struct {
	TrnRequestUID    string    `gorm:"column:trn_request_uid;primarykey" json:"trn_request_uid" example:"8bd09808-61fa-42fd-8a03-bf961b5678cd"`
	ReferenceNumber  string    `gorm:"column:reference_number" json:"reference_number" example:"REF123456"`
	AttachedDocument string    `gorm:"column:attached_document" json:"attached_document" example:"document.pdf"`
	UpdatedAt        time.Time `gorm:"column:updated_at" json:"-"`
	UpdatedBy        string    `gorm:"column:updated_by" json:"-"`
}

func (VmsTrnRequestDocument) TableName() string {
	return "public.vms_trn_request"
}

// VmsTrnRequestCost
type VmsTrnRequestCost struct {
	TrnRequestUID   string    `gorm:"column:trn_request_uid;primarykey" json:"trn_request_uid" example:"8bd09808-61fa-42fd-8a03-bf961b5678cd"`
	RefCostTypeCode string    `gorm:"column:ref_cost_type_code" json:"ref_cost_type_code" example:"2"`
	CostNo          string    `gorm:"column:cost_no" json:"cost_no" example:"COSTNO123"`
	UpdatedAt       time.Time `gorm:"column:updated_at" json:"-"`
	UpdatedBy       string    `gorm:"column:updated_by" json:"-"`
}

func (VmsTrnRequestCost) TableName() string {
	return "public.vms_trn_request"
}

// VmsTrnRequestVehicleType
type VmsTrnRequestVehicleType struct {
	TrnRequestUID        string    `gorm:"column:trn_request_uid;primarykey" json:"trn_request_uid" example:"8bd09808-61fa-42fd-8a03-bf961b5678cd"`
	RequestVehicleTypeId int       `gorm:"column:requested_vehicle_type_id" json:"requested_vehicle_type_id" example:"1"`
	UpdatedAt            time.Time `gorm:"column:updated_at" json:"-"`
	UpdatedBy            string    `gorm:"column:updated_by" json:"-"`
}

func (VmsTrnRequestVehicleType) TableName() string {
	return "public.vms_trn_request"
}

// VmsTrnRequestApprover
type VmsTrnRequestApprover struct {
	TrnRequestUID               string `gorm:"column:trn_request_uid;primarykey" json:"trn_request_uid" example:"8bd09808-61fa-42fd-8a03-bf961b5678cd"`
	ApprovedRequestEmpID        string `gorm:"column:approved_request_emp_id" json:"approved_request_emp_id"`
	ApprovedRequestEmpName      string `gorm:"column:approved_request_emp_name" json:"-"`
	ApprovedRequestDeptSAP      string `gorm:"column:approved_request_dept_sap" json:"-"`
	ApprovedRequestDeptSAPShort string `gorm:"column:approved_request_dept_sap_short" json:"-"`
	ApprovedRequestDeptSAPFull  string `gorm:"column:approved_request_dept_sap_full" json:"-"`

	UpdatedAt time.Time `gorm:"column:updated_at" json:"-"`
	UpdatedBy string    `gorm:"column:updated_by" json:"-"`
}

func (VmsTrnRequestApprover) TableName() string {
	return "public.vms_trn_request"
}

// VmsTrnRequestApproved
type VmsTrnRequestApproved struct {
	TrnRequestUID               string    `gorm:"column:trn_request_uid;primaryKey" json:"trn_request_uid" example:"8bd09808-61fa-42fd-8a03-bf961b5678cd"`
	ApprovedRequestEmpID        string    `gorm:"column:approved_request_emp_id" json:"-"`
	ApprovedRequestEmpName      string    `gorm:"column:approved_request_emp_name" json:"-"`
	ApprovedRequestDeptSAP      string    `gorm:"column:approved_request_dept_sap" json:"-"`
	ApprovedRequestDeptSAPShort string    `gorm:"column:approved_request_dept_sap_short" json:"-"`
	ApprovedRequestDeptSAPFull  string    `gorm:"column:approved_request_dept_sap_full" json:"-"`
	RefRequestStatusCode        string    `gorm:"column:ref_request_status_code" json:"-"`
	UpdatedAt                   time.Time `gorm:"column:updated_at" json:"-"`
	UpdatedBy                   string    `gorm:"column:updated_by" json:"-"`
}

func (VmsTrnRequestApproved) TableName() string {
	return "public.vms_trn_request"
}

// VmsTrnRequestApprovedWithRecieiveKey
type VmsTrnRequestApprovedWithRecieiveKey struct {
	TrnRequestUID               string    `gorm:"column:trn_request_uid;primaryKey" json:"trn_request_uid" example:"8bd09808-61fa-42fd-8a03-bf961b5678cd"`
	ApprovedRequestEmpID        string    `gorm:"column:approved_request_emp_id" json:"approved_request_emp_id" example:"990001"`
	ApprovedRequestEmpName      string    `gorm:"column:approved_request_emp_name" json:"-"`
	ApprovedRequestDeptSAP      string    `gorm:"column:approved_request_dept_sap" json:"-"`
	ApprovedRequestDeptSAPShort string    `gorm:"column:approved_request_dept_sap_short" json:"-"`
	ApprovedRequestDeptSAPFull  string    `gorm:"column:approved_request_dept_sap_full" json:"-"`
	RefRequestStatusCode        string    `gorm:"column:ref_request_status_code" json:"-"`
	ReceivedKeyPlace            string    `gorm:"column:received_key_place" json:"received_key_place" example:"Main Office"`
	ReceivedKeyStartDatetime    time.Time `gorm:"column:received_key_start_datetime" json:"received_key_start_datetime" example:"2025-02-16T08:00:00Z"`
	ReceivedKeyEndDatetime      time.Time `gorm:"column:received_key_end_datetime" json:"-"`
	UpdatedBy                   string    `gorm:"column:updated_by" json:"-"`
	UpdatedAt                   time.Time `gorm:"column:updated_at" json:"-"`
}

func (VmsTrnRequestApprovedWithRecieiveKey) TableName() string {
	return "public.vms_trn_request"
}

// VmsTrnRequestCanceled
type VmsTrnRequestCanceled struct {
	TrnRequestUID               string    `gorm:"column:trn_request_uid;primarykey" json:"trn_request_uid" example:"8bd09808-61fa-42fd-8a03-bf961b5678cd"`
	CanceledRequestReason       string    `gorm:"column:canceled_request_reason;" json:"canceled_request_reason" example:"Test Cancel"`
	CanceledRequestEmpID        string    `gorm:"column:canceled_request_emp_id" json:"canceled_request_emp_id"`
	CanceledRequestEmpName      string    `gorm:"column:canceled_request_emp_name" json:"-"`
	CanceledRequestDeptSAP      string    `gorm:"column:canceled_request_dept_sap" json:"-"`
	CanceledRequestDeptSAPShort string    `gorm:"column:canceled_request_dept_sap_short" json:"-"`
	CanceledRequestDeptSAPFull  string    `gorm:"column:canceled_request_dept_sap_full" json:"-"`
	RefRequestStatusCode        string    `gorm:"column:ref_request_status_code" json:"-"`
	UpdatedAt                   time.Time `gorm:"column:updated_at" json:"-"`
	UpdatedBy                   string    `gorm:"column:updated_by" json:"-"`
}

func (VmsTrnRequestCanceled) TableName() string {
	return "public.vms_trn_request"
}

// VmsTrnRequestSendedBack
type VmsTrnRequestSendedBack struct {
	TrnRequestUID                 string    `gorm:"column:trn_request_uid;primarykey" json:"trn_request_uid" example:"8bd09808-61fa-42fd-8a03-bf961b5678cd"`
	SendedBackRequestReason       string    `gorm:"column:sended_back_request_reason;" json:"sended_back_request_reason" example:"Test Send Back"`
	RefRequestStatusCode          string    `gorm:"column:ref_request_status_code" json:"-"`
	SendedBackRequestEmpID        string    `gorm:"column:sended_back_request_emp_id" json:"-"`
	SendedBackRequestEmpName      string    `gorm:"column:sended_back_request_emp_name" json:"-"`
	SendedBackRequestDeptSAP      string    `gorm:"column:sended_back_request_dept_sap" json:"-"`
	SendedBackRequestDeptSAPShort string    `gorm:"column:sended_back_request_dept_sap_short" json:"-"`
	SendedBackRequestDeptSAPFull  string    `gorm:"column:sended_back_request_dept_sap_full" json:"-"`
	UpdatedAt                     time.Time `gorm:"column:updated_at" json:"-"`
	UpdatedBy                     string    `gorm:"column:updated_by" json:"-"`
}

func (VmsTrnRequestSendedBack) TableName() string {
	return "public.vms_trn_request"
}

// VmsTrnRequestDriver
type VmsTrnRequestDriver struct {
	TrnRequestUID       string    `gorm:"column:trn_request_uid;primarykey" json:"trn_request_uid" example:"8bd09808-61fa-42fd-8a03-bf961b5678cd"`
	MasCarPoolDriverUID string    `gorm:"column:mas_carpool_driver_uid" json:"mas_carpool_driver_uid" example:"a6c8a34b-9245-49c8-a12b-45fae77a4e7d"`
	UpdatedAt           time.Time `gorm:"column:updated_at" json:"-"`
	UpdatedBy           string    `gorm:"column:updated_by" json:"-"`
}

func (VmsTrnRequestDriver) TableName() string {
	return "public.vms_trn_request"
}

// VmsTrnRequestVehicle
type VmsTrnRequestVehicle struct {
	TrnRequestUID string    `gorm:"column:trn_request_uid;primarykey" json:"trn_request_uid" example:"8bd09808-61fa-42fd-8a03-bf961b5678cd"`
	MasVehicleUID string    `gorm:"column:mas_vehicle_uid" json:"mas_vehicle_uid"  example:"a6c8a34b-9245-49c8-a12b-45fae77a4e7d"`
	UpdatedAt     time.Time `gorm:"column:updated_at" json:"-"`
	UpdatedBy     string    `gorm:"column:updated_by" json:"-"`
}

func (VmsTrnRequestVehicle) TableName() string {
	return "public.vms_trn_request"
}

// VmsTrnRequestVehicleInfo
type VmsTrnRequestVehicleInfo struct {
	NumberOfAvailableDrivers int `gorm:"-" json:"number_of_available_drivers" example:"2"`
}
