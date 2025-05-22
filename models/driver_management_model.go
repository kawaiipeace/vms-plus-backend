package models

import "time"

// VmsMasDriverList is a struct that represents a list of drivers in the VMS system.
type VmsMasDriverList struct {
	MasDriverUID                   string             `gorm:"primaryKey;column:mas_driver_uid;type:uuid" json:"mas_driver_uid"`
	DriverID                       string             `gorm:"column:driver_id" json:"driver_id"`
	DriverName                     string             `gorm:"column:driver_name" json:"driver_name"`
	DriverImage                    string             `gorm:"column:driver_image" json:"driver_image"`
	DriverNickname                 string             `gorm:"column:driver_nickname" json:"driver_nickname"`
	DriverDeptSapWork              string             `gorm:"column:driver_dept_sap_work" json:"driver_dept_sap_work"`
	DriverDeptSapShortWork         string             `gorm:"column:driver_dept_sap_short_work" json:"driver_dept_sap_short_name_work"`
	DriverContactNumber            string             `gorm:"column:driver_contact_number" json:"driver_contact_number"`
	DriverAverageSatisfactionScore float64            `gorm:"column:driver_average_satisfaction_score" json:"driver_average_satisfaction_score"`
	DriverTotalSatisfactionReview  int                `gorm:"column:driver_total_satisfaction_review" json:"driver_total_satisfaction_review"`
	WorkType                       int                `gorm:"column:work_type" json:"work_type"`
	IsActive                       int                `gorm:"column:is_active" json:"is_active"`
	DriverLicenseEndDate           string             `gorm:"column:driver_license_end_date" json:"driver_license_end_date"`
	ApprovedJobDriverEndDate       time.Time          `gorm:"column:end_date" json:"approved_job_driver_end_date"`
	RefDriverStatusCode            int                `gorm:"column:ref_driver_status_code" json:"-"`
	DriverStatus                   VmsRefDriverStatus `gorm:"foreignKey:RefDriverStatusCode;references:RefDriverStatusCode" json:"driver_status"`
}

func (VmsMasDriverList) TableName() string {
	return "vms_mas_driver"
}

// VmsMasDriverRequest
type VmsMasDriverRequest struct {
	MasDriverUID               string                     `gorm:"primaryKey;column:mas_driver_uid" json:"-"`
	DriverImage                string                     `gorm:"column:driver_image" json:"driver_image" example:"https://example.com/driver_image.jpg"`
	DriverName                 string                     `gorm:"column:driver_name" json:"driver_name" example:"John Doe"`
	DriverID                   string                     `gorm:"column:driver_id" json:"-"`
	DriverNickname             string                     `gorm:"column:driver_nickname" json:"driver_nickname" example:"Johnny"`
	DriverContactNumber        string                     `gorm:"column:driver_contact_number" json:"driver_contact_number" example:"+1234567890"`
	DriverIdentificationNo     string                     `gorm:"column:driver_identification_no" json:"driver_identification_no" example:"ID123456789"`
	DriverBirthdate            time.Time                  `gorm:"column:driver_birthdate" json:"driver_birthdate" example:"1990-01-01T00:00:00Z"`
	WorkType                   int                        `gorm:"column:work_type" json:"work_type" example:"1"`
	IsReplacement              string                     `gorm:"column:is_replacement" json:"is_replacement" example:"1"`
	ContractNo                 string                     `gorm:"column:contract_no" json:"contract_no" example:"CON123456"`
	DriverDeptSapHire          string                     `gorm:"column:driver_dept_sap_hire" json:"driver_dept_sap_hire" example:"1000"`
	DriverDeptSapShortNameHire string                     `gorm:"column:driver_dept_sap_short_name_hire" json:"-"`
	MasVendorCode              string                     `gorm:"column:mas_vendor_code" json:"mas_vendor_code" example:"VENDOR123"`
	DriverDeptSapWork          string                     `gorm:"column:driver_dept_sap_work" json:"driver_dept_sap_work" example:"10001"`
	DriverDeptSapShortWork     string                     `gorm:"column:driver_dept_sap_short_work" json:"-"`
	DriverDeptSapFullWork      string                     `gorm:"column:driver_dept_sap_full_work" json:"-"`
	ApprovedJobDriverStartDate time.Time                  `gorm:"column:approved_job_driver_start_date" json:"approved_job_driver_start_date" example:"2023-01-01T00:00:00Z"`
	ApprovedJobDriverEndDate   time.Time                  `gorm:"column:approved_job_driver_end_date" json:"approved_job_driver_end_date" example:"2023-12-31T23:59:59Z"`
	RefOtherUseCode            string                     `gorm:"column:ref_other_use_code" json:"ref_other_use_code" example:"1"`
	DriverLicense              VmsMasDriverLicenseRequest `gorm:"foreignKey:MasDriverUID;references:MasDriverUID" json:"driver_license"`
	DriverDocuments            []VmsMasDriverDocument     `gorm:"foreignKey:MasDriverUID;references:MasDriverUID" json:"driver_documents"`
	CreatedAt                  time.Time                  `gorm:"column:created_at" json:"-"`
	CreatedBy                  string                     `gorm:"column:created_by" json:"-"`
	UpdatedAt                  time.Time                  `gorm:"column:updated_at" json:"-"`
	UpdatedBy                  string                     `gorm:"column:updated_by" json:"-"`
	IsDeleted                  string                     `gorm:"column:is_deleted" json:"-"`
	IsActive                   string                     `gorm:"column:is_active" json:"-"`
}

