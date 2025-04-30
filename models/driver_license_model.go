package models

import "time"

//DriverLicense
type VmsDriverLicenseCard struct {
	MasDriverUID               string                      `json:"mas_driver_uid" gorm:"column:mas_driver_uid"`
	DriverName                 string                      `gorm:"column:driver_name" json:"driver_name" example:"John Doe"`
	DriverDeptSapShortNameWork string                      `gorm:"column:driver_dept_sap_short_name_work" json:"driver_dept_sap_short_name_work" example:"กยจ."`
	DriverLicense              VmsDriverLicenseCardLicense `gorm:"foreignKey:MasDriverUID;references:MasDriverUID" json:"driver_license"`
	LicenseStatus              string                      `gorm:"-" json:"license_status"`
	AnnualYYYY                 int                         `gorm:"column:-" json:"annual_yyyy" example:"2568"`
	IsNoExpiryDate             bool                        `gorm:"column:is_no_expiry_date" json:"is_no_expiry_date"`
}

func (VmsDriverLicenseCard) TableName() string {
	return "vms_mas_driver"
}

//VmsDriverLicenseCardLicense
type VmsDriverLicenseCardLicense struct {
	MasDriverUID             string                  `json:"mas_driver_uid" gorm:"column:mas_driver_uid"`
	DriverLicenseNo          string                  `json:"driver_license_no" gorm:"column:driver_license_no"`
	RefDriverLicenseTypeCode string                  `json:"ref_driver_license_type_code" gorm:"column:ref_driver_license_type_code"`
	DriverLicenseStartDate   string                  `json:"driver_license_start_date" gorm:"column:driver_license_start_date"`
	DriverLicenseEndDate     string                  `json:"driver_license_end_date" gorm:"column:driver_license_end_date"`
	DriverLicenseType        VmsRefDriverLicenseType `gorm:"foreignKey:RefDriverLicenseTypeCode;references:RefDriverLicenseTypeCode" json:"driver_license_type"`
}

func (VmsDriverLicenseCardLicense) TableName() string {
	return "vms_mas_driver_license"
}

//VmsTrnRequestAnnualDriverSummary
type VmsTrnRequestAnnualDriverSummary struct {
	RefRequestAnnualDriverStatusCode string `gorm:"column:ref_request_annual_driver_status_code" json:"ref_request_annual_driver_status_code"`
	RefRequestAnnualDriverStatusName string `json:"ref_request_annual_driver_status_name"`
	Count                            int    `gorm:"column:count" json:"count"`
}

//VmsDriverLicenseAnnualList
type VmsDriverLicenseAnnualList struct {
	TrnRequestAnnualDriverUID        string    `gorm:"column:trn_request_annual_driver_uid;primaryKey" json:"trn_request_annual_driver_uid"`
	RequestAnnualDriverNo            string    `gorm:"column:request_annual_driver_no" json:"request_annual_driver_no"`
	AnnualYYYY                       int       `gorm:"column:annual_yyyy" json:"annual_yyyy" example:"2568"`
	CreatedRequestDatetime           time.Time `gorm:"column:created_request_datetime" json:"created_request_datetime"`
	CreatedRequestEmpID              string    `gorm:"column:created_request_emp_id" json:"created_request_emp_id"`
	CreatedRequestEmpName            string    `gorm:"column:created_request_emp_name" json:"created_request_emp_name"`
	CreatedRequestDeptSapNameShort   string    `gorm:"column:created_request_dept_sap_name_short" json:"created_request_dept_sap_name_short"`
	CreatedRequestDeptSapNameFull    string    `gorm:"column:created_request_dept_sap_name_full" json:"created_request_dept_sap_name_full"`
	RefRequestAnnualDriverStatusCode string    `gorm:"column:ref_request_annual_driver_status_code" json:"ref_request_annual_driver_status_code"`
	RefRequestAnnualDriverStatusName string    `gorm:"column:-" json:"ref_request_annual_driver_status_name"`
	RefDriverLicenseTypeCode         string    `gorm:"column:ref_driver_license_type_code" json:"ref_driver_license_type_code" example:"1"`
	RefDriverLicenseTypeName         string    `gorm:"column:ref_driver_license_type_name" json:"ref_driver_license_type_name" example:"1"`
	DriverLicenseExpireDate          time.Time `gorm:"column:driver_license_expire_date" json:"driver_license_expire_date" example:"2025-12-31T00:00:00Z"`
}

