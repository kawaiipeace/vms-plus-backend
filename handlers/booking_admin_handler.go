package handlers

import (
	"fmt"
	"net/http"
	"time"
	"vms_plus_be/config"
	"vms_plus_be/funcs"
	"vms_plus_be/models"

	"github.com/gin-gonic/gin"
)

type BookingAdminHandler struct {
	Role string
}

// SearchRequests godoc
// @Summary Search booking requests and get summary counts by request status code
// @Description Search for requests using a keyword and get the summary of counts grouped by request status code
// @Tags Booking-admin
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
// @Router /api/booking-admin/search-requests [get]
func (h *BookingAdminHandler) SearchRequests(c *gin.Context) {
	funcs.GetAuthenUser(c, h.Role)
	//funcs.GetAuthenUser(c, h.Role)
	statusNameMap := map[string]string{
		"30": "รออนุมัติ",
		"31": "ตีกลับ",
		"40": "อนุมัติแล้ว",
		"90": "ยกเลิกคำขอ",
	}

	var requests []models.VmsTrnRequest_List
	var summary []models.VmsTrnRequest_Summary

	// Use the keys from statusNameMap as the list of valid status codes
	statusCodes := make([]string, 0, len(statusNameMap))
	for code := range statusNameMap {
		statusCodes = append(statusCodes, code)
	}

	// Build the main query
	query := config.DB.Table("public.vms_trn_request AS req").
		Select("req.*, status.ref_request_status_desc").
		Joins("LEFT JOIN public.vms_ref_request_status AS status ON req.ref_request_status_code = status.ref_request_status_code").
		Where("req.ref_request_status_code IN (?)", statusCodes)

	// Apply additional filters (search, date range, etc.)
	if search := c.Query("search"); search != "" {
		query = query.Where("req.request_no LIKE ? OR req.vehicle_license_plate LIKE ? OR req.vehicle_user_emp_name LIKE ? OR req.work_place LIKE ?", "%"+search+"%", "%"+search+"%", "%"+search+"%", "%"+search+"%")
	}
	if startDate := c.Query("startdate"); startDate != "" {
		query = query.Where("req.start_datetime >= ?", startDate)
	}
	if endDate := c.Query("enddate"); endDate != "" {
		query = query.Where("req.start_datetime <= ?", endDate)
	}

	// Ordering
	orderBy := c.Query("order_by")
	orderDir := c.Query("order_dir")
	if orderDir != "desc" {
		orderDir = "asc"
	}
	switch orderBy {
	case "request_no":
		query = query.Order("req.request_no " + orderDir)
	case "start_datetime":
		query = query.Order("req.start_datetime " + orderDir)
	case "ref_request_status_code":
		query = query.Order("req.ref_request_status_code " + orderDir)
	}

	// Pagination
	page := c.DefaultQuery("page", "1")
	limit := c.DefaultQuery("limit", "10")
	var pageInt, pageSizeInt int
	fmt.Sscanf(page, "%d", &pageInt)
	fmt.Sscanf(limit, "%d", &pageSizeInt)
	if pageInt < 1 {
		pageInt = 1
	}
	if pageSizeInt < 1 {
		pageSizeInt = 10
	}
	offset := (pageInt - 1) * pageSizeInt
	var total int64
	if err := query.Count(&total).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	query = query.Offset(offset).Limit(pageSizeInt)

	// Execute the main query
	if err := query.Scan(&requests).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	for i := range requests {
		requests[i].RefRequestStatusName = statusNameMap[requests[i].RefRequestStatusCode]
	}

	// Build the summary query
	summaryQuery := config.DB.Table("public.vms_trn_request AS req").
		Select("req.ref_request_status_code, COUNT(*) as count").
		Where("req.ref_request_status_code IN (?)", statusCodes).
		Group("req.ref_request_status_code")

	// Execute the summary query
	dbSummary := []struct {
		RefRequestStatusCode string `gorm:"column:ref_request_status_code"`
		Count                int    `gorm:"column:count"`
	}{}
	if err := summaryQuery.Scan(&dbSummary).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Create a complete summary with all statuses from statusNameMap
	for code, name := range statusNameMap {
		count := 0
		for _, dbItem := range dbSummary {
			if dbItem.RefRequestStatusCode == code {
				count = dbItem.Count
				break
			}
		}
		summary = append(summary, models.VmsTrnRequest_Summary{
			RefRequestStatusCode: code,
			RefRequestStatusName: name,
			Count:                count,
		})
	}

	// Return both the filtered requests and the complete summary
	c.JSON(http.StatusOK, gin.H{
		"pagination": gin.H{
			"total":      total,
			"page":       page,
			"limit":      pageSizeInt,
			"totalPages": (total + int64(pageSizeInt) - 1) / int64(pageSizeInt), // Calculate total pages
		},
		"requests": requests,
		"summary":  summary,
	})
}

// GetRequest godoc
// @Summary Retrieve a specific booking request
// @Description This endpoint fetches details of a specific booking request using its unique identifier (TrnRequestUID).
// @Tags Booking-admin
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Param id path string true "TrnRequestUID (trn_request_uid)"
// @Router /api/booking-admin/request/{id} [get]
func (h *BookingAdminHandler) GetRequest(c *gin.Context) {
	funcs.GetAuthenUser(c, h.Role)
	funcs.GetRequest(c)
}

