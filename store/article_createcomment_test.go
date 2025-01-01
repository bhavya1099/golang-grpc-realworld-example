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
func (m *MockDB) Create(value interface{}) *store.DB {
	args := m.Called(value)
	return args.Get(0).(*store.DB)
}
func TestArticleStoreCreateComment(t *testing.T) {
	t.Run("CreateComment_Success", func(t *testing.T) {

		mockDB := new(MockDB)
		articleStore := store.ArticleStore{db: mockDB}
		comment := &model.Comment{
			Body:      "Test Comment",
			UserID:    1,
			ArticleID: 1,
		}
		mockDB.On("Create", comment).Return(&store.DB{})

		err := articleStore.CreateComment(comment)

		assert.NoError(t, err, "Expected no error for successful comment creation")
	})

	t.Run("CreateComment_EmptyBody", func(t *testing.T) {

		mockDB := new(MockDB)
		articleStore := store.ArticleStore{db: mockDB}
		comment := &model.Comment{
			Body:      "",
			UserID:    1,
			ArticleID: 1,
		}

		err := articleStore.CreateComment(comment)

		assert.Error(t, err, "Expected error for empty comment body")
	})

	t.Run("CreateComment_EmptyUserID", func(t *testing.T) {

		mockDB := new(MockDB)
		articleStore := store.ArticleStore{db: mockDB}
		comment := &model.Comment{
			Body:      "Test Comment",
			UserID:    0,
			ArticleID: 1,
		}

		err := articleStore.CreateComment(comment)

		assert.Error(t, err, "Expected error for empty UserID")
	})

	t.Run("CreateComment_EmptyArticleID", func(t *testing.T) {

		mockDB := new(MockDB)
		articleStore := store.ArticleStore{db: mockDB}
		comment := &model.Comment{
			Body:      "Test Comment",
			UserID:    1,
			ArticleID: 0,
		}

		err := articleStore.CreateComment(comment)

		assert.Error(t, err, "Expected error for empty ArticleID")
	})

	t.Run("CreateComment_DatabaseError", func(t *testing.T) {

		mockDB := new(MockDB)
		articleStore := store.ArticleStore{db: mockDB}
		comment := &model.Comment{
			Body:      "Test Comment",
			UserID:    1,
			ArticleID: 1,
		}
		mockDB.On("Create", comment).Return(errors.New("database error"))

		err := articleStore.CreateComment(comment)

		assert.Error(t, err, "Expected error for database error during comment creation")
	})
}