func (VmsDriverLicenseAnnualList) TableName() string {
	return "vms_trn_request_annual_driver"
}

//VmsDriverLicenseAnnualRequest
type VmsDriverLicenseAnnualRequest struct {
	TrnRequestAnnualDriverUID        string    `gorm:"column:trn_request_annual_driver_uid;primaryKey" json:"-"`
	RequestAnnualDriverNo            string    `gorm:"column:request_annual_driver_no" json:"-"`
	AnnualYYYY                       int       `gorm:"column:annual_yyyy" json:"annual_yyyy" example:"2568"`
	CreatedRequestDatetime           time.Time `gorm:"column:created_request_datetime" json:"-"`
	CreatedRequestEmpID              string    `gorm:"column:created_request_emp_id" json:"-"`
	CreatedRequestEmpName            string    `gorm:"column:created_request_emp_name" json:"-"`
	CreatedRequestEmpPosition        string    `gorm:"column:created_request_emp_position" json:"-"`
	CreatedRequestDeptSap            string    `gorm:"column:created_request_dept_sap" json:"-"`
	CreatedRequestDeptSapNameShort   string    `gorm:"column:created_request_dept_sap_name_short" json:"-"`
	CreatedRequestDeptSapNameFull    string    `gorm:"column:created_request_dept_sap_name_full" json:"-"`
	CreatedRequestPhoneNumber        string    `gorm:"column:created_request_phone_number" json:"-"`
	CreatedRequestMobileNumber       string    `gorm:"column:created_request_mobile_number" json:"-"`
	RefRequestAnnualDriverStatusCode string    `gorm:"column:ref_request_annual_driver_status_code" json:"-"`
	DriverLicenseNo                  string    `gorm:"column:driver_license_no" json:"driver_license_no" example:"DL12345678"`
	RefDriverLicenseTypeCode         string    `gorm:"column:ref_driver_license_type_code" json:"ref_driver_license_type_code" example:"1"`
	DriverLicenseExpireDate          time.Time `gorm:"column:driver_license_expire_date" json:"driver_license_expire_date" example:"2025-12-31T00:00:00Z"`
	DriverLicenseImg                 string    `gorm:"column:driver_license_img" json:"driver_license_img" example:"image_url"`
	DriverCertificateNo              string    `gorm:"column:driver_certificate_no" json:"driver_certificate_no" example:"CERT12345"`
	DriverCertificateName            string    `gorm:"column:driver_certificate_name" json:"driver_certificate_name" example:"Safety Certificate"`
	DriverCertificateTypeCode        int       `gorm:"column:driver_certificate_type_code" json:"driver_certificate_type_code" example:"1"`
	DriverCertificateIssueDate       time.Time `gorm:"column:driver_certificate_issue_date" json:"driver_certificate_issue_date" example:"2023-01-01T00:00:00Z"`
	DriverCertificateExpireDate      time.Time `gorm:"column:driver_certificate_expire_date" json:"driver_certificate_expire_date" example:"2024-12-31T00:00:00Z"`
	DriverCertificateImg             string    `gorm:"column:driver_certificate_img" json:"driver_certificate_img" example:"certificate_image_url"`
	RequestIssueDate                 time.Time `gorm:"column:request_issue_date" json:"request_issue_date" example:"2023-01-01T00:00:00Z"`
	RequestExpireDate                time.Time `gorm:"column:request_expire_date" json:"request_expire_date" example:"2023-12-31T00:00:00Z"`
	UpdatedAt                        time.Time `gorm:"column:updated_at" json:"-"`
	UpdatedBy                        string    `gorm:"column:updated_by" json:"-"`

	ConfirmedRequestEmpID        string `gorm:"column:confirmed_request_emp_id" json:"confirmed_request_emp_id" example:"990002"`
	ConfirmedRequestEmpName      string `gorm:"column:confirmed_request_emp_name" json:"-"`
	ConfirmedRequestEmpPosition  string `gorm:"column:confirmed_request_emp_position" json:"-"`
	ConfirmedRequestDeptSap      string `gorm:"column:confirmed_request_dept_sap" json:"-"`
	ConfirmedRequestDeptSapShort string `gorm:"column:confirmed_request_dept_sap_short" json:"-"`
	ConfirmedRequestDeptSapFull  string `gorm:"column:confirmed_request_dept_sap_full" json:"-"`
	ConfirmedRequestPhoneNumber  string `gorm:"column:confirmed_request_phone_number" json:"-"`
	ConfirmedRequestMobileNumber string `gorm:"column:confirmed_request_mobile_number" json:"-"`

	ApprovedRequestEmpID        string `gorm:"column:approved_request_emp_id" json:"approved_request_emp_id" example:"990001"`
	ApprovedRequestEmpName      string `gorm:"column:approved_request_emp_name" json:"-"`
	ApprovedRequestEmpPosition  string `gorm:"column:approved_request_emp_position" json:"-"`
	ApprovedRequestDeptSap      string `gorm:"column:approved_request_dept_sap" json:"-"`
	ApprovedRequestDeptSapShort string `gorm:"column:approved_request_dept_sap_short" json:"-"`
	ApprovedRequestDeptSapFull  string `gorm:"column:approved_request_dept_sap_full" json:"-"`
	ApprovedRequestPhoneNumber  string `gorm:"column:approved_request_phone_number" json:"-"`
	ApprovedRequestMobileNumber string `gorm:"column:approved_request_mobile_number" json:"-"`

	RejectedRequestEmpPosition string `gorm:"column:rejected_request_emp_position" json:"-"`
	CanceledRequestEmpPosition string `gorm:"column:canceled_request_emp_position" json:"-"`
}

