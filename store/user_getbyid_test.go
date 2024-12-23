package store

import (
	"testing"
	"errors"
	"github.com/jinzhu/gorm"
	"github.com/raahii/golang-grpc-realworld-example/model"
)

func TestUserStoreGetByID(t *testing.T) {

	t.Run("GetByID Success Scenario", func(t *testing.T) {

		mockUser := &model.User{ID: 1, Name: "John Doe"}
		mockDB := &gorm.DB{}
		mockDB.On("Find", &model.User{}, mockUser.ID).Return(nil)

		userStore := &UserStore{db: mockDB}

		user, err := userStore.GetByID(mockUser.ID)

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if user.ID != mockUser.ID {
			t.Errorf("Expected user ID to be %d, got %d", mockUser.ID, user.ID)
		}
		if user.Name != mockUser.Name {
			t.Errorf("Expected user Name to be %s, got %s", mockUser.Name, user.Name)
		}
	})

	t.Run("GetByID User Not Found", func(t *testing.T) {

		invalidUserID := uint(999)
		mockDB := &gorm.DB{}
		mockDB.On("Find", &model.User{}, invalidUserID).Return(gorm.ErrRecordNotFound)

		userStore := &UserStore{db: mockDB}

		_, err := userStore.GetByID(invalidUserID)

		if !errors.Is(err, gorm.ErrRecordNotFound) {
			t.Errorf("Expected user not found error, got %v", err)
		}
	})

	t.Run("GetByID Database Error", func(t *testing.T) {

		mockDB := &gorm.DB{}
		mockDB.On("Find").Return(errors.New("database error"))

		userStore := &UserStore{db: mockDB}

		_, err := userStore.GetByID(1)

		if err == nil {
			t.Error("Expected database error, got nil")
		}
	})

	t.Run("GetByID Empty ID", func(t *testing.T) {

		emptyUserID := uint(0)

		userStore := &UserStore{}

		_, err := userStore.GetByID(emptyUserID)

		if err == nil {
			t.Error("Expected error for empty ID, got nil")
		}
	})

	t.Run("GetByID Performance Test", func(t *testing.T) {

		mockDB := &gorm.DB{}
		largeUserID := uint(1000000)
		mockDB.On("Find", &model.User{}, largeUserID).Return(nil)

		userStore := &UserStore{db: mockDB}

		user, err := userStore.GetByID(largeUserID)

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if user.ID != largeUserID {
			t.Errorf("Expected user ID to be %d, got %d", largeUserID, user.ID)
		}
	})
}
