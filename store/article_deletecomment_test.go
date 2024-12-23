package store

import (
	"errors"
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/jinzhu/gorm"
	"github.com/raahii/golang-grpc-realworld-example/model"
	"github.com/your_package/store"
)

func TestArticleStoreDeleteComment(t *testing.T) {

	t.Run("DeleteComment_SuccessfulDeletion", func(t *testing.T) {

		mockDB, _ := gorm.Open("sqlite3", ":memory:")
		defer mockDB.Close()
		articleStore := &store.ArticleStore{db: mockDB}
		comment := &model.Comment{ID: 1}

		err := articleStore.DeleteComment(comment)

		assert.NoError(t, err, "Expected no error for successful deletion")
		t.Log("Comment deleted successfully")
	})

	t.Run("DeleteComment_NilComment", func(t *testing.T) {

		mockDB, _ := gorm.Open("sqlite3", ":memory:")
		defer mockDB.Close()
		articleStore := &store.ArticleStore{db: mockDB}

		err := articleStore.DeleteComment(nil)

		assert.Error(t, err, "Expected error for nil comment")
		t.Log("Error received as expected for nil comment")
	})

	t.Run("DeleteComment_ErrorOnDeletion", func(t *testing.T) {

		mockDB, _ := gorm.Open("sqlite3", ":memory:")
		defer mockDB.Close()
		articleStore := &store.ArticleStore{db: mockDB}
		comment := &model.Comment{ID: 999}

		err := articleStore.DeleteComment(comment)

		assert.Error(t, err, "Expected error for deletion failure")
		assert.True(t, errors.Is(err, gorm.ErrRecordNotFound), "Expected error to be of type gorm.ErrRecordNotFound")
		t.Log("Deletion error handled correctly")
	})
}
