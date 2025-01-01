package handler

import (
	"context"
	"errors"
	"testing"
	"github.com/raahii/golang-grpc-realworld-example/handler"
	"github.com/raahii/golang-grpc-realworld-example/model"
	"github.com/raahii/golang-grpc-realworld-example/proto"
	"github.com/stretchr/testify/assert"
)






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

func (m *MockArticleStore) GetArticles(tagName, username string, favoritedBy *model.User, limit, offset int64) ([]model.Article, error) {
	return m.MockGetArticles(tagName, username, favoritedBy, limit, offset)
}
func (m *MockUserStore) GetByID(id uint) (*model.User, error) {
	return m.MockGetByID(id)
}
func (m *MockArticleStore) IsFavorited(a *model.Article, u *model.User) (bool, error) {
	return m.MockIsFavorited(a, u)
}
func (m *MockUserStore) IsFollowing(a *model.User, b *model.User) (bool, error) {
	return m.MockIsFollowing(a, b)
}
func TestHandlerGetArticles(t *testing.T) {

	t.Run("GetArticles with default limit when limit is not provided", func(t *testing.T) {
		mockHandler := &handler.Handler{}
		request := &proto.GetArticlesRequest{}

		response, err := mockHandler.GetArticles(context.Background(), request)

		assert.Equal(t, int64(20), request.GetLimit())
		assert.Nil(t, err)
		assert.NotNil(t, response)
	})

	t.Run("GetArticles with favorited user not found", func(t *testing.T) {
		mockHandler := &handler.Handler{}
		request := &proto.GetArticlesRequest{Favorited: "nonexistentuser"}

		response, err := mockHandler.GetArticles(context.Background(), request)

		assert.Nil(t, response)
		assert.EqualError(t, err, "rpc error: code = NotFound desc = user not found")
	})

	t.Run("GetArticles with error searching articles in the database", func(t *testing.T) {
		mockHandler := &handler.Handler{}
		request := &proto.GetArticlesRequest{}

		mockHandler.ArticleStore = &MockArticleStore{MockGetArticles: func(tagName, username string, favoritedBy *model.User, limit, offset int64) ([]model.Article, error) {
			return nil, errors.New("database error")
		}}

		response, err := mockHandler.GetArticles(context.Background(), request)

		assert.Nil(t, response)
		assert.EqualError(t, err, "rpc error: code = Aborted desc = internal server error")
	})

	t.Run("GetArticles with current user not found", func(t *testing.T) {
		mockHandler := &handler.Handler{}
		request := &proto.GetArticlesRequest{}

		mockHandler.UserStore = &MockUserStore{MockGetByID: func(id uint) (*model.User, error) {
			return nil, errors.New("user not found")
		}}

		response, err := mockHandler.GetArticles(context.Background(), request)

		assert.Nil(t, response)
		assert.EqualError(t, err, "rpc error: code = NotFound desc = user not found")
	})

	t.Run("GetArticles with favorited status retrieval failure", func(t *testing.T) {
		mockHandler := &handler.Handler{}
		request := &proto.GetArticlesRequest{}

		mockHandler.ArticleStore = &MockArticleStore{MockIsFavorited: func(a *model.Article, u *model.User) (bool, error) {
			return false, errors.New("favorited status retrieval failure")
		}}

		response, err := mockHandler.GetArticles(context.Background(), request)

		assert.Nil(t, response)
		assert.EqualError(t, err, "rpc error: code = Aborted desc = internal server error")
	})

	t.Run("GetArticles with following status retrieval failure", func(t *testing.T) {
		mockHandler := &handler.Handler{}
		request := &proto.GetArticlesRequest{}

		mockHandler.UserStore = &MockUserStore{MockIsFollowing: func(a *model.User, b *model.User) (bool, error) {
			return false, errors.New("following status retrieval failure")
		}}

		response, err := mockHandler.GetArticles(context.Background(), request)

		assert.Nil(t, response)
		assert.EqualError(t, err, "rpc error: code = NotFound desc = internal server error")
	})
}
