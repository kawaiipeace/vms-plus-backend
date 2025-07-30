package models

//MasUserEmp
type MasUserEmp struct {
	EmpID         string `gorm:"column:emp_id" json:"emp_id"`
	FullName      string `gorm:"column:full_name" json:"full_name"`
	DeptSAP       string `gorm:"column:dept_sap" json:"dept_sap"`
	DeptSAPShort  string `gorm:"column:dept_sap_short" json:"dept_sap_short"`
	DeptSAPFull   string `gorm:"column:dept_sap_full" json:"dept_sap_full"`
	Position      string `gorm:"column:posi_text" json:"posi_text"`
	TelMobile     string `gorm:"column:tel_mobile" json:"tel_mobile"`
	TelInternal   string `gorm:"column:tel_internal" json:"tel_internal"`
	BureauDeptSap string `gorm:"column:bureau_dept_sap" json:"bureau_dept_sap"`
	BusinessArea  string `gorm:"column:business_area" json:"business_area"`
	ImageUrl      string `gorm:"column:image_url" json:"image_url"`
	LevelCode     string `gorm:"column:level_code" json:"level_code"`
	IsEmployee    bool   `gorm:"column:is_employee" json:"is_employee"`
}

func (MasUserEmp) TableName() string {
	return "mas_employee"
}

// MasUserDriver
type MasUserDriver struct {
	EmpID        string             `gorm:"column:emp_id" json:"emp_id"`
	FullName     string             `gorm:"column:full_name" json:"full_name"`
	Position     string             `gorm:"column:posi_text" json:"posi_text"`
	DeptSAP      string             `gorm:"column:dept_sap" json:"dept_sap"`
	DeptSAPShort string             `gorm:"column:dept_sap_short" json:"dept_sap_short"`
	DeptSAPFull  string             `gorm:"column:dept_sap_full" json:"dept_sap_full"`
	TelMobile    string             `gorm:"column:tel_mobile" json:"tel_mobile"`
	TelInternal  string             `gorm:"column:tel_internal" json:"tel_internal"`
	ImageUrl     string             `gorm:"column:image_url" json:"image_url"`
	AnnualDriver VmsTrnAnnualDriver `gorm:"foreignKey:EmpID;references:CreatedRequestEmpId" json:"annual_driver"`
}

func (MasUserDriver) TableName() string {
	return "mas_employee"
}

// VmsMasSatisfactionSurveyQuestions
type VmsMasSatisfactionSurveyQuestions struct {
	MasSatisfactionSurveyQuestionsUID string `gorm:"column:mas_satisfaction_survey_questions_uid;" json:"mas_satisfaction_survey_questions_uid"`
	QuestionTitle                     string `gorm:"column:question_title" json:"question_title"`
	QuestionsDescription              string `gorm:"column:questions_description" json:"questions_description"`
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

type VmsMasDepartment struct {
	DeptSAP        string `gorm:"column:dept_sap;primaryKey" json:"dept_sap"`
	DeptShort      string `gorm:"column:dept_short" json:"dept_short"`
	DeptFull       string `gorm:"column:dept_full" json:"dept_full"`
	CostCenterCode string `gorm:"column:cost_center_code" json:"cost_center_code"`
	CostCenterName string `gorm:"column:cost_center_name" json:"cost_center_name"`
	BusinessArea   string `gorm:"column:business_area" json:"business_area"`
}

func (VmsMasDepartment) TableName() string {
	return "public.vms_mas_department"
}

type VmsMasVehicleArray struct {
	MasVehicleUID string `gorm:"column:mas_vehicle_uid" json:"mas_vehicle_uid" example:"f3b29096-140e-49dc-97ee-17fa9352aff6"`
}

type VmsMasDriverArray struct {
	MasDriverUID string `gorm:"column:mas_driver_uid" json:"mas_driver_uid" example:"ec4a2cee-aded-47bd-9d93-4a1a74cb58a4"`
}

type VmsMasDepartmentTree struct {
	DeptSAP      string                 `gorm:"column:dept_sap;primaryKey" json:"-"`
	DeptUpper    string                 `gorm:"column:dept_upper" json:"-"`
	DeptShort    string                 `gorm:"column:dept_short" json:"-"`
	DeptFull     string                 `gorm:"column:dept_full" json:"text"`
	ResourceName string                 `gorm:"column:resource_name" json:"resource_name"`
	DeptUnder    []VmsMasDepartmentTree `gorm:"foreignKey:DeptSAP;references:DeptUpper" json:"children"`
}

func (VmsMasDepartmentTree) TableName() string {
	return "mas_department"
}

type VmsMasHolidays struct {
	HolidaysDate   TimeWithZone `gorm:"column:mas_holidays_date" json:"mas_holidays_date"`
	HolidaysDetail string       `gorm:"column:mas_holidays_detail" json:"mas_holidays_detail"`
}

func (VmsMasHolidays) TableName() string {
	return "vms_mas_holidays"
}

type VmsMasManager struct {
	Type           string `gorm:"column:type" json:"type"`
	EmpIDLeader    int    `gorm:"column:emp_id_leader" json:"emp_id_leader"`
	EmpName        string `gorm:"column:emp_name" json:"emp_name"`
	PosiCode       int    `gorm:"column:posi_code" json:"posi_code"`
	LevelCode      string `gorm:"column:level_code" json:"level_code"`
	DeptShort      string `gorm:"column:dept_short" json:"dept_short"`
	DeptSAP        int    `gorm:"column:dept_sap" json:"dept_sap"`
	DeptUpper      int    `gorm:"column:dept_upper" json:"dept_upper"`
	PlansTextShort string `gorm:"column:plans_text_short" json:"plans_text_short"`
	Email          string `gorm:"column:email" json:"email"`
}
