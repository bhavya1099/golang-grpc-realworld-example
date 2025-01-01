package store_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/raahii/golang-grpc-realworld-example/model"
	"github.com/raahii/golang-grpc-realworld-example/store"
)

// MockDB is a mock implementation of *gorm.DB
type MockDB struct {
	mock.Mock
}

// Find is a mock implementation of Find method in *gorm.DB
func (m *MockDB) Find(out interface{}, where ...interface{}) *store.DB {
	args := m.Called(out, where)
	return args.Get(0).(*store.DB)
}

func TestArticleStoreGetArticles(t *testing.T) {
	// Mocked ArticleStore with a mock DB
	mockDB := new(MockDB)
	articleStore := store.ArticleStore{db: mockDB}

	// Test Scenario 1: GetArticles with no filtering criteria
	t.Run("GetArticles with no filtering criteria", func(t *testing.T) {
		mockDB.On("Find", mock.Anything, mock.Anything).Return(&store.DB{})

		articles, err := articleStore.GetArticles("", "", nil, 0, 0)

		assert.NoError(t, err)
		assert.NotNil(t, articles)
	})

	// Test Scenario 2: GetArticles with username filtering
	t.Run("GetArticles with username filtering", func(t *testing.T) {
		mockDB.On("Find", mock.Anything, mock.Anything).
			Return(&store.DB{})

		articles, err := articleStore.GetArticles("testUser", "", nil, 0, 0)

		assert.NoError(t, err)
		assert.NotNil(t, articles)
	})

	// Test Scenario 3: GetArticles with tagName filtering
	t.Run("GetArticles with tagName filtering", func(t *testing.T) {
		mockDB.On("Find", mock.Anything, mock.Anything).
			Return(&store.DB{})

		articles, err := articleStore.GetArticles("", "testTag", nil, 0, 0)

		assert.NoError(t, err)
		assert.NotNil(t, articles)
	})

	// Test Scenario 4: GetArticles with favoritedBy filtering
	t.Run("GetArticles with favoritedBy filtering", func(t *testing.T) {
		mockDB.On("Find", mock.Anything, mock.Anything).
			Return(&store.DB{})

		favoritedBy := &model.User{ID: 1}
		articles, err := articleStore.GetArticles("", "", favoritedBy, 0, 0)

		assert.NoError(t, err)
		assert.NotNil(t, articles)
	})

	// Test Scenario 5: GetArticles with limit and offset
	t.Run("GetArticles with limit and offset", func(t *testing.T) {
		mockDB.On("Find", mock.Anything, mock.Anything).
			Return(&store.DB{})

		articles, err := articleStore.GetArticles("", "", nil, 10, 5)

		assert.NoError(t, err)
		assert.NotNil(t, articles)
	})

	// Test Scenario 6: GetArticles error handling
	t.Run("GetArticles error handling", func(t *testing.T) {
		mockDB.On("Find", mock.Anything, mock.Anything).
			Return(&store.DB{Error: assert.AnError})

		articles, err := articleStore.GetArticles("", "", nil, 0, 0)

		assert.Error(t, err)
		assert.Empty(t, articles)
	})
}
