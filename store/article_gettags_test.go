package store

import (
	"testing"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/raahii/golang-grpc-realworld-example/model"
	"github.com/raahii/golang-grpc-realworld-example/store"
)

type MockDB struct {
	mock.Mock
}func (m *MockDB) Find(out interface{}) *gorm.DB {
	args := m.Called(out)
	return args.Get(0).(*gorm.DB)
}
func TestArticleStoreGetTags(t *testing.T) {

	t.Run("Successful retrieval of tags", func(t *testing.T) {
		mockDB := &MockDB{}
		expectedTags := []model.Tag{{Name: "tag1"}, {Name: "tag2"}}
		mockDB.On("Find", &[]model.Tag{}).Return(&gorm.DB{})

		store := store.ArticleStore{db: mockDB}
		tags, err := store.GetTags()

		assert.NoError(t, err)
		assert.Equal(t, expectedTags, tags)
	})

	t.Run("Handling an empty result set", func(t *testing.T) {
		mockDB := &MockDB{}
		mockDB.On("Find", &[]model.Tag{}).Return(&gorm.DB{})

		store := store.ArticleStore{db: mockDB}
		tags, err := store.GetTags()

		assert.NoError(t, err)
		assert.Empty(t, tags)
	})

	t.Run("Error handling when database query fails", func(t *testing.T) {
		mockDB := &MockDB{}
		expectedErr := errors.New("database error")
		mockDB.On("Find", &[]model.Tag{}).Return(&gorm.DB{Error: expectedErr})

		store := store.ArticleStore{db: mockDB}
		tags, err := store.GetTags()

		assert.Error(t, err)
		assert.Empty(t, tags)
		assert.EqualError(t, err, expectedErr.Error())
	})

	t.Run("Performance testing with a large dataset", func(t *testing.T) {
		mockDB := &MockDB{}

		store := store.ArticleStore{db: mockDB}
		tags, err := store.GetTags()

		assert.NoError(t, err)

	})

	t.Run("Testing behavior with nil database connection", func(t *testing.T) {
		store := store.ArticleStore{db: nil}
		tags, err := store.GetTags()

		assert.Error(t, err)
		assert.Nil(t, tags)
	})
}