func (VmsDriverLicenseAnnualRequest) TableName() string {
	return "vms_trn_request_annual_driver"
}

type VmsTrnRequestAnnualDriverNo struct {
	RequestAnnualDriverNo string `gorm:"column:request_annual_driver_no" json:"request_annual_driver_no"`
}

//VmsDriverLicenseAnnualResponse
type VmsDriverLicenseAnnualResponse struct {
	TrnRequestAnnualDriverUID        string    `gorm:"column:trn_request_annual_driver_uid;primaryKey" json:"trn_request_annual_driver_uid"`
	RequestAnnualDriverNo            string    `gorm:"column:request_annual_driver_no" json:"request_annual_driver_no"`
	AnnualYYYY                       int       `gorm:"column:annual_yyyy" json:"annual_yyyy" example:"2568"`
	CreatedRequestDatetime           time.Time `gorm:"column:created_request_datetime" json:"created_request_datetime"`
	CreatedRequestEmpID              string    `gorm:"column:created_request_emp_id" json:"created_request_emp_id"`
	CreatedRequestEmpName            string    `gorm:"column:created_request_emp_name" json:"created_request_emp_name"`
	CreatedRequestEmpPosition        string    `gorm:"column:created_request_emp_position" json:"created_request_emp_position"`
	CreatedRequestDeptSap            string    `gorm:"column:created_request_dept_sap" json:"created_request_dept_sap"`
	CreatedRequestDeptSapNameShort   string    `gorm:"column:created_request_dept_sap_name_short" json:"created_request_dept_sap_name_short"`
	CreatedRequestDeptSapNameFull    string    `gorm:"column:created_request_dept_sap_name_full" json:"created_request_dept_sap_name_full"`
	CreatedRequestPhoneNumber        string    `gorm:"column:created_request_phone_number" json:"created_request_phone_number"`
	CreatedRequestMobileNumber       string    `gorm:"column:created_request_mobile_number" json:"created_request_mobile_number"`
	RefRequestAnnualDriverStatusCode string    `gorm:"column:ref_request_annual_driver_status_code" json:"ref_request_annual_driver_status_code"`
	DriverLicenseNo                  string    `gorm:"column:driver_license_no" json:"driver_license_no" example:"DL12345678"`
	RefDriverLicenseTypeCode         string    `gorm:"column:ref_driver_license_type_code" json:"ref_driver_license_type_code" example:"1"`
	DriverLicenseExpireDate          time.Time `gorm:"column:driver_license_expire_date" json:"driver_license_expire_date" example:"2025-12-31T00:00:00Z"`
	DriverLicenseImg                 string    `gorm:"column:driver_license_img" json:"driver_license_img" example:"image_url"`
	DriverCertificateNo              string    `gorm:"column:driver_certificate_no" json:"driver_certificate_no" example:"CERT12345"`
	DriverCertificateName            string    `gorm:"column:driver_certificate_name" json:"driver_certificate_name" example:"Safety Certificate"`
	DriverCertificateTypeCode        int       `gorm:"column:driver_certificate_type_code" json:"driver_certificate_type_code" example:"1"`
	DriverCertificateIssueDate       time.Time `gorm:"column:driver_certificate_issue_date" json:"driver_certificate_issue_date" example:"2023-01-01T00:00:00Z"`
	DriverCertificateExpireDate      time.Time `gorm:"column:driver_certificate_expire_date" json:"driver_certificate_expire_date" example:"2024-12-31T00:00:00Z"`
	DriverCertificateImg             string    `gorm:"column:driver_certificate_img" json:"driver_certificate_img" example:"certificate_image_url"`
	RequestIssueDate                 time.Time `gorm:"column:request_issue_date" json:"request_issue_date" example:"2023-01-01T00:00:00Z"`
	RequestExpireDate                time.Time `gorm:"column:request_expire_date" json:"request_expire_date" example:"2023-12-31T00:00:00Z"`
	UpdatedAt                        time.Time `gorm:"column:updated_at" json:"-"`
	UpdatedBy                        string    `gorm:"column:updated_by" json:"-"`

	ConfirmedRequestEmpID        string    `gorm:"column:confirmed_request_emp_id" json:"confirmed_request_emp_id" example:"990002"`
	ConfirmedRequestEmpName      string    `gorm:"column:confirmed_request_emp_name" json:"confirmed_request_emp_name"`
	ConfirmedRequestEmpPosition  string    `gorm:"column:confirmed_request_emp_position" json:"confirmed_request_emp_position"`
	ConfirmedRequestDeptSap      string    `gorm:"column:confirmed_request_dept_sap" json:"confirmed_request_dept_sap"`
	ConfirmedRequestDeptSapShort string    `gorm:"column:confirmed_request_dept_sap_short" json:"confirmed_request_dept_sap_short"`
	ConfirmedRequestDeptSapFull  string    `gorm:"column:confirmed_request_dept_sap_full" json:"confirmed_request_dept_sap_full"`
	ConfirmedRequestPhoneNumber  string    `gorm:"column:confirmed_request_phone_number" json:"confirmed_request_phone_number"`
	ConfirmedRequestMobileNumber string    `gorm:"column:confirmed_request_mobile_number" json:"confirmed_request_mobile_number"`
	ConfirmedRequestDatetime     time.Time `gorm:"column:confirmed_request_datetime" json:"confirmed_request_datetime"`

	ApprovedRequestEmpID        string    `gorm:"column:approved_request_emp_id" json:"approved_request_emp_id" example:"990001"`
	ApprovedRequestEmpName      string    `gorm:"column:approved_request_emp_name" json:"approved_request_emp_name"`
	ApprovedRequestEmpPosition  string    `gorm:"column:approved_request_emp_position" json:"approved_request_emp_position"`
	ApprovedRequestDeptSap      string    `gorm:"column:approved_request_dept_sap" json:"approved_request_dept_sap"`
	ApprovedRequestDeptSapShort string    `gorm:"column:approved_request_dept_sap_short" json:"approved_request_dept_sap_short"`
	ApprovedRequestDeptSapFull  string    `gorm:"column:approved_request_dept_sap_full" json:"approved_request_dept_sap_full"`
	ApprovedRequestPhoneNumber  string    `gorm:"column:approved_request_phone_number" json:"approved_request_phone_number"`
	ApprovedRequestMobileNumber string    `gorm:"column:approved_request_mobile_number" json:"approved_request_mobile_number"`
	ApprovedRequestDatetime     time.Time `gorm:"column:approved_request_datetime" json:"approved_request_datetime"`

	RejectedRequestEmpID        string    `gorm:"column:rejected_request_emp_id" json:"rejected_request_emp_id"`
	RejectedRequestEmpName      string    `gorm:"column:rejected_request_emp_name" json:"rejected_request_emp_name"`
	RejectedRequestEmpPosition  string    `gorm:"column:rejected_request_emp_position" json:"rejected_request_emp_position"`
	RejectedRequestDeptSap      string    `gorm:"column:rejected_request_dept_sap" json:"rejected_request_dept_sap"`
	RejectedRequestDeptSapShort string    `gorm:"column:rejected_request_dept_sap_short" json:"rejected_request_dept_sap_short"`
	RejectedRequestDeptSapFull  string    `gorm:"column:rejected_request_dept_sap_full" json:"rejected_request_dept_sap_full"`
	RejectedRequestPhoneNumber  string    `gorm:"column:rejected_request_phone_number" json:"rejected_request_phone_number"`
	RejectedRequestMobileNumber string    `gorm:"column:rejected_request_mobile_number" json:"rejected_request_mobile_number"`
	RejectedRequestReason       string    `gorm:"column:rejected_request_reason" json:"rejected_request_reason"`
	RejectedRequestDatetime     time.Time `gorm:"column:rejected_request_datetime" json:"rejected_request_datetime"`

	CanceledRequestEmpID        string                  `gorm:"column:canceled_request_emp_id" json:"canceled_request_emp_id"`
	CanceledRequestEmpName      string                  `gorm:"column:canceled_request_emp_name" json:"canceled_request_emp_name"`
	CanceledRequestEmpPosition  string                  `gorm:"column:canceled_request_emp_position" json:"canceled_request_emp_position"`
	CanceledRequestDeptSap      string                  `gorm:"column:canceled_request_dept_sap" json:"canceled_request_dept_sap"`
	CanceledRequestDeptSapShort string                  `gorm:"column:canceled_request_dept_sap_short" json:"canceled_request_dept_sap_short"`
	CanceledRequestDeptSapFull  string                  `gorm:"column:canceled_request_dept_sap_full" json:"canceled_request_dept_sap_full"`
	CanceledRequestPhoneNumber  string                  `gorm:"column:canceled_request_phone_number" json:"canceled_request_phone_number"`
	CanceledRequestMobileNumber string                  `gorm:"column:canceled_request_mobile_number" json:"canceled_request_mobile_number"`
	CanceledRequestDatetime     time.Time               `gorm:"column:canceled_request_datetime" json:"canceled_request_datetime"`
	CanceledRequestReason       string                  `gorm:"column:canceled_request_reason;" json:"canceled_request_reason" example:"Test Cancel"`
	DriverLicenseType           VmsRefDriverLicenseType `gorm:"foreignKey:RefDriverLicenseTypeCode;references:RefDriverLicenseTypeCode" json:"driver_license_type"`
	ProgressRequestStatus       []ProgressRequestStatus `gorm:"-" json:"progress_request_status"`
}

