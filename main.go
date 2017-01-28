package main

import (
	"github.com/gorilla/mux"
	"iceroad/codenight/login"
	"iceroad/codenight/user"
	_ "iceroad/codenight/util/db"
	"iceroad/codenight/util/env"
	_ "iceroad/codenight/util/session"
	"log"
	"net/http"
	"os"
)

func main() {
	log.SetOutput(os.Stdout)
	environment := env.Get()

	router := mux.NewRouter()

	user.AddRoutes(router)
	login.AddRoutes(router)

	log.Printf("Starting server on port: %s \n", environment.Port)
	log.Fatal(http.ListenAndServe(":"+environment.Port, router))
}
