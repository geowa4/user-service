package health

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestGoodDB(t *testing.T) {
	assert := assert.New(t)
	db, err := sqlmock.New()
	assert.Nil(err)
	sqlmock.ExpectExec("SELECT 1 AS ping").
		WithArgs().
		WillReturnResult(sqlmock.NewResult(1, 1))
	hr := NewRouter(db, nil)
	measurements := hr.GetSimpleMeasurements()
	assert.True(measurements.NumCPU > 0)
	assert.True(measurements.NumGoroutine > 0)
	assert.True(measurements.MemoryAllocated > 0)
	assert.True(measurements.DBConnected)
	assert.Empty(measurements.DBPingError)
	assert.Nil(db.Close())
}

func TestBadDB(t *testing.T) {
	assert := assert.New(t)
	dbErrMsg := "test db error"
	db, err := sqlmock.New()
	assert.Nil(err)
	sqlmock.ExpectExec("SELECT 1 AS ping").
		WithArgs().
		WillReturnError(errors.New(dbErrMsg))
	hr := NewRouter(db, nil)
	measurements := hr.GetSimpleMeasurements()
	assert.True(measurements.NumCPU > 0)
	assert.True(measurements.NumGoroutine > 0)
	assert.True(measurements.MemoryAllocated > 0)
	assert.False(measurements.DBConnected)
	assert.Equal(measurements.DBPingError, dbErrMsg)
	assert.Nil(db.Close())
}
