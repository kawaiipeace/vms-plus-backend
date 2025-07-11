package handlers

import (
	"fmt"
	"net/http"
	"time"
	"vms_plus_be/config"
	"vms_plus_be/funcs"
	"vms_plus_be/messages"
	"vms_plus_be/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type DriverLicenseUserHandler struct {
	Role string
}

func (h *DriverLicenseUserHandler) SetQueryRole(user *models.AuthenUserEmp, query *gorm.DB) *gorm.DB {
	return query.Where("created_request_emp_id = ?", user.EmpID)
}

func (h *DriverLicenseUserHandler) SetQueryRoleDept(user *models.AuthenUserEmp, query *gorm.DB) *gorm.DB {
	return query
}
func (h *DriverLicenseUserHandler) SetQueryStatusCanUpdate(query *gorm.DB) *gorm.DB {
	return query.Where("ref_request_status_code in ('11') and is_deleted = '0'")
}

var StatusDriverAnnualLicense = map[string]string{
	"00": "ไม่มี",
	"10": "รออนุมัติ",
	"11": "ตีกลับ",
	"20": "รออนุมัติ",
	"21": "ตีกลับ",
	"30": "อนุมัติแล้ว",
	"31": "มีผลปีถัดไป",
	"80": "หมดอายุ",
	"90": "ยกเลิก",
}

func GetProgressRequestHistory(request models.VmsDriverLicenseAnnualResponse) []models.ProgressRequestHistory {
	var progressRequestHistory []models.ProgressRequestHistory
	if request.RefRequestAnnualDriverStatusCode >= "10" {
		progressRequestHistory = append(progressRequestHistory, models.ProgressRequestHistory{
			ProgressIcon:     "3",
			ProgressName:     "ขออนุมัติ",
			ProgressDatetime: models.TimeWithZone{Time: request.CreatedRequestDatetime.Time},
		})
	}
	if request.RefRequestAnnualDriverStatusCode >= "20" {
		progressRequestHistory = append(progressRequestHistory, models.ProgressRequestHistory{
			ProgressIcon:     "3",
			ProgressName:     "อนุมัติจากต้นสังกัด",
			ProgressDatetime: models.TimeWithZone{Time: request.ConfirmedRequestDatetime.Time},
		})
	}
	if request.RefRequestAnnualDriverStatusCode >= "30" {
		progressRequestHistory = append(progressRequestHistory, models.ProgressRequestHistory{
			ProgressIcon:     "3",
			ProgressName:     "อนุมัติให้ทำหน้าที่ขับรถยนต์",
			ProgressDatetime: models.TimeWithZone{Time: request.ApprovedRequestDatetime.Time},
		})
	}
	return progressRequestHistory
}

