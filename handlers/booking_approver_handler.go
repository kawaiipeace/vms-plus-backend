package handlers

import (
	"net/http"
	"time"
	"vms_plus_be/config"
	"vms_plus_be/funcs"
	"vms_plus_be/models"

	"github.com/gin-gonic/gin"
)

type BookingApproverHandler struct {
	Role string
}

// SearchRequests godoc
// @Summary Search booking requests and get summary counts by request status code
// @Description Search for requests using a keyword and get the summary of counts grouped by request status code
// @Tags Booking-approver
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Param search query string false "Search keyword (matches request_no, vehicle_license_plate, vehicle_user_emp_name, or work_place)"
// @Param ref_request_status_code query string false "Filter by multiple request status codes (comma-separated, e.g., 'A,B,C')"
// @Param startdate query string false "Filter by start datetime (YYYY-MM-DD format)"
// @Param enddate query string false "Filter by end datetime (YYYY-MM-DD format)"
// @Param order_by query string false "Order by request_no, start_datetime, ref_request_status_code"
// @Param order_dir query string false "Order direction: asc or desc"
// @Param page query int false "Page number (default: 1)"
// @Param page_size query int false "Number of records per page (default: 10)"
// @Router /api/booking-approver/search-requests [get]
func (h *BookingApproverHandler) SearchRequests(c *gin.Context) {
	funcs.GetAuthenUser(c, h.Role)
	funcs.SearchRequests(c)
}

// GetRequest godoc
// @Summary Retrieve a specific booking request
// @Description This endpoint fetches details of a specific booking request using its unique identifier (TrnRequestUID).
// @Tags Booking-approver
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Param id path string true "TrnRequestUID (trn_request_uid)"
// @Router /api/booking-approver/request/{id} [get]
func (h *BookingApproverHandler) GetRequest(c *gin.Context) {
	funcs.GetAuthenUser(c, h.Role)
	funcs.GetRequest(c)
}

// UpdateSendedBack godoc
// @Summary Update sended back status for an item
// @Description This endpoint allows users to update the sended back status of an item.
// @Tags Booking-approver
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Param data body models.VmsTrnRequest_SendedBack true "VmsTrnRequest_SendedBack data"
// @Router /api/booking-approver/update-sended-back [put]
func (h *BookingApproverHandler) UpdateSendedBack(c *gin.Context) {
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
			RefRequestStatusCode:          "21", // ผู้มีอำนาจไม่ยืนยัน ตีกลับคำขอ
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

// UpdateApproved godoc
// @Summary Update sended back status for an item
// @Description This endpoint allows users to update the sended back status of an item.
// @Tags Booking-approver
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Param data body models.VmsTrnRequest_Approved true "VmsTrnRequest_Approved data"
// @Router /api/booking-approver/update-approved [put]
func (h *BookingApproverHandler) UpdateApproved(c *gin.Context) {
	user := funcs.GetAuthenUser(c, h.Role)
	var request models.VmsTrnRequest_Approved

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var existing models.VmsTrnRequest_Approved
	if err := config.DB.First(&existing, "trn_request_uid = ?", request.TrnRequestUID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Booking not found"})
		return
	}
	logUpdate := models.LogUpdate{
		UpdatedAt: time.Now(),
		UpdatedBy: user.EmpID,
	}
	if err := config.DB.Model(&models.VmsTrnRequest_Approved_Update{}).
		Where("trn_request_uid = ?", request.TrnRequestUID).
		Updates(models.VmsTrnRequest_Approved_Update{
			VmsTrnRequest_Approved:      request,
			RefRequestStatusCode:        "30", // ยืนยันคำขอแล้ว รอตรวจสอบคำขอ
			ApprovedRequestEmpID:        user.EmpID,
			ApprovedRequestEmpName:      user.FullName(),
			ApprovedRequestDeptSAP:      user.DeptSAP,
			ApprovedRequestDeptSAPShort: user.DeptSAPShort,
			ApprovedRequestDeptSAPFull:  user.DeptSAPFull,
			LogUpdate:                   logUpdate,
		}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update"})
		return
	}
	var result models.VmsTrnRequest_Approved_Update
	if err := config.DB.First(&result, "trn_request_uid = ?", request.TrnRequestUID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Booking not found"})
		return
	}
	funcs.CreateTrnLog(result.TrnRequestUID,
		result.RefRequestStatusCode,
		"",
		user.EmpID)

	c.JSON(http.StatusOK, gin.H{"message": "Updated successfully", "result": result})
}

// UpdateCanceled godoc
// @Summary Update cancel status for an item
// @Description This endpoint allows users to update the cancel status of an item.
// @Tags Booking-approver
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Param data body models.VmsTrnRequest_Canceled true "VmsTrnRequest_Canceled data"
// @Router /api/booking-approver/update-canceled [put]
func (h *BookingApproverHandler) UpdateCanceled(c *gin.Context) {
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
	var result models.VmsTrnRequest_Canceled_Update
	if err := config.DB.First(&result, "trn_request_uid = ?", request.TrnRequestUID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Booking not found"})
		return
	}
	funcs.CreateTrnLog(result.TrnRequestUID,
		result.RefRequestStatusCode,
		result.CanceledRequestReason,
		user.EmpID)

	c.JSON(http.StatusOK, gin.H{"message": "Updated successfully", "result": result})
}