func (VmsMasDriverRequest) TableName() string {
	return "vms_mas_driver"
}

// VmsMasDriverImport
type VmsMasDriverImport struct {
	MasDriverUID               string                     `gorm:"primaryKey;column:mas_driver_uid" json:"-"`
	DriverName                 string                     `gorm:"column:driver_name" json:"driver_name" example:"John Doe"`
	DriverID                   string                     `gorm:"column:driver_id" json:"-"`
	DriverNickname             string                     `gorm:"column:driver_nickname" json:"driver_nickname" example:"Johnny"`
	DriverContactNumber        string                     `gorm:"column:driver_contact_number" json:"driver_contact_number" example:"+1234567890"`
	DriverIdentificationNo     string                     `gorm:"column:driver_identification_no" json:"driver_identification_no" example:"ID123456789"`
	DriverBirthdate            time.Time                  `gorm:"column:driver_birthdate" json:"driver_birthdate" example:"1990-01-01T00:00:00Z"`
	WorkType                   int                        `gorm:"column:work_type" json:"work_type" example:"1"`
	IsReplacement              string                     `gorm:"column:is_replacement" json:"is_replacement" example:"1"`
	ContractNo                 string                     `gorm:"column:contract_no" json:"contract_no" example:"CON123456"`
	DriverDeptSapHire          string                     `gorm:"column:driver_dept_sap_hire" json:"driver_dept_sap_hire" example:"HR"`
	DriverDeptSapShortNameHire string                     `gorm:"column:driver_dept_sap_short_name_hire" json:"driver_dept_sap_short_name_hire" example:"HR"`
	MasVendorCode              string                     `gorm:"column:mas_vendor_code" json:"mas_vendor_code" example:"VENDOR123"`
	DriverDeptSapWork          string                     `gorm:"column:driver_dept_sap_work" json:"driver_dept_sap_work" example:"กยจ."`
	DriverDeptSapShortNameWork string                     `gorm:"column:driver_dept_sap_short_work" json:"driver_dept_sap_short_name_work" example:"กยจ."`
	StartDate                  time.Time                  `gorm:"column:start_date" json:"start_date" example:"2023-01-01T00:00:00Z"`
	EndDate                    time.Time                  `gorm:"column:end_date" json:"end_date" example:"2023-12-31T23:59:59Z"`
	RefOtherUseCode            string                     `gorm:"column:ref_other_use_code" json:"ref_other_use_code" example:"1"`
	DriverLicense              VmsMasDriverLicenseRequest `gorm:"foreignKey:MasDriverUID;references:MasDriverUID" json:"driver_license"`
	CreatedAt                  time.Time                  `gorm:"column:created_at" json:"-"`
	CreatedBy                  string                     `gorm:"column:created_by" json:"-"`
	UpdatedAt                  time.Time                  `gorm:"column:updated_at" json:"-"`
	UpdatedBy                  string                     `gorm:"column:updated_by" json:"-"`
	IsDeleted                  string                     `gorm:"column:is_deleted" json:"-"`
	IsActive                   string                     `gorm:"column:is_active" json:"-"`
}

