package handlers

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
	"vms_plus_be/config"
	"vms_plus_be/funcs"
	"vms_plus_be/models"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

type LoginHandler struct {
}

// RequestKeyCloak godoc
// @Summary Request Keycloak authentication token
// @Description This endpoint allows a user to request an authentication token from Keycloak for secure access.
// @Tags Login
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param data body models.KeyCloak_Request true "KeyCloak_Request data"
// @Router /api/login/request-keycloak [post]
func (h *LoginHandler) RequestKeyCloak(c *gin.Context) {
	var req models.KeyCloak_Request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Invalid JSON input"})
		return
	}

	c.JSON(200, gin.H{
		"url": config.AppConfig.KeyCloakEndPoint + "/auth?response_type=code&client_id=" + config.AppConfig.KeyCloakClientID + "&redirect_uri=" + req.RedirectUri + "&scope=openid&state=001",
	})
}

// AuthenKeyCloak godoc
// @Summary Authenticate user via Keycloak
// @Description This endpoint authenticates a user using Keycloak, providing secure access to protected resources.
// @Tags Login
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param data body models.KeyCloak_Authen true "KeyCloak_Authen data"
// @Router /api/login/authen-keycloak [post]
func (h *LoginHandler) AuthenKeyCloak(c *gin.Context) {
	var req models.KeyCloak_Authen
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Invalid JSON input"})
		return
	}
	endpoint := config.AppConfig.KeyCloakEndPoint + "/token"
	// Define form data
	data := url.Values{}
	data.Set("grant_type", "authorization_code")
	data.Set("code", req.Code)
	data.Set("redirect_uri", req.RedirectUri)
	data.Set("client_id", config.AppConfig.KeyCloakClientID)
	data.Set("client_secret", config.AppConfig.KeyCloakSecret)

	// Create request
	creq, err := http.NewRequest("POST", endpoint, strings.NewReader(data.Encode()))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Set headers
	creq.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Execute request
	client := &http.Client{}
	resp, err := client.Do(creq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer resp.Body.Close()

	// Read response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var errorResponse models.KeyCloak_Error_Response
	if err := json.Unmarshal(body, &errorResponse); err == nil {
		if errorResponse.Error != "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": errorResponse.ErrorDescription})
			return
		}
	}

	var successResponse models.KeyCloak_Success_Response
	if err := json.Unmarshal(body, &successResponse); err == nil && successResponse.AccessToken != "" {
		userInfo, err := GetUserInfo(successResponse.AccessToken)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		} else {
			var user models.AuthenUserEmp
			err = config.DB.Where("emp_id = ?", userInfo.Username).
				First(&user).Error

			if err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
				} else {
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				}
				return
			}
			newAccessToken, err := funcs.GenerateJWT(user, "access", time.Duration(config.AppConfig.JwtAccessTokenTime)*time.Minute, successResponse.AccessToken, successResponse.RefreshToken)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Error generating access token"})
				return
			}

			newRefreshToken, err := funcs.GenerateJWT(user, "refresh", time.Duration(config.AppConfig.JwtRefreshTokenTime)*time.Hour, successResponse.AccessToken, successResponse.RefreshToken)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Error generating refresh token"})
				return
			}

			c.JSON(http.StatusOK, models.Login_Response{
				AccessToken:  newAccessToken,
				RefreshToken: newRefreshToken,
			})
			return
		}
	}
	c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
}

// RequestThaiID godoc
// @Summary Request ThaiID authentication token
// @Description This endpoint allows a user to request an authentication token from ThaiID for secure access.
// @Tags Login
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param data body models.ThaiID_Request true "ThaiID_Request data"
// @Router /api/login/request-thaiid [post]
func (h *LoginHandler) RequestThaiID(c *gin.Context) {
	var req models.ThaiID_Request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Invalid JSON input"})
		return
	}

	c.JSON(200, gin.H{
		"url": config.AppConfig.ThaiIDEndPoint + "/auth/?response_type=code&client_id=" + config.AppConfig.ThaiIDClientID + "&redirect_uri=" + req.RedirectUri + "&scope=pid%20name%20birthdate%20openid&state=001",
	})
}

