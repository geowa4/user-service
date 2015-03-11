package health

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
	ConfigureResource(container)
	httpRequest, _ := http.NewRequest("GET", "/health", nil)
	httpRequest.Header.Set("Accept", restful.MIME_JSON)
	httpWriter := httptest.NewRecorder()
	container.ServeHTTP(httpWriter, httpRequest)

	assert := assert.New(t)
	assert.Equal(200, httpWriter.Code)
	m := new(measurements)
	json.NewDecoder(httpWriter.Body).Decode(m)
	assert.True(m.NumCPU > 0)
	assert.True(m.NumGoroutine > 0)
	assert.True(m.MemoryAllocated > 0)
}
