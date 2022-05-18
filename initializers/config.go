package initializers

type Config struct {
	ServerHost string `envconfig:"SERVER_HOST"`
	DbHost     string `envconfig:"DB_HOST"`
	DbPort     string `envconfig:"DB_PORT"`
	DbUser     string `envconfig:"DB_USER"`
	DbPass     string `envconfig:"DB_PASS"`
	DbName     string `envconfig:"DB_NAME"`
	MqttHost   string `envconfig:"MQTT_HOST"`
	MqttPort   string `envconfig:"MQTT_PORT"`
	MqttClient string `envconfig:"MQTT_CLIENT"`
	SvLogPath  string `envconfig:"SV_LOG_FILE"`
}
