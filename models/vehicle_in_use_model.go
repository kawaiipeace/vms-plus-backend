package models

import (
	"time"
)

type VmsTrnRequestVehicleInUseList struct {
	TrnRequestUID                    string               `gorm:"column:trn_request_uid;primaryKey;" json:"trn_request_uid"`
	RequestNo                        string               `gorm:"column:request_no" json:"request_no"`
	VehicleUserEmpID                 string               `gorm:"column:vehicle_user_emp_id" json:"vehicle_user_emp_id"`
	VehicleUserEmpName               string               `gorm:"column:vehicle_user_emp_name" json:"vehicle_user_emp_name"`
	VehicleUserDeptSAPShort          string               `gorm:"column:vehicle_user_dept_sap_name_short" json:"vehicle_user_dept_sap_short" example:"Finance"`
	VehicleLicensePlate              string               `gorm:"column:vehicle_license_plate" json:"vehicle_license_plate"`
	VehicleLicensePlateProvinceShort string               `gorm:"column:vehicle_license_plate_province_short" json:"vehicle_license_plate_province_short"`
	VehicleLicensePlateProvinceFull  string               `gorm:"column:vehicle_license_plate_province_full" json:"vehicle_license_plate_province_full"`
	VehicleDepartmentDeptSapShort    string               `gorm:"column:vehicle_department_dept_sap_short" json:"vehicle_department_dept_sap_short"`
	WorkPlace                        string               `gorm:"column:work_place" json:"work_place"`
	StartDatetime                    string               `gorm:"column:start_datetime" json:"start_datetime"`
	EndDatetime                      string               `gorm:"column:end_datetime" json:"end_datetime"`
	RefRequestStatusCode             string               `gorm:"column:ref_request_status_code" json:"ref_request_status_code"`
	RefRequestStatusName             string               `json:"ref_request_status_name"`
	IsHaveSubRequest                 string               `gorm:"column:is_have_sub_request" json:"is_have_sub_request" example:"0"`
	ReceivedKeyPlace                 string               `gorm:"column:received_key_place" json:"received_key_place"`
	ReceivedKeyStartDatetime         time.Time            `gorm:"column:received_key_start_datetime" json:"received_key_start_datetime"`
	ReceivedKeyEndDatetime           time.Time            `gorm:"column:received_key_end_datetime" json:"received_key_end_datetime"`
	CanceledRequestDatetime          time.Time            `gorm:"column:canceled_request_datetime" json:"canceled_request_datetime"`
	RefVehicleKeyTypeCode            int                  `gorm:"column:ref_vehicle_key_type_code" json:"ref_vehicle_key_type_code" example:"1"`
	RefVehicleKeyType                VmsRefVehicleKeyType `gorm:"foreignKey:RefVehicleKeyTypeCode;references:RefVehicleKeyTypeCode" json:"ref_vehicle_key_type"`
	CommentOnReturnedVehicle         string               `gorm:"column:comment_on_returned_vehicle" json:"comment_on_returned_vehicle" example:"OK"`
}

