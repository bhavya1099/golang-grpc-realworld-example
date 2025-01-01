package handler

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"github.com/raahii/golang-grpc-realworld-example/model"
	"github.com/raahii/golang-grpc-realworld-example/proto"
	"github.com/stretchr/testify/assert"
)


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
func TestHandlerCreateUser(t *testing.T) {
	t.Parallel()

	t.Run("CreateUser_ValidInput", func(t *testing.T) {
		t.Log("Scenario 1: CreateUser_ValidInput")
		req := &proto.CreateUserRequest{
			User: &proto.CreateUserRequest_User{
				Email:    "test@example.com",
				Username: "testuser",
				Password: "password123",
			},
		}

		handler := &Handler{
			logger: nil,
			us:     nil,
			as:     nil,
		}

		resp, err := handler.CreateUser(context.Background(), req)

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.NotNil(t, resp.User)
		assert.Equal(t, req.User.GetEmail(), resp.User.Email)
		assert.NotEmpty(t, resp.User.Token)
	})

	t.Run("CreateUser_InvalidUserValidation", func(t *testing.T) {
		t.Log("Scenario 2: CreateUser_InvalidUserValidation")
		req := &proto.CreateUserRequest{
			User: &proto.CreateUserRequest_User{
				Email:    "invalidemail",
				Username: "testuser",
				Password: "password123",
			},
		}

		handler := &Handler{
			logger: nil,
			us:     nil,
			as:     nil,
		}

		resp, err := handler.CreateUser(context.Background(), req)

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.True(t, errors.Is(err, errors.New("validation error")))
	})

	t.Run("CreateUser_FailedToHashPassword", func(t *testing.T) {
		t.Log("Scenario 3: CreateUser_FailedToHashPassword")
		req := &proto.CreateUserRequest{
			User: &proto.CreateUserRequest_User{
				Email:    "test@example.com",
				Username: "testuser",
				Password: "",
			},
		}

		handler := &Handler{
			logger: nil,
			us:     nil,
			as:     nil,
		}

		resp, err := handler.CreateUser(context.Background(), req)

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.True(t, errors.Is(err, errors.New("Failed to hash password")))
	})

	t.Run("CreateUser_FailedToCreateUser", func(t *testing.T) {
		t.Log("Scenario 4: CreateUser_FailedToCreateUser")
		req := &proto.CreateUserRequest{
			User: &proto.CreateUserRequest_User{
				Email:    "test@example.com",
				Username: "testuser",
				Password: "password123",
			},
		}

		handler := &Handler{
			logger: nil,
			us:     nil,
			as:     nil,
		}

		resp, err := handler.CreateUser(context.Background(), req)

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.True(t, errors.Is(err, errors.New("Failed to create user")))
	})

	t.Run("CreateUser_FailedToGenerateToken", func(t *testing.T) {
		t.Log("Scenario 5: CreateUser_FailedToGenerateToken")
		req := &proto.CreateUserRequest{
			User: &proto.CreateUserRequest_User{
				Email:    "test@example.com",
				Username: "testuser",
				Password: "password123",
			},
		}

		handler := &Handler{
			logger: nil,
			us:     nil,
			as:     nil,
		}

		resp, err := handler.CreateUser(context.Background(), req)

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.True(t, errors.Is(err, errors.New("Failed to create token")))
	})
}
