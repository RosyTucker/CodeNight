package user

import (
	"github.com/gorilla/mux"
	"github.com/rosytucker/codenight/config"
	"github.com/rosytucker/codenight/github"
	"github.com/rosytucker/codenight/session"
	"github.com/rosytucker/codenight/web"
	"log"
	"net/http"
)

var environment = config.GetEnv()

func AddRoutes(router *mux.Router) {
	router.HandleFunc("/user/current", getCurrentUserHandler).Methods(http.MethodGet)
	router.HandleFunc("/user/{userId:[A-Za-z0-9]+}", getUserHandler).Methods(http.MethodGet)
	router.HandleFunc("/user/{userId:[A-Za-z0-9]+}", putUserHandler).Methods(http.MethodPut)
	router.HandleFunc("/login", loginHandler).Methods(http.MethodGet)
	router.HandleFunc("/oauthCallback", oauthCallbackHandler).Methods(http.MethodGet)
}

func getUserHandler(res http.ResponseWriter, req *http.Request) {
	userId := mux.Vars(req)["userId"]

	log.Printf("Finding user with id: %+v \n", userId)

	validUserId := ValidateId(userId)

	if !validUserId {
		log.Printf("ERROR: Invalid user id format '%s'\n", userId)
		httpError := web.HttpError{Code: web.ErrorCodeInvalidFormat}
		web.JsonResponse(res, httpError, http.StatusNotFound)
		return
	}

	user, err := GetPublicById(userId)

	if err != nil {
		log.Printf("ERROR: Failed to find user %+v \n", err)
		httpError := web.HttpError{Code: web.ErrorCodeNotFound}
		web.JsonResponse(res, httpError, http.StatusNotFound)
		return
	}
	log.Printf("SUCCESS: Fetched user with id: %+v \n", user)
	web.JsonResponse(res, user, http.StatusOK)
}

func getCurrentUserHandler(res http.ResponseWriter, req *http.Request) {
	userId := session.Get(req, "userId")
	user, err := GetById(userId)

	log.Printf("Finding user with id: %+v \n", userId)

	if err != nil {
		log.Printf("ERROR: Failed to find current user %+v \n", err)
		httpError := web.HttpError{Code: web.ErrorCodeNotFound}
		web.JsonResponse(res, httpError, http.StatusNotFound)
		return
	}
	log.Printf("SUCCESS: Fetched current user with id: %+v \n", user)
	web.JsonResponse(res, user, http.StatusOK)
}

func putUserHandler(res http.ResponseWriter, req *http.Request) {
	userId := mux.Vars(req)["userId"]

	validUserId := ValidateId(userId)

	if validUserId {
		log.Printf("ERROR: Invalid user id format '%s'\n", userId)
		httpError := web.HttpError{Code: web.ErrorCodeInvalidFormat}
		web.JsonResponse(res, httpError, http.StatusNotFound)
		return
	}

	loggedInUserId := session.Get(req, "userId")

	if loggedInUserId != userId {
		log.Printf("ERROR: User %s attempted to update user %s \n", loggedInUserId, userId)
		httpError := web.HttpError{Code: web.ErrorCodeForbidden, Message: "you can only update yourself"}
		web.JsonResponse(res, httpError, http.StatusForbidden)
		return
	}

	log.Printf("PUTTING User with Id '%+v' \n", userId)

	user, err := FromJsonBody(req.Body)

	if err != nil {
		log.Printf("ERROR: Failed to read body as json user %+v \n", err)
		httpError := web.HttpError{Code: web.ErrorCodeInvalidFormat, Message: err.Error()}
		web.JsonResponse(res, httpError, http.StatusBadRequest)
		return
	}

	err = Replace(userId, user)

	if err != nil {
		log.Printf("ERROR: Failed to PUT user %+v \n", err)
		httpError := web.HttpError{Code: web.ErrorCodeServerError}
		web.JsonResponse(res, httpError, http.StatusInternalServerError)
		return
	}

	log.Printf("SUCCESS: PUT User with Id '%+v' \n", userId)
	res.Header().Set("Location", "/user/"+userId)
	res.WriteHeader(http.StatusNoContent)
}

func oauthCallbackHandler(res http.ResponseWriter, req *http.Request) {
	token, err := github.GetToken(req)

	if err != nil {
		log.Printf("ERROR: Failed to get github token '%s'\n", err)
		http.Redirect(res, req, environment.PostLoginRedirect, http.StatusTemporaryRedirect)
		return
	}

	githubUser, err := github.GetUser(token)

	if err != nil {
		log.Printf("ERROR: Failed to get github user", err)
		http.Redirect(res, req, environment.PostLoginRedirect, http.StatusTemporaryRedirect)
		return
	}

	isAdmin := *githubUser.Login == environment.MasterUser

	user := &User{
		Name:      githubUser.Name,
		Token:     web.EncodeJson(token),
		UserName:  *githubUser.Login,
		Email:     githubUser.Email,
		Blog:      githubUser.Blog,
		Location:  githubUser.Location,
		AvatarUrl: githubUser.AvatarURL,
		IsAdmin:   isAdmin}

	log.Printf("Creating User with username '%s' if they dont already exist \n", user.UserName)
	userId, err := CreateIfNotExists(user)

	if err != nil {
		log.Printf("ERROR: Failed to create user", err)
		http.Redirect(res, req, environment.PostLoginRedirect, http.StatusTemporaryRedirect)
		return
	}

	session.Set(res, req, "userId", userId)
	log.Printf("Updated User with Id '%+v' \n", userId)

	res.Header().Set("Location", "/user/"+userId)
	http.Redirect(res, req, environment.PostLoginRedirect, http.StatusTemporaryRedirect)
}

func loginHandler(res http.ResponseWriter, req *http.Request) {
	url := github.LoginRedirectUrl()
	http.Redirect(res, req, url, http.StatusTemporaryRedirect)
}
