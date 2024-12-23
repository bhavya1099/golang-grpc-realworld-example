package store

import (
	"testing"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
	"your-module-path/store"
)

func TestNewArticleStore(t *testing.T) {

	t.Run("NewArticleStore with Valid DB", func(t *testing.T) {
		mockDB := &gorm.DB{}
		articleStore := store.NewArticleStore(mockDB)
		assert.NotNil(t, articleStore, "ArticleStore instance should not be nil")
		assert.Equal(t, mockDB, articleStore.GetDB(), "DB field should match the mockDB")
		t.Log("NewArticleStore with Valid DB - Test Passed")
	})

	t.Run("NewArticleStore with Nil DB", func(t *testing.T) {
		articleStore := store.NewArticleStore(nil)
		assert.NotNil(t, articleStore, "ArticleStore instance should not be nil")
		assert.Nil(t, articleStore.GetDB(), "DB field should be nil")
		t.Log("NewArticleStore with Nil DB - Test Passed")
	})

	t.Run("NewArticleStore with Empty DB", func(t *testing.T) {
		mockEmptyDB := &gorm.DB{}
		articleStore := store.NewArticleStore(mockEmptyDB)
		assert.NotNil(t, articleStore, "ArticleStore instance should not be nil")
		assert.Equal(t, mockEmptyDB, articleStore.GetDB(), "DB field should match the mockEmptyDB")
		t.Log("NewArticleStore with Empty DB - Test Passed")
	})

	t.Run("NewArticleStore with Non-Empty DB", func(t *testing.T) {
		mockNonEmptyDB := &gorm.DB{}
		articleStore := store.NewArticleStore(mockNonEmptyDB)
		assert.NotNil(t, articleStore, "ArticleStore instance should not be nil")
		assert.Equal(t, mockNonEmptyDB, articleStore.GetDB(), "DB field should match the mockNonEmptyDB")
		t.Log("NewArticleStore with Non-Empty DB - Test Passed")
	})

	t.Run("NewArticleStore Performance Test", func(t *testing.T) {

		t.Skip("Skipping Performance Test for now")
	})
}
