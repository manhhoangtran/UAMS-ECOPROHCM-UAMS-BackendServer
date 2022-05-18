// Package models provides business entity models and related business logics for the app
package models

import (
	"time"
)

type GormModel struct {
	ID        uint      `gorm:"primaryKey"`
	CreatedAt time.Time `swaggerignore:"true"`
	UpdatedAt time.Time `swaggerignore:"true"`
}

type DeleteID struct {
	ID uint `json:"id"`
}

// Struct defines user's password types
type UserPass struct {
	RfidPass   string `gorm:"type:varchar(256)" json:"rfidPass"`
	KeypadPass string `gorm:"type:varchar(256)" json:"keypadPass"`
}

// Struct defines HTTP request payload for creating open doorlock scheduler for users
type UserSchedulerReq struct {
	Scheduler       `json:"scheduler" binding:"required"`
	GatewayID       string `json:"gatewayId" binding:"required"`
	DoorlockAddress string `json:"doorlockAddress" binding:"required"`
}

// Struct defines all services for our IoC
type ServiceOptions struct {
	StudentSvc           *StudentSvc
	CustomerSvc          *CustomerSvc
	EmployeeSvc          *EmployeeSvc
	GatewaySvc           *GatewaySvc
	DoorlockSvc          *DoorlockSvc
	AreaSvc              *AreaSvc
	LogSvc               *LogSvc
	GwNetworkSvc         *GwNetworkSvc
	SchedulerSvc         *SchedulerSvc
	SecretKeySvc         *SecretKeySvc
	DoorlockStatusLogSvc *DoorlockStatusLogSvc
}
