package main

import (
	"github.com/gorilla/mux"
	"iceroad/codenight/config"
	"iceroad/codenight/db"
	"iceroad/codenight/github"
	"iceroad/codenight/session"
	"iceroad/codenight/user"
	"log"
	"net/http"
	"os"
)

func main() {
	log.SetOutput(os.Stdout)

	environment := config.GetEnv()

	session.SetupStore(environment)
	db.EstablishInitialConnection(environment)
	github.ConfigureClient(environment)

	router := mux.NewRouter()

	user.AddRoutes(router)

	log.Printf("Starting server on port: %s \n", environment.Port)
	log.Fatal(http.ListenAndServe(":"+environment.Port, router))
}
