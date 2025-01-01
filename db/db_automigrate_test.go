package db

import (
	"errors"
	"testing"
	"github.com/jinzhu/gorm"
	"github.com/raahii/golang-grpc-realworld-example/model"
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
func TestAutoMigrate(t *testing.T) {
	t.Run("AutoMigrate_SuccessfulMigration", func(t *testing.T) {

		db, _ := gorm.Open("sqlite3", ":memory:")

		err := AutoMigrate(db)

		if err != nil {
			t.Errorf("AutoMigrate failed with error: %v", err)
		}
	})

	t.Run("AutoMigrate_ErrorDuringMigration", func(t *testing.T) {

		db, _ := gorm.Open("sqlite3", ":memory:")

		errModel := struct{}{}

		err := db.AutoMigrate(&errModel).Error

		if err == nil {
			t.Errorf("Expected error during migration, but got nil")
		}
	})

	t.Run("AutoMigrate_NilDBInstance", func(t *testing.T) {

		var db *gorm.DB

		err := AutoMigrate(db)

		if err == nil {
			t.Error("Expected error with nil DB instance, but got nil")
		}
	})

	t.Run("AutoMigrate_EmptyModelList", func(t *testing.T) {

		db, _ := gorm.Open("sqlite3", ":memory:")

		err := AutoMigrate(db)

		if err != nil {
			t.Errorf("AutoMigrate should succeed with an empty model list, but got error: %v", err)
		}
	})

	t.Run("AutoMigrate_ConcurrentMigrations", func(t *testing.T) {

		db, _ := gorm.Open("sqlite3", ":memory:")
		numRoutines := 5
		done := make(chan struct{})

		for i := 0; i < numRoutines; i++ {
			go func() {

				err := AutoMigrate(db)

				if err != nil {
					t.Errorf("AutoMigrate failed during concurrent migration: %v", err)
				}
				done <- struct{}{}
			}()
		}

		for i := 0; i < numRoutines; i++ {
			<-done
		}
	})
}
