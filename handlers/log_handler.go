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

// GetLogRequest godoc
// @Summary Get log requests by trnRequestUID
// @Tags Log
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Param id path string true "Transaction Request UID"
// @Param page query int false "Page number"
// @Param limit query int false "Limit per page"
// @Router /api/log/request/{id} [get]
func (h *LogHandler) GetLogRequest(c *gin.Context) {
	var logRequests []models.LogRequest
	trnRequestUID := c.Param("id")
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

	c.JSON(http.StatusOK, gin.H{
		"total":      total,
		"page":       page,
		"limit":      limit,
		"totalPages": (total + int64(limit) - 1) / int64(limit), // Calculate total pages
		"logs":       logRequests,
	})
}
