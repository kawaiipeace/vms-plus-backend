package funcs

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
	"vms_plus_be/config"
	"vms_plus_be/models"

	"github.com/google/uuid"
)

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
			SendNotificationPEA(notifyEmpID, notifyTemplate.NotifyTitle+" "+notifyMessage)
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
			SendNotificationPEA(notifyEmpID, notifyTemplate.NotifyTitle+" "+notifyMessage)
		}

	}
}

func SendNotificationPEA(empID, message string) {
	if config.AppConfig.PEANotificationEndPoint == "" || config.AppConfig.PEANotificationToken == "" {
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
