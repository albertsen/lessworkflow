package server

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type errorMessage struct {
	Message string `json:"message"`
}

func SendOK(w http.ResponseWriter, data interface{}) {
	SendResponse(w, http.StatusOK, data)
}

func SendError(w http.ResponseWriter, status int, message interface{}) {
	messageString := fmt.Sprintf("%s", message)
	SendResponse(w, status, errorMessage{Message: messageString})
}

func SendResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(statusCode)
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}

func SendResponseOrError(w http.ResponseWriter, expectedStatusCode int, actualStatusCode int, data interface{}, err error) {
	if actualStatusCode != expectedStatusCode {
		SendError(w, actualStatusCode, err)
	} else {
		SendResponse(w, expectedStatusCode, data)
	}
}
