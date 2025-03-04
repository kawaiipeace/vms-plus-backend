package handlers

import (
	"net/http"
	"time"
	"vms_plus_be/config"
	"vms_plus_be/funcs"
	"vms_plus_be/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type BookingUserHandler struct {
	Role string
}

// CreateRequest godoc
// @Summary Create a new booking request
// @Description This endpoint allows a booking user to create a new request.
// @Tags Booking-user
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Param data body models.VmsTrnRequest_Request true "VmsTrnRequest_Request data"
// @Router /api/booking-user/create-request [post]
func (h *BookingUserHandler) CreateRequest(c *gin.Context) {
	user := funcs.GetAuthenUser(c, h.Role)

	var req models.VmsTrnRequest_Request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON input"})
		return
	}

	logCreate :=
		models.LogCreate{
			CreatedAt: time.Now(),
			CreatedBy: user.EmpID,
		}
	vms_trn_req := models.VmsTrnRequest_Create{
		TrnRequestUID:              uuid.New().String(),
		VmsTrnRequest_Request:      req,
		CreatedRequestEmpID:        user.EmpID,
		CreatedRequestEmpName:      user.FullName(),
		CreatedRequestDeptSAP:      user.DeptSAP,
		CreatedRequestDeptSAPShort: user.DeptSAPShort,
		CreatedRequestDeptSAPFull:  user.DeptSAPFull,
		LogCreate:                  logCreate,
	}

	if err := config.DB.Create(&vms_trn_req).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create request"})
		return
	}
	funcs.CreateTrnLog(vms_trn_req.TrnRequestUID,
		vms_trn_req.RefRequestStatusCode,
		"Create new request",
		user.EmpID)

	c.JSON(http.StatusCreated, gin.H{"message": "Request created successfully", "data": req})
}

// ListRequest godoc
// @Summary Get a list of booking requests
// @Description This endpoint retrieves a list of booking requests for the authenticated user.
// @Tags Booking-user
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Router /api/booking-user/requests [get]
func (h *BookingUserHandler) ListRequest(c *gin.Context) {
	funcs.GetAuthenUser(c, h.Role)
	funcs.ListRequest(c)
}

// SearchRequests godoc
// @Summary Search booking requests and get summary counts by request status code
// @Description Search for requests using a keyword and get the summary of counts grouped by request status code
// @Tags Booking-user
// @Accept  json
// @Produce  json
// @Param search query string false "Search keyword (matches request_no, vehicle_license_plate, vehicle_user_emp_name, or work_place)"
// @Param ref_request_status_code query string false "Filter by multiple request status codes (comma-separated, e.g., 'A,B,C')"
// @Param startdate query string false "Filter by start datetime (YYYY-MM-DD format)"
// @Param enddate query string false "Filter by end datetime (YYYY-MM-DD format)"
// @Param order_by query string false "Order by request_no, start_datetime, ref_request_status_code"
// @Param order_dir query string false "Order direction: asc or desc"
// @Param page query int false "Page number (default: 1)"
// @Param limit query int false "Number of records per page (default: 10)"
// @Router /api/booking-user/search-requests [get]
func (h *BookingUserHandler) SearchRequests(c *gin.Context) {
	funcs.GetAuthenUser(c, h.Role)
	funcs.SearchRequests(c)
}

// GetRequest godoc
// @Summary Retrieve a specific booking request
// @Description This endpoint fetches details of a specific booking request using its unique identifier (TrnRequestUID).
// @Tags Booking-user
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Param id path string true "TrnRequestUID (trn_request_uid)"
// @Router /api/booking-user/request/{id} [get]
func (h *BookingUserHandler) GetRequest(c *gin.Context) {
	funcs.GetAuthenUser(c, h.Role)
	funcs.GetRequest(c)
}