// VmsTrnRequestVehicleInUseResponse
type VmsTrnRequestVehicleInUseResponse struct {
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
	RefVehicleKeyTypeCode    int       `gorm:"column:ref_vehicle_key_type_code" json:"ref_vehicle_key_type_code" example:"1"`
	ReceivedKeyDatetime      time.Time `gorm:"column:received_key_datetime" json:"received_key_datetime" example:"2025-02-16T08:00:00Z"`

	ReceiverKeyType               int                    `gorm:"column:receiver_key_type" json:"receiver_key_type" example:"3"`
	ReceivedKeyEmpID              string                 `gorm:"column:received_key_emp_id" json:"received_key_emp_id" example:"1234567890"`
	ReceivedKeyEmpName            string                 `gorm:"column:received_key_emp_name" json:"received_key_emp_name"`
	ReceivedKeyDeptSAP            string                 `gorm:"column:received_key_dept_sap" json:"received_key_dept_sap"`
	ReceivedKeyDeptSAPShort       string                 `gorm:"column:received_key_dept_sap_short" json:"received_key_dept_sap_short"`
	ReceivedKeyDeptSAPFull        string                 `gorm:"column:received_key_dept_sap_full" json:"received_key_dept_sap_full"`
	ReceivedKeyInternalContactNum string                 `gorm:"column:received_key_internal_contact_number" json:"received_key_internal_contact_number" example:"5551234"`
	ReceivedKeyMobileContactNum   string                 `gorm:"column:received_key_mobile_contact_number" json:"received_key_mobile_contact_number" example:"0812345678"`
	ReceivedKeyRemark             string                 `gorm:"column:received_key_remark" json:"received_key_remark" example:"Employee received the key"`
	ReceivedKeyImageURL           string                 `gorm:"-" json:"received_key_image_url"`
	VehicleImagesReceived         []VehicleImageReceived `gorm:"foreignKey:TrnRequestUID;references:TrnRequestUID" json:"vehicle_images_received"`
	ReceivedVehicleEmpID          string                 `gorm:"column:received_vehicle_emp_id" json:"received_vehicle_emp_id"`
	ReceivedVehicleEmpName        string                 `gorm:"column:received_vehicle_emp_name" json:"received_vehicle_emp_name"`
	ReceivedVehicleDeptSAP        string                 `gorm:"column:received_vehicle_dept_sap" json:"received_vehicle_dept_sap"`
	ReceivedVehicleDeptSAPShort   string                 `gorm:"column:received_vehicle_dept_sap_short" json:"received_vehicle_dept_sap_short"`
	ReceivedVehicleDeptSAPFull    string                 `gorm:"column:received_vehicle_dept_sap_full" json:"received_vehicle_dept_sap_full"`

	ReturnedVehicleDatetime     time.Time              `gorm:"column:returned_vehicle_datetime" json:"returned_vehicle_datetime" example:"2025-04-16T14:30:00Z"`
	MileEnd                     int                    `gorm:"column:mile_end" json:"mile_end" example:"12000"`
	FuelEnd                     int                    `gorm:"column:fuel_end" json:"fuel_end" example:"70"`
	ReturnedCleanlinessLevel    int                    `gorm:"column:returned_cleanliness_level" json:"returned_cleanliness_level" example:"1"`
	CommentOnReturnedVehicle    string                 `gorm:"column:comment_on_returned_vehicle" json:"comment_on_returned_vehicle" example:"OK"`
	VehicleImagesReturned       []VehicleImageReturned `gorm:"foreignKey:TrnRequestUID;references:TrnRequestUID" json:"vehicle_images_returned"`
	ReturnedVehicleEmpID        string                 `gorm:"column:returned_vehicle_emp_id" json:"returned_vehicle_emp_id"`
	ReturnedVehicleEmpName      string                 `gorm:"column:returned_vehicle_emp_name" json:"returned_vehicle_emp_name"`
	ReturnedVehicleDeptSAP      string                 `gorm:"column:returned_vehicle_dept_sap" json:"returned_vehicle_dept_sap"`
	ReturnedVehicleDeptSAPShort string                 `gorm:"column:returned_vehicle_dept_sap_short" json:"returned_vehicle_dept_sap_short"`
	ReturnedVehicleDeptSAPFull  string                 `gorm:"column:returned_vehicle_dept_sap_full" json:"returned_vehicle_dept_sap_full"`

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
	ProgressRequestStatus       []ProgressRequestStatus `gorm:"-" json:"progress_request_status"`
}

func (VmsTrnRequestVehicleInUseResponse) TableName() string {
	return "public.vms_trn_request"
}

// VmsTrnTripDetail_List
type VmsTrnTripDetailList struct {
	TrnTripDetailUID     string    `gorm:"column:trn_trip_detail_uid;primaryKey" json:"trn_trip_detail_uid" example:"123e4567-e89b-12d3-a456-426614174000"`
	TrnRequestUID        string    `gorm:"column:trn_request_uid" json:"trn_request_uid" example:"8bd09808-61fa-42fd-8a03-bf961b5678cd"`
	TripStartDatetime    time.Time `gorm:"column:trip_start_datetime" json:"trip_start_datetime" example:"2025-03-26T08:00:00Z"`
	TripEndDatetime      time.Time `gorm:"column:trip_end_datetime" json:"trip_end_datetime" example:"2025-03-26T10:00:00Z"`
	TripDeparturePlace   string    `gorm:"column:trip_departure_place" json:"trip_departure_place" example:"Changi Airport"`
	TripDestinationPlace string    `gorm:"column:trip_destination_place" json:"trip_destination_place" example:"Marina Bay Sands"`
	TripStartMiles       int       `gorm:"column:trip_start_miles" json:"trip_start_miles" example:"5000"`
	TripEndMiles         int       `gorm:"column:trip_end_miles" json:"trip_end_miles" example:"5050"`
	TripDetail           string    `gorm:"column:trip_detail" json:"trip_detail" example:"Routine transport between airport and hotel."`
}

func (VmsTrnTripDetailList) TableName() string {
	return "public.vms_trn_trip_detail"
}

