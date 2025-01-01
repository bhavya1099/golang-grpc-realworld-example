package model

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/raahii/golang-grpc-realworld-example/model"
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
func TestArticleOverwrite(t *testing.T) {
	tt := []struct {
		name        string
		title       string
		description string
		body        string
	}{
		{
			name:        "Overwrite with Non-Empty Title, Description, and Body",
			title:       "New Title",
			description: "New Description",
			body:        "New Body",
		},
		{
			name:        "Overwrite with Empty Title, Non-Empty Description, and Body",
			title:       "",
			description: "New Description",
			body:        "New Body",
		},
		{
			name:        "Overwrite with Empty Description, Non-Empty Title, and Body",
			title:       "New Title",
			description: "",
			body:        "New Body",
		},
		{
			name:        "Overwrite with Empty Body, Non-Empty Title and Description",
			title:       "New Title",
			description: "New Description",
			body:        "",
		},
		{
			name:        "Overwrite with Empty Title, Description, and Body",
			title:       "",
			description: "",
			body:        "",
		},
		{
			name:        "Overwrite with Empty Title and Non-Empty Description and Body",
			title:       "",
			description: "New Description",
			body:        "New Body",
		},
		{
			name:        "Overwrite with Non-Empty Title and Empty Description and Body",
			title:       "New Title",
			description: "",
			body:        "",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			article := &model.Article{
				Title:       "Initial Title",
				Description: "Initial Description",
				Body:        "Initial Body",
			}

			article.Overwrite(tc.title, tc.description, tc.body)

			assert.Equal(t, tc.title, article.Title, "Title should be updated correctly")
			assert.Equal(t, tc.description, article.Description, "Description should be updated correctly")
			assert.Equal(t, tc.body, article.Body, "Body should be updated correctly")
		})
	}
}
