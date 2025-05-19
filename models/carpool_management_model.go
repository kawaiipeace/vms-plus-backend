package models

import (
	"time"
)

// VmsMasCarpoolList
type VmsMasCarpoolList struct {
	MasCarpoolUID        string `gorm:"column:mas_carpool_uid" json:"mas_carpool_uid"`
	CarpoolName          string `gorm:"column:carpool_name" json:"carpool_name"`
	CarpoolDeptSap       string `gorm:"column:carpool_dept_sap" json:"carpool_dept_sap"`
	CarpoolContactPlace  string `gorm:"column:carpool_contact_place" json:"carpool_contact_place"`
	CarpoolContactNumber string `gorm:"column:carpool_contact_number" json:"carpool_contact_number" `
	NumberOfDrivers      int    `gorm:"number_of_drivers" json:"number_of_drivers" `
	NumberOfVehicles     int    `gorm:"number_of_vehicles" json:"number_of_vehicles" `
	NumberOfApprovers    int    `gorm:"numberOfApprovers" json:"numberOfApprovers" `
	IsActive             string `gorm:"column:is_active" json:"is_active"`
	CarpoolStatus        string `gorm:"-" json:"carpool_status"`
}

func (VmsMasCarpoolList) TableName() string {
	return "vms_mas_carpool"
}

// VmsMasCarpoolRequest
type VmsMasCarpoolRequest struct {
	MasCarpoolUID            string                        `gorm:"column:mas_carpool_uid;primaryKey" json:"-"`
	CarpoolName              string                        `gorm:"column:carpool_name" json:"carpool_name" example:"carpool_name"`
	CarpoolContactPlace      string                        `gorm:"column:carpool_contact_place" json:"carpool_contact_place" example:"city"`
	CarpoolDeptSap           string                        `gorm:"column:carpool_dept_sap"  json:"carpool_dept_sap" example:"10001"`
	CarpoolContactNumber     string                        `gorm:"column:carpool_contact_number" json:"carpool_contact_number" example:"111"`
	CarpoolMainBusinessArea  string                        `gorm:"column:carpool_main_business_area" json:"-"`
	Remark                   string                        `gorm:"column:remark" json:"remark" example:"remark"`
	RefCarpoolChooseCarID    int                           `gorm:"column:ref_carpool_choose_car_id" json:"ref_carpool_choose_car_id" example:"1"`
	RefCarpoolChooseDriverID int                           `gorm:"column:ref_carpool_choose_driver_id" json:"ref_carpool_choose_driver_id" example:"1"`
	CarpoolType              string                        `gorm:"column:carpool_type" json:"carpool_type" example:"1"`
	CarpoolAuthorizedDepts   []VmsMasCarpoolAuthorizedDept `gorm:"foreignKey:MasCarpoolUID;references:MasCarpoolUID" json:"carpool_authorized_depts"`
	IsHaveDriverForCarpool   string                        `gorm:"column:is_have_driver_for_carpool" json:"-"`
	IsMustPassStatus30       string                        `gorm:"column:is_must_pass_status_30" json:"is_must_pass_status_30" example:"0"`
	IsMustPassStatus40       string                        `gorm:"column:is_must_pass_status_40" json:"is_must_pass_status_40" example:"0"`
	IsMustPassStatus50       string                        `gorm:"column:is_must_pass_status_50" json:"is_must_pass_status_50" example:"0"`
	IsActive                 string                        `gorm:"column:is_active" json:"is_active"`
	CreatedAt                time.Time                     `gorm:"column:created_at;autoCreateTime" json:"-"`
	UpdatedAt                time.Time                     `gorm:"column:updated_at;autoUpdateTime" json:"-"`
	CreatedBy                string                        `gorm:"column:created_by" json:"-"`
	UpdatedBy                string                        `gorm:"column:updated_by" json:"-"`
}

func (VmsMasCarpoolRequest) TableName() string {
	return "vms_mas_carpool"
}

// VmsMasCarpoolRequest
type VmsMasCarpoolActive struct {
	MasCarpoolUID string `gorm:"column:mas_carpool_uid;primaryKey" json:"mas_carpool_uid" example:"164632c9-1d33-477e-b335-97a4e79a5845"`
	IsActive      string `gorm:"column:is_active" json:"is_active" example:"1"`
	CreatedBy     string `gorm:"column:created_by" json:"-"`
	UpdatedBy     string `gorm:"column:updated_by" json:"-"`
}

