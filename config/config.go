package config

import (
	"encoding/json"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// DB instance
var DB *gorm.DB

type ContextKey string

const ClaimsKey ContextKey = "claims"

// Config structure for application settings
type Config struct {
	AppName             string `json:"app_name"`
	Version             string `json:"version"`
	IsDev               bool   `json:"isdev"`
	Host                string `json:"host"`
	Port                int    `json:"port"`
	LogLevel            string `json:"log_level"`
	Dsn_DB              string `json:"dsn_db"`
	ApiKey              string `json:"api_key"`
	JWTSecret           string `json:"jwt_secret"`
	JwtAccessTokenTime  int    `json:"jwt_access_token_time"`
	JwtRefreshTokenTime int    `json:"jwt_refresh_token_time"`
	OtpExpired          int    `json:"otp_expired"`
	KeyCloakClientID    string `json:"keycloak_client_id"`
	KeyCloakSecret      string `json:"keycloak_client_secret"`
	KeyCloakEndPoint    string `json:"keycloak_end_point"`
	ThaiIDClientID      string `json:"thaiid_client_id"`
	ThaiIDSecret        string `json:"thaiid_client_secret"`
	ThaiIDEndPoint      string `json:"thaiid_end_point"`
	SaveFilePath        string `json:"save_file_path"`
	SaveFileUrl         string `json:"save_file_url"`
}

// AppConfig is a globally accessible configuration variable
var AppConfig Config

func init() {
	// Load the configuration from config.json
	data, err := os.ReadFile("config.json")
	if err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}

	// Unmarshal the JSON into the AppConfig variable
	err = json.Unmarshal(data, &AppConfig)
	if err != nil {
		log.Fatalf("Error parsing config file: %v", err)
	}
}

func InitDB() {
	var err error
	DB, err = gorm.Open(postgres.Open(AppConfig.Dsn_DB), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	log.Println("Database connected")
}
