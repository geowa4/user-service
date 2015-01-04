package github

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"
)

func makeGitHubRequest(method string, url string, body io.Reader, accessToken string) (req *http.Request, err error) {
	req, err = http.NewRequest(
		method,
		url,
		body,
	)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	if accessToken != "" {
		req.Header.Set("Authorization", "token "+accessToken)
	}
	if method == "POST" {
		req.Header.Set("Content-Type", "application/json")
	}
	return
}

// TokenRequest represents the data that must be sent to GitHub to get a token.
type TokenRequest struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	Code         string `json:"code"`
}

// TokenResponse represents the response from GitHub.
type TokenResponse struct {
	AccessToken string `json:"access_token"`
	Scope       string `json:"scope"`
	TokenType   string `json:"token_type"`
}

func postCodeToGitHub(code string) (*TokenResponse, error) {
	tokenReqBytes, err := json.Marshal(TokenRequest{
		os.Getenv("GITHUB_CLIENT_ID"),
		os.Getenv("GITHUB_CLIENT_SECRET"),
		code,
	})
	if err != nil {
		return nil, err
	}
	req, err := makeGitHubRequest(
		"POST",
		"https://github.com/login/oauth/access_token",
		bytes.NewReader(tokenReqBytes),
		"",
	)
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	token := new(TokenResponse)
	err = json.NewDecoder(resp.Body).Decode(token)
	if err != nil {
		return nil, err
	}

	return token, err
}

// EmailResponse represents the response for an email entity
type EmailResponse struct {
	Email    string `json:"email"`
	Primary  bool   `json:"primary"`
	Verified bool   `json:"verified"`
}

func getPrimaryEmail(accessToken string) (string, error) {
	req, err := makeGitHubRequest(
		"GET",
		"https://api.github.com/user",
		nil,
		accessToken,
	)
	if err != nil {
		return "", err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var emails []EmailResponse
	err = json.NewDecoder(resp.Body).Decode(emails)
	if err != nil {
		return "", err
	}
	for _, email := range emails {
		if email.Primary {
			return email.Email, nil
		}
	}
	return "", errors.New("No primary email found.")
}

// UserResponse represents the GitHub user response.
type UserResponse struct {
	Email     string `json:"email"`
	Name      string `json:"name"`
	Login     string `json:"login"`
	ID        int64  `json:"id"`
	AvatarURL string `json:"avatar_url"`
	HTMLURL   string `json:"html_url"`
}

func getGitHubUser(accessToken string) (*UserResponse, error) {
	req, err := makeGitHubRequest(
		"GET",
		"https://api.github.com/user",
		nil,
		accessToken,
	)
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	user := new(UserResponse)
	err = json.NewDecoder(resp.Body).Decode(user)
	if err != nil {
		return nil, err
	}
	if user.Email == "" {
		// Ignore errors since there's nothing we can do.
		user.Email, _ = getPrimaryEmail(accessToken)
	}
	return user, err
}
