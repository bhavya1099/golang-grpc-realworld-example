package store_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/raahii/golang-grpc-realworld-example/model"
	"github.com/raahii/golang-grpc-realworld-example/store"
)

type mockDB struct {
	mock.Mock
}

func (m *mockDB) Find(out interface{}, where ...interface{}) *gorm.DB {
	args := m.Called(out, where...)
	return args.Get(0).(*gorm.DB)
}

func TestUserStoreGetByID(t *testing.T) {
	t.Parallel()

	// Scenario 1: Valid User ID Returns User Details Successfully
	t.Run("Valid User ID Returns User Details Successfully", func(t *testing.T) {
		expectedUser := &model.User{
			Model:    gorm.Model{ID: 1},
			Username: "testuser",
			Email:    "test@example.com",
			Password: "password",
			Bio:      "Test Bio",
			Image:    "test.jpg",
		}

		mockDB := &mockDB{}
		mockDB.On("Find", &model.User{}, uint(1)).Return(nil)

		userStore := store.UserStore{db: mockDB}
		user, err := userStore.GetByID(1)

		assert.NoError(t, err)
		assert.Equal(t, expectedUser, user)
	})

	// Scenario 2: Invalid User ID Returns Error
	t.Run("Invalid User ID Returns Error", func(t *testing.T) {
		mockDB := &mockDB{}
		mockDB.On("Find", &model.User{}, uint(2)).Return(errors.New("record not found"))

		userStore := store.UserStore{db: mockDB}
		user, err := userStore.GetByID(2)

		assert.Error(t, err)
		assert.Nil(t, user)
	})

	// Scenario 3: Error Handling for Database Query Failure
	t.Run("Error Handling for Database Query Failure", func(t *testing.T) {
		mockDB := &mockDB{}
		mockDB.On("Find", &model.User{}, uint(3)).Return(errors.New("database error"))

		userStore := store.UserStore{db: mockDB}
		user, err := userStore.GetByID(3)

		assert.Error(t, err)
		assert.Nil(t, user)
	})

	// Scenario 4: Edge Case - Maximum User ID Value
	t.Run("Edge Case - Maximum User ID Value", func(t *testing.T) {
		expectedUser := &model.User{
			Model:    gorm.Model{ID: ^uint(0)},
			Username: "maxuser",
			Email:    "max@example.com",
			Password: "maxpassword",
			Bio:      "Max Bio",
			Image:    "max.jpg",
		}

		mockDB := &mockDB{}
		mockDB.On("Find", &model.User{}, ^uint(0)).Return(nil)

		userStore := store.UserStore{db: mockDB}
		user, err := userStore.GetByID(^uint(0))

		assert.NoError(t, err)
		assert.Equal(t, expectedUser, user)
	})
}