// VmsTrnTripDetail_Request
type VmsTrnTripDetailRequest struct {
	TrnRequestUID        string    `gorm:"column:trn_request_uid" json:"trn_request_uid" example:"8bd09808-61fa-42fd-8a03-bf961b5678cd"`
	TripStartDatetime    time.Time `gorm:"column:trip_start_datetime" json:"trip_start_datetime" example:"2025-03-26T08:00:00Z"`
	TripEndDatetime      time.Time `gorm:"column:trip_end_datetime" json:"trip_end_datetime" example:"2025-03-26T10:00:00Z"`
	TripDeparturePlace   string    `gorm:"column:trip_departure_place" json:"trip_departure_place" example:"Changi Airport"`
	TripDestinationPlace string    `gorm:"column:trip_destination_place" json:"trip_destination_place" example:"Marina Bay Sands"`
	TripStartMiles       int       `gorm:"column:trip_start_miles" json:"trip_start_miles" example:"5000"`
	TripEndMiles         int       `gorm:"column:trip_end_miles" json:"trip_end_miles" example:"5050"`
	TripDetail           string    `gorm:"column:trip_detail" json:"trip_detail" example:"Routine transport between airport and hotel."`
}

// VmsTrnTripDetail
type VmsTrnTripDetail struct {
	TrnTripDetailUID string `gorm:"column:trn_trip_detail_uid;primaryKey" json:"trn_trip_detail_uid" example:"123e4567-e89b-12d3-a456-426614174000"`
	VmsTrnTripDetailRequest
	MasVehicleUID                    string    `gorm:"column:mas_vehicle_uid" json:"mas_vehicle_uid" example:"789e4567-e89b-12d3-a456-426614174002"`
	VehicleLicensePlate              string    `gorm:"column:vehicle_license_plate" json:"vehicle_license_plate" example:"SGP1234"`
	VehicleLicensePlateProvinceShort string    `gorm:"column:vehicle_license_plate_province_short" json:"vehicle_license_plate_province_short" example:"SG"`
	VehicleLicensePlateProvinceFull  string    `gorm:"column:vehicle_license_plate_province_full" json:"vehicle_license_plate_province_full" example:"Singapore"`
	MasVehicleDepartmentUID          string    `gorm:"column:mas_vehicle_department_uid" json:"mas_vehicle_department_uid" example:"abc12345-6789-1234-5678-abcdef012345"`
	MasCarpoolUID                    string    `gorm:"column:mas_carpool_uid" json:"mas_carpool_uid" example:"xyz12345-6789-1234-5678-abcdef012345"`
	EmployeeOrDriverID               string    `gorm:"column:employee_or_driver_id" json:"employee_or_driver_id" example:"driver001"`
	CreatedAt                        time.Time `gorm:"column:created_at" json:"-"`
	CreatedBy                        string    `gorm:"column:created_by" json:"-"`
	UpdatedAt                        time.Time `gorm:"column:updated_at" json:"-"`
	UpdatedBy                        string    `gorm:"column:updated_by" json:"-"`
	IsDeleted                        string    `gorm:"column:is_deleted" json:"-"`
}

func (VmsTrnTripDetail) TableName() string {
	return "public.vms_trn_trip_detail"
}

// VmsTrnAddFuel_List
type VmsTrnAddFuel_List struct {
	TrnAddFuelUid        string    `gorm:"column:trn_add_fuel_uid" json:"trn_add_fuel_uid" example:"123e4567-e89b-12d3-a456-426614174000"`
	TrnRequestUID        string    `gorm:"column:trn_request_uid" json:"trn_request_uid" example:"8bd09808-61fa-42fd-8a03-bf961b5678cd"`
	TripStartDatetime    time.Time `gorm:"column:trip_start_datetime" json:"trip_start_datetime" example:"2025-03-26T08:00:00"`
	TripEndDatetime      time.Time `gorm:"column:trip_end_datetime" json:"trip_end_datetime" example:"2025-03-26T10:00:00"`
	TripDeparturePlace   string    `gorm:"column:trip_departure_place" json:"trip_departure_place" example:"Changi Airport"`
	TripDestinationPlace string    `gorm:"column:trip_destination_place" json:"trip_destination_place" example:"Marina Bay Sands"`
	TripStartMiles       int       `gorm:"column:trip_start_miles" json:"trip_start_miles" example:"5000"`
	TripEndMiles         int       `gorm:"column:trip_end_miles" json:"trip_end_miles" example:"5050"`
	TripDetail           string    `gorm:"column:trip_detail" json:"trip_detail" example:"Routine transport between airport and hotel."`
}

