package main

import (
	"cnuser"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

const port = ":4000"
const userPath = "/user"

func main() {
	fmt.Printf("Starting server on port: %s", port)
	http.HandleFunc(userPath, getUser)
	log.Fatal(http.ListenAndServe(port, nil))
}

func getUser(responseWriter http.ResponseWriter, request *http.Request) {
	user := cnuser.User{Id: 1, Name: "bob", Email: "bob@bob.com", Description: "Bob makes burgers"}
	jsonResponse(responseWriter, user)
}

func jsonResponse(responseWriter http.ResponseWriter, bodyObj interface{}) {
	responseWriter.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(responseWriter).Encode(bodyObj)
}
