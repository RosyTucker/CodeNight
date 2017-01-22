package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"iceroad/codenight/login"
	"iceroad/codenight/user"
	"iceroad/codenight/util"
	"log"
	"net/http"
)

func main() {
	env := util.GetEnv()

	router := mux.NewRouter()

	user.AddRoutes(router)
	login.AddRoutes(router)

	fmt.Printf("Starting server on port: %s \n", env.Port)
	log.Fatal(http.ListenAndServe(":"+env.Port, router))
}