// VmsTrnAddFuel_Request
type VmsTrnAddFuelRequest struct {
	TrnRequestUID        string    `gorm:"column:trn_request_uid" json:"trn_request_uid" example:"8bd09808-61fa-42fd-8a03-bf961b5678cd"`
	RefOilStationBrandId int       `gorm:"column:ref_oil_station_brand_id" json:"ref_oil_station_brand_id" example:"1"`
	RefFuelTypeId        int       `gorm:"column:ref_fuel_type_id" json:"ref_fuel_type_id" example:"1"`
	Mile                 int       `gorm:"column:mile" json:"mile" example:"12000"`
	TaxInvoiceDate       time.Time `gorm:"column:tax_invoice_date;type:timestamp" json:"tax_invoice_date" example:"2025-03-26T08:00:00Z"`
	TaxInvoiceNo         string    `gorm:"column:tax_invoice_no;type:varchar(20)" json:"tax_invoice_no" example:"INV1234567890"`
	PricePerLiter        float64   `gorm:"column:price_per_liter;type:numeric(10,2)" json:"price_per_liter" example:"35.50"`
	SumLiter             float64   `gorm:"column:sum_liter;type:numeric(10,2)" json:"sum_liter" example:"50.00"`
	SumPrice             float64   `gorm:"column:sum_price;type:numeric(10,2)" json:"sum_price" example:"1872.50"`
	ReceiptImg           string    `gorm:"column:receipt_img;type:varchar(100)" json:"receipt_img" example:"http://vms.pea.co.th/receipt.jpg"`
	RefPaymentTypeCode   int       `gorm:"column:ref_payment_type_code" json:"ref_payment_type_code" example:"1"`
}

// VmsTrnAddFuel
type VmsTrnAddFuel struct {
	TrnAddFuelUID string `gorm:"column:trn_add_fuel_uid;primaryKey" json:"trn_add_fuel_uid" example:"123e4567-e89b-12d3-a456-426614174000"`
	VmsTrnAddFuelRequest
	MasVehicleUID                    string    `gorm:"column:mas_vehicle_uid" json:"mas_vehicle_uid"`
	VehicleLicensePlate              string    `gorm:"column:vehicle_license_plate" json:"vehicle_license_plate"`
	VehicleLicensePlateProvinceShort string    `gorm:"column:vehicle_license_plate_province_short" json:"vehicle_license_plate_province_short"`
	VehicleLicensePlateProvinceFull  string    `gorm:"column:vehicle_license_plate_province_full" json:"vehicle_license_plate_province_full"`
	MasVehicleDepartmentUID          string    `gorm:"column:mas_vehicle_department_uid" json:"mas_vehicle_department_uid"`
	AddFuelDateTime                  time.Time `gorm:"column:add_fuel_date_time" json:"add_fuel_date_time" example:"2025-03-26T08:00:00Z"`
	RefCostTypeCode                  int       `gorm:"column:ref_cost_type_code" json:"ref_cost_type_code" example:"1"`
	CreatedAt                        time.Time `gorm:"column:created_at" json:"-"`
	CreatedBy                        string    `gorm:"column:created_by" json:"-"`
	UpdatedAt                        time.Time `gorm:"column:updated_at" json:"-"`
	UpdatedBy                        string    `gorm:"column:updated_by" json:"-"`
	IsDeleted                        string    `gorm:"column:is_deleted" json:"-"`
}

func (VmsTrnAddFuel) TableName() string {
	return "public.vms_trn_add_fuel"
}

// VmsTrnSatisfactionSurveyAnswers
type VmsTrnSatisfactionSurveyAnswers struct {
	TrnSatisfactionSurveyAnswersUID    string    `gorm:"column:trn_satisfaction_survey_answers_uid;primaryKey" json:"-"`
	TrnRequestUID                      string    `gorm:"column:trn_request_uid" json:"-"`
	MasSatisfactionSurveyQuestionsCode int       `gorm:"column:mas_satisfaction_survey_questions_code" json:"mas_satisfaction_survey_questions_code" example:"1"`
	SurveyAnswer                       int       `gorm:"column:survey_answer" json:"survey_answer" example:"5"`
	SurveyAnswerDate                   time.Time `gorm:"column:survey_answer_date" json:"-"`
	SurveyAnswerEmpID                  string    `gorm:"column:survey_answer_emp_id" json:"-"`
}

func (VmsTrnSatisfactionSurveyAnswers) TableName() string {
	return "public.vms_trn_satisfaction_survey_answers"
}

