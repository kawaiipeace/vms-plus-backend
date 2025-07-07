package funcs

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
	"vms_plus_be/config"
	"vms_plus_be/models"
	"vms_plus_be/userhub"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

// Claims for JWT
type Claims struct {
	EmpID         string   `json:"emp_id"`
	TokenType     string   `json:"token_type"`
	FirstName     string   `json:"first_name"`
	LastName      string   `json:"last_name"`
	FullName      string   `json:"full_name"`
	Position      string   `json:"position"`
	DeptSAP       string   `json:"dept_sap"`
	DeptSAPShort  string   `json:"dept_sap_short"`
	DeptSAPFull   string   `json:"dept_sap_full"`
	BureauDeptSap string   `json:"bureau_dept_sap"`
	MobilePhone   string   `json:"mobile_number"`
	DeskPhone     string   `json:"internal_number"`
	BusinessArea  string   `json:"business_area"`
	ImageUrl      string   `json:"image_url"`
	Roles         []string `json:"roles"`
	IsEmployee    bool     `json:"is_employee"`
	LevelCode     string   `json:"level_code"`
	LoginBy       string   `json:"login_by"`
	jwt.RegisteredClaims
}

var (
	jwtSecret = []byte(config.AppConfig.JWTSecret)
)

func CheckConfirmerRole(user *models.AuthenUserEmp) {
	if Contains(user.Roles, "level1-approval") {
		//remove level1-approval
		user.Roles = RemoveFromSlice(user.Roles, "level1-approval")
	}
	//check if vms_trn_request has confirmed_request_emp_id
	var trnRequestList models.VmsTrnRequestList
	err := config.DB.Where("confirmed_request_emp_id = ?", user.EmpID).First(&trnRequestList).Error
	if err == nil {
		user.Roles = append(user.Roles, "level1-approval")
		return
	}

	//check if vms_trn_request_annual_driver has confirmed_request_emp_id
	var driverLicenseAnnualList models.VmsDriverLicenseAnnualList
	err = config.DB.Where("confirmed_request_emp_id = ?", user.EmpID).First(&driverLicenseAnnualList).Error
	if err == nil {
		user.Roles = append(user.Roles, "level1-approval")
		return
	}

}
func CheckApproverRole(user *models.AuthenUserEmp) {
	if Contains(user.Roles, "license-approval") {
		//remove license-approval
		user.Roles = RemoveFromSlice(user.Roles, "license-approval")
	}

	//check if exist vms_trn_request_annual_driver has approved_request_emp_id=user.EmpID
	var driverLicenseAnnualList models.VmsDriverLicenseAnnualList
	err := config.DB.Where("approved_request_emp_id = ?", user.EmpID).First(&driverLicenseAnnualList).Error
	if err == nil {
		user.Roles = append(user.Roles, "license-approval")
		return
	}
}

func CheckAdminApprovalRole(user *models.AuthenUserEmp) {
	if Contains(user.Roles, "admin-approval") {
		//remove admin-approval
		user.Roles = RemoveFromSlice(user.Roles, "admin-approval")
	}
	//check if vms_trn_request has confirmed_request_emp_id
	var adminApproval models.VmsMasCarpoolAdmin
	err := config.DB.Where("admin_emp_no = ? AND is_deleted = '0' AND is_active = '1'", user.EmpID).
		Where("mas_carpool_uid IN (SELECT mas_carpool_uid FROM vms_mas_carpool WHERE is_deleted = '0' AND is_active = '1')").
		First(&adminApproval).Error
	if err == nil {
		user.Roles = append(user.Roles, "admin-approval")
		return
	}
}
func CheckFinalApprovalRole(user *models.AuthenUserEmp) {
	if Contains(user.Roles, "final-approval") {
		//remove final-approval
		user.Roles = RemoveFromSlice(user.Roles, "final-approval")
	}
	//check if vms_trn_request has confirmed_request_emp_id
	var finalApproval models.VmsMasCarpoolApprover
	err := config.DB.Where("approver_emp_no = ? AND is_deleted = '0' AND is_active = '1'", user.EmpID).
		Where("mas_carpool_uid IN (SELECT mas_carpool_uid FROM vms_mas_carpool WHERE is_deleted = '0' AND is_active = '1')").
		First(&finalApproval).Error
	if err == nil {
		user.Roles = append(user.Roles, "final-approval")
		return
	}
}

func RemoveFromSlice(slice []string, value string) []string {
	for i, v := range slice {
		if v == value {
			return append(slice[:i], slice[i+1:]...)
		}
	}
	return slice
}