func (VmsDriverLicenseAnnualResponse) TableName() string {
	return "vms_trn_request_annual_driver"
}

// VmsDriverLicenseAnnualCanceled
type VmsDriverLicenseAnnualCanceled struct {
	TrnRequestAnnualDriverUID        string    `gorm:"column:trn_request_annual_driver_uid;primaryKey" json:"trn_request_annual_driver_uid" example:"095fbfbf-378e-4507-b15f-e53ac60370e7"`
	CanceledRequestReason            string    `gorm:"column:canceled_request_reason;" json:"canceled_request_reason" example:"Test Cancel"`
	CanceledRequestEmpID             string    `gorm:"column:canceled_request_emp_id" json:"-"`
	CanceledRequestEmpName           string    `gorm:"column:canceled_request_emp_name" json:"-"`
	CanceledRequestDeptSAP           string    `gorm:"column:canceled_request_dept_sap" json:"-"`
	CanceledRequestDeptSAPShort      string    `gorm:"column:canceled_request_dept_sap_short" json:"-"`
	CanceledRequestDeptSAPFull       string    `gorm:"column:canceled_request_dept_sap_full" json:"-"`
	RefRequestAnnualDriverStatusCode string    `gorm:"column:ref_request_annual_driver_status_code" json:"-"`
	CanceledRequestDatetime          time.Time `gorm:"column:canceled_request_datetime" json:"-" `
	UpdatedAt                        time.Time `gorm:"column:updated_at" json:"-"`
	UpdatedBy                        string    `gorm:"column:updated_by" json:"-"`
}

