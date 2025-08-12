package models

import (
	"time"
)

// VmsTrnRequestAdminList
type VmsTrnRequestAdminList struct {
	VmsTrnRequestList
	RefVehicleTypeName       string       `gorm:"column:ref_vehicle_type_name" json:"ref_vehicle_type_name"`
	DriverEmpId              string       `gorm:"column:driver_emp_id" json:"driver_emp_id"`
	MasVehicleUID            *string      `gorm:"column:mas_vehicle_uid" json:"mas_vehicle_uid"`
	MasCarpoolDriverUID      *string      `gorm:"column:mas_carpool_driver_uid" json:"mas_carpool_driver_uid"`
	DriverName               string       `gorm:"column:driver_name" json:"driver_name"`
	DriverDeptName           string       `gorm:"column:driver_dept_name" json:"driver_dept_name"`
	VehicleDeptName          string       `gorm:"column:vehicle_dept_name" json:"vehicle_dept_name"`
	VehicleCarpoolName       string       `gorm:"column:vehicle_carpool_name" json:"vehicle_carpool_name"`
	VehicleCarpoolText       string       `gorm:"column:vehicle_carpool_text" json:"vehicle_carpool_text"`
	IsAdminChooseDriver      int          `gorm:"column:is_admin_choose_driver" json:"-"`
	IsAdminChooseVehicle     int          `gorm:"column:is_admin_choose_vehicle" json:"-"`
	IsPEAEmployeeDriver      int          `gorm:"column:is_pea_employee_driver" json:"is_pea_employee_driver"`
	TripType                 int          `gorm:"column:ref_trip_type_code" json:"trip_type" example:"1"`
	TripTypeName             string       `gorm:"-" json:"trip_type_name" example:"1"`
	CanChooseVehicle         bool         `gorm:"-" json:"can_choose_vehicle"`
	CanChooseDriver          bool         `gorm:"-" json:"can_choose_driver"`
	CanceledRequestDatetime  TimeWithZone `gorm:"column:canceled_request_datetime" json:"canceled_request_datetime"`
	DriverCarpoolName        string       `gorm:"column:driver_carpool_name" json:"-"`
	RefCarpoolChooseCarID    int          `gorm:"column:ref_carpool_choose_car_id" json:"-"`
	RefCarpoolChooseDriverID int          `gorm:"column:ref_carpool_choose_driver_id" json:"-"`
	WorkDescription          string       `gorm:"column:work_description" json:"work_description"`
	KeyReceiverFullName      string       `gorm:"column:receiver_fullname" json:"key_receiver_fullname"`
	KeyReceiverDeptNameShort string       `gorm:"column:receiver_dept_name_short" json:"key_receiver_dept_name_short"`
}

func (VmsTrnRequestAdminList) TableName() string {
	return "public.vms_trn_request"
}

// VmsTrnRequestList
type VmsTrnRequestList struct {
	TrnRequestUid                    string                 `gorm:"column:trn_request_uid;primaryKey;" json:"trn_request_uid"`
	RequestNo                        string                 `gorm:"column:request_no" json:"request_no"`
	VehicleUserEmpID                 string                 `gorm:"column:vehicle_user_emp_id" json:"vehicle_user_emp_id"`
	VehicleUserPosition              string                 `gorm:"column:vehicle_user_position" json:"vehicle_user_position"`
	VehicleUserEmpName               string                 `gorm:"column:vehicle_user_emp_name" json:"vehicle_user_emp_name"`
	VehicleUserDeptNameShort         string                 `gorm:"column:vehicle_user_dept_name_short" json:"vehicle_user_dept_name_short" example:"Finance"`
	VehicleLicensePlate              string                 `gorm:"column:vehicle_license_plate" json:"vehicle_license_plate"`
	VehicleLicensePlateProvinceShort string                 `gorm:"column:vehicle_license_plate_province_short" json:"vehicle_license_plate_province_short"`
	VehicleLicensePlateProvinceFull  string                 `gorm:"column:vehicle_license_plate_province_full" json:"vehicle_license_plate_province_full"`
	VehicleDepartmentDeptSapShort    string                 `gorm:"column:vehicle_department_dept_sap_short" json:"vehicle_department_dept_sap_short"`
	WorkPlace                        string                 `gorm:"column:work_place" json:"work_place"`
	ReserveStartDatetime             TimeWithZone           `gorm:"column:reserve_start_datetime" json:"start_datetime"`
	ReserveEndDatetime               TimeWithZone           `gorm:"column:reserve_end_datetime" json:"end_datetime"`
	RefRequestStatusCode             string                 `gorm:"column:ref_request_status_code" json:"ref_request_status_code"`
	RefRequestStatusName             string                 `json:"ref_request_status_name"`
	IsHaveSubRequest                 string                 `gorm:"column:is_have_sub_request" json:"is_have_sub_request" example:"0"`
	ReceivedKeyPlace                 string                 `gorm:"column:appointment_key_handover_place" json:"received_key_place" example:"Main Office"`
	ReceivedKeyStartDatetime         TimeWithZone           `gorm:"column:appointment_key_handover_start_datetime" json:"received_key_start_datetime" swaggertype:"string" example:"2025-02-16T08:00:00Z"`
	ReceivedKeyEndDatetime           TimeWithZone           `gorm:"column:appointment_key_handover_end_datetime" json:"received_key_end_datetime" swaggertype:"string" example:"2025-02-16T09:30:00Z"`
	CanceledRequestDatetime          TimeWithZone           `gorm:"column:canceled_request_datetime" json:"canceled_request_datetime"`
	IsPEAEmployeeDriver              string                 `gorm:"column:is_pea_employee_driver" json:"is_pea_employee_driver"`
	CarpoolName                      string                 `gorm:"column:vehicle_carpool_name" json:"vehicle_carpool_name"`
	WorkDescription                  string                 `gorm:"column:work_description" json:"work_description"`
	ActionDetail                     string                 `gorm:"column:action_detail" json:"action_detail"`
	CanPickupButton                  bool                   `gorm:"-" json:"can_pickup_button"`
	CanScoreButton                   bool                   `gorm:"-" json:"can_score_button"`
	CanTravelCardButton              bool                   `gorm:"-" json:"can_travel_card_button"`
	TravelDetails                    []VmsTrnTripDetailList `gorm:"foreignKey:TrnRequestUID;references:TrnRequestUid" json:"-"`
	TripType                         int                    `gorm:"column:ref_trip_type_code" json:"trip_type" example:"1"`
	TripTypeName                     string                 `gorm:"-" json:"trip_type_name" example:"1"`
}

