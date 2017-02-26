package main

import (
	"github.com/gorilla/mux"
	"github.com/rosytucker/codenight/config"
	"github.com/rosytucker/codenight/db"
	"github.com/rosytucker/codenight/github"
	"github.com/rosytucker/codenight/user"
	"github.com/rs/cors"
	"log"
	"net/http"
)

func main() {
	environment := config.GetEnv()
	db.EstablishInitialConnection(environment)
	github.ConfigureClient(environment)

	router := mux.NewRouter()
	user.AddRoutes(router)

	corsOptions := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000", "https://codenight-ldn-ui.herokuapp.com"},
		AllowedHeaders:   []string{"*"},
		AllowedMethods:   []string{"GET", "OPTIONS", "POST", "PUT"},
		AllowCredentials: true,
	})
	config.Log.DebugF("Starting server on port: %s \n", environment.Port)
	log.Fatal(http.ListenAndServe(":"+environment.Port, corsOptions.Handler(router)))
}
