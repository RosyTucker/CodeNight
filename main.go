package main

import (
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"iceroad/codenight/login"
	"iceroad/codenight/user"
	"iceroad/codenight/util"
	"log"
	"net/http"
)

const userPath = "/user"

func getUser(responseWriter http.ResponseWriter, request *http.Request, params httprouter.Params) {
	user := user.User{Id: 1, Name: "bob", Email: "bob@bob.com", Description: "Bob makes burgers"}
	util.JsonResponse(responseWriter, user, http.StatusOK)
}

func addUser(responseWriter http.ResponseWriter, request *http.Request, params httprouter.Params) {
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

	router := httprouter.New()

	login.AddRoutes(router)
	router.GET(userPath, getUser)
	router.POST(userPath, addUser)

	fmt.Printf("Starting server on port: %s", env.Port)
	log.Fatal(http.ListenAndServe(":"+env.Port, router))
}
