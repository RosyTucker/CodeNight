package env

import (
	"encoding/json"
	"net/http"
)

func JsonResponse(responseWriter http.ResponseWriter, bodyObj interface{}, statusCode int) {
	responseWriter.WriteHeader(statusCode)
	responseWriter.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(responseWriter).Encode(bodyObj)
}

type HttpError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}
