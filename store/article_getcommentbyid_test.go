package store_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/raahii/golang-grpc-realworld-example/model"
	"github.com/raahii/golang-grpc-realworld-example/store"
)

type MockDB struct {
	mock.Mock
}

func (m *MockDB) Find(out interface{}, where ...interface{}) *store.DB {
	args := m.Called(out, where)
	return args.Get(0).(*store.DB)
}

func TestArticleStoreGetCommentByID(t *testing.T) {
	t.Parallel()

	// Scenario 1: Retrieve Comment Successfully
	t.Run("Retrieve Comment Successfully", func(t *testing.T) {
		mockDB := new(MockDB)
		commentID := uint(1)
		expectedComment := &model.Comment{ID: commentID, Body: "Test Comment"}

		mockDB.On("Find", &model.Comment{}, commentID).Return(nil)

		articleStore := store.ArticleStore{db: mockDB}
		comment, err := articleStore.GetCommentByID(commentID)

		assert.Equal(t, expectedComment, comment, "Returned comment should match expected comment")
		assert.Nil(t, err, "Error should be nil")
	})

	// Scenario 2: Comment Not Found
	t.Run("Comment Not Found", func(t *testing.T) {
		mockDB := new(MockDB)
		commentID := uint(999)

		mockDB.On("Find", &model.Comment{}, commentID).Return(errors.New("record not found"))

		articleStore := store.ArticleStore{db: mockDB}
		comment, err := article	commentID)

		assert.Nil(t, comment, "Comment should be nil")
		assert.Error(t, err, "Error should be returned for non-existent comment")
	})

	// Scenario 3: Database Error Handling
	t.Run("Database Error Handling", func(t *testing.T) {
		mockDB := new(MockDB)
		commentID := uint(2)
		expectedError := errors.New("database error")

		mockDB.On("Find", &model.Comment{}, commentID).Return(expectedError)

		articleStore := store.ArticleStore{db: mockDB}
		comment, err := articleStore.GetCommentByID(commentID)

		assert.Nil(t, comment, "Comment should be nil")
		assert.EqualError(t, err, expectedError.Error(), "Error message should match expected error")
	})

	// Scenario 4: Empty Comment ID
	t.Run("Empty Comment ID", func(t *testing.T) {
		mockDB := new(MockDB)
		commentID := uint(0)
		expectedError := errors.New("invalid comment ID")

		articleStore := store.ArticleStore{db: mockDB}
		comment, err := articleStore.GetCommentByID(commentID)

		assert.Nil(t, comment, "Comment should be nil")
		assert.EqualError(t, err, expectedError.Error(), "Error message should indicate invalid comment ID")
	})
}
