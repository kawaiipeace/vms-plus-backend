package models

import (
	"time"
)

type VmsTrnRequestVehicleInUseList struct {
	TrnRequestUID                    string       `gorm:"column:trn_request_uid;primaryKey;" json:"trn_request_uid"`
	RequestNo                        string       `gorm:"column:request_no" json:"request_no"`
	VehicleUserEmpID                 string       `gorm:"column:vehicle_user_emp_id" json:"vehicle_user_emp_id"`
	VehicleUserEmpName               string       `gorm:"column:vehicle_user_emp_name" json:"vehicle_user_emp_name"`
	VehicleUserDeptSAPShort          string       `gorm:"column:vehicle_user_dept_name_short" json:"vehicle_user_dept_sap_short" example:"Finance"`
	VehicleUserPosition              string       `gorm:"column:vehicle_user_position" json:"vehicle_user_position"`
	VehicleLicensePlate              string       `gorm:"column:vehicle_license_plate" json:"vehicle_license_plate"`
	VehicleLicensePlateProvinceShort string       `gorm:"column:vehicle_license_plate_province_short" json:"vehicle_license_plate_province_short"`
	VehicleLicensePlateProvinceFull  string       `gorm:"column:vehicle_license_plate_province_full" json:"vehicle_license_plate_province_full"`
	VehicleDepartmentDeptSapShort    string       `gorm:"column:vehicle_department_dept_sap_short" json:"vehicle_department_dept_sap_short"`
	WorkPlace                        string       `gorm:"column:work_place" json:"work_place"`
	ReserveStartDatetime             TimeWithZone `gorm:"column:reserve_start_datetime" json:"start_datetime"`
	ReserveEndDatetime               TimeWithZone `gorm:"column:reserve_end_datetime" json:"end_datetime"`
	RefRequestStatusCode             string       `gorm:"column:ref_request_status_code" json:"ref_request_status_code"`
	RefRequestStatusName             string       `json:"ref_request_status_name"`
	IsHaveSubRequest                 string       `gorm:"column:is_have_sub_request" json:"is_have_sub_request" example:"0"`
	CanceledRequestDatetime          TimeWithZone `gorm:"column:canceled_request_datetime" json:"canceled_request_datetime"`
	ReceivedKeyPlace                 string       `gorm:"column:appointment_key_handover_place" json:"received_key_place" example:"Main Office"`
	ReceivedKeyStartDatetime         TimeWithZone `gorm:"column:appointment_key_handover_start_datetime" json:"received_key_start_datetime" example:"2025-02-16T08:00:00Z"`
	ReceivedKeyEndDatetime           TimeWithZone `gorm:"column:appointment_key_handover_end_datetime" json:"received_key_end_datetime" example:"2025-02-16T09:30:00Z"`
	KeyReceiverPersonalID            string       `gorm:"column:receiver_personal_id" json:"key_receiver_personal_id"`
	KeyReceiverFullName              string       `gorm:"column:receiver_fullname" json:"key_receiver_fullname"`
	KeyReceiverDeptNameShort         string       `gorm:"column:receiver_dept_name_short" json:"key_receiver_dept_name_short"`
	KeyReceiverDeptNameFull          string       `gorm:"column:receiver_dept_name_full" json:"key_receiver_dept_name_full"`
	KeyReceiverDeskPhone             string       `gorm:"column:receiver_desk_phone" json:"key_receiver_desk_phone"`
	KeyReceiverMobilePhone           string       `gorm:"column:receiver_mobile_phone" json:"key_receiver_mobile_phone"`
	KeyReceiverPosition              string       `gorm:"column:receiver_position" json:"key_receiver_position"`

	RefVehicleKeyTypeCode   int                  `gorm:"column:ref_vehicle_key_type_code" json:"ref_vehicle_key_type_code" example:"1"`
	RefVehicleKeyType       VmsRefVehicleKeyType `gorm:"foreignKey:RefVehicleKeyTypeCode;references:RefVehicleKeyTypeCode" json:"ref_vehicle_key_type"`
	ReturnedVehicleDatetime TimeWithZone         `gorm:"column:returned_vehicle_datetime" json:"returned_vehicle_datetime"`
	ReturnedVehicleRemark   string               `gorm:"column:returned_vehicle_remark" json:"returned_vehicle_remark" example:"OK"`

	ParkingPlace        string `gorm:"column:parking_place" json:"parking_place"`
	NextStartDatetime   string `gorm:"-" json:"next_start_datetime"`
	WorkDescription     string `gorm:"column:work_description" json:"work_description"`
	CanPickupButton     bool   `gorm:"-" json:"can_pickup_button"`
	CanScoreButton      bool   `gorm:"-" json:"can_score_button"`
	CanTravelCardButton bool   `gorm:"-" json:"can_travel_card_button"`

	IsPEAEmployeeDriver int    `gorm:"column:is_pea_employee_driver" json:"is_pea_employee_driver" example:"1"`
	DriverCarpoolName   string `gorm:"column:driver_carpool_name" json:"driver_carpool_name"`
	VehicleCarpoolName  string `gorm:"column:vehicle_carpool_name" json:"vehicle_carpool_name"`
	VehicleCarpoolText  string `gorm:"column:vehicle_carpool_text" json:"vehicle_carpool_text"`
	DriverDeptName      string `gorm:"column:driver_dept_name" json:"driver_dept_name"`
	VehicleDeptName     string `gorm:"column:vehicle_dept_name" json:"vehicle_dept_name"`
}

