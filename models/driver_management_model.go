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
	DriverDeptSapShortNameWork     string             `gorm:"column:driver_dept_sap_short_name_work" json:"driver_dept_sap_short_name_work"`
	DriverContactNumber            string             `gorm:"column:driver_contact_number" json:"driver_contact_number"`
	DriverAverageSatisfactionScore float64            `gorm:"column:driver_average_satisfaction_score" json:"driver_average_satisfaction_score"`
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

// VmsMasDriverRequest is a struct that represents a request for driver information in the VMS system.
type VmsMasDriverRequest struct {
	MasDriverUID               string                         `gorm:"primaryKey;column:mas_driver_uid" json:"-"`
	DriverImage                string                         `gorm:"column:driver_image" json:"driver_image" example:"https://example.com/driver_image.jpg"`
	DriverName                 string                         `gorm:"column:driver_name" json:"driver_name" example:"John Doe"`
	DriverNickname             string                         `gorm:"column:driver_nickname" json:"driver_nickname" example:"Johnny"`
	DriverContactNumber        string                         `gorm:"column:driver_contact_number" json:"driver_contact_number" example:"+1234567890"`
	DriverIdentificationNo     string                         `gorm:"column:driver_identification_no" json:"driver_identification_no" example:"ID123456789"`
	DriverBirthdate            time.Time                      `gorm:"column:driver_birthdate" json:"driver_birthdate" example:"1990-01-01T00:00:00Z"`
	WorkType                   int                            `gorm:"column:work_type" json:"work_type" example:"1"`
	ContractNo                 string                         `gorm:"column:contract_no" json:"contract_no" example:"CON123456"`
	DriverDeptSapShortNameHire string                         `gorm:"column:driver_dept_sap_short_name_hire" json:"driver_dept_sap_short_name_hire" example:"HR"`
	MasVendorCode              string                         `gorm:"column:mas_vendor_code" json:"mas_vendor_code" example:"VENDOR123"`
	DriverDeptSapShortNameWork string                         `gorm:"column:driver_dept_sap_short_name_work" json:"driver_dept_sap_short_name_work" example:"กยจ."`
	ApprovedJobDriverStartDate time.Time                      `gorm:"column:approved_job_driver_start_date" json:"approved_job_driver_start_date" example:"2023-01-01T00:00:00Z"`
	ApprovedJobDriverEndDate   time.Time                      `gorm:"column:approved_job_driver_end_date" json:"approved_job_driver_end_date" example:"2023-12-31T23:59:59Z"`
	RefOtherUseCode            string                         `gorm:"column:ref_other_use_code" json:"ref_other_use_code" example:"1"`
	DriverLicense              VmsMasDriverLicenseRequest     `gorm:"foreignKey:MasDriverUID;references:MasDriverUID" json:"driver_license"`
	DriverCertificate          VmsMasDriverCertificateRequest `gorm:"foreignKey:MasDriverUID;references:MasDriverUID" json:"driver_certificate"`
	CreatedAt                  time.Time                      `gorm:"column:created_at" json:"-"`
	CreatedBy                  string                         `gorm:"column:created_by" json:"-"`
	UpdatedAt                  time.Time                      `gorm:"column:updated_at" json:"-"`
	UpdatedBy                  string                         `gorm:"column:updated_by" json:"-"`
	IsDeleted                  string                         `gorm:"column:is_deleted" json:"-"`
	IsActive                   string                         `gorm:"column:is_active" json:"-"`

	IsReplacement string `gorm:"column:is_replacement" json:"-"`
}

