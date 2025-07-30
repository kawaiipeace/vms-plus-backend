package handlers

import (
	"net/http"
	"vms_plus_be/config"
	"vms_plus_be/messages"
	"vms_plus_be/models"
	"vms_plus_be/userhub"

	"github.com/gin-gonic/gin"
)

type ServiceHandler struct {
	Role string
}

func (h *ServiceHandler) checkServiceKey(c *gin.Context, serviceCode string) {
	serviceKey := c.GetHeader("ServiceKey")
	isValid, err := userhub.CheckServiceKey(serviceKey, serviceCode)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized" + err.Error(), "message": messages.ErrUnauthorized.Error()})
		c.Abort()
		return
	}

	if !isValid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized", "message": messages.ErrUnauthorized.Error()})
		c.Abort()
		return
	}
}

// GetVMSToEEMS
// @Summary Get VMS to EEMS
// @Description Get VMS to EEMS
// @Tags Service
// @Accept json
// @Produce json
// @Security ServiceKey
// @Param request_no path string true "RequestNo"
// @router /api/service/vms-to-eems/{request_no} [get]
func (h *ServiceHandler) GetVMSToEEMS(c *gin.Context) {
	h.checkServiceKey(c, "vms")
	if c.IsAborted() {
		return
	}
	requestNo := c.Param("request_no")
	if requestNo == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "RequestNo is required", "message": messages.ErrBadRequest.Error()})
		return
	}
	var request models.VmsToEEMS
	if err := config.DB.
		Table("public.vms_trn_request AS req").
		Select(`req.*, v.vehicle_license_plate,v.vehicle_license_plate_province_short,v.vehicle_license_plate_province_full`).
		Joins("LEFT JOIN vms_mas_vehicle v on v.mas_vehicle_uid = req.mas_vehicle_uid").
		Where("request_no = ?", requestNo).First(&request).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Request not found", "message": messages.ErrBookingNotFound.Error()})
		return
	}

	c.JSON(http.StatusOK, request)
}