// VmsTrnRequestVehicleInUseResponse
type VmsTrnRequestVehicleInUseResponse struct {
	TrnRequestUID            string `gorm:"column:trn_request_uid;primaryKey;" json:"trn_request_uid"`
	RequestNo                string `gorm:"column:request_no" json:"request_no"`
	MasCarpoolUID            string `gorm:"column:mas_carpool_uid" json:"mas_carpool_uid"`
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

	ReserveStartDatetime TimeWithZone   `gorm:"column:reserve_start_datetime" json:"start_datetime" example:"2025-01-01T08:00:00Z"`
	ReserveEndDatetime   TimeWithZone   `gorm:"column:reserve_end_datetime" json:"end_datetime" example:"2025-01-01T10:00:00Z"`
	RefTripTypeCode      *int           `gorm:"ref_trip_type_code" json:"trip_type" example:"1"`
	RefTripType          VmsRefTripType `gorm:"foreignKey:RefTripTypeCode;references:RefTripTypeCode" json:"trip_type_name"`

	WorkPlace          string `gorm:"column:work_place" json:"work_place" example:"Head Office"`
	WorkDescription    string `gorm:"column:work_description" json:"work_description" example:"Business Meeting"`
	NumberOfPassengers int    `gorm:"column:number_of_passengers" json:"number_of_passengers" example:"3"`
	Remark             string `gorm:"column:remark" json:"remark" example:"Urgent request"`
	DocNo              string `gorm:"column:doc_no" json:"doc_no" example:"REF123456"`
	DocFile            string `gorm:"column:doc_file" json:"doc_file" example:"document.pdf"`

	NumberOfAvailableDrivers int `gorm:"-" json:"number_of_available_drivers" example:"2"`

	RefCostTypeCode int            `gorm:"column:ref_cost_type_code" json:"ref_cost_type_code" example:"1"`
	RefCostType     VmsRefCostType `gorm:"foreignKey:RefCostTypeCode;references:RefCostTypeCode" json:"cost_type"`
	CostCenter      string         `gorm:"column:cost_center" json:"cost_center" example:"B0002211"`
	WbsNo           string         `gorm:"column:wbs_no" json:"wbs_no" example:"WBS12345"`
	NetworkNo       string         `gorm:"column:network_no" json:"network_no" example:"NET12345"`
	ActivityNo      string         `gorm:"column:activity_no" json:"activity_no" example:"A12345"`
	ProjectNo       string         `gorm:"column:project_no" json:"project_no" example:"PROJ12345"`
	PmOrderNo       string         `gorm:"column:pm_order_no" json:"pm_order_no" example:"PM123456"`

	MasCarpoolDriverUID  string            `gorm:"column:mas_carpool_driver_uid;type:uuid" json:"mas_carpool_driver_uid"`
	MasDriver            VmsMasDriver      `gorm:"foreignKey:MasCarpoolDriverUID;references:MasDriverUID" json:"driver"`
	IsAdminChooseVehicle string            `gorm:"column:is_admin_choose_vehicle" json:"is_admin_choose_vehicle" example:"0"`
	RequestVehicleTypeID int               `gorm:"column:requested_vehicle_type_id" json:"requested_vehicle_type_id" example:"1"`
	RequestVehicleType   VmsRefVehicleType `gorm:"foreignKey:RequestVehicleTypeID;references:RefVehicleTypeCode" json:"request_vehicle_type"`

	DriverEmpID            string `gorm:"column:driver_emp_id" json:"driver_emp_id" example:"700001"`
	DriverEmpName          string `gorm:"column:driver_emp_name" json:"driver_emp_name" example:"John Doe"`
	DriverEmpDeptSAP       string `gorm:"column:driver_emp_dept_sap" json:"driver_emp_dept_sap"`
	DriverEmpPosition      string `gorm:"column:driver_emp_position" json:"driver_emp_position" example:""`
	DriverEmpDeptNameShort string `gorm:"column:driver_emp_dept_name_short" json:"driver_emp_dept_name_short"`
	DriverEmpDeptNameFull  string `gorm:"column:driver_emp_dept_name_full" json:"driver_emp_dept_name_full"`
	DriverInternalContact  string `gorm:"column:driver_emp_desk_phone" json:"driver_internal_contact_number" example:"1234567890"`
	DriverMobileContact    string `gorm:"column:driver_emp_mobile_phone" json:"driver_mobile_contact_number" example:"0987654321"`
	DriverImageURL         string `gorm:"-" json:"driver_image_url"`

	PickupPlace    string       `gorm:"column:pickup_place" json:"pickup_place" example:"Main Office"`
	PickupDateTime TimeWithZone `gorm:"column:pickup_datetime" json:"pickup_datetime" example:"2025-02-16T08:30:00Z"`

	MasVehicleUID                 string        `gorm:"column:mas_vehicle_uid;type:uuid" json:"mas_vehicle_uid"`
	VehicleDepartmentDeptSap      string        `gorm:"column:vehicle_department_dept_sap" json:"vehicle_department_dept_sap"`
	VehicleDepartmentDeptSapShort string        `gorm:"column:vehicle_department_dept_sap_short" json:"mas_vehicle_department_dept_sap_short"`
	VehicleDepartmentDeptSapFull  string        `gorm:"column:vehicle_department_dept_sap_full" json:"mas_vehicle_department_dept_sap_full"`
	MasVehicle                    VmsMasVehicle `gorm:"foreignKey:MasVehicleUID;references:MasVehicleUID" json:"vehicle"`

	ReceivedKeyPlace         string       `gorm:"column:appointment_key_handover_place" json:"received_key_place" example:"Main Office"`
	ReceivedKeyStartDatetime TimeWithZone `gorm:"column:appointment_key_handover_start_datetime" json:"received_key_start_datetime" swaggertype:"string" example:"2025-02-16T08:00:00Z"`
	ReceivedKeyEndDatetime   TimeWithZone `gorm:"column:appointment_key_handover_end_datetime" json:"received_key_end_datetime" example:"2025-02-16T09:30:00Z"`

	RefVehicleKeyTypeCode int                  `gorm:"column:ref_vehicle_key_type_code" json:"ref_vehicle_key_type_code" example:"1"`
	ReceivedKeyDatetime   TimeWithZone         `gorm:"column:received_key_datetime" json:"received_key_datetime" example:"2025-02-16T08:00:00Z"`
	ReceiverKeyType       int                  `gorm:"column:receiver_type" json:"receiver_key_type" example:"3"`
	ReceiverKeyTypeDetail VmsRefVehicleKeyType `gorm:"foreignKey:RefVehicleKeyTypeCode;references:RefVehicleKeyTypeCode" json:"receiver_key_type_detail"`
	FleetCardNo           string               `gorm:"column:fleet_card_no" json:"fleet_card_no"`

	IsPeaEmployeeDriver string `gorm:"column:is_pea_employee_driver" json:"is_pea_employee_driver" example:"1"`

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

	ReturnedVehicleDatetime     TimeWithZone           `gorm:"column:returned_vehicle_datetime" json:"returned_vehicle_datetime" example:"2025-04-16T14:30:00Z"`
	ReturnedParkingPlace        string                 `gorm:"column:returned_parking_place" json:"returned_parking_place" example:"Parking Lot 1"`
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

	AcceptedVehicleDatetime     TimeWithZone            `gorm:"column:accepted_vehicle_datetime" json:"accepted_vehicle_datetime" example:"2025-04-16T14:30:00Z"`
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
	RejectedRequestReason       string                  `gorm:"column:rejected_request_reason;" json:"rejected_request_reason" example:"Test Send Back"`
	CanceledRequestReason       string                  `gorm:"column:canceled_request_reason;" json:"canceled_request_reason" example:"Test Cancel"`
	CanceledRequestDatetime     TimeWithZone            `gorm:"canceled_request_datetime" json:"canceled_request_datetime"`
	ProgressRequestStatus       []ProgressRequestStatus `gorm:"-" json:"progress_request_status"`
	ParkingPlace                string                  `gorm:"column:parking_place" json:"parking_place"`
	IsReturnOverDue             bool                    `gorm:"-" json:"is_return_overdue"`
	NextRequest                 VmsTrnNextRequest       `gorm:"-" json:"next_request"`

	SatisfactionSurveyAnswers []VmsTrnSatisfactionSurveyAnswersResponse `gorm:"foreignKey:TrnRequestUID;references:TrnRequestUID" json:"satisfaction_survey_answers"`
	CanPickupButton           bool                                      `gorm:"-" json:"can_pickup_button"`
	CanScoreButton            bool                                      `gorm:"-" json:"can_score_button"`
	CanTravelCardButton       bool                                      `gorm:"-" json:"can_travel_card_button"`
}

