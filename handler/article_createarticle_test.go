package handler

import (
	"context"
	"errors"
	"testing"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	pb "github.com/raahii/golang-grpc-realworld-example/proto"
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
func TestHandlerCreateArticle(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserStore := NewMockUserStore(ctrl)
	mockArticleStore := NewMockArticleStore(ctrl)

	handler := &Handler{
		us: mockUserStore,
		as: mockArticleStore,
	}

	t.Run("Create valid article successfully", func(t *testing.T) {
		validRequest := &pb.CreateAritcleRequest{
			Article: &pb.CreateAritcleRequest_Article{
				Title:       "Test Article",
				Description: "Test Description",
				Body:        "Test Body",
				TagList:     []string{"tag1", "tag2"},
			},
		}

		mockUserStore.EXPECT().GetByID(gomock.Any()).Return(&model.User{}, nil)
		mockArticleStore.EXPECT().Create(gomock.Any()).Return(nil)

		resp, err := handler.CreateArticle(context.Background(), validRequest)

		assert.NoError(t, err)
		assert.NotNil(t, resp)
	})

	t.Run("Unauthenticated user attempts to create an article", func(t *testing.T) {
		unauthRequest := &pb.CreateAritcleRequest{
			Article: &pb.CreateAritcleRequest_Article{
				Title:       "Test Article",
				Description: "Test Description",
				Body:        "Test Body",
				TagList:     []string{"tag1", "tag2"},
			},
		}

		mockUserStore.EXPECT().GetByID(gomock.Any()).Return(nil, errors.New("user not found"))

		_, err := handler.CreateArticle(context.Background(), unauthRequest)

		assert.Error(t, err)
		statusErr, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.Unauthenticated, statusErr.Code())
	})

}
