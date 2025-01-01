package db

import (
	"database/sql"
	"errors"
	"testing"
	"github.com/DATA-DOG/go-txdb"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"github.com/raahii/golang-grpc-realworld-example/db"
)

type DB struct {
	// Total time waited for new connections.
	waitDuration atomic.Int64

	connector driver.Connector
	// numClosed is an atomic counter which represents a total number of
	// closed connections. Stmt.openStmt checks it before cleaning closed
	// connections in Stmt.css.
	numClosed atomic.Uint64

	mu           sync.Mutex    // protects following fields
	freeConn     []*driverConn // free connections ordered by returnedAt oldest to newest
	connRequests map[uint64]chan connRequest
	nextRequest  uint64 // Next key to use in connRequests.
	numOpen      int    // number of opened and pending open connections
	// Used to signal the need for new connections
	// a goroutine running connectionOpener() reads on this chan and
	// maybeOpenNewConnections sends on the chan (one send per needed connection)
	// It is closed during db.Close(). The close tells the connectionOpener
	// goroutine to exit.
	openerCh          chan struct{}
	closed            bool
	dep               map[finalCloser]depSet
	lastPut           map[*driverConn]string // stacktrace of last conn's put; debug only
	maxIdleCount      int                    // zero means defaultMaxIdleConns; negative means 0
	maxOpen           int                    // <= 0 means unlimited
	maxLifetime       time.Duration          // maximum amount of time a connection may be reused
	maxIdleTime       time.Duration          // maximum amount of time a connection may be idle before being closed
	cleanerCh         chan struct{}
	waitCount         int64 // Total number of connections waited for.
	maxIdleClosed     int64 // Total number of connections closed due to idle count.
	maxIdleTimeClosed int64 // Total number of connections closed due to idle time.
	maxLifetimeClosed int64 // Total number of connections closed due to max connection lifetime limit.

	stop func() // stop cancels the connection opener.
}

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
func TestNewTestDB(t *testing.T) {

	t.Run("Scenario 1: Test NewTestDB with successful initialization", func(t *testing.T) {
		err := godotenv.Load("../env/test.env")
		assert.NoError(t, err, "loading test environment variables should not return an error")

		dbInstance, err := db.NewTestDB()

		assert.NotNil(t, dbInstance, "database connection should not be nil")
		assert.NoError(t, err, "error should be nil")
	})

	t.Run("Scenario 2: Test NewTestDB when failing to load test environment variables", func(t *testing.T) {

	})

	t.Run("Scenario 3: Test NewTestDB with error initializing the test database connection", func(t *testing.T) {

	})

	t.Run("Scenario 4: Test NewTestDB with max idle connections set to 3", func(t *testing.T) {
		dbInstance, _ := db.NewTestDB()

		sqlDB := dbInstance.DB()
		assert.Equal(t, 3, sqlDB.Stats().MaxIdleConns, "maximum idle connections should be set to 3")
	})

	t.Run("Scenario 5: Test NewTestDB with logging mode disabled", func(t *testing.T) {
		dbInstance, _ := db.NewTestDB()

		assert.False(t, dbInstance.LogMode(gorm.LogWarn).logMode, "logging mode should be disabled")
	})

	t.Run("Scenario 6: Test NewTestDB with successful transaction database registration", func(t *testing.T) {

	})

	t.Run("Scenario 7: Test NewTestDB with unique database connection names", func(t *testing.T) {

	})

	t.Run("Scenario 8: Test NewTestDB with concurrent calls", func(t *testing.T) {

	})
}
