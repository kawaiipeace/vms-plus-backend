package funcs

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"slices"
	"strings"
	"time"
	"vms_plus_be/config"
	"vms_plus_be/models"
	"vms_plus_be/userhub"

	"github.com/google/uuid"
)

func IsAllowNotifyEmpID(empID string) bool {
	empIDs := []string{
		"465056", //นางสาวพจนา พานิชนิตินนท์
		"499910", //นางสาวนพรัตน์ อภิชาตสิริธรรม
		"460137", //นายวิทยา สว่างวงษ์
		"505291", //นายศรัญยู บริรัตน์ฤทธิ์
		"511181", //นายจอมภูภพ อิศโร
		"514285", //นายธนพล วิจารณ์ปรีชา
	}
	return slices.Contains(empIDs, empID)
}

func GetNotifyURL(notify models.Notification) string {
	if notify.NotifyRole == "vehicle-user" && notify.NotifyType == "request-booking" &&
		Contains([]string{"20", "21", "30", "31", "40", "41", "90"}, notify.RefRequestStatusCode) {
		return "vehicle-booking/request-list/" + notify.RecordUID
	}
	if notify.NotifyRole == "vehicle-user" && notify.NotifyType == "request-booking" &&
		Contains([]string{"50", "51", "60", "70", "80", "90"}, notify.RefRequestStatusCode) {
		return "vehicle-in-use/user/" + notify.RecordUID
	}
	if notify.NotifyRole == "vehicle-user" && notify.NotifyType == "request-annual-driver" &&
		Contains([]string{"10", "11", "20", "21", "30", "90"}, notify.RefRequestStatusCode) {
		return "vehicle-booking/request-list/" + notify.RecordUID
	}
	if notify.NotifyRole == "driver" {
		return "vehicle-booking/request-list/" + notify.RecordUID
	}
	if notify.NotifyRole == "level1-approval" && notify.NotifyType == "request-annual-driver" &&
		Contains([]string{"10", "11", "20"}, notify.RefRequestStatusCode) {
		return "/administrator/driver-license-confirmer/" + notify.RecordUID
	}
	if notify.NotifyRole == "level1-approval" && notify.NotifyType == "request-booking" &&
		Contains([]string{"20", "21", "30"}, notify.RefRequestStatusCode) {
		return "/administrator/booking-approver/" + notify.RecordUID
	}
	if notify.NotifyRole == "admin-department" && notify.NotifyType == "request-booking" &&
		Contains([]string{"30", "31", "40"}, notify.RefRequestStatusCode) {
		return "/administrator/booking-approver/" + notify.RecordUID
	}
	if notify.NotifyRole == "admin-department" && notify.NotifyType == "request-booking" &&
		Contains([]string{"50", "51", "60", "70", "80"}, notify.RefRequestStatusCode) {
		return "/administrator/vehicle-in-use/" + notify.RecordUID
	}
	if notify.NotifyRole == "final-approval" && notify.NotifyType == "request-booking" &&
		Contains([]string{"40", "41", "50"}, notify.RefRequestStatusCode) {
		return "/administrator/booking-final/" + notify.RecordUID
	}
	if notify.NotifyRole == "license-approver" && notify.NotifyType == "request-annual-driver" &&
		Contains([]string{"20", "21", "30"}, notify.RefRequestStatusCode) {
		return "/administrator/driver-license-approver/" + notify.RecordUID
	}
	return ""
}