// UpdateVehicleUser godoc
// @Summary Update vehicle information for a booking user
// @Description This endpoint allows a booking user to update the vehicle details associated with their request.
// @Tags Booking-user
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Param data body models.VmsTrnRequest_Update_VehicleUser true "VmsTrnRequest_Update_VehicleUser data"
// @Router /api/booking-user/update-vehicle-user [put]
func (h *BookingUserHandler) UpdateVehicleUser(c *gin.Context) {
	user := funcs.GetAuthenUser(c, h.Role)

	var request models.VmsTrnRequest_Update_VehicleUser

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var existing models.VmsTrnRequest_Update_VehicleUser
	if err := config.DB.First(&existing, "trn_request_uid = ?", request.TrnRequestUID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Booking not found"})
		return
	}
	type Data_Update struct {
		models.VmsTrnRequest_Update_VehicleUser
		models.LogUpdate
	}
	logUpdate := models.LogUpdate{
		UpdatedAt: time.Now(),
		UpdatedBy: user.EmpID,
	}
	if err := config.DB.Model(&Data_Update{}).
		Where("trn_request_uid = ?", request.TrnRequestUID).
		Updates(Data_Update{
			VmsTrnRequest_Update_VehicleUser: request,
			LogUpdate:                        logUpdate,
		}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update"})
		return
	}
	var result Data_Update
	if err := config.DB.First(&result, "trn_request_uid = ?", request.TrnRequestUID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Booking not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Updated successfully", "result": result})
}

// UpdateTrip godoc
// @Summary Update trip details for a booking request
// @Description This endpoint allows a booking user to update the details of an existing trip associated with their request.
// @Tags Booking-user
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Param data body models.VmsTrnRequest_Update_Trip true "VmsTrnRequest_Update_Trip data"
// @Router /api/booking-user/update-trip [put]
func (h *BookingUserHandler) UpdateTrip(c *gin.Context) {
	user := funcs.GetAuthenUser(c, h.Role)

	var request models.VmsTrnRequest_Update_Trip

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var existing models.VmsTrnRequest_Update_Trip
	if err := config.DB.First(&existing, "trn_request_uid = ?", request.TrnRequestUID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Booking not found"})
		return
	}

	type Data_Update struct {
		models.VmsTrnRequest_Update_Trip
		models.LogUpdate
	}
	logUpdate := models.LogUpdate{
		UpdatedAt: time.Now(),
		UpdatedBy: user.EmpID,
	}
	if err := config.DB.Model(&Data_Update{}).
		Where("trn_request_uid = ?", request.TrnRequestUID).
		Updates(Data_Update{
			VmsTrnRequest_Update_Trip: request,
			LogUpdate:                 logUpdate,
		}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update"})
		return
	}
	var result Data_Update
	if err := config.DB.First(&result, "trn_request_uid = ?", request.TrnRequestUID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Booking not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Updated successfully", "result": result})
}

// UpdatePickup godoc
// @Summary Update pickup details for a booking request
// @Description This endpoint allows a booking user to update the pickup information for an existing booking.
// @Tags Booking-user
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Param data body models.VmsTrnRequest_Update_Pickup true "VmsTrnRequest_Update_Pickup data"
// @Router /api/booking-user/update-pickup [put]
func (h *BookingUserHandler) UpdatePickup(c *gin.Context) {
	user := funcs.GetAuthenUser(c, h.Role)
	var request models.VmsTrnRequest_Update_Pickup

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	type Data_Update struct {
		models.VmsTrnRequest_Update_Pickup
		models.LogUpdate
	}
	logUpdate := models.LogUpdate{
		UpdatedAt: time.Now(),
		UpdatedBy: user.FullName(),
	}
	if err := config.DB.Model(&Data_Update{}).
		Where("trn_request_uid = ?", request.TrnRequestUID).
		Updates(Data_Update{
			VmsTrnRequest_Update_Pickup: request,
			LogUpdate:                   logUpdate,
		}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update"})
		return
	}
	var result Data_Update
	if err := config.DB.First(&result, "trn_request_uid = ?", request.TrnRequestUID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Booking not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Updated successfully", "result": result})
}

// UpdateDocument godoc
// @Summary Update document details for a booking request
// @Description This endpoint allows a booking user to update the document associated with their booking request.
// @Tags Booking-user
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Param data body models.VmsTrnRequest_Update_Document true "VmsTrnRequest_Update_Document data"
// @Router /api/booking-user/update-document [put]
func (h *BookingUserHandler) UpdateDocument(c *gin.Context) {
	user := funcs.GetAuthenUser(c, h.Role)
	var request models.VmsTrnRequest_Update_Document

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var existing models.VmsTrnRequest_Update_Document
	if err := config.DB.First(&existing, "trn_request_uid = ?", request.TrnRequestUID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Booking not found"})
		return
	}

	type Data_Update struct {
		models.VmsTrnRequest_Update_Document
		models.LogUpdate
	}
	logUpdate := models.LogUpdate{
		UpdatedAt: time.Now(),
		UpdatedBy: user.EmpID,
	}
	if err := config.DB.Model(&Data_Update{}).
		Where("trn_request_uid = ?", request.TrnRequestUID).
		Updates(Data_Update{
			VmsTrnRequest_Update_Document: request,
			LogUpdate:                     logUpdate,
		}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update"})
		return
	}
	var result Data_Update
	if err := config.DB.First(&result, "trn_request_uid = ?", request.TrnRequestUID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Booking not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Updated successfully", "result": result})
}

// UpdateCost godoc
// @Summary Update cost details for a booking request
// @Description This endpoint allows a booking user to update the cost information for an existing booking request.
// @Tags Booking-user
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Param data body models.VmsTrnRequest_Update_Document true "VmsTrnRequest_Update_Document data"
// @Router /api/booking-user/update-cost [put]
func (h *BookingUserHandler) UpdateCost(c *gin.Context) {
	user := funcs.GetAuthenUser(c, h.Role)
	var request models.VmsTrnRequest_Update_Cost

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var existing models.VmsTrnRequest_Update_Cost
	if err := config.DB.First(&existing, "trn_request_uid = ?", request.TrnRequestUID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Booking not found"})
		return
	}

	type Data_Update struct {
		models.VmsTrnRequest_Update_Cost
		models.LogUpdate
	}
	logUpdate := models.LogUpdate{
		UpdatedAt: time.Now(),
		UpdatedBy: user.EmpID,
	}
	if err := config.DB.Model(&Data_Update{}).
		Where("trn_request_uid = ?", request.TrnRequestUID).
		Updates(Data_Update{
			VmsTrnRequest_Update_Cost: request,
			LogUpdate:                 logUpdate,
		}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update"})
		return
	}
	var result Data_Update
	if err := config.DB.First(&result, "trn_request_uid = ?", request.TrnRequestUID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Booking not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Updated successfully", "result": result})
}

// UpdateVehicleType godoc
// @Summary Update vehicle type for a booking request
// @Description This endpoint allows a booking user to update the vehicle type associated with their booking request.
// @Tags Booking-user
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Param data body models.VmsTrnRequest_Update_VehicleType true "VmsTrnRequest_Update_VehicleType data"
// @Router /api/booking-user/update-vehicle-type [put]
func (h *BookingUserHandler) UpdateVehicleType(c *gin.Context) {
	user := funcs.GetAuthenUser(c, h.Role)

	var request models.VmsTrnRequest_Update_VehicleType

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var existing models.VmsTrnRequest_Update_VehicleType
	if err := config.DB.First(&existing, "trn_request_uid = ?", request.TrnRequestUID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Booking not found"})
		return
	}
	type Data_Update struct {
		models.VmsTrnRequest_Update_VehicleType
		models.LogUpdate
	}
	logUpdate := models.LogUpdate{
		UpdatedAt: time.Now(),
		UpdatedBy: user.EmpID,
	}
	if err := config.DB.Model(&Data_Update{}).
		Where("trn_request_uid = ?", request.TrnRequestUID).
		Updates(Data_Update{
			VmsTrnRequest_Update_VehicleType: request,
			LogUpdate:                        logUpdate,
		}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update"})
		return
	}
	var result Data_Update
	if err := config.DB.First(&result, "trn_request_uid = ?", request.TrnRequestUID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Booking not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Updated successfully", "result": result})
}

// UpdateApproved godoc
// @Summary Update approval status for a booking request
// @Description This endpoint allows a booking user to update the approval status of their booking request.
// @Tags Booking-user
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Param data body models.VmsTrnRequest_Update_Approver true "VmsTrnRequest_Update_Approver data"
// @Router /api/booking-user/update-approved [put]
func (h *BookingUserHandler) UpdateApproved(c *gin.Context) {
	user := funcs.GetAuthenUser(c, h.Role)

	var request models.VmsTrnRequest_Update_Approver

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var existing models.VmsTrnRequest_Update_Approver
	if err := config.DB.First(&existing, "trn_request_uid = ?", request.TrnRequestUID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Booking not found"})
		return
	}

	type Data_Update struct {
		models.VmsTrnRequest_Update_Approver
		models.LogUpdate
	}
	logUpdate := models.LogUpdate{
		UpdatedAt: time.Now(),
		UpdatedBy: user.EmpID,
	}
	if err := config.DB.Model(&Data_Update{}).
		Where("trn_request_uid = ?", request.TrnRequestUID).
		Updates(Data_Update{
			VmsTrnRequest_Update_Approver: request,
			LogUpdate:                     logUpdate,
		}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update"})
		return
	}
	var result Data_Update
	if err := config.DB.First(&result, "trn_request_uid = ?", request.TrnRequestUID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Booking not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Updated successfully", "result": result})
}

// UpdateSendedBack godoc
// @Summary Update sended back status for an item
// @Description This endpoint allows users to update the sended back status of an item.
// @Tags Booking-user
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Param data body models.VmsTrnRequest_SendedBack true "VmsTrnRequest_SendedBack data"
// @Router /api/booking-user/update-sended-back [put]
func (h *BookingUserHandler) UpdateSendedBack(c *gin.Context) {
	user := funcs.GetAuthenUser(c, h.Role)
	var request models.VmsTrnRequest_SendedBack

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var existing models.VmsTrnRequest_SendedBack
	if err := config.DB.First(&existing, "trn_request_uid = ?", request.TrnRequestUID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Booking not found"})
		return
	}

	logUpdate := models.LogUpdate{
		UpdatedAt: time.Now(),
		UpdatedBy: user.EmpID,
	}

	if err := config.DB.Model(&models.VmsTrnRequest_SendedBack_Update{}).
		Where("trn_request_uid = ?", request.TrnRequestUID).
		Updates(models.VmsTrnRequest_SendedBack_Update{
			VmsTrnRequest_SendedBack:      request,
			RefRequestStatusCode:          "20", //
			SendedBackRequestEmpID:        user.EmpID,
			SendedBackRequestEmpName:      user.FullName(),
			SendedBackRequestDeptSAP:      user.DeptSAP,
			SendedBackRequestDeptSAPShort: user.DeptSAPShort,
			SendedBackRequestDeptSAPFull:  user.DeptSAPFull,
			LogUpdate:                     logUpdate,
		}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update"})
		return
	}
	var result models.VmsTrnRequest_SendedBack_Update
	if err := config.DB.First(&result, "trn_request_uid = ?", request.TrnRequestUID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Booking not found"})
		return
	}
	funcs.CreateTrnLog(result.TrnRequestUID,
		result.RefRequestStatusCode,
		result.SendedBackRequestReason,
		user.EmpID)
	c.JSON(http.StatusOK, gin.H{"message": "Updated successfully", "result": result})
}

// UpdateCanceled godoc
// @Summary Update cancel status for an item
// @Description This endpoint allows users to update the cancel status of an item.
// @Tags Booking-user
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Param data body models.VmsTrnRequest_Canceled true "VmsTrnRequest_Canceled data"
// @Router /api/booking-user/update-canceled [put]
func (h *BookingUserHandler) UpdateCanceled(c *gin.Context) {
	user := funcs.GetAuthenUser(c, h.Role)
	var request models.VmsTrnRequest_Canceled

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var existing models.VmsTrnRequest_Canceled
	if err := config.DB.First(&existing, "trn_request_uid = ?", request.TrnRequestUID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Booking not found"})
		return
	}
	logUpdate := models.LogUpdate{
		UpdatedAt: time.Now(),
		UpdatedBy: user.EmpID,
	}
	if err := config.DB.Model(&models.VmsTrnRequest_Canceled_Update{}).
		Where("trn_request_uid = ?", request.TrnRequestUID).
		Updates(models.VmsTrnRequest_Canceled_Update{
			VmsTrnRequest_Canceled:      request,
			RefRequestStatusCode:        "90", // ยกเลิกคำขอ
			CanceledRequestEmpID:        user.EmpID,
			CanceledRequestEmpName:      user.FullName(),
			CanceledRequestDeptSAP:      user.DeptSAP,
			CanceledRequestDeptSAPShort: user.DeptSAPShort,
			CanceledRequestDeptSAPFull:  user.DeptSAPFull,
			LogUpdate:                   logUpdate,
		}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update"})
		return
	}
	var result models.VmsTrnRequest_SendedBack_Update
	if err := config.DB.First(&result, "trn_request_uid = ?", request.TrnRequestUID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Booking not found"})
		return
	}
	funcs.CreateTrnLog(result.TrnRequestUID,
		result.RefRequestStatusCode,
		result.SendedBackRequestReason,
		user.EmpID)
	c.JSON(http.StatusOK, gin.H{"message": "Updated successfully", "result": result})
}