func (VmsMasDriverImport) TableName() string {
	return "vms_mas_driver"
}

// VmsMasDriverLicense is a struct that represents a driver's license information in the VMS system.
type VmsMasDriverLicenseRequest struct {
	MasDriverLicenseUID      string    `gorm:"column:mas_driver_license_uid;primaryKey" json:"-"`
	MasDriverUID             string    `gorm:"column:mas_driver_uid;type:uuid" json:"-"`
	RefDriverLicenseTypeCode string    `gorm:"column:ref_driver_license_type_code" json:"ref_driver_license_type_code" example:"1"`
	DriverLicenseNo          string    `gorm:"column:driver_license_no" json:"driver_license_no" example:"D123456789"`
	DriverLicenseEndDate     time.Time `gorm:"column:driver_license_end_date" json:"driver_license_end_date" example:"2025-12-31T23:59:59Z"`
	DriverLicenseImage       string    `gorm:"column:driver_license_image" json:"driver_license_image" example:"https://example.com/license_image.jpg"`
	DriverLicenseStartDate   time.Time `gorm:"column:driver_license_start_date" json:"driver_license_start_date" example:"2020-01-01T00:00:00Z"`

	CreatedAt time.Time `gorm:"column:created_at" json:"-"`
	CreatedBy string    `gorm:"column:created_by" json:"-"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"-"`
	UpdatedBy string    `gorm:"column:updated_by" json:"-"`
	IsDeleted string    `gorm:"column:is_deleted" json:"-"`
	IsActive  string    `gorm:"column:is_active" json:"-"`
}

func (VmsMasDriverLicenseRequest) TableName() string {
	return "vms_mas_driver_license"
}

// VmsMasDriverLicense is a struct that represents a driver's license information in the VMS system.
type VmsMasDriverCertificateRequest struct {
	MasDriverCertificateUID      string    `gorm:"column:mas_driver_certificate_uid;primaryKey" json:"-"`
	MasDriverUID                 string    `gorm:"column:mas_driver_uid;type:uuid" json:"-"`
	DriverCertificateImage       string    `gorm:"column:driver_certificate_image" json:"driver_certificate_image" example:"https://example.com/certificate_image.jpg"`
	RefDriverCertificateTypeCode string    `gorm:"column:ref_driver_certificate_type_code" json:"ref_driver_certificate_type_code" example:"1"`
	DriverCertificateIssueDate   time.Time `gorm:"column:driver_certificate_issue_date" json:"driver_certificate_issue_date" example:"2023-01-01T00:00:00Z"`
	DriverCertificateExpireDate  time.Time `gorm:"column:driver_certificate_expire_date" json:"driver_certificate_expire_date" example:"2025-12-31T23:59:59Z"`
	CreatedAt                    time.Time `gorm:"column:created_at" json:"-"`
	CreatedBy                    string    `gorm:"column:created_by" json:"-"`
	UpdatedAt                    time.Time `gorm:"column:updated_at" json:"-"`
	UpdatedBy                    string    `gorm:"column:updated_by" json:"-"`
	IsDeleted                    string    `gorm:"column:is_deleted" json:"-"`
	IsActive                     string    `gorm:"column:is_active" json:"-"`
}

func (VmsMasDriverCertificateRequest) TableName() string {
	return "vms_mas_driver_certificate"
}

