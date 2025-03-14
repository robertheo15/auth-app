package http

import (
	"encoding/json"
	"log"
	"net/http"
)

type DataTemplate struct {
	Message string      `json:"message"`
	Status  string      `json:"status"`
	Data    interface{} `json:"data"`
}

type ErrorMsg struct {
	Error string `json:"error"`
}

func SetResponse(w http.ResponseWriter, httpStatus int, data interface{}, message string, success bool) {
	response := DataTemplate{
		Message: message,
		Status:  http.StatusText(httpStatus),
		Data:    data,
	}
	encoded, err := json.Marshal(response)
	if err != nil {
		log.Printf("[ERROR][Util][WriteSuccess] marshalling response, %+v\n", err)
	}

	w.WriteHeader(httpStatus)
	w.Write(encoded)
}
func SetError(w http.ResponseWriter, httpStatus int, err error) {
	response := ErrorMsg{
		Error: "Internal Server Error",
	}

	encoded, err := json.Marshal(response)
	if err != nil {
		log.Printf("[ERROR][Util][WriteError] marshalling response, %+v\n", err)
	}

	w.WriteHeader(httpStatus)
	w.Write(encoded)
}
func ReturnInternalServerError(w http.ResponseWriter) {
	response := ErrorMsg{
		Error: "Internal Server Error",
	}

	encoded, err := json.Marshal(response)
	if err != nil {
		log.Printf("[ERROR][Util][WriteError] marshalling response, %+v\n", err)
	}

	w.WriteHeader(http.StatusInternalServerError)
	w.Write(encoded)
}
