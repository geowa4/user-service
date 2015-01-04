package user

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestSave(t *testing.T) {
	assert := assert.New(t)
	db, err := sqlmock.New()
	assert.Nil(err)
	columns := []string{"id"}
	var expectedID int64
	expectedID = 1
	sqlmock.ExpectQuery("SELECT upsert_user\\((.+)\\)").
		WithArgs("email@foo.bar", "somename", 11, "ghlogin", "ghavatar", "ghhtml", "ghtoken", "ghscope").
		WillReturnRows(sqlmock.NewRows(columns).AddRow(expectedID))
	user := User{
		ID:              expectedID,
		Email:           "email@foo.bar",
		Name:            "somename",
		AccessToken:     "ghtoken",
		Scope:           "ghscope",
		GitHubID:        11,
		GitHubLogin:     "ghlogin",
		GitHubAvatarURL: "ghavatar",
		GitHubHTMLURL:   "ghhtml",
	}
	id, err := user.Save(db)
	assert.Nil(err)
	assert.Equal(expectedID, id)
	assert.Equal(expectedID, user.ID)
	assert.Nil(db.Close())
}

func TestLoadAll(t *testing.T) {
	assert := assert.New(t)
	db, err := sqlmock.New()
	assert.Nil(err)
	columns := []string{"user_id", "email", "name", "github_id", "github_login", "github_avatar_url", "github_html_url", "github_access_token", "github_scope"}
	sqlmock.ExpectQuery("SELECT (.+) FROM users").
		WithArgs().
		WillReturnRows(sqlmock.NewRows(columns).FromCSVString("1,email@foo.bar,somename,11,ghlogin,ghavatar,ghhtml,ghtoken,ghscope\n2,email2@foo.bar,othername,22,ghlogin2,ghavatar2,ghhtml2,ghtoken2,ghscope2"))
	users, err := LoadAll(db)
	assert.Nil(err)
	assert.Len(users, 2)
	assert.Equal(1, users[0].ID)
	assert.Equal(2, users[1].ID)
	assert.Nil(db.Close())
}

func TestLoadByGitHubID(t *testing.T) {
	assert := assert.New(t)
	db, err := sqlmock.New()
	assert.Nil(err)
	columns := []string{"user_id", "email", "name", "github_id", "github_login", "github_avatar_url", "github_html_url", "github_access_token", "github_scope"}
	sqlmock.ExpectQuery("SELECT (.+) FROM users WHERE github_id = \\$1 LIMIT 1").
		WithArgs(11).
		WillReturnRows(sqlmock.NewRows(columns).FromCSVString("1,email@foo.bar,somename,11,ghlogin,ghavatar,ghhtml,ghtoken,ghscope"))
	user, err := LoadByGitHubID(db, 11)
	assert.Nil(err)
	assert.Equal(1, user.ID)
	assert.Nil(db.Close())
}
