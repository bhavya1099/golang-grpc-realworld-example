package store

import (
	"errors"
	"testing"
	"github.com/jinzhu/gorm"
	"github.com/raahii/golang-grpc-realworld-example/model"
	"github.com/stretchr/testify/assert"
	"your-package-path/store"
)

func TestUserStoreUnfollow(t *testing.T) {

	t.Run("Unfollow successful with valid users", func(t *testing.T) {

		mockDB := &gorm.DB{}
		userStore := store.UserStore{db: mockDB}
		userA := &model.User{}
		userB := &model.User{}

		err := userStore.Unfollow(userA, userB)

		assert.NoError(t, err, "Expected no error for successful unfollow")
		t.Log("Unfollow successful with valid users test passed")
	})

	t.Run("Unfollow failure with invalid user to unfollow", func(t *testing.T) {

		mockDB := &gorm.DB{}
		userStore := store.UserStore{db: mockDB}
		userA := &model.User{}
		invalidUser := &model.User{}

		err := userStore.Unfollow(userA, invalidUser)

		assert.Error(t, err, "Expected error for unfollowing invalid user")
		t.Log("Unfollow failure with invalid user test passed")
	})

	t.Run("Unfollow failure with self-unfollow attempt", func(t *testing.T) {

		mockDB := &gorm.DB{}
		userStore := store.UserStore{db: mockDB}
		userA := &model.User{}

		err := userStore.Unfollow(userA, userA)

		assert.Error(t, err, "Expected error for self-unfollow attempt")
		t.Log("Unfollow failure with self-unfollow attempt test passed")
	})

	t.Run("Unfollow failure with database error", func(t *testing.T) {

		mockDB := &gorm.DB{}
		mockDB.Error = errors.New("simulated database error")
		userStore := store.UserStore{db: mockDB}
		userA := &model.User{}
		userB := &model.User{}

		err := userStore.Unfollow(userA, userB)

		assert.Error(t, err, "Expected error for database error during unfollow")
		t.Log("Unfollow failure with database error test passed")
	})
}
