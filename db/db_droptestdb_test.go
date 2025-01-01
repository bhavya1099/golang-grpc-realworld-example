package db

import (
	"errors"
	"sync"
	"testing"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
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
func TestDropTestDB(t *testing.T) {
	var mutex sync.Mutex
	var txdbInitialized bool

	t.Run("DropTestDB_ClosesDBSuccessfully", func(t *testing.T) {
		t.Log("Scenario 1: DropTestDB_ClosesDBSuccessfully")
		db, _ := gorm.Open("sqlite3", "test.db")
		err := DropTestDB(db)
		assert.NoError(t, err)
	})

	t.Run("DropTestDB_NilDBInstance", func(t *testing.T) {
		t.Log("Scenario 2: DropTestDB_NilDBInstance")
		err := DropTestDB(nil)
		assert.NoError(t, err)
	})

	t.Run("DropTestDB_ErrorOnClose", func(t *testing.T) {
		t.Log("Scenario 3: DropTestDB_ErrorOnClose")
		mockDB := &gorm.DB{}
		mockDB.Close = func() error {
			return errors.New("error on close")
		}
		err := DropTestDB(mockDB)
		assert.Error(t, err)
	})

	t.Run("DropTestDB_MultipleCallsToClose", func(t *testing.T) {
		t.Log("Scenario 4: DropTestDB_MultipleCallsToClose")
		db, _ := gorm.Open("sqlite3", "test.db")
		DropTestDB(db)
		err := DropTestDB(db)
		assert.NoError(t, err)
	})

	t.Run("DropTestDB_ConcurrentCalls", func(t *testing.T) {
		t.Log("Scenario 5: DropTestDB_ConcurrentCalls")
		db, _ := gorm.Open("sqlite3", "test.db")
		var wg sync.WaitGroup
		wg.Add(2)
		go func() {
			defer wg.Done()
			err := DropTestDB(db)
			assert.NoError(t, err)
		}()
		go func() {
			defer wg.Done()
			err := DropTestDB(db)
			assert.NoError(t, err)
		}()
		wg.Wait()
	})
}
