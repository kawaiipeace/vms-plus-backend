package models

import "time"

type VmsTrnVehicleKeyHandover struct {
	HandoverUid      string       `gorm:"column:handover_uid;primaryKey;" json:"handover_uid" example:"0b07440c-ab04-49d0-8730-d62ce0a9bab9"`
	TrnRequestUID    string       `gorm:"column:trn_request_uid;" json:"trn_request_uid" example:"0b07440c-ab04-49d0-8730-d62ce0a9bab9"`
	AppointmentStart TimeWithZone `gorm:"column:appointment_start" json:"appointment_start" swaggertype:"string" example:"2025-02-16T08:00:00Z"`
	AppointmentEnd   TimeWithZone `gorm:"column:appointment_end" json:"appointment_end" swaggertype:"string" example:"2025-02-16T17:00:00Z"`
	ReceiverType     int          `gorm:"column:receiver_type" json:"-" example:"0"`
	CreatedAt        time.Time    `gorm:"column:created_at" json:"created_at"`
	CreatedBy        string       `gorm:"column:created_by" json:"created_by"`
	UpdatedAt        time.Time    `gorm:"column:updated_at" json:"-"`
	UpdatedBy        string       `gorm:"column:updated_by" json:"-"`
}

func (VmsTrnVehicleKeyHandover) TableName() string {
	return "public.vms_trn_vehicle_key_handover"
}

// VmsTrnReceivedKeyDriver
type VmsTrnReceivedKeyDriver struct {
	TrnRequestUID         string    `gorm:"column:trn_request_uid;primaryKey;" json:"trn_request_uid" example:"0b07440c-ab04-49d0-8730-d62ce0a9bab9"`
	ReceiverType          int       `gorm:"column:receiver_type" json:"-" example:"1"`
	ReceiverPersonalId    string    `gorm:"column:receiver_personal_id" json:"-"`
	ReceiverFullname      string    `gorm:"column:receiver_fullname" json:"-"`
	ReceiverDeptSAP       string    `gorm:"column:receiver_dept_sap" json:"-"`
	ReceiverDeptNameShort string    `gorm:"column:receiver_dept_name_short" json:"-"`
	ReceiverDeptNameFull  string    `gorm:"column:receiver_dept_name_full" json:"-"`
	ReceiverDeskPhone     string    `gorm:"column:receiver_desk_phone" json:"-"`
	ReceiverMobilePhone   string    `gorm:"column:receiver_mobile_phone" json:"-"`
	ReceiverPosition      string    `gorm:"column:receiver_position" json:"-"`
	UpdatedAt             time.Time `gorm:"column:updated_at" json:"-"`
	UpdatedBy             string    `gorm:"column:updated_by" json:"-"`
}

func (VmsTrnReceivedKeyDriver) TableName() string {
	return "public.vms_trn_vehicle_key_handover"
}

// VmsTrnReceivedKeyPEA
type VmsTrnReceivedKeyPEA struct {
	TrnRequestUID         string    `gorm:"column:trn_request_uid;primaryKey;" json:"trn_request_uid" example:"0b07440c-ab04-49d0-8730-d62ce0a9bab9"`
	ReceiverType          int       `gorm:"column:receiver_type" json:"-" example:"2"`
	ReceiverPersonalId    string    `gorm:"column:receiver_personal_id" json:"received_key_emp_id" example:"990001"`
	ReceiverFullname      string    `gorm:"column:receiver_fullname" json:"-"`
	ReceiverDeptSAP       string    `gorm:"column:receiver_dept_sap" json:"-"`
	ReceiverDeptNameShort string    `gorm:"column:receiver_dept_name_short" json:"-"`
	ReceiverDeptNameFull  string    `gorm:"column:receiver_dept_name_full" json:"-"`
	ReceiverDeskPhone     string    `gorm:"column:receiver_desk_phone" json:"received_key_internal_contact_number" example:"5551234"`
	ReceiverMobilePhone   string    `gorm:"column:receiver_mobile_phone" json:"received_key_mobile_contact_number" example:"0812345678"`
	ReceiverPosition      string    `gorm:"column:receiver_position" json:"-"`
	Remark                string    `gorm:"column:remark" json:"received_key_remark" example:"Employee received the key"`
	UpdatedAt             time.Time `gorm:"column:updated_at" json:"-"`
	UpdatedBy             string    `gorm:"column:updated_by" json:"-"`
}

func (VmsTrnReceivedKeyPEA) TableName() string {
	return "public.vms_trn_vehicle_key_handover"
}

// VmsTrnReceivedKeyOutSider
type VmsTrnReceivedKeyOutSider struct {
	TrnRequestUID         string    `gorm:"column:trn_request_uid;primaryKey;" json:"trn_request_uid" example:"0b07440c-ab04-49d0-8730-d62ce0a9bab9"`
	ReceiverType          int       `gorm:"column:receiver_type" json:"-" example:"3"`
	ReceiverPersonalId    string    `gorm:"column:receiver_personal_id" json:"received_key_emp_id" example:"3101000000026"`
	ReceiverFullname      string    `gorm:"column:receiver_fullname" json:"received_key_fullname" example:"Somchai Prasert"`
	ReceiverMobilePhone   string    `gorm:"column:receiver_mobile_phone" json:"received_key_mobile_contact_number" example:"0812345678"`
	ReceiverDeptSAP       string    `gorm:"column:receiver_dept_sap" json:"-"`
	ReceiverDeptNameShort string    `gorm:"column:receiver_dept_name_short" json:"-"`
	ReceiverDeptNameFull  string    `gorm:"column:receiver_dept_name_full" json:"-"`
	ReceiverDeskPhone     string    `gorm:"column:receiver_desk_phone" json:"-"`
	ReceiverPosition      string    `gorm:"column:receiver_position" json:"-"`
	Remark                string    `gorm:"column:remark" json:"received_key_remark" example:"Outsider received the key"`
	UpdatedAt             time.Time `gorm:"column:updated_at" json:"-"`
	UpdatedBy             string    `gorm:"column:updated_by" json:"-"`
}

