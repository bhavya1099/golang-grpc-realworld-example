package handler

import (
	"context"
	"errors"
	"testing"
	"github.com/raahii/golang-grpc-realworld-example/handler"
	"github.com/raahii/golang-grpc-realworld-example/model"
	pb "github.com/raahii/golang-grpc-realworld-example/proto"
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




type T struct {
	common
	isEnvSet bool
	context  *testContext // For running tests and subtests.
}

func (m *mockUserStore) GetByID(id uint) (*model.User, error) {
	if id == 1 {
		return &model.User{ID: 1, Username: "mockUser", Email: "mock@example.com", Password: "mockPass", Image: "mockImage", Bio: "mockBio"}, nil
	}
	return nil, errors.New("user not found")
}
func TestHandlerUpdateUser(t *testing.T) {

	t.Run("Update User With Valid Input", func(t *testing.T) {

		userStoreMock := &mockUserStore{}
		handler := handler.Handler{us: userStoreMock}
		request := &pb.UpdateUserRequest{
			User: &pb.UpdateUserRequest_User{
				Username: "newUsername",
				Email:    "new@example.com",
				Password: "newPassword",
				Image:    "newImage",
				Bio:      "newBio",
			},
		}

		response, err := handler.UpdateUser(context.Background(), request)

		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, "newUsername", response.User.Username)
		assert.Equal(t, "new@example.com", response.User.Email)
		assert.Equal(t, "newImage", response.User.Image)
		assert.Equal(t, "newBio", response.User.Bio)
	})

	t.Run("Update User With Empty Request", func(t *testing.T) {

		userStoreMock := &mockUserStore{}
		handler := handler.Handler{us: userStoreMock}
		request := &pb.UpdateUserRequest{}

		response, err := handler.UpdateUser(context.Background(), request)

		assert.Error(t, err)
		assert.Nil(t, response)
		assert.EqualError(t, err, "missing user information")
	})

}
func (m *mockUserStore) Update(user *model.User) error {
	if user.ID == 1 {
		return nil
	}
	return errors.New("update failed")
}