// AuthenThaiID godoc
// @Summary Authenticate user via ThaiID
// @Description This endpoint authenticates a user using ThaiID, providing secure access to protected resources.
// @Tags Login
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param data body models.ThaiID_Authen true "ThaiID_Authen data"
// @Router /api/login/authen-thaiid [post]
func (h *LoginHandler) AuthenThaiID(c *gin.Context) {
	var req models.ThaiID_Authen
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Invalid JSON input"})
		return
	}
	endpoint := config.AppConfig.ThaiIDEndPoint + "/token/"
	// Define form data
	data := url.Values{}
	data.Set("grant_type", "authorization_code")
	data.Set("code", req.Code)
	data.Set("redirect_uri", req.RedirectUri)
	data.Set("client_id", config.AppConfig.ThaiIDClientID)
	data.Set("client_secret", config.AppConfig.ThaiIDSecret)

	// Create request
	creq, err := http.NewRequest("POST", endpoint, strings.NewReader(data.Encode()))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Set headers
	creq.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Execute request
	client := &http.Client{}
	resp, err := client.Do(creq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer resp.Body.Close()

	// Read response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var errorResponse models.KeyCloak_Error_Response
	if err := json.Unmarshal(body, &errorResponse); err == nil {
		if errorResponse.Error != "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": errorResponse.ErrorDescription})
			return
		}
	}

	var successResponse models.KeyCloak_Success_Response
	if err := json.Unmarshal(body, &successResponse); err == nil && successResponse.AccessToken != "" {
		userInfo, err := GetUserInfo(successResponse.AccessToken)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		} else {
			var user models.AuthenUserEmp
			err = config.DB.Where("tel_mobile = ?", userInfo.UserID).
				First(&user).Error

			if err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
				} else {
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				}
				return
			}

			newAccessToken, err := funcs.GenerateJWT(user, "access", time.Duration(config.AppConfig.JwtAccessTokenTime)*time.Minute, "", "")
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Error generating access token"})
				return
			}

			newRefreshToken, err := funcs.GenerateJWT(user, "refresh", time.Duration(config.AppConfig.JwtRefreshTokenTime)*time.Hour, "", "")
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Error generating refresh token"})
				return
			}

			c.JSON(http.StatusOK, models.Login_Response{
				AccessToken:  newAccessToken,
				RefreshToken: newRefreshToken,
			})
			return
		}
	}
	c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
}

func GetUserInfo(accessToken string) (models.KeyCloak_UserInfo, error) {
	userInfoEndpoint := config.AppConfig.KeyCloakEndPoint + "/userinfo"

	req, err := http.NewRequest("GET", userInfoEndpoint, nil)
	if err != nil {
		return models.KeyCloak_UserInfo{}, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return models.KeyCloak_UserInfo{}, fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return models.KeyCloak_UserInfo{}, fmt.Errorf("error reading response: %w", err)
	}

	fmt.Println(string(body))

	var userInfo models.KeyCloak_UserInfo
	if err := json.Unmarshal(body, &userInfo); err != nil {
		return models.KeyCloak_UserInfo{}, fmt.Errorf("error parsing user info: %w", err)
	}

	return userInfo, nil
}

// RequestOTP godoc
// @Summary Request One-Time Password (OTP) for user authentication
// @Description This endpoint allows a user to request a One-Time Password (OTP) for authentication purposes.
// @Tags Login
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param data body models.OTP_Request true "OTP_Request data"
// @Router /api/login/request-otp [post]
func (h *LoginHandler) RequestOTP(c *gin.Context) {
	var req models.OTP_Request
	if err := c.ShouldBindJSON(&req); err != nil || req.Phone == "" {
		c.JSON(http.StatusNotFound, gin.H{"error": "Invalid JSON input", "message": "เกิดความผิดพลาดโปรดลองใหม่"})
		return
	}
	refCode := funcs.RandomRefCode(4)
	expiry := time.Minute * time.Duration(config.AppConfig.OtpExpired)
	otpID, otpErr := SendOTP(req.Phone, refCode, expiry)
	if otpErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "OTP sending failed", "message": "เกิดความผิดพลาดโปรดลองใหม่"})
		return
	}

	expiresAt := time.Now().Add(expiry)
	otpRequest := models.OTP_Request_Create{
		PhoneNo:   req.Phone,
		OTPID:     otpID,
		ExpiresAt: expiresAt,
	}

	if err := config.DB.Create(&otpRequest).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "message": "เกิดความผิดพลาดโปรดลองใหม่"})
		return
	}

	c.JSON(200, gin.H{
		"otpId":   otpID,
		"refCode": refCode,
		"message": "OTP sent successfully",
	})
}

