package db

import (
	errors "errors"
	fmt "fmt"
	os "os"
	testing "testing"
	sql "database/sql"
	time "time"
	gosqlmock "github.com/DATA-DOG/go-sqlmock"
	gorm "github.com/jinzhu/gorm"
	model "github.com/raahii/golang-grpc-realworld-example/model"
	sync "sync"
	assert "github.com/stretchr/testify/assert"
	ioutil "io/ioutil"
	uuid "github.com/google/uuid"
	godotenv "github.com/joho/godotenv"
	gotxdb "github.com/DATA-DOG/go-txdb"
	strings "strings"
	toml "github.com/BurntSushi/toml"
	debug "runtime/debug"
)



var invalidDSN = "invalid_dsn"
var openFunc = func(dialect string, args ...interface{}) (*gorm.DB, error) {
	if args[0] == validDSN {
		db, _, _ := gosqlmock.New()
		return &gorm.DB{DB: db}, nil
	}
	return nil, errors.New("failed to connect")
}
var validDSN = "valid_dsn"




/*
ROOST_METHOD_HASH=dsn_e202d1c4f9
ROOST_METHOD_SIG_HASH=dsn_b336e03d64

FUNCTION_DEF=func dsn() (string, error) 

*/
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
func TestAutoMigrate(t *testing.T) {
	type testCase struct {
		name        string
		setup       func() (*gorm.DB, sqlmock.Sqlmock)
		expectedErr error
	}

	tests := []testCase{
		{
			name: "Successful Auto Migration of All Models",
			setup: func() (*gorm.DB, sqlmock.Sqlmock) {
				db, mock, err := sqlmock.New()
				if err != nil {
					t.Fatalf("failed to open sqlmock database: %v", err)
				}
				gormDB, err := gorm.Open("mysql", db)
				if err != nil {
					t.Fatalf("failed to initialize gorm DB: %v", err)
				}
				mock.ExpectExec("CREATE TABLE").WillReturnResult(sqlmock.NewResult(1, 1))
				return gormDB, mock
			},
			expectedErr: nil,
		},
		{
			name: "Handling Database Connection Error",
			setup: func() (*gorm.DB, sqlmock.Sqlmock) {
				db, mock, err := sqlmock.New()
				if err != nil {
					t.Fatalf("failed to open sqlmock database: %v", err)
				}
				gormDB, err := gorm.Open("mysql", db)
				if err != nil {
					t.Fatalf("failed to initialize gorm DB: %v", err)
				}
				mock.ExpectExec("CREATE TABLE").WillReturnError(fmt.Errorf("connection error"))
				return gormDB, mock
			},
			expectedErr: fmt.Errorf("connection error"),
		},
		{
			name: "Partial Migration Due to Model Error",
			setup: func() (*gorm.DB, sqlmock.Sqlmock) {
				db, mock, err := sqlmock.New()
				if err != nil {
					t.Fatalf("failed to open sqlmock database: %v", err)
				}
				gormDB, err := gorm.Open("mysql", db)
				if err != nil {
					t.Fatalf("failed to initialize gorm DB: %v", err)
				}

				mock.ExpectExec("CREATE TABLE").WillReturnError(fmt.Errorf("unsupported data type"))
				return gormDB, mock
			},
			expectedErr: fmt.Errorf("unsupported data type"),
		},
		{
			name: "Verifying Idempotency of Auto Migration",
			setup: func() (*gorm.DB, sqlmock.Sqlmock) {
				db, mock, err := sqlmock.New()
				if err != nil {
					t.Fatalf("failed to open sqlmock database: %v", err)
				}
				gormDB, err := gorm.Open("mysql", db)
				if err != nil {
					t.Fatalf("failed to initialize gorm DB: %v", err)
				}
				mock.ExpectExec("CREATE TABLE").WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectExec("CREATE TABLE").WillReturnResult(sqlmock.NewResult(1, 1))
				return gormDB, mock
			},
			expectedErr: nil,
		},
		{
			name: "Handling Migration of Already Existing Tables",
			setup: func() (*gorm.DB, sqlmock.Sqlmock) {
				db, mock, err := sqlmock.New()
				if err != nil {
					t.Fatalf("failed to open sqlmock database: %v", err)
				}
				gormDB, err := gorm.Open("mysql", db)
				if err != nil {
					t.Fatalf("failed to initialize gorm DB: %v", err)
				}
				mock.ExpectExec("CREATE TABLE").WillReturnError(fmt.Errorf("table already exists"))
				return gormDB, mock
			},
			expectedErr: fmt.Errorf("table already exists"),
		},
		{
			name: "Migration with Complex Relationships",
			setup: func() (*gorm.DB, sqlmock.Sqlmock) {
				db, mock, err := sqlmock.New()
				if err != nil {
					t.Fatalf("failed to open sqlmock database: %v", err)
				}
				gormDB, err := gorm.Open("mysql", db)
				if err != nil {
					t.Fatalf("failed to initialize gorm DB: %v", err)
				}
				mock.ExpectExec("CREATE TABLE").WillReturnResult(sqlmock.NewResult(1, 1))
				return gormDB, mock
			},
			expectedErr: nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					t.Logf("Panic encountered so failing test. %v\n%s", r, string(debug.Stack()))
					t.Fail()
				}
			}()

			db, mock := tc.setup()
			defer db.Close()

			err := AutoMigrate(db)
			if err != nil && tc.expectedErr == nil {
				t.Errorf("expected no error, but got %v", err)
			}
			if err == nil && tc.expectedErr != nil {
				t.Errorf("expected error %v, but got none", tc.expectedErr)
			}
			if err != nil && tc.expectedErr != nil && err.Error() != tc.expectedErr.Error() {
				t.Errorf("expected error %v, but got %v", tc.expectedErr, err)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}

			t.Logf("Test case '%s' executed successfully", tc.name)
		})
	}
}