func (VmsTrnRequestVehicleInUseResponse) TableName() string {
	return "public.vms_trn_request"
}

// VmsTrnTripDetail_List
type VmsTrnTripDetailList struct {
	TrnTripDetailUID     string       `gorm:"column:trn_trip_detail_uid;primaryKey" json:"trn_trip_detail_uid" example:"123e4567-e89b-12d3-a456-426614174000"`
	TrnRequestUID        string       `gorm:"column:trn_request_uid" json:"trn_request_uid" example:"0b07440c-ab04-49d0-8730-d62ce0a9bab9"`
	TripStartDatetime    TimeWithZone `gorm:"column:trip_start_datetime" json:"trip_start_datetime" example:"2025-03-26T08:00:00Z"`
	TripEndDatetime      TimeWithZone `gorm:"column:trip_end_datetime" json:"trip_end_datetime" example:"2025-03-26T10:00:00Z"`
	TripDeparturePlace   string       `gorm:"column:trip_departure_place" json:"trip_departure_place" example:"Changi Airport"`
	TripDestinationPlace string       `gorm:"column:trip_destination_place" json:"trip_destination_place" example:"Marina Bay Sands"`
	TripStartMiles       int          `gorm:"column:trip_start_miles" json:"trip_start_miles" example:"5000"`
	TripEndMiles         int          `gorm:"column:trip_end_miles" json:"trip_end_miles" example:"5050"`
	TripDetail           string       `gorm:"column:trip_detail" json:"trip_detail" example:"Routine transport between airport and hotel."`
}

