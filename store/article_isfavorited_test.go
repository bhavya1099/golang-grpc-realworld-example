package store

import (
	"errors"
	"testing"
	"github.com/golang/mock/gomock"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
	"github.com/raahii/golang-grpc-realworld-example/model"
	"github.com/raahii/golang-grpc-realworld-example/store"
)

func TestArticleStoreIsFavorited(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := gorm.NewMockDB(ctrl)

	article := &model.Article{ID: 1}
	user := &model.User{ID: 1}

	t.Run("Valid Favorite Article", func(t *testing.T) {
		mockDB.EXPECT().Table("favorite_articles").
			Where("article_id = ? AND user_id = ?", article.ID, user.ID).
			Count(gomock.Any()).Return(nil)

		store := store.ArticleStore{db: mockDB}
		isFavorited, err := store.IsFavorited(article, user)
		assert.NoError(t, err)
		assert.True(t, isFavorited)
	})

	t.Run("Invalid Article and User", func(t *testing.T) {
		store := store.ArticleStore{db: mockDB}
		isFavorited, err := store.IsFavorited(nil, user)
		assert.NoError(t, err)
		assert.False(t, isFavorited)

		isFavorited, err = store.IsFavorited(article, nil)
		assert.NoError(t, err)
		assert.False(t, isFavorited)
	})

	t.Run("Non-Favorited Article", func(t *testing.T) {
		mockDB.EXPECT().Table("favorite_articles").
			Where("article_id = ? AND user_id = ?", article.ID, user.ID).
			Count(gomock.Any()).Return(nil)

		store := store.ArticleStore{db: mockDB}
		isFavorited, err := store.IsFavorited(article, user)
		assert.NoError(t, err)
		assert.False(t, isFavorited)
	})

	t.Run("Error Handling - Database Query Error", func(t *testing.T) {
		mockDB.EXPECT().Table("favorite_articles").
			Where("article_id = ? AND user_id = ?", article.ID, user.ID).
			Count(gomock.Any()).Return(errors.New("database error"))

		store := store.ArticleStore{db: mockDB}
		isFavorited, err := store.IsFavorited(article, user)
		assert.Error(t, err)
		assert.False(t, isFavorited)
	})

	t.Run("Edge Case - Empty Database Result", func(t *testing.T) {
		mockDB.EXPECT().Table("favorite_articles").
			Where("article_id = ? AND user_id = ?", article.ID, user.ID).
			Count(gomock.Any()).Return(nil)

		store := store.ArticleStore{db: mockDB}
		isFavorited, err := store.IsFavorited(article, user)
		assert.NoError(t, err)
		assert.False(t, isFavorited)
	})
}