func (VmsTrnRequestList) TableName() string {
	return "public.vms_trn_request"
}

type VmsTrnRequestSummary struct {
	RefRequestStatusCode string `gorm:"column:ref_request_status_code" json:"ref_request_status_code"`
	RefRequestStatusName string `json:"ref_request_status_name"`
	Count                int    `gorm:"column:count" json:"count"`
}

// VmsTrnRequestRequest
type VmsTrnRequestRequest struct {
	TrnRequestUID        string `gorm:"column:trn_request_uid" json:"-"`
	RequestNo            string `gorm:"column:request_no" json:"request_no"`
	RefRequestStatusCode string `gorm:"column:ref_request_status_code" json:"-"`
	RefRequestTypeCode   int    `gorm:"column:ref_request_type_code" json:"-"`
	IsHaveSubRequest     string `gorm:"column:is_have_sub_request" json:"-" example:"0"`

	CreatedRequestDatetime      TimeWithZone `gorm:"column:created_request_datetime" json:"-"`
	CreatedRequestEmpID         string       `gorm:"column:created_request_emp_id" json:"-"`
	CreatedRequestEmpName       string       `gorm:"column:created_request_emp_name" json:"-"`
	CreatedRequestDeskPhone     string       `gorm:"column:created_request_desk_phone" json:"-"`
	CreatedRequestMobilePhone   string       `gorm:"column:created_request_mobile_phone" json:"-"`
	CreatedRequestPosition      string       `gorm:"column:created_request_position" json:"-"`
	CreatedRequestDeptSAP       string       `gorm:"column:created_request_dept_sap" json:"-"`
	CreatedRequestDeptNameShort string       `gorm:"column:created_request_dept_name_short" json:"-"`
	CreatedRequestDeptNameFull  string       `gorm:"column:created_request_dept_name_full" json:"-"`
	CreatedRequestRemark        string       `gorm:"column:created_request_remark" json:"-"`
	//Step1
	VehicleUserEmpID         string `gorm:"column:vehicle_user_emp_id" json:"vehicle_user_emp_id" example:"990001"`
	VehicleUserEmpName       string `gorm:"column:vehicle_user_emp_name" json:"-"`
	VehicleUserDeptSAP       string `gorm:"column:vehicle_user_dept_sap" json:"-"`
	VehicleUserDeskPhone     string `gorm:"column:vehicle_user_desk_phone" json:"car_user_internal_contact_number" example:"1122"`
	VehicleUserMobilePhone   string `gorm:"column:vehicle_user_mobile_phone" json:"car_user_mobile_contact_number" example:"0987654321"`
	VehicleUserPosition      string `gorm:"column:vehicle_user_position" json:"-"`
	VehicleUserDeptSap       string `gorm:"column:vehicle_user_dept_sap" json:"-"`
	VehicleUserDeptNameShort string `gorm:"column:vehicle_user_dept_name_short" json:"-"`
	VehicleUserDeptNameFull  string `gorm:"column:vehicle_user_dept_name_full" json:"-"`

	ReserveStartDatetime TimeWithZone `gorm:"column:reserve_start_datetime" json:"start_datetime" swaggertype:"string" example:"2025-01-01T08:00:00Z"`
	ReserveEndDatetime   TimeWithZone `gorm:"column:reserve_end_datetime" json:"end_datetime" swaggertype:"string" example:"2025-01-01T10:00:00Z"`
	RefTripTypeCode      int          `gorm:"ref_trip_type_code" json:"trip_type" example:"1"`

	WorkPlace          string `gorm:"column:work_place" json:"work_place" example:"Head Office"`
	WorkDescription    string `gorm:"column:work_description" json:"work_description" example:"Business Meeting"`
	NumberOfPassengers int    `gorm:"column:number_of_passengers" json:"number_of_passengers" example:"3"`
	Remark             string `gorm:"column:remark" json:"remark" example:"Urgent request"`
	DocNo              string `gorm:"column:doc_no" json:"doc_no" example:"REF123456"`
	DocFile            string `gorm:"column:doc_file" json:"doc_file" example:"https://vms-plus.pea.co.th/files/document.pdf"`
	DocFileName        string `gorm:"column:doc_file_name" json:"doc_file_name" example:"document.pdf"`

	RefCostTypeCode int    `gorm:"column:ref_cost_type_code" json:"ref_cost_type_code" example:"1"`
	CostCenter      string `gorm:"column:cost_center" json:"cost_center" example:"B0002211"`
	WbsNo           string `gorm:"column:wbs_no" json:"wbs_no" example:"WBS12345"`
	NetworkNo       string `gorm:"column:network_no" json:"network_no" example:"NET12345"`
	ActivityNo      string `gorm:"column:activity_no" json:"activity_no" example:"A12345"`
	PmOrderNo       string `gorm:"column:pm_order_no" json:"pm_order_no" example:"PM123456"`

	//Step 2
	MasCarpoolUID           *string `gorm:"column:mas_carpool_uid" json:"mas_carpool_uid" example:"389b0f63-4195-4ece-bf35-0011c2f5f28c"`
	MasVehicleDepartmentUID *string `gorm:"column:mas_vehicle_department_uid" json:"-" example:"21d2ea5a-4ad6-4a95-a64d-73b72d43bd55"`
	RequestedVehicleType    string  `gorm:"column:requested_vehicle_type" json:"requested_vehicle_type" example:"Sedan"`
	MasVehicleUID           *string `gorm:"column:mas_vehicle_uid" json:"mas_vehicle_uid" example:"21d2ea5a-4ad6-4a95-a64d-73b72d43bd55"`

	//MasVehicleDepartmentUID string `gorm:"column:mas_vehicle_department_uid" json:"-"`
	MasVehicleEvUID       string `gorm:"column:mas_vehicle_ev_uid" json:"-"`
	VehicleOwnerDeptSAP   string `gorm:"column:vehicle_owner_dept_sap" json:"-"`
	IsAdminChooseVehicle  string `gorm:"-" json:"is_admin_choose_vehicle" example:"0"`
	IsSystemChooseVehicle string `gorm:"-" json:"is_system_choose_vehicle" example:"0"`

	//Step 3
	MasCarPoolDriverUID *string `gorm:"column:mas_carpool_driver_uid" json:"mas_carpool_driver_uid" example:"a6c8a34b-9245-49c8-a12b-45fae77a4e7d"`
	IsPEAEmployeeDriver string  `gorm:"column:is_pea_employee_driver" json:"is_pea_employee_driver" example:"1"`
	IsAdminChooseDriver string  `gorm:"-" json:"is_admin_choose_driver" example:"0"`

	DriverEmpID            string `gorm:"column:driver_emp_id" json:"driver_emp_id" example:"700001"`
	DriverEmpName          string `gorm:"column:driver_emp_name" json:"-"`
	DriverDeptSAP          string `gorm:"column:driver_emp_dept_sap" json:"-"`
	DriverEmpDeskPhone     string `gorm:"column:driver_emp_desk_phone" json:"driver_internal_contact_number" example:"1221"`
	DriverEmpMobilePhone   string `gorm:"column:driver_emp_mobile_phone" json:"driver_mobile_contact_number" example:"0987654321"`
	DriverEmpPosition      string `gorm:"column:driver_emp_position" json:"-"`
	DriverEmpDeptSAP       string `gorm:"column:driver_emp_dept_sap" json:"-"`
	DriverEmpDeptNameShort string `gorm:"column:driver_emp_dept_name_short" json:"-"`
	DriverEmpDeptNameFull  string `gorm:"column:driver_emp_dept_name_full" json:"-"`

	PickupPlace    *string      `gorm:"column:pickup_place" json:"pickup_place" example:"Main Office"`
	PickupDateTime TimeWithZone `gorm:"column:pickup_datetime" json:"pickup_datetime" swaggertype:"string" example:"2025-02-16T08:30:00Z"`

	//Step 4
	ConfirmedRequestEmpID         string `gorm:"column:confirmed_request_emp_id" json:"confirmed_request_emp_id" example:"501621"`
	ConfirmedRequestEmpName       string `gorm:"column:confirmed_request_emp_name" json:"-"`
	ConfirmedRequestDeskPhone     string `gorm:"column:cconfirmed_request_desk_phone" json:"-"`
	ConfirmedRequestMobilePhone   string `gorm:"column:confirmed_request_mobile_phone" json:"-"`
	ConfirmedRequestPosition      string `gorm:"column:confirmed_request_position" json:"-"`
	ConfirmedRequestDeptSAP       string `gorm:"column:confirmed_request_dept_sap" json:"-"`
	ConfirmedRequestDeptNameShort string `gorm:"column:confirmed_request_dept_name_short" json:"-"`
	ConfirmedRequestDeptNameFull  string `gorm:"column:confirmed_request_dept_name_full" json:"-"`
	//
	CreatedAt time.Time `gorm:"column:created_at" json:"-"`
	CreatedBy string    `gorm:"column:created_by" json:"-"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"-"`
	UpdatedBy string    `gorm:"column:updated_by" json:"-"`
	IsDeleted string    `gorm:"column:is_deleted" json:"-"`
}

