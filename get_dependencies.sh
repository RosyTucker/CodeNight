set -e
# Auth
go get -v github.com/dgrijalva/jwt-go
go get -v github.com/dgrijalva/jwt-go/request
go get -v github.com/google/go-github/github
go get -v golang.org/x/oauth2

# Util
go get -v github.com/apsdehal/go-logger
go get -v github.com/pkg/errors
go get -v github.com/gorilla/mux
go get -v github.com/asaskevich/govalidator
go get -v github.com/gorilla/handlers

# DB
go get -v gopkg.in/mgo.v2
go get -v gopkg.in/mgo.v2/bson