func (VmsTrnTripDetailList) TableName() string {
	return "public.vms_trn_trip_detail"
}

// VmsTrnTripDetail_Request
type VmsTrnTripDetailRequest struct {
	TrnRequestUID        string       `gorm:"column:trn_request_uid" json:"trn_request_uid" example:"0b07440c-ab04-49d0-8730-d62ce0a9bab9"`
	TripStartDatetime    TimeWithZone `gorm:"column:trip_start_datetime" json:"trip_start_datetime" swaggertype:"string" example:"2025-03-26T08:00:00Z"`
	TripEndDatetime      TimeWithZone `gorm:"column:trip_end_datetime" json:"trip_end_datetime" swaggertype:"string" example:"2025-03-26T10:00:00Z"`
	TripDeparturePlace   string       `gorm:"column:trip_departure_place" json:"trip_departure_place" example:"Changi Airport"`
	TripDestinationPlace string       `gorm:"column:trip_destination_place" json:"trip_destination_place" example:"Marina Bay Sands"`
	TripStartMiles       int          `gorm:"column:trip_start_miles" json:"trip_start_miles" example:"5000"`
	TripEndMiles         int          `gorm:"column:trip_end_miles" json:"trip_end_miles" example:"5050"`
	TripDetail           string       `gorm:"column:trip_detail" json:"trip_detail" example:"Routine transport between airport and hotel."`
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
	EmployeeOrDriverID               string    `gorm:"column:driver_emp_id" json:"employee_or_driver_id" example:"driver001"`
	CreatedAt                        time.Time `gorm:"column:created_at" json:"-"`
	CreatedBy                        string    `gorm:"column:created_by" json:"-"`
	UpdatedAt                        time.Time `gorm:"column:updated_at" json:"-"`
	UpdatedBy                        string    `gorm:"column:updated_by" json:"-"`
	IsDeleted                        string    `gorm:"column:is_deleted" json:"-"`
}