func (VmsTrnRequestRequest) TableName() string {
	return "public.vms_trn_request"
}

type VmsTrnRequestResponse struct {
	TrnRequestUID            string `gorm:"column:trn_request_uid;type:uuid;" json:"trn_request_uid"`
	RequestNo                string `gorm:"column:request_no" json:"request_no"`
	VehicleUserEmpID         string `gorm:"column:vehicle_user_emp_id" json:"vehicle_user_emp_id" example:"990001"`
	VehicleUserEmpName       string `gorm:"column:vehicle_user_emp_name" json:"vehicle_user_emp_name"`
	VehicleUserDeptSAP       string `gorm:"column:vehicle_user_dept_sap" json:"vehicle_user_dept_sap"`
	VehicleUserDeskPhone     string `gorm:"column:vehicle_user_desk_phone" json:"car_user_internal_contact_number" example:"1122"`
	VehicleUserMobilePhone   string `gorm:"column:vehicle_user_mobile_phone" json:"car_user_mobile_contact_number" example:"0987654321"`
	VehicleUserPosition      string `gorm:"column:vehicle_user_position" json:"vehicle_user_position"`
	VehicleUserDeptNameShort string `gorm:"column:vehicle_user_dept_name_short" json:"vehicle_user_dept_name_short"`
	VehicleUserDeptNameFull  string `gorm:"column:vehicle_user_dept_name_full" json:"vehicle_user_dept_name_full"`
	VehicleUserImageUrl      string `gorm:"-" json:"vehicle_user_image_url"`

	VehicleLicensePlate              string `gorm:"column:vehicle_license_plate" json:"vehicle_license_plate" example:"ABC1234"`
	VehicleLicensePlateProvinceShort string `gorm:"column:vehicle_license_plate_province_short" json:"vehicle_license_plate_province_short"`
	VehicleLicensePlateProvinceFull  string `gorm:"column:vehicle_license_plate_province_full" json:"vehicle_license_plate_province_full"`

	ReserveStartDatetime TimeWithZone   `gorm:"column:reserve_start_datetime" json:"start_datetime" swaggertype:"string" example:"2025-01-01T08:00:00Z"`
	ReserveEndDatetime   TimeWithZone   `gorm:"column:reserve_end_datetime" json:"end_datetime" swaggertype:"string" example:"2025-01-01T10:00:00Z"`
	RefTripTypeCode      *int           `gorm:"ref_trip_type_code" json:"trip_type" example:"1"`
	RefTripType          VmsRefTripType `gorm:"foreignKey:RefTripTypeCode;references:RefTripTypeCode" json:"trip_type_name"`

	WorkPlace          string `gorm:"column:work_place" json:"work_place" example:"Head Office"`
	WorkDescription    string `gorm:"column:work_description" json:"work_description" example:"Business Meeting"`
	NumberOfPassengers int    `gorm:"column:number_of_passengers" json:"number_of_passengers" example:"3"`
	Remark             string `gorm:"column:remark" json:"remark" example:"Urgent request"`
	DocNo              string `gorm:"column:doc_no" json:"doc_no" example:"REF123456"`
	DocFile            string `gorm:"column:doc_file" json:"doc_file" example:"https://vms-plus.pea.co.th/files/document.pdf"`
	DocFileName        string `gorm:"column:doc_file_name" json:"doc_file_name" example:"document.pdf"`

	NumberOfAvailableDrivers int `gorm:"-" json:"number_of_available_drivers" example:"2"`

	RefCostTypeCode int            `gorm:"column:ref_cost_type_code" json:"ref_cost_type_code" example:"1"`
	RefCostType     VmsRefCostType `gorm:"foreignKey:RefCostTypeCode;references:RefCostTypeCode" json:"cost_type"`
	CostCenter      string         `gorm:"column:cost_center" json:"cost_center" example:"B0002211"`
	WbsNo           string         `gorm:"column:wbs_no" json:"wbs_no" example:"WBS12345"`
	NetworkNo       string         `gorm:"column:network_no" json:"network_no" example:"NET12345"`
	ActivityNo      string         `gorm:"column:activity_no" json:"activity_no" example:"A12345"`
	PmOrderNo       string         `gorm:"column:pm_order_no" json:"pm_order_no" example:"PM123456"`

	MasVehicleUID                 string        `gorm:"column:mas_vehicle_uid;type:uuid" json:"mas_vehicle_uid"`
	VehicleDepartmentDeptSap      string        `gorm:"column:vehicle_department_dept_sap" json:"vehicle_department_dept_sap"`
	VehicleDepartmentDeptSapShort string        `gorm:"column:vehicle_department_dept_sap_short" json:"mas_vehicle_department_dept_sap_short"`
	VehicleDepartmentDeptSapFull  string        `gorm:"column:vehicle_department_dept_sap_full" json:"mas_vehicle_department_dept_sap_full"`
	MasVehicle                    VmsMasVehicle `gorm:"foreignKey:MasVehicleUID;references:MasVehicleUID" json:"vehicle"`

	IsAdminChooseVehicle  string `gorm:"-" json:"is_admin_choose_vehicle" example:"0"`
	IsSystemChooseVehicle string `gorm:"-" json:"is_system_choose_vehicle" example:"0"`
	RequestedVehicleType  string `gorm:"column:requested_vehicle_type" json:"requested_vehicle_type" example:"Sedan"`

	IsPEAEmployeeDriver string `gorm:"column:is_pea_employee_driver" json:"is_pea_employee_driver" example:"1"`
	IsAdminChooseDriver string `gorm:"column:is_admin_choose_driver" json:"is_admin_choose_driver" example:"1"`

	MasCarpoolDriverUID    string       `gorm:"column:mas_carpool_driver_uid;type:uuid" json:"mas_carpool_driver_uid"`
	MasDriver              VmsMasDriver `gorm:"foreignKey:MasCarpoolDriverUID;references:MasDriverUID" json:"driver"`
	IsUseDriver            bool         `gorm:"-" json:"is_use_driver"`
	DriverEmpID            string       `gorm:"column:driver_emp_id" json:"driver_emp_id" example:"700001"`
	DriverEmpName          string       `gorm:"column:driver_emp_name" json:"driver_emp_name" example:"John Doe"`
	DriverEmpPosition      string       `gorm:"column:driver_emp_position" json:"driver_emp_position" example:""`
	DriverEmpDeptSAP       string       `gorm:"column:driver_emp_dept_sap" json:"driver_emp_dept_sap" example:"DPT001"`
	DriverEmpDeptNameShort string       `gorm:"column:driver_emp_dept_name_short" json:"driver_emp_dept_name_short"`
	DriverEmpDeptNameFull  string       `gorm:"column:driver_emp_dept_name_full" json:"driver_emp_dept_name_full"`
	DriverEmpImageUrl      string       `gorm:"-" json:"driver_emp_image_url"`
	DriverInternalContact  string       `gorm:"column:driver_emp_desk_phone" json:"driver_internal_contact_number" example:"1234567890"`
	DriverMobileContact    string       `gorm:"column:driver_emp_mobile_phone" json:"driver_mobile_contact_number" example:"0987654321"`
	DriverImageURL         string       `gorm:"-" json:"driver_image_url"`
	PickupPlace            string       `gorm:"column:pickup_place" json:"pickup_place" example:"Main Office"`
	PickupDateTime         TimeWithZone `gorm:"column:pickup_datetime" json:"pickup_datetime" swaggertype:"string" example:"2025-02-16T08:30:00Z"`

	ReceivedKeyPlace         string       `gorm:"column:appointment_key_handover_place" json:"received_key_place" example:"Main Office"`
	ReceivedKeyStartDatetime TimeWithZone `gorm:"column:appointment_key_handover_start_datetime" json:"received_key_start_datetime" swaggertype:"string" example:"2025-02-16T08:00:00Z"`
	ReceivedKeyEndDatetime   TimeWithZone `gorm:"column:appointment_key_handover_end_datetime" json:"received_key_end_datetime" swaggertype:"string" example:"2025-02-16T09:30:00Z"`

	ConfirmedRequestEmpID         string `gorm:"column:confirmed_request_emp_id" json:"confirmed_request_emp_id" example:"501621"`
	ConfirmedRequestEmpName       string `gorm:"column:confirmed_request_emp_name" json:"confirmed_request_emp_name"`
	ConfirmedRequestDeskPhone     string `gorm:"column:cconfirmed_request_desk_phone" json:"confirmed_request_desk_phone"`
	ConfirmedRequestMobilePhone   string `gorm:"column:confirmed_request_mobile_phone" json:"confirmed_request_mobile_phone"`
	ConfirmedRequestPosition      string `gorm:"column:confirmed_request_position" json:"confirmed_request_position"`
	ConfirmedRequestDeptSAP       string `gorm:"column:confirmed_request_dept_sap" json:"confirmed_request_dept_sap"`
	ConfirmedRequestDeptNameShort string `gorm:"column:confirmed_request_dept_name_short" json:"confirmed_request_dept_name_short"`
	ConfirmedRequestDeptNameFull  string `gorm:"column:confirmed_request_dept_name_full" json:"confirmed_request_dept_name_full"`
	ConfirmedRequestImageUrl      string `gorm:"-" json:"confirmed_request_image_url"`

	CanCancelRequest         bool                    `gorm:"-" json:"can_cancel_request"`
	CancelReqeustEmpID       string                  `gorm:"column:cancel_reqeust_emp_id" json:"cancel_reqeust_emp_id"`
	CanceledRequestDatetime  TimeWithZone            `gorm:"canceled_request_datetime" json:"canceled_request_datetime"`
	CanceledRequestRole      string                  `gorm:"-" json:"canceled_request_role"`
	RefRequestStatusCode     string                  `gorm:"column:ref_request_status_code" json:"ref_request_status_code"`
	RefRequestStatus         VmsRefRequestStatus     `gorm:"foreignKey:RefRequestStatusCode;references:RefRequestStatusCode" json:"ref_request_status"`
	RefRequestStatusName     string                  `json:"ref_request_status_name"`
	RejectRequestReason      string                  `gorm:"column:rejected_request_reason;" json:"rejected_request_reason" example:"Test Send Back"`
	CanceledRequestReason    string                  `gorm:"column:canceled_request_reason;" json:"canceled_request_reason" example:"Test Cancel"`
	ProgressRequestStatus    []ProgressRequestStatus `gorm:"-" json:"progress_request_status"`
	ProgressRequestStatusEmp `gorm:"-" json:"progress_request_status_emp"`
	MasCarpoolUID            *string `gorm:"column:mas_carpool_uid" json:"mas_carpool_uid" example:"389b0f63-4195-4ece-bf35-0011c2f5f28c"`
	CarpoolName              string  `gorm:"column:carpool_name" json:"carpool_name"`
	CanChooseVehicle         bool    `gorm:"-" json:"can_choose_vehicle"`
	CanChooseDriver          bool    `gorm:"-" json:"can_choose_driver"`
}

