package store

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/jinzhu/gorm"
	"github.com/raahii/golang-grpc-realworld-example/model"
	"github.com/raahii/golang-grpc-realworld-example/store"
)

func TestArticleStoreGetByID(t *testing.T) {

	mockDB, _, _ := gorm.Open("sqlite3", ":memory:")
	defer mockDB.Close()

	articleStore := store.ArticleStore{db: mockDB}

	t.Run("Successful retrieval of an article by ID", func(t *testing.T) {

		expectedArticle := &model.Article{ID: 1, Title: "Test Article", Content: "Sample content"}
		mockDB.Create(expectedArticle)

		retrievedArticle, err := articleStore.GetByID(1)

		assert.NoError(t, err)
		assert.NotNil(t, retrievedArticle)
		assert.Equal(t, expectedArticle.ID, retrievedArticle.ID)
		assert.Equal(t, expectedArticle.Title, retrievedArticle.Title)
		assert.Equal(t, expectedArticle.Content, retrievedArticle.Content)
	})

	t.Run("Article not found for a non-existent ID", func(t *testing.T) {

		_, err := articleStore.GetByID(999)

		assert.Error(t, err)
		assert.EqualError(t, err, gorm.ErrRecordNotFound.Error())
	})

	t.Run("Error handling when database query fails", func(t *testing.T) {

		mockDB.Close()

		_, err := articleStore.GetByID(1)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "database is closed")
	})

	t.Run("Retrieval of an article with associated tags and author", func(t *testing.T) {

		expectedArticle := &model.Article{ID: 1, Title: "Test Article", Content: "Sample content"}
		mockDB.Create(expectedArticle)
		tag1 := model.Tag{Name: "Tag1"}
		tag2 := model.Tag{Name: "Tag2"}
		mockDB.Create(&tag1)
		mockDB.Create(&tag2)
		mockDB.Model(expectedArticle).Association("Tags").Append(tag1, tag2)

		author := model.User{Name: "John Doe"}
		mockDB.Create(&author)
		expectedArticle.Author = author
		mockDB.Model(expectedArticle).Update("AuthorID", author.ID)

		retrievedArticle, err := articleStore.GetByID(1)

		assert.NoError(t, err)
		assert.NotNil(t, retrievedArticle)
		assert.Equal(t, 2, len(retrievedArticle.Tags))
		assert.Equal(t, tag1.Name, retrievedArticle.Tags[0].Name)
		assert.Equal(t, tag2.Name, retrievedArticle.Tags[1].Name)
		assert.Equal(t, author.Name, retrievedArticle.Author.Name)
	})

}
