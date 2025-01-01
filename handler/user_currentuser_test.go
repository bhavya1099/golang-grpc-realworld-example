package handler

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"github.com/golang/mock/gomock"
	"github.com/raahii/golang-grpc-realworld-example/auth"
	"github.com/raahii/golang-grpc-realworld-example/handler"
	"github.com/raahii/golang-grpc-realworld-example/model"
	pb "github.com/raahii/golang-grpc-realworld-example/proto"
	"github.com/raahii/golang-grpc-realworld-example/store"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Controller struct {
	// T should only be called within a generated mock. It is not intended to
	// be used in user code and may be changed in future versions. T is the
	// TestReporter passed in when creating the Controller via NewController.
	// If the TestReporter does not implement a TestHelper it will be wrapped
	// with a nopTestHelper.
	T             TestHelper
	mu            sync.Mutex
	expectedCalls *callSet
	finished      bool
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
func TestHandlerCurrentUser(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockLogger := NewMockEvent(mockCtrl)

	mockUserStore := store.NewMockUserStore(mockCtrl)

	h := &handler.Handler{
		logger: mockLogger,
		us:     mockUserStore,
	}

	t.Run("Successful CurrentUser Request", func(t *testing.T) {
		expectedUser := &model.User{
			ID:       1,
			Email:    "test@example.com",
			Username: "testuser",
			Bio:      "Test Bio",
			Image:    "test.jpg",
		}
		mockUserStore.EXPECT().GetByID(expectedUser.ID).Return(expectedUser, nil)

		ctx := context.Background()
		req := &pb.Empty{}

		userResponse, err := h.CurrentUser(ctx, req)

		assert.NoError(t, err)
		assert.NotNil(t, userResponse)
		assert.Equal(t, expectedUser.Email, userResponse.User.Email)
		assert.NotEmpty(t, userResponse.User.Token)
	})

	t.Run("Unauthenticated CurrentUser Request", func(t *testing.T) {
		mockUserStore.EXPECT().GetByID(gomock.Any()).Times(0)

		ctx := context.Background()
		req := &pb.Empty{}

		authError := errors.New("unauthenticated")
		mockLogger.EXPECT().Error().Err(authError).AnyTimes()

		userResponse, err := h.CurrentUser(ctx, req)

		assert.Error(t, err)
		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.Unauthenticated, st.Code())
	})

	t.Run("User Not Found", func(t *testing.T) {
		expectedUserID := uint(1)
		mockUserStore.EXPECT().GetByID(expectedUserID).Return(nil, errors.New("user not found"))

		ctx := context.Background()
		req := &pb.Empty{}

		notFoundError := errors.New("user not found")
		mockLogger.EXPECT().Error().Err(gomock.Any()).AnyTimes()

		userResponse, err := h.CurrentUser(ctx, req)

		assert.Error(t, err)
		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.NotFound, st.Code())
	})

	t.Run("Token Generation Failure", func(t *testing.T) {
		expectedUser := &model.User{
			ID:       1,
			Email:    "test@example.com",
			Username: "testuser",
			Bio:      "Test Bio",
			Image:    "test.jpg",
		}
		mockUserStore.EXPECT().GetByID(expectedUser.ID).Return(expectedUser, nil)

		ctx := context.Background()
		req := &pb.Empty{}

		tokenGenError := errors.New("token generation error")
		mockLogger.EXPECT().Error().Err(tokenGenError).AnyTimes()

		userResponse, err := h.CurrentUser(ctx, req)

		assert.Error(t, err)
		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.Aborted, st.Code())
	})

	t.Run("Error Getting User Information", func(t *testing.T) {
		expectedUserID := uint(1)
		mockUserStore.EXPECT().GetByID(expectedUserID).Return(nil, errors.New("error getting user info"))

		ctx := context.Background()
		req := &pb.Empty{}

		mockLogger.EXPECT().Error().Err(gomock.Any()).AnyTimes()

		userResponse, err := h.CurrentUser(ctx, req)

		assert.Error(t, err)
		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.NotFound, st.Code())
	})

	t.Run("Empty Request", func(t *testing.T) {
		mockUserStore.EXPECT().GetByID(gomock.Any()).Times(0)

		ctx := context.Background()
		req := &pb.Empty{}

		authError := errors.New("unauthenticated")
		mockLogger.EXPECT().Error().Err(authError).AnyTimes()

		userResponse, err := h.CurrentUser(ctx, req)

		assert.Error(t, err)
		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.Unauthenticated, st.Code())
	})
}
