package store

import (
	"testing"
	"errors"
	"github.com/jinzhu/gorm"
	"github.com/raahii/golang-grpc-realworld-example/model"
)

func TestUserStoreUpdate(t *testing.T) {

	t.Run("Update User Successfully", func(t *testing.T) {
		mockDB := &gorm.DB{}
		userModel := &model.User{
			Username: "testuser",
			Email:    "test@example.com",
			Password: "password",
			Bio:      "Test Bio",
			Image:    "profile.jpg",
		}

		userStore := UserStore{db: mockDB}
		err := userStore.Update(userModel)

		if err != nil {
			t.Errorf("Scenario 1: Expected no error, got %v", err)
		}
	})

	t.Run("Update User with Empty Model", func(t *testing.T) {
		mockDB := &gorm.DB{}
		emptyUserModel := &model.User{}

		userStore := UserStore{db: mockDB}
		err := userStore.Update(emptyUserModel)

		if err == nil {
			t.Error("Scenario 2: Expected error for empty model, got nil")
		}
	})

	t.Run("Update User with Database Error", func(t *testing.T) {
		mockDB, _, _ := gorm.Open(&gorm.Dialect{}, "")
		userModel := &model.User{
			Username: "testuser",
			Email:    "test@example.com",
			Password: "password",
			Bio:      "Test Bio",
			Image:    "profile.jpg",
		}

		userStore := UserStore{db: mockDB}
		err := userStore.Update(userModel)

		if err == nil {
			t.Error("Scenario 3: Expected error due to database error, got nil")
		}
	})

	t.Run("Update User with Nil UserStore", func(t *testing.T) {
		var userStore *UserStore
		userModel := &model.User{
			Username: "testuser",
			Email:    "test@example.com",
			Password: "password",
			Bio:      "Test Bio",
			Image:    "profile.jpg",
		}

		err := userStore.Update(userModel)

		if err == nil {
			t.Error("Scenario 4: Expected error for nil UserStore, got nil")
		}
	})

	t.Run("Update User with Invalid Model Fields", func(t *testing.T) {
		mockDB := &gorm.DB{}
		invalidUserModel := &model.User{}

		userStore := UserStore{db: mockDB}
		err := userStore.Update(invalidUserModel)

		if err == nil {
			t.Error("Scenario 5: Expected error for invalid model fields, got nil")
		}
	})

	t.Run("Update User with Unchanged Model", func(t *testing.T) {
		mockDB := &gorm.DB{}
		unchangedUserModel := &model.User{
			Username: "testuser",
			Email:    "test@example.com",
			Password: "password",
			Bio:      "Test Bio",
			Image:    "profile.jpg",
		}

		userStore := UserStore{db: mockDB}
		err := userStore.Update(unchangedUserModel)

		if err != nil {
			t.Errorf("Scenario 6: Expected no error for unchanged model, got %v", err)
		}
	})

	t.Run("Update User with Model Not Found in Database", func(t *testing.T) {
		mockDB := &gorm.DB{}
		notFoundUserModel := &model.User{
			Username: "nonexistent",
			Email:    "notfound@example.com",
			Password: "password",
			Bio:      "Nonexistent user",
			Image:    "missing.jpg",
		}

		userStore := UserStore{db: mockDB}
		err := userStore.Update(notFoundUserModel)

		if err == nil {
			t.Error("Scenario 7: Expected error for user not found, got nil")
		}
	})
}
