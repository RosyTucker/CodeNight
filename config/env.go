package config

import (
	"crypto/rsa"
	"github.com/dgrijalva/jwt-go"
	"os"
	"strconv"
	"time"
)

type Env struct {
	Port                  string
	GithubKey             string
	GithubSecret          string
	GithubCallbackUrl     string
	GithubStateString     string
	JwtPrivateKey         *rsa.PrivateKey
	JwtPublicKey          *rsa.PublicKey
	JwtExpiryHours        time.Duration
	PostLoginRedirect     string
	MongoConnectionString string
	MasterUser            string
}

func GetEnv() Env {
	jwtPrivateBytes := []byte(noDefault("JWT_PRIVATE_KEY_BYTES"))
	jwtPublicBytes := []byte(noDefault("JWT_PUBLIC_KEY_BYTES"))

	jwtPrivateKey, err := jwt.ParseRSAPrivateKeyFromPEM(jwtPrivateBytes)
	if err != nil {
		Log.ErrorF("Failed to get JWT private key '%+v'", err)
	}

	jwtPublicKey, err := jwt.ParseRSAPublicKeyFromPEM(jwtPublicBytes)
	if err != nil {
		Log.ErrorF("Failed to get JWT public key '%+v'", err)
	}

	return Env{
		Port:                  noDefault("PORT"),
		GithubKey:             noDefault("GITHUB_KEY"),
		GithubSecret:          noDefault("GITHUB_SECRET"),
		GithubCallbackUrl:     noDefault("GITHUB_CALLBACK_URL"),
		GithubStateString:     noDefault("GITHUB_STATE_STRING"),
		JwtPrivateKey:         jwtPrivateKey,
		JwtPublicKey:          jwtPublicKey,
		JwtExpiryHours:        time.Duration(noDefaultInt("JWT_EXPIRY_HOURS")),
		PostLoginRedirect:     noDefault("POST_LOGIN_REDIRECT"),
		MongoConnectionString: noDefault("MONGO_CONNECTION_STRING"),
		MasterUser:            noDefault("MASTER_USER")}
}

func defaultWhenEmpty(key string, defaultValue string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return defaultValue
	}
	return value
}

func noDefault(key string) string {
	value := os.Getenv(key)

	if len(value) != 0 {
		return value
	}

	panic("Missing environment variable for key: " + key)
}

func noDefaultInt(key string) int {
	value := os.Getenv(key)

	if len(value) != 0 {
		intValue, err := strconv.Atoi(value)
		if err != nil {
			panic(err)
		}
		return intValue
	}

	panic("Missing integer environment variable for key: " + key)
}
