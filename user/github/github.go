package github

// User organizes the data we care about for the GitHub user.
type User struct {
	AccessToken string
	Scope       string
	Email       string
	Name        string
	ID          int64
	Login       string
	AvatarURL   string
	HTMLURL     string
}

// GetUserFromCode retrieves the access token and user from GitHub
// https://github.com/login/oauth/authorize?client_id=&scope=user:email&state=random
func GetUserFromCode(code string) (*User, error) {
	tokenResp, err := postCodeToGitHub(code)
	if err != nil {
		return nil, err
	}
	userResp, err := getGitHubUser(tokenResp.AccessToken)
	if err != nil {
		return nil, err
	}

	user := new(User)
	user.AccessToken = tokenResp.AccessToken
	user.Scope = tokenResp.Scope
	user.Email = userResp.Email
	user.Name = userResp.Name
	user.ID = userResp.ID
	user.Login = userResp.Login
	user.AvatarURL = userResp.AvatarURL
	user.HTMLURL = userResp.HTMLURL
	return user, nil
}