func CreateRequestBookingNotification(trnRequestUID string) {
	var request models.RequestBookingNotification
	fmt.Println("trnRequestUID:", trnRequestUID)
	if err := config.DB.Where("trn_request_uid = ?", trnRequestUID).First(&request).Error; err != nil {
		fmt.Println("Error getting request booking:", err)
		return
	}

	var notifyTemplates []models.NotificationTemplate

	if err := config.DB.Where("ref_request_status_code = ? AND is_deleted = false AND notify_type = 'request-booking'", request.RefRequestStatusCode).Find(&notifyTemplates).Error; err != nil {
		fmt.Println("Error getting notify templates:", err)
		return
	}

	for _, notifyTemplate := range notifyTemplates {
		var notifyEmpID string
		notifyMessage := notifyTemplate.NotifyMessage
		notifyMessage = strings.Replace(notifyMessage, "**request_no**", request.RequestNo, -1)
		fmt.Println("Notify role:", notifyTemplate.NotifyRole)
		switch notifyTemplate.NotifyRole {
		case "vehicle-user":
			notifyEmpID = request.CreatedRequestEmpID
		case "driver":
			notifyEmpID = request.DriverEmpID
		case "level1-approval":
			notifyEmpID = request.ConfirmedRequestEmpID
		case "approval-department":
			notifyEmpID = request.ApprovedRequestEmpID
		case "approval-carpool":
			notifyEmpID = request.ApprovedRequestEmpID
		}
		if notifyEmpID != "" {
			//create notification
			notification := models.Notification{
				TrnNotifyUID:         uuid.New().String(),
				EmpID:                notifyEmpID,
				Title:                notifyTemplate.NotifyTitle,
				Message:              notifyMessage,
				RecordUID:            request.TrnRequestUID,
				NotifyType:           notifyTemplate.NotifyType,
				NotifyRole:           notifyTemplate.NotifyRole,
				RefRequestStatusCode: request.RefRequestStatusCode,
				IsRead:               false,
				CreatedAt:            time.Now(),
			}
			if err := config.DB.Create(&notification).Error; err != nil {
				fmt.Println("Error creating notification:", err)
				return
			}
			userInfo, err := userhub.GetUserInfo(notifyEmpID)
			if err != nil {
				fmt.Println("Error getting user info:", err)
				return
			}
			go SendNotificationWorkD(notifyEmpID, notifyTemplate.NotifyTitle, notifyMessage, "", GetNotifyURL(notification), userInfo.DeptSAP, userInfo.DeptSAPShort)
			go SendNotificationPEA(notifyEmpID, notifyTemplate.NotifyTitle+" "+notifyMessage)
			go SendNotificationSMS(notifyEmpID, notifyTemplate.NotifyTitle+" "+notifyMessage)
		}

	}
}

func CreateRequestAnnualLicenseNotification(trnAnnualLicenseUID string) {
	var request models.RequestAnnualLicenseNotification
	fmt.Println("trnRequestUID:", trnAnnualLicenseUID)
	if err := config.DB.Where("trn_request_annual_driver_uid = ?", trnAnnualLicenseUID).First(&request).Error; err != nil {
		fmt.Println("Error getting request booking:", err)
		return
	}

	var notifyTemplates []models.NotificationTemplate

	if err := config.DB.Where("ref_request_status_code = ? AND is_deleted = false AND notify_type = 'request-annual-driver'", request.RefRequestAnnualDriverStatusCode).Find(&notifyTemplates).Error; err != nil {
		fmt.Println("Error getting notify templates:", err)
		return
	}

	for _, notifyTemplate := range notifyTemplates {
		var notifyEmpID string
		notifyMessage := notifyTemplate.NotifyMessage
		notifyMessage = strings.Replace(notifyMessage, "**request_no**", request.RequestAnnualDriverNo, -1)
		fmt.Println("Notify role:", notifyTemplate.NotifyRole)
		switch notifyTemplate.NotifyRole {
		case "vehicle-user":
			notifyEmpID = request.CreatedRequestEmpID
		case "level1-approval":
			notifyEmpID = request.ConfirmedRequestEmpID
		case "license-approval":
			notifyEmpID = request.ApprovedRequestEmpID
		}
		if notifyEmpID != "" {
			//create notification
			notification := models.Notification{
				TrnNotifyUID:         uuid.New().String(),
				EmpID:                notifyEmpID,
				Title:                notifyTemplate.NotifyTitle,
				Message:              notifyMessage,
				RecordUID:            request.TrnRequestAnnualDriverUID,
				NotifyType:           notifyTemplate.NotifyType,
				NotifyRole:           notifyTemplate.NotifyRole,
				RefRequestStatusCode: request.RefRequestAnnualDriverStatusCode,
				IsRead:               false,
				CreatedAt:            time.Now(),
			}
			if err := config.DB.Create(&notification).Error; err != nil {
				fmt.Println("Error creating notification:", err)
				return
			}
			userInfo, err := userhub.GetUserInfo(notifyEmpID)
			if err != nil {
				fmt.Println("Error getting user info:", err)
				return
			}
			go SendNotificationPEA(notifyEmpID, notifyTemplate.NotifyTitle+" "+notifyMessage)
			go SendNotificationWorkD(notifyEmpID, notifyTemplate.NotifyTitle, notifyMessage, "", GetNotifyURL(notification), userInfo.DeptSAP, userInfo.DeptSAPShort)
			go SendNotificationSMS(notifyEmpID, notifyTemplate.NotifyTitle+" "+notifyMessage)
		}

	}
}

