package handlers

import (
	"net/http"
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