// UpdateSendedBack godoc
// @Summary Update sended back status for an item
// @Description This endpoint allows users to update the sended back status of an item.
// @Tags Booking-admin
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Param data body models.VmsTrnRequest_SendedBack true "VmsTrnRequest_SendedBack data"
// @Router /api/booking-admin/update-sended-back [put]
func (h *BookingAdminHandler) UpdateSendedBack(c *gin.Context) {
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
			RefRequestStatusCode:          "31", // ตรวจสอบไม่ผ่าน รอแก้ไขคำขอ
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
// @Summary Update Approved status for an item
// @Description This endpoint allows users to update the sended back status of an item.
// @Tags Booking-admin
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Param data body models.VmsTrnRequest_Approved true "VmsTrnRequest_Approved data"
// @Router /api/booking-admin/update-approved [put]
func (h *BookingAdminHandler) UpdateApproved(c *gin.Context) {
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
			RefRequestStatusCode:        "40", // ตรวจสอบผ่านแล้ว รออนุมัติคำขอ
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
// @Tags Booking-admin
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Param data body models.VmsTrnRequest_Canceled true "VmsTrnRequest_Canceled data"
// @Router /api/booking-admin/update-canceled [put]
func (h *BookingAdminHandler) UpdateCanceled(c *gin.Context) {
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

// UpdateVehicleUser godoc
// @Summary Update vehicle information for a booking user
// @Description This endpoint allows a booking user to update the vehicle details associated with their request.
// @Tags Booking-admin
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Param data body models.VmsTrnRequest_Update_VehicleUser true "VmsTrnRequest_Update_VehicleUser data"
// @Router /api/booking-admin/update-vehicle-user [put]
func (h *BookingAdminHandler) UpdateVehicleUser(c *gin.Context) {
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
// @Tags Booking-admin
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Param data body models.VmsTrnRequest_Update_Trip true "VmsTrnRequest_Update_Trip data"
// @Router /api/booking-admin/update-trip [put]
func (h *BookingAdminHandler) UpdateTrip(c *gin.Context) {
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
// @Tags Booking-admin
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Param data body models.VmsTrnRequest_Update_Pickup true "VmsTrnRequest_Update_Pickup data"
// @Router /api/booking-admin/update-pickup [put]
func (h *BookingAdminHandler) UpdatePickup(c *gin.Context) {
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
		UpdatedBy: user.EmpID,
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
// @Tags Booking-admin
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Param data body models.VmsTrnRequest_Update_Document true "VmsTrnRequest_Update_Document data"
// @Router /api/booking-admin/update-document [put]
func (h *BookingAdminHandler) UpdateDocument(c *gin.Context) {
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
// @Tags Booking-admin
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Param data body models.VmsTrnRequest_Update_Document true "VmsTrnRequest_Update_Document data"
// @Router /api/booking-admin/update-cost [put]
func (h *BookingAdminHandler) UpdateCost(c *gin.Context) {
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

// UpdateDriver godoc
// @Summary Update driver for a booking request
// @Description This endpoint allows a booking user to update the cost information for an existing booking request.
// @Tags Booking-admin
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Param data body models.VmsTrnRequest_Update_Document true "VmsTrnRequest_Update_Document data"
// @Router /api/booking-admin/update-driver [put]
func (h *BookingAdminHandler) UpdateDriver(c *gin.Context) {
	user := funcs.GetAuthenUser(c, h.Role)
	var request models.VmsTrnRequest_Update_Driver

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var existing models.VmsTrnRequest_Update_Driver
	if err := config.DB.First(&existing, "trn_request_uid = ?", request.TrnRequestUID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Booking not found"})
		return
	}
	type Data_Update struct {
		models.VmsTrnRequest_Update_Driver
		models.LogUpdate
	}
	logUpdate := models.LogUpdate{
		UpdatedAt: time.Now(),
		UpdatedBy: user.EmpID,
	}
	if err := config.DB.Model(&Data_Update{}).
		Where("trn_request_uid = ?", request.TrnRequestUID).
		Updates(Data_Update{
			VmsTrnRequest_Update_Driver: request,
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

// UpdateVehicle godoc
// @Summary Update vehicle for a booking request
// @Description This endpoint allows a booking user to update the cost information for an existing booking request.
// @Tags Booking-admin
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Param data body models.VmsTrnRequest_Update_Document true "VmsTrnRequest_Update_Document data"
// @Router /api/booking-admin/update-vehicle [put]
func (h *BookingAdminHandler) UpdateVehicle(c *gin.Context) {
	user := funcs.GetAuthenUser(c, h.Role)
	var request models.VmsTrnRequest_Update_Vehicle

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var existing models.VmsTrnRequest_Update_Vehicle
	if err := config.DB.First(&existing, "trn_request_uid = ?", request.TrnRequestUID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Booking not found"})
		return
	}
	type Data_Update struct {
		models.VmsTrnRequest_Update_Vehicle
		models.LogUpdate
	}
	logUpdate := models.LogUpdate{
		UpdatedAt: time.Now(),
		UpdatedBy: user.EmpID,
	}
	if err := config.DB.Model(&Data_Update{}).
		Where("trn_request_uid = ?", request.TrnRequestUID).
		Updates(Data_Update{
			VmsTrnRequest_Update_Vehicle: request,
			LogUpdate:                    logUpdate,
		}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update"})
		return
	}
	if err := funcs.UpdateTrnRequestData(request.TrnRequestUID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update data request"})
	}
	var result Data_Update
	if err := config.DB.First(&result, "trn_request_uid = ?", request.TrnRequestUID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Booking not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Updated successfully", "result": result})
}
