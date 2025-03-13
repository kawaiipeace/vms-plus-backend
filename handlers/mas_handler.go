package handlers

import (
	"net/http"
	"vms_plus_be/config"
	"vms_plus_be/models"

	"github.com/gin-gonic/gin"
)

type MasHandler struct {
}

// ListVehicleUser godoc
// @Summary Retrieve the Vehicle User
// @Description This endpoint allows a user to retrieve Vehicle User.
// @Tags MAS
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Param search query string false "Search by Employee ID or Full Name"
// @Router /api/mas/user-vehicle-users [get]
func (h *MasHandler) ListVehicleUser(c *gin.Context) {
	var lists []models.MasUserEmp
	search := c.Query("search")

	query := config.DB

	// Apply search filter if provided
	if search != "" {
		query = query.Where("emp_id = ? OR full_name ILIKE ?", search, "%"+search+"%")
	}

	// Execute query
	if err := query.Find(&lists).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Not found"})
		return
	}

	c.JSON(http.StatusOK, lists)
}
