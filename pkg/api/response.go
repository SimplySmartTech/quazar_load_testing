package api

import (
	"encoding/json"
	"net/http"
)

// swagger:model Response
type Response struct {
	//in:body
	//example: "success message"
	Message string `json:"message,omitempty"`
	//in:body
	Data interface{} `json:"data,omitempty"`
	//in:body
	//example: "error message"
	Error string `json:"error,omitempty"`
}

func RespondWithJSON(rw http.ResponseWriter, status int, response Response) {
	respBytes, err := json.Marshal(response)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	rw.Header().Add("Content-Type", "application/json")
	rw.WriteHeader(status)
	rw.Write(respBytes)
}
