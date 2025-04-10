package funcs

import (
	"fmt"
	"net/http"
	"strings"
	"time"
	"vms_plus_be/config"
	"vms_plus_be/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func SearchRequests(c *gin.Context) {
	var requests []struct {
		TrnRequestUID        string `gorm:"column:trn_request_uid;type:uuid;" json:"trn_request_uid"`
		RequestNo            string `gorm:"column:request_no" json:"request_no"`
		VehicleUserEmpID     string `gorm:"column:vehicle_user_emp_id" json:"vehicle_user_emp_id"`
		VehicleUserEmpName   string `gorm:"column:vehicle_user_emp_name" json:"vehicle_user_emp_name"`
		VehicleLicensePlate  string `gorm:"column:vehicle_license_plate" json:"vehicle_license_plate"`
		WorkPlace            string `gorm:"column:work_place" json:"work_place"`
		StartDatetime        string `gorm:"column:start_datetime" json:"start_datetime"`
		EndDatetime          string `gorm:"column:end_datetime" json:"end_datetime"`
		RefRequestStatusCode string `gorm:"column:ref_request_status_code" json:"ref_request_status_code"`
		RefRequestStatusDesc string `gorm:"column:ref_request_status_desc" json:"ref_request_status_desc"`
	}
	var summary []struct {
		RefRequestStatusCode string `gorm:"column:ref_request_status_code" json:"ref_request_status_code"`
		RefRequestStatusDesc string `gorm:"column:ref_request_status_name_1" json:"ref_request_status_desc"`
		Count                int    `gorm:"column:count" json:"count"`
	}

	query := config.DB.Table("public.vms_trn_request AS req").
		Select("req.*, status.ref_request_status_desc").
		Joins("LEFT JOIN public.vms_ref_request_status AS status ON req.ref_request_status_code = status.ref_request_status_code")

	// Apply filters to both the main query and summary query
	if search := c.Query("search"); search != "" {
		query = query.Where("req.request_no LIKE ? OR req.vehicle_license_plate LIKE ? OR req.vehicle_user_emp_name LIKE ? OR req.work_place LIKE ?", "%"+search+"%", "%"+search+"%", "%"+search+"%", "%"+search+"%")
	}

	// Filter by ref_request_status_code
	if statusCodes := c.Query("ref_request_status_code"); statusCodes != "" {
		statusCodeList := strings.Split(statusCodes, ",") // Split by comma
		query = query.Where("req.ref_request_status_code IN (?)", statusCodeList)
	}

	// Filter by date range
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

	// Execute the main search query
	if err := query.Scan(&requests).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// **Build Summary Query with Same Filters**
	summaryQuery := config.DB.Table("public.vms_trn_request AS req").
		Select("req.ref_request_status_code, status.ref_request_status_name_1, COUNT(*) as count").
		Joins("LEFT JOIN public.vms_ref_request_status AS status ON req.ref_request_status_code = status.ref_request_status_code")

	// Grouping to get count per status
	summaryQuery = summaryQuery.Group("req.ref_request_status_code, status.ref_request_status_name_1")

	// Execute summary query
	if err := summaryQuery.Scan(&summary).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return both the filtered requests and the summary
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

func ListRequest(c *gin.Context) {
	var requests []models.VmsTrnRequest_Response
	if err := config.DB.
		Preload("VmsMasVehicle.RefFuelType").
		Preload("VMSMasDriver").
		Find(&requests).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Request not found"})
		return
	}
	for i := range requests {
		if requests[i].VMSMasDriver.DriverBirthdate != (time.Time{}) {
			requests[i].VMSMasDriver.Age = requests[i].VMSMasDriver.CalculateAgeInYearsMonths()
		}
	}
	c.JSON(http.StatusOK, requests)
}

func GetRequest(c *gin.Context) {
	TrnRequestUID := c.Param("id")
	parsedID, err := uuid.Parse(TrnRequestUID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid TrnRequestUID"})
		return
	}
	var request models.VmsTrnRequest_Response
	if err := config.DB.
		Preload("VmsMasVehicle.RefFuelType").
		Preload("VMSMasDriver").
		First(&request, "trn_request_uid = ?", parsedID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Request not found"})
		return
	}
	if request.VMSMasDriver.DriverBirthdate != (time.Time{}) {
		request.VMSMasDriver.Age = request.VMSMasDriver.CalculateAgeInYearsMonths()
	}
	request.NumberOfAvailableDrivers = 2

	request.CanCancelRequest = true
	c.JSON(http.StatusOK, request)
}