func GenerateJWT(user models.AuthenUserEmp, tokenType string, expiration time.Duration) (string, error) {
	CheckConfirmerRole(&user)
	CheckApproverRole(&user)
	CheckAdminApprovalRole(&user)
	CheckFinalApprovalRole(&user)
	jwtSecret = []byte(config.AppConfig.JWTSecret)
	claims := Claims{
		EmpID:         user.EmpID,
		FullName:      user.FullName,
		TokenType:     tokenType,
		LoginBy:       user.LoginBy,
		FirstName:     user.FirstName,
		LastName:      user.LastName,
		Position:      user.Position,
		DeptSAP:       user.DeptSAP,
		DeptSAPShort:  user.DeptSAPShort,
		DeptSAPFull:   user.DeptSAPFull,
		BureauDeptSap: user.BureauDeptSap,
		MobilePhone:   user.MobilePhone,
		DeskPhone:     user.DeskPhone,
		BusinessArea:  user.BusinessArea,
		ImageUrl:      user.ImageUrl,
		Roles:         user.Roles,
		IsEmployee:    user.IsEmployee,
		LevelCode:     user.LevelCode,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiration)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}
func GenerateRefreshJWT(user models.AuthenUserEmp, tokenType string, expiration time.Duration) (string, error) {
	jwtSecret = []byte(config.AppConfig.JWTSecret)
	claims := Claims{
		EmpID:     user.EmpID,
		FullName:  user.FullName,
		TokenType: tokenType,
		LoginBy:   user.LoginBy,

		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiration)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

func ExtractUserFromJWT(c *gin.Context) (*models.AuthenUserEmp, error) {
	// Extract JWT token from Authorization header
	authHeader := c.GetHeader("Authorization")
	secretKey := config.AppConfig.JWTSecret
	if authHeader == "" {
		return nil, errors.New("missing Authorization header")
	}

	// Remove "Bearer " prefix
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	if tokenString == authHeader { // If "Bearer " prefix is missing
		return nil, errors.New("invalid token format")
	}

	// Parse and validate the JWT token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil // Use your secret key
	})

	if err != nil || !token.Valid {
		return nil, errors.New("invalid or expired token")
	}

	// Extract claims from token
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}

	rolesInterface := claims["roles"].([]interface{})
	roles := make([]string, len(rolesInterface))

	for i, v := range rolesInterface {
		roles[i] = v.(string) // Perform type assertion
	}

	// Map claims to UserEmp struct
	user := &models.AuthenUserEmp{
		EmpID:         claims["emp_id"].(string),
		FullName:      claims["full_name"].(string),
		LoginBy:       claims["login_by"].(string),
		FirstName:     claims["first_name"].(string),
		LastName:      claims["last_name"].(string),
		Position:      claims["position"].(string),
		DeptSAP:       claims["dept_sap"].(string),
		DeptSAPShort:  claims["dept_sap_short"].(string),
		DeptSAPFull:   claims["dept_sap_full"].(string),
		BureauDeptSap: claims["bureau_dept_sap"].(string),
		MobilePhone:   claims["mobile_number"].(string),
		DeskPhone:     claims["internal_number"].(string),
		BusinessArea:  claims["business_area"].(string),
		ImageUrl:      claims["image_url"].(string),
		IsEmployee:    claims["is_employee"].(bool),
		LevelCode:     claims["level_code"].(string),
		Roles:         roles,
	}

	return user, nil
}

func GetAuthenUser(c *gin.Context, roles string) *models.AuthenUserEmp {
	// Extract user from JWT
	var empUser models.AuthenUserEmp
	//501621 //510683
	if config.AppConfig.IsDev && c.Request.Header.Get("Authorization") == "" {
		user, err := userhub.GetUserInfo("700001")
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			c.Abort()
			return &empUser
		}
		empUser = user
		empUser.LoginBy = "keycloak"
		empUser.IsEmployee = true
		CheckConfirmerRole(&empUser)
		CheckApproverRole(&empUser)
		CheckAdminApprovalRole(&empUser)
		CheckFinalApprovalRole(&empUser)
		//empUser.Roles = append(empUser.Roles, "license-approval")
		//empUser.Roles = []string{"admin-region"}
		if empUser.LevelCode == "M5" {
			empUser.IsLevelM5 = "1"
		} else {
			empUser.IsLevelM5 = "0"
		}

		if roles == "*" {
			return &empUser
		}
		for _, role := range strings.Split(roles, ",") {
			if Contains(empUser.Roles, role) {
				return &empUser
			}
		}

		c.JSON(http.StatusForbidden, gin.H{"error": "insufficient permissions"})
		c.Abort()
		return &models.AuthenUserEmp{}
	}
	jwt, err := ExtractUserFromJWT(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error(), "message": "Please login again"})
		c.Abort()
	}
	empUser = models.AuthenUserEmp{
		EmpID:         jwt.EmpID,
		FirstName:     jwt.FirstName,
		LastName:      jwt.LastName,
		FullName:      jwt.FullName,
		Position:      jwt.Position,
		DeptSAP:       jwt.DeptSAP,
		DeptSAPShort:  jwt.DeptSAPShort,
		DeptSAPFull:   jwt.DeptSAPFull,
		BureauDeptSap: jwt.BureauDeptSap,
		MobilePhone:   jwt.MobilePhone,
		DeskPhone:     jwt.DeskPhone,
		BusinessArea:  jwt.BusinessArea,
		ImageUrl:      jwt.ImageUrl,
		Roles:         jwt.Roles,
		LoginBy:       jwt.LoginBy,
		IsEmployee:    jwt.IsEmployee,
		LevelCode:     jwt.LevelCode,
	}
	if roles == "level1-approval" {
		CheckConfirmerRole(&empUser)
	}
	if roles == "license-approval" {
		CheckApproverRole(&empUser)
	}
	if roles == "admin-approval" {
		CheckAdminApprovalRole(&empUser)
	}
	if roles == "final-approval" {
		CheckFinalApprovalRole(&empUser)
	}

	if empUser.LevelCode == "M5" {
		empUser.IsLevelM5 = "1"
	} else {
		empUser.IsLevelM5 = "0"
	}
	if roles == "*" {
		return &empUser
	}
	for _, role := range strings.Split(roles, ",") {
		if Contains(empUser.Roles, role) {
			return &empUser
		}
	}

	c.JSON(http.StatusForbidden, gin.H{"error": "insufficient permissions"})
	c.Abort()

	return &empUser
}

