package db

import (
	"bytes"
	"database/sql"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"runtime/debug"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/DATA-DOG/go-txdb"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"github.com/joho/godotenv"
	"github.com/raahii/golang-grpc-realworld-example/model"
	"github.com/stretchr/testify/assert"
)

var invalidDSN = "invalid_dsn"
var validDSN = "valid_dsn"

var openFunc = func(dialect string, args ...interface{}) (*gorm.DB, error) {
	if args[0] == validDSN {
		db, _, _ := sqlmock.New()
		gormDB, err := gorm.Open("sqlmock", db)
		return gormDB, err
	}
	return nil, errors.New("failed to connect")
}

/*
ROOST_METHOD_HASH=dsn_e202d1c4f9
ROOST_METHOD_SIG_HASH=dsn_b336e03d64

FUNCTION_DEF=func dsn() (string, error)

*/
func dsn() (string, error) {
	host := os.Getenv("DB_HOST")
	if host == "" {
		return "", errors.New("$DB_HOST is not set")
	}

	user := os.Getenv("DB_USER")
	if user == "" {
		return "", errors.New("$DB_USER is not set")
	}

	password := os.Getenv("DB_PASSWORD")
	if password == "" {
		return "", errors.New("$DB_PASSWORD is not set")
	}

	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		return "", errors.New("$DB_NAME is not set")
	}

	port := os.Getenv("DB_PORT")
	if port == "" {
		return "", errors.New("$DB_PORT is not set")
	}

	dsn := fmt.Sprintf("%s:%s@(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", user, password, host, port, dbName)
	return dsn, nil
}

func TestDsn(t *testing.T) {
	type testScenario struct {
		name          string
		envSetup      func()
		expectedDSN   string
		expectedError string
	}

	tests := []testScenario{
		{
			name: "Scenario 1: Successful DSN Generation",
			envSetup: func() {
				os.Setenv("DB_HOST", "localhost")
				os.Setenv("DB_USER", "user")
				os.Setenv("DB_PASSWORD", "password")
				os.Setenv("DB_NAME", "dbname")
				os.Setenv("DB_PORT", "3306")
			},
			expectedDSN:   "user:password@(localhost:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local",
			expectedError: "",
		},
		{
			name: "Scenario 2: Missing DB_HOST Environment Variable",
			envSetup: func() {
				os.Setenv("DB_USER", "user")
				os.Setenv("DB_PASSWORD", "password")
				os.Setenv("DB_NAME", "dbname")
				os.Setenv("DB_PORT", "3306")
				os.Unsetenv("DB_HOST")
			},
			expectedDSN:   "",
			expectedError: "$DB_HOST is not set",
		},
		{
			name: "Scenario 3: Missing DB_USER Environment Variable",
			envSetup: func() {
				os.Setenv("DB_HOST", "localhost")
				os.Setenv("DB_PASSWORD", "password")
				os.Setenv("DB_NAME", "dbname")
				os.Setenv("DB_PORT", "3306")
				os.Unsetenv("DB_USER")
			},
			expectedDSN:   "",
			expectedError: "$DB_USER is not set",
		},
		{
			name: "Scenario 4: Missing DB_PASSWORD Environment Variable",
			envSetup: func() {
				os.Setenv("DB_HOST", "localhost")
				os.Setenv("DB_USER", "user")
				os.Setenv("DB_NAME", "dbname")
				os.Setenv("DB_PORT", "3306")
				os.Unsetenv("DB_PASSWORD")
			},
			expectedDSN:   "",
			expectedError: "$DB_PASSWORD is not set",
		},
		{
			name: "Scenario 5: Missing DB_NAME Environment Variable",
			envSetup: func() {
				os.Setenv("DB_HOST", "localhost")
				os.Setenv("DB_USER", "user")
				os.Setenv("DB_PASSWORD", "password")
				os.Setenv("DB_PORT", "3306")
				os.Unsetenv("DB_NAME")
			},
			expectedDSN:   "",
			expectedError: "$DB_NAME is not set",
		},
		{
			name: "Scenario 6: Missing DB_PORT Environment Variable",
			envSetup: func() {
				os.Setenv("DB_HOST", "localhost")
				os.Setenv("DB_USER", "user")
				os.Setenv("DB_PASSWORD", "password")
				os.Setenv("DB_NAME", "dbname")
				os.Unsetenv("DB_PORT")
			},
			expectedDSN:   "",
			expectedError: "$DB_PORT is not set",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					t.Logf("Panic encountered so failing test. %v\n%s", r, string(debug.Stack()))
					t.Fail()
				}
			}()

			tt.envSetup()

			dsn, err := dsn()

			if tt.expectedError != "" {
				if err == nil || err.Error() != tt.expectedError {
					t.Errorf("Expected error: %s, but got: %v", tt.expectedError, err)
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if dsn != tt.expectedDSN {
					t.Errorf("Expected DSN: %s, but got: %s", tt.expectedDSN, dsn)
				}
			}
			t.Logf("Success: %s", tt.name)
		})
	}
}