func (VmsDriverLicenseAnnualCanceled) TableName() string {
	return "public.vms_trn_request_annual_driver"
}

// VmsDriverLicenseAnnualRejected
type VmsDriverLicenseAnnualRejected struct {
	TrnRequestAnnualDriverUID        string    `gorm:"column:trn_request_annual_driver_uid;primaryKey" json:"trn_request_annual_driver_uid" example:"095fbfbf-378e-4507-b15f-e53ac60370e7"`
	RejectedRequestReason            string    `gorm:"column:rejected_request_reason;" json:"rejected_request_reason" example:"Test Reject"`
	RejectedRequestEmpID             string    `gorm:"column:rejected_request_emp_id" json:"-"`
	RejectedRequestEmpName           string    `gorm:"column:rejected_request_emp_name" json:"-"`
	RejectedRequestDeptSAP           string    `gorm:"column:rejected_request_dept_sap" json:"-"`
	RejectedRequestDeptSAPShort      string    `gorm:"column:rejected_request_dept_sap_short" json:"-"`
	RejectedRequestDeptSAPFull       string    `gorm:"column:rejected_request_dept_sap_full" json:"-"`
	RefRequestAnnualDriverStatusCode string    `gorm:"column:ref_request_annual_driver_status_code" json:"-"`
	RejectedRequestDatetime          time.Time `gorm:"column:rejected_request_datetime" json:"-" `
	UpdatedAt                        time.Time `gorm:"column:updated_at" json:"-"`
	UpdatedBy                        string    `gorm:"column:updated_by" json:"-"`
}

