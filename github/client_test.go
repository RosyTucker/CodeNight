package github

import (
	"github.com/rosytucker/codenight/config"
	"github.com/stretchr/testify/assert"
	"testing"
	"net/http/httptest"
	"net/http"
)

var env = config.Env{
	GithubKey:         "some_key",
	GithubSecret:      "some_secret",
	GithubStateString: "some_state_string",
}

func TestClient_LoginRedirectUrl(t *testing.T) {
	client := NewClient(env)

	expected := "https://github.com/login/oauth/authorize?access_type=online&client_id=some_key&response_type=code&scope=user&state=some_state_string"
	actual := client.LoginRedirectUrl()

	assert.Equal(t, expected, actual)
}

func TestClient_GetUser(t *testing.T) {
	client := NewClient(env)

	req := httptest.NewRequest(http.MethodGet, "/some target", nil)
	user, error := client.GetUser(req)

	assert.Nil(t, error)

}