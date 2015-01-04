package wrappers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefaults(t *testing.T) {
	var (
		assert      = assert.New(t)
		innerCalled = false
		inner       = func(w http.ResponseWriter, r *http.Request) {
			innerCalled = true
			w.Write([]byte("Test"))
		}
		recorder = httptest.NewRecorder()
	)
	wrapped := Defaults(inner, "Test")
	req, err := http.NewRequest("GET", "http://www.example.com", nil)
	assert.Nil(err)
	wrapped.ServeHTTP(recorder, req)
	assert.True(innerCalled)
	assert.Equal("*", recorder.Header().Get("Access-Control-Allow-Origin"))
	assert.Equal("HEAD,OPTIONS,GET,PUT,PATCH,POST,DELETE",
		recorder.Header().Get("Access-Control-Allow-Methods"))
	assert.Equal("Accept,Content-Type,Content-Length,Accept-Encoding,"+
		"Authorization,X-Requested-With",
		recorder.Header().Get("Access-Control-Allow-Headers"))
	assert.Equal("Test", recorder.Body.String())
}

func TestDefaultsAsPreflight(t *testing.T) {
	var (
		assert      = assert.New(t)
		innerCalled = false
		inner       = func(w http.ResponseWriter, r *http.Request) {
			innerCalled = true
			w.Write([]byte("Test"))
		}
		recorder = httptest.NewRecorder()
	)
	wrapped := Defaults(inner, "Test")
	req, err := http.NewRequest("OPTIONS", "http://www.example.com", nil)
	assert.Nil(err)
	wrapped.ServeHTTP(recorder, req)
	assert.False(innerCalled)
	assert.Equal("*", recorder.Header().Get("Access-Control-Allow-Origin"))
	assert.Equal("HEAD,OPTIONS,GET,PUT,PATCH,POST,DELETE",
		recorder.Header().Get("Access-Control-Allow-Methods"))
	assert.Equal("Accept,Content-Type,Content-Length,Accept-Encoding,"+
		"Authorization,X-Requested-With",
		recorder.Header().Get("Access-Control-Allow-Headers"))
	assert.Equal("", recorder.Body.String())
}
