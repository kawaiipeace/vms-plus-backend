package models

import (
	"fmt"
	"time"
)

// VmsMasCarpoolList
type VmsMasCarpoolList struct {
	MasCarpoolUID            string `gorm:"column:mas_carpool_uid" json:"mas_carpool_uid"`
	CarpoolName              string `gorm:"column:carpool_name" json:"carpool_name"`
	CarpoolDeptSap           string `gorm:"column:carpool_dept_sap" json:"carpool_dept_sap"`
	CarpoolType              string `gorm:"column:carpool_type" json:"carpool_type"`
	CarpoolTypeName          string `gorm:"-" json:"carpool_type_name"`
	CarpoolAuthorizedDepts   string `gorm:"carpool_authorized_depts" json:"carpool_authorized_depts"`
	CarpoolAdminEmpName      string `gorm:"carpool_admin_emp_name" json:"carpool_admin_emp_name"`
	CarpoolAdminDeptSapShort string `gorm:"carpool_admin_dept_sap_short" json:"carpool_admin_dept_sap_short"`
	AdminDeptSapShort        string `gorm:"column:admin_dept_sap_short" json:"admin_dept_sap_short"`
	AdminPosition            string `gorm:"column:admin_position" json:"admin_position"`
	CarpoolContactPlace      string `gorm:"column:carpool_contact_place" json:"carpool_contact_place"`
	CarpoolContactNumber     string `gorm:"column:carpool_contact_number" json:"carpool_contact_number" `
	NumberOfDrivers          int    `gorm:"number_of_drivers" json:"number_of_drivers" `
	NumberOfVehicles         int    `gorm:"number_of_vehicles" json:"number_of_vehicles" `
	NumberOfApprovers        int    `gorm:"numberOfApprovers" json:"numberOfApprovers" `
	IsActive                 string `gorm:"column:is_active" json:"is_active"`
	CarpoolStatus            string `gorm:"-" json:"carpool_status"`
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

	CarPoolAdmins    []VmsMasCarpoolAdminCreate    `gorm:"-" json:"carpool_admins"`
	CarPoolApprovers []VmsMasCarpoolApproverCreate `gorm:"-" json:"carpool_approvers"`
	CarPoolVehicles  []VmsMasCarpoolVehicleCreate  `gorm:"-" json:"carpool_vehicles"`
	CarPoolDrivers   []VmsMasCarpoolDriverCreate   `gorm:"-" json:"carpool_drivers"`
}

func (VmsMasCarpoolRequest) TableName() string {
	return "vms_mas_carpool"
}

// VmsMasCarpoolUpdate
type VmsMasCarpoolUpdate struct {
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
	CreatedAt                time.Time                     `gorm:"column:created_at;autoCreateTime" json:"-"`
	UpdatedAt                time.Time                     `gorm:"column:updated_at;autoUpdateTime" json:"-"`
	CreatedBy                string                        `gorm:"column:created_by" json:"-"`
	UpdatedBy                string                        `gorm:"column:updated_by" json:"-"`
}

func (VmsMasCarpoolUpdate) TableName() string {
	return "vms_mas_carpool"
}

// VmsMasCarpoolActive
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
	CarpoolType              string                                `gorm:"column:carpool_type" json:"carpool_type"`
	CarpoolTypeName          string                                `gorm:"-" json:"carpool_type_name"`
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
	AdminPosition         string `gorm:"column:admin_position" json:"admin_position"`
	AdminDeptSap          string `gorm:"column:admin_dept_sap" json:"admin_dept_sap"`
	AdminDeptSapShort     string `gorm:"column:admin_dept_sap_short" json:"admin_dept_sap_short"`
	InternalContactNumber string `gorm:"column:internal_contact_number" json:"internal_contact_number" example:"1234"`
	MobileContactNumber   string `gorm:"column:mobile_contact_number" json:"mobile_contact_number" example:"9876543210"`
	IsMainAdmin           string `gorm:"column:is_main_admin" json:"is_main_admin"`
	IsActive              string `gorm:"column:is_active" json:"is_active"`
	ImageUrl              string `gorm:"-" json:"image_url"`
}

// TableName sets the table name for the VmsMasCarpoolAdmin model
func (VmsMasCarpoolAdminList) TableName() string {
	return "vms_mas_carpool_admin"
}

