package models

import "time"

//VmsTrnReceivedKeyDriver
type VmsTrnReceivedKeyDriver struct {
	TrnRequestUID        string    `gorm:"column:trn_request_uid;primaryKey;" json:"trn_request_uid" example:"3045a994-ba0b-431d-acf2-98768a9c5fc9"`
	ReceiverKeyType      int       `gorm:"column:receiver_key_type" json:"-" example:"1"`
	RefRequestStatusCode string    `gorm:"column:ref_request_status_code" json:"ref_request_status_code"`
	UpdatedAt            time.Time `gorm:"column:updated_at" json:"-"`
	UpdatedBy            string    `gorm:"column:updated_by" json:"-"`
}

func (VmsTrnReceivedKeyDriver) TableName() string {
	return "public.vms_trn_request"
}

//VmsTrnReceivedKeyPEA
type VmsTrnReceivedKeyPEA struct {
	TrnRequestUID                 string    `gorm:"column:trn_request_uid;primaryKey;" json:"trn_request_uid" example:"3045a994-ba0b-431d-acf2-98768a9c5fc9"`
	ReceiverKeyType               int       `gorm:"column:receiver_key_type" json:"-" example:"2"`
	ReceivedKeyEmpID              string    `gorm:"column:received_key_emp_id" json:"received_key_emp_id" example:"990001"`
	ReceivedKeyEmpName            string    `gorm:"column:received_key_emp_name" json:"-"`
	ReceivedKeyDeptSAP            string    `gorm:"column:received_key_dept_sap" json:"-"`
	ReceivedKeyDeptSAPShort       string    `gorm:"column:received_key_dept_sap_short" json:"-"`
	ReceivedKeyDeptSAPFull        string    `gorm:"column:received_key_dept_sap_full" json:"-"`
	ReceivedKeyInternalContactNum string    `gorm:"column:received_key_internal_contact_number" json:"received_key_internal_contact_number" example:"5551234"`
	ReceivedKeyMobileContactNum   string    `gorm:"column:received_key_mobile_contact_number" json:"received_key_mobile_contact_number" example:"0812345678"`
	ReceivedKeyRemark             string    `gorm:"column:received_key_remark" json:"received_key_remark" example:"Employee received the key"`
	RefRequestStatusCode          string    `gorm:"column:ref_request_status_code" json:"-"`
	UpdatedAt                     time.Time `gorm:"column:updated_at" json:"-"`
	UpdatedBy                     string    `gorm:"column:updated_by" json:"-"`
}

func (VmsTrnReceivedKeyPEA) TableName() string {
	return "public.vms_trn_request"
}

//VmsTrnReceivedKeyOutSider
type VmsTrnReceivedKeyOutSider struct {
	TrnRequestUID               string    `gorm:"column:trn_request_uid;primaryKey;" json:"trn_request_uid" example:"3045a994-ba0b-431d-acf2-98768a9c5fc9"`
	ReceiverKeyType             int       `gorm:"column:receiver_key_type" json:"-" example:"3"`
	ReceivedKeyEmpName          string    `gorm:"column:received_key_emp_name" json:"outsider_name" example:"John Doe"`
	ReceivedKeyMobileContactNum string    `gorm:"column:received_key_mobile_contact_number" json:"received_key_mobile_contact_number" example:"0812345678"`
	ReceivedKeyRemark           string    `gorm:"column:received_key_remark" json:"received_key_remark" example:"OutSider received the key"`
	RefRequestStatusCode        string    `gorm:"column:ref_request_status_code" json:"-"`
	UpdatedAt                   time.Time `gorm:"column:updated_at" json:"-"`
	UpdatedBy                   string    `gorm:"column:updated_by" json:"-"`
}

func (VmsTrnReceivedKeyOutSider) TableName() string {
	return "public.vms_trn_request"
}

// VmsTrnRequestUpdateRecieiveKey
type VmsTrnRequestUpdateRecieivedKey struct {
	TrnRequestUID            string    `gorm:"column:trn_request_uid;primaryKey" json:"trn_request_uid" example:"8bd09808-61fa-42fd-8a03-bf961b5678cd"`
	ReceivedKeyPlace         string    `gorm:"column:received_key_place" json:"received_key_place" example:"Main Office"`
	ReceivedKeyStartDatetime time.Time `gorm:"column:received_key_start_datetime" json:"received_key_start_datetime" example:"2025-02-16T08:00:00Z"`
	ReceivedKeyEndDatetime   time.Time `gorm:"column:received_key_end_datetime" json:"received_key_end_datetime" example:"2025-02-16T17:00:00Z"`
	UpdatedAt                time.Time `gorm:"column:updated_at" json:"-"`
	UpdatedBy                string    `gorm:"column:updated_by" json:"-"`
}

func (VmsTrnRequestUpdateRecieivedKey) TableName() string {
	return "public.vms_trn_request"
}

// VmsTrnRequestUpdateRecieivedKeyDetail
type VmsTrnRequestUpdateRecieivedKeyDetail struct {
	TrnRequestUID         string    `gorm:"column:trn_request_uid;primaryKey" json:"trn_request_uid" example:"8bd09808-61fa-42fd-8a03-bf961b5678cd"`
	RefVehicleKeyTypeCode int       `gorm:"column:ref_vehicle_key_type_code" json:"ref_vehicle_key_type_code" example:"1"`
	ReceivedKeyDatetime   time.Time `gorm:"column:received_key_datetime" json:"received_key_datetime" example:"2025-02-16T08:00:00Z"`
	UpdatedAt             time.Time `gorm:"column:updated_at" json:"-"`
	UpdatedBy             string    `gorm:"column:updated_by" json:"-"`
}

func (VmsTrnRequestUpdateRecieivedKeyDetail) TableName() string {
	return "public.vms_trn_request"
}
