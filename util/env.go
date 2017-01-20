package util

import (
	"os"
)

type Env struct {
	Port              string
	GithubKey         string
	GithubSecret      string
	GithubCallbackUrl string
	GithubStateString string
}

func GetEnv() Env {
	const defaultPort = "4000"
	return Env{
		Port:              defaultWhenEmpty("PORT", defaultPort),
		GithubKey:         noDefault("GITHUB_KEY"),
		GithubSecret:      noDefault("GITHUB_SECRET"),
		GithubCallbackUrl: noDefault("GITHUB_CALLBACK_URL"),
		GithubStateString: noDefault("GITHUB_STATE_STRING")}
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
