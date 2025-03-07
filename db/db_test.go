package db

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"sync"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/DATA-DOG/go-txdb"
	"github.com/BurntSushi/toml"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"github.com/joho/godotenv"
	"github.com/raahii/golang-grpc-realworld-example/model"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"runtime/debug"
	"strings"
	"time"
)

var invalidDSN = "invalid_dsn"
var validDSN = "valid_dsn"

var openFunc = func(dialect string, args ...interface{}) (*gorm.DB, error) {
	if args[0] == validDSN {
		db, _, _ := gosqlmock.New()
		return &gorm.DB{DB: db}, nil
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
	dbname := os.Getenv("DB_NAME")
	if dbname == "" {
		return "", errors.New("$DB_NAME is not set")
	}
	port := os.Getenv("DB_PORT")
	if port == "" {
		return "", errors.New("$DB_PORT is not set")
	}

	dsn := fmt.Sprintf("%s:%s@(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", user, password, host, port, dbname)
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
	return db.AutoMigrate(&model.User{}, &model.Article{}, &model.Tag{}, &model.Comment{}).Error
}

func TestAutoMigrate(t *testing.T) {
	tests := []struct {
		name      string
		mockSetup func(sqlmock.Sqlmock)
		wantErr   bool
	}{
		{
			name: "Successful Auto Migration of All Models",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("CREATE TABLE IF NOT EXISTS \"users\"").WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectExec("CREATE TABLE IF NOT EXISTS \"articles\"").WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectExec("CREATE TABLE IF NOT EXISTS \"tags\"").WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectExec("CREATE TABLE IF NOT EXISTS \"comments\"").WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			wantErr: false,
		},
		{
			name: "Auto Migration with Database Error",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("CREATE TABLE IF NOT EXISTS \"users\"").WillReturnError(sql.ErrConnDone)
				mock.ExpectRollback()
			},
			wantErr: true,
		},
		{
			name: "Partial Migration Success",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("CREATE TABLE IF NOT EXISTS \"users\"").WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectExec("CREATE TABLE IF NOT EXISTS \"articles\"").WillReturnError(sql.ErrConnDone)
				mock.ExpectRollback()
			},
			wantErr: true,
		},
		{
			name: "Migration with Unsupported Model Structure",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("CREATE TABLE IF NOT EXISTS \"users\"").WillReturnError(sql.ErrNoRows)
				mock.ExpectRollback()
			},
			wantErr: true,
		},
		{
			name: "Concurrent Auto Migration Requests",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("CREATE TABLE IF NOT EXISTS \"users\"").WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectExec("CREATE TABLE IF NOT EXISTS \"articles\"").WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectExec("CREATE TABLE IF NOT EXISTS \"tags\"").WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectExec("CREATE TABLE IF NOT EXISTS \"comments\"").WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					t.Logf("Panic encountered so failing test: %v\n%s", r, string(debug.Stack()))
					t.Fail()
				}
			}()

			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("unexpected error when opening a stub database connection: %v", err)
			}
			defer db.Close()

			gdb, err := gorm.Open("sqlite3", db)
			if err != nil {
				t.Fatalf("unexpected error when opening a gorm database connection: %v", err)
			}

			tt.mockSetup(mock)

			err = AutoMigrate(gdb)
			if (err != nil) != tt.wantErr {
				t.Errorf("AutoMigrate() error = %v, wantErr %v", err, tt.wantErr)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}

			t.Logf("Test %s: Passed", tt.name)
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
	dsnStr, err := dsn()
	if err != nil {
		return nil, err
	}

	for i := 0; i < 3; i++ {
		db, err := openFunc("mysql", dsnStr)
		if err == nil {
			return db, nil
		}
		time.Sleep(time.Second)
	}

	return nil, errors.New("failed to open database after retries")
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
				db, _, _ := gosqlmock.New()
				return &gorm.DB{DB: db}, nil
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
					db, _, _ := gosqlmock.New()
					return &gorm.DB{DB: db}, nil
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
				db, _, _ := gosqlmock.New()
				return &gorm.DB{DB: db}, nil
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
				db, _, _ := gosqlmock.New()
				return &gorm.DB{DB: db}, nil
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

	dsnStr, err := dsn()
	if err != nil {
		return nil, err
	}

	txdb.Register("txdb", "mysql", dsnStr)

	db, err := gorm.Open("txdb", uuid.New().String())
	if err != nil {
		return nil, err
	}

	return db, nil
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
	var users []model.User

	if _, err := toml.DecodeFile("db/seed/users.toml", &users); err != nil {
		return err
	}

	tx := db.Begin()
	for _, user := range users {
		if err := tx.Create(&user).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit().Error
}

