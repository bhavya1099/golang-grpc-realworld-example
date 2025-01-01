package store_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/raahii/golang-grpc-realworld-example/model"
	"github.com/raahii/golang-grpc-realworld-example/store"
)

type MockDB struct {
	mock.Mock
}

func (m *MockDB) Preload(column string, conditions ...interface{}) *store.DB {
	args := m.Called(column, conditions)
	return args.Get(0).(*store.DB)
}

func (m *MockDB) Where(query interface{}, args ...interface{}) *store.DB {
	args := m.Called(query, args)
	return args.Get(0).(*store.DB)
}

func TestArticleStoreGetComments(t *testing.T) {
	// Test Scenario 1: Retrieve Comments for a Valid Article ID
	t.Run("Retrieve Comments for a Valid Article ID", func(t *testing.T) {
		mockArticle := &model.Article{ID: 1}
		mockComments := []model.Comment{{Body: "First comment"}, {Body: "Second comment"}}

		mockDB := new(MockDB)
		mockDB.On("Preload", "Author").Return(mockDB)
		mockDB.On("Where", "article_id = ?", mockArticle.ID).Return(mockDB)
		mockDB.On("Find", &mockComments).Return(nil)

		articleStore := store.ArticleStore{db: mockDB}
		comments, err := articleStore.GetComments(mockArticle)

		assert.NoError(t, err, "No error expected")
		assert.Equal(t, mockComments, comments, "Comments should match expected")
	})

	// Test Scenario 2: No Comments Found for Non-Existent Article ID
	t.Run("No Comments Found for Non-Existent Article ID", func(t *testing.T) {
		mockArticle := &model.Article{ID: 999}

		mockDB := new(MockDB)
		mockDB.On("Preload", "Author").Return(mockDB)
		mockDB.On("Where", "article_id = ?", mockArticle.ID).Return(mockDB)
		mockDB.On("Find", &[]model.Comment{}).Return(nil)

		articleStore := store.ArticleStore{db: mockDB}
		comments, err := articleStore.GetComments(mockArticle)

		assert.NoError(t, err, "No error expected")
		assert.Empty(t, comments, "No comments should be found")
	})

	// Test Scenario 3: Error Handling for Database Query Failure
	t.Run("Error Handling for Database Query Failure", func(t *testing.T) {
		mockArticle := &model.Article{ID: 1}

		mockDB := new(MockDB)
		mockDB.On("Preload", "Author").Return(mockDB)
		mockDB.On("Where", "article_id = ?", mockArticle.ID).Return(mockDB)
		mockDB.On("Find", &[]model.Comment{}).Return(mockDB, errors.New("database error"))

		articleStore := store.ArticleStore{db: mockDB}
		comments, err := articleStore.GetComments(mockArticle)

		assert.Error(t, err, "Error should be returned")
		assert.Empty(t, comments, "No comments should be found due to error")
	})
}
