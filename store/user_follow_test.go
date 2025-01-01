package store_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/raahii/golang-grpc-realworld-example/model"
	"github.com/raahii/golang-grpc-realworld-example/store"
)

func TestUserStoreFollow(t *testing.T) {
	t.Parallel()

	// Scenario 1: Follow Successful
	t.Run("Follow Successful", func(t *testing.T) {
		t.Log("Testing successful follow operation")
		// Arrange
		a := &model.User{}
		b := &model.User{}
		mockDB := &gorm.DB{} // Mocked gorm.DB
		userStore := store.UserStore{db: mockDB}

		// Act
		err := userStore.Follow(a, b)

		// Assert
		assert.NoError(t, err, "Expected no error during follow operation")
		// Validation: Check if user 'a' follows user 'b'
		assert.Contains(t, a.Follows, *b, "User 'a' should follow user 'b'")
	})

	// Scenario 2: Follow Same User
	t.Run("Follow Same User", func(t *testing.T) {
		t.Log("Testing follow operation with the same user")
		// Arrange
		a := &model.User{}
		mockDB := &gorm.DB{} // Mocked gorm.DB
		userStore := store.UserStore{db: mockDB}

		// Act
		err := userStore.Follow(a, a)

		// Assert
		assert.NoError(t, err, "Expected no error when trying to follow the same user")
		// Validation: Check if user 'a' is not duplicated in follows list
		assert.Equal(t, 1, len(a.Follows), "User 'a' should not be duplicated in follows list")
	})

	// Scenario 3: Follow Error
	t.Run("Follow Error", func(t *testing.T) {
		t.Log("Testing error handling during follow operation")
		// Arrange
		a := &model.User{}
		b := &model.User{}
		mockError := errors.New("mock error")
		mockDB := &gorm.DB{Error: mockError} // Mocked gorm.DB with error
		userStore := store.UserStore{db: mockDB}

		// Act
		err := userStore.Follow(a, b)

		// Assert
		assert.Error(t, err, "Expected error during follow operation")
		// Validation: Check if the error returned matches the mock error
		assert.EqualError(t, err, "mock error", "Error should match the expected mock error")
	})
}
