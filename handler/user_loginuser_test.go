package handler

import (
	"context"
	"errors"
	"testing"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/your/module/handler"
	"github.com/your/module/model"
	pb "github.com/your/module/proto"
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
func TestHandlerLoginUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLogger := NewMockLogger(ctrl)
	mockUserStore := NewMockUserStore(ctrl)

	h := handler.Handler{
		logger: mockLogger,
		us:     mockUserStore,
	}

	t.Run("Successful User Login", func(t *testing.T) {
		email := "user@example.com"
		password := "password123"
		user := &model.User{
			Email:    email,
			Password: "hashedpassword123",
			Username: "user123",
			Bio:      "bio",
			Image:    "image.jpg",
		}

		req := &pb.LoginUserRequest{
			User: &pb.LoginUserRequest_User{
				Email:    email,
				Password: password,
			},
		}

		mockUserStore.EXPECT().GetByEmail(email).Return(user, nil)

		mockLogger.EXPECT().Info().Times(1)
		mockLogger.EXPECT().Error().Times(0)

		res, err := h.LoginUser(context.Background(), req)

		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, user.Email, res.GetUser().GetEmail())
		assert.NotEmpty(t, res.GetUser().GetToken())
	})

	t.Run("Invalid Email Login", func(t *testing.T) {
		email := "invalid@example.com"
		password := "password123"

		req := &pb.LoginUserRequest{
			User: &pb.LoginUserRequest_User{
				Email:    email,
				Password: password,
			},
		}

		mockUserStore.EXPECT().GetByEmail(email).Return(nil, errors.New("user not found"))

		mockLogger.EXPECT().Info().Times(1)
		mockLogger.EXPECT().Error().Times(1)

		_, err := h.LoginUser(context.Background(), req)

		assert.Error(t, err)
		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.InvalidArgument, st.Code())
	})

	t.Run("Incorrect Password Login", func(t *testing.T) {
		email := "user@example.com"
		password := "wrongpassword"

		req := &pb.LoginUserRequest{
			User: &pb.LoginUserRequest_User{
				Email:    email,
				Password: password,
			},
		}

		user := &model.User{
			Email:    email,
			Password: "hashedpassword123",
			Username: "user123",
			Bio:      "bio",
			Image:    "image.jpg",
		}

		mockUserStore.EXPECT().GetByEmail(email).Return(user, nil)

		mockLogger.EXPECT().Info().Times(1)
		mockLogger.EXPECT().Error().Times(1)

		_, err := h.LoginUser(context.Background(), req)

		assert.Error(t, err)
		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.InvalidArgument, st.Code())
	})

	t.Run("Internal Server Error", func(t *testing.T) {
		email := "user@example.com"
		password := "password123"
		user := &model.User{
			Email:    email,
			Password: "hashedpassword123",
			Username: "user123",
			Bio:      "bio",
			Image:    "image.jpg",
		}

		req := &pb.LoginUserRequest{
			User: &pb.LoginUserRequest_User{
				Email:    email,
				Password: password,
			},
		}

		mockUserStore.EXPECT().GetByEmail(email).Return(user, nil)
		mockLogger.EXPECT().Info().Times(1)
		mockLogger.EXPECT().Error().Times(1)

		mockUserStore.EXPECT().GetByEmail(email).Return(user, nil)
		mockUserStore.EXPECT().GenerateToken(user.ID).Return("", errors.New("token generation error"))

		_, err := h.LoginUser(context.Background(), req)

		assert.Error(t, err)
		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.Aborted, st.Code())
	})
}