func (VmsTrnRequestResponse) TableName() string {
	return "public.vms_trn_request"
}

type VmsTrnRequestRequestNo struct {
	RequestNo string `gorm:"column:request_no" json:"request_no"`
}

type ProgressRequestStatus struct {
	ProgressIcon     string       `gorm:"column:progress_icon" json:"progress_icon"`
	ProgressName     string       `gorm:"column:progress_name" json:"progress_name"`
	ProgressDatetime TimeWithZone `gorm:"column:progress_datetime" json:"progress_datetime"`
}

type ProgressRequestHistory struct {
	ProgressIcon     string       `gorm:"column:progress_icon" json:"progress_icon"`
	ProgressName     string       `gorm:"column:progress_name" json:"progress_name"`
	ProgressDatetime TimeWithZone `gorm:"column:progress_datetime" json:"progress_datetime"`
}

// VmsTrnRequestVehicleUser
type VmsTrnRequestVehicleUser struct {
	TrnRequestUID            string `gorm:"column:trn_request_uid;primarykey" json:"trn_request_uid" example:"0b07440c-ab04-49d0-8730-d62ce0a9bab9"`
	VehicleUserEmpID         string `gorm:"column:vehicle_user_emp_id" json:"vehicle_user_emp_id" example:"990001"`
	VehicleUserEmpName       string `gorm:"column:vehicle_user_emp_name" json:"-"`
	VehicleUserDeptSAP       string `gorm:"column:vehicle_user_dept_sap" json:"-"`
	VehicleUserDeskPhone     string `gorm:"column:vehicle_user_desk_phone" json:"car_user_internal_contact_number" example:"1122"`
	VehicleUserMobilePhone   string `gorm:"column:vehicle_user_mobile_phone" json:"car_user_mobile_contact_number" example:"0987654321"`
	VehicleUserPosition      string `gorm:"column:vehicle_user_position" json:"-"`
	VehicleUserDeptSap       string `gorm:"column:vehicle_user_dept_sap" json:"-"`
	VehicleUserDeptNameShort string `gorm:"column:vehicle_user_dept_name_short" json:"-"`
	VehicleUserDeptNameFull  string `gorm:"column:vehicle_user_dept_name_full" json:"-"`

	UpdatedAt time.Time `gorm:"column:updated_at" json:"-"`
	UpdatedBy string    `gorm:"column:updated_by" json:"-"`
}

