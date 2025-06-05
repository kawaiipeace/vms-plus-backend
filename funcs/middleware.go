package funcs

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"vms_plus_be/config"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func ApiKeyMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		apiKey := c.GetHeader("X-ApiKey")

		if config.AppConfig.IsDev {
			apiKey = config.AppConfig.ApiKey
		}

		if apiKey == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing API key"})
			c.Abort()
			return
		}

		if apiKey != config.AppConfig.ApiKey {
			c.JSON(http.StatusForbidden, gin.H{"error": "Invalid API key"})
			c.Abort()
			return
		}

		c.Next()
	}
}

func ApiKeyAuthenMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if config.AppConfig.IsDev {
			c.Next()
			return
		}
		apiKey := c.GetHeader("X-ApiKey")
		if config.AppConfig.IsDev {
			apiKey = config.AppConfig.ApiKey
		}
		if apiKey == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing API key"})
			c.Abort()
			return
		}

		if apiKey != config.AppConfig.ApiKey {
			c.JSON(http.StatusForbidden, gin.H{"error": "Invalid API key"})
			c.Abort()
			return
		}

		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing or invalid Authorization header"})
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("unexpected signing method")
			}
			return []byte(config.AppConfig.JWTSecret), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid JWT token", "message": config.AppConfig.JWTSecret})
			c.Abort()
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			ctx := context.WithValue(c.Request.Context(), config.ClaimsKey, claims)
			c.Request = c.Request.WithContext(ctx)
			c.Next()
			return
		}
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid JWT token claims"})
		c.Abort()
	}
}
