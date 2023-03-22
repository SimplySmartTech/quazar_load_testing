package loadtest

import (
	"encoding/json"
	"log"
	"net/http"
	"simply_smart_mqtt_load_testing/pkg/api"
	"simply_smart_mqtt_load_testing/pkg/model"

	paho_mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/go-playground/validator/v10"
)

func InitializeLoadTest(lt LoadTesting, client paho_mqtt.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get request body from request
		var loadParameter model.InitializeLoadTestRequest
		err := json.NewDecoder(r.Body).Decode(&loadParameter)
		if err != nil {
			log.Fatal(err)
		}

		// Validate input parameters
		err = validator.New().Struct(loadParameter)
		if err != nil {
			log.Fatal(err)
		}

		// Generate payload and send : thingsies
		lt.GeneratePayload(r.Context(), loadParameter, client)
	}
}

func CountThingsData(service LoadTesting) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		data, err := service.CountTemplateThingKeysData(ctx)
		if err != nil {
			api.RespondWithJSON(w, http.StatusInternalServerError, api.Response{Error: err.Error()})
			return
		}
		api.RespondWithJSON(w, http.StatusOK, api.Response{Data: data, Message: "successfully fetch"})
	}
}

func RemoveThingsDataFromTemplate(service LoadTesting) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		data, err := service.RemoveThingsData(ctx)
		if err != nil {
			api.RespondWithJSON(w, http.StatusInternalServerError, api.Response{Error: err.Error()})
			return
		}
		api.RespondWithJSON(w, http.StatusOK, api.Response{Data: data, Message: "successfully fetch"})
	}
}
