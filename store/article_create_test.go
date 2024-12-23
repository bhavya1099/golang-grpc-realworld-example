package store

import (
	"testing"
	"github.com/jinzhu/gorm"
	"github.com/raahii/golang-grpc-realworld-example/model"
	"github.com/yourpackage/store"
)

func TestArticleStoreCreate(t *testing.T) {

	t.Run("Create new article successfully", func(t *testing.T) {

		article := &model.Article{
			Title:       "Test Title",
			Description: "Test Description",
			Body:        "Test Body",
			UserID:      1,
		}

		err := store.ArticleStore{}.Create(article)

		if err != nil {
			t.Errorf("Scenario 1: Expected no error, got %v", err)
		}
	})

	t.Run("Error handling when creating article fails", func(t *testing.T) {

		article := &model.Article{}

		err := store.ArticleStore{}.Create(article)

		if err == nil {
			t.Error("Scenario 2: Expected an error, got nil")
		}
	})

	t.Run("Create article with missing required fields", func(t *testing.T) {

		article := &model.Article{
			Title: "Test Title",
		}

		err := store.ArticleStore{}.Create(article)

		if err == nil {
			t.Error("Scenario 3: Expected an error, got nil")
		}
	})

	t.Run("Create article with duplicate unique constraints", func(t *testing.T) {

		article := &model.Article{
			Title:       "Test Title",
			Description: "Test Description",
			Body:        "Test Body",
			UserID:      1,
		}

		err := store.ArticleStore{}.Create(article)

		if err == nil {
			t.Error("Scenario 4: Expected an error, got nil")
		}
	})

	t.Run("Create article with maximum field lengths", func(t *testing.T) {

		article := &model.Article{

			Title:       "Lorem ipsum dolor sit amet, consectetur adipiscing elit.",
			Description: "Lorem ipsum dolor sit amet, consectetur adipiscing elit.",
			Body:        "Lorem ipsum dolor sit amet, consectetur adipiscing elit.",
			UserID:      1,
		}

		err := store.ArticleStore{}.Create(article)

		if err != nil {
			t.Errorf("Scenario 5: Expected no error, got %v", err)
		}
	})
}
