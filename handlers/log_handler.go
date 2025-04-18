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

func GetRoleOfCreater(refRequestStatusCode string) string {
	switch refRequestStatusCode {
	case "20":
		return "ผู้สร้างคำขอ"
	case "21":
	case "30":
		return "ผู้อนุมัติต้นสังกัด"
	case "40":
		return "ผู้ดูแลยานพาหนะ"
	case "90":
		return "ผู้สร้างคำขอ"
	default:
		return ""
	}
	return ""
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
	config.DB.Model(&models.LogRequest{}).Where("trn_request_uid = ?", trnRequestUID).Count(&total)

	if err := config.DB.Preload("CreatedByEmp").
		Preload("Status").
		Where("trn_request_uid = ?", trnRequestUID).
		Order("created_at desc").
		Limit(limit).
		Offset(offset).
		Find(&logRequests).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	for i := range logRequests {
		logRequests[i].RoleOfCreater = GetRoleOfCreater(logRequests[i].Status.RefRequestStatusCode)
	}

	c.JSON(http.StatusOK, gin.H{
		"total":      total,
		"page":       page,
		"limit":      limit,
		"totalPages": (total + int64(limit) - 1) / int64(limit), // Calculate total pages
		"logs":       logRequests,
	})
}
