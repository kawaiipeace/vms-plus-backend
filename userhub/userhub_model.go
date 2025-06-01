package userhub

type ServiceCheckPhoneNumberRequest struct {
	ServiceCode string `json:"service_code" example:"vms"`
	Phone       string `json:"phone" example:"0818088770"`
}
type ServiceCheckPhoneNumberResponse struct {
	IsValid bool   `json:"is_valid" example:"true"`
	Message string `json:"message" example:"success"`
}

type ServiceLoginUserRequest struct {
	ServiceCode string `json:"service_code" example:"vms"`
	LoginBy     string `json:"login_by" example:"keycloak"`
	EmpID       string `json:"emp_id" example:"700001"`
	IdentityNo  string `json:"identity_no" example:"1234567890"`
	Phone       string `json:"phone" example:"0818088770"`
	IpAddress   string `json:"ip_address" example:"192.168.1.1"`
}
type ServiceUserInfoRequest struct {
	ServiceCode string `json:"service_code" example:"vms"`
	EmpID       string `json:"emp_id" example:"700001"`
}

type ServiceUserInfoResponse struct {
	EmpID         string   `json:"emp_id" example:"700001"`
	FirstName     string   `json:"first_name" example:"John"`
	LastName      string   `json:"last_name" example:"Doe"`
	FullName      string   `json:"full_name" example:"John Doe"`
	IdentityNo    string   `json:"identity_no" example:"1234567890"`
	Position      string   `json:"posi_text" example:"Manager"`
	DeptSAP       string   `json:"dept_sap" example:"1234567890"`
	DeptSAPShort  string   `json:"dept_sap_short" example:"1234567890"`
	DeptSAPFull   string   `json:"dept_sap_full" example:"1234567890"`
	BureauDeptSap string   `json:"bureau_dept_sap" example:"1234567890"`
	MobilePhone   string   `json:"mobile_number" example:"0818088770"`
	DeskPhone     string   `json:"internal_number" example:"0818088770"`
	BusinessArea  string   `json:"business_area" example:"1234567890"`
	LevelCode     string   `json:"level_code" example:"1234567890"`
	ImageUrl      string   `json:"image_url" example:"https://example.com/image.jpg"`
	Roles         []string `json:"roles" example:"['admin', 'user']"`
}

type ServiceListUserRequest struct {
	ServiceCode   string   `json:"service_code" example:"vms"`
	Search        string   `json:"search" example:"700001"`
	UpperDeptSap  string   `json:"upper_dept_sap" example:"4455"`
	BureauDeptSap string   `json:"bureau_dept_sap" example:"4455"`
	BusinessArea  string   `json:"business_area" example:"Z00"`
	LevelCodes    string   `json:"level_codes" example:"M1,M2,M3"`
	EmpIDs        []string `json:"emp_ids" example:"['700001', '700002']"`
	Role          string   `json:"role" example:"admin_approval"`
	Limit         int      `json:"limit" example:"10"`
}