/*
ROOST_METHOD_HASH=AutoMigrate_68b002ff76
ROOST_METHOD_SIG_HASH=AutoMigrate_19557fe223

FUNCTION_DEF=func AutoMigrate(db *gorm.DB) error // AutoMigrate is a wrapper of (*gorm.DB).AutoMigrate

*/
func AutoMigrate(db *gorm.DB) error {
	// Assuming we're auto-migrating a model named User
	return db.AutoMigrate(&model.User{}).Error
}

func TestAutoMigrate(t *testing.T) {
	tests := []struct {
		name           string
		setupDB        func() (*gorm.DB, sqlmock.Sqlmock, error)
		expectedErr    bool
		expectedErrMsg string
	}{
		{
			name: "Successful Auto Migration of All Models",
			setupDB: func() (*gorm.DB, sqlmock.Sqlmock, error) {
				db, mock, err := sqlmock.New()
				if err != nil {
					return nil, nil, err
				}
				gormDB, err := gorm.Open("sqlmock", db)
				return gormDB, mock, err
			},
			expectedErr: false,
		},
		{
			name: "Auto Migration Fails Due to Database Connection Error",
			setupDB: func() (*gorm.DB, sqlmock.Sqlmock, error) {
				return nil, nil, errors.New("database connection error")
			},
			expectedErr:    true,
			expectedErrMsg: "database connection error",
		},
		{
			name: "Partial Migration Success with One Model Failing",
			setupDB: func() (*gorm.DB, sqlmock.Sqlmock, error) {
				db, mock, err := sqlmock.New()
				if err != nil {
					return nil, nil, err
				}
				gormDB, err := gorm.Open("sqlmock", db)
				if err != nil {
					return nil, nil, err
				}
				mock.ExpectExec("CREATE TABLE").WillReturnError(errors.New("schema issue"))
				return gormDB, mock, nil
			},
			expectedErr:    true,
			expectedErrMsg: "schema issue",
		},
		{
			name: "Auto Migration with Pre-existing Tables",
			setupDB: func() (*gorm.DB, sqlmock.Sqlmock, error) {
				db, mock, err := sqlmock.New()
				if err != nil {
					return nil, nil, err
				}
				gormDB, err := gorm.Open("sqlmock", db)
				return gormDB, mock, err
			},
			expectedErr: false,
		},
		{
			name: "Auto Migration with Unsupported Database Dialect",
			setupDB: func() (*gorm.DB, sqlmock.Sqlmock, error) {
				db, mock, err := sqlmock.New()
				if err != nil {
					return nil, nil, err
				}
				gormDB, err := gorm.Open("sqlmock", db)
				if err != nil {
					return nil, nil, err
				}
				mock.ExpectExec("CREATE TABLE").WillReturnError(errors.New("unsupported dialect"))
				return gormDB, mock, nil
			},
			expectedErr:    true,
			expectedErrMsg: "unsupported dialect",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					t.Logf("Panic encountered so failing test. %v\n%s", r, string(debug.Stack()))
					t.Fail()
				}
			}()

			db, mock, err := tt.setupDB()
			if err != nil {
				if !tt.expectedErr {
					t.Fatalf("unexpected error: %v", err)
				}
				if err.Error() != tt.expectedErrMsg {
					t.Fatalf("expected error message: %v, got: %v", tt.expectedErrMsg, err.Error())
				}
				return
			}
			defer db.Close()

			err = AutoMigrate(db)

			if tt.expectedErr {
				if err == nil {
					t.Errorf("expected an error but got none")
				} else if err.Error() != tt.expectedErrMsg {
					t.Errorf("expected error message: %v, got: %v", tt.expectedErrMsg, err.Error())
				}
			} else if err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

/*
ROOST_METHOD_HASH=DropTestDB_a4e5029e5a
ROOST_METHOD_SIG_HASH=DropTestDB_362a330faf

FUNCTION_DEF=func DropTestDB(d *gorm.DB) error // DropTestDB close connection

*/
func DropTestDB(d *gorm.DB) error {
	if d == nil {
		return nil
	}
	return d.Close()
}

func TestDropTestDb(t *testing.T) {
	type testScenario struct {
		description string
		setup       func() *gorm.DB
		expectedErr error
	}

	scenarios := []testScenario{
		{
			description: "Successful Database Connection Closure",
			setup: func() *gorm.DB {
				db, _, _ := sqlmock.New()
				gormDB, _ := gorm.Open("sqlmock", db)
				return gormDB
			},
			expectedErr: nil,
		},
		{
			description: "Handling of Nil Database Connection",
			setup: func() *gorm.DB {
				return nil
			},
			expectedErr: nil,
		},
		{
			description: "Database Connection Already Closed",
			setup: func() *gorm.DB {
				db, mock, _ := sqlmock.New()
				gormDB, _ := gorm.Open("sqlmock", db)
				mock.ExpectClose()
				gormDB.Close()
				return gormDB
			},
			expectedErr: nil,
		},
		{
			description: "Concurrent Database Connection Closures",
			setup: func() *gorm.DB {
				db, _, _ := sqlmock.New()
				gormDB, _ := gorm.Open("sqlmock", db)
				return gormDB
			},
			expectedErr: nil,
		},
		{
			description: "Mocked Database Connection with Error on Close",
			setup: func() *gorm.DB {
				db, mock, _ := sqlmock.New()
				gormDB, _ := gorm.Open("sqlmock", db)
				mock.ExpectClose().WillReturnError(fmt.Errorf("close error"))
				return gormDB
			},
			expectedErr: nil,
		},
	}

	for _, scenario := range scenarios {
		t.Run(scenario.description, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					t.Logf("Panic encountered so failing test. %v\n%s", r, string(debug.Stack()))
					t.Fail()
				}
			}()

			db := scenario.setup()

			var wg sync.WaitGroup
			if scenario.description == "Concurrent Database Connection Closures" {
				for i := 0; i < 10; i++ {
					wg.Add(1)
					go func() {
						defer wg.Done()
						err := DropTestDB(db)
						assert.Equal(t, scenario.expectedErr, err)
					}()
				}
				wg.Wait()
			} else {
				err := DropTestDB(db)
				assert.Equal(t, scenario.expectedErr, err)
			}

			t.Logf("Scenario '%s' passed.", scenario.description)
		})
	}
}

