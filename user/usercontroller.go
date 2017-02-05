package user

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
	"github.com/gorilla/mux"
	"github.com/rosytucker/codenight/config"
	"github.com/rosytucker/codenight/github"
	"github.com/rosytucker/codenight/web"
	"log"
	"net/http"
	"time"
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
	// userId := mux.Vars(req)["userId"]

	// log.Printf("Getting User with Id '%+v' \n", userId)

	// isHandled := handleInvalidUserIdForRequest(userId, false, res, req)
	// if isHandled {
	// 	return
	// }

	// user, err := GetPublicById(userId)

	// if err != nil {
	// 	log.Printf("ERROR: Failed to find user %+v \n", err)
	// 	httpError := web.HttpError{Code: web.ErrorCodeNotFound}
	// 	web.JsonResponse(res, httpError, http.StatusNotFound)
	// 	return
	// }

	// log.Printf("SUCCESS: Fetched user with id: %+v \n", user)
	// web.JsonResponse(res, user, http.StatusOK)
}

func getCurrentUserHandler(res http.ResponseWriter, req *http.Request) {
	token, err := request.ParseFromRequestWithClaims(req, request.OAuth2Extractor, &JwtClaims{}, func(token *jwt.Token) (interface{}, error) {
		return environment.JwtPublicKey, nil
	})

	if err != nil {
		log.Printf("ERROR: User not logged in or Invalid JWT '%+v' \n", err)
		httpError := web.HttpError{
			Code:    web.ErrorCodeUnauthorized,
			Message: "you must be logged in try and view user information"}
		web.JsonResponse(res, httpError, http.StatusUnauthorized)
		return
	}

	userId := token.Claims.(*JwtClaims).UserId

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
	// userId := mux.Vars(req)["userId"]

	// log.Printf("Putting User with Id '%+v' \n", userId)

	// isHandled := handleInvalidUserIdForRequest(userId, true, res, req)
	// if isHandled {
	// 	return
	// }

	// user, err := FromJsonBody(req.Body)

	// if err != nil {
	// 	log.Printf("ERROR: Failed to read body as json user %+v \n", err)
	// 	httpError := web.HttpError{Code: web.ErrorCodeInvalidFormat, Message: err.Error()}
	// 	web.JsonResponse(res, httpError, http.StatusBadRequest)
	// 	return
	// }

	// err = Replace(userId, user)

	// if err != nil {
	// 	log.Printf("ERROR: Failed to PUT user %+v \n", err)
	// 	httpError := web.HttpError{Code: web.ErrorCodeServerError}
	// 	web.JsonResponse(res, httpError, http.StatusInternalServerError)
	// 	return
	// }

	// log.Printf("SUCCESS: PUT User with Id '%+v' \n", userId)
	// res.Header().Set("Location", "/user/"+userId)
	// res.WriteHeader(http.StatusNoContent)
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

	jwt, err := createJwt(userId, user)

	if err != nil {
		log.Printf("ERROR: Failed to create jwt '%+v' \n", err)
		http.Redirect(res, req, environment.PostLoginRedirect, http.StatusTemporaryRedirect)
		return
	}

	cookie := http.Cookie{Name: "Auth", Value: jwt, Expires: time.Now().Add(time.Hour * environment.JwtExpiryHours), HttpOnly: true}
	http.SetCookie(res, &cookie)

	res.Header().Set("Location", "/user/"+userId)
	http.Redirect(res, req, environment.PostLoginRedirect, http.StatusTemporaryRedirect)
}

func loginHandler(res http.ResponseWriter, req *http.Request) {
	url := github.LoginRedirectUrl()
	http.Redirect(res, req, url, http.StatusTemporaryRedirect)
}

type JwtClaims struct {
	UserId  string
	IsAdmin bool
	jwt.StandardClaims
}

func createJwt(userId string, user *User) (string, error) {
	// create the token
	token := jwt.New(jwt.SigningMethodRS256)

	log.Printf("User id in token '%+v' \n", userId)

	token.Claims = JwtClaims{
		userId,
		user.IsAdmin,
		jwt.StandardClaims{ExpiresAt: time.Now().Add(time.Hour * environment.JwtExpiryHours).Unix()}}

	//Sign and get the complete encoded token as string
	return token.SignedString(environment.JwtPrivateKey)
}
