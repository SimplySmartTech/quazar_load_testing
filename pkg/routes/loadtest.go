package routes

import (
	"net/http"
	"simply_smart_mqtt_load_testing/pkg/dependencies"
	"simply_smart_mqtt_load_testing/pkg/loadtest"

	paho_mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/gorilla/mux"
)

func LoadTest(router *mux.Router, deps dependencies.Dependency, client paho_mqtt.Client) {
	router.HandleFunc("/count-things-data-in-template", loadtest.CountThingsData(deps.Lt)).Methods(http.MethodGet)
	router.HandleFunc("/init-load-test", loadtest.InitializeLoadTest(deps.Lt, client)).Methods(http.MethodPost)
	router.HandleFunc("/remove-all-things-reading", loadtest.RemoveThingsDataFromTemplate(deps.Lt)).Methods(http.MethodGet)
	router.HandleFunc("/create-thingkeys", loadtest.CreateThingKeys(deps.Lt)).Methods(http.MethodPost).Methods(http.MethodPost)
}
