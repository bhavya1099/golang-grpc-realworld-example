package model

import (
	"errors"
	"testing"
	"github.com/golang/mock/gomock"
	"golang.org/x/crypto/bcrypt"
	"github.com/stretchr/testify/assert"
	"github.com/raahii/golang-grpc-realworld-example/model"
)

type Controller struct {
	// T should only be called within a generated mock. It is not intended to
	// be used in user code and may be changed in future versions. T is the
	// TestReporter passed in when creating the Controller via NewController.
	// If the TestReporter does not implement a TestHelper it will be wrapped
	// with a nopTestHelper.
	T             TestHelper
	mu            sync.Mutex
	expectedCalls *callSet
	finished      bool
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
func TestUserHashPassword(t *testing.T) {
	t.Parallel()

	tt := []struct {
		name     string
		user     *model.User
		mockHash func()
		errMsg   string
	}{
		{
			name: "HashPassword_ValidPassword",
			user: &model.User{Password: "securePassword"},
			mockHash: func() {
				bcrypt.GenerateFromPassword = func(password []byte, cost int) ([]byte, error) {
					return []byte("hashedPassword"), nil
				}
			},
			errMsg: "",
		},
		{
			name:     "HashPassword_EmptyPassword",
			user:     &model.User{Password: ""},
			mockHash: func() {},
			errMsg:   "password should not be empty",
		},
		{
			name: "HashPassword_ErrorGeneratingHash",
			user: &model.User{Password: "testPassword"},
			mockHash: func() {
				bcrypt.GenerateFromPassword = func(password []byte, cost int) ([]byte, error) {
					return nil, errors.New("error generating hash")
				}
			},
			errMsg: "error generating hash",
		},
		{
			name: "HashPassword_LongPassword",
			user: &model.User{Password: "veryLongSecurePassword"},
			mockHash: func() {
				bcrypt.GenerateFromPassword = func(password []byte, cost int) ([]byte, error) {
					return []byte("hashedPassword"), nil
				}
			},
			errMsg: "",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			tc.mockHash()
			err := tc.user.HashPassword()

			if tc.errMsg != "" {
				assert.NotNil(t, err, "Expected an error but got nil")
				assert.EqualError(t, err, tc.errMsg, "Error message mismatch")
			} else {
				assert.Nil(t, err, "Expected no error but got one")
				assert.NotEmpty(t, tc.user.Password, "Password should not be empty after hashing")
			}
		})
	}
}