// GetLicenseCard godoc
// @Summary Retrieve driver's license card
// @Description Get the driver's license card details
// @Tags Driver-license-user
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Router /api/driver-license-user/card [get]
func (h *DriverLicenseUserHandler) GetLicenseCard(c *gin.Context) {
	user := funcs.GetAuthenUser(c, h.Role)
	if c.IsAborted() {
		return
	}
	driver := models.VmsDriverLicenseCard{
		EmpID:          user.EmpID,
		DriverName:     user.FullName,
		DeptSAPShort:   user.DeptSAPShort,
		IsNoExpiryDate: false,
	}

	annualYear := time.Now().Year() + 543

	var license models.VmsDriverLicenseAnnualResponse
	query := config.DB.Where("created_request_emp_id = ? and is_deleted = ? and annual_yyyy = ? and ref_request_annual_driver_status_code <> ?", user.EmpID, "0", annualYear, "90")
	err := query.Preload("DriverLicenseType").
		Preload("DriverCertificateType").
		Order("ref_request_annual_driver_status_code").
		First(&license).Error

	if err == nil {
		driver.LicenseStatusCode = license.RefRequestAnnualDriverStatusCode
		driver.LicenseStatus = StatusDriverAnnualLicense[driver.LicenseStatusCode]
	} else {
		driver.LicenseStatusCode = "00"
		driver.LicenseStatus = StatusDriverAnnualLicense[driver.LicenseStatusCode]
	}
	driver.TrnRequestAnnualDriverUID = license.TrnRequestAnnualDriverUID
	driver.RequestAnnualDriverNo = license.RequestAnnualDriverNo
	driver.AnnualYYYY = license.AnnualYYYY
	driver.DriverLicense = models.VmsDriverLicenseCardLicense{
		EmpID:                    license.CreatedRequestEmpID,
		DriverLicenseNo:          license.DriverLicenseNo,
		RefDriverLicenseTypeCode: license.RefDriverLicenseTypeCode,
		DriverLicenseStartDate:   license.ApprovedRequestDatetime,
		DriverLicenseImage:       license.DriverLicenseImg,
		DriverLicenseType:        license.DriverLicenseType,
	}

	// Calculate the end of the year for driver.AnnualYYYY - 543
	driver.DriverLicense.DriverLicenseEndDate = license.RequestExpireDate
	driver.DriverCertificate = models.VmsDriverLicenseCardCertificate{
		EmpID:                       license.CreatedRequestEmpID,
		DriverCertificateNo:         license.DriverCertificateNo,
		DriverCertificateName:       license.DriverCertificateName,
		DriverCertificateIssueDate:  license.DriverCertificateIssueDate,
		DriverCertificateExpireDate: license.DriverCertificateExpireDate,
		DriverCertificateImg:        license.DriverCertificateImg,
		DriverCertificateTypeCode:   license.DriverCertificateTypeCode,
		DriverCertificateType:       license.DriverCertificateType,
	}
	driver.ProgressRequestHistory = GetProgressRequestHistory(license)

	//next annual
	var licenseNext models.VmsDriverLicenseAnnualResponse
	queryNext := config.DB.Where("created_request_emp_id = ? and is_deleted = ? and annual_yyyy = ? and ref_request_annual_driver_status_code <> ?", user.EmpID, "0", annualYear+1, "90")
	errNext := queryNext.Order("ref_request_annual_driver_status_code").
		First(&licenseNext).Error

	if errNext == nil {
		driver.NextTrnRequestAnnualDriverUID = licenseNext.TrnRequestAnnualDriverUID
		driver.NextAnnualYYYY = licenseNext.AnnualYYYY
		driver.NextLicenseStatusCode = licenseNext.RefRequestAnnualDriverStatusCode
		driver.NextLicenseStatus = StatusDriverAnnualLicense[driver.NextLicenseStatusCode]
	} else {
		driver.NextAnnualYYYY = annualYear + 1
		driver.NextLicenseStatusCode = "00"
		driver.NextLicenseStatus = StatusDriverAnnualLicense[driver.NextLicenseStatusCode]
	}

	//prev annual
	var licensePrev models.VmsDriverLicenseAnnualResponse
	queryPrev := config.DB.Where("created_request_emp_id = ? and is_deleted = ? and annual_yyyy = ? and ref_request_annual_driver_status_code <> ?", user.EmpID, "0", annualYear-1, "90")
	errPrev := queryPrev.Order("ref_request_annual_driver_status_code").
		First(&licensePrev).Error

	if errPrev == nil {
		driver.PrevTrnRequestAnnualDriverUID = licensePrev.TrnRequestAnnualDriverUID
		driver.PrevAnnualYYYY = licensePrev.AnnualYYYY
		driver.PrevLicenseStatusCode = licensePrev.RefRequestAnnualDriverStatusCode
		driver.PrevLicenseStatus = StatusDriverAnnualLicense[driver.PrevLicenseStatusCode]
	} else {
		driver.PrevAnnualYYYY = annualYear - 1
		driver.PrevLicenseStatusCode = "00"
		driver.PrevLicenseStatus = StatusDriverAnnualLicense[driver.PrevLicenseStatusCode]
	}

	c.JSON(http.StatusOK, gin.H{"driver": driver})
}