func (VmsTrnRequestVehicleUser) TableName() string {
	return "public.vms_trn_request"
}

// VmsTrnRequestTrip
type VmsTrnRequestTrip struct {
	TrnRequestUID string `gorm:"column:trn_request_uid;primarykey" json:"trn_request_uid" example:"0b07440c-ab04-49d0-8730-d62ce0a9bab9"`

	ReserveStartDatetime TimeWithZone `gorm:"column:reserve_start_datetime" json:"start_datetime" swaggertype:"string" example:"2025-01-01T08:00:00Z"`
	ReserveEndDatetime   TimeWithZone `gorm:"column:reserve_end_datetime" json:"end_datetime" swaggertype:"string" example:"2025-01-01T10:00:00Z"`
	RefTripTypeCode      int          `gorm:"ref_trip_type_code" json:"trip_type" example:"1"`
	WorkPlace            string       `gorm:"column:work_place" json:"work_place" example:"Head Office"`
	WorkDescription      string       `gorm:"column:work_description" json:"work_description" example:"Business Meeting"`
	NumberOfPassengers   int          `gorm:"column:number_of_passengers" json:"number_of_passengers" example:"3"`
	Remark               string       `gorm:"column:remark" json:"remark" example:"Urgent request"`
	UpdatedAt            time.Time    `gorm:"column:updated_at" json:"-"`
	UpdatedBy            string       `gorm:"column:updated_by" json:"-"`
}

