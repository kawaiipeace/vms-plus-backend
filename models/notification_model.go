package models

import "time"

type Notification struct {
	TrnNotifyUID string    `gorm:"column:trn_notify_uid;primaryKey" json:"trn_notify_uid"`
	EmpID        string    `gorm:"column:emp_id;not null" json:"emp_id"`
	Title        string    `gorm:"column:title;not null" json:"title"`
	Message      string    `gorm:"column:message;not null" json:"message"`
	IsRead       bool      `gorm:"column:is_read;default:false" json:"is_read"`
	CreatedAt    time.Time `gorm:"column:created_at;default:CURRENT_TIMESTAMP" json:"created_at"`
	ReadAt       time.Time `gorm:"column:read_at" json:"read_at"`
	Duration     string    `gorm:"-" json:"duration"`
}

func (Notification) TableName() string {
	return "vms_trn_notifications"
}
