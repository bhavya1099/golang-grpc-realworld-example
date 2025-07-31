package store

import (
	"fmt"
	"strings"
	"testing"
	"os"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"github.com/raahii/golang-grpc-realworld-example/model"
	"runtime/debug"
)

func TestArticleStoreAddFavorite(t *testing.T) {
	// Define table-driven test cases and scenarios
	testCases := []struct {
		name          string
		mockFn        func(mock sqlmock.Sqlmock)
		article       *model.Article
		user          *model.User
		expectedError string
	}{
		{
			name: "Successfully Add a Favorite to an Article",
			mockFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO .*favorited_users.*").WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectExec("UPDATE .*articles.*").WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			article: &model.Article{FavoritesCount: 0},
			user:    &model.User{},
			expectedError: "",
		},
		{
			name: "Fail to Add Favorite Due to Database Transaction Error",
			mockFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO .*favorited_users.*").WillReturnError(fmt.Errorf("database error"))
				mock.ExpectRollback()
			},
			article: &model.Article{},
			user:    &model.User{},
			expectedError: "database error",
		},
		{
			name: "Fail to Update Favorites Count Due to Database Transaction Error",
			mockFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO .*favorited_users.*").WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectExec("UPDATE .*articles.*").WillReturnError(fmt.Errorf("database error"))
				mock.ExpectRollback()
			},
			article: &model.Article{},
			user:    &model.User{},
			expectedError: "database error",
		},
		{
			name: "Validate Transaction Commit on Successful Operation",
			mockFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO .*favorited_users.*").WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectExec("UPDATE .*articles.*").WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			article: &model.Article{},
			user:    &model.User{},
			expectedError: "",
		},
		{
			name: "Article or User is nil",
			mockFn: func(mock sqlmock.Sqlmock) {},
			article: nil,
			user: nil,
			expectedError: "invalid input",
		},
		{
			name: "Adding a Favorite to Already Favorited Article",
			mockFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("SELECT .* FROM .*").WillReturnRows(
					sqlmock.NewRows([]string{"id"}).AddRow(1),
				)
				mock.ExpectRollback()
			},
			article: &model.Article{},
			user:    &model.User{},
			expectedError: "already favorited",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					t.Logf("Panic encountered so failing test. %v\n%s", r, string(debug.Stack()))
					t.Fail()
				}
			}()

			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("failed to create mock DB: %v", err)
			}
			defer db.Close()

			gormDB, err := gorm.Open("sqlite3", db)
			if err != nil {
				t.Fatalf("failed to open Gorm DB: %v", err)
			}
			defer gormDB.Close()

			mock.ExpectQuery("PRAGMA foreign_keys").WillReturnRows(sqlmock.NewRows([]string{"foreign_keys"}))
			tc.mockFn(mock)

			store := ArticleStore{db: gormDB}
			err = store.AddFavorite(tc.article, tc.user)

			if tc.expectedError == "" && err != nil {
				t.Errorf("unexpected error returned: %v", err)
			} else if tc.expectedError != "" && err == nil {
				t.Errorf("expected error %v but got none", tc.expectedError)
			} else if tc.expectedError != "" && !strings.Contains(err.Error(), tc.expectedError) {
				t.Errorf("expected error %v but got %v", tc.expectedError, err.Error())
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("expectations weren't met: %v", err)
			}
		})
	}
}

// TODO: Additional scenarios like concurrent behaviors could be simulated if code supports goroutines.