func GetUserEmpInfo(empID string) models.MasUserEmp {
	user, err := userhub.GetUserInfo(empID)
	if err != nil {
		return models.MasUserEmp{}
	}
	empUser := models.MasUserEmp{
		EmpID:         user.EmpID,
		FullName:      user.FullName,
		Position:      user.Position,
		DeptSAP:       user.DeptSAP,
		DeptSAPShort:  user.DeptSAPShort,
		DeptSAPFull:   user.DeptSAPFull,
		TelMobile:     user.MobilePhone,
		TelInternal:   user.DeskPhone,
		BusinessArea:  user.BusinessArea,
		BureauDeptSap: user.BureauDeptSap,
		IsEmployee:    user.IsEmployee,
	}
	return empUser
}

func SetQueryAdminRole(user *models.AuthenUserEmp, query *gorm.DB) *gorm.DB {
	query = query.Where(
		`exists (
			select 1 from vms_mas_carpool_admin ca 
			where ca.mas_carpool_uid = vms_trn_request.mas_carpool_uid 
			and ca.admin_emp_no = ? and ca.is_deleted = '0' and ca.is_active = '1' 
		) or exists (
			select 1 from vms_mas_vehicle_department vd 
			where vd.mas_vehicle_uid = vms_trn_request.mas_vehicle_uid 
			and vd.bureau_dept_sap in (?)	
		)`,
		user.EmpID,
		user.BureauDeptSap,
	)
	return query
}

func SetQueryApproverRole(user *models.AuthenUserEmp, query *gorm.DB) *gorm.DB {
	query = query.Where(
		`exists (
			select 1 from vms_mas_carpool_approver ca 
			where ca.mas_carpool_uid = vms_trn_request.mas_carpool_uid 
			and ca.approver_emp_no = ? and ca.is_deleted = '0' and ca.is_active = '1' 
		) or exists (
			select 1 from vms_mas_vehicle_department vd 
			where vd.mas_vehicle_uid = vms_trn_request.mas_vehicle_uid 
			and vd.bureau_dept_sap in (?)	
		)`,
		user.EmpID,
		user.BureauDeptSap,
	)
	return query
}

func GetDeptSAPShort(deptSAP string) string {
	//call db public.fn_get_long_short_dept_name_by_dept_sap
	var deptSAPShort string
	err := config.DB.Raw("SELECT public.fn_get_long_short_dept_name_by_dept_sap(?)", deptSAP).Scan(&deptSAPShort).Error
	if err != nil {
		return ""
	}
	return deptSAPShort
}
func GetDeptSAPFull(deptSAP string) string {
	//call db public.fn_get_long_full_dept_name_by_dept_sap
	var deptSAPFull string
	err := config.DB.Raw("SELECT public.fn_get_long_full_dept_name_by_dept_sap(?)", deptSAP).Scan(&deptSAPFull).Error
	if err != nil {
		return ""
	}
	return deptSAPFull
}

func GetUserManager(DeptSAP string) []models.VmsMasManager {
	url := config.AppConfig.HrPlatformEndPoint + "/get-manager?dept_sap=" + DeptSAP
	resp, err := http.Get(url)
	if err != nil {
		return []models.VmsMasManager{}
	}
	fmt.Println(url)
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return []models.VmsMasManager{}
	}
	type ResponseData struct {
		Data struct {
			Data struct {
				DataDetail []models.VmsMasManager `json:"dataDetail"`
			} `json:"data"`
		} `json:"data"`
	}
	var response ResponseData
	if err := json.Unmarshal([]byte(body), &response); err != nil {
		fmt.Println(err)
	}
	return response.Data.Data.DataDetail
}

func GetBusinessAreaCodeFromDeptSap(deptSAP string) string {
	//call db public.fn_get_business_area_code_from_dept_sap
	var businessArea string
	err := config.DB.Table("vms_mas_department").Where("dept_sap = ?", deptSAP).
		Select("business_area").First(&businessArea).Error
	if err != nil {
		return "Z000"
	}
	return businessArea
}
