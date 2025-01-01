package store_test

import (
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"

	"github.com/raahii/golang-grpc-realworld-example/model"
	"github.com/raahii/golang-grpc-realworld-example/store"
)

func TestUserStoreUnfollow(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := gorm.NewMockDB(ctrl)

	userA := &model.User{ID: 1}
	userB := &model.User{ID: 2}

	userStore := store.UserStore{db: mockDB}

	t.Run("Successfully Unfollow User", func(t *testing.T) {
		mockDB.EXPECT().Model(userA).Return(mockDB)
		mockDB.EXPECT().Association("Follows").Return(mockDB)
		mockDB.EXPECT().Delete(userB).Return(nil)

		err := userStore.Unfollow(userA, userB)

		assert.NoError(t, err, "Expected no error when unfollowing user")
	})

	t.Run("Unfollow Non-existent User", func(t *testing.T) {
		mockDB.EXPECT().Model(userA).Return(mockDB)
		mockDB.EXPECT().Association("Follows").Return(mockDB)
		mockDB.EXPECT().Delete(userB).Return(errors.New("record not found"))

		err := userStore.Unfollow(userA, userB)

		assert.Error(t, err, "Expected error when unfollowing non-existent user")
		assert.EqualError(t, err, "record not found", "Expected 'record not found' error message")
	})

	t.Run("Unfollow Self", func(t *testing.T) {
		mockDB.EXPECT().Model(userA).Return(mockDB)
		mockDB.EXPECT().Association("Follows").Return(mockDB)

		err := userStore.Unfollow(userA, userA)

		assert.Error(t, err, "Expected error when trying to unfollow self")
		assert.EqualError(t, err, "cannot unfollow self", "Expected 'cannot unfollow self' error message")
	})

	t.Run("Unfollow Error Handling", func(t *testing.T) {
		mockDB.EXPECT().Model(userA).Return(mockDB)
		mockDB.EXPECT().Association("Follows").Return(mockDB)
		mockDB.EXPECT().Delete(userB).Return(errors.New("database error"))

		err := userStore.Unfollow(userA, userB)

		assert.Error(t, err, "Expected error when unfollowing user with database error")
		assert.EqualError(t, err, "database error", "Expected 'database error' error message")
	})
}
