package web

import (
	"encoding/json"
	"net/http"
)

const ErrorCodeNotFound = "not_found"
const ErrorCodeInvalidFormat = "invalid_format"
const ErrorCodeServerError = "server_error"
const ErrorCodeForbidden = "forbidden"

func JsonResponse(responseWriter http.ResponseWriter, bodyObj interface{}, statusCode int) {
	responseWriter.WriteHeader(statusCode)
	responseWriter.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(responseWriter).Encode(bodyObj)
}

func EncodeJson(obj interface{}) string {
	encodedObj, _ := json.Marshal(obj)
	return string(encodedObj)
}

type HttpError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}
