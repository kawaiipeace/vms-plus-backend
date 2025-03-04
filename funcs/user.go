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
	EmpID     string `json:"emp_id"`
	FullName  string `json:"full_name"`
	TokenType string `json:"token_type"`
	Role      string `json:"role"`
	jwt.RegisteredClaims
}

var (
	jwtSecret = []byte(config.AppConfig.JWTSecret)
)

func GenerateJWT(user models.AuthenUserEmp, tokenType string, expiration time.Duration) (string, error) {
	claims := Claims{
		EmpID:     user.EmpID,
		FullName:  user.FirstName + " " + user.LastName,
		TokenType: tokenType,
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

	// Map claims to UserEmp struct
	user := &models.AuthenJwtUsr{
		EmpID:    claims["emp_id"].(string),
		FullName: claims["full_name"].(string),
	}

	return user, nil
}

func GetAuthenUser(c *gin.Context, role string) *models.AuthenUserEmp {
	// Extract user from JWT
	jwt, err := ExtractUserFromJWT(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		c.Abort()
	}

	var user models.AuthenUserEmp
	result := config.DB.Where("emp_id = ?", jwt.EmpID).First(&user)
	if result.Error != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		c.Abort()
	}

	//check role

	return &user
}
