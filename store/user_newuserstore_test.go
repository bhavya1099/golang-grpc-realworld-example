package store

import (
	"testing"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
	"your-package-name/store"
)

func TestNewUserStore(t *testing.T) {
	t.Parallel()

	tt := []struct {
		name       string
		dbInstance *gorm.DB
	}{
		{
			name:       "Valid DB Instance",
			dbInstance: &gorm.DB{},
		},
		{
			name:       "Nil DB Instance",
			dbInstance: nil,
		},
		{
			name:       "Empty DB Instance",
			dbInstance: &gorm.DB{},
		},
		{
			name:       "Initialized but No Connection",
			dbInstance: &gorm.DB{Value: "Initialized"},
		},
		{
			name:       "DB Instance with Data",
			dbInstance: &gorm.DB{Value: "Data"},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			t.Log("Scenario:", tc.name)

			userStore := store.NewUserStore(tc.dbInstance)

			switch tc.name {
			case "Valid DB Instance":
				assert.NotNil(t, userStore, "UserStore should not be nil")
				assert.Equal(t, tc.dbInstance, userStore.DB, "UserStore's db field should match the provided gorm.DB instance")
			case "Nil DB Instance":
				assert.Nil(t, userStore, "UserStore should be nil")
			case "Empty DB Instance":
				assert.NotNil(t, userStore, "UserStore should not be nil")
				assert.Equal(t, tc.dbInstance, userStore.DB, "UserStore's db field should match the provided empty gorm.DB instance")
			case "Initialized but No Connection":
				assert.NotNil(t, userStore, "UserStore should not be nil")
				assert.Equal(t, tc.dbInstance, userStore.DB, "UserStore's db field should match the provided initialized gorm.DB instance")
			case "DB Instance with Data":
				assert.NotNil(t, userStore, "UserStore should not be nil")
				assert.Equal(t, tc.dbInstance, userStore.DB, "UserStore's db field should match the provided gorm.DB instance with data")
			default:
				t.Error("Unknown test scenario")
			}
		})
	}
}