func (VmsMasDriverRequest) TableName() string {
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
	MasDriverUID               string                          `gorm:"primaryKey;column:mas_driver_uid" json:"mas_driver_uid"`
	DriverImage                string                          `gorm:"column:driver_image" json:"driver_image" example:"https://example.com/driver_image.jpg"`
	DriverName                 string                          `gorm:"column:driver_name" json:"driver_name" example:"John Doe"`
	DriverNickname             string                          `gorm:"column:driver_nickname" json:"driver_nickname" example:"Johnny"`
	DriverContactNumber        string                          `gorm:"column:driver_contact_number" json:"driver_contact_number" example:"+1234567890"`
	DriverIdentificationNo     string                          `gorm:"column:driver_identification_no" json:"driver_identification_no" example:"ID123456789"`
	DriverBirthdate            time.Time                       `gorm:"column:driver_birthdate" json:"driver_birthdate" example:"1990-01-01T00:00:00Z"`
	WorkType                   int                             `gorm:"column:work_type" json:"work_type" example:"1"`
	ContractNo                 string                          `gorm:"column:contract_no" json:"contract_no" example:"CON123456"`
	DriverDeptSapShortNameHire string                          `gorm:"column:driver_dept_sap_short_name_hire" json:"driver_dept_sap_short_name_hire" example:"HR"`
	MasVendorCode              string                          `gorm:"column:mas_vendor_code" json:"mas_vendor_code" example:"VENDOR123"`
	DriverDeptSapShortNameWork string                          `gorm:"column:driver_dept_sap_short_name_work" json:"driver_dept_sap_short_name_work" example:"กยจ."`
	ApprovedJobDriverStartDate time.Time                       `gorm:"column:approved_job_driver_start_date" json:"approved_job_driver_start_date" example:"2023-01-01T00:00:00Z"`
	ApprovedJobDriverEndDate   time.Time                       `gorm:"column:approved_job_driver_end_date" json:"approved_job_driver_end_date" example:"2023-12-31T23:59:59Z"`
	DriverLicense              VmsMasDriverLicenseResponse     `gorm:"foreignKey:MasDriverUID;references:MasDriverUID" json:"driver_license"`
	DriverCertificate          VmsMasDriverCertificateResponse `gorm:"foreignKey:MasDriverUID;references:MasDriverUID" json:"driver_certificate"`
	CreatedAt                  time.Time                       `gorm:"column:created_at" json:"-"`
	CreatedBy                  string                          `gorm:"column:created_by" json:"-"`
	UpdatedAt                  time.Time                       `gorm:"column:updated_at" json:"-"`
	UpdatedBy                  string                          `gorm:"column:updated_by" json:"-"`
	IsDeleted                  string                          `gorm:"column:is_deleted" json:"-"`
	IsActive                   string                          `gorm:"column:is_active" json:"-"`

	IsReplacement       string             `gorm:"column:is_replacement" json:"is_replacement"`
	RefOtherUseCode     string             `gorm:"column:ref_other_use_code" json:"ref_other_use_code"`
	RefDriverStatusCode int                `gorm:"column:ref_driver_status_code" json:"-"`
	DriverStatus        VmsRefDriverStatus `gorm:"foreignKey:RefDriverStatusCode;references:RefDriverStatusCode" json:"driver_status"`
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
	DriverDeptSapShortNameHire string    `gorm:"column:driver_dept_sap_short_name_hire" json:"driver_dept_sap_short_name_hire" example:"HR"`
	MasVendorCode              string    `gorm:"column:mas_vendor_code" json:"mas_vendor_code" example:"VENDOR123"`
	DriverDeptSapShortNameWork string    `gorm:"column:driver_dept_sap_short_name_work" json:"driver_dept_sap_short_name_work" example:"กยจ."`
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
	MasDriverLicenseUID      string    `gorm:"column:mas_driver_license_uid;primaryKey" json:"mas_driver_license_uid" example:"3e89ebe5-d597-4ee2-b0a1-c3a5628cf131"`
	MasDriverUID             string    `gorm:"column:mas_driver_uid;type:uuid" json:"-"`
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

// VmsMasDriverLeaveStatusUpdate
type VmsMasDriverLeaveStatusUpdate struct {
	MasDriverUID        string    `gorm:"primaryKey;column:mas_driver_uid" json:"mas_driver_uid" example:"8d14e6df-5d65-486e-b079-393d9c817a09"`
	LeaveStartDate      time.Time `gorm:"column:leave_start_date" json:"leave_start_date" example:"2023-01-01T00:00:00Z"`
	LeaveEndDate        time.Time `gorm:"column:leave_end_date" json:"leave_end_date" example:"2023-12-31T23:59:59Z"`
	TimeType            string    `gorm:"column:time_type" json:"time_type" example:"1"`
	LeaveReson          string    `gorm:"column:leave_reason" json:"leave_reason" example:"Vacation"`
	ReplaceMasDriverUID string    `gorm:"column:replace_mas_driver_uid" json:"replace_mas_driver_uid" example:"8d14e6df-5d65-486e-b079-393d9c817a09"`
	RefDriverStatusCode int       `gorm:"column:ref_driver_status_code" json:"ref_driver_status_code" example:"1"`
	UpdatedAt           time.Time `gorm:"column:updated_at" json:"-"`
	UpdatedBy           string    `gorm:"column:updated_by" json:"-"`
}

func (VmsMasDriverLeaveStatusUpdate) TableName() string {
	return "vms_mas_driver"
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
	ReplaceMasDriverUID  string    `gorm:"column:replace_mas_driver_uid" json:"replace_mas_driver_uid" example:"8d14e6df-5d65-486e-b079-393d9c817a09"`
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