func (VmsTrnTripDetail) TableName() string {
	return "public.vms_trn_trip_detail"
}

// VmsTrnAddFuel_Request
type VmsTrnAddFuelRequest struct {
	TrnRequestUID        string       `gorm:"column:trn_request_uid" json:"trn_request_uid" example:"0b07440c-ab04-49d0-8730-d62ce0a9bab9"`
	RefOilStationBrandId int          `gorm:"column:ref_oil_station_brand_id" json:"ref_oil_station_brand_id" example:"1"`
	RefFuelTypeId        int          `gorm:"column:ref_fuel_type_id" json:"ref_fuel_type_id" example:"1"`
	Mile                 int          `gorm:"column:mile" json:"mile" example:"12000"`
	TaxInvoiceDate       TimeWithZone `gorm:"column:tax_invoice_date;type:timestamp" json:"tax_invoice_date" swaggertype:"string" example:"2025-03-26T08:00:00Z"`
	TaxInvoiceNo         string       `gorm:"column:tax_invoice_no;type:varchar(20)" json:"tax_invoice_no" example:"INV1234567890"`
	PricePerLiter        float64      `gorm:"column:price_per_liter;type:numeric(10,2)" json:"price_per_liter" example:"35.50"`
	SumLiter             float64      `gorm:"column:sum_liter;type:numeric(10,2)" json:"sum_liter" example:"50.00"`
	Vat                  float64      `gorm:"column:vat;type:numeric(10,2)" json:"vat" example:"3.50"`
	BeforeVatPrice       float64      `gorm:"column:before_vat_price;type:numeric(10,2)" json:"before_vat_price" example:"46.50"`
	SumPrice             float64      `gorm:"column:sum_price;type:numeric(10,2)" json:"sum_price" example:"1872.50"`
	ReceiptImg           string       `gorm:"column:receipt_img;type:varchar(100)" json:"receipt_img" example:"http://vms.pea.co.th/receipt.jpg"`
	RefPaymentTypeCode   int          `gorm:"column:ref_payment_type_code" json:"ref_payment_type_code" example:"1"`
}