// VmsMasDriverResponse
type VmsMasDriverResponse struct {
	MasDriverUID               string                      `gorm:"primaryKey;column:mas_driver_uid" json:"mas_driver_uid"`
	DriverImage                string                      `gorm:"column:driver_image" json:"driver_image" example:"https://example.com/driver_image.jpg"`
	DriverName                 string                      `gorm:"column:driver_name" json:"driver_name" example:"John Doe"`
	DriverNickname             string                      `gorm:"column:driver_nickname" json:"driver_nickname" example:"Johnny"`
	DriverContactNumber        string                      `gorm:"column:driver_contact_number" json:"driver_contact_number" example:"+1234567890"`
	DriverIdentificationNo     string                      `gorm:"column:driver_identification_no" json:"driver_identification_no" example:"ID123456789"`
	DriverBirthdate            time.Time                   `gorm:"column:driver_birthdate" json:"driver_birthdate" example:"1990-01-01T00:00:00Z"`
	WorkType                   int                         `gorm:"column:work_type" json:"work_type" example:"1"`
	ContractNo                 string                      `gorm:"column:contract_no" json:"contract_no" example:"CON123456"`
	DriverDeptSapShortNameHire string                      `gorm:"column:driver_dept_sap_short_name_hire" json:"driver_dept_sap_short_name_hire" example:"HR"`
	MasVendorCode              string                      `gorm:"column:mas_vendor_code" json:"mas_vendor_code" example:"VENDOR123"`
	DriverDeptSapShortNameWork string                      `gorm:"column:driver_dept_sap_short_name_work" json:"driver_dept_sap_short_name_work" example:"กยจ."`
	ApprovedJobDriverStartDate time.Time                   `gorm:"column:approved_job_driver_start_date" json:"approved_job_driver_start_date" example:"2023-01-01T00:00:00Z"`
	ApprovedJobDriverEndDate   time.Time                   `gorm:"column:approved_job_driver_end_date" json:"approved_job_driver_end_date" example:"2023-12-31T23:59:59Z"`
	DriverLicense              VmsMasDriverLicenseResponse `gorm:"foreignKey:MasDriverUID;references:MasDriverUID" json:"driver_license"`
	DriverDocuments            []VmsMasDriverDocument      `gorm:"foreignKey:MasDriverUID;references:MasDriverUID" json:"driver_documents"`

	CreatedAt time.Time `gorm:"column:created_at" json:"-"`
	CreatedBy string    `gorm:"column:created_by" json:"-"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"-"`
	UpdatedBy string    `gorm:"column:updated_by" json:"-"`
	IsDeleted string    `gorm:"column:is_deleted" json:"-"`
	IsActive  string    `gorm:"column:is_active" json:"-"`

	IsReplacement         string             `gorm:"column:is_replacement" json:"is_replacement"`
	RefOtherUseCode       string             `gorm:"column:ref_other_use_code" json:"ref_other_use_code"`
	RefDriverStatusCode   int                `gorm:"column:ref_driver_status_code" json:"-"`
	DriverStatus          VmsRefDriverStatus `gorm:"foreignKey:RefDriverStatusCode;references:RefDriverStatusCode" json:"driver_status"`
	AlertDriverStatus     string             `gorm:"column:alert_driver_status" json:"alert_driver_status"`
	AlertDriverStatusDesc string             `gorm:"column:alert_driver_status_desc" json:"alert_driver_status_desc"`
}

func (VmsMasDriverResponse) TableName() string {
	return "vms_mas_driver"
}

// VmsMasDriverLicenseResponse
type VmsMasDriverLicenseResponse struct {
	MasDriverLicenseUID      string                  `gorm:"column:mas_driver_license_uid;primaryKey" json:"mas_driver_license_uid"`
	MasDriverUID             string                  `gorm:"column:mas_driver_uid;type:uuid" json:"-"`
	RefDriverLicenseTypeCode string                  `gorm:"column:ref_driver_license_type_code" json:"ref_driver_license_type_code" example:"1"`
	DriverLicenseNo          string                  `gorm:"column:driver_license_no" json:"driver_license_no" example:"D123456789"`
	DriverLicenseEndDate     time.Time               `gorm:"column:driver_license_end_date" json:"driver_license_end_date" example:"2025-12-31T23:59:59Z"`
	DriverLicenseImage       string                  `gorm:"column:driver_license_image" json:"driver_license_image" example:"https://example.com/license_image.jpg"`
	DriverLicenseStartDate   time.Time               `gorm:"column:driver_license_start_date" json:"driver_license_start_date" example:"2020-01-01T00:00:00Z"`
	DriverLicenseType        VmsRefDriverLicenseType `gorm:"foreignKey:RefDriverLicenseTypeCode;references:RefDriverLicenseTypeCode" json:"driver_license_type"`
}

func (VmsMasDriverLicenseResponse) TableName() string {
	return "vms_mas_driver_license"
}

