package handlers

import (
	"fmt"
	"net/http"
	"vms_plus_be/config"
	"vms_plus_be/messages"
	"vms_plus_be/models"

	"github.com/gin-gonic/gin"
)

type ServiceHandler struct {
	Role string
}

func (h *ServiceHandler) checkServiceKey(c *gin.Context) {
	serviceKey := c.GetHeader("ServiceKey")
	checkSum := 0
	for _, char := range serviceKey {
		checkSum += int(char)
	}
	fmt.Println(checkSum)
	if checkSum != 7340 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized", "message": messages.ErrUnauthorized.Error()})
		c.Abort()
		return
	}
}

// GetRequestBooking
// @Summary Get request booking
// @Description Get request booking
// @Tags Service
// @Accept json
// @Produce json
// @Security ServiceKey
// @Param request_no path string true "RequestNo"
// @router /api/service/request-booking/{request_no} [get]
func (h *ServiceHandler) GetRequestBooking(c *gin.Context) {
	h.checkServiceKey(c)
	if c.IsAborted() {
		return
	}
	requestNo := c.Param("request_no")
	if requestNo == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "RequestNo is required", "message": messages.ErrBadRequest.Error()})
		return
	}
	var request models.VmsTrnRequesService
	if err := config.DB.
		Preload("MasVehicle.RefFuelType").
		Preload("MasVehicle.VehicleDepartment").
		Preload("RefCostType").
		Preload("MasDriver").
		Preload("RefRequestStatus").
		Preload("RequestVehicleType").
		Preload("RefTripType").
		Preload("TripDetails").
		Preload("AddFuels").
		Preload("AddFuels.RefCostType").
		Preload("AddFuels.RefOilStationBrand").
		Preload("AddFuels.RefFuelType").
		Preload("AddFuels.RefPaymentType").
		Preload("VehicleImagesReceived").
		Preload("VehicleImagesReturned").
		Preload("VehicleImagesReturned").
		Preload("VehicleImageInspect").
		Preload("ReceiverKeyTypeDetail").
		Where("request_no = ?", requestNo).First(&request).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Request not found", "message": messages.ErrBookingNotFound.Error()})
		return
	}

	c.JSON(http.StatusOK, request)
}

func (h *ServiceHandler) GetVMSToEEMS(c *gin.Context) {
	h.checkServiceKey(c)
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
		Preload("MasVehicle.RefFuelType").
		Preload("MasVehicle.VehicleDepartment").
		Preload("RefCostType").
		Preload("MasDriver").
		Preload("RefRequestStatus").
		Preload("RequestVehicleType").
		Preload("RefTripType").
		Preload("TripDetails").
		Preload("AddFuels").
		Preload("AddFuels.RefCostType").
		Preload("AddFuels.RefOilStationBrand").
		Preload("AddFuels.RefFuelType").
		Preload("AddFuels.RefPaymentType").
		Preload("VehicleImagesReceived").
		Preload("VehicleImagesReturned").
		Preload("VehicleImagesReturned").
		Preload("VehicleImageInspect").
		Preload("ReceiverKeyTypeDetail").
		Where("request_no = ?", requestNo).First(&request).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Request not found", "message": messages.ErrBookingNotFound.Error()})
		return
	}

	c.JSON(http.StatusOK, request)
}