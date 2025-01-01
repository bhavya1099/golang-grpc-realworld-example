package handler

import (
	"context"
	"errors"
	"testing"
	"github.com/raahii/golang-grpc-realworld-example/handler"
	"github.com/raahii/golang-grpc-realworld-example/model"
	"github.com/raahii/golang-grpc-realworld-example/proto"
	"github.com/raahii/golang-grpc-realworld-example/store"
	"github.com/stretchr/testify/assert"
)



type User struct {
	gorm.Model
	Username         string    `gorm:"unique_index;not null"`
	Email            string    `gorm:"unique_index;not null"`
	Password         string    `gorm:"not null"`
	Bio              string    `gorm:"not null"`
	Image            string    `gorm:"not null"`
	Follows          []User    `gorm:"many2many:follows;jointable_foreignkey:from_user_id;association_jointable_foreignkey:to_user_id"`
	FavoriteArticles []Article `gorm:"many2many:favorite_articles;"`
}


type T struct {
	common
	isEnvSet bool
	context  *testContext // For running tests and subtests.
}



type User struct {
	gorm.Model
	Username         string    `gorm:"unique_index;not null"`
	Email            string    `gorm:"unique_index;not null"`
	Password         string    `gorm:"not null"`
	Bio              string    `gorm:"not null"`
	Image            string    `gorm:"not null"`
	Follows          []User    `gorm:"many2many:follows;jointable_foreignkey:from_user_id;association_jointable_foreignkey:to_user_id"`
	FavoriteArticles []Article `gorm:"many2many:favorite_articles;"`
}

type mockUserStore struct{}

type T struct {
	common
	isEnvSet bool
	context  *testContext // For running tests and subtests.
}

func (m *mockUserStore) GetByID(id uint) (*model.User, error) {
	if id == 1 {
		return &model.User{ID: 1, Username: "user1"}, nil
	}
	return nil, errors.New("user not found")
}
func (m *mockUserStore) GetByUsername(username string) (*model.User, error) {
	if username == "user2" {
		return &model.User{ID: 2, Username: "user2"}, nil
	}
	return nil, errors.New("user not found")
}
func (m *mockUserStore) IsFollowing(a *model.User, b *model.User) (bool, error) {
	if a.ID == 1 && b.ID == 2 {
		return true, nil
	}
	return false, errors.New("following status error")
}
func TestHandlerUnfollowUser(t *testing.T) {
	h := &handler.Handler{
		us: &mockUserStore{},
	}

	t.Run("UnfollowUser_SuccessfulUnfollow", func(t *testing.T) {
		req := &proto.UnfollowRequest{Username: "user2"}
		ctx := context.Background()

		profileResp, err := h.UnfollowUser(ctx, req)

		assert.NoError(t, err)
		assert.NotNil(t, profileResp)
		assert.Equal(t, "user2", profileResp.Profile.Username)
		assert.False(t, profileResp.Profile.Following)
	})

	t.Run("UnfollowUser_UnauthenticatedError", func(t *testing.T) {
		req := &proto.UnfollowRequest{Username: "user2"}
		ctx := context.Background()

		_, err := h.UnfollowUser(ctx, req)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "unauthenticated")
	})

	t.Run("UnfollowUser_UserNotFound", func(t *testing.T) {
		req := &proto.UnfollowRequest{Username: "user3"}
		ctx := context.Background()

		_, err := h.UnfollowUser(ctx, req)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "user was not found")
	})

	t.Run("UnfollowUser_SelfUnfollowError", func(t *testing.T) {
		req := &proto.UnfollowRequest{Username: "user1"}
		ctx := context.Background()

		_, err := h.UnfollowUser(ctx, req)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "cannot follow yourself")
	})

	t.Run("UnfollowUser_InternalServerError", func(t *testing.T) {
		req := &proto.UnfollowRequest{Username: "user2"}
		ctx := context.Background()

		_, err := h.UnfollowUser(ctx, req)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "internal server error")
	})

	t.Run("UnfollowUser_FollowingStatusError", func(t *testing.T) {
		req := &proto.UnfollowRequest{Username: "user2"}
		ctx := context.Background()

		_, err := h.UnfollowUser(ctx, req)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "following status error")
	})

	t.Run("UnfollowUser_FailedUnfollowError", func(t *testing.T) {
		req := &proto.UnfollowRequest{Username: "user2"}
		ctx := context.Background()

		_, err := h.UnfollowUser(ctx, req)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to unfollow")
	})

	t.Run("UnfollowUser_NotFollowingError", func(t *testing.T) {
		req := &proto.UnfollowRequest{Username: "user2"}
		ctx := context.Background()

		_, err := h.UnfollowUser(ctx, req)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "you are not following the user")
	})
}
func (m *mockUserStore) Unfollow(a *model.User, b *model.User) error {
	if a.ID == 1 && b.ID == 2 {
		return nil
	}
	return errors.New("failed to unfollow")
}
