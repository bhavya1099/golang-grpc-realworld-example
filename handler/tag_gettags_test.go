package handler

import (
	"context"
	"errors"
	"testing"
	"github.com/raahii/golang-grpc-realworld-example/handler"
	"github.com/raahii/golang-grpc-realworld-example/model"
	"github.com/raahii/golang-grpc-realworld-example/store"
	"github.com/stretchr/testify/assert"
)



type T struct {
	common
	isEnvSet bool
	context  *testContext // For running tests and subtests.
}

type MockArticleStore struct{}


type Tag struct {
	gorm.Model
	Name string `gorm:"not null"`
}

type T struct {
	common
	isEnvSet bool
	context  *testContext // For running tests and subtests.
}
func (mas *MockArticleStore) GetTags() ([]model.Tag, error) {
	tags := []model.Tag{
		{Name: "tag1"},
		{Name: "tag2"},
		{Name: "tag3"},
	}
	return tags, nil
}
func TestHandlerGetTags(t *testing.T) {
	h := &handler.Handler{
		as: &MockArticleStore{},
	}

	t.Run("Successful retrieval of tags", func(t *testing.T) {
		ctx := context.Background()
		req := &pb.Empty{}
		tagsResp, err := h.GetTags(ctx, req)
		assert.NoError(t, err)
		assert.NotNil(t, tagsResp)
		assert.Equal(t, []string{"tag1", "tag2", "tag3"}, tagsResp.Tags)
	})

	t.Run("Empty tag list", func(t *testing.T) {
		h.as = &MockArticleStore{}
		ctx := context.Background()
		req := &pb.Empty{}
		tagsResp, err := h.GetTags(ctx, req)
		assert.NoError(t, err)
		assert.NotNil(t, tagsResp)
		assert.Empty(t, tagsResp.Tags)
	})

	t.Run("Error handling - Failed to get tags", func(t *testing.T) {
		h.as = &MockArticleStore{}
		h.as = &MockArticleStore{err: errors.New("mock error")}
		ctx := context.Background()
		req := &pb.Empty{}
		tagsResp, err := h.GetTags(ctx, req)
		assert.Error(t, err)
		assert.Nil(t, tagsResp)
	})

	t.Run("Tag name extraction", func(t *testing.T) {
		h.as = &MockArticleStore{}
		ctx := context.Background()
		req := &pb.Empty{}
		tagsResp, err := h.GetTags(ctx, req)
		assert.NoError(t, err)
		assert.NotNil(t, tagsResp)
		assert.ElementsMatch(t, []string{"tag1", "tag2", "tag3"}, tagsResp.Tags)
	})

	t.Run("Large number of tags", func(t *testing.T) {

		largeTags := make([]model.Tag, 1000)
		for i := 0; i < 1000; i++ {
			largeTags[i] = model.Tag{Name: "tag" + string(i)}
		}
		h.as = &MockArticleStore{tags: largeTags}

		ctx := context.Background()
		req := &pb.Empty{}
		tagsResp, err := h.GetTags(ctx, req)
		assert.NoError(t, err)
		assert.NotNil(t, tagsResp)
		assert.Len(t, tagsResp.Tags, 1000)
	})
}
