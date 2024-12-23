package store

import (
	"errors"
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/jinzhu/gorm"
	"github.com/raahii/golang-grpc-realworld-example/model"
	"github.com/your/package/store"
)

type MockDB struct {
	mock.Mock
}func (m *MockDB) Begin() *gorm.DB {
	args := m.Called()
	return args.Get(0).(*gorm.DB)
}
func TestArticleStoreDeleteFavorite(t *testing.T) {
	t.Parallel()

	mockArticle := &model.Article{
		Model:          gorm.Model{ID: 1},
		Title:          "Test Article",
		FavoritesCount: 1,
	}
	mockUser := &model.User{
		Model:    gorm.Model{ID: 1},
		Username: "testuser",
	}

	t.Run("DeleteFavorite_NormalOperation", func(t *testing.T) {

		mockDB := new(MockDB)
		mockDB.On("Begin").Return(mockDB)
		mockDB.On("Model", mockArticle).Return(mockDB)
		mockDB.On("Association", "FavoritedUsers").Return(mockDB)
		mockDB.On("Delete", mockUser).Return(nil)
		mockDB.On("Model", mockArticle).Return(mockDB)
		mockDB.On("Update", "favorites_count", mock.Anything).Return(nil)
		mockDB.On("Commit").Return(mockDB)

		articleStore := store.ArticleStore{db: mockDB}

		err := articleStore.DeleteFavorite(mockArticle, mockUser)

		assert.NoError(t, err)
		assert.Equal(t, int32(0), mockArticle.FavoritesCount)
		mockDB.AssertExpectations(t)
	})

	t.Run("DeleteFavorite_FailToDeleteFavorite", func(t *testing.T) {

		mockDB := new(MockDB)
		mockDB.On("Begin").Return(mockDB)
		mockDB.On("Model", mockArticle).Return(mockDB)
		mockDB.On("Association", "FavoritedUsers").Return(mockDB)
		mockDB.On("Delete", mockUser).Return(errors.New("delete error"))
		mockDB.On("Rollback").Return(mockDB)

		articleStore := store.ArticleStore{db: mockDB}

		err := articleStore.DeleteFavorite(mockArticle, mockUser)

		assert.Error(t, err)
		mockDB.AssertExpectations(t)
	})

}
