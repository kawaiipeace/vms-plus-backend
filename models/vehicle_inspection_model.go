package models

import (
	"time"
)

// VmsTrnReturnedVehicle
type VmsTrnReturnedVehicleNoImage struct {
	TrnRequestUID            string    `gorm:"column:trn_request_uid;primaryKey" json:"trn_request_uid" example:"8bd09808-61fa-42fd-8a03-bf961b5678cd"`
	ReturnedVehicleDatetime  time.Time `gorm:"column:returned_vehicle_datetime" json:"returned_vehicle_datetime" example:"2025-04-16T14:30:00Z"`
	MileEnd                  int       `gorm:"column:mile_end" json:"mile_end" example:"12000"`
	FuelEnd                  int       `gorm:"column:fuel_end" json:"fuel_end" example:"70"`
	ReturnedCleanlinessLevel int       `gorm:"column:returned_cleanliness_level" json:"returned_cleanliness_level" example:"1"`
	CommentOnReturnedVehicle string    `gorm:"column:comment_on_returned_vehicle" json:"comment_on_returned_vehicle" example:"OK"`
	UpdatedAt                time.Time `gorm:"column:updated_at" json:"-"`
	UpdatedBy                string    `gorm:"column:updated_by" json:"-"`
}

func (VmsTrnReturnedVehicleNoImage) TableName() string {
	return "public.vms_trn_request"
}

// VmsTrnReturnedVehicle
type VmsTrnReturnedVehicleImages struct {
	TrnRequestUID string                 `gorm:"column:trn_request_uid;primaryKey" json:"trn_request_uid" example:"8bd09808-61fa-42fd-8a03-bf961b5678cd"`
	VehicleImages []VehicleImageReturned `gorm:"foreignKey:TrnRequestUID;references:TrnRequestUID" json:"vehicle_images"`
	UpdatedAt     time.Time              `gorm:"column:updated_at" json:"-"`
	UpdatedBy     string                 `gorm:"column:updated_by" json:"-"`
}

func (VmsTrnReturnedVehicleImages) TableName() string {
	return "public.vms_trn_request"
}

// VmsTrnSatisfactionSurveyAnswersResponse
type VmsTrnSatisfactionSurveyAnswersResponse struct {
	TrnSatisfactionSurveyAnswersUID    string                            `gorm:"column:trn_satisfaction_survey_answers_uid;primaryKey" json:"-"`
	TrnRequestUID                      string                            `gorm:"column:trn_request_uid" json:"-"`
	MasSatisfactionSurveyQuestionsCode int                               `gorm:"column:mas_satisfaction_survey_questions_code" json:"mas_satisfaction_survey_questions_code" example:"1"`
	SurveyAnswer                       int                               `gorm:"column:survey_answer" json:"survey_answer" example:"5"`
	SatisfactionSurveyQuestions        VmsMasSatisfactionSurveyQuestions `gorm:"foreignKey:MasSatisfactionSurveyQuestionsCode;references:MasSatisfactionSurveyQuestionsCode" json:"satisfaction_survey_questions"`
}

func (VmsTrnSatisfactionSurveyAnswersResponse) TableName() string {
	return "public.vms_trn_satisfaction_survey_answers"
}

// VmsTrnRequestAccepted
type VmsTrnRequestAccepted struct {
	TrnRequestUID               string    `gorm:"column:trn_request_uid;primaryKey" json:"trn_request_uid" example:"8bd09808-61fa-42fd-8a03-bf961b5678cd"`
	AcceptedVehicleDatetime     time.Time `gorm:"column:accepted_vehicle_datetime" json:"accepted_vehicle_datetime" example:"2025-04-16T14:30:00Z"`
	AcceptedVehicleEmpID        string    `gorm:"column:accepted_vehicle_emp_id" json:"-"`
	AcceptedVehicleEmpName      string    `gorm:"column:accepted_vehicle_emp_name" json:"-"`
	AcceptedVehicleDeptSAP      string    `gorm:"column:accepted_vehicle_dept_sap" json:"-"`
	AcceptedVehicleDeptSAPShort string    `gorm:"column:accepted_vehicle_dept_sap_short" json:"-"`
	AcceptedVehicleDeptSAPFull  string    `gorm:"column:accepted_vehicle_dept_sap_full" json:"-"`
	RefRequestStatusCode        string    `gorm:"column:ref_request_status_code" json:"-"`
	UpdatedAt                   time.Time `gorm:"column:updated_at" json:"-"`
	UpdatedBy                   string    `gorm:"column:updated_by" json:"-"`
}

func (VmsTrnRequestAccepted) TableName() string {
	return "public.vms_trn_request"
}

// VmsTrnInsVehicle
type VmsTrnInspectVehicleImages struct {
	TrnRequestUID string                `gorm:"column:trn_request_uid;primaryKey" json:"trn_request_uid" example:"8bd09808-61fa-42fd-8a03-bf961b5678cd"`
	VehicleImages []VehicleImageInspect `gorm:"foreignKey:TrnRequestUID;references:TrnRequestUID" json:"vehicle_images"`
	UpdatedAt     time.Time             `gorm:"column:updated_at" json:"-"`
	UpdatedBy     string                `gorm:"column:updated_by" json:"-"`
}

func (VmsTrnInspectVehicleImages) TableName() string {
	return "public.vms_trn_request"
}

// VehicleImageReturned
type VehicleImageInspect struct {
	TrnVehicleImgReturnedUID string `gorm:"column:trn_vehicle_img_returned_uid;primaryKey" json:"-"`
	TrnRequestUID            string `gorm:"column:trn_request_uid;" json:"-"`
	RefVehicleImgSideCode    int    `gorm:"column:ref_vehicle_img_side_code" json:"ref_vehicle_img_side_code" example:"1"`
	VehicleImgFile           string `gorm:"column:vehicle_img_file" json:"vehicle_img_file" example:"http://vms.pea.co.th/side_image.jpg"`
}

func (VehicleImageInspect) TableName() string {
	return "public.vms_trn_vehicle_img_returned"
}
