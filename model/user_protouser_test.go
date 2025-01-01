package model

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/raahii/golang-grpc-realworld-example/model"
	pb "github.com/raahii/golang-grpc-realworld-example/proto"
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
func TestUserProtoUser(t *testing.T) {

	t.Log("Scenario 1: Successful ProtoUser Conversion")
	mockUser := &model.User{
		Username: "john_doe",
		Email:    "john.doe@example.com",
		Password: "hashed_password",
		Bio:      "Sample bio",
		Image:    "profile.jpg",
	}
	token := "sample_token"
	expectedProtoUser := &pb.User{
		Email:    mockUser.Email,
		Token:    token,
		Username: mockUser.Username,
		Bio:      mockUser.Bio,
		Image:    mockUser.Image,
	}
	assert.Equal(t, expectedProtoUser, mockUser.ProtoUser(token))

	t.Log("Scenario 2: Empty Token ProtoUser Conversion")
	mockUserEmptyToken := &model.User{
		Username: "jane_smith",
		Email:    "jane.smith@example.com",
		Password: "hashed_password",
		Bio:      "Another bio",
		Image:    "image.jpg",
	}
	emptyToken := ""
	expectedProtoUserEmptyToken := &pb.User{
		Email:    mockUserEmptyToken.Email,
		Token:    emptyToken,
		Username: mockUserEmptyToken.Username,
		Bio:      mockUserEmptyToken.Bio,
		Image:    mockUserEmptyToken.Image,
	}
	assert.Equal(t, expectedProtoUserEmptyToken, mockUserEmptyToken.ProtoUser(emptyToken))

	t.Log("Scenario 3: Nil User Input")
	var nilUser *model.User
	token = "test_token"
	assert.Nil(t, nilUser.ProtoUser(token))

	t.Log("Scenario 4: Custom User Fields ProtoUser Conversion")
	customUser := &model.User{
		Username: "custom_user",
		Email:    "custom.user@example.com",
		Password: "hashed_password",
		Bio:      "Custom bio",
		Image:    "custom.jpg",
	}
	token = "custom_token"
	expectedProtoUserCustom := &pb.User{
		Email:    customUser.Email,
		Token:    token,
		Username: customUser.Username,
		Bio:      customUser.Bio,
		Image:    customUser.Image,
	}
	assert.Equal(t, expectedProtoUserCustom, customUser.ProtoUser(token))

	t.Log("Scenario 5: Edge Case - Long Username ProtoUser Conversion")
	longUsernameUser := &model.User{
		Username: "very_long_username_12345678901234567890123456789012345678901234567890",
		Email:    "long.user@example.com",
		Password: "hashed_password",
		Bio:      "Long username bio",
		Image:    "long.jpg",
	}
	token = "long_token"
	expectedProtoUserLongUsername := &pb.User{
		Email:    longUsernameUser.Email,
		Token:    token,
		Username: "very_long_username_12345678901234567890123456789012345678",
		Bio:      longUsernameUser.Bio,
		Image:    longUsernameUser.Image,
	}
	assert.Equal(t, expectedProtoUserLongUsername, longUsernameUser.ProtoUser(token))

	t.Log("Scenario 6: Error Handling - Invalid Token")
	invalidTokenUser := &model.User{
		Username: "invalid_user",
		Email:    "invalid.user@example.com",
		Password: "hashed_password",
		Bio:      "Invalid token bio",
		Image:    "invalid.jpg",
	}
	invalidToken := "invalid@token"
	expectedProtoUserInvalidToken := &pb.User{
		Email:    invalidTokenUser.Email,
		Token:    invalidTokenUser.Token,
		Username: invalidTokenUser.Username,
		Bio:      invalidTokenUser.Bio,
		Image:    invalidTokenUser.Image,
	}
	assert.Equal(t, expectedProtoUserInvalidToken, invalidTokenUser.ProtoUser(invalidToken))

	t.Log("Scenario 7: Edge Case - Empty User Fields")
	emptyFieldsUser := &model.User{
		Username: "",
		Email:    "empty.user@example.com",
		Password: "hashed_password",
		Bio:      "",
		Image:    "",
	}
	token = "empty_token"
	expectedProtoUserEmptyFields := &pb.User{
		Email:    emptyFieldsUser.Email,
		Token:    token,
		Username: "",
		Bio:      "",
		Image:    "",
	}
	assert.Equal(t, expectedProtoUserEmptyFields, emptyFieldsUser.ProtoUser(token))

	t.Log("Scenario 8: Boundary Check - Maximum Field Lengths")
	maxLengthUser := &model.User{
		Username: "max_username_12345678901234567890123456789012345678901234567890",
		Email:    "max.email.123456789012345678901234567890123456789012345678901234567890@example.com",
		Password: "max_password_123456789012345678901234567890123456789012345678901234567890",
		Bio:      "max_bio_123456789012345678901234567890123456789012345678901234567890",
		Image:    "max_image_123456789012345678901234567890123456789012345678901234567890.jpg",
	}
	token = "max_token_123456789012345678901234567890123456789012345678901234567890"
	expectedProtoUserMaxLength := &pb.User{
		Email:    maxLengthUser.Email,
		Token:    token,
		Username: "max_username_12345678901234567890123456789012345678901234567890",
		Bio:      "max_bio_123456789012345678901234567890123456789012345678901234567890",
		Image:    "max_image_123456789012345678901234567890123456789012345678901234567890.jpg",
	}
	assert.Equal(t, expectedProtoUserMaxLength, maxLengthUser.ProtoUser(token))

	t.Log("Scenario 9: Performance Testing")
	numUsers := 1000
	users := make([]*model.User, numUsers)
	commonToken := "common_token"
	for i := 0; i < numUsers; i++ {
		users[i] = &model.User{
			Username: "user" + string(i),
			Email:    "user" + string(i) + "@example.com",
			Password: "hashed_password",
			Bio:      "User bio",
			Image:    "user.jpg",
		}
	}
	for i := 0; i < numUsers; i++ {
		assert.NotNil(t, users[i].ProtoUser(commonToken))
	}
}
