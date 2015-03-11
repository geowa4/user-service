package user

import (
	"bytes"
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
	recorder := httptest.NewRecorder()
	container.ServeHTTP(recorder, httpRequest)

	assert := assert.New(t)
	assert.Equal(200, recorder.Code)
	userList := []user{}
	json.NewDecoder(recorder.Body).Decode(&userList)
	assert.Equal("1", userList[0].ID)
	assert.Equal("Hay-bay-bay", userList[0].Name)
}

func TestListNoUsers(t *testing.T) {
	container := restful.NewContainer()
	ConfigureResource(container)
	httpReq, _ := http.NewRequest("GET", "/users", nil)
	httpReq.Header.Set("Accept", restful.MIME_JSON)
	recorder := httptest.NewRecorder()
	container.ServeHTTP(recorder, httpReq)

	assert := assert.New(t)
	assert.Equal(404, recorder.Code)
	assert.Equal("Users could not be found.", recorder.Body.String())
}

func TestGetUser(t *testing.T) {
	container := restful.NewContainer()
	u := ConfigureResource(container)
	u.dao = map[string]user{
		"1": user{
			ID:   "1",
			Name: "Hay-bay-bay",
		},
	}
	httpReq, _ := http.NewRequest("GET", "/users/1", nil)
	httpReq.Header.Set("Accept", restful.MIME_JSON)
	recorder := httptest.NewRecorder()
	container.ServeHTTP(recorder, httpReq)

	assert := assert.New(t)
	assert.Equal(200, recorder.Code)
	assert.Equal(restful.MIME_JSON, recorder.Header().Get("Content-Type"))
	user := user{}
	json.NewDecoder(recorder.Body).Decode(&user)
	assert.Equal("1", user.ID)
}

func TestGetNonexistentUser(t *testing.T) {
	container := restful.NewContainer()
	ConfigureResource(container)
	httpReq, _ := http.NewRequest("GET", "/users/2", nil)
	httpReq.Header.Set("Accept", restful.MIME_JSON)
	recorder := httptest.NewRecorder()
	container.ServeHTTP(recorder, httpReq)

	assert := assert.New(t)
	assert.Equal(404, recorder.Code)
	assert.Equal("text/plain", recorder.Header().Get("Content-Type"))
	assert.Equal("User could not be found.", recorder.Body.String())
}

func TestCreateUser(t *testing.T) {
	container := restful.NewContainer()
	ConfigureResource(container)
	userJSONBuffer := bytes.NewBuffer([]byte{})
	json.NewEncoder(userJSONBuffer).Encode(user{
		Name: "Hay-bay-bay",
	})
	httpReq, err := http.NewRequest("POST", "/users", userJSONBuffer)
	if err != nil {
		t.Fatal(err)
	}
	httpReq.Header.Set("Content-Type", restful.MIME_JSON)
	httpReq.Header.Set("Accept", restful.MIME_JSON)
	recorder := httptest.NewRecorder()
	container.ServeHTTP(recorder, httpReq)

	assert := assert.New(t)
	assert.Equal(201, recorder.Code)
	assert.Equal(restful.MIME_JSON, recorder.Header().Get("Content-Type"))
	user := user{}
	json.NewDecoder(recorder.Body).Decode(&user)
	assert.Equal("1", user.ID)
	assert.Equal("Hay-bay-bay", user.Name)
}

func TestDeleteUser(t *testing.T) {
	container := restful.NewContainer()
	u := ConfigureResource(container)
	u.dao = map[string]user{
		"1": user{
			ID:   "1",
			Name: "Hay-bay-bay",
		},
	}
	httpReq, _ := http.NewRequest("DELETE", "/users/1", nil)
	httpReq.Header.Set("Accept", restful.MIME_JSON)
	recorder := httptest.NewRecorder()
	container.ServeHTTP(recorder, httpReq)

	assert := assert.New(t)
	assert.Equal(200, recorder.Code)
	assert.Equal(restful.MIME_JSON, recorder.Header().Get("Content-Type"))
	user := user{}
	json.NewDecoder(recorder.Body).Decode(&user)
	assert.Equal("1", user.ID)
	assert.Equal("Hay-bay-bay", user.Name)
}

func TestDeleteNonexistentUser(t *testing.T) {
	container := restful.NewContainer()
	ConfigureResource(container)
	httpReq, _ := http.NewRequest("DELETE", "/users/1", nil)
	httpReq.Header.Set("Accept", restful.MIME_JSON)
	recorder := httptest.NewRecorder()
	container.ServeHTTP(recorder, httpReq)

	assert := assert.New(t)
	assert.Equal(404, recorder.Code)
	assert.Equal("text/plain", recorder.Header().Get("Content-Type"))
	assert.Equal("User could not be found.", recorder.Body.String())
}
