package store

import (
	"errors"
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/raahii/golang-grpc-realworld-example/model"
	"github.com/raahii/golang-grpc-realworld-example/store"
)

func TestUserStoreFollow(t *testing.T) {
	t.Parallel()

	t.Run("Append User B to User A's Follows List Successfully", func(t *testing.T) {
		userStore := &store.UserStore{}
		userA := &model.User{ID: 1}
		userB := &model.User{ID: 2}

		err := userStore.Follow(userA, userB)

		assert.NoError(t, err, "User B should be successfully appended to User A's Follows list")
		t.Log("User B successfully appended to User A's Follows list")
	})

	t.Run("Attempt to Append User B to User A's Follows List with Duplicate User B", func(t *testing.T) {
		userStore := &store.UserStore{}
		userA := &model.User{ID: 1}
		userB := &model.User{ID: 2}

		err := userStore.Follow(userA, userB)

		assert.Error(t, err, "Appending duplicate User B should result in an error")
		assert.True(t, errors.Is(err, store.ErrDuplicateEntry), "Error should indicate unique constraint violation")
		t.Log("Error correctly indicates unique constraint violation for duplicate User B")
	})

	t.Run("Attempt to Append User Itself to Its Follows List", func(t *testing.T) {
		userStore := &store.UserStore{}
		userA := &model.User{ID: 1}

		err := userStore.Follow(userA, userA)

		assert.Error(t, err, "Appending User itself should result in an error")
		assert.True(t, errors.Is(err, store.ErrSelfReference), "Error should indicate self-referencing relationship")
		t.Log("Error correctly indicates self-referencing relationship for User A")
	})

	t.Run("Error Handling - Database Connection Failure", func(t *testing.T) {
		userStore := &store.UserStore{}
		userA := &model.User{ID: 1}
		userB := &model.User{ID: 2}

		err := userStore.Follow(userA, userB)

		assert.Error(t, err, "Database connection failure should result in an error")
		assert.True(t, errors.Is(err, store.ErrDatabaseFailure), "Error should indicate database connection error")
		t.Log("Error correctly indicates database connection failure")
	})
}
