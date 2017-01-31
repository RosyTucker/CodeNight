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

	log.Printf("Getting User with Id '%+v' \n", userId)

	isHandled := handleInvalidUserIdForRequest(userId, false, res, req)
	if isHandled {
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
	loggedInUserId, err := session.Get(req, "userId")

	isHandled := handleInvalidUserIdForRequest(loggedInUserId, false, res, req)
	if isHandled {
		return
	}

	user, err := GetById(loggedInUserId)

	log.Printf("Finding current user with id: %+v \n", loggedInUserId)

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

	log.Printf("Putting User with Id '%+v' \n", userId)

	isHandled := handleInvalidUserIdForRequest(userId, true, res, req)
	if isHandled {
		return
	}

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
		log.Printf("ERROR: Failed to create user %+v \n", err)
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

func handleInvalidUserIdForRequest(requestedUserId string, requiresSelf bool, res http.ResponseWriter, req *http.Request) bool {
	log.Printf("Validating request for user with id: %+v \n", requestedUserId)

	loggedInUserId, err := session.Get(req, "userId")

	if err != nil {
		log.Printf("ERROR: User not logged in and attempted to access user %s \n", requestedUserId)
		httpError := web.HttpError{
			Code:    web.ErrorCodeUnauthorized,
			Message: "you must be logged in try and view user information"}
		web.JsonResponse(res, httpError, http.StatusUnauthorized)
		return true
	}

	if requiresSelf && loggedInUserId != requestedUserId {
		log.Printf("ERROR: User %s attempted to access user %s \n", loggedInUserId, requestedUserId)
		httpError := web.HttpError{
			Code:    web.ErrorCodeForbidden,
			Message: "you can only update yourself"}
		web.JsonResponse(res, httpError, http.StatusForbidden)
		return true
	}

	validUserId := ValidateId(requestedUserId)

	if !validUserId {
		log.Printf("ERROR: Invalid user id format '%s'\n", requestedUserId)
		httpError := web.HttpError{Code: web.ErrorCodeInvalidFormat}
		web.JsonResponse(res, httpError, http.StatusNotFound)
		return true
	}
	return false
}
