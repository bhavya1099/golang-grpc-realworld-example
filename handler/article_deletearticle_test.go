package handler

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/raahii/golang-grpc-realworld-example/auth"
	"github.com/raahii/golang-grpc-realworld-example/handler"
	"github.com/raahii/golang-grpc-realworld-example/model"
	pb "github.com/raahii/golang-grpc-realworld-example/proto"
	"github.com/raahii/golang-grpc-realworld-example/store"
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
func TestHandlerDeleteArticle(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	logger := &zerolog.Logger{}

	mockUserStore := store.NewMockUserStore(mockCtrl)
	mockUserStore.EXPECT().GetByID(gomock.Any()).Return(&model.User{ID: 1}, nil).AnyTimes()

	mockArticleStore := store.NewMockArticleStore(mockCtrl)
	mockArticleStore.EXPECT().GetByID(gomock.Any()).Return(&model.Article{ID: 1, Author: &model.User{ID: 1}}, nil).AnyTimes()
	mockArticleStore.EXPECT().Delete(gomock.Any()).Return(nil).AnyTimes()

	h := &handler.Handler{
		logger: logger,
		us:     mockUserStore,
		as:     mockArticleStore,
	}

	tt := []struct {
		name          string
		req           *pb.DeleteArticleRequest
		expectedError error
	}{
		{
			name: "Successful deletion by author",
			req:  &pb.DeleteArticleRequest{Slug: "1"},
		},
		{
			name:          "Attempt to delete with invalid slug format",
			req:           &pb.DeleteArticleRequest{Slug: "invalid"},
			expectedError: status.Error(codes.InvalidArgument, "invalid article id"),
		},
		{
			name:          "Deleting non-existent article",
			req:           &pb.DeleteArticleRequest{Slug: "999"},
			expectedError: status.Error(codes.InvalidArgument, "invalid article id"),
		},
		{
			name:          "Unauthorized user deletion attempt",
			req:           &pb.DeleteArticleRequest{Slug: "1"},
			expectedError: status.Error(codes.Unauthenticated, "forbidden"),
		},
		{
			name:          "Error during article deletion process",
			req:           &pb.DeleteArticleRequest{Slug: "1"},
			expectedError: status.Error(codes.Unauthenticated, "failed to delete article"),
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			_, err := h.DeleteArticle(ctx, tc.req)

			if tc.expectedError != nil {
				assert.Error(t, err)
				assert.EqualError(t, err, tc.expectedError.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