func (VmsMasCarpoolActive) TableName() string {
	return "vms_mas_carpool"
}

// VmsMasCarpoolResponse
type VmsMasCarpoolResponse struct {
	MasCarpoolUID            string                                `gorm:"column:mas_carpool_uid;primaryKey" json:"mas_carpool_uid"`
	CarpoolName              string                                `gorm:"column:carpool_name" json:"carpool_name" example:"carpool_name"`
	CarpoolContactPlace      string                                `gorm:"column:carpool_contact_place" json:"carpool_contact_place" example:"city"`
	CarpoolDeptSap           string                                `gorm:"column:carpool_dept_sap"`
	CarpoolContactNumber     string                                `gorm:"column:carpool_contact_number" json:"carpool_contact_number" example:"111"`
	CarpoolMainBusinessArea  string                                `gorm:"column:carpool_main_business_area" json:"carpool_main_business_area" example:"0000"`
	Remark                   string                                `gorm:"column:remark" json:"remark" example:"remark"`
	RefCarpoolChooseCarID    int                                   `gorm:"column:ref_carpool_choose_car_id" json:"ref_carpool_choose_car_id" example:"1"`
	RefCarpoolChooseDriverID int                                   `gorm:"column:ref_carpool_choose_driver_id" json:"ref_carpool_choose_driver_id" example:"1"`
	CarpoolType              string                                `gorm:"column:carpool_type" json:"-"`
	IsHaveDriverForCarpool   string                                `gorm:"column:is_have_driver_for_carpool" json:"-"`
	IsMustPassStatus30       string                                `gorm:"column:is_must_pass_status_30" json:"is_must_pass_status_30" example:"0"`
	IsMustPassStatus40       string                                `gorm:"column:is_must_pass_status_40" json:"is_must_pass_status_40" example:"0"`
	IsMustPassStatus50       string                                `gorm:"column:is_must_pass_status_50" json:"is_must_pass_status_50" example:"0"`
	IsActive                 string                                `gorm:"column:is_active" json:"is_active"`
	CarpoolChooseDriver      VmsRefCarpoolChooseDriver             `gorm:"foreignKey:RefCarpoolChooseDriverID;references:RefCarpoolChooseDriverID" json:"carpool_choose_driver"`
	CarpoolChooseCar         VmsRefCarpoolChooseCar                `gorm:"foreignKey:RefCarpoolChooseCarID;references:RefCarpoolChooseCarID" json:"carpool_choose_car"`
	CarpoolAuthorizedDepts   []VmsMasCarpoolAuthorizedDeptResponse `gorm:"foreignKey:MasCarpoolUID;references:MasCarpoolUID" json:"carpool_authorized_depts"`
	CarpoolAdmins            []VmsMasCarpoolAdmin                  `gorm:"foreignKey:MasCarpoolUID;references:MasCarpoolUID" json:"carpool_admins"`
	CarpoolApprovers         []VmsMasCarpoolApprover               `gorm:"foreignKey:MasCarpoolUID;references:MasCarpoolUID" json:"carpool_approvers"`
	CarpoolVehicles          []VmsMasCarpoolVehicle                `gorm:"foreignKey:MasCarpoolUID;references:MasCarpoolUID" json:"carpool_vehicles"`
	CarpoolDrivers           []VmsMasCarpoolDriver                 `gorm:"foreignKey:MasCarpoolUID;references:MasCarpoolUID" json:"carpool_drivers"`

	IsCarpoolReady        bool `gorm:"column:is_carpool_ready" json:"is_carpool_ready"`
	IsCarpoolChooseDriver bool `gorm:"column:is_carpool_choose_driver" json:"is_carpool_choose_driver"`
}

func (VmsMasCarpoolResponse) TableName() string {
	return "vms_mas_carpool"
}

