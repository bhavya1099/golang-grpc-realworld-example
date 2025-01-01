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

func (m *MockEvent) Err(err error) *MockEvent {
	return m
}
func (m *MockLogger) Error() *MockEvent {
	return &MockEvent{}
}
func (m *MockUserStore) GetByID(id uint) (*model.User, error) {
	args := m.Called(id)
	user, _ := args.Get(0).(*model.User)
	err, _ := args.Error(1)
	return user, err
}
func (m *MockArticleStore) GetFeedArticles(userIDs []uint, limit, offset int64) ([]model.Article, error) {
	args := m.Called(userIDs, limit, offset)
	articles, _ := args.Get(0).([]model.Article)
	err, _ := args.Error(1)
	return articles, err
}
func (m *MockUserStore) GetFollowingUserIDs(user *model.User) ([]uint, error) {
	args := m.Called(user)
	userIDs, _ := args.Get(0).([]uint)
	err, _ := args.Error(1)
	return userIDs, err
}
func (m *MockLogger) Info() *MockEvent {
	return &MockEvent{}
}
func (m *MockEvent) Interface(key string, i interface{}) *MockEvent {
	return m
}
func (m *MockArticleStore) IsFavorited(a *model.Article, u *model.User) (bool, error) {
	args := m.Called(a, u)
	favorited, _ := args.Get(0).(bool)
	err, _ := args.Error(1)
	return favorited, err
}
func (m *MockUserStore) IsFollowing(a, b *model.User) (bool, error) {
	args := m.Called(a, b)
	following, _ := args.Get(0).(bool)
	err, _ := args.Error(1)
	return following, err
}
func (m *MockEvent) Msg(msg string) {}
func TestHandlerGetFeedArticles(t *testing.T) {

	mockUserStore := &MockUserStore{}
	mockArticleStore := &MockArticleStore{}
	logger := &MockLogger{}

	h := handler.Handler{
		logger: logger,
		us:     mockUserStore,
		as:     mockArticleStore,
	}

	t.Run("Test Successful Retrieval of Feed Articles", func(t *testing.T) {
		validRequest := &pb.GetFeedArticlesRequest{Limit: 10, Offset: 0}
		ctx := context.Background()

		mockUser := &model.User{ID: 1}
		mockUserStore.On("GetByID", mockUser.ID).Return(mockUser, nil)

		mockUserStore.On("GetFollowingUserIDs", mockUser).Return([]uint{2, 3}, nil)

		mockArticles := []model.Article{
			{ID: 1, Title: "Article 1", Description: "Description 1", Body: "Body 1", Author: model.User{ID: 2}},
			{ID: 2, Title: "Article 2", Description: "Description 2", Body: "Body 2", Author: model.User{ID: 3}},
		}
		mockArticleStore.On("GetFeedArticles", []uint{2, 3}, validRequest.Limit, validRequest.Offset).Return(mockArticles, nil)

		mockArticleStore.On("IsFavorited", &mockArticles[0], mockUser).Return(true, nil)
		mockArticleStore.On("IsFavorited", &mockArticles[1], mockUser).Return(false, nil)

		mockUserStore.On("IsFollowing", mockUser, &mockArticles[0].Author).Return(true, nil)
		mockUserStore.On("IsFollowing", mockUser, &mockArticles[1].Author).Return(false, nil)

		resp, err := h.GetFeedArticles(ctx, validRequest)

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, int32(len(mockArticles)), resp.ArticlesCount)
	})

}
