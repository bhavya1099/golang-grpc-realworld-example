package handler

import (
	"context"
	"errors"
	"testing"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/raahii/golang-grpc-realworld-example/handler"
	"github.com/raahii/golang-grpc-realworld-example/model"
	"github.com/raahii/golang-grpc-realworld-example/proto"
	"github.com/raahii/golang-grpc-realworld-example/store"
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
func TestHandlerFavoriteArticle(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserStore := store.NewMockUserStore(ctrl)
	mockArticleStore := store.NewMockArticleStore(ctrl)

	h := &handler.Handler{
		us: mockUserStore,
		as: mockArticleStore,
	}

	t.Run("Successful favorite article request", func(t *testing.T) {
		mockUser := &model.User{ID: 1, Username: "testuser"}
		mockArticle := &model.Article{ID: 1, Title: "Test Article", Author: mockUser}
		mockUserStore.EXPECT().GetByID(gomock.Any()).Return(mockUser, nil)
		mockArticleStore.EXPECT().GetByID(gomock.Any()).Return(mockArticle, nil)
		mockArticleStore.EXPECT().AddFavorite(mockArticle, mockUser).Return(nil)

		req := &proto.FavoriteArticleRequest{Slug: "1"}
		resp, err := h.FavoriteArticle(context.Background(), req)

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, mockArticle.Title, resp.Article.Title)
		assert.True(t, resp.Article.Favorited)
	})

	t.Run("Unauthenticated user attempting to favorite an article", func(t *testing.T) {
		mockUserStore.EXPECT().GetByID(gomock.Any()).Return(nil, errors.New("user not found"))

		req := &proto.FavoriteArticleRequest{Slug: "1"}
		resp, err := h.FavoriteArticle(context.Background(), req)

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Equal(t, "unauthenticated", err.Error())
	})

	t.Run("User favoriting a non-existent article", func(t *testing.T) {
		mockUser := &model.User{ID: 1, Username: "testuser"}
		mockUserStore.EXPECT().GetByID(gomock.Any()).Return(mockUser, nil)
		mockArticleStore.EXPECT().GetByID(gomock.Any()).Return(nil, errors.New("article not found"))

		req := &proto.FavoriteArticleRequest{Slug: "999"}
		resp, err := h.FavoriteArticle(context.Background(), req)

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Equal(t, "requested article (slug=999) not found", err.Error())
	})

	t.Run("Error adding favorite to the article", func(t *testing.T) {
		mockUser := &model.User{ID: 1, Username: "testuser"}
		mockArticle := &model.Article{ID: 1, Title: "Test Article", Author: mockUser}
		mockUserStore.EXPECT().GetByID(gomock.Any()).Return(mockUser, nil)
		mockArticleStore.EXPECT().GetByID(gomock.Any()).Return(mockArticle, nil)
		mockArticleStore.EXPECT().AddFavorite(mockArticle, mockUser).Return(errors.New("error adding favorite"))

		req := &proto.FavoriteArticleRequest{Slug: "1"}
		resp, err := h.FavoriteArticle(context.Background(), req)

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Equal(t, "failed to add favorite", err.Error())
	})

	t.Run("Error retrieving following status for the article author", func(t *testing.T) {
		mockUser := &model.User{ID: 1, Username: "testuser"}
		mockArticle := &model.Article{ID: 1, Title: "Test Article", Author: mockUser}
		mockUserStore.EXPECT().GetByID(gomock.Any()).Return(mockUser, nil)
		mockArticleStore.EXPECT().GetByID(gomock.Any()).Return(mockArticle, nil)
		mockArticleStore.EXPECT().AddFavorite(mockArticle, mockUser).Return(nil)
		mockUserStore.EXPECT().IsFollowing(mockUser, gomock.Any()).Return(false, errors.New("error retrieving following status"))

		req := &proto.FavoriteArticleRequest{Slug: "1"}
		resp, err := h.FavoriteArticle(context.Background(), req)

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Equal(t, "failed to get following status", err.Error())
	})

	t.Run("Successful favorite article request with author following status", func(t *testing.T) {
		mockUser := &model.User{ID: 1, Username: "testuser"}
		mockArticle := &model.Article{ID: 1, Title: "Test Article", Author: mockUser}
		mockUserStore.EXPECT().GetByID(gomock.Any()).Return(mockUser, nil)
		mockArticleStore.EXPECT().GetByID(gomock.Any()).Return(mockArticle, nil)
		mockArticleStore.EXPECT().AddFavorite(mockArticle, mockUser).Return(nil)
		mockUserStore.EXPECT().IsFollowing(mockUser, gomock.Any()).Return(true, nil)

		req := &proto.FavoriteArticleRequest{Slug: "1"}
		resp, err := h.FavoriteArticle(context.Background(), req)

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, mockArticle.Title, resp.Article.Title)
		assert.True(t, resp.Article.Favorited)
		assert.True(t, resp.Article.Author.Following)
	})
}