// VmsMasCarpoolDelete
type VmsMasCarpoolDelete struct {
	MasCarpoolUID string    `gorm:"column:mas_carpool_uid;primaryKey" json:"mas_carpool_uid" example:"164632c9-1d33-477e-b335-97a4e79a5845"`
	CarpoolName   string    `gorm:"column:carpool_name" json:"carpool_name" example:"carpool name 1"`
	UpdatedAt     time.Time `gorm:"column:updated_at" json:"-"`
	UpdatedBy     string    `gorm:"column:updated_by" json:"-"`
}

func (VmsMasCarpoolDelete) TableName() string {
	return "vms_mas_carpool"
}

// VmsMasCarpoolAdminList
type VmsMasCarpoolAdminList struct {
	MasCarpoolAdminUID    string `gorm:"column:mas_carpool_admin_uid;primaryKey" json:"mas_carpool_admin_uid"`
	MasCarpoolUID         string `gorm:"column:mas_carpool_uid" json:"mas_carpool_uid" example:"164632c9-1d33-477e-b335-97a4e79a5845"`
	AdminEmpNo            string `gorm:"column:admin_emp_no" json:"admin_emp_no" example:"990003"`
	AdminEmpName          string `gorm:"column:admin_emp_name" json:"admin_emp_name" example:"emp name"`
	AdminDeptSap          string `gorm:"column:admin_dept_sap" json:"admin_dept_sap"`
	AdminDeptSapShort     string `gorm:"column:admin_dept_sap_short" json:"admin_dept_sap_short"`
	InternalContactNumber string `gorm:"column:internal_contact_number" json:"internal_contact_number" example:"1234"`
	MobileContactNumber   string `gorm:"column:mobile_contact_number" json:"mobile_contact_number" example:"9876543210"`
	IsMainAdmin           string `gorm:"column:is_main_admin" json:"is_main_admin"`
	IsActive              string `gorm:"column:is_active" json:"is_active"`
}

// TableName sets the table name for the VmsMasCarpoolAdmin model
func (VmsMasCarpoolAdminList) TableName() string {
	return "vms_mas_carpool_admin"
}

// VmsMasCarpoolAdmin
type VmsMasCarpoolAdmin struct {
	MasCarpoolAdminUID    string    `gorm:"column:mas_carpool_admin_uid;primaryKey" json:"-"`
	MasCarpoolUID         string    `gorm:"column:mas_carpool_uid" json:"mas_carpool_uid" example:"164632c9-1d33-477e-b335-97a4e79a5845"`
	AdminEmpNo            string    `gorm:"column:admin_emp_no" json:"admin_emp_no" example:"990003"`
	AdminDeptSap          string    `gorm:"column:admin_dept_sap" json:"-"`
	InternalContactNumber string    `gorm:"column:internal_contact_number" json:"internal_contact_number" example:"1234"`
	MobileContactNumber   string    `gorm:"column:mobile_contact_number" json:"mobile_contact_number" example:"9876543210"`
	IsMainAdmin           string    `gorm:"column:is_main_admin" json:"-"`
	IsActive              string    `gorm:"column:is_active" json:"-"`
	IsDeleted             string    `gorm:"column:is_deleted" json:"-"`
	CreatedAt             time.Time `gorm:"column:created_at;autoCreateTime" json:"-"`
	UpdatedAt             time.Time `gorm:"column:updated_at;autoUpdateTime" json:"-"`
	CreatedBy             string    `gorm:"column:created_by" json:"-"`
	UpdatedBy             string    `gorm:"column:updated_by" json:"-"`
}

// TableName sets the table name for the VmsMasCarpoolAdmin model
func (VmsMasCarpoolAdmin) TableName() string {
	return "vms_mas_carpool_admin"
}

// VmsMasCarpoolApproverList
type VmsMasCarpoolApproverList struct {
	MasCarpoolApproverUID string `gorm:"column:mas_carpool_approver_uid;primaryKey" json:"mas_carpool_approver_uid"`
	MasCarpoolUID         string `gorm:"column:mas_carpool_uid" json:"mas_carpool_uid" example:"164632c9-1d33-477e-b335-97a4e79a5845"`
	ApproverEmpNo         string `gorm:"column:approver_emp_no" json:"approver_emp_no" example:"990004"`
	ApproverEmpName       string `gorm:"column:approver_emp_name" json:"approver_emp_name" example:"approver name"`
	ApproverDeptSap       string `gorm:"column:approver_dept_sap" json:"approver_dept_sap"`
	ApproverDeptSapShort  string `gorm:"column:approver_dept_sap_short" json:"approver_dept_sap_short"`
	InternalContactNumber string `gorm:"column:internal_contact_number" json:"internal_contact_number" example:"5678"`
	MobileContactNumber   string `gorm:"column:mobile_contact_number" json:"mobile_contact_number" example:"9876543211"`
	IsMainApprover        string `gorm:"column:is_main_approver" json:"is_main_approver"`
	IsActive              string `gorm:"column:is_active" json:"is_active"`
}

