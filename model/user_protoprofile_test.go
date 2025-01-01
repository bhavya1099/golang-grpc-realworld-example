package model

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/your/package/model"
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
func TestUserProtoProfile(t *testing.T) {
	tt := []struct {
		name      string
		user      model.User
		following bool
		expected  model.Profile
	}{
		{
			name: "User Profile with Following True",
			user: model.User{
				Username: "john_doe",
				Bio:      "Sample bio text.",
				Image:    "profile.jpg",
			},
			following: true,
			expected: model.Profile{
				Username:  "john_doe",
				Bio:       "Sample bio text.",
				Image:     "profile.jpg",
				Following: true,
			},
		},
		{
			name: "User Profile with Following False",
			user: model.User{
				Username: "jane_smith",
				Bio:      "Another bio text.",
				Image:    "avatar.jpg",
			},
			following: false,
			expected: model.Profile{
				Username:  "jane_smith",
				Bio:       "Another bio text.",
				Image:     "avatar.jpg",
				Following: false,
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			profile := tc.user.ProtoProfile(tc.following)
			assert.Equal(t, tc.expected.Username, profile.Username, "Username mismatch")
			assert.Equal(t, tc.expected.Bio, profile.Bio, "Bio mismatch")
			assert.Equal(t, tc.expected.Image, profile.Image, "Image mismatch")
			assert.Equal(t, tc.expected.Following, profile.Following, "Following status mismatch")
		})
	}
}
