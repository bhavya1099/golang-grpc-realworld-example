package store_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/raahii/golang-grpc-realworld-example/model"
	"github.com/raahii/golang-grpc-realworld-example/store"
)

func TestUserStoreGetFollowingUserIDs(t *testing.T) {
	t.Parallel()

	// Scenario 1: GetFollowingUserIDs_Success
	t.Run("GetFollowingUserIDs_Success", func(t *testing.T) {
		t.Log("Scenario: GetFollowingUserIDs_Success")

		// Arrange
		mockUser := &model.User{ID: 1} // Mock user with ID 1
		userStore := store.UserStore{}  // Initialize UserStore

		// Act
		userIDs, err := userStore.GetFollowingUserIDs(mockUser)

		// Assert
		require.NoError(t, err, "Error not expected in GetFollowingUserIDs_Success")
		assert.Equal(t, []uint{2, 3, 4}, userIDs, "Returned user IDs do not match expected IDs")
	})

	// Scenario 2: GetFollowingUserIDs_NoRows
	t.Run("GetFollowingUserIDs_NoRows", func(t *testing.T) {
		t.Log("Scenario: GetFollowingUserIDs_NoRows")

		// Arrange
		mockUser := &model.User{ID: 5} // Mock user with ID 5 (no following users)
		userStore := store.UserStore{}  // Initialize UserStore

		// Act
		userIDs, err := userStore.GetFollowingUserIDs(mockUser)

		// Assert
		require.NoError(t, err, "Error not expected in GetFollowingUserIDs_NoRows")
		assert.Empty(t, userIDs, "Returned user IDs should be empty")
	})

	// Scenario 3: GetFollowingUserIDs_DBError
	t.Run("GetFollowingUserIDs_DBError", func(t *testing.T) {
		t.Log("Scenario: GetFollowingUserIDs_DBError")

		// Arrange
		mockUser := &model.User{ID: 7} // Mock user with ID 7
		userStore := store.UserStore{}  // Initialize UserStore with DB error setup

		// Act
		userIDs, err := userStore.GetFollowingUserIDs(mockUser)

		// Assert
		require.Error(t, err, "Error expected in GetFollowingUserIDs_DBError")
		assert.Empty(t, userIDs, "Returned user IDs should be empty")
	})

	// Scenario 4: GetFollowingUserIDs_EmptyUser
	t.Run("GetFollowingUserIDs_EmptyUser", func(t *testing.T) {
		t.Log("Scenario: GetFollowingUserIDs_EmptyUser")

		// Arrange
		mockUser := &model.User{}      // Empty mock user
		userStore := store.UserStore{} // Initialize UserStore

		// Act
		_, err := userStore.GetFollowingUserIDs(mockUser)

		// Assert
		require.Error(t, err, "Error expected in GetFollowingUserIDs_EmptyUser")
	})

	// Scenario 5: GetFollowingUserIDs_MultipleRows
	t.Run("GetFollowingUserIDs_MultipleRows", func(t *testing.T) {
		t.Log("Scenario: GetFollowingUserIDs_MultipleRows")

		// Arrange
		mockUser := &model.User{ID: 10} // Mock user with ID 10 (follows multiple users)
		userStore := store.UserStore{}   // Initialize UserStore

		// Act
		userIDs, err := userStore.GetFollowingUserIDs(mockUser)

		// Assert
		require.NoError(t, err, "Error not expected in GetFollowingUserIDs_MultipleRows")
		assert.ElementsMatch(t, []uint{11, 12, 13}, userIDs, "Returned user IDs do not match expected IDs")
	})
}