// VmsMasCarpoolAdminCreate
type VmsMasCarpoolAdminCreate struct {
	MasCarpoolAdminUID    string    `gorm:"column:mas_carpool_admin_uid;primaryKey" json:"-"`
	MasCarpoolUID         string    `gorm:"column:mas_carpool_uid" json:"-"`
	AdminEmpNo            string    `gorm:"column:admin_emp_no" json:"admin_emp_no" example:"990003"`
	AdminEmpName          string    `gorm:"column:admin_emp_name" json:"-"`
	AdminDeptSap          string    `gorm:"column:admin_dept_sap" json:"-"`
	AdminPosition         string    `gorm:"column:admin_position" json:"-"`
	InternalContactNumber string    `gorm:"column:internal_contact_number" json:"internal_contact_number" example:"1234"`
	MobileContactNumber   string    `gorm:"column:mobile_contact_number" json:"mobile_contact_number" example:"9876543210"`
	IsMainAdmin           string    `gorm:"column:is_main_admin" json:"is_main_admin" example:"0"`
	IsActive              string    `gorm:"column:is_active" json:"-"`
	IsDeleted             string    `gorm:"column:is_deleted" json:"-"`
	CreatedAt             time.Time `gorm:"column:created_at;autoCreateTime" json:"-"`
	UpdatedAt             time.Time `gorm:"column:updated_at;autoUpdateTime" json:"-"`
	CreatedBy             string    `gorm:"column:created_by" json:"-"`
	UpdatedBy             string    `gorm:"column:updated_by" json:"-"`
}

func (VmsMasCarpoolAdminCreate) TableName() string {
	return "vms_mas_carpool_admin"
}

// VmsMasCarpoolAdmin
type VmsMasCarpoolAdmin struct {
	MasCarpoolAdminUID    string    `gorm:"column:mas_carpool_admin_uid;primaryKey" json:"-"`
	MasCarpoolUID         string    `gorm:"column:mas_carpool_uid" json:"mas_carpool_uid" example:"164632c9-1d33-477e-b335-97a4e79a5845"`
	AdminEmpNo            string    `gorm:"column:admin_emp_no" json:"admin_emp_no" example:"990003"`
	AdminEmpName          string    `gorm:"column:admin_emp_name" json:"-"`
	AdminDeptSap          string    `gorm:"column:admin_dept_sap" json:"-"`
	AdminPosition         string    `gorm:"column:admin_position" json:"-"`
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
	ApproverPosition      string `gorm:"column:approver_position" json:"approver_position"`
	ApproverDeptSap       string `gorm:"column:approver_dept_sap" json:"approver_dept_sap"`
	ApproverDeptSapShort  string `gorm:"column:approver_dept_sap_short" json:"approver_dept_sap_short"`
	InternalContactNumber string `gorm:"column:internal_contact_number" json:"internal_contact_number" example:"5678"`
	MobileContactNumber   string `gorm:"column:mobile_contact_number" json:"mobile_contact_number" example:"9876543211"`
	IsMainApprover        string `gorm:"column:is_main_approver" json:"is_main_approver"`
	IsActive              string `gorm:"column:is_active" json:"is_active"`
	ImageUrl              string `gorm:"-" json:"image_url"`
}

func (VmsMasCarpoolApproverList) TableName() string {
	return "vms_mas_carpool_approver"
}

// VmsMasCarpoolApproverCreate
type VmsMasCarpoolApproverCreate struct {
	MasCarpoolApproverUID string    `gorm:"column:mas_carpool_approver_uid;primaryKey" json:"-"`
	MasCarpoolUID         string    `gorm:"column:mas_carpool_uid" json:"-"`
	ApproverEmpNo         string    `gorm:"column:approver_emp_no" json:"approver_emp_no" example:"990004"`
	ApproverEmpName       string    `gorm:"column:approver_emp_name" json:"-"`
	ApproverPosition      string    `gorm:"column:approver_position" json:"-"`
	ApproverDeptSap       string    `gorm:"column:approver_dept_sap" json:"-"`
	InternalContactNumber string    `gorm:"column:internal_contact_number" json:"internal_contact_number" example:"5678"`
	MobileContactNumber   string    `gorm:"column:mobile_contact_number" json:"mobile_contact_number" example:"9876543211"`
	IsMainApprover        string    `gorm:"column:is_main_approver" json:"is_main_approver" example:"0"`
	IsActive              string    `gorm:"column:is_active" json:"-"`
	IsDeleted             string    `gorm:"column:is_deleted" json:"-"`
	CreatedAt             time.Time `gorm:"column:created_at;autoCreateTime" json:"-"`
	UpdatedAt             time.Time `gorm:"column:updated_at;autoUpdateTime" json:"-"`
	CreatedBy             string    `gorm:"column:created_by" json:"-"`
	UpdatedBy             string    `gorm:"column:updated_by" json:"-"`
}

