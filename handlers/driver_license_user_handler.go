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
	if user.EmpID == "" {
		return query
	}
	return query
}

func (h *DriverLicenseUserHandler) SetQueryRoleDept(user *models.AuthenUserEmp, query *gorm.DB) *gorm.DB {
	if user.EmpID == "" {
		return query
	}
	return query
}
func (h *DriverLicenseUserHandler) SetQueryStatusCanUpdate(query *gorm.DB) *gorm.DB {
	return query.Where("ref_request_status_code in ('11') and is_deleted = '0'")
}

var StatusDriverAnnualLicense = map[string]string{
	"10": "ไม่มีใบอนุญาต",
	"20": "กำลังดำเนินการ",
	"30": "อนุมัติแล้ว",
	"31": "มีผลปีถัดไป",
	"80": "หมดอายุ",
	"90": "ยกเลิก",
}

func GetProgressRequestHistory(request models.VmsDriverLicenseAnnualResponse) []models.ProgressRequestHistory {
	var progressRequestHistory []models.ProgressRequestHistory
	if request.RefRequestAnnualDriverStatusCode == "10" {
		progressRequestHistory = append(progressRequestHistory, models.ProgressRequestHistory{
			ProgressIcon:     "3",
			ProgressName:     "ขออนุมัติ",
			ProgressDateTime: request.CreatedRequestDatetime,
		})
	}
	if request.RefRequestAnnualDriverStatusCode == "20" {
		progressRequestHistory = append(progressRequestHistory, models.ProgressRequestHistory{
			ProgressIcon:     "3",
			ProgressName:     "อนุมัติจากต้นสังกัด",
			ProgressDateTime: request.ConfirmedRequestDatetime,
		})
	}
	if request.RefRequestAnnualDriverStatusCode == "30" {
		progressRequestHistory = append(progressRequestHistory, models.ProgressRequestHistory{
			ProgressIcon:     "3",
			ProgressName:     "อนุมัติให้ทำหน้าที่ขับรถยนต์",
			ProgressDateTime: request.ApprovedRequestDatetime,
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

	//Check VmsDriverLicenseAnnualList
	var license models.VmsDriverLicenseAnnualResponse
	err := config.DB.Where("created_request_emp_id = ? and is_deleted = ?", user.EmpID, "0").
		Preload("DriverLicenseType").
		Preload("DriverCertificateType").
		Order("created_request_datetime DESC").
		Find(&license).Error

	if err == nil {
		if license.RefRequestAnnualDriverStatusCode == "30" {
			driver.LicenseStatusCode = "30"
		} else if license.RefRequestAnnualDriverStatusCode == "90" {
			driver.LicenseStatusCode = "90"
		} else {
			driver.LicenseStatusCode = "20"
		}
		driver.LicenseStatus = StatusDriverAnnualLicense[driver.LicenseStatusCode]
	} else {
		driver.LicenseStatusCode = "10"
		driver.LicenseStatus = StatusDriverAnnualLicense[driver.LicenseStatusCode]
	}

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
	annualYearEnd := time.Date(driver.AnnualYYYY-543, 12, 31, 23, 59, 59, 0, time.UTC)

	// Set DriverLicenseEndDate to the minimum of license.DriverLicenseExpireDate and annualYearEnd
	if license.DriverLicenseExpireDate.Before(annualYearEnd) {
		driver.DriverLicense.DriverLicenseEndDate = license.DriverLicenseExpireDate
	} else {
		driver.DriverLicense.DriverLicenseEndDate = annualYearEnd
	}
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
	request.TrnRequestAnnualDriverUID = uuid.New().String()
	request.CreatedRequestEmpID = user.EmpID
	empUser := funcs.GetUserEmpInfo(request.CreatedRequestEmpID)
	request.CreatedRequestEmpName = empUser.FullName
	request.CreatedRequestEmpPosition = empUser.Position
	request.CreatedRequestDeptSap = empUser.DeptSAP
	request.CreatedRequestDeptSapNameShort = empUser.DeptSAPShort
	request.CreatedRequestDeptSapNameFull = empUser.DeptSAPFull
	request.CreatedRequestMobileNumber = empUser.MobilePhone
	request.CreatedRequestPhoneNumber = empUser.DeskPhone
	request.CreatedRequestDatetime = time.Now()

	confirmUser := funcs.GetUserEmpInfo(request.ConfirmedRequestEmpID)
	request.ConfirmedRequestEmpName = confirmUser.FullName
	request.ConfirmedRequestEmpPosition = confirmUser.Position
	request.ConfirmedRequestDeptSap = confirmUser.DeptSAP
	request.ConfirmedRequestDeptSapShort = confirmUser.DeptSAPShort
	request.ConfirmedRequestDeptSapFull = confirmUser.DeptSAPFull
	request.ConfirmedRequestMobileNumber = empUser.MobilePhone
	request.ConfirmedRequestPhoneNumber = empUser.DeskPhone

	approveUser := funcs.GetUserEmpInfo(request.ApprovedRequestEmpID)
	request.ApprovedRequestEmpName = approveUser.FullName
	request.ApprovedRequestEmpPosition = approveUser.Position
	request.ApprovedRequestDeptSap = approveUser.DeptSAP
	request.ApprovedRequestDeptSapShort = approveUser.DeptSAPShort
	request.ApprovedRequestDeptSapFull = approveUser.DeptSAPFull
	request.ApprovedRequestMobileNumber = empUser.MobilePhone
	request.ApprovedRequestPhoneNumber = empUser.DeskPhone

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

	if err := config.DB.Create(&request).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create driver license annual record", "message": messages.ErrInternalServer.Error()})
		return
	}
	if err := config.DB.First(&result, "trn_request_annual_driver_uid = ?", request.TrnRequestAnnualDriverUID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "annual not found", "message": messages.ErrNotfound.Error()})
		return
	}

	// Return success response
	c.JSON(http.StatusCreated, gin.H{"message": "Driver license annual record created successfully", "result": result})
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
			{ProgressIcon: "3", ProgressName: "ขออนุมัติ"},
			{ProgressIcon: "1", ProgressName: "รอต้นสังกัดตรวจสอบ"},
			{ProgressIcon: "0", ProgressName: "รออนุมัติให้ทำหน้าที่ขับรถประจำปี"},
		}
	}
	if request.RefRequestAnnualDriverStatusCode == "11" {
		request.ProgressRequestStatus = []models.ProgressRequestStatus{
			{ProgressIcon: "3", ProgressName: "ขออนุมัติ"},
			{ProgressIcon: "2", ProgressName: "ตีกลับจากต้นสังกัด"},
			{ProgressIcon: "0", ProgressName: "รออนุมัติให้ทำหน้าที่ขับรถประจำปี"},
		}
	}
	if request.RefRequestAnnualDriverStatusCode == "20" {
		request.ProgressRequestStatus = []models.ProgressRequestStatus{
			{ProgressIcon: "3", ProgressName: "ขออนุมัติ"},
			{ProgressIcon: "3", ProgressName: "ต้นสังกัดตรวจสอบ"},
			{ProgressIcon: "1", ProgressName: "รออนุมัติให้ทำหน้าที่ขับรถประจำปี"},
		}
	}
	if request.RefRequestAnnualDriverStatusCode == "21" {
		request.ProgressRequestStatus = []models.ProgressRequestStatus{
			{ProgressIcon: "3", ProgressName: "ขออนุมัติ"},
			{ProgressIcon: "3", ProgressName: "ต้นสังกัดตรวจสอบ"},
			{ProgressIcon: "2", ProgressName: "ตีกลับจากผู้อนุมัติ"},
		}
	}
	if request.RefRequestAnnualDriverStatusCode == "30" {
		request.ProgressRequestStatus = []models.ProgressRequestStatus{
			{ProgressIcon: "3", ProgressName: "ขออนุมัติ"},
			{ProgressIcon: "3", ProgressName: "ต้นสังกัดตรวจสอบ"},
			{ProgressIcon: "3", ProgressName: "อนุมัติให้ทำหน้าที่ขับรถประจำปี"},
		}

	}
	if request.RefRequestAnnualDriverStatusCode == "90" {
		request.ProgressRequestStatus = []models.ProgressRequestStatus{
			{ProgressIcon: "2", ProgressName: "ยกเลิก"},
		}
	}
	if request.RefRequestAnnualDriverStatusCode == "91" {
		request.ProgressRequestStatus = []models.ProgressRequestStatus{
			{ProgressIcon: "2", ProgressName: "ยกเลิกจากผู้ขอ"},
		}
	}
	if request.RefRequestAnnualDriverStatusCode == "92" {
		request.ProgressRequestStatus = []models.ProgressRequestStatus{
			{ProgressIcon: "2", ProgressName: "ยกเลิกจากต้นสังกัด"},
		}
	}
	if request.RefRequestAnnualDriverStatusCode == "93" {
		request.ProgressRequestStatus = []models.ProgressRequestStatus{
			{ProgressIcon: "3", ProgressName: "อนุมัติจากต้นสังกัด"},
			{ProgressIcon: "2", ProgressName: "ยกเลิกจากผู้อนุมัติ"},
		}
	}
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
	request.CanceledRequestDatetime = time.Now()

	if err := config.DB.Save(&request).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to update: %v", err), "message": messages.ErrInternalServer.Error()})
		return
	}

	if err := config.DB.First(&result, "trn_request_annual_driver_uid = ? AND is_deleted = ?", request.TrnRequestAnnualDriverUID, "0").Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Driver license annual record not found", "message": messages.ErrNotfound.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Updated successfully", "result": result})
}