func SendOTP(phone string, refCode string, expiry time.Duration) (string, error) {

	//fmt.Printf("Sending OTP to phone %s\n", phone) // Simulate sending

	soapEndpoint := "https://crm.pea.co.th/Modules/SMS/WebServices/SmsGatewayService.asmx"
	soapAction := "http://crm.pea.co.th/modules/sms/smsgatewayservice/RequestOtpBySmsService"

	soapRequest := fmt.Sprintf(`<?xml version="1.0" encoding="utf-8"?>
<soap:Envelope xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" 
               xmlns:xsd="http://www.w3.org/2001/XMLSchema" 
               xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
  <soap:Body>
    <RequestOtpBySmsService xmlns="http://crm.pea.co.th/modules/sms/smsgatewayservice/">
      <authenKey>545653AA-19E0-41BB-B89F-8485559CD0A7</authenKey>
      <smsServiceId>ae9d5c1b-7ed8-444e-8bb0-707ab7e3e68a</smsServiceId>
      <telephoneNumber>%s</telephoneNumber>
      <messageTemplate>หมายเลข OTP ของท่านคือ **pw** (รหัสอ้างอิง %s) โปรดป้อน ภายใน %d นาที</messageTemplate>
      <timeoutSecond>%d</timeoutSecond>
    </RequestOtpBySmsService>
  </soap:Body>
</soap:Envelope>`, phone, refCode, int(expiry.Minutes()), int(expiry.Seconds()))

	req, err := http.NewRequest("POST", soapEndpoint, bytes.NewBuffer([]byte(soapRequest)))
	if err != nil {
		return "", fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("Content-Type", "text/xml; charset=utf-8")
	req.Header.Set("SOAPAction", soapAction)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error sending request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response: %v", err)
	}

	var envelope models.Envelope
	err = xml.Unmarshal(body, &envelope)
	if err != nil {
		return "", fmt.Errorf("error parsing SOAP response: %v", err)
	}
	// Return the extracted result
	return envelope.Body.RequestOtpBySmsServiceResponse.RequestOtpBySmsServiceResult, nil
}

// VerifyOTP godoc
// @Summary Verify One-Time Password (OTP) for user authentication
// @Description This endpoint allows a user to verify the One-Time Password (OTP) they received.
// @Tags Login
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param data body models.OTPVerify_Request true "OTPVerify_Request data"
// @Router /api/login/verify-otp [post]
func (h *LoginHandler) VerifyOTP(c *gin.Context) {

	var req models.OTPVerify_Request
	if err := c.ShouldBindJSON(&req); err != nil || req.OtpId == "" || req.OTP == "" {
		c.JSON(400, gin.H{"error": "Invalid JSON input", "message": "เกิดความผิดพลาดโปรดลองใหม่"})
		return
	}

	var otpRequest models.OTP_Request_Create
	if err := config.DB.Where("otp_id = ? and status='pending'", req.OtpId).First(&otpRequest).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "OTP request not found", "message": "เกิดความผิดพลาดโปรดลองใหม่"})
			return
		}
		return
	}
	if otpRequest.ExpiresAt.Before(time.Now()) {
		c.JSON(403, gin.H{"error": "รหัส OTP หมดอายุแล้ว กรุณากด 'ขอรหัส OTP ใหม่' เพื่อกรอกอีกครั้ง"})
		return
	}

	// Check OTP
	result, err := CheckOTP(req.OtpId, req.OTP)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "message": "เกิดความผิดพลาดโปรดลองใหม่"})
		return
	}
	if !result {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid OTP", "message": "กรอก OTP ไม่ถูกต้อง ตรวจสอบให้แน่ใจว่าคุณใช้รหัส OTP ที่ได้รับล่าสุด"})
		return
	}

	otpRequest.Status = "verified"
	if err := config.DB.Save(&otpRequest).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to update OTP status: %v", err), "message": "เกิดความผิดพลาดโปรดลองใหม่"})
		return
	}

	var user models.AuthenUserEmp
	err = config.DB.Where("tel_mobile = ?", otpRequest.PhoneNo).
		First(&user).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found", "message": "เกิดความผิดพลาดโปรดลองใหม่"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	// Generate JWT tokens
	accessToken, err := funcs.GenerateJWT(user, "access", time.Duration(config.AppConfig.JwtAccessTokenTime)*time.Minute, "", "")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error generating access token", "message": "เกิดความผิดพลาดโปรดลองใหม่"})
		return
	}

	refreshToken, err := funcs.GenerateJWT(user, "refresh", time.Duration(config.AppConfig.JwtRefreshTokenTime)*time.Hour, "", "")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error generating refresh token", "message": "เกิดความผิดพลาดโปรดลองใหม่"})
		return
	}

	c.JSON(http.StatusOK, models.Login_Response{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})
}

