package store

import (
	"errors"
	"testing"
	"github.com/jinzhu/gorm"
	"github.com/raahii/golang-grpc-realworld-example/model"
	"github.com/stretchr/testify/assert"
)

func TestUserStoreIsFollowing(t *testing.T) {

	t.Run("Valid Follow Relationship Exists", func(t *testing.T) {
		userA := &model.User{ID: 1}
		userB := &model.User{ID: 2}
		userA.Follows = append(userA.Follows, *userB)
		userStore := &store.UserStore{db: &gorm.DB{}}

		isFollowing, err := userStore.IsFollowing(userA, userB)

		assert.NoError(t, err)
		assert.True(t, isFollowing)
		t.Log("Valid follow relationship exists test passed")
	})

	t.Run("No Follow Relationship Exists", func(t *testing.T) {
		userA := &model.User{ID: 1}
		userB := &model.User{ID: 2}
		userStore := &store.UserStore{db: &gorm.DB{}}

		isFollowing, err := userStore.IsFollowing(userA, userB)

		assert.NoError(t, err)
		assert.False(t, isFollowing)
		t.Log("No follow relationship exists test passed")
	})

	t.Run("One User is Nil", func(t *testing.T) {
		userA := &model.User{ID: 1}
		var userB *model.User
		userStore := &store.UserStore{db: &gorm.DB{}}

		isFollowing, err := userStore.IsFollowing(userA, userB)

		assert.NoError(t, err)
		assert.False(t, isFollowing)
		t.Log("One user is nil test passed")
	})

	t.Run("Both Users are Nil", func(t *testing.T) {
		var userA, userB *model.User
		userStore := &store.UserStore{db: &gorm.DB{}}

		isFollowing, err := userStore.IsFollowing(userA, userB)

		assert.NoError(t, err)
		assert.False(t, isFollowing)
		t.Log("Both users are nil test passed")
	})

	t.Run("Database Error Occurs", func(t *testing.T) {
		userA := &model.User{ID: 1}
		userB := &model.User{ID: 2}
		mockDB := &gorm.DB{}
		mockDB.Error = errors.New("database error")
		userStore := &store.UserStore{db: mockDB}

		isFollowing, err := userStore.IsFollowing(userA, userB)

		assert.Error(t, err)
		assert.False(t, isFollowing)
		t.Log("Database error occurs test passed")
	})
}
