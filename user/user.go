package user

import (
	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2/bson"
	"iceroad/codenight/util/db"
	"log"
	"net/http"
)

type User struct {
	Id          bson.ObjectId `bson:"_id,omitempty"`
	Name        *string       `json:"name", bson:"name"`
	UserName    *string       `json:"username", bson:"username"`
	Email       *string       `json:"email", bson:"email"`
	Description *string       `json:"description", bson:"description"`
	Blog        *string       `json:"blog", bson:"blog"`
	Location    *string       `json:"location", bson:"location"`
	AvatarUrl   *string       `json:"avatar", bson: "avatar"`
}

func getUserHandler(res http.ResponseWriter, req *http.Request) {
	// 	vars := mux.Vars(req)
	// 	userId := vars["userId"]
	// 	result := db.Users().Find(bson.M{"_id": userId})

	// 	log.Printf("RESULT: user: %+v \n", result)
}

func Upsert(user *User) error {
	session := db.Connect()
	defer session.Close()

	usersColl := session.DB("codenight").C("users")

	user.Id = bson.NewObjectId()
	err := usersColl.Insert(&user)

	if err != nil {
		log.Fatalf("Failed to upsert user %+v \n", err)
	}
	log.Printf("SUCCESS: user upserted: %+v \n", user)

	return err
}

func AddRoutes(router *mux.Router) {
	router.HandleFunc("/user/{userId:[0-9]+}", getUserHandler).Methods(http.MethodGet)
}
