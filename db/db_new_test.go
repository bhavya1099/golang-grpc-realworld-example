package db

import (
	"testing"
	"time"
	"github.com/stretchr/testify/assert"
	"github.com/jinzhu/gorm"
	"github.com/raahii/golang-grpc-realworld-example/db"
)

type T struct {
	common
	isEnvSet bool
	context  *testContext // For running tests and subtests.
}

type T struct {
	common
	isEnvSet bool
	context  *testContext // For running tests and subtests.
}


type Time struct {
	// wall and ext encode the wall time seconds, wall time nanoseconds,
	// and optional monotonic clock reading in nanoseconds.
	//
	// From high to low bit position, wall encodes a 1-bit flag (hasMonotonic),
	// a 33-bit seconds field, and a 30-bit wall time nanoseconds field.
	// The nanoseconds field is in the range [0, 999999999].
	// If the hasMonotonic bit is 0, then the 33-bit field must be zero
	// and the full signed 64-bit wall seconds since Jan 1 year 1 is stored in ext.
	// If the hasMonotonic bit is 1, then the 33-bit field holds a 33-bit
	// unsigned wall seconds since Jan 1 year 1885, and ext holds a
	// signed 64-bit monotonic clock reading, nanoseconds since process start.
	wall uint64
	ext  int64

	// loc specifies the Location that should be used to
	// determine the minute, hour, month, day, and year
	// that correspond to this Time.
	// The nil location means UTC.
	// All UTC times are represented with loc==nil, never loc==&utcLoc.
	loc *Location
}
func TestNew(t *testing.T) {
	t.Run("Successful database connection on first attempt", func(t *testing.T) {

		dbConn, err := db.New()

		assert.NotNil(t, dbConn)
		assert.NoError(t, err)
	})

	t.Run("Database connection attempt retries 10 times", func(t *testing.T) {

		start := time.Now()
		_, err := db.New()
		elapsed := time.Since(start)

		assert.Error(t, err)
		assert.GreaterOrEqual(t, elapsed.Milliseconds(), 10000)
	})

	t.Run("Maximum idle connections set to 3 after successful connection", func(t *testing.T) {

		dbConn, _ := db.New()

		maxIdleConns := dbConn.DB().MaxIdleConns()
		assert.Equal(t, 3, maxIdleConns)
	})

	t.Run("Logging mode disabled by default after connection", func(t *testing.T) {

		dbConn, _ := db.New()

		assert.False(t, dbConn.LogModeEnabled())
	})

	t.Run("Error returned when unable to establish a database connection", func(t *testing.T) {

		_, err := db.New()

		assert.Error(t, err)
	})
}
