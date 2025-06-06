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

// @securityDefinitions.apikey ServiceKey
// @in header
// @name ServiceKey

func main() {
	config.InitConfig()
	config.InitDB()
	handlers.InitMinIO(config.AppConfig.MinIoEndPoint, config.AppConfig.MinIoAccessKey, config.AppConfig.MinIoSecretKey, true)

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
	//LoginHandler
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

	//VehicleHandler
	vehicleHandler := handlers.VehicleHandler{Role: "*"}
	router.GET("/api/vehicle/search", funcs.ApiKeyAuthenMiddleware(), vehicleHandler.SearchVehicles)
	router.GET("/api/vehicle/search-booking", funcs.ApiKeyAuthenMiddleware(), vehicleHandler.SearchBookingVehicles)
	router.GET("/api/vehicle/types", funcs.ApiKeyAuthenMiddleware(), vehicleHandler.GetTypes)
	router.GET("/api/vehicle/departments", funcs.ApiKeyAuthenMiddleware(), vehicleHandler.GetDepartments)
	router.GET("/api/vehicle/:mas_vehicle_uid", funcs.ApiKeyAuthenMiddleware(), vehicleHandler.GetVehicle)
	router.GET("/api/vehicle-info/:mas_vehicle_uid", funcs.ApiKeyAuthenMiddleware(), vehicleHandler.GetVehicleInfo)
	router.GET("/api/vehicle/car-types-by-detail", funcs.ApiKeyAuthenMiddleware(), vehicleHandler.GetCarTypeDetails)

	//DriverHandler
	driverHandler := handlers.DriverHandler{Role: "*"}
	router.GET("/api/driver/search", funcs.ApiKeyAuthenMiddleware(), driverHandler.GetDrivers)
	router.GET("/api/driver/search-booking", funcs.ApiKeyAuthenMiddleware(), driverHandler.GetBookingDrivers)
	router.GET("/api/driver/:mas_driver_uid", funcs.ApiKeyAuthenMiddleware(), driverHandler.GetDriver)
	router.GET("/api/driver/search-other-dept", funcs.ApiKeyAuthenMiddleware(), driverHandler.GetDriversOtherDept)
	router.GET("/api/driver/work-type", funcs.ApiKeyAuthenMiddleware(), driverHandler.GetWorkType)

	//BookingUserHandler
	bookingUserHandler := handlers.BookingUserHandler{Role: "vehicle-user"}
	router.GET("/api/booking-user/menu-requests", funcs.ApiKeyAuthenMiddleware(), bookingUserHandler.MenuRequests)
	router.POST("/api/booking-user/create-request", funcs.ApiKeyAuthenMiddleware(), bookingUserHandler.CreateRequest)
	router.GET("/api/booking-user/request/:trn_request_uid", funcs.ApiKeyAuthenMiddleware(), bookingUserHandler.GetRequest)
	router.PUT("/api/booking-user/update-vehicle-user", funcs.ApiKeyAuthenMiddleware(), bookingUserHandler.UpdateVehicleUser)
	router.PUT("/api/booking-user/update-trip", funcs.ApiKeyAuthenMiddleware(), bookingUserHandler.UpdateTrip)
	router.PUT("/api/booking-user/update-pickup", funcs.ApiKeyAuthenMiddleware(), bookingUserHandler.UpdatePickup)
	router.PUT("/api/booking-user/update-document", funcs.ApiKeyAuthenMiddleware(), bookingUserHandler.UpdateDocument)
	router.PUT("/api/booking-user/update-cost", funcs.ApiKeyAuthenMiddleware(), bookingUserHandler.UpdateCost)
	router.PUT("/api/booking-user/update-vehicle-type", funcs.ApiKeyAuthenMiddleware(), bookingUserHandler.UpdateVehicleType)
	router.PUT("/api/booking-user/update-confirmer", funcs.ApiKeyAuthenMiddleware(), bookingUserHandler.UpdateConfirmer)
	router.GET("/api/booking-user/search-requests", funcs.ApiKeyAuthenMiddleware(), bookingUserHandler.SearchRequests)
	router.PUT("/api/booking-user/update-canceled", funcs.ApiKeyAuthenMiddleware(), bookingUserHandler.UpdateCanceled)
	router.PUT("/api/booking-user/update-resend", funcs.ApiKeyAuthenMiddleware(), bookingUserHandler.UpdateResend)

	//BookingConfirmerHandler
	bookingConfirmerHandler := handlers.BookingConfirmerHandler{Role: "level1-approval"}
	router.GET("/api/booking-confirmer/menu-requests", funcs.ApiKeyAuthenMiddleware(), bookingConfirmerHandler.MenuRequests)
	router.GET("/api/booking-confirmer/mmenu-requests", funcs.ApiKeyAuthenMiddleware(), bookingConfirmerHandler.MenuRequests)
	router.GET("/api/booking-confirmer/search-requests", funcs.ApiKeyAuthenMiddleware(), bookingConfirmerHandler.SearchRequests)
	router.GET("/api/booking-confirmer/request/:trn_request_uid", funcs.ApiKeyAuthenMiddleware(), bookingConfirmerHandler.GetRequest)
	router.PUT("/api/booking-confirmer/update-rejected", funcs.ApiKeyAuthenMiddleware(), bookingConfirmerHandler.UpdateRejected)
	router.PUT("/api/booking-confirmer/update-approved", funcs.ApiKeyAuthenMiddleware(), bookingConfirmerHandler.UpdateApproved)
	router.PUT("/api/booking-confirmer/update-canceled", funcs.ApiKeyAuthenMiddleware(), bookingConfirmerHandler.UpdateCanceled)

	//BookingAdminHandler
	bookinAdminHandler := handlers.BookingAdminHandler{Role: "admin-approval"}
	router.GET("/api/booking-admin/menu-requests", funcs.ApiKeyAuthenMiddleware(), bookinAdminHandler.MenuRequests)
	router.GET("/api/booking-admin/search-requests", funcs.ApiKeyAuthenMiddleware(), bookinAdminHandler.SearchRequests)
	router.GET("/api/booking-admin/request/:trn_request_uid", funcs.ApiKeyAuthenMiddleware(), bookinAdminHandler.GetRequest)
	router.PUT("/api/booking-admin/update-sended-back", funcs.ApiKeyAuthenMiddleware(), bookinAdminHandler.UpdateRejected)
	router.PUT("/api/booking-admin/update-approved", funcs.ApiKeyAuthenMiddleware(), bookinAdminHandler.UpdateApproved)
	router.PUT("/api/booking-admin/update-canceled", funcs.ApiKeyAuthenMiddleware(), bookinAdminHandler.UpdateCanceled)
	router.PUT("/api/booking-admin/update-rejected", funcs.ApiKeyAuthenMiddleware(), bookinAdminHandler.UpdateRejected)
	router.PUT("/api/booking-admin/update-vehicle-user", funcs.ApiKeyAuthenMiddleware(), bookinAdminHandler.UpdateVehicleUser)
	router.PUT("/api/booking-admin/update-trip", funcs.ApiKeyAuthenMiddleware(), bookinAdminHandler.UpdateTrip)
	router.PUT("/api/booking-admin/update-pickup", funcs.ApiKeyAuthenMiddleware(), bookinAdminHandler.UpdatePickup)
	router.PUT("/api/booking-admin/update-document", funcs.ApiKeyAuthenMiddleware(), bookinAdminHandler.UpdateDocument)
	router.PUT("/api/booking-admin/update-cost", funcs.ApiKeyAuthenMiddleware(), bookinAdminHandler.UpdateCost)
	router.PUT("/api/booking-admin/update-vehicle", funcs.ApiKeyAuthenMiddleware(), bookinAdminHandler.UpdateVehicle)
	router.PUT("/api/booking-admin/update-driver", funcs.ApiKeyAuthenMiddleware(), bookinAdminHandler.UpdateDriver)

	//BookingFinalHandler
	bookingFinalHandler := handlers.BookingFinalHandler{Role: "final-approval"}
	router.GET("/api/booking-final/menu-requests", funcs.ApiKeyAuthenMiddleware(), bookingFinalHandler.MenuRequests)
	router.GET("/api/booking-final/search-requests", funcs.ApiKeyAuthenMiddleware(), bookingFinalHandler.SearchRequests)
	router.GET("/api/booking-final/request/:trn_request_uid", funcs.ApiKeyAuthenMiddleware(), bookingFinalHandler.GetRequest)
	router.PUT("/api/booking-final/update-rejected", funcs.ApiKeyAuthenMiddleware(), bookingFinalHandler.UpdateRejected)
	router.PUT("/api/booking-final/update-approved", funcs.ApiKeyAuthenMiddleware(), bookingFinalHandler.UpdateApproved)
	router.PUT("/api/booking-final/update-canceled", funcs.ApiKeyAuthenMiddleware(), bookingFinalHandler.UpdateCanceled)

	//ReceivedKeyUserHandler
	receivedKeyUserHandler := handlers.ReceivedKeyUserHandler{Role: "vehicle-user"}
	router.GET("/api/received-key-user/search-requests", funcs.ApiKeyAuthenMiddleware(), receivedKeyUserHandler.SearchRequests)
	router.GET("/api/received-key-user/request/:trn_request_uid", funcs.ApiKeyAuthenMiddleware(), receivedKeyUserHandler.GetRequest)
	router.PUT("/api/received-key-user/update-key-pickup-pea", funcs.ApiKeyAuthenMiddleware(), receivedKeyUserHandler.UpdateKeyPickupPEA)
	router.PUT("/api/received-key-user/update-key-pickup-outsider", funcs.ApiKeyAuthenMiddleware(), receivedKeyUserHandler.UpdateKeyPickupOutSider)
	router.PUT("/api/received-key-user/update-key-pickup-driver", funcs.ApiKeyAuthenMiddleware(), receivedKeyUserHandler.UpdateKeyPickupDriver)
	router.PUT("/api/received-key-user/update-canceled", funcs.ApiKeyAuthenMiddleware(), receivedKeyUserHandler.UpdateCanceled)
	router.PUT("/api/received-key-user/update-recieived-key-confirmed", funcs.ApiKeyAuthenMiddleware(), receivedKeyUserHandler.UpdateRecieivedKeyConfirmed)

	//ReceivedKeyAdminHandler
	receivedKeyAdminHandler := handlers.ReceivedKeyAdminHandler{Role: "admin-approval,admin-dept-approval"}
	router.GET("/api/received-key-admin/search-requests", funcs.ApiKeyAuthenMiddleware(), receivedKeyAdminHandler.SearchRequests)
	router.GET("/api/received-key-admin/request/:trn_request_uid", funcs.ApiKeyAuthenMiddleware(), receivedKeyAdminHandler.GetRequest)
	router.PUT("/api/received-key-admin/update-recieived-key", funcs.ApiKeyAuthenMiddleware(), receivedKeyAdminHandler.UpdateRecieivedKey)
	router.PUT("/api/received-key-admin/update-key-pickup-pea", funcs.ApiKeyAuthenMiddleware(), receivedKeyAdminHandler.UpdateKeyPickupPEA)
	router.PUT("/api/received-key-admin/update-key-pickup-outsider", funcs.ApiKeyAuthenMiddleware(), receivedKeyAdminHandler.UpdateKeyPickupOutSider)
	router.PUT("/api/received-key-admin/update-key-pickup-driver", funcs.ApiKeyAuthenMiddleware(), receivedKeyAdminHandler.UpdateKeyPickupDriver)
	router.PUT("/api/received-key-admin/update-canceled", funcs.ApiKeyAuthenMiddleware(), receivedKeyAdminHandler.UpdateCanceled)
	router.PUT("/api/received-key-admin/update-recieived-key-detail", funcs.ApiKeyAuthenMiddleware(), receivedKeyAdminHandler.UpdateRecieivedKeyDetail)

	//ReceivedKeyDriverHandler
	receivedKeyDriverHandler := handlers.ReceivedKeyDriverHandler{Role: "driver"}
	router.GET("/api/received-key-driver/search-requests", funcs.ApiKeyAuthenMiddleware(), receivedKeyDriverHandler.SearchRequests)
	router.GET("/api/received-key-driver/request/:trn_request_uid", funcs.ApiKeyAuthenMiddleware(), receivedKeyDriverHandler.GetRequest)
	router.PUT("/api/received-key-driver/update-recieived-key-confirmed", funcs.ApiKeyAuthenMiddleware(), receivedKeyDriverHandler.UpdateRecieivedKeyConfirmed)
	router.GET("/api/booking-driver/menu-requests", funcs.ApiKeyAuthenMiddleware(), receivedKeyDriverHandler.MenuRequests)

	//ReceivedVehicleUserHandler
	receivedVehicleUserHandler := handlers.ReceivedVehicleUserHandler{Role: "vehicle-user"}
	router.GET("/api/received-vehicle-user/search-requests", funcs.ApiKeyAuthenMiddleware(), receivedVehicleUserHandler.SearchRequests)
	router.GET("/api/received-vehicle-user/request/:trn_request_uid", funcs.ApiKeyAuthenMiddleware(), receivedVehicleUserHandler.GetRequest)
	router.PUT("/api/received-vehicle-user/received-vehicle", funcs.ApiKeyAuthenMiddleware(), receivedVehicleUserHandler.ReceivedVehicle)
	router.GET("/api/received-vehicle-user/travel-card/:trn_request_uid", funcs.ApiKeyAuthenMiddleware(), receivedVehicleUserHandler.GetTravelCard)

	//ReceivedVehicleAdminHandler
	receivedVehicleAdminHandler := handlers.ReceivedVehicleAdminHandler{Role: "admin-approval,admin-dept-approval"}
	router.GET("/api/received-vehicle-admin/search-requests", funcs.ApiKeyAuthenMiddleware(), receivedVehicleAdminHandler.SearchRequests)
	router.GET("/api/received-vehicle-admin/request/:trn_request_uid", funcs.ApiKeyAuthenMiddleware(), receivedVehicleAdminHandler.GetRequest)
	router.PUT("/api/received-vehicle-admin/received-vehicle", funcs.ApiKeyAuthenMiddleware(), receivedVehicleAdminHandler.ReceivedVehicle)
	router.GET("/api/received-vehicle-admin/travel-card/:trn_request_uid", funcs.ApiKeyAuthenMiddleware(), receivedVehicleAdminHandler.GetTravelCard)

	//ReceivedVehicleDriverHandler
	receivedVehicleDriverHandler := handlers.ReceivedVehicleDriverHandler{Role: "driver"}
	router.GET("/api/received-vehicle-driver/search-requests", funcs.ApiKeyAuthenMiddleware(), receivedVehicleDriverHandler.SearchRequests)
	router.GET("/api/received-vehicle-driver/request/:trn_request_uid", funcs.ApiKeyAuthenMiddleware(), receivedVehicleDriverHandler.GetRequest)
	router.PUT("/api/received-vehicle-driver/received-vehicle", funcs.ApiKeyAuthenMiddleware(), receivedVehicleDriverHandler.ReceivedVehicle)
	router.GET("/api/received-vehicle-driver/travel-card/:trn_request_uid", funcs.ApiKeyAuthenMiddleware(), receivedVehicleDriverHandler.GetTravelCard)

	//VehicleInUseUserHandler
	vehicleInUseUserHandler := handlers.VehicleInUseUserHandler{Role: "vehicle-user"}
	router.GET("/api/vehicle-in-use-user/search-requests", funcs.ApiKeyAuthenMiddleware(), vehicleInUseUserHandler.SearchRequests)
	router.GET("/api/vehicle-in-use-user/request/:trn_request_uid", funcs.ApiKeyAuthenMiddleware(), vehicleInUseUserHandler.GetRequest)
	router.GET("/api/vehicle-in-use-user/travel-details/:trn_request_uid", funcs.ApiKeyAuthenMiddleware(), vehicleInUseUserHandler.GetVehicleTripDetails)
	router.GET("/api/vehicle-in-use-user/travel-detail/:trn_trip_detail_uid", funcs.ApiKeyAuthenMiddleware(), vehicleInUseUserHandler.GetVehicleTripDetail)
	router.POST("/api/vehicle-in-use-user/create-travel-detail", funcs.ApiKeyAuthenMiddleware(), vehicleInUseUserHandler.CreateVehicleTripDetail)
	router.PUT("/api/vehicle-in-use-user/update-travel-detail/:trn_trip_detail_uid", funcs.ApiKeyAuthenMiddleware(), vehicleInUseUserHandler.UpdateVehicleTripDetail)
	router.DELETE("/api/vehicle-in-use-user/delete-travel-detail/:trn_trip_detail_uid", funcs.ApiKeyAuthenMiddleware(), vehicleInUseUserHandler.DeleteVehicleTripDetail)
	router.GET("/api/vehicle-in-use-user/add-fuel-details/:trn_request_uid", funcs.ApiKeyAuthenMiddleware(), vehicleInUseUserHandler.GetVehicleAddFuelDetails)
	router.GET("/api/vehicle-in-use-user/add-fuel-detail/:trn_add_fuel_uid", funcs.ApiKeyAuthenMiddleware(), vehicleInUseUserHandler.GetVehicleAddFuelDetail)
	router.POST("/api/vehicle-in-use-user/create-add-fuel", funcs.ApiKeyAuthenMiddleware(), vehicleInUseUserHandler.CreateVehicleAddFuel)
	router.PUT("/api/vehicle-in-use-user/update-add-fuel/:trn_add_fuel_uid", funcs.ApiKeyAuthenMiddleware(), vehicleInUseUserHandler.UpdateVehicleAddFuel)
	router.DELETE("/api/vehicle-in-use-user/delete-add-fuel/:trn_add_fuel_uid", funcs.ApiKeyAuthenMiddleware(), vehicleInUseUserHandler.DeleteVehicleAddFuel)
	router.GET("/api/vehicle-in-use-user/travel-card/:trn_request_uid", funcs.ApiKeyAuthenMiddleware(), vehicleInUseUserHandler.GetTravelCard)
	router.PUT("/api/vehicle-in-use-user/returned-vehicle", funcs.ApiKeyAuthenMiddleware(), vehicleInUseUserHandler.ReturnedVehicle)
	router.PUT("/api/vehicle-in-use-user/update-satisfaction-survey/:trn_request_uid", funcs.ApiKeyAuthenMiddleware(), vehicleInUseUserHandler.UpdateSatisfactionSurvey)

	//VehicleInUseAdminHandler
	vehicleInUseAdminHandler := handlers.VehicleInUseAdminHandler{Role: "admin-approval,admin-dept-approval"}
	router.GET("/api/vehicle-in-use-admin/search-requests", funcs.ApiKeyAuthenMiddleware(), vehicleInUseAdminHandler.SearchRequests)
	router.GET("/api/vehicle-in-use-admin/request/:trn_request_uid", funcs.ApiKeyAuthenMiddleware(), vehicleInUseAdminHandler.GetRequest)
	router.GET("/api/vehicle-in-use-admin/travel-details/:trn_request_uid", funcs.ApiKeyAuthenMiddleware(), vehicleInUseAdminHandler.GetVehicleTripDetails)
	router.GET("/api/vehicle-in-use-admin/travel-detail/:trn_trip_detail_uid", funcs.ApiKeyAuthenMiddleware(), vehicleInUseAdminHandler.GetVehicleTripDetail)
	router.POST("/api/vehicle-in-use-admin/create-travel-detail", funcs.ApiKeyAuthenMiddleware(), vehicleInUseAdminHandler.CreateVehicleTripDetail)
	router.PUT("/api/vehicle-in-use-admin/update-travel-detail/:trn_trip_detail_uid", funcs.ApiKeyAuthenMiddleware(), vehicleInUseAdminHandler.UpdateVehicleTripDetail)
	router.DELETE("/api/vehicle-in-use-admin/delete-travel-detail/:trn_trip_detail_uid", funcs.ApiKeyAuthenMiddleware(), vehicleInUseAdminHandler.DeleteVehicleTripDetail)
	router.GET("/api/vehicle-in-use-admin/add-fuel-details/:trn_request_uid", funcs.ApiKeyAuthenMiddleware(), vehicleInUseAdminHandler.GetVehicleAddFuelDetails)
	router.GET("/api/vehicle-in-use-admin/add-fuel-detail/:trn_add_fuel_uid", funcs.ApiKeyAuthenMiddleware(), vehicleInUseAdminHandler.GetVehicleAddFuelDetail)
	router.POST("/api/vehicle-in-use-admin/create-add-fuel", funcs.ApiKeyAuthenMiddleware(), vehicleInUseAdminHandler.CreateVehicleAddFuel)
	router.PUT("/api/vehicle-in-use-admin/update-add-fuel/:trn_add_fuel_uid", funcs.ApiKeyAuthenMiddleware(), vehicleInUseAdminHandler.UpdateVehicleAddFuel)
	router.DELETE("/api/vehicle-in-use-admin/delete-add-fuel/:trn_add_fuel_uid", funcs.ApiKeyAuthenMiddleware(), vehicleInUseAdminHandler.DeleteVehicleAddFuel)
	router.GET("/api/vehicle-in-use-admin/travel-card/:trn_request_uid", funcs.ApiKeyAuthenMiddleware(), vehicleInUseAdminHandler.GetTravelCard)
	router.PUT("/api/vehicle-in-use-admin/returned-vehicle", funcs.ApiKeyAuthenMiddleware(), vehicleInUseAdminHandler.ReturnedVehicle)
	router.PUT("/api/vehicle-in-use-admin/update-received-vehicle", funcs.ApiKeyAuthenMiddleware(), vehicleInUseAdminHandler.UpdateReceivedVehicle)
	router.PUT("/api/vehicle-in-use-admin/update-received-vehicle-images", funcs.ApiKeyAuthenMiddleware(), vehicleInUseAdminHandler.UpdateReceivedVehicleImages)

	//VehicleInUseDriverHandler
	vehicleInUseDriverHandler := handlers.VehicleInUseDriverHandler{Role: "driver,vehicle-user,admin-approval"}
	router.GET("/api/vehicle-in-use-driver/search-requests", funcs.ApiKeyAuthenMiddleware(), vehicleInUseDriverHandler.SearchRequests)
	router.GET("/api/vehicle-in-use-driver/request/:trn_request_uid", funcs.ApiKeyAuthenMiddleware(), vehicleInUseDriverHandler.GetRequest)
	router.GET("/api/vehicle-in-use-driver/travel-details/:trn_request_uid", funcs.ApiKeyAuthenMiddleware(), vehicleInUseDriverHandler.GetVehicleTripDetails)
	router.GET("/api/vehicle-in-use-driver/travel-detail/:trn_trip_detail_uid", funcs.ApiKeyAuthenMiddleware(), vehicleInUseDriverHandler.GetVehicleTripDetail)
	router.POST("/api/vehicle-in-use-driver/create-travel-detail", funcs.ApiKeyAuthenMiddleware(), vehicleInUseDriverHandler.CreateVehicleTripDetail)
	router.PUT("/api/vehicle-in-use-driver/update-travel-detail/:trn_trip_detail_uid", funcs.ApiKeyAuthenMiddleware(), vehicleInUseDriverHandler.UpdateVehicleTripDetail)
	router.DELETE("/api/vehicle-in-use-driver/delete-travel-detail/:trn_trip_detail_uid", funcs.ApiKeyAuthenMiddleware(), vehicleInUseDriverHandler.DeleteVehicleTripDetail)
	router.GET("/api/vehicle-in-use-driver/add-fuel-details/:trn_request_uid", funcs.ApiKeyAuthenMiddleware(), vehicleInUseDriverHandler.GetVehicleAddFuelDetails)
	router.GET("/api/vehicle-in-use-driver/add-fuel-detail/:trn_add_fuel_uid", funcs.ApiKeyAuthenMiddleware(), vehicleInUseDriverHandler.GetVehicleAddFuelDetail)
	router.POST("/api/vehicle-in-use-driver/create-add-fuel", funcs.ApiKeyAuthenMiddleware(), vehicleInUseDriverHandler.CreateVehicleAddFuel)
	router.PUT("/api/vehicle-in-use-driver/update-add-fuel/:trn_add_fuel_uid", funcs.ApiKeyAuthenMiddleware(), vehicleInUseDriverHandler.UpdateVehicleAddFuel)
	router.DELETE("/api/vehicle-in-use-driver/delete-add-fuel/:trn_add_fuel_uid", funcs.ApiKeyAuthenMiddleware(), vehicleInUseDriverHandler.DeleteVehicleAddFuel)
	router.GET("/api/vehicle-in-use-driver/travel-card/:trn_request_uid", funcs.ApiKeyAuthenMiddleware(), vehicleInUseDriverHandler.GetTravelCard)
	router.PUT("/api/vehicle-in-use-driver/returned-vehicle", funcs.ApiKeyAuthenMiddleware(), vehicleInUseDriverHandler.ReturnedVehicle)
	router.PUT("/api/vehicle-in-use-driver/update-received-vehicle", funcs.ApiKeyAuthenMiddleware(), vehicleInUseDriverHandler.UpdateReceivedVehicle)
	router.PUT("/api/vehicle-in-use-driver/update-received-vehicle-images", funcs.ApiKeyAuthenMiddleware(), vehicleInUseDriverHandler.UpdateReceivedVehicleImages)

	//VehicleInspectionAdminHandler
	vehicleInspectionAdminHandler := handlers.VehicleInspectionAdminHandler{Role: "admin-approval"}
	router.GET("/api/vehicle-inspection-admin/search-requests", funcs.ApiKeyAuthenMiddleware(), vehicleInspectionAdminHandler.SearchRequests)
	router.GET("/api/vehicle-inspection-admin/request/:trn_request_uid", funcs.ApiKeyAuthenMiddleware(), vehicleInspectionAdminHandler.GetRequest)
	router.GET("/api/vehicle-inspection-admin/travel-details/:trn_request_uid", funcs.ApiKeyAuthenMiddleware(), vehicleInspectionAdminHandler.GetVehicleTripDetails)
	router.GET("/api/vehicle-inspection-admin/travel-detail/:trn_trip_detail_uid", funcs.ApiKeyAuthenMiddleware(), vehicleInspectionAdminHandler.GetVehicleTripDetail)
	router.POST("/api/vehicle-inspection-admin/create-travel-detail", funcs.ApiKeyAuthenMiddleware(), vehicleInspectionAdminHandler.CreateVehicleTripDetail)
	router.PUT("/api/vehicle-inspection-admin/update-travel-detail/:trn_trip_detail_uid", funcs.ApiKeyAuthenMiddleware(), vehicleInspectionAdminHandler.UpdateVehicleTripDetail)
	router.DELETE("/api/vehicle-inspection-admin/delete-travel-detail/:trn_trip_detail_uid", funcs.ApiKeyAuthenMiddleware(), vehicleInspectionAdminHandler.DeleteVehicleTripDetail)
	router.GET("/api/vehicle-inspection-admin/add-fuel-details/:trn_request_uid", funcs.ApiKeyAuthenMiddleware(), vehicleInspectionAdminHandler.GetVehicleAddFuelDetails)
	router.GET("/api/vehicle-inspection-admin/add-fuel-detail/:trn_add_fuel_uid", funcs.ApiKeyAuthenMiddleware(), vehicleInspectionAdminHandler.GetVehicleAddFuelDetail)
	router.POST("/api/vehicle-inspection-admin/create-add-fuel", funcs.ApiKeyAuthenMiddleware(), vehicleInspectionAdminHandler.CreateVehicleAddFuel)
	router.PUT("/api/vehicle-inspection-admin/update-add-fuel/:trn_add_fuel_uid", funcs.ApiKeyAuthenMiddleware(), vehicleInspectionAdminHandler.UpdateVehicleAddFuel)
	router.DELETE("/api/vehicle-inspection-admin/delete-add-fuel/:trn_add_fuel_uid", funcs.ApiKeyAuthenMiddleware(), vehicleInspectionAdminHandler.DeleteVehicleAddFuel)
	router.GET("/api/vehicle-inspection-admin/travel-card/:trn_request_uid", funcs.ApiKeyAuthenMiddleware(), vehicleInspectionAdminHandler.GetTravelCard)
	router.PUT("/api/vehicle-inspection-admin/update-returned-vehicle", funcs.ApiKeyAuthenMiddleware(), vehicleInspectionAdminHandler.UpdateReturnedVehicle)
	router.PUT("/api/vehicle-inspection-admin/update-returned-vehicle-images", funcs.ApiKeyAuthenMiddleware(), vehicleInspectionAdminHandler.UpdateReturnedVehicleImages)
	router.GET("/api/vehicle-inspection-admin/satisfaction-survey/:trn_request_uid", funcs.ApiKeyAuthenMiddleware(), vehicleInspectionAdminHandler.GetSatisfactionSurvey)
	router.PUT("/api/vehicle-inspection-admin/update-rejected", funcs.ApiKeyAuthenMiddleware(), vehicleInspectionAdminHandler.UpdateRejected)
	router.PUT("/api/vehicle-inspection-admin/update-accepted", funcs.ApiKeyAuthenMiddleware(), vehicleInspectionAdminHandler.UpdateAccepted)
	router.PUT("/api/vehicle-inspection-admin/update-inspect-vehicle-images", funcs.ApiKeyAuthenMiddleware(), vehicleInspectionAdminHandler.UpdateInspectVehicleImages)

	//VehicleManagementHandler
	vehicleManagementHandler := handlers.VehicleManagementHandler{Role: "admin-super,admin-region,admin-dept"}
	router.GET("/api/vehicle-management/search", funcs.ApiKeyAuthenMiddleware(), vehicleManagementHandler.SearchVehicles)
	router.PUT("/api/vehicle-management/update-vehicle-is-active", funcs.ApiKeyAuthenMiddleware(), vehicleManagementHandler.UpdateVehicleIsActive)
	router.GET("/api/vehicle-management/timeline", funcs.ApiKeyAuthenMiddleware(), vehicleManagementHandler.GetVehicleTimeLine)
	router.POST("/api/vehicle-management/report-trip-detail", funcs.ApiKeyAuthenMiddleware(), vehicleManagementHandler.ReportTripDetail)
	router.POST("/api/vehicle-management/report-add-fuel", funcs.ApiKeyAuthenMiddleware(), vehicleManagementHandler.ReportAddFuel)

	//DriverManagementHandler
	driverManagementHandler := handlers.DriverManagementHandler{Role: "admin-super,admin-region,admin-dept"}
	router.GET("/api/driver-management/search", funcs.ApiKeyAuthenMiddleware(), driverManagementHandler.SearchDrivers)
	router.POST("/api/driver-management/create-driver", funcs.ApiKeyAuthenMiddleware(), driverManagementHandler.CreateDriver)
	router.GET("/api/driver-management/driver/:mas_driver_uid", funcs.ApiKeyAuthenMiddleware(), driverManagementHandler.GetDriver)
	router.PUT("/api/driver-management/update-driver-detail", funcs.ApiKeyAuthenMiddleware(), driverManagementHandler.UpdateDriverDetail)
	router.PUT("/api/driver-management/update-driver-contract", funcs.ApiKeyAuthenMiddleware(), driverManagementHandler.UpdateDriverContract)
	router.PUT("/api/driver-management/update-driver-license", funcs.ApiKeyAuthenMiddleware(), driverManagementHandler.UpdateDriverLicense)
	router.PUT("/api/driver-management/update-driver-documents", funcs.ApiKeyAuthenMiddleware(), driverManagementHandler.UpdateDriverDocuments)
	router.PUT("/api/driver-management/update-driver-leave-status", funcs.ApiKeyAuthenMiddleware(), driverManagementHandler.UpdateDriverLeaveStatus)
	router.PUT("/api/driver-management/update-driver-is-active", funcs.ApiKeyAuthenMiddleware(), driverManagementHandler.UpdateDriverIsActive)
	router.DELETE("/api/driver-management/delete-driver", funcs.ApiKeyAuthenMiddleware(), driverManagementHandler.DeleteDriver)
	router.PUT("/api/driver-management/update-driver-layoff-status", funcs.ApiKeyAuthenMiddleware(), driverManagementHandler.UpdateDriverLayoffStatus)
	router.PUT("/api/driver-management/update-driver-resign-status", funcs.ApiKeyAuthenMiddleware(), driverManagementHandler.UpdateDriverResignStatus)
	router.GET("/api/driver-management/replacement-drivers", funcs.ApiKeyAuthenMiddleware(), driverManagementHandler.GetReplacementDrivers)
	router.GET("/api/driver-management/timeline", funcs.ApiKeyAuthenMiddleware(), driverManagementHandler.GetDriverTimeLine)
	router.POST("/api/driver-management/import-driver", funcs.ApiKeyAuthenMiddleware(), driverManagementHandler.ImportDriver)
	router.POST("/api/driver-management/work-report", funcs.ApiKeyAuthenMiddleware(), driverManagementHandler.GetDriverWorkReport)

	//DriverLicenseUserHandler
	driverLicenseUserHandler := handlers.DriverLicenseUserHandler{Role: "vehicle-user"}
	router.GET("/api/driver-license-user/card", funcs.ApiKeyAuthenMiddleware(), driverLicenseUserHandler.GetLicenseCard)
	router.POST("/api/driver-license-user/create-license-annual", funcs.ApiKeyAuthenMiddleware(), driverLicenseUserHandler.CreateDriverLicenseAnnual)
	router.GET("/api/driver-license-user/license-annual/:trn_request_annual_driver_uid", funcs.ApiKeyAuthenMiddleware(), driverLicenseUserHandler.GetDriverLicenseAnnual)
	router.PUT("/api/driver-license-user/update-license-annual-canceled", funcs.ApiKeyAuthenMiddleware(), driverLicenseUserHandler.UpdateDriverLicenseAnnualCanceled)
	router.PUT("/api/driver-license-user/resend-license-annual/:trn_request_annual_driver_uid", funcs.ApiKeyAuthenMiddleware(), driverLicenseUserHandler.ResendDriverLicenseAnnual)

	//DriverLicenseConfirmerHandler
	driverLicenseConfirmerHandler := handlers.DriverLicenseConfirmerHandler{Role: "level1-approval"}
	router.GET("/api/driver-license-confirmer/search-requests", funcs.ApiKeyAuthenMiddleware(), driverLicenseConfirmerHandler.SearchRequests)
	router.GET("/api/driver-license-confirmer/license-annual/:trn_request_annual_driver_uid", funcs.ApiKeyAuthenMiddleware(), driverLicenseConfirmerHandler.GetDriverLicenseAnnual)
	router.PUT("/api/driver-license-confirmer/update-license-annual-canceled", funcs.ApiKeyAuthenMiddleware(), driverLicenseConfirmerHandler.UpdateDriverLicenseAnnualCanceled)
	router.PUT("/api/driver-license-confirmer/update-license-annual-confirmed", funcs.ApiKeyAuthenMiddleware(), driverLicenseConfirmerHandler.UpdateDriverLicenseAnnualConfirmed)
	router.PUT("/api/driver-license-confirmer/update-license-annual-rejected", funcs.ApiKeyAuthenMiddleware(), driverLicenseConfirmerHandler.UpdateDriverLicenseAnnualRejected)
	router.PUT("/api/driver-license-confirmer/update-license-annual-approver", funcs.ApiKeyAuthenMiddleware(), driverLicenseConfirmerHandler.UpdateDriverLicenseAnnualApprover)

	//DriverLicenseApproverHandler
	driverLicenseApproverHandler := handlers.DriverLicenseApproverHandler{Role: "license-approval"}
	router.GET("/api/driver-license-approver/menu-requests", funcs.ApiKeyAuthenMiddleware(), driverLicenseApproverHandler.MenuRequests)
	router.GET("/api/driver-license-approver/search-requests", funcs.ApiKeyAuthenMiddleware(), driverLicenseApproverHandler.SearchRequests)
	router.GET("/api/driver-license-approver/license-annual/:trn_request_annual_driver_uid", funcs.ApiKeyAuthenMiddleware(), driverLicenseApproverHandler.GetDriverLicenseAnnual)
	router.PUT("/api/driver-license-approver/update-license-annual-canceled", funcs.ApiKeyAuthenMiddleware(), driverLicenseApproverHandler.UpdateDriverLicenseAnnualCanceled)
	router.PUT("/api/driver-license-approver/update-license-annual-approved", funcs.ApiKeyAuthenMiddleware(), driverLicenseApproverHandler.UpdateDriverLicenseAnnualApproved)
	router.PUT("/api/driver-license-approver/update-license-annual-rejected", funcs.ApiKeyAuthenMiddleware(), driverLicenseApproverHandler.UpdateDriverLicenseAnnualRejected)

	//CarpoolManagementHandler
	carpoolManagementHandler := handlers.CarpoolManagementHandler{Role: "admin-super,admin-region,admin-dept"}
	router.GET("/api/carpool-management/search", funcs.ApiKeyAuthenMiddleware(), carpoolManagementHandler.SearchCarpools)
	router.GET("/api/carpool-management/export", funcs.ApiKeyAuthenMiddleware(), carpoolManagementHandler.ExportCarpools)
	router.POST("/api/carpool-management/create", funcs.ApiKeyAuthenMiddleware(), carpoolManagementHandler.CreateCarpool)
	router.GET("/api/carpool-management/carpool/:mas_carpool_uid", funcs.ApiKeyAuthenMiddleware(), carpoolManagementHandler.GetCarpool)
	router.PUT("/api/carpool-management/update/:mas_carpool_uid", funcs.ApiKeyAuthenMiddleware(), carpoolManagementHandler.UpdateCarpool)
	router.DELETE("/api/carpool-management/delete", funcs.ApiKeyAuthenMiddleware(), carpoolManagementHandler.DeleteCarpool)
	router.GET("/api/carpool-management/mas-department", funcs.ApiKeyAuthenMiddleware(), carpoolManagementHandler.GetMasDepartment)
	router.GET("/api/carpool-management/mas-department/:carpool_type", funcs.ApiKeyAuthenMiddleware(), carpoolManagementHandler.GetMasDepartment)

	router.GET("/api/carpool-management/admin-mas-search", funcs.ApiKeyAuthenMiddleware(), carpoolManagementHandler.SearchMasAdminUser)
	router.GET("/api/carpool-management/admin-search/:mas_carpool_uid", funcs.ApiKeyAuthenMiddleware(), carpoolManagementHandler.SearchCarpoolAdmin)
	router.GET("/api/carpool-management/admin-detail/:mas_carpool_admin_uid", funcs.ApiKeyAuthenMiddleware(), carpoolManagementHandler.GetCarpoolAdmin)
	router.POST("/api/carpool-management/admin-create", funcs.ApiKeyAuthenMiddleware(), carpoolManagementHandler.CreateCarpoolAdmin)
	router.PUT("/api/carpool-management/admin-update/:mas_carpool_admin_uid", funcs.ApiKeyAuthenMiddleware(), carpoolManagementHandler.UpdateCarpoolAdmin)
	router.DELETE("/api/carpool-management/admin-delete/:mas_carpool_admin_uid", funcs.ApiKeyAuthenMiddleware(), carpoolManagementHandler.DeleteCarpoolAdmin)
	router.PUT("/api/carpool-management/set-active", funcs.ApiKeyAuthenMiddleware(), carpoolManagementHandler.SetActiveCarpool)
	router.PUT("/api/carpool-management/admin-update-main-admin/:mas_carpool_admin_uid", funcs.ApiKeyAuthenMiddleware(), carpoolManagementHandler.UpdateCarpoolMainAdmin)

	router.GET("/api/carpool-management/approver-mas-search", funcs.ApiKeyAuthenMiddleware(), carpoolManagementHandler.SearchMasApprovalUser)
	router.GET("/api/carpool-management/approver-search/:mas_carpool_uid", funcs.ApiKeyAuthenMiddleware(), carpoolManagementHandler.SearchCarpoolApprover)
	router.GET("/api/carpool-management/approver-detail/:mas_carpool_approver_uid", funcs.ApiKeyAuthenMiddleware(), carpoolManagementHandler.GetCarpoolApprover)
	router.POST("/api/carpool-management/approver-create", funcs.ApiKeyAuthenMiddleware(), carpoolManagementHandler.CreateCarpoolApprover)
	router.PUT("/api/carpool-management/approver-update/:mas_carpool_approver_uid", funcs.ApiKeyAuthenMiddleware(), carpoolManagementHandler.UpdateCarpoolApprover)
	router.DELETE("/api/carpool-management/approver-delete/:mas_carpool_approver_uid", funcs.ApiKeyAuthenMiddleware(), carpoolManagementHandler.DeleteCarpoolApprover)
	router.PUT("/api/carpool-management/approver-update-main-approver/:mas_carpool_approver_uid", funcs.ApiKeyAuthenMiddleware(), carpoolManagementHandler.UpdateCarpoolMainApprover)

	router.GET("/api/carpool-management/vehicle-search/:mas_carpool_uid", funcs.ApiKeyAuthenMiddleware(), carpoolManagementHandler.SearchCarpoolVehicle)
	router.POST("/api/carpool-management/vehicle-create", funcs.ApiKeyAuthenMiddleware(), carpoolManagementHandler.CreateCarpoolVehicle)
	router.PUT("/api/carpool-management/vehicle-set-active", funcs.ApiKeyAuthenMiddleware(), carpoolManagementHandler.SetActiveCarpoolVehicle)
	router.DELETE("/api/carpool-management/vehicle-delete/:mas_carpool_vehicle_uid", funcs.ApiKeyAuthenMiddleware(), carpoolManagementHandler.DeleteCarpoolVehicle)
	router.GET("/api/carpool-management/vehicle-mas-search", funcs.ApiKeyAuthenMiddleware(), carpoolManagementHandler.SearchMasVehicles)
	router.POST("/api/carpool-management/vehicle-mas-details", funcs.ApiKeyAuthenMiddleware(), carpoolManagementHandler.GetMasVehicleDetail)
	router.GET("/api/carpool-management/vehicle-timeline/:mas_carpool_uid", funcs.ApiKeyAuthenMiddleware(), carpoolManagementHandler.GetCarpoolVehicleTimeLine)

	router.GET("/api/carpool-management/driver-search/:mas_carpool_uid", funcs.ApiKeyAuthenMiddleware(), carpoolManagementHandler.SearchCarpoolDriver)
	router.POST("/api/carpool-management/driver-create", funcs.ApiKeyAuthenMiddleware(), carpoolManagementHandler.CreateCarpoolDriver)
	router.PUT("/api/carpool-management/driver-set-active", funcs.ApiKeyAuthenMiddleware(), carpoolManagementHandler.SetActiveCarpoolDriver)
	router.DELETE("/api/carpool-management/driver-delete/:mas_carpool_driver_uid", funcs.ApiKeyAuthenMiddleware(), carpoolManagementHandler.DeleteCarpoolDriver)
	router.GET("/api/carpool-management/driver-mas-search", funcs.ApiKeyAuthenMiddleware(), carpoolManagementHandler.SearchMasDrivers)
	router.POST("/api/carpool-management/driver-mas-details", funcs.ApiKeyAuthenMiddleware(), carpoolManagementHandler.GetMasDriverDetails)
	router.GET("/api/carpool-management/driver-timeline/:mas_carpool_uid", funcs.ApiKeyAuthenMiddleware(), carpoolManagementHandler.GetCarpoolDriverTimeLine)

	//MasHandler
	masHandler := handlers.MasHandler{}
	router.GET("/api/mas/user-vehicle-users", funcs.ApiKeyAuthenMiddleware(), masHandler.ListVehicleUser)
	router.GET("/api/mas/user-driver-users", funcs.ApiKeyAuthenMiddleware(), masHandler.ListDriverUser)
	router.GET("/api/mas/user-confirmer-users", funcs.ApiKeyAuthenMiddleware(), masHandler.ListConfirmerUser)
	router.GET("/api/mas/user-admin-approval-users", funcs.ApiKeyAuthenMiddleware(), masHandler.ListAdminApprovalUser)
	router.GET("/api/mas/user-final-approval-users", funcs.ApiKeyAuthenMiddleware(), masHandler.ListFinalApprovalUser)
	router.GET("/api/mas/user-received-key-users", funcs.ApiKeyAuthenMiddleware(), masHandler.ListReceivedKeyUser)
	router.GET("/api/mas/user/:emp_id", funcs.ApiKeyAuthenMiddleware(), masHandler.GetUserEmp)
	router.GET("/api/mas/satisfaction_survey_questions", funcs.ApiKeyAuthenMiddleware(), masHandler.ListVmsMasSatisfactionSurveyQuestions)
	router.GET("/api/mas/vehicle-departments", funcs.ApiKeyAuthenMiddleware(), masHandler.ListVehicleDepartment)
	router.GET("/api/mas/department-tree", funcs.ApiKeyAuthenMiddleware(), masHandler.GetDepartmentTree)
	router.GET("/api/mas/driver-departments", funcs.ApiKeyAuthenMiddleware(), masHandler.ListDriverDepartment)
	router.GET("/api/mas/user-confirmer-license-users", funcs.ApiKeyAuthenMiddleware(), masHandler.ListConfirmerLicenseUser)
	router.GET("/api/mas/user-approval-license-users", funcs.ApiKeyAuthenMiddleware(), masHandler.ListApprovalLicenseUser)
	router.GET("/api/mas/holidays", funcs.ApiKeyAuthenMiddleware(), masHandler.ListHoliday)

	//RefHandler
	refHandler := handlers.RefHandler{}
	router.GET("/api/ref/cost-type", funcs.ApiKeyAuthenMiddleware(), refHandler.ListCostType)
	router.GET("/api/ref/cost-type/:code", funcs.ApiKeyAuthenMiddleware(), refHandler.GetCostType)
	router.GET("/api/ref/request-status", funcs.ApiKeyAuthenMiddleware(), refHandler.ListRequestStatus)
	router.GET("/api/ref/fuel-type", funcs.ApiKeyAuthenMiddleware(), refHandler.ListFuelType)
	router.GET("/api/ref/oil-station-brand", funcs.ApiKeyAuthenMiddleware(), refHandler.ListOilStationBrand)
	router.GET("/api/ref/vehicle-img-side", funcs.ApiKeyAuthenMiddleware(), refHandler.ListVehicleImgSide)
	router.GET("/api/ref/payment-type-code", funcs.ApiKeyAuthenMiddleware(), refHandler.ListPaymentTypeCode)
	router.GET("/api/ref/driver-other-use", funcs.ApiKeyAuthenMiddleware(), refHandler.ListDriverOtherUse)
	router.GET("/api/ref/driver-license-type", funcs.ApiKeyAuthenMiddleware(), refHandler.ListDriverLicenseType)
	router.GET("/api/ref/driver-certificate-type", funcs.ApiKeyAuthenMiddleware(), refHandler.ListDriverCertificateType)
	router.GET("/api/ref/carpool-choose-driver", funcs.ApiKeyAuthenMiddleware(), refHandler.ListCarpoolChooseDriver)
	router.GET("/api/ref/carpool-choose-car", funcs.ApiKeyAuthenMiddleware(), refHandler.ListCarpoolChooseCar)
	router.GET("/api/ref/vehicle-key-type", funcs.ApiKeyAuthenMiddleware(), refHandler.ListVehicleKeyType)
	router.GET("/api/ref/leave-time-type", funcs.ApiKeyAuthenMiddleware(), refHandler.ListLeaveTimeType)
	router.GET("/api/ref/driver-status", funcs.ApiKeyAuthenMiddleware(), refHandler.ListDriverStatus)
	router.GET("/api/ref/cost-center", funcs.ApiKeyAuthenMiddleware(), refHandler.ListCostCenter)
	router.GET("/api/ref/vehicle-status", funcs.ApiKeyAuthenMiddleware(), refHandler.ListVehicleStatus)

	//NotificationHandler
	notificationHandler := handlers.NotificationHandler{}
	router.GET("/api/notification", funcs.ApiKeyAuthenMiddleware(), notificationHandler.GetNotification)
	router.PUT("/api/notification/read/:notification_uid", funcs.ApiKeyAuthenMiddleware(), notificationHandler.UpdateReadNotification)
	//LogHandler
	logHandler := handlers.LogHandler{}
	router.GET("/api/log/request/:trn_request_uid", funcs.ApiKeyAuthenMiddleware(), logHandler.GetLogRequest)

	//ServiceHandler
	serviceHandler := handlers.ServiceHandler{}
	router.GET("/api/service/request-booking/:request_no", serviceHandler.GetRequestBooking)
	router.GET("/api/service/vms-to-eems/:request_no", serviceHandler.GetVMSToEEMS)

	//UploadHandler
	uploadHandler := handlers.UploadHandler{}
	router.POST("/api/upload", funcs.ApiKeyMiddleware(), uploadHandler.UploadFile)
	router.GET("/api/upload/files/:bucket", uploadHandler.ListFiles)
	router.GET("/api/files/:bucket/:file", uploadHandler.ViewFile)

	// Swagger documentation
	router.GET("/api/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Start server
	port := strconv.Itoa(config.AppConfig.Port)
	log.Println("Server started at " + config.AppConfig.Host + ":" + port)

	router.Run(config.AppConfig.Host + ":" + port)
}
