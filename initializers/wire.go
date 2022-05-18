//go:build wireinject
// +build wireinject

package initializers

import (
	"github.com/google/wire"
)

var ApplicationSet = wire.NewSet(
	ProvideConfig,
	ProvideGormDb,
	ProvideSvcOptions,
	ProvideMqttClient,
	ProvideHandlerOptions,
	ProvideAppInfrastructure,
)

func InitApplication(envFilePath string) (*ContextContainer, func(), error) {
	wire.Build(
		ApplicationSet,
	)
	return nil, nil, nil
}
