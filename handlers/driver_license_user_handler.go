package handlers

import (
	"fmt"
	"net/http"
	"time"
	"vms_plus_be/config"
	"vms_plus_be/funcs"
	"vms_plus_be/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type DriverLicenseUserHandler struct {
	Role string
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
	//user := funcs.GetAuthenUser(c, h.Role)
	masDriverUID := "ed9ccc24-2dd6-4294-8136-a78e1bdc6362"
	var driver models.VmsDriverLicenseCard

	if err := config.DB.Where("mas_driver_uid = ? AND is_deleted = ?", masDriverUID, "0").
		Preload("DriverLicense", func(db *gorm.DB) *gorm.DB {
			return db.Order("driver_license_end_date DESC").Limit(1)
		}).
		Preload("DriverLicense.DriverLicenseType").
		First(&driver).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Driver not found"})
		return
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
	var request models.VmsDriverLicenseAnnualRequest
	var result struct {
		models.VmsDriverLicenseAnnualRequest
		RequestAnnualDriverNo     string `gorm:"column:request_annual_driver_no" json:"request_annual_driver_no"`
		TrnRequestAnnualDriverUID string `gorm:"column:trn_request_annual_driver_uid" json:"trn_request_annual_driver_uid"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON input"})
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
	request.CreatedRequestMobileNumber = empUser.MobileNumber
	request.CreatedRequestPhoneNumber = empUser.InternalNumber
	request.CreatedRequestDatetime = time.Now()

	confirmUser := funcs.GetUserEmpInfo(request.ConfirmedRequestEmpID)
	request.ConfirmedRequestEmpName = confirmUser.FullName
	request.ConfirmedRequestEmpPosition = confirmUser.Position
	request.ConfirmedRequestDeptSap = confirmUser.DeptSAP
	request.ConfirmedRequestDeptSapShort = confirmUser.DeptSAPShort
	request.ConfirmedRequestDeptSapFull = confirmUser.DeptSAPFull
	request.ConfirmedRequestMobileNumber = empUser.MobileNumber
	request.ConfirmedRequestPhoneNumber = empUser.InternalNumber

	approveUser := funcs.GetUserEmpInfo(request.ApprovedRequestEmpID)
	request.ApprovedRequestEmpName = approveUser.FullName
	request.ApprovedRequestEmpPosition = approveUser.Position
	request.ApprovedRequestDeptSap = approveUser.DeptSAP
	request.ApprovedRequestDeptSapShort = approveUser.DeptSAPShort
	request.ApprovedRequestDeptSapFull = approveUser.DeptSAPFull
	request.ApprovedRequestMobileNumber = empUser.MobileNumber
	request.ApprovedRequestPhoneNumber = empUser.InternalNumber

	request.RefRequestAnnualDriverStatusCode = "10"
	request.RejectedRequestEmpPosition = ""
	request.CanceledRequestEmpPosition = ""

	request.UpdatedAt = time.Now()
	request.UpdatedBy = user.EmpID

	var maxRequestNo string
	config.DB.Table("vms_trn_request_annual_driver").
		Select("MAX(request_annual_driver_no)").
		Where("request_annual_driver_no LIKE ?", "RAD%").
		Scan(&maxRequestNo)

	var nextNumber int
	if maxRequestNo != "" {
		fmt.Sscanf(maxRequestNo, "RAD%d", &nextNumber)
	}
	nextNumber++

	request.RequestAnnualDriverNo = fmt.Sprintf("RAD%09d", nextNumber)

	if err := config.DB.Create(&request).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create driver license annual record"})
		return
	}
	if err := config.DB.First(&result, "trn_request_annual_driver_uid = ?", request.TrnRequestAnnualDriverUID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "annual not found"})
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
	//user := funcs.GetAuthenUser(c, h.Role)
	trnRequestAnnualDriverUID := c.Param("trn_request_annual_driver_uid")
	var request models.VmsDriverLicenseAnnualResponse

	if err := config.DB.
		Preload("DriverLicenseType").
		First(&request, "trn_request_annual_driver_uid = ? and is_deleted = ?", trnRequestAnnualDriverUID, "0").Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "annual not found"})
		return
	}
	if request.RefRequestAnnualDriverStatusCode == "10" {
		request.ProgressRequestStatus = []models.ProgressRequestStatus{
			{ProgressIcon: "3", ProgressName: "ขออนุมัติ"},
			{ProgressIcon: "1", ProgressName: "รออนุมัติจากต้นสังกัด"},
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
			{ProgressIcon: "3", ProgressName: "อนุมัติจากต้นสังกัด"},
			{ProgressIcon: "1", ProgressName: "รออนุมัติให้ทำหน้าที่ขับรถประจำปี"},
		}
	}
	if request.RefRequestAnnualDriverStatusCode == "21" {
		request.ProgressRequestStatus = []models.ProgressRequestStatus{
			{ProgressIcon: "3", ProgressName: "ขออนุมัติ"},
			{ProgressIcon: "3", ProgressName: "อนุมัติจากต้นสังกัด"},
			{ProgressIcon: "2", ProgressName: "ผู้อนุมัติตีกลับ"},
		}
	}
	if request.RefRequestAnnualDriverStatusCode == "30" {
		request.ProgressRequestStatus = []models.ProgressRequestStatus{
			{ProgressIcon: "3", ProgressName: "ขออนุมัติ"},
			{ProgressIcon: "3", ProgressName: "อนุมัติจากต้นสังกัด"},
			{ProgressIcon: "3", ProgressName: "อนุมัติให้ทำหน้าที่ขับรถประจำปี"},
		}
	}
	if request.RefRequestAnnualDriverStatusCode == "90" {
		request.ProgressRequestStatus = []models.ProgressRequestStatus{
			{ProgressIcon: "3", ProgressName: "ขออนุมัติ"},
			{ProgressIcon: "3", ProgressName: "ยกเลิกอนุมัติจากต้นสังกัด"},
			{ProgressIcon: "0", ProgressName: "รออนุมัติให้ทำหน้าที่ขับรถประจำปี"},
		}
	}
	// Return success response
	c.JSON(http.StatusCreated, gin.H{"message": "Driver license annual record created successfully", "result": request})
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
	var request, driverLicenseAnnual models.VmsDriverLicenseAnnualCanceled
	var result struct {
		models.VmsDriverLicenseAnnualCanceled
		models.VmsTrnRequestAnnualDriverNo
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := config.DB.First(&driverLicenseAnnual, "trn_request_annual_driver_uid = ? AND is_deleted = ?", request.TrnRequestAnnualDriverUID, "0").Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Driver license annual record not found"})
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to update: %v", err)})
		return
	}

	if err := config.DB.First(&result, "trn_request_annual_driver_uid = ? AND is_deleted = ?", request.TrnRequestAnnualDriverUID, "0").Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Driver license annual record not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Updated successfully", "result": result})
}
