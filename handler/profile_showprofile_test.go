package handler

import (
	"context"
	"errors"
	"testing"
	"github.com/raahii/golang-grpc-realworld-example/auth"
	"github.com/raahii/golang-grpc-realworld-example/handler"
	"github.com/raahii/golang-grpc-realworld-example/model"
	"github.com/raahii/golang-grpc-realworld-example/proto"
	"github.com/stretchr/testify/assert"
)



type User struct {
	gorm.Model
	Username         string    `gorm:"unique_index;not null"`
	Email            string    `gorm:"unique_index;not null"`
	Password         string    `gorm:"not null"`
	Bio              string    `gorm:"not null"`
	Image            string    `gorm:"not null"`
	Follows          []User    `gorm:"many2many:follows;jointable_foreignkey:from_user_id;association_jointable_foreignkey:to_user_id"`
	FavoriteArticles []Article `gorm:"many2many:favorite_articles;"`
}


type T struct {
	common
	isEnvSet bool
	context  *testContext // For running tests and subtests.
}

type MockUserStore struct {
	UserData map[uint]*model.User
}



type User struct {
	gorm.Model
	Username         string    `gorm:"unique_index;not null"`
	Email            string    `gorm:"unique_index;not null"`
	Password         string    `gorm:"not null"`
	Bio              string    `gorm:"not null"`
	Image            string    `gorm:"not null"`
	Follows          []User    `gorm:"many2many:follows;jointable_foreignkey:from_user_id;association_jointable_foreignkey:to_user_id"`
	FavoriteArticles []Article `gorm:"many2many:favorite_articles;"`
}


type T struct {
	common
	isEnvSet bool
	context  *testContext // For running tests and subtests.
}

func (m *MockUserStore) GetByID(id uint) (*model.User, error) {
	user, ok := m.UserData[id]
	if !ok {
		return nil, errors.New("user not found")
	}
	return user, nil
}
func (m *MockUserStore) GetByUsername(username string) (*model.User, error) {
	for _, user := range m.UserData {
		if user.Username == username {
			return user, nil
		}
	}
	return nil, errors.New("user not found")
}
func (m *MockUserStore) IsFollowing(a *model.User, b *model.User) (bool, error) {

	return a.ID == b.ID, nil
}
func TestHandlerShowProfile(t *testing.T) {

	mockUserStore := &MockUserStore{
		UserData: map[uint]*model.User{
			1: {ID: 1, Username: "user1", Bio: "Bio1", Image: "image1"},
			2: {ID: 2, Username: "user2", Bio: "Bio2", Image: "image2"},
		},
	}

	h := &handler.Handler{
		us: mockUserStore,
	}

	ctx := context.Background()
	validRequest := &proto.ShowProfileRequest{Username: "user2"}
	invalidRequest := &proto.ShowProfileRequest{Username: "user3"}

	t.Run("Successful profile retrieval for an existing user", func(t *testing.T) {
		res, err := h.ShowProfile(ctx, validRequest)
		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, "user2", res.Profile.Username)
		assert.Equal(t, "Bio2", res.Profile.Bio)
		assert.Equal(t, "image2", res.Profile.Image)
		assert.True(t, res.Profile.Following)
		t.Log("Successful profile retrieval test passed")
	})

	t.Run("Unauthenticated request handling", func(t *testing.T) {
		res, err := h.ShowProfile(ctx, invalidRequest)
		assert.Error(t, err)
		assert.Nil(t, res)
		assert.Equal(t, "unauthenticated", err.Error())
		t.Log("Unauthenticated request test passed")
	})

	t.Run("Requested user not found", func(t *testing.T) {
		res, err := h.ShowProfile(ctx, invalidRequest)
		assert.Error(t, err)
		assert.Nil(t, res)
		assert.Equal(t, "user was not found", err.Error())
		t.Log("Requested user not found test passed")
	})

	t.Run("Internal server error during following check", func(t *testing.T) {

		mockUserStore.IsFollowing = func(a *model.User, b *model.User) (bool, error) {
			return false, errors.New("internal server error")
		}

		res, err := h.ShowProfile(ctx, validRequest)
		assert.Error(t, err)
		assert.Nil(t, res)
		assert.Equal(t, "internal server error", err.Error())
		t.Log("Internal server error during following check test passed")
	})

	t.Run("Boundary case - Empty username in the request", func(t *testing.T) {
		emptyRequest := &proto.ShowProfileRequest{Username: ""}
		res, err := h.ShowProfile(ctx, emptyRequest)
		assert.Error(t, err)
		assert.Nil(t, res)
		assert.Equal(t, "user was not found", err.Error())
		t.Log("Boundary case - Empty username in the request test passed")
	})
}
