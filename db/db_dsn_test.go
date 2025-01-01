package db

import (
	"errors"
	"os"
	"testing"
	"github.com/stretchr/testify/assert"
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
func Testdsn(t *testing.T) {

	t.Run("Test DSN function with all environment variables set", func(t *testing.T) {

		os.Setenv("DB_HOST", "localhost")
		os.Setenv("DB_USER", "testuser")
		os.Setenv("DB_PASSWORD", "testpassword")
		os.Setenv("DB_NAME", "testdb")
		os.Setenv("DB_PORT", "3306")

		dsn, err := dsn()

		assert.NoError(t, err, "No error expected")
		assert.Equal(t, "testuser:testpassword@(localhost:3306)/testdb?charset=utf8mb4&parseTime=True&loc=Local", dsn)
	})

	t.Run("Test DSN function with missing DB_HOST environment variable", func(t *testing.T) {

		os.Unsetenv("DB_HOST")
		os.Setenv("DB_USER", "testuser")
		os.Setenv("DB_PASSWORD", "testpassword")
		os.Setenv("DB_NAME", "testdb")
		os.Setenv("DB_PORT", "3306")

		_, err := dsn()

		assert.Error(t, err, "$DB_HOST is not set")
	})

	t.Run("Test DSN function with missing DB_USER environment variable", func(t *testing.T) {

		os.Setenv("DB_HOST", "localhost")
		os.Unsetenv("DB_USER")
		os.Setenv("DB_PASSWORD", "testpassword")
		os.Setenv("DB_NAME", "testdb")
		os.Setenv("DB_PORT", "3306")

		_, err := dsn()

		assert.Error(t, err, "$DB_USER is not set")
	})

	t.Run("Test DSN function with missing DB_PASSWORD environment variable", func(t *testing.T) {

		os.Setenv("DB_HOST", "localhost")
		os.Setenv("DB_USER", "testuser")
		os.Unsetenv("DB_PASSWORD")
		os.Setenv("DB_NAME", "testdb")
		os.Setenv("DB_PORT", "3306")

		_, err := dsn()

		assert.Error(t, err, "$DB_PASSWORD is not set")
	})

	t.Run("Test DSN function with missing DB_NAME environment variable", func(t *testing.T) {

		os.Setenv("DB_HOST", "localhost")
		os.Setenv("DB_USER", "testuser")
		os.Setenv("DB_PASSWORD", "testpassword")
		os.Unsetenv("DB_NAME")
		os.Setenv("DB_PORT", "3306")

		_, err := dsn()

		assert.Error(t, err, "$DB_NAME is not set")
	})

	t.Run("Test DSN function with missing DB_PORT environment variable", func(t *testing.T) {

		os.Setenv("DB_HOST", "localhost")
		os.Setenv("DB_USER", "testuser")
		os.Setenv("DB_PASSWORD", "testpassword")
		os.Setenv("DB_NAME", "testdb")
		os.Unsetenv("DB_PORT")

		_, err := dsn()

		assert.Error(t, err, "$DB_PORT is not set")
	})
}
