package models

//MasUserEmp
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

// MasUserDriver
type MasUserDriver struct {
	EmpID        string             `gorm:"column:emp_id" json:"emp_id"`
	FullName     string             `gorm:"column:full_name" json:"full_name"`
	DeptSAP      string             `gorm:"column:dept_sap" json:"dept_sap"`
	DeptSAPShort string             `gorm:"column:dept_sap_short" json:"dept_sap_short"`
	DeptSAPFull  string             `gorm:"column:dept_sap_full" json:"dept_sap_full"`
	TelMobile    string             `gorm:"column:tel_mobile" json:"tel_mobile"`
	TelInternal  string             `gorm:"column:tel_internal" json:"tel_internal"`
	ImageURL     string             `gorm:"column:image_url" json:"image_url"`
	AnnualDriver VmsTrnAnnualDriver `gorm:"foreignKey:EmpID;references:CreatedRequestEmpId" json:"annual_driver"`
}

func (MasUserDriver) TableName() string {
	return "vms_user.mas_employee"
}

// VmsMasSatisfactionSurveyQuestions
type VmsMasSatisfactionSurveyQuestions struct {
	MasSatisfactionSurveyQuestionsCode  string `gorm:"column:mas_satisfaction_survey_questions_code;" json:"mas_satisfaction_survey_questions_code"`
	MasSatisfactionSurveyQuestionsTitle string `gorm:"column:mas_satisfaction_survey_questions_title" json:"mas_satisfaction_survey_questions_title"`
	MasSatisfactionSurveyQuestionsDesc  string `gorm:"column:mas_satisfaction_survey_questions_desc" json:"mas_satisfaction_survey_questions_desc"`
}

func (VmsMasSatisfactionSurveyQuestions) TableName() string {
	return "public.vms_mas_satisfaction_survey_questions"
}

//VmsMasVehicleDepartmentList
type VmsMasVehicleDepartmentList struct {
	VehicleOwnerDeptSap string `gorm:"column:vehicle_owner_dept_sap" json:"vehicle_owner_dept_sap"`
	DeptSapShort        string `gorm:"column:dept_sap_short" json:"dept_sap_short"`
	DeptSapFull         string `gorm:"column:dept_sap_full" json:"dept_sap_full"`
	DeptType            string `gorm:"column:dept_type" json:"dept_type"`
}
