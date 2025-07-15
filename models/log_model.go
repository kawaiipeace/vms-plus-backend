package models

import "time"

type LogCreate struct {
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	CreatedBy string    `gorm:"column:created_by;not null" json:"created_by" binding:"required"`
}
type LogUpdate struct {
	UpdatedAt time.Time `gorm:"column:updated_at;autoCreateTime" json:"updated_at"`
	UpdatedBy string    `gorm:"column:updated_by;not null" json:"updated_by" binding:"required"`
}

type LogRequest struct {
	LogRequestActionUID      string       `gorm:"primaryKey;column:log_request_action_uid" json:"log_request_action_uid"`
	TrnRequestUID            string       `gorm:"column:trn_request_uid;not null" json:"trn_request_uid"`
	RefRequestStatusCode     string       `gorm:"column:ref_request_status_code" json:"ref_request_status_code"`
	LogRequestActionDatetime TimeWithZone `gorm:"column:log_request_action_datetime;default:CURRENT_TIMESTAMP" json:"log_request_action_datetime"`
	ActionByPersonalID       string       `gorm:"column:action_by_personal_id" json:"action_by_personal_id"`
	ActionByFullname         string       `gorm:"column:action_by_fullname" json:"action_by_fullname"`
	ActionByRole             string       `gorm:"column:action_by_role" json:"action_by_role"`
	ActionByPosition         string       `gorm:"column:action_by_position" json:"action_by_position"`
	ActionByDepartment       string       `gorm:"column:action_by_department" json:"action_by_department"`
	ActionDetail             string       `gorm:"column:action_detail" json:"action_detail"`
	Remark                   string       `gorm:"column:remark" json:"remark"`
	RoleOfCreater            string       `gorm:"-" json:"role_of_creater"`
}

func (LogRequest) TableName() string {
	return "public.vms_log_request_action"
}

type RequestStatus struct {
	RefRequestStatusCode string `gorm:"primaryKey;column:ref_request_status_code" json:"ref_request_status_code"`
	RefRequestStatusDesc string `gorm:"column:ref_request_status_desc" json:"ref_request_status_desc"`
}

func (RequestStatus) TableName() string {
	return "public.vms_ref_request_status"
}

type VmsLogRequest struct {
	LogRequestActionUID      string       `gorm:"primaryKey;column:log_request_action_uid" json:"log_request_action_uid"`
	TrnRequestUID            string       `gorm:"column:trn_request_uid;not null" json:"trn_request_uid"`
	RefRequestStatusCode     string       `gorm:"column:ref_request_status_code" json:"ref_request_status_code"`
	LogRequestActionDatetime TimeWithZone `gorm:"column:log_request_action_datetime;default:CURRENT_TIMESTAMP" json:"log_request_action_datetime"`
	ActionByPersonalID       string       `gorm:"column:action_by_personal_id" json:"action_by_personal_id"`
	ActionByFullname         string       `gorm:"column:action_by_fullname" json:"action_by_fullname"`
	ActionByRole             string       `gorm:"column:action_by_role" json:"action_by_role"`
	ActionByPosition         string       `gorm:"column:action_by_position" json:"action_by_position"`
	ActionByDepartment       string       `gorm:"column:action_by_department" json:"action_by_department"`
	ActionDetail             string       `gorm:"column:action_detail" json:"action_detail"`
	Remark                   string       `gorm:"column:remark" json:"remark"`
	IsDeleted                string       `gorm:"column:is_deleted" json:"is_deleted"`
}

// TableName overrides the default table name
func (VmsLogRequest) TableName() string {
	return "vms_log_request_action"
}
