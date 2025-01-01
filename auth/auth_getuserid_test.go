package auth

import (
	"context"
	"errors"
	"fmt"
	"os"
	"testing"
	"time"
	"github.com/dgrijalva/jwt-go"
	"github.com/grpc-ecosystem/go-grpc-middleware/auth"
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
func TestGetUserID(t *testing.T) {

	t.Run("Successful Token Validation", func(t *testing.T) {
		token := jwt.New(jwt.SigningMethodHS256)
		claims := &claims{
			UserID:    123,
			ExpiresAt: time.Now().Add(time.Hour).Unix(),
		}
		token.Claims = claims
		tokenString, _ := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
		ctx := context.WithValue(context.Background(), "token", tokenString)

		userID, err := auth.GetUserID(ctx)

		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}

		if userID != claims.UserID {
			t.Errorf("Expected user ID %d, got: %d", claims.UserID, userID)
		}
	})

	t.Run("Expired Token Handling", func(t *testing.T) {
		token := jwt.New(jwt.SigningMethodHS256)
		claims := &claims{
			UserID:    456,
			ExpiresAt: time.Now().Add(-time.Hour).Unix(),
		}
		token.Claims = claims
		tokenString, _ := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
		ctx := context.WithValue(context.Background(), "token", tokenString)

		_, err := auth.GetUserID(ctx)

		if err == nil || err.Error() != "token expired" {
			t.Errorf("Expected 'token expired' error, got: %v", err)
		}
	})

	t.Run("Invalid Token Format", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), "token", "invalid_token")

		_, err := auth.GetUserID(ctx)

		if err == nil || err.Error() != "invalid token: couldn't handle this token" {
			t.Errorf("Expected 'invalid token' error, got: %v", err)
		}
	})

	t.Run("Token with Unrecognized Claims", func(t *testing.T) {
		token := jwt.New(jwt.SigningMethodHS256)
		claims := jwt.MapClaims{"user_id": 789}
		token.Claims = claims
		tokenString, _ := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
		ctx := context.WithValue(context.Background(), "token", tokenString)

		_, err := auth.GetUserID(ctx)

		if err == nil || err.Error() != "invalid token: cannot map token to claims" {
			t.Errorf("Expected 'cannot map token to claims' error, got: %v", err)
		}
	})

	t.Run("Token Validation Failure", func(t *testing.T) {
		token := jwt.New(jwt.SigningMethodHS256)
		claims := &claims{
			UserID:    101,
			ExpiresAt: time.Now().Add(time.Hour).Unix(),
		}
		token.Claims = claims
		tokenString, _ := token.SignedString([]byte("incorrect_secret"))
		ctx := context.WithValue(context.Background(), "token", tokenString)

		_, err := auth.GetUserID(ctx)

		if err == nil || err.Error() != "invalid token: couldn't handle this token" {
			t.Errorf("Expected 'couldn't handle this token' error, got: %v", err)
		}
	})

	t.Run("Token Expiration Check", func(t *testing.T) {
		token := jwt.New(jwt.SigningMethodHS256)
		claims := &claims{
			UserID:    999,
			ExpiresAt: time.Now().Add(-time.Hour).Unix(),
		}
		token.Claims = claims
		tokenString, _ := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
		ctx := context.WithValue(context.Background(), "token", tokenString)

		_, err := auth.GetUserID(ctx)

		if err == nil || err.Error() != "token expired" {
			t.Errorf("Expected 'token expired' error, got: %v", err)
		}
	})

	t.Run("Invalid JWT Secret", func(t *testing.T) {
		token := jwt.New(jwt.SigningMethodHS256)
		claims := &claims{
			UserID:    777,
			ExpiresAt: time.Now().Add(time.Hour).Unix(),
		}
		token.Claims = claims
		tokenString, _ := token.SignedString([]byte("invalid_secret"))
		ctx := context.WithValue(context.Background(), "token", tokenString)

		_, err := auth.GetUserID(ctx)

		if err == nil || err.Error() != "invalid token: couldn't handle this token" {
			t.Errorf("Expected 'couldn't handle this token' error, got: %v", err)
		}
	})
}
