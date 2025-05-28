package funcs

import (
	"errors"
	"net/http"
	"strings"
	"time"
	"vms_plus_be/config"
	"vms_plus_be/models"
	"vms_plus_be/userhub"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
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

	LoginBy string `json:"login_by"`
	jwt.RegisteredClaims
}

var (
	jwtSecret = []byte(config.AppConfig.JWTSecret)
)

func GenerateJWT(user models.AuthenUserEmp, tokenType string, expiration time.Duration) (string, error) {
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
		empUser.Roles = []string{"vehicle-user", "level1-approval", "admin-approval", "admin-dept-approval", "final-approval", "license-confirmer",
			"license-approval", "driver", "admin-super"}
		empUser.LoginBy = "keycloak"
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
	}
	empUser.Roles = []string{"vehicle-user", "level1-approval", "admin-approval", "admin-dept-approval", "final-approval", "license-confirmer",
		"license-approval", "driver", "admin-super"}
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
		EmpID:        user.EmpID,
		FullName:     user.FullName,
		Position:     user.Position,
		DeptSAP:      user.DeptSAP,
		DeptSAPShort: user.DeptSAPShort,
		DeptSAPFull:  user.DeptSAPFull,
		TelMobile:    user.MobilePhone,
		TelInternal:  user.DeskPhone,
		BusinessArea: user.BusinessArea,
	}
	return empUser
}
