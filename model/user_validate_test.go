package model

import (
	"regexp"
	"testing"
	"github.com/stretchr/testify/assert"
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
func TestUserValidate(t *testing.T) {
	t.Parallel()

	t.Run("Scenario 1: Successful validation of a User with valid input", func(t *testing.T) {
		user := model.User{
			Username: "john_doe",
			Email:    "john.doe@example.com",
			Password: "password123",
		}

		err := user.Validate()

		assert.NoError(t, err, "Validation should pass for valid user data")
	})

	t.Run("Scenario 2: Validation failure due to missing Username", func(t *testing.T) {
		user := model.User{
			Email:    "jane.doe@example.com",
			Password: "password456",
		}

		err := user.Validate()

		assert.Error(t, err, "Validation should fail for missing Username")
		assert.Contains(t, err.Error(), "Username", "Error message should indicate missing Username")
	})

	t.Run("Scenario 3: Validation failure due to invalid Email format", func(t *testing.T) {
		user := model.User{
			Username: "alice",
			Email:    "invalid_email",
			Password: "password789",
		}

		err := user.Validate()

		assert.Error(t, err, "Validation should fail for invalid Email format")
		assert.Contains(t, err.Error(), "Email", "Error message should indicate invalid Email format")
	})

	t.Run("Scenario 4: Validation failure due to missing Password", func(t *testing.T) {
		user := model.User{
			Username: "bob",
			Email:    "bob@example.com",
		}

		err := user.Validate()

		assert.Error(t, err, "Validation should fail for missing Password")
		assert.Contains(t, err.Error(), "Password", "Error message should indicate missing Password")
	})

	t.Run("Scenario 5: Validation failure due to invalid characters in Username", func(t *testing.T) {
		user := model.User{
			Username: "invalid_username@",
			Email:    "invalid@example.com",
			Password: "password123",
		}

		err := user.Validate()

		assert.Error(t, err, "Validation should fail for invalid characters in Username")
		assert.Contains(t, err.Error(), "Username", "Error message should indicate invalid characters in Username")
	})

	t.Run("Scenario 6: Validation failure due to multiple validation errors", func(t *testing.T) {
		user := model.User{
			Email: "invalid_email",
		}

		err := user.Validate()

		assert.Error(t, err, "Validation should fail for multiple validation errors")
		assert.Contains(t, err.Error(), "Username", "Error message should indicate missing Username")
		assert.Contains(t, err.Error(), "Password", "Error message should indicate missing Password")
	})

	t.Run("Scenario 7: Validation success with all optional fields empty", func(t *testing.T) {
		user := model.User{
			Username: "eve",
			Email:    "eve@example.com",
			Password: "password123",
		}

		err := user.Validate()

		assert.NoError(t, err, "Validation should pass for empty optional fields")
	})
}
