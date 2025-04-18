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
	config.InitConfig()
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

	//BookingUserHandler
	bookingUserHandler := handlers.BookingUserHandler{}
	router.GET("/api/booking-user/menu-requests", funcs.ApiKeyAuthenMiddleware(), bookingUserHandler.MenuRequests)
	router.POST("/api/booking-user/create-request", funcs.ApiKeyAuthenMiddleware(), bookingUserHandler.CreateRequest)
	router.GET("/api/booking-user/request/:trn_request_uid", funcs.ApiKeyAuthenMiddleware(), bookingUserHandler.GetRequest)
	router.PUT("/api/booking-user/update-vehicle-user", funcs.ApiKeyAuthenMiddleware(), bookingUserHandler.UpdateVehicleUser)
	router.PUT("/api/booking-user/update-trip", funcs.ApiKeyAuthenMiddleware(), bookingUserHandler.UpdateTrip)
	router.PUT("/api/booking-user/update-pickup", funcs.ApiKeyAuthenMiddleware(), bookingUserHandler.UpdatePickup)
	router.PUT("/api/booking-user/update-document", funcs.ApiKeyAuthenMiddleware(), bookingUserHandler.UpdateDocument)
	router.PUT("/api/booking-user/update-cost", funcs.ApiKeyAuthenMiddleware(), bookingUserHandler.UpdateCost)
	router.PUT("/api/booking-user/update-vehicle-type", funcs.ApiKeyAuthenMiddleware(), bookingUserHandler.UpdateVehicleType)
	router.PUT("/api/booking-user/update-approver", funcs.ApiKeyAuthenMiddleware(), bookingUserHandler.UpdateApprover)
	router.GET("/api/booking-user/search-requests", funcs.ApiKeyAuthenMiddleware(), bookingUserHandler.SearchRequests)
	router.PUT("/api/booking-user/update-canceled", funcs.ApiKeyAuthenMiddleware(), bookingUserHandler.UpdateCanceled)
	router.PUT("/api/booking-user/update-sended-back", funcs.ApiKeyAuthenMiddleware(), bookingUserHandler.UpdateSendedBack)

	//BookingApproverHandler
	bookingApproverHandler := handlers.BookingApproverHandler{}
	router.GET("/api/booking-approver/menu-requests", funcs.ApiKeyAuthenMiddleware(), bookingApproverHandler.MenuRequests)
	router.GET("/api/booking-approver/mmenu-requests", funcs.ApiKeyAuthenMiddleware(), bookingApproverHandler.MenuRequests)
	router.GET("/api/booking-approver/search-requests", funcs.ApiKeyAuthenMiddleware(), bookingApproverHandler.SearchRequests)
	router.GET("/api/booking-approver/request/:trn_request_uid", funcs.ApiKeyAuthenMiddleware(), bookingApproverHandler.GetRequest)
	router.PUT("/api/booking-approver/update-sended-back", funcs.ApiKeyAuthenMiddleware(), bookingApproverHandler.UpdateSendedBack)
	router.PUT("/api/booking-approver/update-approved", funcs.ApiKeyAuthenMiddleware(), bookingApproverHandler.UpdateApproved)
	router.PUT("/api/booking-approver/update-canceled", funcs.ApiKeyAuthenMiddleware(), bookingApproverHandler.UpdateCanceled)

	//BookingAdminHandler
	bookinAdminHandler := handlers.BookingAdminHandler{}
	router.GET("/api/booking-admin/menu-requests", funcs.ApiKeyAuthenMiddleware(), bookinAdminHandler.MenuRequests)
	router.GET("/api/booking-admin/search-requests", funcs.ApiKeyAuthenMiddleware(), bookinAdminHandler.SearchRequests)
	router.GET("/api/booking-admin/request/:trn_request_uid", funcs.ApiKeyAuthenMiddleware(), bookinAdminHandler.GetRequest)
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

	//BookingAdminDeptHandler
	bookinAdminDeptHandler := handlers.BookingAdminDeptHandler{}
	router.GET("/api/booking-admin-dept/menu-requests", funcs.ApiKeyAuthenMiddleware(), bookinAdminDeptHandler.MenuRequests)
	router.GET("/api/booking-admin-dept/search-requests", funcs.ApiKeyAuthenMiddleware(), bookinAdminDeptHandler.SearchRequests)
	router.GET("/api/booking-admin-dept/request/:trn_request_uid", funcs.ApiKeyAuthenMiddleware(), bookinAdminDeptHandler.GetRequest)
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

	//BookingFinalHandler
	bookingFinalHandler := handlers.BookingFinalHandler{}
	router.GET("/api/booking-final/menu-requests", funcs.ApiKeyAuthenMiddleware(), bookingFinalHandler.MenuRequests)
	router.GET("/api/booking-final/search-requests", funcs.ApiKeyAuthenMiddleware(), bookingFinalHandler.SearchRequests)
	router.GET("/api/booking-final/request/:trn_request_uid", funcs.ApiKeyAuthenMiddleware(), bookingFinalHandler.GetRequest)
	router.PUT("/api/booking-final/update-sended-back", funcs.ApiKeyAuthenMiddleware(), bookingFinalHandler.UpdateSendedBack)
	router.PUT("/api/booking-final/update-approved", funcs.ApiKeyAuthenMiddleware(), bookingFinalHandler.UpdateApproved)
	router.PUT("/api/booking-final/update-canceled", funcs.ApiKeyAuthenMiddleware(), bookingFinalHandler.UpdateCanceled)

	//ReceivedKeyUserHandler
	receivedKeyUserHandler := handlers.ReceivedKeyUserHandler{}
	router.GET("/api/received-key-user/search-requests", funcs.ApiKeyAuthenMiddleware(), receivedKeyUserHandler.SearchRequests)
	router.GET("/api/received-key-user/request/:trn_request_uid", funcs.ApiKeyAuthenMiddleware(), receivedKeyUserHandler.GetRequest)
	router.PUT("/api/received-key-user/update-key-pickup-pea", funcs.ApiKeyAuthenMiddleware(), receivedKeyUserHandler.UpdateKeyPickupPEA)
	router.PUT("/api/received-key-user/update-key-pickup-outsider", funcs.ApiKeyAuthenMiddleware(), receivedKeyUserHandler.UpdateKeyPickupOutSider)
	router.PUT("/api/received-key-user/update-key-pickup-driver", funcs.ApiKeyAuthenMiddleware(), receivedKeyUserHandler.UpdateKeyPickupDriver)
	router.PUT("/api/received-key-user/update-canceled", funcs.ApiKeyAuthenMiddleware(), receivedKeyUserHandler.UpdateCanceled)

	//ReceivedKeyAdminHandler
	receivedKeyAdminHandler := handlers.ReceivedKeyAdminHandler{}
	router.GET("/api/received-key-admin/search-requests", funcs.ApiKeyAuthenMiddleware(), receivedKeyAdminHandler.SearchRequests)
	router.GET("/api/received-key-admin/request/:trn_request_uid", funcs.ApiKeyAuthenMiddleware(), receivedKeyAdminHandler.GetRequest)
	router.PUT("/api/received-key-admin/update-recieived-key", funcs.ApiKeyAuthenMiddleware(), receivedKeyAdminHandler.UpdateRecieivedKey)
	router.PUT("/api/received-key-admin/update-key-pickup-pea", funcs.ApiKeyAuthenMiddleware(), receivedKeyAdminHandler.UpdateKeyPickupPEA)
	router.PUT("/api/received-key-admin/update-key-pickup-outsider", funcs.ApiKeyAuthenMiddleware(), receivedKeyAdminHandler.UpdateKeyPickupOutSider)
	router.PUT("/api/received-key-admin/update-key-pickup-driver", funcs.ApiKeyAuthenMiddleware(), receivedKeyAdminHandler.UpdateKeyPickupDriver)
	router.PUT("/api/received-key-admin/update-canceled", funcs.ApiKeyAuthenMiddleware(), receivedKeyAdminHandler.UpdateCanceled)
	router.PUT("/api/received-key-admin/update-recieived-key-detail", funcs.ApiKeyAuthenMiddleware(), receivedKeyAdminHandler.UpdateRecieivedKeyDetail)

	//ReceivedKeyDriverHandler
	receivedKeyDriverHandler := handlers.ReceivedKeyDriverHandler{}
	router.GET("/api/received-key-driver/search-requests", funcs.ApiKeyAuthenMiddleware(), receivedKeyDriverHandler.SearchRequests)
	router.GET("/api/received-key-driver/request/:trn_request_uid", funcs.ApiKeyAuthenMiddleware(), receivedKeyDriverHandler.GetRequest)
	router.PUT("/api/received-key-driver/update-recieived-key-detail", funcs.ApiKeyAuthenMiddleware(), receivedKeyDriverHandler.UpdateRecieivedKeyDetail)

	//ReceivedVehicleUserHandler
	receivedVehicleUserHandler := handlers.ReceivedVehicleUserHandler{}
	router.GET("/api/received-vehicle-user/search-requests", funcs.ApiKeyAuthenMiddleware(), receivedVehicleUserHandler.SearchRequests)
	router.GET("/api/received-vehicle-user/request/:trn_request_uid", funcs.ApiKeyAuthenMiddleware(), receivedVehicleUserHandler.GetRequest)
	router.PUT("/api/received-vehicle-user/received-vehicle", funcs.ApiKeyAuthenMiddleware(), receivedVehicleUserHandler.ReceivedVehicle)
	router.GET("/api/received-vehicle-user/travel-card/:trn_request_uid", funcs.ApiKeyAuthenMiddleware(), receivedVehicleUserHandler.GetTravelCard)

	//ReceivedVehicleAdminHandler
	receivedVehicleAdminHandler := handlers.ReceivedVehicleAdminHandler{}
	router.GET("/api/received-vehicle-admin/search-requests", funcs.ApiKeyAuthenMiddleware(), receivedVehicleAdminHandler.SearchRequests)
	router.GET("/api/received-vehicle-admin/request/:trn_request_uid", funcs.ApiKeyAuthenMiddleware(), receivedVehicleAdminHandler.GetRequest)
	router.PUT("/api/received-vehicle-admin/received-vehicle", funcs.ApiKeyAuthenMiddleware(), receivedVehicleAdminHandler.ReceivedVehicle)
	router.GET("/api/received-vehicle-admin/travel-card/:trn_request_uid", funcs.ApiKeyAuthenMiddleware(), receivedVehicleAdminHandler.GetTravelCard)

	//ReceivedVehicleDriverHandler
	receivedVehicleDriverHandler := handlers.ReceivedVehicleDriverHandler{}
	router.GET("/api/received-vehicle-driver/search-requests", funcs.ApiKeyAuthenMiddleware(), receivedVehicleDriverHandler.SearchRequests)
	router.GET("/api/received-vehicle-driver/request/:trn_request_uid", funcs.ApiKeyAuthenMiddleware(), receivedVehicleDriverHandler.GetRequest)
	router.PUT("/api/received-vehicle-driver/received-vehicle", funcs.ApiKeyAuthenMiddleware(), receivedVehicleDriverHandler.ReceivedVehicle)
	router.GET("/api/received-vehicle-driver/travel-card/:trn_request_uid", funcs.ApiKeyAuthenMiddleware(), receivedVehicleDriverHandler.GetTravelCard)

	//VehicleInUseUserHandler
	vehicleInUseUserHandler := handlers.VehicleInUseUserHandler{}
	router.GET("/api/vehicle-in-use-user/search-requests", funcs.ApiKeyAuthenMiddleware(), vehicleInUseUserHandler.SearchRequests)
	router.GET("/api/vehicle-in-use-user/request/:trn_request_uid", funcs.ApiKeyAuthenMiddleware(), vehicleInUseUserHandler.GetRequest)
	router.GET("/api/vehicle-in-use-user/travel-details/:trn_request_uid", funcs.ApiKeyAuthenMiddleware(), vehicleInUseUserHandler.GetVehicleTripDetails)
	router.GET("/api/vehicle-in-use-user/travel-detail/:trn_trip_detail_uid", funcs.ApiKeyAuthenMiddleware(), vehicleInUseUserHandler.GetVehicleTripDetail)
	router.POST("/api/vehicle-in-use-user/create-travel-detail", funcs.ApiKeyAuthenMiddleware(), vehicleInUseUserHandler.CreateVehicleTripDetail)
	router.PUT("/api/vehicle-in-use-user/update-travel-detail/:trn_trip_detail_uid", funcs.ApiKeyAuthenMiddleware(), vehicleInUseUserHandler.UpdateVehicleTripDetail)
	router.DELETE("/api/vehicle-in-use-user/delete-travel-detail/:trn_trip_detail_uid", funcs.ApiKeyAuthenMiddleware(), vehicleInUseUserHandler.DeleteVehicleTripDetail)
	router.GET("/api/vehicle-in-use-user/add-fuel-details/:trn_request_uid", funcs.ApiKeyAuthenMiddleware(), vehicleInUseUserHandler.GetVehicleAddFuelDetails)
	router.GET("/api/vehicle-in-use/add-fuel-detail/:trn_add_fuel_uid", funcs.ApiKeyAuthenMiddleware(), vehicleInUseUserHandler.GetVehicleAddFuelDetail)
	router.POST("/api/vehicle-in-use-user/create-add-fuel", funcs.ApiKeyAuthenMiddleware(), vehicleInUseUserHandler.CreateVehicleAddFuel)
	router.PUT("/api/vehicle-in-use-user/update-add-fuel/:trn_add_fuel_uid", funcs.ApiKeyAuthenMiddleware(), vehicleInUseUserHandler.UpdateVehicleAddFuel)
	router.DELETE("/api/vehicle-in-use-user/delete-add-fuel/:trn_add_fuel_uid", funcs.ApiKeyAuthenMiddleware(), vehicleInUseUserHandler.DeleteVehicleAddFuel)
	router.GET("/api/vehicle-in-use-user/travel-card/:trn_request_uid", funcs.ApiKeyAuthenMiddleware(), vehicleInUseUserHandler.GetTravelCard)
	router.PUT("/api/vehicle-in-use-user/returned-vehicle", funcs.ApiKeyAuthenMiddleware(), vehicleInUseUserHandler.ReturnedVehicle)

	//VehicleInUseAdminHandler
	vehicleInUseAdminHandler := handlers.VehicleInUseAdminHandler{}
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
	vehicleInUseDriverHandler := handlers.VehicleInUseDriverHandler{}
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
	vehicleInspectionAdminHandler := handlers.VehicleInspectionAdminHandler{}
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
	router.PUT("/api/vehicle-inspection-admin/update-sended-back", funcs.ApiKeyAuthenMiddleware(), vehicleInspectionAdminHandler.UpdateSendedBack)
	router.PUT("/api/vehicle-inspection-admin/update-accepted", funcs.ApiKeyAuthenMiddleware(), vehicleInspectionAdminHandler.UpdateAccepted)

	//VehicleHandler
	vehicleHandler := handlers.VehicleHandler{}
	router.GET("/api/vehicle/search", funcs.ApiKeyAuthenMiddleware(), vehicleHandler.SearchVehicles)
	router.GET("/api/vehicle/types", funcs.ApiKeyAuthenMiddleware(), vehicleHandler.GetTypes)
	router.GET("/api/vehicle/departments", funcs.ApiKeyAuthenMiddleware(), vehicleHandler.GetDepartments)
	router.GET("/api/vehicle/:mas_vehicle_uid", funcs.ApiKeyAuthenMiddleware(), vehicleHandler.GetVehicle)
	router.GET("/api/vehicle-info/:mas_vehicle_uid", funcs.ApiKeyAuthenMiddleware(), vehicleHandler.GetVehicleInfo)

	//DriverHandler
	driverHandler := handlers.DriverHandler{}
	router.GET("/api/driver/search", funcs.ApiKeyAuthenMiddleware(), driverHandler.GetDrivers)
	router.GET("/api/driver/:mas_driver_uid", funcs.ApiKeyAuthenMiddleware(), driverHandler.GetDriver)
	router.GET("/api/driver/search-other-dept", funcs.ApiKeyAuthenMiddleware(), driverHandler.GetDriversOtherDept)
	router.GET("/api/driver/work-type", funcs.ApiKeyAuthenMiddleware(), driverHandler.GetWorkType)

	//MasHandler
	masHandler := handlers.MasHandler{}
	router.GET("/api/mas/user-vehicle-users", funcs.ApiKeyAuthenMiddleware(), masHandler.ListVehicleUser)
	router.GET("/api/mas/user-driver-users", funcs.ApiKeyAuthenMiddleware(), masHandler.ListDriverUser)
	router.GET("/api/mas/user-approval-users", funcs.ApiKeyAuthenMiddleware(), masHandler.ListApprovalUser)
	router.GET("/api/mas/user-admin-approval-users", funcs.ApiKeyAuthenMiddleware(), masHandler.ListAdminApprovalUser)
	router.GET("/api/mas/user-final-approval-users", funcs.ApiKeyAuthenMiddleware(), masHandler.ListFinalApprovalUser)
	router.GET("/api/mas/user/:emp_id", funcs.ApiKeyAuthenMiddleware(), masHandler.GetUserEmp)
	router.GET("/api/mas/satisfaction_survey_questions", funcs.ApiKeyAuthenMiddleware(), masHandler.ListVmsMasSatisfactionSurveyQuestions)
	router.GET("/api/mas/vehicle-departments", funcs.ApiKeyAuthenMiddleware(), masHandler.ListVehicleDepartment)

	//RefHandler
	refHandler := handlers.RefHandler{}
	router.GET("/api/ref/cost-type", funcs.ApiKeyAuthenMiddleware(), refHandler.ListCostType)
	router.GET("/api/ref/cost-type/:code", funcs.ApiKeyAuthenMiddleware(), refHandler.GetCostType)
	router.GET("/api/ref/request-status", funcs.ApiKeyAuthenMiddleware(), refHandler.ListRequestStatus)
	router.GET("/api/ref/fuel-type", funcs.ApiKeyAuthenMiddleware(), refHandler.ListFuelType)
	router.GET("/api/ref/oil-station-brand", funcs.ApiKeyAuthenMiddleware(), refHandler.ListOilStationBrand)
	router.GET("/api/ref/vehicle-img-side", funcs.ApiKeyAuthenMiddleware(), refHandler.ListVehicleImgSide)
	router.GET("/api/ref/payment-type-code", funcs.ApiKeyAuthenMiddleware(), refHandler.ListPaymentTypeCode)

	logHandler := handlers.LogHandler{}
	router.GET("/api/log/request/:trn_request_uid", funcs.ApiKeyAuthenMiddleware(), logHandler.GetLogRequest)

	//UploadHandler
	uploadHandler := handlers.UploadHandler{}
	router.POST("/api/upload", funcs.ApiKeyMiddleware(), uploadHandler.UploadFile)

	// Swagger documentation
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Start server
	port := strconv.Itoa(config.AppConfig.Port)
	log.Println("Server started at " + config.AppConfig.Host + ":" + port)

	router.Run(config.AppConfig.Host + ":" + port)
}