// VmsMasDriverCertificateResponse
type VmsMasDriverCertificateResponse struct {
	MasDriverCertificateUid      string    `gorm:"column:mas_driver_certificate_uid;primaryKey" json:"mas_driver_certificate_uid"`
	MasDriverUID                 string    `gorm:"column:mas_driver_uid;type:uuid" json:"-"`
	DriverCertificateImage       string    `gorm:"column:driver_certificate_image" json:"driver_certificate_image" example:"https://example.com/certificate_image.jpg"`
	RefDriverCertificateTypeCode string    `gorm:"column:ref_driver_certificate_type_code" json:"ref_driver_certificate_type_code" example:"1"`
	DriverCertificateIssueDate   time.Time `gorm:"column:driver_certificate_issue_date" json:"driver_certificate_issue_date" example:"2023-01-01T00:00:00Z"`
	DriverCertificateExpireDate  time.Time `gorm:"column:driver_certificate_expire_date" json:"driver_certificate_expire_date" example:"2025-12-31T23:59:59Z"`
}

func (VmsMasDriverCertificateResponse) TableName() string {
	return "vms_mas_driver_certificate"
}

// VmsMasDriverDetail
type VmsMasDriverDetailUpdate struct {
	MasDriverUID           string    `gorm:"primaryKey;column:mas_driver_uid" json:"mas_driver_uid" example:"8d14e6df-5d65-486e-b079-393d9c817a09"`
	DriverImage            string    `gorm:"column:driver_image" json:"driver_image" example:"https://example.com/driver_image.jpg"`
	DriverName             string    `gorm:"column:driver_name" json:"driver_name" example:"John Doe"`
	DriverNickname         string    `gorm:"column:driver_nickname" json:"driver_nickname" example:"Johnny"`
	DriverContactNumber    string    `gorm:"column:driver_contact_number" json:"driver_contact_number" example:"+1234567890"`
	DriverIdentificationNo string    `gorm:"column:driver_identification_no" json:"driver_identification_no" example:"ID123456789"`
	DriverBirthdate        time.Time `gorm:"column:driver_birthdate" json:"driver_birthdate" example:"1990-01-01T00:00:00Z"`
	WorkType               int       `gorm:"column:work_type" json:"work_type" example:"1"`
	UpdatedAt              time.Time `gorm:"column:updated_at" json:"-"`
	UpdatedBy              string    `gorm:"column:updated_by" json:"-"`
}

func (VmsMasDriverDetailUpdate) TableName() string {
	return "vms_mas_driver"
}

// VmsMasDriverContract
type VmsMasDriverContractUpdate struct {
	MasDriverUID               string    `gorm:"primaryKey;column:mas_driver_uid" json:"mas_driver_uid" example:"8d14e6df-5d65-486e-b079-393d9c817a09"`
	ContractNo                 string    `gorm:"column:contract_no" json:"contract_no" example:"CON123456"`
	DriverDeptSapHire          string    `gorm:"column:driver_dept_sap_hire" json:"driver_dept_sap_hire" example:"1000"`
	DriverDeptSapShortNameHire string    `gorm:"column:driver_dept_sap_short_name_hire" json:"-"`
	MasVendorCode              string    `gorm:"column:mas_vendor_code" json:"mas_vendor_code" example:"VENDOR123"`
	DriverDeptSapWork          string    `gorm:"column:driver_dept_sap_work" json:"driver_dept_sap_work" example:"10001"`
	DriverDeptSapShortWork     string    `gorm:"column:driver_dept_sap_short_work" json:"-"`
	DriverDeptSapFullWork      string    `gorm:"column:driver_dept_sap_full_work" json:"-"`
	ApprovedJobDriverStartDate time.Time `gorm:"column:approved_job_driver_start_date" json:"approved_job_driver_start_date" example:"2023-01-01T00:00:00Z"`
	ApprovedJobDriverEndDate   time.Time `gorm:"column:approved_job_driver_end_date" json:"approved_job_driver_end_date" example:"2023-12-31T23:59:59Z"`
	RefOtherUseCode            int       `gorm:"column:ref_other_use_code" json:"ref_other_use_code" example:"1"`
	UpdatedAt                  time.Time `gorm:"column:updated_at" json:"-"`
	UpdatedBy                  string    `gorm:"column:updated_by" json:"-"`
}

func (VmsMasDriverContractUpdate) TableName() string {
	return "vms_mas_driver"
}

