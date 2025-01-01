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
func TestArticleValidate(t *testing.T) {
	t.Parallel()

	tt := []struct {
		name     string
		article  model.Article
		expected error
	}{
		{
			name: "Test Validation Success",
			article: model.Article{
				Title: "Sample Title",
				Body:  "Sample Body",
				Tags:  []model.Tag{{Name: "tag1"}, {Name: "tag2"}},
			},
			expected: nil,
		},
		{
			name: "Test Validation Failure - Missing Title",
			article: model.Article{
				Body: "Sample Body",
				Tags: []model.Tag{{Name: "tag1"}, {Name: "tag2"}},
			},
			expected: validation.NewErrors(validation.Field(&model.Article{}, "Title", validation.Required.Error())),
		},
		{
			name: "Test Validation Failure - Missing Body",
			article: model.Article{
				Title: "Sample Title",
				Tags:  []model.Tag{{Name: "tag1"}, {Name: "tag2"}},
			},
			expected: validation.NewErrors(validation.Field(&model.Article{}, "Body", validation.Required.Error())),
		},
		{
			name: "Test Validation Failure - Missing Tags",
			article: model.Article{
				Title: "Sample Title",
				Body:  "Sample Body",
			},
			expected: validation.NewErrors(validation.Field(&model.Article{}, "Tags", validation.Required.Error())),
		},
		{
			name:    "Test Validation Failure - Multiple Missing Fields",
			article: model.Article{},
			expected: validation.NewErrors(
				validation.Field(&model.Article{}, "Title", validation.Required.Error()),
				validation.Field(&model.Article{}, "Body", validation.Required.Error()),
				validation.Field(&model.Article{}, "Tags", validation.Required.Error()),
			),
		},
		{
			name:    "Test Validation Failure - Empty Article",
			article: model.Article{},
			expected: validation.NewErrors(
				validation.Field(&model.Article{}, "Title", validation.Required.Error()),
				validation.Field(&model.Article{}, "Body", validation.Required.Error()),
				validation.Field(&model.Article{}, "Tags", validation.Required.Error()),
			),
		},
		{
			name:    "Test Validation Failure - Nil Article",
			article: model.Article{},
			expected: validation.NewErrors(
				validation.Field(&model.Article{}, "Title", validation.Required.Error()),
				validation.Field(&model.Article{}, "Body", validation.Required.Error()),
				validation.Field(&model.Article{}, "Tags", validation.Required.Error()),
			),
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.article.Validate()
			if tc.expected != nil {
				assert.Error(t, err, "Expected error but got nil")
				assert.Equal(t, tc.expected.Error(), err.Error(), "Error message mismatch")
			} else {
				assert.NoError(t, err, "Expected no error but got error")
			}
		})
	}
}