func (VmsMasCarpoolApproverList) TableName() string {
	return "vms_mas_carpool_approver"
}

// VmsMasCarpoolApprover
type VmsMasCarpoolApprover struct {
	MasCarpoolApproverUID string    `gorm:"column:mas_carpool_approver_uid;primaryKey" json:"-"`
	MasCarpoolUID         string    `gorm:"column:mas_carpool_uid" json:"mas_carpool_uid" example:"164632c9-1d33-477e-b335-97a4e79a5845"`
	ApproverEmpNo         string    `gorm:"column:approver_emp_no" json:"approver_emp_no" example:"990004"`
	ApproverDeptSap       string    `gorm:"column:approver_dept_sap" json:"-"`
	InternalContactNumber string    `gorm:"column:internal_contact_number" json:"internal_contact_number" example:"5678"`
	MobileContactNumber   string    `gorm:"column:mobile_contact_number" json:"mobile_contact_number" example:"9876543211"`
	IsMainApprover        string    `gorm:"column:is_main_approver" json:"-"`
	IsActive              string    `gorm:"column:is_active" json:"-"`
	IsDeleted             string    `gorm:"column:is_deleted" json:"-"`
	CreatedAt             time.Time `gorm:"column:created_at;autoCreateTime" json:"-"`
	UpdatedAt             time.Time `gorm:"column:updated_at;autoUpdateTime" json:"-"`
	CreatedBy             string    `gorm:"column:created_by" json:"-"`
	UpdatedBy             string    `gorm:"column:updated_by" json:"-"`
}

func (VmsMasCarpoolApprover) TableName() string {
	return "vms_mas_carpool_approver"
}

// VmsMasCarpoolVehicleList
type VmsMasCarpoolVehicleList struct {
	MasCarpoolVehicleUID string `gorm:"column:mas_carpool_vehicle_uid;primaryKey" json:"mas_carpool_vehicle_uid"`
	MasCarpoolUID        string `gorm:"column:mas_carpool_uid" json:"mas_carpool_uid" example:"164632c9-1d33-477e-b335-97a4e79a5845"`
	VmsMasCarpoolVehicleDetail
}

func (VmsMasCarpoolVehicleList) TableName() string {
	return "vms_mas_carpool_vehicle"
}

// VmsMasCarpoolVehicleDetail
type VmsMasCarpoolVehicleDetail struct {
	MasVehicleUID        string    `gorm:"column:mas_vehicle_uid;primary" json:"mas_vehicle_uid" example:"334632c9-1d33-477e-b335-97a4e79a5845"`
	IsActive             string    `gorm:"column:is_active" json:"is_active"`
	VehicleLicensePlate  string    `gorm:"column:vehicle_license_plate" json:"vehicle_license_plate"`
	VehicleBrandName     string    `gorm:"column:vehicle_brand_name" json:"vehicle_brand_name"`
	VehicleModelName     string    `gorm:"column:vehicle_model_name" json:"vehicle_model_name"`
	RefVehicleTypeCode   string    `gorm:"column:ref_vehicle_type_code" json:"ref_vehicle_type_code"`
	RefVehicleTypeName   string    `gorm:"column:ref_vehicle_type_name" json:"ref_vehicle_type_name"`
	VehicleOwnerDeptSAP  string    `gorm:"column:vehicle_owner_dept_short" json:"vehicle_owner_dept_short"`
	FleetCardNo          string    `gorm:"column:fleet_card_no" json:"fleet_card_no"`
	IsTaxCredit          bool      `gorm:"column:is_tax_credit" json:"is_tax_credit"`
	VehicleMileage       float64   `gorm:"column:vehicle_mileage" json:"vehicle_mileage"`
	VehicleGetDate       time.Time `gorm:"column:vehicle_get_date" json:"vehicle_get_date"` // Changed to time.Time
	RefVehicleStatusCode string    `gorm:"column:ref_vehicle_status_code" json:"ref_vehicle_status_code"`
	RefVehicleStatusName string    `gorm:"column:ref_vehicle_status_name" json:"ref_vehicle_status_name"`
	Age                  string    `json:"age"`
}

