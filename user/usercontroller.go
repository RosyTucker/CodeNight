package user

import (
	"github.com/gorilla/mux"
	"github.com/rosytucker/codenight/config"
	"github.com/rosytucker/codenight/github"
	"github.com/rosytucker/codenight/web"
	"net/http"
)

var environment = config.GetEnv()

func AddRoutes(router *mux.Router) {
	router.HandleFunc("/user", web.RequiresAuth(getCurrentUserHandler)).Methods(http.MethodGet)
	router.HandleFunc("/user/{userId:[A-Za-z0-9]+}", getUserHandler).Methods(http.MethodGet)
	router.HandleFunc("/user/{userId:[A-Za-z0-9]+}", web.RequiresAuth(putUserHandler)).Methods(http.MethodPut)
	router.HandleFunc("/login", loginHandler).Methods(http.MethodGet)
	router.HandleFunc("/oauthCallback", oauthCallbackHandler).Methods(http.MethodGet)
}

func getUserHandler(res http.ResponseWriter, req *http.Request) {
	requestedUserId := mux.Vars(req)["userId"]

	config.Log.InfoF("Getting user with Id '%+v'", requestedUserId)

	if !ValidUserId(requestedUserId) {
		config.Log.ErrorF("Invalid User Id format: '%+v'", requestedUserId)
		httpError := web.HttpError{Code: web.ErrorCodeInvalidFormat, Message: "userId was not formatted correctly"}
		web.JsonResponse(res, httpError, http.StatusBadRequest)
		return
	}

	foundUser, err := GetPublicById(requestedUserId)

	if err != nil {
		config.Log.ErrorF("Failed to find user '%+v'", err)
		httpError := web.HttpError{Code: web.ErrorCodeNotFound}
		web.JsonResponse(res, httpError, http.StatusNotFound)
		return
	}

	config.Log.InfoF(" Fetched user with id: '%+v'", requestedUserId)

	web.JsonResponse(res, foundUser, http.StatusOK)
}

func getCurrentUserHandler(res http.ResponseWriter, req *http.Request, claims *web.JwtClaims) {
	requestedUserId := claims.UserId

	config.Log.InfoF("Finding current user with id: '%+v'", requestedUserId)

	foundUser, err := GetById(requestedUserId)

	if err != nil {
		config.Log.ErrorF("Failed to find current user: '%+v'", err)
		httpError := web.HttpError{Code: web.ErrorCodeNotFound}
		web.JsonResponse(res, httpError, http.StatusNotFound)
		return
	}
	config.Log.InfoF("Fetched current user with id: '%+v'", requestedUserId)
	web.JsonResponse(res, foundUser, http.StatusOK)
}

func putUserHandler(res http.ResponseWriter, req *http.Request, claims *web.JwtClaims) {
	requestedUserId := mux.Vars(req)["userId"]

	config.Log.InfoF("Request for Put User with Id '%s' from '%s'", requestedUserId, claims.UserId)

	if !ValidUserId(requestedUserId) {
		config.Log.ErrorF("Invalid User Id format: '%+v'", requestedUserId)
		httpError := web.HttpError{Code: web.ErrorCodeInvalidFormat, Message: "userId was not formatted correctly"}
		web.JsonResponse(res, httpError, http.StatusBadRequest)
		return
	}

	if requestedUserId != claims.UserId && !claims.IsAdmin {
		config.Log.ErrorF("User '%s' tried to edit user '%s'", claims.UserId, requestedUserId)
		httpError := web.HttpError{Code: web.ErrorCodeForbidden, Message: "you can only update yourself"}
		web.JsonResponse(res, httpError, http.StatusForbidden)
		return
	}

	putUser, err := PublicFromJsonBody(req.Body)

	if err != nil {
		config.Log.ErrorF("Failed to read body as json user '%+v'", err)
		httpError := web.HttpError{Code: web.ErrorCodeInvalidFormat, Message: err.Error()}
		web.JsonResponse(res, httpError, http.StatusBadRequest)
		return
	}

	errors := putUser.Validate()
	if len(errors) != 0 {
		config.Log.ErrorF("PUT user format is invalid '%+v' Errors: '%+v'", putUser, errors)
		httpError := web.HttpError{Code: web.ErrorCodeInvalidFormat, ValidationErrors: errors}
		web.JsonResponse(res, httpError, http.StatusBadRequest)
		return
	}

	err = Replace(requestedUserId, putUser)

	if err != nil {
		config.Log.ErrorF("Failed to PUT user: '%+v'", err)
		httpError := web.HttpError{Code: web.ErrorCodeServerError}
		web.JsonResponse(res, httpError, http.StatusInternalServerError)
		return
	}

	config.Log.InfoF("PUT User with Id '%+v'", requestedUserId)
	res.Header().Set("Location", "/user/"+requestedUserId)
	res.WriteHeader(http.StatusNoContent)
}

func loginHandler(res http.ResponseWriter, req *http.Request) {
	url := github.LoginRedirectUrl()
	http.Redirect(res, req, url, http.StatusTemporaryRedirect)
}

func oauthCallbackHandler(res http.ResponseWriter, req *http.Request) {
	token, err := github.GetToken(req)

	if err != nil {
		config.Log.ErrorF("Failed to get github token '%+v'", err)
		http.Redirect(res, req, environment.PostLoginRedirect, http.StatusTemporaryRedirect)
		return
	}

	githubUser, err := github.GetUser(token)

	if err != nil {
		config.Log.ErrorF("Failed to get github user '%+v'", err)
		http.Redirect(res, req, environment.PostLoginRedirect, http.StatusTemporaryRedirect)
		return
	}

	newUser := &User{
		Name:      defaultString(githubUser.Name),
		Token:     web.EncodeJson(token),
		UserName:  defaultString(githubUser.Login),
		Email:     defaultString(githubUser.Email),
		Blog:      defaultString(githubUser.Blog),
		Location:  defaultString(githubUser.Location),
		AvatarUrl: defaultString(githubUser.AvatarURL),
		IsAdmin:   defaultString(githubUser.Login) == environment.MasterUser}

	config.Log.InfoF("Creating User with username '%s' if they dont already exist", newUser.UserName)
	userId, err := CreateIfNotExists(newUser)

	if err != nil {
		config.Log.ErrorF("Failed to create user '%+v'", err)
		http.Redirect(res, req, environment.PostLoginRedirect, http.StatusTemporaryRedirect)
		return
	}

	config.Log.InfoF("Updated User with Id '%+v', adding JWT", userId)

	res.Header().Set("Location", "/user/"+userId)
	web.SetJwt(res, req, userId, newUser.IsAdmin)
}

func defaultString(point *string) string {
	if point == nil {
		return ""
	}
	return *point
}
