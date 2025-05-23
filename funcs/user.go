package funcs

import (
	"errors"
	"net/http"
	"strings"
	"time"
	"vms_plus_be/config"
	"vms_plus_be/models"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// Claims for JWT
type Claims struct {
	EmpID     string   `json:"emp_id"`
	FullName  string   `json:"full_name"`
	TokenType string   `json:"token_type"`
	Roles     []string `json:"roles"`
	LoginBy   string   `json:"login_by"`
	jwt.RegisteredClaims
}

var (
	jwtSecret = []byte(config.AppConfig.JWTSecret)
)

func GenerateJWT(user models.AuthenUserEmp, tokenType string, expiration time.Duration, accessToken string, refreshToken string) (string, error) {
	jwtSecret = []byte(config.AppConfig.JWTSecret)
	claims := Claims{
		EmpID:     user.EmpID,
		FullName:  user.FirstName + " " + user.LastName,
		TokenType: tokenType,
		LoginBy:   user.LoginBy,
		Roles:     []string{"vehicle-user", "level1-approval", "admin-approval", "admin-dept-approval", "final-approval", "driver", "admin-super"},
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiration)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

func ExtractUserFromJWT(c *gin.Context) (*models.AuthenJwtUsr, error) {
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
	user := &models.AuthenJwtUsr{
		EmpID:    claims["emp_id"].(string),
		FullName: claims["full_name"].(string),
		Roles:    roles,
	}

	return user, nil
}

func GetAuthenUser(c *gin.Context, roles string) *models.AuthenUserEmp {
	// Extract user from JWT
	var empUser models.AuthenUserEmp
	//501621
	if config.AppConfig.IsDev {
		if err := config.DBu.First(&empUser, "emp_id = ?", "700001").Error; err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			c.Abort()
			return &empUser
		}
		empUser.BusinessArea = "Z000"
		empUser.Roles = []string{"vehicle-user", "level1-approval", "admin-approval", "admin-dept-approval", "final-approval", "driver", "admin-super"}
		empUser.LoginBy = "keycloak"
		empUser.ImageUrl = config.DefaultAvatarURL
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

	if err := config.DBu.First(&empUser, "emp_id = ?", jwt.EmpID).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		c.Abort()
		return &empUser
	}
	empUser.Roles = jwt.Roles
	empUser.LoginBy = jwt.LoginBy
	empUser.ImageUrl = GetEmpImage(empUser.EmpID)
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

func GetUserEmpInfo(empID string) models.MasUser {
	var empUser models.MasUser
	if err := config.DBu.First(&empUser, "emp_id = ?", empID).Error; err != nil {
		return models.MasUser{}
	}
	return empUser
}
