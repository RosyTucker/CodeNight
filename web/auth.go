package web

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
	"github.com/rosytucker/codenight/config"
	"log"
	"net/http"
	"time"
)

var environment = config.GetEnv()

func RequiresAuth(next func(http.ResponseWriter, *http.Request, *JwtClaims)) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {

		token, err := request.ParseFromRequestWithClaims(req, request.OAuth2Extractor, &JwtClaims{}, func(token *jwt.Token) (interface{}, error) {
			return environment.JwtPublicKey, nil
		})

		if err != nil {
			log.Printf("ERROR: User not logged in or Invalid JWT '%+v' \n", err)
			httpError := HttpError{
				Code:    ErrorCodeUnauthorized,
				Message: "you must be logged in try and view user information"}
			JsonResponse(res, httpError, http.StatusUnauthorized)
			return
		}
		next(res, req, token.Claims.(*JwtClaims))
	}
}

type JwtClaims struct {
	UserId  string
	IsAdmin bool
	jwt.StandardClaims
}

func SetJwt(res http.ResponseWriter, req *http.Request, userId string, isAdmin bool) {
	token := jwt.New(jwt.SigningMethodRS256)
	expiry := time.Now().Add(time.Hour * environment.JwtExpiryHours)

	log.Printf("Creating JWT for '%s' \n", userId)

	token.Claims = JwtClaims{
		userId,
		isAdmin,
		jwt.StandardClaims{ExpiresAt: expiry.Unix()}}

	jwt, err := token.SignedString(environment.JwtPrivateKey)

	if err != nil {
		log.Printf("ERROR: Failed to create jwt '%+v' \n", err)
		http.Redirect(res, req, environment.PostLoginRedirect, http.StatusTemporaryRedirect)
		return
	}

	cookie := http.Cookie{Name: "Auth", Value: jwt, Expires: expiry, HttpOnly: true}
	http.SetCookie(res, &cookie)
}
