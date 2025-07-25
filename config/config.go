package config

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// DB instance
var DB *gorm.DB

//var DBu *gorm.DB

type ContextKey string

const ClaimsKey ContextKey = "claims"

// Config structure for application settings
type Config struct {
	AppName             string
	Version             string
	IsDev               bool
	Host                string
	Port                int
	LogLevel            string
	Dsn_DB              string
	Dsn_DBu             string
	ApiKey              string
	JWTSecret           string
	JwtAccessTokenTime  int
	JwtRefreshTokenTime int
	OtpExpired          int
	KeyCloakClientID    string
	KeyCloakSecret      string
	KeyCloakEndPoint    string
	ThaiIDClientID      string
	ThaiIDSecret        string
	ThaiIDEndPoint      string

	MinIoEndPoint  string
	MinIoAccessKey string
	MinIoSecretKey string
	MinIoNotUseSSL bool

	DevSaveFilePath string
	DevSaveFileUrl  string

	UserHubEndPoint              string
	UserHubServiceKey            string
	HrPlatformEndPoint           string
	PEANotificationEndPoint      string
	PEANotificationToken         string
	PEAWorkDNotificationEndPoint string
	PEAWorkDNotificationToken    string
}

// AppConfig is a globally accessible configuration variable
var AppConfig Config
var DefaultAvatarURL string

func InitConfig() {
	DefaultAvatarURL = "https://vms-plus.pea.co.th/images/user-avatar.jpg"
	// Load the configuration from config.json
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}

	AppConfig = Config{
		AppName:             "VMS_PLUS",
		Version:             "1.0.0",
		IsDev:               os.Getenv("ISDEV") == "true",
		Host:                os.Getenv("HOST"),
		Port:                getEnvAsInt("PORT", 28080),
		LogLevel:            os.Getenv("LOG_LEVEL"),
		Dsn_DB:              os.Getenv("DSN_DB"),
		Dsn_DBu:             os.Getenv("DSN_DB_USER"),
		ApiKey:              os.Getenv("API_KEY"),
		JWTSecret:           os.Getenv("JWT_SECRET"),
		JwtAccessTokenTime:  600,  // Default: 60 minutes
		JwtRefreshTokenTime: 1440, // Default: 1440 minutes
		OtpExpired:          1,    // Default: 1 minutes
		KeyCloakClientID:    os.Getenv("KEYCLOAK_CLIENT_ID"),
		KeyCloakSecret:      os.Getenv("KEYCLOAK_CLIENT_SECRET"),
		KeyCloakEndPoint:    os.Getenv("KEYCLOAK_END_POINT"),

		ThaiIDClientID: os.Getenv("THAIID_CLIENT_ID"),
		ThaiIDSecret:   os.Getenv("THAIID_CLIENT_SECRET"),
		ThaiIDEndPoint: os.Getenv("THAIID_END_POINT"),

		MinIoEndPoint:  os.Getenv("MINIO_END_POINT"),
		MinIoAccessKey: os.Getenv("MINIO_ACCESS_KEY"),
		MinIoSecretKey: os.Getenv("MINIO_SECRET_KEY"),
		MinIoNotUseSSL: os.Getenv("MINIO_NOT_USE_SSL") == "true",

		DevSaveFilePath:              os.Getenv("DEV_SAVE_FILE_PATH"),
		DevSaveFileUrl:               os.Getenv("DEV_SAVE_FILE_URL"),
		UserHubEndPoint:              os.Getenv("USER_HUB_END_POINT"),
		UserHubServiceKey:            os.Getenv("USER_HUB_SERVICE_KEY"),
		HrPlatformEndPoint:           os.Getenv("HR_PLATFORM_END_POINT"),
		PEANotificationEndPoint:      os.Getenv("PEA_NOTIFICATION_END_POINT"),
		PEANotificationToken:         os.Getenv("PEA_NOTIFICATION_TOKEN"),
		PEAWorkDNotificationEndPoint: os.Getenv("PEA_WORK_D_NOTIFICATION_END_POINT"),
		PEAWorkDNotificationToken:    os.Getenv("PEA_WORK_D_NOTIFICATION_TOKEN"),
	}
	fmt.Printf("load AppConfig: %s %d\n", AppConfig.AppName, AppConfig.Port)

}
func getEnvAsInt(key string, defaultValue int) int {
	valueStr := os.Getenv(key)
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return defaultValue
	}
	return value
}
func InitDB() {
	var err error
	DB, err = gorm.Open(postgres.Open(AppConfig.Dsn_DB), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database vms_plus:", err)
	}

}
