package store_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/jinzhu/gorm"
	"github.com/raahii/golang-grpc-realworld-example/store"
)

func TestNewUserStore(t *testing.T) {
	t.Parallel()

	// Scenario 1: NewUserStore_WithValidDB_ReturnsUserStore
	t.Run("NewUserStore_WithValidDB_ReturnsUserStore", func(t *testing.T) {
		// Arrange
		validDB := &gorm.DB{}

		// Act
		userStore := store.NewUserStore(validDB)

		// Assert
		assert.NotNil(t, userStore, "UserStore should not be nil")
		assert.Equal(t, validDB, userStore.DB, "UserStore's DB should match the provided DB")
	})

	// Scenario 2: NewUserStore_WithNilDB_ReturnsNilUserStore
	t.Run("NewUserStore_WithNilDB_ReturnsNilUserStore", func(t *testing.T) {
		// Arrange
		var nilDB *gorm.DB

		// Act
		userStore := store.NewUserStore(nilDB)

		// Assert
		assert.Nil(t, userStore, "UserStore should be nil when provided with nil DB")
	})

	// Scenario 3: NewUserStore_WithEmptyDB_ReturnsUserStoreWithEmptyDB
	t.Run("NewUserStore_WithEmptyDB_ReturnsUserStoreWithEmptyDB", func(t *testing.T) {
		// Arrange
		emptyDB := &gorm.DB{}

		// Act
		userStore := store.NewUserStore(emptyDB)

		// Assert
		assert.NotNil(t, userStore, "UserStore should not be nil")
		assert.Equal(t, emptyDB, userStore.DB, "UserStore's DB should match the provided empty DB")
	})

	// Scenario 4: NewUserStore_WithNonNilDB_ReturnsUserStoreWithSameDB
	t.Run("NewUserStore_WithNonNilDB_ReturnsUserStoreWithSameDB", func(t *testing.T) {
		// Arrange
		predefinedDB := &gorm.DB{Value: "test"}

		// Act
		userStore := store.NewUserStore(predefinedDB)

		// Assert
		assert.NotNil(t, userStore, "UserStore should not be nil")
		assert.Equal(t, predefinedDB, userStore.DB, "UserStore's DB should match the provided predefined DB")
	})

	// Scenario 5: NewUserStore_WithInvalidDB_ReturnsUserStoreWithInvalidDB
	t.Run("NewUserStore_WithInvalidDB_ReturnsUserStoreWithInvalidDB", func(t *testing.T) {
		// Arrange
		invalidDB := &gorm.DB{Value: nil} // Simulating an invalid DB with uninitialized field

		// Act
		userStore := store.NewUserStore(invalidDB)

		// Assert
		assert.NotNil(t, userStore, "UserStore should not be nil")
		assert.Equal(t, invalidDB, userStore.DB, "UserStore's DB should match the provided invalid DB")
	})
}
