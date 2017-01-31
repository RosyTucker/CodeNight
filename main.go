package main

import (
	"github.com/gorilla/mux"
	"github.com/rosytucker/codenight/config"
	"github.com/rosytucker/codenight/db"
	"github.com/rosytucker/codenight/github"
	"github.com/rosytucker/codenight/session"
	"github.com/rosytucker/codenight/user"
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
