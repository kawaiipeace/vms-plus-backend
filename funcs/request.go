package funcs

import (
	"fmt"
	"net/http"
	"sort"
	"strings"
	"time"
	"vms_plus_be/config"
	"vms_plus_be/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

var StatusNameMap = map[string]string{
	"20": "รออนุมัติ",
	"21": "ถูกตีกลับ",
	"30": "รอตรวจสอบ",
	"31": "ถูกตีกลับ",
	"40": "รออนุมัติ",
	"41": "ถูกตีกลับ",
	"50": "รอรับกุญแจ",
	"51": "รอรับยานพาหนะ",
	"60": "เดินทาง",
	"70": "รอตรวจสอบ",
	"71": "คืนยานพาหนะไม่สำเร็จ",
	"80": "เสร็จสิ้น",
	"90": "ยกเลิกคำขอ",
}

func MenuRequests(statusMenuMap map[string]string) ([]models.VmsTrnRequestSummary, error) {
	var summary []models.VmsTrnRequestSummary

	// Group the request counts by statusMenuMap
	groupedSummary := make(map[string]int)
	for key := range statusMenuMap {
		statusCodes := strings.Split(key, ",")
		var count int64
		if err := config.DB.Table("vms_trn_request").
			Where("ref_request_status_code IN ?", statusCodes).
			Count(&count).Error; err != nil {
			return nil, err
		}
		groupedSummary[key] += int(count)
	}

	// Build the summary from the grouped data
	for key, count := range groupedSummary {
		summary = append(summary, models.VmsTrnRequestSummary{
			RefRequestStatusCode: key,
			RefRequestStatusName: statusMenuMap[key],
			Count:                count,
		})
	}
	// Sort the summary by RefRequestStatusCode
	sort.Slice(summary, func(i, j int) bool {
		return summary[i].RefRequestStatusCode < summary[j].RefRequestStatusCode
	})
	return summary, nil
}

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
		query = query.Where("req.request_no ILIKE ? OR req.vehicle_license_plate ILIKE ? OR req.vehicle_user_emp_name ILIKE ? OR req.work_place ILIKE ?", "%"+search+"%", "%"+search+"%", "%"+search+"%", "%"+search+"%")
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
	var requests []models.VmsTrnRequestResponse
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

func GetRequest(c *gin.Context, statusNameMap map[string]string) (models.VmsTrnRequestResponse, error) {
	id := c.Param("trn_request_uid")
	var request models.VmsTrnRequestResponse
	trnRequestUID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid TrnRequestUID"})
		return request, err
	}

	if err := config.DB.
		Preload("VmsMasVehicle.RefFuelType").
		Preload("VMSMasDriver").
		Preload("RefRequestStatus").
		Preload("RequestVehicleType").
		First(&request, "trn_request_uid = ?", trnRequestUID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Request not found"})
		return request, err
	}
	if request.VMSMasDriver.DriverBirthdate != (time.Time{}) {
		request.VMSMasDriver.Age = request.VMSMasDriver.CalculateAgeInYearsMonths()
	}
	request.NumberOfAvailableDrivers = 2
	request.DriverImageURL = config.DefaultAvatarURL
	request.CanCancelRequest = true
	request.IsUseDriver = request.MasCarpoolDriverUID != ""
	request.RefRequestStatusName = statusNameMap[request.RefRequestStatusCode]
	UpdateTrnRequestData(request.TrnRequestUID)
	//c.JSON(http.StatusOK, request)
	return request, nil
}

func GetRequestVehicelInUse(c *gin.Context, statusNameMap map[string]string) (models.VmsTrnRequestVehicleInUseResponse, error) {
	id := c.Param("trn_request_uid")
	var request models.VmsTrnRequestVehicleInUseResponse
	trnRequestUID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid TrnRequestUID"})
		return request, err
	}

	if err := config.DB.
		Preload("VmsMasVehicle.RefFuelType").
		Preload("VMSMasDriver").
		Preload("RefRequestStatus").
		Preload("RequestVehicleType").
		Preload("VehicleImagesReceived").
		Preload("VehicleImagesReturned").
		Preload("VehicleImagesReturned").
		Preload("VehicleImageInspect").
		Preload("ReceiverKeyTypeDetail").
		Preload("SatisfactionSurveyAnswers.SatisfactionSurveyQuestions").
		First(&request, "trn_request_uid = ?", trnRequestUID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Request not found"})
		return request, err
	}
	if request.VMSMasDriver.DriverBirthdate != (time.Time{}) {
		request.VMSMasDriver.Age = request.VMSMasDriver.CalculateAgeInYearsMonths()
	}
	request.NumberOfAvailableDrivers = 2
	request.DriverImageURL = config.DefaultAvatarURL
	request.ReceivedKeyImageURL = config.DefaultAvatarURL
	request.CanCancelRequest = true
	request.IsUseDriver = request.MasCarpoolDriverUID != ""
	request.RefRequestStatusName = StatusNameMap[request.RefRequestStatusCode]
	request.FleetCardNo = request.VmsMasVehicle.VehicleDepartment.FleetCardNo
	for i := range request.SatisfactionSurveyAnswers {
		desc := request.SatisfactionSurveyAnswers[i].SatisfactionSurveyQuestions.MasSatisfactionSurveyQuestionsDesc
		parts := strings.SplitN(desc, ":", 2)
		request.SatisfactionSurveyAnswers[i].SatisfactionSurveyQuestions.MasSatisfactionSurveyQuestionsTitle = parts[0] // Title before colon
		if len(parts) > 1 {
			request.SatisfactionSurveyAnswers[i].SatisfactionSurveyQuestions.MasSatisfactionSurveyQuestionsDesc = parts[1] // Remaining description after colon
		} else {
			request.SatisfactionSurveyAnswers[i].SatisfactionSurveyQuestions.MasSatisfactionSurveyQuestionsDesc = "" // Empty if no colon found
		}
	}
	if err := config.DB.
		Preload("RefTripType").
		Where("trn_request_uid <> ?", request.TrnRequestUID).
		Order("created_at DESC").
		First(&request.NextRequest).Error; err == nil {
		request.NextRequest.RefRequestStatusName = StatusNameMap[request.NextRequest.RefRequestStatusCode]
	}
	request.MileUsed = request.MileEnd - request.MileStart
	if err := config.DB.
		Table("vms_trn_add_fuel").
		Where("trn_request_uid = ? AND is_deleted = '0'", request.TrnRequestUID).
		Count(&request.AddFuelsCount).Error; err != nil {
		request.AddFuelsCount = 0
	}
	if err := config.DB.
		Table("vms_trn_trip_detail").
		Where("trn_request_uid = ? AND is_deleted = '0'", request.TrnRequestUID).
		Count(&request.TripDetailsCount).Error; err != nil {
		request.TripDetailsCount = 0
	}

	request.IsReturnOverDue = false
	if time.Now().Truncate(24 * time.Hour).After(request.EndDateTime.Truncate(24 * time.Hour)) {
		request.IsReturnOverDue = true
	}

	//UpdateTrnRequestData(request.TrnRequestUID)
	//c.JSON(http.StatusOK, request)
	return request, nil
}