func (VmsTrnReceivedKeyOutSider) TableName() string {
	return "public.vms_trn_vehicle_key_handover"
}

// VmsTrnRequestUpdateRecieiveKey
type VmsTrnRequestUpdateRecieivedKey struct {
	TrnRequestUID            string       `gorm:"column:trn_request_uid;primaryKey" json:"trn_request_uid" example:"0b07440c-ab04-49d0-8730-d62ce0a9bab9"`
	ReceivedKeyPlace         string       `gorm:"column:received_key_place" json:"received_key_place" example:"Main Office"`
	ReceivedKeyStartDatetime TimeWithZone `gorm:"column:received_key_start_datetime" json:"received_key_start_datetime" swaggertype:"string" example:"2025-02-16T08:00:00Z"`
	ReceivedKeyEndDatetime   TimeWithZone `gorm:"column:received_key_end_datetime" json:"received_key_end_datetime" swaggertype:"string" example:"2025-02-16T17:00:00Z"`
	UpdatedAt                time.Time    `gorm:"column:updated_at" json:"-"`
	UpdatedBy                string       `gorm:"column:updated_by" json:"-"`
}

func (VmsTrnRequestUpdateRecieivedKey) TableName() string {
	return "public.vms_trn_request"
}

// VmsTrnRequestUpdateRecieivedKeyDetail
type VmsTrnRequestUpdateRecieivedKeyDetail struct {
	TrnRequestUID         string       `gorm:"column:trn_request_uid;primaryKey" json:"trn_request_uid" example:"0b07440c-ab04-49d0-8730-d62ce0a9bab9"`
	RefVehicleKeyTypeCode int          `gorm:"column:ref_vehicle_key_type_code" json:"ref_vehicle_key_type_code" example:"1"`
	ReceivedKeyDatetime   TimeWithZone `gorm:"column:actual_receive_time" json:"received_key_datetime" swaggertype:"string" example:"2025-02-16T08:00:00Z"`
	UpdatedAt             time.Time    `gorm:"column:updated_at" json:"-"`
	UpdatedBy             string       `gorm:"column:updated_by" json:"-"`
}

func (VmsTrnRequestUpdateRecieivedKeyDetail) TableName() string {
	return "public.vms_trn_vehicle_key_handover"
}

// VmsTrnRequestUpdateRecieivedKeyConfirmed
type VmsTrnRequestUpdateRecieivedKeyConfirmed struct {
	TrnRequestUID         string       `gorm:"column:trn_request_uid;primaryKey" json:"trn_request_uid" example:"0b07440c-ab04-49d0-8730-d62ce0a9bab9"`
	RefVehicleKeyTypeCode int          `gorm:"column:ref_vehicle_key_type_code" json:"ref_vehicle_key_type_code" example:"1"`
	ReceivedKeyDatetime   TimeWithZone `gorm:"column:actual_receive_time" json:"received_key_datetime" swaggertype:"string" example:"2025-02-16T08:00:00Z"`
	UpdatedAt             time.Time    `gorm:"column:updated_at" json:"-"`
	UpdatedBy             string       `gorm:"column:updated_by" json:"-"`
}

func (VmsTrnRequestUpdateRecieivedKeyConfirmed) TableName() string {
	return "public.vms_trn_vehicle_key_handover"
}

// VmsTrnRequestUpdateRecieivedKeyConfirmed
type VmsTrnRequestUpdateRecieivedKeyStatus struct {
	TrnRequestUID                string    `gorm:"column:trn_request_uid;primaryKey" json:"trn_request_uid" example:"0b07440c-ab04-49d0-8730-d62ce0a9bab9"`
	RefRequestStatusCode         string    `gorm:"column:ref_request_status_code" json:"-"`
	ApprovedRequestEmpID         string    `gorm:"column:approved_request_emp_id" json:"-"`
	ApprovedRequestEmpName       string    `gorm:"column:approved_request_emp_name" json:"-"`
	ApprovedRequestDeptSAP       string    `gorm:"column:approved_request_dept_sap" json:"-"`
	ApprovedRequestDeptNameShort string    `gorm:"column:approved_request_dept_name_short" json:"-"`
	ApprovedRequestDeptNameFull  string    `gorm:"column:approved_request_dept_name_full" json:"-"`
	ApprovedRequestDeskPhone     string    `gorm:"column:approved_request_desk_phone" json:"-"`
	ApprovedRequestMobilePhone   string    `gorm:"column:approved_request_mobile_phone" json:"-"`
	ApprovedRequestPosition      string    `gorm:"column:approved_request_position" json:"-"`
	UpdatedAt                    time.Time `gorm:"column:updated_at" json:"-"`
	UpdatedBy                    string    `gorm:"column:updated_by" json:"-"`
}

func (VmsTrnRequestUpdateRecieivedKeyStatus) TableName() string {
	return "public.vms_trn_request"
}