/*
ROOST_METHOD_HASH=New_9360a2bab3
ROOST_METHOD_SIG_HASH=New_9365800ea3

FUNCTION_DEF=func New() (*gorm.DB, error) // New return mysql connection

*/
func New() (*gorm.DB, error) {
	dsn, err := dsn()
	if err != nil {
		return nil, err
	}

	var db *gorm.DB
	for i := 0; i < 5; i++ {
		db, err = openFunc("mysql", dsn)
		if err == nil {
			return db, nil
		}
		time.Sleep(2 * time.Second)
	}

	return nil, err
}

func TestNew(t *testing.T) {
	type testCase struct {
		name        string
		dsn         func() (string, error)
		open        func(dialect string, args ...interface{}) (*gorm.DB, error)
		expectedErr bool
		expectedNil bool
	}

	tests := []testCase{
		{
			name: "Successful Database Connection",
			dsn: func() (string, error) {
				return validDSN, nil
			},
			open: func(dialect string, args ...interface{}) (*gorm.DB, error) {
				db, _, _ := sqlmock.New()
				return gorm.Open("sqlmock", db)
			},
			expectedErr: false,
			expectedNil: false,
		},
		{
			name: "Failed Database Connection Due to Invalid DSN",
			dsn: func() (string, error) {
				return invalidDSN, nil
			},
			open:        openFunc,
			expectedErr: true,
			expectedNil: true,
		},
		{
			name: "Retry Logic on Temporary Connection Failure",
			dsn: func() (string, error) {
				return validDSN, nil
			},
			open: func(dialect string, args ...interface{}) (*gorm.DB, error) {
				staticCounter := 0
				return func(dialect string, args ...interface{}) (*gorm.DB, error) {
					if staticCounter < 3 {
						staticCounter++
						return nil, errors.New("temporary error")
					}
					db, _, _ := sqlmock.New()
					return gorm.Open("sqlmock", db)
				}(dialect, args...)
			},
			expectedErr: false,
			expectedNil: false,
		},
		{
			name: "Max Retry Limit Reached",
			dsn: func() (string, error) {
				return validDSN, nil
			},
			open: func(dialect string, args ...interface{}) (*gorm.DB, error) {
				return nil, errors.New("persistent error")
			},
			expectedErr: true,
			expectedNil: true,
		},
		{
			name: "Check Database Connection Pool Configuration",
			dsn: func() (string, error) {
				return validDSN, nil
			},
			open: func(dialect string, args ...interface{}) (*gorm.DB, error) {
				db, _, _ := sqlmock.New()
				return gorm.Open("sqlmock", db)
			},
			expectedErr: false,
			expectedNil: false,
		},
		{
			name: "Log Mode Configuration",
			dsn: func() (string, error) {
				return validDSN, nil
			},
			open: func(dialect string, args ...interface{}) (*gorm.DB, error) {
				db, _, _ := sqlmock.New()
				return gorm.Open("sqlmock", db)
			},
			expectedErr: false,
			expectedNil: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					t.Logf("Panic encountered so failing test. %v\n%s", r, string(debug.Stack()))
					t.Fail()
				}
			}()

			dsn = tt.dsn
			openFunc = tt.open

			db, err := New()

			if tt.expectedErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			if tt.expectedNil {
				assert.Nil(t, db)
			} else {
				assert.NotNil(t, db)
			}
		})
	}
}

