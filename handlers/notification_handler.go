package handlers

import (
	"net/http"
	"time"
	"vms_plus_be/config"
	"vms_plus_be/funcs"
	"vms_plus_be/models"

	"github.com/gin-gonic/gin"
)

type NotificationHandler struct {
	Role string
}

// GetNotification godoc
// @Summary Get Notification
// @Description Get Notification
// @Tags Notification
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Router /api/notification [get]
func (h *NotificationHandler) GetNotification(c *gin.Context) {
	user := funcs.GetAuthenUser(c, "*")
	if c.IsAborted() {
		return
	}

	var notifys []models.Notification
	var total, unread int64
	err := config.DB.Where("emp_id = ?", user.EmpID).Order("created_at DESC").Find(&notifys)
	if err.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error, "message": "Failed to get notifications"})
		return
	}
	config.DB.Model(&models.Notification{}).Where("emp_id = ?", user.EmpID).Count(&total)
	config.DB.Model(&models.Notification{}).Where("emp_id = ? AND is_read = ?", user.EmpID, "0").Count(&unread)

	for i, notify := range notifys {
		notifys[i].Duration = funcs.GetDuration(notify.CreatedAt)
		if notify.NotifyType == "request-booking" {
			notifys[i].NotifyURL = notify.NotifyURL + "?trn_request_uid=" + notify.RecordUID
		}
		notifys[i].NotifyURL = funcs.GetNotifyURL(notify)
	}
	if len(notifys) == 0 {
		notifys = []models.Notification{}
	}

	c.JSON(http.StatusOK, gin.H{
		"notifications": notifys,
		"total":         total,
		"unread":        unread,
	})
}

// UpdateReadNotification godoc
// @Summary Update Read Notification
// @Description Update Read Notification
// @Tags Notification
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Param notification_uid path string true "Notification UID"
// @Router /api/notification/read/{notification_uid} [put]
func (h *NotificationHandler) UpdateReadNotification(c *gin.Context) {
	user := funcs.GetAuthenUser(c, "*")
	if c.IsAborted() {
		return
	}
	notificationUID := c.Param("notification_uid")

	var notification models.Notification
	err := config.DB.Where("trn_notify_uid = ? AND emp_id = ?", notificationUID, user.EmpID).First(&notification).Error
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error, "message": "Notification not found"})
		return
	}
	//update set is_read = true,read_at = now()
	err = config.DB.Model(&models.Notification{}).Where("trn_notify_uid = ? AND emp_id = ?", notificationUID, user.EmpID).
		Update("is_read", true).
		Update("read_at", time.Now()).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error, "message": "Failed to update notification"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Notification updated"})
}