// VmsTrnAddFuel
type VmsTrnAddFuel struct {
	TrnAddFuelUID string `gorm:"column:trn_add_fuel_uid;primaryKey" json:"trn_add_fuel_uid" example:"123e4567-e89b-12d3-a456-426614174000"`
	VmsTrnAddFuelRequest
	MasVehicleUID                    string                `gorm:"column:mas_vehicle_uid" json:"mas_vehicle_uid"`
	VehicleLicensePlate              string                `gorm:"column:vehicle_license_plate" json:"vehicle_license_plate"`
	VehicleLicensePlateProvinceShort string                `gorm:"column:vehicle_license_plate_province_short" json:"vehicle_license_plate_province_short"`
	VehicleLicensePlateProvinceFull  string                `gorm:"column:vehicle_license_plate_province_full" json:"vehicle_license_plate_province_full"`
	MasVehicleDepartmentUID          string                `gorm:"column:mas_vehicle_department_uid" json:"mas_vehicle_department_uid"`
	AddFuelDateTime                  time.Time             `gorm:"column:add_fuel_date_time" json:"add_fuel_date_time" example:"2025-03-26T08:00:00Z"`
	RefCostTypeCode                  int                   `gorm:"column:ref_cost_type_code" json:"ref_cost_type_code" example:"1"`
	RefCostType                      VmsRefCostType        `gorm:"foreignKey:RefCostTypeCode;references:RefCostTypeCode" json:"ref_cost_type"`
	RefOilStationBrandID             int                   `gorm:"column:ref_oil_station_brand_id" json:"ref_oil_station_brand_id" example:"1"`
	RefOilStationBrand               VmsRefOilStationBrand `gorm:"foreignKey:RefOilStationBrandId;references:RefOilStationBrandId" json:"ref_oil_station_brand"`
	RefFuelTypeID                    int                   `gorm:"column:ref_fuel_type_id" json:"ref_fuel_type_id" example:"1"`
	RefFuelType                      VmsRefFuelType        `gorm:"foreignKey:RefFuelTypeID;references:RefFuelTypeID" json:"ref_fuel_type"`
	RefPaymentTypeCode               int                   `gorm:"column:ref_payment_type_code" json:"ref_payment_type_code" example:"1"`
	RefPaymentType                   VmsRefPaymentType     `gorm:"foreignKey:RefPaymentTypeCode;references:RefPaymentTypeCode" json:"ref_payment_type"`
	CreatedAt                        time.Time             `gorm:"column:created_at" json:"-"`
	CreatedBy                        string                `gorm:"column:created_by" json:"-"`
	UpdatedAt                        time.Time             `gorm:"column:updated_at" json:"-"`
	UpdatedBy                        string                `gorm:"column:updated_by" json:"-"`
	IsDeleted                        string                `gorm:"column:is_deleted" json:"-"`
}

func (VmsTrnAddFuel) TableName() string {
	return "public.vms_trn_add_fuel"
}

// VmsTrnSatisfactionSurveyAnswers
type VmsTrnSatisfactionSurveyAnswers struct {
	TrnSatisfactionSurveyAnswersUID    string       `gorm:"column:trn_satisfaction_survey_answers_uid;primaryKey" json:"-"`
	TrnRequestUID                      string       `gorm:"column:trn_request_uid" json:"-"`
	MasSatisfactionSurveyQuestionsCode string       `gorm:"column:mas_satisfaction_survey_questions_uid" json:"mas_satisfaction_survey_questions_code" example:"1"`
	SurveyAnswerScore                  int          `gorm:"column:survey_answer_score" json:"survey_answer" example:"5"`
	SurveyAnswerDate                   TimeWithZone `gorm:"column:survey_answer_date" json:"-"`
	SurveyAnswerEmpID                  string       `gorm:"column:survey_answer_emp_id" json:"-"`
	DriverID                           string       `gorm:"column:driver_id" json:"-"`
}

func (VmsTrnSatisfactionSurveyAnswers) TableName() string {
	return "public.vms_trn_satisfaction_survey_answers"
}