// VmsMasCarpoolVehicle
type VmsMasCarpoolVehicle struct {
	MasCarpoolVehicleUID string    `gorm:"column:mas_carpool_vehicle_uid;primaryKey" json:"-"`
	MasCarpoolUID        string    `gorm:"column:mas_carpool_uid" json:"mas_carpool_uid" example:"164632c9-1d33-477e-b335-97a4e79a5845"`
	MasVehicleUID        string    `gorm:"column:mas_vehicle_uid" json:"mas_vehicle_uid" example:"334632c9-1d33-477e-b335-97a4e79a5845"`
	IsActive             string    `gorm:"column:is_active" json:"-"`
	IsDeleted            string    `gorm:"column:is_deleted" json:"-"`
	CreatedAt            time.Time `gorm:"column:created_at;autoCreateTime" json:"-"`
	UpdatedAt            time.Time `gorm:"column:updated_at;autoUpdateTime" json:"-"`
	CreatedBy            string    `gorm:"column:created_by" json:"-"`
	UpdatedBy            string    `gorm:"column:updated_by" json:"-"`

	StartDate               time.Time `gorm:"column:start_date" json:"-"`
	EndDate                 time.Time `gorm:"column:end_date" json:"-"`
	MasVehicleDepartmentUID string    `gorm:"column:mas_vehicle_department_uid" json:"-"`
}

func (VmsMasCarpoolVehicle) TableName() string {
	return "vms_mas_carpool_vehicle"
}

// VmsMasCarpoolVehicleActive
type VmsMasCarpoolVehicleActive struct {
	MasCarpoolVehicleUID string    `gorm:"column:mas_carpool_vehicle_uid;primaryKey" json:"mas_carpool_vehicle_uid" example:"164632c9-1d33-477e-b335-97a4e79a5845"`
	IsActive             string    `gorm:"column:is_active" json:"is_active" example:"1"`
	UpdatedAt            time.Time `gorm:"column:updated_at;autoUpdateTime" json:"-"`
	UpdatedBy            string    `gorm:"column:updated_by" json:"-"`
}

func (VmsMasCarpoolVehicleActive) TableName() string {
	return "vms_mas_carpool_vehicle"
}

// VmsMasCarpoolDriverList
type VmsMasCarpoolDriverList struct {
	MasCarpoolDriverUID string `gorm:"column:mas_carpool_driver_uid;primaryKey" json:"mas_carpool_driver_uid"`
	MasCarpoolUID       string `gorm:"column:mas_carpool_uid" json:"mas_carpool_uid" example:"164632c9-1d33-477e-b335-97a4e79a5845"`
	VmsMasCarpoolDriverDetail
}

func (VmsMasCarpoolDriverList) TableName() string {
	return "vms_mas_carpool_driver"
}

// VmsMasCarpoolDriverDetail
type VmsMasCarpoolDriverDetail struct {
	MasDriverUID                   string    `gorm:"column:mas_driver_uid" json:"mas_driver_uid" example:"334632c9-1d33-477e-b335-97a4e79a5845"`
	DriverImage                    string    `gorm:"column:driver_image" json:"driver_image"`
	DriverName                     string    `gorm:"column:driver_name" json:"driver_name"`
	DriverNickname                 string    `gorm:"column:driver_nickname" json:"driver_nickname"`
	DriverDeptSapShortNameHire     string    `gorm:"column:driver_dept_sap_short_name_hire" json:"driver_dept_sap_short_name_hire"`
	DriverContactNumber            string    `gorm:"column:driver_contact_number" json:"driver_contact_number"`
	DriverLicenseEndDate           time.Time `gorm:"column:driver_license_end_date" json:"driver_license_end_date"`
	ApprovedJobDriverEndDate       time.Time `gorm:"column:approved_job_driver_end_date" json:"approved_job_driver_end_date"`
	DriverAverageSatisfactionScore float64   `gorm:"column:driver_average_satisfaction_score" json:"driver_average_satisfaction_score"`
	RefDriverStatusCode            string    `gorm:"column:ref_driver_status_code" json:"ref_driver_status_code"`
	RefDriverStatusName            string    `gorm:"column:driver_status_name" json:"driver_status_name"`
	IsActive                       string    `gorm:"column:is_active" json:"-"`
}