func (VmsMasCarpoolApproverCreate) TableName() string {
	return "vms_mas_carpool_approver"
}

// VmsMasCarpoolApprover
type VmsMasCarpoolApprover struct {
	MasCarpoolApproverUID string    `gorm:"column:mas_carpool_approver_uid;primaryKey" json:"-"`
	MasCarpoolUID         string    `gorm:"column:mas_carpool_uid" json:"mas_carpool_uid" example:"164632c9-1d33-477e-b335-97a4e79a5845"`
	ApproverEmpNo         string    `gorm:"column:approver_emp_no" json:"approver_emp_no" example:"990004"`
	ApproverEmpName       string    `gorm:"column:approver_emp_name" json:"-"`
	ApproverPosition      string    `gorm:"column:approver_position" json:"-"`
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
	MasVehicleUID                    string       `gorm:"column:mas_vehicle_uid;primary" json:"mas_vehicle_uid" example:"334632c9-1d33-477e-b335-97a4e79a5845"`
	IsActive                         string       `gorm:"column:is_active" json:"is_active"`
	VehicleLicensePlate              string       `gorm:"column:vehicle_license_plate" json:"vehicle_license_plate"`
	VehicleLicensePlateProvinceShort string       `gorm:"column:vehicle_license_plate_province_short" json:"vehicle_license_plate_province_short"`
	VehicleBrandName                 string       `gorm:"column:vehicle_brand_name" json:"vehicle_brand_name"`
	VehicleModelName                 string       `gorm:"column:vehicle_model_name" json:"vehicle_model_name"`
	RefVehicleTypeCode               string       `gorm:"column:ref_vehicle_type_code" json:"ref_vehicle_type_code"`
	RefVehicleTypeName               string       `gorm:"column:ref_vehicle_type_name" json:"ref_vehicle_type_name"`
	VehicleOwnerDeptSAP              string       `gorm:"column:vehicle_owner_dept_short" json:"vehicle_owner_dept_short"`
	FleetCardNo                      string       `gorm:"column:fleet_card_no" json:"fleet_card_no"`
	IsTaxCredit                      bool         `gorm:"column:is_tax_credit" json:"is_tax_credit"`
	VehicleMileage                   float64      `gorm:"column:vehicle_mileage" json:"vehicle_mileage"`
	VehicleRegistrationDate          TimeWithZone `gorm:"column:vehicle_registration_date" json:"vehicle_registration_date"` // Changed to time.Time
	RefVehicleStatusCode             string       `gorm:"column:ref_vehicle_status_code" json:"ref_vehicle_status_code"`
	RefVehicleStatusName             string       `gorm:"column:ref_vehicle_status_name" json:"ref_vehicle_status_name"`
	Age                              string       `json:"age"`
	Seat                             int          `gorm:"column:seat" json:"seat"`
	VehicleColor                     string       `gorm:"column:vehicle_color" json:"vehicle_color"`
	VehicleGear                      string       `gorm:"column:vehicle_gear" json:"vehicle_gear"`
	RefFuelTypeID                    int          `gorm:"column:ref_fuel_type_id" json:"ref_fuel_type_id"`
	FuelTypeName                     string       `gorm:"column:fuel_type_name" json:"fuel_type_name"`
	VehiclePeaID                     string       `gorm:"column:vehicle_pea_id" json:"vehicle_pea_id"`
	ParkingPlace                     string       `gorm:"column:parking_place" json:"parking_place"`
	VehicleImgs                      []string     `gorm:"-" json:"vehicle_imgs"`
}

