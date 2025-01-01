package store_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/raahii/golang-grpc-realworld-example/model"
	"github.com/raahii/golang-grpc-realworld-example/store"
)

func TestUserStoreIsFollowing(t *testing.T) {
	// Scenario 1
	t.Run("Test when both input users are not nil and a follows b", func(t *testing.T) {
		// Arrange
		userA := &model.User{ID: 1}
		userB := &model.User{ID: 2}
		mockDB := &store.UserStore{} // TODO: Mock the database with follow relationship data between userA and userB
		
		// Act
		result, err := mockDB.IsFollowing(userA, userB)
		
		// Assert
		assert.NoError(t, err, "Error should be nil")
		assert.True(t, result, "User A should follow User B")
		t.Log("Scenario 1 executed successfully")
	})

	// Scenario 2
	t.Run("Test when both input users are not nil and a does not follow b", func(t *testing.T) {
		// Arrange
		userA := &model.User{ID: 1}
		userB := &model.User{ID: 2}
		mockDB := &store.UserStore{} // TODO: Mock the database without follow relationship between userA and userB
		
		// Act
		result, err := mockDB.IsFollowing(userA, userB)
		
		// Assert
		assert.NoError(t, err, "Error should be nil")
		assert.False(t, result, "User A should not follow User B")
		t.Log("Scenario 2 executed successfully")
	})

	// Scenario 3
	t.Run("Test when user a is nil", func(t *testing.T) {
		// Arrange
		var userA *model.User
		userB := &model.User{ID: 2}
		mockDB := &store.UserStore{}
		
		// Act
		result, err := mockDB.IsFollowing(userA, userB)
		
		// Assert
		assert.NoError(t, err, "Error should be nil")
		assert.False(t, result, "User A is nil, should not follow User B")
		t.Log("Scenario 3 executed successfully")
	})

	// Scenario 4
	t.Run("Test when user b is nil", func(t *testing.T) {
		// Arrange
		userA := &model.User{ID: 1}
		var userB *model.User
		mockDB := &store.UserStore{}
		
		// Act
		result, err := mockDB.IsFollowing(userA, userB)
		
		// Assert
		assert.NoError(t, err, "Error should be nil")
		assert.False(t, result, "User B is nil, should not be followed by User A")
		t.Log("Scenario 4 executed successfully")
	})

	// Scenario 5
	t.Run("Test when an error occurs during database query", func(t *testing.T) {
		// Arrange
		userA := &model.User{ID: 1}
		userB := &model.User{ID: 2}
		mockDB := &store.UserStore{} // TODO: Mock the database to simulate an error during query execution
		
		// Act
		result, err := mockDB.IsFollowing(userA, userB)
		
		// Assert
		assert.Error(t, err, "Error should occur during database query")
		assert.False(t, result, "Function should return false")
		t.Log("Scenario 5 executed successfully")
	})
}
