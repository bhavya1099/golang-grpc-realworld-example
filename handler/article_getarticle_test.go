package handler

import (
	"context"
	"errors"
	"testing"
	"github.com/golang/mock/gomock"
	"github.com/raahii/golang-grpc-realworld-example/handler"
	"github.com/raahii/golang-grpc-realworld-example/model"
	"github.com/raahii/golang-grpc-realworld-example/proto"
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


type Article struct {
	gorm.Model
	Title          string `gorm:"not null"`
	Description    string `gorm:"not null"`
	Body           string `gorm:"not null"`
	Tags           []Tag  `gorm:"many2many:article_tags"`
	Author         User   `gorm:"foreignkey:UserID"`
	UserID         uint   `gorm:"not null"`
	FavoritesCount int32  `gorm:"not null;default=0"`
	FavoritedUsers []User `gorm:"many2many:favorite_articles"`
	Comments       []Comment
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
func TestHandlerGetArticle(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLogger := NewMockLogger(ctrl)
	mockUserStore := NewMockUserStore(ctrl)
	mockArticleStore := NewMockArticleStore(ctrl)

	h := &handler.Handler{
		logger: mockLogger,
		us:     mockUserStore,
		as:     mockArticleStore,
	}

	validSlug := "123"
	invalidSlug := "abc"

	t.Run("GetArticle returns article response for a valid article request", func(t *testing.T) {
		mockArticle := &model.Article{ID: 123, Title: "Test Article", Description: "Test Description", Body: "Test Body", FavoritesCount: 0}
		mockUser := &model.User{ID: 1, Username: "testuser", Bio: "Test Bio", Image: "test.jpg"}

		mockArticleStore.EXPECT().GetByID(gomock.Any()).Return(mockArticle, nil)
		mockUserStore.EXPECT().GetByID(gomock.Any()).Return(mockUser, nil)
		mockArticleStore.EXPECT().IsFavorited(mockArticle, mockUser).Return(false, nil)
		mockUserStore.EXPECT().IsFollowing(mockUser, mockArticle.Author).Return(false, nil)

		ctx := context.Background()
		req := &proto.GetArticleRequest{Slug: validSlug}
		resp, err := h.GetArticle(ctx, req)

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.NotNil(t, resp.Article)
		assert.Equal(t, mockArticle.Title, resp.Article.Title)
	})

	t.Run("GetArticle returns error for invalid article slug", func(t *testing.T) {
		ctx := context.Background()
		req := &proto.GetArticleRequest{Slug: invalidSlug}

		mockLogger.EXPECT().Error().Times(2)
		resp, err := h.GetArticle(ctx, req)

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Equal(t, codes.InvalidArgument, status.Code(err))
	})

	t.Run("GetArticle handles article not found error", func(t *testing.T) {
		mockArticleStore.EXPECT().GetByID(gomock.Any()).Return(nil, errors.New("article not found"))

		ctx := context.Background()
		req := &proto.GetArticleRequest{Slug: validSlug}
		resp, err := h.GetArticle(ctx, req)

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Equal(t, codes.InvalidArgument, status.Code(err))
	})

	t.Run("GetArticle returns article with favorited status for authenticated user", func(t *testing.T) {
		mockArticle := &model.Article{ID: 123, Title: "Test Article", Description: "Test Description", Body: "Test Body", FavoritesCount: 0}
		mockUser := &model.User{ID: 1, Username: "testuser", Bio: "Test Bio", Image: "test.jpg"}

		mockArticleStore.EXPECT().GetByID(gomock.Any()).Return(mockArticle, nil)
		mockUserStore.EXPECT().GetByID(gomock.Any()).Return(mockUser, nil)
		mockArticleStore.EXPECT().IsFavorited(mockArticle, mockUser).Return(true, nil)
		mockUserStore.EXPECT().IsFollowing(mockUser, mockArticle.Author).Return(false, nil)

		ctx := context.Background()
		req := &proto.GetArticleRequest{Slug: validSlug}
		resp, err := h.GetArticle(ctx, req)

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.NotNil(t, resp.Article)
		assert.True(t, resp.Article.Favorited)
	})

	t.Run("GetArticle returns article with following status for authenticated user", func(t *testing.T) {
		mockArticle := &model.Article{ID: 123, Title: "Test Article", Description: "Test Description", Body: "Test Body", FavoritesCount: 0}
		mockUser := &model.User{ID: 1, Username: "testuser", Bio: "Test Bio", Image: "test.jpg"}

		mockArticleStore.EXPECT().GetByID(gomock.Any()).Return(mockArticle, nil)
		mockUserStore.EXPECT().GetByID(gomock.Any()).Return(mockUser, nil)
		mockArticleStore.EXPECT().IsFavorited(mockArticle, mockUser).Return(false, nil)
		mockUserStore.EXPECT().IsFollowing(mockUser, mockArticle.Author).Return(true, nil)

		ctx := context.Background()
		req := &proto.GetArticleRequest{Slug: validSlug}
		resp, err := h.GetArticle(ctx, req)

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.NotNil(t, resp.Article)
		assert.True(t, resp.Article.Author.Following)
	})

	t.Run("GetArticle returns internal server error for failed favorited check", func(t *testing.T) {
		mockArticle := &model.Article{ID: 123, Title: "Test Article", Description: "Test Description", Body: "Test Body", FavoritesCount: 0}
		mockUser := &model.User{ID: 1, Username: "testuser", Bio: "Test Bio", Image: "test.jpg"}

		mockArticleStore.EXPECT().GetByID(gomock.Any()).Return(mockArticle, nil)
		mockUserStore.EXPECT().GetByID(gomock.Any()).Return(mockUser, nil)
		mockArticleStore.EXPECT().IsFavorited(mockArticle, mockUser).Return(false, errors.New("favorited check failed"))

		ctx := context.Background()
		req := &proto.GetArticleRequest{Slug: validSlug}
		resp, err := h.GetArticle(ctx, req)

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Equal(t, codes.Aborted, status.Code(err))
	})

	t.Run("GetArticle returns internal server error for failed following check", func(t *testing.T) {
		mockArticle := &model.Article{ID: 123, Title: "Test Article", Description: "Test Description", Body: "Test Body", FavoritesCount: 0}
		mockUser := &model.User{ID: 1, Username: "testuser", Bio: "Test Bio", Image: "test.jpg"}

		mockArticleStore.EXPECT().GetByID(gomock.Any()).Return(mockArticle, nil)
		mockUserStore.EXPECT().GetByID(gomock.Any()).Return(mockUser, nil)
		mockArticleStore.EXPECT().IsFavorited(mockArticle, mockUser).Return(false, nil)
		mockUserStore.EXPECT().IsFollowing(mockUser, mockArticle.Author).Return(false, errors.New("following check failed"))

		ctx := context.Background()
		req := &proto.GetArticleRequest{Slug: validSlug}
		resp, err := h.GetArticle(ctx, req)

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Equal(t, codes.NotFound, status.Code(err))
	})
}
