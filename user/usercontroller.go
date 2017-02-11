package user

import (
	"github.com/gorilla/mux"
	"github.com/rosytucker/codenight/config"
	"github.com/rosytucker/codenight/github"
	"github.com/rosytucker/codenight/web"
	"log"
	"net/http"
)

var environment = config.GetEnv()

func AddRoutes(router *mux.Router) {
	router.HandleFunc("/user/current", web.RequiresAuth(getCurrentUserHandler)).Methods(http.MethodGet)
	router.HandleFunc("/user/{userId:[A-Za-z0-9]+}", getUserHandler).Methods(http.MethodGet)
	router.HandleFunc("/user/{userId:[A-Za-z0-9]+}", web.RequiresAuth(putUserHandler)).Methods(http.MethodPut)
	router.HandleFunc("/login", loginHandler).Methods(http.MethodGet)
	router.HandleFunc("/oauthCallback", oauthCallbackHandler).Methods(http.MethodGet)
}

func getUserHandler(res http.ResponseWriter, req *http.Request) {
	requestedUserId := mux.Vars(req)["userId"]

	log.Printf("Getting user with Id '%+v' \n", requestedUserId)

	if !ValidUserId(requestedUserId) {
		log.Printf("ERROR: Invalid User Id format: %+v \n", requestedUserId)
		httpError := web.HttpError{Code: web.ErrorCodeInvalidFormat, Message: "userId was not formatted correctly"}
		web.JsonResponse(res, httpError, http.StatusBadRequest)
		return
	}

	foundUser, err := GetPublicById(requestedUserId)

	if err != nil {
		log.Printf("ERROR: Failed to find user %+v \n", err)
		httpError := web.HttpError{Code: web.ErrorCodeNotFound}
		web.JsonResponse(res, httpError, http.StatusNotFound)
		return
	}

	log.Printf("SUCCESS: Fetched user with id: %+v \n", requestedUserId)
	web.JsonResponse(res, foundUser, http.StatusOK)
}

func getCurrentUserHandler(res http.ResponseWriter, req *http.Request, claims *web.JwtClaims) {
	requestedUserId := claims.UserId

	log.Printf("Finding current foundUser with id: %+v \n", requestedUserId)

	foundUser, err := GetById(requestedUserId)

	if err != nil {
		log.Printf("ERROR: Failed to find current foundUser %+v \n", err)
		httpError := web.HttpError{Code: web.ErrorCodeNotFound}
		web.JsonResponse(res, httpError, http.StatusNotFound)
		return
	}
	log.Printf("SUCCESS: Fetched current user with id: %+v \n", foundUser)
	web.JsonResponse(res, foundUser, http.StatusOK)
}

func putUserHandler(res http.ResponseWriter, req *http.Request, claims *web.JwtClaims) {
	requestedUserId := mux.Vars(req)["userId"]

	log.Printf("Request for Put User with Id '%s' from '%s' \n", requestedUserId, claims.UserId)

	if !ValidUserId(requestedUserId) {
		log.Printf("ERROR: Invalid User Id format: %+v \n", requestedUserId)
		httpError := web.HttpError{Code: web.ErrorCodeInvalidFormat, Message: "userId was not formatted correctly"}
		web.JsonResponse(res, httpError, http.StatusBadRequest)
		return
	}

	if requestedUserId != claims.UserId || !claims.IsAdmin {
		log.Printf("ERROR: User '%s' tried to edit user '%s' \n", claims.UserId, requestedUserId)
		httpError := web.HttpError{Code: web.ErrorCodeForbidden, Message: "you can only update yourself"}
		web.JsonResponse(res, httpError, http.StatusForbidden)
		return
	}

	putUser, err := FromJsonBody(req.Body)

	if err != nil {
		log.Printf("ERROR: Failed to read body as json user %+v \n", err)
		httpError := web.HttpError{Code: web.ErrorCodeInvalidFormat, Message: err.Error()}
		web.JsonResponse(res, httpError, http.StatusBadRequest)
		return
	}

	err = Replace(requestedUserId, putUser)

	if err != nil {
		log.Printf("ERROR: Failed to PUT user %+v \n", err)
		httpError := web.HttpError{Code: web.ErrorCodeServerError}
		web.JsonResponse(res, httpError, http.StatusInternalServerError)
		return
	}

	log.Printf("SUCCESS: PUT User with Id '%+v' \n", requestedUserId)
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

	user := &User{
		Name:      githubUser.Name,
		Token:     web.EncodeJson(token),
		UserName:  *githubUser.Login,
		Email:     githubUser.Email,
		Blog:      githubUser.Blog,
		Location:  githubUser.Location,
		AvatarUrl: githubUser.AvatarURL,
		IsAdmin:   *githubUser.Login == environment.MasterUser}

	log.Printf("Creating User with username '%s' if they dont already exist \n", user.UserName)
	userId, err := CreateIfNotExists(user)

	if err != nil {
		log.Printf("ERROR: Failed to create user %+v \n", err)
		http.Redirect(res, req, environment.PostLoginRedirect, http.StatusTemporaryRedirect)
		return
	}
	log.Printf("Updated User with Id '%+v', adding JWT \n", userId)

	web.SetJwt(res, req, userId, user.IsAdmin)

	res.Header().Set("Location", "/user/"+userId)
	http.Redirect(res, req, environment.PostLoginRedirect, http.StatusTemporaryRedirect)
}
