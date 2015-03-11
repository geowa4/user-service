package user

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/emicklei/go-restful"
	"github.com/stretchr/testify/assert"
)

func TestListUsers(t *testing.T) {
	container := restful.NewContainer()
	u := ConfigureResource(container)
	u.dao = map[string]user{
		"1": user{
			ID:   "1",
			Name: "Hay-bay-bay",
		},
	}
	httpRequest, _ := http.NewRequest("GET", "/users", nil)
	httpRequest.Header.Set("Accept", restful.MIME_JSON)
	httpWriter := httptest.NewRecorder()
	container.ServeHTTP(httpWriter, httpRequest)

	assert := assert.New(t)
	assert.Equal(200, httpWriter.Code)
	userList := []user{}
	json.NewDecoder(httpWriter.Body).Decode(&userList)
	assert.Equal("1", userList[0].ID)
	assert.Equal("Hay-bay-bay", userList[0].Name)
}

func TestListNoUsers(t *testing.T) {
	container := restful.NewContainer()
	ConfigureResource(container)
	httpRequest, _ := http.NewRequest("GET", "/users", nil)
	httpRequest.Header.Set("Accept", restful.MIME_JSON)
	httpWriter := httptest.NewRecorder()
	container.ServeHTTP(httpWriter, httpRequest)

	assert := assert.New(t)
	assert.Equal(404, httpWriter.Code)
	assert.Equal("Users could not be found.", httpWriter.Body.String())
}
