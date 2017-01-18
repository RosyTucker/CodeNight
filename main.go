package main

import (
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"iceroad/codenight/user"
	"log"
	"net/http"
	"os"
)

const userPath = "/user"

func main() {
	env := getEnv()
	fmt.Printf("Starting server on port: %s", env.Port)

	router := httprouter.New()
	router.GET(userPath, getUser)
	router.POST(userPath, addUser)

	log.Fatal(http.ListenAndServe(":"+env.Port, router))
}

func getUser(responseWriter http.ResponseWriter, request *http.Request, params httprouter.Params) {
	user := user.User{Id: 1, Name: "bob", Email: "bob@bob.com", Description: "Bob makes burgers"}
	jsonResponse(responseWriter, user, http.StatusOK)
}

func addUser(responseWriter http.ResponseWriter, request *http.Request, params httprouter.Params) {
	userToAdd := user.User{}
	err := json.NewDecoder(request.Body).Decode(&userToAdd)

	if err != nil {
		responseBody := HttpError{Code: "invalid_format", Message: err.Error()}
		statusCode := http.StatusUnprocessableEntity
		jsonResponse(responseWriter, responseBody, statusCode)
		return
	}

	jsonResponse(responseWriter, userToAdd, http.StatusOK)
}

func jsonResponse(responseWriter http.ResponseWriter, bodyObj interface{}, statusCode int) {
	responseWriter.WriteHeader(statusCode)
	responseWriter.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(responseWriter).Encode(bodyObj)
}

type env struct {
	Port string
}

func getEnv() env {
	const defaultPort = "4000"
	return env{Port: defaultWhenEmpty(os.Getenv("PORT"), defaultPort)}
}

func defaultWhenEmpty(value string, defaultValue string) string {
	if len(value) == 0 {
		return defaultValue
	}
	return value
}

type HttpError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}
