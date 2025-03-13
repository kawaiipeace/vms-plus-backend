package models

type MasUserEmp struct {
	EmpID        string `gorm:"column:emp_id" json:"emp_id"`
	FullName     string `gorm:"column:full_name" json:"full_name"`
	DeptSAP      string `gorm:"column:dept_sap" json:"dept_sap"`
	DeptSAPShort string `gorm:"column:dept_sap_short" json:"dept_sap_short"`
	DeptSAPFull  string `gorm:"column:dept_sap_full" json:"dept_sap_full"`
	TelMobile    string `gorm:"column:tel_mobile" json:"tel_mobile"`
	TelInternal  string `gorm:"column:tel_internal" json:"tel_internal"`
	Image_url    string `gorm:"column:image_url" json:"image_url"`
}

func (MasUserEmp) TableName() string {
	return "vms_user.mas_employee"
}