// VmsMasDriverLicenseUpdate
type VmsMasDriverLicenseUpdate struct {
	MasDriverLicenseUID      string    `gorm:"column:mas_driver_license_uid;primaryKey" json:"-"`
	MasDriverUID             string    `gorm:"column:mas_driver_uid" json:"mas_driver_uid"  example:"3e89ebe5-d597-4ee2-b0a1-c3a5628cf131"`
	RefDriverLicenseTypeCode string    `gorm:"column:ref_driver_license_type_code" json:"ref_driver_license_type_code" example:"1"`
	DriverLicenseNo          string    `gorm:"column:driver_license_no" json:"driver_license_no" example:"D123456789"`
	DriverLicenseStartDate   time.Time `gorm:"column:driver_license_start_date" json:"driver_license_start_date" example:"2020-01-01T00:00:00Z"`
	DriverLicenseEndDate     time.Time `gorm:"column:driver_license_end_date" json:"driver_license_end_date" example:"2025-12-31T23:59:59Z"`

	UpdatedAt time.Time `gorm:"column:updated_at" json:"-"`
	UpdatedBy string    `gorm:"column:updated_by" json:"-"`
}

func (VmsMasDriverLicenseUpdate) TableName() string {
	return "vms_mas_driver_license"
}

// VmsMasDriverRequest is a struct that represents a request for driver information in the VMS system.
type VmsMasDriverDocumentUpdate struct {
	MasDriverUID    string                 `gorm:"primaryKey;column:mas_driver_uid" json:"mas_driver_uid"  example:"d6c76da4-fbf9-44e1-a22b-ddc887b1939e"`
	DriverLicense   VmsMasDriverDocument   `gorm:"foreignKey:MasDriverUID;references:MasDriverUID" json:"driver_license"`
	DriverDocuments []VmsMasDriverDocument `gorm:"foreignKey:MasDriverUID;references:MasDriverUID" json:"driver_documents"`
	UpdatedAt       time.Time              `gorm:"column:updated_at" json:"-"`
	UpdatedBy       string                 `gorm:"column:updated_by" json:"-"`
}

func (VmsMasDriverDocumentUpdate) TableName() string {
	return "vms_mas_driver"
}

// VmsMasDriverLeaveStatusUpdate
type VmsMasDriverLeaveStatusUpdate struct {
	TrnDriverLeaveUID    string    `gorm:"column:trn_driver_leave_uid;primaryKey" json:"-"`
	MasDriverUID         string    `gorm:"column:mas_driver_uid;type:uuid" json:"mas_driver_uid" example:"8d14e6df-5d65-486e-b079-393d9c817a09"`
	RefDriverStatusCode  int       `gorm:"column:ref_driver_status_code" json:"-"`
	LeaveStartDate       time.Time `gorm:"column:leave_start_date" json:"leave_start_date" example:"2025-01-25T00:00:00Z"`
	LeaveEndDate         time.Time `gorm:"column:leave_end_date" json:"leave_end_date" example:"2025-01-30T23:59:59Z"`
	LeaveTimeTypeCode    int16     `gorm:"column:leave_time_type_code" json:"leave_time_type_code" example:"1"`
	LeaveReason          string    `gorm:"column:leave_reason" json:"leave_reason" example:"Sick leave"`
	ReplacementDriverUID string    `gorm:"column:replacement_driver_uid;type:uuid" json:"replacement_driver_uid" example:"0a33f4df-5da8-4831-b3e4-27b5c6134c7c"`
	CreatedAt            time.Time `gorm:"column:created_at" json:"-"`
	CreatedBy            string    `gorm:"column:created_by" json:"-"`
	UpdatedAt            time.Time `gorm:"column:updated_at" json:"-"`
	UpdatedBy            string    `gorm:"column:updated_by" json:"-"`
	IsDeleted            string    `gorm:"column:is_deleted" json:"-"`
}

func (VmsMasDriverLeaveStatusUpdate) TableName() string {
	return "vms_trn_driver_leave"
}

