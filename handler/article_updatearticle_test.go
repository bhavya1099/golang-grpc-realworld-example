package handler

import (
	"context"
	"errors"
	"testing"
	"github.com/raahii/golang-grpc-realworld-example/handler"
	"github.com/raahii/golang-grpc-realworld-example/model"
	pb "github.com/raahii/golang-grpc-realworld-example/proto"
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




type T struct {
	common
	isEnvSet bool
	context  *testContext // For running tests and subtests.
}

func (m *mockUserStore) GetByID(id uint) (*model.User, error) {
	if m.user.ID == id {
		return m.user, nil
	}
	return nil, errors.New("user not found")
}
func TestHandlerUpdateArticle(t *testing.T) {

	mockUserID := uint(1)
	mockArticleID := uint(10)
	mockSlug := "10"
	mockTitle := "Updated Title"
	mockDescription := "Updated Description"
	mockBody := "Updated Body"
	mockRequest := &pb.UpdateArticleRequest{
		Article: &pb.UpdateArticleRequest_Article{
			Slug:        mockSlug,
			Title:       mockTitle,
			Description: mockDescription,
			Body:        mockBody,
		},
	}

	ctx := context.TODO()

	mockHandler := &handler.Handler{
		us: &mockUserStore{},
		as: &mockArticleStore{},
	}

	t.Run("Successful article update by the author", func(t *testing.T) {
		mockUser := &model.User{ID: mockUserID}
		mockArticle := &model.Article{ID: mockArticleID, Author: mockUser}

		mockUserStore := &mockUserStore{mockUser}
		mockArticleStore := &mockArticleStore{mockArticle}

		mockHandler.us = mockUserStore
		mockHandler.as = mockArticleStore

		response, err := mockHandler.UpdateArticle(ctx, mockRequest)

		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, mockTitle, response.Article.Title)
		assert.Equal(t, mockDescription, response.Article.Description)
		assert.Equal(t, mockBody, response.Article.Body)
	})

}
