package user

import (
	"github.com/gorilla/mux"
	"net/http"
)

type User struct {
	Id          *int    `json:"id"`
	Name        *string `json:"name"`
	UserName    *string `json:"username"`
	Email       *string `json:"email"`
	Description *string `json:"description"`
	Blog        *string `json:"blog"`
	Location    *string `json:"location"`
	AvatarUrl   *string `json:"avatar"`
}

func getUserHandler(res http.ResponseWriter, req *http.Request) {
	// vars := mux.Vars(req)
	// userId := vars["userId"]
}

func Upsert(user User, token string) error {
	return nil
}

func AddRoutes(router *mux.Router) {
	router.HandleFunc("/user/{userId:[0-9]+}", getUserHandler).Methods(http.MethodGet)
}
