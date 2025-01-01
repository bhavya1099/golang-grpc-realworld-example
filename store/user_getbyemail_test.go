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

func (m *MockDB) First(out interface{}, where ...interface{}) *store.DB {
	args := m.Called(out, where)
	return args.Get(0).(*store.DB)
}

func (m *MockDB) Where(query interface{}, args ...interface{}) *store.DB {
	args := m.Called(query, args)
	return args.Get(0).(*store.DB)
}

func TestUserStoreGetByEmail(t *testing.T) {
	// Mock database instance
	mockDB := new(MockDB)

	// Initialize UserStore with mock DB
	userStore := &store.UserStore{db: mockDB}

	t.Run("GetByEmail_ValidEmail_ReturnsUser", func(t *testing.T) {
		// Define expected user
		expectedUser := &model.User{
			Username: "testuser",
			Email:    "test@example.com",
			Password: "password",
			Bio:      "Test bio",
			Image:    "image.png",
		}

		// Mock DB behavior for successful retrieval
		mockDB.On("First", &model.User{}, "email = ?", "test@example.com").Return(nil)

		// Call GetByEmail function
		user, err := userStore.GetByEmail("test@example.com")

		// Assertion
		assert.NoError(t, err)
		assert.Equal(t, expectedUser, user)
	})

	t.Run("GetByEmail_NonExistentEmail_ReturnsError", func(t *testing.T) {
		// Mock DB behavior for non-existent email
		mockDB.On("First", &model.User{}, "email = ?", "nonexistent@example.com").Return(errors.New("record not found"))

		// Call GetByEmail function
		user, err := userStore.GetByEmail("nonexistent@example.com")

		// Assertion
		assert.Error(t, err)
		assert.Nil(t, user)
	})

	t.Run("GetByEmail_EmptyEmail_ReturnsError", func(t *testing.T) {
		// Call GetByEmail with empty email
		user, err := userStore.GetByEmail("")

		// Assertion
		assert.Error(t, err)
		assert.Nil(t, user)
	})

	t.Run("GetByEmail_DBError_ReturnsError", func(t *testing.T) {
		// Mock DB behavior for database error
		mockDB.On("First", &model.User{}, "email = ?", "error@example.com").Return(errors.New("database error"))

		// Call GetByEmail function
		user, err := userStore.GetByEmail("error@example.com")

		// Assertion
		assert.Error(t, err)
		assert.Nil(t, user)
	})
}