func (VmsTrnRequestTrip) TableName() string {
	return "public.vms_trn_request"
}

// VmsTrnRequestPickup
type VmsTrnRequestPickup struct {
	TrnRequestUID  string       `gorm:"column:trn_request_uid;primarykey" json:"trn_request_uid" example:"0b07440c-ab04-49d0-8730-d62ce0a9bab9"`
	PickupPlace    string       `gorm:"column:pickup_place" json:"pickup_place" example:"Main Office"`
	PickupDateTime TimeWithZone `gorm:"column:pickup_datetime" json:"pickup_datetime" swaggertype:"string" example:"2025-02-16T08:30:00Z"`
	UpdatedAt      time.Time    `gorm:"column:updated_at" json:"-"`
	UpdatedBy      string       `gorm:"column:updated_by" json:"-"`
}

func (VmsTrnRequestPickup) TableName() string {
	return "public.vms_trn_request"
}

// VmsTrnRequestDocument
type VmsTrnRequestDocument struct {
	TrnRequestUID string    `gorm:"column:trn_request_uid;primarykey" json:"trn_request_uid" example:"0b07440c-ab04-49d0-8730-d62ce0a9bab9"`
	DocNo         string    `gorm:"column:doc_no" json:"doc_no" example:"REF123456"`
	DocFile       string    `gorm:"column:doc_file" json:"doc_file" example:"https://vms-plus.pea.co.th/files/document.pdf"`
	DocFileName   string    `gorm:"column:doc_file_name" json:"doc_file_name" example:"document.pdf"`
	UpdatedAt     time.Time `gorm:"column:updated_at" json:"-"`
	UpdatedBy     string    `gorm:"column:updated_by" json:"-"`
}

func (VmsTrnRequestDocument) TableName() string {
	return "public.vms_trn_request"
}

// VmsTrnRequestCost
type VmsTrnRequestCost struct {
	TrnRequestUID   string `gorm:"column:trn_request_uid;primarykey" json:"trn_request_uid" example:"0b07440c-ab04-49d0-8730-d62ce0a9bab9"`
	RefCostTypeCode int    `gorm:"column:ref_cost_type_code" json:"ref_cost_type_code" example:"1"`
	CostCenter      string `gorm:"column:cost_center" json:"cost_center" example:"B0002211"`
	WbsNo           string `gorm:"column:wbs_no" json:"wbs_no" example:"WBS12345"`
	NetworkNo       string `gorm:"column:network_no" json:"network_no" example:"NET12345"`
	ActivityNo      string `gorm:"column:activity_no" json:"activity_no" example:"A12345"`
	PmOrderNo       string `gorm:"column:pm_order_no" json:"pm_order_no" example:"PM123456"`

	UpdatedAt time.Time `gorm:"column:updated_at" json:"-"`
	UpdatedBy string    `gorm:"column:updated_by" json:"-"`
}

func (VmsTrnRequestCost) TableName() string {
	return "public.vms_trn_request"
}

// VmsTrnRequestVehicleType
type VmsTrnRequestVehicleType struct {
	TrnRequestUID        string    `gorm:"column:trn_request_uid;primarykey" json:"trn_request_uid" example:"0b07440c-ab04-49d0-8730-d62ce0a9bab9"`
	RequestedVehicleType string    `gorm:"column:requested_vehicle_type" json:"requested_vehicle_type" example:"Sedan"`
	UpdatedAt            time.Time `gorm:"column:updated_at" json:"-"`
	UpdatedBy            string    `gorm:"column:updated_by" json:"-"`
}

func (VmsTrnRequestVehicleType) TableName() string {
	return "public.vms_trn_request"
}

// type VmsTrnRequestConfirmer struct {

