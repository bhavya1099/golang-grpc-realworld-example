package store_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/raahii/golang-grpc-realworld-example/model"
	"github.com/raahii/golang-grpc-realworld-example/store"
)

type MockDB struct {
	mock.Mock
}

func (m *MockDB) Model(value interface{}) *store.DB {
	args := m.Called(value)
	return args.Get(0).(*store.DB)
}

func (m *MockDB) Update(attrs ...interface{}) *store.DB {
	args := m.Called(attrs...)
	return args.Get(0).(*store.DB)
}

func TestUserStoreUpdate(t *testing.T) {
	mockDB := &MockDB{}

	// Scenario 1: Update User Successfully
	t.Run("Update User Successfully", func(t *testing.T) {
		user := &model.User{
			Username: "john_doe",
			Email:    "john.doe@example.com",
			Password: "password123",
			Bio:      "Test bio",
			Image:    "profile.jpg",
		}

		mockDB.On("Model", user).Return(mockDB)
		mockDB.On("Update", user).Return(nil)

		userStore := &store.UserStore{db: mockDB}
		err := userStore.Update(user)

		assert.NoError(t, err, "Update User should not return an error")
	})

	// Scenario 2: Update User with Empty User Object
	t.Run("Update User with Empty User Object", func(t *testing.T) {
		emptyUser := &model.User{}

		userStore := &store.UserStore{db: mockDB}
		err := userStore.Update(emptyUser)

		assert.Error(t, err, "Update User with empty object should return an error")
	})

	// Scenario 3: Update User with Invalid Data
	t.Run("Update User with Invalid Data", func(t *testing.T) {
		invalidUser := &model.User{
			Username: "", // Invalid: empty username
		}

		userStore := &store.UserStore{db: mockDB}
		err := userStore.Update(invalidUser)

		assert.Error(t, err, "Update User with invalid data should return an error")
	})

	// Scenario 4: Update User - Database Error
	t.Run("Update User - Database Error", func(t *testing.T) {
		user := &model.User{
			Username: "jane_doe",
			Email:    "jane.doe@example.com",
			Password: "password456",
			Bio:      "Test bio 2",
			Image:    "profile2.jpg",
		}

		mockDB.On("Model", user).Return(mockDB)
		mockDB.On("Update", user).Return(errors.New("database error"))

		userStore := &store.UserStore{db: mockDB}
		err := userStore.Update(user)

		assert.Error(t, err, "Update User should return an error on database error")
	})

	// Scenario 5: Update User - Existing Username Conflict
	t.Run("Update User - Existing Username Conflict", func(t *testing.T) {
		existingUser := &model.User{
			Username: "john_doe", // Conflicts with existing user
		}

		mockDB.On("Model", existingUser).Return(mockDB)
		mockDB.On("Update", existingUser).Return(errors.New("username conflict"))

		userStore := &store.UserStore{db: mockDB}
		err := userStore.Update(existingUser)

		assert.Error(t, err, "Update User should return an error on username conflict")
	})

	mockDB.AssertExpectations(t)
}
