package store

import (
	"testing"
	"time"
	"github.com/jinzhu/gorm"
	"github.com/raahii/golang-grpc-realworld-example/model"
	"github.com/raahii/golang-grpc-realworld-example/store"
)

func TestUserStoreGetByEmail(t *testing.T) {

	db, _ := gorm.Open("sqlite3", ":memory:")
	defer db.Close()
	db.AutoMigrate(&model.User{})

	userStore := store.UserStore{db: db}

	t.Run("Valid email provided, user exists in the database", func(t *testing.T) {
		email := "test@example.com"
		expectedUser := &model.User{Email: email}
		db.Create(expectedUser)

		user, err := userStore.GetByEmail(email)

		if err != nil {
			t.Errorf("Scenario 1: Expected no error, got %v", err)
		}

		if user.Email != email {
			t.Errorf("Scenario 1: Expected user with email %s, got %s", email, user.Email)
		}
	})

	t.Run("Valid email provided, user does not exist in the database", func(t *testing.T) {
		email := "nonexistent@example.com"

		_, err := userStore.GetByEmail(email)

		if err == nil {
			t.Error("Scenario 2: Expected error, got nil")
		}
	})

	t.Run("Empty email string provided", func(t *testing.T) {
		_, err := userStore.GetByEmail("")

		if err == nil {
			t.Error("Scenario 3: Expected error, got nil")
		}
	})

	t.Run("Database query error occurs", func(t *testing.T) {

		db.Close()

		_, err := userStore.GetByEmail("test@example.com")

		if err == nil {
			t.Error("Scenario 4: Expected error, got nil")
		}
	})

	t.Run("Performance testing with multiple user records", func(t *testing.T) {

		for i := 0; i < 1000; i++ {
			db.Create(&model.User{
				Email:  "user" + string(i) + "@example.com",
				Name:   "User" + string(i),
				Active: true,
				Joined: time.Now(),
			})
		}

		start := time.Now()
		_, _ = userStore.GetByEmail("user500@example.com")
		elapsed := time.Since(start).Seconds()

		t.Logf("Scenario 5: Time taken for retrieval: %f seconds", elapsed)
	})

	t.Run("Case sensitivity in email search", func(t *testing.T) {
		email := "CaseSensitive@example.com"
		expectedUser := &model.User{Email: email}
		db.Create(expectedUser)

		testCases := []struct {
			email    string
			expected *model.User
		}{
			{"casesensitive@example.com", expectedUser},
			{"CASESENSITIVE@example.com", expectedUser},
			{"CaseSensitive@EXAMPLE.com", expectedUser},
		}

		for _, tc := range testCases {
			user, err := userStore.GetByEmail(tc.email)

			if err != nil {
				t.Errorf("Scenario 6: Expected no error, got %v", err)
			}

			if user.Email != tc.expected.Email {
				t.Errorf("Scenario 6: Expected user with email %s, got %s", tc.expected.Email, user.Email)
			}
		}
	})
}
