package user

import (
	"database/sql"
	"fmt"

	"github.com/kisielk/sqlstruct"
)

const (
	loadByGitHubID  = "SELECT %s FROM users WHERE github_id = $1 LIMIT 1"
	upsertStatement = "SELECT upsert_user(CAST ($1 AS TEXT), CAST ($2 AS TEXT), $3, CAST ($4 AS TEXT), CAST ($5 AS TEXT), CAST ($6 AS TEXT), CAST ($7 AS TEXT), CAST ($8 AS TEXT))"
)

// User represents all pertinent user data to this service.
type User struct {
	ID              int64  `json:"user_id"           sql:"user_id"`
	Email           string `json:"email"             sql:"email"`
	Name            string `json:"name"              sql:"name"`
	AccessToken     string `json:"-"                 sql:"github_access_token"`
	Scope           string `json:"-"                 sql:"github_scope"`
	GitHubID        int64  `json:"github_id"         sql:"github_id"`
	GitHubLogin     string `json:"github_login"      sql:"github_login"`
	GitHubAvatarURL string `json:"github_avatar_url" sql:"github_avatar_url"`
	GitHubHTMLURL   string `json:"github_html_url"   sql:"github_html_url"`
}

// Users represents a slice of user objects.
type Users []User

// Save saves the user.
func (u *User) Save(db *sql.DB) (int64, error) {
	var id int64
	err := db.QueryRow(
		upsertStatement,
		u.Email,
		u.Name,
		u.GitHubID,
		u.GitHubLogin,
		u.GitHubAvatarURL,
		u.GitHubHTMLURL,
		u.AccessToken,
		u.Scope,
	).
		Scan(&id)
	if err != nil {
		return 0, err
	}
	u.ID = id
	return id, nil
}

// LoadAll returns all users in the database.
func LoadAll(db *sql.DB) (Users, error) {
	users := []User{}
	statement := fmt.Sprintf("SELECT %s FROM users", sqlstruct.Columns(User{}))
	rows, err := db.Query(statement)
	if err != nil {
		return users, err
	}
	defer rows.Close()
	for rows.Next() {
		var user User
		if err = sqlstruct.Scan(&user, rows); err != nil {
			return users, err
		}
		users = append(users, user)
	}
	return users, nil
}

// LoadByGitHubID loads a user by the GitHub ID
func LoadByGitHubID(db *sql.DB, githubID int64) (*User, error) {
	var user User
	statement := fmt.Sprintf(loadByGitHubID, sqlstruct.Columns(user))
	rows, err := db.Query(statement, githubID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		if err := sqlstruct.Scan(&user, rows); err != nil {
			return nil, err
		}
	}
	return &user, nil
}
