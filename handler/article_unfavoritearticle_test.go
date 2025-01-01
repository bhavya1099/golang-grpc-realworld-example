package handler

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"testing"
	"github.com/raahii/golang-grpc-realworld-example/model"
	pb "github.com/raahii/golang-grpc-realworld-example/proto"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
func TestHandlerUnfavoriteArticle(t *testing.T) {

	t.Run("Unfavorite an Article Successfully", func(t *testing.T) {
		mockUserID := uint(1)
		mockArticleID := uint(1)

		mockHandler := &Handler{
			logger: nil,
			us: &store.UserStore{

				GetByID: func(id uint) (*model.User, error) {
					if id == mockUserID {
						return &model.User{}, nil
					}
					return nil, errors.New("user not found")
				},
			},
			as: &store.ArticleStore{

				GetByID: func(id uint) (*model.Article, error) {
					if id == mockArticleID {
						return &model.Article{}, nil
					}
					return nil, errors.New("article not found")
				},
				DeleteFavorite: func(a *model.Article, u *model.User) error {
					return nil
				},
			},
		}

		req := &pb.UnfavoriteArticleRequest{Slug: strconv.Itoa(int(mockArticleID))}
		articleResponse, err := mockHandler.UnfavoriteArticle(context.Background(), req)

		assert.NoError(t, err)
		assert.NotNil(t, articleResponse)
		assert.False(t, articleResponse.Article.Favorited)
	})

}
