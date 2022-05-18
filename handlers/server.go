package handlers

import (
	logger "github.com/ecoprohcm/DMS_BackendServer/logs"
	"github.com/gin-gonic/gin"
)

func SetupRouter(
	hOpts *HandlerOptions,
) *gin.Engine {
	r := gin.Default()

	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(logger.GinLogger())
	r.Use(CORSMiddleware())
	v1R := r.Group("/v1")
	{
		// Gateway routes
		v1R.GET("/gateways", hOpts.GatewayHandler.FindAllGateway)
		v1R.GET("/gateway/:id", hOpts.GatewayHandler.FindGatewayByID)
		v1R.POST("/gateway", hOpts.GatewayHandler.CreateGateway)
		v1R.PATCH("/gateway", hOpts.GatewayHandler.UpdateGateway)
		v1R.DELETE("/gateway", hOpts.GatewayHandler.DeleteGateway)
		v1R.DELETE("/gateway/:id/doorlock", hOpts.GatewayHandler.DeleteGatewayDoorlock)
		v1R.POST("/block/cmd", hOpts.GatewayHandler.UpdateGatewayCmdByBlockID)

		// Area routes
		v1R.GET("/areas", hOpts.AreaHandler.FindAllArea)
		v1R.GET("/area/:id", hOpts.AreaHandler.FindAreaByID)
		v1R.POST("/area", hOpts.AreaHandler.CreateArea)
		v1R.PATCH("/area", hOpts.AreaHandler.UpdateArea)
		v1R.DELETE("/area", hOpts.AreaHandler.DeleteArea)

		// Doorlock routes
		v1R.GET("/doorlocks", hOpts.DoorlockHandler.FindAllDoorlock)
		v1R.GET("/doorlock/:id", hOpts.DoorlockHandler.FindDoorlockByID)
		v1R.GET("/doorlock/status/:id", hOpts.DoorlockHandler.GetDoorlockStatusByID)
		// v1R.GET("/doorlock/status/serial/:id", hOpts.DoorlockHandler.GetDoorlockStatusBySerialID)
		v1R.POST("/doorlock", hOpts.DoorlockHandler.CreateDoorlock)
		v1R.PATCH("/doorlock", hOpts.DoorlockHandler.UpdateDoorlock)
		v1R.PATCH("/doorlock/cmd", hOpts.DoorlockHandler.UpdateDoorlockCmd)
		v1R.PATCH("/doorlock/state/cmd", hOpts.DoorlockHandler.UpdateDoorlockStateCmd)
		v1R.DELETE("/doorlock", hOpts.DoorlockHandler.DeleteDoorlock)

		// Doorlock log route
		v1R.GET("/doorlockStatusLogs", hOpts.DoorlockStatusLogHandler.GetAllDoorlockStatusLogs)
		v1R.GET("/doorlockStatusLog/:doorId", hOpts.DoorlockStatusLogHandler.GetDoorlockStatusLogByDoorID)
		v1R.GET("/doorlockStatusLog/date/:fromTime/:toTime", hOpts.DoorlockStatusLogHandler.GetDoorlockStatusLogInTimeRange)
		v1R.DELETE("/doorlockStatusLog/:doorId", hOpts.DoorlockStatusLogHandler.DeleteDoorlockStatusLogByDoorID)
		v1R.DELETE("/doorlockStatusLog/date/:fromTime/:toTime", hOpts.DoorlockStatusLogHandler.DeleteDoorlockStatusLogInTimeRange)

		// Student routes
		v1R.GET("/students", hOpts.StudentHandler.FindAllStudent)
		v1R.GET("/student/:mssv", hOpts.StudentHandler.FindStudentByMSSV)
		v1R.POST("/student", hOpts.StudentHandler.CreateStudent)
		v1R.PATCH("/student", hOpts.StudentHandler.UpdateStudent)
		v1R.DELETE("/student", hOpts.StudentHandler.DeleteStudent)
		v1R.POST("/student/:mssv/scheduler", hOpts.StudentHandler.AppendStudentScheduler)

		// Employee routes
		v1R.GET("/employees", hOpts.EmployeeHandler.FindAllEmployee)
		v1R.GET("/employee/:msnv", hOpts.EmployeeHandler.FindEmployeeByMSNV)
		v1R.POST("/employee", hOpts.EmployeeHandler.CreateEmployee)
		v1R.PATCH("/employee", hOpts.EmployeeHandler.UpdateEmployee)
		v1R.DELETE("/employee", hOpts.EmployeeHandler.DeleteEmployee)
		v1R.POST("/employee/:msnv/scheduler", hOpts.EmployeeHandler.AppendEmployeeScheduler)

		// Customer routes
		v1R.GET("/customers", hOpts.CustomerHandler.FindAllCustomer)
		v1R.GET("/customer/:cccd", hOpts.CustomerHandler.FindCustomerByCCCD)
		v1R.POST("/customer", hOpts.CustomerHandler.CreateCustomer)
		v1R.PATCH("/customer", hOpts.CustomerHandler.UpdateCustomer)
		v1R.DELETE("/customer", hOpts.CustomerHandler.DeleteCustomer)
		v1R.POST("/customer/:cccd/scheduler", hOpts.CustomerHandler.AppendCustomerScheduler)

		// Scheduler routes
		v1R.GET("/schedulers", hOpts.SchedulerHandler.FindAllScheduler)
		v1R.GET("/scheduler/:id", hOpts.SchedulerHandler.FindSchedulerByID)
		v1R.POST("/scheduler", hOpts.SchedulerHandler.CreateScheduler)
		v1R.PATCH("/scheduler", hOpts.SchedulerHandler.UpdateScheduler)
		v1R.DELETE("/scheduler", hOpts.SchedulerHandler.DeleteScheduler)
		v1R.POST("/scheduler/excel", hOpts.SchedulerHandler.AppendSchedulerOnExcel)
		v1R.PATCH("/scheduler/excel", hOpts.SchedulerHandler.UpdateSchedulerOnExcel)
		// Gateway log routes
		v1R.GET("/gatewayLogs", hOpts.LogHandler.FindAllGatewayLog)
		v1R.GET("/gatewayLog/:id", hOpts.LogHandler.FindGatewayLogByID)
		v1R.GET("/gatewayLogs/period/:id/date/:from/:to", hOpts.LogHandler.FindGatewayLogsByTime)
		v1R.POST("/gatewayLogs/period", hOpts.LogHandler.UpdateGatewayLogCleanPeriod)
		v1R.GET("/gatewayLogs/:id/period", hOpts.LogHandler.FindGatewayLogsTypeByTime)

		// Secret key routes
		v1R.GET("/secretkeys", hOpts.SecretKeyHandler.FindSecretKey)
		v1R.POST("/secretkey", hOpts.SecretKeyHandler.CreateSecretKey)
		v1R.PATCH("/secretkey", hOpts.SecretKeyHandler.UpdateSecretKey)
	}
	return r
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Accept, Origin, Cache-Control, X-Requested-With, User-Agent, Accept-Language, Accept-Encoding")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
