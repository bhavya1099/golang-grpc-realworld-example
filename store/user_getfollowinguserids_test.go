package store

import (
	"errors"
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/jinzhu/gorm"
	"github.com/raahii/golang-grpc-realworld-example/model"
	"github.com/raahii/golang-grpc-realworld-example/store"
)

type mockDB struct {
	mock.Mock
}func (m *mockDB) Rows() *gorm.Rows {
	args := m.Called()
	return args.Get(0).(*gorm.Rows)
}
func (m *mockDB) Select(query interface{}, args ...interface{}) *gorm.DB {

	return nil
}
func (m *mockDB) Table(name string) *gorm.DB {
	args := m.Called(name)
	return args.Get(0).(*gorm.DB)
}
func TestUserStoreGetFollowingUserIDs(t *testing.T) {
	t.Parallel()

	t.Run("GetFollowingUserIDs_ExistingFollowers", func(t *testing.T) {
		mockDBInstance := new(mockDB)
		userStore := store.UserStore{db: mockDBInstance}

		mockDBInstance.On("Table", "follows").Return(mockDBInstance)
		mockDBInstance.On("Rows").Return(&gorm.Rows{})

		expectedIDs := []uint{1, 2, 3}
		mockDBInstance.On("Scan", mock.Anything).Return(nil).Run(func(args mock.Arguments) {
			rows := args.Get(0).(*gorm.Rows)
			for _, id := range expectedIDs {
				rows.Scan(&id)
			}
		})

		user := &model.User{ID: 123}
		ids, err := userStore.GetFollowingUserIDs(user)

		assert.NoError(t, err)
		assert.ElementsMatch(t, expectedIDs, ids)
	})

	t.Run("GetFollowingUserIDs_NoFollowers", func(t *testing.T) {
		mockDBInstance := new(mockDB)
		userStore := store.UserStore{db: mockDBInstance}

		mockDBInstance.On("Table", "follows").Return(mockDBInstance)
		mockDBInstance.On("Rows").Return(&gorm.Rows{})

		user := &model.User{ID: 456}
		ids, err := userStore.GetFollowingUserIDs(user)

		assert.NoError(t, err)
		assert.Empty(t, ids)
	})

	t.Run("GetFollowingUserIDs_ErrorFetchingData", func(t *testing.T) {
		mockDBInstance := new(mockDB)
		userStore := store.UserStore{db: mockDBInstance}

		mockDBInstance.On("Table", "follows").Return(mockDBInstance)
		mockDBInstance.On("Rows").Return(nil, errors.New("database error"))

		user := &model.User{ID: 789}
		ids, err := userStore.GetFollowingUserIDs(user)

		assert.Error(t, err)
		assert.Empty(t, ids)
	})

	t.Run("GetFollowingUserIDs_EmptyUser", func(t *testing.T) {
		mockDBInstance := new(mockDB)
		userStore := store.UserStore{db: mockDBInstance}

		user := &model.User{}
		ids, err := userStore.GetFollowingUserIDs(user)

		assert.Error(t, err)
		assert.Empty(t, ids)
	})
}
func (m *mockDB) Where(query interface{}, args ...interface{}) *gorm.DB {

	return nil
}
