package github

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMakeGetRequest(t *testing.T) {
	assert := assert.New(t)
	req, err := makeGitHubRequest(
		"GET",
		"https://api.github.com/user",
		nil,
		"accessToken",
	)
	assert.Nil(err)
	h := req.Header
	assert.Equal("application/json", h.Get("Accept"))
	assert.Equal("token accessToken", h.Get("Authorization"))
	assert.Equal("", h.Get("Content-Type"))
	assert.Equal("GET", req.Method)
	assert.Equal("api.github.com", req.URL.Host)
}

func TestMakeGetRequestWithNoToken(t *testing.T) {
	assert := assert.New(t)
	req, err := makeGitHubRequest(
		"GET",
		"https://api.github.com/user",
		nil,
		"",
	)
	assert.Nil(err)
	h := req.Header
	assert.Equal("", h.Get("Authorization"))
}

func TestMakePostRequest(t *testing.T) {
	assert := assert.New(t)
	req, err := makeGitHubRequest(
		"POST",
		"https://github.com/login/oauth/access_token",
		bytes.NewReader([]byte{}),
		"",
	)
	assert.Nil(err)
	h := req.Header
	assert.Equal("application/json", h.Get("Accept"))
	assert.Equal("", h.Get("Authorization"))
	assert.Equal("application/json", h.Get("Content-Type"))
	assert.Equal("POST", req.Method)
	assert.Equal("github.com", req.URL.Host)
}
