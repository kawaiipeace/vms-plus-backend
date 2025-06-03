package models

type AuthenJwtUsr struct {
	EmpID    string   `gorm:"column:emp_uid" json:"emp_id"`
	FullName string   `gorm:"column:full_name" json:"full_name"`
	DeptSAP  string   `gorm:"column:dept_sap" json:"dept_sap"`
	Roles    []string `gorm:"column:roles" json:"roles"`
	LoginBy  string   `gorm:"column:login_by" json:"login_by"`
}

// AuthenUserEmp
type AuthenUserEmp struct {
	EmpID                     string   `gorm:"column:emp_id" json:"emp_id"`
	FirstName                 string   `gorm:"column:first_name" json:"first_name"`
	LastName                  string   `gorm:"column:last_name" json:"last_name"`
	FullName                  string   `gorm:"column:full_name" json:"full_name"`
	Position                  string   `gorm:"column:posi_text" json:"posi_text"`
	DeptSAP                   string   `gorm:"column:dept_sap" json:"dept_sap"`
	DeptSAPShort              string   `gorm:"column:dept_sap_short" json:"dept_sap_short"`
	DeptSAPFull               string   `gorm:"column:dept_sap_full" json:"dept_sap_full"`
	BureauDeptSap             string   `gorm:"column:bureau_dept_sap" json:"bureau_dept_sap"`
	MobilePhone               string   `gorm:"column:mobile_number" json:"mobile_number"`
	DeskPhone                 string   `gorm:"column:internal_number" json:"internal_number"`
	BusinessArea              string   `gorm:"column:business_area" json:"business_area"`
	ImageUrl                  string   `gorm:"-" json:"image_url"`
	LicenseStatusCode         string   `gorm:"-" json:"license_status_code"`
	LicenseStatus             string   `gorm:"-" json:"license_status"`
	TrnRequestAnnualDriverUID string   `gorm:"-" json:"trn_request_annual_driver_uid"`
	AnnualYYYY                int      `gorm:"-" json:"annual_yyyy" example:"2568"`
	Roles                     []string `gorm:"-" json:"roles"`
	LoginBy                   string   `gorm:"-" json:"login_by"`
	IsEmployee                bool     `gorm:"-" json:"is_employee"`
	LevelCode                 string   `gorm:"-" json:"level_code"`
}

func (AuthenUserEmp) TableName() string {
	return "mas_employee"
}