type VmsTrnRequestConfirmer struct {
	TrnRequestUID                 string `gorm:"column:trn_request_uid;primarykey" json:"trn_request_uid" example:"0b07440c-ab04-49d0-8730-d62ce0a9bab9"`
	ConfirmedRequestEmpID         string `gorm:"column:confirmed_request_emp_id" json:"confirmed_request_emp_id" example:"700001"`
	ConfirmedRequestEmpName       string `gorm:"column:confirmed_request_emp_name" json:"-"`
	ConfirmedRequestDeskPhone     string `gorm:"column:cconfirmed_request_desk_phone" json:"-"`
	ConfirmedRequestMobilePhone   string `gorm:"column:confirmed_request_mobile_phone" json:"-"`
	ConfirmedRequestPosition      string `gorm:"column:confirmed_request_position" json:"-"`
	ConfirmedRequestDeptSAP       string `gorm:"column:confirmed_request_dept_sap" json:"-"`
	ConfirmedRequestDeptNameShort string `gorm:"column:confirmed_request_dept_name_short" json:"-"`
	ConfirmedRequestDeptNameFull  string `gorm:"column:confirmed_request_dept_name_full" json:"-"`

	UpdatedAt time.Time `gorm:"column:updated_at" json:"-"`
	UpdatedBy string    `gorm:"column:updated_by" json:"-"`
}

func (VmsTrnRequestConfirmer) TableName() string {
	return "public.vms_trn_request"
}

// VmsTrnRequestConfirmed
type VmsTrnRequestConfirmed struct {
	TrnRequestUID        string    `gorm:"column:trn_request_uid;primaryKey" json:"trn_request_uid" example:"0b07440c-ab04-49d0-8730-d62ce0a9bab9"`
	RefRequestStatusCode string    `gorm:"column:ref_request_status_code" json:"-"`
	UpdatedAt            time.Time `gorm:"column:updated_at" json:"-"`
	UpdatedBy            string    `gorm:"column:updated_by" json:"-"`
}

func (VmsTrnRequestConfirmed) TableName() string {
	return "public.vms_trn_request"
}

// VmsTrnRequestRejected
type VmsTrnRequestRejected struct {
	TrnRequestUID                string       `gorm:"column:trn_request_uid;primarykey" json:"trn_request_uid" example:"0b07440c-ab04-49d0-8730-d62ce0a9bab9"`
	RejectedRequestReason        string       `gorm:"column:rejected_request_reason;" json:"rejected_request_reason" example:"Test Reject"`
	RefRequestStatusCode         string       `gorm:"column:ref_request_status_code" json:"-"`
	RejectedRequestDatetime      TimeWithZone `gorm:"column:rejected_request_datetime" json:"-"`
	RejectedRequestEmpID         string       `gorm:"column:rejected_request_emp_id" json:"-"`
	RejectedRequestEmpName       string       `gorm:"column:rejected_request_emp_name" json:"-"`
	RejectedRequestDeskPhone     string       `gorm:"column:rejected_request_desk_phone" json:"-"`
	RejectedRequestMobilePhone   string       `gorm:"column:rejected_request_mobile_phone" json:"-"`
	RejectedRequestPosition      string       `gorm:"column:rejected_request_position" json:"-"`
	RejectedRequestDeptSAP       string       `gorm:"column:rejected_request_dept_sap" json:"-"`
	RejectedRequestDeptNameShort string       `gorm:"column:rejected_request_dept_name_short" json:"-"`
	RejectedRequestDeptNameFull  string       `gorm:"column:rejected_request_dept_name_full" json:"-"`

	UpdatedAt time.Time `gorm:"column:updated_at" json:"-"`
	UpdatedBy string    `gorm:"column:updated_by" json:"-"`
}

func (VmsTrnRequestRejected) TableName() string {
	return "public.vms_trn_request"
}

// VmsTrnRequestResend
type VmsTrnRequestResend struct {
	TrnRequestUID        string    `gorm:"column:trn_request_uid;primarykey" json:"trn_request_uid" example:"0b07440c-ab04-49d0-8730-d62ce0a9bab9"`
	RefRequestStatusCode string    `gorm:"column:ref_request_status_code" json:"-"`
	UpdatedAt            time.Time `gorm:"column:updated_at" json:"-"`
	UpdatedBy            string    `gorm:"column:updated_by" json:"-"`
}

func (VmsTrnRequestResend) TableName() string {
	return "public.vms_trn_request"
}

// VmsTrnRequestApproved
type VmsTrnRequestApproved struct {
	TrnRequestUID                string       `gorm:"column:trn_request_uid;primaryKey" json:"trn_request_uid" example:"0b07440c-ab04-49d0-8730-d62ce0a9bab9"`
	RefRequestStatusCode         string       `gorm:"column:ref_request_status_code" json:"-"`
	ApprovedRequestEmpID         string       `gorm:"column:approved_request_emp_id" json:"-"`
	ApprovedRequestEmpName       string       `gorm:"column:approved_request_emp_name" json:"-"`
	ApprovedRequestDeskPhone     string       `gorm:"column:approved_request_desk_phone" json:"-"`
	ApprovedRequestMobilePhone   string       `gorm:"column:approved_request_mobile_phone" json:"-"`
	ApprovedRequestPosition      string       `gorm:"column:approved_request_position" json:"-"`
	ApprovedRequestDeptSAP       string       `gorm:"column:approved_request_dept_sap" json:"-"`
	ApprovedRequestDeptNameShort string       `gorm:"column:approved_request_dept_name_short" json:"-"`
	ApprovedRequestDeptNameFull  string       `gorm:"column:approved_request_dept_name_full" json:"-"`
	ApprovedRequestDatetime      TimeWithZone `gorm:"column:approved_request_datetime" json:"-"`
	UpdatedAt                    time.Time    `gorm:"column:updated_at" json:"-"`
	UpdatedBy                    string       `gorm:"column:updated_by" json:"-"`
}

func (VmsTrnRequestApproved) TableName() string {
	return "public.vms_trn_request"
}