/*
ROOST_METHOD_HASH=DropTestDB_a4e5029e5a
ROOST_METHOD_SIG_HASH=DropTestDB_362a330faf

FUNCTION_DEF=func DropTestDB(d *gorm.DB) error // DropTestDB close connection


*/
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
func TestSeed(t *testing.T) {
	type testCase struct {
		name           string
		setup          func() (*gorm.DB, func(), error)
		expectedErr    bool
		expectedErrMsg string
		expectedCount  int
	}

	tests := []testCase{
		{
			name: "Successful Seeding of Users from TOML File",
			setup: func() (*gorm.DB, func(), error) {
				mockDB, mock, err := sqlmock.New()
				if err != nil {
					return nil, nil, err
				}
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO").WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()

				db, err := gorm.Open("sqlmock", mockDB)
				if err != nil {
					return nil, nil, err
				}

				tomlData := `[[Users]]
Username = "testuser"
Email = "test@example.com"
Password = "password"
Bio = "bio"
Image = "image"`

				filePath := "db/seed/users.toml"
				err = ioutil.WriteFile(filePath, []byte(tomlData), 0644)
				if err != nil {
					return nil, nil, err
				}

				cleanup := func() {
					_ = os.Remove(filePath)
				}

				return db, cleanup, nil
			},
			expectedErr:   false,
			expectedCount: 1,
		},
		{
			name: "Error Reading TOML File",
			setup: func() (*gorm.DB, func(), error) {
				mockDB, _, err := sqlmock.New()
				if err != nil {
					return nil, nil, err
				}

				db, err := gorm.Open("sqlmock", mockDB)
				if err != nil {
					return nil, nil, err
				}

				filePath := "db/seed/nonexistent.toml"

				cleanup := func() {}

				return db, cleanup, nil
			},
			expectedErr:    true,
			expectedErrMsg: "open db/seed/nonexistent.toml: no such file or directory",
		},
		{
			name: "Error Decoding TOML File",
			setup: func() (*gorm.DB, func(), error) {
				mockDB, _, err := sqlmock.New()
				if err != nil {
					return nil, nil, err
				}

				db, err := gorm.Open("sqlmock", mockDB)
				if err != nil {
					return nil, nil, err
				}

				tomlData := `[[Users]
Username = "testuser"`

				filePath := "db/seed/malformed.toml"
				err = ioutil.WriteFile(filePath, []byte(tomlData), 0644)
				if err != nil {
					return nil, nil, err
				}

				cleanup := func() {
					_ = os.Remove(filePath)
				}

				return db, cleanup, nil
			},
			expectedErr:    true,
			expectedErrMsg: "toml: line 1: error: unexpected end of table",
		},
		{
			name: "Database Insertion Error",
			setup: func() (*gorm.DB, func(), error) {
				mockDB, mock, err := sqlmock.New()
				if err != nil {
					return nil, nil, err
				}
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO").WillReturnError(fmt.Errorf("insertion error"))
				mock.ExpectRollback()

				db, err := gorm.Open("sqlmock", mockDB)
				if err != nil {
					return nil, nil, err
				}

				tomlData := `[[Users]]
Username = "testuser"
Email = "test@example.com"
Password = "password"
Bio = "bio"
Image = "image"`

				filePath := "db/seed/users.toml"
				err = ioutil.WriteFile(filePath, []byte(tomlData), 0644)
				if err != nil {
					return nil, nil, err
				}

				cleanup := func() {
					_ = os.Remove(filePath)
				}

				return db, cleanup, nil
			},
			expectedErr:    true,
			expectedErrMsg: "insertion error",
		},
		{
			name: "Empty TOML File Handling",
			setup: func() (*gorm.DB, func(), error) {
				mockDB, _, err := sqlmock.New()
				if err != nil {
					return nil, nil, err
				}

				db, err := gorm.Open("sqlmock", mockDB)
				if err != nil {
					return nil, nil, err
				}

				filePath := "db/seed/empty.toml"
				err = ioutil.WriteFile(filePath, []byte{}, 0644)
				if err != nil {
					return nil, nil, err
				}

				cleanup := func() {
					_ = os.Remove(filePath)
				}

				return db, cleanup, nil
			},
			expectedErr:   false,
			expectedCount: 0,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					t.Logf("Panic encountered so failing test. %v\n%s", r, string(debug.Stack()))
					t.Fail()
				}
			}()

			db, cleanup, err := tc.setup()
			if err != nil {
				t.Fatalf("Setup failed: %v", err)
			}
			defer cleanup()

			err = Seed(db)

			if tc.expectedErr {
				if err == nil {
					t.Fatalf("Expected error but got none")
				}
				if !strings.Contains(err.Error(), tc.expectedErrMsg) {
					t.Fatalf("Expected error message to contain %q, got %q", tc.expectedErrMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Fatalf("Expected no error but got %v", err)
				}

			}

			t.Logf("Test case %s passed", tc.name)
		})
	}
}

