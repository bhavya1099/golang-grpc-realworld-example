package store_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/raahii/golang-grpc-realworld-example/model"
	"github.com/raahii/golang-grpc-realworld-example/store"
)

// Mock DB connection for testing
type MockDB struct{}

func (db *MockDB) Model(value interface{}) *MockDB {
	return db
}

func (db *MockDB) Update(attrs ...interface{}) *MockDB {
	return db
}

func TestArticleStoreUpdate(t *testing.T) {
	// Scenario 1: Update Article Successfully
	t.Run("Update Article Successfully", func(t *testing.T) {
		// Arrange
		article := &model.Article{
			Title:       "Test Title",
			Description: "Test Description",
			Body:        "Test Body",
			UserID:      1,
		}
		articleStore := store.ArticleStore{db: &MockDB{}}

		// Act
		err := articleStore.Update(article)

		// Assert
		assert.NoError(t, err, "Update should be successful")
		t.Log("Update Article Successfully test passed")
	})

	// Scenario 2: Update Article with Nil Pointer
	t.Run("Update Article with Nil Pointer", func(t *testing.T) {
		// Arrange
		var article *model.Article
		articleStore := store.ArticleStore{db: &MockDB{}}

		// Act
		err := articleStore.Update(article)

		// Assert
		assert.Error(t, err, "Update with nil pointer should return an error")
		t.Log("Update Article with Nil Pointer test passed")
	})

	// Scenario 3: Update Article Error Handling
	t.Run("Update Article Error Handling", func(t *testing.T) {
		// Arrange
		faultyArticle := &model.Article{
			Title: "Faulty Title",
			// Simulating a faulty article that triggers an error
		}
		articleStore := store.ArticleStore{db: &MockDB{}}

		// Act
		err := articleStore.Update(faultyArticle)

		// Assert
		assert.Error(t, err, "Update with faulty article should return an error")
		t.Log("Update Article Error Handling test passed")
	})

	// Scenario 4: Update Article with Empty Fields
	t.Run("Update Article with Empty Fields", func(t *testing.T) {
		// Arrange
		emptyArticle := &model.Article{
			// Empty article with no required fields filled
		}
		articleStore := store.ArticleStore{db: &MockDB{}}

		// Act
		err := articleStore.Update(emptyArticle)

		// Assert
		assert.Error(t, err, "Update with empty fields should return an error")
		t.Log("Update Article with Empty Fields test passed")
	})
}
