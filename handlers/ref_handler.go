package handlers

import (
	"net/http"
	"vms_plus_be/config"
	"vms_plus_be/funcs"
	"vms_plus_be/messages"
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
		c.JSON(http.StatusNotFound, gin.H{"error": "Not found", "message": messages.ErrNotfound.Error()})
		return
	}
	c.JSON(http.StatusOK, lists)
}

// ListVehicleStatus godoc
// @Summary Retrieve all vehicle statuses
// @Description This endpoint retrieves all vehicle statuses.
// @Tags REF
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Router /api/ref/vehicle-status [get]
func (h *RefHandler) ListVehicleStatus(c *gin.Context) {
	var lists []models.VmsRefVehicleStatus
	if err := config.DB.
		Find(&lists).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Not found", "message": messages.ErrNotfound.Error()})
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
// @Param emp_id query string false "Employee ID (emp_id) default(700001)"
// @Router /api/ref/cost-type [get]
func (h *RefHandler) ListCostType(c *gin.Context) {
	var lists []models.VmsRefCostType
	if err := config.DB.
		Find(&lists).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Not found", "message": messages.ErrNotfound.Error()})
		return
	}
	for i := range lists {
		if lists[i].RefCostTypeCode == "1" {
			empID := c.Query("emp_id")
			empUser := funcs.GetUserEmpInfo(empID)
			if empUser.DeptSAP == "" {
				lists[i].CostCenter = ""
				continue
			}
			var department models.VmsMasDepartment
			if err := config.DB.
				Where("dept_sap = ?", empUser.DeptSAP).
				First(&department).Error; err == nil {
				lists[i].CostCenter = department.CostCenterCode + "  " + department.CostCenterName
			} else {
				lists[i].CostCenter = ""
			}
		}

	}
	c.JSON(http.StatusOK, lists)
}

// ListCostCenter godoc
// @Summary Retrieve all cost centers
// @Description This endpoint retrieves all cost centers.
// @Tags REF
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Param search query string false "Search cost_center_code,cost_center_name"
// @Router /api/ref/cost-center [get]
func (h *RefHandler) ListCostCenter(c *gin.Context) {
	var lists []models.VmsRefCostCenter
	search := c.Query("search")
	query := config.DB.
		Select("DISTINCT cost_center_code || ' ' || cost_center_name as cost_center").
		Where("cost_center_code IS NOT NULL AND cost_center_code != ''")
	if search != "" {
		query = query.Where("cost_center_code || ' ' || cost_center_name LIKE ?", "%"+search+"%")
	}
	if err := query.
		Order("cost_center").
		Find(&lists).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Not found", "message": messages.ErrNotfound.Error()})
		return
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
	if c.IsAborted() {
		return
	}
	code := c.Param("code")
	user := funcs.GetAuthenUser(c, "*")
	var costType models.VmsRefCostType
	if err := config.DB.
		First(&costType, "ref_cost_type_code = ?", code).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Cost type not found", "message": messages.ErrNotfound.Error()})
		return
	}
	if costType.RefCostTypeCode == "1" {
		var department models.VmsMasDepartment
		if err := config.DB.
			Where("dept_sap = ?", user.DeptSAP).
			First(&department).Error; err == nil {
			costType.CostCenter = department.CostCenterCode + "  " + department.CostCenterName
		} else {
			costType.CostCenter = ""
		}
	}
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

// ListVehicleImgSide godoc
// @Summary Retrieve all vehicle image sides
// @Description This endpoint retrieves all vehicle image sides.
// @Tags REF
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Router /api/ref/vehicle-img-side [get]
func (h *RefHandler) ListVehicleImgSide(c *gin.Context) {
	var lists []models.VmsRefVehicleImgSide
	if err := config.DB.
		Find(&lists).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Not found", "message": messages.ErrNotfound.Error()})
		return
	}
	c.JSON(http.StatusOK, lists)
}

// ListPaymentTypeCode godoc
// @Summary Retrieve all payment type codes
// @Description This endpoint retrieves all payment type codes.
// @Tags REF
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Router /api/ref/payment-type-code [get]
func (h *RefHandler) ListPaymentTypeCode(c *gin.Context) {
	var lists []models.VmsRefPaymentType
	if err := config.DB.
		Find(&lists).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Not found", "message": messages.ErrNotfound.Error()})
		return
	}
	c.JSON(http.StatusOK, lists)
}

// ListDriverOtherUse godoc
// @Summary Retrieve all payment type codes
// @Description This endpoint retrieves all payment type codes.
// @Tags REF
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Router /api/ref/driver-other-use [get]
func (h *RefHandler) ListDriverOtherUse(c *gin.Context) {
	var lists []models.VmsRefOtherUse
	if err := config.DB.
		Find(&lists).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Not found", "message": messages.ErrNotfound.Error()})
		return
	}
	c.JSON(http.StatusOK, lists)
}

