package models

import "time"

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

type VmsTrnReceivedKeyPEA struct {
	TrnRequestUID                 string    `gorm:"column:trn_request_uid;primaryKey;" json:"trn_request_uid" example:"3 w"`
	ReceiverKeyType               int       `gorm:"column:receiver_key_type" json:"-" example:"2"`
	ReceivedKeyEmpID              string    `gorm:"column:received_key_emp_id" json:"received_key_emp_id" example:"1234567890"`
	ReceivedKeyInternalContactNum string    `gorm:"column:received_key_internal_contact_number" json:"received_key_internal_contact_number" example:"5551234"`
	ReceivedKeyMobileContactNum   string    `gorm:"column:received_key_mobile_contact_number" json:"received_key_mobile_contact_number" example:"0812345678"`
	ReceivedKeyRemark             string    `gorm:"column:received_key_remark" json:"received_key_remark" example:"Employee received the key"`
	RefRequestStatusCode          string    `gorm:"column:ref_request_status_code" json:"ref_request_status_code"`
	UpdatedAt                     time.Time `gorm:"column:updated_at" json:"-"`
	UpdatedBy                     string    `gorm:"column:updated_by" json:"-"`
}

func (VmsTrnReceivedKeyPEA) TableName() string {
	return "public.vms_trn_request"
}

type VmsTrnReceivedKeyOutSider struct {
	TrnRequestUID               string    `gorm:"column:trn_request_uid;primaryKey;" json:"trn_request_uid" example:"3045a994-ba0b-431d-acf2-98768a9c5fc9"`
	ReceiverKeyType             int       `gorm:"column:receiver_key_type" json:"-" example:"3"`
	ReceivedKeyEmpName          string    `gorm:"column:received_key_emp_name" json:"outsider_name" example:"John Doe"`
	ReceivedKeyMobileContactNum string    `gorm:"column:received_key_mobile_contact_number" json:"received_key_mobile_contact_number" example:"0812345678"`
	ReceivedKeyRemark           string    `gorm:"column:received_key_remark" json:"received_key_remark" example:"OutSider received the key"`
	RefRequestStatusCode        string    `gorm:"column:ref_request_status_code" json:"ref_request_status_code"`
	UpdatedAt                   time.Time `gorm:"column:updated_at" json:"-"`
	UpdatedBy                   string    `gorm:"column:updated_by" json:"-"`
}

func (VmsTrnReceivedKeyOutSider) TableName() string {
	return "public.vms_trn_request"
}