// VmsMasDriverIsActiveUpdate
type VmsMasDriverIsActiveUpdate struct {
	MasDriverUID string    `gorm:"primaryKey;column:mas_driver_uid" json:"mas_driver_uid" example:"8d14e6df-5d65-486e-b079-393d9c817a09"`
	IsActive     string    `gorm:"column:is_active" json:"is_active" example:"1"`
	UpdatedAt    time.Time `gorm:"column:updated_at" json:"-"`
	UpdatedBy    string    `gorm:"column:updated_by" json:"-"`
}

func (VmsMasDriverIsActiveUpdate) TableName() string {
	return "vms_mas_driver"
}

// VmsMasDriverDelete
type VmsMasDriverDelete struct {
	MasDriverUID string    `gorm:"primaryKey;column:mas_driver_uid" json:"mas_driver_uid" example:"8d14e6df-5d65-486e-b079-393d9c817a09"`
	DriverName   string    `gorm:"column:driver_name" json:"driver_name" example:"John Doe"`
	UpdatedAt    time.Time `gorm:"column:updated_at" json:"-"`
	UpdatedBy    string    `gorm:"column:updated_by" json:"-"`
}

func (VmsMasDriverDelete) TableName() string {
	return "vms_mas_driver"
}

// VmsMasDriverLayoff
type VmsMasDriverLayoffStatusUpdate struct {
	MasDriverUID         string    `gorm:"primaryKey;column:mas_driver_uid" json:"mas_driver_uid" example:"8d14e6df-5d65-486e-b079-393d9c817a09"`
	RefDriverStatusCode  int       `gorm:"column:ref_driver_status_code" json:"-"`
	ReplacedMasDriverUid string    `gorm:"column:replaced_mas_driver_uid" json:"replaced_mas_driver_uid" example:"8d14e6df-5d65-486e-b079-393d9c817a09"`
	UpdatedAt            time.Time `gorm:"column:updated_at" json:"-"`
	UpdatedBy            string    `gorm:"column:updated_by" json:"-"`
}

func (VmsMasDriverLayoffStatusUpdate) TableName() string {
	return "vms_mas_driver"
}

// VmsMasDriverResignStatusUpdate
type VmsMasDriverResignStatusUpdate struct {
	MasDriverUID         string    `gorm:"primaryKey;column:mas_driver_uid" json:"mas_driver_uid" example:"8d14e6df-5d65-486e-b079-393d9c817a09"`
	RefDriverStatusCode  int       `gorm:"column:ref_driver_status_code" json:"-"`
	ReplacedMasDriverUid string    `gorm:"column:replaced_mas_driver_uid" json:"replaced_mas_driver_uid" example:"8d14e6df-5d65-486e-b079-393d9c817a09"`
	UpdatedAt            time.Time `gorm:"column:updated_at" json:"-"`
	UpdatedBy            string    `gorm:"column:updated_by" json:"-"`
}

func (VmsMasDriverResignStatusUpdate) TableName() string {
	return "vms_mas_driver"
}

// VmsMasDriverDocument is a struct that represents a driver's document information in the VMS system.
type VmsMasDriverDocument struct {
	MasDriverDocumentUID string    `gorm:"column:mas_driver_document_uid;primaryKey;type:uuid" json:"-"`
	MasDriverUID         string    `gorm:"column:mas_driver_uid;type:uuid" json:"-"`
	DriverDocumentNo     int       `gorm:"column:driver_document_no" json:"driver_document_no" example:"1"`
	DriverDocumentName   string    `gorm:"column:driver_document_name" json:"driver_document_name" example:"CardID.pdf"`
	DriverDocumentFile   string    `gorm:"column:driver_document_file" json:"driver_document_file" example:"https://example.com/document.pdf"`
	CreatedAt            time.Time `gorm:"column:created_at" json:"-"`
	CreatedBy            string    `gorm:"column:created_by" json:"-"`
	UpdatedAt            time.Time `gorm:"column:updated_at" json:"-"`
	UpdatedBy            string    `gorm:"column:updated_by" json:"-"`
	IsDeleted            string    `gorm:"column:is_deleted" json:"-"`
}

func (VmsMasDriverDocument) TableName() string {
	return "vms_mas_driver_document"
}

