package health

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTakeMeasurements(t *testing.T) {
	m := takeMeasurements()
	assert := assert.New(t)
	assert.True(m.NumCPU > 0)
	assert.True(m.NumGoroutine > 0)
	assert.True(m.MemoryAllocated > 0)
}