// ListDriverLicenseType godoc
// @Summary Retrieve all driver license types
// @Description This endpoint retrieves all driver license types.
// @Tags REF
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Router /api/ref/driver-license-type [get]
func (h *RefHandler) ListDriverLicenseType(c *gin.Context) {
	var lists []models.VmsRefDriverLicenseType
	if err := config.DB.
		Find(&lists).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Not found", "message": messages.ErrNotfound.Error()})
		return
	}
	c.JSON(http.StatusOK, lists)
}

// ListDriverCertificateType godoc
// @Summary Retrieve all driver certificate types
// @Description This endpoint retrieves all driver certificate types.
// @Tags REF
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Router /api/ref/driver-certificate-type [get]
func (h *RefHandler) ListDriverCertificateType(c *gin.Context) {
	var lists []models.VmsRefDriverCertificateType
	if err := config.DB.
		Find(&lists).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Not found", "message": messages.ErrNotfound.Error()})
		return
	}
	c.JSON(http.StatusOK, lists)
}

// ListCarpoolChooseCar godoc
// @Summary Retrieve all carpool choose car options
// @Description This endpoint retrieves all carpool choose car options.
// @Tags REF
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Router /api/ref/carpool-choose-car [get]
func (h *RefHandler) ListCarpoolChooseCar(c *gin.Context) {
	var lists []models.VmsRefCarpoolChooseCar
	if err := config.DB.
		Find(&lists).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Not found", "message": messages.ErrNotfound.Error()})
		return
	}
	c.JSON(http.StatusOK, lists)
}

// ListCarpoolChooseDriver godoc
// @Summary Retrieve all carpool choose driver options
// @Description This endpoint retrieves all carpool choose driver options.
// @Tags REF
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Router /api/ref/carpool-choose-driver [get]
func (h *RefHandler) ListCarpoolChooseDriver(c *gin.Context) {
	var lists []models.VmsRefCarpoolChooseDriver
	if err := config.DB.
		Find(&lists).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Not found", "message": messages.ErrNotfound.Error()})
		return
	}
	c.JSON(http.StatusOK, lists)
}

// ListVehicleKeyType godoc
// @Summary Retrieve all vehicle key types
// @Description This endpoint retrieves all vehicle key types.
// @Tags REF
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Router /api/ref/vehicle-key-type [get]
func (h *RefHandler) ListVehicleKeyType(c *gin.Context) {
	var lists []models.VmsRefVehicleKeyType
	if err := config.DB.
		Find(&lists).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Not found", "message": messages.ErrNotfound.Error()})
		return
	}
	c.JSON(http.StatusOK, lists)
}

// ListLeaveTimeType godoc
// @Summary Retrieve all leave time types
// @Description This endpoint retrieves all leave time types.
// @Tags REF
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Router /api/ref/leave-time-type [get]
func (h *RefHandler) ListLeaveTimeType(c *gin.Context) {
	var lists []models.VmsRefLeaveTimeType
	if err := config.DB.
		Find(&lists).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Not found", "message": messages.ErrNotfound.Error()})
		return
	}
	c.JSON(http.StatusOK, lists)
}

// ListDriverStatus godoc
// @Summary Retrieve all driver statuses
// @Description This endpoint retrieves all driver statuses.
// @Tags REF
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Router /api/ref/driver-status [get]
func (h *RefHandler) ListDriverStatus(c *gin.Context) {
	var lists []models.VmsRefDriverStatus
	if err := config.DB.
		Find(&lists).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Not found", "message": messages.ErrNotfound.Error()})
		return
	}
	c.JSON(http.StatusOK, lists)
}

// ListTimelineStatus godoc
// @Summary Retrieve all timeline statuses
// @Description This endpoint retrieves all timeline statuses.
// @Tags REF
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Router /api/ref/timeline-status [get]
func (h *RefHandler) ListTimelineStatus(c *gin.Context) {
	var lists []models.VmsRefTimelineStatus
	lists = append(lists, models.VmsRefTimelineStatus{
		RefTimelineStatusID:   "1",
		RefTimelineStatusName: "รออนุมัติ",
	})
	lists = append(lists, models.VmsRefTimelineStatus{
		RefTimelineStatusID:   "2",
		RefTimelineStatusName: "ไป-กลับ",
	})
	lists = append(lists, models.VmsRefTimelineStatus{
		RefTimelineStatusID:   "3",
		RefTimelineStatusName: "ค้างแรม",
	})
	lists = append(lists, models.VmsRefTimelineStatus{
		RefTimelineStatusID:   "4",
		RefTimelineStatusName: "เสร็จสิ้น",
	})

	c.JSON(http.StatusOK, lists)
}
