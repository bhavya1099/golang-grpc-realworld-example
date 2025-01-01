package store

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/jinzhu/gorm"
	"github.com/raahii/golang-grpc-realworld-example/store"
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
func TestNewArticleStore(t *testing.T) {
	t.Parallel()

	t.Run("NewArticleStore_WithValidDBInstance", func(t *testing.T) {

		validDB := &gorm.DB{}

		articleStore := store.NewArticleStore(validDB)

		assert.NotNil(t, articleStore)
		assert.Equal(t, validDB, articleStore.DB, "DB instance match")
		t.Log("Scenario 1: NewArticleStore_WithValidDBInstance - Test Passed")
	})

	t.Run("NewArticleStore_WithNilDBInstance", func(t *testing.T) {

		var nilDB *gorm.DB

		articleStore := store.NewArticleStore(nilDB)

		assert.Nil(t, articleStore, "ArticleStore instance is nil")
		t.Log("Scenario 2: NewArticleStore_WithNilDBInstance - Test Passed")
	})

	t.Run("NewArticleStore_WithEmptyDBInstance", func(t *testing.T) {

		emptyDB := &gorm.DB{}

		articleStore := store.NewArticleStore(emptyDB)

		assert.NotNil(t, articleStore)
		assert.Equal(t, emptyDB, articleStore.DB, "DB instance match")
		t.Log("Scenario 3: NewArticleStore_WithEmptyDBInstance - Test Passed")
	})

	t.Run("NewArticleStore_WithPopulatedDBInstance", func(t *testing.T) {

		populatedDB := &gorm.DB{Value: "Data"}

		articleStore := store.NewArticleStore(populatedDB)

		assert.NotNil(t, articleStore)
		assert.Equal(t, populatedDB, articleStore.DB, "DB instance match")
		t.Log("Scenario 4: NewArticleStore_WithPopulatedDBInstance - Test Passed")
	})
}
