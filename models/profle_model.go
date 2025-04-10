package models

type AuthenJwtUsr struct {
	EmpID        string   `gorm:"column:emp_uid" json:"emp_id"`
	FullName     string   `gorm:"column:full_name" json:"full_name"`
	DeptSAP      string   `gorm:"column:dept_sap" json:"dept_sap"`
	Roles        []string `gorm:"column:roles" json:"roles"`
	AccessToken  string   `gorm:"column:access_token" json:"access_token"`
	RefreshToken string   `gorm:"column:refresh_token" json:"refresh_token"`
}

type AuthenUserEmp struct {
	EmpID        string `gorm:"column:emp_id" json:"emp_id"`
	FirstName    string `gorm:"column:first_name" json:"first_name"`
	LastName     string `gorm:"column:last_name" json:"last_name"`
	DeptSAP      string `gorm:"column:dept_sap" json:"dept_sap"`
	DeptSAPShort string `gorm:"column:dept_sap_short" json:"dept_sap_short"`
	DeptSAPFull  string `gorm:"column:dept_sap_full" json:"dept_sap_full"`
}

func (e *AuthenUserEmp) FullName() string {
	return e.FirstName + " " + e.LastName
}

func (AuthenUserEmp) TableName() string {
	return "vms_user.mas_employee"
}