func SendNotificationPEA(empID, message string) {
	if config.AppConfig.PEANotificationEndPoint == "" || config.AppConfig.PEANotificationToken == "" {
		return
	}
	if !IsAllowNotifyEmpID(empID) {
		return
	}
	body := models.NotificationRequestBodyPEA{
		EmployeeId:    empID,
		MessageTypeID: "11",
		Message:       message,
	}

	jsonData, err := json.Marshal(body)
	if err != nil {
		fmt.Printf("Error marshalling JSON: %v", err)
	}

	// Create HTTP request
	req, err := http.NewRequest("POST", config.AppConfig.PEANotificationEndPoint, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Printf("Error creating request: %v", err)
	}

	req.Header.Set("Authorization", config.AppConfig.PEANotificationToken)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "PostmanRuntime/7.43.4")
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")

	// Create custom HTTP client with insecure TLS config
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	// Send HTTP request
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error sending request: %v", err)
		return
	}
	defer resp.Body.Close()

	// Read response
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading response: %v", err)
	}

	fmt.Printf("Response Status: %s\n", resp.Status)
	fmt.Printf("Response Body: %s\n", responseBody)
}

func SendNotificationWorkD(empID, headline, subHeadline, content, url, deptSap, targetName string) {
	if config.AppConfig.PEAWorkDNotificationEndPoint == "" || config.AppConfig.PEAWorkDNotificationToken == "" {
		return
	}

	if !IsAllowNotifyEmpID(empID) {
		return
	}

	payloadStr := fmt.Sprintf(`{
		"notificationType": "MESSAGE",
		"id": null,
		"headline": "%s",	
		"subHeadline": "%s",
		"content": "%s",
		"mobileUrlScheme": %s,
		"desktopUrlScheme": %s,
		"tabletUrlScheme": %s,
		"system": "VMS Plus",
		"notificationTargetType": "SPECIFIC",
		"notificationSpecificGroup": [
			{
				"deptSap": "%s",
				"targetName": "%s"
			}
		],
		"publishedDate": "%s",
		"expiredDate": "%s",
		"isModal": false
	}`, headline, subHeadline, content, url, url, url, deptSap, targetName,
		time.Now().Format("2006-01-02 15:04:05"),
		time.Now().Add(24*time.Hour).Format("2006-01-02 15:04:05"))

	payload := []byte(payloadStr)

	req, err := http.NewRequest("POST", config.AppConfig.PEAWorkDNotificationEndPoint, bytes.NewBuffer(payload))
	if err != nil {
		panic(err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Public-Key", config.AppConfig.PEAWorkDNotificationToken)
	req.Header.Set("username", empID)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	fmt.Println("Response Status:", resp.Status)
	fmt.Println("Response Body:", string(body))
}

func SendNotificationSMS(empID, message string) {

	userInfo, err := userhub.GetUserInfo(empID)
	if err != nil {
		fmt.Printf("Error getting user info: %v", err)
		return
	}

	if userInfo.MobilePhone == "" {
		return
	}

	soapEndpoint := "https://crm.pea.co.th/Modules/SMS/WebServices/SmsGatewayService.asmx"
	soapAction := "http://crm.pea.co.th/modules/sms/smsgatewayservice/SendSms"

	soapRequest := fmt.Sprintf(`<?xml version="1.0" encoding="utf-8"?>
<soap:Envelope xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" 
               xmlns:xsd="http://www.w3.org/2001/XMLSchema" 
               xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
  <soap:Body>
    <SendSms xmlns="http://crm.pea.co.th/modules/sms/smsgatewayservice/">
      <authenKey>545653AA-19E0-41BB-B89F-8485559CD0A7</authenKey>
      <smsServiceId>ae9d5c1b-7ed8-444e-8bb0-707ab7e3e68a</smsServiceId>
      <telephoneNumber>%s</telephoneNumber>
      <message>%s</message>
    </SendSms>
  </soap:Body>
</soap:Envelope>`, userInfo.MobilePhone, message)

	req, err := http.NewRequest("POST", soapEndpoint, bytes.NewBuffer([]byte(soapRequest)))
	if err != nil {
		fmt.Printf("Error creating request: %v", err)
		return
	}

	req.Header.Set("Content-Type", "text/xml; charset=utf-8")
	req.Header.Set("SOAPAction", soapAction)

	// Create custom HTTP client with timeout and TLS config
	client := &http.Client{
		Timeout: 30 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error sending request: %v", err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading response: %v", err)
		return
	}

	var envelope models.Envelope
	err = xml.Unmarshal(body, &envelope)
	if err != nil {
		fmt.Printf("Error parsing SOAP response: %v", err)
		return
	}
}
