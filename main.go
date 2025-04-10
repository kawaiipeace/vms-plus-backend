package main

import (
	"log"
	"strconv"
	"time"
	"vms_plus_be/config"
	_ "vms_plus_be/docs"
	"vms_plus_be/funcs"
	"vms_plus_be/handlers"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title VMS_PLUS
// @version 1.0
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name X-ApiKey

// @securityDefinitions.apikey AuthorizationAuth
// @Description Bearer [your Authorization]
// @in header
// @name Authorization

func main() {
	config.InitDB()
	router := gin.Default()
	router.SetTrustedProxies([]string{"192.168.1.1", "192.168.1.2"})
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},                                         // Allowed domains
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},              // Allowed methods
		AllowHeaders:     []string{"Content-Type", "Authorization", "X-ApiKey"}, // Allowed headers
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,           // Allow cookies
		MaxAge:           12 * time.Hour, // Cache preflight request
	}))

	// Initialize handler
	loginHandler := handlers.LoginHandler{}
	router.POST("/api/login/request-otp", funcs.ApiKeyMiddleware(), loginHandler.RequestOTP)
	router.POST("/api/login/verify-otp", funcs.ApiKeyMiddleware(), loginHandler.VerifyOTP)
	router.POST("/api/login/refresh-token", funcs.ApiKeyMiddleware(), loginHandler.RefreshToken)
	router.POST("/api/login/request-keycloak", funcs.ApiKeyMiddleware(), loginHandler.RequestKeyCloak)
	router.POST("/api/login/authen-keycloak", funcs.ApiKeyMiddleware(), loginHandler.AuthenKeyCloak)
	router.POST("/api/login/request-thaiid", funcs.ApiKeyMiddleware(), loginHandler.RequestThaiID)
	router.POST("/api/login/authen-thaiid", funcs.ApiKeyMiddleware(), loginHandler.AuthenThaiID)
	router.GET("/api/login/profile", funcs.ApiKeyAuthenMiddleware(), loginHandler.Profile)
	router.GET("/api/logout", funcs.ApiKeyAuthenMiddleware(), loginHandler.Logout)

	bookingUserHandler := handlers.BookingUserHandler{}
	router.POST("/api/booking-user/create-request", funcs.ApiKeyAuthenMiddleware(), bookingUserHandler.CreateRequest)
	router.GET("/api/booking-user/request/:id", funcs.ApiKeyAuthenMiddleware(), bookingUserHandler.GetRequest)
	router.GET("/api/booking-user/requests", funcs.ApiKeyAuthenMiddleware(), bookingUserHandler.ListRequest)
	router.PUT("/api/booking-user/update-vehicle-user", funcs.ApiKeyAuthenMiddleware(), bookingUserHandler.UpdateVehicleUser)
	router.PUT("/api/booking-user/update-trip", funcs.ApiKeyAuthenMiddleware(), bookingUserHandler.UpdateTrip)
	router.PUT("/api/booking-user/update-pickup", funcs.ApiKeyAuthenMiddleware(), bookingUserHandler.UpdatePickup)
	router.PUT("/api/booking-user/update-document", funcs.ApiKeyAuthenMiddleware(), bookingUserHandler.UpdateDocument)
	router.PUT("/api/booking-user/update-cost", funcs.ApiKeyAuthenMiddleware(), bookingUserHandler.UpdateCost)
	router.PUT("/api/booking-user/update-vehicle-type", funcs.ApiKeyAuthenMiddleware(), bookingUserHandler.UpdateVehicleType)
	router.PUT("/api/booking-user/update-approver", funcs.ApiKeyAuthenMiddleware(), bookingUserHandler.UpdateApprover)
	router.GET("/api/booking-user/search-requests", funcs.ApiKeyAuthenMiddleware(), bookingUserHandler.SearchRequests)

	bookingApproverHandler := handlers.BookingApproverHandler{}
	router.GET("/api/booking-approver/search-requests", funcs.ApiKeyAuthenMiddleware(), bookingApproverHandler.SearchRequests)
	router.GET("/api/booking-approver/request/:id", funcs.ApiKeyAuthenMiddleware(), bookingApproverHandler.GetRequest)
	router.PUT("/api/booking-approver/update-sended-back", funcs.ApiKeyAuthenMiddleware(), bookingApproverHandler.UpdateSendedBack)
	router.PUT("/api/booking-approver/update-approved", funcs.ApiKeyAuthenMiddleware(), bookingApproverHandler.UpdateApproved)
	router.PUT("/api/booking-approver/update-canceled", funcs.ApiKeyAuthenMiddleware(), bookingApproverHandler.UpdateCanceled)

	bookinAdminHandler := handlers.BookingAdminHandler{}
	router.GET("/api/booking-admin/search-requests", funcs.ApiKeyAuthenMiddleware(), bookinAdminHandler.SearchRequests)
	router.GET("/api/booking-admin/request/:id", funcs.ApiKeyAuthenMiddleware(), bookinAdminHandler.GetRequest)
	router.PUT("/api/booking-admin/update-sended-back", funcs.ApiKeyAuthenMiddleware(), bookinAdminHandler.UpdateSendedBack)
	router.PUT("/api/booking-admin/update-approved", funcs.ApiKeyAuthenMiddleware(), bookinAdminHandler.UpdateApproved)
	router.PUT("/api/booking-admin/update-canceled", funcs.ApiKeyAuthenMiddleware(), bookinAdminHandler.UpdateCanceled)
	router.PUT("/api/booking-admin/update-vehicle-user", funcs.ApiKeyAuthenMiddleware(), bookinAdminHandler.UpdateVehicleUser)
	router.PUT("/api/booking-admin/update-trip", funcs.ApiKeyAuthenMiddleware(), bookinAdminHandler.UpdateTrip)
	router.PUT("/api/booking-admin/update-pickup", funcs.ApiKeyAuthenMiddleware(), bookinAdminHandler.UpdatePickup)
	router.PUT("/api/booking-admin/update-document", funcs.ApiKeyAuthenMiddleware(), bookinAdminHandler.UpdateDocument)
	router.PUT("/api/booking-admin/update-cost", funcs.ApiKeyAuthenMiddleware(), bookinAdminHandler.UpdateCost)
	router.PUT("/api/booking-admin/update-vehicle", funcs.ApiKeyAuthenMiddleware(), bookinAdminHandler.UpdateVehicle)
	router.PUT("/api/booking-admin/update-driver", funcs.ApiKeyAuthenMiddleware(), bookinAdminHandler.UpdateDriver)

	bookinAdminDeptHandler := handlers.BookingAdminDeptHandler{}
	router.GET("/api/booking-admin-dept/search-requests", funcs.ApiKeyAuthenMiddleware(), bookinAdminDeptHandler.SearchRequests)
	router.GET("/api/booking-admin-dept/request/:id", funcs.ApiKeyAuthenMiddleware(), bookinAdminDeptHandler.GetRequest)
	router.PUT("/api/booking-admin-dept/update-sended-back", funcs.ApiKeyAuthenMiddleware(), bookinAdminDeptHandler.UpdateSendedBack)
	router.PUT("/api/booking-admin-dept/update-approved", funcs.ApiKeyAuthenMiddleware(), bookinAdminDeptHandler.UpdateApproved)
	router.PUT("/api/booking-admin-dept/update-canceled", funcs.ApiKeyAuthenMiddleware(), bookinAdminDeptHandler.UpdateCanceled)
	router.PUT("/api/booking-admin-dept/update-vehicle-user", funcs.ApiKeyAuthenMiddleware(), bookinAdminDeptHandler.UpdateVehicleUser)
	router.PUT("/api/booking-admin-dept/update-trip", funcs.ApiKeyAuthenMiddleware(), bookinAdminDeptHandler.UpdateTrip)
	router.PUT("/api/booking-admin-dept/update-pickup", funcs.ApiKeyAuthenMiddleware(), bookinAdminDeptHandler.UpdatePickup)
	router.PUT("/api/booking-admin-dept/update-document", funcs.ApiKeyAuthenMiddleware(), bookinAdminDeptHandler.UpdateDocument)
	router.PUT("/api/booking-admin-dept/update-cost", funcs.ApiKeyAuthenMiddleware(), bookinAdminDeptHandler.UpdateCost)
	router.PUT("/api/booking-admin-dept/update-vehicle", funcs.ApiKeyAuthenMiddleware(), bookinAdminDeptHandler.UpdateVehicle)
	router.PUT("/api/booking-admin-dept/update-driver", funcs.ApiKeyAuthenMiddleware(), bookinAdminDeptHandler.UpdateDriver)

	bookingFinalHandler := handlers.BookingFinalHandler{}
	router.GET("/api/booking-final/search-requests", funcs.ApiKeyAuthenMiddleware(), bookingFinalHandler.SearchRequests)
	router.GET("/api/booking-final/request/:id", funcs.ApiKeyAuthenMiddleware(), bookingFinalHandler.GetRequest)
	router.PUT("/api/booking-final/update-sended-back", funcs.ApiKeyAuthenMiddleware(), bookingFinalHandler.UpdateSendedBack)
	router.PUT("/api/booking-final/update-approved", funcs.ApiKeyAuthenMiddleware(), bookingFinalHandler.UpdateApproved)
	router.PUT("/api/booking-final/update-canceled", funcs.ApiKeyAuthenMiddleware(), bookingFinalHandler.UpdateCanceled)

	receivedKeyHandler := handlers.ReceivedKeyHandler{}
	router.GET("/api/received-key/search-requests", funcs.ApiKeyAuthenMiddleware(), receivedKeyHandler.SearchRequests)
	router.GET("/api/received-key/request/:id", funcs.ApiKeyAuthenMiddleware(), receivedKeyHandler.GetRequest)
	router.PUT("/api/received-key/update-key-pickup-emp", funcs.ApiKeyAuthenMiddleware(), receivedKeyHandler.UpdateKeyPickup_Emp)
	router.PUT("/api/received-key/update-key-pickup-outsource", funcs.ApiKeyAuthenMiddleware(), receivedKeyHandler.UpdateKeyPickup_OutSource)
	router.PUT("/api/received-key/update-canceled", funcs.ApiKeyAuthenMiddleware(), receivedKeyHandler.UpdateCanceled)

	receivedVehicleHandler := handlers.ReceivedVehicleHandler{}
	router.GET("/api/received-vehicle/search-requests", funcs.ApiKeyAuthenMiddleware(), receivedVehicleHandler.SearchRequests)
	router.GET("/api/received-vehicle/request/:id", funcs.ApiKeyAuthenMiddleware(), receivedVehicleHandler.GetRequest)
	router.PUT("/api/received-vehicle/update-vehicle-pickup", funcs.ApiKeyAuthenMiddleware(), receivedVehicleHandler.UpdateVehiclePickup)

	vehicleInUseHandler := handlers.VehicleInUseHandler{}
	router.GET("/api/vehicle-in-use/search-requests", funcs.ApiKeyAuthenMiddleware(), vehicleInUseHandler.SearchRequests)
	router.GET("/api/vehicle-in-use/request/:id", funcs.ApiKeyAuthenMiddleware(), vehicleInUseHandler.GetRequest)
	router.GET("/api/vehicle-in-use/travel-details", funcs.ApiKeyAuthenMiddleware(), vehicleInUseHandler.GetVehicleTripDetails)
	router.GET("/api/vehicle-in-use/travel-detail/:id", funcs.ApiKeyAuthenMiddleware(), vehicleInUseHandler.GetVehicleTripDetail)
	router.POST("/api/vehicle-in-use/create-travel-detail", funcs.ApiKeyAuthenMiddleware(), vehicleInUseHandler.CreateVehicleTripDetail)
	router.PUT("/api/vehicle-in-use/update-travel-detail/:id", funcs.ApiKeyAuthenMiddleware(), vehicleInUseHandler.UpdateVehicleTripDetail)
	router.DELETE("/api/vehicle-in-use/delete-travel-detail/:id", funcs.ApiKeyAuthenMiddleware(), vehicleInUseHandler.DeleteVehicleTripDetail)
	router.GET("/api/vehicle-in-use/add-fuel-details", funcs.ApiKeyAuthenMiddleware(), vehicleInUseHandler.GetVehicleAddFuelDetails)
	router.GET("/api/vehicle-in-use/add-fuel-detail/:id", funcs.ApiKeyAuthenMiddleware(), vehicleInUseHandler.GetVehicleAddFuelDetail)
	router.POST("/api/vehicle-in-use/create-add-fuel", funcs.ApiKeyAuthenMiddleware(), vehicleInUseHandler.CreateVehicleAddFuel)
	router.PUT("/api/vehicle-in-use/update-add-fuel/:id", funcs.ApiKeyAuthenMiddleware(), vehicleInUseHandler.UpdateVehicleAddFuel)
	router.DELETE("/api/vehicle-in-use/delete-add-fuel/:id", funcs.ApiKeyAuthenMiddleware(), vehicleInUseHandler.DeleteVehicleAddFuel)

	vehicleHandler := handlers.VehicleHandler{}
	router.GET("/api/vehicle/search", funcs.ApiKeyAuthenMiddleware(), vehicleHandler.SearchVehicles)
	router.GET("/api/vehicle/types", funcs.ApiKeyAuthenMiddleware(), vehicleHandler.GetTypes)
	router.GET("/api/vehicle/departments", funcs.ApiKeyAuthenMiddleware(), vehicleHandler.GetDepartments)
	router.GET("/api/vehicle/:id", funcs.ApiKeyAuthenMiddleware(), vehicleHandler.GetVehicle)
	router.GET("/api/vehicle-info/:id", funcs.ApiKeyAuthenMiddleware(), vehicleHandler.GetVehicleInfo)

	driverHandler := handlers.DriverHandler{}
	router.GET("/api/driver/search", funcs.ApiKeyAuthenMiddleware(), driverHandler.GetDrivers)
	router.GET("/api/driver/:id", funcs.ApiKeyAuthenMiddleware(), driverHandler.GetDriver)
	router.GET("/api/driver/search-other-dept", funcs.ApiKeyAuthenMiddleware(), driverHandler.GetDriversOtherDept)

	masHandler := handlers.MasHandler{}
	router.GET("/api/mas/user-vehicle-users", funcs.ApiKeyAuthenMiddleware(), masHandler.ListVehicleUser)
	router.GET("/api/mas/user-driver-users", funcs.ApiKeyAuthenMiddleware(), masHandler.ListDriverUser)
	router.GET("/api/mas/user-approval-users", funcs.ApiKeyAuthenMiddleware(), masHandler.ListApprovalUser)
	router.GET("/api/mas/user-admin-approval-users", funcs.ApiKeyAuthenMiddleware(), masHandler.ListAdminApprovalUser)
	router.GET("/api/mas/user-final-approval-users", funcs.ApiKeyAuthenMiddleware(), masHandler.ListFinalApprovalUser)
	router.GET("/api/mas/user/:id", funcs.ApiKeyAuthenMiddleware(), masHandler.GetUserEmp)

	refHandler := handlers.RefHandler{}
	router.GET("/api/ref/cost-type", funcs.ApiKeyAuthenMiddleware(), refHandler.ListCostType)
	router.GET("/api/ref/cost-type/:code", funcs.ApiKeyAuthenMiddleware(), refHandler.GetCostType)
	router.GET("/api/ref/request-status", funcs.ApiKeyAuthenMiddleware(), refHandler.ListRequestStatus)
	router.GET("/api/ref/fuel-type", funcs.ApiKeyAuthenMiddleware(), refHandler.ListFuelType)
	router.GET("/api/ref/oil-station-brand", funcs.ApiKeyAuthenMiddleware(), refHandler.ListOilStationBrand)

	logHandler := handlers.LogHandler{}
	router.GET("/api/log/request/:id", funcs.ApiKeyAuthenMiddleware(), logHandler.GetLogRequest)

	uploadHandler := handlers.UploadHandler{}
	router.POST("/api/upload", funcs.ApiKeyMiddleware(), uploadHandler.UploadFile)

	// Swagger documentation
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Start server
	port := strconv.Itoa(config.AppConfig.Port)
	log.Println("Server started at " + config.AppConfig.Host + ":" + port)

	router.Run(config.AppConfig.Host + ":" + port)
}
