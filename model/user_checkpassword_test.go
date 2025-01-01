package model

import (
	"testing"
	"golang.org/x/crypto/bcrypt"
	"github.com/stretchr/testify/assert"
	model "your-package-path/model"
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
func TestUserCheckPassword(t *testing.T) {

	t.Run("Valid Password Comparison", func(t *testing.T) {
		user := model.User{
			Username: "testuser",
			Email:    "test@example.com",
			Password: "$2a$10$123456789012345678901u6h2v7xX2W5H9tQ1UHb9Gd3Y2L4n.2NS",
		}
		plainPassword := "password123"
		result := user.CheckPassword(plainPassword)
		assert.True(t, result, "Expected password match")
	})

	t.Run("Invalid Password Comparison", func(t *testing.T) {
		user := model.User{
			Username: "testuser",
			Email:    "test@example.com",
			Password: "$2a$10$123456789012345678901u6h2v7xX2W5H9tQ1UHb9Gd3Y2L4n.2NS",
		}
		plainPassword := "wrongpassword"
		result := user.CheckPassword(plainPassword)
		assert.False(t, result, "Expected password mismatch")
	})

	t.Run("Empty Password Comparison", func(t *testing.T) {
		user := model.User{
			Username: "testuser",
			Email:    "test@example.com",
			Password: "$2a$10$123456789012345678901u6h2v7xX2W5H9tQ1UHb9Gd3Y2L4n.2NS",
		}
		plainPassword := ""
		result := user.CheckPassword(plainPassword)
		assert.False(t, result, "Expected password mismatch for empty password")
	})

	t.Run("Hash Comparison Error", func(t *testing.T) {
		user := model.User{
			Username: "testuser",
			Email:    "test@example.com",
			Password: "$2a$10$123456789012345678901u6h2v7xX2W5H9tQ1UHb9Gd3Y2L4n.2NS",
		}

		plainPassword := "password123"
		result := user.CheckPassword(plainPassword)
		assert.False(t, result, "Expected password mismatch due to hash comparison error")
	})

	t.Run("Boundary Test", func(t *testing.T) {
		user := model.User{
			Username: "testuser",
			Email:    "test@example.com",
			Password: "$2a$10$123456789012345678901u6h2v7xX2W5H9tQ1UHb9Gd3Y2L4n.2NS",
		}

		maxLengthPassword := "Aa1!Aa1!Aa1!Aa1!Aa1!Aa1!Aa1!Aa1!Aa1!Aa1!Aa1!Aa1!Aa1!Aa1!Aa1!"
		result := user.CheckPassword(maxLengthPassword)
		assert.True(t, result, "Expected password match for maximum length password")
	})

	t.Run("Incorrect User Password", func(t *testing.T) {
		user := model.User{
			Username: "testuser",
			Email:    "test@example.com",
			Password: "$2a$10$123456789012345678901u6h2v7xX2W5H9tQ1UHb9Gd3Y2L4n.2NS",
		}
		plainPassword := "differentpassword"
		result := user.CheckPassword(plainPassword)
		assert.False(t, result, "Expected password mismatch for incorrect password")
	})
}
