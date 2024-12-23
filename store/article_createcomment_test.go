package store

import (
	"errors"
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/jinzhu/gorm"
	"github.com/raahii/golang-grpc-realworld-example/model"
	"github.com/raahii/golang-grpc-realworld-example/store"
)

func TestArticleStoreCreateComment(t *testing.T) {

	t.Run("CreateComment_SuccessfulCreation", func(t *testing.T) {

		mockDB, _ := gorm.Open("sqlite3", ":memory:")
		defer mockDB.Close()
		articleStore := store.ArticleStore{db: mockDB}
		comment := &model.Comment{Body: "Test Comment", UserID: 1, ArticleID: 1}

		err := articleStore.CreateComment(comment)

		assert.NoError(t, err, "Expected no error for successful comment creation")
		t.Log("CreateComment_SuccessfulCreation test passed")
	})

	t.Run("CreateComment_EmptyComment", func(t *testing.T) {

		mockDB, _ := gorm.Open("sqlite3", ":memory:")
		defer mockDB.Close()
		articleStore := store.ArticleStore{db: mockDB}
		emptyComment := &model.Comment{}

		err := articleStore.CreateComment(emptyComment)

		assert.Error(t, err, "Expected error for empty comment creation")
		t.Log("CreateComment_EmptyComment test passed")
	})

	t.Run("CreateComment_DBError", func(t *testing.T) {

		mockDB, _ := gorm.Open("sqlite3", ":memory:")
		mockDB.Close()
		articleStore := store.ArticleStore{db: mockDB}
		comment := &model.Comment{Body: "Test Comment", UserID: 1, ArticleID: 1}

		err := articleStore.CreateComment(comment)

		assert.Error(t, err, "Expected error for DB error during comment creation")
		t.Log("CreateComment_DBError test passed")
	})

	t.Run("CreateComment_NilStore", func(t *testing.T) {

		var articleStore *store.ArticleStore
		comment := &model.Comment{Body: "Test Comment", UserID: 1, ArticleID: 1}

		err := articleStore.CreateComment(comment)

		assert.Error(t, err, "Expected error for nil store during comment creation")
		t.Log("CreateComment_NilStore test passed")
	})

	t.Run("CreateComment_NilComment", func(t *testing.T) {

		mockDB, _ := gorm.Open("sqlite3", ":memory:")
		defer mockDB.Close()
		articleStore := store.ArticleStore{db: mockDB}
		var nilComment *model.Comment

		err := articleStore.CreateComment(nilComment)

		assert.Error(t, err, "Expected error for nil comment during creation")
		t.Log("CreateComment_NilComment test passed")
	})

	t.Run("CreateComment_DuplicateComment", func(t *testing.T) {

		mockDB, _ := gorm.Open("sqlite3", ":memory:")
		defer mockDB.Close()
		articleStore := store.ArticleStore{db: mockDB}
		comment := &model.Comment{Body: "Test Comment", UserID: 1, ArticleID: 1}
		_ = articleStore.CreateComment(comment)

		err := articleStore.CreateComment(comment)

		assert.Error(t, err, "Expected error for duplicate comment creation")
		t.Log("CreateComment_DuplicateComment test passed")
	})
}
