package auth_test

import (
	"os"
	"testing"
	"time"
	"github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/assert"
	"github.com/your/package/auth" // TODO: Update the import path with the correct package path

	var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

func TestgenerateToken(t *testing.T) {
	
	tests := []struct {
		name         string
		userID       uint
		expiration   time.Duration
		expectError  bool
		expectedFunc func(token string, err error)
	}{
		{
			name:       "Generate Token Successfully",
			userID:     1,
			expiration: time.Hour * 72,
			expectError: false,
			expectedFunc: func(token string, err error) {
				assert.NotNil(t, token, "Token should not be nil")
				assert.NoError(t, err, "Error should be nil")
			},
		},
		{
			name:       "Generate Token with Invalid Expiration Time",
			userID:     2,
			expiration: -time.Hour,
			expectError: true,
			expectedFunc: func(token string, err error) {
				assert.Empty(t, token, "Token should be empty")
				assert.Error(t, err, "Error should not be nil")
			},
		},
		{
			name:       "Generate Token with Empty User ID",
			userID:     0,
			expiration: time.Hour * 24,
			expectError: true,
			expectedFunc: func(token string, err error) {
				assert.Empty(t, token, "Token should be empty")
				assert.Error(t, err, "Error should not be nil")
			},
		},
		{
			name:       "Generate Token Error Handling",
			userID:     3,
			expiration: time.Hour * 48,
			expectError: true,
			expectedFunc: func(token string, err error) {
				assert.Empty(t, token, "Token should be empty")
				assert.Error(t, err, "Error should not be nil")
			},
		},
		{
			name:       "Generate Token with Maximum Expiration Time",
			userID:     4,
			expiration: time.Duration(int64(^uint64(0) >> 1)),
			expectError: false,
			expectedFunc: func(token string, err error) {
				assert.NotNil(t, token, "Token should not be nil")
				assert.NoError(t, err, "Error should be nil")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			now := time.Now()
			expirationTime := now.Add(tt.expiration)

			token, err := auth.generateToken(tt.userID, expirationTime)

			tt.expectedFunc(token, err)
		})
	}
}
