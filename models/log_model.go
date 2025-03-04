package models

import (
	"time"
)

type LogCreate struct {
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	CreatedBy string    `gorm:"column:created_by;not null" json:"created_by" binding:"required"`
}
type LogUpdate struct {
	UpdatedAt time.Time `gorm:"column:updated_at;autoCreateTime" json:"updated_at"`
	UpdatedBy string    `gorm:"column:updated_by;not null" json:"updated_by" binding:"required"`
}

type LogRequest struct {
	LogRequestUID string        `gorm:"primaryKey" json:"log_request_uid"`
	TrnRequestUID string        `gorm:"column:trn_request_uid" json:"trn_request_uid"`
	RefStatusCode string        `gorm:"column:ref_status_code" json:"ref_status_code"`
	LogRemark     string        `gorm:"column:log_remark" json:"log_remark"`
	CreatedAt     string        `gorm:"autoCreateTime" json:"created_at"`
	CreatedBy     string        `gorm:"column:created_by" json:"created_by"`
	CreatedByEmp  EmpUsr        `gorm:"foreignKey:CreatedBy;references:EmpNo" json:"created_by_emp"`
	Status        RequestStatus `gorm:"foreignKey:RefStatusCode;references:RefRequestStatusCode" json:"status"`
}

func (LogRequest) TableName() string {
	return "public.vms_log_request"
}

type EmpUsr struct {
	EmpNo   string `gorm:"primaryKey;column:emp_no" json:"emp_no"`
	EmpName string `gorm:"column:emp_name" json:"emp_name"`
	DeptSAP string `gorm:"column:dept_sap" json:"dept_sap"`
}

func (EmpUsr) TableName() string {
	return "vms_user.mas_emp"
}

type RequestStatus struct {
	RefRequestStatusCode string `gorm:"primaryKey;column:ref_request_status_code" json:"ref_request_status_code"`
	RefRequestStatusDesc string `gorm:"column:ref_request_status_desc" json:"ref_request_status_desc"`
}

func (RequestStatus) TableName() string {
	return "public.vms_ref_request_status"
}

type VmsLogRequest struct {
	LogRequestUID string    `gorm:"column:log_request_uid;primaryKey" json:"log_request_uid"`
	TrnRequestUID string    `gorm:"column:trn_request_uid;not null" json:"trn_request_uid" binding:"required"`
	RefStatusCode string    `gorm:"column:ref_status_code" json:"ref_status_code"`
	LogRemark     string    `gorm:"column:log_remark" json:"log_remark"`
	CreatedAt     time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	CreatedBy     string    `gorm:"column:created_by;not null" json:"created_by" binding:"required"`
}

// TableName overrides the default table name
func (VmsLogRequest) TableName() string {
	return "vms_log_request"
}
