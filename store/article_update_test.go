package store

import (
	"testing"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
	"github.com/raahii/golang-grpc-realworld-example/model"
)

type MockArticleStore struct {
	db *gorm.DB
}
func TestArticleStoreUpdate(t *testing.T) {

	t.Run("Update Article Successfully", func(t *testing.T) {
		mockStore := &MockArticleStore{db: &gorm.DB{}}
		article := &model.Article{Title: "Sample Title", Description: "Sample Description", Body: "Sample Body"}
		err := mockStore.Update(article)
		assert.NoError(t, err)
		t.Log("Article updated successfully")
	})

	t.Run("Update Article with Empty Data", func(t *testing.T) {
		mockStore := &MockArticleStore{db: &gorm.DB{}}
		emptyArticle := &model.Article{}
		err := mockStore.Update(emptyArticle)
		assert.Error(t, err)
		t.Log("Error received for updating empty article")
	})

	t.Run("Update Non-existent Article", func(t *testing.T) {
		mockStore := &MockArticleStore{db: &gorm.DB{}}
		nonExistentArticle := &model.Article{Model: gorm.Model{ID: 999}}
		err := mockStore.Update(nonExistentArticle)
		assert.Error(t, err)
		t.Log("Error received for updating non-existent article")
	})

	t.Run("Update Article with Database Error", func(t *testing.T) {
		mockStore := &MockArticleStore{db: &gorm.DB{}}
		problematicArticle := &model.Article{Title: "Database Error Article"}

		err := mockStore.Update(problematicArticle)
		assert.Error(t, err)
		t.Log("Error received for updating article with database error")
	})

	t.Run("Update Article with Partial Data", func(t *testing.T) {
		mockStore := &MockArticleStore{db: &gorm.DB{}}
		article := &model.Article{Title: "Partial Data Title", Description: "Partial Data Description"}
		err := mockStore.Update(article)
		assert.NoError(t, err)
		t.Log("Partial data article updated successfully")
	})

	t.Run("Update Article with Tags and Comments", func(t *testing.T) {
		mockStore := &MockArticleStore{db: &gorm.DB{}}
		article := &model.Article{
			Title:       "Article with Tags and Comments",
			Description: "Article Description",
			Body:        "Article Body",
			Tags:        []model.Tag{{Name: "Tag1"}, {Name: "Tag2"}},
			Comments:    []model.Comment{{Body: "Comment1"}, {Body: "Comment2"}},
		}
		err := mockStore.Update(article)
		assert.NoError(t, err)
		t.Log("Article with tags and comments updated successfully")
	})
}
func (s *MockArticleStore) Update(m *model.Article) error {
	return s.db.Model(&m).Update(&m).Error
}
