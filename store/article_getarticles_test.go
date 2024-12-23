package store

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/raahii/golang-grpc-realworld-example/model"
	"github.com/raahii/golang-grpc-realworld-example/store"
)

func TestArticleStoreGetArticles(t *testing.T) {

	t.Run("GetArticles with no filtering criteria", func(t *testing.T) {

		articleStore := store.ArticleStore{db: &gorm.DB{}}

		articles, err := articleStore.GetArticles("", "", nil, 10, 0)

		require.NoError(t, err)
		assert.NotNil(t, articles)
		assert.NotEmpty(t, articles)
	})

	t.Run("GetArticles with tag name filtering", func(t *testing.T) {

		articleStore := store.ArticleStore{db: &gorm.DB{}}
		tagName := "technology"

		articles, err := articleStore.GetArticles(tagName, "", nil, 10, 0)

		require.NoError(t, err)
		assert.NotNil(t, articles)
		for _, article := range articles {
			assert.Contains(t, article.Tags, tagName)
		}
	})

	t.Run("GetArticles with username filtering", func(t *testing.T) {

		articleStore := store.ArticleStore{db: &gorm.DB{}}
		username := "testuser"

		articles, err := articleStore.GetArticles("", username, nil, 10, 0)

		require.NoError(t, err)
		assert.NotNil(t, articles)
		for _, article := range articles {
			assert.Equal(t, username, article.Author.Username)
		}
	})

	t.Run("GetArticles with favorite filtering", func(t *testing.T) {

		articleStore := store.ArticleStore{db: &gorm.DB{}}
		favoritedBy := &model.User{ID: 123}

		articles, err := articleStore.GetArticles("", "", favoritedBy, 10, 0)

		require.NoError(t, err)
		assert.NotNil(t, articles)
		for _, article := range articles {
			assert.Contains(t, article.FavoritedUsers, *favoritedBy)
		}
	})

	t.Run("GetArticles with pagination", func(t *testing.T) {

		articleStore := store.ArticleStore{db: &gorm.DB{}}
		limit := 5
		offset := 5

		articles, err := articleStore.GetArticles("", "", nil, int64(limit), int64(offset))

		require.NoError(t, err)
		assert.NotNil(t, articles)
		assert.Equal(t, limit, len(articles))
	})

	t.Run("GetArticles error handling", func(t *testing.T) {

		articleStore := store.ArticleStore{db: &gorm.DB{}}

		articles, err := articleStore.GetArticles("", "", nil, 10, 0)

		assert.Error(t, err)
		assert.Empty(t, articles)
	})
}
