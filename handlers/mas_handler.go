package handlers

import (
	"net/http"
	"time"
	"vms_plus_be/config"
	"vms_plus_be/models"

	"github.com/gin-gonic/gin"
)

type MasHandler struct {
}

// ListVehicleUser godoc
// @Summary Retrieve the Vehicle Users
// @Description This endpoint allows a user to retrieve Vehicle Users.
// @Tags MAS
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Param search query string false "Search by Employee ID or Full Name"
// @Router /api/mas/user-vehicle-users [get]
func (h *MasHandler) ListVehicleUser(c *gin.Context) {
	var lists []models.MasUserDriver
	search := c.Query("search")

	query := config.DB

	// Apply search filter if provided
	if search != "" {
		query = query.Where("emp_id = ? OR full_name ILIKE ?", search, "%"+search+"%")
	}

	// Execute query
	if err := query.
		Find(&lists).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Not found"})
		return
	}

	defaultURL := "http://pntdev.ddns.net:28089/VMS_PLUS/PIX/user-avatar.jpg"

	// Loop to modify or set AnnualDriver
	for i := range lists {
		lists[i].Image_url = defaultURL
	}

	c.JSON(http.StatusOK, lists)
}

// ListDriverUser godoc
// @Summary Retrieve the Driver Users
// @Description This endpoint allows a user to retrieve Driver Users.
// @Tags MAS
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Param search query string false "Search by Employee ID or Full Name"
// @Router /api/mas/user-driver-users [get]
func (h *MasHandler) ListDriverUser(c *gin.Context) {
	var lists []models.MasUserDriver
	search := c.Query("search")

	query := config.DB

	// Apply search filter if provided
	if search != "" {
		query = query.Where("emp_id = ? OR full_name ILIKE ?", search, "%"+search+"%")
	}

	// Execute query
	if err := query.Preload("AnnualDriver").
		Find(&lists).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Not found"})
		return
	}

	defaultURL := "http://pntdev.ddns.net:28089/VMS_PLUS/PIX/user-avatar.jpg"

	// Loop to modify or set AnnualDriver
	for i := range lists {
		lists[i].Image_url = defaultURL
		lists[i].AnnualDriver.AnnualYYYY = 2568
		lists[i].AnnualDriver.DriverLicenseNo = "A123456"
		lists[i].AnnualDriver.RequestAnnualDriverNo = "B00001"
		lists[i].AnnualDriver.DriverLicenseExpireDate = time.Now().AddDate(1, 0, 0)
		lists[i].AnnualDriver.RequestIssueDate = time.Now().AddDate(0, -6, 0)
		lists[i].AnnualDriver.RequestExpireDate = time.Now().AddDate(0, 6, 0)
	}

	c.JSON(http.StatusOK, lists)
}

// ListApprovalUser godoc
// @Summary Retrieve the Approval Users
// @Description This endpoint allows a user to retrieve Approval Users.
// @Tags MAS
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Param search query string false "Search by Employee ID or Full Name"
// @Router /api/mas/user-approval-users [get]
func (h *MasHandler) ListApprovalUser(c *gin.Context) {
	var lists []models.MasUserEmp
	search := c.Query("search")

	query := config.DB

	// Apply search filter if provided
	if search != "" {
		query = query.Where("emp_id = ? OR full_name ILIKE ?", search, "%"+search+"%")
	}

	// Execute query
	if err := query.
		Find(&lists).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Not found"})
		return
	}

	defaultURL := "http://pntdev.ddns.net:28089/VMS_PLUS/PIX/user-avatar.jpg"
	// For loop to set Image_url for each element in the slice
	for i := range lists {
		lists[i].Image_url = defaultURL
	}

	c.JSON(http.StatusOK, lists)
}

// ListAdminApprovalUser godoc
// @Summary Retrieve the Admin Approval Users
// @Description This endpoint allows a user to retrieve Admin Approval Users.
// @Tags MAS
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Param search query string false "Search by Employee ID or Full Name"
// @Router /api/mas/user-admin-approval-users [get]
func (h *MasHandler) ListAdminApprovalUser(c *gin.Context) {
	var lists []models.MasUserEmp
	search := c.Query("search")

	query := config.DB

	// Apply search filter if provided
	if search != "" {
		query = query.Where("emp_id = ? OR full_name ILIKE ?", search, "%"+search+"%")
	}

	// Execute query
	if err := query.
		Find(&lists).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Not found"})
		return
	}

	defaultURL := "http://pntdev.ddns.net:28089/VMS_PLUS/PIX/user-avatar.jpg"
	// For loop to set Image_url for each element in the slice
	for i := range lists {
		lists[i].Image_url = defaultURL
	}

	c.JSON(http.StatusOK, lists)
}

// ListFinalApprovalUser godoc
// @Summary Retrieve the Final Approval User
// @Description This endpoint allows a user to retrieve Final Approval User.
// @Tags MAS
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Param search query string false "Search by Employee ID or Full Name"
// @Router /api/mas/user-final-approval-users [get]
func (h *MasHandler) ListFinalApprovalUser(c *gin.Context) {
	var lists []models.MasUserEmp
	search := c.Query("search")

	query := config.DB

	// Apply search filter if provided
	if search != "" {
		query = query.Where("emp_id = ? OR full_name ILIKE ?", search, "%"+search+"%")
	}

	// Execute query
	if err := query.
		Find(&lists).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Not found"})
		return
	}

	defaultURL := "http://pntdev.ddns.net:28089/VMS_PLUS/PIX/user-avatar.jpg"
	// For loop to set Image_url for each element in the slice
	for i := range lists {
		lists[i].Image_url = defaultURL
	}

	c.JSON(http.StatusOK, lists)
}

// GetUserEmp godoc
// @Summary Retrieve a specific user
// @Description This endpoint fetches details of a user.
// @Tags MAS
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security AuthorizationAuth
// @Param id path string true "EmpID (EmpID)"
// @Router /api/mas/user/{id} [get]
func (h *MasHandler) GetUserEmp(c *gin.Context) {
	//funcs.GetAuthenUser(c, h.Role)
	EmpID := c.Param("id")

	var userEmp models.MasUserEmp
	if err := config.DB.
		First(&userEmp, "emp_id = ?", EmpID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	c.JSON(http.StatusOK, userEmp)
}
