package funcs

import (
	"fmt"
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
		if notifyTemplate.NotifyRole == "vehicle-user" {
			notifyEmpID = request.CreatedRequestEmpID
		} else if notifyTemplate.NotifyRole == "driver" {
			notifyEmpID = request.DriverEmpID

		} else if notifyTemplate.NotifyRole == "level1-approval" {
			notifyEmpID = request.ConfirmedRequestEmpID
		} else if notifyTemplate.NotifyRole == "final-approval" {
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
		if notifyTemplate.NotifyRole == "vehicle-user" {
			notifyEmpID = request.CreatedRequestEmpID
		} else if notifyTemplate.NotifyRole == "level1-approval" {
			notifyEmpID = request.ConfirmedRequestEmpID
		} else if notifyTemplate.NotifyRole == "final-approval" {
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
		}

	}
}
