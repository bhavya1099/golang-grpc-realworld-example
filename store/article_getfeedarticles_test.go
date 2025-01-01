package store_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/raahii/golang-grpc-realworld-example/model"
	"github.com/raahii/golang-grpc-realworld-example/store"
)

func TestArticleStoreGetFeedArticles(t *testing.T) {
	mockDB := &store.DB{} // Mock DB for testing
	articleStore := &store.ArticleStore{db: mockDB}

	t.Run("GetFeedArticles retrieves articles for valid user IDs with limit and offset", func(t *testing.T) {
		userIDs := []uint{1, 2, 3}
		limit := int64(10)
		offset := int64(0)

		articles, err := articleStore.GetFeedArticles(userIDs, limit, offset)

		assert.NoError(t, err)
		assert.NotNil(t, articles)
		assert.NotEmpty(t, articles)
		t.Log("Retrieved articles:", articles)
	})

	t.Run("GetFeedArticles handles an empty list of user IDs gracefully", func(t *testing.T) {
		userIDs := []uint{}
		limit := int64(10)
		offset := int64(0)

		articles, err := articleStore.GetFeedArticles(userIDs, limit, offset)

		assert.NoError(t, err)
		assert.NotNil(t, articles)
		assert.Empty(t, articles)
		t.Log("Retrieved articles for empty user IDs:", articles)
	})

	t.Run("GetFeedArticles correctly applies the provided limit and offset", func(t *testing.T) {
		userIDs := []uint{1}
		limit := int64(2)
		offset := int64(1)

		articles, err := articleStore.GetFeedArticles(userIDs, limit, offset)

		assert.NoError(t, err)
		assert.NotNil(t, articles)
		assert.Equal(t, 2, len(articles))
		t.Log("Retrieved articles with limit and offset:", articles)
	})

	t.Run("GetFeedArticles handles database errors gracefully", func(t *testing.T) {
		mockDB.Error = assert.AnError
		userIDs := []uint{1}
		limit := int64(10)
		offset := int64(0)

		articles, err := articleStore.GetFeedArticles(userIDs, limit, offset)

		assert.Error(t, err)
		assert.Nil(t, articles)
		t.Log("Encountered database error:", err)
	})
}