// VmsTrnReturnedVehicle
type VmsTrnReturnedVehicle struct {
	TrnRequestUID               string                 `gorm:"column:trn_request_uid;primaryKey" json:"trn_request_uid" example:"8bd09808-61fa-42fd-8a03-bf961b5678cd"`
	ReturnedVehicleDatetime     time.Time              `gorm:"column:returned_vehicle_datetime" json:"returned_vehicle_datetime" example:"2025-04-16T14:30:00Z"`
	MileEnd                     int                    `gorm:"column:mile_end" json:"mile_end" example:"12000"`
	FuelEnd                     int                    `gorm:"column:fuel_end" json:"fuel_end" example:"70"`
	ReturnedCleanlinessLevel    int                    `gorm:"column:returned_cleanliness_level" json:"returned_cleanliness_level" example:"1"`
	CommentOnReturnedVehicle    string                 `gorm:"column:comment_on_returned_vehicle" json:"comment_on_returned_vehicle" example:"OK"`
	VehicleImages               []VehicleImageReturned `gorm:"foreignKey:TrnRequestUID;references:TrnRequestUID" json:"vehicle_images"`
	ReturnedVehicleEmpID        string                 `gorm:"column:returned_vehicle_emp_id" json:"returned_vehicle_emp_id"`
	ReturnedVehicleEmpName      string                 `gorm:"column:returned_vehicle_emp_name" json:"-"`
	ReturnedVehicleDeptSAP      string                 `gorm:"column:returned_vehicle_dept_sap" json:"-"`
	ReturnedVehicleDeptSAPShort string                 `gorm:"column:returned_vehicle_dept_sap_short" json:"-"`
	ReturnedVehicleDeptSAPFull  string                 `gorm:"column:returned_vehicle_dept_sap_full" json:"-"`
	RefRequestStatusCode        string                 `gorm:"column:ref_request_status_code" json:"-"`
	UpdatedAt                   time.Time              `gorm:"column:updated_at" json:"-"`
	UpdatedBy                   string                 `gorm:"column:updated_by" json:"-"`
}

func (VmsTrnReturnedVehicle) TableName() string {
	return "public.vms_trn_request"
}

// VehicleImageReturned
type VehicleImageReturned struct {
	TrnVehicleImgReturnedUID string `gorm:"column:trn_vehicle_img_returned_uid;primaryKey" json:"-"`
	TrnRequestUID            string `gorm:"column:trn_request_uid;" json:"-"`
	RefVehicleImgSideCode    int    `gorm:"column:ref_vehicle_img_side_code" json:"ref_vehicle_img_side_code" example:"1"`
	VehicleImgFile           string `gorm:"column:vehicle_img_file" json:"vehicle_img_file" example:"http://vms.pea.co.th/side_image.jpg"`
}

func (VehicleImageReturned) TableName() string {
	return "public.vms_trn_vehicle_img_returned"
}

type VmsTrnReceivedVehicleNoImgage struct {
	TrnRequestUID            string    `gorm:"column:trn_request_uid;primaryKey" json:"trn_request_uid" example:"8bd09808-61fa-42fd-8a03-bf961b5678cd"`
	PickupDatetime           time.Time `gorm:"column:pickup_datetime" json:"pickup_datetime" example:"2025-03-26T14:30:00Z"`
	MileStart                int       `gorm:"column:mile_start" json:"mile_start" example:"10000"`
	FuelStart                int       `gorm:"column:fuel_start" json:"fuel_start" example:"50"`
	CommentOnReceivedVehicle string    `gorm:"column:comment_on_received_vehicle" json:"comment_on_received_vehicle" example:"Minor scratch on bumper"`
	UpdatedAt                time.Time `gorm:"column:updated_at" json:"-"`
	UpdatedBy                string    `gorm:"column:updated_by" json:"-"`
}

func (VmsTrnReceivedVehicleNoImgage) TableName() string {
	return "public.vms_trn_request"
}

// VmsTrnReceivedVehicleImages
type VmsTrnReceivedVehicleImages struct {
	TrnRequestUID string                 `gorm:"column:trn_request_uid;primaryKey" json:"trn_request_uid" example:"8bd09808-61fa-42fd-8a03-bf961b5678cd"`
	VehicleImages []VehicleImageReceived `gorm:"foreignKey:TrnRequestUID;references:TrnRequestUID" json:"vehicle_images"`
	UpdatedAt     time.Time              `gorm:"column:updated_at" json:"-"`
	UpdatedBy     string                 `gorm:"column:updated_by" json:"-"`
}

func (VmsTrnReceivedVehicleImages) TableName() string {
	return "public.vms_trn_request"
}
