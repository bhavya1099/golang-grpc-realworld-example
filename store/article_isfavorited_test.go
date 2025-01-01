package store

import (
	"testing"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
	"github.com/raahii/golang-grpc-realworld-example/model"
)

// TestArticleStoreIsFavorited tests the IsFavorited function in ArticleStore
func TestArticleStoreIsFavorited(t *testing.T) {
	// Scenario 1: Test IsFavorited when both article and user are provided
	t.Run("IsFavorited with valid article and user", func(t *testing.T) {
		// Arrange
		article := &model.Article{ID: 1}
		user := &model.User{ID: 1}

		// Act
		favorited, err := (&ArticleStore{db: &gorm.DB{}}).IsFavorited(article, user)

		// Assert
		assert.NoError(t, err)
		assert.True(t, favorited)
		t.Log("IsFavorited with valid article and user passed")
	})

	// Scenario 2: Test IsFavorited with nil article and user
	t.Run("IsFavorited with nil article and user", func(t *testing.T) {
		// Act
		favorited, err := (&ArticleStore{db: &gorm.DB{}}).IsFavorited(nil, nil)

		// Assert
		assert.NoError(t, err)
		assert.False(t, favorited)
		t.Log("IsFavorited with nil article and user passed")
	})

	// TODO: Add tests for remaining scenarios
}
