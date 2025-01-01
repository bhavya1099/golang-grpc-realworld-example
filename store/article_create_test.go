package store

import (
	"errors"
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/raahii/golang-grpc-realworld-example/model"
	"github.com/raahii/golang-grpc-realworld-example/store"
)



type Call struct {
	Parent *Mock

	// The name of the method that was or will be called.
	Method string

	// Holds the arguments of the method.
	Arguments Arguments

	// Holds the arguments that should be returned when
	// this method is called.
	ReturnArguments Arguments

	// Holds the caller info for the On() call
	callerInfo []string

	// The number of times to return the return arguments when setting
	// expectations. 0 means to always return the value.
	Repeatability int

	// Amount of times this call has been called
	totalCalls int

	// Call to this method can be optional
	optional bool

	// Holds a channel that will be used to block the Return until it either
	// receives a message or is closed. nil means it returns immediately.
	WaitFor <-chan time.Time

	waitTime time.Duration

	// Holds a handler used to manipulate arguments content that are passed by
	// reference. It's useful when mocking methods such as unmarshalers or
	// decoders.
	RunFn func(Arguments)

	// PanicMsg holds msg to be used to mock panic on the function call
	//  if the PanicMsg is set to a non nil string the function call will panic
	// irrespective of other settings
	PanicMsg *string

	// Calls which must be satisfied before this call can be
	requires []*Call
}

type T struct {
	common
	isEnvSet bool
	context  *testContext // For running tests and subtests.
}


type MockDB struct {
	mock.Mock
}



type T struct {
	common
	isEnvSet bool
	context  *testContext // For running tests and subtests.
}
func (m *MockDB) Create(value interface{}) *MockDB {
	args := m.Called(value)
	return args.Get(0).(*MockDB)
}
func TestArticleStoreCreate(t *testing.T) {
	t.Parallel()

	t.Run("Create a new article successfully", func(t *testing.T) {
		mockDB := new(MockDB)
		mockArticle := &model.Article{
			Title:       "Test Title",
			Description: "Test Description",
			Body:        "Test Body",
			UserID:      1,
		}

		mockDB.On("Create", mockArticle).Return(&MockDB{}, nil)

		articleStore := store.ArticleStore{db: mockDB}
		err := articleStore.Create(mockArticle)

		assert.NoError(t, err, "Error should be nil for successful article creation")
	})

	t.Run("Create an article with missing mandatory fields", func(t *testing.T) {
		mockDB := new(MockDB)
		mockArticle := &model.Article{
			Title: "Test Title",
		}

		articleStore := store.ArticleStore{db: mockDB}
		err := articleStore.Create(mockArticle)

		assert.Error(t, err, "Error should be returned for missing mandatory fields")
		assert.EqualError(t, err, "missing mandatory fields: Description, Body, UserID")
	})

	t.Run("Create an article with duplicate title", func(t *testing.T) {
		mockDB := new(MockDB)
		mockArticle := &model.Article{
			Title:       "Existing Title",
			Description: "Test Description",
			Body:        "Test Body",
			UserID:      1,
		}

		mockDB.On("Create", mockArticle).Return(&MockDB{}, errors.New("duplicate key violation"))

		articleStore := store.ArticleStore{db: mockDB}
		err := articleStore.Create(mockArticle)

		assert.Error(t, err, "Error should be returned for duplicate title")
		assert.EqualError(t, err, "article with title 'Existing Title' already exists")
	})

	t.Run("Create an article with excessively long fields", func(t *testing.T) {
		mockDB := new(MockDB)
		mockArticle := &model.Article{
			Title:       "Very Long Title Exceeding Maximum Length Allowed",
			Description: "Test Description",
			Body:        "Test Body",
			UserID:      1,
		}

		articleStore := store.ArticleStore{db: mockDB}
		err := articleStore.Create(mockArticle)

		assert.Error(t, err, "Error should be returned for excessively long fields")
		assert.EqualError(t, err, "field length exceeds maximum allowed")
	})

	t.Run("Create an article with invalid author ID", func(t *testing.T) {
		mockDB := new(MockDB)
		mockArticle := &model.Article{
			Title:       "Test Title",
			Description: "Test Description",
			Body:        "Test Body",
			UserID:      999,
		}

		mockDB.On("Create", mockArticle).Return(&MockDB{}, errors.New("invalid author ID"))

		articleStore := store.ArticleStore{db: mockDB}
		err := articleStore.Create(mockArticle)

		assert.Error(t, err, "Error should be returned for invalid author ID")
		assert.EqualError(t, err, "author with ID 999 does not exist")
	})
}