func (VmsDriverLicenseAnnualRejected) TableName() string {
	return "public.vms_trn_request_annual_driver"
}

// VmsDriverLicenseAnnualApproved
type VmsDriverLicenseAnnualApproved struct {
	TrnRequestAnnualDriverUID        string    `gorm:"column:trn_request_annual_driver_uid;primaryKey" json:"trn_request_annual_driver_uid" example:"095fbfbf-378e-4507-b15f-e53ac60370e7"`
	ApprovedRequestEmpID             string    `gorm:"column:approved_request_emp_id" json:"-" example:"990001"`
	ApprovedRequestEmpName           string    `gorm:"column:approved_request_emp_name" json:"-"`
	ApprovedRequestDeptSAP           string    `gorm:"column:approved_request_dept_sap" json:"-"`
	ApprovedRequestDeptSAPShort      string    `gorm:"column:approved_request_dept_sap_short" json:"-"`
	ApprovedRequestDeptSAPFull       string    `gorm:"column:approved_request_dept_sap_full" json:"-"`
	RefRequestAnnualDriverStatusCode string    `gorm:"column:ref_request_annual_driver_status_code" json:"-"`
	ApprovedRequestDatetime          time.Time `gorm:"column:approved_request_datetime" json:"-"`
	UpdatedAt                        time.Time `gorm:"column:updated_at" json:"-"`
	UpdatedBy                        string    `gorm:"column:updated_by" json:"-"`
}

func (VmsDriverLicenseAnnualApproved) TableName() string {
	return "public.vms_trn_request_annual_driver"
}

//VmsDriverLicenseAnnualApprover
type VmsDriverLicenseAnnualApprover struct {
	TrnRequestAnnualDriverUID string    `gorm:"column:trn_request_annual_driver_uid;primaryKey" json:"trn_request_annual_driver_uid" example:"095fbfbf-378e-4507-b15f-e53ac60370e7"`
	ApproveddRequestEmpId     string    `gorm:"column:approved_request_emp_id" json:"approved_request_emp_id" example:"990003"`
	UpdatedAt                 time.Time `gorm:"column:updated_at" json:"-"`
	UpdatedBy                 string    `gorm:"column:updated_by" json:"-"`
}

func (VmsDriverLicenseAnnualApprover) TableName() string {
	return "public.vms_trn_request_annual_driver"
}

// VmsDriverLicenseAnnualConfirmed
type VmsDriverLicenseAnnualConfirmed struct {
	TrnRequestAnnualDriverUID        string    `gorm:"column:trn_request_annual_driver_uid;primaryKey" json:"trn_request_annual_driver_uid" example:"095fbfbf-378e-4507-b15f-e53ac60370e7"`
	ConfirmedRequestEmpID            string    `gorm:"column:confirmed_request_emp_id" json:"-"`
	ConfirmedRequestEmpName          string    `gorm:"column:confirmed_request_emp_name" json:"-"`
	ConfirmedRequestDeptSAP          string    `gorm:"column:confirmed_request_dept_sap" json:"-"`
	ConfirmedRequestDeptSAPShort     string    `gorm:"column:confirmed_request_dept_sap_short" json:"-"`
	ConfirmedRequestDeptSAPFull      string    `gorm:"column:confirmed_request_dept_sap_full" json:"-"`
	ConfirmedRequestDatetime         time.Time `gorm:"column:confirmed_request_datetime" json:"-"`
	RefRequestAnnualDriverStatusCode string    `gorm:"column:ref_request_annual_driver_status_code" json:"-"`
	UpdatedAt                        time.Time `gorm:"column:updated_at" json:"-"`
	UpdatedBy                        string    `gorm:"column:updated_by" json:"-"`
}

func (VmsDriverLicenseAnnualConfirmed) TableName() string {
	return "public.vms_trn_request_annual_driver"
}
