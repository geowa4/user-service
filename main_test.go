package main

import (
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPing(t *testing.T) {
	assert := assert.New(t)
	recorder := httptest.NewRecorder()
	ping(recorder, nil)
	assert.Equal(200, recorder.Code)
	assert.Equal("text/plain; charset=utf-8", recorder.Header().Get("Content-Type"))
	assert.Equal("OK", recorder.Body.String())
}
