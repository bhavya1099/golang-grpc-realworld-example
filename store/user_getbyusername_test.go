package store

import (
	"errors"
	"testing"

	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
)

type mockDB struct{}

func (m *mockDB) First(out interface{}, where ...interface{}) *gorm.DB {
	return &gorm.DB{}
}

func TestUserStoreGetByUsername(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		username string
		mockDB   *mockDB
		expected *model.User
		errMsg   string
	}{
		{
			name:     "GetByUsername_ValidUsername",
			username: "validUser",
			mockDB:   &mockDB{},
			expected: &model.User{Username: "validUser"},
			errMsg:   "",
		},
		{
			name:     "GetByUsername_EmptyUsername",
			username: "",
			mockDB:   &mockDB{},
			expected: nil,
			errMsg:   "record not found",
		},
		{
			name:     "GetByUsername_NonExistentUsername",
			username: "nonExistentUser",
			mockDB:   &mockDB{},
			expected: nil,
			errMsg:   "record not found",
		},
		{
			name:     "GetByUsername_DBError",
			username: "errorUser",
			mockDB:   &mockDB{},
			expected: nil,
			errMsg:   "database error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userStore := UserStore{db: tt.mockDB}
			user, err := userStore.GetByUsername(tt.username)

			if tt.errMsg != "" {
				assert.Error(t, err)
				assert.Nil(t, user)
				assert.Equal(t, tt.errMsg, err.Error())
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, user)
				assert.Equal(t, tt.expected, user)
			}
		})
	}
}
