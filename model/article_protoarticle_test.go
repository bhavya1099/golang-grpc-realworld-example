package model

import (
	"testing"
	"time"
	"github.com/stretchr/testify/assert"
	model "your-package-path/model"
	pb "your-package-path/proto"
)

type T struct {
	common
	isEnvSet bool
	context  *testContext // For running tests and subtests.
}

type T struct {
	common
	isEnvSet bool
	context  *testContext // For running tests and subtests.
}
func TestArticleProtoArticle(t *testing.T) {

	t.Run("ProtoArticle with Favorited True", func(t *testing.T) {

		article := model.Article{
			Title:          "Sample Title",
			Description:    "Sample Description",
			Body:           "Sample Body",
			FavoritesCount: 10,
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
			Tags:           []model.Tag{{Name: "tag1"}, {Name: "tag2"}},
		}

		protoArticle := article.ProtoArticle(true)

		assert.True(t, protoArticle.Favorited, "Favorited field should be true")
	})

	t.Run("ProtoArticle with Favorited False", func(t *testing.T) {

		article := model.Article{
			Title:          "Sample Title",
			Description:    "Sample Description",
			Body:           "Sample Body",
			FavoritesCount: 5,
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
			Tags:           []model.Tag{{Name: "tag1"}},
		}

		protoArticle := article.ProtoArticle(false)

		assert.False(t, protoArticle.Favorited, "Favorited field should be false")
	})

	t.Run("ProtoArticle with Empty Tags", func(t *testing.T) {

		article := model.Article{
			Title:       "Sample Title",
			Description: "Sample Description",
			Body:        "Sample Body",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		protoArticle := article.ProtoArticle(false)

		assert.Empty(t, protoArticle.TagList, "TagList should be empty")
	})

}