// VmsTrnRequestCanceled
type VmsTrnRequestCanceled struct {
	TrnRequestUID                string       `gorm:"column:trn_request_uid;primarykey" json:"trn_request_uid" example:"0b07440c-ab04-49d0-8730-d62ce0a9bab9"`
	CanceledRequestReason        string       `gorm:"column:canceled_request_reason;" json:"canceled_request_reason" example:"Test Cancel"`
	CanceledRequestEmpID         string       `gorm:"column:canceled_request_emp_id" json:"-"`
	CanceledRequestEmpName       string       `gorm:"column:canceled_request_emp_name" json:"-"`
	CanceledRequestDeskPhone     string       `gorm:"column:canceled_request_desk_phone" json:"-"`
	CanceledRequestMobilePhone   string       `gorm:"column:canceled_request_mobile_phone" json:"-"`
	CanceledRequestPosition      string       `gorm:"column:canceled_request_position" json:"-"`
	CanceledRequestDeptSAP       string       `gorm:"column:canceled_request_dept_sap" json:"-"`
	CanceledRequestDeptNameShort string       `gorm:"column:canceled_request_dept_name_short" json:"-"`
	CanceledRequestDeptNameFull  string       `gorm:"column:canceled_request_dept_name_full" json:"-"`
	CanceledRequestDatetime      TimeWithZone `gorm:"column:canceled_request_datetime" json:"-"`
	RefRequestStatusCode         string       `gorm:"column:ref_request_status_code" json:"-"`
	UpdatedAt                    time.Time    `gorm:"column:updated_at" json:"-"`
	UpdatedBy                    string       `gorm:"column:updated_by" json:"-"`
}

func (VmsTrnRequestCanceled) TableName() string {
	return "public.vms_trn_request"
}

// VmsTrnRequestVehicle
type VmsTrnRequestVehicle struct {
	TrnRequestUID           string    `gorm:"column:trn_request_uid;primarykey" json:"trn_request_uid" example:"0b07440c-ab04-49d0-8730-d62ce0a9bab9"`
	MasVehicleUID           string    `gorm:"column:mas_vehicle_uid" json:"mas_vehicle_uid"  example:"a6c8a34b-9245-49c8-a12b-45fae77a4e7d"`
	MasVehicleDepartmentUID string    `gorm:"column:mas_vehicle_department_uid" json:"-"`
	UpdatedAt               time.Time `gorm:"column:updated_at" json:"-"`
	UpdatedBy               string    `gorm:"column:updated_by" json:"-"`
}

func (VmsTrnRequestVehicle) TableName() string {
	return "public.vms_trn_request"
}

// VmsTrnRequestApprovedWithRecieiveKey
type VmsTrnRequestApprovedWithRecieiveKey struct {
	HandoverUID              string       `gorm:"column:handover_uid;primaryKey" json:"-"`
	TrnRequestUID            string       `gorm:"column:trn_request_uid" json:"trn_request_uid" example:"0b07440c-ab04-49d0-8730-d62ce0a9bab9"`
	ReceivedKeyPlace         string       `gorm:"column:appointment_location" json:"received_key_place" example:"Main Office"`
	ReceivedKeyStartDatetime TimeWithZone `gorm:"column:appointment_start" json:"received_key_start_datetime" swaggertype:"string" example:"2025-02-16T08:00:00Z"`
	ReceivedKeyEndDatetime   TimeWithZone `gorm:"column:appointment_end" json:"received_key_end_datetime" swaggertype:"string" example:"2025-02-16T09:30:00Z"`
	ReceiverType             int          `gorm:"column:receiver_type" json:"receiver_type" example:"0"`
	ApprovedRequestEmpID     string       `gorm:"-" json:"approved_request_emp_id" example:"700001"`

	CreatedBy string    `gorm:"column:created_by" json:"-"`
	CreatedAt time.Time `gorm:"column:created_at" json:"-"`
	UpdatedBy string    `gorm:"column:updated_by" json:"-"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"-"`
}

func (VmsTrnRequestApprovedWithRecieiveKey) TableName() string {
	return "public.vms_trn_vehicle_key_handover"
}

// VmsTrnRequestDriver
type VmsTrnRequestDriver struct {
	TrnRequestUID          string    `gorm:"column:trn_request_uid;primarykey" json:"trn_request_uid" example:"0b07440c-ab04-49d0-8730-d62ce0a9bab9"`
	MasCarPoolDriverUID    string    `gorm:"column:mas_carpool_driver_uid" json:"mas_carpool_driver_uid" example:"a6c8a34b-9245-49c8-a12b-45fae77a4e7d"`
	DriverEmpID            string    `gorm:"column:driver_emp_id" json:"driver_emp_id" example:"700001"`
	DriverEmpName          string    `gorm:"column:driver_emp_name" json:"driver_emp_name" example:"John Doe"`
	DriverEmpDeptSAP       string    `gorm:"column:driver_emp_dept_sap" json:"driver_emp_dept_sap" example:"B0002211"`
	DriverEmpDeptNameShort string    `gorm:"column:driver_emp_dept_name_short" json:"driver_emp_dept_name_short" example:"B0002211"`
	DriverEmpDeptNameFull  string    `gorm:"column:driver_emp_dept_name_full" json:"driver_emp_dept_name_full" example:"B0002211"`
	UpdatedAt              time.Time `gorm:"column:updated_at" json:"-"`
	UpdatedBy              string    `gorm:"column:updated_by" json:"-"`
}

func (VmsTrnRequestDriver) TableName() string {
	return "public.vms_trn_request"
}

// VmsTrnRequestVehicleInfo
type VmsTrnRequestVehicleInfo struct {
	NumberOfAvailableDrivers int `gorm:"-" json:"number_of_available_drivers" example:"2"`
}
