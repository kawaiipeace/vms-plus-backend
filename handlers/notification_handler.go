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

	config.DB.Where("emp_id = ?", user.EmpID).Find(&notifys).
		Order("created_at DESC")
	for i, notify := range notifys {
		notifys[i].Duration = funcs.GetDuration(notify.CreatedAt)
	}
	if len(notifys) == 0 {
		notifys = []models.Notification{}
	}
	c.JSON(http.StatusOK, notifys)
}