// VmsTrnReturnedVehicle
type VmsTrnReturnedVehicle struct {
	TrnRequestUID                string                 `gorm:"column:trn_request_uid;primaryKey" json:"trn_request_uid" example:"0b07440c-ab04-49d0-8730-d62ce0a9bab9"`
	ReturnedVehicleDatetime      TimeWithZone           `gorm:"column:returned_vehicle_datetime" json:"returned_vehicle_datetime" swaggertype:"string" example:"2025-04-16T14:30:00Z"`
	MileEnd                      int                    `gorm:"column:mile_end" json:"mile_end" example:"12000"`
	FuelEnd                      int                    `gorm:"column:fuel_end" json:"fuel_end" example:"70"`
	ReturnedCleanlinessLevel     int                    `gorm:"column:ref_cleanliness_code" json:"returned_cleanliness_level" example:"1"`
	ReturnedParkingPlace         string                 `gorm:"column:returned_parking_place" json:"returned_vehicle_parking" example:"Parking Lot 1"`
	ReturnedVehicleRemark        string                 `gorm:"column:returned_vehicle_remark" json:"returned_vehicle_remark" example:"OK"`
	VehicleImages                []VehicleImageReturned `gorm:"foreignKey:TrnRequestUID;references:TrnRequestUID" json:"vehicle_images"`
	ReturnedVehicleEmpID         string                 `gorm:"column:returned_vehicle_emp_id" json:"returned_vehicle_emp_id" example:"700001"`
	ReturnedVehicleEmpName       string                 `gorm:"column:returned_vehicle_emp_name" json:"-"`
	ReturnedVehicleDeptSAP       string                 `gorm:"column:returned_vehicle_dept_sap" json:"-"`
	ReturnedVehicleDeptNameShort string                 `gorm:"column:returned_vehicle_dept_name_short" json:"-"`
	ReturnedVehicleDeptNameFull  string                 `gorm:"column:returned_vehicle_dept_name_full" json:"-"`
	RefRequestStatusCode         string                 `gorm:"column:ref_request_status_code" json:"-"`
	UpdatedAt                    time.Time              `gorm:"column:updated_at" json:"-"`
	UpdatedBy                    string                 `gorm:"column:updated_by" json:"-"`
}

func (VmsTrnReturnedVehicle) TableName() string {
	return "public.vms_trn_request"
}

// VehicleImageReturned
type VehicleImageReturned struct {
	TrnVehicleImgReturnedUID string    `gorm:"column:trn_vehicle_img_returned_uid;primaryKey" json:"-"`
	TrnRequestUID            string    `gorm:"column:trn_request_uid;" json:"-"`
	RefVehicleImgSideCode    int       `gorm:"column:ref_vehicle_img_side_code" json:"ref_vehicle_img_side_code" example:"1"`
	VehicleImgFile           string    `gorm:"column:vehicle_img_file" json:"vehicle_img_file" example:"http://vms.pea.co.th/side_image.jpg"`
	CreatedAt                time.Time `gorm:"column:created_at" json:"-"`
	CreatedBy                string    `gorm:"column:created_by" json:"-"`
	UpdatedAt                time.Time `gorm:"column:updated_at" json:"-"`
	UpdatedBy                string    `gorm:"column:updated_by" json:"-"`
	IsDeleted                string    `gorm:"column:is_deleted" json:"-"`
}

func (VehicleImageReturned) TableName() string {
	return "public.vms_trn_vehicle_img_returned"
}

type VmsTrnReceivedVehicleNoImgage struct {
	TrnRequestUID            string       `gorm:"column:trn_request_uid;primaryKey" json:"trn_request_uid" example:"0b07440c-ab04-49d0-8730-d62ce0a9bab9"`
	PickupDatetime           TimeWithZone `gorm:"column:pickup_datetime" json:"pickup_datetime" swaggertype:"string" example:"2025-03-26T14:30:00Z"`
	MileStart                int          `gorm:"column:mile_start" json:"mile_start" example:"10000"`
	FuelStart                int          `gorm:"column:fuel_start" json:"fuel_start" example:"50"`
	ReturnedCleanlinessLevel int          `gorm:"column:ref_cleanliness_code" json:"returned_cleanliness_level" example:"1"`
	ReceivedVehicleRemark    string       `gorm:"column:received_vehicle_remark" json:"received_vehicle_remark" example:"Minor scratch on bumper"`
	ReturnedVehicleRemark    string       `gorm:"column:returned_vehicle_remark" json:"returned_vehicle_remark" example:"OK"`
	UpdatedAt                time.Time    `gorm:"column:updated_at" json:"-"`
	UpdatedBy                string       `gorm:"column:updated_by" json:"-"`
}

