package models

import "time"

type KeyCloak_Request struct {
	RedirectUri string `json:"redirect_uri" example:"http://localhost/callback_code_token" binding:"required"`
}
type KeyCloak_Authen struct {
	Code        string `json:"code" example:"b7d287b3-8ded-40f2-9060-a60f1207b378.6c99a927-ae7f-4d15-a69c-ea8aed9b612f.88e696e4-d20f-47f5-a391-26854ff28848" binding:"required"`
	RedirectUri string `json:"redirect_uri" example:"http://localhost/callback_code_token" binding:"required"`
}
type KeyCloak_Error_Response struct {
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
}
type KeyCloak_Success_Response struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
}
type KeyCloak_UserInfo struct {
	Username string `json:"preferred_username"`
	FullName string `json:"hr_fullname_th"`
	UserID   string `json:"sub"`
}

type ThaiID_Request struct {
	RedirectUri string `json:"redirect_uri" example:"http://localhost/callback_code_token" binding:"required"`
}
type ThaiID_Authen struct {
	Code        string `json:"code" example:"b7d287b3-8ded-40f2-9060-a60f1207b378.6c99a927-ae7f-4d15-a69c-ea8aed9b612f.88e696e4-d20f-47f5-a391-26854ff28848" binding:"required"`
	RedirectUri string `json:"redirect_uri" example:"http://localhost/callback_code_token" binding:"required"`
}

type ThaiID_Error_Response struct {
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
}
type ThaiID_Success_Response struct {
	AccessToken string `json:"access_token"`
	PID         string `json:"pid"`
}

type OTP_Request struct {
	Phone string `json:"phone" example:"0818088770" binding:"required"`
}

type OTPVerify_Request struct {
	OtpId string `json:"otpId" example:"35801ccf-43cc-4e21-ae23-fc40471d9789" binding:"required"`
	OTP   string `json:"otp" example:"000000" binding:"required"`
}

type Login_Response struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type RefreshToken_Request struct {
	RefreshToken string `json:"refreshToken" binding:"required"`
}

type OTP_Request_Create struct {
	ReqID      uint      `gorm:"primaryKey;column:req_id"`
	PhoneNo    string    `gorm:"column:phone_no;size:10;not null"`
	OTPID      string    `gorm:"column:otp_id;size:36;not null"`
	CreatedAt  time.Time `gorm:"column:created_at;autoCreateTime"`
	ExpiresAt  time.Time `gorm:"column:expires_at;not null"`
	Status     string    `gorm:"column:status;size:20;default:pending"`
	IsEmployee bool      `gorm:"column:is_employee"`
}

func (OTP_Request_Create) TableName() string {
	return "trn_otp_request"
}
