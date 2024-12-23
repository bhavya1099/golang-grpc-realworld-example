package store

import (
	"testing"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
	"github.com/raahii/golang-grpc-realworld-example/model"
)

func TestArticleStoreGetFeedArticles(t *testing.T) {

	mockDB, _, _ := gorm.Open("sqlite3", ":memory:")
	articleStore := ArticleStore{db: mockDB}

	t.Run("GetFeedArticles_NormalCase", func(t *testing.T) {
		userIDs := []uint{1, 2, 3}
		limit := int64(10)
		offset := int64(0)

		articles, err := articleStore.GetFeedArticles(userIDs, limit, offset)

		assert.NoError(t, err)
		assert.NotNil(t, articles)
		assert.Equal(t, 0, len(articles))
		t.Log("GetFeedArticles_NormalCase test passed")
	})

	t.Run("GetFeedArticles_EmptyUserIDs", func(t *testing.T) {
		userIDs := []uint{}
		limit := int64(10)
		offset := int64(0)

		articles, err := articleStore.GetFeedArticles(userIDs, limit, offset)

		assert.NoError(t, err)
		assert.NotNil(t, articles)
		assert.Equal(t, 0, len(articles))
		t.Log("GetFeedArticles_EmptyUserIDs test passed")
	})

	t.Run("GetFeedArticles_LimitZero", func(t *testing.T) {
		userIDs := []uint{1, 2, 3}
		limit := int64(0)
		offset := int64(0)

		articles, err := articleStore.GetFeedArticles(userIDs, limit, offset)

		assert.NoError(t, err)
		assert.NotNil(t, articles)
		assert.Equal(t, 0, len(articles))
		t.Log("GetFeedArticles_LimitZero test passed")
	})

	t.Run("GetFeedArticles_OffsetGreaterThanTotalArticles", func(t *testing.T) {
		userIDs := []uint{1, 2, 3}
		limit := int64(10)
		offset := int64(100)

		articles, err := articleStore.GetFeedArticles(userIDs, limit, offset)

		assert.NoError(t, err)
		assert.NotNil(t, articles)
		assert.Equal(t, 0, len(articles))
		t.Log("GetFeedArticles_OffsetGreaterThanTotalArticles test passed")
	})

	t.Run("GetFeedArticles_DatabaseError", func(t *testing.T) {

		mockDB, _, _ := gorm.Open("sqlite3", ":memory:")
		mockDB.Close()
		articleStore.db = mockDB

		userIDs := []uint{1, 2, 3}
		limit := int64(10)
		offset := int64(0)

		articles, err := articleStore.GetFeedArticles(userIDs, limit, offset)

		assert.Error(t, err)
		assert.Nil(t, articles)
		t.Log("GetFeedArticles_DatabaseError test passed")
	})
}
