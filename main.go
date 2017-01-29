package main

import (
	"github.com/gorilla/mux"
	_ "iceroad/codenight/db"
	"iceroad/codenight/env"
	_ "iceroad/codenight/session"
	"iceroad/codenight/user"
	"log"
	"net/http"
	"os"
)

func main() {
	log.SetOutput(os.Stdout)

	environment := env.Get()

	router := mux.NewRouter()

	user.AddRoutes(router)

	log.Printf("Starting server on port: %s \n", environment.Port)
	log.Fatal(http.ListenAndServe(":"+environment.Port, router))
}
