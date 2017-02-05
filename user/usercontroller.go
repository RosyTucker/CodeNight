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
	router.HandleFunc("/user/{userId:[A-Za-z0-9]+}", putUserHandler).Methods(http.MethodPut)
	router.HandleFunc("/login", loginHandler).Methods(http.MethodGet)
	router.HandleFunc("/oauthCallback", oauthCallbackHandler).Methods(http.MethodGet)
}

func getUserHandler(res http.ResponseWriter, req *http.Request) {
	userId := mux.Vars(req)["userId"]

	log.Printf("Getting User with Id '%+v' \n", userId)

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

func getCurrentUserHandler(res http.ResponseWriter, req *http.Request, claims *web.JwtClaims) {
	userId := claims.UserId

	log.Printf("Finding current user with id: %+v \n", userId)

	user, err := GetById(userId)

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
