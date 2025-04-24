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
	jwt.RegisteredClaims
}

var (
	jwtSecret = []byte(config.AppConfig.JWTSecret)
)

func GenerateJWT(user models.AuthenUserEmp, tokenType string, expiration time.Duration, accessToken string, refreshToken string) (string, error) {
	claims := Claims{
		EmpID:     user.EmpID,
		FullName:  user.FirstName + " " + user.LastName,
		TokenType: tokenType,
		Roles:     []string{"vehicel-user", "level1-approval", "admin-approval", "final-approval"},
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

func GetAuthenUser(c *gin.Context, role string) *models.AuthenUserEmp {
	if config.AppConfig.IsDev {
		duser := models.AuthenUserEmp{
			EmpID: "700001",
		}
		return &duser
	}

	// Extract user from JWT
	var empUser models.AuthenUserEmp

	if config.AppConfig.IsDev {
		if err := config.DB.First(&empUser, "emp_id = ?", "700001").Error; err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			c.Abort()
			return &empUser
		}
		empUser.Roles = []string{"vehicel-user", "level1-approval", "admin-approval", "final-approval"}
		return &empUser
	}
	jwt, err := ExtractUserFromJWT(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		c.Abort()
	}

	if err := config.DB.First(&empUser, "emp_id = ?", jwt.EmpID).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		c.Abort()
		return &empUser
	}
	empUser.Roles = jwt.Roles
	//check role

	return &empUser
}

func GetUserEmpInfo(empID string) models.AuthenUserEmp {
	var empUser models.AuthenUserEmp
	if err := config.DB.First(&empUser, "emp_id = ?", empID).Error; err != nil {
		return models.AuthenUserEmp{}
	}
	return empUser
}
