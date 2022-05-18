package handlers

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/ecoprohcm/DMS_BackendServer/models"
)

type HandlerOptions struct {
	AreaHandler              *AreaHandler
	CustomerHandler          *CustomerHandler
	DoorlockHandler          *DoorlockHandler
	EmployeeHandler          *EmployeeHandler
	GatewayHandler           *GatewayHandler
	LogHandler               *GatewayLogHandler
	StudentHandler           *StudentHandler
	SchedulerHandler         *SchedulerHandler
	SecretKeyHandler         *SecretKeyHandler
	DoorlockStatusLogHandler *DoorlockStatusLogHandler
}

type HandlerDependencies struct {
	SvcOpts    *models.ServiceOptions
	MqttClient mqtt.Client
}