func (VmsTrnReceivedVehicleNoImgage) TableName() string {
	return "public.vms_trn_request"
}

// VmsTrnReceivedVehicleImages
type VmsTrnReceivedVehicleImages struct {
	TrnRequestUID string                 `gorm:"column:trn_request_uid;primaryKey" json:"trn_request_uid" example:"0b07440c-ab04-49d0-8730-d62ce0a9bab9"`
	VehicleImages []VehicleImageReceived `gorm:"foreignKey:TrnRequestUID;references:TrnRequestUID" json:"vehicle_images"`
	UpdatedAt     time.Time              `gorm:"column:updated_at" json:"-"`
	UpdatedBy     string                 `gorm:"column:updated_by" json:"-"`
}

func (VmsTrnReceivedVehicleImages) TableName() string {
	return "public.vms_trn_request"
}

// VmsTrnNextRequest
type VmsTrnNextRequest struct {
	TrnRequestUid                    string `gorm:"column:trn_request_uid;type:uuid;" json:"trn_request_uid"`
	RequestNo                        string `gorm:"column:request_no" json:"request_no"`
	VehicleLicensePlate              string `gorm:"column:vehicle_license_plate" json:"vehicle_license_plate"`
	VehicleLicensePlateProvinceShort string `gorm:"column:vehicle_license_plate_province_short" json:"vehicle_license_plate_province_short"`
	VehicleLicensePlateProvinceFull  string `gorm:"column:vehicle_license_plate_province_full" json:"vehicle_license_plate_province_full"`

	VehicleUserEmpID         string         `gorm:"column:vehicle_user_emp_id" json:"vehicle_user_emp_id"`
	VehicleUserEmpName       string         `gorm:"column:vehicle_user_emp_name" json:"vehicle_user_emp_name"`
	VehicleUserDeptSAP       string         `gorm:"column:vehicle_user_dept_sap" json:"vehicle_user_dept_sap"`
	VehicleUserDeskPhone     string         `gorm:"column:vehicle_user_desk_phone" json:"car_user_internal_contact_number" example:"1122"`
	VehicleUserMobilePhone   string         `gorm:"column:vehicle_user_mobile_phone" json:"car_user_mobile_contact_number" example:"0987654321"`
	VehicleUserPosition      string         `gorm:"column:vehicle_user_position" json:"vehicle_user_position"`
	VehicleUserDeptNameShort string         `gorm:"column:vehicle_user_dept_name_short" json:"vehicle_user_dept_name_short"`
	VehicleUserDeptNameFull  string         `gorm:"column:vehicle_user_dept_name_full" json:"vehicle_user_dept_name_full"`
	WorkPlace                string         `gorm:"column:work_place" json:"work_place"`
	ReserveStartDatetime     TimeWithZone   `gorm:"column:reserve_start_datetime" json:"start_datetime"`
	ReserveEndDatetime       TimeWithZone   `gorm:"column:reserve_end_datetime" json:"end_datetime"`
	RefRequestStatusCode     string         `gorm:"column:ref_request_status_code" json:"ref_request_status_code"`
	RefRequestStatusName     string         `json:"ref_request_status_name"`
	RefTripTypeCode          int            `gorm:"ref_trip_type_code" json:"trip_type" example:"1"`
	RefTripType              VmsRefTripType `gorm:"foreignKey:RefTripTypeCode;references:RefTripTypeCode" json:"trip_type_name"`
}

func (VmsTrnNextRequest) TableName() string {
	return "public.vms_trn_request"
}