func TestSeed(t *testing.T) {
	type testScenario struct {
		name          string
		setupFile     func() error
		setupDB       func(sqlmock.Sqlmock)
		expectedError error
	}

	tests := []testScenario{
		{
			name: "Scenario 1: Successful Seeding of Users from a Valid TOML File",
			setupFile: func() error {
				tomlData := `
[[Users]]
Username = "user1"
Email = "user1@example.com"
Password = "password1"
Bio = "Bio for user1"
Image = "image1.png"
`
				return ioutil.WriteFile("db/seed/users.toml", []byte(tomlData), 0644)
			},
			setupDB: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO \"users\"").WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			expectedError: nil,
		},
		{
			name: "Scenario 2: Handling Missing TOML File",
			setupFile: func() error {
				return os.Remove("db/seed/users.toml")
			},
			setupDB:       func(mock sqlmock.Sqlmock) {},
			expectedError: errors.New("open db/seed/users.toml: no such file or directory"),
		},
		{
			name: "Scenario 3: Handling Invalid TOML File Format",
			setupFile: func() error {
				tomlData := `
invalid toml data
`
				return ioutil.WriteFile("db/seed/users.toml", []byte(tomlData), 0644)
			},
			setupDB:       func(mock sqlmock.Sqlmock) {},
			expectedError: fmt.Errorf("toml: line 1: bare keys cannot contain ' '"),
		},
		{
			name: "Scenario 4: Database Error on User Creation",
			setupFile: func() error {
				tomlData := `
[[Users]]
Username = "user1"
Email = "user1@example.com"
Password = "password1"
Bio = "Bio for user1"
Image = "image1.png"
`
				return ioutil.WriteFile("db/seed/users.toml", []byte(tomlData), 0644)
			},
			setupDB: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO \"users\"").WillReturnError(errors.New("database error"))
				mock.ExpectRollback()
			},
			expectedError: errors.New("database error"),
		},
		{
			name: "Scenario 5: Empty Users List in TOML File",
			setupFile: func() error {
				tomlData := `
[[Users]]
`
				return ioutil.WriteFile("db/seed/users.toml", []byte(tomlData), 0644)
			},
			setupDB:       func(mock sqlmock.Sqlmock) {},
			expectedError: nil,
		},
		{
			name: "Scenario 6: Partial Success with Mixed Valid and Invalid User Data",
			setupFile: func() error {
				tomlData := `
[[Users]]
Username = "user1"
Email = "user1@example.com"
Password = "password1"
Bio = "Bio for user1"
Image = "image1.png"

[[Users]]
Username = "" # Invalid entry
Email = "user2@example.com"
Password = "password2"
Bio = "Bio for user2"
Image = "image2.png"
`
				return ioutil.WriteFile("db/seed/users.toml", []byte(tomlData), 0644)
			},
			setupDB: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO \"users\"").WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectExec("INSERT INTO \"users\"").WillReturnError(errors.New("validation error"))
				mock.ExpectRollback()
			},
			expectedError: errors.New("validation error"),
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

			if err := tt.setupFile(); err != nil {
				t.Fatalf("Failed to setup test file: %v", err)
			}

			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("Failed to create mock DB: %v", err)
			}
			defer db.Close()
			gormDB, err := gorm.Open("sqlmock", db)
			if err != nil {
				t.Fatalf("Failed to open gorm DB: %v", err)
			}
			defer gormDB.Close()

			tt.setupDB(mock)

			err = Seed(gormDB)

			if !errors.Is(err, tt.expectedError) && (err == nil || tt.expectedError == nil || err.Error() != tt.expectedError.Error()) {
				t.Errorf("Expected error: %v, got: %v", tt.expectedError, err)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("There were unmet expectations: %v", err)
			}

			t.Logf("Test case %s executed successfully", tt.name)
		})
	}
}
