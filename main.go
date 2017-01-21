package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"iceroad/codenight/login"
	"iceroad/codenight/user"
	"iceroad/codenight/util"
	"log"
	"net/http"
)

const userPath = "/user"

func getUser(responseWriter http.ResponseWriter, request *http.Request) {
	user := user.User{Id: 1, Name: "bob", Email: "bob@bob.com", Description: "Bob makes burgers"}
	util.JsonResponse(responseWriter, user, http.StatusOK)
}

func addUser(responseWriter http.ResponseWriter, request *http.Request) {
	userToAdd := user.User{}
	err := json.NewDecoder(request.Body).Decode(&userToAdd)

	if err != nil {
		responseBody := util.HttpError{Code: "invalid_format", Message: err.Error()}
		statusCode := http.StatusUnprocessableEntity
		util.JsonResponse(responseWriter, responseBody, statusCode)
		return
	}

	util.JsonResponse(responseWriter, userToAdd, http.StatusOK)
}

func main() {
	env := util.GetEnv()

	router := mux.NewRouter()

	login.AddRoutes(router)
	router.HandleFunc(userPath, getUser).Methods(http.MethodGet)
	router.HandleFunc(userPath, addUser).Methods(http.MethodPost)

	fmt.Printf("Starting server on port: %s \n", env.Port)
	log.Fatal(http.ListenAndServe(":"+env.Port, router))
}