// CreateDriverLicenseAnnual godoc
// @Summary Create a new driver license annual record
// @Description This endpoint allows creating a new driver license annual record.
// @Tags Driver-license-user
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Param data body models.VmsDriverLicenseAnnualRequest true "VmsDriverLicenseAnnualRequest data"
// @Router /api/driver-license-user/create-license-annual [post]
func (h *DriverLicenseUserHandler) CreateDriverLicenseAnnual(c *gin.Context) {
	user := funcs.GetAuthenUser(c, h.Role)
	if c.IsAborted() {
		return
	}
	var request models.VmsDriverLicenseAnnualRequest
	var result struct {
		models.VmsDriverLicenseAnnualRequest
		RequestAnnualDriverNo     string `gorm:"column:request_annual_driver_no" json:"request_annual_driver_no"`
		TrnRequestAnnualDriverUID string `gorm:"column:trn_request_annual_driver_uid" json:"trn_request_annual_driver_uid"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON input", "message": messages.ErrInvalidJSONInput.Error()})
		return
	}

	checkQuery := config.DB.Where("created_request_emp_id = ? and annual_yyyy = ? and is_deleted = ? and ref_request_annual_driver_status_code <> ?", user.EmpID, request.AnnualYYYY, "0", "90")
	var checkRequest models.VmsDriverLicenseAnnualRequest
	err := checkQuery.First(&checkRequest).Error
	if err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Request already exists", "message": messages.ErrBadRequest.Error()})
		return
	}
	request.TrnRequestAnnualDriverUID = uuid.New().String()
	request.CreatedRequestEmpID = user.EmpID
	empUser := funcs.GetUserEmpInfo(request.CreatedRequestEmpID)
	request.CreatedRequestEmpName = empUser.FullName
	request.CreatedRequestEmpPosition = empUser.Position
	request.CreatedRequestDeptSap = empUser.DeptSAP
	request.CreatedRequestDeptSapNameShort = empUser.DeptSAPShort
	request.CreatedRequestDeptSapNameFull = empUser.DeptSAPFull
	request.CreatedRequestMobileNumber = empUser.TelMobile
	request.CreatedRequestPhoneNumber = empUser.TelInternal
	request.CreatedRequestDatetime = models.TimeWithZone{Time: time.Now()}

	confirmUser := funcs.GetUserEmpInfo(request.ConfirmedRequestEmpID)
	request.ConfirmedRequestEmpName = confirmUser.FullName
	request.ConfirmedRequestEmpPosition = confirmUser.Position
	request.ConfirmedRequestDeptSap = confirmUser.DeptSAP
	request.ConfirmedRequestDeptSapShort = confirmUser.DeptSAPShort
	request.ConfirmedRequestDeptSapFull = confirmUser.DeptSAPFull
	request.ConfirmedRequestMobileNumber = confirmUser.TelMobile
	request.ConfirmedRequestPhoneNumber = confirmUser.TelInternal

	approveUser := funcs.GetUserEmpInfo(request.ApprovedRequestEmpID)
	request.ApprovedRequestEmpName = approveUser.FullName
	request.ApprovedRequestEmpPosition = approveUser.Position
	request.ApprovedRequestDeptSap = approveUser.DeptSAP
	request.ApprovedRequestDeptSapShort = approveUser.DeptSAPShort
	request.ApprovedRequestDeptSapFull = approveUser.DeptSAPFull
	request.ApprovedRequestMobileNumber = approveUser.TelMobile
	request.ApprovedRequestPhoneNumber = approveUser.TelInternal

	request.RefRequestAnnualDriverStatusCode = "10"
	request.RejectedRequestEmpPosition = ""
	request.CanceledRequestEmpPosition = ""

	request.UpdatedAt = time.Now()
	request.UpdatedBy = user.EmpID

	var maxRequestNo string
	config.DB.Table("vms_trn_request_annual_driver").
		Select("MAX(request_annual_driver_no)").
		Where("request_annual_driver_no ILIKE ?", "RAD%").
		Scan(&maxRequestNo)

	var nextNumber int
	if maxRequestNo != "" {
		fmt.Sscanf(maxRequestNo, "RAD%d", &nextNumber)
	}
	nextNumber++

	request.RequestAnnualDriverNo = fmt.Sprintf("RAD%09d", nextNumber)

	if user.IsLevelM5 == "1" {
		request.RefRequestAnnualDriverStatusCode = "20"
	}

	if err := config.DB.Create(&request).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create driver license annual record", "message": messages.ErrInternalServer.Error()})
		return
	}
	if err := config.DB.First(&result, "trn_request_annual_driver_uid = ?", request.TrnRequestAnnualDriverUID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "annual not found", "message": messages.ErrNotfound.Error()})
		return
	}
	funcs.CreateRequestAnnualLicenseNotification(request.TrnRequestAnnualDriverUID)
	// Return success response
	c.JSON(http.StatusCreated, gin.H{"message": "Driver license annual record created successfully", "result": result})
}

// ResendDriverLicenseAnnual godoc
// @Summary Resend driver license annual record
// @Description This endpoint allows resending a driver license annual record.
// @Tags Driver-license-user
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Param trn_request_annual_driver_uid path string true "trnRequestAnnualDriverUID (trn_request_annual_driver_uid)"
// @Param data body models.VmsDriverLicenseAnnualRequest true "VmsDriverLicenseAnnualRequest data"
// @Router /api/driver-license-user/resend-license-annual/{trn_request_annual_driver_uid} [put]
func (h *DriverLicenseUserHandler) ResendDriverLicenseAnnual(c *gin.Context) {
	user := funcs.GetAuthenUser(c, h.Role)
	if c.IsAborted() {
		return
	}
	trnRequestAnnualDriverUID := c.Param("trn_request_annual_driver_uid")
	var request models.VmsDriverLicenseAnnualRequest
	var result struct {
		models.VmsDriverLicenseAnnualRequest
		RequestAnnualDriverNo     string `gorm:"column:request_annual_driver_no" json:"request_annual_driver_no"`
		TrnRequestAnnualDriverUID string `gorm:"column:trn_request_annual_driver_uid" json:"trn_request_annual_driver_uid"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON input", "message": messages.ErrInvalidJSONInput.Error()})
		return
	}
	var existsRequest models.VmsDriverLicenseAnnualRequest
	query := h.SetQueryRole(user, config.DB)
	query.First(&existsRequest, "trn_request_annual_driver_uid = ? AND is_deleted = ? AND ref_request_annual_driver_status_code in (?)", trnRequestAnnualDriverUID, "0", []string{"11", "21"})
	if existsRequest.TrnRequestAnnualDriverUID == "" {
		c.JSON(http.StatusNotFound, gin.H{"error": "annual not found", "message": messages.ErrNotfound.Error()})
		return
	}

	request.TrnRequestAnnualDriverUID = trnRequestAnnualDriverUID
	request.CreatedRequestEmpID = user.EmpID
	empUser := funcs.GetUserEmpInfo(request.CreatedRequestEmpID)
	request.CreatedRequestEmpName = empUser.FullName
	request.CreatedRequestEmpPosition = empUser.Position
	request.CreatedRequestDeptSap = empUser.DeptSAP
	request.CreatedRequestDeptSapNameShort = empUser.DeptSAPShort
	request.CreatedRequestDeptSapNameFull = empUser.DeptSAPFull
	request.CreatedRequestMobileNumber = empUser.TelMobile
	request.CreatedRequestPhoneNumber = empUser.TelInternal
	request.CreatedRequestDatetime = existsRequest.CreatedRequestDatetime

	confirmUser := funcs.GetUserEmpInfo(request.ConfirmedRequestEmpID)
	request.ConfirmedRequestEmpName = confirmUser.FullName
	request.ConfirmedRequestEmpPosition = confirmUser.Position
	request.ConfirmedRequestDeptSap = confirmUser.DeptSAP
	request.ConfirmedRequestDeptSapShort = confirmUser.DeptSAPShort
	request.ConfirmedRequestDeptSapFull = confirmUser.DeptSAPFull
	request.ConfirmedRequestMobileNumber = confirmUser.TelMobile
	request.ConfirmedRequestPhoneNumber = confirmUser.TelInternal

	approveUser := funcs.GetUserEmpInfo(request.ApprovedRequestEmpID)
	request.ApprovedRequestEmpName = approveUser.FullName
	request.ApprovedRequestEmpPosition = approveUser.Position
	request.ApprovedRequestDeptSap = approveUser.DeptSAP
	request.ApprovedRequestDeptSapShort = approveUser.DeptSAPShort
	request.ApprovedRequestDeptSapFull = approveUser.DeptSAPFull
	request.ApprovedRequestMobileNumber = approveUser.TelMobile
	request.ApprovedRequestPhoneNumber = approveUser.TelInternal

	request.RefRequestAnnualDriverStatusCode = "10"
	if user.IsLevelM5 == "1" {
		request.RefRequestAnnualDriverStatusCode = "20"
	}
	request.RejectedRequestEmpPosition = existsRequest.RejectedRequestEmpPosition
	request.CanceledRequestEmpPosition = existsRequest.CanceledRequestEmpPosition

	request.UpdatedAt = time.Now()
	request.UpdatedBy = user.EmpID
	request.RequestAnnualDriverNo = existsRequest.RequestAnnualDriverNo

	if err := config.DB.Save(&request).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to resend driver license annual record", "message": messages.ErrInternalServer.Error()})
		return
	}
	if err := config.DB.First(&result, "trn_request_annual_driver_uid = ?", request.TrnRequestAnnualDriverUID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "annual not found", "message": messages.ErrNotfound.Error()})
		return
	}
	funcs.CreateRequestAnnualLicenseNotification(request.TrnRequestAnnualDriverUID)
	// Return success response
	c.JSON(http.StatusOK, gin.H{"message": "Driver license annual record resend successfully", "result": result})
}