type DriverTimeLine struct {
	MasDriverUID               string             `gorm:"column:mas_driver_uid;primaryKey" json:"mas_driver_uid" example:"8d14e6df-5d65-486e-b079-393d9c817a09"`
	DriverName                 string             `gorm:"column:driver_name" json:"driver_name" example:"John Doe"`
	DriverNickname             string             `gorm:"column:driver_nickname" json:"driver_nickname" example:"Johnny"`
	DriverContactNumber        string             `gorm:"column:driver_contact_number" json:"driver_contact_number" example:"+1234567890"`
	DriverDeptSapShortNameWork string             `gorm:"column:driver_dept_sap_short_name_work" json:"driver_dept_sap_short_name_work" example:"กยจ."`
	WorkLastMonth              string             `gorm:"column:work_last_month" json:"work_last_month" example:"22 วัน/3 งาน"`
	WorkThisMonth              string             `gorm:"column:work_this_month" json:"work_this_month" example:"16 วัน/2 งาน"`
	DriverTrnRequests          []DriverTrnRequest `gorm:"-" json:"driver_trn_requests"`
}

type DriverTrnRequest struct {
	TrnRequestUID          string             `gorm:"column:trn_request_uid" json:"trn_request_uid"`
	MasDriverUID           string             `gorm:"column:mas_carpool_driver_uid" json:"mas_carpool_driver_uid"`
	RequestNo              string             `gorm:"column:request_no" json:"request_no"`
	ReserveStartDatetime   string             `gorm:"column:reserve_start_datetime" json:"start_datetime"`
	ReserveEndDatetime     string             `gorm:"column:reserve_end_datetime" json:"end_datetime"`
	RefRequestStatusCode   string             `gorm:"column:ref_request_status_code" json:"ref_request_status_code"`
	RefRequestStatusName   string             `json:"ref_request_status_name"`
	RefTripTypeCode        int                `gorm:"column:ref_trip_type_code" json:"ref_trip_type_code"`
	VehicleUserEmpID       string             `gorm:"column:vehicle_user_emp_id" json:"vehicle_user_emp_id" example:"990001"`
	VehicleUserEmpName     string             `gorm:"column:vehicle_user_emp_name" json:"vehicle_user_emp_name"`
	VehicleUserDeptSAP     string             `gorm:"column:vehicle_user_dept_sap" json:"vehicle_user_dept_sap"`
	VehicleUserDeskPhone   string             `gorm:"column:vehicle_user_desk_phone" json:"car_user_internal_contact_number" example:"1122"`
	VehicleUserMobilePhone string             `gorm:"column:vehicle_user_mobile_phone" json:"car_user_mobile_contact_number" example:"0987654321"`
	TripDetails            []VmsTrnTripDetail `gorm:"foreignKey:TrnRequestUID;references:TrnRequestUID" json:"trip_details"`
	TimeLineStatus         string             `gorm:"-" json:"time_line_status"`
}

func (DriverTrnRequest) TableName() string {
	return "public.vms_trn_request"
}

type DriverWorkReport struct {
	MasDriverUID                     string    `json:"mas_driver_uid"`
	DriverName                       string    `json:"driver_name"`
	DriverNickname                   string    `json:"driver_nickname"`
	DriverID                         string    `json:"driver_id"`
	DriverDeptSapShortWork           string    `json:"driver_dept_sap_short_work"`
	DriverDeptSapFullWork            string    `json:"driver_dept_sap_full_work"`
	ReserveStartDatetime             time.Time `json:"reserve_start_datetime"`
	ReserveEndDatetime               time.Time `json:"reserve_end_datetime"`
	WorkType                         string    `json:"work_type"`
	VehicleLicensePlate              string    `json:"vehicle_license_plate"`
	VehicleLicensePlateProvinceShort string    `json:"vehicle_license_plate_province_short"`
	VehicleLicensePlateProvinceFull  string    `json:"vehicle_license_plate_province_full"`
	VehicleCarTypeDetail             string    `json:"vehicle_car_type_detail"`
	TripStartDatetime                time.Time `json:"trip_start_datetime"`
	TripEndDatetime                  time.Time `json:"trip_end_datetime"`
	TripDeparturePlace               string    `json:"trip_departure_place"`
	TripDestinationPlace             string    `json:"trip_destination_place"`
	TripStartMiles                   float64   `json:"trip_start_miles"`
	TripEndMiles                     float64   `json:"trip_end_miles"`
	TripDetail                       string    `json:"trip_detail"`
}
