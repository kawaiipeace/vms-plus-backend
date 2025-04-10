package models

type VmsTrnReceivedKey_Emp struct {
	TrnRequestUID                 string `gorm:"column:trn_request_uid;type:uuid;" json:"trn_request_uid" example:"a7de5318-1e05-4511-abe7-8c1c6374ab29"`
	IsPEAEmployeeReceivedKey      bool   `gorm:"column:is_pea_employee_received_key" json:"-" example:"true"`
	ReceivedKeyEmpID              string `gorm:"column:received_key_emp_id" json:"received_key_emp_id" example:"1234567890"`
	ReceivedKeyInternalContactNum string `gorm:"column:received_key_internal_contact_number" json:"received_key_internal_contact_number" example:"5551234"`
	ReceivedKeyMobileContactNum   string `gorm:"column:received_key_mobile_contact_number" json:"received_key_mobile_contact_number" example:"0812345678"`
	ReceivedKeyRemark             string `gorm:"column:received_key_remark" json:"received_key_remark" example:"Employee received the key"`
}

func (VmsTrnReceivedKey_Emp) TableName() string {
	return "public.vms_trn_request"
}

type VmsTrnReceivedKey_OutSource struct {
	TrnRequestUID               string `gorm:"column:trn_request_uid;type:uuid;" json:"trn_request_uid" example:"a7de5318-1e05-4511-abe7-8c1c6374ab29"`
	IsPEAEmployeeReceivedKey    bool   `gorm:"column:is_pea_employee_received_key" json:"-" example:"true"`
	ReceivedKeyEmpID            string `gorm:"column:received_key_emp_id" json:"received_key_emp_id" example:"1234567890"`
	ReceivedKeyEmpName          string `gorm:"column:received_key_emp_name" json:"received_key_emp_name" example:"John Doe"`
	ReceivedKeyMobileContactNum string `gorm:"column:received_key_mobile_contact_number" json:"received_key_mobile_contact_number" example:"0812345678"`
	ReceivedKeyRemark           string `gorm:"column:received_key_remark" json:"received_key_remark" example:"Employee received the key"`
}

func (VmsTrnReceivedKey_OutSource) TableName() string {
	return "public.vms_trn_request"
}