// GetDriverLicenseAnnual godoc
// @Summary Retrieve a specific driver license annual record
// @Description Get the details of a driver license annual record by its unique identifier
// @Tags Driver-license-user
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Param trn_request_annual_driver_uid path string true "trnRequestAnnualDriverUID (trn_request_annual_driver_uid)"
// @Router /api/driver-license-user/license-annual/{trn_request_annual_driver_uid} [get]
func (h *DriverLicenseUserHandler) GetDriverLicenseAnnual(c *gin.Context) {
	user := funcs.GetAuthenUser(c, h.Role)
	if c.IsAborted() {
		return
	}
	trnRequestAnnualDriverUID := c.Param("trn_request_annual_driver_uid")
	var request models.VmsDriverLicenseAnnualResponse

	query := h.SetQueryRole(user, config.DB)
	if err := query.
		Preload("DriverLicenseType").
		Preload("DriverCertificateType").
		First(&request, "trn_request_annual_driver_uid = ? and is_deleted = ?", trnRequestAnnualDriverUID, "0").Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "annual not found", "message": messages.ErrNotfound.Error()})
		return
	}
	if request.RefRequestAnnualDriverStatusCode == "10" {
		request.ProgressRequestStatus = []models.ProgressRequestStatus{
			{ProgressIcon: "3", ProgressName: "ขออนุมัติ", ProgressDatetime: request.CreatedRequestDatetime},
			{ProgressIcon: "1", ProgressName: "รอต้นสังกัดตรวจสอบ", ProgressDatetime: request.ConfirmedRequestDatetime},
			{ProgressIcon: "0", ProgressName: "รออนุมัติให้ทำหน้าที่ขับรถประจำปี", ProgressDatetime: request.ApprovedRequestDatetime},
		}
		request.ProgressRequestStatusEmp = models.ProgressRequestStatusEmp{
			ActionRole:   "ผู้อนุมัติต้นสังกัด",
			EmpID:        request.ConfirmedRequestEmpID,
			EmpName:      request.ConfirmedRequestEmpName,
			EmpPosition:  request.ConfirmedRequestEmpPosition,
			DeptSAP:      request.ConfirmedRequestDeptSap,
			DeptSAPShort: request.ConfirmedRequestDeptSapShort,
			DeptSAPFull:  request.ConfirmedRequestDeptSapFull,
			PhoneNumber:  request.ConfirmedRequestPhoneNumber,
			MobileNumber: request.ConfirmedRequestMobileNumber,
		}
	}
	if request.RefRequestAnnualDriverStatusCode == "11" {
		request.ProgressRequestStatus = []models.ProgressRequestStatus{
			{ProgressIcon: "3", ProgressName: "ขออนุมัติ", ProgressDatetime: request.CreatedRequestDatetime},
			{ProgressIcon: "2", ProgressName: "ตีกลับจากต้นสังกัด", ProgressDatetime: request.RejectedRequestDatetime},
			{ProgressIcon: "0", ProgressName: "รออนุมัติให้ทำหน้าที่ขับรถประจำปี", ProgressDatetime: request.ApprovedRequestDatetime},
		}
		request.ProgressRequestStatusEmp = models.ProgressRequestStatusEmp{
			ActionRole:   "ผู้อนุมัติต้นสังกัด",
			EmpID:        request.ConfirmedRequestEmpID,
			EmpName:      request.ConfirmedRequestEmpName,
			EmpPosition:  request.ConfirmedRequestEmpPosition,
			DeptSAP:      request.ConfirmedRequestDeptSap,
			DeptSAPShort: request.ConfirmedRequestDeptSapShort,
			DeptSAPFull:  request.ConfirmedRequestDeptSapFull,
			PhoneNumber:  request.ConfirmedRequestPhoneNumber,
			MobileNumber: request.ConfirmedRequestMobileNumber,
		}
	}
	if request.RefRequestAnnualDriverStatusCode == "20" {
		request.ProgressRequestStatus = []models.ProgressRequestStatus{
			{ProgressIcon: "3", ProgressName: "ขออนุมัติ", ProgressDatetime: request.CreatedRequestDatetime},
			{ProgressIcon: "3", ProgressName: "ต้นสังกัดตรวจสอบ", ProgressDatetime: request.ConfirmedRequestDatetime},
			{ProgressIcon: "1", ProgressName: "รออนุมัติให้ทำหน้าที่ขับรถประจำปี", ProgressDatetime: request.ApprovedRequestDatetime},
		}
		request.ProgressRequestStatusEmp = models.ProgressRequestStatusEmp{
			ActionRole:   "ผู้อนุมัติให้ทำหน้าที่ขับรถประจำปี",
			EmpID:        request.ApprovedRequestEmpID,
			EmpName:      request.ApprovedRequestEmpName,
			EmpPosition:  request.ApprovedRequestEmpPosition,
			DeptSAP:      request.ApprovedRequestDeptSap,
			DeptSAPShort: request.ApprovedRequestDeptSapShort,
			DeptSAPFull:  request.ApprovedRequestDeptSapFull,
			PhoneNumber:  request.ApprovedRequestPhoneNumber,
			MobileNumber: request.ApprovedRequestMobileNumber,
		}
	}
	if request.RefRequestAnnualDriverStatusCode == "21" {
		request.ProgressRequestStatus = []models.ProgressRequestStatus{
			{ProgressIcon: "3", ProgressName: "ขออนุมัติ", ProgressDatetime: request.CreatedRequestDatetime},
			{ProgressIcon: "3", ProgressName: "ต้นสังกัดตรวจสอบ", ProgressDatetime: request.ConfirmedRequestDatetime},
			{ProgressIcon: "2", ProgressName: "ตีกลับจากผู้อนุมัติ", ProgressDatetime: request.RejectedRequestDatetime},
		}
		request.ProgressRequestStatusEmp = models.ProgressRequestStatusEmp{
			ActionRole:   "ผู้อนุมัติต้นสังกัด",
			EmpID:        request.ConfirmedRequestEmpID,
			EmpName:      request.ConfirmedRequestEmpName,
			EmpPosition:  request.ConfirmedRequestEmpPosition,
			DeptSAP:      request.ConfirmedRequestDeptSap,
			DeptSAPShort: request.ConfirmedRequestDeptSapShort,
			DeptSAPFull:  request.ConfirmedRequestDeptSapFull,
			PhoneNumber:  request.ConfirmedRequestPhoneNumber,
			MobileNumber: request.ConfirmedRequestMobileNumber,
		}
	}
	if request.RefRequestAnnualDriverStatusCode == "30" {
		request.ProgressRequestStatus = []models.ProgressRequestStatus{
			{ProgressIcon: "3", ProgressName: "ขออนุมัติ", ProgressDatetime: request.CreatedRequestDatetime},
			{ProgressIcon: "3", ProgressName: "ต้นสังกัดตรวจสอบ", ProgressDatetime: request.ConfirmedRequestDatetime},
			{ProgressIcon: "3", ProgressName: "อนุมัติให้ทำหน้าที่ขับรถประจำปี", ProgressDatetime: request.ApprovedRequestDatetime},
		}
		request.ProgressRequestStatusEmp = models.ProgressRequestStatusEmp{
			ActionRole:   "ผู้อนุมัติให้ทำหน้าที่ขับรถประจำปี",
			EmpID:        request.ApprovedRequestEmpID,
			EmpName:      request.ApprovedRequestEmpName,
			EmpPosition:  request.ApprovedRequestEmpPosition,
			DeptSAP:      request.ApprovedRequestDeptSap,
			DeptSAPShort: request.ApprovedRequestDeptSapShort,
			DeptSAPFull:  request.ApprovedRequestDeptSapFull,
			PhoneNumber:  request.ApprovedRequestPhoneNumber,
			MobileNumber: request.ApprovedRequestMobileNumber,
		}
	}

	if request.RefRequestAnnualDriverStatusCode == "90" && request.CanceledRequestEmpID == request.CreatedRequestEmpID {
		request.ProgressRequestStatus = []models.ProgressRequestStatus{
			{ProgressIcon: "2", ProgressName: "ยกเลิก", ProgressDatetime: request.CanceledRequestDatetime},
		}
		request.ProgressRequestStatusEmp = models.ProgressRequestStatusEmp{
			ActionRole:   "ผู้ขออนุมัติ",
			EmpID:        request.CreatedRequestEmpID,
			EmpName:      request.CreatedRequestEmpName,
			EmpPosition:  request.CreatedRequestEmpPosition,
			DeptSAP:      request.CreatedRequestDeptSap,
			DeptSAPShort: request.CreatedRequestDeptSapNameShort,
			DeptSAPFull:  request.CreatedRequestDeptSapNameFull,
			PhoneNumber:  request.CreatedRequestPhoneNumber,
			MobileNumber: request.CreatedRequestMobileNumber,
		}
	} else if request.RefRequestAnnualDriverStatusCode == "90" && request.CanceledRequestEmpID == request.ConfirmedRequestEmpID {
		request.ProgressRequestStatus = []models.ProgressRequestStatus{
			{ProgressIcon: "2", ProgressName: "ยกเลิกจากต้นสังกัด", ProgressDatetime: request.CanceledRequestDatetime},
		}
		request.ProgressRequestStatusEmp = models.ProgressRequestStatusEmp{
			ActionRole:   "ผู้อนุมัติต้นสังกัด",
			EmpID:        request.ConfirmedRequestEmpID,
			EmpName:      request.ConfirmedRequestEmpName,
			EmpPosition:  request.ConfirmedRequestEmpPosition,
			DeptSAP:      request.ConfirmedRequestDeptSap,
			DeptSAPShort: request.ConfirmedRequestDeptSapShort,
			DeptSAPFull:  request.ConfirmedRequestDeptSapFull,
			PhoneNumber:  request.ConfirmedRequestPhoneNumber,
			MobileNumber: request.ConfirmedRequestMobileNumber,
		}
	}

	if request.RefRequestAnnualDriverStatusCode == "93" && request.CanceledRequestEmpID == request.ApprovedRequestEmpID {
		request.ProgressRequestStatus = []models.ProgressRequestStatus{
			{ProgressIcon: "3", ProgressName: "อนุมัติจากต้นสังกัด", ProgressDatetime: request.ApprovedRequestDatetime},
			{ProgressIcon: "2", ProgressName: "ยกเลิกจากผู้อนุมัติ", ProgressDatetime: request.CanceledRequestDatetime},
		}
		request.ProgressRequestStatusEmp = models.ProgressRequestStatusEmp{
			ActionRole:   "ผู้อนุมัติให้ทำหน้าที่ขับรถประจำปี",
			EmpID:        request.ApprovedRequestEmpID,
			EmpName:      request.ApprovedRequestEmpName,
			EmpPosition:  request.ApprovedRequestEmpPosition,
			DeptSAP:      request.ApprovedRequestDeptSap,
			DeptSAPShort: request.ApprovedRequestDeptSapShort,
			DeptSAPFull:  request.ApprovedRequestDeptSapFull,
			PhoneNumber:  request.ApprovedRequestPhoneNumber,
			MobileNumber: request.ApprovedRequestMobileNumber,
		}
	}
	request.RefRequestAnnualDriverStatusName = LicenseStatusNameMapConfirmer[request.RefRequestAnnualDriverStatusCode]
	request.CreatedRequestImageUrl = funcs.GetEmpImage(request.CreatedRequestEmpID)
	request.ConfirmedRequestImageUrl = funcs.GetEmpImage(request.ConfirmedRequestEmpID)
	request.ApprovedRequestImageUrl = funcs.GetEmpImage(request.ApprovedRequestEmpID)
	request.ProgressRequestHistory = GetProgressRequestHistory(request)
	c.JSON(http.StatusOK, request)
}

