package store

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/raahii/golang-grpc-realworld-example/model"
	"github.com/raahii/golang-grpc-realworld-example/store"
)

func TestArticleStoreGetComments(t *testing.T) {

	t.Run("Successful retrieval of comments for a valid article ID", func(t *testing.T) {

		mockArticle := &model.Article{ID: 1}
		expectedComments := []model.Comment{{Body: "Great article", UserID: 1, ArticleID: 1}}

		articleStore := store.ArticleStore{}
		comments, err := articleStore.GetComments(mockArticle)

		assert.NoError(t, err, "No error expected")
		assert.Equal(t, expectedComments, comments, "Comments should match expected")
		t.Log("Successfully retrieved comments for a valid article ID")
	})

	t.Run("No comments found for a non-existent article ID", func(t *testing.T) {

		mockArticle := &model.Article{ID: 999}

		articleStore := store.ArticleStore{}
		comments, err := articleStore.GetComments(mockArticle)

		assert.Error(t, err, "Error expected for non-existent article ID")
		assert.Empty(t, comments, "Comments should be empty")
		t.Log("No comments found for a non-existent article ID")
	})

	t.Run("Error handling when database query fails", func(t *testing.T) {

		mockArticle := &model.Article{ID: 1}
		expectedError := "Database error"

		articleStore := store.ArticleStore{}

		comments, err := articleStore.GetComments(mockArticle)

		assert.Error(t, err, "Error expected for database query failure")
		assert.Empty(t, comments, "Comments should be empty")
		assert.Contains(t, err.Error(), expectedError, "Error message should indicate database error")
		t.Log("Error handling when database query fails")
	})

	t.Run("Handling empty comments for a valid article ID", func(t *testing.T) {

		mockArticle := &model.Article{ID: 2}

		articleStore := store.ArticleStore{}
		comments, err := articleStore.GetComments(mockArticle)

		assert.NoError(t, err, "No error expected")
		assert.Empty(t, comments, "Comments should be empty")
		t.Log("Successfully handled empty comments for a valid article ID")
	})

	t.Run("Multiple comments retrieval for an article with duplicate comments", func(t *testing.T) {

		mockArticle := &model.Article{ID: 3}
		expectedComments := []model.Comment{{Body: "Interesting article", UserID: 1, ArticleID: 3},
			{Body: "Great insights", UserID: 2, ArticleID: 3},
			{Body: "Great insights", UserID: 2, ArticleID: 3}}

		articleStore := store.ArticleStore{}
		comments, err := articleStore.GetComments(mockArticle)

		assert.NoError(t, err, "No error expected")
		assert.ElementsMatch(t, expectedComments, comments, "Comments should include duplicates")
		t.Log("Successfully retrieved multiple comments for an article with duplicate comments")
	})
}
