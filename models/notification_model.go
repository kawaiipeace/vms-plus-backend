package models

import "time"

type Notification struct {
	TrnNotifyUID         string    `gorm:"column:trn_notify_uid;primaryKey" json:"trn_notify_uid"`
	NotifyType           string    `gorm:"column:notify_type;not null" json:"notify_type"`
	NotifyRole           string    `gorm:"column:notify_role;not null" json:"notify_role"`
	RecordUID            string    `gorm:"column:record_uid;not null" json:"record_uid"`
	EmpID                string    `gorm:"column:emp_id;not null" json:"emp_id"`
	Title                string    `gorm:"column:title;not null" json:"title"`
	Message              string    `gorm:"column:message;not null" json:"message"`
	IsRead               bool      `gorm:"column:is_read;default:false" json:"is_read"`
	CreatedAt            time.Time `gorm:"column:created_at;default:CURRENT_TIMESTAMP" json:"created_at"`
	ReadAt               time.Time `gorm:"column:read_at" json:"read_at"`
	Duration             string    `gorm:"-" json:"duration"`
	NotifyURL            string    `gorm:"-" json:"notify_url"`
	RefRequestStatusCode string    `gorm:"column:ref_request_status_code;not null" json:"ref_request_status_code"`
}

func (Notification) TableName() string {
	return "vms_trn_notifications"
}

type NotificationTemplate struct {
	MasTemplateUID       string    `gorm:"column:mas_template_uid;primaryKey" json:"mas_template_uid"`
	NotifyType           string    `gorm:"column:notify_type;not null" json:"notify_type"`
	RefRequestStatusCode string    `gorm:"column:ref_request_status_code;not null" json:"ref_request_status_code"`
	NotifyRole           string    `gorm:"column:notify_role;not null" json:"notify_role"`
	NotifyTitle          string    `gorm:"column:notify_title;not null" json:"notify_title"`
	NotifyMessage        string    `gorm:"column:notify_message;not null" json:"notify_message"`
	CreatedAt            time.Time `gorm:"column:created_at;default:CURRENT_TIMESTAMP" json:"created_at"`
}

func (NotificationTemplate) TableName() string {
	return "vms_mas_notification_template"
}

type RequestBookingNotification struct {
	TrnRequestUID         string `gorm:"column:trn_request_uid;primaryKey" json:"trn_request_uid"`
	RefRequestStatusCode  string `gorm:"column:ref_request_status_code" json:"-"`
	RequestNo             string `gorm:"column:request_no" json:"request_no" example:"123456"`
	CreatedRequestEmpID   string `gorm:"column:created_request_emp_id" json:"-"`
	VehicleUserEmpID      string `gorm:"column:vehicle_user_emp_id" json:"vehicle_user_emp_id" example:"990001"`
	DriverEmpID           string `gorm:"column:driver_emp_id" json:"driver_emp_id" example:"700001"`
	ConfirmedRequestEmpID string `gorm:"column:confirmed_request_emp_id" json:"confirmed_request_emp_id" example:"501621"`
	ApprovedRequestEmpID  string `gorm:"column:approved_request_emp_id" json:"approved_request_emp_id" example:"501621"`
}

func (RequestBookingNotification) TableName() string {
	return "vms_trn_request"
}

type RequestAnnualLicenseNotification struct {
	TrnRequestAnnualDriverUID        string `gorm:"column:trn_request_annual_driver_uid;primaryKey" json:"trn_request_annual_driver_uid"`
	RefRequestAnnualDriverStatusCode string `gorm:"column:ref_request_annual_driver_status_code" json:"-"`
	RequestAnnualDriverNo            string `gorm:"column:request_annual_driver_no" json:"-"`
	AnnualYYYY                       int    `gorm:"column:annual_yyyy" json:"annual_yyyy" example:"2568"`
	CreatedRequestEmpID              string `gorm:"column:created_request_emp_id" json:"-"`
	ConfirmedRequestEmpID            string `gorm:"column:confirmed_request_emp_id" json:"confirmed_request_emp_id" example:"501621"`
	ApprovedRequestEmpID             string `gorm:"column:approved_request_emp_id" json:"approved_request_emp_id" example:"501621"`
}

func (RequestAnnualLicenseNotification) TableName() string {
	return "vms_trn_request_annual_driver"
}