// VmsMasCarpoolDriver
type VmsMasCarpoolDriver struct {
	MasCarpoolDriverUID string    `gorm:"column:mas_carpool_driver_uid;primaryKey" json:"-"`
	MasCarpoolUID       string    `gorm:"column:mas_carpool_uid" json:"mas_carpool_uid" example:"164632c9-1d33-477e-b335-97a4e79a5845"`
	MasDriverUID        string    `gorm:"column:mas_driver_uid" json:"mas_driver_uid" example:"334632c9-1d33-477e-b335-97a4e79a5845"`
	StartDate           time.Time `gorm:"column:start_date" json:"-"`
	EndDate             time.Time `gorm:"column:end_date" json:"-"`
	IsActive            string    `gorm:"column:is_active" json:"-"`
	IsDeleted           string    `gorm:"column:is_deleted" json:"-"`
	CreatedAt           time.Time `gorm:"column:created_at;autoCreateTime" json:"-"`
	UpdatedAt           time.Time `gorm:"column:updated_at;autoUpdateTime" json:"-"`
	CreatedBy           string    `gorm:"column:created_by" json:"-"`
	UpdatedBy           string    `gorm:"column:updated_by" json:"-"`
}

func (VmsMasCarpoolDriver) TableName() string {
	return "vms_mas_carpool_driver"
}

// VmsMasCarpoolVehicleActive
type VmsMasCarpoolDriverActive struct {
	MasCarpoolDriverUID string    `gorm:"column:mas_carpool_driver_uid;primaryKey" json:"mas_carpool_driver_uid" example:"164632c9-1d33-477e-b335-97a4e79a5845"`
	IsActive            string    `gorm:"column:is_active" json:"is_active" example:"1"`
	UpdatedAt           time.Time `gorm:"column:updated_at;autoUpdateTime" json:"-"`
	UpdatedBy           string    `gorm:"column:updated_by" json:"-"`
}

func (VmsMasCarpoolDriverActive) TableName() string {
	return "vms_mas_carpool_driver"
}

// VmsMasCarpoolAuthorizedDept
type VmsMasCarpoolAuthorizedDept struct {
	MasCarpoolAuthorizedDeptUID string    `gorm:"column:mas_carpool_authorized_dept_uid;primaryKey" json:"-"`
	MasCarpoolUID               string    `gorm:"column:mas_carpool_uid" json:"-"`
	DeptSap                     string    `gorm:"column:dept_sap" json:"dept_sap" example:"10001"`
	CreatedAt                   time.Time `gorm:"column:created_at;autoCreateTime" json:"-"`
	CreatedBy                   string    `gorm:"column:created_by" json:"-"`
	IsDeleted                   string    `gorm:"column:is_deleted" json:"-"`
}

func (VmsMasCarpoolAuthorizedDept) TableName() string {
	return "vms_mas_carpool_authorized_dept"
}

type VmsMasCarpoolAuthorizedDeptResponse struct {
	MasCarpoolAuthorizedDeptUID string           `gorm:"column:mas_carpool_authorized_dept_uid;primaryKey" json:"-"`
	MasCarpoolUID               string           `gorm:"column:mas_carpool_uid" json:"-"`
	DeptSap                     string           `gorm:"column:dept_sap" json:"dept_sap" example:"10001"`
	CreatedAt                   time.Time        `gorm:"column:created_at;autoCreateTime" json:"-"`
	CreatedBy                   string           `gorm:"column:created_by" json:"-"`
	IsDeleted                   string           `gorm:"column:is_deleted" json:"-"`
	MasDepartment               VmsMasDepartment `gorm:"foreignKey:DeptSap;references:DeptSAP" json:"mas_department"`
}

func (VmsMasCarpoolAuthorizedDeptResponse) TableName() string {
	return "vms_mas_carpool_authorized_dept"
}
