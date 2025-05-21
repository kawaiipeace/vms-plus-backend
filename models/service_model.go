package models

import "time"

type VmsTrnRequesService struct {
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

	VehicleLicensePlate              string `gorm:"column:vehicle_license_plate" json:"vehicle_license_plate" example:"ABC1234"`
	VehicleLicensePlateProvinceShort string `gorm:"column:vehicle_license_plate_province_short" json:"vehicle_license_plate_province_short"`
	VehicleLicensePlateProvinceFull  string `gorm:"column:vehicle_license_plate_province_full" json:"vehicle_license_plate_province_full"`

	ReserveStartDatetime time.Time      `gorm:"column:reserve_start_datetime" json:"start_datetime" example:"2025-01-01T08:00:00Z"`
	ReserveEndDatetime   time.Time      `gorm:"column:reserve_end_datetime" json:"end_datetime" example:"2025-01-01T10:00:00Z"`
	RefTripTypeCode      int            `gorm:"ref_trip_type_code" json:"trip_type" example:"1"`
	RefTripType          VmsRefTripType `gorm:"foreignKey:RefTripTypeCode;references:RefTripTypeCode" json:"trip_type_name"`

	WorkPlace          string `gorm:"column:work_place" json:"work_place" example:"Head Office"`
	WorkDescription    string `gorm:"column:work_description" json:"objective" example:"Business Meeting"`
	NumberOfPassengers int    `gorm:"column:number_of_passengers" json:"number_of_passengers" example:"3"`
	Remark             string `gorm:"column:remark" json:"remark" example:"Urgent request"`
	DocNo              string `gorm:"column:doc_no" json:"reference_number" example:"REF123456"`
	DocFile            string `gorm:"column:doc_file" json:"attached_document" example:"document.pdf"`

	NumberOfAvailableDrivers int `gorm:"-" json:"number_of_available_drivers" example:"2"`

	RefCostTypeCode int            `gorm:"column:ref_cost_type_code" json:"ref_cost_type_code" example:"1"`
	RefCostType     VmsRefCostType `gorm:"foreignKey:RefCostTypeCode;references:RefCostTypeCode" json:"cost_type"`
	CostCenter      string         `gorm:"column:cost_center" json:"cost_center" example:"B0002211"`
	WbsNo           string         `gorm:"column:wbs_no" json:"wbs_no" example:"WBS12345"`
	NetworkNo       string         `gorm:"column:network_no" json:"network_no" example:"NET12345"`
	ProjectNo       string         `gorm:"column:project_no" json:"project_no" example:"PROJ12345"`

	MasCarpoolDriverUID  string            `gorm:"column:mas_carpool_driver_uid;type:uuid" json:"mas_carpool_driver_uid"`
	MasDriver            VmsMasDriver      `gorm:"foreignKey:MasCarpoolDriverUID;references:MasDriverUID" json:"driver"`
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
	MasVehicle                    VmsMasVehicle `gorm:"foreignKey:MasVehicleUID;references:MasVehicleUID" json:"vehicle"`

	ReceivedKeyPlace         string    `gorm:"column:appointment_key_handover_place" json:"received_key_place" example:"Main Office"`
	ReceivedKeyStartDatetime time.Time `gorm:"column:appointment_key_handover_start_datetime" json:"received_key_start_datetime" example:"2025-02-16T08:00:00Z"`
	ReceivedKeyEndDatetime   time.Time `gorm:"column:appointment_key_handover_end_datetime" json:"received_key_end_datetime" example:"2025-02-16T09:30:00Z"`

	RefVehicleKeyTypeCode int                  `gorm:"column:ref_vehicle_key_type_code" json:"ref_vehicle_key_type_code" example:"1"`
	ReceivedKeyDatetime   time.Time            `gorm:"column:received_key_datetime" json:"received_key_datetime" example:"2025-02-16T08:00:00Z"`
	ReceiverKeyType       int                  `gorm:"column:receiver_key_type" json:"receiver_key_type" example:"3"`
	ReceiverKeyTypeDetail VmsRefVehicleKeyType `gorm:"foreignKey:ReceiverKeyType;references:RefVehicleKeyTypeCode" json:"receiver_key_type_detail"`
	FleetCardNo           string               `gorm:"column:fleet_card_no" json:"fleet_card_no"`

	ReceivedKeyEmpID         string `gorm:"column:receiver_personal_id" json:"received_key_emp_id" example:"990001"`
	ReceivedKeyEmpName       string `gorm:"column:receiver_fullname" json:"received_key_emp_name"`
	ReceivedKeyDeptSAP       string `gorm:"column:receiver_dept_sap" json:"received_key_dept_sap"`
	ReceivedKeyDeptNameShort string `gorm:"column:receiver_dept_name_short" json:"received_key_dept_sap_short"`
	ReceivedKeyDeptNameFull  string `gorm:"column:receiver_dept_name_full" json:"received_key_dept_sap_full"`
	ReceivedKeyDeskPhone     string `gorm:"column:receiver_desk_phone" json:"received_key_internal_contact_number" example:"5551234"`
	ReceivedKeyMobilePhone   string `gorm:"column:receiver_mobile_phone" json:"received_key_mobile_contact_number" example:"0812345678"`
	ReceiverKeyPosition      string `gorm:"column:receiver_position" json:"received_key_position"`
	ReceivedKeyRemark        string `gorm:"column:receiver_key_remark" json:"received_key_remark" example:"Employee received the key"`
	ReceivedKeyImageURL      string `gorm:"-" json:"received_key_image_url"`

	VehicleImagesReceived       []VehicleImageReceived `gorm:"foreignKey:TrnRequestUID;references:TrnRequestUID" json:"vehicle_images_received"`
	ReceivedVehicleEmpID        string                 `gorm:"column:received_vehicle_emp_id" json:"received_vehicle_emp_id"`
	ReceivedVehicleEmpName      string                 `gorm:"column:received_vehicle_emp_name" json:"received_vehicle_emp_name"`
	ReceivedVehicleDeptSAP      string                 `gorm:"column:received_vehicle_dept_sap" json:"received_vehicle_dept_sap"`
	ReceivedVehicleDeptSAPShort string                 `gorm:"column:received_vehicle_dept_sap_short" json:"received_vehicle_dept_sap_short"`
	ReceivedVehicleDeptSAPFull  string                 `gorm:"column:received_vehicle_dept_sap_full" json:"received_vehicle_dept_sap_full"`
	MileStart                   int                    `gorm:"column:mile_start" json:"mile_start" example:"10000"`
	FuelStart                   int                    `gorm:"column:fuel_start" json:"fuel_start" example:"50"`
	ReceivedVehicleRemark       string                 `gorm:"column:received_vehicle_remark" json:"received_vehicle_remark" example:"Minor scratch on bumper"`

	ReturnedVehicleDatetime     time.Time              `gorm:"column:returned_vehicle_datetime" json:"returned_vehicle_datetime" example:"2025-04-16T14:30:00Z"`
	MileEnd                     int                    `gorm:"column:mile_end" json:"mile_end" example:"12000"`
	FuelEnd                     int                    `gorm:"column:fuel_end" json:"fuel_end" example:"70"`
	MileUsed                    int                    `gorm:"-" json:"mile_used" example:"200"`
	AddFuelsCount               int64                  `gorm:"-" json:"add_fuels_count" example:"1"`
	TripDetailsCount            int64                  `gorm:"-" json:"trip_details_count" example:"2"`
	ReturnedCleanlinessLevel    int                    `gorm:"column:ref_cleanliness_code" json:"returned_cleanliness_level" example:"1"`
	ReturnedVehicleRemark       string                 `gorm:"column:returned_vehicle_remark" json:"returned_vehicle_remark" example:"OK"`
	VehicleImagesReturned       []VehicleImageReturned `gorm:"foreignKey:TrnRequestUID;references:TrnRequestUID" json:"vehicle_images_returned"`
	ReturnedVehicleEmpID        string                 `gorm:"column:returned_vehicle_emp_id" json:"returned_vehicle_emp_id"`
	ReturnedVehicleEmpName      string                 `gorm:"column:returned_vehicle_emp_name" json:"returned_vehicle_emp_name"`
	ReturnedVehicleDeptSAP      string                 `gorm:"column:returned_vehicle_dept_sap" json:"returned_vehicle_dept_sap"`
	ReturnedVehicleDeptSAPShort string                 `gorm:"column:returned_vehicle_dept_sap_short" json:"returned_vehicle_dept_sap_short"`
	ReturnedVehicleDeptSAPFull  string                 `gorm:"column:returned_vehicle_dept_sap_full" json:"returned_vehicle_dept_sap_full"`
	VehicleImageInspect         []VehicleImageInspect  `gorm:"foreignKey:TrnRequestUID;references:TrnRequestUID" json:"vehicle_image_inspect"`

	AcceptedVehicleDatetime     time.Time               `gorm:"column:accepted_vehicle_datetime" json:"accepted_vehicle_datetime" example:"2025-04-16T14:30:00Z"`
	AcceptedVehicleEmpID        string                  `gorm:"column:accepted_vehicle_emp_id" json:"accepted_vehicle_emp_id"`
	AcceptedVehicleEmpName      string                  `gorm:"column:accepted_vehicle_emp_name" json:"accepted_vehicle_emp_name"`
	AcceptedVehicleDeptSAP      string                  `gorm:"column:accepted_vehicle_dept_sap" json:"accepted_vehicle_dept_sap"`
	AcceptedVehicleDeptSAPShort string                  `gorm:"column:accepted_vehicle_dept_sap_short" json:"accepted_vehicle_dept_sap_short"`
	AcceptedVehicleDeptSAPFull  string                  `gorm:"column:accepted_vehicle_dept_sap_full" json:"accepted_vehicle_dept_sap_full"`
	IsUseDriver                 bool                    `gorm:"column:is_use_driver" json:"is_use_driver"`
	CanCancelRequest            bool                    `gorm:"-" json:"can_cancel_request"`
	RefRequestStatusCode        string                  `gorm:"column:ref_request_status_code" json:"ref_request_status_code"`
	RefRequestStatus            VmsRefRequestStatus     `gorm:"foreignKey:RefRequestStatusCode;references:RefRequestStatusCode" json:"ref_request_status"`
	RefRequestStatusName        string                  `json:"ref_request_status_name"`
	SendedBackRequestReason     string                  `gorm:"column:sended_back_request_reason;" json:"sended_back_request_reason" example:"Test Send Back"`
	CanceledRequestReason       string                  `gorm:"column:canceled_request_reason;" json:"canceled_request_reason" example:"Test Cancel"`
	CanceledRequestDatetime     time.Time               `gorm:"canceled_request_datetime" json:"canceled_request_datetime"`
	ProgressRequestStatus       []ProgressRequestStatus `gorm:"-" json:"progress_request_status"`
	ParkingPlace                string                  `gorm:"column:parking_place" json:"parking_place"`
	IsReturnOverDue             bool                    `gorm:"-" json:"is_return_overdue"`
	NextRequest                 VmsTrnNextRequest       `gorm:"-" json:"next_request"`

	SatisfactionSurveyAnswers []VmsTrnSatisfactionSurveyAnswersResponse `gorm:"foreignKey:TrnRequestUID;references:TrnRequestUID" json:"satisfaction_survey_answers"`

	TripDetails []VmsTrnTripDetail `gorm:"foreignKey:TrnRequestUID;references:TrnRequestUID" json:"trip_details"`
	AddFuels    []VmsTrnAddFuel    `gorm:"foreignKey:TrnRequestUID;references:TrnRequestUID" json:"add_fuels"`
}

func (VmsTrnRequesService) TableName() string {
	return "public.vms_trn_request"
}

// MADE BY PEACE
type VmsToEEMS struct {
	TrnRequestUID           string 			`gorm:"column:trn_request_uid;type:uuid;" json:"trn_request_uid"`
	RequestNo               string 			`gorm:"column:request_no" json:"request_no"`
	VehicleLicensePlate     string 			`gorm:"column:vehicle_license_plate" json:"vehicle_license_plate" example:"ABC1234"`
	ReserveStartDatetime 	time.Time		`gorm:"column:reserve_start_datetime" json:"start_datetime" example:"2025-01-01T08:00:00Z"`
	ReserveEndDatetime   	time.Time  		`gorm:"column:reserve_end_datetime" json:"end_datetime" example:"2025-01-01T10:00:00Z"`
	WorkPlace          		string 			`gorm:"column:work_place" json:"work_place" example:"Head Office"`
	WorkDescription    		string 			`gorm:"column:work_description" json:"objective" example:"Business Meeting"`
	DocNo              		string 			`gorm:"column:doc_no" json:"reference_number" example:"REF123456"`
}

func (VmsToEEMS) TableName() string {
	return "public.vms_trn_request"
}