/*
ROOST_METHOD_HASH=NewTestDB_4ca8dd9747
ROOST_METHOD_SIG_HASH=NewTestDB_ab2b828b16

FUNCTION_DEF=func NewTestDB() (*gorm.DB, error) // NewTestDB return mysql connection wrapped txdb

*/
func NewTestDB() (*gorm.DB, error) {
	err := godotenv.Load("../env/test.env")
	if err != nil {
		return nil, err
	}

	dsn, err := dsn()
	if err != nil {
		return nil, err
	}

	txdb.Register("txdb", "mysql", dsn)
	db, err := sql.Open("txdb", uuid.New().String())
	if err != nil {
		return nil, err
	}

	return gorm.Open("mysql", db)
}

func TestNewTestDb(t *testing.T) {
	tests := []struct {
		name           string
		envFileExists  bool
		mockDSNSuccess bool
		txdbError      bool
		expectError    bool
	}{
		{
			name:           "Successful Database Connection Initialization",
			envFileExists:  true,
			mockDSNSuccess: true,
			txdbError:      false,
			expectError:    false,
		},
		{
			name:           "Missing Environment File",
			envFileExists:  false,
			mockDSNSuccess: true,
			txdbError:      false,
			expectError:    true,
		},
		{
			name:           "Invalid DSN String",
			envFileExists:  true,
			mockDSNSuccess: false,
			txdbError:      false,
			expectError:    true,
		},
		{
			name:           "Concurrency and Mutex Locking",
			envFileExists:  true,
			mockDSNSuccess: true,
			txdbError:      false,
			expectError:    false,
		},
		{
			name:           "Maximum Idle Connections Setting",
			envFileExists:  true,
			mockDSNSuccess: true,
			txdbError:      false,
			expectError:    false,
		},
		{
			name:           "Error During txdb Registration",
			envFileExists:  true,
			mockDSNSuccess: true,
			txdbError:      true,
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					t.Logf("Panic encountered so failing test. %v\n%s", r, string(debug.Stack()))
					t.Fail()
				}
			}()

			if tt.envFileExists {
				_ = godotenv.Load("../env/test.env")
			} else {
				_ = os.Remove("../env/test.env")
			}

			dsn = mockDSN(tt.mockDSNSuccess)

			if tt.txdbError {
				txdb.Register = func(driverName, dataSourceName, dsn string) error {
					return errors.New("txdb registration failed")
				}
			}

			db, err := NewTestDB()

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				} else {
					t.Logf("Received expected error: %v", err)
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if db == nil {
					t.Error("Expected non-nil gorm.DB object")
				} else {
					t.Log("Received valid gorm.DB object")
				}
			}

			if tt.name == "Maximum Idle Connections Setting" && db != nil {
				if db.DB().Stats().Idle > 3 {
					t.Errorf("Expected maximum idle connections to be 3, got: %v", db.DB().Stats().Idle)
				}
			}
		})
	}
}

func mockDSN(success bool) func() (string, error) {
	return func() (string, error) {
		if success {
			return "valid_dsn", nil
		}
		return "", errors.New("invalid DSN")
	}
}

