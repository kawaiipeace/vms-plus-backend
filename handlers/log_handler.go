package handlers

import (
	"net/http"
	"strconv"
	"vms_plus_be/config"
	"vms_plus_be/models"

	"github.com/gin-gonic/gin"
)

type LogHandler struct {
}

func GetRoleOfCreater(role string) string {
	switch role {
	case "vehicle-user":
		return "ผู้ใช้ยานพาหนะ"
	case "level1-approval":
		return "ผู้อนุมัติต้นสังกัด"
	case "admin-department":
		return "ผู้ดูแลยานพาหนะ"
	case "admin-carpool":
		return "ผู้ดูแลยานพาหนะ"
	case "approval-carpool":
		return "ผู้อนุมัติขั้นสุดท้าย"
	case "approval-department":
		return "ผู้อนุมัติขั้นสุดท้าย"
	case "driver":
		return "คนขับรถ"
	default:
		return ""
	}
}

// GetLogRequest godoc
// @Summary Get log requests by trnRequestUID
// @Tags Log
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Param trn_request_uid path string true "trn_request_uid"
// @Param page query int false "Page number"
// @Param limit query int false "Limit per page"
// @Router /api/log/request/{trn_request_uid} [get]
func (h *LogHandler) GetLogRequest(c *gin.Context) {
	var logRequests []models.LogRequest
	trnRequestUID := c.Param("trn_request_uid")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	offset := (page - 1) * limit
	var total int64
	config.DB.Model(&models.LogRequest{}).Where("trn_request_uid = ? AND is_deleted = ?", trnRequestUID, "0").Count(&total)

	if err := config.DB.
		Where("trn_request_uid = ? AND is_deleted = ?", trnRequestUID, "0").
		Order("log_request_action_datetime").
		Limit(limit).
		Offset(offset).
		Find(&logRequests).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	for i := range logRequests {
		logRequests[i].RoleOfCreater = GetRoleOfCreater(logRequests[i].ActionByRole)
	}

	c.JSON(http.StatusOK, gin.H{
		"total":      total,
		"page":       page,
		"limit":      limit,
		"totalPages": (total + int64(limit) - 1) / int64(limit), // Calculate total pages
		"logs":       logRequests,
	})
}
