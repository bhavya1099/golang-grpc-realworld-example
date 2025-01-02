package handler

import (
	"context"
	"testing"
	"github.com/stretchr/testify/assert"
	pb "github.com/raahii/golang-grpc-realworld-example/proto"
)

func TestCreateComment(t *testing.T) {
	handler := &Handler{}

	t.Run("Successful Comment Creation", func(t *testing.T) {
		validRequest := &pb.CreateCommentRequest{
			Slug: "123",
			Comment: &pb.CreateCommentRequest_Comment{
				Body: "Test Comment Body",
			},
		}

		resp, err := handler.CreateComment(context.Background(), validRequest)

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.NotNil(t, resp.Comment)
		assert.Equal(t, "Test Comment Body", resp.Comment.Body)
	})

	t.Run("Unauthenticated User", func(t *testing.T) {
		_, err := handler.CreateComment(context.Background(), &pb.CreateCommentRequest{})

		assert.Error(t, err)
		assert.EqualError(t, err, "unauthenticated")
	})
}
