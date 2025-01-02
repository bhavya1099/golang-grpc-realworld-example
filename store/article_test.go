package store

import (
	"testing"
	"time"
	"bytes"
	"errors"
	"fmt"
	"os"
	"github.com/stretchr/testify/assert"
	"github.com/jinzhu/gorm"
	"github.com/raahii/golang-grpc-realworld-example/model"
)

type common struct{}

type testContext struct{}

type Location struct{}

type readOp int

type T struct {
	common
	isEnvSet bool
	context  *testContext
}

type Time struct {
	wall uint64
	ext  int64
	loc  *Location
}

type Buffer struct {
	buf      []byte
	off      int
	lastRead readOp
}

func TestNewArticleStore(t *testing.T) {
	// Remaining test functions...
}

func TestCreateComment(t *testing.T) {
	var buf bytes.Buffer
	os.Stdout = &buf

	tests := []struct {
		name     string
		comment  *model.Comment
		mockFunc func(comment *model.Comment) *gorm.DB
		expErr   error
	}{
		{
			name: "Successfully Create a Comment",
			comment: &model.Comment{
				Body:      "Great article!",
				UserID:    1,
				ArticleID: 123,
			},
			mockFunc: func(comment *model.Comment) *gorm.DB {
				return &gorm.DB{}
			},
			expErr: nil,
		},
		{
			name: "Error Creating Comment Due to Empty Body",
			comment: &model.Comment{
				Body:      "",
				UserID:    1,
				ArticleID: 123,
			},
			mockFunc: func(comment *model.Comment) *gorm.DB {
				return &gorm.DB{}
			},
			expErr: errors.New("comment body is empty"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := tt.mockFunc(tt.comment)
			articleStore := &ArticleStore{db: mockDB}

			err := articleStore.CreateComment(tt.comment)

			if tt.expErr != nil {
				assert.NotNil(t, err)
				assert.EqualError(t, err, tt.expErr.Error())
			} else {
				assert.Nil(t, err)
			}
		})
		t.Log(fmt.Sprintf("Test Case: %s - Log: %s", tt.name, buf.String()))
		buf.Reset()
	}
}
