package models

import (
	"time"
)

// VmsTrnReturnedVehicle
type VmsTrnReturnedVehicleNoImage struct {
	TrnRequestUID           string    `gorm:"column:trn_request_uid;primaryKey" json:"trn_request_uid" example:"0b07440c-ab04-49d0-8730-d62ce0a9bab9"`
	ReturnedVehicleDatetime time.Time `gorm:"column:returned_vehicle_datetime" json:"returned_vehicle_datetime" example:"2025-04-16T14:30:00Z"`
	MileEnd                 int       `gorm:"column:mile_end" json:"mile_end" example:"12000"`
	FuelEnd                 int       `gorm:"column:fuel_end" json:"fuel_end" example:"70"`
	ReceivedVehicleRemark   string    `gorm:"column:received_vehicle_remark" json:"received_vehicle_remark" example:"Minor scratch on bumper"`
	ReturnedVehicleRemark   string    `gorm:"column:returned_vehicle_remark" json:"returned_vehicle_remark" example:"OK"`
	UpdatedAt               time.Time `gorm:"column:updated_at" json:"-"`
	UpdatedBy               string    `gorm:"column:updated_by" json:"-"`
}

func (VmsTrnReturnedVehicleNoImage) TableName() string {
	return "public.vms_trn_request"
}

// VmsTrnReturnedVehicle
type VmsTrnReturnedVehicleImages struct {
	TrnRequestUID string                 `gorm:"column:trn_request_uid;primaryKey" json:"trn_request_uid" example:"0b07440c-ab04-49d0-8730-d62ce0a9bab9"`
	VehicleImages []VehicleImageReturned `gorm:"foreignKey:TrnRequestUID;references:TrnRequestUID" json:"vehicle_images"`
	UpdatedAt     time.Time              `gorm:"column:updated_at" json:"-"`
	UpdatedBy     string                 `gorm:"column:updated_by" json:"-"`
}

func (VmsTrnReturnedVehicleImages) TableName() string {
	return "public.vms_trn_request"
}

// VmsTrnSatisfactionSurveyAnswersResponse
type VmsTrnSatisfactionSurveyAnswersResponse struct {
	TrnSatisfactionSurveyAnswersUID   string                            `gorm:"column:trn_satisfaction_survey_answers_uid;primaryKey" json:"-"`
	TrnRequestUID                     string                            `gorm:"column:trn_request_uid" json:"-"`
	MasSatisfactionSurveyQuestionsUID string                            `gorm:"column:mas_satisfaction_survey_questions_uid" json:"mas_satisfaction_survey_questions_uid" example:"1"`
	SurveyAnswer                      int                               `gorm:"column:survey_answer_score" json:"survey_answer" example:"5"`
	SatisfactionSurveyQuestions       VmsMasSatisfactionSurveyQuestions `gorm:"foreignKey:MasSatisfactionSurveyQuestionsUID;references:MasSatisfactionSurveyQuestionsUID" json:"satisfaction_survey_questions"`
}

func (VmsTrnSatisfactionSurveyAnswersResponse) TableName() string {
	return "public.vms_trn_satisfaction_survey_answers"
}

// VmsTrnRequestAccepted
type VmsTrnRequestAccepted struct {
	TrnRequestUID              string    `gorm:"column:trn_request_uid;primaryKey" json:"trn_request_uid" example:"0b07440c-ab04-49d0-8730-d62ce0a9bab9"`
	InspectVehicleDatetime     time.Time `gorm:"column:inspected_vehicle_datetime" json:"accepted_vehicle_datetime" example:"2025-04-16T14:30:00Z"`
	InspectVehicleEmpID        string    `gorm:"column:inspected_vehicle_emp_id" json:"-"`
	InspectVehicleEmpName      string    `gorm:"column:inspected_vehicle_emp_name" json:"-"`
	InspectVehicleDeptSAP      string    `gorm:"column:inspected_vehicle_dept_sap" json:"-"`
	InspectVehicleDeptSAPShort string    `gorm:"column:inspected_vehicle_dept_name_short" json:"-"`
	InspectVehicleDeptSAPFull  string    `gorm:"column:inspected_vehicle_dept_name_full" json:"-"`
	RefRequestStatusCode       string    `gorm:"column:ref_request_status_code" json:"-"`
	UpdatedAt                  time.Time `gorm:"column:updated_at" json:"-"`
	UpdatedBy                  string    `gorm:"column:updated_by" json:"-"`
}

func (VmsTrnRequestAccepted) TableName() string {
	return "public.vms_trn_request"
}

// VmsTrnInsVehicle
type VmsTrnInspectVehicleImages struct {
	TrnRequestUID string                `gorm:"column:trn_request_uid;primaryKey" json:"trn_request_uid" example:"0b07440c-ab04-49d0-8730-d62ce0a9bab9"`
	VehicleImages []VehicleImageInspect `gorm:"foreignKey:TrnRequestUID;references:TrnRequestUID" json:"vehicle_images"`
	UpdatedAt     time.Time             `gorm:"column:updated_at" json:"-"`
	UpdatedBy     string                `gorm:"column:updated_by" json:"-"`
}

func (VmsTrnInspectVehicleImages) TableName() string {
	return "public.vms_trn_request"
}

// VehicleImageInspect
type VehicleImageInspect struct {
	TrnVehicleImgReturnedUID string    `gorm:"column:trn_vehicle_img_inspect_uid;primaryKey" json:"-"`
	TrnRequestUID            string    `gorm:"column:trn_request_uid;" json:"-"`
	RefVehicleImgSideCode    int       `gorm:"column:ref_vehicle_img_side_code" json:"ref_vehicle_img_side_code" example:"1"`
	VehicleImgFile           string    `gorm:"column:vehicle_img_file" json:"vehicle_img_file" example:"http://vms.pea.co.th/side_image.jpg"`
	CreatedAt                time.Time `gorm:"column:created_at" json:"-"`
	CreatedBy                string    `gorm:"column:created_by" json:"-"`
	UpdatedAt                time.Time `gorm:"column:updated_at" json:"-"`
	UpdatedBy                string    `gorm:"column:updated_by" json:"-"`
	IsDeleted                string    `gorm:"column:is_deleted" json:"-"`
}

func (VehicleImageInspect) TableName() string {
	return "public.vms_trn_vehicle_img_inspect"

}
