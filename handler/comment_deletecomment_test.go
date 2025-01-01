package handler

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/raahii/golang-grpc-realworld-example/model"
	pb "github.com/raahii/golang-grpc-realworld-example/proto"
	"github.com/stretchr/testify/assert"
)

func TestHandlerDeleteComment(t *testing.T) {
	// Mock Logger, UserStore, and ArticleStore for testing purposes
	logger := // mock logger
	us := // mock UserStore
	as := // mock ArticleStore

	h := &Handler{logger: logger, us: us, as: as}

	tt := []struct {
		name     string
		ctx      context.Context
		req      *pb.DeleteCommentRequest
		userID   uint
		userErr  error
		comment  *model.Comment
		commentErr error
		expected error
	}{
		{
			name: "Delete comment successfully",
			ctx: context.TODO(),
			req: &pb.DeleteCommentRequest{Slug: "article-1", Id: "1"},
			userID: 1,
			userErr: nil,
			comment: &model.Comment{ID: 1, ArticleID: 1, UserID: 1},
			commentErr: nil,
			expected: nil,
		},
		{
			name: "Unauthenticated user deletion attempt",
			ctx: context.TODO(),
			req: &pb.DeleteCommentRequest{Slug: "article-1", Id: "1"},
			userID: 0,
			userErr: errors.New("unauthenticated"),
			comment: nil,
			commentErr: nil,
			expected: status.Error(codes.Unauthenticated, "unauthenticated"),
		},
		{
			name: "Delete comment with invalid ID",
			ctx: context.TODO(),
			req: &pb.DeleteCommentRequest{Slug: "article-1", Id: "invalid"},
			userID: 1,
			userErr: nil,
			comment: nil,
			commentErr: errors.New("invalid ID"),
			expected: status.Error(codes.InvalidArgument, "invalid article id"),
		},
		{
			name: "Delete comment not belonging to the user",
			ctx: context.TODO(),
			req: &pb.DeleteCommentRequest{Slug: "article-1", Id: "2"},
			userID: 1,
			userErr: nil,
			comment: &model.Comment{ID: 2, ArticleID: 1, UserID: 2},
			commentErr: nil,
			expected: status.Error(codes.InvalidArgument, "forbidden"),
		},
		{
			name: "Delete comment in a different article",
			ctx: context.TODO(),
			req: &pb.DeleteCommentRequest{Slug: "article-1", Id: "1"},
			userID: 1,
			userErr: nil,
			comment: &model.Comment{ID: 1, ArticleID: 2, UserID: 1},
			commentErr: nil,
			expected: status.Error(codes.InvalidArgument, "the comment is not in the article"),
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			// Mock GetUserID and GetByID functions
			auth.GetUserID = func(ctx context.Context) (uint, error) {
				return tc.userID, tc.userErr
			}
			us.GetByID = func(id uint) (*model.User, error) {
				return &model.User{ID: id}, nil
			}
			as.GetCommentByID = func(id uint) (*model.Comment, error) {
				return tc.comment, tc.commentErr
			}
			as.DeleteComment = func(m *model.Comment) error {
				return nil
			}

			_, err := h.DeleteComment(tc.ctx, tc.req)
			assert.Equal(t, tc.expected, err)
		})
	}
}
