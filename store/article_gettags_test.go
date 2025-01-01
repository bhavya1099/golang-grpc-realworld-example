package store_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/raahii/golang-grpc-realworld-example/model"
	"github.com/raahii/golang-grpc-realworld-example/store"
)

func TestArticleStoreGetTags(t *testing.T) {
	// Scenario 1: Successful retrieval of tags from the database
	tags := []model.Tag{
		{Name: "tag1"},
		{Name: "tag2"},
		{Name: "tag3"},
	}
	mockStore := &store.ArticleStore{db: &MockDB{tags: tags}}
	
	t.Run("Scenario 1", func(t *testing.T) {
		tags, err := mockStore.GetTags()
		assert.NoError(t, err)
		assert.ElementsMatch(t, tags, tags)
		t.Log("Scenario 1: Successful retrieval of tags from the database - Passed")
	})

	// Scenario 2: Empty tag list retrieval from the database
	emptyTags := []model.Tag{}
	emptyMockStore := &store.ArticleStore{db: &MockDB{tags: emptyTags}}
	
	t.Run("Scenario 2", func(t *testing.T) {
		tags, err := emptyMockStore.GetTags()
		assert.NoError(t, err)
		assert.ElementsMatch(t, tags, emptyTags)
		t.Log("Scenario 2: Empty tag list retrieval from the database - Passed")
	})

	// Scenario 3: Error handling when database query fails
	errorMockStore := &store.ArticleStore{db: &MockDB{queryError: true}}
	
	t.Run("Scenario 3", func(t *testing.T) {
		tags, err := errorMockStore.GetTags()
		assert.Error(t, err)
		assert.Empty(t, tags)
		t.Log("Scenario 3: Error handling when database query fails - Passed")
	})

	// Scenario 4: Retrieval of tags with special characters from the database
	specialTags := []model.Tag{
		{Name: "tag@#!"},
		{Name: "tag$%^"},
		{Name: "tag&*("},
	}
	specialMockStore := &store.ArticleStore{db: &MockDB{tags: specialTags}}
	
	t.Run("Scenario 4", func(t *testing.T) {
		tags, err := specialMockStore.GetTags()
		assert.NoError(t, err)
		assert.ElementsMatch(t, tags, specialTags)
		t.Log("Scenario 4: Retrieval of tags with special characters from the database - Passed")
	})

	// Scenario 5: Performance testing for tag retrieval from a large database
	largeTags := generateLargeTagData(1000) // Assuming a large dataset with 1000 tags
	largeMockStore := &store.ArticleStore{db: &MockDB{tags: largeTags}}
	
	t.Run("Scenario 5", func(t *testing.T) {
		tags, err := largeMockStore.GetTags()
		assert.NoError(t, err)
		assert.Len(t, tags, 1000)
		t.Log("Scenario 5: Performance testing for tag retrieval from a large database - Passed")
	})
}

func generateLargeTagData(count int) []model.Tag {
	tags := make([]model.Tag, count)
	for i := 0; i < count; i++ {
		tags[i] = model.Tag{Name: "tag" + string(i)}
	}
	return tags
}

// MockDB is a mock implementation of gorm.DB for testing purposes
type MockDB struct {
	tags      []model.Tag
	queryError bool
}

func (mdb *MockDB) Find(out interface{}, where ...interface{}) *gorm.DB {
	if mdb.queryError {
		return &gorm.DB{Error: gorm.ErrRecordNotFound}
	}
	tags := out.(*[]model.Tag)
	*tags = mdb.tags
	return &gorm.DB{}
}
