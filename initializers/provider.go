package initializers

import (
	"fmt"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/ecoprohcm/DMS_BackendServer/handlers"
	logger "github.com/ecoprohcm/DMS_BackendServer/logs"
	"github.com/ecoprohcm/DMS_BackendServer/models"
	"github.com/ecoprohcm/DMS_BackendServer/mqttSvc"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

type ContextContainer struct {
	Config         Config
	Db             *gorm.DB
	MqttClient     mqtt.Client
	HandlerOptions *handlers.HandlerOptions
}

func ProvideConfig(envFilePath string) (Config, error) {
	cfg := Config{}
	err := godotenv.Load(envFilePath) //use env.local for localhost
	if err != nil {
		fmt.Printf("Error loading .env file %s", err)
		return Config{}, err
	}
	err = envconfig.Process("", &cfg)
	if err != nil {
		return cfg, err
	}

	logger.InitLogger(cfg.SvLogPath)
	return cfg, nil
}

func ProvideGormDb(config Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf("sqlserver://%s:%s@%s:%s?database=%s", config.DbUser, config.DbPass, config.ServerHost, config.DbPort, config.DbName)
	db, err := gorm.Open(sqlserver.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println("failed to connect database")
		return nil, err
	}
	models.Migrate(db)
	return db, nil
}

func ProvideSvcOptions(db *gorm.DB) *models.ServiceOptions {
	return &models.ServiceOptions{
		GatewaySvc:           models.NewGatewaySvc(db),
		GwNetworkSvc:         models.NewGwNetworkSvc(db),
		AreaSvc:              models.NewAreaSvc(db),
		DoorlockSvc:          models.NewDoorlockSvc(db),
		LogSvc:               models.NewLogSvc(db),
		StudentSvc:           models.NewStudentSvc(db),
		EmployeeSvc:          models.NewEmployeeSvc(db),
		SchedulerSvc:         models.NewSchedulerSvc(db),
		CustomerSvc:          models.NewCustomerSvc(db),
		SecretKeySvc:         models.NewSecretKeySvc(db),
		DoorlockStatusLogSvc: models.NewDoorlockStatusLogSvc(db),
	}
}

func ProvideMqttClient(config Config, svcOptions *models.ServiceOptions) mqtt.Client {
	return mqttSvc.MqttClient(
		config.MqttClient,
		config.ServerHost,
		config.MqttPort,
		svcOptions,
	)
}

func ProvideHandlerOptions(svcOptions *models.ServiceOptions, mqttClient mqtt.Client) *handlers.HandlerOptions {
	deps := &handlers.HandlerDependencies{
		SvcOpts:    svcOptions,
		MqttClient: mqttClient,
	}

	return &handlers.HandlerOptions{
		AreaHandler:              handlers.NewAreaHandler(deps),
		CustomerHandler:          handlers.NewCustomerHandler(deps),
		DoorlockHandler:          handlers.NewDoorlockHandler(deps),
		EmployeeHandler:          handlers.NewEmployeeHandler(deps),
		GatewayHandler:           handlers.NewGatewayHandler(deps),
		LogHandler:               handlers.NewGatewayLogHandler(deps),
		StudentHandler:           handlers.NewStudentHandler(deps),
		SchedulerHandler:         handlers.NewSchedulerHandler(deps),
		SecretKeyHandler:         handlers.NewSecretKeyHandler(deps),
		DoorlockStatusLogHandler: handlers.NewDoorlockStatusLogHandler(deps),
	}
}

func ProvideAppInfrastructure(config Config, db *gorm.DB, mqttClient mqtt.Client, handlerOpts *handlers.HandlerOptions) *ContextContainer {
	return &ContextContainer{
		Config:         config,
		Db:             db,
		MqttClient:     mqttClient,
		HandlerOptions: handlerOpts,
	}
}