// VmsMasCarpoolVehicleCreate
type VmsMasCarpoolVehicleCreate struct {
	MasCarpoolVehicleUID string    `gorm:"column:mas_carpool_vehicle_uid;primaryKey" json:"-"`
	MasCarpoolUID        string    `gorm:"column:mas_carpool_uid" json:"-"`
	MasVehicleUID        string    `gorm:"column:mas_vehicle_uid" json:"mas_vehicle_uid" example:"770ea678-c586-4d6b-9df2-c7d756bd4cc4"`
	IsActive             string    `gorm:"column:is_active" json:"-"`
	IsDeleted            string    `gorm:"column:is_deleted" json:"-"`
	CreatedAt            time.Time `gorm:"column:created_at;autoCreateTime" json:"-"`
	UpdatedAt            time.Time `gorm:"column:updated_at;autoUpdateTime" json:"-"`
	CreatedBy            string    `gorm:"column:created_by" json:"-"`
	UpdatedBy            string    `gorm:"column:updated_by" json:"-"`

	StartDate               TimeWithZone `gorm:"column:start_date" json:"-"`
	EndDate                 TimeWithZone `gorm:"column:end_date" json:"-"`
	MasVehicleDepartmentUID string       `gorm:"column:mas_vehicle_department_uid" json:"-"`
}

func (VmsMasCarpoolVehicleCreate) TableName() string {
	return "vms_mas_carpool_vehicle"
}

// VmsMasCarpoolVehicle
type VmsMasCarpoolVehicle struct {
	MasCarpoolVehicleUID string    `gorm:"column:mas_carpool_vehicle_uid;primaryKey" json:"-"`
	MasCarpoolUID        string    `gorm:"column:mas_carpool_uid" json:"mas_carpool_uid" example:"164632c9-1d33-477e-b335-97a4e79a5845"`
	MasVehicleUID        string    `gorm:"column:mas_vehicle_uid" json:"mas_vehicle_uid" example:"770ea678-c586-4d6b-9df2-c7d756bd4cc4"`
	IsActive             string    `gorm:"column:is_active" json:"-"`
	IsDeleted            string    `gorm:"column:is_deleted" json:"-"`
	CreatedAt            time.Time `gorm:"column:created_at;autoCreateTime" json:"-"`
	UpdatedAt            time.Time `gorm:"column:updated_at;autoUpdateTime" json:"-"`
	CreatedBy            string    `gorm:"column:created_by" json:"-"`
	UpdatedBy            string    `gorm:"column:updated_by" json:"-"`

	StartDate               TimeWithZone `gorm:"column:start_date" json:"-"`
	EndDate                 TimeWithZone `gorm:"column:end_date" json:"-"`
	MasVehicleDepartmentUID string       `gorm:"column:mas_vehicle_department_uid" json:"-"`
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
	MasDriverUID                   string       `gorm:"column:mas_driver_uid" json:"mas_driver_uid" example:"334632c9-1d33-477e-b335-97a4e79a5845"`
	DriverImage                    string       `gorm:"column:driver_image" json:"driver_image"`
	DriverName                     string       `gorm:"column:driver_name" json:"driver_name"`
	DriverNickname                 string       `gorm:"column:driver_nickname" json:"driver_nickname"`
	DriverBirthdate                TimeWithZone `gorm:"column:driver_birthdate" json:"driver_birthdate"`
	Age                            string       `gorm:"-" json:"age"`
	DriverDeptSapHire              string       `gorm:"column:driver_dept_sap_hire" json:"driver_dept_sap_hire"`
	DriverDeptSapShortNameHire     string       `gorm:"column:driver_dept_sap_short_name_hire" json:"driver_dept_sap_short_name_hire"`
	DriverDeptSapWork              string       `gorm:"column:driver_dept_sap_work" json:"driver_dept_sap_work"`
	DriverDeptSapShortNameWork     string       `gorm:"column:driver_dept_sap_short_work" json:"driver_dept_sap_short_work"`
	DriverContactNumber            string       `gorm:"column:driver_contact_number" json:"driver_contact_number"`
	ApprovedJobDriverEndDate       TimeWithZone `gorm:"column:approved_job_driver_end_date" json:"approved_job_driver_end_date"`
	DriverAverageSatisfactionScore float64      `gorm:"column:driver_average_satisfaction_score" json:"driver_average_satisfaction_score"`
	DriverTotalSatisfactionReview  int          `gorm:"column:driver_total_satisfaction_review" json:"driver_total_satisfaction_review"`
	RefDriverStatusName            string       `gorm:"column:driver_status_name" json:"driver_status_name"`
	ContractNo                     string       `gorm:"column:contract_no" json:"contract_no"`
	EndDate                        TimeWithZone `gorm:"column:end_date" json:"end_date"`
	DriverLicenseNo                string       `gorm:"column:driver_license_no" json:"driver_license_no"`
	DriverLicenseEndDate           TimeWithZone `gorm:"column:driver_license_end_date" json:"driver_license_end_date"`
	IsActive                       string       `gorm:"column:is_active" json:"is_active"`
	VendorName                     string       `gorm:"column:vendor_name" json:"vendor_name"`
}

