package store

import (
	"errors"
	"testing"
	"github.com/jinzhu/gorm"
	"github.com/raahii/golang-grpc-realworld-example/model"
	"github.com/stretchr/testify/assert"
)

type MockDB struct{}func (m *MockDB) Find(out interface{}, where ...interface{}) *gorm.DB {
	comment := out.(*model.Comment)

	if comment.ID == 1 {
		comment.ID = 1
		comment.Body = "Mock Comment Body"
		comment.UserID = 1
		comment.ArticleID = 1
		return &gorm.DB{}
	}
	return &gorm.DB{Error: errors.New("record not found")}
}
func TestArticleStoreGetCommentByID(t *testing.T) {

	store := &store.ArticleStore{db: &MockDB{}}

	tests := []struct {
		name      string
		commentID uint
		expected  *model.Comment
		err       error
	}{
		{
			name:      "Valid comment ID provided",
			commentID: 1,
			expected:  &model.Comment{ID: 1, Body: "Mock Comment Body", UserID: 1, ArticleID: 1},
			err:       nil,
		},
		{
			name:      "Invalid comment ID provided",
			commentID: 10,
			expected:  nil,
			err:       errors.New("record not found"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			comment, err := store.GetCommentByID(tt.commentID)

			if tt.err != nil {
				assert.NotNil(t, err, "Expected error but got nil")
				assert.EqualError(t, err, tt.err.Error(), "Unexpected error")
			} else {
				assert.Nil(t, err, "Unexpected error")
				assert.NotNil(t, comment, "Expected comment but got nil")
				assert.Equal(t, tt.expected.ID, comment.ID, "Incorrect comment ID")
				assert.Equal(t, tt.expected.Body, comment.Body, "Incorrect comment Body")
				assert.Equal(t, tt.expected.UserID, comment.UserID, "Incorrect comment UserID")
				assert.Equal(t, tt.expected.ArticleID, comment.ArticleID, "Incorrect comment ArticleID")
			}
		})
	}
}
