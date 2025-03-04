package handlers

import (
	"net/http"
	"vms_plus_be/config"
	"vms_plus_be/models"

	"github.com/gin-gonic/gin"
)

type RefHandler struct {
}

// ListRequestStatus godoc
// @Summary Retrieve the status of booking requests
// @Description This endpoint allows a booking user to retrieve the status of their booking requests.
// @Tags REF
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Router /api/ref/request-status [get]
func (h *RefHandler) ListRequestStatus(c *gin.Context) {
	var lists []models.VmsRefRequestStatus
	if err := config.DB.
		Find(&lists).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Not found"})
		return
	}
	c.JSON(http.StatusOK, lists)
}

// ListCostType godoc
// @Summary Retrieve available cost types
// @Description This endpoint allows a user to retrieve a list of available cost types for booking requests.
// @Tags REF
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Router /api/ref/cost-type [get]
func (h *RefHandler) ListCostType(c *gin.Context) {
	var lists []models.VmsRefCostType
	if err := config.DB.
		Find(&lists).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Not found"})
		return
	}
	c.JSON(http.StatusOK, lists)
}
