package model

import (
	"testing"
	"github.com/go-ozzo/ozzo-validation"
	"github.com/raahii/golang-grpc-realworld-example/model"
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
func TestCommentValidate(t *testing.T) {
	commentEmptyBody := model.Comment{
		Body:      "",
		UserID:    1,
		ArticleID: 1,
	}

	t.Log("Scenario 1: Validate Comment Body Required")
	err := commentEmptyBody.Validate()
	if err == nil {
		t.Error("Expected validation error for empty Body field, but got nil")
	}

	commentValidBody := model.Comment{
		Body:      "This is a valid comment body",
		UserID:    1,
		ArticleID: 1,
	}

	t.Log("Scenario 2: Validate Comment Body Provided")
	err = commentValidBody.Validate()
	if err != nil {
		t.Errorf("Unexpected validation error for valid Body field: %v", err)
	}

	commentMissingUserID := model.Comment{
		Body:      "Comment with missing UserID",
		UserID:    0,
		ArticleID: 1,
	}

	t.Log("Scenario 3: Validate Comment Missing User ID")
	err = commentMissingUserID.Validate()
	if err == nil {
		t.Error("Expected validation error for missing UserID field, but got nil")
	}

	commentWithAuthorAndArticle := model.Comment{
		Body:      "Comment with Author and Article",
		UserID:    1,
		ArticleID: 1,
		Author:    model.User{},
		Article:   model.Article{},
	}

	t.Log("Scenario 4: Validate Comment with Author and Article")
	err = commentWithAuthorAndArticle.Validate()
	if err != nil {
		t.Errorf("Unexpected validation error for Comment with Author and Article: %v", err)
	}

	commentInvalidBodyLength := model.Comment{
		Body:      "This is a very long comment body that exceeds the maximum allowed length",
		UserID:    1,
		ArticleID: 1,
	}

	t.Log("Scenario 5: Validate Comment with Invalid Body Length")
	err = commentInvalidBodyLength.Validate()
	if err == nil {
		t.Error("Expected validation error for invalid Body length, but got nil")
	}
}
