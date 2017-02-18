package web

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
	"github.com/rosytucker/codenight/config"
	"html/template"
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
			config.Log.InfoF("User not logged in or Invalid JWT '%+v'", err)
			httpError := HttpError{
				Code:    ErrorCodeUnauthorized,
				Message: "You must be logged in try and view user information!"}
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

type formTemplateViewModel struct {
	RedirectUri string
	Token       string
}

func SetJwt(res http.ResponseWriter, req *http.Request, userId string, isAdmin bool) {
	token := jwt.New(jwt.SigningMethodRS256)
	expiry := time.Now().Add(time.Hour * environment.JwtExpiryHours)

	config.Log.InfoF("Creating JWT for '%s'", userId)

	token.Claims = JwtClaims{
		userId,
		isAdmin,
		jwt.StandardClaims{ExpiresAt: expiry.Unix()}}

	jwtToken, err := token.SignedString(environment.JwtPrivateKey)

	if err != nil {
		config.Log.ErrorF("Failed to create jwt '%+v'", err)
		http.Redirect(res, req, environment.PostLoginRedirect, http.StatusTemporaryRedirect)
		return
	}

	formTemp, err := template.ParseFiles("./web/formPost.html")

	if err != nil {
		config.Log.ErrorF("Failed to create post form template '%+v'", err)
		http.Redirect(res, req, environment.PostLoginRedirect, http.StatusTemporaryRedirect)
		return
	}

	formTemp.Execute(res, formTemplateViewModel{RedirectUri: environment.PostLoginRedirect, Token: jwtToken})
}
