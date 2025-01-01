package db

import (
	"errors"
	"fmt"
	"io/ioutil"
	"testing"
	"github.com/BurntSushi/toml"
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

func TestSeed(t *testing.T) {
	type User struct {
		ID   int
		Name string
		Age  int
	}

	t.Run("Successful Seed with Valid User Data", func(t *testing.T) {

		validUserData := `
			[[users]]
			ID = 1
			Name = "Alice"
			Age = 30
			[[users]]
			ID = 2
			Name = "Bob"
			Age = 25
		`

		dbMock := &gorm.DB{}

		err := Seed(dbMock, validUserData)

		assert.NoError(t, err, "Expected Seed function to complete without errors")

	})

	t.Run("Seed with Empty User Data", func(t *testing.T) {

		emptyUserData := ``

		dbMock := &gorm.DB{}

		err := Seed(dbMock, emptyUserData)

		assert.NoError(t, err, "Expected Seed function to return nil with empty user data")

	})

	t.Run("Seed with Invalid User Data", func(t *testing.T) {

		invalidUserData := `
			[[users]]
			ID: 1
			Name: "Alice"
			Age: 30
			[[users]]
			ID: 2
			Name: "Bob"
			Age: 25
		`

		dbMock := &gorm.DB{}

		err := Seed(dbMock, invalidUserData)

		assert.Error(t, err, "Expected Seed function to return an error with invalid user data")

	})

	t.Run("Seed with Database Creation Error", func(t *testing.T) {

		failingUserData := `
			[[users]]
			ID = 1
			Name = "Alice"
			Age = 30
			[[users]]
			ID = 2
			Name = "Bob"
			Age = 25
		`

		dbMock := &gorm.DB{
			Error: errors.New("Database creation error"),
		}

		err := Seed(dbMock, failingUserData)

		assert.Error(t, err, "Expected Seed function to return an error due to database creation failure")

	})

	t.Run("Seed with Missing Users TOML File", func(t *testing.T) {

		dbMock := &gorm.DB{}

		err := Seed(dbMock, "")

		assert.Error(t, err, "Expected Seed function to return an error with missing users.toml file")

	})
}
