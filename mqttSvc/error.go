package mqttSvc

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
	logger "github.com/ecoprohcm/DMS_BackendServer/logs"
)

func HandleMqttErr(t mqtt.Token) error {
	if t == nil || t.Error() == nil {
		return nil
	}

	logger.LogWithoutFields(logger.MQTT, logger.ErrorLevel, t.Error())
	return t.Error()
}