/*
ROOST_METHOD_HASH=Seed_a96d5d4dee
ROOST_METHOD_SIG_HASH=Seed_878933cebc

FUNCTION_DEF=func Seed(db *gorm.DB) error

*/
func Seed(db *gorm.DB) error {
	var users struct {
		Users []model.User `toml:"Users"`
	}

	if _, err := toml.DecodeFile("db/seed/users.toml", &users); err != nil {
		return err
	}

	tx := db.Begin()
	for _, user := range users.Users {
		if err := tx.Create(&user).Error; err != nil {
			tx.Rollback()
			return err
		}
	}
	return tx.Commit().Error
}

func TestSeed(t *testing.T) {
	tests := []struct {
		name        string
		setup       func() (*gorm.DB, func(), error)
		expectedErr error
	}{
		{
			name: "Successfully Seed Users from TOML File",
			setup: func() (*gorm.DB, func(), error) {
				mockDB, mock, err := sqlmock.New()
				if err != nil {
					return nil, nil, err
				}

				db, err := gorm.Open("sqlmock", mockDB)
				if err != nil {
					return nil, nil, err
				}

				usersTOML := `
				[[Users]]
				Username = "user1"
				Email = "user1@example.com"
				Password = "password1"
				Bio = "Bio1"
				Image = "Image1"
				`

				ioutil.WriteFile("db/seed/users.toml", []byte(usersTOML), 0644)
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `users`").WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()

				return db, func() {
					mockDB.Close()
					os.Remove("db/seed/users.toml")
				}, nil
			},
			expectedErr: nil,
		},
		{
			name: "Handle Missing TOML File Gracefully",
			setup: func() (*gorm.DB, func(), error) {
				mockDB, _, err := sqlmock.New()
				if err != nil {
					return nil, nil, err
				}

				db, err := gorm.Open("sqlmock", mockDB)
				if err != nil {
					return nil, nil, err
				}

				return db, func() {
					mockDB.Close()
				}, nil
			},
			expectedErr: errors.New("open db/seed/users.toml: no such file or directory"),
		},
		{
			name: "Handle Invalid TOML Format",
			setup: func() (*gorm.DB, func(), error) {
				mockDB, _, err := sqlmock.New()
				if err != nil {
					return nil, nil, err
				}

				db, err := gorm.Open("sqlmock", mockDB)
				if err != nil {
					return nil, nil, err
				}

				ioutil.WriteFile("db/seed/users.toml", []byte("invalid toml"), 0644)

				return db, func() {
					mockDB.Close()
					os.Remove("db/seed/users.toml")
				}, nil
			},
			expectedErr: errors.New("toml: line 1: unexpected EOF"),
		},
		{
			name: "Handle Database Insertion Error",
			setup: func() (*gorm.DB, func(), error) {
				mockDB, mock, err := sqlmock.New()
				if err != nil {
					return nil, nil, err
				}

				db, err := gorm.Open("sqlmock", mockDB)
				if err != nil {
					return nil, nil, err
				}

				usersTOML := `
				[[Users]]
				Username = "user2"
				Email = "user2@example.com"
				Password = "password2"
				Bio = "Bio2"
				Image = "Image2"
				`

				ioutil.WriteFile("db/seed/users.toml", []byte(usersTOML), 0644)
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `users`").WillReturnError(errors.New("insertion error"))
				mock.ExpectRollback()

				return db, func() {
					mockDB.Close()
					os.Remove("db/seed/users.toml")
				}, nil
			},
			expectedErr: errors.New("insertion error"),
		},
		{
			name: "Seed with Empty User List",
			setup: func() (*gorm.DB, func(), error) {
				mockDB, _, err := sqlmock.New()
				if err != nil {
					return nil, nil, err
				}

				db, err := gorm.Open("sqlmock", mockDB)
				if err != nil {
					return nil, nil, err
				}

				ioutil.WriteFile("db/seed/users.toml", []byte("Users = []"), 0644)

				return db, func() {
					mockDB.Close()
					os.Remove("db/seed/users.toml")
				}, nil
			},
			expectedErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					t.Logf("Panic encountered so failing test. %v", r)
					t.Fail()
				}
			}()

			db, teardown, err := tt.setup()
			if err != nil {
				t.Fatalf("Failed to set up test: %v", err)
			}
			defer teardown()

			err = Seed(db)

			if tt.expectedErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedErr.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}

			t.Logf("Completed test case: %s", tt.name)
		})
	}
}