// UpdateDriverLicenseAnnualCanceled godoc
// @Summary Update cancel status for a driver license annual record
// @Description This endpoint allows users to update the cancel status of a driver license annual record.
// @Tags Driver-license-user
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Param data body models.VmsDriverLicenseAnnualCanceled true "VmsDriverLicenseAnnualCanceled data"
// @Router /api/driver-license-user/update-license-annual-canceled [put]
func (h *DriverLicenseUserHandler) UpdateDriverLicenseAnnualCanceled(c *gin.Context) {
	user := funcs.GetAuthenUser(c, h.Role)
	if c.IsAborted() {
		return
	}
	var request, driverLicenseAnnual models.VmsDriverLicenseAnnualCanceled
	var result struct {
		models.VmsDriverLicenseAnnualCanceled
		models.VmsTrnRequestAnnualDriverNo
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "message": messages.ErrInvalidJSONInput.Error()})
		return
	}
	query := h.SetQueryRole(user, config.DB)
	if err := query.First(&driverLicenseAnnual, "trn_request_annual_driver_uid = ? AND is_deleted = ?", request.TrnRequestAnnualDriverUID, "0").Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Driver license annual record not found", "message": messages.ErrNotfound.Error()})
		return
	}
	request.RefRequestAnnualDriverStatusCode = "90"
	request.UpdatedAt = time.Now()
	request.UpdatedBy = user.EmpID

	empUser := funcs.GetUserEmpInfo(user.EmpID)
	request.CanceledRequestEmpID = empUser.EmpID
	request.CanceledRequestEmpName = empUser.FullName
	request.CanceledRequestDeptSAP = empUser.DeptSAP
	request.CanceledRequestDeptSAPShort = empUser.DeptSAPShort
	request.CanceledRequestDeptSAPFull = empUser.DeptSAPFull
	request.CanceledRequestDatetime = models.TimeWithZone{Time: time.Now()}

	if err := config.DB.Save(&request).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to update: %v", err), "message": messages.ErrInternalServer.Error()})
		return
	}

	if err := config.DB.First(&result, "trn_request_annual_driver_uid = ? AND is_deleted = ?", request.TrnRequestAnnualDriverUID, "0").Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Driver license annual record not found", "message": messages.ErrNotfound.Error()})
		return
	}
	funcs.CreateRequestAnnualLicenseNotification(request.TrnRequestAnnualDriverUID)
	c.JSON(http.StatusOK, gin.H{"message": "Updated successfully", "result": result})
}
