package store_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/raahii/golang-grpc-realworld-example/model"
	"github.com/raahii/golang-grpc-realworld-example/store"
)

type mockDB struct {
	mock.Mock
}

func (m *mockDB) Begin() *store.DB {
	args := m.Called()
	return args.Get(0).(*store.DB)
}

func (m *mockDB) Model(value interface{}) *store.DB {
	args := m.Called(value)
	return args.Get(0).(*store.DB)
}

func (m *mockDB) Commit() *store.DB {
	args := m.Called()
	return args.Get(0).(*store.DB)
}

func (m *mockDB) Rollback() *store.DB {
	args := m.Called()
	return args.Get(0).(*store.DB)
}

func TestArticleStoreDeleteFavorite(t *testing.T) {
	// Setup
	mockDB := new(mockDB)
	articleStore := store.ArticleStore{db: mockDB}

	article := &model.Article{FavoritesCount: 1}
	user := &model.User{}

	mockTx := new(store.DB)
	mockDB.On("Begin").Return(mockTx)
	mockDB.On("Model", article).Return(mockTx)
	mockTx.On("Model", article).Return(mockTx)
	mockTx.On("Commit").Return(mockTx)

	// Test Scenario 1: DeleteFavorite_Success
	t.Run("DeleteFavorite_Success", func(t *testing.T) {
		// Arrange
		mockTx.On("Model", article).Return(mockTx)
		mockTx.On("Model", article).Return(mockTx)
		mockTx.On("Commit").Return(mockTx)

		// Act
		err := articleStore.DeleteFavorite(article, user)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, int32(0), article.FavoritesCount)
	})

	// Test Scenario 2: DeleteFavorite_ErrorOnDelete
	t.Run("DeleteFavorite_ErrorOnDelete", func(t *testing.T) {
		// Arrange
		mockTx.On("Model", article).Return(mockTx)
		mockTx.On("Model", article).Return(mockTx)
		mockTx.On("Rollback").Return(mockTx)

		// Act
		mockTx.On("Model", article).Return(mockTx)
		mockTx.On("Delete", user).Return(errors.New("delete error"))
		err := articleStore.DeleteFavorite(article, user)

		// Assert
		assert.Error(t, err)
		assert.Equal(t, int32(1), article.FavoritesCount)
	})

	// Test Scenario 3: DeleteFavorite_ErrorOnUpdate
	t.Run("DeleteFavorite_ErrorOnUpdate", func(t *testing.T) {
		// Arrange
		mockTx.On("Model", article).Return(mockTx)
		mockTx.On("Model", article).Return(mockTx)
		mockTx.On("Rollback").Return(mockTx)

		// Act
		mockTx.On("Model", article).Return(mockTx)
		mockTx.On("Delete", user).Return(nil)
		mockTx.On("Update", "favorites_count", mock.Anything).Return(errors.New("update error"))
		err := articleStore.DeleteFavorite(article, user)

		// Assert
		assert.Error(t, err)
		assert.Equal(t, int32(1), article.FavoritesCount)
	})

	// Test Scenario 4: DeleteFavorite_NoFavoritedUsers
	t.Run("DeleteFavorite_NoFavoritedUsers", func(t *testing.T) {
		// Arrange
		article.FavoritedUsers = []model.User{}
		initialFavoritesCount := article.FavoritesCount

		// Act
		err := articleStore.DeleteFavorite(article, user)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, initialFavoritesCount, article.FavoritesCount)
	})

	// Test Scenario 5: DeleteFavorite_NilArticleOrUser
	t.Run("DeleteFavorite_NilArticleOrUser", func(t *testing.T) {
		// Arrange
		nilArticle := (*model.Article)(nil)
		nilUser := (*model.User)(nil)

		// Act
		err1 := articleStore.DeleteFavorite(nilArticle, user)
		err2 := articleStore.DeleteFavorite(article, nilUser)

		// Assert
		assert.Error(t, err1)
		assert.Error(t, err2)
		assert.Equal(t, int32(1), article.FavoritesCount)
	})
}
