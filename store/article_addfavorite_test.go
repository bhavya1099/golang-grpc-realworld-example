package store

import (
	"errors"
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/jinzhu/gorm"
	"github.com/raahii/golang-grpc-realworld-example/model"
	"github.com/raahii/golang-grpc-realworld-example/store"
)

type MockDB struct {
	mock.Mock
}func (m *MockDB) Begin() *gorm.DB {
	args := m.Called()
	return args.Get(0).(*gorm.DB)
}
func (m *MockDB) Commit() *gorm.DB {
	args := m.Called()
	return args.Get(0).(*gorm.DB)
}
func (m *MockDB) Model(value interface{}) *gorm.DB {
	args := m.Called(value)
	return args.Get(0).(*gorm.DB)
}
func (m *MockDB) Rollback() *gorm.DB {
	args := m.Called()
	return args.Get(0).(*gorm.DB)
}
func TestArticleStoreAddFavorite(t *testing.T) {

	t.Run("AddFavorite_Success", func(t *testing.T) {
		mockDB := new(MockDB)
		mockDB.On("Begin").Return(&gorm.DB{})
		mockDB.On("Model", mock.Anything).Return(&gorm.DB{})
		mockDB.On("Commit").Return(&gorm.DB{})

		s := &store.ArticleStore{db: mockDB}
		article := &model.Article{ID: 1, FavoritesCount: 0}
		user := &model.User{ID: 1}

		err := s.AddFavorite(article, user)

		assert.NoError(t, err)
		assert.Contains(t, article.FavoritedUsers, *user)
		assert.Equal(t, int32(1), article.FavoritesCount)
	})

	t.Run("AddFavorite_Fail_AssociationError", func(t *testing.T) {
		mockDB := new(MockDB)
		mockDB.On("Begin").Return(&gorm.DB{})
		mockDB.On("Model", mock.Anything).Return(&gorm.DB{})
		mockDB.On("Rollback").Return(&gorm.DB{})

		s := &store.ArticleStore{db: mockDB}
		article := &model.Article{ID: 1}
		user := &model.User{ID: 1}

		expectedError := errors.New("association error")
		mockDB.On("Model", article).Return(&gorm.DB{}).Once().Return(&gorm.DB{}).Once().Return(&gorm.DB{}).Once().Return(&gorm.DB{}).Once().Return(&gorm.DB{}).Once().Return(&gorm.DB{}).Return(&gorm.DB{})

		err := s.AddFavorite(article, user)

		assert.Error(t, err)
		assert.NotContains(t, article.FavoritedUsers, *user)
		assert.Equal(t, int32(0), article.FavoritesCount)
	})

	t.Run("AddFavorite_Fail_UpdateError", func(t *testing.T) {
		mockDB := new(MockDB)
		mockDB.On("Begin").Return(&gorm.DB{})
		mockDB.On("Model", mock.Anything).Return(&gorm.DB{})
		mockDB.On("Rollback").Return(&gorm.DB{})

		s := &store.ArticleStore{db: mockDB}
		article := &model.Article{ID: 1, FavoritesCount: 0}
		user := &model.User{ID: 1}

		expectedError := errors.New("update error")
		mockDB.On("Model", article).Return(&gorm.DB{}).Once().Return(&gorm.DB{}).Once().Return(&gorm.DB{}).Once().Return(&gorm.DB{}).Once().Return(&gorm.DB{}).Return(&gorm.DB{})

		err := s.AddFavorite(article, user)

		assert.Error(t, err)
		assert.Contains(t, article.FavoritedUsers, *user)
		assert.Equal(t, int32(0), article.FavoritesCount)
	})

	t.Run("AddFavorite_NullArticle", func(t *testing.T) {
		mockDB := new(MockDB)

		s := &store.ArticleStore{db: mockDB}
		user := &model.User{ID: 1}

		err := s.AddFavorite(nil, user)

		assert.Error(t, err)
	})

	t.Run("AddFavorite_NullUser", func(t *testing.T) {
		mockDB := new(MockDB)

		s := &store.ArticleStore{db: mockDB}
		article := &model.Article{ID: 1}

		err := s.AddFavorite(article, nil)

		assert.Error(t, err)
	})

	t.Run("AddFavorite_DatabaseError", func(t *testing.T) {
		mockDB := new(MockDB)
		mockDB.On("Begin").Return(&gorm.DB{})
		mockDB.On("Rollback").Return(&gorm.DB{})

		s := &store.ArticleStore{db: mockDB}
		article := &model.Article{ID: 1, FavoritesCount: 0}
		user := &model.User{ID: 1}

		expectedError := errors.New("database error")
		mockDB.On("Model", article).Return(&gorm.DB{}).Once().Return(&gorm.DB{}).Return(&gorm.DB{})

		err := s.AddFavorite(article, user)

		assert.Error(t, err)
		assert.Equal(t, int32(0), article.FavoritesCount)
	})
}
