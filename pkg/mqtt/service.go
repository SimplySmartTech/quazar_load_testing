package mqtt

import (
	"fmt"
	"log"
	"os"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type MqttClient interface {
	Publish(topic string, qos byte, retained bool, payload interface{}) error
}

var mc mqtt.Client

type mqttClient struct {
	mClient mqtt.Client
}

func InitMqttClient() {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s:%s", os.Getenv("MQTT_BROKER"), os.Getenv("MQTT_PORT")))
	opts.SetClientID("go_mqtt_client")
	opts.SetUsername(os.Getenv("MQTT_USER"))
	opts.SetPassword(os.Getenv("MQTT_PASS"))
	opts.SetAutoReconnect(true)
	opts.SetDefaultPublishHandler(messagePubHandler)
	opts.OnConnect = connectHandler
	opts.OnConnectionLost = connectLostHandler
	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Println("MQTT client NOT CONNECTED", token.Error().Error())
	}
	mc = client
}

func GetMqttClient() MqttClient {
	log.Println("MQTT client loaded successfully!")
	return &mqttClient{
		mClient: mc,
	}
}

var messagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("Topic: %s | %s\n", msg.Topic(), msg.Payload())
}

var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	fmt.Println("Connected")
}

var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	fmt.Printf("Connect lost: %+v", err)
}

func (m *mqttClient) Publish(topic string, qos byte, retained bool, payload interface{}) error {
	token := m.mClient.Publish(topic, qos, retained, payload)
	log.Println(token.Error())
	return nil
}
