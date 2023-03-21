package main

import (
	"fmt"
	"log"
	"net/http"
	"simply_smart_mqtt_load_testing/pkg/config"
	"simply_smart_mqtt_load_testing/pkg/database"
	"simply_smart_mqtt_load_testing/pkg/dependencies"
	"simply_smart_mqtt_load_testing/pkg/loadtest"
	"simply_smart_mqtt_load_testing/pkg/mqtt"
	"simply_smart_mqtt_load_testing/pkg/routes"

	paho_mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/gorilla/mux"
)

var client paho_mqtt.Client

func init() {
	config.LoadEnvironment()
	database.InitDatabase()

	var broker = "127.0.0.1"
	var port = 1883

	opts := paho_mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s:%d", broker, port))
	opts.SetClientID("go_mqtt_client")
	opts.SetUsername("78d9fa4teebt5add59ctb86e1a286477cb147392")
	opts.SetPassword("78d9fa4teebt5add59ctb86e1a286477cb147391")
	opts.KeepAlive = 100
	opts.SetDefaultPublishHandler(messagePubHandler)
	opts.OnConnect = connectHandler
	opts.OnConnectionLost = connectLostHandler
	client = paho_mqtt.NewClient(opts)
	token := client.Connect()
	if token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	// mqtt.InitMqttClient()
}

func main() {
	r := mux.NewRouter()

	// Initialize dependencies
	db := database.NewDatabase()
	mClient := mqtt.GetMqttClient()
	lt := loadtest.GetLoadTesting(db, mClient)
	deps := dependencies.Dependency{
		Lt: lt,
	}

	// Initialize routes
	routes.LoadTest(r, deps, client)

	// Start server
	if err := http.ListenAndServe(":8088", r); err != nil {
		log.Fatal(err)
	}
}

var messagePubHandler paho_mqtt.MessageHandler = func(client paho_mqtt.Client, msg paho_mqtt.Message) {
	fmt.Printf("Topic: %s | %s\n", msg.Topic(), msg.Payload())
}

var connectHandler paho_mqtt.OnConnectHandler = func(client paho_mqtt.Client) {
	fmt.Println("Connected")
}

var connectLostHandler paho_mqtt.ConnectionLostHandler = func(client paho_mqtt.Client, err error) {
	fmt.Printf("Connect lost: %+v", err)
}
