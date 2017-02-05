package main

import (
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/rosytucker/codenight/config"
	"github.com/rosytucker/codenight/db"
	"github.com/rosytucker/codenight/github"
	"github.com/rosytucker/codenight/user"
	"log"
	"net/http"
	"os"
)

func main() {
	environment := config.GetEnv()

	db.EstablishInitialConnection(environment)
	github.ConfigureClient(environment)

	router := mux.NewRouter()
	user.AddRoutes(router)

	log.Printf("Starting server on port: %s \n", environment.Port)
	log.Fatal(http.ListenAndServe(":"+environment.Port, handlers.LoggingHandler(os.Stdout, router)))
}
