package store

import (
	"errors"
	"testing"
	"github.com/jinzhu/gorm"
	"github.com/raahii/golang-grpc-realworld-example/model"
	"github.com/stretchr/testify/assert"
)

type MockArticleStore struct {
	db *gorm.DB
}
func (s *MockArticleStore) Delete(m *model.Article) error {
	if m == nil {
		return errors.New("article is nil")
	}
	return s.db.Delete(m).Error
}
func TestArticleStoreDelete(t *testing.T) {
	mockDB, _ := gorm.Open("sqlite3", "test.db")
	mockArticleStore := &MockArticleStore{db: mockDB}

	t.Run("DeleteExistingArticle_Success", func(t *testing.T) {

		mockArticle := &model.Article{ID: 1, Title: "Test Article", Description: "Test Description", Body: "Test Body"}

		err := mockArticleStore.Delete(mockArticle)

		assert.NoError(t, err, "Error deleting existing article")
		t.Log("DeleteExistingArticle_Success: Article deleted successfully")
	})

	t.Run("DeleteNonExistingArticle_Success", func(t *testing.T) {

		mockArticle := &model.Article{ID: 2, Title: "Non-existing Article", Description: "Non-existing Description", Body: "Non-existing Body"}

		err := mockArticleStore.Delete(mockArticle)

		assert.NoError(t, err, "Error deleting non-existing article")
		t.Log("DeleteNonExistingArticle_Success: No error returned for non-existing article")
	})

	t.Run("DeleteArticle_Error", func(t *testing.T) {

		mockArticle := &model.Article{ID: 3, Title: "Error Article", Description: "Error Description", Body: "Error Body"}
		mockDB.Close()

		err := mockArticleStore.Delete(mockArticle)

		assert.Error(t, err, "Expected error deleting article")
		t.Log("DeleteArticle_Error: Error encountered during deletion")
	})
}
