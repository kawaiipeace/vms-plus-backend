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
	for i := range lists {
		lists[i].RefCostNo = lists[i].RefCostTypeCode + "-00000-1"
	}
	c.JSON(http.StatusOK, lists)
}

// GetCostType godoc
// @Summary Retrieve a specific cost type
// @Description This endpoint fetches details of a cost type.
// @Tags REF
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Param code path string true "ref_cost_type_code (ref_cost_type_code)"
// @Router /api/ref/cost-type/{code} [get]
func (h *RefHandler) GetCostType(c *gin.Context) {
	//funcs.GetAuthenUser(c, h.Role)
	code := c.Param("code")

	var costType models.VmsRefCostType
	if err := config.DB.
		First(&costType, "ref_cost_type_code = ?", code).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Cost type not found"})
		return
	}
	costType.RefCostNo = costType.RefCostTypeCode + "-00000-1"
	c.JSON(http.StatusOK, costType)
}

// ListFuelType godoc
// @Summary Retrieve the all fuel type
// @Description This endpoint retrieve all fuel type
// @Tags REF
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Router /api/ref/fuel-type [get]
func (h *RefHandler) ListFuelType(c *gin.Context) {
	var lists []models.VmsRefFuelType
	if err := config.DB.
		Find(&lists).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Not found"})
		return
	}
	c.JSON(http.StatusOK, lists)
}

// ListOilStationBrand godoc
// @Summary Retrieve the all oil station brand
// @Description This endpoint retrieve all oil station brand.
// @Tags REF
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Router /api/ref/oil-station-brand [get]
func (h *RefHandler) ListOilStationBrand(c *gin.Context) {
	var lists []models.VmsRefOilStationBrand
	if err := config.DB.
		Find(&lists).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Not found"})
		return
	}
	c.JSON(http.StatusOK, lists)
}
