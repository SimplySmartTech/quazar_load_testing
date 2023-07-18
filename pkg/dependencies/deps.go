package dependencies

import (
	"simply_smart_mqtt_load_testing/pkg/loadtest"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type Dependency struct {
	Lt         loadtest.LoadTesting
	MqttClient mqtt.Client
}
