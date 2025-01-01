package handler

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"github.com/golang/mock/gomock"
	"github.com/raahii/golang-grpc-realworld-example/model"
	pb "github.com/raahii/golang-grpc-realworld-example/proto"
	"github.com/raahii/golang-grpc-realworld-example/store"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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


type T struct {
	common
	isEnvSet bool
	context  *testContext // For running tests and subtests.
}
func TestHandlerFollowUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserStore := store.NewMockUserStore(ctrl)

	h := &Handler{
		us: mockUserStore,
	}

	ctx := context.Background()

	t.Run("FollowUser_SuccessfulFollow", func(t *testing.T) {
		expectedUser := &model.User{
			ID:       1,
			Username: "user1",
		}
		mockUserStore.EXPECT().GetByID(gomock.Any()).Return(expectedUser, nil)
		mockUserStore.EXPECT().GetByUsername("user2").Return(&model.User{ID: 2, Username: "user2"}, nil)
		mockUserStore.EXPECT().Follow(expectedUser, &model.User{ID: 2, Username: "user2"}).Return(nil)

		req := &pb.FollowRequest{Username: "user2"}
		profileResp, err := h.FollowUser(ctx, req)

		assert.NoError(t, err)
		assert.NotNil(t, profileResp)
		assert.Equal(t, "user2", profileResp.Profile.Username)
		assert.True(t, profileResp.Profile.Following)
	})

	t.Run("FollowUser_UnauthenticatedUser", func(t *testing.T) {
		mockUserStore.EXPECT().GetByID(gomock.Any()).Return(nil, errors.New("unauthenticated"))

		req := &pb.FollowRequest{Username: "user2"}
		_, err := h.FollowUser(ctx, req)

		assert.Error(t, err)
		st, _ := status.FromError(err)
		assert.Equal(t, codes.Unauthenticated, st.Code())
	})

	t.Run("FollowUser_SelfFollowAttempt", func(t *testing.T) {
		expectedUser := &model.User{
			ID:       1,
			Username: "user1",
		}
		mockUserStore.EXPECT().GetByID(gomock.Any()).Return(expectedUser, nil)

		req := &pb.FollowRequest{Username: "user1"}
		_, err := h.FollowUser(ctx, req)

		assert.Error(t, err)
		st, _ := status.FromError(err)
		assert.Equal(t, codes.InvalidArgument, st.Code())
	})

	t.Run("FollowUser_UserNotFound", func(t *testing.T) {
		expectedUser := &model.User{
			ID:       1,
			Username: "user1",
		}
		mockUserStore.EXPECT().GetByID(gomock.Any()).Return(expectedUser, nil)
		mockUserStore.EXPECT().GetByUsername("unknown_user").Return(nil, errors.New("user not found"))

		req := &pb.FollowRequest{Username: "unknown_user"}
		_, err := h.FollowUser(ctx, req)

		assert.Error(t, err)
		st, _ := status.FromError(err)
		assert.Equal(t, codes.NotFound, st.Code())
	})

	t.Run("FollowUser_FollowingError", func(t *testing.T) {
		expectedUser := &model.User{
			ID:       1,
			Username: "user1",
		}
		mockUserStore.EXPECT().GetByID(gomock.Any()).Return(expectedUser, nil)
		mockUserStore.EXPECT().GetByUsername("user2").Return(&model.User{ID: 2, Username: "user2"}, nil)
		mockUserStore.EXPECT().Follow(expectedUser, &model.User{ID: 2, Username: "user2"}).Return(errors.New("follow error"))

		req := &pb.FollowRequest{Username: "user2"}
		_, err := h.FollowUser(ctx, req)

		assert.Error(t, err)
		st, _ := status.FromError(err)
		assert.Equal(t, codes.Aborted, st.Code())
	})

	t.Run("FollowUser_CurrentUserNotFound", func(t *testing.T) {
		mockUserStore.EXPECT().GetByID(gomock.Any()).Return(nil, errors.New("user not found"))

		req := &pb.FollowRequest{Username: "user2"}
		_, err := h.FollowUser(ctx, req)

		assert.Error(t, err)
		st, _ := status.FromError(err)
		assert.Equal(t, codes.NotFound, st.Code())
	})

	t.Run("FollowUser_InvalidArgument", func(t *testing.T) {
		expectedUser := &model.User{
			ID:       1,
			Username: "user1",
		}
		mockUserStore.EXPECT().GetByID(gomock.Any()).Return(expectedUser, nil)

		req := &pb.FollowRequest{Username: "invalid_username"}
		_, err := h.FollowUser(ctx, req)

		assert.Error(t, err)
		st, _ := status.FromError(err)
		assert.Equal(t, codes.InvalidArgument, st.Code())
	})
}
