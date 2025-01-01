package store

import (
	"testing"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
	"github.com/raahii/golang-grpc-realworld-example/model"
)

func TestArticleStoreGetByID(t *testing.T) {
	// Mock database connection
	db, _ := gorm.Open("sqlite3", ":memory:")
	defer db.Close()

	// Auto migrate the model for testing
	db.AutoMigrate(&model.Article{})

	// Initialize ArticleStore with the mock DB
	store := ArticleStore{db: db}

	t.Run("GetByID Success", func(t *testing.T) {
		// Arrange
		expectedArticle := model.Article{
			Title:       "Test Article",
			Description: "Test Description",
			Body:        "Test Body",
			UserID:      1,
		}
		db.Create(&expectedArticle)

		// Act
		article, err := store.GetByID(expectedArticle.ID)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, article)
		assert.Equal(t, expectedArticle.Title, article.Title)
		assert.Equal(t, expectedArticle.Description, article.Description)
		assert.Equal(t, expectedArticle.Body, article.Body)
		assert.Equal(t, expectedArticle.UserID, article.UserID)
	})

	t.Run("GetByID Non-Existent ID", func(t *testing.T) {
		// Act
		article, err := store.GetByID(999)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, article)
	})

	t.Run("GetByID Empty Database", func(t *testing.T) {
		// Act
		article, err := store.GetByID(1)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, article)
	})

	t.Run("GetByID Error Handling", func(t *testing.T) {
		// Arrange
		db.Close() // Simulate DB error by closing the connection

		// Act
		article, err := store.GetByID(1)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, article)
	})

	t.Run("GetByID Performance Test", func(t *testing.T) {
		// Arrange
		numArticles := 1000
		for i := 1; i <= numArticles; i++ {
			article := model.Article{
				Title:       "Article Title",
				Description: "Article Description",
				Body:        "Article Body",
				UserID:      1,
				CreatedAt:   time.Now(),
			}
			db.Create(&article)
		}

		// Act
		start := time.Now()
		article, err := store.GetByID(1)
		elapsed := time.Since(start)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, article)
		t.Logf("GetByID Performance Test took %v for %d articles", elapsed, numArticles)
	})
}
