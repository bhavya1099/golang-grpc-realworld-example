package store_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/jinzhu/gorm"
	"github.com/raahii/golang-grpc-realworld-example/model"
	"github.com/raahii/golang-grpc-realworld-example/store"
)

// MockDB is a mock implementation of *gorm.DB
type MockDB struct {
	mock.Mock
}

func (m *MockDB) Begin() *gorm.DB {
	args := m.Called()
	return args.Get(0).(*gorm.DB)
}

func (m *MockDB) Rollback() *gorm.DB {
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

func TestArticleStoreAddFavorite(t *testing.T) {
	// Test Scenario 1: AddFavorite Successfully Increases Favorites Count
	t.Run("AddFavorite Successfully Increases Favorites Count", func(t *testing.T) {
		// Arrange
		article := &model.Article{FavoritesCount: 0}
		user := &model.User{}
		mockDB := new(MockDB)
		mockDB.On("Begin").Return(mockDB)
		mockDB.On("Model", article).Return(mockDB)
		mockDB.On("Model", article).Return(mockDB)
		mockDB.On("Commit").Return(mockDB)

		articleStore := &store.ArticleStore{db: mockDB}

		// Act
		err := articleStore.AddFavorite(article, user)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, int32(1), article.FavoritesCount)
		mockDB.AssertExpectations(t)
	})

	// Test Scenario 2: AddFavorite Fails to Add User to Favorited Users List
	t.Run("AddFavorite Fails to Add User to Favorited Users List", func(t *testing.T) {
		// Arrange
		article := &model.Article{FavoritesCount: 0}
		user := &model.User{}
		mockDB := new(MockDB)
		mockDB.On("Begin").Return(mockDB)
		mockDB.On("Model", article).Return(mockDB)
		mockDB.On("Rollback").Return(mockDB)

		articleStore := &store.ArticleStore{db: mockDB}

		// Act
		err := errors.New("error adding user to favorited users list")
		mockDB.On("Model", article).Return(mockDB)
		mockDB.On("Rollback").Return(mockDB)

		// Assert
		assert.Error(t, articleStore.AddFavorite(article, user))
		assert.Equal(t, int32(0), article.FavoritesCount)
		mockDB.AssertExpectations(t)
	})

	// Test Scenario 3: AddFavorite Fails to Update Favorites Count
	t.Run("AddFavorite Fails to Update Favorites Count", func(t *testing.T) {
		// Arrange
		article := &model.Article{FavoritesCount: 0}
		user := &model.User{}
		mockDB := new(MockDB)
		mockDB.On("Begin").Return(mockDB)
		mockDB.On("Model", article).Return(mockDB)
		mockDB.On("Rollback").Return(mockDB)

		articleStore := &store.ArticleStore{db: mockDB}

		// Act
		err := errors.New("error updating favorites count")
		mockDB.On("Model", article).Return(mockDB)
		mockDB.On("Rollback").Return(mockDB)

		// Assert
		assert.Error(t, articleStore.AddFavorite(article, user))
		assert.Equal(t, int32(0), article.FavoritesCount)
		mockDB.AssertExpectations(t)
	})

	// Test Scenario 4: AddFavorite Rolls Back Transaction on Error
	t.Run("AddFavorite Rolls Back Transaction on Error", func(t *testing.T) {
		// Arrange
		article := &model.Article{FavoritesCount: 0}
		user := &model.User{}
		mockDB := new(MockDB)
		mockDB.On("Begin").Return(mockDB)
		mockDB.On("Model", article).Return(mockDB)
		mockDB.On("Rollback").Return(mockDB)

		articleStore := &store.ArticleStore{db: mockDB}

		// Act
		err := errors.New("error during transaction")
		mockDB.On("Model", article).Return(mockDB)
		mockDB.On("Rollback").Return(mockDB)

		// Assert
		assert.Error(t, articleStore.AddFavorite(article, user))
		assert.Equal(t, int32(0), article.FavoritesCount)
		mockDB.AssertExpectations(t)
	})

	// Test Scenario 5: AddFavorite Handles Concurrent Requests Safely
	t.Run("AddFavorite Handles Concurrent Requests Safely", func(t *testing.T) {
		// Arrange
		article := &model.Article{FavoritesCount: 0}
		user1 := &model.User{}
		user2 := &model.User{}
		mockDB := new(MockDB)
		mockDB.On("Begin").Return(mockDB)
		mockDB.On("Model", article).Return(mockDB)
		mockDB.On("Model", article).Return(mockDB)
		mockDB.On("Commit").Return(mockDB)

		articleStore := &store.ArticleStore{db: mockDB}

		// Act
		go func() {
			articleStore.AddFavorite(article, user1)
		}()
		go func() {
			articleStore.AddFavorite(article, user2)
		}()

		// Assert
		assert.Eventually(t, func() bool {
			return article.FavoritesCount == 2
		}, "favorites count not updated correctly", "100ms", "10ms")
		mockDB.AssertExpectations(t)
	})
}