func CheckOTP(otpId string, otp string) (bool, error) {
	// Define the SOAP request body
	soapRequest := fmt.Sprintf(`<?xml version="1.0" encoding="utf-8"?>
<soap:Envelope xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" 
               xmlns:xsd="http://www.w3.org/2001/XMLSchema" 
               xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
  <soap:Body>
    <VerifyOtp xmlns="http://crm.pea.co.th/modules/sms/smsgatewayservice/">
       <authenKey>545653AA-19E0-41BB-B89F-8485559CD0A7</authenKey>
      <otpId>%s</otpId>
      <otp>%s</otp>
    </VerifyOtp>
  </soap:Body>
</soap:Envelope>`, otpId, otp)

	// Send the request to the SOAP API
	url := "https://crm.pea.co.th/Modules/SMS/WebServices/SmsGatewayService.asmx" // Replace with actual endpoint
	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(soapRequest)))
	if err != nil {
		return false, err
	}

	// Set headers
	req.Header.Set("Content-Type", "text/xml; charset=utf-8")
	req.Header.Set("SOAPAction", "http://crm.pea.co.th/modules/sms/smsgatewayservice/VerifyOtp")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}
	// Parse the SOAP response
	var soapResponse models.VerifyOtpSOAPResponse
	if err := xml.Unmarshal(body, &soapResponse); err != nil {
		return false, err
	}

	// Extract result
	return soapResponse.Body.VerifyOtpResponse.VerifyOtpResult == "true", nil
}

// RefreshToken godoc
// @Summary Refresh authentication token
// @Description This endpoint allows a user to refresh their authentication token.
// @Tags Login
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param data body models.RefreshToken_Request true "RefreshToken_Request data"
// @Router /api/login/refresh-token [post]
func (h *LoginHandler) RefreshToken(c *gin.Context) {
	var req models.RefreshToken_Request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON input"})
		return
	}

	claims := &funcs.Claims{}
	token, err := jwt.ParseWithClaims(req.RefreshToken, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.AppConfig.JWTSecret), nil
	})

	if err != nil || !token.Valid || claims.TokenType != "refresh" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid refresh token"})
		return
	}
	var user models.AuthenUserEmp
	err = config.DB.Where("emp_id = ?", claims.EmpID).
		First(&user).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	newAccessToken, err := funcs.GenerateJWT(user, "access", time.Duration(config.AppConfig.JwtAccessTokenTime)*time.Minute, "", "")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error generating access token"})
		return
	}

	newRefreshToken, err := funcs.GenerateJWT(user, "refresh", time.Duration(config.AppConfig.JwtRefreshTokenTime)*time.Hour, "", "")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error generating refresh token"})
		return
	}

	c.JSON(http.StatusOK, models.Login_Response{
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
	})
}

// Logout godoc
// @Summary Log out the current user
// @Description This endpoint allows a user to log out of their session, invalidating their authentication token.
// @Tags Login
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Router /api/logout [get]
func (h *LoginHandler) Logout(c *gin.Context) {
	/*userInfoEndpoint := config.AppConfig.KeyCloakEndPoint + "/userinfo"
	accessToken := ""
	req, err := http.NewRequest("GET", userInfoEndpoint, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	fmt.Println(string(body))
	*/
	c.JSON(http.StatusCreated, gin.H{"message": "Logout successfully"})
}

// Profile godoc
// @Summary Get user profile
// @Description This endpoint retrieves a user profile for the authenticated user.
// @Tags Login
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Router /api/login/profile [get]
func (h *LoginHandler) Profile(c *gin.Context) {
	user := funcs.GetAuthenUser(c, "")
	c.JSON(http.StatusOK, user)
}
