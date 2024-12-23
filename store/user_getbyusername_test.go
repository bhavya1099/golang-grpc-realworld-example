package store

import (
	"errors"
	"testing"
	"github.com/jinzhu/gorm"
	"github.com/raahii/golang-grpc-realworld-example/model"
	"github.com/raahii/golang-grpc-realworld-example/store"
)

func TestUserStoreGetByUsername(t *testing.T) {

	mockUser := model.User{
		Username: "testuser",
		Email:    "testuser@example.com",
		Password: "testpassword",
	}

	mockDB, _ := gorm.Open("sqlite3", ":memory:")
	userStore := store.UserStore{db: mockDB}
	mockDB.AutoMigrate(&model.User{})
	mockDB.Create(&mockUser)

	testCases := []struct {
		name     string
		username string
		expected *model.User
		wantErr  error
	}{
		{
			name:     "Valid Username Provided",
			username: "testuser",
			expected: &mockUser,
			wantErr:  nil,
		},
		{
			name:     "Empty Username Provided",
			username: "",
			expected: nil,
			wantErr:  errors.New("empty username provided"),
		},
		{
			name:     "Username Not Found",
			username: "nonexistentuser",
			expected: nil,
			wantErr:  errors.New("user not found"),
		},
		{
			name:     "Database Error Handling",
			username: "testuser",
			expected: nil,
			wantErr:  errors.New("database error"),
		},
		{
			name:     "Username with Special Characters",
			username: "test$user",
			expected: &mockUser,
			wantErr:  nil,
		},
		{
			name:     "Boundary Testing - Maximum Username Length",
			username: "verylongusernameverylongusernameverylongusernameverylongusername",
			expected: &mockUser,
			wantErr:  nil,
		},
		{
			name:     "Boundary Testing - Minimum Username Length",
			username: "u",
			expected: &mockUser,
			wantErr:  nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			user, err := userStore.GetByUsername(tc.username)

			if err != nil && tc.wantErr != nil {
				if err.Error() != tc.wantErr.Error() {
					t.Errorf("Test case %s failed: %v, expected: %v", tc.name, err, tc.wantErr)
				}
			} else if err != nil || tc.wantErr != nil {
				t.Errorf("Test case %s failed: unexpected error state", tc.name)
			}

			if user != nil && tc.expected != nil {
				if user.Username != tc.expected.Username || user.Email != tc.expected.Email || user.Password != tc.expected.Password {
					t.Errorf("Test case %s failed: retrieved user does not match expected", tc.name)
				}
			} else if user != nil || tc.expected != nil {
				t.Errorf("Test case %s failed: unexpected user data state", tc.name)
			}
		})
	}
}
