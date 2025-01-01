package store_test

import (
	"errors"
	"testing"

	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"

	"github.com/raahii/golang-grpc-realworld-example/model"
	"github.com/raahii/golang-grpc-realworld-example/store"
)

func TestUserStoreCreate(t *testing.T) {
	// Scenario 1: Test successful creation of a new user
	t.Run("Test successful creation of a new user", func(t *testing.T) {
		// Arrange
		mockDB := &gorm.DB{}
		userStore := store.UserStore{db: mockDB}
		user := &model.User{
			Username: "testuser",
			Email:    "testuser@example.com",
			Password: "password123",
			Bio:      "Test bio",
			Image:    "test.jpg",
		}

		// Act
		err := userStore.Create(user)

		// Assert
		assert.NoError(t, err, "Error should be nil for successful user creation")
	})

	// Scenario 2: Test creation failure due to duplicate username
	t.Run("Test creation failure due to duplicate username", func(t *testing.T) {
		// Arrange
		mockDB := &gorm.DB{}
		userStore := store.UserStore{db: mockDB}
		user := &model.User{
			Username: "existinguser",
			Email:    "newuser@example.com",
			Password: "password123",
			Bio:      "Test bio",
			Image:    "test.jpg",
		}

		// Act
		err := userStore.Create(user)

		// Assert
		assert.Error(t, err, "Error should be returned for duplicate username creation")
		assert.Contains(t, err.Error(), "unique constraint", "Error message should indicate unique constraint violation")
	})

	// Scenario 3: Test creation failure due to missing mandatory field
	t.Run("Test creation failure due to missing mandatory field", func(t *testing.T) {
		// Arrange
		mockDB := &gorm.DB{}
		userStore := store.UserStore{db: mockDB}
		user := &model.User{
			Username: "missingfielduser",
			Password: "password123",
			Bio:      "Test bio",
			Image:    "test.jpg",
		}

		// Act
		err := userStore.Create(user)

		// Assert
		assert.Error(t, err, "Error should be returned for missing mandatory field")
		assert.Contains(t, err.Error(), "missing field", "Error message should indicate missing field")
	})

	// Scenario 4: Test creation failure due to database error
	t.Run("Test creation failure due to database error", func(t *testing.T) {
		// Arrange
		mockDB := &gorm.DB{}
		mockDB.Error = errors.New("database error")
		userStore := store.UserStore{db: mockDB}
		user := &model.User{
			Username: "testuser",
			Email:    "testuser@example.com",
			Password: "password123",
			Bio:      "Test bio",
			Image:    "test.jpg",
		}

		// Act
		err := userStore.Create(user)

		// Assert
		assert.Error(t, err, "Error should be returned for database error")
		assert.EqualError(t, err, "database error", "Error message should indicate database error")
	})

	// Scenario 5: Test creation failure with nil User pointer
	t.Run("Test creation failure with nil User pointer", func(t *testing.T) {
		// Arrange
		mockDB := &gorm.DB{}
		userStore := store.UserStore{db: mockDB}

		// Act
		err := userStore.Create(nil)

		// Assert
		assert.Error(t, err, "Error should be returned for nil User pointer")
		assert.Contains(t, err.Error(), "nil pointer", "Error message should indicate nil pointer")
	})
}