func (d *VmsMasCarpoolDriverDetail) CalculateAgeInYearsMonths() string {
	if d.DriverBirthdate.IsZero() {
		return "ไม่ระบุ"
	}

	today := time.Now()
	years := today.Year() - d.DriverBirthdate.Year()
	months := today.Month() - d.DriverBirthdate.Month()

	// Adjust if birthday hasn't occurred yet this year
	if today.Day() < d.DriverBirthdate.Day() {
		months--
	}

	if months < 0 {
		years--
		months += 12
	}

	return fmt.Sprintf("%d ปี %d เดือน", years, months)
}

// VmsMasCarpoolDriverCreate
type VmsMasCarpoolDriverCreate struct {
	MasCarpoolDriverUID string       `gorm:"column:mas_carpool_driver_uid;primaryKey" json:"-"`
	MasCarpoolUID       string       `gorm:"column:mas_carpool_uid" json:"-"`
	MasDriverUID        string       `gorm:"column:mas_driver_uid" json:"mas_driver_uid" example:"9c63bdfa-1a80-4238-8a3a-cbc1891161a2"`
	StartDate           TimeWithZone `gorm:"column:start_date" json:"-"`
	EndDate             TimeWithZone `gorm:"column:end_date" json:"-"`
	IsActive            string       `gorm:"column:is_active" json:"-"`
	IsDeleted           string       `gorm:"column:is_deleted" json:"-"`
	CreatedAt           time.Time    `gorm:"column:created_at;autoCreateTime" json:"-"`
	UpdatedAt           time.Time    `gorm:"column:updated_at;autoUpdateTime" json:"-"`
	CreatedBy           string       `gorm:"column:created_by" json:"-"`
	UpdatedBy           string       `gorm:"column:updated_by" json:"-"`
}

func (VmsMasCarpoolDriverCreate) TableName() string {
	return "vms_mas_carpool_driver"
}

// VmsMasCarpoolDriver
type VmsMasCarpoolDriver struct {
	MasCarpoolDriverUID string       `gorm:"column:mas_carpool_driver_uid;primaryKey" json:"-"`
	MasCarpoolUID       string       `gorm:"column:mas_carpool_uid" json:"mas_carpool_uid" example:"164632c9-1d33-477e-b335-97a4e79a5845"`
	MasDriverUID        string       `gorm:"column:mas_driver_uid" json:"mas_driver_uid" example:"9c63bdfa-1a80-4238-8a3a-cbc1891161a2"`
	StartDate           TimeWithZone `gorm:"column:start_date" json:"-"`
	EndDate             TimeWithZone `gorm:"column:end_date" json:"-"`
	IsActive            string       `gorm:"column:is_active" json:"-"`
	IsDeleted           string       `gorm:"column:is_deleted" json:"-"`
	CreatedAt           time.Time    `gorm:"column:created_at;autoCreateTime" json:"-"`
	UpdatedAt           time.Time    `gorm:"column:updated_at;autoUpdateTime" json:"-"`
	CreatedBy           string       `gorm:"column:created_by" json:"-"`
	UpdatedBy           string       `gorm:"column:updated_by" json:"-"`
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
	UpdatedAt                   time.Time `gorm:"column:updated_at" json:"-"`
	UpdatedBy                   string    `gorm:"column:updated_by" json:"-"`
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
